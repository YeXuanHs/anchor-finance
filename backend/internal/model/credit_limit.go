package model

import "time"

// CreditLimit 信用额度
type CreditLimit struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	Limit       float64   `gorm:"type:decimal(10,2);default:0" json:"limit"`
	Used        float64   `gorm:"type:decimal(10,2);default:0" json:"used"`
	Available   float64   `gorm:"type:decimal(10,2);default:0" json:"available"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	User        *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// CreditLog 信用额度变动日志
type CreditLog struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`
	Type        string    `gorm:"type:varchar(20);not null" json:"type"` // adjust, use, repay
	Amount      float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	Balance     float64   `gorm:"type:decimal(10,2);not null" json:"balance"`
	RelatedID   uint      `gorm:"index" json:"related_id"`
	RelatedType string    `gorm:"type:varchar(20)" json:"related_type"`
	AdminID     *uint     `gorm:"index" json:"admin_id"`
	Remark      string    `gorm:"type:text" json:"remark"`
	CreatedAt   time.Time `json:"created_at"`
}
