package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	Username     string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"username"`
	Email        string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"type:varchar(255);not null" json:"-"`
	Salt         string         `gorm:"type:varchar(32);not null" json:"-"`
	Nickname     string         `gorm:"type:varchar(64)" json:"nickname"`
	Avatar       string         `gorm:"type:varchar(512)" json:"avatar"`
	Phone        string         `gorm:"type:varchar(20);index" json:"phone"`
	CountryCode  string         `gorm:"type:varchar(8)" json:"country_code"`
	Language     string         `gorm:"type:varchar(16);default:'zh-CN'" json:"language"`
	Timezone     string         `gorm:"type:varchar(64);default:'Asia/Shanghai'" json:"timezone"`
	Currency     string         `gorm:"type:varchar(8);default:'CNY'" json:"currency"`
	Balance      datatypes.Decimal `gorm:"type:decimal(20,4);default:0;not null" json:"balance"`
	Commission   datatypes.Decimal `gorm:"type:decimal(20,4);default:0;not null" json:"commission"`
	GroupID      uint           `gorm:"index;not null;default:1" json:"group_id"`
	Group        UserGroup      `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	Status       int16          `gorm:"type:smallint;default:1;not null;index" json:"status"` // 1=正常 0=禁用 2=待验证
	IsAdmin      bool           `gorm:"default:false" json:"is_admin"`
	IsStaff      bool           `gorm:"default:false" json:"is_staff"`
	LastLoginAt  *time.Time     `json:"last_login_at"`
	LastLoginIP  string         `gorm:"type:varchar(64)" json:"last_login_ip"`
	TwoFactorKey string         `gorm:"type:varchar(128)" json:"-"`
	InviteCode   string         `gorm:"type:varchar(32);index" json:"invite_code"`
	InvitedBy    *uint          `gorm:"index" json:"invited_by"`
	Inviter      *User          `gorm:"foreignKey:InvitedBy" json:"inviter,omitempty"`
	CommissionRate datatypes.Decimal `gorm:"type:decimal(5,4);default:0" json:"commission_rate"`
	Remark       string         `gorm:"type:text" json:"remark"`
	Addresses    []UserAddress  `gorm:"foreignKey:UserID" json:"addresses,omitempty"`
	LoginLogs    []LoginLog     `gorm:"foreignKey:UserID" json:"login_logs,omitempty"`
}

// UserGroup 用户组
type UserGroup struct {
	gorm.Model
	Name           string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"name"`
	Description    string         `gorm:"type:text" json:"description"`
	Discount       datatypes.Decimal `gorm:"type:decimal(5,4);default:1.0000;not null" json:"discount"` // 折扣比例
	CommissionRate datatypes.Decimal `gorm:"type:decimal(5,4);default:0" json:"commission_rate"`
	IsDefault      bool           `gorm:"default:false" json:"is_default"`
	SortOrder      int            `gorm:"default:0" json:"sort_order"`
	Users          []User         `gorm:"foreignKey:GroupID" json:"users,omitempty"`
}

// UserAddress 用户地址
type UserAddress struct {
	gorm.Model
	UserID      uint   `gorm:"index;not null" json:"user_id"`
	User        User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Country     string `gorm:"type:varchar(64);not null" json:"country"`
	Province    string `gorm:"type:varchar(128)" json:"province"`
	City        string `gorm:"type:varchar(128)" json:"city"`
	District    string `gorm:"type:varchar(128)" json:"district"`
	Street      string `gorm:"type:varchar(256)" json:"street"`
	ZipCode     string `gorm:"type:varchar(16)" json:"zip_code"`
	FirstName   string `gorm:"type:varchar(64)" json:"first_name"`
	LastName    string `gorm:"type:varchar(64)" json:"last_name"`
	Company     string `gorm:"type:varchar(128)" json:"company"`
	TaxID       string `gorm:"type:varchar(64)" json:"tax_id"`
	Phone       string `gorm:"type:varchar(32)" json:"phone"`
	IsDefault   bool   `gorm:"default:false" json:"is_default"`
	AddressType string `gorm:"type:varchar(16);default:'personal'" json:"address_type"` // personal/business
}

// LoginLog 登录日志
type LoginLog struct {
	gorm.Model
	UserID    uint   `gorm:"index;not null" json:"user_id"`
	User      User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	IP        string `gorm:"type:varchar(64);not null" json:"ip"`
	UserAgent string `gorm:"type:varchar(512)" json:"user_agent"`
	Location  string `gorm:"type:varchar(256)" json:"location"`
	Device    string `gorm:"type:varchar(128)" json:"device"`
	Status    int16  `gorm:"type:smallint;not null" json:"status"` // 1=成功 0=失败
	Remark    string `gorm:"type:varchar(256)" json:"remark"`
}
