package service

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/upstream"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type UpstreamService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewUpstreamService(db *gorm.DB, log *logger.Logger) *UpstreamService {
	return &UpstreamService{db: db, log: log}
}

// GetDB 获取数据库连接
func (s *UpstreamService) GetDB() *gorm.DB {
	return s.db
}

// GetProviderByID returns an upstream provider by its primary key.
func (s *UpstreamService) GetProviderByID(id uint) (*model.UpstreamProvider, error) {
	var provider model.UpstreamProvider
	if err := s.db.First(&provider, id).Error; err != nil {
		return nil, err
	}
	return &provider, nil
}

// TestConnection performs a real connection test against the upstream provider.
func (s *UpstreamService) TestConnection(id uint) (*upstream.ConnectionResult, error) {
	provider, err := s.GetProviderByID(id)
	if err != nil {
		return nil, fmt.Errorf("provider not found: %w", err)
	}

	client, err := upstream.NewClient(provider)
	if err != nil {
		return nil, err
	}

	result, err := client.TestConnection()

	status := "success"
	message := result.Message
	if err != nil {
		status = "failed"
		message = err.Error()
	}
	s.db.Create(&model.UpstreamSyncLog{
		UpstreamID: provider.ID,
		Action:     "test_connection",
		Status:     status,
		Message:    message,
	})

	return result, err
}

// generateSlug creates a URL-friendly slug from a product name and remote ID.
func generateSlug(name string, remoteID string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	var b strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		}
	}
	slug = b.String()
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = "product"
	}
	return fmt.Sprintf("upstream-%s-%s", slug, remoteID)
}

// ensureGroupID returns a product group ID for synced products.
func (s *UpstreamService) ensureGroupID(provider *model.UpstreamProvider) (uint, error) {
	if provider.Config != nil {
		if gid, ok := provider.Config["group_id"]; ok {
			switch v := gid.(type) {
			case float64:
				return uint(v), nil
			case int:
				return uint(v), nil
			case uint:
				return v, nil
			}
		}
	}

	var group model.ProductGroup
	err := s.db.Where("slug = ?", "upstream-products").First(&group).Error
	if err == gorm.ErrRecordNotFound {
		group = model.ProductGroup{
			Name:      "上游产品",
			Slug:      "upstream-products",
			SortOrder: 999,
			Status:    1,
		}
		if err := s.db.Create(&group).Error; err != nil {
			return 0, fmt.Errorf("create default group: %w", err)
		}
		return group.ID, nil
	} else if err != nil {
		return 0, fmt.Errorf("query group: %w", err)
	}
	return group.ID, nil
}

// GetUpstreamProducts 获取上游产品列表（含分组结构）
func (s *UpstreamService) GetUpstreamProducts(providerID uint) (*upstream.UpstreamProductsResult, error) {
	provider, err := s.GetProviderByID(providerID)
	if err != nil {
		return nil, fmt.Errorf("供应商不存在: %w", err)
	}

	client, err := upstream.NewClient(provider)
	if err != nil {
		return nil, err
	}

	return client.FetchProductsWithGroups()
}

// GetUpstreamGroups 获取上游分组列表
func (s *UpstreamService) GetUpstreamGroups(providerID uint) ([]upstream.RemoteProductGroup, error) {
	result, err := s.GetUpstreamProducts(providerID)
	if err != nil {
		return nil, err
	}
	return result.Groups, nil
}

// GetLocalGroups 获取本地分组列表
func (s *UpstreamService) GetLocalGroups() ([]model.ProductGroup, error) {
	var groups []model.ProductGroup
	err := s.db.Where("status = ?", 1).Order("sort_order ASC").Find(&groups).Error
	return groups, err
}

// CreateLocalGroup 创建本地分组
func (s *UpstreamService) CreateLocalGroup(name string) (*model.ProductGroup, error) {
	group := &model.ProductGroup{
		Name:      name,
		Status:    1,
		SortOrder: 0,
	}
	if err := s.db.Create(group).Error; err != nil {
		return nil, err
	}
	return group, nil
}

// DockConfig 对接配置
type DockConfig struct {
	ProviderID   uint    `json:"provider_id"`
	LocalGroupID uint    `json:"local_group_id"`
	Percent      float64 `json:"percent"`
	Concurrency  int     `json:"concurrency"`
}

