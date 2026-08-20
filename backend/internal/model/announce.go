package model

import (
	"time"

	"gorm.io/gorm"
)

// Announce 公告管理
type Announce struct {
	gorm.Model
	Title     string     `gorm:"type:varchar(256);not null" json:"title"`
	Content   string     `gorm:"type:text" json:"content"`
	Type      string     `gorm:"type:varchar(16);default:notice" json:"type"` // notice/maintenance/urgent
	Priority  int        `gorm:"default:0" json:"priority"`
	Status    int16      `gorm:"type:smallint;default:1;index" json:"status"` // 1=发布 0=草稿
	StartTime *time.Time `gorm:"index" json:"start_time"`
	EndTime   *time.Time `gorm:"index" json:"end_time"`
	Views     int        `gorm:"default:0" json:"views"`
	IsTop     bool       `gorm:"default:false;index" json:"is_top"`
	AuthorID  uint       `gorm:"index" json:"author_id"`
}
