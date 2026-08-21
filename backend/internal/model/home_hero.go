package model

import (
	"time"

	"gorm.io/gorm"
)

// HomeHero 首页英雄区模型
type HomeHero struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"size:200;not null" json:"title"`
	Subtitle  string         `gorm:"size:500" json:"subtitle"`
	ImageURL  string         `gorm:"size:500" json:"image_url"`
	LinkURL   string         `gorm:"size:500" json:"link_url"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	Status    string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (HomeHero) TableName() string {
	return "home_heroes"
}
