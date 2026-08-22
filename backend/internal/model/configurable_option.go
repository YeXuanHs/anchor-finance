package model

import (
	"time"

	"gorm.io/gorm"
)

// ConfigurableOption 全局可配置项模型
type ConfigurableOption struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Type      string         `gorm:"size:50;not null" json:"type"` // text, select, radio, checkbox, textarea
	Options   string         `gorm:"type:text" json:"options"`     // JSON格式的选项列表
	Default   string         `gorm:"size:500" json:"default_value"`
	Required  bool           `gorm:"default:false" json:"required"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	Status    string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (ConfigurableOption) TableName() string {
	return "configurable_options"
}
