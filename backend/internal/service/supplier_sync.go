package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/YeXuanHs/anchor-finance/internal/pluginengine"
)

// UpstreamDriver 上游驱动接口（MD 7.3 参考创欧UpstreamDriver）
type UpstreamDriver interface {
	// Key 驱动唯一标识
	Key() string
	// Name 驱动显示名称
	Name() string
	// Capabilities 支持的能力
	Capabilities() []string
	// FetchProducts 拉取上游商品列表
	FetchProducts() ([]RemoteProduct, error)
	// FetchProductGroups 拉取上游商品分组
	FetchProductGroups() ([]RemoteGroup, error)
	// GetProductStructure 获取完整商品结构（分组+商品树）
	GetProductStructure() (*ProductStructure, error)
	// SyncStatus 同步服务状态
	SyncStatus(serviceID string) (*StatusResult, error)
	// CreateService 创建服务（开通）
	CreateService(params CreateServiceParams) (*ServiceResult, error)
	// SuspendService 暂停服务
	SuspendService(serviceID string) error
	// UnsuspendService 取消暂停
	UnsuspendService(serviceID string) error
	// TerminateService 终止服务
	TerminateService(serviceID string) error
	// RenewService 续费服务
	RenewService(serviceID string, cycle string) error
}

// ProductStructure 商品结构（MD 7.3）
type ProductStructure struct {
	Groups []RemoteGroup `json:"groups"`
}

// RemoteProduct 上游商品
type RemoteProduct struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	GroupID     string  `json:"group_id"`
	GroupName   string  `json:"group_name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Cycle       string  `json:"cycle"`
	Stock       int     `json:"stock"`
	Status      string  `json:"status"`
}

// RemoteGroup 上游商品分组（MD 7.3）
type RemoteGroup struct {
	ID       string        `json:"id"`
	Name     string        `json:"name"`
	ParentID string        `json:"parent_id"`
	Children []RemoteGroup `json:"children,omitempty"`
	Products []RemoteProduct `json:"products,omitempty"`
}

// StatusResult 状态同步结果
type StatusResult struct {
	Status    string `json:"status"`
	IPAddress string `json:"ip_address"`
	Hostname  string `json:"hostname"`
}

// CreateServiceParams 创建服务参数
type CreateServiceParams struct {
	ProductID string `json:"product_id"`
	Cycle     string `json:"cycle"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

