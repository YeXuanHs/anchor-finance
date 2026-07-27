package model

import (
	"time"

	"gorm.io/gorm"
)

// DcimDatacenter 机房
type DcimDatacenter struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(128);not null" json:"name"`
	Code        string         `gorm:"type:varchar(32);uniqueIndex;not null" json:"code"`
	Country     string         `gorm:"type:varchar(64)" json:"country"`
	City        string         `gorm:"type:varchar(64)" json:"city"`
	Address     string         `gorm:"type:varchar(256)" json:"address"`
	Provider    string         `gorm:"type:varchar(128)" json:"provider"`
	Network     string         `gorm:"type:varchar(128)" json:"network"` // 运营商线路
	Bandwidth   string         `gorm:"type:varchar(64)" json:"bandwidth"`
	Status      int8           `gorm:"type:smallint;default:1" json:"status"` // 1正常 0维护
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// DcimServer 物理服务器
type DcimServer struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Name           string         `gorm:"type:varchar(128);not null" json:"name"`
	Hostname       string         `gorm:"type:varchar(128)" json:"hostname"`
	IP             string         `gorm:"type:varchar(45);uniqueIndex;not null" json:"ip"`
	IPv6           string         `gorm:"type:varchar(128)" json:"ipv6"`
	MAC            string         `gorm:"type:varchar(17)" json:"mac"`
	DatacenterID   uint           `gorm:"index" json:"datacenter_id"`
	Rack           string         `gorm:"type:varchar(32)" json:"rack"`
	RackPosition   string         `gorm:"type:varchar(16)" json:"rack_position"`
	CPU            string         `gorm:"type:varchar(128)" json:"cpu"`
	CPUCores       int            `gorm:"default:0" json:"cpu_cores"`
	MemoryMB       int            `gorm:"default:0" json:"memory_mb"`
	DiskType       string         `gorm:"type:varchar(32)" json:"disk_type"` // SSD/HDD/NVMe
	DiskSizeGB     int            `gorm:"default:0" json:"disk_size_gb"`
	BandwidthMbps  int            `gorm:"default:0" json:"bandwidth_mbps"`
	TrafficGB      int            `gorm:"default:0" json:"traffic_gb"` // 月流量GB，0=不限
	OS             string         `gorm:"type:varchar(64)" json:"os"`
	Status         int8           `gorm:"type:smallint;default:0;index" json:"status"` // 0=关机 1=运行中 2=故障 3=维护
	PowerStatus    int8           `gorm:"type:smallint;default:0" json:"power_status"` // 0=off 1=on
	OwnerID        *uint          `gorm:"index" json:"owner_id"`
	AssignedAt     *time.Time     `json:"assigned_at"`
	ExpiredAt      *time.Time     `gorm:"index" json:"expired_at"`
	Remark         string         `gorm:"type:text" json:"remark"`
	Tags           string         `gorm:"type:varchar(256)" json:"tags"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	Datacenter     *DcimDatacenter `gorm:"foreignKey:DatacenterID" json:"datacenter,omitempty"`
	Owner          *User          `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
}

// DcimCloud 云服务器
type DcimCloud struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Name           string         `gorm:"type:varchar(128);not null" json:"name"`
	Hostname       string         `gorm:"type:varchar(128)" json:"hostname"`
	IP             string         `gorm:"type:varchar(45);uniqueIndex;not null" json:"ip"`
	IPv6           string         `gorm:"type:varchar(128)" json:"ipv6"`
	DatacenterID   uint           `gorm:"index" json:"datacenter_id"`
	CPU            int            `gorm:"default:0" json:"cpu"`        // vCPU数
	MemoryMB       int            `gorm:"default:0" json:"memory_mb"`
	DiskSizeGB     int            `gorm:"default:0" json:"disk_size_gb"`
	BandwidthMbps  int            `gorm:"default:0" json:"bandwidth_mbps"`
	TrafficGB      int            `gorm:"default:0" json:"traffic_gb"`
	OS             string         `gorm:"type:varchar(64)" json:"os"`
	VirtualType    string         `gorm:"type:varchar(32)" json:"virtual_type"` // KVM/VMware/LXC
	Status         int8           `gorm:"type:smallint;default:0;index" json:"status"` // 0=关机 1=运行中 2=故障 3=创建中 4=重装中
	PowerStatus    int8           `gorm:"type:smallint;default:0" json:"power_status"`
	OwnerID        *uint          `gorm:"index" json:"owner_id"`
	ParentServerID *uint          `gorm:"index" json:"parent_server_id"` // 宿主机
	PlanID         *uint          `gorm:"index" json:"plan_id"`
	PriceMonthly   float64        `gorm:"type:decimal(10,2);default:0" json:"price_monthly"`
	ExpiredAt      *time.Time     `gorm:"index" json:"expired_at"`
	ProvisionedAt  *time.Time     `json:"provisioned_at"`
	Remark         string         `gorm:"type:text" json:"remark"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	Datacenter     *DcimDatacenter `gorm:"foreignKey:DatacenterID" json:"datacenter,omitempty"`
	Owner          *User          `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	ParentServer   *DcimServer    `gorm:"foreignKey:ParentServerID" json:"parent_server,omitempty"`
}

// DcimOperationLog 服务器操作日志
type DcimOperationLog struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ServerType  string    `gorm:"type:varchar(16);not null;index" json:"server_type"` // physical/cloud
	ServerID    uint      `gorm:"index;not null" json:"server_id"`
	OperatorID  uint      `gorm:"index;not null" json:"operator_id"`
	Action      string    `gorm:"type:varchar(32);not null" json:"action"` // boot/shutdown/reboot/reinstall/renew/migrate
	Params      string    `gorm:"type:text" json:"params"`
	Status      int8      `gorm:"type:smallint;default:1" json:"status"` // 1=执行中 2=成功 3=失败
	Result      string    `gorm:"type:text" json:"result"`
	ErrorMsg    string    `gorm:"type:text" json:"error_msg"`
	CreatedAt   time.Time `json:"created_at"`
	FinishedAt  *time.Time `json:"finished_at"`
}
