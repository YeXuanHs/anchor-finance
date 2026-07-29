package service

import (
	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type LinkCauseService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewLinkCauseService(db *gorm.DB, log *logger.Logger) *LinkCauseService {
	return &LinkCauseService{db: db, log: log}
}

// GetList returns paginated link causes.
func (s *LinkCauseService) GetList(page, pageSize int, linkType string, keyword string) ([]model.LinkCause, int64, error) {
	var causes []model.LinkCause
	var total int64

	query := s.db.Model(&model.LinkCause{})
	if linkType != "" {
		query = query.Where("link_type = ?", linkType)
	}
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("sort_order ASC, id ASC").Find(&causes).Error; err != nil {
		return nil, 0, err
	}
	return causes, total, nil
}

// GetByID returns a single cause by ID.
func (s *LinkCauseService) GetByID(id uint) (*model.LinkCause, error) {
	var cause model.LinkCause
	if err := s.db.First(&cause, id).Error; err != nil {
		return nil, err
	}
	return &cause, nil
}

// GetTree returns link causes as a tree structure.
func (s *LinkCauseService) GetTree(linkType string) ([]model.LinkCause, error) {
	var causes []model.LinkCause
	query := s.db.Where("parent_id IS NULL AND status = 1")
	if linkType != "" {
		query = query.Where("link_type = ?", linkType)
	}
	if err := query.Preload("Children", "status = 1").Order("sort_order ASC").Find(&causes).Error; err != nil {
		return nil, err
	}
	return causes, nil
}

// Create creates a new link cause.
func (s *LinkCauseService) Create(cause *model.LinkCause) error {
	return s.db.Create(cause).Error
}

// Update updates a link cause.
func (s *LinkCauseService) Update(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.LinkCause{}).Where("id = ?", id).Updates(updates).Error
}

// Delete deletes a link cause.
func (s *LinkCauseService) Delete(id uint) error {
	return s.db.Delete(&model.LinkCause{}, id).Error
}
