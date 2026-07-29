package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// UserTaste 用户偏好设置
type UserTaste struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"uniqueIndex;not null" json:"user_id"`
	Theme     string         `gorm:"type:varchar(32);default:default" json:"theme"` // 主题
	Language  string         `gorm:"type:varchar(16);default:zh-CN" json:"language"` // 语言
	Layout    string         `gorm:"type:varchar(32);default:default" json:"layout"` // 布局
	Settings  datatypes.JSON `gorm:"type:jsonb" json:"settings"` // 其他偏好设置JSON
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
