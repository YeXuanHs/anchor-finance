package model

import (
	"time"

	"gorm.io/gorm"
)

// Product 产品模型
type Product struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"size:100;not null" json:"name"`
	ProductTypeID uint           `gorm:"index" json:"product_type_id"`
	GroupID       uint           `gorm:"index" json:"group_id"`
	Type          string         `gorm:"size:50" json:"type"` // host, server, domain, etc.
	Description   string         `gorm:"size:500" json:"description"`
	Amount        float64        `gorm:"type:decimal(10,2)" json:"-"`      // 兼容字段
	Price         float64        `gorm:"type:decimal(10,2)" json:"price"`  // 唯一价格字段
	BillingCycle  string         `gorm:"size:20" json:"billing_cycle"`
	ConfigOptions string         `gorm:"type:text" json:"config_options"` // JSON配置选项
	Status        string         `gorm:"size:20;default:active" json:"status"` // active, hidden, deleted
	Hidden        int            `gorm:"default:0" json:"hidden"`        // 0=显示, 1=隐藏（zjmf兼容）
	StockControl  int            `gorm:"default:0" json:"stock_control"` // 0=不限制, 1=限制库存
	Qty           int            `gorm:"default:999" json:"qty"`         // 库存数量
	SetupFee      float64        `gorm:"type:decimal(10,2);default:0" json:"setup_fee"` // 设置费
	PayType       string         `gorm:"size:50" json:"pay_type"`        // 付款类型（月付/年付等）
	AutoSetup     string         `gorm:"size:50" json:"auto_setup"`      // 自动开通方式
	IsFeatured    int            `gorm:"default:0" json:"is_featured"`   // 是否推荐
	SortOrder     int            `gorm:"default:0" json:"sort_order"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Product) TableName() string {
	return "products"
}

// ProductGroup 产品分组模型
type ProductGroup struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	ParentID    uint           `gorm:"index;default:0" json:"parent_id"`
	Description string         `gorm:"size:500" json:"description"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	Status      string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (ProductGroup) TableName() string {
	return "product_groups"
}
