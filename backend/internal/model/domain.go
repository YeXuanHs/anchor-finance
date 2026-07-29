package model

import (
	"time"

	"gorm.io/gorm"
)

// Domain 域名模型
type Domain struct {
	gorm.Model
	UserID           uint       `gorm:"index;not null" json:"user_id"`
	DomainName       string     `gorm:"type:varchar(255);not null;index" json:"domain_name"`
	RegistrarID      *uint      `gorm:"index" json:"registrar_id"`
	RegistrationDate *time.Time `json:"registration_date"`
	ExpiryDate       *time.Time `gorm:"index" json:"expiry_date"`
	NextDueDate      *time.Time `json:"next_due_date"`
	Nameservers      string     `gorm:"type:text" json:"nameservers"`
	Status           string     `gorm:"type:varchar(32);default:'active'" json:"status"`
	AutoRenew        bool       `gorm:"default:false" json:"auto_renew"`
	WhoisPrivacy     bool       `gorm:"default:false" json:"whois_privacy"`
	TransferLock     bool       `gorm:"default:true" json:"transfer_lock"`
	EPPCode          string     `gorm:"type:varchar(128)" json:"epp_code"`
	DNSManaged       bool       `gorm:"default:true" json:"dns_managed"`
	Price            float64    `gorm:"type:decimal(20,4)" json:"price"`
	Currency         string     `gorm:"type:varchar(8);default:'CNY'" json:"currency"`
	Metadata         string     `gorm:"type:jsonb" json:"metadata"`
	User             *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// DomainDNSRecord 域名DNS记录
type DomainDNSRecord struct {
	gorm.Model
	DomainID uint    `gorm:"index;not null" json:"domain_id"`
	Type     string  `gorm:"type:varchar(10);not null" json:"type"`
	Name     string  `gorm:"type:varchar(255);not null" json:"name"`
	Value    string  `gorm:"type:text;not null" json:"value"`
	TTL      int     `gorm:"default:3600" json:"ttl"`
	Priority *int    `json:"priority"`
	Domain   *Domain `gorm:"foreignKey:DomainID" json:"domain,omitempty"`
}

// DomainTransfer 域名转移
type DomainTransfer struct {
	gorm.Model
	UserID      uint       `gorm:"index;not null" json:"user_id"`
	DomainName  string     `gorm:"type:varchar(255);not null" json:"domain_name"`
	EPPCode     string     `gorm:"type:varchar(128)" json:"epp_code"`
	Status      string     `gorm:"type:varchar(32);default:'pending'" json:"status"`
	RegistrarID *uint      `json:"registrar_id"`
	Price       float64    `gorm:"type:decimal(20,4)" json:"price"`
	AdminID     *uint      `json:"admin_id"`
	Remark      string     `gorm:"type:text" json:"remark"`
	CompletedAt *time.Time `json:"completed_at"`
	User        *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
