package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/upstream"

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
	Group        string  `json:"group" binding:"required,max=64"`
	Name         string  `json:"name" binding:"required,max=128"`
	Code         string  `json:"code" binding:"required,max=64"`
	Type         string  `json:"type" binding:"required,oneof=text textarea number select radio checkbox switch color date json slider quantity"`
	Value        string  `json:"value"`
	DefaultValue string  `json:"default_value"`
	Options      string  `json:"options"` // JSON array string
	Placeholder  string  `json:"placeholder" binding:"max=256"`
	Tip          string  `json:"tip" binding:"max=512"`
	Validation   string  `json:"validation" binding:"max=256"`
	IsRequired   bool    `json:"is_required"`
	IsPublic     bool    `json:"is_public"`
	IsReadOnly   bool    `json:"is_read_only"`
	SortOrder    int     `json:"sort_order"`
	Status       int16   `json:"status"`
	GroupID      uint    `json:"group_id"`
	Min          float64 `json:"min"`
	Max          float64 `json:"max"`
	Step         float64 `json:"step"`
	UpstreamID   *uint   `json:"upstream_id"`
}

type UpdateConfigOptionRequest struct {
	Group        *string  `json:"group"`
	Name         *string  `json:"name"`
	Code         *string  `json:"code"`
	Type         *string  `json:"type"`
	Value        *string  `json:"value"`
	DefaultValue *string  `json:"default_value"`
	Options      *string  `json:"options"`
	Placeholder  *string  `json:"placeholder"`
	Tip          *string  `json:"tip"`
	Validation   *string  `json:"validation"`
	IsRequired   *bool    `json:"is_required"`
	IsPublic     *bool    `json:"is_public"`
	IsReadOnly   *bool    `json:"is_read_only"`
	SortOrder    *int     `json:"sort_order"`
	Status       *int16   `json:"status"`
	GroupID      *uint    `json:"group_id"`
	Min          *float64 `json:"min"`
	Max          *float64 `json:"max"`
	Step         *float64 `json:"step"`
	UpstreamID   *uint    `json:"upstream_id"`
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
		GroupID:      req.GroupID,
		Min:          req.Min,
		Max:          req.Max,
		Step:         req.Step,
		UpstreamID:   req.UpstreamID,
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
	if req.GroupID != nil {
		updates["group_id"] = *req.GroupID
	}
	if req.Min != nil {
		updates["min"] = *req.Min
	}
	if req.Max != nil {
		updates["max"] = *req.Max
	}
	if req.Step != nil {
		updates["step"] = *req.Step
	}
	if req.UpstreamID != nil {
		updates["upstream_id"] = *req.UpstreamID
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

// ─── Product Config Groups ───

type CreateProductConfigGroupRequest struct {
	Name        string `json:"name" binding:"required,max=255"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
	Enabled     *bool  `json:"enabled"`
}

type UpdateProductConfigGroupRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	SortOrder   *int    `json:"sort_order"`
	Enabled     *bool   `json:"enabled"`
}

// GetProductConfigGroups returns all product config groups.
func (s *ConfigOptionService) GetProductConfigGroups() ([]model.ProductConfigGroup, error) {
	var groups []model.ProductConfigGroup
	if err := s.db.Order("sort_order ASC, id ASC").Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

// CreateProductConfigGroup creates a new product config group.
func (s *ConfigOptionService) CreateProductConfigGroup(req CreateProductConfigGroupRequest) (*model.ProductConfigGroup, error) {
	group := &model.ProductConfigGroup{
		Name:        req.Name,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		Enabled:     true,
	}
	if req.Enabled != nil {
		group.Enabled = *req.Enabled
	}
	if err := s.db.Create(group).Error; err != nil {
		return nil, err
	}
	s.log.Infof("product config group created: id=%d name=%s", group.ID, group.Name)
	return group, nil
}

// UpdateProductConfigGroup updates a product config group.
func (s *ConfigOptionService) UpdateProductConfigGroup(id uint, req UpdateProductConfigGroupRequest) (*model.ProductConfigGroup, error) {
	var group model.ProductConfigGroup
	if err := s.db.First(&group, id).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}

	if len(updates) > 0 {
		if err := s.db.Model(&group).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	if err := s.db.First(&group, id).Error; err != nil {
		return nil, err
	}
	s.log.Infof("product config group updated: id=%d", id)
	return &group, nil
}

// DeleteProductConfigGroup deletes a product config group.
func (s *ConfigOptionService) DeleteProductConfigGroup(id uint) error {
	// Check if linked to any products
	var linkCount int64
	s.db.Model(&model.ProductConfigOptionLink{}).Where("group_id = ?", id).Count(&linkCount)
	if linkCount > 0 {
		return errors.New("cannot delete group: still linked to products")
	}

	// Clear group_id from config options
	s.db.Model(&model.ConfigOption{}).Where("group_id = ?", id).Update("group_id", 0)

	if err := s.db.Delete(&model.ProductConfigGroup{}, id).Error; err != nil {
		return err
	}
	s.log.Infof("product config group deleted: id=%d", id)
	return nil
}

// LinkGroupToProduct links a config group to a product.
func (s *ConfigOptionService) LinkGroupToProduct(groupID, productID uint) error {
	var existing model.ProductConfigOptionLink
	err := s.db.Where("group_id = ? AND product_id = ?", groupID, productID).First(&existing).Error
	if err == nil {
		return errors.New("group already linked to this product")
	}

	link := &model.ProductConfigOptionLink{
		ProductID: productID,
		GroupID:   groupID,
	}
	if err := s.db.Create(link).Error; err != nil {
		return err
	}
	s.log.Infof("config group %d linked to product %d", groupID, productID)
	return nil
}

// UnlinkGroupFromProduct removes the link between a config group and a product.
func (s *ConfigOptionService) UnlinkGroupFromProduct(groupID, productID uint) error {
	result := s.db.Where("group_id = ? AND product_id = ?", groupID, productID).Delete(&model.ProductConfigOptionLink{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("link not found")
	}
	s.log.Infof("config group %d unlinked from product %d", groupID, productID)
	return nil
}

// GetProductConfigOptions returns all config options for a product via its linked groups.
func (s *ConfigOptionService) GetProductConfigOptions(productID uint) ([]model.ProductConfigGroup, error) {
	var links []model.ProductConfigOptionLink
	if err := s.db.Where("product_id = ?", productID).Preload("Group").Order("sort_order ASC").Find(&links).Error; err != nil {
		return nil, err
	}

	var groups []model.ProductConfigGroup
	for _, link := range links {
		if link.Group != nil {
			group := *link.Group
			var options []model.ConfigOption
			s.db.Where("group_id = ? AND status = 1", group.ID).Order("sort_order ASC").Find(&options)
			groups = append(groups, group)
			_ = options // options could be attached as needed
		}
	}
	return groups, nil
}

// GetProductConfigOptionsByProduct returns config groups with their options for a product.
func (s *ConfigOptionService) GetProductConfigOptionsByProduct(productID uint) ([]map[string]interface{}, error) {
	var links []model.ProductConfigOptionLink
	if err := s.db.Where("product_id = ?", productID).Preload("Group").Order("sort_order ASC").Find(&links).Error; err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, link := range links {
		if link.Group == nil || !link.Group.Enabled {
			continue
		}
		var options []model.ConfigOption
		s.db.Where("group_id = ? AND status = 1", link.Group.ID).Order("sort_order ASC").Find(&options)

		item := map[string]interface{}{
			"group":   link.Group,
			"options": options,
		}
		result = append(result, item)
	}
	return result, nil
}

// SyncUpstreamOptions syncs config options from the upstream provider for a given product.
// It fetches the remote product's config options, then creates/updates local ConfigOption
// records to match. Returns the number of options synced.
func (s *ConfigOptionService) SyncUpstreamOptions(productID uint) (int, error) {
	// 1. Find the upstream product mapping for this local product.
	var mapping model.UpstreamProduct
	if err := s.db.Where("local_product_id = ?", productID).First(&mapping).Error; err != nil {
		return 0, fmt.Errorf("no upstream mapping for product %d: %w", productID, err)
	}

	// 2. Get the upstream provider.
	var provider model.UpstreamProvider
	if err := s.db.First(&provider, mapping.UpstreamID).Error; err != nil {
		return 0, fmt.Errorf("upstream provider not found: %w", err)
	}
	if !provider.IsActive {
		return 0, fmt.Errorf("upstream provider %d is inactive", provider.ID)
	}

	// 3. Create upstream client and fetch products (config options come from the product data).
	client, err := upstream.NewClient(&provider)
	if err != nil {
		return 0, fmt.Errorf("create upstream client: %w", err)
	}

	remoteProducts, err := client.FetchProducts()
	if err != nil {
		return 0, fmt.Errorf("fetch upstream products: %w", err)
	}

	// 4. Find the matching remote product by RemoteProductID.
	var remoteProduct *upstream.RemoteProduct
	for _, rp := range remoteProducts {
		if rp.RemoteID == mapping.RemoteProductID {
			remoteProduct = &rp
			break
		}
	}
	if remoteProduct == nil {
		return 0, fmt.Errorf("remote product %s not found in upstream", mapping.RemoteProductID)
	}

	// 5. Sync config options from the remote product.
	synced := 0
	if remoteProduct.ConfigOptions == nil || len(remoteProduct.ConfigOptions) == 0 {
		s.log.Infof("product %d: upstream product has no config options", productID)
		return 0, nil
	}

	// Load existing config options that have an upstream_id for this product's groups.
	var existingOptions []model.ConfigOption
	if err := s.db.Where("upstream_id IS NOT NULL AND upstream_id > 0").Find(&existingOptions).Error; err != nil {
		return 0, fmt.Errorf("load existing options: %w", err)
	}

	// Index existing options by upstream_id for fast lookup.
	existingByUpstreamID := make(map[uint]*model.ConfigOption)
	for i := range existingOptions {
		if existingOptions[i].UpstreamID != nil {
			existingByUpstreamID[*existingOptions[i].UpstreamID] = &existingOptions[i]
		}
	}

	// Process each remote config option.
	seenUpstreamIDs := make(map[uint]bool)
	for key, val := range remoteProduct.ConfigOptions {
		optMap, ok := val.(map[string]interface{})
		if !ok {
			continue
		}

		remoteID, _ := optMap["id"].(float64)
		if remoteID == 0 {
			continue
		}
		upstreamID := uint(remoteID)
		seenUpstreamIDs[upstreamID] = true

		name, _ := optMap["name"].(string)
		if name == "" {
			name = key
		}
		code, _ := optMap["code"].(string)
		if code == "" {
			code = fmt.Sprintf("upstream_%d", upstreamID)
		}
		optType, _ := optMap["type"].(string)
		if optType == "" {
			optType = "text"
		}
		value, _ := optMap["value"].(string)
		defaultValue, _ := optMap["default_value"].(string)
		placeholder, _ := optMap["placeholder"].(string)
		tip, _ := optMap["tip"].(string)

		// Serialize options array if present.
		var optionsJSON []byte
		if optArr, ok := optMap["options"]; ok {
			optionsJSON, _ = json.Marshal(optArr)
		}

		if existing, found := existingByUpstreamID[upstreamID]; found {
			// Update existing option.
			updates := map[string]interface{}{
				"name":          name,
				"code":          code,
				"type":          optType,
				"value":         value,
				"default_value": defaultValue,
				"placeholder":   placeholder,
				"tip":           tip,
			}
			if optionsJSON != nil {
				updates["options"] = string(optionsJSON)
			}
			if err := s.db.Model(existing).Updates(updates).Error; err != nil {
				s.log.Warnf("update config option upstream_id=%d: %v", upstreamID, err)
				continue
			}
		} else {
			// Create new option.
			newOpt := model.ConfigOption{
				Name:         name,
				Code:         code,
				Type:         optType,
				Value:        value,
				DefaultValue: defaultValue,
				Placeholder:  placeholder,
				Tip:          tip,
				Status:       1,
				UpstreamID:   &upstreamID,
			}
			if optionsJSON != nil {
				newOpt.Options = optionsJSON
			}
			if err := s.db.Create(&newOpt).Error; err != nil {
				s.log.Warnf("create config option code=%s: %v", code, err)
				continue
			}
		}
		synced++
	}

	// 6. Delete local options with upstream_id that no longer exist upstream.
	for upstreamID, existing := range existingByUpstreamID {
		if !seenUpstreamIDs[upstreamID] {
			if existing.IsReadOnly {
				continue
			}
			if err := s.db.Delete(existing).Error; err != nil {
				s.log.Warnf("delete stale config option id=%d: %v", existing.ID, err)
			}
		}
	}

	s.log.Infof("sync upstream options for product %d: %d options synced", productID, synced)
	return synced, nil
}
