package model

import (
	"time"

	"gorm.io/gorm"
)

// Admin represents an administrator user.
type Admin struct {
	gorm.Model
	Username  string     `gorm:"type:varchar(64);uniqueIndex;not null" json:"username"`
	Password  string     `gorm:"type:varchar(255);not null" json:"-"`
	Email     string     `gorm:"type:varchar(128)" json:"email"`
	RealName  string     `gorm:"type:varchar(64)" json:"real_name"`
	RoleID    uint       `gorm:"index" json:"role_id"`
	Status    int        `gorm:"default:1" json:"status"` // 1=active, 0=disabled
	LastLogin *time.Time `json:"last_login"`
	Avatar    string     `gorm:"type:varchar(255)" json:"avatar"`
}
