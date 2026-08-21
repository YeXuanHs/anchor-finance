package model

import (
	"time"

	"gorm.io/gorm"
)

// CreditLimit 信用额度模型
type CreditLimit struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"uniqueIndex;not null" json:"user_id"`
	Amount    float64        `gorm:"type:decimal(10,2);default:0" json:"amount"`
	Used      float64        `gorm:"type:decimal(10,2);default:0" json:"used"`
	Status    string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (CreditLimit) TableName() string {
	return "credit_limits"
}

// CreditLimitLog 信用额日志模型
type CreditLimitLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Type      string    `gorm:"size:20;not null" json:"type"` // add, reduce, use, refund
	Amount    float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	Remark    string    `gorm:"size:500" json:"remark"`
	AdminID   uint      `gorm:"index" json:"admin_id"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (CreditLimitLog) TableName() string {
	return "credit_limit_logs"
}
