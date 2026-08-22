package model

import (
	"time"
)

// Captcha 验证码模型
type Captcha struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Target    string    `gorm:"size:100;not null;index" json:"target"` // 手机号或邮箱
	Code      string    `gorm:"size:10;not null" json:"code"`
	Type      string    `gorm:"size:20;not null" json:"type"` // register, login, reset_password
	Used      bool      `gorm:"default:false" json:"used"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (Captcha) TableName() string {
	return "captchas"
}
