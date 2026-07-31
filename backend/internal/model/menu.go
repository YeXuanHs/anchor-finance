package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Menu 菜单
type Menu struct {
	gorm.Model
	Name      string         `gorm:"type:varchar(64);not null" json:"name"`           // 菜单名称
	Icon      string         `gorm:"type:varchar(128)" json:"icon"`                   // 图标
	URL       string         `gorm:"type:varchar(256)" json:"url"`                    // 跳转链接
	ParentID  *uint          `gorm:"index" json:"parent_id"`                          // 父级ID
	Parent    *Menu          `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children  []Menu         `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	SortOrder int            `gorm:"default:0;index" json:"sort_order"`               // 排序
	Type      string         `gorm:"type:varchar(16);not null;index" json:"type"`     // top/side/bottom
	IsVisible bool           `gorm:"default:true" json:"is_visible"`                  // 是否显示
	Permission string        `gorm:"type:varchar(128)" json:"permission"`             // 所需权限标识
	Target    string         `gorm:"type:varchar(16);default:'_self'" json:"target"`  // _self/_blank
	Badge     string         `gorm:"type:varchar(32)" json:"badge"`                   // 角标文字
	BadgeType string         `gorm:"type:varchar(16)" json:"badge_type"`              // 角标类型: dot/number/text
	Extra     datatypes.JSON `gorm:"type:json" json:"extra"`                         // 扩展配置
	IsActive  bool           `gorm:"default:true;index" json:"is_active"`             // 是否启用
}
