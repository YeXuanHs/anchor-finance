package model

import (
	"time"

	"gorm.io/gorm"
)

// PromoCode 优惠码模型
// 注意：PromoCode 与 Coupon 功能有重叠但定位不同，暂不合并：
// - PromoCode 面向营销推广，通常通过渠道分发，绑定活动场景
// - Coupon 面向用户发放，通常绑定 CouponCampaign 活动，支持更灵活的发放策略
// 贸然合并可能导致业务逻辑混乱，后续如需统一应做完整的数据迁移和业务梳理。
type PromoCode struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Code        string         `gorm:"size:50;uniqueIndex;not null" json:"code"`
	Name        string         `gorm:"size:100" json:"name"`
	Type        string         `gorm:"size:20;not null" json:"type"` // percent, fixed
	Value       float64        `gorm:"type:decimal(10,2)" json:"value"` // 折扣百分比或固定金额
	MinAmount   float64        `gorm:"type:decimal(10,2);default:0" json:"min_amount"` // 最低消费
	MaxDiscount float64        `gorm:"type:decimal(10,2);default:0" json:"max_discount"` // 最大折扣
	UsageLimit  int            `gorm:"default:0" json:"usage_limit"` // 使用次数限制，0表示不限
	UsedCount   int            `gorm:"default:0" json:"used_count"`
	StartDate   *time.Time     `json:"start_date"`
	EndDate     *time.Time     `json:"end_date"`
	Status      string         `gorm:"size:20;default:active" json:"status"` // active, disabled, expired
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (PromoCode) TableName() string {
	return "promo_codes"
}
