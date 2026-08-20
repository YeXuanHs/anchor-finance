package model

import (
	"time"

	"gorm.io/gorm"
)

// PromoCode 优惠码/促销码
type PromoCode struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Code           string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"code"`
	Type           string         `gorm:"type:varchar(20);not null;comment:percent/fixed/override/free" json:"type"`
	Value          float64        `gorm:"type:decimal(10,2);default:0" json:"value"`
	Cycles         string         `gorm:"type:varchar(255);comment:适用周期逗号分隔" json:"cycles"`
	AppliesTo      string         `gorm:"type:varchar(500);comment:适用产品ID逗号分隔" json:"applies_to"`
	Requires       string         `gorm:"type:varchar(500);comment:要求产品ID逗号分隔" json:"requires"`
	Recurring      bool           `gorm:"default:false;comment:是否循环优惠" json:"recurring"`
	RecurFor       int            `gorm:"default:0;comment:循环次数" json:"recur_for"`
	RequiresExist  bool           `gorm:"default:false;comment:是否要求已有产品" json:"requires_exist"`
	MaxTimes       int            `gorm:"default:0;comment:最大使用次数" json:"max_times"`
	Lifelong       bool           `gorm:"default:false;comment:是否终身优惠" json:"lifelong"`
	OneTime        bool           `gorm:"default:false;comment:是否一次性" json:"one_time"`
	OnlyNewClient  bool           `gorm:"default:false;comment:仅新客户" json:"only_new_client"`
	OnlyOldClient  bool           `gorm:"default:false;comment:仅老客户" json:"only_old_client"`
	OncePerClient  bool           `gorm:"default:false;comment:每客户仅一次" json:"once_per_client"`
	Upgrades       bool           `gorm:"default:false;comment:是否支持升降级" json:"upgrades"`
	UpgradeConfig  string         `gorm:"type:text;comment:升降级配置JSON" json:"upgrade_config"`
	IsDiscount     bool           `gorm:"default:false;comment:是否为折扣码" json:"is_discount"`
	Notes          string         `gorm:"type:text" json:"notes"`
	UsedCount      int            `gorm:"default:0" json:"used_count"`
	StartTime      int64          `gorm:"not null" json:"start_time"`
	ExpirationTime int64          `gorm:"not null" json:"expiration_time"`
	Status         int8           `gorm:"default:1;comment:1=启用 0=禁用" json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// PromoCodeLog 优惠码使用记录
type PromoCodeLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PromoID   uint      `gorm:"index;not null" json:"promo_id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	OrderID   uint      `gorm:"index" json:"order_id"`
	Amount    float64   `gorm:"type:decimal(10,2)" json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

func (PromoCode) TableName() string    { return "promo_codes" }
func (PromoCodeLog) TableName() string { return "promo_code_logs" }
