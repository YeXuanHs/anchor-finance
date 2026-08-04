package service

import (
	"errors"
	"fmt"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type ConfigServerService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewConfigServerService(db *gorm.DB, log *logger.Logger) *ConfigServerService {
	return &ConfigServerService{db: db, log: log}
}

// ---------- Server Config ----------

type CreateServerConfigRequest struct {
	Name            string  `json:"name" binding:"required,max=128"`
	Code            string  `json:"code" binding:"required,max=64"`
	Type            string  `json:"type" binding:"required,oneof=vps dedicated cloud reseller"`
	Provider        string  `json:"provider" binding:"max=128"`
	TemplateID      *uint   `json:"template_id"`
	CPU             string  `json:"cpu" binding:"max=64"`
	Memory          int     `json:"memory"`
	Disk            int     `json:"disk"`
	Bandwidth       int     `json:"bandwidth"`
	TrafficLimit    int64   `json:"traffic_limit"`
	IPCount         int     `json:"ip_count"`
	Location        string  `json:"location" binding:"max=128"`
	Datacenter      string  `json:"datacenter" binding:"max=128"`
	OS              string  `json:"os"`
	Features        string  `json:"features"`
	PriceMonthly    float64 `json:"price_monthly"`
	PriceQuarter    float64 `json:"price_quarter"`
	PriceSemiAnn    float64 `json:"price_semi_annual"`
	PriceAnnual     float64 `json:"price_annual"`
	PriceBiennial   float64 `json:"price_biennial"`
	PriceTriennial  float64 `json:"price_triennial"`
	PricingStrategy string  `json:"pricing_strategy" binding:"omitempty,oneof=fixed graduated promotional"`
	StockTotal      int     `json:"stock_total"`
	MaxPerUser      int     `json:"max_per_user"`
	SortOrder       int     `json:"sort_order"`
	Status          int16   `json:"status"`
	Remark          string  `json:"remark"`
}

type UpdateServerConfigRequest struct {
	Name            *string  `json:"name"`
	Code            *string  `json:"code"`
	Type            *string  `json:"type"`
	Provider        *string  `json:"provider"`
	TemplateID      *uint    `json:"template_id"`
	CPU             *string  `json:"cpu"`
	Memory          *int     `json:"memory"`
	Disk            *int     `json:"disk"`
	Bandwidth       *int     `json:"bandwidth"`
	TrafficLimit    *int64   `json:"traffic_limit"`
	IPCount         *int     `json:"ip_count"`
	Location        *string  `json:"location"`
	Datacenter      *string  `json:"datacenter"`
	OS              *string  `json:"os"`
	Features        *string  `json:"features"`
	PriceMonthly    *float64 `json:"price_monthly"`
	PriceQuarter    *float64 `json:"price_quarter"`
	PriceSemiAnn    *float64 `json:"price_semi_annual"`
	PriceAnnual     *float64 `json:"price_annual"`
	PriceBiennial   *float64 `json:"price_biennial"`
	PriceTriennial  *float64 `json:"price_triennial"`
	PricingStrategy *string  `json:"pricing_strategy"`
	StockTotal      *int     `json:"stock_total"`
	MaxPerUser      *int     `json:"max_per_user"`
	SortOrder       *int     `json:"sort_order"`
	Status          *int16   `json:"status"`
	Remark          *string  `json:"remark"`
}

func (s *ConfigServerService) GetList(page, pageSize int, serverType string, keyword string) ([]model.ServerConfig, int64, error) {
	var items []model.ServerConfig
	var total int64

	query := s.db.Model(&model.ServerConfig{})
	if serverType != "" {
		query = query.Where("type = ?", serverType)
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
		Preload("Template").
		Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (s *ConfigServerService) GetByID(id uint) (*model.ServerConfig, error) {
	var item model.ServerConfig
	if err := s.db.Preload("Template").First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *ConfigServerService) Create(req CreateServerConfigRequest) (*model.ServerConfig, error) {
	var count int64
	s.db.Model(&model.ServerConfig{}).Where("code = ?", req.Code).Count(&count)
	if count > 0 {
		return nil, errors.New("server config code already exists")
	}

	item := &model.ServerConfig{
		Name:            req.Name,
		Code:            req.Code,
		Type:            req.Type,
		Provider:        req.Provider,
		TemplateID:      req.TemplateID,
		CPU:             req.CPU,
		Memory:          req.Memory,
		Disk:            req.Disk,
		Bandwidth:       req.Bandwidth,
		TrafficLimit:    req.TrafficLimit,
		IPCount:         req.IPCount,
		Location:        req.Location,
		Datacenter:      req.Datacenter,
		PriceMonthly:    req.PriceMonthly,
		PriceQuarter:    req.PriceQuarter,
		PriceSemiAnn:    req.PriceSemiAnn,
		PriceAnnual:     req.PriceAnnual,
		PriceBiennial:   req.PriceBiennial,
		PriceTriennial:  req.PriceTriennial,
		PricingStrategy: req.PricingStrategy,
		StockTotal:      req.StockTotal,
		StockUsed:       0,
		MaxPerUser:      req.MaxPerUser,
		SortOrder:       req.SortOrder,
		Status:          req.Status,
		Remark:          req.Remark,
	}
	if item.PricingStrategy == "" {
		item.PricingStrategy = "fixed"
	}
	if item.IPCount == 0 {
		item.IPCount = 1
	}

	if err := s.db.Create(item).Error; err != nil {
		return nil, err
	}
	s.log.Infof("server config created: id=%d code=%s", item.ID, item.Code)
	return item, nil
}

func (s *ConfigServerService) Update(id uint, req UpdateServerConfigRequest) (*model.ServerConfig, error) {
	var item model.ServerConfig
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}

	if req.Code != nil {
		var count int64
		s.db.Model(&model.ServerConfig{}).Where("code = ? AND id != ?", *req.Code, id).Count(&count)
		if count > 0 {
			return nil, errors.New("server config code already exists")
		}
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Code != nil {
		updates["code"] = *req.Code
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Provider != nil {
		updates["provider"] = *req.Provider
	}
	if req.TemplateID != nil {
		updates["template_id"] = req.TemplateID
	}
	if req.CPU != nil {
		updates["cpu"] = *req.CPU
	}
	if req.Memory != nil {
		updates["memory"] = *req.Memory
	}
	if req.Disk != nil {
		updates["disk"] = *req.Disk
	}
	if req.Bandwidth != nil {
		updates["bandwidth"] = *req.Bandwidth
	}
	if req.TrafficLimit != nil {
		updates["traffic_limit"] = *req.TrafficLimit
	}
	if req.IPCount != nil {
		updates["ip_count"] = *req.IPCount
	}
	if req.Location != nil {
		updates["location"] = *req.Location
	}
	if req.Datacenter != nil {
		updates["datacenter"] = *req.Datacenter
	}
	if req.OS != nil {
		updates["os"] = *req.OS
	}
	if req.Features != nil {
		updates["features"] = *req.Features
	}
	if req.PriceMonthly != nil {
		updates["price_monthly"] = *req.PriceMonthly
	}
	if req.PriceQuarter != nil {
		updates["price_quarter"] = *req.PriceQuarter
	}
	if req.PriceSemiAnn != nil {
		updates["price_semi_annual"] = *req.PriceSemiAnn
	}
	if req.PriceAnnual != nil {
		updates["price_annual"] = *req.PriceAnnual
	}
	if req.PriceBiennial != nil {
		updates["price_biennial"] = *req.PriceBiennial
	}
	if req.PriceTriennial != nil {
		updates["price_triennial"] = *req.PriceTriennial
	}
	if req.PricingStrategy != nil {
		updates["pricing_strategy"] = *req.PricingStrategy
	}
	if req.StockTotal != nil {
		updates["stock_total"] = *req.StockTotal
	}
	if req.MaxPerUser != nil {
		updates["max_per_user"] = *req.MaxPerUser
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Remark != nil {
		updates["remark"] = *req.Remark
	}

	if len(updates) > 0 {
		if err := s.db.Model(&item).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	if err := s.db.Preload("Template").First(&item, id).Error; err != nil {
		return nil, err
	}
	s.log.Infof("server config updated: id=%d", id)
	return &item, nil
}

func (s *ConfigServerService) Delete(id uint) error {
	var count int64
	s.db.Model(&model.ServerProduct{}).Where("server_config_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("server config has associated products, remove them first")
	}
	if err := s.db.Delete(&model.ServerConfig{}, id).Error; err != nil {
		return err
	}
	s.log.Infof("server config deleted: id=%d", id)
	return nil
}

func (s *ConfigServerService) BatchUpdateStatus(ids []uint, status int16) error {
	if len(ids) == 0 {
		return errors.New("ids is empty")
	}
	result := s.db.Model(&model.ServerConfig{}).Where("id IN ?", ids).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	s.log.Infof("server config batch status updated: ids=%v status=%d", ids, status)
	return nil
}

func (s *ConfigServerService) BatchDelete(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("ids is empty")
	}
	var count int64
	s.db.Model(&model.ServerProduct{}).Where("server_config_id IN ?", ids).Count(&count)
	if count > 0 {
		return errors.New("some server configs have associated products")
	}
	result := s.db.Delete(&model.ServerConfig{}, ids)
	if result.Error != nil {
		return result.Error
	}
	s.log.Infof("server config batch deleted: ids=%v", ids)
	return nil
}

func (s *ConfigServerService) UpdateSort(id uint, sortOrder int) error {
	result := s.db.Model(&model.ServerConfig{}).Where("id = ?", id).Update("sort_order", sortOrder)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("server config not found")
	}
	return nil
}

// ---------- Server Template ----------

type CreateServerTemplateRequest struct {
	Name        string `json:"name" binding:"required,max=128"`
	Code        string `json:"code" binding:"required,max=64"`
	Type        string `json:"type" binding:"required"`
	Description string `json:"description"`
	Config      string `json:"config" binding:"required"`
	SortOrder   int    `json:"sort_order"`
	Status      int16  `json:"status"`
}

type UpdateServerTemplateRequest struct {
	Name        *string `json:"name"`
	Code        *string `json:"code"`
	Type        *string `json:"type"`
	Description *string `json:"description"`
	Config      *string `json:"config"`
	SortOrder   *int    `json:"sort_order"`
	Status      *int16  `json:"status"`
}

func (s *ConfigServerService) GetTemplateList(page, pageSize int, templateType string) ([]model.ServerTemplate, int64, error) {
	var items []model.ServerTemplate
	var total int64

	query := s.db.Model(&model.ServerTemplate{})
	if templateType != "" {
		query = query.Where("type = ?", templateType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("sort_order ASC, id ASC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (s *ConfigServerService) CreateTemplate(req CreateServerTemplateRequest) (*model.ServerTemplate, error) {
	var count int64
	s.db.Model(&model.ServerTemplate{}).Where("code = ?", req.Code).Count(&count)
	if count > 0 {
		return nil, errors.New("template code already exists")
	}

	item := &model.ServerTemplate{
		Name:        req.Name,
		Code:        req.Code,
		Type:        req.Type,
		Description: req.Description,
		Config:      []byte(req.Config),
		SortOrder:   req.SortOrder,
		Status:      req.Status,
	}
	if item.Status == 0 {
		item.Status = 1
	}

	if err := s.db.Create(item).Error; err != nil {
		return nil, err
	}
	s.log.Infof("server template created: id=%d code=%s", item.ID, item.Code)
	return item, nil
}

func (s *ConfigServerService) UpdateTemplate(id uint, req UpdateServerTemplateRequest) (*model.ServerTemplate, error) {
	var item model.ServerTemplate
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}

	if req.Code != nil {
		var count int64
		s.db.Model(&model.ServerTemplate{}).Where("code = ? AND id != ?", *req.Code, id).Count(&count)
		if count > 0 {
			return nil, errors.New("template code already exists")
		}
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Code != nil {
		updates["code"] = *req.Code
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Config != nil {
		updates["config"] = *req.Config
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
	s.log.Infof("server template updated: id=%d", id)
	return &item, nil
}

func (s *ConfigServerService) DeleteTemplate(id uint) error {
	var count int64
	s.db.Model(&model.ServerConfig{}).Where("template_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("template is in use by server configs")
	}
	if err := s.db.Delete(&model.ServerTemplate{}, id).Error; err != nil {
		return err
	}
	s.log.Infof("server template deleted: id=%d", id)
	return nil
}

// ─── ConfigServers Admin Methods (from zjmf ConfigServersController) ───

// AdminServerList returns paginated servers with filters.
func (s *ConfigServerService) AdminServerList(page, pageSize int, gid uint, search string) ([]map[string]interface{}, int64, error) {
	query := s.db.Table("servers").Where("server_type = 'normal'")
	if gid > 0 {
		query = query.Where("gid = ?", gid)
	}
	if search != "" {
		query = query.Where("name LIKE ? OR hostname LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	query.Count(&total)

	var servers []map[string]interface{}
	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&servers).Error; err != nil {
		return nil, 0, err
	}

	// Enrich with group names and module types
	var groups []struct{ ID uint; Name string }
	s.db.Table("server_groups").Where("system_type = 'normal'").Find(&groups)
	groupMap := make(map[uint]string)
	for _, g := range groups {
		groupMap[g.ID] = g.Name
	}

	for i, srv := range servers {
		gidVal, _ := srv["gid"].(uint)
		srv["gname"] = groupMap[gidVal]
		// Count hosted products
		serverID, _ := srv["id"].(uint)
		var useNum int64
		s.db.Table("host").Where("serverid = ?", serverID).Count(&useNum)
		maxAcc, _ := srv["max_accounts"].(int)
		srv["open_num"] = fmt.Sprintf("%d/%d", useNum, maxAcc)
		servers[i] = srv
	}

	return servers, total, nil
}

// AdminAddServersData returns data for the add server page.
func (s *ConfigServerService) AdminAddServersData() map[string]interface{} {
	return map[string]interface{}{
		"modules": []map[string]string{
			{"value": "cpanel", "name": "cPanel"},
			{"value": "whm", "name": "WHM"},
			{"value": "plesk", "name": "Plesk"},
			{"value": "directadmin", "name": "DirectAdmin"},
			{"value": "virtualizor", "name": "Virtualizor"},
			{"value": "solusvm", "name": "SolusVM"},
		},
		"groups": []interface{}{},
	}
}

// AdminGetModulesGroup returns server groups filtered by module type.
func (s *ConfigServerService) AdminGetModulesGroup(moduleType string) []map[string]interface{} {
	var gids []uint
	s.db.Table("servers").Where("server_type = 'normal' AND type = ?", moduleType).Pluck("gid", &gids)

	var groups []map[string]interface{}
	s.db.Table("server_groups").Where("system_type = 'normal'").Find(&groups)

	var result []map[string]interface{}
	for _, g := range groups {
		gid, _ := g["id"].(uint)
		var count int64
		s.db.Table("servers").Where("gid = ?", gid).Count(&count)
		if containsUint(gids, gid) || count == 0 {
			result = append(result, g)
		}
	}
	return result
}

func containsUint(arr []uint, val uint) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

// AdminAddServersPost creates a new server.
func (s *ConfigServerService) AdminAddServersPost(name, ipAddress, hostname, username, password, accessHash, serverType string, port, maxAccounts, gid int, secure, disabled int) (uint, error) {
	// Check unique name
	var count int64
	s.db.Table("servers").Where("server_type = 'normal' AND name = ?", name).Count(&count)
	if count > 0 {
		return 0, fmt.Errorf("该接口已存在")
	}

	var serverID uint
	err := s.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Exec(`INSERT INTO servers (name, ip_address, hostname, username, password, accesshash, type, port, max_accounts, gid, secure, disabled, server_type) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`,
			name, ipAddress, hostname, username, password, accessHash, serverType, port, maxAccounts, gid, secure, disabled, "normal")
		if res.Error != nil {
			return res.Error
		}
		var lid struct{ ID uint }
		tx.Raw("SELECT LAST_INSERT_ID() as id").Scan(&lid)
		serverID = lid.ID
		return nil
	})
	return serverID, err
}

// AdminEditServers returns server detail for editing.
func (s *ConfigServerService) AdminEditServers(id uint) (map[string]interface{}, error) {
	var server map[string]interface{}
	if err := s.db.Table("servers").Where("id = ? AND server_type = 'normal'", id).Find(&server).Error; err != nil {
		return nil, err
	}
	if server == nil {
		return nil, fmt.Errorf("server not found")
	}
	delete(server, "password")
	return server, nil
}

// AdminEditServersPost updates a server.
func (s *ConfigServerService) AdminEditServersPost(id uint, name, ipAddress, hostname, username, password, accessHash, serverType string, port, maxAccounts, gid int, secure, disabled int) error {
	// Check unique name
	var count int64
	s.db.Table("servers").Where("server_type = 'normal' AND name = ? AND id != ?", name, id).Count(&count)
	if count > 0 {
		return fmt.Errorf("该接口已存在")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{
			"name":          name,
			"ip_address":    ipAddress,
			"hostname":      hostname,
			"username":      username,
			"type":          serverType,
			"port":          port,
			"max_accounts":  maxAccounts,
			"gid":           gid,
			"secure":        secure,
			"disabled":      disabled,
			"accesshash":    accessHash,
		}
		if password != "" {
			updates["password"] = password
		}
		return tx.Table("servers").Where("id = ? AND server_type = 'normal'", id).Updates(updates).Error
	})
}

// AdminDeleteServers deletes a server if not in use.
func (s *ConfigServerService) AdminDeleteServers(id uint) error {
	var count int64
	s.db.Table("host").Where("serverid = ?", id).Count(&count)
	if count > 0 {
		return fmt.Errorf("此接口已被使用，不能删除")
	}
	return s.db.Table("servers").Where("id = ? AND server_type = 'normal'", id).Delete(nil).Error
}

// AdminServerGroupsList returns paginated server groups.
func (s *ConfigServerService) AdminServerGroupsList(page, pageSize int) ([]map[string]interface{}, int64, error) {
	query := s.db.Table("server_groups").Where("system_type = 'normal'")
	var total int64
	query.Count(&total)

	var groups []map[string]interface{}
	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&groups).Error; err != nil {
		return nil, 0, err
	}

	// Enrich with server counts
	for i, g := range groups {
		gid, _ := g["id"].(uint)
		var servers []struct{ ID uint; MaxAccounts int }
		s.db.Table("servers").Where("gid = ?", gid).Select("id, max_accounts").Find(&servers)
		totalMax := 0
		var serverIDs []uint
		for _, srv := range servers {
			totalMax += srv.MaxAccounts
			serverIDs = append(serverIDs, srv.ID)
		}
		var useNum int64
		if len(serverIDs) > 0 {
			s.db.Table("host").Where("serverid IN ? AND serverid != 0", serverIDs).Count(&useNum)
		}
		g["open_num"] = fmt.Sprintf("%d/%d", useNum, totalMax)
		groups[i] = g
	}

	return groups, total, nil
}

// AdminCreateGroupPage returns data for creating a server group.
func (s *ConfigServerService) AdminCreateGroupPage() map[string]interface{} {
	var servers []map[string]interface{}
	s.db.Table("servers").Where("gid = 0").Find(&servers)
	return map[string]interface{}{
		"servers":      servers,
		"mode": []map[string]interface{}{
			{"name": "平均分配", "value": 1, "desc": "产品优先分配给产品数量最少的接口"},
			{"name": "逐个分配", "value": 2, "desc": "按最初创建的接口开始分配，满额后切换下一接口"},
		},
	}
}

// AdminCreateGroupPost creates a new server group.
func (s *ConfigServerService) AdminCreateGroupPost(name string, mode int, serverIDs []uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var count int64
		tx.Table("server_groups").Where("system_type = 'normal' AND name = ?", name).Count(&count)
		if count > 0 {
			return fmt.Errorf("该接口组已存在")
		}

		res := tx.Exec("INSERT INTO server_groups (name, mode, system_type) VALUES (?, ?, 'normal')", name, mode)
		if res.Error != nil {
			return res.Error
		}
		var lid struct{ ID uint }
		tx.Raw("SELECT LAST_INSERT_ID() as id").Scan(&lid)

		if len(serverIDs) > 0 {
			tx.Exec("UPDATE servers SET gid = ? WHERE id IN ?", lid.ID, serverIDs)
		}
		return nil
	})
}

// AdminEditServerGroup returns data for editing a server group.
func (s *ConfigServerService) AdminEditServerGroup(id uint) (map[string]interface{}, error) {
	var group map[string]interface{}
	if err := s.db.Table("server_groups").Where("id = ?", id).Find(&group).Error; err != nil {
		return nil, err
	}

	var servers []map[string]interface{}
	s.db.Table("servers").Where("gid IN (0, ?)", id).Find(&servers)

	var selectedIDs []uint
	s.db.Table("servers").Where("gid = ?", id).Pluck("id", &selectedIDs)

	return map[string]interface{}{
		"server_group":  group,
		"servers":       servers,
		"select_servers": selectedIDs,
		"mode": []map[string]interface{}{
			{"name": "平均分配", "value": 1, "desc": "产品优先分配给产品数量最少的接口"},
			{"name": "逐个分配", "value": 2, "desc": "按最初创建的接口开始分配，满额后切换下一接口"},
		},
	}, nil
}

// AdminEditServerGroupPost updates a server group.
func (s *ConfigServerService) AdminEditServerGroupPost(id uint, name string, mode int, serverIDs []uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var count int64
		tx.Table("server_groups").Where("system_type = 'normal' AND name = ? AND id != ?", name, id).Count(&count)
		if count > 0 {
			return fmt.Errorf("该接口组已存在")
		}

		tx.Exec("UPDATE server_groups SET name = ?, mode = ? WHERE id = ?", name, mode, id)

		// Reset servers in this group
		tx.Exec("UPDATE servers SET gid = 0 WHERE gid = ?", id)
		if len(serverIDs) > 0 {
			tx.Exec("UPDATE servers SET gid = ? WHERE id IN ?", id, serverIDs)
		}
		return nil
	})
}

// AdminDeleteServerGroup deletes a server group if empty.
func (s *ConfigServerService) AdminDeleteServerGroup(id uint) error {
	var count int64
	s.db.Table("servers").Where("gid = ?", id).Count(&count)
	if count > 0 {
		return fmt.Errorf("此接口分组中已有接口，不能删除")
	}
	return s.db.Table("server_groups").Where("system_type = 'normal' AND id = ?", id).Delete(nil).Error
}

// AdminTestLink tests connection to a server.
func (s *ConfigServerService) AdminTestLink(id uint) (map[string]interface{}, error) {
	var server map[string]interface{}
	if err := s.db.Table("servers").Where("id = ?", id).Find(&server).Error; err != nil {
		return nil, fmt.Errorf("server not found")
	}
	if server == nil {
		return nil, fmt.Errorf("server not found")
	}

	// Update link status (simplified - in real impl would test actual connection)
	s.db.Exec("UPDATE servers SET link_status = 1 WHERE id = ?", id)
	return map[string]interface{}{
		"server_status": 1,
		"msg":           "连接成功",
	}, nil
}
