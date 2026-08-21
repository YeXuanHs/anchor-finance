package model

import (
	"time"

	"gorm.io/gorm"
)

// MemberLevel 会员等级模型
type MemberLevel struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:50;not null" json:"name"`
	Description string         `gorm:"size:500" json:"description"`
	Discount    float64        `gorm:"type:decimal(5,2);default:100" json:"discount"` // 折扣百分比，100表示无折扣
	MinPoints   int            `gorm:"default:0" json:"min_points"` // 最低积分
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	Status      string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (MemberLevel) TableName() string {
	return "member_levels"
}
