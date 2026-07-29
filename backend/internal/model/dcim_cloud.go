package model

import (
	"time"

	"gorm.io/gorm"
)

// DcimCloudServer 魔方云服务器
type DcimCloudServer struct {
	gorm.Model
	Name          string     `gorm:"type:varchar(128);not null" json:"name"`
	Hostname      string     `gorm:"type:varchar(128)" json:"hostname"`
	IP            string     `gorm:"type:varchar(64);not null;index" json:"ip"`
	IPv6          string     `gorm:"type:varchar(128)" json:"ipv6"`
	Username      string     `gorm:"type:varchar(64)" json:"username"`
	Password      string     `gorm:"type:varchar(256)" json:"password"`
	Secure        int8       `gorm:"type:smallint;default:0" json:"secure"` // 0=不安全 1=安全
	Disabled      int8       `gorm:"type:smallint;default:0;index" json:"disabled"`
	UserPrefix    string     `gorm:"type:varchar(32)" json:"user_prefix"`
	AccountType   string     `gorm:"type:varchar(32)" json:"account_type"` // admin/reseller/user
	DatacenterID  uint       `gorm:"index" json:"datacenter_id"`
	CPU           int        `gorm:"default:0" json:"cpu"`
	MemoryMB      int        `gorm:"default:0" json:"memory_mb"`
	DiskSizeGB    int        `gorm:"default:0" json:"disk_size_gb"`
	BandwidthMbps int        `gorm:"default:0" json:"bandwidth_mbps"`
	TrafficGB     int        `gorm:"default:0" json:"traffic_gb"`
	OS            string     `gorm:"type:varchar(64)" json:"os"`
	VirtualType   string     `gorm:"type:varchar(32)" json:"virtual_type"` // KVM/VMware/LXC
	OwnerID       *uint      `gorm:"index" json:"owner_id"`
	Status        int8       `gorm:"type:smallint;default:1;index" json:"status"` // 1=运行中 2=已停止 3=创建中 4=错误
	ExpiredAt     *time.Time `json:"expired_at"`
	Remark        string     `gorm:"type:text" json:"remark"`
	Tags          string     `gorm:"type:varchar(256)" json:"tags"`
	ParentServerID *uint     `gorm:"index" json:"parent_server_id"`
	PriceMonthly  float64    `gorm:"type:decimal(10,2);default:0" json:"price_monthly"`
	LastSyncAt    *time.Time `json:"last_sync_at"`
}

// DcimCloudOperationLog 魔方云操作日志
type DcimCloudOperationLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ServerID   uint      `gorm:"index;not null" json:"server_id"`
	Action     string    `gorm:"type:varchar(32);not null" json:"action"` // boot/shutdown/reboot/reinstall/resize
	OperatorID uint      `gorm:"index" json:"operator_id"`
	Status     int8      `gorm:"type:smallint;not null" json:"status"` // 1=成功 2=失败
	ErrorMsg   string    `gorm:"type:text" json:"error_msg"`
	CreatedAt  time.Time `json:"created_at"`
}
