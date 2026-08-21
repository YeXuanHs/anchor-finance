package model

import (
	"time"

	"gorm.io/gorm"
)

// MediaFile 媒体文件模型
type MediaFile struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:200;not null" json:"name"`
	Path      string         `gorm:"size:500;not null" json:"path"`
	Size      int64          `gorm:"not null" json:"size"`
	MimeType  string         `gorm:"size:100" json:"mime_type"`
	Extension string         `gorm:"size:20" json:"extension"`
	UploadedBy uint          `gorm:"index" json:"uploaded_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (MediaFile) TableName() string {
	return "media_files"
}
