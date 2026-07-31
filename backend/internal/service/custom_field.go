package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type CustomFieldService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewCustomFieldService(db *gorm.DB, log *logger.Logger) *CustomFieldService {
	return &CustomFieldService{db: db, log: log}
}

// ────────────────────────── 请求体 ──────────────────────────

type CreateFieldRequest struct {
	Name        string   `json:"name" binding:"required,max=128"`
	Label       string   `json:"label" binding:"required,max=256"`
	Type        string   `json:"type" binding:"required,oneof=text textarea select checkbox radio file date number"`
	Group       string   `json:"group" binding:"required,oneof=product cart client host"`
	Required    bool     `json:"required"`
	DefaultVal  string   `json:"default_val"`
	Options     []string `json:"options"`
	Validation  string   `json:"validation"`
	Placeholder string   `json:"placeholder"`
	HelpText    string   `json:"help_text"`
	SortOrder   int      `json:"sort_order"`
}

type UpdateFieldRequest struct {
	Name        *string  `json:"name"`
	Label       *string  `json:"label"`
	Type        *string  `json:"type"`
	Group       *string  `json:"group"`
	Required    *bool    `json:"required"`
	DefaultVal  *string  `json:"default_val"`
	Options     []string `json:"options"`
	Validation  *string  `json:"validation"`
	Placeholder *string  `json:"placeholder"`
	HelpText    *string  `json:"help_text"`
	SortOrder   *int     `json:"sort_order"`
	Enabled     *bool    `json:"enabled"`
}

type SaveValuesRequest struct {
	Values map[string]string `json:"values" binding:"required"`
}

type CreateGroupRequest struct {
	Name      string `json:"name" binding:"required,max=128"`
	Label     string `json:"label" binding:"required,max=256"`
	Type      string `json:"type" binding:"required,oneof=product cart client host"`
	SortOrder int    `json:"sort_order"`
}

type UpdateGroupRequest struct {
	Name      *string `json:"name"`
	Label     *string `json:"label"`
	Type      *string `json:"type"`
	SortOrder *int    `json:"sort_order"`
	Enabled   *bool   `json:"enabled"`
}

type ReorderRequest struct {
	IDs []uint `json:"ids" binding:"required"`
}

type CopyFieldsRequest struct {
	FromGroup string `json:"from_group" binding:"required"`
	ToGroup   string `json:"to_group" binding:"required"`
}

type ImportFieldsRequest struct {
	Fields []CreateFieldRequest `json:"fields" binding:"required"`
}

type ValidateRequest struct {
	Group  string            `json:"group" binding:"required"`
	Values map[string]string `json:"values" binding:"required"`
}

// ────────────────────────── Field Management ──────────────────────────

// GetFields lists fields by group.
func (s *CustomFieldService) GetFields(group string) ([]model.CustomField, error) {
	var fields []model.CustomField
	q := s.db.Order("sort_order ASC, id ASC")
	if group != "" {
		q = q.Where("group = ?", group)
	}
	if err := q.Find(&fields).Error; err != nil {
		return nil, err
	}
	return fields, nil
}

// GetFieldByID returns a single field.
func (s *CustomFieldService) GetFieldByID(id uint) (*model.CustomField, error) {
	var field model.CustomField
	if err := s.db.First(&field, id).Error; err != nil {
		return nil, err
	}
	return &field, nil
}

// CreateField creates a new custom field.
func (s *CustomFieldService) CreateField(req CreateFieldRequest) (*model.CustomField, error) {
	field := model.CustomField{
		Name:        req.Name,
		Label:       req.Label,
		Type:        req.Type,
		Group:       req.Group,
		Required:    req.Required,
		DefaultVal:  req.DefaultVal,
		Validation:  req.Validation,
		Placeholder: req.Placeholder,
		HelpText:    req.HelpText,
		SortOrder:   req.SortOrder,
		Enabled:     true,
	}
	if len(req.Options) > 0 {
		optJSON, err := json.Marshal(req.Options)
		if err != nil {
			return nil, errors.New("invalid options format")
		}
		field.Options = string(optJSON)
	}
	if err := s.db.Create(&field).Error; err != nil {
		return nil, err
	}
	return &field, nil
}

