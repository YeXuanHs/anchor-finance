package model

import (
	"time"

	"gorm.io/gorm"
)

// Order 订单模型
type Order struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"index;not null" json:"user_id"`
	OrderNo       string         `gorm:"size:50;uniqueIndex;not null" json:"order_no"`
	ProductID     uint           `gorm:"index" json:"product_id"`
	ProductName   string         `gorm:"size:100" json:"product_name"`
	Quantity      int            `gorm:"default:1" json:"quantity"`
	Amount        float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status        string         `gorm:"size:20;default:pending" json:"status"` // pending, paid, active, cancelled, refunded
	Type          string         `gorm:"size:20;default:new" json:"type"`       // new, renew, upgrade
	PaymentMethod string         `gorm:"size:50" json:"payment_method"`
	PaidAt        *time.Time     `json:"paid_at"`
	Remark        string         `gorm:"size:500" json:"remark"`
	Note          string         `gorm:"type:text" json:"note"` // 管理员备注
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Order) TableName() string {
	return "orders"
}

// OrderItem 订单项
type OrderItem struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	OrderID     uint    `gorm:"index;not null" json:"order_id"`
	ProductID   uint    `gorm:"index" json:"product_id"`
	ProductName string  `gorm:"size:200" json:"product_name"`
	Quantity    int     `gorm:"default:1" json:"quantity"`
	Cycle       string  `gorm:"size:20" json:"cycle"`
	Amount      float64 `gorm:"type:decimal(10,2)" json:"amount"`
}

// TableName 指定表名
func (OrderItem) TableName() string {
	return "order_items"
}
