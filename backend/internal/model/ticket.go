package model

import (
	"time"

	"gorm.io/gorm"
)

// Ticket 工单模型
type Ticket struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `gorm:"index;not null" json:"user_id"`
	TicketNo   string         `gorm:"size:50;uniqueIndex;not null" json:"ticket_no"`
	Subject    string         `gorm:"size:200;not null" json:"subject"`
	Status     string         `gorm:"size:20;default:open" json:"status"` // open, pending, closed
	Priority   string         `gorm:"size:20;default:normal" json:"priority"` // low, normal, high, urgent
	Department string         `gorm:"size:50" json:"department"`
	AssignedTo *uint          `gorm:"index" json:"assigned_to"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Ticket) TableName() string {
	return "tickets"
}

// TicketDepartment 工单部门模型
type TicketDepartment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	Status    string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (TicketDepartment) TableName() string {
	return "ticket_departments"
}

// TicketStatus 工单状态模型
type TicketStatus struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Value     string         `gorm:"size:20;uniqueIndex;not null" json:"value"`
	Label     string         `gorm:"size:50;not null" json:"label"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (TicketStatus) TableName() string {
	return "ticket_statuses"
}
