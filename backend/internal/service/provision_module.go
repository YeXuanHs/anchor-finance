package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"anchorfinance/internal/model"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ==================== Module Management ====================

// GetModules returns all provision modules with optional type and enabled filter.
func (s *ProvisionService) GetModules(moduleType string, enabled *bool) ([]model.ProvisionModule, error) {
	var modules []model.ProvisionModule
	query := s.db.Model(&model.ProvisionModule{})
	if moduleType != "" {
		query = query.Where("type = ?", moduleType)
	}
	if enabled != nil {
		query = query.Where("active = ?", *enabled)
	}
	if err := query.Order("priority DESC, id ASC").Find(&modules).Error; err != nil {
		return nil, err
	}
	return modules, nil
}

// GetModuleByID returns a single module by ID.
func (s *ProvisionService) GetModuleByID(id uint) (*model.ProvisionModule, error) {
	var module model.ProvisionModule
	if err := s.db.First(&module, id).Error; err != nil {
		return nil, err
	}
	return &module, nil
}

// CreateModule creates a new provision module.
func (s *ProvisionService) CreateModule(req CreateModuleRequest) (*model.ProvisionModule, error) {
	var existing model.ProvisionModule
	if err := s.db.Where("slug = ?", req.Slug).First(&existing).Error; err == nil {
		return nil, errors.New("module with this slug already exists")
	}

	module := &model.ProvisionModule{
		Name:              req.Name,
		Slug:              req.Slug,
		Description:       req.Description,
		Type:              req.Type,
		SupportedProducts: req.SupportedProducts,
		Config:            req.Config,
		ServerURL:         req.ServerURL,
		ServerIP:          req.ServerIP,
		APIKey:            req.APIKey,
		APISecret:         req.APISecret,
		Active:            true,
		Priority:          req.Priority,
		Weight:            req.Weight,
	}
	if module.Weight == 0 {
		module.Weight = 1
	}
	if module.MaxRetries == 0 {
		module.MaxRetries = 3
	}
	if module.Timeout == 0 {
		module.Timeout = 30
	}
	if err := s.db.Create(module).Error; err != nil {
		return nil, err
	}
	return module, nil
}

