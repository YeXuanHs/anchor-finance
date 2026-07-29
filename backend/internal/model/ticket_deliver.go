package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// TicketDeliverRule 工单传递规则
type TicketDeliverRule struct {
	gorm.Model
	Name        string         `gorm:"type:varchar(128);not null" json:"name"`
	Departments datatypes.JSON `gorm:"type:jsonb;not null" json:"departments"` // 部门ID列表
	Products    datatypes.JSON `gorm:"type:jsonb" json:"products"` // 产品ID列表
	Priority    int            `gorm:"default:0" json:"priority"`
	Status      int16          `gorm:"type:smallint;default:1;not null;index" json:"status"`
	Description string         `gorm:"type:text" json:"description"`
	Condition   datatypes.JSON `gorm:"type:jsonb" json:"condition"` // 额外条件
}

// TicketDeliverLog 工单传递记录
type TicketDeliverLog struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	TicketID   uint   `gorm:"index;not null" json:"ticket_id"`
	RuleID     *uint  `gorm:"index" json:"rule_id"`
	FromDept   uint   `gorm:"index" json:"from_dept"`
	ToDept     uint   `gorm:"index" json:"to_dept"`
	OperatorID uint   `gorm:"index" json:"operator_id"`
	Reason     string `gorm:"type:text" json:"reason"`
	CreatedAt  string `gorm:"autoCreateTime" json:"created_at"`
}
