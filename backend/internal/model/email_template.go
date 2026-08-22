package model

import (
	"time"

	"gorm.io/gorm"
)

// EmailTemplate 邮件模板模型
type EmailTemplate struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Subject   string         `gorm:"size:200" json:"subject"`
	Content   string         `gorm:"type:text" json:"content"`
	Variables string         `gorm:"type:text" json:"variables"` // JSON格式的变量说明
	Status    string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (EmailTemplate) TableName() string {
	return "email_templates"
}