// DockResult 对接结果
type DockResult struct {
	TaskID       uint   `json:"task_id"`
	TotalCount   int    `json:"total_count"`
	SuccessCount int    `json:"success_count"`
	FailedCount  int    `json:"failed_count"`
}

// DockProducts 对接指定产品（支持多线程）
func (s *UpstreamService) DockProducts(config DockConfig, productIDs []string) (*DockResult, error) {
	provider, err := s.GetProviderByID(config.ProviderID)
	if err != nil {
		return nil, fmt.Errorf("供应商不存在")
	}

	if len(productIDs) == 0 {
		return nil, fmt.Errorf("没有要对接的产品")
	}

	if config.Percent <= 0 {
		config.Percent = 120
	}
	if config.Concurrency <= 0 {
		config.Concurrency = 10
	}

	// 获取上游产品信息
	client, err := upstream.NewClient(provider)
	if err != nil {
		return nil, err
	}
	result, _ := client.FetchProductsWithGroups()
	productMap := make(map[string]upstream.RemoteProduct)
	if result != nil {
		for _, p := range result.Products {
			productMap[p.RemoteID] = p
		}
	}

	// 多线程对接
	var (
		wg      sync.WaitGroup
		sem     = make(chan struct{}, config.Concurrency)
		success int32
		failed  int32
	)

	for _, pid := range productIDs {
		wg.Add(1)
		sem <- struct{}{}

		go func(remoteID string) {
			defer wg.Done()
			defer func() { <-sem }()

			remoteProduct := productMap[remoteID]
			if err := s.dockSingleProduct(config.ProviderID, config.LocalGroupID, config.Percent, remoteID, remoteProduct); err != nil {
				atomic.AddInt32(&failed, 1)
				s.log.Errorf("对接产品 %s 失败: %v", remoteID, err)
			} else {
				atomic.AddInt32(&success, 1)
			}
		}(pid)
	}

	wg.Wait()

	return &DockResult{
		TotalCount:   len(productIDs),
		SuccessCount: int(success),
		FailedCount:  int(failed),
	}, nil
}

// DockGroup 对接整个分组
func (s *UpstreamService) DockGroup(providerID uint, groupID string, localGroupID uint, percent float64, concurrency int) (*DockResult, error) {
	result, err := s.GetUpstreamProducts(providerID)
	if err != nil {
		return nil, err
	}

	var productIDs []string
	for _, p := range result.Products {
		if p.GroupID == groupID {
			productIDs = append(productIDs, p.RemoteID)
		}
	}

	if len(productIDs) == 0 {
		return nil, fmt.Errorf("该分组下没有产品")
	}

	return s.DockProducts(DockConfig{
		ProviderID:   providerID,
		LocalGroupID: localGroupID,
		Percent:      percent,
		Concurrency:  concurrency,
	}, productIDs)
}

// dockSingleProduct 对接单个产品
func (s *UpstreamService) dockSingleProduct(providerID, localGroupID uint, percent float64, remoteID string, remoteProduct upstream.RemoteProduct) error {
	// 检查是否已对接
	var existing model.UpstreamProduct
	if err := s.db.Where("upstream_id = ? AND remote_product_id = ?", providerID, remoteID).First(&existing).Error; err == nil {
		return nil // 已对接
	}

	// 计算价格
	price := remoteProduct.Price * percent / 100

	// 创建本地产品
	product := model.Product{
		GroupID:  localGroupID,
		Name:     remoteProduct.Name,
		Slug:     generateSlug(remoteProduct.Name, remoteID),
		Price:    decimal.NewFromFloat(price),
		Currency: remoteProduct.Currency,
		Type:     "upstream",
		Status:   1,
	}
	if err := s.db.Create(&product).Error; err != nil {
		return err
	}

	// 创建映射
	mapping := model.UpstreamProduct{
		LocalProductID:  product.ID,
		UpstreamID:      providerID,
		RemoteProductID: remoteID,
		Config: model.JSON{
			"remote_name": remoteProduct.Name,
			"percent":     percent,
		},
	}
	return s.db.Create(&mapping).Error
}

