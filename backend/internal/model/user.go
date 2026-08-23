package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Username       string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email          string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
	PasswordHash   string         `gorm:"size:255;not null" json:"-"`
	Phone          string         `gorm:"size:20" json:"phone"`
	Company        string         `gorm:"size:100" json:"company"`
	Status         string         `gorm:"size:20;default:active" json:"status"` // active, suspended, closed
	Balance        float64        `gorm:"type:decimal(10,2);default:0" json:"balance"`
	CreditLimit    float64        `gorm:"type:decimal(10,2);default:0" json:"credit_limit"`
	GroupID        uint           `gorm:"default:0" json:"group_id"`
	LevelID        uint           `gorm:"default:0" json:"level_id"`
	Avatar         string         `gorm:"size:255" json:"avatar"`
	RealName       string         `gorm:"size:50" json:"real_name"`
	IDCard         string         `gorm:"size:20" json:"-"`
	IsVerified     bool           `gorm:"default:false" json:"is_verified"`
	LoginFailCount int            `gorm:"default:0" json:"-"`
	LockedUntil    *time.Time     `json:"-"`
	LastLoginAt    *time.Time     `json:"last_login_at"`
	LastLoginIP    string         `gorm:"size:45" json:"last_login_ip"`
	APIKey         string         `gorm:"size:64" json:"-"`           // API密钥（AES加密存储）
	APIEnabled     bool           `gorm:"default:false" json:"api_enabled"` // API是否开通
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// 可以在这里添加创建前的逻辑
	return nil
}
