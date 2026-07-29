package model

import (
	"time"

	"gorm.io/gorm"
)

// BatchSyncTask 批量同步任务
type BatchSyncTask struct {
	gorm.Model
	Name        string         `gorm:"type:varchar(128);not null" json:"name"`
	Type        string         `gorm:"type:varchar(32);not null;index" json:"type"` // server/product/user
	Status      int8           `gorm:"type:smallint;default:0;index" json:"status"` // 0=待执行 1=执行中 2=已完成 3=失败
	Total       int            `gorm:"default:0" json:"total"`
	Success     int            `gorm:"default:0" json:"success"`
	Failed      int            `gorm:"default:0" json:"failed"`
	Skipped     int            `gorm:"default:0" json:"skipped"`
	ErrorMsg    string         `gorm:"type:text" json:"error_msg"`
	StartedAt   *time.Time     `json:"started_at"`
	FinishedAt  *time.Time     `json:"finished_at"`
	Duration    int            `gorm:"default:0" json:"duration"` // 耗时秒
	CreatedBy   uint           `gorm:"index" json:"created_by"`
}

// BatchSyncLog 批量同步日志
type BatchSyncLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TaskID     uint      `gorm:"index;not null" json:"task_id"`
	TargetID   uint      `gorm:"index" json:"target_id"` // 同步目标ID
	TargetType string    `gorm:"type:varchar(32)" json:"target_type"`
	Status     int8      `gorm:"type:smallint;not null" json:"status"` // 1=成功 2=失败 3=跳过
	ErrorMsg   string    `gorm:"type:text" json:"error_msg"`
	CreatedAt  time.Time `json:"created_at"`
}
