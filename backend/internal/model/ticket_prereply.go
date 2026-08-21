package model

import (
	"time"

	"gorm.io/gorm"
)

// TicketPrereply 工单预回复模型
type TicketPrereply struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	CategoryID uint           `gorm:"index" json:"category_id"`
	Title      string         `gorm:"size:200;not null" json:"title"`
	Content    string         `gorm:"type:text" json:"content"`
	SortOrder  int            `gorm:"default:0" json:"sort_order"`
	Status     string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (TicketPrereply) TableName() string {
	return "ticket_prereplies"
}

// TicketPrereplyCategory 工单预回复分类模型
type TicketPrereplyCategory struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	Status    string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (TicketPrereplyCategory) TableName() string {
	return "ticket_prereply_categories"
}
