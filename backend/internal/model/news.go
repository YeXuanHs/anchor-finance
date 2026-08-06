package model

import (
	"time"

	"gorm.io/gorm"
)

// NewsCategory 新闻分类
type NewsCategory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ParentID  uint      `gorm:"index;default:0" json:"parent_id"`
	Title     string    `gorm:"type:varchar(100);not null" json:"title"`
	Name      string    `gorm:"type:varchar(50);not null" json:"name"`
	Slug      string    `gorm:"type:varchar(50);uniqueIndex" json:"slug"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// News 新闻/文章
type News struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CategoryID  uint           `gorm:"index;not null" json:"category_id"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Slug        string         `gorm:"type:varchar(255);uniqueIndex" json:"slug"`
	Summary     string         `gorm:"type:varchar(500)" json:"summary"`
	Content     string         `gorm:"type:text" json:"content"`
	CoverImage  string         `gorm:"type:varchar(255)" json:"cover_image"`
	Keywords    string         `gorm:"type:varchar(255)" json:"keywords"`
	ViewCount   int            `gorm:"default:0" json:"view_count"`
	IsPublished bool           `gorm:"default:true;index" json:"is_published"`
	IsSticky    bool           `gorm:"default:false" json:"is_sticky"`
	PublishedAt *time.Time     `json:"published_at"`
	AdminID     uint           `gorm:"index" json:"admin_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Category    *NewsCategory  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}
