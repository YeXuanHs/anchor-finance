package model

import "gorm.io/gorm"

// V10Cart V10购物车
type V10Cart struct {
	gorm.Model
	UserID       uint    `gorm:"index;not null" json:"user_id"`
	ProductID    uint    `gorm:"index;not null" json:"product_id"`
	ConfigOption string  `gorm:"type:json" json:"config_option"`
	Cycle        string  `gorm:"type:varchar(16);not null" json:"cycle"` // monthly/quarterly/semi-annually/annually
	Quantity     int     `gorm:"default:1" json:"quantity"`
	UnitPrice    float64 `gorm:"type:decimal(12,2)" json:"unit_price"`
	CouponCode   string  `gorm:"type:varchar(64)" json:"coupon_code"`
}
