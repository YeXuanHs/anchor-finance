package model

import "gorm.io/gorm"

// Contact 联系人
type Contact struct {
	gorm.Model
	UserID    uint   `gorm:"index;not null" json:"user_id"`
	Name      string `gorm:"type:varchar(64);not null" json:"name"`
	Email     string `gorm:"type:varchar(128)" json:"email"`
	Phone     string `gorm:"type:varchar(20)" json:"phone"`
	Company   string `gorm:"type:varchar(128)" json:"company"`
	Address   string `gorm:"type:varchar(256)" json:"address"`
	City      string `gorm:"type:varchar(64)" json:"city"`
	State     string `gorm:"type:varchar(64)" json:"state"`
	Zip       string `gorm:"type:varchar(16)" json:"zip"`
	Country   string `gorm:"type:varchar(8)" json:"country"`
	IsDefault bool   `gorm:"default:false;index" json:"is_default"`
	Type      string `gorm:"type:varchar(16);default:personal" json:"type"` // personal/company
}
