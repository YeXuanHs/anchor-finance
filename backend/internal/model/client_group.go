package model

import (
	"time"

	"gorm.io/datatypes"
)

// ClientGroup 客户分组
type ClientGroup struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Name            string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Description     string         `gorm:"type:text" json:"description"`
	Discount        float64        `gorm:"type:decimal(5,4);default:1.0000;not null" json:"discount"`
	DiscountPercent float64        `gorm:"type:decimal(5,2);default:0" json:"discount_percent"` // 折扣百分比 0-100
	TaxRate         float64        `gorm:"type:decimal(5,2);default:0" json:"tax_rate"`          // 税率百分比
	AutoAssignRule  datatypes.JSON `gorm:"type:jsonb" json:"auto_assign_rule"`                    // 自动分配规则
	Priority        int            `gorm:"default:0;index" json:"priority"`
	IsActive        bool           `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

// ClientGroupMember 客户分组成员关联
type ClientGroupMember struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	GroupID   uint      `gorm:"index;not null" json:"group_id"`
	Group     ClientGroup `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
