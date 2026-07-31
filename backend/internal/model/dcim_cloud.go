package model

import (
	"time"

	"gorm.io/datatypes"
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

// CloudNATRule 云服务器NAT端口转发规则
type CloudNATRule struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CloudID   uint      `gorm:"index" json:"cloud_id"`
	Name      string    `gorm:"type:varchar(256)" json:"name"`
	Protocol  string    `gorm:"type:varchar(16)" json:"protocol"` // tcp/udp
	ExtPort   int       `gorm:"not null" json:"ext_port"`
	IntPort   int       `gorm:"not null" json:"int_port"`
	IntIP     string    `gorm:"type:varchar(45)" json:"int_ip"`
	Status    string    `gorm:"type:varchar(32)" json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// CloudSecurityGroup 云服务器安全组
type CloudSecurityGroup struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	CloudID       uint           `gorm:"index" json:"cloud_id"`
	Name          string         `gorm:"type:varchar(256)" json:"name"`
	Rules         datatypes.JSON `gorm:"type:json" json:"rules"`
	DefaultAction string         `gorm:"type:varchar(16)" json:"default_action"` // accept/drop
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// CloudSecurityGroupRule 云服务器安全组规则
type CloudSecurityGroupRule struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	GroupID   uint      `gorm:"index" json:"group_id"`
	Direction string    `gorm:"type:varchar(16)" json:"direction"` // in/out
	Protocol  string    `gorm:"type:varchar(16)" json:"protocol"`
	PortRange string    `gorm:"type:varchar(32)" json:"port_range"`
	Source    string    `gorm:"type:varchar(128)" json:"source"`
	Action    string    `gorm:"type:varchar(16)" json:"action"` // accept/drop
	Priority  int       `gorm:"default:100" json:"priority"`
	CreatedAt time.Time `json:"created_at"`
}

// CloudISO 可用ISO镜像
type CloudISO struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Name   string `gorm:"type:varchar(256)" json:"name"`
	SizeMB int    `json:"size_mb"`
	URL    string `gorm:"type:varchar(512)" json:"url"`
	Status string `gorm:"type:varchar(32)" json:"status"` // available/mounted
}

// CloudFlowPacket 云服务器流量包
type CloudFlowPacket struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CloudID   uint      `gorm:"index" json:"cloud_id"`
	Name      string    `gorm:"type:varchar(256)" json:"name"`
	SizeGB    int       `json:"size_gb"`
	UsedGB    int       `json:"used_gb"`
	ExpiredAt time.Time `json:"expired_at"`
	Status    string    `gorm:"type:varchar(32)" json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// CloudChart 云服务器监控数据点
type CloudChart struct {
	Timestamp  time.Time `gorm:"primaryKey" json:"timestamp"`
	CloudID    uint      `gorm:"primaryKey;index" json:"cloud_id"`
	CPURate    float64   `json:"cpu_rate"`
	MemoryRate float64   `json:"memory_rate"`
	DiskRate   float64   `json:"disk_rate"`
	NetIn      int64     `json:"net_in"`
	NetOut     int64     `json:"net_out"`
}
