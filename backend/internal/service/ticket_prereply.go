package service

import (
	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type TicketPrereplyService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewTicketPrereplyService(db *gorm.DB, log *logger.Logger) *TicketPrereplyService {
	return &TicketPrereplyService{db: db, log: log}
}

// GetReplyList returns all prereply categories with replies.
func (s *TicketPrereplyService) GetReplyList() ([]model.TicketPrereplyCategory, error) {
	var categories []model.TicketPrereplyCategory
	if err := s.db.Preload("Replies", "status = 1").
		Where("status = 1").
		Order("sort_order ASC").
		Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// AddCategory adds a new prereply category.
func (s *TicketPrereplyService) AddCategory(cat *model.TicketPrereplyCategory) error {
	return s.db.Create(cat).Error
}

// UpdateCategory updates a prereply category.
func (s *TicketPrereplyService) UpdateCategory(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.TicketPrereplyCategory{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteCategory deletes a prereply category.
func (s *TicketPrereplyService) DeleteCategory(id uint) error {
	return s.db.Delete(&model.TicketPrereplyCategory{}, id).Error
}

// AddReply adds a new prereply.
func (s *TicketPrereplyService) AddReply(reply *model.TicketPrereply) error {
	return s.db.Create(reply).Error
}

// UpdateReply updates a prereply.
func (s *TicketPrereplyService) UpdateReply(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.TicketPrereply{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteReply deletes a prereply.
func (s *TicketPrereplyService) DeleteReply(id uint) error {
	return s.db.Delete(&model.TicketPrereply{}, id).Error
}

// GetReplyByID returns a single prereply by ID.
func (s *TicketPrereplyService) GetReplyByID(id uint) (*model.TicketPrereply, error) {
	var reply model.TicketPrereply
	if err := s.db.First(&reply, id).Error; err != nil {
		return nil, err
	}
	return &reply, nil
}

// IncrementUseCount increments the use count of a prereply.
func (s *TicketPrereplyService) IncrementUseCount(id uint) error {
	return s.db.Model(&model.TicketPrereply{}).Where("id = ?", id).
		UpdateColumn("use_count", gorm.Expr("use_count + 1")).Error
}

// GetCategoryByID returns a single category by ID.
func (s *TicketPrereplyService) GetCategoryByID(id uint) (*model.TicketPrereplyCategory, error) {
	var category model.TicketPrereplyCategory
	if err := s.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// GetAllCategories returns all categories for dropdown selection.
func (s *TicketPrereplyService) GetAllCategories() ([]model.TicketPrereplyCategory, error) {
	var categories []model.TicketPrereplyCategory
	if err := s.db.Where("status = 1").Order("sort_order ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// SearchReplies searches prereplies by title and content.
func (s *TicketPrereplyService) SearchReplies(title, content string) ([]model.TicketPrereply, error) {
	var replies []model.TicketPrereply
	query := s.db.Where("status = 1")
	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	if content != "" {
		query = query.Where("content LIKE ?", "%"+content+"%")
	}
	if err := query.Order("use_count DESC").Find(&replies).Error; err != nil {
		return nil, err
	}
	return replies, nil
}
