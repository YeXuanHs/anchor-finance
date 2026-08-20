package model

import "time"

// Server 服务器管理
type Server struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	Name            string     `gorm:"type:varchar(128);not null" json:"name"`
	Hostname        string     `gorm:"type:varchar(128)" json:"hostname"`
	IP              string     `gorm:"type:varchar(45);not null;index" json:"ip"`
	Port            int        `gorm:"default:22" json:"port"`
	Type            string     `gorm:"type:varchar(16);not null" json:"type"` // physical/cloud/vps
	DatacenterID    uint       `gorm:"index" json:"datacenter_id"`
	GroupID         uint       `gorm:"index" json:"group_id"`
	OperatingSystem string     `gorm:"type:varchar(64)" json:"operating_system"`
	CPU             string     `gorm:"type:varchar(64)" json:"cpu"`
	Memory          int        `gorm:"default:0" json:"memory"`      // MB
	Disk            int        `gorm:"default:0" json:"disk"`        // GB
	Bandwidth       int        `gorm:"default:0" json:"bandwidth"`   // Mbps
	MonthlyCost     float64    `gorm:"type:decimal(12,2)" json:"monthly_cost"`
	Status          string     `gorm:"type:varchar(16);default:active;index" json:"status"` // active/suspended/maintenance/offline
	Virtualization  string     `gorm:"type:varchar(32)" json:"virtualization"`
	Username        string     `gorm:"type:varchar(64)" json:"username"`
	Password        string     `gorm:"type:varchar(256)" json:"-"`
	Notes           string     `gorm:"type:text" json:"notes"`
	ProvisionedAt   *time.Time `json:"provisioned_at"`
	ExpiresAt       *time.Time `json:"expires_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
