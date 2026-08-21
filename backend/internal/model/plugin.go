package model

import (
	"time"

	"gorm.io/gorm"
)

// Plugin 插件模型
type Plugin struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Slug        string         `gorm:"size:50;uniqueIndex" json:"slug"`
	Domain      string         `gorm:"size:50;index" json:"domain"` // payment, sms, mail, etc.
	Description string         `gorm:"size:500" json:"description"`
	Version     string         `gorm:"size:20" json:"version"`
	Author      string         `gorm:"size:100" json:"author"`
	Status      string         `gorm:"size:20;default:inactive" json:"status"` // active, inactive
	ConfigJSON  string         `gorm:"type:text" json:"config_json"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Plugin) TableName() string {
	return "plugins"
}
