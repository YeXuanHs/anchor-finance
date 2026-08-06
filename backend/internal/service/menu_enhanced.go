package service

import (
	"encoding/json"
	"fmt"
	"time"

	"anchorfinance/pkg/logger"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// MenuEnhancedService 菜单增强服务
type MenuEnhancedService struct {
	db *gorm.DB
}

// NewMenuEnhancedService 创建菜单增强服务
func NewMenuEnhancedService(db *gorm.DB) *MenuEnhancedService {
	return &MenuEnhancedService{db: db}
}

// WebNav 网站导航
type WebNav struct {
	ID        uint           `gorm:"primaryKey"`
	Name      string         `gorm:"size:128;not null"`
	Slug      string         `gorm:"size:64;uniqueIndex"`
	Type      string         `gorm:"size:32"` // header/footer/sidebar
	ParentID  *uint          `gorm:"index"`
	URL       string         `gorm:"size:512"`
	Icon      string         `gorm:"size:64"`
	Target    string         `gorm:"size:16"` // _self/_blank
	SortOrder int
	Enabled   bool           `gorm:"default:true"`
	Children  []WebNav       `gorm:"foreignKey:ParentID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// HookMenu Hook菜单
type HookMenu struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:128"`
	HookPoint string `gorm:"size:64"` // header/footer/sidebar/content
	Content   string `gorm:"type:text"`
	Position  string `gorm:"size:32"`
	SortOrder int
	Enabled   bool   `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ProductMenu 产品菜单
type ProductMenu struct {
	ID         uint   `gorm:"primaryKey"`
	GroupID    uint   `gorm:"index"`
	Name       string `gorm:"size:128"`
	Slug       string `gorm:"size:64;uniqueIndex"`
	Icon       string `gorm:"size:64"`
	Template   string `gorm:"size:64"`
	ShowInNav  bool   `gorm:"default:true"`
	ShowInHome bool   `gorm:"default:true"`
	SortOrder  int
	Enabled    bool   `gorm:"default:true"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// GetWebNavs 获取网站导航列表
func (s *MenuEnhancedService) GetWebNavs(navType string) ([]WebNav, error) {
	var navs []WebNav
	query := s.db.Where("enabled = ?", true)

	if navType != "" {
		query = query.Where("type = ?", navType)
	}

	err := query.Order("sort_order").Find(&navs).Error
	if err != nil {
		return nil, err
	}

	// 构建树形结构
	return s.buildNavTree(navs, nil), nil
}

// buildNavTree 构建导航树
func (s *MenuEnhancedService) buildNavTree(navs []WebNav, parentID *uint) []WebNav {
	var result []WebNav
	for _, nav := range navs {
		if (parentID == nil && nav.ParentID == nil) ||
			(parentID != nil && nav.ParentID != nil && *nav.ParentID == *parentID) {
			children := s.buildNavTree(navs, &nav.ID)
			nav.Children = children
			result = append(result, nav)
		}
	}
	return result
}

// CreateWebNav 创建网站导航
func (s *MenuEnhancedService) CreateWebNav(nav *WebNav) error {
	return s.db.Create(nav).Error
}

// UpdateWebNav 更新网站导航
func (s *MenuEnhancedService) UpdateWebNav(id uint, updates map[string]interface{}) error {
	return s.db.Model(&WebNav{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteWebNav 删除网站导航
func (s *MenuEnhancedService) DeleteWebNav(id uint) error {
	// 检查是否有子菜单
	var count int64
	s.db.Model(&WebNav{}).Where("parent_id = ?", id).Count(&count)
	if count > 0 {
		return fmt.Errorf("cannot delete menu with children")
	}
	return s.db.Delete(&WebNav{}, id).Error
}

// GetDefaultSenior 获取默认高级导航
func (s *MenuEnhancedService) GetDefaultSenior() ([]map[string]interface{}, error) {
	// 返回系统默认的导航结构
	defaults := []map[string]interface{}{
		{"name": "首页", "url": "/", "icon": "home", "sort_order": 1},
		{"name": "产品", "url": "/products", "icon": "server", "sort_order": 2},
		{"name": "解决方案", "url": "/solutions", "icon": "lightbulb", "sort_order": 3},
		{"name": "新闻", "url": "/news", "icon": "newspaper", "sort_order": 4},
		{"name": "帮助中心", "url": "/help", "icon": "question-circle", "sort_order": 5},
		{"name": "联系我们", "url": "/contact", "icon": "envelope", "sort_order": 6},
	}
	return defaults, nil
}

// AddCustomPage 添加自定义页面
func (s *MenuEnhancedService) AddCustomPage(name, slug, content string) error {
	// 创建自定义页面
	page := map[string]interface{}{
		"name":    name,
		"slug":    slug,
		"content": content,
		"type":    "page",
	}

	logger.Info("Creating custom page", "name", name, "slug", slug)

	// 存入数据库
	return s.db.Table("custom_pages").Create(page).Error
}

// AddProductPage 添加产品页面
func (s *MenuEnhancedService) AddProductPage(groupID uint, name, slug, template string) error {
	menu := &ProductMenu{
		GroupID:    groupID,
		Name:       name,
		Slug:       slug,
		Template:   template,
		ShowInNav:  true,
		ShowInHome: true,
		Enabled:    true,
	}
	return s.db.Create(menu).Error
}

// GetSystemNav 获取系统导航
func (s *MenuEnhancedService) GetSystemNav() ([]map[string]interface{}, error) {
	var navs []WebNav
	s.db.Where("type = ? AND enabled = ?", "system", true).Order("sort_order").Find(&navs)

	result := make([]map[string]interface{}, len(navs))
	for i, nav := range navs {
		result[i] = map[string]interface{}{
			"id":    nav.ID,
			"name":  nav.Name,
			"url":   nav.URL,
			"icon":  nav.Icon,
			"target": nav.Target,
		}
	}

	return result, nil
}

// GetProductList 获取产品菜单列表
func (s *MenuEnhancedService) GetProductList() ([]ProductMenu, error) {
	var menus []ProductMenu
	err := s.db.Where("enabled = ? AND show_in_nav = ?", true, true).
		Order("sort_order").Find(&menus).Error
	return menus, err
}

// CreateWebPage 创建网页
func (s *MenuEnhancedService) CreateWebPage(data map[string]interface{}) error {
	return s.db.Table("custom_pages").Create(data).Error
}

// GetMenuType 获取菜单类型
func (s *MenuEnhancedService) GetMenuType() ([]map[string]string, error) {
	types := []map[string]string{
		{"value": "header", "label": "顶部导航"},
		{"value": "footer", "label": "底部导航"},
		{"value": "sidebar", "label": "侧边栏"},
		{"value": "system", "label": "系统菜单"},
		{"value": "product", "label": "产品菜单"},
	}
	return types, nil
}

// GetOtherMenu 获取其他菜单
func (s *MenuEnhancedService) GetOtherMenu(menuType string) ([]WebNav, error) {
	var navs []WebNav
	err := s.db.Where("type = ? AND enabled = ?", menuType, true).
		Order("sort_order").Find(&navs).Error
	return navs, err
}

// DelTwoMenu 删除二级菜单
func (s *MenuEnhancedService) DelTwoMenu(id uint) error {
	return s.DeleteWebNav(id)
}

// GetTypeAllMenu 获取所有类型菜单
func (s *MenuEnhancedService) GetTypeAllMenu() (map[string][]WebNav, error) {
	var navs []WebNav
	s.db.Where("enabled = ?", true).Order("sort_order").Find(&navs)

	result := make(map[string][]WebNav)
	for _, nav := range navs {
		result[nav.Type] = append(result[nav.Type], nav)
	}

	return result, nil
}

// EditMenuActive 编辑菜单激活状态
func (s *MenuEnhancedService) EditMenuActive(id uint, active bool) error {
	return s.db.Model(&WebNav{}).Where("id = ?", id).Update("enabled", active).Error
}

// GetNavType 获取导航类型
func (s *MenuEnhancedService) GetNavType() ([]string, error) {
	return []string{"header", "footer", "sidebar", "system", "product"}, nil
}

// GetCreateWebData 获取创建网页数据
func (s *MenuEnhancedService) GetCreateWebData() (map[string]interface{}, error) {
	var groups []struct {
		ID   uint
		Name string
	}
	s.db.Table("product_groups").Select("id, name").Find(&groups)

	result := map[string]interface{}{
		"product_groups": groups,
		"templates":      []string{"default", "sidebar", "full_width", "landing"},
	}

	return result, nil
}

// GetLang 获取语言列表
func (s *MenuEnhancedService) GetLang() ([]map[string]string, error) {
	return []map[string]string{
		{"code": "zh-CN", "name": "简体中文"},
		{"code": "zh-TW", "name": "繁体中文"},
		{"code": "en", "name": "English"},
	}, nil
}

// DirectDel 直接删除
func (s *MenuEnhancedService) DirectDel(id uint) error {
	return s.db.Unscoped().Delete(&WebNav{}, id).Error
}

// AddHookMenu 添加Hook菜单
func (s *MenuEnhancedService) AddHookMenu(menu *HookMenu) error {
	return s.db.Create(menu).Error
}

// DelHookMenu 删除Hook菜单
func (s *MenuEnhancedService) DelHookMenu(id uint) error {
	return s.db.Delete(&HookMenu{}, id).Error
}

// GetHookMenus 获取Hook菜单列表
func (s *MenuEnhancedService) GetHookMenus() ([]HookMenu, error) {
	var menus []HookMenu
	err := s.db.Order("sort_order").Find(&menus).Error
	return menus, err
}

// GetOneNavs 获取单个导航
func (s *MenuEnhancedService) GetOneNavs(id uint) (*WebNav, error) {
	var nav WebNav
	err := s.db.First(&nav, id).Error
	if err != nil {
		return nil, err
	}
	return &nav, nil
}

// SaveLinks 保存链接
func (s *MenuEnhancedService) SaveLinks(navID uint, links []map[string]interface{}) error {
	// 删除旧链接
	s.db.Where("nav_id = ?", navID).Delete(&struct{ ID uint }{})

	// 创建新链接
	for _, link := range links {
		link["nav_id"] = navID
		s.db.Table("nav_links").Create(link)
	}

	return nil
}

// DeleteLinks 删除链接
func (s *MenuEnhancedService) DeleteLinks(navID uint) error {
	return s.db.Where("nav_id = ?", navID).Delete(&struct{ ID uint }{}).Error
}

// AllLinks 获取所有链接
func (s *MenuEnhancedService) AllLinks() ([]map[string]interface{}, error) {
	var links []map[string]interface{}
	s.db.Table("nav_links").Find(&links)
	return links, nil
}

// GetWebNavList 获取网站导航列表（管理用）
func (s *MenuEnhancedService) GetWebNavList() ([]WebNav, error) {
	var navs []WebNav
	err := s.db.Order("sort_order").Find(&navs).Error
	return navs, err
}

// SetWebNavList 设置网站导航列表
func (s *MenuEnhancedService) SetWebNavList(navs []WebNav) error {
	// 清空现有导航
	s.db.Where("1 = 1").Delete(&WebNav{})

	// 批量创建
	for i := range navs {
		navs[i].SortOrder = i
	}
	return s.db.Create(&navs).Error
}

// JSON helper
func menuToJSON(v interface{}) datatypes.JSON {
	b, _ := json.Marshal(v)
	return datatypes.JSON(b)
}
