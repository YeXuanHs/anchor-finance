package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

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

type ConfigSortItem struct {
	ID        uint `json:"id"`
	SortOrder int  `json:"sort_order"`
}

func (s *ConfigOptionService) BatchUpdateSort(items []ConfigSortItem) error {
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

// ─── ConfigOptions Admin Methods (from zjmf ConfigOptionsController) ───

// AdminGroupsList returns all global product config groups with linked products.
func (s *ConfigOptionService) AdminGroupsList(order, sort, keywords string) ([]map[string]interface{}, error) {
	query := s.db.Table("product_config_groups").
		Select("id, name, description").
		Where("`global` = 1")
	if keywords != "" {
		query = query.Where("name LIKE ?", "%"+keywords+"%")
	}
	if order == "" {
		order = "id"
	}
	if sort == "" {
		sort = "desc"
	}

	var groups []map[string]interface{}
	if err := query.Order(order + " " + sort).Find(&groups).Error; err != nil {
		return nil, err
	}

	for i, g := range groups {
		gid, _ := g["id"].(uint)
		var products []struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		}
		s.db.Table("product_config_links a").
			Select("b.id, b.name").
			Joins("LEFT JOIN products b ON a.pid = b.id").
			Where("a.gid = ?", gid).
			Find(&products)
		g["products"] = products
		groups[i] = g
	}
	return groups, nil
}

// AdminSearchPage returns available product types for search.
func (s *ConfigOptionService) AdminSearchPage() map[string]string {
	types := map[string]string{
		"hosting":   "虚拟主机",
		"server":    "独立服务器",
		"vps":       "VPS",
		"reseller":  "分销",
		"other":     "其他",
	}
	return types
}

// AdminCreateGroupsData returns products grouped by product groups for the create page.
func (s *ConfigOptionService) AdminCreateGroupsData(typeFilter int) ([]map[string]interface{}, error) {
	query := s.db.Table("products p").
		Select("pg.name as pg_name, p.name as p_name, p.id as p_id, CONCAT(pg.name, '：', p.name) as link, pg.type").
		Joins("LEFT JOIN product_groups pg ON p.gid = pg.id")
	if typeFilter > 0 && typeFilter < 4 {
		query = query.Where("pg.type = ?", typeFilter)
	}

	var result []map[string]interface{}
	if err := query.Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

// AdminCreateGroup creates a new product config group with product links.
func (s *ConfigOptionService) AdminCreateGroup(name, description string, global int, productIDs []uint) (uint, error) {
	var groupID uint
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Insert group
		result := tx.Exec(
			"INSERT INTO product_config_groups (name, description, `global`) VALUES (?, ?, ?)",
			name, description, global,
		)
		if result.Error != nil {
			return result.Error
		}
		groupID = uint(result.RowsAffected)
		// Get the actual inserted ID
		var lastID struct{ ID uint }
		tx.Raw("SELECT LAST_INSERT_ID() as id").Scan(&lastID)
		groupID = lastID.ID

		if len(productIDs) > 0 {
			// Verify products exist
			var validIDs []uint
			tx.Table("products").Where("id IN ?", productIDs).Pluck("id", &validIDs)

			for _, pid := range validIDs {
				tx.Exec("INSERT INTO product_config_links (gid, pid) VALUES (?, ?)", groupID, pid)
			}
			// Increment product version
			tx.Exec("UPDATE products SET location_version = location_version + 1 WHERE id IN ?", validIDs)
		}
		return nil
	})
	return groupID, err
}

