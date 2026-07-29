package service

import (
	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type LinkKnowledgeService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewLinkKnowledgeService(db *gorm.DB, log *logger.Logger) *LinkKnowledgeService {
	return &LinkKnowledgeService{db: db, log: log}
}

// GetList returns paginated link knowledges.
func (s *LinkKnowledgeService) GetList(page, pageSize int, knowledgeType string, category string, keyword string) ([]model.LinkKnowledge, int64, error) {
	var items []model.LinkKnowledge
	var total int64

	query := s.db.Model(&model.LinkKnowledge{})
	if knowledgeType != "" {
		query = query.Where("type = ?", knowledgeType)
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("sort_order ASC, id DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// GetByID returns a single knowledge by ID.
func (s *LinkKnowledgeService) GetByID(id uint) (*model.LinkKnowledge, error) {
	var item model.LinkKnowledge
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// GetByCause returns knowledges linked to a cause.
func (s *LinkKnowledgeService) GetByCause(causeID uint) ([]model.LinkKnowledge, error) {
	var items []model.LinkKnowledge
	if err := s.db.Where("link_cause = ? AND status = 1", causeID).Order("sort_order ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// Create creates a new link knowledge.
func (s *LinkKnowledgeService) Create(item *model.LinkKnowledge) error {
	return s.db.Create(item).Error
}

// Update updates a link knowledge.
func (s *LinkKnowledgeService) Update(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.LinkKnowledge{}).Where("id = ?", id).Updates(updates).Error
}

// Delete deletes a link knowledge.
func (s *LinkKnowledgeService) Delete(id uint) error {
	return s.db.Delete(&model.LinkKnowledge{}, id).Error
}
