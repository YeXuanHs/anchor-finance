package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// OAuthProvider OAuth提供商配置
type OAuthProvider struct {
	gorm.Model
	Name        string         `gorm:"type:varchar(32);not null;uniqueIndex:idx_oauth_provider_name" json:"name"` // wechat/qq/github/google/alipay/weibo/dingtalk
	DisplayName string         `gorm:"type:varchar(64);not null" json:"display_name"`
	Icon        string         `gorm:"type:varchar(512)" json:"icon"`
	ClientID    string         `gorm:"type:varchar(256);not null" json:"client_id"`
	ClientSecret string        `gorm:"type:varchar(256);not null" json:"client_secret"`
	RedirectURL string         `gorm:"type:varchar(512)" json:"redirect_url"`
	AuthURL     string         `gorm:"type:varchar(512)" json:"auth_url"`
	TokenURL    string         `gorm:"type:varchar(512)" json:"token_url"`
	UserInfoURL string         `gorm:"type:varchar(512)" json:"user_info_url"`
	Scopes      string         `gorm:"type:varchar(512)" json:"scopes"`
	Extra       datatypes.JSON `gorm:"type:jsonb" json:"extra"` // 额外配置参数
	IsEnabled   bool           `gorm:"default:true" json:"is_enabled"`
	SortOrder   int            `gorm:"default:0;index" json:"sort_order"`
	Status      int16          `gorm:"type:smallint;default:1;not null;index" json:"status"` // 1=启用 0=禁用
	UserCount   int            `gorm:"default:0" json:"user_count"` // 绑定用户数
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
