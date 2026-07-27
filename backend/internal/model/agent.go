package model

import (
	"time"

	"gorm.io/gorm"
)

// Agent 代理
type Agent struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	UserID         uint           `gorm:"uniqueIndex;not null" json:"user_id"`
	AgentNo        string         `gorm:"type:varchar(20);uniqueIndex;not null" json:"agent_no"`
	ParentID       *uint          `gorm:"index" json:"parent_id"`
	Level          int            `gorm:"default:1" json:"level"`
	CommissionRate float64        `gorm:"type:decimal(5,2);default:10" json:"commission_rate"`
	Balance        float64        `gorm:"type:decimal(10,2);default:0" json:"balance"`
	TotalEarned    float64        `gorm:"type:decimal(10,2);default:0" json:"total_earned"`
	Status         int8           `gorm:"type:smallint;default:1" json:"status"` // 1正常 2冻结
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	User           *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Parent         *Agent         `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
}

// AgentCommission 代理佣金记录
type AgentCommission struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AgentID   uint      `gorm:"index;not null" json:"agent_id"`
	OrderID   uint      `gorm:"index" json:"order_id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Amount    float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	Rate      float64   `gorm:"type:decimal(5,2)" json:"rate"`
	Status    int8      `gorm:"type:smallint;default:1" json:"status"` // 1待确认 2已确认 3已拒绝
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
