package model

import (
	"time"

	"gorm.io/gorm"
)

// News 新闻模型
type News struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Title      string         `gorm:"size:200;not null" json:"title"`
	Content    string         `gorm:"type:text" json:"content"`
	CategoryID uint           `gorm:"index" json:"category_id"`
	Author     string         `gorm:"size:50" json:"author"`
	Status     string         `gorm:"size:20;default:draft" json:"status"` // draft, published
	ViewCount  int            `gorm:"default:0" json:"view_count"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (News) TableName() string {
	return "news"
}

// NewsCategory 新闻分类模型
type NewsCategory struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	Status    string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (NewsCategory) TableName() string {
	return "news_categories"
}

// KnowledgeCategory 知识库分类模型
type KnowledgeCategory struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	Status    string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (KnowledgeCategory) TableName() string {
	return "knowledge_categories"
}

// KnowledgeArticle 知识库文章模型
type KnowledgeArticle struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Title      string         `gorm:"size:200;not null" json:"title"`
	Content    string         `gorm:"type:text" json:"content"`
	CategoryID uint           `gorm:"index" json:"category_id"`
	Status     string         `gorm:"size:20;default:draft" json:"status"` // draft, published
	ViewCount  int            `gorm:"default:0" json:"view_count"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (KnowledgeArticle) TableName() string {
	return "knowledge_articles"
}

// Download 下载模型
type Download struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Title      string         `gorm:"size:200;not null" json:"title"`
	FileName   string         `gorm:"size:200" json:"file_name"`
	FileURL    string         `gorm:"size:500" json:"file_url"`
	CategoryID uint           `gorm:"index" json:"category_id"`
	DownloadCount int         `gorm:"default:0" json:"download_count"`
	Status     string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Download) TableName() string {
	return "downloads"
}

// DownloadCategory 下载分类模型
type DownloadCategory struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	Status    string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (DownloadCategory) TableName() string {
	return "download_categories"
}
