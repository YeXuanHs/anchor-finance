package service

import (
	"fmt"
	"strings"

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
	// Remove non-alphanumeric characters except hyphens.
	var b strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		}
	}
	slug = b.String()
	// Collapse multiple hyphens.
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
// It checks the provider config for "group_id", otherwise creates/reuses a
// default "上游产品" (Upstream Products) group.
func (s *UpstreamService) ensureGroupID(provider *model.UpstreamProvider) (uint, error) {
	// Check if the provider config specifies a group_id.
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

	// Find or create the default upstream product group.
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

// SyncProducts fetches products from the upstream and upserts them as
// UpstreamProduct mappings. Returns the number of products synced.
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
