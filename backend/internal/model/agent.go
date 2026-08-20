package model

import (
	"time"

	"gorm.io/datatypes"
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

// AgentResource 代理资源
type AgentResource struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AgentID   uint      `gorm:"index" json:"agent_id"`
	Type      string    `gorm:"size:32" json:"type"` // product/host/domain
	Name      string    `gorm:"size:256" json:"name"`
	Price     float64   `gorm:"type:decimal(20,4)" json:"price"`
	Stock     int       `gorm:"default:0" json:"stock"`
	Status    string    `gorm:"size:32" json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// AgentInspection 代理巡检
type AgentInspection struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	AgentID     uint           `gorm:"index" json:"agent_id"`
	HostID      uint           `gorm:"index" json:"host_id"`
	Host        *Host          `gorm:"foreignKey:HostID" json:"host,omitempty"`
	Type        string         `gorm:"size:32" json:"type"` // routine/emergency
	Status      string         `gorm:"size:32" json:"status"` // pending/running/completed
	Result      string         `gorm:"type:text" json:"result"`
	Images      datatypes.JSON `gorm:"type:json" json:"images"`
	CreatedAt   time.Time      `json:"created_at"`
	CompletedAt *time.Time     `json:"completed_at"`
}

// AgentConsumption 代理消费记录
type AgentConsumption struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	AgentID    uint      `gorm:"index" json:"agent_id"`
	OrderID    uint      `gorm:"index" json:"order_id"`
	Amount     float64   `gorm:"type:decimal(20,4)" json:"amount"`
	Commission float64   `gorm:"type:decimal(20,4)" json:"commission"`
	Type       string    `gorm:"size:32" json:"type"` // order/renewal/upgrade
	CreatedAt  time.Time `json:"created_at"`
}

// AgentToken 代理API令牌
type AgentToken struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	AgentID    uint       `gorm:"index" json:"agent_id"`
	Token      string     `gorm:"size:256;uniqueIndex" json:"token"`
	Name       string     `gorm:"size:128" json:"name"`
	ExpiresAt  time.Time  `json:"expires_at"`
	LastUsedAt *time.Time `json:"last_used_at"`
	Enabled    bool       `gorm:"default:true" json:"enabled"`
	CreatedAt  time.Time  `json:"created_at"`
}

// AgentLog 代理操作日志
type AgentLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AgentID   uint      `gorm:"index" json:"agent_id"`
	Action    string    `gorm:"size:64" json:"action"`
	Target    string    `gorm:"size:64" json:"target"`
	TargetID  uint      `json:"target_id"`
	Detail    string    `gorm:"type:text" json:"detail"`
	IP        string    `gorm:"size:64" json:"ip"`
	CreatedAt time.Time `json:"created_at"`
}

// AgentAfterSale 代理售后
type AgentAfterSale struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	AgentID     uint       `gorm:"index" json:"agent_id"`
	OrderID     uint       `gorm:"index" json:"order_id"`
	Type        string     `gorm:"size:32" json:"type"` // refund/replace/repair
	Status      string     `gorm:"size:32" json:"status"` // pending/approved/rejected/completed
	Reason      string     `gorm:"type:text" json:"reason"`
	Amount      float64    `gorm:"type:decimal(20,4)" json:"amount"`
	ProcessedAt *time.Time `json:"processed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
