package service

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// APIManageService manages API keys for admin.
type APIManageService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewAPIManageService(db *gorm.DB, log *logger.Logger) *APIManageService {
	return &APIManageService{db: db, log: log}
}

func (s *APIManageService) List(page, pageSize int, keyword string, status *int16) ([]model.APIKey, int64, error) {
	var items []model.APIKey
	var total int64

	query := s.db.Model(&model.APIKey{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (s *APIManageService) GetByID(id uint) (*model.APIKey, error) {
	var item model.APIKey
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *APIManageService) Create(apiKey *model.APIKey) error {
	apiKey.Secret = generateAPISecret()
	apiKey.CreatedAt = time.Now()
	return s.db.Create(apiKey).Error
}

func (s *APIManageService) Update(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.APIKey{}).Where("id = ?", id).Updates(updates).Error
}

func (s *APIManageService) Delete(id uint) error {
	return s.db.Delete(&model.APIKey{}, id).Error
}

func (s *APIManageService) SetStatus(id uint, status int16) error {
	return s.db.Model(&model.APIKey{}).Where("id = ?", id).Update("status", status).Error
}

func (s *APIManageService) Regenerate(id uint) (string, error) {
	newSecret := generateAPISecret()
	if err := s.db.Model(&model.APIKey{}).Where("id = ?", id).Update("secret", newSecret).Error; err != nil {
		return "", err
	}
	return newSecret, nil
}

func (s *APIManageService) GetPermissions() ([]string, error) {
	return []string{
		"products:read", "products:write",
		"orders:read", "orders:write",
		"users:read", "users:write",
		"tickets:read", "tickets:write",
		"invoices:read", "invoices:write",
		"hosts:read", "hosts:write",
	}, nil
}

func generateAPISecret() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
