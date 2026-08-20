package model

import (
	"gorm.io/gorm"
)

// HomepageFeature 首页特色区域
type HomepageFeature struct {
	gorm.Model
	Title       string `gorm:"type:varchar(128);not null" json:"title"`
	Description string `gorm:"type:varchar(512)" json:"description"`
	Icon        string `gorm:"type:varchar(128)" json:"icon"`      // 图标名称或URL
	LinkURL     string `gorm:"type:varchar(512)" json:"link_url"`  // 点击跳转链接
	SortOrder   int    `gorm:"default:0;index" json:"sort_order"`
	Status      int16  `gorm:"type:smallint;default:1;not null" json:"status"` // 1=启用 0=禁用
	Position    string `gorm:"type:varchar(32);default:'home';index" json:"position"` // home/product
}
