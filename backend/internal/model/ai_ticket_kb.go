package model

import (
	"time"

	"gorm.io/gorm"
)

// AITicketKnowledge AI工单知识库
// 对应 mianyu_ai_ticket 的 mianyu_ai_ticket_knowledge
// 简单结构：标题、关键词、问题、答案
type AITicketKnowledge struct {
	gorm.Model
	Title    string `gorm:"type:varchar(255);not null" json:"title"`
	Keywords string `gorm:"type:varchar(500)" json:"keywords"`
	Question string `gorm:"type:text" json:"question"`
	Answer   string `gorm:"type:text" json:"answer"`
	Sort     int    `gorm:"default:0" json:"sort"`
	Status   int8   `gorm:"default:1;comment:1=启用 0=禁用" json:"status"`
}

// AITicketRule AI工单自动化规则
// 对应 mianyu_ai_ticket 的 mianyu_ai_ticket_rules
type AITicketRule struct {
	gorm.Model
	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	Keywords    string `gorm:"type:varchar(500)" json:"keywords"`
	DeptFilter  string `gorm:"type:varchar(255);comment:适用部门ID逗号分隔" json:"dept_filter"`
	TargetDept  int    `gorm:"default:0;comment:转派目标部门" json:"target_dept"`
	Priority    int    `gorm:"default:100" json:"priority"`
	Action      string `gorm:"type:varchar(50);default:reply_only;comment:reply_only/reply_and_transfer/transfer_only/close" json:"action"`
	PromptExtra string `gorm:"type:text;comment:额外提示词" json:"prompt_extra"`
	SampleReply string `gorm:"type:text;comment:示例回复" json:"sample_reply"`
	Status      int8   `gorm:"default:1" json:"status"`
}

// AITicketQueue AI工单处理队列
// 对应 mianyu_ai_ticket 的 mianyu_ai_ticket_queue
type AITicketQueue struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	TicketID   uint       `gorm:"index" json:"ticket_id"`
	ReplyID    uint       `gorm:"default:0" json:"reply_id"`
	EventType  string     `gorm:"type:varchar(32)" json:"event_type"`
	UserID     uint       `gorm:"default:0" json:"user_id"`
	DeptID     uint       `gorm:"default:0" json:"dept_id"`
	Status     string     `gorm:"type:varchar(16);default:pending;index;comment:pending/processing/completed/failed" json:"status"`
	Attempts   int8       `gorm:"default:0" json:"attempts"`
	ErrorMsg   string     `gorm:"type:text" json:"error_msg"`
	ModelUsed  string     `gorm:"type:varchar(128)" json:"model_used"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// AITicketProcessLog AI工单处理日志
type AITicketProcessLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TicketID   uint      `gorm:"index" json:"ticket_id"`
	Event      string    `gorm:"type:varchar(50)" json:"event"`
	Decision   string    `gorm:"type:varchar(50)" json:"decision"`
	Confidence float64   `gorm:"type:decimal(5,4)" json:"confidence"`
	Status     string    `gorm:"type:varchar(16);default:success" json:"status"`
	Message    string    `gorm:"type:text" json:"message"`
	ModelUsed  string    `gorm:"type:varchar(128)" json:"model_used"`
	CreatedAt  time.Time `json:"created_at"`
}

// AITicketNotifyLog AI工单通知日志
type AITicketNotifyLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TicketID  uint      `gorm:"index" json:"ticket_id"`
	Channel   string    `gorm:"type:varchar(32);default:mail;comment:mail/telegram/webhook/wechat" json:"channel"`
	Status    string    `gorm:"type:varchar(16);default:success" json:"status"`
	Content   string    `gorm:"type:text" json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// AITicketMode AI工单模式（每个工单独立控制AI/人工）
type AITicketMode struct {
	TicketID  uint      `gorm:"primaryKey" json:"ticket_id"`
	Mode      string    `gorm:"type:varchar(16);default:ai;comment:ai/human" json:"mode"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AITicketConfig AI工单全局配置
// 将所有配置存为单条 JSON 记录
type AITicketConfig struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Key   string `gorm:"type:varchar(50);uniqueIndex;not null" json:"key"`
	Value string `gorm:"type:text" json:"value"`
}
