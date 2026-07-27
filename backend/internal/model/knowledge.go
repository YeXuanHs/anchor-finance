package model

import (
	"time"

	"gorm.io/gorm"
)

// KnowledgeCategory 知识库分类
type KnowledgeCategory struct {
	ID        uint                `gorm:"primaryKey" json:"id"`
	Name      string              `gorm:"type:varchar(100);not null" json:"name"`
	Slug      string              `gorm:"type:varchar(100);uniqueIndex" json:"slug"`
	ParentID  *uint               `gorm:"index" json:"parent_id"`
	SortOrder int                 `gorm:"default:0" json:"sort_order"`
	IsActive  bool                `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
	Parent    *KnowledgeCategory  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children  []KnowledgeCategory `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Articles  []KnowledgeArticle  `gorm:"foreignKey:CategoryID" json:"articles,omitempty"`
}

// KnowledgeArticle 知识库文章
type KnowledgeArticle struct {
	ID           uint               `gorm:"primaryKey" json:"id"`
	CategoryID   uint               `gorm:"index;not null" json:"category_id"`
	Title        string             `gorm:"type:varchar(255);not null" json:"title"`
	Slug         string             `gorm:"type:varchar(255);uniqueIndex" json:"slug"`
	Content      string             `gorm:"type:text" json:"content"`
	Summary      string             `gorm:"type:varchar(500)" json:"summary"`
	Keywords     string             `gorm:"type:varchar(255)" json:"keywords"`
	ViewCount    int                `gorm:"default:0" json:"view_count"`
	HelpCount    int                `gorm:"default:0" json:"help_count"`
	NotHelpCount int                `gorm:"default:0" json:"not_help_count"`
	IsPublished  bool               `gorm:"default:true;index" json:"is_published"`
	SortOrder    int                `gorm:"default:0" json:"sort_order"`
	AdminID      uint               `gorm:"index" json:"admin_id"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
	DeletedAt    gorm.DeletedAt     `gorm:"index" json:"-"`
	Category     *KnowledgeCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}
