package model

import (
	"gorm.io/gorm"
)

// LinkCause 关联原因
type LinkCause struct {
	gorm.Model
	Name     string `gorm:"type:varchar(128);not null" json:"name"`
	ParentID *uint  `gorm:"index" json:"parent_id"`
	Parent   *LinkCause `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []LinkCause `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	LinkType string `gorm:"type:varchar(32);not null;index" json:"link_type"` // ticket/order/product
	Level    int    `gorm:"default:0" json:"level"`
	Status   int16  `gorm:"type:smallint;default:1;not null;index" json:"status"`
	SortOrder int   `gorm:"default:0" json:"sort_order"`
}