// SyncProducts fetches products from the upstream and upserts them.
func (s *UpstreamService) SyncProducts(id uint) (int, error) {
	provider, err := s.GetProviderByID(id)
	if err != nil {
		return 0, fmt.Errorf("provider not found: %w", err)
	}

	client, err := upstream.NewClient(provider)
	if err != nil {
		return 0, err
	}

	groupID, err := s.ensureGroupID(provider)
	if err != nil {
		return 0, fmt.Errorf("resolve product group: %w", err)
	}

	remoteProducts, err := client.FetchProducts()
	if err != nil {
		s.db.Create(&model.UpstreamSyncLog{
			UpstreamID: provider.ID,
			Action:     "sync_products",
			Status:     "failed",
			Message:    err.Error(),
		})
		return 0, fmt.Errorf("fetch products: %w", err)
	}

	synced := 0
	for _, rp := range remoteProducts {
		var mapping model.UpstreamProduct
		err := s.db.Where("upstream_id = ? AND remote_product_id = ?", provider.ID, rp.RemoteID).
			First(&mapping).Error

		price := decimal.NewFromFloat(rp.Price)

		if err == gorm.ErrRecordNotFound {
			localProduct := model.Product{
				GroupID:  groupID,
				Name:     rp.Name,
				Slug:     generateSlug(rp.Name, rp.RemoteID),
				Price:    price,
				Currency: rp.Currency,
				Type:     rp.Type,
				Status:   1,
			}
			if err := s.db.Create(&localProduct).Error; err != nil {
				s.log.WithField("remote_id", rp.RemoteID).Errorf("create local product: %v", err)
				continue
			}

			mapping = model.UpstreamProduct{
				LocalProductID:  localProduct.ID,
				UpstreamID:      provider.ID,
				RemoteProductID: rp.RemoteID,
				Config: model.JSON{
					"remote_name": rp.Name,
				},
			}
			if err := s.db.Create(&mapping).Error; err != nil {
				s.log.WithField("remote_id", rp.RemoteID).Errorf("create mapping: %v", err)
				continue
			}
		} else if err != nil {
			s.log.WithField("remote_id", rp.RemoteID).Errorf("query mapping: %v", err)
			continue
		} else {
			s.db.Model(&model.Product{}).Where("id = ?", mapping.LocalProductID).Updates(map[string]interface{}{
				"name":     rp.Name,
				"price":    price,
				"currency": rp.Currency,
				"type":     rp.Type,
			})
			s.db.Model(&mapping).Updates(map[string]interface{}{
				"config": model.JSON{"remote_name": rp.Name},
			})
		}
		synced++
	}

	s.db.Create(&model.UpstreamSyncLog{
		UpstreamID: provider.ID,
		Action:     "sync_products",
		Status:     "success",
		Message:    fmt.Sprintf("synced %d products", synced),
	})

	s.log.WithField("provider_id", id).Infof("upstream product sync completed: %d products", synced)
	return synced, nil
}

// SyncAllProducts 同步所有对接产品
func (s *UpstreamService) SyncAllProducts(providerID uint) (int, error) {
	var mappings []model.UpstreamProduct
	if err := s.db.Where("upstream_id = ?", providerID).Find(&mappings).Error; err != nil {
		return 0, err
	}

	synced := 0
	for _, m := range mappings {
		if err := s.SyncSingleProduct(m.LocalProductID); err == nil {
			synced++
		}
	}
	return synced, nil
}

// SyncSingleProduct 同步单个产品
func (s *UpstreamService) SyncSingleProduct(localProductID uint) error {
	var mapping model.UpstreamProduct
	if err := s.db.Where("local_product_id = ?", localProductID).First(&mapping).Error; err != nil {
		return fmt.Errorf("产品未对接上游")
	}

	var provider model.UpstreamProvider
	if err := s.db.First(&provider, mapping.UpstreamID).Error; err != nil {
		return fmt.Errorf("供应商不存在")
	}

	client, err := upstream.NewClient(&provider)
	if err != nil {
		return err
	}

	products, err := client.FetchProducts()
	if err != nil {
		return err
	}

	for _, p := range products {
		if p.RemoteID == mapping.RemoteProductID {
			price := decimal.NewFromFloat(p.Price)
			s.db.Model(&model.Product{}).Where("id = ?", localProductID).Updates(map[string]interface{}{
				"name":     p.Name,
				"price":    price,
				"currency": p.Currency,
			})
			return nil
		}
	}

	return fmt.Errorf("上游产品不存在")
}

// GetProviderList 获取供应商列表
func (s *UpstreamService) GetProviderList() ([]model.UpstreamProvider, error) {
	var providers []model.UpstreamProvider
	err := s.db.Find(&providers).Error
	return providers, err
}