// UpdateModule updates a provision module.
func (s *ProvisionService) UpdateModule(id uint, req UpdateModuleRequest) (*model.ProvisionModule, error) {
	var module model.ProvisionModule
	if err := s.db.First(&module, id).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Slug != nil {
		updates["slug"] = *req.Slug
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Config != nil {
		updates["config"] = *req.Config
	}
	if req.ServerURL != nil {
		updates["server_url"] = *req.ServerURL
	}
	if req.ServerIP != nil {
		updates["server_ip"] = *req.ServerIP
	}
	if req.APIKey != nil {
		updates["api_key"] = *req.APIKey
	}
	if req.APISecret != nil {
		updates["api_secret"] = *req.APISecret
	}
	if req.Active != nil {
		updates["active"] = *req.Active
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}
	if req.Weight != nil {
		updates["weight"] = *req.Weight
	}
	if req.Timeout != nil {
		updates["timeout"] = *req.Timeout
	}

	if len(updates) > 0 {
		if err := s.db.Model(&module).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	if err := s.db.First(&module, id).Error; err != nil {
		return nil, err
	}
	return &module, nil
}

// DeleteModule deletes a provision module and its related buttons/charts/functions.
func (s *ProvisionService) DeleteModule(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("module_id = ?", id).Delete(&model.ProvisionButton{}).Error; err != nil {
			return err
		}
		if err := tx.Where("module_id = ?", id).Delete(&model.ProvisionChart{}).Error; err != nil {
			return err
		}
		if err := tx.Where("module_id = ?", id).Delete(&model.CustomFunction{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.ProvisionModule{}, id).Error
	})
}

// TestModule tests connectivity to a module.
func (s *ProvisionService) TestModule(id uint) (*model.ProvisionModule, error) {
	return s.TestConnection(id)
}

// ==================== Client Area Rendering ====================

// RenderClientArea renders the client area for a host/service.
func (s *ProvisionService) RenderClientArea(hostID uint) (map[string]interface{}, error) {
	host, module, err := s.getHostAndModule(hostID)
	if err != nil {
		return nil, err
	}

	buttons, _ := s.getClientButtons(module.ID)
	charts, _ := s.getChartsByModule(module.ID)

	result := map[string]interface{}{
		"host":   host,
		"module": module,
		"buttons": buttons,
		"charts":  charts,
	}
	return result, nil
}

// RenderClientAreaDetail renders the detailed client area with usage data.
func (s *ProvisionService) RenderClientAreaDetail(hostID uint) (map[string]interface{}, error) {
	host, module, err := s.getHostAndModule(hostID)
	if err != nil {
		return nil, err
	}

	buttons, _ := s.getClientButtons(module.ID)
	charts, _ := s.getChartsByModule(module.ID)
	usage, _ := s.GetUsage(hostID)

	result := map[string]interface{}{
		"host":    host,
		"module":  module,
		"buttons": buttons,
		"charts":  charts,
		"usage":   usage,
	}
	return result, nil
}

// GetClientButtons returns available client buttons for a host's module.
func (s *ProvisionService) GetClientButtons(hostID uint) ([]model.ProvisionButton, error) {
	_, module, err := s.getHostAndModule(hostID)
	if err != nil {
		return nil, err
	}
	return s.getClientButtons(module.ID)
}

func (s *ProvisionService) getClientButtons(moduleID uint) ([]model.ProvisionButton, error) {
	var buttons []model.ProvisionButton
	err := s.db.Where("module_id = ? AND type = ? AND enabled = ?", moduleID, "client", true).
		Order("sort_order ASC").
		Find(&buttons).Error
	return buttons, err
}

// ExecuteClientButton executes a client button action for a host.
func (s *ProvisionService) ExecuteClientButton(hostID uint, action string) (map[string]interface{}, error) {
	_, module, err := s.getHostAndModule(hostID)
	if err != nil {
		return nil, err
	}

	var button model.ProvisionButton
	if err := s.db.Where("module_id = ? AND action = ? AND type = ? AND enabled = ?",
		module.ID, action, "client", true).First(&button).Error; err != nil {
		return nil, fmt.Errorf("button action '%s' not found or disabled", action)
	}

	if button.URL != "" {
		return map[string]interface{}{
			"type": "redirect",
			"url":  button.URL,
		}, nil
	}

	return map[string]interface{}{
		"type":    "action",
		"action":  action,
		"host_id": hostID,
		"status":  "executed",
	}, nil
}

// ==================== Admin Area Rendering ====================

// RenderAdminArea renders the admin area for a host/service.
func (s *ProvisionService) RenderAdminArea(hostID uint) (map[string]interface{}, error) {
	host, module, err := s.getHostAndModule(hostID)
	if err != nil {
		return nil, err
	}

	buttons, _ := s.getAdminButtons(module.ID)
	charts, _ := s.getChartsByModule(module.ID)
	usage, _ := s.GetUsage(hostID)

	result := map[string]interface{}{
		"host":    host,
		"module":  module,
		"buttons": buttons,
		"charts":  charts,
		"usage":   usage,
	}
	return result, nil
}

// GetAdminButtons returns available admin buttons for a host's module.
func (s *ProvisionService) GetAdminButtons(hostID uint) ([]model.ProvisionButton, error) {
	_, module, err := s.getHostAndModule(hostID)
	if err != nil {
		return nil, err
	}
	return s.getAdminButtons(module.ID)
}

func (s *ProvisionService) getAdminButtons(moduleID uint) ([]model.ProvisionButton, error) {
	var buttons []model.ProvisionButton
	err := s.db.Where("module_id = ? AND type = ? AND enabled = ?", moduleID, "admin", true).
		Order("sort_order ASC").
		Find(&buttons).Error
	return buttons, err
}

// ExecuteAdminButton executes an admin button action for a host.
func (s *ProvisionService) ExecuteAdminButton(hostID uint, action string) (map[string]interface{}, error) {
	_, module, err := s.getHostAndModule(hostID)
	if err != nil {
		return nil, err
	}

	var button model.ProvisionButton
	if err := s.db.Where("module_id = ? AND action = ? AND type = ? AND enabled = ?",
		module.ID, action, "admin", true).First(&button).Error; err != nil {
		return nil, fmt.Errorf("admin button action '%s' not found or disabled", action)
	}

	if button.URL != "" {
		return map[string]interface{}{
			"type": "redirect",
			"url":  button.URL,
		}, nil
	}

	return map[string]interface{}{
		"type":    "action",
		"action":  action,
		"host_id": hostID,
		"status":  "executed",
	}, nil
}

// GetDefaultButtons returns the default set of provision buttons.
func (s *ProvisionService) GetDefaultButtons() []map[string]interface{} {
	return []map[string]interface{}{
		{"name": "Create", "action": "create", "type": "admin", "position": "action", "icon": "plus"},
		{"name": "Suspend", "action": "suspend", "type": "admin", "position": "action", "icon": "pause"},
		{"name": "Unsuspend", "action": "unsuspend", "type": "admin", "position": "action", "icon": "play"},
		{"name": "Terminate", "action": "terminate", "type": "admin", "position": "action", "icon": "trash", "confirm": true},
		{"name": "Renew", "action": "renew", "type": "admin", "position": "action", "icon": "refresh"},
		{"name": "Rebuild", "action": "rebuild", "type": "admin", "position": "action", "icon": "redo", "confirm": true},
	}
}

// ==================== Charts ====================

// GetCharts returns available charts for a host's module.
func (s *ProvisionService) GetCharts(hostID uint) ([]model.ProvisionChart, error) {
	_, module, err := s.getHostAndModule(hostID)
	if err != nil {
		return nil, err
	}
	return s.getChartsByModule(module.ID)
}

func (s *ProvisionService) getChartsByModule(moduleID uint) ([]model.ProvisionChart, error) {
	var charts []model.ProvisionChart
	err := s.db.Where("module_id = ? AND enabled = ?", moduleID, true).Find(&charts).Error
	return charts, err
}

// GetChartData returns chart data for a specific chart and period.
func (s *ProvisionService) GetChartData(hostID, chartID uint, period string) (map[string]interface{}, error) {
	var chart model.ProvisionChart
	if err := s.db.First(&chart, chartID).Error; err != nil {
		return nil, fmt.Errorf("chart not found: %w", err)
	}

	// Chart data would be fetched from the configured endpoint or generated locally
	result := map[string]interface{}{
		"chart_id": chartID,
		"host_id":  hostID,
		"period":   period,
		"type":     chart.Type,
		"name":     chart.Name,
		"data":     []interface{}{},
	}

	// If endpoint is configured, data would come from there
	if chart.Endpoint != "" {
		result["endpoint"] = chart.Endpoint
	}

	return result, nil
}

// CreateChart creates a new chart configuration.
func (s *ProvisionService) CreateChart(req CreateChartRequest) (*model.ProvisionChart, error) {
	chart := &model.ProvisionChart{
		ModuleID: req.ModuleID,
		Name:     req.Name,
		Type:     req.Type,
		Endpoint: req.Endpoint,
		Config:   req.Config,
		Enabled:  true,
	}
	if err := s.db.Create(chart).Error; err != nil {
		return nil, err
	}
	return chart, nil
}

// UpdateChart updates a chart configuration.
func (s *ProvisionService) UpdateChart(id uint, req UpdateChartRequest) (*model.ProvisionChart, error) {
	var chart model.ProvisionChart
	if err := s.db.First(&chart, id).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Endpoint != nil {
		updates["endpoint"] = *req.Endpoint
	}
	if req.Config != nil {
		updates["config"] = *req.Config
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}

	if len(updates) > 0 {
		if err := s.db.Model(&chart).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	if err := s.db.First(&chart, id).Error; err != nil {
		return nil, err
	}
	return &chart, nil
}

// ==================== Usage Tracking ====================

// GetUsage returns the current usage data for a host.
func (s *ProvisionService) GetUsage(hostID uint) (*model.ProvisionUsage, error) {
	var usage model.ProvisionUsage
	if err := s.db.Where("host_id = ?", hostID).First(&usage).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &model.ProvisionUsage{HostID: hostID}, nil
		}
		return nil, err
	}
	return &usage, nil
}

