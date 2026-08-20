package model

import (
	"time"

	"gorm.io/gorm"
)

// ProductDivert 产品转移
type ProductDivert struct {
	gorm.Model
	ServiceID   uint       `gorm:"index;not null" json:"service_id"`
	FromUserID  uint       `gorm:"index;not null" json:"from_user_id"`
	ToUserID    uint       `gorm:"index;not null" json:"to_user_id"`
	ToEmail     string     `gorm:"type:varchar(128)" json:"to_email"`
	Status      string     `gorm:"type:varchar(16);default:pending" json:"status"` // pending/approved/rejected/completed
	Reason      string     `gorm:"type:text" json:"reason"`
	AdminNote   string     `gorm:"type:text" json:"admin_note"`
	CompletedAt *time.Time `json:"completed_at"`
}
