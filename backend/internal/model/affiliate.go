package model

import (
	"time"

	"gorm.io/gorm"
)

// Affiliate 推荐返利计划
type Affiliate struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	UserID         uint           `gorm:"uniqueIndex;not null" json:"user_id"`
	Code           string         `gorm:"type:varchar(20);uniqueIndex;not null" json:"code"`
	Balance        float64        `gorm:"type:decimal(10,2);default:0" json:"balance"`
	TotalEarned    float64        `gorm:"type:decimal(10,2);default:0" json:"total_earned"`
	TotalWithdrawn float64        `gorm:"type:decimal(10,2);default:0" json:"total_withdrawn"`
	CommissionRate float64        `gorm:"type:decimal(5,2);default:10" json:"commission_rate"`
	WithdrawMin    float64        `gorm:"type:decimal(10,2);default:100" json:"withdraw_min"`
	IsActive       bool           `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	User           *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// AffiliateRecord 推荐记录
type AffiliateRecord struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	AffiliateID uint       `gorm:"index;not null" json:"affiliate_id"`
	UserID      uint       `gorm:"index;not null" json:"user_id"`
	OrderID     uint       `gorm:"index" json:"order_id"`
	Amount      float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status      int8       `gorm:"type:smallint;default:1" json:"status"` // 1待确认 2已确认 3已拒绝
	Description string     `gorm:"type:text" json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Affiliate   *Affiliate `gorm:"foreignKey:AffiliateID" json:"affiliate,omitempty"`
}

// AffiliateVisit 推荐访问记录
type AffiliateVisit struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	AffiliateID uint      `gorm:"index;not null" json:"affiliate_id"`
	IP          string    `gorm:"type:varchar(64)" json:"ip"`
	UserAgent   string    `gorm:"type:text" json:"user_agent"`
	RefererURL  string    `gorm:"type:varchar(512)" json:"referer_url"`
	LandingURL  string    `gorm:"type:varchar(512)" json:"landing_url"`
	Converted   bool      `gorm:"default:false" json:"converted"` // 是否转化为注册
	UserID      *uint     `gorm:"index" json:"user_id"`          // 转化后的用户ID
	CreatedAt   time.Time `json:"created_at"`
	Affiliate   *Affiliate `gorm:"foreignKey:AffiliateID" json:"affiliate,omitempty"`
}

// AffiliateWithdraw 提现记录
type AffiliateWithdraw struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	AffiliateID uint       `gorm:"index;not null" json:"affiliate_id"`
	Amount      float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	Fee         float64    `gorm:"type:decimal(10,2);default:0" json:"fee"`
	Actual      float64    `gorm:"type:decimal(10,2);not null" json:"actual"`
	Method      string     `gorm:"type:varchar(50)" json:"method"`  // alipay, bank, balance
	Account     string     `gorm:"type:varchar(100)" json:"account"` // 收款账号
	Status      int8       `gorm:"type:smallint;default:1" json:"status"` // 1待审核 2已打款 3已拒绝
	AdminNote   string     `gorm:"type:text" json:"admin_note"`
	ProcessedAt *time.Time `json:"processed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Affiliate   *Affiliate `gorm:"foreignKey:AffiliateID" json:"affiliate,omitempty"`
}
