package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Plugin 插件
type Plugin struct {
	gorm.Model
	Name        string         `gorm:"type:varchar(128);not null;index" json:"name"`
	Slug        string         `gorm:"type:varchar(128);uniqueIndex" json:"slug"`
	Description string         `gorm:"type:text" json:"description"`
	Author      string         `gorm:"type:varchar(128)" json:"author"`
	Website     string         `gorm:"type:varchar(512)" json:"website"`
	Version     string         `gorm:"type:varchar(32);not null" json:"version"`
	License     string         `gorm:"type:varchar(64)" json:"license"`
	Config      datatypes.JSON `gorm:"type:jsonb" json:"config"`
	Enabled     bool           `gorm:"default:false;index" json:"enabled"`
	IsCore      bool           `gorm:"default:false" json:"is_core"`
	Hooks       datatypes.JSON `gorm:"type:jsonb" json:"hooks"`
	Dependencies datatypes.JSON `gorm:"type:jsonb" json:"dependencies"`
	SortOrder   int            `gorm:"default:0;index" json:"sort_order"`
	Status      int16          `gorm:"type:smallint;default:1;not null;index" json:"status"` // 1=正常 0=异常
	Path        string         `gorm:"type:varchar(512)" json:"path"`
	Metadata    datatypes.JSON `gorm:"type:jsonb" json:"metadata"`
	InstalledAt time.Time      `json:"installed_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// PluginLog 插件日志
type PluginLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PluginID  uint      `gorm:"index;not null" json:"plugin_id"`
	Plugin    Plugin    `gorm:"foreignKey:PluginID" json:"plugin,omitempty"`
	Action    string    `gorm:"type:varchar(32);not null" json:"action"` // install/uninstall/enable/disable/config/error
	Detail    string    `gorm:"type:text" json:"detail"`
	AdminID   uint      `gorm:"index" json:"admin_id"`
	CreatedAt time.Time `json:"created_at"`
}
