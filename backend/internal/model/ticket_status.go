package model

import (
	"gorm.io/gorm"
)

// TicketStatus 工单状态
type TicketStatus struct {
	gorm.Model
	Title      string `gorm:"type:varchar(64);not null" json:"title"`
	Code       string `gorm:"type:varchar(32);uniqueIndex;not null" json:"code"` // open/answered/customer_reply/closed 等
	Color      string `gorm:"type:varchar(16)" json:"color"` // 前端展示颜色
	ShowActive int8   `gorm:"type:smallint;default:1" json:"show_active"` // 1=在活跃列表显示 0=不显示
	IsDefault  bool   `gorm:"default:false" json:"is_default"`
	IsSystem   bool   `gorm:"default:false" json:"is_system"` // 系统状态不可删除
	Order      int    `gorm:"default:0" json:"order"`
	Status     int16  `gorm:"type:smallint;default:1;not null;index" json:"status"` // 1=启用 0=禁用
}
