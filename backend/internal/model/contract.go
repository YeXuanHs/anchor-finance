package model

import (
	"time"

	"gorm.io/gorm"
)

// Contract 合同
type Contract struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	ContractNo string         `gorm:"type:varchar(32);uniqueIndex;not null" json:"contract_no"`
	UserID     uint           `gorm:"index;not null" json:"user_id"`
	Title      string         `gorm:"type:varchar(255);not null" json:"title"`
	Content    string         `gorm:"type:text" json:"content"`
	Type       string         `gorm:"type:varchar(50)" json:"type"`               // service, sales, custom
	Status     int8           `gorm:"type:smallint;default:1" json:"status"`      // 1草稿 2待签署 3已签署 4已过期 5已作废
	Amount     float64        `gorm:"type:decimal(10,2)" json:"amount"`
	StartDate  *time.Time     `json:"start_date"`
	EndDate    *time.Time     `json:"end_date"`
	SignedAt   *time.Time     `json:"signed_at"`
	AdminID    uint           `gorm:"index" json:"admin_id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	User       *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
