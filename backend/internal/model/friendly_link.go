package model

import "gorm.io/gorm"

// FriendlyLink 友情链接
type FriendlyLink struct {
	gorm.Model
	Name      string `gorm:"type:varchar(128);not null" json:"name"`
	URL       string `gorm:"type:varchar(512);not null" json:"url"`
	Logo      string `gorm:"type:varchar(256)" json:"logo"`
	SortOrder int    `gorm:"default:0;index" json:"sort_order"`
	Status    int16  `gorm:"type:smallint;default:1;index" json:"status"`
	Group     string `gorm:"type:varchar(32);default:default" json:"group"`
}
