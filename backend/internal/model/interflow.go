package model

import "gorm.io/gorm"

// Interflow 工单流转
type Interflow struct {
	gorm.Model
	Name        string `gorm:"type:varchar(128);not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	IsEnabled   bool   `gorm:"default:true;index" json:"is_enabled"`
	Priority    int    `gorm:"default:0" json:"priority"`
	Config      string `gorm:"type:json" json:"config"`
}

// InterflowChatRecord 工单流转聊天记录
type InterflowChatRecord struct {
	gorm.Model
	InterflowID uint   `gorm:"index;not null" json:"interflow_id"`
	TicketID    uint   `gorm:"index;not null" json:"ticket_id"`
	UserID      uint   `gorm:"index" json:"user_id"`
	Role        string `gorm:"type:varchar(16);not null" json:"role"` // user/admin/system
	Content     string `gorm:"type:text;not null" json:"content"`
	Attachments string `gorm:"type:json" json:"attachments"`
}

// InterflowFunc 工单流转功能
type InterflowFunc struct {
	gorm.Model
	Name        string `gorm:"type:varchar(64);not null" json:"name"`
	Code        string `gorm:"type:varchar(32);uniqueIndex;not null" json:"code"`
	Type        string `gorm:"type:varchar(16)" json:"type"`
	Config      string `gorm:"type:json" json:"config"`
	IsEnabled   bool   `gorm:"default:true" json:"is_enabled"`
	Description string `gorm:"type:varchar(256)" json:"description"`
}

// InterflowKeyword 工单流转关键词
type InterflowKeyword struct {
	gorm.Model
	Keyword  string `gorm:"type:varchar(64);not null;index" json:"keyword"`
	Type     string `gorm:"type:varchar(16);default:auto_reply" json:"type"` // auto_reply/block/flag
	Response string `gorm:"type:text" json:"response"`
	IsEnabled bool  `gorm:"default:true" json:"is_enabled"`
}

// InterflowMatchingExecute 工单流转匹配执行
type InterflowMatchingExecute struct {
	gorm.Model
	Name      string `gorm:"type:varchar(64);not null" json:"name"`
	Condition string `gorm:"type:json;not null" json:"condition"`
	Action    string `gorm:"type:json;not null" json:"action"`
	Priority  int    `gorm:"default:0" json:"priority"`
	IsEnabled bool   `gorm:"default:true" json:"is_enabled"`
}

// InterflowClients 工单流转客户关联
type InterflowClients struct {
	gorm.Model
	InterflowID uint   `gorm:"index;not null" json:"interflow_id"`
	ClientID    uint   `gorm:"index;not null" json:"client_id"`
	Role        string `gorm:"type:varchar(16)" json:"role"`
}
