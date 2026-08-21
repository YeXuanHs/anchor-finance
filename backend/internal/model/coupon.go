package model

import (
	"time"

	"gorm.io/gorm"
)

// Coupon 优惠券模型
type Coupon struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Code        string         `gorm:"size:50;uniqueIndex" json:"code"`
	Type        string         `gorm:"size:20;not null" json:"type"` // percent, fixed
	Value       float64        `gorm:"type:decimal(10,2)" json:"value"`
	MinAmount   float64        `gorm:"type:decimal(10,2);default:0" json:"min_amount"`
	MaxDiscount float64        `gorm:"type:decimal(10,2);default:0" json:"max_discount"`
	UsageLimit  int            `gorm:"default:0" json:"usage_limit"` // 0表示不限
	UsedCount   int            `gorm:"default:0" json:"used_count"`
	StartDate   *time.Time     `json:"start_date"`
	EndDate     *time.Time     `json:"end_date"`
	Status      string         `gorm:"size:20;default:active" json:"status"` // active, disabled, expired
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Coupon) TableName() string {
	return "coupons"
}

// CouponCampaign 优惠券活动模型
type CouponCampaign struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Description string         `gorm:"size:500" json:"description"`
	CouponID    uint           `gorm:"index" json:"coupon_id"`
	StartDate   *time.Time     `json:"start_date"`
	EndDate     *time.Time     `json:"end_date"`
	Status      string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (CouponCampaign) TableName() string {
	return "coupon_campaigns"
}
