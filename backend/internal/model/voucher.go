package model

import (
	"time"

	"gorm.io/gorm"
)

// Voucher 发票申请
type Voucher struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UID        uint           `gorm:"index;not null" json:"uid"`
	InvoiceID  uint           `gorm:"index;not null;default:0" json:"invoice_id"`
	PostID     uint           `gorm:"index;not null;default:0" json:"post_id"`
	TypeID     uint           `gorm:"index;not null;default:0" json:"type_id"`
	ExpressID  uint           `gorm:"index;not null;default:0" json:"express_id"`
	Amount     float64        `gorm:"type:decimal(10,2);not null;default:0" json:"amount"`
	Status     string         `gorm:"type:varchar(25);not null;default:'Pending'" json:"status"`
	Notes      string         `gorm:"type:varchar(500);default:''" json:"notes"`
	CheckTime  int64          `gorm:"not null;default:0" json:"check_time"`
	CreateTime int64          `gorm:"not null;default:0" json:"create_time"`
	UpdateTime int64          `gorm:"not null;default:0" json:"update_time"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	User       *User          `gorm:"foreignKey:UID" json:"user,omitempty"`
	VoucherType *VoucherType  `gorm:"foreignKey:TypeID" json:"voucher_type,omitempty"`
	Post       *VoucherPost  `gorm:"foreignKey:PostID" json:"post,omitempty"`
	Express    *Express       `gorm:"foreignKey:ExpressID" json:"express,omitempty"`
	Invoice    *Invoice       `gorm:"foreignKey:InvoiceID" json:"invoice,omitempty"`
}

// VoucherType 发票抬头
type VoucherType struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UID         uint           `gorm:"index;not null" json:"uid"`
	Title       string         `gorm:"type:varchar(50);not null;default:''" json:"title"`
	IssueType   string         `gorm:"type:varchar(50);not null;default:'person'" json:"issue_type"`
	InvoiceType string         `gorm:"type:varchar(50);not null;default:'common'" json:"invoice_type"`
	TaxID       string         `gorm:"type:varchar(100);not null;default:''" json:"tax_id"`
	Bank        string         `gorm:"type:varchar(100);not null;default:''" json:"bank"`
	Account     string         `gorm:"type:varchar(100);not null;default:''" json:"account"`
	Address     string         `gorm:"type:varchar(100);not null;default:''" json:"address"`
	Phone       string         `gorm:"type:varchar(100);not null;default:''" json:"phone"`
	CreateTime  int64          `gorm:"not null;default:0" json:"create_time"`
	UpdateTime  int64          `gorm:"not null;default:0" json:"update_time"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// VoucherPost 发票收件地址
type VoucherPost struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UID        uint           `gorm:"index;not null" json:"uid"`
	Username   string         `gorm:"type:varchar(50);not null;default:''" json:"username"`
	Phone      string         `gorm:"type:varchar(50);not null;default:''" json:"phone"`
	Province   string         `gorm:"type:varchar(100);not null;default:''" json:"province"`
	City       string         `gorm:"type:varchar(100);not null;default:''" json:"city"`
	Region     string         `gorm:"type:varchar(100);not null;default:''" json:"region"`
	Detail     string         `gorm:"type:varchar(500);not null;default:''" json:"detail"`
	Post       string         `gorm:"type:varchar(50);not null;default:''" json:"post"`
	IsDefault  bool           `gorm:"column:\"default\";not null;default:false" json:"is_default"`
	CreateTime int64          `gorm:"not null;default:0" json:"create_time"`
	UpdateTime int64          `gorm:"not null;default:0" json:"update_time"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// VoucherInvoice 发票与账单关联
type VoucherInvoice struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	VoucherID uint      `gorm:"index;not null" json:"voucher_id"`
	InvoiceID uint      `gorm:"index;not null" json:"invoice_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (Voucher) TableName() string        { return "voucher" }
func (VoucherType) TableName() string    { return "voucher_type" }
func (VoucherPost) TableName() string    { return "voucher_post" }
func (VoucherInvoice) TableName() string { return "voucher_invoices" }
