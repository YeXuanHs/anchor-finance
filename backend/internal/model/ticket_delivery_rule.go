package model

import "time"

// TicketDeliveryRule 工单投递规则
type TicketDeliveryRule struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"size:100;not null" json:"name"`
	DepartmentID uint  `gorm:"index" json:"department_id"`
	ProductID   uint   `gorm:"index" json:"product_id"`
	Keyword     string `gorm:"size:200" json:"keyword"`
	UpstreamURL string `gorm:"size:500" json:"upstream_url"`
	UpstreamKey string `gorm:"size:200" json:"upstream_key"`
	Status      string `gorm:"size:20;default:active" json:"status"`
	SortOrder   int    `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 指定表名
func (TicketDeliveryRule) TableName() string { return "ticket_delivery_rules" }
