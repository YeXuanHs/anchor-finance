package model

import "time"

// UserLevel 用户等级
type UserLevel struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(50);not null" json:"name"`
	MinAmount   float64   `gorm:"type:decimal(10,2);default:0" json:"min_amount"`
	Discount    float64   `gorm:"type:decimal(5,2);default:100" json:"discount"`
	Priority    int       `gorm:"default:0" json:"priority"`
	Icon        string    `gorm:"type:varchar(255)" json:"icon"`
	Description string    `gorm:"type:text" json:"description"`
	IsDefault   bool      `gorm:"default:false" json:"is_default"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
