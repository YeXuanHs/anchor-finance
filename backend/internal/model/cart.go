package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ShoppingCart 购物车
type ShoppingCart struct {
	gorm.Model
	UserID       uint           `gorm:"index;not null" json:"user_id"`
	ProductID    uint           `gorm:"index;not null" json:"product_id"`
	BillingCycle string         `gorm:"type:varchar(32);not null" json:"billing_cycle"`
	Quantity     int            `gorm:"default:1" json:"quantity"`
	Config       datatypes.JSON `gorm:"type:jsonb" json:"config"`
	Domain       string         `gorm:"type:varchar(255)" json:"domain"`
}

// TableName overrides the default table name.
func (ShoppingCart) TableName() string {
	return "shopping_carts"
}
