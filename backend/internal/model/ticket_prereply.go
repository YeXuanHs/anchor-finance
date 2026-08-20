package model

import (
	"gorm.io/gorm"
)

// TicketPrereplyCategory 工单预设回复分类
type TicketPrereplyCategory struct {
	gorm.Model
	Name      string              `gorm:"type:varchar(128);not null" json:"name"`
	Status    int16               `gorm:"type:smallint;default:1;not null;index" json:"status"`
	SortOrder int                 `gorm:"default:0" json:"sort_order"`
	Replies   []TicketPrereply    `gorm:"foreignKey:CategoryID" json:"replies,omitempty"`
}

// TicketPrereply 工单预设回复
type TicketPrereply struct {
	gorm.Model
	CategoryID uint   `gorm:"index;not null" json:"category_id"`
	Category   TicketPrereplyCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Title      string `gorm:"type:varchar(256);not null" json:"title"`
	Content    string `gorm:"type:text;not null" json:"content"`
	Status     int16  `gorm:"type:smallint;default:1;not null;index" json:"status"`
	SortOrder  int    `gorm:"default:0" json:"sort_order"`
	UseCount   int    `gorm:"default:0" json:"use_count"`
}
