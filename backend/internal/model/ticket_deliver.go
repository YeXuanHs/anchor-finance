package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// TicketDeliverRule 工单传递规则
type TicketDeliverRule struct {
	gorm.Model
	Name        string         `gorm:"type:varchar(128);not null" json:"name"`
	Departments datatypes.JSON `gorm:"type:json;not null" json:"departments"` // 部门ID列表
	Products    datatypes.JSON `gorm:"type:json" json:"products"`             // 产品ID列表
	Priority    int            `gorm:"default:0" json:"priority"`
	Status      int16          `gorm:"type:smallint;default:1;not null;index" json:"status"`
	Description string         `gorm:"type:text" json:"description"`
	Condition   datatypes.JSON `gorm:"type:json" json:"condition"` // 额外条件

	// 上游透传配置
	UpstreamType     string `gorm:"type:varchar(32)" json:"upstream_type"`      // zjmf/v10/anchorfinance/none
	UpstreamURL      string `gorm:"type:varchar(512)" json:"upstream_url"`      // 上游系统URL
	UpstreamAPIKey   string `gorm:"type:varchar(256)" json:"upstream_api_key"`  // 上游API密钥
	UpstreamDeptID   uint   `gorm:"index" json:"upstream_dept_id"`              // 上游部门ID
	UpstreamProductID uint  `gorm:"index" json:"upstream_product_id"`           // 上游产品ID
	MaskKeywords     string `gorm:"type:text" json:"mask_keywords"`             // 敏感关键词（换行分隔）
	AutoReply        int16  `gorm:"type:smallint;default:0" json:"auto_reply"`  // 是否自动回复上游响应
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
	Direction  string `gorm:"type:varchar(16)" json:"direction"` // upstream/downstream
	Status     string `gorm:"type:varchar(16)" json:"status"`    // pending/success/failed
	ErrorMsg   string `gorm:"type:text" json:"error_msg"`
	CreatedAt  string `gorm:"autoCreateTime" json:"created_at"`
}

// TicketUpstream 上游工单配置
type TicketUpstream struct {
	gorm.Model
	Name        string `gorm:"type:varchar(128);not null" json:"name"`
	Type        string `gorm:"type:varchar(32);not null" json:"type"` // zjmf/v10/anchorfinance
	URL         string `gorm:"type:varchar(512);not null" json:"url"`
	APIKey      string `gorm:"type:varchar(256)" json:"api_key"`
	Username    string `gorm:"type:varchar(64)" json:"username"`
	Status      int16  `gorm:"type:smallint;default:1;not null;index" json:"status"`
	Description string `gorm:"type:text" json:"description"`
}

// TicketUpstreamMapping 上游工单映射
type TicketUpstreamMapping struct {
	ID              uint   `gorm:"primaryKey" json:"id"`
	LocalTicketID   uint   `gorm:"index;not null" json:"local_ticket_id"`
	UpstreamTicketID string `gorm:"type:varchar(64);not null" json:"upstream_ticket_id"`
	UpstreamID      uint   `gorm:"index;not null" json:"upstream_id"`
	Direction       string `gorm:"type:varchar(16)" json:"direction"` // upstream/downstream
	CreatedAt       string `gorm:"autoCreateTime" json:"created_at"`
}

// TableName 表名
func (TicketDeliverRule) TableName() string {
	return "ticket_deliver_rules"
}

// TableName 表名
func (TicketDeliverLog) TableName() string {
	return "ticket_deliver_logs"
}

// TableName 表名
func (TicketUpstream) TableName() string {
	return "ticket_upstreams"
}

// TableName 表名
func (TicketUpstreamMapping) TableName() string {
	return "ticket_upstream_mappings"
}
