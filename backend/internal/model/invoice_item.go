package model

import (
	"gorm.io/gorm"
)

// InvoiceItem 账单项目（独立模型文件，扩展 model/invoice.go 中已有的 InvoiceItem）
// 注意：此文件仅定义补充结构，InvoiceItem 主定义已在 invoice.go 中。
// 下面定义 InvoiceItemAdditional 用于存储额外的账单项信息。

// InvoiceItemDiscount 账单项折扣记录
type InvoiceItemDiscount struct {
	gorm.Model
	ItemID       uint    `gorm:"index;not null" json:"item_id"`
	Item         InvoiceItem `gorm:"foreignKey:ItemID" json:"item,omitempty"`
	PromotionID  *uint   `gorm:"index" json:"promotion_id"`
	Promotion    *SalePromotion `gorm:"foreignKey:PromotionID" json:"promotion,omitempty"`
	Type         string  `gorm:"type:varchar(32);not null" json:"type"` // coupon/promotion/bulk/manual
	DiscountType string  `gorm:"type:varchar(16);not null" json:"discount_type"` // amount/percent
	Value        float64 `gorm:"type:decimal(12,2);not null" json:"value"`
	AppliedAmount float64 `gorm:"type:decimal(12,2);not null" json:"applied_amount"` // 实际减免金额
	Description  string  `gorm:"type:varchar(256)" json:"description"`
}

// InvoiceItemTax 账单项税费
type InvoiceItemTax struct {
	gorm.Model
	ItemID      uint    `gorm:"index;not null" json:"item_id"`
	Item        InvoiceItem `gorm:"foreignKey:ItemID" json:"item,omitempty"`
	TaxName     string  `gorm:"type:varchar(64);not null" json:"tax_name"`
	TaxRate     float64 `gorm:"type:decimal(5,4);not null" json:"tax_rate"`
	TaxAmount   float64 `gorm:"type:decimal(12,2);not null" json:"tax_amount"`
	Description string  `gorm:"type:varchar(256)" json:"description"`
}
