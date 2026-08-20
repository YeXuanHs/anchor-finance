package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// Developer 开发者
type Developer struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	UserID      uint       `gorm:"uniqueIndex;not null" json:"user_id"`
	Name        string     `gorm:"size:128;not null" json:"name"`
	Company     string     `gorm:"size:256" json:"company"`
	Email       string     `gorm:"size:256;not null" json:"email"`
	Phone       string     `gorm:"size:32" json:"phone"`
	Website     string     `gorm:"size:512" json:"website"`
	Description string     `gorm:"type:text" json:"description"`
	Status      string     `gorm:"size:32;default:pending;comment:pending/approved/rejected/disabled" json:"status"`
	ApprovedAt  *time.Time `json:"approved_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// DeveloperAPIKey 开发者API密钥
type DeveloperAPIKey struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	DeveloperID uint       `gorm:"index;not null" json:"developer_id"`
	Name        string     `gorm:"size:128;not null" json:"name"`
	KeyID       string     `gorm:"size:64;uniqueIndex;not null" json:"key_id"`
	KeySecret   string     `gorm:"size:256;not null" json:"key_secret"`
	Scopes      string     `gorm:"type:text" json:"scopes"` // JSON array
	ExpiresAt   *time.Time `json:"expires_at"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	Status      string     `gorm:"size:32;default:active;comment:active/revoked" json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
}

// DeveloperBilling 开发者计费
type DeveloperBilling struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	DeveloperID uint      `gorm:"index;not null" json:"developer_id"`
	Amount      float64   `gorm:"type:decimal(12,2);not null" json:"amount"`
	Type        string    `gorm:"size:32;comment:api_call/storage/transfer" json:"type"`
	Description string    `gorm:"size:512" json:"description"`
	Status      string    `gorm:"size:32;default:pending;comment:pending/paid" json:"status"`
	Period      string    `gorm:"size:32" json:"period"` // e.g. "2024-01"
	CreatedAt   time.Time `json:"created_at"`
	PaidAt      *time.Time `json:"paid_at"`
}

// DeveloperDoc 开发者文档
type DeveloperDoc struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:256;not null" json:"title"`
	Slug      string    `gorm:"size:256;uniqueIndex" json:"slug"`
	Content   string    `gorm:"type:longtext" json:"content"`
	Category  string    `gorm:"size:64" json:"category"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	Status    string    `gorm:"size:32;default:published;comment:draft/published" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DeveloperService 开发者服务
type DeveloperService struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewDeveloperService creates a new DeveloperService.
func NewDeveloperService(db *gorm.DB, log *logger.Logger) *DeveloperService {
	return &DeveloperService{db: db, log: log}
}

// List returns paginated developers.
func (s *DeveloperService) List(page, pageSize int) ([]Developer, int64, error) {
	var items []Developer
	var total int64

	query := s.db.Model(&Developer{})
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// GetDetail returns a single developer with related info.
func (s *DeveloperService) GetDetail(id uint) (map[string]interface{}, error) {
	var dev Developer
	if err := s.db.First(&dev, id).Error; err != nil {
		return nil, err
	}

	var apiKeyCount int64
	s.db.Model(&DeveloperAPIKey{}).Where("developer_id = ? AND status = 'active'", id).Count(&apiKeyCount)

	var totalBilling float64
	s.db.Model(&DeveloperBilling{}).Where("developer_id = ?", id).
		Select("COALESCE(SUM(amount), 0)").Scan(&totalBilling)

	return map[string]interface{}{
		"developer":        dev,
		"active_api_keys":  apiKeyCount,
		"total_billing":    totalBilling,
	}, nil
}

// Create creates a new developer.
func (s *DeveloperService) Create(dev *Developer) error {
	dev.Status = "pending"
	return s.db.Create(dev).Error
}

// Update updates a developer.
func (s *DeveloperService) Update(id uint, updates map[string]interface{}) error {
	return s.db.Model(&Developer{}).Where("id = ?", id).Updates(updates).Error
}

// Delete deletes a developer.
func (s *DeveloperService) Delete(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Revoke all API keys
		if err := tx.Model(&DeveloperAPIKey{}).Where("developer_id = ?", id).
			Update("status", "revoked").Error; err != nil {
			return err
		}
		return tx.Delete(&Developer{}, id).Error
	})
}

