package model

import (
	"gorm.io/gorm"
)

// LinkKnowledge 知识关联
type LinkKnowledge struct {
	gorm.Model
	Title      string `gorm:"type:varchar(256);not null" json:"title"`
	Content    string `gorm:"type:text" json:"content"`
	LinkCause  uint   `gorm:"index" json:"link_cause"` // 关联原因ID
	Type       string `gorm:"type:varchar(32);index" json:"type"` // article/faq/guide
	Category   string `gorm:"type:varchar(64);index" json:"category"`
	Status     int16  `gorm:"type:smallint;default:1;not null;index" json:"status"`
	SortOrder  int    `gorm:"default:0" json:"sort_order"`
	ViewCount  int    `gorm:"default:0" json:"view_count"`
	HelpfulYes int    `gorm:"default:0" json:"helpful_yes"`
	HelpfulNo  int    `gorm:"default:0" json:"helpful_no"`
}
