package model

import (
	"time"

	"gorm.io/gorm"
)

// NavGroup 前台导航分组
type NavGroup struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Groupname    string         `gorm:"size:255;not null" json:"groupname"`   // 分组名称
	FaIcon       string         `gorm:"size:255" json:"fa_icon"`             // 图标
	Order        int            `gorm:"default:0" json:"order"`              // 排序
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	
	// 关联
	Products     []Product      `gorm:"many2many:nav_group_products;" json:"products,omitempty"`
	ProductCount int            `gorm:"-" json:"product_count"`
}

func (NavGroup) TableName() string {
	return "nav_groups"
}

// NavGroupProduct 导航分组产品关联
type NavGroupProduct struct {
	NavGroupID uint `gorm:"primaryKey"`
	ProductID  uint `gorm:"primaryKey"`
}

func (NavGroupProduct) TableName() string {
	return "nav_group_products"
}

// NavGroupUser 用户导航分组配置
type NavGroupUser struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	Uid       uint `gorm:"index;not null" json:"uid"`       // 用户ID
	GroupID   uint `gorm:"index" json:"group_id"`           // 分组ID
	IsShow    bool `gorm:"default:true" json:"is_show"`     // 是否显示
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (NavGroupUser) TableName() string {
	return "nav_group_users"
}
