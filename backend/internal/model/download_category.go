package model

import "gorm.io/gorm"

// DownloadCategory 下载分类
type DownloadCategory struct {
	gorm.Model
	Name        string `gorm:"type:varchar(128);not null" json:"name"`
	Description string `gorm:"type:varchar(256)" json:"description"`
	ParentID    uint   `gorm:"index;default:0" json:"parent_id"`
	SortOrder   int    `gorm:"default:0" json:"sort_order"`
	Status      int16  `gorm:"type:smallint;default:1" json:"status"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	Icon        string `gorm:"type:varchar(256)" json:"icon"`
}
