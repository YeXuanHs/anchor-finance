package model

import (
	"time"

	"gorm.io/gorm"
)

// DownloadFile 下载文件
type DownloadFile struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	CategoryID    uint           `gorm:"index;not null" json:"category_id"`
	Title         string         `gorm:"type:varchar(255);not null" json:"title"`
	Description   string         `gorm:"type:text" json:"description"`
	FilePath      string         `gorm:"type:varchar(500);not null" json:"file_path"`
	FileSize      int64          `gorm:"not null" json:"file_size"`
	FileType      string         `gorm:"type:varchar(50)" json:"file_type"`
	DownloadCount int            `gorm:"default:0" json:"download_count"`
	IsPublished   bool           `gorm:"default:true;index" json:"is_published"`
	SortOrder     int            `gorm:"default:0" json:"sort_order"`
	AdminID       uint           `gorm:"index" json:"admin_id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Category      *DownloadCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}