// UpdateField updates an existing field.
func (s *CustomFieldService) UpdateField(id uint, req UpdateFieldRequest) (*model.CustomField, error) {
	var field model.CustomField
	if err := s.db.First(&field, id).Error; err != nil {
		return nil, errors.New("field not found")
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Label != nil {
		updates["label"] = *req.Label
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Group != nil {
		updates["group"] = *req.Group
	}
	if req.Required != nil {
		updates["required"] = *req.Required
	}
	if req.DefaultVal != nil {
		updates["default_val"] = *req.DefaultVal
	}
	if req.Validation != nil {
		updates["validation"] = *req.Validation
	}
	if req.Placeholder != nil {
		updates["placeholder"] = *req.Placeholder
	}
	if req.HelpText != nil {
		updates["help_text"] = *req.HelpText
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}
	if req.Options != nil {
		optJSON, err := json.Marshal(req.Options)
		if err != nil {
			return nil, errors.New("invalid options format")
		}
		updates["options"] = string(optJSON)
	}

	if err := s.db.Model(&field).Updates(updates).Error; err != nil {
		return nil, err
	}
	if err := s.db.First(&field, id).Error; err != nil {
		return nil, err
	}
	return &field, nil
}

// DeleteField deletes a field and its values.
func (s *CustomFieldService) DeleteField(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("field_id = ?", id).Delete(&model.CustomFieldValue{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.CustomField{}, id).Error
	})
}

// ReorderFields reorders fields by the given ID list.
func (s *CustomFieldService) ReorderFields(ids []uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for i, id := range ids {
			if err := tx.Model(&model.CustomField{}).Where("id = ?", id).Update("sort_order", i).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ────────────────────────── Field Values ──────────────────────────

// GetValues returns all values for an owner.
func (s *CustomFieldService) GetValues(ownerID uint, ownerType string) ([]model.CustomFieldValue, error) {
	var values []model.CustomFieldValue
	if err := s.db.Where("owner_id = ? AND owner_type = ?", ownerID, ownerType).Find(&values).Error; err != nil {
		return nil, err
	}
	return values, nil
}

// SaveValues saves field values for an owner (upsert).
func (s *CustomFieldService) SaveValues(ownerID uint, ownerType string, values map[string]string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for fieldName, val := range values {
			var field model.CustomField
			if err := tx.Where("name = ? AND enabled = true", fieldName).First(&field).Error; err != nil {
				continue
			}
			var existing model.CustomFieldValue
			err := tx.Where("field_id = ? AND owner_id = ? AND owner_type = ?", field.ID, ownerID, ownerType).First(&existing).Error
			if err == nil {
				existing.Value = val
				if err := tx.Save(&existing).Error; err != nil {
					return err
				}
			} else {
				newVal := model.CustomFieldValue{
					FieldID:   field.ID,
					OwnerID:   ownerID,
					OwnerType: ownerType,
					Value:     val,
				}
				if err := tx.Create(&newVal).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// DeleteValues deletes all values for an owner.
func (s *CustomFieldService) DeleteValues(ownerID uint, ownerType string) error {
	return s.db.Where("owner_id = ? AND owner_type = ?", ownerID, ownerType).Delete(&model.CustomFieldValue{}).Error
}

// GetValue returns a single field value.
func (s *CustomFieldService) GetValue(fieldID, ownerID uint, ownerType string) (*model.CustomFieldValue, error) {
	var val model.CustomFieldValue
	if err := s.db.Where("field_id = ? AND owner_id = ? AND owner_type = ?", fieldID, ownerID, ownerType).First(&val).Error; err != nil {
		return nil, err
	}
	return &val, nil
}

// ────────────────────────── Group Management ──────────────────────────

// GetGroups lists groups by type.
func (s *CustomFieldService) GetGroups(typ string) ([]model.CustomFieldGroup, error) {
	var groups []model.CustomFieldGroup
	q := s.db.Order("sort_order ASC, id ASC")
	if typ != "" {
		q = q.Where("type = ?", typ)
	}
	if err := q.Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

// CreateGroup creates a new field group.
func (s *CustomFieldService) CreateGroup(req CreateGroupRequest) (*model.CustomFieldGroup, error) {
	group := model.CustomFieldGroup{
		Name:      req.Name,
		Label:     req.Label,
		Type:      req.Type,
		SortOrder: req.SortOrder,
		Enabled:   true,
	}
	if err := s.db.Create(&group).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

// UpdateGroup updates an existing group.
func (s *CustomFieldService) UpdateGroup(id uint, req UpdateGroupRequest) (*model.CustomFieldGroup, error) {
	var group model.CustomFieldGroup
	if err := s.db.First(&group, id).Error; err != nil {
		return nil, errors.New("group not found")
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Label != nil {
		updates["label"] = *req.Label
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}

	if err := s.db.Model(&group).Updates(updates).Error; err != nil {
		return nil, err
	}
	if err := s.db.First(&group, id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

// DeleteGroup deletes a group.
func (s *CustomFieldService) DeleteGroup(id uint) error {
	return s.db.Delete(&model.CustomFieldGroup{}, id).Error
}

// ────────────────────────── Validation ──────────────────────────

// ValidateFields validates field values against field definitions.
func (s *CustomFieldService) ValidateFields(group string, values map[string]string) ([]string, error) {
	var fields []model.CustomField
	if err := s.db.Where("`group` = ? AND enabled = true", group).Find(&fields).Error; err != nil {
		return nil, err
	}

	var errs []string
	for _, field := range fields {
		val, exists := values[field.Name]

		if field.Required && (!exists || val == "") {
			errs = append(errs, fmt.Sprintf("%s is required", field.Label))
			continue
		}

		if !exists || val == "" {
			continue
		}

		if field.Validation != "" {
			re, err := regexp.Compile(field.Validation)
			if err != nil {
				s.log.Warnf("invalid validation regex for field %s: %v", field.Name, err)
				continue
			}
			if !re.MatchString(val) {
				errs = append(errs, fmt.Sprintf("%s format is invalid", field.Label))
			}
		}
	}
	return errs, nil
}

// GetCartCustomFields returns cart-specific fields.
func (s *CustomFieldService) GetCartCustomFields() ([]model.CustomField, error) {
	return s.GetFields("cart")
}

// GetProductCustomFields returns product fields with values.
func (s *CustomFieldService) GetProductCustomFields(productID uint) ([]model.CustomField, map[string]string, error) {
	fields, err := s.GetFields("product")
	if err != nil {
		return nil, nil, err
	}
	valMap, err := s.getValueMap(productID, "product")
	if err != nil {
		return nil, nil, err
	}
	return fields, valMap, nil
}

// GetClientCustomFields returns client fields.
func (s *CustomFieldService) GetClientCustomFields() ([]model.CustomField, error) {
	return s.GetFields("client")
}

// GetHostCustomFields returns host fields with values.
func (s *CustomFieldService) GetHostCustomFields(hostID uint) ([]model.CustomField, map[string]string, error) {
	fields, err := s.GetFields("host")
	if err != nil {
		return nil, nil, err
	}
	valMap, err := s.getValueMap(hostID, "host")
	if err != nil {
		return nil, nil, err
	}
	return fields, valMap, nil
}

// ────────────────────────── Bulk Operations ──────────────────────────

// CopyFields copies fields from one group to another.
func (s *CustomFieldService) CopyFields(fromGroup, toGroup string) error {
	var fields []model.CustomField
	if err := s.db.Where("`group` = ?", fromGroup).Order("sort_order ASC").Find(&fields).Error; err != nil {
		return err
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, f := range fields {
			newField := f
			newField.ID = 0
			newField.Group = toGroup
			if err := tx.Create(&newField).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ImportFields imports fields from a request.
func (s *CustomFieldService) ImportFields(req ImportFieldsRequest) (int, error) {
	count := 0
	err := s.db.Transaction(func(tx *gorm.DB) error {
		for _, r := range req.Fields {
			field := model.CustomField{
				Name:        r.Name,
				Label:       r.Label,
				Type:        r.Type,
				Group:       r.Group,
				Required:    r.Required,
				DefaultVal:  r.DefaultVal,
				Validation:  r.Validation,
				Placeholder: r.Placeholder,
				HelpText:    r.HelpText,
				SortOrder:   r.SortOrder,
				Enabled:     true,
			}
			if len(r.Options) > 0 {
				optJSON, err := json.Marshal(r.Options)
				if err != nil {
					return err
				}
				field.Options = string(optJSON)
			}
			if err := tx.Create(&field).Error; err != nil {
				return err
			}
			count++
		}
		return nil
	})
	return count, err
}

// ExportFields exports fields of a group as JSON.
func (s *CustomFieldService) ExportFields(group string) ([]map[string]interface{}, error) {
	fields, err := s.GetFields(group)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, f := range fields {
		item := map[string]interface{}{
			"name":         f.Name,
			"label":        f.Label,
			"type":         f.Type,
			"group":        f.Group,
			"required":     f.Required,
			"default_val":  f.DefaultVal,
			"validation":   f.Validation,
			"placeholder":  f.Placeholder,
			"help_text":    f.HelpText,
			"sort_order":   f.SortOrder,
		}
		if f.Options != "" {
			var opts []string
			if err := json.Unmarshal([]byte(f.Options), &opts); err == nil {
				item["options"] = opts
			}
		}
		result = append(result, item)
	}
	return result, nil
}

// ────────────────────────── helpers ──────────────────────────

func (s *CustomFieldService) getValueMap(ownerID uint, ownerType string) (map[string]string, error) {
	var values []model.CustomFieldValue
	if err := s.db.Where("owner_id = ? AND owner_type = ?", ownerID, ownerType).Find(&values).Error; err != nil {
		return nil, err
	}
	var fieldIDs []uint
	for _, v := range values {
		fieldIDs = append(fieldIDs, v.FieldID)
	}
	var fields []model.CustomField
	if len(fieldIDs) > 0 {
		s.db.Where("id IN ?", fieldIDs).Find(&fields)
	}
	idNameMap := make(map[uint]string)
	for _, f := range fields {
		idNameMap[f.ID] = f.Name
	}
	m := make(map[string]string)
	for _, v := range values {
		if name, ok := idNameMap[v.FieldID]; ok {
			m[name] = v.Value
		}
	}
	return m, nil
}