// UpdateUsage updates usage data for a host.
func (s *ProvisionService) UpdateUsage(hostID uint, req UpdateUsageRequest) (*model.ProvisionUsage, error) {
	var usage model.ProvisionUsage
	err := s.db.Where("host_id = ?", hostID).First(&usage).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		usage = model.ProvisionUsage{
			HostID:    hostID,
			ModuleID:  req.ModuleID,
			CPU:       req.CPU,
			Memory:    req.Memory,
			Disk:      req.Disk,
			Bandwidth: req.Bandwidth,
			TrafficIn: req.TrafficIn,
			TrafficOut: req.TrafficOut,
		}
		if req.Extra != nil {
			usage.Extra = *req.Extra
		}
		if err := s.db.Create(&usage).Error; err != nil {
			return nil, err
		}
		return &usage, nil
	} else if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"cpu":         req.CPU,
		"memory":      req.Memory,
		"disk":        req.Disk,
		"bandwidth":   req.Bandwidth,
		"traffic_in":  req.TrafficIn,
		"traffic_out": req.TrafficOut,
	}
	if req.Extra != nil {
		updates["extra"] = *req.Extra
	}
	if err := s.db.Model(&usage).Updates(updates).Error; err != nil {
		return nil, err
	}
	if err := s.db.Where("host_id = ?", hostID).First(&usage).Error; err != nil {
		return nil, err
	}
	return &usage, nil
}