// AdminEditGroupData returns data for editing a config group.
func (s *ConfigOptionService) AdminEditGroupData(gid uint, typeFilter int) (map[string]interface{}, error) {
	var group map[string]interface{}
	if err := s.db.Table("product_config_groups").Where("id = ?", gid).Find(&group).Error; err != nil {
		return nil, err
	}

	var pids []struct{ PID uint }
	s.db.Table("product_config_links").Where("gid = ?", gid).Find(&pids)

	query := s.db.Table("products p").
		Select("pg.name as pg_name, p.name as p_name, p.id as p_id, CONCAT(pg.name, '：', p.name) as link, pg.type").
		Joins("LEFT JOIN product_groups pg ON p.gid = pg.id")
	if typeFilter > 0 && typeFilter < 4 {
		query = query.Where("pg.type = ?", typeFilter)
	}
	var products []map[string]interface{}
	query.Find(&products)

	var options []map[string]interface{}
	s.db.Table("product_config_options").
		Where("gid = ? AND linkage_pid = 0 AND linkage_top_pid = 0", gid).
		Find(&options)

	return map[string]interface{}{
		"group":   group,
		"pid":     pids,
		"product": products,
		"options": options,
	}, nil
}

// AdminEditGroupPost updates a config group.
func (s *ConfigOptionService) AdminEditGroupPost(gid uint, name, description string, productIDs []uint, orders map[string]int, hidden map[string]int, upgrade map[string]int) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		tx.Exec("UPDATE product_config_groups SET name = ?, description = ? WHERE id = ?",
			name, description, gid)

		// Update product links
		tx.Exec("DELETE FROM product_config_links WHERE gid = ?", gid)
		for _, pid := range productIDs {
			tx.Exec("INSERT INTO product_config_links (gid, pid) VALUES (?, ?)", gid, pid)
		}
		if len(productIDs) > 0 {
			tx.Exec("UPDATE products SET location_version = location_version + 1 WHERE id IN ?", productIDs)
		}

		// Update option orders/hidden/upgrade
		for idStr, sortOrder := range orders {
			optID, _ := strconv.ParseUint(idStr, 10, 64)
			if optID == 0 {
				continue
			}
			h := 0
			if v, ok := hidden[idStr]; ok && v != 0 {
				h = 1
			}
			u := 0
			if v, ok := upgrade[idStr]; ok && v != 0 {
				u = 1
			}
			tx.Exec("UPDATE product_config_options SET `order` = ?, hidden = ?, upgrade = ? WHERE id = ?",
				sortOrder, h, u, optID)
		}
		return nil
	})
}

// AdminAddOptionsPage returns data for the add options page.
func (s *ConfigOptionService) AdminAddOptionsPage(gid uint, pid uint) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"type": map[string]string{
			"1":  "下拉",
			"2":  "单选",
			"3":  "是/否",
			"4":  "复选",
			"5":  "多选",
			"6":  "文本",
			"7":  "文本区域",
			"8":  "密码",
			"9":  "链接",
			"10": "数量",
			"11": "日期",
			"12": "颜色",
		},
		"paytype":              "",
		"pay_type_recurring":   []string{},
		"pay_ontrial_status":   0,
		"cycle":                []string{},
	}
	return data, nil
}

// AdminAddOption creates a new config option with a sub-option.
func (s *ConfigOptionService) AdminAddOption(gid uint, optionName string, optionType int, notes string, qtyStage int, linkagePID uint, subOptionName string, subSortOrder int, subHidden int) (uint, error) {
	var cid uint
	err := s.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Exec(
			"INSERT INTO product_config_options (gid, option_name, option_type, notes, qty_stage, linkage_pid) VALUES (?, ?, ?, ?, ?, ?)",
			gid, optionName, optionType, notes, qtyStage, linkagePID,
		)
		if result.Error != nil {
			return result.Error
		}
		var lastID struct{ ID uint }
		tx.Raw("SELECT LAST_INSERT_ID() as id").Scan(&lastID)
		cid = lastID.ID

		// Update linkage level
		if linkagePID == 0 {
			tx.Exec("UPDATE product_config_options SET linkage_level = ?, linkage_top_pid = 0 WHERE id = ?",
				"0-"+strconv.FormatUint(uint64(cid), 10), cid)
		} else {
			var parent struct{ LinkageLevel string; LinkageTopPID uint; LinkagePID uint }
			tx.Raw("SELECT linkage_level, linkage_top_pid, linkage_pid FROM product_config_options WHERE id = ?", linkagePID).Scan(&parent)
			topPID := parent.LinkageTopPID
			if parent.LinkagePID == 0 {
				topPID = linkagePID
			}
			tx.Exec("UPDATE product_config_options SET linkage_top_pid = ?, linkage_level = ? WHERE id = ?",
				topPID, parent.LinkageLevel+"-"+strconv.FormatUint(uint64(cid), 10), cid)
		}

		// Create sub-option
		var subID uint
		subResult := tx.Exec(
			"INSERT INTO product_config_options_sub (config_id, option_name, sort_order, hidden, qty_minimum, qty_maximum) VALUES (?, ?, ?, ?, 0, 0)",
			cid, subOptionName, subSortOrder, subHidden,
		)
		if subResult.Error != nil {
			return subResult.Error
		}
		var subLastID struct{ ID uint }
		tx.Raw("SELECT LAST_INSERT_ID() as id").Scan(&subLastID)
		subID = subLastID.ID

		// Create pricing entries for all currencies
		var currencyIDs []uint
		tx.Table("currencies").Pluck("id", &currencyIDs)
		for _, currencyID := range currencyIDs {
			tx.Exec("INSERT INTO pricing (type, relid, currency, monthly) VALUES ('configoptions', ?, ?, 0)", subID, currencyID)
		}

		// Update product versions
		tx.Exec("UPDATE products SET location_version = location_version + 1 WHERE id IN (SELECT pid FROM product_config_links WHERE gid = ?)", gid)
		return nil
	})
	return cid, err
}

