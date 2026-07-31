package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// MessageConfig 消息通道配置
type MessageConfig struct {
	gorm.Model
	Channel     string         `gorm:"type:varchar(32);not null;uniqueIndex:idx_msg_channel" json:"channel"` // email/sms/site
	Name        string         `gorm:"type:varchar(128);not null" json:"name"`
	Provider    string         `gorm:"type:varchar(64)" json:"provider"` // smtp/aliyun_sms/tencent_sms/custom
	Config      datatypes.JSON `gorm:"type:json;not null" json:"config"` // 通道配置（加密存储）
	SenderName  string         `gorm:"type:varchar(128)" json:"sender_name"`
	SenderAddr  string         `gorm:"type:varchar(255)" json:"sender_addr"`
	Signature   string         `gorm:"type:varchar(128)" json:"signature"` // 短信签名
	RateLimit   int            `gorm:"default:0" json:"rate_limit"`       // 每分钟发送限制, 0=不限
	DailyLimit  int            `gorm:"default:0" json:"daily_limit"`      // 每日发送限制, 0=不限
	IsEnabled   bool           `gorm:"default:true" json:"is_enabled"`
	TestAddress string         `gorm:"type:varchar(255)" json:"test_address"` // 测试接收地址
	LastTestAt  *time.Time     `json:"last_test_at"`
	LastTestOK  bool           `gorm:"default:false" json:"last_test_ok"`
	Status      int16          `gorm:"type:smallint;default:1;not null;index" json:"status"` // 1=启用 0=禁用
	Remark      string         `gorm:"type:text" json:"remark"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
