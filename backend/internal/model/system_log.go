package model

import "time"

// SystemLog 系统日志
type SystemLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Level     string    `gorm:"type:varchar(16);not null;index" json:"level"` // debug/info/warn/error
	Module    string    `gorm:"type:varchar(64);index" json:"module"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	Details   string    `gorm:"type:jsonb" json:"details"`
	UserID    uint      `gorm:"index" json:"user_id"`
	IP        string    `gorm:"type:varchar(45)" json:"ip"`
	UserAgent string    `gorm:"type:varchar(256)" json:"user_agent"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}
