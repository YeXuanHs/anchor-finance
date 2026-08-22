package model

import (
	"time"

	"gorm.io/gorm"
)

// Blacklist 黑名单模型
type Blacklist struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Type      string         `gorm:"size:20;not null;index" json:"type"` // ip, email, phone, username
	Value     string         `gorm:"size:200;not null;uniqueIndex" json:"value"`
	Reason    string         `gorm:"size:500" json:"reason"`
	AdminID   uint           `gorm:"index" json:"admin_id"`
	ExpiresAt *time.Time     `gorm:"index" json:"expires_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Blacklist) TableName() string {
	return "blacklists"
}
