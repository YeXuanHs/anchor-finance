package model

import "time"

// AdminNote 管理员备注
type AdminNote struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	AdminID   uint      `gorm:"index;not null" json:"admin_id"`
	Admin     Admin     `gorm:"foreignKey:AdminID" json:"admin,omitempty"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	IsPrivate bool      `gorm:"default:true" json:"is_private"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName overrides the default table name.
func (AdminNote) TableName() string {
	return "admin_notes"
}
