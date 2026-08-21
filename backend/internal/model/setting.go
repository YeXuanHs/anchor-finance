package model

import (
	"time"

	"gorm.io/gorm"
)

// Setting 系统设置模型
type Setting struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Group     string         `gorm:"size:50;index" json:"group"` // general, finance, security, etc.
	Key       string         `gorm:"size:100;uniqueIndex" json:"key"`
	Value     string         `gorm:"type:text" json:"value"`
	Type      string         `gorm:"size:20;default:string" json:"type"` // string, int, bool, json
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Setting) TableName() string {
	return "settings"
}

// Menu 菜单模型
type Menu struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ParentID  uint           `gorm:"index;default:0" json:"parent_id"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	Path      string         `gorm:"size:100" json:"path"`
	Icon      string         `gorm:"size:50" json:"icon"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	IsVisible bool           `gorm:"default:true" json:"is_visible"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Menu) TableName() string {
	return "menus"
}
