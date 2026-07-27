package model

import (
	"time"
)

// BalanceLog 余额变动日志
type BalanceLog struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`
	Amount      float64   `gorm:"type:decimal(12,2);not null" json:"amount"` // positive=credit, negative=debit
	Balance     float64   `gorm:"type:decimal(12,2);not null" json:"balance"` // balance after transaction
	RelatedID   uint      `json:"related_id"`
	RelatedType string    `gorm:"type:varchar(32)" json:"related_type"` // order/recharge/refund/admin
	Description string    `gorm:"type:varchar(255)" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// TableName overrides the default table name.
func (BalanceLog) TableName() string {
	return "balance_logs"
}
