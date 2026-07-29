package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Hook 钩子定义
type Hook struct {
	gorm.Model
	Name        string         `gorm:"type:varchar(128);not null;uniqueIndex" json:"name"`
	Code        string         `gorm:"type:varchar(64);not null;uniqueIndex" json:"code"`
	Event       string         `gorm:"type:varchar(64);not null;index" json:"event"` // 触发事件: user.created/order.paid/ticket.replied 等
	Type        string         `gorm:"type:varchar(32);not null" json:"type"` // url/script/plugin
	URL         string         `gorm:"type:varchar(512)" json:"url"` // webhook URL
	Script      string         `gorm:"type:text" json:"script"` // 脚本内容
	PluginID    *uint          `gorm:"index" json:"plugin_id"`
	Headers     datatypes.JSON `gorm:"type:jsonb" json:"headers"` // 自定义请求头
	Params      datatypes.JSON `gorm:"type:jsonb" json:"params"` // 附加参数
	Status      int16          `gorm:"type:smallint;default:1;index" json:"status"` // 1=启用 0=禁用
	IsSystem    bool           `gorm:"default:false" json:"is_system"`
	Timeout     int            `gorm:"default:30" json:"timeout"` // 超时秒数
	RetryCount  int            `gorm:"default:0" json:"retry_count"` // 重试次数
	Description string         `gorm:"type:text" json:"description"`
	LastRunAt   *time.Time     `json:"last_run_at"`
	RunCount    int64          `gorm:"default:0" json:"run_count"`
	FailCount   int64          `gorm:"default:0" json:"fail_count"`
}

// HookLog 钩子执行日志
type HookLog struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	HookID     uint           `gorm:"index;not null" json:"hook_id"`
	Hook       Hook           `gorm:"foreignKey:HookID" json:"hook,omitempty"`
	Event      string         `gorm:"type:varchar(64);not null;index" json:"event"`
	Request    datatypes.JSON `gorm:"type:jsonb" json:"request"` // 请求数据
	Response   string         `gorm:"type:text" json:"response"` // 响应内容
	StatusCode int            `gorm:"default:0" json:"status_code"`
	Status     int8           `gorm:"type:smallint;not null" json:"status"` // 1=成功 2=失败 3=超时
	ErrorMsg   string         `gorm:"type:text" json:"error_msg"`
	Duration   int            `gorm:"default:0" json:"duration"` // 耗时毫秒
	RelType    string         `gorm:"type:varchar(32)" json:"rel_type"`
	RelID      uint           `gorm:"index" json:"rel_id"`
	CreatedAt  time.Time      `json:"created_at"`
}
