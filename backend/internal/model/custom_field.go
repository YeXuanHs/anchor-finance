package model

import "gorm.io/gorm"

// CustomField 自定义字段
type CustomField struct {
	gorm.Model
	Name        string `gorm:"type:varchar(64);not null" json:"name"`
	Label       string `gorm:"type:varchar(128);not null" json:"label"`
	Type        string `gorm:"type:varchar(16);not null" json:"type"` // text/textarea/select/checkbox/radio/file
	Options     string `gorm:"type:jsonb" json:"options"`
	Default     string `gorm:"type:varchar(256)" json:"default"`
	Required    bool   `gorm:"default:false" json:"required"`
	SortOrder   int    `gorm:"default:0" json:"sort_order"`
	RelType     string `gorm:"type:varchar(32);index" json:"rel_type"` // product/service/order
	RelID       uint   `gorm:"index" json:"rel_id"`
	Description string `gorm:"type:varchar(256)" json:"description"`
}
