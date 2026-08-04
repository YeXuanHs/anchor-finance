package model

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"gorm.io/gorm"
)

// UserAPIKey represents a user's API key for programmatic access.
type UserAPIKey struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Key       string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"-"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	LastUsed  *time.Time     `json:"last_used"`
}

// GenerateAPIKey generates a random hex API key.
func GenerateAPIKey() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
