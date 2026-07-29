package model

import (
	"gorm.io/gorm"
)

// CancelRequest 账号注销申请
type CancelRequest struct {
	gorm.Model
	UserID  uint   `gorm:"index;not null" json:"user_id"`
	User    *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Reason  string `gorm:"type:text" json:"reason"`
	Status  string `gorm:"type:varchar(20);default:'pending';index" json:"status"` // pending, approved, rejected
	AdminID *uint  `gorm:"index" json:"admin_id"`
	Remark  string `gorm:"type:text" json:"remark"`
}

// TableName 指定表名
func (CancelRequest) TableName() string {
	return "cancel_requests"
}
