package model

import (
	"time"

	"gorm.io/gorm"
)

// AIShoppingAssistantConfig AI 购物助手配置
// 移植自 mahiru_ai_shopping
type AIShoppingAssistantConfig struct {
	ID              uint    `gorm:"primaryKey" json:"id"`
	Enabled         bool    `gorm:"default:false" json:"enabled"`
	AIConfigID      uint    `gorm:"not null" json:"ai_config_id"`
	WelcomeMessage  string  `gorm:"type:text;default:您好！我是AI购物助手，可以帮您推荐合适的产品和服务。请问您有什么需求？" json:"welcome_message"`
	SystemPrompt    string  `gorm:"type:text" json:"system_prompt"`
	MaxRecommendations int  `gorm:"default:5;comment:最大推荐产品数" json:"max_recommendations"`
	IncludePricing  bool    `gorm:"default:true;comment:是否展示价格" json:"include_pricing"`
	ShowOnAllPages  bool    `gorm:"default:true;comment:是否在所有页面显示" json:"show_on_all_pages"`
	TriggerKeywords string  `gorm:"type:text;comment:触发关键词，逗号分隔" json:"trigger_keywords"`
}

// AIShoppingSession AI 购物助手会话
type AIShoppingSession struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	SessionID string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"session_id"`
	Context   string    `gorm:"type:text;comment:会话上下文JSON" json:"context"`
	Status    string    `gorm:"type:varchar(20);default:active;comment:active/closed" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ClosedAt  *time.Time `json:"closed_at"`
}

// AIShoppingMessage AI 购物助手消息
type AIShoppingMessage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SessionID string    `gorm:"type:varchar(64);index;not null" json:"session_id"`
	Role      string    `gorm:"type:varchar(20);not null;comment:user/assistant/system" json:"role"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Products  string    `gorm:"type:text;comment:推荐的产品JSON" json:"products"`
	CreatedAt time.Time `json:"created_at"`
}

// ProductCatalogConfig 商品目录配置
// 移植自 mahiru_ai_shopping
type ProductCatalogConfig struct {
	ID               uint   `gorm:"primaryKey" json:"id"`
	LayoutStyle      string `gorm:"type:varchar(20);default:grid;comment:grid/list/card" json:"layout_style"`
	ShowFilters      bool   `gorm:"default:true" json:"show_filters"`
	ShowComparison   bool   `gorm:"default:true;comment:显示产品对比" json:"show_comparison"`
	ShowReviews      bool   `gorm:"default:true;comment:显示用户评价" json:"show_reviews"`
	ShowTechSpecs    bool   `gorm:"default:true;comment:显示技术规格" json:"show_tech_specs"`
	EnableSort       bool   `gorm:"default:true" json:"enable_sort"`
	DefaultSort      string `gorm:"type:varchar(30);default:recommend;comment:recommend/price_asc/price_desc/newest/popular" json:"default_sort"`
	ProductsPerPage  int    `gorm:"default:12" json:"products_per_page"`
	ShowCategoryTree bool   `gorm:"default:true" json:"show_category_tree"`
}
