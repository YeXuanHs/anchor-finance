package model

import (
	"time"

	"gorm.io/gorm"
)

// ProductConfigGroup 产品配置选项分组
// zjmf表: product_config_groups
type ProductConfigGroup struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Description string         `gorm:"size:500" json:"description"`
	Order       int            `gorm:"default:0" json:"order"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ProductConfigGroup) TableName() string { return "product_config_groups" }

// ProductConfigOption 产品配置选项
// zjmf表: product_config_options
type ProductConfigOption struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	GID            uint           `gorm:"column:gid;index;not null" json:"gid"`             // 关联config_group
	OptionName     string         `gorm:"size:200;not null" json:"option_name"`  // 选项名（支持"价格|显示名"格式）
	OptionType     int            `gorm:"default:0" json:"option_type"`          // 类型: 1=单选, 2=多选, 3=下拉, 5=操作系统, 6=数量, 8=文本
	Order          int            `gorm:"default:0" json:"order"`
	Hidden         bool           `gorm:"default:false" json:"hidden"`
	Auto           bool           `gorm:"default:false" json:"auto"`             // 自动选择
	IsDiscount     bool           `gorm:"default:false" json:"is_discount"`      // 是否打折
	IsRebate       bool           `gorm:"default:false" json:"is_rebate"`        // 是否返利
	QtyMinimum     int            `gorm:"default:0" json:"qty_minimum"`
	QtyMaximum     int            `gorm:"default:0" json:"qty_maximum"`
	QtyStage       int            `gorm:"default:1" json:"qty_stage"`
	Unit           string         `gorm:"size:20" json:"unit"`
	Upgrade        bool           `gorm:"default:false" json:"upgrade"`
	Notes          string         `gorm:"size:500" json:"notes"`
	UpstreamID     int            `gorm:"default:0" json:"upstream_id"`          // 上游ID
	LinkagePID     int            `gorm:"column:linkage_pid;default:0" json:"linkage_pid"`
	LinkageTopPID  int            `gorm:"column:linkage_top_pid;default:0" json:"linkage_top_pid"`
	LinkageLevel   int            `gorm:"column:linkage_level;default:0" json:"linkage_level"`
	CopyID         int            `gorm:"column:copyid;default:0" json:"copy_id"`
	Senior         bool           `gorm:"default:false" json:"senior"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ProductConfigOption) TableName() string { return "product_config_options" }

// ProductConfigOptionSub 配置选项子项（下拉选项、操作系统列表等）
// zjmf表: product_config_options_sub
type ProductConfigOptionSub struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	ConfigID   uint           `gorm:"index;not null" json:"config_id"`    // 关联config_option
	OptionName string         `gorm:"size:200;not null" json:"option_name"`
	SortOrder  int            `gorm:"default:0" json:"sort_order"`
	Hidden     bool           `gorm:"default:false" json:"hidden"`
	CopyID     int            `gorm:"column:copyid;default:0" json:"copy_id"`
	UpstreamID int            `gorm:"default:0" json:"upstream_id"`
	LinkageLevel int          `gorm:"column:linkage_level;default:0" json:"linkage_level"`
	LinkagePID int            `gorm:"column:linkage_pid;default:0" json:"linkage_pid"`
	LinkageTopPID int         `gorm:"column:linkage_top_pid;default:0" json:"linkage_top_pid"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ProductConfigOptionSub) TableName() string { return "product_config_options_sub" }

// ProductConfigLink 产品↔配置分组关联
// zjmf表: product_config_links
type ProductConfigLink struct {
	ID  uint `gorm:"primaryKey" json:"id"`
	PID uint `gorm:"column:pid;index;not null" json:"pid"` // 产品ID
	GID uint `gorm:"column:gid;index;not null" json:"gid"` // 分组ID
}

func (ProductConfigLink) TableName() string { return "product_config_links" }

// ProductConfigPricing 配置选项子项定价
// zjmf表: pricing（type=config_option）
type ProductConfigPricing struct {
	ID         uint    `gorm:"primaryKey" json:"id"`
	RelID      uint    `gorm:"index;not null" json:"relid"`      // 子选项ID
	Type       string  `gorm:"size:20;not null" json:"type"`     // config_option
	Currency   int     `gorm:"not null" json:"currency"`         // 货币ID
	Monthly    float64 `gorm:"type:decimal(10,2);default:0" json:"monthly"`
	Quarterly  float64 `gorm:"type:decimal(10,2);default:0" json:"quarterly"`
	Semiannual float64 `gorm:"type:decimal(10,2);default:0" json:"semiannually"`
	Annually   float64 `gorm:"type:decimal(10,2);default:0" json:"annually"`
	Biennially float64 `gorm:"type:decimal(10,2);default:0" json:"biennially"`
	Triennially float64 `gorm:"type:decimal(10,2);default:0" json:"triennially"`
	MonthlySetup    float64 `gorm:"type:decimal(10,2);default:0" json:"monthlysetupfee"`
	QuarterlySetup  float64 `gorm:"type:decimal(10,2);default:0" json:"quarterlysetupfee"`
	SemiannualSetup float64 `gorm:"type:decimal(10,2);default:0" json:"semiannuallysetupfee"`
	AnnuallySetup   float64 `gorm:"type:decimal(10,2);default:0" json:"annuallysetupfee"`
	BienniallySetup float64 `gorm:"type:decimal(10,2);default:0" json:"bienniallysetupfee"`
	TrienniallySetup float64 `gorm:"type:decimal(10,2);default:0" json:"trienniallysetupfee"`
}

func (ProductConfigPricing) TableName() string { return "product_config_pricings" }
