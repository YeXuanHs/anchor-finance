package model

import (
	"time"

	"gorm.io/gorm"
)

// Certification 实名认证
type Certification struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	UserID          uint           `gorm:"uniqueIndex;not null" json:"user_id"`
	Type            string         `gorm:"type:varchar(20);not null" json:"type"` // individual, enterprise
	RealName        string         `gorm:"type:varchar(50);not null" json:"real_name"`
	IDCard          string         `gorm:"type:varchar(50)" json:"id_card"`
	FrontImage      string         `gorm:"type:varchar(255)" json:"front_image"`
	BackImage       string         `gorm:"type:varchar(255)" json:"back_image"`
	HandImage       string         `gorm:"type:varchar(255)" json:"hand_image"`
	EnterpriseName  string         `gorm:"type:varchar(100)" json:"enterprise_name"`
	BusinessLicense string         `gorm:"type:varchar(255)" json:"business_license"`
	Status          int8           `gorm:"type:smallint;default:1" json:"status"` // 1待审核 2已通过 3已拒绝
	RejectReason    string         `gorm:"type:text" json:"reject_reason"`
	ReviewedBy      *uint          `gorm:"index" json:"reviewed_by"`
	ReviewedAt      *time.Time     `json:"reviewed_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	User            *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
