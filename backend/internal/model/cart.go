package model

import (
	"time"

	"gorm.io/gorm"
)

// CartItem 购物车项
type CartItem struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"index;not null" json:"user_id"`
	ProductID   uint           `gorm:"index;not null" json:"product_id"`
	ProductName string         `gorm:"size:200" json:"product_name"`
	ConfigID    uint           `json:"config_id"`                              // 可配置项ID
	ConfigJSON  string         `gorm:"type:text" json:"config_json"`           // 可配置项选择(JSON)
	Quantity    int            `gorm:"default:1" json:"quantity"`
	Cycle       string         `gorm:"size:20" json:"cycle"`                   // monthly, quarterly, semiannually, annually, biennially, triennially
	Amount      float64        `gorm:"type:decimal(10,2)" json:"amount"`       // 服务端计算的价格，不信任前端
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CartItem) TableName() string {
	return "cart_items"
}