// AdminDeleteSubOption deletes a config sub-option and its pricing.
func (s *ConfigOptionService) AdminDeleteSubOption(subID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var sub struct{ UpstreamID uint }
		tx.Raw("SELECT upstream_id FROM product_config_options_sub WHERE id = ?", subID).Scan(&sub)
		if sub.UpstreamID > 0 {
			return fmt.Errorf("上游配置子项,不可删除")
		}

		// Delete sub-option and its linkage children
		tx.Exec("DELETE FROM pricing WHERE type = 'configoptions' AND relid = ?", subID)
		tx.Exec("DELETE FROM product_config_options_sub WHERE id = ?", subID)

		// Update product versions
		tx.Exec(`UPDATE products SET location_version = location_version + 1 WHERE id IN (
			SELECT pid FROM product_config_links WHERE gid IN (
				SELECT gid FROM product_config_options WHERE id IN (
					SELECT config_id FROM product_config_options_sub WHERE id = ?
				)
			)
		)`, subID)
		return nil
	})
}

// AdminDeleteOption deletes a config option and all its sub-options.
func (s *ConfigOptionService) AdminDeleteOption(cid uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var opt struct{ UpstreamID uint }
		tx.Raw("SELECT upstream_id FROM product_config_options WHERE id = ?", cid).Scan(&opt)
		if opt.UpstreamID > 0 {
			return fmt.Errorf("上游配置项,不可删除")
		}

		// Get sub-option IDs for pricing cleanup
		var subIDs []uint
		tx.Table("product_config_options_sub").Where("config_id = ?", cid).Pluck("id", &subIDs)
		if len(subIDs) > 0 {
			tx.Exec("DELETE FROM pricing WHERE type = 'configoptions' AND relid IN ?", subIDs)
		}
		tx.Exec("DELETE FROM product_config_options_sub WHERE config_id = ?", cid)
		tx.Exec("DELETE FROM product_config_options WHERE id = ?", cid)

		// Update product versions
		tx.Exec(`UPDATE products SET location_version = location_version + 1 WHERE id IN (
			SELECT pid FROM product_config_links WHERE gid IN (
				SELECT gid FROM product_config_options WHERE id = ?
			)
		)`, cid)
		return nil
	})
}

