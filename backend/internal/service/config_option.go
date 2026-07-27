package service

import (
	"errors"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type ConfigOptionService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewConfigOptionService(db *gorm.DB, log *logger.Logger) *ConfigOptionService {
	return &ConfigOptionService{db: db, log: log}
}

type CreateConfigOptionRequest struct {
	Group        string `json:"group" binding:"required,max=64"`
	Name         string `json:"name" binding:"required,max=128"`
	Code         string `json:"code" binding:"required,max=64"`
	Type         string `json:"type" binding:"required,oneof=text textarea number select radio checkbox switch color date json"`
	Value        string `json:"value"`
	DefaultValue string `json:"default_value"`
	Options      string `json:"options"` // JSON array string
	Placeholder  string `json:"placeholder" binding:"max=256"`
	Tip          string `json:"tip" binding:"max=512"`
	Validation   string `json:"validation" binding:"max=256"`
	IsRequired   bool   `json:"is_required"`
	IsPublic     bool   `json:"is_public"`
	IsReadOnly   bool   `json:"is_read_only"`
	SortOrder    int    `json:"sort_order"`
	Status       int16  `json:"status"`
}

type UpdateConfigOptionRequest struct {
	Group        *string `json:"group"`
	Name         *string `json:"name"`
	Code         *string `json:"code"`
	Type         *string `json:"type"`
	Value        *string `json:"value"`
	DefaultValue *string `json:"default_value"`
	Options      *string `json:"options"`
	Placeholder  *string `json:"placeholder"`
	Tip          *string `json:"tip"`
	Validation   *string `json:"validation"`
	IsRequired   *bool   `json:"is_required"`
	IsPublic     *bool   `json:"is_public"`
	IsReadOnly   *bool   `json:"is_read_only"`
	SortOrder    *int    `json:"sort_order"`
	Status       *int16  `json:"status"`
}

func (s *ConfigOptionService) GetList(page, pageSize int, group string, keyword string) ([]model.ConfigOption, int64, error) {
	var items []model.ConfigOption
	var total int64

	query := s.db.Model(&model.ConfigOption{})
	if group != "" {
		query = query.Where("`group` = ?", group)
	}
	if keyword != "" {
		q := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR code LIKE ?", q, q)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).
		Order("sort_order ASC, id ASC").
		Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (s *ConfigOptionService) GetByGroup(group string) ([]model.ConfigOption, error) {
	var items []model.ConfigOption
	if err := s.db.Where("`group` = ? AND status = 1", group).
		Order("sort_order ASC, id ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (s *ConfigOptionService) GetByID(id uint) (*model.ConfigOption, error) {
	var item model.ConfigOption
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *ConfigOptionService) Create(req CreateConfigOptionRequest) (*model.ConfigOption, error) {
	var count int64
	s.db.Model(&model.ConfigOption{}).Where("code = ?", req.Code).Count(&count)
	if count > 0 {
		return nil, errors.New("option code already exists")
	}

	item := &model.ConfigOption{
		Group:        req.Group,
		Name:         req.Name,
		Code:         req.Code,
		Type:         req.Type,
		Value:        req.Value,
		DefaultValue: req.DefaultValue,
		Placeholder:  req.Placeholder,
		Tip:          req.Tip,
		Validation:   req.Validation,
		IsRequired:   req.IsRequired,
		IsPublic:     req.IsPublic,
		IsReadOnly:   req.IsReadOnly,
		SortOrder:    req.SortOrder,
		Status:       req.Status,
	}
	if item.Status == 0 {
		item.Status = 1
	}
	if req.Options != "" {
		item.Options = []byte(req.Options)
	}

	if err := s.db.Create(item).Error; err != nil {
		return nil, err
	}
	s.log.Infof("config option created: id=%d code=%s", item.ID, item.Code)
	return item, nil
}

func (s *ConfigOptionService) Update(id uint, req UpdateConfigOptionRequest) (*model.ConfigOption, error) {
	var item model.ConfigOption
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}

	if req.Code != nil {
		var count int64
		s.db.Model(&model.ConfigOption{}).Where("code = ? AND id != ?", *req.Code, id).Count(&count)
		if count > 0 {
			return nil, errors.New("option code already exists")
		}
	}

	updates := map[string]interface{}{}
	if req.Group != nil {
		updates["group"] = *req.Group
	}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Code != nil {
		updates["code"] = *req.Code
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Value != nil {
		updates["value"] = *req.Value
	}
	if req.DefaultValue != nil {
		updates["default_value"] = *req.DefaultValue
	}
	if req.Options != nil {
		updates["options"] = *req.Options
	}
	if req.Placeholder != nil {
		updates["placeholder"] = *req.Placeholder
	}
	if req.Tip != nil {
		updates["tip"] = *req.Tip
	}
	if req.Validation != nil {
		updates["validation"] = *req.Validation
	}
	if req.IsRequired != nil {
		updates["is_required"] = *req.IsRequired
	}
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}
	if req.IsReadOnly != nil {
		updates["is_read_only"] = *req.IsReadOnly
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) > 0 {
		if err := s.db.Model(&item).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	s.log.Infof("config option updated: id=%d", id)
	return &item, nil
}

func (s *ConfigOptionService) Delete(id uint) error {
	var item model.ConfigOption
	if err := s.db.First(&item, id).Error; err != nil {
		return err
	}
	if item.IsReadOnly {
		return errors.New("cannot delete a read-only option")
	}
	if err := s.db.Delete(&model.ConfigOption{}, id).Error; err != nil {
		return err
	}
	s.log.Infof("config option deleted: id=%d", id)
	return nil
}

type SortItem struct {
	ID        uint `json:"id"`
	SortOrder int  `json:"sort_order"`
}

func (s *ConfigOptionService) BatchUpdateSort(items []SortItem) error {
	if len(items) == 0 {
		return errors.New("items is empty")
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			result := tx.Model(&model.ConfigOption{}).Where("id = ?", item.ID).Update("sort_order", item.SortOrder)
			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})
}

func (s *ConfigOptionService) UpdateSort(id uint, sortOrder int) error {
	result := s.db.Model(&model.ConfigOption{}).Where("id = ?", id).Update("sort_order", sortOrder)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("config option not found")
	}
	return nil
}

func (s *ConfigOptionService) UpdateValue(code string, value string) error {
	result := s.db.Model(&model.ConfigOption{}).Where("code = ?", code).Update("value", value)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("config option not found")
	}
	return nil
}

func (s *ConfigOptionService) BatchUpdateValue(items map[string]string) error {
	if len(items) == 0 {
		return errors.New("items is empty")
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		for code, value := range items {
			result := tx.Model(&model.ConfigOption{}).Where("code = ?", code).Update("value", value)
			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})
}

func (s *ConfigOptionService) GetGroups() ([]string, error) {
	var groups []string
	if err := s.db.Model(&model.ConfigOption{}).Distinct().Pluck("`group`", &groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}
