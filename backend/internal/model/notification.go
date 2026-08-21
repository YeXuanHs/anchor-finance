package model

import (
	"time"

	"gorm.io/gorm"
)

// NotificationTemplate 通知模板模型
type NotificationTemplate struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Type      string         `gorm:"size:20;not null" json:"type"` // email, sms
	Subject   string         `gorm:"size:200" json:"subject"`
	Content   string         `gorm:"type:text" json:"content"`
	Variables string         `gorm:"type:text" json:"variables"` // JSON格式的变量说明
	Status    string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (NotificationTemplate) TableName() string {
	return "notification_templates"
}
