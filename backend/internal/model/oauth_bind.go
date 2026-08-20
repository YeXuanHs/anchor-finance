package model

import (
	"time"

	"gorm.io/gorm"
)

// OAuthBind OAuth绑定
type OAuthBind struct {
	gorm.Model
	UserID       uint       `gorm:"uniqueIndex:idx_user_provider;not null" json:"user_id"`
	Provider     string     `gorm:"type:varchar(32);uniqueIndex:idx_user_provider;not null" json:"provider"`
	ProviderUID  string     `gorm:"type:varchar(128);not null" json:"provider_uid"`
	AccessToken  string     `gorm:"type:varchar(512)" json:"-"`
	RefreshToken string     `gorm:"type:varchar(512)" json:"-"`
	ExpiresAt    *time.Time `json:"expires_at"`
	Nickname     string     `gorm:"type:varchar(64)" json:"nickname"`
	Avatar       string     `gorm:"type:varchar(256)" json:"avatar"`
	Extra        string     `gorm:"type:json" json:"extra"`
}

// TableName overrides the table name.
func (OAuthBind) TableName() string {
	return "oauth_binds"
}
