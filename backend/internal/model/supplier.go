package model

import (
	"time"

	"gorm.io/gorm"
)

// Supplier 供应商模型
type Supplier struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Slug        string         `gorm:"size:50" json:"slug"`        // 插件标识
	Domain      string         `gorm:"size:50" json:"domain"`      // 插件域（server/payment/sms等）
	Type        string         `gorm:"size:50" json:"type"`        // manual, zjmf, v10, anchor
	APIURL      string         `gorm:"size:500" json:"api_url"`
	APIKey      string         `gorm:"size:500" json:"api_key"`    // API密钥（AES加密存储，json返回脱敏值）
	APISecret   string         `gorm:"size:500" json:"-"`          // API密钥（AES加密存储，json不返回）
	Description string         `gorm:"size:500" json:"description"`
	Status      string         `gorm:"size:20;default:active" json:"status"` // active, disabled
	Balance     float64        `gorm:"type:decimal(10,2);default:0" json:"balance"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Supplier) TableName() string {
	return "suppliers"
}

// SupplierProduct 供应商产品模型
type SupplierProduct struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	SupplierID        uint           `gorm:"index;not null" json:"supplier_id"`
	RemoteProductID   string         `gorm:"size:100" json:"remote_product_id"`
	RemoteProductName string         `gorm:"size:200" json:"remote_product_name"`
	Name              string         `gorm:"size:200" json:"name"`           // 商品名称
	LocalProductID    uint           `gorm:"index" json:"local_product_id"`
	ProfitRate        float64        `gorm:"type:decimal(5,2);default:25" json:"profit_rate"` // 利润率
	LocalPrice        float64        `gorm:"type:decimal(10,2)" json:"local_price"`
	RemotePrice       float64        `gorm:"type:decimal(10,2)" json:"remote_price"`
	Stock             int            `gorm:"default:0" json:"stock"`         // 库存
	LocalStock        int            `gorm:"default:0" json:"local_stock"`
	Status            string         `gorm:"size:20;default:active" json:"status"`
	DisableReason     string         `gorm:"size:100" json:"disable_reason"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (SupplierProduct) TableName() string {
	return "supplier_products"
}
