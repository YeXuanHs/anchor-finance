package model

import (
	"time"

	"gorm.io/gorm"
)

// CustomField 自定义字段
type CustomField struct {
	gorm.Model
	Name        string `gorm:"type:varchar(128);not null" json:"name"`
	Label       string `gorm:"type:varchar(256);not null" json:"label"`
	Type        string `gorm:"type:varchar(32);not null" json:"type"` // text/textarea/select/checkbox/radio/file/date/number
	Group       string `gorm:"type:varchar(32);index" json:"group"`   // product/cart/client/host
	Required    bool   `gorm:"default:false" json:"required"`
	DefaultVal  string `gorm:"type:text" json:"default_val"`
	Options     string `gorm:"type:json" json:"options"` // JSON array for select/radio/checkbox
	Validation  string `gorm:"type:varchar(256)" json:"validation"`
	Placeholder string `gorm:"type:varchar(256)" json:"placeholder"`
	HelpText    string `gorm:"type:varchar(512)" json:"help_text"`
	SortOrder   int    `gorm:"default:0" json:"sort_order"`
	Enabled     bool   `gorm:"default:true" json:"enabled"`
}

// CustomFieldValue 自定义字段值
type CustomFieldValue struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	FieldID   uint           `gorm:"index" json:"field_id"`
	OwnerID   uint           `gorm:"index" json:"owner_id"`
	OwnerType string         `gorm:"type:varchar(32)" json:"owner_type"` // product/cart/client/host
	Value     string         `gorm:"type:text" json:"value"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// CustomFieldGroup 自定义字段分组
type CustomFieldGroup struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(128);uniqueIndex" json:"name"`
	Label     string    `gorm:"type:varchar(256)" json:"label"`
	Type      string    `gorm:"type:varchar(32)" json:"type"` // product/cart/client/host
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	Enabled   bool      `gorm:"default:true" json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
}
