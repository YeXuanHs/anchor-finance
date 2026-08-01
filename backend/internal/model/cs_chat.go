package model

import (
	"time"

	"gorm.io/gorm"
)

// CSChatSession 客服聊天会话
// 对应 anchor_cloud_finance_pro 的 acfp_chat_sessions
type CSChatSession struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	UserID          uint       `gorm:"index;default:0" json:"user_id"`
	VisitorID       string     `gorm:"type:varchar(64);index" json:"visitor_id"`
	ClientName      string     `gorm:"type:varchar(100)" json:"client_name"`
	ClientEmail     string     `gorm:"type:varchar(128)" json:"client_email"`
	ClientPhone     string     `gorm:"type:varchar(32)" json:"client_phone"`
	ClientAvatar    string     `gorm:"type:varchar(500)" json:"client_avatar"`
	Mode            string     `gorm:"type:varchar(16);default:ai;comment:ai/human" json:"mode"`
	Status          string     `gorm:"type:varchar(16);default:open;index;comment:open/closed" json:"status"`
	Rating          int        `gorm:"default:0;comment:1-5星评分" json:"rating"`
	RatingComment   string     `gorm:"type:varchar(500)" json:"rating_comment"`
	StaffID         uint       `gorm:"default:0" json:"staff_id"`
	ThreadSessionID uint       `gorm:"default:0;comment:关联的主线程会话ID" json:"thread_session_id"`
	LastMessage     string     `gorm:"type:text" json:"last_message"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	ClosedAt        *time.Time `json:"closed_at"`
}

// CSChatMessage 客服聊天消息
type CSChatMessage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SessionID uint      `gorm:"index;not null" json:"session_id"`
	Role      string    `gorm:"type:varchar(20);not null;comment:user/assistant/system/human" json:"role"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// CSChatConfig 客服聊天配置
// 将所有配置存储为 JSON
type CSChatConfig struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Key    string `gorm:"type:varchar(50);uniqueIndex;not null" json:"key"`
	Value  string `gorm:"type:text" json:"value"`
	Remark string `gorm:"type:varchar(255)" json:"remark"`
}

// CSChatQuickReply 快捷回复
type CSChatQuickReply struct {
	gorm.Model
	Content   string `gorm:"type:text;not null" json:"content"`
	SortOrder int    `gorm:"default:0" json:"sort_order"`
	IsActive  bool   `gorm:"default:true" json:"is_active"`
}
