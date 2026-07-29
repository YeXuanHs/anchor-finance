package model

import "gorm.io/gorm"

// PaymentGateway 支付网关配置
type PaymentGateway struct {
	gorm.Model
	Name                string  `gorm:"type:varchar(64);not null" json:"name"`
	Code                string  `gorm:"type:varchar(32);uniqueIndex;not null" json:"code"`
	Description         string  `gorm:"type:varchar(256)" json:"description"`
	Icon                string  `gorm:"type:varchar(256)" json:"icon"`
	Config              string  `gorm:"type:jsonb" json:"config"`
	FeeRate             float64 `gorm:"type:decimal(5,4);default:0" json:"fee_rate"`
	MinAmount           float64 `gorm:"type:decimal(12,2);default:0" json:"min_amount"`
	MaxAmount           float64 `gorm:"type:decimal(12,2);default:0" json:"max_amount"`
	SortOrder           int     `gorm:"default:0" json:"sort_order"`
	IsEnabled           bool    `gorm:"default:true;index" json:"is_enabled"`
	SupportedCurrencies string  `gorm:"type:varchar(256)" json:"supported_currencies"`
}
