package model

import (
	"time"

	"gorm.io/gorm"
)

// CancelRequest 取消请求模型
type CancelRequest struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	ServiceID uint           `gorm:"index" json:"service_id"`
	Reason    string         `gorm:"type:text" json:"reason"`
	Status    string         `gorm:"size:20;default:pending" json:"status"` // pending, approved, rejected
	Remark    string         `gorm:"size:500" json:"remark"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (CancelRequest) TableName() string {
	return "cancel_requests"
}
