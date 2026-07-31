package model

import "time"

// SystemConfig is a key-value setting stored in the system_settings table.
type SystemConfig struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Key       string    `gorm:"type:varchar(128);uniqueIndex;not null" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	Group     string    `gorm:"type:varchar(64);default:general" json:"group"`
	Name      string    `gorm:"type:varchar(128)" json:"name"`
	Type      string    `gorm:"type:varchar(32);default:string" json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (SystemConfig) TableName() string {
	return "system_settings"
}
