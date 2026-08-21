package model

import (
	"time"

	"gorm.io/gorm"
)

// SystemLog 系统日志模型
type SystemLog struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Type      string         `gorm:"size:50;index" json:"type"`
	Content   string         `gorm:"type:text" json:"content"`
	UserID    uint           `gorm:"index" json:"user_id"`
	Username  string         `gorm:"size:50" json:"username"`
	IP        string         `gorm:"size:45" json:"ip"`
	UserAgent string         `gorm:"size:500" json:"user_agent"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (SystemLog) TableName() string {
	return "system_logs"
}

// OperationLog 操作日志模型
type OperationLog struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index" json:"user_id"`
	Username  string         `gorm:"size:50" json:"username"`
	Action    string         `gorm:"size:50" json:"action"`
	Resource  string         `gorm:"size:50" json:"resource"`
	ResourceID uint          `gorm:"index" json:"resource_id"`
	Detail    string         `gorm:"type:text" json:"detail"`
	IP        string         `gorm:"size:45" json:"ip"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (OperationLog) TableName() string {
	return "operation_logs"
}

// LoginLog 登录日志模型
type LoginLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Username  string    `gorm:"size:50" json:"username"`
	IsAdmin   bool      `gorm:"default:false" json:"is_admin"`
	IP        string    `gorm:"size:45" json:"ip"`
	UserAgent string    `gorm:"size:500" json:"user_agent"`
	Status    string    `gorm:"size:20" json:"status"` // success, failed
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (LoginLog) TableName() string {
	return "login_logs"
}
