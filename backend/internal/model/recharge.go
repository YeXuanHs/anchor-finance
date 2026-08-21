package model

import (
	"time"

	"gorm.io/gorm"
)

// Recharge 充值记录模型
type Recharge struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"index;not null" json:"user_id"`
	Amount        float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
	Gateway       string         `gorm:"size:50" json:"gateway"` // alipay, wxpay, etc.
	TransactionNo string         `gorm:"size:100" json:"transaction_no"`
	Status        string         `gorm:"size:20;default:pending" json:"status"` // pending, success, failed
	PaidAt        *time.Time     `json:"paid_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Recharge) TableName() string {
	return "recharges"
}
