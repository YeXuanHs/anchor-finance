package service

import (
	"encoding/json"
	"errors"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"anchorfinance/pkg/logger"
)

// TicketDepartment 工单部门
type TicketDepartment struct {
	ID           uint               `gorm:"primaryKey" json:"id"`
	Name         string             `gorm:"size:128;not null" json:"name"`
	Description  string             `gorm:"type:text" json:"description"`
	Slug         string             `gorm:"uniqueIndex;size:128" json:"slug"`
	ParentID     *uint              `gorm:"index" json:"parent_id"`
	Children     []TicketDepartment `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	ManagerIDs   datatypes.JSON     `gorm:"type:json" json:"manager_ids"`
	MemberIDs    datatypes.JSON     `gorm:"type:json" json:"member_ids"`
	SortOrder    int                `gorm:"default:0;index" json:"sort_order"`
	Status       int                `gorm:"default:1;index" json:"status"` // 1=启用 0=禁用
	AutoAssign   bool               `gorm:"default:false" json:"auto_assign"`
	AssignRule   datatypes.JSON     `gorm:"type:json" json:"assign_rule"`
	EmailNotify  bool               `gorm:"default:true" json:"email_notify"`
	SMSNotify    bool               `gorm:"default:false" json:"sms_notify"`
	AutoReply    string             `gorm:"type:text" json:"auto_reply"`
	TicketPrefix string             `gorm:"size:16" json:"ticket_prefix"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type TicketDepartmentService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewTicketDepartmentService(db *gorm.DB, log *logger.Logger) *TicketDepartmentService {
	return &TicketDepartmentService{db: db, log: log}
}

type CreateTicketDepartmentRequest struct {
	Name         string                 `json:"name" binding:"required,max=128"`
	Description  string                 `json:"description"`
	Slug         string                 `json:"slug" binding:"omitempty,max=128"`
	ParentID     *uint                  `json:"parent_id"`
	ManagerIDs   []uint                 `json:"manager_ids"`
	MemberIDs    []uint                 `json:"member_ids"`
	SortOrder    int                    `json:"sort_order"`
	AutoAssign   bool                   `json:"auto_assign"`
	AssignRule   map[string]interface{} `json:"assign_rule"`
	EmailNotify  bool                   `json:"email_notify"`
	SMSNotify    bool                   `json:"sms_notify"`
	AutoReply    string                 `json:"auto_reply"`
	TicketPrefix string                 `json:"ticket_prefix"`
}

type UpdateTicketDepartmentRequest struct {
	Name         *string                `json:"name"`
	Description  *string                `json:"description"`
	Slug         *string                `json:"slug"`
	ParentID     *uint                  `json:"parent_id"`
	ManagerIDs   []uint                 `json:"manager_ids"`
	MemberIDs    []uint                 `json:"member_ids"`
	SortOrder    *int                   `json:"sort_order"`
	AutoAssign   *bool                  `json:"auto_assign"`
	AssignRule   map[string]interface{} `json:"assign_rule"`
	EmailNotify  *bool                  `json:"email_notify"`
	SMSNotify    *bool                  `json:"sms_notify"`
	AutoReply    *string                `json:"auto_reply"`
	TicketPrefix *string                `json:"ticket_prefix"`
}

// Create creates a new ticket department.
func (s *TicketDepartmentService) Create(req CreateTicketDepartmentRequest) (*TicketDepartment, error) {
	managerJSON, _ := json.Marshal(req.ManagerIDs)
	memberJSON, _ := json.Marshal(req.MemberIDs)

	var ruleJSON datatypes.JSON
	if req.AssignRule != nil {
		ruleJSON, _ = json.Marshal(req.AssignRule)
	}

	dept := &TicketDepartment{
		Name:         req.Name,
		Description:  req.Description,
		Slug:         req.Slug,
		ParentID:     req.ParentID,
		ManagerIDs:   datatypes.JSON(managerJSON),
		MemberIDs:    datatypes.JSON(memberJSON),
		SortOrder:    req.SortOrder,
		Status:       1,
		AutoAssign:   req.AutoAssign,
		AssignRule:   ruleJSON,
		EmailNotify:  req.EmailNotify,
		SMSNotify:    req.SMSNotify,
		AutoReply:    req.AutoReply,
		TicketPrefix: req.TicketPrefix,
	}

	if err := s.db.Create(dept).Error; err != nil {
		return nil, err
	}

	s.log.Infof("ticket department created: %s", dept.Name)
	return dept, nil
}

// GetByID returns a single department by ID.
func (s *TicketDepartmentService) GetByID(id uint) (*TicketDepartment, error) {
	var dept TicketDepartment
	if err := s.db.First(&dept, id).Error; err != nil {
		return nil, err
	}
	return &dept, nil
}

// Update modifies an existing department.
func (s *TicketDepartmentService) Update(id uint, req UpdateTicketDepartmentRequest) (*TicketDepartment, error) {
	dept, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Slug != nil {
		updates["slug"] = *req.Slug
	}
	if req.ParentID != nil {
		updates["parent_id"] = *req.ParentID
	}
	if req.ManagerIDs != nil {
		managerJSON, _ := json.Marshal(req.ManagerIDs)
		updates["manager_ids"] = datatypes.JSON(managerJSON)
	}
	if req.MemberIDs != nil {
		memberJSON, _ := json.Marshal(req.MemberIDs)
		updates["member_ids"] = datatypes.JSON(memberJSON)
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.AutoAssign != nil {
		updates["auto_assign"] = *req.AutoAssign
	}
	if req.AssignRule != nil {
		ruleJSON, _ := json.Marshal(req.AssignRule)
		updates["assign_rule"] = datatypes.JSON(ruleJSON)
	}
	if req.EmailNotify != nil {
		updates["email_notify"] = *req.EmailNotify
	}
	if req.SMSNotify != nil {
		updates["sms_notify"] = *req.SMSNotify
	}
	if req.AutoReply != nil {
		updates["auto_reply"] = *req.AutoReply
	}
	if req.TicketPrefix != nil {
		updates["ticket_prefix"] = *req.TicketPrefix
	}

	if err := s.db.Model(dept).Updates(updates).Error; err != nil {
		return nil, err
	}
	return s.GetByID(id)
}

// Delete soft-deletes a department.
func (s *TicketDepartmentService) Delete(id uint) error {
	// Check if department has children
	var childCount int64
	s.db.Model(&TicketDepartment{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		return errors.New("cannot delete department with children")
	}

	return s.db.Delete(&TicketDepartment{}, id).Error
}

// GetList returns all departments with pagination.
func (s *TicketDepartmentService) GetList(page, pageSize int, status *int, keyword string) ([]TicketDepartment, int64, error) {
	var depts []TicketDepartment
	var total int64

	query := s.db.Model(&TicketDepartment{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if keyword != "" {
		q := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR slug LIKE ?", q, q)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("sort_order ASC, id ASC").Find(&depts).Error; err != nil {
		return nil, 0, err
	}
	return depts, total, nil
}

// GetTree returns departments as a tree structure.
func (s *TicketDepartmentService) GetTree() ([]TicketDepartment, error) {
	var depts []TicketDepartment
	if err := s.db.Where("status = 1").Order("sort_order ASC, id ASC").Find(&depts).Error; err != nil {
		return nil, err
	}

	deptMap := make(map[uint]*TicketDepartment)
	for i := range depts {
		deptMap[depts[i].ID] = &depts[i]
	}

	var roots []TicketDepartment
	for _, dept := range depts {
		if dept.ParentID == nil {
			roots = append(roots, dept)
		} else if parent, ok := deptMap[*dept.ParentID]; ok {
			parent.Children = append(parent.Children, dept)
		}
	}

	return roots, nil
}

// AddMember adds a member to the department.
func (s *TicketDepartmentService) AddMember(deptID, userID uint) error {
	dept, err := s.GetByID(deptID)
	if err != nil {
		return err
	}

	var memberIDs []uint
	json.Unmarshal(dept.MemberIDs, &memberIDs)

	for _, id := range memberIDs {
		if id == userID {
			return errors.New("user already a member")
		}
	}

	memberIDs = append(memberIDs, userID)
	memberJSON, _ := json.Marshal(memberIDs)
	return s.db.Model(dept).Update("member_ids", datatypes.JSON(memberJSON)).Error
}

// RemoveMember removes a member from the department.
func (s *TicketDepartmentService) RemoveMember(deptID, userID uint) error {
	dept, err := s.GetByID(deptID)
	if err != nil {
		return err
	}

	var memberIDs []uint
	json.Unmarshal(dept.MemberIDs, &memberIDs)

	newIDs := make([]uint, 0, len(memberIDs))
	for _, id := range memberIDs {
		if id != userID {
			newIDs = append(newIDs, id)
		}
	}

	memberJSON, _ := json.Marshal(newIDs)
	return s.db.Model(dept).Update("member_ids", datatypes.JSON(memberJSON)).Error
}

// SetManagers sets the managers for a department.
func (s *TicketDepartmentService) SetManagers(deptID uint, managerIDs []uint) error {
	dept, err := s.GetByID(deptID)
	if err != nil {
		return err
	}

	managerJSON, _ := json.Marshal(managerIDs)
	return s.db.Model(dept).Update("manager_ids", datatypes.JSON(managerJSON)).Error
}

// GetMembers returns the member IDs of a department.
func (s *TicketDepartmentService) GetMembers(deptID uint) ([]uint, error) {
	dept, err := s.GetByID(deptID)
	if err != nil {
		return nil, err
	}

	var memberIDs []uint
	json.Unmarshal(dept.MemberIDs, &memberIDs)
	return memberIDs, nil
}

// Enable sets department status to 1.
func (s *TicketDepartmentService) Enable(id uint) error {
	return s.db.Model(&TicketDepartment{}).Where("id = ?", id).Update("status", 1).Error
}

// Disable sets department status to 0.
func (s *TicketDepartmentService) Disable(id uint) error {
	return s.db.Model(&TicketDepartment{}).Where("id = ?", id).Update("status", 0).Error
}

// CustomParam represents a custom field definition.
type CustomParam struct {
	ID          uint   `json:"id"`
	Type        string `json:"type"`
	RelID       uint   `json:"rel_id"`
	FieldName   string `json:"field_name"`
	FieldType   string `json:"field_type"`
	Description string `json:"description"`
	FieldOption string `json:"field_option"`
	RegExpr     string `json:"reg_expr"`
	AdminOnly   int8   `json:"admin_only"`
	Required    int8   `json:"required"`
	SortOrder   int    `json:"sort_order"`
	CreatedAt   string `json:"created_at"`
}

// AddCustomParam adds a custom field to a department.
func (s *TicketDepartmentService) AddCustomParam(deptID uint, fieldName, fieldType, description, fieldOption, regExpr string, adminOnly, required int8, sortOrder int) (*CustomParam, error) {
	// Check department exists
	var dept TicketDepartment
	if err := s.db.First(&dept, deptID).Error; err != nil {
		return nil, errors.New("department not found")
	}

	// Check duplicate field name
	var count int64
	s.db.Table("custom_fields").Where("field_name = ? AND type = ? AND rel_id = ?", fieldName, "ticket", deptID).Count(&count)
	if count > 0 {
		return nil, errors.New("field name already exists")
	}

	field := map[string]interface{}{
		"type":         "ticket",
		"rel_id":       deptID,
		"field_name":   fieldName,
		"field_type":   fieldType,
		"description":  description,
		"field_option": fieldOption,
		"reg_expr":     regExpr,
		"admin_only":   adminOnly,
		"required":     required,
		"sort_order":   sortOrder,
		"created_at":   time.Now(),
		"updated_at":   time.Now(),
	}
	if err := s.db.Table("custom_fields").Create(&field).Error; err != nil {
		return nil, err
	}

	return &CustomParam{
		Type:        "ticket",
		RelID:       deptID,
		FieldName:   fieldName,
		FieldType:   fieldType,
		Description: description,
		FieldOption: fieldOption,
		RegExpr:     regExpr,
		AdminOnly:   adminOnly,
		Required:    required,
		SortOrder:   sortOrder,
	}, nil
}

// GetCustomParam returns a custom field by ID.
func (s *TicketDepartmentService) GetCustomParam(fieldID uint) (*CustomParam, error) {
	var field CustomParam
	if err := s.db.Table("custom_fields").Where("id = ?", fieldID).First(&field).Error; err != nil {
		return nil, err
	}
	return &field, nil
}

// EditCustomParam updates a custom field.
func (s *TicketDepartmentService) EditCustomParam(fieldID, deptID uint, fieldName, fieldType, description, fieldOption, regExpr string, adminOnly, required int8, sortOrder int) error {
	// Check duplicate field name (excluding current)
	var count int64
	s.db.Table("custom_fields").Where("field_name = ? AND type = ? AND rel_id = ? AND id <> ?", fieldName, "ticket", deptID, fieldID).Count(&count)
	if count > 0 {
		return errors.New("field name already exists")
	}

	updates := map[string]interface{}{
		"field_name":   fieldName,
		"field_type":   fieldType,
		"description":  description,
		"field_option": fieldOption,
		"reg_expr":     regExpr,
		"admin_only":   adminOnly,
		"required":     required,
		"sort_order":   sortOrder,
		"updated_at":   time.Now(),
	}
	return s.db.Table("custom_fields").Where("id = ?", fieldID).Updates(updates).Error
}

// DeleteCustomParam deletes a custom field and its values.
func (s *TicketDepartmentService) DeleteCustomParam(fieldID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Delete field values
		if err := tx.Table("custom_field_values").Where("field_id = ?", fieldID).Delete(nil).Error; err != nil {
			return err
		}
		// Delete field definition
		return tx.Table("custom_fields").Where("id = ?", fieldID).Delete(nil).Error
	})
}

// MoveUp moves a department up in sort order.
func (s *TicketDepartmentService) MoveUp(deptID uint) error {
	var dept TicketDepartment
	if err := s.db.First(&dept, deptID).Error; err != nil {
		return errors.New("department not found")
	}

	// Find the previous department
	var prev TicketDepartment
	if err := s.db.Where("sort_order < ?", dept.SortOrder).Order("sort_order DESC").First(&prev).Error; err != nil {
		return nil // Already at top
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Swap sort orders
		if err := tx.Model(&dept).Update("sort_order", prev.SortOrder).Error; err != nil {
			return err
		}
		return tx.Model(&prev).Update("sort_order", dept.SortOrder).Error
	})
}

// MoveDown moves a department down in sort order.
func (s *TicketDepartmentService) MoveDown(deptID uint) error {
	var dept TicketDepartment
	if err := s.db.First(&dept, deptID).Error; err != nil {
		return errors.New("department not found")
	}

	// Find the next department
	var next TicketDepartment
	if err := s.db.Where("sort_order > ?", dept.SortOrder).Order("sort_order ASC").First(&next).Error; err != nil {
		return nil // Already at bottom
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Swap sort orders
		if err := tx.Model(&dept).Update("sort_order", next.SortOrder).Error; err != nil {
			return err
		}
		return tx.Model(&next).Update("sort_order", dept.SortOrder).Error
	})
}
