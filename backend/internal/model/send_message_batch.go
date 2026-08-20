package model

import (
	"time"

	"gorm.io/gorm"
)

// SendMessageBatch 批量发送消息任务
type SendMessageBatch struct {
	gorm.Model
	BatchNo    string     `gorm:"type:varchar(64);uniqueIndex;not null" json:"batch_no"`
	SendMethod string     `gorm:"type:varchar(32);not null;index" json:"send_method"` // email/mobile/system
	Subject    string     `gorm:"type:varchar(256)" json:"subject"`
	Content    string     `gorm:"type:text;not null" json:"content"`
	Total      int        `gorm:"default:0" json:"total"`
	Success    int        `gorm:"default:0" json:"success"`
	Failed     int        `gorm:"default:0" json:"failed"`
	Status     int8       `gorm:"type:smallint;default:0;index" json:"status"` // 0=待发送 1=发送中 2=已完成 3=失败
	ErrorMsg   string     `gorm:"type:text" json:"error_msg"`
	StartedAt  *time.Time `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`
	CreatedBy  uint       `gorm:"index" json:"created_by"`
}

// SendMessageBatchRecord 批量发送消息记录
type SendMessageBatchRecord struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	BatchID   uint       `gorm:"index;not null" json:"batch_id"`
	Batch     SendMessageBatch `gorm:"foreignKey:BatchID" json:"batch,omitempty"`
	UserID    uint       `gorm:"index;not null" json:"user_id"`
	Target    string     `gorm:"type:varchar(256)" json:"target"` // 发送目标（邮箱/手机号）
	Status    int8       `gorm:"type:smallint;not null" json:"status"` // 1=成功 2=失败 3=待发送
	ErrorMsg  string     `gorm:"type:text" json:"error_msg"`
	SentAt    *time.Time `json:"sent_at"`
	CreatedAt time.Time  `json:"created_at"`
}
