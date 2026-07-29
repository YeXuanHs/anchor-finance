package model

import (
	"time"

	"gorm.io/gorm"
)

// Captcha 验证码
type Captcha struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CaptchaID string    `gorm:"type:varchar(128);uniqueIndex;not null" json:"captcha_id"`
	Code      string    `gorm:"type:varchar(16);not null" json:"code"`
	IP        string    `gorm:"type:varchar(64);index" json:"ip"`
	Used      bool      `gorm:"default:false" json:"used"`
	ExpiredAt time.Time `gorm:"index" json:"expired_at"`
	CreatedAt time.Time `json:"created_at"`
}

// PublicConfig 公共配置
type PublicConfig struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Key   string `gorm:"type:varchar(128);uniqueIndex;not null" json:"key"`
	Value string `gorm:"type:text" json:"value"`
	Group string `gorm:"type:varchar(64);default:general" json:"group"`
}
