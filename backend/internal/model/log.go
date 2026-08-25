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

// SecurityLog 安全审计日志模型（MD 9.1 功能9：安全审计日志）
// 记录所有安全事件：0元购/SQL注入/暴力破解/越权等
type SecurityLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	AttackType string    `gorm:"size:50;index" json:"attack_type"`   // zero_price/sql_inject/brute_force/idor/csrf/registration_abuse
	UserID     uint      `gorm:"index" json:"user_id"`
	Username   string    `gorm:"size:100" json:"username"`
	IP         string    `gorm:"size:45;index" json:"ip"`
	RealIP     string    `gorm:"size:45" json:"real_ip"`             // X-Forwarded-For
	Path       string    `gorm:"size:500" json:"path"`
	Method     string    `gorm:"size:10" json:"method"`
	UserAgent  string    `gorm:"size:500" json:"user_agent"`
	SessionID  string    `gorm:"size:100" json:"session_id"`
	Referer    string    `gorm:"size:500" json:"referer"`
	Params     string    `gorm:"type:text" json:"params"`            // 脱敏后的请求参数
	Detail     string    `gorm:"type:text" json:"detail"`
	CreatedAt  time.Time `gorm:"index" json:"created_at"`
}

// TableName 指定表名
func (SecurityLog) TableName() string {
	return "security_logs"
}
