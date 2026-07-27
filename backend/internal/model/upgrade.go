package model

import (
	"time"

	"gorm.io/gorm"
)

// UpgradeOrder 升降级订单
type UpgradeOrder struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	OrderNo         string         `gorm:"type:varchar(32);uniqueIndex;not null" json:"order_no"`
	UserID          uint           `gorm:"index;not null" json:"user_id"`
	UserProductID   uint           `gorm:"index;not null" json:"user_product_id"`
	TargetProductID uint           `gorm:"index;not null" json:"target_product_id"`
	Type            int8           `gorm:"type:smallint;not null" json:"type"` // 1升级 2降级
	BillingCycle    string         `gorm:"type:varchar(20)" json:"billing_cycle"`
	Amount          float64        `gorm:"type:decimal(10,2);not null" json:"amount"` // 差价
	Status          int8           `gorm:"type:smallint;default:1" json:"status"`     // 1待支付 2已支付 3已开通 4已取消
	PaymentMethod   string         `gorm:"type:varchar(50)" json:"payment_method"`
	PaidAt          *time.Time     `json:"paid_at"`
	CompletedAt     *time.Time     `json:"completed_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
