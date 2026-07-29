package model

import (
	"time"

	"gorm.io/gorm"
)

// SendMessage 发送消息记录
type SendMessage struct {
	gorm.Model
	UserID  uint       `gorm:"index" json:"user_id"`
	Channel string     `gorm:"type:varchar(16);not null" json:"channel"` // email/sms/wechat
	To      string     `gorm:"type:varchar(128);not null" json:"to"`
	Subject string     `gorm:"type:varchar(256)" json:"subject"`
	Content string     `gorm:"type:text;not null" json:"content"`
	Status  string     `gorm:"type:varchar(16);default:sent" json:"status"` // sent/failed/pending
	Error   string     `gorm:"type:text" json:"error"`
	SentAt  *time.Time `json:"sent_at"`
}
