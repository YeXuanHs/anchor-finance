package model

import (
	"time"

	"gorm.io/gorm"
)

// AIToolConfig AI工具配置（每个工具的启用/禁用状态）
type AIToolConfig struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Enabled   int8           `gorm:"default:1;comment:1=启用 0=禁用" json:"enabled"`
	RiskLevel string         `gorm:"type:varchar(20);default:low;comment:low/medium/high" json:"risk_level"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// AIToolExecutionLog AI工具执行日志
type AIToolExecutionLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TicketID   uint      `gorm:"index" json:"ticket_id"`
	ToolName   string    `gorm:"type:varchar(100);index" json:"tool_name"`
	Args       string    `gorm:"type:text;comment:JSON参数" json:"args"`
	Result     string    `gorm:"type:text;comment:执行结果" json:"result"`
	Success    int8      `gorm:"default:1;comment:1=成功 0=失败" json:"success"`
	ElapsedMs  int       `gorm:"default:0;comment:耗时毫秒" json:"elapsed_ms"`
	RoundIndex int       `gorm:"default:0;comment:工具调用轮次" json:"round_index"`
	CreatedAt  time.Time `json:"created_at"`
}
