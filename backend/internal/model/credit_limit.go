package model

import (
	"time"

	"gorm.io/gorm"
)

// CreditLimit 信用额度
type CreditLimit struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	UserID             uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	Limit              float64   `gorm:"type:decimal(20,4);default:0" json:"limit"`
	Used               float64   `gorm:"type:decimal(20,4);default:0" json:"used"`
	Available          float64   `gorm:"type:decimal(20,4);default:0" json:"available"`
	BillGenerationDay  int       `gorm:"default:1" json:"bill_generation_day"`  // day of month to generate bill (1-28)
	RepaymentPeriod    int       `gorm:"default:15" json:"repayment_period"`    // days after billing to repay
	PrepaymentBalance  float64   `gorm:"type:decimal(20,4);default:0" json:"prepayment_balance"`
	Description        string    `gorm:"type:text" json:"description"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	User               *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// CreditLog 信用额度变动日志
type CreditLog struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`
	Type        string    `gorm:"type:varchar(20);not null" json:"type"`   // adjust, use, repay, apply
	Status      string    `gorm:"type:varchar(20);default:completed" json:"status"` // pending, approved, rejected, completed
	Amount      float64   `gorm:"type:decimal(20,4);not null" json:"amount"`
	Balance     float64   `gorm:"type:decimal(20,4);not null" json:"balance"`
	RelatedID   uint      `gorm:"index" json:"related_id"`
	RelatedType string    `gorm:"type:varchar(20)" json:"related_type"`
	AdminID     *uint     `gorm:"index" json:"admin_id"`
	Remark      string    `gorm:"type:text" json:"remark"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreditBill 信用账单
type CreditBill struct {
	gorm.Model
	UserID          uint      `gorm:"index;not null" json:"user_id"`
	BillMonth       string    `gorm:"type:varchar(7);not null;index" json:"bill_month"` // "2026-07"
	BillingDate     time.Time `gorm:"not null" json:"billing_date"`
	DueDate         time.Time `gorm:"not null" json:"due_date"`
	TotalAmount     float64   `gorm:"type:decimal(20,4);not null" json:"total_amount"`
	PaidAmount      float64   `gorm:"type:decimal(20,4);default:0" json:"paid_amount"`
	RemainingAmount float64   `gorm:"type:decimal(20,4);not null" json:"remaining_amount"`
	LateFee         float64   `gorm:"type:decimal(20,4);default:0" json:"late_fee"`
	Status          string    `gorm:"type:varchar(32);default:'unpaid'" json:"status"` // unpaid, partial, paid, overdue, written_off
	User            *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// CreditBillItem 信用账单明细
type CreditBillItem struct {
	gorm.Model
	BillID      uint    `gorm:"index;not null" json:"bill_id"`
	Type        string  `gorm:"type:varchar(32);not null" json:"type"` // usage, fee, payment
	Description string  `gorm:"type:varchar(255)" json:"description"`
	Amount      float64 `gorm:"type:decimal(20,4);not null" json:"amount"`
	RelatedID   uint    `gorm:"index" json:"related_id"`
	RelatedType string  `gorm:"type:varchar(32)" json:"related_type"`
	Bill        *CreditBill `gorm:"foreignKey:BillID" json:"bill,omitempty"`
}
