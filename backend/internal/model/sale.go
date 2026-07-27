package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// SalePromotion 促销活动
type SalePromotion struct {
	gorm.Model
	Name          string         `gorm:"type:varchar(256);not null" json:"name"`
	Code          string         `gorm:"type:varchar(64);uniqueIndex" json:"code"` // 优惠码
	Type          string         `gorm:"type:varchar(32);not null;index" json:"type"` // amount_off(满减)/percent_off(折扣)/first_order(首单优惠)/free_trial(免费试用)
	Condition     datatypes.JSON `gorm:"type:jsonb;not null" json:"condition"` // 条件JSON：min_amount(最低消费)/min_quantity(最低数量)/applicable_products(适用商品)/applicable_groups(适用分组)
	Discount      datatypes.JSON `gorm:"type:jsonb;not null" json:"discount"`  // 优惠JSON：value(优惠值)/max_discount(最大折扣额)/free_product_id(赠送商品)
	StartAt       time.Time      `gorm:"index;not null" json:"start_at"`
	EndAt         time.Time      `gorm:"index;not null" json:"end_at"`
	MaxUses       int            `gorm:"default:0" json:"maxUses"` // 0=不限
	UsedCount     int            `gorm:"default:0" json:"usedCount"`
	MaxUsesPerUser int           `gorm:"default:1" json:"maxUsesPerUser"` // 每用户限用次数，0=不限
	Stackable     bool           `gorm:"default:false" json:"stackable"` // 是否可叠加
	AutoApply     bool           `gorm:"default:false" json:"autoApply"` // 自动应用（无需优惠码）
	Status        int16          `gorm:"type:smallint;default:1;not null;index" json:"status"` // 1=启用 0=禁用 2=已过期
	Description   string         `gorm:"type:text" json:"description"`
	AdminNotes    string         `gorm:"type:text" json:"admin_notes"`
	UsageLogs     []SalePromotionLog `gorm:"foreignKey:PromotionID" json:"usage_logs,omitempty"`
}

// SalePromotionLog 促销使用记录
type SalePromotionLog struct {
	gorm.Model
	PromotionID uint    `gorm:"index;not null" json:"promotion_id"`
	Promotion   SalePromotion `gorm:"foreignKey:PromotionID" json:"promotion,omitempty"`
	UserID      uint    `gorm:"index;not null" json:"user_id"`
	User        User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	OrderID     *uint   `gorm:"index" json:"order_id"`
	Order       *Order  `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	DiscountAmount float64 `gorm:"type:decimal(12,2);not null" json:"discount_amount"`
	UsedAt      time.Time `gorm:"not null" json:"used_at"`
}
