package model

import (
	"time"

	"gorm.io/gorm"
)

// TwoFactorConfig 二次验证配置模型
type TwoFactorConfig struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Key       string         `gorm:"size:100;uniqueIndex;not null" json:"key"`
	Value     string         `gorm:"type:text" json:"value"`
	Group     string         `gorm:"size:50;default:two_factor" json:"group"`
	Name      string         `gorm:"size:100" json:"name"`
	Type      string         `gorm:"size:20;default:text" json:"type"` // text, number, select, switch
	Options   string         `gorm:"type:text" json:"options"`         // JSON格式的选项
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (TwoFactorConfig) TableName() string {
	return "two_factor_configs"
}

// UserTwoFactor 用户二次验证模型
type UserTwoFactor struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"uniqueIndex;not null" json:"user_id"`
	Type      string         `gorm:"size:20;not null" json:"type"` // sms, email, totp
	Secret    string         `gorm:"size:200" json:"secret"`
	Enabled   bool           `gorm:"default:false" json:"enabled"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (UserTwoFactor) TableName() string {
	return "user_two_factors"
}
