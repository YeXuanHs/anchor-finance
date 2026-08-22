package model

import (
	"time"

	"gorm.io/gorm"
)

// OrderConfig 商品订购设置模型
type OrderConfig struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Key       string         `gorm:"size:100;uniqueIndex;not null" json:"key"`
	Value     string         `gorm:"type:text" json:"value"`
	Group     string         `gorm:"size:50;default:order" json:"group"`
	Name      string         `gorm:"size:100" json:"name"`
	Type      string         `gorm:"size:20;default:text" json:"type"` // text, number, select, switch
	Options   string         `gorm:"type:text" json:"options"`         // JSON格式的选项
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (OrderConfig) TableName() string {
	return "order_configs"
}
