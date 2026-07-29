package model

import "time"

// RunMap 运行映射
type RunMap struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Name      string     `gorm:"type:varchar(64);not null" json:"name"`
	Code      string     `gorm:"type:varchar(32);uniqueIndex;not null" json:"code"`
	Type      string     `gorm:"type:varchar(16);not null" json:"type"` // script/api/webhook
	Config    string     `gorm:"type:jsonb" json:"config"`
	IsEnabled bool       `gorm:"default:true;index" json:"is_enabled"`
	LastRunAt *time.Time `json:"last_run_at"`
	RunCount  int        `gorm:"default:0" json:"run_count"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
