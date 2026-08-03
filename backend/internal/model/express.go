package model

import (
	"time"

	"gorm.io/gorm"
)

// Express 快递公司
type Express struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:50;not null" json:"name"`      // 快递名称
	Code      string         `gorm:"size:20" json:"code"`               // 快递编码
	Price     float64        `gorm:"type:decimal(10,2);default:0" json:"price"` // 价格
	IsActive  bool           `gorm:"default:true" json:"is_active"`     // 是否启用
	SortOrder int            `gorm:"default:0" json:"sort_order"`       // 排序
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Express) TableName() string {
	return "express"
}
