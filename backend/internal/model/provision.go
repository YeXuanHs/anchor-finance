package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ProvisionModule 模块供给配置
type ProvisionModule struct {
	gorm.Model
	Name               string         `gorm:"type:varchar(128);not null;index" json:"name"`
	Slug               string         `gorm:"type:varchar(128);uniqueIndex" json:"slug"`
	Description        string         `gorm:"type:text" json:"description"`
	Type               string         `gorm:"type:varchar(32);not null;index" json:"type"` // hosting/domain/ssl/vpn/cdn/server/other
	SupportedProducts  datatypes.JSON `gorm:"type:json" json:"supported_products"`
	Config             datatypes.JSON `gorm:"type:json" json:"config"`
	ServerURL          string         `gorm:"type:varchar(512)" json:"server_url"`
	ServerIP           string         `gorm:"type:varchar(64)" json:"server_ip"`
	APIKey             string         `gorm:"type:varchar(256)" json:"-"`
	APISecret          string         `gorm:"type:varchar(256)" json:"-"`
	Username           string         `gorm:"type:varchar(128)" json:"-"`
	Password           string         `gorm:"type:varchar(256)" json:"-"`
	Hash               string         `gorm:"type:varchar(128)" json:"hash"`
	Active             bool           `gorm:"default:true;index" json:"active"`
	Priority           int            `gorm:"default:0;index" json:"priority"`
	Weight             int            `gorm:"default:1" json:"weight"`
	MaxRetries         int            `gorm:"default:3" json:"max_retries"`
	Timeout            int            `gorm:"default:30;comment:seconds" json:"timeout"`
	LastTestAt         *time.Time     `json:"last_test_at"`
	LastTestOK         bool           `gorm:"default:false" json:"last_test_ok"`
	LastError          string         `gorm:"type:text" json:"last_error"`
	ProvisionCount     int            `gorm:"default:0" json:"provision_count"`
	SuccessCount       int            `gorm:"default:0" json:"success_count"`
	FailCount          int            `gorm:"default:0" json:"fail_count"`
	Metadata           datatypes.JSON `gorm:"type:json" json:"metadata"`
	UpdatedAt          time.Time      `json:"updated_at"`
}

// ProvisionLog 模块操作日志
type ProvisionLog struct {
	ID         uint             `gorm:"primaryKey" json:"id"`
	ModuleID   uint             `gorm:"index;not null" json:"module_id"`
	Module     ProvisionModule  `gorm:"foreignKey:ModuleID" json:"module,omitempty"`
	Action     string           `gorm:"type:varchar(32);not null" json:"action"` // create/suspend/unsuspend/terminate/test
	Status     int8             `gorm:"type:smallint;default:1;not null" json:"status"` // 1=处理中 2=成功 3=失败
	Request    datatypes.JSON   `gorm:"type:json" json:"request"`
	Response   datatypes.JSON   `gorm:"type:json" json:"response"`
	Error      string           `gorm:"type:text" json:"error"`
	Duration   int              `gorm:"default:0;comment:ms" json:"duration"`
	AdminID    uint             `gorm:"index" json:"admin_id"`
	CreatedAt  time.Time        `json:"created_at"`
}
