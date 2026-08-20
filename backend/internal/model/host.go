package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Host 主机
type Host struct {
	gorm.Model
	Hostname     string         `gorm:"type:varchar(256);not null;index" json:"hostname"`
	IP           string         `gorm:"type:varchar(45);index" json:"ip"`
	IPv6         string         `gorm:"type:varchar(128)" json:"ipv6"`
	OS           string         `gorm:"type:varchar(128)" json:"os"` // 操作系统
	OSVersion    string         `gorm:"type:varchar(64)" json:"os_version"`
	CPU          string         `gorm:"type:varchar(128)" json:"cpu"`
	CPUCores     int            `gorm:"default:0" json:"cpu_cores"`
	MemoryMB     int            `gorm:"default:0" json:"memory_mb"`
	DiskSizeGB   int            `gorm:"default:0" json:"disk_size_gb"`
	DiskType     string         `gorm:"type:varchar(32)" json:"disk_type"` // SSD/HDD/NVMe
	BandwidthMbps int           `gorm:"default:0" json:"bandwidth_mbps"`
	TrafficGB    int            `gorm:"default:0" json:"traffic_gb"` // 月流量GB，0=不限
	Location     string         `gorm:"type:varchar(256)" json:"location"` // 机房位置
	Rack         string         `gorm:"type:varchar(32)" json:"rack"`
	RackPosition string         `gorm:"type:varchar(16)" json:"rack_position"`
	Status       int16          `gorm:"type:smallint;default:0;not null;index" json:"status"` // 0=关机 1=运行中 2=故障 3=维护 4=创建中
	PowerStatus  int8           `gorm:"type:smallint;default:0" json:"power_status"` // 0=off 1=on
	OwnerID      *uint          `gorm:"index" json:"owner_id"`
	Owner        *User          `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	ProductID    *uint          `gorm:"index" json:"product_id"` // 所属产品
	Product      *Product       `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	OrderID      *uint          `gorm:"index" json:"order_id"`
	Order        *Order         `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	ExpiredAt    *time.Time     `gorm:"index" json:"expired_at"` // 到期时间
	ProvisionedAt *time.Time    `json:"provisioned_at"`
	Remark       string         `gorm:"type:text" json:"remark"`
	AdminNotes   string         `gorm:"type:text" json:"admin_notes"`
	Tags         datatypes.JSON `gorm:"type:json" json:"tags"`
	Config       datatypes.JSON `gorm:"type:json" json:"config"` // 主机配置JSON
	Metadata     datatypes.JSON `gorm:"type:json" json:"metadata"`
	Operations   []HostOperation `gorm:"foreignKey:HostID" json:"operations,omitempty"`
}

// HostOperation 主机操作记录
type HostOperation struct {
	gorm.Model
	HostID     uint       `gorm:"index;not null" json:"host_id"`
	Host       Host       `gorm:"foreignKey:HostID" json:"host,omitempty"`
	OperatorID uint       `gorm:"index;not null" json:"operator_id"`
	Action     string     `gorm:"type:varchar(32);not null" json:"action"` // boot/shutdown/reboot/reinstall/rescue/migrate
	Params     string     `gorm:"type:text" json:"params"`
	Status     int8       `gorm:"type:smallint;default:1" json:"status"` // 1=执行中 2=成功 3=失败
	Result     string     `gorm:"type:text" json:"result"`
	ErrorMsg   string     `gorm:"type:text" json:"error_msg"`
	StartedAt  time.Time  `gorm:"not null" json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`
}
