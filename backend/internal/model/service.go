package model

import (
	"time"

	"gorm.io/gorm"
)

// Service 服务模型
type Service struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	UserID         uint           `gorm:"index;not null" json:"user_id"`
	ProductID      uint           `gorm:"index" json:"product_id"`
	ProductTypeID  uint           `gorm:"index" json:"product_type_id"`
	ServerID       uint           `gorm:"index" json:"server_id"`       // 关联servers表（DCIM硬件）
	ProductName    string         `gorm:"size:100" json:"product_name"`
	Domain         string         `gorm:"size:100" json:"domain"`
	Username       string         `gorm:"size:50" json:"username"`
	PasswordHash   string         `gorm:"size:100" json:"-"`
	Config         string         `gorm:"type:text" json:"config"` // JSON配置
	Status         string         `gorm:"size:20;default:active" json:"status"` // active, suspended, terminated, pending
	BillingCycle   string         `gorm:"size:20" json:"billing_cycle"`
	Amount         float64        `gorm:"type:decimal(10,2)" json:"amount"`
	Remark         string         `gorm:"type:text" json:"remark"`
	AutoRenew      bool           `gorm:"default:false" json:"auto_renew"`
	NextDueDate    *time.Time     `json:"next_due_date"`
	ActivatedAt    *time.Time     `json:"activated_at"`
	TerminatedAt   *time.Time     `json:"terminated_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Service) TableName() string {
	return "services"
}
