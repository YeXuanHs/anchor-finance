package model

import (
	"time"

	"gorm.io/gorm"
)

// TicketReply 工单回复模型
type TicketReply struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TicketID  uint           `gorm:"index;not null" json:"ticket_id"`
	UserID    uint           `gorm:"index" json:"user_id"` // 0表示管理员回复
	Content   string         `gorm:"type:text;not null" json:"content"`
	IsAdmin   bool           `gorm:"default:false" json:"is_admin"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (TicketReply) TableName() string {
	return "ticket_replies"
}
