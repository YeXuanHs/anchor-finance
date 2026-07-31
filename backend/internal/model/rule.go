package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Rule 规则
type Rule struct {
	gorm.Model
	Name       string         `gorm:"type:varchar(128);not null" json:"name"`         // 规则名称
	Code       string         `gorm:"type:varchar(64);uniqueIndex" json:"code"`       // 规则编码
	Type       string         `gorm:"type:varchar(32);not null;index" json:"type"`    // 规则类型: pricing/promotion/notification/limit/filter/routing
	Condition  datatypes.JSON `gorm:"type:json;not null" json:"condition"`           // 条件JSON
	Action     datatypes.JSON `gorm:"type:json;not null" json:"action"`              // 动作JSON
	Priority   int            `gorm:"default:0;index" json:"priority"`                // 优先级，数值越大越优先
	Status     int16          `gorm:"type:smallint;default:1;not null;index" json:"status"` // 1=启用 0=禁用
	IsSystem   bool           `gorm:"default:false" json:"is_system"`                 // 系统规则不可删除
	Description string        `gorm:"type:text" json:"description"`                   // 规则描述
	StartDate  *time.Time      `gorm:"index" json:"start_date"`                       // 生效开始时间
	EndDate    *time.Time      `gorm:"index" json:"end_date"`                         // 生效结束时间
	HitCount   int64           `gorm:"default:0" json:"hit_count"`                    // 命中次数
	LastHitAt  *time.Time      `json:"last_hit_at"`                                   // 最近命中时间
	Version    int            `gorm:"default:1" json:"version"`                       // 规则版本
	Extra      datatypes.JSON `gorm:"type:json" json:"extra"`                        // 扩展配置
}

// RuleLog 规则执行日志
type RuleLog struct {
	gorm.Model
	RuleID    uint           `gorm:"index;not null" json:"rule_id"`
	Rule      Rule           `gorm:"foreignKey:RuleID" json:"rule,omitempty"`
	MatchData datatypes.JSON `gorm:"type:json" json:"match_data"`               // 匹配时的输入数据
	Result    datatypes.JSON `gorm:"type:json" json:"result"`                   // 执行结果
	Actions   datatypes.JSON `gorm:"type:json" json:"actions"`                  // 实际执行的动作
	Duration  int64          `gorm:"default:0" json:"duration"`                   // 耗时(ms)
	Success   bool           `gorm:"default:true" json:"success"`                 // 是否执行成功
	ErrorMsg  string         `gorm:"type:text" json:"error_msg"`                  // 错误信息
	RelType   string         `gorm:"type:varchar(32)" json:"rel_type"`            // 关联类型
	RelID     uint           `gorm:"index" json:"rel_id"`                         // 关联ID
}

// RuleTestResult 规则测试结果
type RuleTestResult struct {
	Matched  bool                   `json:"matched"`
	Actions  []map[string]interface{} `json:"actions,omitempty"`
	Duration int64                  `json:"duration"`
	Error    string                 `json:"error,omitempty"`
	Details  map[string]interface{} `json:"details,omitempty"`
}
