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
	PaymentMethod string         `gorm:"size:50" json:"payment_method"`
	PaidAt        *time.Time     `json:"paid_at"`
	Remark        string         `gorm:"size:500" json:"remark"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Order) TableName() string {
	return "orders"
}
