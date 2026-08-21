package model

import (
	"time"

	"gorm.io/gorm"
)

// Verification 实名认证模型
type Verification struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	Type      string         `gorm:"size:20;not null" json:"type"` // person, company
	Name      string         `gorm:"size:50" json:"name"`
	IDCard    string         `gorm:"size:20" json:"id_card"`
	Company   string         `gorm:"size:100" json:"company"`
	License   string         `gorm:"size:100" json:"license"`
	Status    string         `gorm:"size:20;default:pending" json:"status"` // pending, approved, rejected
	Remark    string         `gorm:"size:500" json:"remark"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Verification) TableName() string {
	return "verifications"
}
