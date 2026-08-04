package service

import (
	"errors"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// UserRemarkService manages user remark operations.
type UserRemarkService struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewUserRemarkService creates a new UserRemarkService.
func NewUserRemarkService(db *gorm.DB, log *logger.Logger) *UserRemarkService {
	return &UserRemarkService{db: db, log: log}
}

// AddRemarkRequest is the payload for adding a remark.
type AddRemarkRequest struct {
	UserID  uint   `json:"user_id" binding:"required"`
	Content string `json:"content" binding:"required"`
	Stick   int    `json:"stick"`
	Type    string `json:"type" binding:"omitempty,oneof=general billing support vip risk"`
}

// UpdateRemarkRequest is the payload for updating a remark.
type UpdateRemarkRequest struct {
	Content string `json:"content"`
	Stick   *int   `json:"stick"`
}

// Add creates a new remark on a user.
func (s *UserRemarkService) Add(adminID uint, req AddRemarkRequest) (*model.UserRemark, error) {
	if req.Type == "" {
		req.Type = "general"
	}
	remark := &model.UserRemark{
		UserID:  req.UserID,
		AdminID: adminID,
		Content: req.Content,
		Stick:   req.Stick,
		Type:    req.Type,
	}
	if err := s.db.Create(remark).Error; err != nil {
		return nil, err
	}
	return remark, nil
}

// GetByID fetches a remark by ID.
func (s *UserRemarkService) GetByID(id uint) (*model.UserRemark, error) {
	var remark model.UserRemark
	if err := s.db.Preload("Admin").First(&remark, id).Error; err != nil {
		return nil, err
	}
	return &remark, nil
}

// GetByUser returns a paginated list of remarks for a user, optionally filtered by type.
func (s *UserRemarkService) GetByUser(userID uint, remarkType string, page, pageSize int) ([]model.UserRemark, int64, error) {
	var items []model.UserRemark
	var total int64

	query := s.db.Model(&model.UserRemark{}).Where("user_id = ?", userID)
	if remarkType != "" {
		query = query.Where("type = ?", remarkType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("Admin").
		Offset(offset).Limit(pageSize).Order("stick DESC, id DESC").
		Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// Update updates a remark's content and/or stick value.
func (s *UserRemarkService) Update(id uint, adminID uint, req UpdateRemarkRequest) (*model.UserRemark, error) {
	var remark model.UserRemark
	if err := s.db.Where("id = ? AND admin_id = ?", id, adminID).First(&remark).Error; err != nil {
		return nil, errors.New("remark not found or not owned by this admin")
	}

	updates := map[string]interface{}{}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Stick != nil {
		updates["stick"] = *req.Stick
	}
	if len(updates) > 0 {
		if err := s.db.Model(&remark).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	s.db.First(&remark, id)
	return &remark, nil
}

// Delete removes a remark.
func (s *UserRemarkService) Delete(id uint, adminID uint) error {
	result := s.db.Where("id = ? AND admin_id = ?", id, adminID).Delete(&model.UserRemark{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("remark not found or not owned by this admin")
	}
	return nil
}

// AdminDelete removes any remark (super-admin).
func (s *UserRemarkService) AdminDelete(id uint) error {
	result := s.db.Delete(&model.UserRemark{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("remark not found")
	}
	return nil
}
