package model

import (
	"time"

	"gorm.io/gorm"
)

// CronURLTask URL定时任务
type CronURLTask struct {
	gorm.Model
	Name        string         `gorm:"type:varchar(128);not null" json:"name"`
	URL         string         `gorm:"type:varchar(512);not null" json:"url"`
	Method      string         `gorm:"type:varchar(16);default:GET" json:"method"` // GET/POST
	Headers     string         `gorm:"type:text" json:"headers"` // JSON格式
	Body        string         `gorm:"type:text" json:"body"`
	CronExpr    string         `gorm:"type:varchar(64);not null" json:"cron_expr"`
	Status      int8           `gorm:"type:smallint;default:1;index" json:"status"` // 1=启用 0=禁用
	Timeout     int            `gorm:"default:30" json:"timeout"` // 超时秒数
	LastRunAt   *time.Time     `json:"last_run_at"`
	NextRunAt   *time.Time     `gorm:"index" json:"next_run_at"`
	LastResult  string         `gorm:"type:text" json:"last_result"`
	LastError   string         `gorm:"type:text" json:"last_error"`
	RunCount    int            `gorm:"default:0" json:"run_count"`
	FailCount   int            `gorm:"default:0" json:"fail_count"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedBy   uint           `gorm:"index" json:"created_by"`
}

// CronURLTaskLog URL定时任务执行日志
type CronURLTaskLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	TaskID       uint      `gorm:"index;not null" json:"task_id"`
	TaskName     string    `gorm:"type:varchar(128)" json:"task_name"`
	Status       int8      `gorm:"type:smallint;not null" json:"status"` // 1=成功 2=失败 3=超时
	StatusCode   int       `gorm:"default:0" json:"status_code"` // HTTP状态码
	Response     string    `gorm:"type:text" json:"response"` // 响应内容
	ErrorMsg     string    `gorm:"type:text" json:"error_msg"`
	Duration     int       `gorm:"default:0" json:"duration"` // 耗时毫秒
	StartedAt    time.Time `gorm:"not null" json:"started_at"`
	FinishedAt   *time.Time `json:"finished_at"`
}
