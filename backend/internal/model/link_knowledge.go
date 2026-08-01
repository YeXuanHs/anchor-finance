package model

import (
	"gorm.io/gorm"
)

// LinkKnowledge 知识关联
type LinkKnowledge struct {
	gorm.Model
	Title      string `gorm:"type:varchar(256);not null" json:"title"`
	Content    string `gorm:"type:text" json:"content"`
	Reply      string `gorm:"type:text" json:"reply"`
	LinkCause  uint   `gorm:"index" json:"link_cause"` // 关联原因ID
	Type       string `gorm:"type:varchar(32);index" json:"type"` // 1=文本回复 2=图片回复
	Module     string `gorm:"type:varchar(64)" json:"module"`
	Category   string `gorm:"type:varchar(64);index" json:"category"`
	Status     int16  `gorm:"type:smallint;default:1;not null;index" json:"status"`
	SortOrder  int    `gorm:"default:0" json:"sort_order"`
	ViewCount  int    `gorm:"default:0" json:"view_count"`
	HelpfulYes int    `gorm:"default:0" json:"helpful_yes"`
	HelpfulNo  int    `gorm:"default:0" json:"helpful_no"`
}

// LinkCause 知识库分类
type LinkCause struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Name      string `gorm:"type:varchar(128);not null" json:"name"`
	PID       uint   `gorm:"index;default:0" json:"pid"`
	LevelView string `gorm:"type:varchar(256)" json:"level_view"`
	Status    int16  `gorm:"type:smallint;default:1" json:"status"`
	Order     int    `gorm:"default:0" json:"order"`
}

func (LinkCause) TableName() string {
	return "link_cause"
}

// LinkKeyword 知识库关键词
type LinkKeyword struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Keyword    string `gorm:"type:varchar(128);not null" json:"keyword"`
	Belong     string `gorm:"type:varchar(32);not null" json:"belong"` // knowledge/cause
	RelID      uint   `gorm:"index;not null" json:"relid"`
	Status     int16  `gorm:"type:smallint;default:1" json:"status"`
	CreateTime int64  `json:"create_time"`
}

func (LinkKeyword) TableName() string {
	return "link_keywords"
}
