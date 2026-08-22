package model

import (
	"time"

	"gorm.io/gorm"
)

// TicketRule 工单传递规则模型
type TicketRule struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Name           string         `gorm:"size:100;not null" json:"name"`
	ConditionType  string         `gorm:"size:50;not null" json:"condition_type"` // department, priority, keyword
	ConditionValue string         `gorm:"size:200;not null" json:"condition_value"`
	ActionType     string         `gorm:"size:50;not null" json:"action_type"` // assign, transfer, notify
	ActionValue    string         `gorm:"size:200;not null" json:"action_value"`
	Priority       int            `gorm:"default:0" json:"priority"`
	Status         string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (TicketRule) TableName() string {
	return "ticket_rules"
}
