package model

import (
	"time"

	"gorm.io/gorm"
)

// UserCoupon 用户优惠券领取记录（学创欧：coupon_id + user_id 唯一索引防重复领取）
type UserCoupon struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	CouponID   uint           `gorm:"index;not null" json:"coupon_id"`
	UserID     uint           `gorm:"index;not null" json:"user_id"`
	ReceiveType string        `gorm:"size:20;default:claim" json:"receive_type"` // claim, grant
	Status     int            `gorm:"default:1" json:"status"`                   // 1=OWNED, 2=USED, 3=REVOKED
	ClaimedAt  *time.Time     `json:"claimed_at"`
	UsedAt     *time.Time     `json:"used_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (UserCoupon) TableName() string {
	return "user_coupons"
}