// AdminDeleteGroup deletes a config group and all its options/sub-options.
func (s *ConfigOptionService) AdminDeleteGroup(gid uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Get all option IDs
		var cids []uint
		tx.Table("product_config_options").Where("gid = ?", gid).Pluck("id", &cids)

		if len(cids) > 0 {
			var subIDs []uint
			tx.Table("product_config_options_sub").Where("config_id IN ?", cids).Pluck("id", &subIDs)
			if len(subIDs) > 0 {
				tx.Exec("DELETE FROM pricing WHERE type = 'configoptions' AND relid IN ?", subIDs)
			}
			tx.Exec("DELETE FROM product_config_options_sub WHERE config_id IN ?", cids)
			tx.Exec("DELETE FROM product_config_options WHERE gid = ?", gid)
		}
		tx.Exec("DELETE FROM product_config_links WHERE gid = ?", gid)
		tx.Exec("DELETE FROM product_config_groups WHERE id = ?", gid)

		// Update product versions
		tx.Exec("UPDATE products SET location_version = location_version + 1 WHERE id IN (SELECT pid FROM product_config_links WHERE gid = ?)", gid)
		return nil
	})
}

// AdminDuplicateGroups returns all global groups for duplication.
func (s *ConfigOptionService) AdminDuplicateGroups() ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	if err := s.db.Table("product_config_groups").
		Select("id, name, description, CONCAT(name, '--', description) as links").
		Where("`global` = 1").
		Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

