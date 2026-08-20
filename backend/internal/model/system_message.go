package model

import (
	"time"

	"gorm.io/gorm"
)

// SystemMessage 站内消息
type SystemMessage struct {
	gorm.Model
	UserID  uint       `gorm:"index;not null" json:"user_id"`
	Title   string     `gorm:"type:varchar(256);not null" json:"title"`
	Content string     `gorm:"type:text" json:"content"`
	Type    string     `gorm:"type:varchar(16);default:system" json:"type"` // system/order/ticket/promotion
	IsRead  bool       `gorm:"default:false;index" json:"is_read"`
	ReadAt  *time.Time `json:"read_at"`
	Link    string     `gorm:"type:varchar(512)" json:"link"`
}
