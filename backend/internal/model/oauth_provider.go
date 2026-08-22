package model

import (
	"time"

	"gorm.io/gorm"
)

// OAuthProvider 第三方登录提供商模型
type OAuthProvider struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Code      string         `gorm:"size:50;uniqueIndex;not null" json:"code"`
	Icon      string         `gorm:"size:500" json:"icon"`
	Config    string         `gorm:"type:text" json:"config"` // JSON格式的配置
	Status    string         `gorm:"size:20;default:active" json:"status"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (OAuthProvider) TableName() string {
	return "oauth_providers"
}