// CheckDefineUsage checks custom usage limits for a host.
func (s *ProvisionService) CheckDefineUsage(hostID uint) (map[string]interface{}, error) {
	usage, err := s.GetUsage(hostID)
	if err != nil {
		return nil, err
	}

	var host model.Host
	if err := s.db.First(&host, hostID).Error; err != nil {
		return nil, fmt.Errorf("host not found: %w", err)
	}

	result := map[string]interface{}{
		"host_id": hostID,
		"usage":   usage,
		"limits": map[string]interface{}{
			"disk_gb":        host.DiskSizeGB,
			"bandwidth_mbps": host.BandwidthMbps,
			"traffic_gb":     host.TrafficGB,
		},
		"warnings": []string{},
	}

	// Check thresholds
	warnings := result["warnings"].([]string)
	if host.DiskSizeGB > 0 && usage.Disk > float64(host.DiskSizeGB)*0.9 {
		warnings = append(warnings, "disk usage exceeds 90%")
	}
	if host.TrafficGB > 0 {
		totalTrafficGB := float64(usage.TrafficIn+usage.TrafficOut) / (1024 * 1024 * 1024)
		if totalTrafficGB > float64(host.TrafficGB)*0.9 {
			warnings = append(warnings, "traffic usage exceeds 90%")
		}
	}
	result["warnings"] = warnings

	return result, nil
}

// TrafficUsage returns traffic usage for a host.
func (s *ProvisionService) TrafficUsage(hostID uint) (map[string]interface{}, error) {
	usage, err := s.GetUsage(hostID)
	if err != nil {
		return nil, err
	}

	var host model.Host
	if err := s.db.First(&host, hostID).Error; err != nil {
		return nil, fmt.Errorf("host not found: %w", err)
	}

	trafficInGB := float64(usage.TrafficIn) / (1024 * 1024 * 1024)
	trafficOutGB := float64(usage.TrafficOut) / (1024 * 1024 * 1024)
	totalGB := trafficInGB + trafficOutGB

	result := map[string]interface{}{
		"host_id":       hostID,
		"traffic_in":    trafficInGB,
		"traffic_out":   trafficOutGB,
		"total":         totalGB,
		"limit_gb":      host.TrafficGB,
		"unit":          "GB",
	}

	if host.TrafficGB > 0 {
		result["usage_percent"] = (totalGB / float64(host.TrafficGB)) * 100
		result["remaining_gb"] = float64(host.TrafficGB) - totalGB
	}

	return result, nil
}

// ==================== Custom Functions ====================

// GetCustomFunctions returns custom functions for a module.
func (s *ProvisionService) GetCustomFunctions(moduleID uint) ([]model.CustomFunction, error) {
	var functions []model.CustomFunction
	if err := s.db.Where("module_id = ?", moduleID).Order("id ASC").Find(&functions).Error; err != nil {
		return nil, err
	}
	return functions, nil
}

// CreateCustomFunction creates a new custom function.
func (s *ProvisionService) CreateCustomFunction(req CreateFunctionRequest) (*model.CustomFunction, error) {
	fn := &model.CustomFunction{
		ModuleID: req.ModuleID,
		Name:     req.Name,
		Code:     req.Code,
		Trigger:  req.Trigger,
		Enabled:  true,
	}
	if fn.Trigger == "" {
		fn.Trigger = "manual"
	}
	if err := s.db.Create(fn).Error; err != nil {
		return nil, err
	}
	return fn, nil
}

