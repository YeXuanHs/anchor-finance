package model

import (
	"time"

	"gorm.io/gorm"
)

// KnowledgeBaseCategory 知识库分类
// 移植自 mianyu_ai_ticket
type KnowledgeBaseCategory struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Description string `gorm:"type:varchar(500)" json:"description"`
	ParentID    uint   `gorm:"default:0;index" json:"parent_id"`
	SortOrder   int    `gorm:"default:0" json:"sort_order"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
}

// KnowledgeBaseArticle 知识库文章
type KnowledgeBaseArticle struct {
	gorm.Model
	CategoryID uint   `gorm:"index;not null" json:"category_id"`
	Title      string `gorm:"type:varchar(255);not null" json:"title"`
	Content    string `gorm:"type:longtext;not null" json:"content"`
	Summary    string `gorm:"type:varchar(500)" json:"summary"`
	Tags       string `gorm:"type:varchar(500);comment:逗号分隔的标签" json:"tags"`
	Keywords   string `gorm:"type:varchar(500);comment:AI匹配关键词" json:"keywords"`
	IsFAQ      bool   `gorm:"default:false;index;comment:是否为FAQ" json:"is_faq"`
	ViewCount  int    `gorm:"default:0" json:"view_count"`
	HelpCount  int    `gorm:"default:0;comment:有帮助次数" json:"help_count"`
	IsActive   bool   `gorm:"default:true" json:"is_active"`
	SortOrder  int    `gorm:"default:0" json:"sort_order"`
}

// AIConfig AI 配置
// 用于工单自动回复和购物助手
type AIConfig struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Provider       string    `gorm:"type:varchar(50);not null;comment:openai/claude/deepseek/custom" json:"provider"`
	APIKey         string    `gorm:"type:varchar(500)" json:"api_key"`
	APIEndpoint    string    `gorm:"type:varchar(500)" json:"api_endpoint"`
	Model          string    `gorm:"type:varchar(100);not null" json:"model"`
	MaxTokens      int       `gorm:"default:2000" json:"max_tokens"`
	Temperature    float64   `gorm:"type:decimal(3,2);default:0.7" json:"temperature"`
	SystemPrompt   string    `gorm:"type:text" json:"system_prompt"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// AITicketAutoReplyConfig 工单自动回复配置
type AITicketAutoReplyConfig struct {
	ID                uint    `gorm:"primaryKey" json:"id"`
	Enabled           bool    `gorm:"default:false" json:"enabled"`
	AIConfigID        uint    `gorm:"not null" json:"ai_config_id"`
	ReplyDelay        int     `gorm:"default:5;comment:回复延迟（秒）" json:"reply_delay"`
	ConfidenceThreshold float64 `gorm:"type:decimal(3,2);default:0.7;comment:置信度阈值" json:"confidence_threshold"`
	MaxAutoReplies    int     `gorm:"default:3;comment:同一工单最大自动回复数" json:"max_auto_replies"`
	IncludeKBContent  bool    `gorm:"default:true;comment:是否包含知识库内容" json:"include_kb_content"`
	KBSearchLimit     int     `gorm:"default:5;comment:知识库搜索结果数量" json:"kb_search_limit"`
	DeptIDs           string  `gorm:"type:varchar(500);comment:适用的工单部门ID，逗号分隔" json:"dept_ids"`
	ExcludeKeywords   string  `gorm:"type:text;comment:排除关键词，逗号分隔" json:"exclude_keywords"`
	AddDisclaimer     bool    `gorm:"default:true;comment:是否添加AI回复声明" json:"add_disclaimer"`
	DisclaimerText    string  `gorm:"type:varchar(500);default:此回复由AI生成，仅供参考。如需人工帮助请回复「转人工」" json:"disclaimer_text"`
}

// AITicketLog AI 工单回复日志
type AITicketLog struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	TicketID    uint      `gorm:"index;not null" json:"ticket_id"`
	ReplyID     uint      `gorm:"comment:关联的回复ID" json:"reply_id"`
	Question    string    `gorm:"type:text" json:"question"`
	Answer      string    `gorm:"type:text" json:"answer"`
	Confidence  float64   `gorm:"type:decimal(3,2)" json:"confidence"`
	KBMatchIDs  string    `gorm:"type:varchar(500);comment:匹配的知识库文章ID" json:"kb_match_ids"`
	TokensUsed  int       `gorm:"default:0" json:"tokens_used"`
	Accepted    *bool     `gorm:"comment:用户是否接受此回复" json:"accepted"`
	CreatedAt   time.Time `json:"created_at"`
}
