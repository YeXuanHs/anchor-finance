package model

import (
	"time"

	"gorm.io/gorm"
)

// Plugin 插件模型
type Plugin struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"size:100;not null" json:"name"`
	Code          string         `gorm:"size:50;uniqueIndex;not null" json:"code"`
	Description   string         `gorm:"type:text" json:"description"`
	Version       string         `gorm:"size:20" json:"version"`
	Author        string         `gorm:"size:100" json:"author"`
	Icon          string         `gorm:"size:255" json:"icon"`
	Category      string         `gorm:"size:50" json:"category"` // oauth, payment, verification, other
	Enabled       bool           `gorm:"default:false" json:"enabled"`
	Config        string         `gorm:"type:json" json:"config"`
	ConfigFields  string         `gorm:"type:json" json:"config_fields"`
	HasUpdate     bool           `gorm:"default:false" json:"has_update"`
	InstalledAt   *time.Time     `json:"installed_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// PluginConfig 插件配置字段定义
type PluginConfig struct {
	Key         string      `json:"key"`
	Label       string      `json:"label"`
	Type        string      `json:"type"` // text, password, switch, select, textarea
	Placeholder string      `json:"placeholder"`
	Default     interface{} `json:"default"`
	Options     []Option    `json:"options,omitempty"`
	Required    bool        `json:"required"`
}

// Option 下拉选项
type Option struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// PluginQueryParams 插件查询参数
type PluginQueryParams struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Category string `json:"category"`
	Keyword  string `json:"keyword"`
}

// TogglePluginRequest 切换插件状态请求
type TogglePluginRequest struct {
	Enabled bool `json:"enabled"`
}

// UpdatePluginConfigRequest 更新插件配置请求
type UpdatePluginConfigRequest struct {
	Config map[string]interface{} `json:"config" binding:"required"`
}

// PluginListResponse 插件列表响应
type PluginListResponse struct {
	List  []Plugin `json:"list"`
	Total int64    `json:"total"`
}
