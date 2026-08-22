package model

import (
	"time"

	"gorm.io/gorm"
)

// MarketingPush 营销推送模型
type MarketingPush struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"size:200;not null" json:"title"`
	Content     string         `gorm:"type:text" json:"content"`
	Type        string         `gorm:"size:20;not null" json:"type"` // email, sms, system
	TargetType  string         `gorm:"size:20;default:all" json:"target_type"` // all, group, user
	TargetIDs   string         `gorm:"type:text" json:"target_ids"`   // JSON的目标ID列表
	Status      string         `gorm:"size:20;default:draft" json:"status"` // draft, sending, sent, failed
	SentAt      *time.Time     `json:"sent_at"`
	SentCount   int            `gorm:"default:0" json:"sent_count"`
	FailedCount int            `gorm:"default:0" json:"failed_count"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (MarketingPush) TableName() string {
	return "marketing_pushes"
}
