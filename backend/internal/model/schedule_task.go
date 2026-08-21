package model

import (
	"time"

	"gorm.io/gorm"
)

// ScheduleTask 定时任务模型
type ScheduleTask struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Description string         `gorm:"size:500" json:"description"`
	Cron        string         `gorm:"size:50" json:"cron"` // cron表达式
	Command     string         `gorm:"size:500" json:"command"`
	Status      string         `gorm:"size:20;default:active" json:"status"` // active, disabled
	LastRunAt   *time.Time     `json:"last_run_at"`
	NextRunAt   *time.Time     `json:"next_run_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (ScheduleTask) TableName() string {
	return "schedule_tasks"
}

// ScheduleRun 定时任务运行记录模型
type ScheduleRun struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TaskID    uint      `gorm:"index;not null" json:"task_id"`
	Status    string    `gorm:"size:20;not null" json:"status"` // running, success, failed
	Output    string    `gorm:"type:text" json:"output"`
	Error     string    `gorm:"type:text" json:"error"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (ScheduleRun) TableName() string {
	return "schedule_runs"
}
