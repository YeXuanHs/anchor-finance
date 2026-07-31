package model

import (
	"time"

	"gorm.io/gorm"
)

// CancelRequest 产品取消请求
type CancelRequest struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	HostID    uint           `gorm:"not null;index" json:"host_id"`   // 主机ID
	UID       uint           `gorm:"not null;index" json:"uid"`       // 用户ID
	Type      string         `gorm:"size:20;not null" json:"type"`    // 类型: Immediate立即取消, Endofbilling等待账单周期结束
	Reason    string         `gorm:"size:255" json:"reason"`          // 取消原因
	Status    string         `gorm:"size:20;default:Pending;index" json:"status"` // 状态: Pending, Approved, Rejected, Completed
	AdminID   uint           `gorm:"default:0" json:"admin_id"`       // 处理管理员ID
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Host      *Host          `gorm:"foreignKey:HostID" json:"host,omitempty"`
	User      *User          `gorm:"foreignKey:UID" json:"user,omitempty"`
}

func (CancelRequest) TableName() string {
	return "cancel_requests"
}
