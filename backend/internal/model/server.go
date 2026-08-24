package model

import (
	"time"

	"gorm.io/gorm"
)

// Server 服务器/硬件连接模型（DCIM IPMI连接信息）
type Server struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Name       string         `gorm:"size:100;not null" json:"name"`
	Hostname   string         `gorm:"size:255;not null" json:"hostname"`            // IPMI IP/域名
	Username   string         `gorm:"size:100" json:"username"`                     // IPMI用户名
	Password   string         `gorm:"size:255" json:"-"`                            // IPMI密码（AES加密）
	Port       int            `gorm:"default:443" json:"port"`                      // 端口
	Secure     bool           `gorm:"default:true" json:"secure"`                   // 是否HTTPS
	AccessHash string         `gorm:"size:255" json:"-"`                            // 访问哈希
	ServerType string         `gorm:"size:20;default:dcim" json:"server_type"`      // dcim/normal
	GroupID    uint           `gorm:"index" json:"group_id"`                        // 分组ID
	Disabled   bool           `gorm:"default:false" json:"disabled"`
	LinkStatus bool           `gorm:"default:false" json:"link_status"`             // 连接状态
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Server) TableName() string { return "servers" }

// DcimServer DCIM服务器扩展配置
type DcimServer struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	ServerID       uint           `gorm:"uniqueIndex" json:"server_id"`
	Auth           string         `gorm:"type:text" json:"auth"`
	Area           string         `gorm:"size:100" json:"area"`
	BillType       string         `gorm:"size:20;default:month" json:"bill_type"`
	FlowRemind     string         `gorm:"size:255" json:"flow_remind"`
	ReinstallTimes int            `gorm:"default:3" json:"reinstall_times"`
	BuyTimes       int            `gorm:"default:1" json:"buy_times"`
	ReinstallPrice float64        `gorm:"type:decimal(10,2);default:0" json:"reinstall_price"`
	APIStatus      int            `gorm:"default:0" json:"api_status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (DcimServer) TableName() string { return "dcim_servers" }
