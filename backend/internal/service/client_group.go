package service

import (
	"errors"
	"fmt"

	"anchorfinance/internal/model"

	"gorm.io/gorm"
)

// ClientGroupService manages client group operations.
type ClientGroupService struct {
	db *gorm.DB
}

// NewClientGroupService creates a new ClientGroupService.
func NewClientGroupService(db *gorm.DB) *ClientGroupService {
	return &ClientGroupService{db: db}
}

// CreateGroupRequest is the payload for creating a client group.
type CreateGroupRequest struct {
	Name        string  `json:"name" binding:"required,max=100"`
	Description string  `json:"description"`
	Discount    float64 `json:"discount" binding:"gte=0,lte=1"`
	Priority    int     `json:"priority"`
}

// UpdateGroupRequest is the payload for updating a client group.
type UpdateGroupRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Discount    *float64 `json:"discount"`
	Priority    *int     `json:"priority"`
	IsActive    *bool    `json:"is_active"`
}

// Create creates a new client group.
func (s *ClientGroupService) Create(req CreateGroupRequest) (*model.ClientGroup, error) {
	group := &model.ClientGroup{
		Name:        req.Name,
		Description: req.Description,
		Discount:    req.Discount,
		Priority:    req.Priority,
		IsActive:    true,
	}
	if group.Discount == 0 {
		group.Discount = 1.0
	}
	if err := s.db.Create(group).Error; err != nil {
		return nil, err
	}
	return group, nil
}

// GetByID fetches a client group by ID.
func (s *ClientGroupService) GetByID(id uint) (*model.ClientGroup, error) {
	var group model.ClientGroup
	if err := s.db.First(&group, id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

// GetList returns a paginated client group list.
func (s *ClientGroupService) GetList(page, pageSize int, keyword string) ([]model.ClientGroup, int64, error) {
	var groups []model.ClientGroup
	var total int64

	query := s.db.Model(&model.ClientGroup{})
	if keyword != "" {
		q := "%" + keyword + "%"
		query = query.Where("name LIKE ?", q)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("priority DESC, id ASC").Find(&groups).Error; err != nil {
		return nil, 0, err
	}
	return groups, total, nil
}

// Update modifies an existing client group.
func (s *ClientGroupService) Update(id uint, req UpdateGroupRequest) error {
	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Discount != nil {
		updates["discount"] = *req.Discount
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if len(updates) == 0 {
		return nil
	}
	return s.db.Model(&model.ClientGroup{}).Where("id = ?", id).Updates(updates).Error
}

// Delete removes a client group by ID.
func (s *ClientGroupService) Delete(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("group_id = ?", id).Delete(&model.ClientGroupMember{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.ClientGroup{}, id).Error
	})
}

// GetMembers returns all user IDs in a group.
func (s *ClientGroupService) GetMembers(groupID uint) ([]model.ClientGroupMember, error) {
	var members []model.ClientGroupMember
	if err := s.db.Preload("User").Where("group_id = ?", groupID).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

// AddMember adds a user to a group. Idempotent: silently skips if already a member.
func (s *ClientGroupService) AddMember(groupID, userID uint) error {
	var count int64
	s.db.Model(&model.ClientGroupMember{}).Where("group_id = ? AND user_id = ?", groupID, userID).Count(&count)
	if count > 0 {
		return nil
	}

	member := &model.ClientGroupMember{
		GroupID: groupID,
		UserID:  userID,
	}
	return s.db.Create(member).Error
}

// RemoveMember removes a user from a group.
func (s *ClientGroupService) RemoveMember(groupID, userID uint) error {
	result := s.db.Where("group_id = ? AND user_id = ?", groupID, userID).Delete(&model.ClientGroupMember{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("member not found")
	}
	return nil
}

// SetMembers replaces all members of a group with the given user IDs.
func (s *ClientGroupService) SetMembers(groupID uint, userIDs []uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("group_id = ?", groupID).Delete(&model.ClientGroupMember{}).Error; err != nil {
			return fmt.Errorf("clear existing members: %w", err)
		}
		for _, uid := range userIDs {
			m := &model.ClientGroupMember{GroupID: groupID, UserID: uid}
			if err := tx.Create(m).Error; err != nil {
				return fmt.Errorf("add member %d: %w", uid, err)
			}
		}
		return nil
	})
}
