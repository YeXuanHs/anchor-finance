package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	gorm.Model
	Name         string         `gorm:"type:varchar(128);not null" json:"name"`
	Code         string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"code"`
	Type         string         `gorm:"type:varchar(32);not null;index" json:"type"` // vps/dedicated/cloud/reseller
	Provider     string         `gorm:"type:varchar(128)" json:"provider"`
	TemplateID   *uint          `gorm:"index" json:"template_id"`
	Template     *ServerTemplate `gorm:"foreignKey:TemplateID" json:"template,omitempty"`
	CPU          string         `gorm:"type:varchar(64)" json:"cpu"`
	Memory       int            `gorm:"default:0" json:"memory"` // MB
	Disk         int            `gorm:"default:0" json:"disk"`   // GB
	Bandwidth    int            `gorm:"default:0" json:"bandwidth"` // Mbps
	TrafficLimit int64          `gorm:"default:0" json:"traffic_limit"` // GB, 0=unlimited
	IPCount      int            `gorm:"default:1" json:"ip_count"`
	Location     string         `gorm:"type:varchar(128)" json:"location"`
	Datacenter   string         `gorm:"type:varchar(128)" json:"datacenter"`
	OS           datatypes.JSON `gorm:"type:json" json:"os"` // 支持的操作系统列表
	Features     datatypes.JSON `gorm:"type:json" json:"features"`
	PriceMonthly float64        `gorm:"type:decimal(10,2);default:0" json:"price_monthly"`
	PriceQuarter float64        `gorm:"type:decimal(10,2);default:0" json:"price_quarter"`
	PriceSemiAnn float64        `gorm:"type:decimal(10,2);default:0" json:"price_semi_annual"`
	PriceAnnual  float64        `gorm:"type:decimal(10,2);default:0" json:"price_annual"`
	PriceBiennial float64       `gorm:"type:decimal(10,2);default:0" json:"price_biennial"`
	PriceTriennial float64      `gorm:"type:decimal(10,2);default:0" json:"price_triennial"`
	PricingStrategy string      `gorm:"type:varchar(32);default:'fixed'" json:"pricing_strategy"` // fixed/graduated/promotional
	StockTotal   int            `gorm:"default:0" json:"stock_total"`
	StockUsed    int            `gorm:"default:0" json:"stock_used"`
	MaxPerUser   int            `gorm:"default:0" json:"max_per_user"` // 0=unlimited
	SortOrder    int            `gorm:"default:0" json:"sort_order"`
	Status       int16          `gorm:"type:smallint;default:1;not null;index" json:"status"` // 1=启用 0=禁用
	Remark       string         `gorm:"type:text" json:"remark"`
	Metadata     datatypes.JSON `gorm:"type:json" json:"metadata"`
	Products     []ServerProduct `gorm:"foreignKey:ServerConfigID" json:"products,omitempty"`
}

// ServerTemplate 服务器配置模板
type ServerTemplate struct {
	gorm.Model
	Name        string         `gorm:"type:varchar(128);not null" json:"name"`
	Code        string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"code"`
	Type        string         `gorm:"type:varchar(32);not null;index" json:"type"`
	Description string         `gorm:"type:text" json:"description"`
	Config      datatypes.JSON `gorm:"type:json;not null" json:"config"` // 模板配置参数
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	Status      int16          `gorm:"type:smallint;default:1;not null" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// ServerProduct 服务器与产品关联
type ServerProduct struct {
	ID           uint `gorm:"primaryKey" json:"id"`
	ServerConfigID uint `gorm:"index;not null" json:"server_config_id"`
	ProductID    uint `gorm:"index;not null" json:"product_id"`
	CreatedAt    time.Time `json:"created_at"`
}
