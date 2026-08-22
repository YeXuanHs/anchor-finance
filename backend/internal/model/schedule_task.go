package model

import (
	"time"

	"gorm.io/gorm"
)

// ScheduleTask 定时任务模型
type ScheduleTask struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Cron      string         `gorm:"size:50" json:"cron"`
	Type      string         `gorm:"size:50" json:"type"` // suspend, renew, report
	Status    string         `gorm:"size:20;default:active" json:"status"` // active, disabled
	LastRunAt *time.Time     `json:"last_run_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (ScheduleTask) TableName() string {
	return "schedule_tasks"
}

// ScheduleRun 定时任务运行记录
type ScheduleRun struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	TaskID     uint       `gorm:"index" json:"task_id"`
	Status     string     `gorm:"size:20" json:"status"` // success, failed
	Detail     string     `gorm:"type:text" json:"detail"`
	StartedAt  time.Time  `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

// TableName 指定表名
func (ScheduleRun) TableName() string {
	return "schedule_runs"
}
