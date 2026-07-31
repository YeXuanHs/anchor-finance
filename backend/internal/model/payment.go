package model

import (
	"time"

	"gorm.io/gorm"
)

// Payment 支付记录
type Payment struct {
	gorm.Model
	UserID        uint       `gorm:"index;not null" json:"user_id"`
	InvoiceID     uint       `gorm:"index" json:"invoice_id"`
	TradeNo       string     `gorm:"type:varchar(64);uniqueIndex;not null" json:"trade_no"`
	Gateway       string     `gorm:"type:varchar(32);not null" json:"gateway"`
	Amount        float64    `gorm:"type:decimal(12,2);not null" json:"amount"`
	ActualAmount  float64    `gorm:"type:decimal(12,2)" json:"actual_amount"`
	Fee           float64    `gorm:"type:decimal(12,2);default:0" json:"fee"`
	Currency      string     `gorm:"type:varchar(8);default:CNY" json:"currency"`
	Status        string     `gorm:"type:varchar(16);default:pending;index" json:"status"` // pending/paid/refunded/expired
	PaymentMethod string     `gorm:"type:varchar(32)" json:"payment_method"`
	CallbackURL   string     `gorm:"type:varchar(512)" json:"callback_url"`
	Description   string     `gorm:"type:varchar(256)" json:"description"`
	PaidAt        *time.Time `json:"paid_at"`
	RefundedAt    *time.Time `json:"refunded_at"`
	Extra         string     `gorm:"type:json" json:"extra"`
}
