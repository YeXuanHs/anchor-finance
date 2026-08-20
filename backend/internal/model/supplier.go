package model

import (
	"time"

	"gorm.io/gorm"
)

// Supplier 供应商模型
type Supplier struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	Name               string         `gorm:"size:100;not null" json:"name"`
	APIType            string         `gorm:"size:20;not null;column:api_type" json:"api_type"` // manual, zjmf, v10, anchor
	APIURL             string         `gorm:"size:255;column:api_url" json:"api_url"`
	APIKey             string         `gorm:"size:255;column:api_key" json:"api_key"`
	APIPassword        string         `gorm:"size:255;column:api_password" json:"-"`
	Description        string         `gorm:"type:text" json:"description"`
	Status             string         `gorm:"size:20;default:active" json:"status"` // active, disabled
	ProductsCount      int            `gorm:"default:0;column:products_count" json:"products_count"`
	LastSyncAt         *time.Time     `gorm:"column:last_sync_at" json:"last_sync_at"`
	Balance            float64        `gorm:"type:decimal(10,2);default:0" json:"balance"`
	AvailableProducts  int            `gorm:"default:0;column:available_products" json:"available_products"`
	ConfiguredProducts int            `gorm:"default:0;column:configured_products" json:"configured_products"`
	NormalProducts     int            `gorm:"default:0;column:normal_products" json:"normal_products"`
	TotalProducts      int            `gorm:"default:0;column:total_products" json:"total_products"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

// SupplierQueryParams 供应商查询参数
type SupplierQueryParams struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Keyword  string `json:"keyword"`
	APIType  string `json:"api_type"`
}

// CreateSupplierRequest 创建供应商请求
type CreateSupplierRequest struct {
	Name        string `json:"name" binding:"required"`
	APIType     string `json:"api_type" binding:"required"`
	APIURL      string `json:"api_url"`
	APIKey      string `json:"api_key"`
	APIPassword string `json:"api_password"`
	Description string `json:"description"`
}

// UpdateSupplierRequest 更新供应商请求
type UpdateSupplierRequest struct {
	Name        string `json:"name"`
	APIType     string `json:"api_type"`
	APIURL      string `json:"api_url"`
	APIKey      string `json:"api_key"`
	APIPassword string `json:"api_password"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// SupplierListResponse 供应商列表响应
type SupplierListResponse struct {
	List  []Supplier `json:"list"`
	Total int64      `json:"total"`
	Page  int        `json:"page"`
	PageSize int     `json:"page_size"`
}
