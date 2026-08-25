package model

import (
	"time"

	"gorm.io/gorm"
)

// HomeHero 首页英雄区模型（MD 11.3：JSON字段存储slides+features）
// 单条记录，config字段为JSON，包含slides数组和features数组
type HomeHero struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;default:default" json:"name"`    // 配置名称
	Config    string         `gorm:"type:text;not null" json:"config"`         // JSON格式的Hero配置
	IsDefault bool           `gorm:"default:true" json:"is_default"`          // 是否默认配置
	Status    string         `gorm:"size:20;default:active" json:"status"`    // active, disabled
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (HomeHero) TableName() string {
	return "home_heroes"
}

// HeroSlide 轮播幻灯片（MD 11.2）
type HeroSlide struct {
	Key            string `json:"key"`             // 唯一标识
	RailTitle      string `json:"rail_title"`      // 左侧标签栏文字
	Title          string `json:"title"`           // 大标题
	Desc           string `json:"desc"`            // 描述文字
	PrimaryText    string `json:"primary_text"`    // 主按钮文字
	PrimaryPath    string `json:"primary_path"`    // 主按钮跳转路径
	SecondaryText  string `json:"secondary_text"`  // 副按钮文字
	SecondaryPath  string `json:"secondary_path"`  // 副按钮跳转路径
	Video          string `json:"video"`           // 背景视频URL
	Poster         string `json:"poster"`          // 视频封面图URL
	Ribbon         string `json:"ribbon"`          // 角标文字
	RibbonType     string `json:"ribbon_type"`     // 角标类型：new/hot/warm
}

// HeroFeature 功能卡片（MD 11.2）
type HeroFeature struct {
	Key   string `json:"key"`   // 唯一标识
	Kicker string `json:"kicker"` // 小标签
	Title string `json:"title"` // 标题
	Desc  string `json:"desc"`  // 描述
	Path  string `json:"path"`  // 跳转链接
}

// HomeHeroConfig Hero配置数据结构（MD 11.3）
type HomeHeroConfig struct {
	Slides   []HeroSlide   `json:"slides"`
	Features []HeroFeature `json:"features"`
}
