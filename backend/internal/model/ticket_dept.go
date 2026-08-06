package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// TicketDepartment 工单部门
type TicketDepartment struct {
	gorm.Model
	Name           string         `gorm:"type:varchar(128);not null" json:"name"`
	Description    string         `gorm:"type:text" json:"description"`
	Slug           string         `gorm:"type:varchar(128);uniqueIndex" json:"slug"`
	ParentID       *uint          `gorm:"index" json:"parent_id"`
	Parent         *TicketDepartment `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children       []TicketDepartment `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	ManagerIDs     datatypes.JSON `gorm:"type:json" json:"manager_ids"` // 负责人ID列表 [1,2,3]
	MemberIDs      datatypes.JSON `gorm:"type:json" json:"member_ids"`  // 成员ID列表
	SortOrder      int            `gorm:"default:0;index" json:"sort_order"`
	Status         int16          `gorm:"type:smallint;default:1;not null;index" json:"status"` // 1=启用 0=禁用
	AutoAssign     bool           `gorm:"default:false" json:"auto_assign"` // 自动分配
	AssignRule     datatypes.JSON `gorm:"type:json" json:"assign_rule"` // 自动分配规则：round_robin(轮询)/least_load(最少工单)/random(随机)
	EmailNotify    bool           `gorm:"default:true" json:"email_notify"` // 邮件通知
	SMSNotify      bool           `gorm:"default:false" json:"sms_notify"`  // 短信通知
	AutoReply      string         `gorm:"type:text" json:"auto_reply"` // 自动回复内容
	TicketPrefix   string         `gorm:"type:varchar(16)" json:"ticket_prefix"` // 工单编号前缀
	Tickets        []Ticket       `gorm:"foreignKey:DepartmentID" json:"tickets,omitempty"`
}
