package model

import (
	"gorm.io/gorm"
)

// AdminLog 管理员操作审计日志
type AdminLog struct {
	gorm.Model
	AdminID    uint   `gorm:"index;not null" json:"admin_id"`
	Admin      *User  `gorm:"foreignKey:AdminID" json:"admin,omitempty"`
	Action     string `gorm:"type:varchar(64);not null;index" json:"action"`
	Module     string `gorm:"type:varchar(64);index" json:"module"`
	TargetID   uint   `gorm:"index" json:"target_id"`
	TargetType string `gorm:"type:varchar(64);index" json:"target_type"`
	Detail     string `gorm:"type:text" json:"detail"`
	IP         string `gorm:"type:varchar(64)" json:"ip"`
	UserAgent  string `gorm:"type:varchar(512)" json:"user_agent"`
}

// TableName 指定表名
func (AdminLog) TableName() string {
	return "admin_logs"
}
