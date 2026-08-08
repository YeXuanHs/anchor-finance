package model

import (
	"time"

	"gorm.io/gorm"
)

// Menu 菜单模型
type Menu struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	ParentID  uint           `gorm:"default:0;index" json:"parent_id"`
	URL       string         `gorm:"size:255" json:"url"`
	Icon      string         `gorm:"size:100" json:"icon"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	IsVisible bool           `gorm:"default:true" json:"is_visible"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	IsSystem  bool           `gorm:"default:false" json:"is_system"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
