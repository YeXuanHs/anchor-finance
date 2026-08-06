package model

import (
	"time"

	"gorm.io/datatypes"
)

// UserTaste 用户偏好设置
type UserTaste struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"uniqueIndex;not null" json:"user_id"`
	UID           uint           `gorm:"index" json:"uid"`
	Theme         string         `gorm:"type:varchar(32);default:default" json:"theme"`
	Language      string         `gorm:"type:varchar(16);default:zh-CN" json:"language"`
	Layout        string         `gorm:"type:varchar(32);default:default" json:"layout"`
	TicketRefresh int            `gorm:"default:0" json:"ticket_refresh"`
	Settings      datatypes.JSON `gorm:"type:json" json:"settings"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}
