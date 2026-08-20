package model

import (
	"time"

	"gorm.io/datatypes"
)

// ProvisionButton 模块按钮配置（客户端/管理端动态按钮）
type ProvisionButton struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	ModuleID   uint   `gorm:"index;not null" json:"module_id"`
	Module     ProvisionModule `gorm:"foreignKey:ModuleID" json:"module,omitempty"`
	Name       string `gorm:"type:varchar(128);not null" json:"name"`
	Action     string `gorm:"type:varchar(64);not null" json:"action"`
	Type       string `gorm:"type:varchar(16);not null;index" json:"type"`       // client/admin
	Position   string `gorm:"type:varchar(16);not null;default:header" json:"position"` // header/sidebar/action
	Icon       string `gorm:"type:varchar(64)" json:"icon"`
	URL        string `gorm:"type:varchar(512)" json:"url"`
	Confirm    bool   `gorm:"default:false" json:"confirm"`
	ConfirmMsg string `gorm:"type:varchar(256)" json:"confirm_msg"`
	Enabled    bool   `gorm:"default:true;index" json:"enabled"`
	SortOrder  int    `gorm:"default:0" json:"sort_order"`
	CreatedAt  time.Time `json:"created_at"`
}

// ProvisionChart 模块图表配置
type ProvisionChart struct {
	ID       uint           `gorm:"primaryKey" json:"id"`
	ModuleID uint           `gorm:"index;not null" json:"module_id"`
	Module   ProvisionModule `gorm:"foreignKey:ModuleID" json:"module,omitempty"`
	Name     string         `gorm:"type:varchar(128);not null" json:"name"`
	Type     string         `gorm:"type:varchar(32);not null" json:"type"` // line/bar/pie/gauge
	Endpoint string         `gorm:"type:varchar(512)" json:"endpoint"`
	Config   datatypes.JSON `gorm:"type:json" json:"config"`
	Enabled  bool           `gorm:"default:true" json:"enabled"`
}

// CustomFunction 模块自定义函数
type CustomFunction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ModuleID  uint      `gorm:"index;not null" json:"module_id"`
	Module    ProvisionModule `gorm:"foreignKey:ModuleID" json:"module,omitempty"`
	Name      string    `gorm:"type:varchar(128);not null" json:"name"`
	Code      string    `gorm:"type:text" json:"code"`
	Trigger   string    `gorm:"type:varchar(32);not null;default:manual" json:"trigger"` // manual/auto/on_create/on_renew
	Enabled   bool      `gorm:"default:true" json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
}

// ProvisionUsage 模块用量跟踪
type ProvisionUsage struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	HostID      uint      `gorm:"index;not null" json:"host_id"`
	ModuleID    uint      `gorm:"index;not null" json:"module_id"`
	CPU         float64   `gorm:"default:0" json:"cpu"`
	Memory      float64   `gorm:"default:0" json:"memory"`
	Disk        float64   `gorm:"default:0" json:"disk"`
	Bandwidth   float64   `gorm:"default:0" json:"bandwidth"`
	TrafficIn   int64     `gorm:"default:0" json:"traffic_in"`
	TrafficOut  int64     `gorm:"default:0" json:"traffic_out"`
	Extra       datatypes.JSON `gorm:"type:json" json:"extra"`
	UpdatedAt   time.Time `json:"updated_at"`
}
