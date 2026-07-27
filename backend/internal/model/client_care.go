package model

import "time"

// ClientCareRule 客户关怀规则
type ClientCareRule struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"type:varchar(100);not null" json:"name"`
	Type       string    `gorm:"type:varchar(50);not null" json:"type"` // birthday, expire, inactive, custom
	Condition  JSON      `gorm:"type:jsonb" json:"condition"`
	TemplateID uint      `gorm:"index" json:"template_id"`
	Channel    string    `gorm:"type:varchar(20)" json:"channel"`    // email, sms, wechat
	DaysBefore int       `gorm:"default:0" json:"days_before"`      // 提前几天
	IsActive   bool      `gorm:"default:true" json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ClientCareLog 客户关怀日志
type ClientCareLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	RuleID    uint      `gorm:"index" json:"rule_id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Channel   string    `gorm:"type:varchar(20)" json:"channel"`
	Status    string    `gorm:"type:varchar(20)" json:"status"` // sent, failed
	CreatedAt time.Time `json:"created_at"`
}
