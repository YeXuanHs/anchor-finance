package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// RenewCycle 续费周期管理
type RenewCycle struct {
	gorm.Model
	UserProductID uint             `gorm:"index;not null" json:"user_product_id"`
	UserProduct   *UserProduct     `gorm:"foreignKey:UserProductID" json:"user_product,omitempty"`
	Cycle         string           `gorm:"type:varchar(32);not null" json:"cycle"`
	Amount        datatypes.Decimal `gorm:"type:decimal(20,4);not null" json:"amount"`
	NextDueDate   *time.Time       `gorm:"index" json:"next_due_date"`
	LastRenewDate *time.Time       `json:"last_renew_date"`
	AutoRenew     bool             `gorm:"default:false;index" json:"auto_renew"`
	Status        string           `gorm:"type:varchar(20);default:'active';index" json:"status"` // active, suspended, cancelled
}

// TableName 指定表名
func (RenewCycle) TableName() string {
	return "renew_cycles"
}
