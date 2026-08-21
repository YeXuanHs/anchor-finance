package model

import (
	"time"

	"gorm.io/gorm"
)

// Staff 员工模型
type Staff struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email        string         `gorm:"size:100;uniqueIndex" json:"email"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	RealName     string         `gorm:"size:50" json:"real_name"`
	Phone        string         `gorm:"size:20" json:"phone"`
	RoleID       uint           `gorm:"index" json:"role_id"`
	Status       string         `gorm:"size:20;default:active" json:"status"` // active, disabled
	LastLoginAt  *time.Time     `json:"last_login_at"`
	LastLoginIP  string         `gorm:"size:45" json:"last_login_ip"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Staff) TableName() string {
	return "staff"
}
