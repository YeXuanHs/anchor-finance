package model

import (
	"time"

	"gorm.io/gorm"
)

// Banner 轮播图
type Banner struct {
	gorm.Model
	Title       string     `gorm:"type:varchar(256);not null" json:"title"`
	Description string     `gorm:"type:varchar(512)" json:"description"`
	Type        string     `gorm:"type:varchar(16);not null;default:'image'" json:"type"` // image/video
	MediaURL    string     `gorm:"type:varchar(512);not null" json:"media_url"`
	LinkURL     string     `gorm:"type:varchar(512)" json:"link_url"`
	ButtonText  string     `gorm:"type:varchar(64)" json:"button_text"`
	OpenNew     bool       `gorm:"default:false" json:"open_new"` // 是否新窗口打开
	SortOrder   int        `gorm:"default:0;index" json:"sort_order"`
	Status      int16      `gorm:"type:smallint;default:1;not null;index" json:"status"` // 1=启用 0=禁用
	StartTime   *time.Time `gorm:"index" json:"start_time"` // 定时开始，null=立即
	EndTime     *time.Time `gorm:"index" json:"end_time"`   // 定时结束，null=不截止
	Position    string     `gorm:"type:varchar(32);default:'home';index" json:"position"` // home/product/promotion
	ClickCount  int        `gorm:"default:0" json:"click_count"`
}
