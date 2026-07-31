package model

import (
	"time"
)

// CancelReason 取消原因
type CancelReason struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Reason    string `gorm:"size:255;not null" json:"reason"` // 原因
	IsActive  bool   `gorm:"default:true" json:"is_active"`
	SortOrder int    `gorm:"default:0" json:"sort_order"`
}

func (CancelReason) TableName() string {
	return "cancel_reasons"
}

// BaseInfo 首页基本信息
type BaseInfo struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Name      string `gorm:"size:50;not null" json:"name"` // 模块名称
	Desc      string `gorm:"size:255" json:"desc"`         // 模块描述
	Icon      string `gorm:"size:100" json:"icon"`         // 图标
	Content   string `gorm:"type:text" json:"content"`     // 内容
	SortOrder int    `gorm:"default:0" json:"sort_order"`  // 排序
	IsActive  bool   `gorm:"default:true" json:"is_active"` // 是否启用
}

func (BaseInfo) TableName() string {
	return "base_infos"
}

// PushHost 推送主机记录（用于系统间同步主机信息）
type PushHost struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	HostID    uint      `gorm:"not null;index" json:"host_id"`  // 主机ID
	Status    string    `gorm:"size:1;default:0" json:"status"` // 状态: 1成功, 0失败
	URL       string    `gorm:"type:text;not null" json:"url"`  // 推送URL
	PostData  string    `gorm:"type:text" json:"post_data"`     // 推送数据
	Response  string    `gorm:"type:text" json:"response"`      // 响应内容
	PushCount int       `gorm:"default:0" json:"push_count"`    // 推送次数
	LastPush  time.Time `json:"last_push"`                      // 最后推送时间
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (PushHost) TableName() string {
	return "push_hosts"
}
