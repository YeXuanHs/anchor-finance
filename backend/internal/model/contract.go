package model

import (
	"time"

	"gorm.io/gorm"
)

// Contract 合同模型
type Contract struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `gorm:"index;not null" json:"user_id"`
	ContractNo string         `gorm:"size:50;uniqueIndex;not null" json:"contract_no"`
	Title      string         `gorm:"size:200;not null" json:"title"`
	Content    string         `gorm:"type:text" json:"content"`
	Type       string         `gorm:"size:50" json:"type"` // service, product, custom
	Status     string         `gorm:"size:20;default:draft" json:"status"` // draft, pending, signed, cancelled
	StartDate  *time.Time     `json:"start_date"`
	EndDate    *time.Time     `json:"end_date"`
	SignedAt   *time.Time     `json:"signed_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Contract) TableName() string {
	return "contracts"
}

// ContractTemplate 合同模板模型
type ContractTemplate struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Content   string         `gorm:"type:text" json:"content"`
	Variables string         `gorm:"type:text" json:"variables"` // JSON格式的变量说明
	Status    string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (ContractTemplate) TableName() string {
	return "contract_templates"
}
