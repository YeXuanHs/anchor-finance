package service

import (
	"errors"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type PluginService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewPluginService(db *gorm.DB, log *logger.Logger) *PluginService {
	return &PluginService{db: db, log: log}
}

// GetList returns all plugins with optional enabled filter.
func (s *PluginService) GetList(page, pageSize int, enabled *bool) ([]model.Plugin, int64, error) {
	var plugins []model.Plugin
	var total int64

	query := s.db.Model(&model.Plugin{})
	if enabled != nil {
		query = query.Where("enabled = ?", *enabled)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).
		Order("sort_order ASC, id ASC").
		Find(&plugins).Error; err != nil {
		return nil, 0, err
	}
	return plugins, total, nil
}

// GetByID returns a single plugin by ID.
func (s *PluginService) GetByID(id uint) (*model.Plugin, error) {
	var plugin model.Plugin
	if err := s.db.First(&plugin, id).Error; err != nil {
		return nil, err
	}
	return &plugin, nil
}

// Install installs a new plugin.
func (s *PluginService) Install(req InstallPluginRequest) (*model.Plugin, error) {
	var existing model.Plugin
	if err := s.db.Where("slug = ?", req.Slug).First(&existing).Error; err == nil {
		return nil, errors.New("plugin with this slug already exists")
	}

	plugin := &model.Plugin{
		Name:         req.Name,
		Slug:         req.Slug,
		Description:  req.Description,
		Author:       req.Author,
		Website:      req.Website,
		Version:      req.Version,
		License:      req.License,
		Config:       req.Config,
		Enabled:      false,
		Hooks:        req.Hooks,
		Dependencies: req.Dependencies,
		SortOrder:    req.SortOrder,
		Path:         req.Path,
		Metadata:     req.Metadata,
		InstalledAt:  time.Now(),
	}
	if plugin.Version == "" {
		plugin.Version = "1.0.0"
	}
	if err := s.db.Create(plugin).Error; err != nil {
		return nil, err
	}

	s.logPlugin(plugin.ID, 0, "install", "plugin installed")
	return plugin, nil
}

// Uninstall soft-deletes a plugin. Core plugins cannot be uninstalled.
func (s *PluginService) Uninstall(id uint) error {
	var plugin model.Plugin
	if err := s.db.First(&plugin, id).Error; err != nil {
		return err
	}
	if plugin.IsCore {
		return errors.New("core plugin cannot be uninstalled")
	}

	if plugin.Enabled {
		if err := s.db.Model(&plugin).Update("enabled", false).Error; err != nil {
			return err
		}
	}

	if err := s.db.Delete(&model.Plugin{}, id).Error; err != nil {
		return err
	}

	s.logPlugin(id, 0, "uninstall", "plugin uninstalled")
	return nil
}

// Enable enables a plugin.
func (s *PluginService) Enable(id uint, adminID uint) error {
	var plugin model.Plugin
	if err := s.db.First(&plugin, id).Error; err != nil {
		return err
	}
	if plugin.Enabled {
		return errors.New("plugin already enabled")
	}
	if plugin.Status != 1 {
		return errors.New("plugin status is abnormal")
	}

	if err := s.db.Model(&plugin).Update("enabled", true).Error; err != nil {
		return err
	}

	s.logPlugin(id, adminID, "enable", "plugin enabled")
	return nil
}

// Disable disables a plugin.
func (s *PluginService) Disable(id uint, adminID uint) error {
	var plugin model.Plugin
	if err := s.db.First(&plugin, id).Error; err != nil {
		return err
	}
	if !plugin.Enabled {
		return errors.New("plugin already disabled")
	}
	if plugin.IsCore {
		return errors.New("core plugin cannot be disabled")
	}

	if err := s.db.Model(&plugin).Update("enabled", false).Error; err != nil {
		return err
	}

	s.logPlugin(id, adminID, "disable", "plugin disabled")
	return nil
}

// UpdateConfig updates a plugin's configuration.
func (s *PluginService) UpdateConfig(id uint, req UpdatePluginConfigRequest) (*model.Plugin, error) {
	var plugin model.Plugin
	if err := s.db.First(&plugin, id).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Config != nil {
		updates["config"] = *req.Config
	}
	if req.Hooks != nil {
		updates["hooks"] = *req.Hooks
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

	s.logPlugin(id, 0, "config", "plugin config updated")
	return &plugin, nil
}

// GetLogs returns plugin operation logs with pagination.
func (s *PluginService) GetLogs(page, pageSize int, pluginID *uint) ([]model.PluginLog, int64, error) {
	var logs []model.PluginLog
	var total int64

	query := s.db.Model(&model.PluginLog{})
	if pluginID != nil {
		query = query.Where("plugin_id = ?", *pluginID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).
		Order("id DESC").
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

// logPlugin creates a plugin operation log entry.
func (s *PluginService) logPlugin(pluginID, adminID uint, action, detail string) {
	log := &model.PluginLog{
		PluginID: pluginID,
		Action:   action,
		Detail:   detail,
		AdminID:  adminID,
	}
	s.db.Create(log)
}

// ---------- Request DTOs ----------

type InstallPluginRequest struct {
	Name         string `json:"name" binding:"required,max=128"`
	Slug         string `json:"slug" binding:"required,max=128"`
	Description  string `json:"description"`
	Author       string `json:"author"`
	Website      string `json:"website"`
	Version      string `json:"version"`
	License      string `json:"license"`
	Config       string `json:"config"`
	Hooks        string `json:"hooks"`
	Dependencies string `json:"dependencies"`
	SortOrder    int    `json:"sort_order"`
	Path         string `json:"path"`
	Metadata     string `json:"metadata"`
}

type UpdatePluginConfigRequest struct {
	Config    *string `json:"config"`
	Hooks     *string `json:"hooks"`
	SortOrder *int    `json:"sort_order"`
}
