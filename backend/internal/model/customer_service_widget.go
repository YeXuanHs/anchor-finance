package model

import "gorm.io/gorm"

// CustomerServiceWidget 客服悬浮窗配置
// 移植自 anchor_cloud_finance_pro
type CustomerServiceWidget struct {
	gorm.Model
	Name      string `gorm:"type:varchar(100);not null" json:"name"`
	Type      string `gorm:"type:varchar(20);not null;comment:qq/wechat/telegram/phone/email/custom" json:"type"`
	Icon      string `gorm:"type:varchar(255)" json:"icon"`
	Content   string `gorm:"type:varchar(500);comment:QQ号/微信号/手机号等" json:"content"`
	URL       string `gorm:"type:varchar(500);comment:跳转链接" json:"url"`
	SortOrder int    `gorm:"default:0" json:"sort_order"`
	IsActive  bool   `gorm:"default:true" json:"is_active"`
	IsDefault bool   `gorm:"default:false;comment:是否为默认客服" json:"is_default"`
}

// CustomerServiceWidgetSetting 客服悬浮窗全局设置
type CustomerServiceWidgetSetting struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Enabled      bool   `gorm:"default:true" json:"enabled"`
	Position     string `gorm:"type:varchar(20);default:right;comment:left/right" json:"position"`
	OffsetBottom int    `gorm:"default:100" json:"offset_bottom"`
	ThemeColor   string `gorm:"type:varchar(20);default:#1890ff" json:"theme_color"`
	ShowOnMobile bool   `gorm:"default:true" json:"show_on_mobile"`
	Title        string `gorm:"type:varchar(100);default:联系客服" json:"title"`
	WelcomeText  string `gorm:"type:varchar(500);default:您好，请问有什么可以帮助您？" json:"welcome_text"`
	WorkingHours string `gorm:"type:varchar(100);comment:工作时间说明" json:"working_hours"`
}
