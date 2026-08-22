package model

import "time"

// TokenBlacklist JWT黑名单（登出后token失效）
type TokenBlacklist struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TokenHash string    `gorm:"size:64;uniqueIndex;not null" json:"token_hash"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (TokenBlacklist) TableName() string {
	return "token_blacklist"
}
