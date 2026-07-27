package model

import (
	"time"

	"gorm.io/gorm"
)

// Role 角色
type Role struct {
	gorm.Model
	Name        string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	IsSystem    bool           `gorm:"default:false" json:"is_system"` // 系统内置角色不可删
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	Permissions []Permission   `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	Users       []User         `gorm:"many2many:user_roles;" json:"users,omitempty"`
}

// Permission 权限
type Permission struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Code      string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"code"` // users.view, orders.edit, etc.
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Module    string    `gorm:"type:varchar(50);index" json:"module"` // users, orders, products, etc.
	Type      string    `gorm:"type:varchar(20)" json:"type"`        // view, create, edit, delete
	CreatedAt time.Time `json:"created_at"`
}

// UserRole 用户角色关联
type UserRole struct {
	UserID    uint      `gorm:"primaryKey" json:"user_id"`
	RoleID    uint      `gorm:"primaryKey" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

// RolePermission 角色权限关联
type RolePermission struct {
	RoleID       uint `gorm:"primaryKey" json:"role_id"`
	PermissionID uint `gorm:"primaryKey" json:"permission_id"`
}
