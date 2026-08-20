package model

import (
	"gorm.io/gorm"
)

// Pricing 商品多周期定价
type Pricing struct {
	gorm.Model
	ProductID uint             `gorm:"index;not null" json:"product_id"`
	Product   *Product         `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Cycle     string           `gorm:"type:varchar(32);not null;index" json:"cycle"` // monthly/quarterly/semi-annually/annually/biennially/triennially/hourly/daily/onetime
	Price     float64 `gorm:"type:decimal(20,4);not null" json:"price"`
	Currency  string           `gorm:"type:varchar(8);default:'CNY'" json:"currency"`
	SortOrder int              `gorm:"default:0" json:"sort_order"`
	Enabled   bool             `gorm:"default:true;index" json:"enabled"`
}

// TableName 指定表名
func (Pricing) TableName() string {
	return "pricings"
}
