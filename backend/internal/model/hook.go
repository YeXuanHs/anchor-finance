package model

import (
	"time"

	"gorm.io/gorm"
)

// HookDefinition Hook点定义模型
// zjmf定义了100+个Hook点，Go端触发时查此表确认Hook存在
type HookDefinition struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;uniqueIndex;not null" json:"name"` // hook名称如admin_login, ticket_open等
	Description string         `gorm:"size:500" json:"description"`
	Module      string         `gorm:"size:50" json:"module"` // 模块：ticket/order/host/admin等
	Params      string         `gorm:"type:text" json:"params"` // JSON参数说明
	Status      string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (HookDefinition) TableName() string { return "hook_definitions" }

// HookBinding Hook绑定模型
// 记录哪些Hook绑定了哪些插件
type HookBinding struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	HookName  string         `gorm:"size:100;index;not null" json:"hook_name"` // 对应hook_definitions.name
	PluginID  uint           `gorm:"index;not null" json:"plugin_id"`         // 对应plugins.id
	Priority  int            `gorm:"default:0" json:"priority"`               // 优先级
	Status    string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (HookBinding) TableName() string { return "hook_bindings" }
