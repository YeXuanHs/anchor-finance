package model

import "time"

// AITicketQueue AI工单队列（参考mianyu_ai_ticket_queue）
type AITicketQueue struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TicketID   uint      `gorm:"index;not null" json:"ticket_id"`
	ReplyID    uint      `gorm:"index" json:"reply_id"`
	EventType  string    `gorm:"size:32;not null" json:"event_type"` // open, user_reply, monitoring
	UserID     uint      `gorm:"index" json:"user_id"`
	DeptID     uint      `gorm:"index" json:"dept_id"`
	Status     string    `gorm:"size:16;default:pending" json:"status"` // pending, processing, done, failed
	Attempts   int       `gorm:"default:0" json:"attempts"`
	ErrorMsg   string    `gorm:"type:text" json:"error_msg"`
	ModelUsed  string    `gorm:"size:128" json:"model_used"`
	Step       string    `gorm:"size:50" json:"step"` // 当前处理步骤
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (AITicketQueue) TableName() string { return "ai_ticket_queue" }

// AITicketMode AI工单模式（每工单独立AI/人工控制）
type AITicketMode struct {
	TicketID  uint      `gorm:"primaryKey" json:"ticket_id"`
	Mode      string    `gorm:"size:16;default:ai" json:"mode"` // ai, human
	UpdatedAt time.Time `json:"updated_at"`
}

func (AITicketMode) TableName() string { return "ai_ticket_modes" }

// AITicketKnowledge AI工单知识库
type AITicketKnowledge struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:255;not null" json:"title"`
	Keywords  string    `gorm:"size:500" json:"keywords"`
	Question  string    `gorm:"type:text" json:"question"`
	Answer    string    `gorm:"type:text" json:"answer"`
	Sort      int       `gorm:"default:0" json:"sort"`
	Status    int       `gorm:"default:1" json:"status"` // 1=启用 0=禁用
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (AITicketKnowledge) TableName() string { return "ai_ticket_knowledge" }

// AITicketRule AI工单自动化规则
type AITicketRule struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Keywords    string    `gorm:"size:500" json:"keywords"`    // 逗号分隔关键词
	DeptFilter  string    `gorm:"size:255" json:"dept_filter"` // 逗号分隔部门ID
	TargetDept  int       `gorm:"default:0" json:"target_dept"`
	Priority    int       `gorm:"default:100" json:"priority"`
	Action      string    `gorm:"size:50;default:reply_only" json:"action"` // reply_only, transfer, close
	PromptExtra string    `gorm:"type:text" json:"prompt_extra"`
	SampleReply string    `gorm:"type:text" json:"sample_reply"`
	Status      int       `gorm:"default:1" json:"status"` // 1=启用 0=禁用
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (AITicketRule) TableName() string { return "ai_ticket_rules" }

// AITicketProcessLog AI工单处理日志
type AITicketProcessLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TicketID   uint      `gorm:"index" json:"ticket_id"`
	Event      string    `gorm:"size:50" json:"event"`
	Decision   string    `gorm:"size:50" json:"decision"`
	Confidence float64   `gorm:"type:decimal(5,4);default:0" json:"confidence"`
	Status     string    `gorm:"size:16;default:success" json:"status"` // success, failed, skipped
	Message    string    `gorm:"type:text" json:"message"`
	ModelUsed  string    `gorm:"size:128" json:"model_used"`
	CreatedAt  time.Time `json:"created_at"`
}

func (AITicketProcessLog) TableName() string { return "ai_ticket_process_logs" }
