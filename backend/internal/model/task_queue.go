package model

import (
	"time"

	"gorm.io/gorm"
)

// TaskQueue 任务队列模型
type TaskQueue struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Type       string         `gorm:"size:50;not null;index" json:"type"` // sync_status, renew, suspend, etc.
	TargetID   uint           `gorm:"index" json:"target_id"`             // 服务ID/订单ID等
	Status     string         `gorm:"size:20;default:pending;index" json:"status"` // pending, running, completed, failed
	Priority   int            `gorm:"default:0" json:"priority"`
	Error      string         `gorm:"type:text" json:"error"`
	RetryCount int            `gorm:"default:0" json:"retry_count"`
	MaxRetry   int            `gorm:"default:3" json:"max_retry"`
	RunAt      *time.Time     `gorm:"index" json:"run_at"`
	StartedAt  *time.Time     `json:"started_at"`
	EndedAt    *time.Time     `json:"ended_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (TaskQueue) TableName() string {
	return "task_queues"
}
