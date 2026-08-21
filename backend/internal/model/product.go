package model

import (
	"time"

	"gorm.io/gorm"
)

// Product 产品模型
type Product struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"size:100;not null" json:"name"`
	GroupID      uint           `gorm:"index" json:"group_id"`
	Type         string         `gorm:"size:50" json:"type"` // host, server, domain, etc.
	Description  string         `gorm:"size:500" json:"description"`
	Price        float64        `gorm:"type:decimal(10,2)" json:"price"`
	BillingCycle string         `gorm:"size:20" json:"billing_cycle"` // monthly, quarterly, yearly, onetime
	Status       string         `gorm:"size:20;default:active" json:"status"` // active, hidden, deleted
	SortOrder    int            `gorm:"default:0" json:"sort_order"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
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
