package model

import (
	"time"

	"gorm.io/gorm"
)

// FriendlyLink 友情链接模型
type FriendlyLink struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	URL       string         `gorm:"size:500;not null" json:"url"`
	Logo      string         `gorm:"size:500" json:"logo"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	Status    string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (FriendlyLink) TableName() string {
	return "friendly_links"
}