// Approve approves a developer.
func (s *DeveloperService) Approve(id uint) error {
	now := time.Now()
	return s.db.Model(&Developer{}).Where("id = ? AND status = 'pending'", id).
		Updates(map[string]interface{}{
			"status":      "approved",
			"approved_at": &now,
		}).Error
}

// Reject rejects a developer.
func (s *DeveloperService) Reject(id uint, reason string) error {
	return s.db.Model(&Developer{}).Where("id = ? AND status = 'pending'", id).
		Update("status", "rejected").Error
}

// GetAPIKeys returns API keys for a developer.
func (s *DeveloperService) GetAPIKeys(developerID uint) ([]DeveloperAPIKey, error) {
	var items []DeveloperAPIKey
	if err := s.db.Where("developer_id = ?", developerID).Order("id DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// GenerateAPIKey generates a new API key for a developer.
func (s *DeveloperService) GenerateAPIKey(developerID uint, name string, scopes string) (*DeveloperAPIKey, error) {
	var dev Developer
	if err := s.db.First(&dev, developerID).Error; err != nil {
		return nil, errors.New("developer not found")
	}
	if dev.Status != "approved" {
		return nil, errors.New("developer is not approved")
	}

	keyIDBytes := make([]byte, 16)
	if _, err := rand.Read(keyIDBytes); err != nil {
		return nil, err
	}
	keyID := "ak_" + hex.EncodeToString(keyIDBytes)

	secretBytes := make([]byte, 32)
	if _, err := rand.Read(secretBytes); err != nil {
		return nil, err
	}
	keySecret := hex.EncodeToString(secretBytes)

	apiKey := &DeveloperAPIKey{
		DeveloperID: developerID,
		Name:        name,
		KeyID:       keyID,
		KeySecret:   keySecret,
		Scopes:      scopes,
		Status:      "active",
	}
	if err := s.db.Create(apiKey).Error; err != nil {
		return nil, err
	}

	s.log.Infof("API key generated for developer %d: %s", developerID, keyID)
	return apiKey, nil
}

// RevokeAPIKey revokes an API key.
func (s *DeveloperService) RevokeAPIKey(developerID, keyID uint) error {
	result := s.db.Model(&DeveloperAPIKey{}).
		Where("id = ? AND developer_id = ?", keyID, developerID).
		Update("status", "revoked")
	if result.RowsAffected == 0 {
		return errors.New("API key not found")
	}
	return result.Error
}

// GetDocs returns all developer docs.
func (s *DeveloperService) GetDocs() ([]DeveloperDoc, error) {
	var items []DeveloperDoc
	if err := s.db.Where("status = 'published'").Order("sort_order ASC, id ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// UpdateDocs batch updates developer docs.
func (s *DeveloperService) UpdateDocs(docs []DeveloperDoc) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for i := range docs {
			if docs[i].ID > 0 {
				if err := tx.Model(&DeveloperDoc{}).Where("id = ?", docs[i].ID).
					Updates(map[string]interface{}{
						"title":      docs[i].Title,
						"content":    docs[i].Content,
						"category":   docs[i].Category,
						"sort_order": docs[i].SortOrder,
						"status":     docs[i].Status,
					}).Error; err != nil {
					return err
				}
			} else {
				if err := tx.Create(&docs[i]).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// GetBilling returns paginated billing records for a developer.
func (s *DeveloperService) GetBilling(page, pageSize int) ([]DeveloperBilling, int64, error) {
	var items []DeveloperBilling
	var total int64

	query := s.db.Model(&DeveloperBilling{})
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// SettleBilling settles pending billing for a developer.
func (s *DeveloperService) SettleBilling(developerID uint, period string) error {
	now := time.Now()
	result := s.db.Model(&DeveloperBilling{}).
		Where("developer_id = ? AND status = 'pending'", developerID)
	if period != "" {
		result = result.Where("period = ?", period)
	}
	result = result.Updates(map[string]interface{}{
		"status":  "paid",
		"paid_at": &now,
	})

	if result.RowsAffected == 0 {
		return fmt.Errorf("no pending billing records to settle")
	}

	s.log.Infof("billing settled for developer %d (period=%s)", developerID, period)
	return result.Error
}
