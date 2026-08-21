package model

import (
	"time"

	"gorm.io/gorm"
)

// Currency 货币模型
type Currency struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Code      string         `gorm:"size:10;uniqueIndex;not null" json:"code"` // CNY, USD, etc.
	Name      string         `gorm:"size:50;not null" json:"name"`
	Symbol    string         `gorm:"size:10" json:"symbol"` // ¥, $, etc.
	Rate      float64        `gorm:"type:decimal(10,4);default:1" json:"rate"` // 汇率
	IsDefault bool           `gorm:"default:false" json:"is_default"`
	Status    string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Currency) TableName() string {
	return "currencies"
}
