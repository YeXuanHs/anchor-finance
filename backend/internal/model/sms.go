package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// SMSLog 短信发送日志
type SMSLog struct {
	gorm.Model
	Phone      string         `gorm:"type:varchar(20);index" json:"phone"`
	Content    string         `gorm:"type:text" json:"content"`
	TemplateID string         `gorm:"type:varchar(64)" json:"template_id"`
	Params     datatypes.JSON `gorm:"type:json" json:"params"`
	Status     string         `gorm:"type:varchar(32);default:pending" json:"status"` // pending/sent/failed
	Response   string         `gorm:"type:text" json:"response"`
	Operator   string         `gorm:"type:varchar(32)" json:"operator"` // mobile/unicom/telecom
	UserID     *uint          `gorm:"index" json:"user_id"`
	BatchID    *uint          `gorm:"index" json:"batch_id"`
	SentAt     *time.Time     `json:"sent_at"`
}

// SMSTemplate 短信模板
type SMSTemplate struct {
	gorm.Model
	Name    string         `gorm:"type:varchar(128);not null" json:"name"`
	Code    string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"code"`
	Content string         `gorm:"type:text;not null" json:"content"`
	Params  datatypes.JSON `gorm:"type:json" json:"params"` // 可用参数列表
	Type    string         `gorm:"type:varchar(32);default:notification" json:"type"` // verify/marketing/notification
	Enabled bool           `gorm:"default:true" json:"enabled"`
}

// SMSBatch 短信批量发送任务
type SMSBatch struct {
	gorm.Model
	Name        string     `gorm:"type:varchar(128);not null" json:"name"`
	TemplateID  uint       `gorm:"index" json:"template_id"`
	TargetGroup string     `gorm:"type:varchar(64)" json:"target_group"` // all/new/active/vip
	TotalCount  int        `gorm:"default:0" json:"total_count"`
	SentCount   int        `gorm:"default:0" json:"sent_count"`
	FailedCount int        `gorm:"default:0" json:"failed_count"`
	Status      string     `gorm:"type:varchar(32);default:pending" json:"status"` // pending/sending/completed/failed
	CreatedBy   uint       `gorm:"index" json:"created_by"`
	CompletedAt *time.Time `json:"completed_at"`
}
