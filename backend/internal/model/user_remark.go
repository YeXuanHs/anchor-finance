package model

import (
	"time"

	"gorm.io/gorm"
)

// UserRemark 用户备注模型
type UserRemark struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	AdminID   uint           `gorm:"index" json:"admin_id"`
	AdminName string         `gorm:"size:50" json:"admin_name"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (UserRemark) TableName() string {
	return "user_remarks"
}
