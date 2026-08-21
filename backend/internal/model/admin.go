package model

import (
	"time"

	"gorm.io/gorm"
)

// Admin 管理员模型
type Admin struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email        string         `gorm:"size:100;uniqueIndex" json:"email"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	RealName     string         `gorm:"size:50" json:"real_name"`
	RoleID       uint           `gorm:"not null" json:"role_id"`
	Status       string         `gorm:"size:20;default:active" json:"status"` // active, disabled
	LastLoginAt  *time.Time     `json:"last_login_at"`
	LastLoginIP  string         `gorm:"size:45" json:"last_login_ip"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Admin) TableName() string {
	return "admins"
}

// Role 角色模型
type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Description string         `gorm:"size:255" json:"description"`
	Permissions string         `gorm:"type:json" json:"permissions"` // JSON格式的权限列表
	IsSuper     bool           `gorm:"default:false" json:"is_super"` // 是否超级管理员
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}
