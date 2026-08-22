package model

import (
	"time"

	"gorm.io/gorm"
)

// SalesConfig 销售设置模型
type SalesConfig struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Key       string         `gorm:"size:100;uniqueIndex;not null" json:"key"`
	Value     string         `gorm:"type:text" json:"value"`
	Group     string         `gorm:"size:50;default:sales" json:"group"`
	Name      string         `gorm:"size:100" json:"name"`
	Type      string         `gorm:"size:20;default:text" json:"type"` // text, number, select, switch
	Options   string         `gorm:"type:text" json:"options"`         // JSON格式的选项
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (SalesConfig) TableName() string {
	return "sales_configs"
}

// SalesGroup 销售分组模型
type SalesGroup struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Name       string         `gorm:"size:100;not null" json:"name"`
	Commission float64        `gorm:"default:0" json:"commission"` // 佣金比例
	Status     string         `gorm:"size:20;default:active" json:"status"`
	SortOrder  int            `gorm:"default:0" json:"sort_order"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (SalesGroup) TableName() string {
	return "sales_groups"
}
