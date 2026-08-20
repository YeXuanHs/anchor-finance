package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"anchorfinance/internal/model"

	"gorm.io/gorm"
)

// AutoAssignRule defines the structure for auto-assignment rules.
type AutoAssignRule struct {
	MinTotalSpent float64 `json:"min_total_spent"` // 最低累计消费
	MaxTotalSpent float64 `json:"max_total_spent"` // 最高累计消费 (0=不限)
	MinOrderCount int     `json:"min_order_count"` // 最低订单数
	RegisteredDays int    `json:"registered_days"` // 注册天数要求
}

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
	Name            string          `json:"name" binding:"required,max=100"`
	Description     string          `json:"description"`
	Discount        float64         `json:"discount" binding:"gte=0,lte=1"`
	DiscountPercent float64         `json:"discount_percent" binding:"gte=0,lte=100"`
	TaxRate         float64         `json:"tax_rate" binding:"gte=0,lte=100"`
	AutoAssignRule  *AutoAssignRule `json:"auto_assign_rule"`
	Priority        int             `json:"priority"`
}

// UpdateGroupRequest is the payload for updating a client group.
type UpdateGroupRequest struct {
	Name            *string         `json:"name"`
	Description     *string         `json:"description"`
	Discount        *float64        `json:"discount"`
	DiscountPercent *float64        `json:"discount_percent"`
	TaxRate         *float64        `json:"tax_rate"`
	AutoAssignRule  *AutoAssignRule `json:"auto_assign_rule"`
	Priority        *int            `json:"priority"`
	IsActive        *bool           `json:"is_active"`
}

// Create creates a new client group.
func (s *ClientGroupService) Create(req CreateGroupRequest) (*model.ClientGroup, error) {
	group := &model.ClientGroup{
		Name:            req.Name,
		Description:     req.Description,
		Discount:        req.Discount,
		DiscountPercent: req.DiscountPercent,
		TaxRate:         req.TaxRate,
		Priority:        req.Priority,
		IsActive:        true,
	}
	if group.Discount == 0 {
		group.Discount = 1.0
	}

	if req.AutoAssignRule != nil {
		ruleJSON, err := json.Marshal(req.AutoAssignRule)
		if err != nil {
			return nil, fmt.Errorf("marshal auto assign rule: %w", err)
		}
		group.AutoAssignRule = ruleJSON
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
	if req.DiscountPercent != nil {
		updates["discount_percent"] = *req.DiscountPercent
	}
	if req.TaxRate != nil {
		updates["tax_rate"] = *req.TaxRate
	}
	if req.AutoAssignRule != nil {
		ruleJSON, err := json.Marshal(req.AutoAssignRule)
		if err != nil {
			return fmt.Errorf("marshal auto assign rule: %w", err)
		}
		updates["auto_assign_rule"] = ruleJSON
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

// AutoAssignGroups evaluates all groups with auto-assign rules and assigns users accordingly.
// This should be called periodically (e.g., daily cron) or after significant user events.
func (s *ClientGroupService) AutoAssignGroups() (string, error) {
	var groups []model.ClientGroup
	if err := s.db.Where("is_active = true AND auto_assign_rule IS NOT NULL").Find(&groups).Error; err != nil {
		return "", fmt.Errorf("query groups with rules: %w", err)
	}

	totalAssigned := 0
	for _, group := range groups {
		if group.AutoAssignRule == nil {
			continue
		}

		var rule AutoAssignRule
		if err := json.Unmarshal(group.AutoAssignRule, &rule); err != nil {
			continue
		}

		// Build query for eligible users based on rules
		query := s.db.Table("users").Select("users.id").
			Joins("LEFT JOIN orders ON orders.user_id = users.id AND orders.deleted_at IS NULL")

		if rule.RegisteredDays > 0 {
			query = query.Where("users.created_at <= NOW() - INTERVAL ? DAY", rule.RegisteredDays)
		}

		query = query.Group("users.id")

		if rule.MinOrderCount > 0 {
			query = query.Having("COUNT(orders.id) >= ?", rule.MinOrderCount)
		}

		if rule.MinTotalSpent > 0 {
			query = query.Having("COALESCE(SUM(orders.total), 0) >= ?", rule.MinTotalSpent)
		}

		if rule.MaxTotalSpent > 0 {
			query = query.Having("COALESCE(SUM(orders.total), 0) <= ?", rule.MaxTotalSpent)
		}

		var userIDs []uint
		if err := query.Pluck("users.id", &userIDs).Error; err != nil {
			continue
		}

		// Assign users to group
		for _, uid := range userIDs {
			if err := s.AddMember(group.ID, uid); err != nil {
				continue
			}
			totalAssigned++
		}
	}

	output := fmt.Sprintf("Evaluated %d groups with auto-assign rules, assigned %d users", len(groups), totalAssigned)
	return output, nil
}

// GetUserGroups returns all groups a user belongs to.
func (s *ClientGroupService) GetUserGroups(userID uint) ([]model.ClientGroup, error) {
	var groups []model.ClientGroup
	err := s.db.Joins("JOIN client_group_members ON client_group_members.group_id = client_groups.id").
		Where("client_group_members.user_id = ? AND client_groups.is_active = true", userID).
		Order("client_groups.priority DESC").
		Find(&groups).Error
	return groups, err
}

// CalculateDiscount calculates the effective discount for a user based on their group membership.
func (s *ClientGroupService) CalculateDiscount(userID uint) (float64, float64, error) {
	groups, err := s.GetUserGroups(userID)
	if err != nil {
		return 0, 0, err
	}

	if len(groups) == 0 {
		return 1.0, 0, nil // no discount, no tax
	}

	// Use the best discount from highest priority group
	bestDiscount := 1.0
	taxRate := 0.0
	for _, g := range groups {
		if g.Discount < bestDiscount {
			bestDiscount = g.Discount
		}
		if g.DiscountPercent > 0 {
			discountFromPercent := 1.0 - g.DiscountPercent/100
			if discountFromPercent < bestDiscount {
				bestDiscount = discountFromPercent
			}
		}
		if g.TaxRate > taxRate {
			taxRate = g.TaxRate
		}
	}

	return bestDiscount, taxRate, nil
}
