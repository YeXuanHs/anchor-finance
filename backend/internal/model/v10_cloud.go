package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// V10CloudProduct V10云产品
type V10CloudProduct struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"size:256;not null" json:"name"`
	GroupID       uint           `gorm:"index" json:"group_id"`
	Region        string         `gorm:"size:64" json:"region"`
	CPU           int            `json:"cpu"`
	MemoryMB      int            `json:"memory_mb"`
	DiskGB        int            `json:"disk_gb"`
	BandwidthMbps int            `json:"bandwidth_mbps"`
	TrafficGB     int            `json:"traffic_gb"`
	OS            string         `gorm:"size:128" json:"os"`
	Price         float64        `gorm:"type:decimal(12,2)" json:"price"`
	Cycle         string         `gorm:"size:32" json:"cycle"` // monthly/quarterly/yearly
	Stock         int            `json:"stock"`
	MaxStock      int            `json:"max_stock"`
	UpstreamID    uint           `gorm:"index" json:"upstream_id"`
	UpstreamType  string         `gorm:"size:32" json:"upstream_type"`
	Config        datatypes.JSON `gorm:"type:json" json:"config"`
	Enabled       bool           `gorm:"default:true" json:"enabled"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func (V10CloudProduct) TableName() string {
	return "v10_cloud_products"
}

// V10CloudOrder V10云订单
type V10CloudOrder struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `gorm:"index" json:"user_id"`
	ProductID  uint           `gorm:"index" json:"product_id"`
	HostID     *uint          `gorm:"index" json:"host_id"`
	Cycle      string         `gorm:"size:32" json:"cycle"`
	Quantity   int            `gorm:"default:1" json:"quantity"`
	Config     datatypes.JSON `gorm:"type:json" json:"config"`
	TotalPrice float64        `gorm:"type:decimal(12,2)" json:"total_price"`
	Status     string         `gorm:"size:32" json:"status"` // pending/paid/provisioning/active/failed
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

func (V10CloudOrder) TableName() string {
	return "v10_cloud_orders"
}

// V10CloudConfigOption V10云配置选项
type V10CloudConfigOption struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ProductID uint           `gorm:"index" json:"product_id"`
	Name      string         `gorm:"size:128;not null" json:"name"`
	Type      string         `gorm:"size:32;not null" json:"type"` // cpu/memory/disk/bandwidth/traffic/os/region
	Value     string         `gorm:"size:256" json:"value"`
	Label     string         `gorm:"size:128" json:"label"`
	ParentID  *uint          `gorm:"index" json:"parent_id"`
	Price     float64        `gorm:"type:decimal(12,2);default:0" json:"price"`
	Stock     int            `gorm:"default:-1" json:"stock"` // -1=unlimited
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	Enabled   bool           `gorm:"default:true" json:"enabled"`
	Config    datatypes.JSON `gorm:"type:json" json:"config"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (V10CloudConfigOption) TableName() string {
	return "v10_cloud_config_options"
}

// V10CloudRegion V10云区域
type V10CloudRegion struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Name      string `gorm:"size:128;not null" json:"name"`
	Code      string `gorm:"size:32;uniqueIndex;not null" json:"code"`
	Country   string `gorm:"size:64" json:"country"`
	City      string `gorm:"size:64" json:"city"`
	Enabled   bool   `gorm:"default:true" json:"enabled"`
	SortOrder int    `gorm:"default:0" json:"sort_order"`
}

func (V10CloudRegion) TableName() string {
	return "v10_cloud_regions"
}

// V10CloudOSType V10云操作系统类型
type V10CloudOSType struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Name      string `gorm:"size:128;not null" json:"name"`
	Version   string `gorm:"size:64" json:"version"`
	Arch      string `gorm:"size:32" json:"arch"` // amd64/arm64
	Type      string `gorm:"size:32" json:"type"` // linux/windows/other
	Enabled   bool   `gorm:"default:true" json:"enabled"`
	SortOrder int    `gorm:"default:0" json:"sort_order"`
}

func (V10CloudOSType) TableName() string {
	return "v10_cloud_os_types"
}
