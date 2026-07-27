package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// EmailTemplate 邮件/短信/站内信模板
type EmailTemplate struct {
	gorm.Model
	Code      string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"code"`     // 模板标识
	Name      string         `gorm:"type:varchar(128);not null" json:"name"`                // 模板名称
	Subject   string         `gorm:"type:varchar(256);not null" json:"subject"`             // 标题
	Body      string         `gorm:"type:text;not null" json:"body"`                        // 内容（支持Go模板语法）
	Type      string         `gorm:"type:varchar(16);not null;index" json:"type"`           // email/sms/notice
	Variables datatypes.JSON `gorm:"type:jsonb" json:"variables"`                          // 可用变量列表
	Format    string         `gorm:"type:varchar(16);default:'html'" json:"format"`         // html/plain
	Language  string         `gorm:"type:varchar(16);default:'zh-CN'" json:"language"`      // 语言
	IsSystem  bool           `gorm:"default:false" json:"is_system"`                        // 系统模板不可删除
	Status    int16          `gorm:"type:smallint;default:1;not null" json:"status"`        // 1=启用 0=禁用
	SendCount int            `gorm:"default:0" json:"send_count"`                           // 发送次数
	LastSentAt *time.Time    `json:"last_sent_at"`                                          // 最近发送时间
}

// EmailTemplateLog 模板发送记录
type EmailTemplateLog struct {
	gorm.Model
	TemplateID uint   `gorm:"index;not null" json:"template_id"`
	Template   EmailTemplate `gorm:"foreignKey:TemplateID" json:"template,omitempty"`
	Recipient  string `gorm:"type:varchar(256);not null" json:"recipient"`
	Subject    string `gorm:"type:varchar(256)" json:"subject"`
	Content    string `gorm:"type:text" json:"content"`
	Type       string `gorm:"type:varchar(16);not null" json:"type"` // email/sms/notice
	Status     int16  `gorm:"type:smallint;default:1;not null" json:"status"` // 1=待发送 2=成功 3=失败
	Error      string `gorm:"type:text" json:"error"`
	SentAt     *time.Time `json:"sent_at"`
	OperatorID uint   `gorm:"index" json:"operator_id"`
}
