package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// SystemConfig 系统配置
type SystemConfig struct {
	gorm.Model
	Key         string `gorm:"type:varchar(128);uniqueIndex;not null" json:"key"`
	Value       string `gorm:"type:text" json:"value"`
	Type        string `gorm:"type:varchar(16);default:'string'" json:"type"` // string/int/bool/json/array
	Group       string `gorm:"type:varchar(64);index;not null" json:"group"`
	Name        string `gorm:"type:varchar(128);not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	SortOrder   int    `gorm:"default:0" json:"sort_order"`
	IsPublic    bool   `gorm:"default:false" json:"is_public"` // 是否前端可见
	IsReadOnly  bool   `gorm:"default:false" json:"is_read_only"`
}

// Announcement 公告
type Announcement struct {
	gorm.Model
	Title       string         `gorm:"type:varchar(256);not null" json:"title"`
	Slug        string         `gorm:"type:varchar(256);uniqueIndex" json:"slug"`
	Content     string         `gorm:"type:text;not null" json:"content"`
	Summary     string         `gorm:"type:varchar(512)" json:"summary"`
	Category    string         `gorm:"type:varchar(64);index" json:"category"`
	Author      string         `gorm:"type:varchar(64)" json:"author"`
	Image       string         `gorm:"type:varchar(512)" json:"image"`
	IsPinned    bool           `gorm:"default:false;index" json:"is_pinned"`
	IsPublished bool           `gorm:"default:false;index" json:"is_published"`
	PublishedAt *time.Time     `gorm:"index" json:"published_at"`
	ViewCount   int            `gorm:"default:0" json:"view_count"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	Language    string         `gorm:"type:varchar(16);default:'zh-CN'" json:"language"`
	Metadata    datatypes.JSON `gorm:"type:jsonb" json:"metadata"`
}

// Notification 站内通知
type Notification struct {
	gorm.Model
	UserID      uint   `gorm:"index;not null" json:"user_id"`
	User        User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Type        string `gorm:"type:varchar(64);not null;index" json:"type"` // system/order/ticket/payment/security/promotion
	Title       string `gorm:"type:varchar(256);not null" json:"title"`
	Content     string `gorm:"type:text" json:"content"`
	Icon        string `gorm:"type:varchar(128)" json:"icon"`
	ActionURL   string `gorm:"type:varchar(512)" json:"action_url"`
	RelType     string `gorm:"type:varchar(32)" json:"rel_type"`
	RelID       uint   `gorm:"index" json:"rel_id"`
	IsRead      bool   `gorm:"default:false;index" json:"is_read"`
	ReadAt      *time.Time `json:"read_at"`
	Channel     string `gorm:"type:varchar(32);default:'web'" json:"channel"` // web/email/sms/push
	IsSent      bool   `gorm:"default:false" json:"is_sent"`
	SentAt      *time.Time `json:"sent_at"`
	Metadata    datatypes.JSON `gorm:"type:jsonb" json:"metadata"`
}

// EmailTemplate moved to internal/model/email_template.go

// PaymentGateway 支付网关
type PaymentGateway struct {
	gorm.Model
	Name        string         `gorm:"type:varchar(64);not null" json:"name"`
	Code        string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"code"`
	Description string         `gorm:"type:text" json:"description"`
	Icon        string         `gorm:"type:varchar(512)" json:"icon"`
	Config      datatypes.JSON `gorm:"type:jsonb;not null" json:"config"` // 网关配置（加密存储）
	Currencies  datatypes.JSON `gorm:"type:jsonb" json:"currencies"`       // 支持的货币列表
	MinAmount   datatypes.Decimal `gorm:"type:decimal(20,4);default:0" json:"min_amount"`
	MaxAmount   datatypes.Decimal `gorm:"type:decimal(20,4)" json:"max_amount"`
	FixedFee    datatypes.Decimal `gorm:"type:decimal(20,4);default:0" json:"fixed_fee"`
	PercentFee  datatypes.Decimal `gorm:"type:decimal(5,4);default:0" json:"percent_fee"`
	FeeCurrency string         `gorm:"type:varchar(8)" json:"fee_currency"`
	IsOnline    bool           `gorm:"default:true" json:"is_online"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	Status      int16          `gorm:"type:smallint;default:1;not null;index" json:"status"` // 1=启用 0=禁用
	TestMode    bool           `gorm:"default:false" json:"test_mode"`
	SupportRefund    bool      `gorm:"default:false" json:"support_refund"`
	SupportRecurring bool      `gorm:"default:false" json:"support_recurring"`
	WebhookURL  string         `gorm:"type:varchar(512)" json:"webhook_url"`
	ReturnURL   string         `gorm:"type:varchar(512)" json:"return_url"`
	CancelURL   string         `gorm:"type:varchar(512)" json:"cancel_url"`
	NotifyURL   string         `gorm:"type:varchar(512)" json:"notify_url"`
}

// Admin 管理员
type Admin struct {
	gorm.Model
	Username     string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"username"`
	Email        string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"type:varchar(255);not null" json:"-"`
	Salt         string         `gorm:"type:varchar(32);not null" json:"-"`
	Name         string         `gorm:"type:varchar(64)" json:"name"`
	Avatar       string         `gorm:"type:varchar(512)" json:"avatar"`
	Role         string         `gorm:"type:varchar(32);not null;default:'staff'" json:"role"` // super_admin/admin/staff
	Permissions  datatypes.JSON `gorm:"type:jsonb" json:"permissions"` // 细粒度权限
	DepartmentID *uint          `gorm:"index" json:"department_id"`
	Department   *Department    `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	LastLoginAt  *time.Time     `json:"last_login_at"`
	LastLoginIP  string         `gorm:"type:varchar(64)" json:"last_login_ip"`
	TwoFactorKey string         `gorm:"type:varchar(128)" json:"-"`
	Status       int16          `gorm:"type:smallint;default:1;not null" json:"status"` // 1=启用 0=禁用
	IsSuperAdmin bool           `gorm:"default:false" json:"is_super_admin"`
	Logs         []AdminLog     `gorm:"foreignKey:AdminID" json:"logs,omitempty"`
}

// AdminLog 管理员操作日志
type AdminLog struct {
	gorm.Model
	AdminID   uint           `gorm:"index;not null" json:"admin_id"`
	Admin     Admin          `gorm:"foreignKey:AdminID" json:"admin,omitempty"`
	Action    string         `gorm:"type:varchar(64);not null;index" json:"action"` // create/update/delete/login/export/config
	Module    string         `gorm:"type:varchar(64);not null;index" json:"module"` // user/order/product/ticket/system/payment
	TargetID  uint           `gorm:"index" json:"target_id"`
	TargetType string        `gorm:"type:varchar(64)" json:"target_type"`
	OldData   datatypes.JSON `gorm:"type:jsonb" json:"old_data"`
	NewData   datatypes.JSON `gorm:"type:jsonb" json:"new_data"`
	IPAddress string         `gorm:"type:varchar(64);not null" json:"ip_address"`
	UserAgent string         `gorm:"type:varchar(512)" json:"user_agent"`
	Remark    string         `gorm:"type:text" json:"remark"`
	Status    int16          `gorm:"type:smallint;default:1" json:"status"` // 1=成功 0=失败
}
