package model

import (
	"time"

	"gorm.io/gorm"
)

// Theme 主题模板模型
type Theme struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Code        string         `gorm:"size:50;uniqueIndex;not null" json:"code"`
	Description string         `gorm:"size:500" json:"description"`
	Thumbnail   string         `gorm:"size:500" json:"thumbnail"`
	IsDefault   bool           `gorm:"default:false" json:"is_default"`
	Status      string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Theme) TableName() string {
	return "themes"
}
