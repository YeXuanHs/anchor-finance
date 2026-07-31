package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Voucher 代金券
type Voucher struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Name       string         `gorm:"type:varchar(100);not null" json:"name"`
	Code       string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Type       int8           `gorm:"type:smallint;not null" json:"type"` // 1固定金额 2百分比
	Value      float64        `gorm:"type:decimal(10,2);not null" json:"value"`
	MaxAmount  float64        `gorm:"type:decimal(10,2);default:0" json:"max_amount"` // 百分比券的最大抵扣
	MinOrder   float64        `gorm:"type:decimal(10,2);default:0" json:"min_order"`  // 最低消费
	TotalCount int            `gorm:"default:0" json:"total_count"`                  // 0不限
	UsedCount  int            `gorm:"default:0" json:"used_count"`
	UserCount  int            `gorm:"default:1" json:"user_count"` // 每用户可用次数
	ProductIDs datatypes.JSON `gorm:"type:json" json:"product_ids"`                 // 适用产品，空=全部
	StartDate  time.Time      `json:"start_date"`
	EndDate    time.Time      `json:"end_date"`
	IsActive   bool           `gorm:"default:true" json:"is_active"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// VoucherRecord 代金券使用记录
type VoucherRecord struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	VoucherID uint      `gorm:"index;not null" json:"voucher_id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	OrderID   uint      `gorm:"index;not null" json:"order_id"`
	Amount    float64   `gorm:"type:decimal(10,2);not null" json:"amount"` // 抵扣金额
	CreatedAt time.Time `json:"created_at"`
	Voucher   *Voucher  `gorm:"foreignKey:VoucherID" json:"voucher,omitempty"`
}

// UserVoucher 用户领取的代金券
type UserVoucher struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"index;not null" json:"user_id"`
	VoucherID uint       `gorm:"index;not null" json:"voucher_id"`
	UsedAt    *time.Time `json:"used_at"`
	OrderID   *uint      `gorm:"index" json:"order_id"`
	CreatedAt time.Time  `json:"created_at"`
	Voucher   *Voucher   `gorm:"foreignKey:VoucherID" json:"voucher,omitempty"`
}
