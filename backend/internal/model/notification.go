package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// NotificationTemplate 消息模板
type NotificationTemplate struct {
	gorm.Model
	Code      string `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Name      string `gorm:"type:varchar(100);not null" json:"name"`
	Channel   string `gorm:"type:varchar(20);not null;index" json:"channel"` // email, sms, wechat, webhook
	Subject   string `gorm:"type:varchar(255)" json:"subject"`
	Content   string `gorm:"type:text" json:"content"`
	Variables string `gorm:"type:text" json:"variables"` // 可用变量说明
	IsActive  bool   `gorm:"default:true" json:"is_active"`
}

// NotificationLog 通知日志
type NotificationLog struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"index" json:"user_id"`
	Channel   string     `gorm:"type:varchar(20);not null;index" json:"channel"`
	Template  string     `gorm:"type:varchar(50)" json:"template"`
	To        string     `gorm:"type:varchar(100)" json:"to"`
	Subject   string     `gorm:"type:varchar(255)" json:"subject"`
	Content   string     `gorm:"type:text" json:"content"`
	Status    int8       `gorm:"type:smallint;default:1" json:"status"` // 1发送中 2成功 3失败
	Error     string     `gorm:"type:text" json:"error"`
	SentAt    *time.Time `json:"sent_at"`
	CreatedAt time.Time  `json:"created_at"`
}

// WebhookConfig Webhook配置
type WebhookConfig struct {
	gorm.Model
	Name     string         `gorm:"type:varchar(100);not null" json:"name"`
	URL      string         `gorm:"type:varchar(500);not null" json:"url"`
	Secret   string         `gorm:"type:varchar(100)" json:"secret"`
	Events   datatypes.JSON `gorm:"type:jsonb" json:"events"` // 监听的事件类型
	IsActive bool           `gorm:"default:true" json:"is_active"`
}
