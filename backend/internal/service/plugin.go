package service

import (
	"errors"
	"fmt"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// PluginService 插件服务
type PluginService struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewPluginService 创建插件服务
func NewPluginService(db *gorm.DB, log *logger.Logger) *PluginService {
	return &PluginService{db: db, log: log}
}

// CreatePluginRequest 创建插件请求
type CreatePluginRequest struct {
	Name        string `json:"name" binding:"required,max=64"`
	Title       string `json:"title" binding:"required,max=64"`
	Type        string `json:"type" binding:"required,oneof=mail sms certification gateway oauth server addon"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Version     string `json:"version"`
	HelpURL     string `json:"help_url"`
	Config      string `json:"config"`
	Module      string `json:"module"`
	IsSystem    bool   `json:"is_system"`
	IsEnabled   bool   `json:"is_enabled"`
	SortOrder   int    `json:"sort_order"`
}

// UpdatePluginRequest 更新插件请求
type UpdatePluginRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Author      *string `json:"author"`
	Version     *string `json:"version"`
	HelpURL     *string `json:"help_url"`
	Config      *string `json:"config"`
	IsEnabled   *bool   `json:"is_enabled"`
	SortOrder   *int    `json:"sort_order"`
}

// Create 创建插件
func (s *PluginService) Create(req CreatePluginRequest) (*model.Plugin, error) {
	var count int64
	s.db.Model(&model.Plugin{}).Where("name = ?", req.Name).Count(&count)
	if count > 0 {
		return nil, errors.New("插件标识名已存在")
	}

	plugin := model.Plugin{
		Name:        req.Name,
		Title:       req.Title,
		Type:        req.Type,
		Description: req.Description,
		Author:      req.Author,
		Version:     req.Version,
		HelpURL:     req.HelpURL,
		Config:      req.Config,
		Module:      req.Module,
		IsSystem:    req.IsSystem,
		IsEnabled:   req.IsEnabled,
		SortOrder:   req.SortOrder,
	}

	if err := s.db.Create(&plugin).Error; err != nil {
		return nil, err
	}

	s.log.Infof("插件创建成功: id=%d name=%s", plugin.ID, plugin.Name)
	return &plugin, nil
}

// Update 更新插件
func (s *PluginService) Update(id uint, req UpdatePluginRequest) (*model.Plugin, error) {
	var plugin model.Plugin
	if err := s.db.First(&plugin, id).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Author != nil {
		updates["author"] = *req.Author
	}
	if req.Version != nil {
		updates["version"] = *req.Version
	}
	if req.HelpURL != nil {
		updates["help_url"] = *req.HelpURL
	}
	if req.Config != nil {
		updates["config"] = *req.Config
	}
	if req.IsEnabled != nil {
		updates["is_enabled"] = *req.IsEnabled
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}

	if len(updates) > 0 {
		if err := s.db.Model(&plugin).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	if err := s.db.First(&plugin, id).Error; err != nil {
		return nil, err
	}

	s.log.Infof("插件更新成功: id=%d", id)
	return &plugin, nil
}

// Delete 删除插件
func (s *PluginService) Delete(id uint) error {
	var plugin model.Plugin
	if err := s.db.First(&plugin, id).Error; err != nil {
		return errors.New("插件不存在")
	}
	if plugin.IsSystem {
		return errors.New("系统内置插件不能删除")
	}

	result := s.db.Delete(&model.Plugin{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("插件不存在")
	}
	s.log.Infof("插件删除成功: id=%d", id)
	return nil
}

// GetByID 根据ID获取插件
func (s *PluginService) GetByID(id uint) (*model.Plugin, error) {
	var plugin model.Plugin
	if err := s.db.First(&plugin, id).Error; err != nil {
		return nil, err
	}
	return &plugin, nil
}

// GetByName 根据名称获取插件
func (s *PluginService) GetByName(name string) (*model.Plugin, error) {
	var plugin model.Plugin
	if err := s.db.Where("name = ?", name).First(&plugin).Error; err != nil {
		return nil, err
	}
	return &plugin, nil
}

// GetList 获取插件列表
func (s *PluginService) GetList(pluginType string, isEnabled *bool) ([]model.Plugin, error) {
	var plugins []model.Plugin
	query := s.db.Model(&model.Plugin{})

	if pluginType != "" {
		query = query.Where("type = ?", pluginType)
	}
	if isEnabled != nil {
		query = query.Where("is_enabled = ?", *isEnabled)
	}

	if err := query.Order("sort_order ASC, id ASC").Find(&plugins).Error; err != nil {
		return nil, err
	}
	return plugins, nil
}

// SetEnabled 设置启用状态（启用时自动建表）
func (s *PluginService) SetEnabled(id uint, enabled bool) error {
	var plugin model.Plugin
	if err := s.db.First(&plugin, id).Error; err != nil {
		return errors.New("插件不存在")
	}

	// 启用插件时自动建表
	if enabled && !plugin.IsEnabled {
		if err := s.createPluginTables(plugin.Name); err != nil {
			s.log.Errorf("插件 %s 建表失败: %v", plugin.Name, err)
			return fmt.Errorf("启用失败: 建表出错 - %v", err)
		}
		s.log.Infof("插件 %s 表结构已创建", plugin.Name)
	}

	result := s.db.Model(&plugin).Update("is_enabled", enabled)
	if result.Error != nil {
		return result.Error
	}
	s.log.Infof("插件状态更新: id=%d name=%s enabled=%v", id, plugin.Name, enabled)
	return nil
}

// ToggleStatus 切换启用状态（带建表逻辑）
func (s *PluginService) ToggleStatus(id uint) error {
	var plugin model.Plugin
	if err := s.db.First(&plugin, id).Error; err != nil {
		return err
	}
	return s.SetEnabled(id, !plugin.IsEnabled)
}

// createPluginTables 根据插件名称创建对应的表
func (s *PluginService) createPluginTables(pluginName string) error {
	pluginTableMap := map[string][]interface{}{
		// 客服聊天系统
		"cs_chat": {
			&model.CSChatSession{},
			&model.CSChatMessage{},
			&model.CSChatConfig{},
			&model.CSChatQuickReply{},
		},
		// AI工单自动回复（含 Agent Function Calling）
		"ai_ticket": {
			&model.AITicketKnowledge{},
			&model.AITicketRule{},
			&model.AITicketQueue{},
			&model.AITicketProcessLog{},
			&model.AITicketNotifyLog{},
			&model.AITicketMode{},
			&model.AITicketConfig{},
			&model.AIToolConfig{},
			&model.AIToolExecutionLog{},
		},
		// AI购物助手
		"ai_shopping": {
			&model.AIShoppingConfig{},
			&model.AIShoppingChatLog{},
			&model.ProductCatalogTool{},
		},
		// 邮箱后缀白名单
		"email_suffix_whitelist": {
			&model.EmailSuffixWhitelist{},
		},
		// anchor_cloud_finance_pro 插件模块
		"acfp": {
			&model.ACFPFailNotifyEvent{},
			&model.ACFPUpstreamCache{},
			&model.ACFPIPHistory{},
			&model.ACFPLimitedSale{},
			&model.ACFPPriceLock{},
			&model.ACFPLog{},
			&model.ACFPCronStatus{},
			&model.ACFPCertProConfig{},
			&model.ACFPCertMinor{},
			&model.ACFPBatchTask{},
		},
	}

	models, ok := pluginTableMap[pluginName]
	if !ok {
		// 没有专属表的插件（如smtp、oauth等）直接跳过
		return nil
	}

	for _, m := range models {
		if err := s.db.AutoMigrate(m); err != nil {
			return fmt.Errorf("建表失败 %T: %w", m, err)
		}
	}
	return nil
}

// UpdateConfig 更新插件配置
func (s *PluginService) UpdateConfig(id uint, config string) error {
	result := s.db.Model(&model.Plugin{}).Where("id = ?", id).Update("config", config)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("插件不存在")
	}
	s.log.Infof("插件配置更新: id=%d", id)
	return nil
}

// GetEnabledByType 获取指定类型的启用插件
func (s *PluginService) GetEnabledByType(pluginType string) (*model.Plugin, error) {
	var plugin model.Plugin
	if err := s.db.Where("type = ? AND is_enabled = true", pluginType).
		Order("sort_order ASC").
		First(&plugin).Error; err != nil {
		return nil, fmt.Errorf("没有启用的%s插件", pluginType)
	}
	return &plugin, nil
}

// GetEnabledList 获取指定类型的启用插件列表
func (s *PluginService) GetEnabledList(pluginType string) ([]model.Plugin, error) {
	var plugins []model.Plugin
	if err := s.db.Where("type = ? AND is_enabled = true", pluginType).
		Order("sort_order ASC").
		Find(&plugins).Error; err != nil {
		return nil, err
	}
	return plugins, nil
}

// InitDefaults 初始化默认插件
func (s *PluginService) InitDefaults() {
	var count int64
	s.db.Model(&model.Plugin{}).Count(&count)
	if count > 0 {
		return
	}

	defaults := []model.Plugin{
		// 邮件插件
		{Name: "Smtp", Title: "SMTP邮件", Type: "mail", Description: "SMTP邮件发送", Author: "锚点财务", Version: "1.0", IsSystem: true, IsEnabled: false, SortOrder: 1},
		{Name: "Subemail", Title: "赛邮邮件", Type: "mail", Description: "赛邮邮件服务", Author: "锚点财务", Version: "1.0", IsSystem: true, IsEnabled: false, SortOrder: 2},
		{Name: "Alimail", Title: "阿里云邮件", Type: "mail", Description: "阿里云邮件推送", Author: "锚点财务", Version: "1.0", IsSystem: true, IsEnabled: false, SortOrder: 3},

		// 短信插件
		{Name: "Submail", Title: "赛邮短信", Type: "sms", Description: "赛邮短信服务", Author: "锚点财务", Version: "1.0", IsSystem: true, IsEnabled: false, SortOrder: 1},
		{Name: "Smsbao", Title: "短信宝", Type: "sms", Description: "短信宝短信服务", Author: "锚点财务", Version: "1.0", IsSystem: true, IsEnabled: false, SortOrder: 2},

		// 实名认证插件
		{Name: "Wechat", Title: "微信实名认证", Type: "certification", Description: "腾讯云人脸核身", Author: "锚点财务", Version: "1.0", IsSystem: true, IsEnabled: false, SortOrder: 1},
		{Name: "Idcsmartali", Title: "阿里云实名认证", Type: "certification", Description: "阿里云实人认证", Author: "锚点财务", Version: "1.0", IsSystem: true, IsEnabled: false, SortOrder: 2},

		// OAuth登录插件
		{Name: "Weixin", Title: "微信登录", Type: "oauth", Description: "微信扫码登录", Author: "锚点财务", Version: "1.0", IsSystem: true, IsEnabled: false, SortOrder: 1},
		{Name: "QQ", Title: "QQ登录", Type: "oauth", Description: "QQ互联登录", Author: "锚点财务", Version: "1.0", IsSystem: true, IsEnabled: false, SortOrder: 2},
		{Name: "Alipay", Title: "支付宝登录", Type: "oauth", Description: "支付宝登录", Author: "锚点财务", Version: "1.0", IsSystem: true, IsEnabled: false, SortOrder: 3},
		{Name: "Weibo", Title: "微博登录", Type: "oauth", Description: "微博登录", Author: "锚点财务", Version: "1.0", IsSystem: true, IsEnabled: false, SortOrder: 4},

		// 服务器模块（自动开通）
		{Name: "ProxmoxVE", Title: "ProxmoxVE", Type: "server", Description: "ProxmoxVE虚拟化", Author: "锚点财务", Version: "1.0", IsSystem: true, IsEnabled: false, SortOrder: 1},
		{Name: "NoKVM", Title: "NoKVM", Type: "server", Description: "NoKVM虚拟化", Author: "锚点财务", Version: "1.0", IsSystem: true, IsEnabled: false, SortOrder: 2},
		{Name: "BtHosts", Title: "宝塔主机", Type: "server", Description: "宝塔主机管理", Author: "锚点财务", Version: "1.0", IsSystem: true, IsEnabled: false, SortOrder: 3},
		{Name: "WlKanglePro", Title: "WlKanglePro", Type: "server", Description: "Kangle主机", Author: "锚点财务", Version: "1.0", IsSystem: true, IsEnabled: false, SortOrder: 4},

		// 扩展插件
		{Name: "ExportExcel", Title: "导出Excel", Type: "addon", Description: "数据导出为Excel", Author: "锚点财务", Version: "1.0", IsSystem: true, IsEnabled: false, SortOrder: 1},
		{Name: "ProductDivert", Title: "产品转移", Type: "addon", Description: "产品转移功能", Author: "锚点财务", Version: "1.0", IsSystem: true, IsEnabled: false, SortOrder: 2},
		{Name: "ExpiredIpLog", Title: "过期IP日志", Type: "addon", Description: "过期IP记录", Author: "锚点财务", Version: "1.0", IsSystem: true, IsEnabled: false, SortOrder: 3},
		{Name: "ClientCare", Title: "客户关怀", Type: "addon", Description: "客户关怀系统", Author: "锚点财务", Version: "1.0", IsSystem: true, IsEnabled: false, SortOrder: 4},
	}

	for _, d := range defaults {
		s.db.Create(&d)
	}

	s.log.Info("默认插件初始化完成")
}

// ==================== 服务器模块管理 ====================

// CreateServerModule 创建服务器模块
func (s *PluginService) CreateServerModule(req model.ServerModule) (*model.ServerModule, error) {
	var count int64
	s.db.Model(&model.ServerModule{}).Where("name = ?", req.Name).Count(&count)
	if count > 0 {
		return nil, errors.New("服务器模块标识名已存在")
	}

	if err := s.db.Create(&req).Error; err != nil {
		return nil, err
	}

	s.log.Infof("服务器模块创建成功: id=%d name=%s", req.ID, req.Name)
	return &req, nil
}

// UpdateServerModule 更新服务器模块
func (s *PluginService) UpdateServerModule(id uint, updates map[string]interface{}) error {
	result := s.db.Model(&model.ServerModule{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("服务器模块不存在")
	}
	return nil
}

// DeleteServerModule 删除服务器模块
func (s *PluginService) DeleteServerModule(id uint) error {
	result := s.db.Delete(&model.ServerModule{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("服务器模块不存在")
	}
	return nil
}

// GetServerModuleByID 根据ID获取服务器模块
func (s *PluginService) GetServerModuleByID(id uint) (*model.ServerModule, error) {
	var module model.ServerModule
	if err := s.db.First(&module, id).Error; err != nil {
		return nil, err
	}
	return &module, nil
}

// GetServerModuleList 获取服务器模块列表
func (s *PluginService) GetServerModuleList(moduleType string) ([]model.ServerModule, error) {
	var modules []model.ServerModule
	query := s.db.Model(&model.ServerModule{})

	if moduleType != "" {
		query = query.Where("module = ?", moduleType)
	}

	if err := query.Order("sort_order ASC, id ASC").Find(&modules).Error; err != nil {
		return nil, err
	}
	return modules, nil
}

// ==================== 服务器分组管理 ====================

// CreateServerGroup 创建服务器分组
func (s *PluginService) CreateServerGroup(req model.ServerGroup) (*model.ServerGroup, error) {
	var count int64
	s.db.Model(&model.ServerGroup{}).Where("name = ?", req.Name).Count(&count)
	if count > 0 {
		return nil, errors.New("分组名称已存在")
	}

	if err := s.db.Create(&req).Error; err != nil {
		return nil, err
	}

	return &req, nil
}

// UpdateServerGroup 更新服务器分组
func (s *PluginService) UpdateServerGroup(id uint, updates map[string]interface{}) error {
	result := s.db.Model(&model.ServerGroup{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("分组不存在")
	}
	return nil
}

// DeleteServerGroup 删除服务器分组
func (s *PluginService) DeleteServerGroup(id uint) error {
	result := s.db.Delete(&model.ServerGroup{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("分组不存在")
	}
	return nil
}

// GetServerGroupList 获取服务器分组列表
func (s *PluginService) GetServerGroupList() ([]model.ServerGroup, error) {
	var groups []model.ServerGroup
	if err := s.db.Order("sort_order ASC, id ASC").Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}
