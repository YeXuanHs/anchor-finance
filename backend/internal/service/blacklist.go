package service

import (
	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"
	"gorm.io/gorm"
)

// BlacklistService manages blacklist operations.
type BlacklistService struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewBlacklistService creates a new BlacklistService.
func NewBlacklistService(db *gorm.DB, log *logger.Logger) *BlacklistService {
	return &BlacklistService{db: db, log: log}
}

// List returns all blacklist entries with pagination.
func (s *BlacklistService) List(page, pageSize int, blacklistType string) ([]model.Blacklist, int64, error) {
	var items []model.Blacklist
	var total int64

	query := s.db.Model(&model.Blacklist{})
	if blacklistType != "" {
		query = query.Where("type = ?", blacklistType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// Create adds a new blacklist entry.
func (s *BlacklistService) Create(item *model.Blacklist) error {
	return s.db.Create(item).Error
}

// Update updates a blacklist entry.
func (s *BlacklistService) Update(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.Blacklist{}).Where("id = ?", id).Updates(updates).Error
}

// Delete removes a blacklist entry.
func (s *BlacklistService) Delete(id uint) error {
	return s.db.Delete(&model.Blacklist{}, id).Error
}

// IsBlacklisted checks if a value is blacklisted.
func (s *BlacklistService) IsBlacklisted(blacklistType, value string) (bool, error) {
	var count int64
	err := s.db.Model(&model.Blacklist{}).
		Where("type = ? AND value = ? AND (expires_at IS NULL OR expires_at > NOW())", blacklistType, value).
		Count(&count).Error
	return count > 0, err
}