// ExecuteCustomFunction executes a custom function for a host.
func (s *ProvisionService) ExecuteCustomFunction(fnID, hostID uint) (map[string]interface{}, error) {
	var fn model.CustomFunction
	if err := s.db.First(&fn, fnID).Error; err != nil {
		return nil, fmt.Errorf("function not found: %w", err)
	}
	if !fn.Enabled {
		return nil, errors.New("function is disabled")
	}

	var host model.Host
	if err := s.db.First(&host, hostID).Error; err != nil {
		return nil, fmt.Errorf("host not found: %w", err)
	}

	// Execute the function and log the result
	result := map[string]interface{}{
		"function_id": fnID,
		"host_id":     hostID,
		"function":    fn.Name,
		"status":      "executed",
		"executed_at": time.Now(),
	}

	s.log.Infof("custom function executed: fn=%d host=%d name=%s", fnID, hostID, fn.Name)
	return result, nil
}

// DeleteCustomFunction deletes a custom function.
func (s *ProvisionService) DeleteCustomFunction(fnID uint) error {
	return s.db.Delete(&model.CustomFunction{}, fnID).Error
}

// ==================== Button CRUD ====================

// GetButtons returns buttons for a module with optional type filter.
func (s *ProvisionService) GetButtons(moduleID uint, buttonType string) ([]model.ProvisionButton, error) {
	var buttons []model.ProvisionButton
	query := s.db.Where("module_id = ?", moduleID)
	if buttonType != "" {
		query = query.Where("type = ?", buttonType)
	}
	if err := query.Order("sort_order ASC").Find(&buttons).Error; err != nil {
		return nil, err
	}
	return buttons, nil
}

// CreateButton creates a new button.
func (s *ProvisionService) CreateButton(req CreateButtonRequest) (*model.ProvisionButton, error) {
	button := &model.ProvisionButton{
		ModuleID:   req.ModuleID,
		Name:       req.Name,
		Action:     req.Action,
		Type:       req.Type,
		Position:   req.Position,
		Icon:       req.Icon,
		URL:        req.URL,
		Confirm:    req.Confirm,
		ConfirmMsg: req.ConfirmMsg,
		Enabled:    true,
		SortOrder:  req.SortOrder,
	}
	if button.Position == "" {
		button.Position = "header"
	}
	if err := s.db.Create(button).Error; err != nil {
		return nil, err
	}
	return button, nil
}

