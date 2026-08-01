package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Nav 导航菜单项
type Nav struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(64);not null" json:"name"`           // 菜单名称
	URL       string         `gorm:"type:varchar(256)" json:"url"`                    // 跳转链接
	ParentID  uint           `gorm:"default:0;index" json:"pid"`                      // 父级ID
	Order     int            `gorm:"default:0" json:"order"`                          // 排序
	FaIcon    string         `gorm:"type:varchar(128)" json:"fa_icon"`                // 图标
	MenuType  int            `gorm:"default:1;comment:1=用户中心 2=www头部 3=www尾部" json:"menu_type"` // 菜单类型
	NavType   int            `gorm:"default:0;comment:0=系统页面 1=URL 2=产品分组" json:"nav_type"` // 导航类型
	MenuID    uint           `gorm:"default:1;comment:关联menus表" json:"menuid"`      // 菜单组ID
	Lang      datatypes.JSON `gorm:"type:json" json:"lang"`                          // 多语言
	Plugin    string         `gorm:"type:varchar(64)" json:"plugin"`                  // 关联插件
	IsDisplay bool           `gorm:"default:true" json:"is_display"`                  // 是否显示
	CreatedAt int64          `json:"created_at"`
	UpdatedAt int64          `json:"updated_at"`
}

// MenuActive 菜单激活配置
type MenuActive struct {
	ID       uint `gorm:"primaryKey" json:"id"`
	MenuType int  `gorm:"uniqueIndex;comment:1=用户中心 2=www头部 3=www尾部" json:"type"`
	MenuID   uint `json:"menuid"`
}

// MenuItem 菜单项（API返回格式）
type MenuItem struct {
	ID       uint         `json:"id"`
	Name     string       `json:"name"`
	URL      string       `json:"url"`
	PID      uint         `json:"pid"`
	FaIcon   string       `json:"fa_icon"`
	NavType  int          `json:"nav_type,omitempty"`
	Children []*MenuItem  `json:"children,omitempty"`
}

// GetNavTree 获取导航树
func GetNavTree(db *gorm.DB, menuType int) ([]*MenuItem, error) {
	// 获取激活的菜单组
	var active MenuActive
	if err := db.Where("menu_type = ?", menuType).First(&active).Error; err != nil {
		// 没有激活菜单，使用默认menuid=1
		active.MenuID = 1
	}

	// 获取所有导航项
	var navs []Nav
	if err := db.Where("menu_id = ? AND menu_type = ?", active.MenuID, menuType).
		Order("`order` ASC, id ASC").
		Find(&navs).Error; err != nil {
		return nil, err
	}

	// 构建树
	return buildNavTree(navs, 0), nil
}

// buildNavTree 构建导航树
func buildNavTree(navs []Nav, parentID uint) []*MenuItem {
	var items []*MenuItem
	for _, nav := range navs {
		if nav.ParentID == parentID {
			item := &MenuItem{
				ID:     nav.ID,
				Name:   nav.Name,
				URL:    nav.URL,
				PID:    nav.ParentID,
				FaIcon: nav.FaIcon,
				NavType: nav.NavType,
			}
			item.Children = buildNavTree(navs, nav.ID)
			items = append(items, item)
		}
	}
	return items
}
