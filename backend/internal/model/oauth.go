package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// OAuthAccount stores third-party OAuth binding information.
type OAuthAccount struct {
	gorm.Model
	UserID   uint           `gorm:"index;not null" json:"user_id"`
	User     User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Provider string         `gorm:"type:varchar(32);index;not null" json:"provider"`
	OpenID   string         `gorm:"type:varchar(128);not null" json:"openid"`
	UnionID  string         `gorm:"type:varchar(128);index" json:"unionid"`
	Username string         `gorm:"type:varchar(128)" json:"username"`
	Email    string         `gorm:"type:varchar(255)" json:"email"`
	Avatar   string         `gorm:"type:varchar(512)" json:"avatar"`
	RawData  datatypes.JSON `gorm:"type:json" json:"raw_data"`
}

// TableName overrides the table name.
func (OAuthAccount) TableName() string {
	return "oauth_accounts"
}

// OAuthState stores temporary OAuth state for CSRF protection.
type OAuthState struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	State     string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"state"`
	Redirect  string    `gorm:"type:varchar(512)" json:"redirect"`
	ExpiresAt time.Time `gorm:"index;not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName overrides the table name.
func (OAuthState) TableName() string {
	return "oauth_states"
}
