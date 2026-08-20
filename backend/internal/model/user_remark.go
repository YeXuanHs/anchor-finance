package model

import "time"

// UserRemark 用户备注
type UserRemark struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	AdminID   uint      `gorm:"index;not null" json:"admin_id"`
	Admin     User      `gorm:"foreignKey:AdminID" json:"admin,omitempty"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Stick     int       `gorm:"default:0;index" json:"stick"` // 置顶排序，越大越靠前
	Type      string    `gorm:"type:varchar(32);default:'general';index" json:"type"` // general, billing, support, vip, risk
	CreatedAt time.Time `json:"created_at"`
}
