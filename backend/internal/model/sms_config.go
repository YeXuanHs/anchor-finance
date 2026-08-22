package model

import (
	"time"

	"gorm.io/gorm"
)

// SMSConfig 短信配置模型
type SMSConfig struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Key       string         `gorm:"size:100;uniqueIndex;not null" json:"key"`
	Value     string         `gorm:"type:text" json:"value"`
	Group     string         `gorm:"size:50;default:sms" json:"group"`
	Name      string         `gorm:"size:100" json:"name"`
	Type      string         `gorm:"size:20;default:text" json:"type"` // text, number, select, switch
	Options   string         `gorm:"type:text" json:"options"`         // JSON格式的选项
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (SMSConfig) TableName() string {
	return "sms_configs"
}
