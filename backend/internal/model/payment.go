package model

import (
	"time"

	"gorm.io/gorm"
)

// Payment 支付记录模型
type Payment struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"index;not null" json:"user_id"`
	InvoiceID     uint           `gorm:"index" json:"invoice_id"`
	Amount        float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
	Gateway       string         `gorm:"size:50" json:"gateway"` // alipay, wxpay, balance, etc.
	TransactionNo string         `gorm:"size:100" json:"transaction_no"`
	Status        string         `gorm:"size:20;default:pending" json:"status"` // pending, success, failed, refunded
	PaidAt        *time.Time     `json:"paid_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Payment) TableName() string {
	return "payments"
}
