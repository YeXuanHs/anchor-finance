package model

import (
	"time"

	"gorm.io/gorm"
)

// SSLCertificate represents an SSL certificate order.
type SSLCertificate struct {
	gorm.Model
	UserID          uint       `gorm:"index;not null" json:"user_id"`
	Domain          string     `gorm:"type:varchar(255);not null" json:"domain"`
	ProductID       *uint      `gorm:"index" json:"product_id"`
	CertificateType string     `gorm:"type:varchar(32);not null" json:"certificate_type"` // dv, ov, ev, wildcard
	Status          string     `gorm:"type:varchar(32);default:'pending'" json:"status"`   // pending, issued, expired, revoked, cancelled
	Issuer          string     `gorm:"type:varchar(255)" json:"issuer"`
	SerialNumber    string     `gorm:"type:varchar(128)" json:"serial_number"`
	IssueDate       *time.Time `json:"issue_date"`
	ExpiryDate      *time.Time `gorm:"index" json:"expiry_date"`
	AutoRenew       bool       `gorm:"default:false" json:"auto_renew"`
	CSR             string     `gorm:"type:text" json:"csr"`
	PrivateKey      string     `gorm:"type:text" json:"-"`
	Certificate     string     `gorm:"type:text" json:"certificate"`
	CaBundle        string     `gorm:"type:text" json:"ca_bundle"`
	SANs            string     `gorm:"type:text;column:sans" json:"sans"` // JSON array of additional domains
	ValidationType  string     `gorm:"type:varchar(32)" json:"validation_type"` // dns, email, http
	ValidationData  string     `gorm:"type:json" json:"validation_data"`
	OrderID         string     `gorm:"type:varchar(128)" json:"order_id"`
	Price           float64    `gorm:"type:decimal(20,4)" json:"price"`
	Metadata        string     `gorm:"type:json" json:"metadata"`
	User            *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Product         *Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// SSLCertificateType represents available SSL product types for pricing.
type SSLCertificateType struct {
	ID              uint    `gorm:"primaryKey" json:"id"`
	Name            string  `gorm:"type:varchar(128);not null" json:"name"`
	Code            string  `gorm:"type:varchar(32);uniqueIndex;not null" json:"code"` // dv, ov, ev, wildcard
	Description     string  `gorm:"type:text" json:"description"`
	Price           float64 `gorm:"type:decimal(20,4);not null" json:"price"`
	Currency        string  `gorm:"type:varchar(8);default:'CNY'" json:"currency"`
	ValidityDays    int     `gorm:"default:365" json:"validity_days"`
	WildcardSupport bool    `gorm:"default:false" json:"wildcard_support"`
	SortOrder       int     `gorm:"default:0" json:"sort_order"`
	Status          int16   `gorm:"type:smallint;default:1;not null" json:"status"` // 1=enabled 0=disabled
}
