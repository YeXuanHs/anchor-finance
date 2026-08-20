package model

import (
	"time"

	"gorm.io/gorm"
)

// CronTask 定时任务
type CronTask struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"type:varchar(128);not null" json:"name"`
	Type         string         `gorm:"type:varchar(32);not null;index" json:"type"` // custom/system/plugin
	CronExpr     string         `gorm:"type:varchar(64);not null" json:"cron_expr"`
	Command      string         `gorm:"type:text" json:"command"`
	Params       string         `gorm:"type:text" json:"params"` // JSON格式参数
	Status       int8           `gorm:"type:smallint;default:1;index" json:"status"` // 1启用 0禁用
	LastRunAt    *time.Time     `json:"last_run_at"`
	NextRunAt    *time.Time     `gorm:"index" json:"next_run_at"`
	LastResult   string         `gorm:"type:text" json:"last_result"` // success/failed + message
	LastError    string         `gorm:"type:text" json:"last_error"`
	RunCount     int            `gorm:"default:0" json:"run_count"`
	Timeout      int            `gorm:"default:60" json:"timeout"` // 超时秒数
	Priority     int            `gorm:"default:0" json:"priority"`
	Description  string         `gorm:"type:text" json:"description"`
	CreatedBy    uint           `gorm:"index" json:"created_by"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// CronTaskLog 定时任务执行日志
type CronTaskLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TaskID     uint      `gorm:"index;not null" json:"task_id"`
	TaskName   string    `gorm:"type:varchar(128)" json:"task_name"`
	Status     int8      `gorm:"type:smallint;not null" json:"status"` // 1执行中 2成功 3失败 4超时
	Output     string    `gorm:"type:text" json:"output"`
	ErrorMsg   string    `gorm:"type:text" json:"error_msg"`
	StartedAt  time.Time `gorm:"not null" json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`
	Duration   int       `gorm:"default:0" json:"duration"` // 耗时毫秒
	Trigger    string    `gorm:"type:varchar(16)" json:"trigger"` // auto/manual
}
