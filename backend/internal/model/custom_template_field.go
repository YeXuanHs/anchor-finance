package model

import (
	"time"

	"gorm.io/gorm"
)

// CustomTemplateField 官网自定义字段模型
type CustomTemplateField struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Page      string         `gorm:"size:50;not null;index" json:"page"` // home, about, contact, etc.
	Name      string         `gorm:"size:100;not null" json:"name"`
	Key       string         `gorm:"size:100;not null" json:"key"`
	Type      string         `gorm:"size:20;default:text" json:"type"` // text, textarea, image, html
	Value     string         `gorm:"type:text" json:"value"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	Status    string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (CustomTemplateField) TableName() string {
	return "custom_template_fields"
}
