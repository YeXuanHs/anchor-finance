package model

import (
	"time"

	"gorm.io/gorm"
)

// Invoice 账单模型
type Invoice struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"index;not null" json:"user_id"`
	InvoiceNo     string         `gorm:"size:50;uniqueIndex" json:"invoice_no"`
	OrderID       uint           `gorm:"index" json:"order_id"`
	Amount        float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status        string         `gorm:"size:20;default:unpaid" json:"status"` // unpaid, paid, cancelled, refunded
	PaymentMethod string         `gorm:"size:50" json:"payment_method"`
	DueDate       *time.Time     `json:"due_date"`
	PaidAt        *time.Time     `json:"paid_at"`
	Note          string         `gorm:"type:text" json:"note"` // 管理员备注
	Remark        string         `gorm:"size:500" json:"remark"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Invoice) TableName() string {
	return "invoices"
}