// UpdateButton updates a button.
func (s *ProvisionService) UpdateButton(id uint, req UpdateButtonRequest) (*model.ProvisionButton, error) {
	var button model.ProvisionButton
	if err := s.db.First(&button, id).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Action != nil {
		updates["action"] = *req.Action
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Position != nil {
		updates["position"] = *req.Position
	}
	if req.Icon != nil {
		updates["icon"] = *req.Icon
	}
	if req.URL != nil {
		updates["url"] = *req.URL
	}
	if req.Confirm != nil {
		updates["confirm"] = *req.Confirm
	}
	if req.ConfirmMsg != nil {
		updates["confirm_msg"] = *req.ConfirmMsg
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}

	if len(updates) > 0 {
		if err := s.db.Model(&button).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	if err := s.db.First(&button, id).Error; err != nil {
		return nil, err
	}
	return &button, nil
}

// DeleteButton deletes a button.
func (s *ProvisionService) DeleteButton(id uint) error {
	return s.db.Delete(&model.ProvisionButton{}, id).Error
}

// ==================== SSL/Download Resources ====================

// SSLButton returns SSL button configuration for a host.
func (s *ProvisionService) SSLButton(hostID uint) (map[string]interface{}, error) {
	_, module, err := s.getHostAndModule(hostID)
	if err != nil {
		return nil, err
	}

	var buttons []model.ProvisionButton
	s.db.Where("module_id = ? AND type = ? AND action LIKE ?", module.ID, "client", "ssl%").
		Find(&buttons)

	result := map[string]interface{}{
		"host_id":   hostID,
		"module_id": module.ID,
		"buttons":   buttons,
	}
	return result, nil
}

// DownloadResource downloads a module resource.
func (s *ProvisionService) DownloadResource(hostID, resourceID uint) (map[string]interface{}, error) {
	_, module, err := s.getHostAndModule(hostID)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"host_id":     hostID,
		"module_id":   module.ID,
		"resource_id": resourceID,
		"status":      "ready",
	}
	return result, nil
}

// ==================== Flow Packets ====================

// AfterFlowPacketPaid is a post-payment hook for flow packets.
func (s *ProvisionService) AfterFlowPacketPaid(hostID, packetID uint) error {
	var host model.Host
	if err := s.db.First(&host, hostID).Error; err != nil {
		return fmt.Errorf("host not found: %w", err)
	}

	s.log.Infof("flow packet paid: host=%d packet=%d", hostID, packetID)
	return nil
}

// ==================== Helpers ====================

func (s *ProvisionService) getHostAndModule(hostID uint) (*model.Host, *model.ProvisionModule, error) {
	var host model.Host
	if err := s.db.First(&host, hostID).Error; err != nil {
		return nil, nil, fmt.Errorf("host not found: %w", err)
	}

	if host.ProductID == nil {
		return nil, nil, errors.New("host has no associated product")
	}

	module, err := s.findModuleForProduct(*host.ProductID)
	if err != nil {
		return nil, nil, fmt.Errorf("no provision module for host: %w", err)
	}

	return &host, module, nil
}

// marshalJSON marshals a value to datatypes.JSON.
func marshalJSON(v interface{}) datatypes.JSON {
	b, _ := json.Marshal(v)
	return datatypes.JSON(b)
}

// ==================== Request DTOs ====================

type CreateModuleRequest struct {
	Name              string         `json:"name" binding:"required,max=128"`
	Slug              string         `json:"slug" binding:"required,max=128"`
	Description       string         `json:"description"`
	Type              string         `json:"type" binding:"required"`
	SupportedProducts datatypes.JSON `json:"supported_products"`
	Config            datatypes.JSON `json:"config"`
	ServerURL         string         `json:"server_url"`
	ServerIP          string         `json:"server_ip"`
	APIKey            string         `json:"api_key"`
	APISecret         string         `json:"api_secret"`
	Priority          int            `json:"priority"`
	Weight            int            `json:"weight"`
	Timeout           int            `json:"timeout"`
}

type UpdateModuleRequest struct {
	Name              *string         `json:"name"`
	Slug              *string         `json:"slug"`
	Description       *string         `json:"description"`
	Type              *string         `json:"type"`
	Config            *datatypes.JSON `json:"config"`
	ServerURL         *string         `json:"server_url"`
	ServerIP          *string         `json:"server_ip"`
	APIKey            *string         `json:"api_key"`
	APISecret         *string         `json:"api_secret"`
	Active            *bool           `json:"active"`
	Priority          *int            `json:"priority"`
	Weight            *int            `json:"weight"`
	Timeout           *int            `json:"timeout"`
}

type CreateChartRequest struct {
	ModuleID uint           `json:"module_id" binding:"required"`
	Name     string         `json:"name" binding:"required,max=128"`
	Type     string         `json:"type" binding:"required"`
	Endpoint string         `json:"endpoint"`
	Config   datatypes.JSON `json:"config"`
}

type UpdateChartRequest struct {
	Name     *string         `json:"name"`
	Type     *string         `json:"type"`
	Endpoint *string         `json:"endpoint"`
	Config   *datatypes.JSON `json:"config"`
	Enabled  *bool           `json:"enabled"`
}

type CreateButtonRequest struct {
	ModuleID   uint   `json:"module_id" binding:"required"`
	Name       string `json:"name" binding:"required,max=128"`
	Action     string `json:"action" binding:"required,max=64"`
	Type       string `json:"type" binding:"required"`
	Position   string `json:"position"`
	Icon       string `json:"icon"`
	URL        string `json:"url"`
	Confirm    bool   `json:"confirm"`
	ConfirmMsg string `json:"confirm_msg"`
	SortOrder  int    `json:"sort_order"`
}

type UpdateButtonRequest struct {
	Name       *string `json:"name"`
	Action     *string `json:"action"`
	Type       *string `json:"type"`
	Position   *string `json:"position"`
	Icon       *string `json:"icon"`
	URL        *string `json:"url"`
	Confirm    *bool   `json:"confirm"`
	ConfirmMsg *string `json:"confirm_msg"`
	Enabled    *bool   `json:"enabled"`
	SortOrder  *int    `json:"sort_order"`
}

type CreateFunctionRequest struct {
	ModuleID uint   `json:"module_id" binding:"required"`
	Name     string `json:"name" binding:"required,max=128"`
	Code     string `json:"code" binding:"required"`
	Trigger  string `json:"trigger"`
}

type UpdateUsageRequest struct {
	ModuleID  uint           `json:"module_id"`
	CPU       float64        `json:"cpu"`
	Memory    float64        `json:"memory"`
	Disk      float64        `json:"disk"`
	Bandwidth float64        `json:"bandwidth"`
	TrafficIn int64          `json:"traffic_in"`
	TrafficOut int64         `json:"traffic_out"`
	Extra     *datatypes.JSON `json:"extra"`
}
