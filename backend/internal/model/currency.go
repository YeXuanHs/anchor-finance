package model

import "time"

// Currency 货币
type Currency struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Code      string    `gorm:"type:varchar(10);uniqueIndex;not null" json:"code"`
	Name      string    `gorm:"type:varchar(50);not null" json:"name"`
	Symbol    string    `gorm:"type:varchar(10);not null" json:"symbol"`
	Rate      float64   `gorm:"type:decimal(10,4);default:1" json:"rate"`
	IsDefault bool      `gorm:"default:false" json:"is_default"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	Precision int       `gorm:"default:2" json:"precision"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