// AdminDuplicateGroupPost duplicates a config group.
func (s *ConfigOptionService) AdminDuplicateGroupPost(gid uint, newName string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var oldGroup struct{ Description string; Global int }
		tx.Raw("SELECT description, `global` FROM product_config_groups WHERE id = ?", gid).Scan(&oldGroup)

		var newGroupID uint
		result := tx.Exec("INSERT INTO product_config_groups (name, description, `global`) VALUES (?, ?, ?)",
			newName, oldGroup.Description, oldGroup.Global)
		if result.Error != nil {
			return result.Error
		}
		var lastID struct{ ID uint }
		tx.Raw("SELECT LAST_INSERT_ID() as id").Scan(&lastID)
		newGroupID = lastID.ID

		// Duplicate product links
		var oldLinks []struct{ PID uint }
		tx.Raw("SELECT pid FROM product_config_links WHERE gid = ?", gid).Scan(&oldLinks)
		for _, l := range oldLinks {
			tx.Exec("INSERT INTO product_config_links (gid, pid) VALUES (?, ?)", newGroupID, l.PID)
		}

		// Duplicate options
		var oldOptions []struct {
			ID          uint
			OptionName  string
			OptionType  int
			QtyMinimum  int
			QtyMaximum  int
			Order       int
			Hidden      int
			Upgrade     int
		}
		tx.Raw("SELECT id, option_name, option_type, qty_minimum, qty_maximum, `order`, hidden, upgrade FROM product_config_options WHERE gid = ?", gid).Scan(&oldOptions)

		for _, oldOpt := range oldOptions {
			_ = tx.Exec(`INSERT INTO product_config_options (gid, option_name, option_type, qty_minimum, qty_maximum, `+"`order`"+`, hidden, upgrade) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
				newGroupID, oldOpt.OptionName, oldOpt.OptionType, oldOpt.QtyMinimum, oldOpt.QtyMaximum, oldOpt.Order, oldOpt.Hidden, oldOpt.Upgrade)
			var newCID struct{ ID uint }
			tx.Raw("SELECT LAST_INSERT_ID() as id").Scan(&newCID)

			// Duplicate sub-options
			var oldSubs []struct {
				ID         uint
				OptionName string
				SortOrder  int
				Hidden     int
				QtyMin     int
				QtyMax     int
			}
			tx.Raw("SELECT id, option_name, sort_order, hidden, qty_minimum, qty_maximum FROM product_config_options_sub WHERE config_id = ?", oldOpt.ID).Scan(&oldSubs)

			for _, oldSub := range oldSubs {
				_ = tx.Exec(`INSERT INTO product_config_options_sub (config_id, option_name, sort_order, hidden, qty_minimum, qty_maximum) VALUES (?, ?, ?, ?, ?, ?)`,
					newCID.ID, oldSub.OptionName, oldSub.SortOrder, oldSub.Hidden, oldSub.QtyMin, oldSub.QtyMax)
				var newSubID struct{ ID uint }
				tx.Raw("SELECT LAST_INSERT_ID() as id").Scan(&newSubID)

				// Duplicate pricing
				var oldPricings []map[string]interface{}
				tx.Raw("SELECT * FROM pricing WHERE type = 'configoptions' AND relid = ?", oldSub.ID).Scan(&oldPricings)
				for _, p := range oldPricings {
					delete(p, "id")
					p["relid"] = newSubID.ID
					cols := make([]string, 0, len(p))
					vals := make([]interface{}, 0, len(p))
					placeholders := make([]string, 0, len(p))
					for k, v := range p {
						cols = append(cols, k)
						vals = append(vals, v)
						placeholders = append(placeholders, "?")
					}
					tx.Exec("INSERT INTO pricing ("+strings.Join(cols, ",")+") VALUES ("+strings.Join(placeholders, ",")+")", vals...)
				}
			}
		}
		return nil
	})
}

// AdminEditConfig returns config option detail with sub-options and pricing.
func (s *ConfigOptionService) AdminEditConfig(cid uint, pid uint) (map[string]interface{}, error) {
	var option map[string]interface{}
	if err := s.db.Table("product_config_options").Where("id = ?", cid).Find(&option).Error; err != nil {
		return nil, err
	}

	var subOptions []map[string]interface{}
	if err := s.db.Table("product_config_options_sub").
		Where("config_id = ?", cid).
		Order("sort_order ASC, id ASC").
		Find(&subOptions).Error; err != nil {
		return nil, err
	}

	for i, sub := range subOptions {
		subID, _ := sub["id"].(uint)
		var pricings []map[string]interface{}
		s.db.Table("pricing p").
			Select("c.code, p.currency, p.monthly, p.quarterly, p.semi_annually, p.annually, p.biennially, p.triennially, p.onetime, p.free, p.msetupfee, p.qsetupfee, p.ssetupfee, p.asetupfee, p.bsetupfee, p.tsetupfee").
			Joins("LEFT JOIN currencies c ON c.id = p.currency").
			Where("p.type = 'configoptions' AND p.relid = ?", subID).
			Find(&pricings)
		sub["child"] = pricings
		subOptions[i] = sub
	}

	return map[string]interface{}{
		"option":     option,
		"suboptions": subOptions,
		"type": map[string]string{
			"1": "下拉", "2": "单选", "3": "是/否", "4": "复选", "5": "多选",
			"6": "文本", "7": "文本区域", "8": "密码", "9": "链接", "10": "数量", "11": "日期", "12": "颜色",
		},
	}, nil
}

// AdminGetNextLinkAgeList returns the next level of linkage options.
func (s *ConfigOptionService) AdminGetNextLinkAgeList(subID uint) (map[string]interface{}, error) {
	var configID uint
	s.db.Table("product_config_options_sub").Where("linkage_pid = ?", subID).Pluck("config_id", &configID)

	if configID == 0 {
		return map[string]interface{}{"sub_pid": subID}, nil
	}

	var option map[string]interface{}
	s.db.Table("product_config_options").Where("id = ?", configID).Find(&option)

	var subOptions []map[string]interface{}
	s.db.Table("product_config_options_sub").
		Where("config_id = ?", configID).
		Order("sort_order DESC").
		Find(&subOptions)
	for i := range subOptions {
		subOptions[i]["disable"] = false
	}

	option["sub_son"] = subOptions
	option["sub_pid"] = subID
	return option, nil
}

// AdminEditConfigPost updates config option and sub-options.
func (s *ConfigOptionService) AdminEditConfigPost(cid uint, optionName string, optionType int, notes string, qtyMin, qtyMax, qtyStage, isDiscount, senior int, unit string, subUpdates map[uint]map[string]interface{}, prices map[string]map[string]map[string]interface{}) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		tx.Exec(`UPDATE product_config_options SET option_name = ?, option_type = ?, notes = ?, qty_minimum = ?, qty_maximum = ?, qty_stage = ?, is_discount = ?, unit = ?, senior = ? WHERE id = ?`,
			optionName, optionType, notes, qtyMin, qtyMax, qtyStage, isDiscount, unit, senior, cid)

		// Update existing sub-options
		for subID, subData := range subUpdates {
			name, _ := subData["option_name"].(string)
			sortOrder, _ := subData["sort_order"].(int)
			qMin, _ := subData["qty_minimum"].(int)
			qMax, _ := subData["qty_maximum"].(int)
			hidden, _ := subData["hidden"].(int)
			tx.Exec("UPDATE product_config_options_sub SET option_name = ?, sort_order = ?, qty_minimum = ?, qty_maximum = ?, hidden = ? WHERE id = ?",
				name, sortOrder, qMin, qMax, hidden, subID)
		}

		// Update pricing
		for currIDStr, subPrices := range prices {
			currID, _ := strconv.ParseUint(currIDStr, 10, 64)
			for subIDStr, priceData := range subPrices {
				subID, _ := strconv.ParseUint(subIDStr, 10, 64)
				sets := []string{}
				vals := []interface{}{}
				for col, val := range priceData {
					sets = append(sets, col+" = ?")
					vals = append(vals, val)
				}
				if len(sets) > 0 {
					vals = append(vals, "configoptions", subID, currID)
					tx.Exec("UPDATE pricing SET "+strings.Join(sets, ", ")+" WHERE type = ? AND relid = ? AND currency = ?", vals...)
				}
			}
		}

		// Update product versions
		tx.Exec(`UPDATE products SET location_version = location_version + 1 WHERE id IN (
			SELECT pid FROM product_config_links WHERE gid IN (
				SELECT gid FROM product_config_options WHERE id = ?
			)
		)`, cid)
		return nil
	})
}

// AdminSaveLinkAgeLevel saves a linkage level config option.
func (s *ConfigOptionService) AdminSaveLinkAgeLevel(gid uint, optionName string, optionType int, notes string, linkagePID uint, subOptionName string, subLinkagePID uint, hidden int, optionID uint, subOptionID uint) (uint, uint, error) {
	var newOptionID, newSubOptionID uint
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Create or update option
		if optionID > 0 {
			tx.Exec(`UPDATE product_config_options SET gid=?, option_name=?, option_type=?, notes=?, linkage_pid=? WHERE id=?`,
				gid, optionName, optionType, notes, linkagePID, optionID)
			newOptionID = optionID
		} else {
			res := tx.Exec(`INSERT INTO product_config_options (gid, option_name, option_type, notes, linkage_pid) VALUES (?,?,?,?,?)`,
				gid, optionName, optionType, notes, linkagePID)
			var lid struct{ ID uint }
			tx.Raw("SELECT LAST_INSERT_ID() as id").Scan(&lid)
			newOptionID = lid.ID
			_ = res
		}

		// Update linkage_level for option
		if linkagePID == 0 {
			tx.Exec("UPDATE product_config_options SET linkage_level = ?, linkage_top_pid = 0 WHERE id = ?",
				"0-"+strconv.FormatUint(uint64(newOptionID), 10), newOptionID)
		} else {
			var parent struct{ LinkageLevel string; LinkageTopPID uint; LinkagePID uint }
			tx.Raw("SELECT linkage_level, linkage_top_pid, linkage_pid FROM product_config_options WHERE id = ?", linkagePID).Scan(&parent)
			topPID := parent.LinkageTopPID
			if parent.LinkagePID == 0 {
				topPID = linkagePID
			}
			tx.Exec("UPDATE product_config_options SET linkage_top_pid = ?, linkage_level = ? WHERE id = ?",
				topPID, parent.LinkageLevel+"-"+strconv.FormatUint(uint64(newOptionID), 10), newOptionID)
		}

		// Create or update sub-option
		subData := map[string]interface{}{
			"option_name":  subOptionName,
			"qty_minimum":  0,
			"qty_maximum":  0,
			"hidden":       hidden,
			"config_id":    newOptionID,
			"linkage_pid":  subLinkagePID,
		}
		if subOptionID > 0 {
			tx.Exec("UPDATE product_config_options_sub SET option_name=?, qty_minimum=?, qty_maximum=?, hidden=?, config_id=?, linkage_pid=? WHERE id=?",
				subData["option_name"], subData["qty_minimum"], subData["qty_maximum"], subData["hidden"], subData["config_id"], subData["linkage_pid"], subOptionID)
			newSubOptionID = subOptionID
		} else {
			res := tx.Exec("INSERT INTO product_config_options_sub (option_name, qty_minimum, qty_maximum, hidden, config_id, linkage_pid) VALUES (?,?,?,?,?,?)",
				subData["option_name"], subData["qty_minimum"], subData["qty_maximum"], subData["hidden"], subData["config_id"], subData["linkage_pid"])
			var sid struct{ ID uint }
			tx.Raw("SELECT LAST_INSERT_ID() as id").Scan(&sid)
			newSubOptionID = sid.ID
			_ = res

			// Create pricing for all currencies
			var currIDs []uint
			tx.Table("currencies").Pluck("id", &currIDs)
			for _, cid := range currIDs {
				tx.Exec("INSERT INTO pricing (type, relid, currency, monthly) VALUES ('configoptions', ?, ?, 0)", newSubOptionID, cid)
			}
		}

		// Update linkage_level for sub-option
		if subLinkagePID == 0 {
			tx.Exec("UPDATE product_config_options_sub SET linkage_level = ? WHERE id = ?",
				"0-"+strconv.FormatUint(uint64(newSubOptionID), 10), newSubOptionID)
		} else {
			var subParent struct{ LinkageLevel string; LinkageTopPID uint; LinkagePID uint }
			tx.Raw("SELECT linkage_level, linkage_top_pid, linkage_pid FROM product_config_options_sub WHERE id = ?", subLinkagePID).Scan(&subParent)
			topPID := subParent.LinkageTopPID
			if subParent.LinkagePID == 0 {
				topPID = subLinkagePID
			}
			tx.Exec("UPDATE product_config_options_sub SET linkage_top_pid = ?, linkage_level = ? WHERE id = ?",
				topPID, subParent.LinkageLevel+"-"+strconv.FormatUint(uint64(newSubOptionID), 10), newSubOptionID)
		}

		// Update product versions
		tx.Exec("UPDATE products SET location_version = location_version + 1 WHERE id IN (SELECT pid FROM product_config_links WHERE gid = ?)", gid)
		return nil
	})
	return newOptionID, newSubOptionID, err
}

// AdminSaveConfigOptionInfo saves basic config option info.
func (s *ConfigOptionService) AdminSaveConfigOptionInfo(gid uint, optionName string, optionType int, notes string, qtyStage int, linkagePID uint, isDiscount int, unit string, senior int, optionID uint) (uint, error) {
	var newOptionID uint
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Check if group is upstream
		var group struct{ UpstreamID uint }
		tx.Raw("SELECT upstream_id FROM product_config_groups WHERE id = ?", gid).Scan(&group)
		if group.UpstreamID > 0 {
			return fmt.Errorf("上游配置项组,不可添加配置项")
		}

		if optionID > 0 {
			tx.Exec(`UPDATE product_config_options SET gid=?, option_name=?, option_type=?, notes=?, qty_stage=?, linkage_pid=?, is_discount=?, unit=?, senior=? WHERE id=?`,
				gid, optionName, optionType, notes, qtyStage, linkagePID, isDiscount, unit, senior, optionID)
			newOptionID = optionID
		} else {
			res := tx.Exec(`INSERT INTO product_config_options (gid, option_name, option_type, notes, qty_stage, linkage_pid, is_discount, unit, senior) VALUES (?,?,?,?,?,?,?,?,?)`,
				gid, optionName, optionType, notes, qtyStage, linkagePID, isDiscount, unit, senior)
			var lid struct{ ID uint }
			tx.Raw("SELECT LAST_INSERT_ID() as id").Scan(&lid)
			newOptionID = lid.ID
			_ = res
		}

		// Update linkage_level
		if linkagePID == 0 {
			tx.Exec("UPDATE product_config_options SET linkage_level = ?, linkage_top_pid = 0 WHERE id = ?",
				"0-"+strconv.FormatUint(uint64(newOptionID), 10), newOptionID)
		} else {
			var parent struct{ LinkageLevel string; LinkageTopPID uint; LinkagePID uint }
			tx.Raw("SELECT linkage_level, linkage_top_pid, linkage_pid FROM product_config_options WHERE id = ?", linkagePID).Scan(&parent)
			topPID := parent.LinkageTopPID
			if parent.LinkagePID == 0 {
				topPID = linkagePID
			}
			tx.Exec("UPDATE product_config_options SET linkage_top_pid = ?, linkage_level = ? WHERE id = ?",
				topPID, parent.LinkageLevel+"-"+strconv.FormatUint(uint64(newOptionID), 10), newOptionID)
		}

		tx.Exec("UPDATE products SET location_version = location_version + 1 WHERE id IN (SELECT pid FROM product_config_links WHERE gid = ?)", gid)
		return nil
	})
	return newOptionID, err
}

// AdminSaveLinkAgeOrder saves sort order for linkage sub-options.
func (s *ConfigOptionService) AdminSaveLinkAgeOrder(subIDs []uint) error {
	for i, subID := range subIDs {
		s.db.Exec("UPDATE product_config_options_sub SET sort_order = ? WHERE id = ?", i, subID)
	}
	return nil
}

// AdminDelLinkAgeSub deletes a linkage sub-option.
func (s *ConfigOptionService) AdminDelLinkAgeSub(subID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var sub struct{ UpstreamID uint }
		tx.Raw("SELECT upstream_id FROM product_config_options_sub WHERE id = ?", subID).Scan(&sub)
		if sub.UpstreamID > 0 {
			return fmt.Errorf("上游配置子项,不可删除")
		}

		tx.Exec("DELETE FROM pricing WHERE type = 'configoptions' AND relid = ?", subID)
		tx.Exec("DELETE FROM product_config_options_sub WHERE id = ?", subID)

		tx.Exec(`UPDATE products SET location_version = location_version + 1 WHERE id IN (
			SELECT pid FROM product_config_links WHERE gid IN (
				SELECT gid FROM product_config_options WHERE id IN (
					SELECT config_id FROM product_config_options_sub WHERE id = ?
				)
			)
		)`, subID)
		return nil
	})
}

// AdminCheckOS returns products that have OS-type config options (option_type=5).
func (s *ConfigOptionService) AdminCheckOS() ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	if err := s.db.Table("products a").
		Select("a.id as pid, a.name, s.server_type as type").
		Joins("LEFT JOIN servers s ON s.gid = a.server_group").
		Joins("LEFT JOIN product_config_links b ON a.id = b.pid").
		Joins("LEFT JOIN product_config_options c ON b.gid = c.gid").
		Where("(s.server_type = 'dcimcloud' OR s.server_type = 'dcim') AND s.gid > 0 AND a.server_group > 0 AND c.option_type = 5 AND (a.api_type = 'normal' OR a.api_type = '')").
		Order("a.id DESC").
		Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

// AdminGetOS returns OS options for a product.
func (s *ConfigOptionService) AdminGetOS(pid uint) (map[string]interface{}, error) {
	var product struct {
		ServerGroup uint
		CID         uint
		APIType     string
		Type        string
	}
	s.db.Table("products a").
		Select("a.server_group, c.id as cid, a.api_type, a.type").
		Joins("LEFT JOIN product_config_links b ON a.id = b.pid").
		Joins("LEFT JOIN product_config_options c ON b.gid = c.gid").
		Where("a.id = ? AND c.option_type = 5", pid).
		Find(&product)

	if product.CID == 0 {
		return nil, fmt.Errorf("未找到OS配置项")
	}
	if product.Type != "dcim" && product.Type != "dcimcloud" {
		return nil, fmt.Errorf("只有本地接口才能拉取")
	}

	// Return current OS sub-options
	var subOptions []map[string]interface{}
	s.db.Table("product_config_options_sub").
		Where("config_id = ?", product.CID).
		Find(&subOptions)

	return map[string]interface{}{
		"pid":         pid,
		"cid":         product.CID,
		"sub_options": subOptions,
	}, nil
}
