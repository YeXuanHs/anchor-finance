package model

import (
	"time"

	"gorm.io/datatypes"
)

// AggregateLoginProvider 聚合登录提供商配置
type AggregateLoginProvider struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(50);not null" json:"name"`        // 显示名称，如"QQ登录"
	Code      string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"` // 唯一标识，如 qq, wx, alipay, github
	Type      string         `gorm:"type:varchar(20);not null" json:"type"`        // 登录类型：oauth, juhe
	APIURL    string         `gorm:"type:varchar(500)" json:"api_url"`             // 聚合登录API地址
	AppID     string         `gorm:"type:varchar(100)" json:"app_id"`             // 应用ID
	AppKey    string         `gorm:"type:varchar(255)" json:"app_key"`            // 应用密钥
	Config    datatypes.JSON `gorm:"type:jsonb" json:"config"`                    // 扩展配置
	IsActive  bool           `gorm:"default:true" json:"is_active"`               // 是否启用
	SortOrder int            `gorm:"default:0" json:"sort_order"`                 // 排序权重
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// TableName overrides the table name.
func (AggregateLoginProvider) TableName() string {
	return "aggregate_login_providers"
}

// AggregateLoginLog 聚合登录日志
type AggregateLoginLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`                   // 关联用户ID（登录成功后填充）
	ProviderID uint      `gorm:"index" json:"provider_id"`               // 关联提供商ID
	Provider   AggregateLoginProvider `gorm:"foreignKey:ProviderID" json:"provider,omitempty"`
	SocialUID  string    `gorm:"type:varchar(128)" json:"social_uid"`    // 第三方用户ID
	Nickname   string    `gorm:"type:varchar(128)" json:"nickname"`      // 第三方昵称
	Avatar     string    `gorm:"type:varchar(512)" json:"avatar"`        // 第三方头像
	IP         string    `gorm:"type:varchar(64)" json:"ip"`             // 登录IP
	UserAgent  string    `gorm:"type:varchar(512)" json:"user_agent"`    // 浏览器UA
	Status     int16     `gorm:"type:smallint;default:1;not null" json:"status"` // 1=成功 0=失败
	Remark     string    `gorm:"type:varchar(256)" json:"remark"`        // 备注/错误信息
	CreatedAt  time.Time `json:"created_at"`
}

// TableName overrides the table name.
func (AggregateLoginLog) TableName() string {
	return "aggregate_login_logs"
}
