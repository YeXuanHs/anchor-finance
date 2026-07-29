package service

import (
	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type AdvancedOptionsService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewAdvancedOptionsService(db *gorm.DB, log *logger.Logger) *AdvancedOptionsService {
	return &AdvancedOptionsService{db: db, log: log}
}

// GetOptions returns advanced config options for a product.
func (s *AdvancedOptionsService) GetOptions(productID uint) ([]model.AdvancedOption, error) {
	var options []model.AdvancedOption
	if err := s.db.Where("product_id = ? AND status = 1", productID).Order("sort_order ASC").Find(&options).Error; err != nil {
		return nil, err
	}
	return options, nil
}

// GetOptionByID returns a single option by ID.
func (s *AdvancedOptionsService) GetOptionByID(id uint) (*model.AdvancedOption, error) {
	var option model.AdvancedOption
	if err := s.db.First(&option, id).Error; err != nil {
		return nil, err
	}
	return &option, nil
}

// CreateOption creates a new advanced config option.
func (s *AdvancedOptionsService) CreateOption(opt *model.AdvancedOption) error {
	return s.db.Create(opt).Error
}

// UpdateOption updates an advanced config option.
func (s *AdvancedOptionsService) UpdateOption(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.AdvancedOption{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteOption deletes an advanced config option.
func (s *AdvancedOptionsService) DeleteOption(id uint) error {
	return s.db.Delete(&model.AdvancedOption{}, id).Error
}

// GetLinks returns config links for a product.
func (s *AdvancedOptionsService) GetLinks(productID uint) ([]model.AdvancedOptionLink, error) {
	var links []model.AdvancedOptionLink
	if err := s.db.Where("product_id = ? AND status = 1", productID).Order("sort_order ASC").Find(&links).Error; err != nil {
		return nil, err
	}
	return links, nil
}

// GetLinkByID returns a single link by ID.
func (s *AdvancedOptionsService) GetLinkByID(id uint) (*model.AdvancedOptionLink, error) {
	var link model.AdvancedOptionLink
	if err := s.db.First(&link, id).Error; err != nil {
		return nil, err
	}
	return &link, nil
}

// CreateLink creates a new config link.
func (s *AdvancedOptionsService) CreateLink(link *model.AdvancedOptionLink) error {
	return s.db.Create(link).Error
}

// UpdateLink updates a config link.
func (s *AdvancedOptionsService) UpdateLink(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.AdvancedOptionLink{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteLink deletes a config link.
func (s *AdvancedOptionsService) DeleteLink(id uint) error {
	return s.db.Delete(&model.AdvancedOptionLink{}, id).Error
}
