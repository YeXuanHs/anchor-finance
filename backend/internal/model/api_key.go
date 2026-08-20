package model

import (
	"time"
)

// APIKey represents an API key for external API access.
type APIKey struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(128);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Secret      string    `gorm:"type:varchar(128);uniqueIndex;not null" json:"-"`
	Status      int16     `gorm:"type:smallint;default:1;index" json:"status"` // 1=active 0=disabled
	RateLimit   int       `gorm:"default:0" json:"rate_limit"` // 0=unlimited
	Permissions string    `gorm:"type:text" json:"permissions"` // JSON array
	UserID      uint      `gorm:"index" json:"user_id"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (APIKey) TableName() string { return "api_keys" }
