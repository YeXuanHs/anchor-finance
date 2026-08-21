package model

import (
	"time"

	"gorm.io/gorm"
)

// Referral 推介记录模型
type Referral struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `gorm:"index;not null" json:"user_id"`
	ReferrerID uint           `gorm:"index;not null" json:"referrer_id"` // 推荐人ID
	Status     string         `gorm:"size:20;default:pending" json:"status"` // pending, completed, cancelled
	Reward     float64        `gorm:"type:decimal(10,2);default:0" json:"reward"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Referral) TableName() string {
	return "referrals"
}

// ReferralWithdrawal 推介提现模型
type ReferralWithdrawal struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	Amount    float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status    string         `gorm:"size:20;default:pending" json:"status"` // pending, approved, rejected, paid
	Remark    string         `gorm:"size:500" json:"remark"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (ReferralWithdrawal) TableName() string {
	return "referral_withdrawals"
}
