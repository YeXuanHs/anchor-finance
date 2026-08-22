package model

import (
	"time"

	"gorm.io/gorm"
)

// SMSTemplate 短信模板模型
type SMSTemplate struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Code      string         `gorm:"size:50;uniqueIndex;not null" json:"code"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	Variables string         `gorm:"type:text" json:"variables"` // JSON格式的变量说明
	Status    string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (SMSTemplate) TableName() string {
	return "sms_templates"
}
