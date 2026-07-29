package model

import (
	"time"

	"gorm.io/gorm"
)

// Blacklist 黑名单
type Blacklist struct {
	gorm.Model
	Type      string     `gorm:"type:varchar(20);not null;index" json:"type"` // ip, email, phone, domain
	Value     string     `gorm:"type:varchar(255);not null;uniqueIndex" json:"value"`
	Reason    string     `gorm:"type:text" json:"reason"`
	AdminID   uint       `gorm:"index" json:"admin_id"`
	ExpiresAt *time.Time `gorm:"index" json:"expires_at"`
}

// TableName 指定表名
func (Blacklist) TableName() string {
	return "blacklists"
}