// ServiceResult 服务创建结果
type ServiceResult struct {
	ServiceID string `json:"service_id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	IPAddress string `json:"ip_address"`
	Hostname  string `json:"hostname"`
}

// PluginDriver 基于PHP插件引擎的上游驱动（通用实现）
type PluginDriver struct {
	supplierID uint
	slug       string
	domain     string
	apiURL     string
	apiKey     string
	apiSecret  string
	apiType    string // manual, zjmf, v10, anchor
}

// NewDriver 根据供应商类型创建对应的驱动（Go原生实现，不走PHP）
func NewDriver(supplier model.Supplier) UpstreamDriver {
	switch supplier.Type {
	case "zjmf":
		return NewZjmfDriver(supplier.ID, supplier.APIURL, supplier.APIKey, supplier.APISecret)
	case "v10":
		return NewV10Driver(supplier.ID, supplier.APIURL, supplier.APIKey, supplier.APISecret)
	case "anchor":
		return NewAnchorDriver(supplier.ID, supplier.APIURL, supplier.APIKey, supplier.APISecret)
	case "manual":
		return NewManualDriver(supplier.ID)
	default:
		// 未知类型走PHP插件引擎（兼容未来扩展）
		return NewPluginDriver(supplier)
	}
}

// NewPluginDriver 创建插件驱动（兼容旧代码，走PHP插件引擎）
func NewPluginDriver(supplier model.Supplier) *PluginDriver {
	return &PluginDriver{
		supplierID: supplier.ID,
		slug:       supplier.Slug,
		domain:     supplier.Domain,
		apiURL:     supplier.APIURL,
		apiKey:     supplier.APIKey,
		apiSecret:  supplier.APISecret,
		apiType:    supplier.Type,
	}
}

func (d *PluginDriver) Key() string   { return d.slug }
func (d *PluginDriver) Name() string  { return d.slug }
func (d *PluginDriver) Capabilities() []string {
	return []string{"provisioning", "renewal", "status_sync", "product_sync"}
}

func (d *PluginDriver) FetchProducts() ([]RemoteProduct, error) {
	results, err := pluginengine.TriggerHook("supplier_fetch_products", map[string]interface{}{
		"supplier_id": d.supplierID,
		"api_url":     d.apiURL,
		"api_key":     d.apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("插件引擎离线: %w", err)
	}

	var products []RemoteProduct
	if len(results) > 0 && results[0].Data != nil {
		if data, ok := results[0].Data.(map[string]interface{}); ok {
			if items, ok := data["products"].([]interface{}); ok {
				for _, item := range items {
					if m, ok := item.(map[string]interface{}); ok {
						p := RemoteProduct{
							ID:   fmt.Sprintf("%v", m["id"]),
							Name: fmt.Sprintf("%v", m["name"]),
						}
						if price, ok := m["price"].(float64); ok {
							p.Price = price
						}
						if stock, ok := m["stock"].(float64); ok {
							p.Stock = int(stock)
						}
						products = append(products, p)
					}
				}
			}
		}
	}
	return products, nil
}

func (d *PluginDriver) FetchProductGroups() ([]RemoteGroup, error) {
	results, err := pluginengine.TriggerHook("supplier_fetch_groups", map[string]interface{}{
		"supplier_id": d.supplierID,
		"api_url":     d.apiURL,
		"api_key":     d.apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("插件引擎离线: %w", err)
	}

	var groups []RemoteGroup
	if len(results) > 0 && results[0].Data != nil {
		if data, ok := results[0].Data.(map[string]interface{}); ok {
			if items, ok := data["groups"].([]interface{}); ok {
				for _, item := range items {
					if m, ok := item.(map[string]interface{}); ok {
						g := RemoteGroup{
							ID:   fmt.Sprintf("%v", m["id"]),
							Name: fmt.Sprintf("%v", m["name"]),
						}
						groups = append(groups, g)
					}
				}
			}
		}
	}
	return groups, nil
}

func (d *PluginDriver) GetProductStructure() (*ProductStructure, error) {
	groups, err := d.FetchProductGroups()
	if err != nil {
		return nil, err
	}
	return &ProductStructure{Groups: groups}, nil
}

func (d *PluginDriver) SyncStatus(serviceID string) (*StatusResult, error) {
	results, err := pluginengine.TriggerHook("service_status_sync", map[string]interface{}{
		"service_id":  serviceID,
		"supplier_id": d.supplierID,
	})
	if err != nil {
		return nil, fmt.Errorf("插件引擎离线: %w", err)
	}
	if len(results) > 0 && results[0].Data != nil {
		if data, ok := results[0].Data.(map[string]interface{}); ok {
			sr := &StatusResult{}
			if s, ok := data["status"].(string); ok {
				sr.Status = s
			}
			if ip, ok := data["ip_address"].(string); ok {
				sr.IPAddress = ip
			}
			return sr, nil
		}
	}
	return &StatusResult{Status: "unknown"}, nil
}

func (d *PluginDriver) CreateService(params CreateServiceParams) (*ServiceResult, error) {
	results, err := pluginengine.TriggerHook("service_create", map[string]interface{}{
		"supplier_id": d.supplierID,
		"product_id":  params.ProductID,
		"cycle":       params.Cycle,
		"username":    params.Username,
		"password":    params.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("插件引擎离线: %w", err)
	}
	if len(results) > 0 && results[0].Data != nil {
		if data, ok := results[0].Data.(map[string]interface{}); ok {
			sr := &ServiceResult{}
			if sid, ok := data["service_id"].(string); ok {
				sr.ServiceID = sid
			}
			if ip, ok := data["ip_address"].(string); ok {
				sr.IPAddress = ip
			}
			return sr, nil
		}
	}
	return &ServiceResult{}, nil
}

func (d *PluginDriver) SuspendService(serviceID string) error {
	_, err := pluginengine.TriggerHook("service_suspend", map[string]interface{}{
		"service_id":  serviceID,
		"supplier_id": d.supplierID,
	})
	return err
}

func (d *PluginDriver) UnsuspendService(serviceID string) error {
	_, err := pluginengine.TriggerHook("service_unsuspend", map[string]interface{}{
		"service_id":  serviceID,
		"supplier_id": d.supplierID,
	})
	return err
}

func (d *PluginDriver) TerminateService(serviceID string) error {
	_, err := pluginengine.TriggerHook("service_terminate", map[string]interface{}{
		"service_id":  serviceID,
		"supplier_id": d.supplierID,
	})
	return err
}

func (d *PluginDriver) RenewService(serviceID string, cycle string) error {
	_, err := pluginengine.TriggerHook("service_renew", map[string]interface{}{
		"service_id":  serviceID,
		"supplier_id": d.supplierID,
		"cycle":       cycle,
	})
	return err
}

// SupplierSyncService 供应商同步服务
type SupplierSyncService struct {
	drivers map[uint]UpstreamDriver
	mu      sync.RWMutex
}

// NewSupplierSyncService 创建供应商同步服务
func NewSupplierSyncService() *SupplierSyncService {
	return &SupplierSyncService{
		drivers: make(map[uint]UpstreamDriver),
	}
}

// RegisterDriver 注册驱动
func (s *SupplierSyncService) RegisterDriver(supplierID uint, driver UpstreamDriver) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.drivers[supplierID] = driver
}

// GetDriver 获取驱动
func (s *SupplierSyncService) GetDriver(supplierID uint) UpstreamDriver {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.drivers[supplierID]
}

// SyncAllProducts 同步所有供应商商品（多线程，自动创建分组+自动上架）
func (s *SupplierSyncService) SyncAllProducts(supplierID uint) error {
	driver := s.GetDriver(supplierID)
	if driver == nil {
		return fmt.Errorf("供应商驱动未注册: %d", supplierID)
	}

	// MD 7.2.2: 拉取上游分组并自动创建本地分组（需开启auto_create_groups）
	if getSettingInt("auto_create_groups", 0) == 1 {
		groups, err := driver.FetchProductGroups()
		if err == nil && len(groups) > 0 {
			s.autoCreateGroups(groups)
		}
	}

	products, err := driver.FetchProducts()
	if err != nil {
		return fmt.Errorf("拉取商品失败: %w", err)
	}

	db := database.GetDB()

	// 检查是否启用自动上架
	autoListing := false
	var listingSetting model.Setting
	if err := db.Where("`key` = ?", "auto_listing_enabled").First(&listingSetting).Error; err == nil && listingSetting.Value == "1" {
		autoListing = true
	}

	// MD 7.2.3: 多线程并发同步（每个商品一个goroutine）
	var wg sync.WaitGroup
	for _, p := range products {
		wg.Add(1)
		go func(product RemoteProduct) {
			defer wg.Done()
			s.syncSingleProduct(supplierID, product, autoListing)
		}(p)
	}
	wg.Wait()
	return nil
}

// syncSingleProduct 同步单个商品（goroutine安全）
func (s *SupplierSyncService) syncSingleProduct(supplierID uint, p RemoteProduct, autoListing bool) {
	db := database.GetDB()

	// 查找或创建本地供应商商品
	var existing model.SupplierProduct
	result := db.Where("supplier_id = ? AND remote_product_id = ?", supplierID, p.ID).First(&existing)
	if result.Error != nil {
		// 新商品，创建供应商商品记录
		profitRate := float64(getSettingInt("default_profit_rate", 25))
		localPrice := p.Price * (1 + profitRate/100)

		sp := model.SupplierProduct{
			SupplierID:      supplierID,
			RemoteProductID: p.ID,
			Name:            p.Name,
			RemotePrice:     p.Price,
			LocalPrice:      localPrice,
			ProfitRate:      profitRate,
			Stock:           p.Stock,
			Status:          "active",
		}
		db.Create(&sp)

		// MD 7.2.5: 自动上架（读取分组映射表，创建本地Product到对应分组）
		if autoListing {
			product := model.Product{
				Name:   p.Name,
				Type:   "server",
				Price:  localPrice,
				Amount: localPrice,
				Status: "active",
			}
			// 查找分组映射（上游分组→本地分组）
			if p.GroupID != "" {
				var mapping model.SupplierGroupMapping
				if err := db.Where("supplier_id = ? AND remote_group_id = ?", supplierID, p.GroupID).First(&mapping).Error; err == nil {
					product.GroupID = mapping.LocalGroupID
					if mapping.ProfitRate > 0 {
						product.Price = p.Price * (1 + mapping.ProfitRate/100)
						product.Amount = product.Price
					}
				}
			}
			db.Create(&product)
		}
	} else {
		// 已有商品，更新价格和库存
		db.Model(&existing).Updates(map[string]interface{}{
			"remote_price": p.Price,
			"stock":        p.Stock,
			"name":         p.Name,
		})
	}
}

// SyncAllPrices 同步所有供应商价格
func (s *SupplierSyncService) SyncAllPrices(supplierID uint) error {
	driver := s.GetDriver(supplierID)
	if driver == nil {
		return fmt.Errorf("供应商驱动未注册: %d", supplierID)
	}

	products, err := driver.FetchProducts()
	if err != nil {
		return err
	}

	db := database.GetDB()
	for _, p := range products {
		var existing model.SupplierProduct
		if err := db.Where("supplier_id = ? AND remote_product_id = ?", supplierID, p.ID).First(&existing).Error; err != nil {
			continue
		}

		// 应用利润率（MD 7.2.4：利润率后台可配置，默认25%）
		profitRate := existing.ProfitRate
		if profitRate <= 0 {
			profitRate = float64(getSettingInt("default_profit_rate", 25))
		}
		localPrice := p.Price * (1 + profitRate/100)

		// 更新
		db.Model(&existing).Updates(map[string]interface{}{
			"remote_price": p.Price,
			"local_price":  localPrice,
		})
	}
	return nil
}

// getSettingInt 从settings表读取整数配置
func getSettingInt(key string, defaultVal int) int {
	db := database.GetDB()
	var setting model.Setting
	if err := db.Where("`key` = ?", key).First(&setting).Error; err == nil {
		var v int
		if _, err := fmt.Sscanf(setting.Value, "%d", &v); err == nil && v > 0 {
			return v
		}
	}
	return defaultVal
}

// StartPriceSyncCron 启动价格同步定时任务（MD 7.2.4：频率后台可配置，默认每小时）
func (s *SupplierSyncService) StartPriceSyncCron() {
	go func() {
		for {
			// 从settings读取同步频率（分钟）
			intervalMin := getSettingInt("price_sync_interval", 60)
			time.Sleep(time.Duration(intervalMin) * time.Minute)

			// 检查是否启用
			if getSettingInt("price_sync_enabled", 1) == 0 {
				continue
			}

			db := database.GetDB()
			var suppliers []model.Supplier
			db.Where("status = ?", "active").Find(&suppliers)
			for _, supplier := range suppliers {
				if s.GetDriver(supplier.ID) != nil {
					s.SyncAllPrices(supplier.ID)
				}
			}
		}
	}()
}

// StartStockSyncCron 启动库存同步定时任务（MD 7.2.4：频率后台可配置，默认每5分钟）
func (s *SupplierSyncService) StartStockSyncCron() {
	go func() {
		for {
			// 从settings读取同步频率（分钟）
			intervalMin := getSettingInt("stock_sync_interval", 5)
			time.Sleep(time.Duration(intervalMin) * time.Minute)

			// 检查是否启用
			if getSettingInt("stock_sync_enabled", 1) == 0 {
				continue
			}

			db := database.GetDB()
			var suppliers []model.Supplier
			db.Where("status = ?", "active").Find(&suppliers)
			for _, supplier := range suppliers {
				if s.GetDriver(supplier.ID) != nil {
					s.SyncAllProducts(supplier.ID)
				}
			}
		}
	}()
}

// autoCreateGroups 自动创建商品分组（MD 7.2.2）
func (s *SupplierSyncService) autoCreateGroups(groups []RemoteGroup) {
	db := database.GetDB()
	for _, g := range groups {
		var existing model.ProductGroup
		if err := db.Where("name = ? AND parent_id = 0", g.Name).First(&existing).Error; err != nil {
			existing = model.ProductGroup{Name: g.Name, ParentID: 0, Status: "active"}
			db.Create(&existing)
		}
		for _, child := range g.Children {
			var existingChild model.ProductGroup
			if err := db.Where("name = ? AND parent_id = ?", child.Name, existing.ID).First(&existingChild).Error; err != nil {
				existingChild = model.ProductGroup{Name: child.Name, ParentID: existing.ID, Status: "active"}
				db.Create(&existingChild)
			}
		}
	}
}

// SyncProductStatus 同步商品状态（MD 7.2.4：上游隐藏/删除时本地自动同步）
func (s *SupplierSyncService) SyncProductStatus(supplierID uint) error {
	driver := s.GetDriver(supplierID)
	if driver == nil {
		return fmt.Errorf("供应商驱动未注册: %d", supplierID)
	}

	remoteProducts, err := driver.FetchProducts()
	if err != nil {
		return err
	}

	remoteMap := make(map[string]bool)
	for _, rp := range remoteProducts {
		remoteMap[rp.ID] = true
	}

	db := database.GetDB()
	var localProducts []model.SupplierProduct
	db.Where("supplier_id = ?", supplierID).Find(&localProducts)

	for _, lp := range localProducts {
		_, exists := remoteMap[lp.RemoteProductID]
		if !exists && lp.Status == "active" {
			db.Model(&lp).Update("status", "disabled")
		} else if exists && lp.Status == "disabled" {
			db.Model(&lp).Update("status", "active")
		}
	}
	return nil
}
