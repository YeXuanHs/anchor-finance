package model

import (
	"time"

	"gorm.io/gorm"
)

// CustomField 自定义字段模型
type CustomField struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Label       string         `gorm:"size:100;not null" json:"label"`
	Type        string         `gorm:"size:50;not null" json:"type"` // text, textarea, select, checkbox, radio, date, number
	Options     string         `gorm:"type:text" json:"options"`     // JSON格式的选项（用于select, checkbox, radio）
	Required    bool           `gorm:"default:false" json:"required"`
	DefaultValue string        `gorm:"size:500" json:"default_value"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	Status      string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (CustomField) TableName() string {
	return "custom_fields"
}
