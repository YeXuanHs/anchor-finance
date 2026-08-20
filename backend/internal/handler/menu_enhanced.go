package handler

import (
	"net/http"
	"strconv"

	"anchorfinance/internal/service"

	"github.com/gin-gonic/gin"
)

// MenuEnhancedHandler 菜单增强处理器
type MenuEnhancedHandler struct {
	menuSvc *service.MenuEnhancedService
}

// NewMenuEnhancedHandler 创建菜单增强处理器
func NewMenuEnhancedHandler(menuSvc *service.MenuEnhancedService) *MenuEnhancedHandler {
	return &MenuEnhancedHandler{menuSvc: menuSvc}
}

// GetWebNavs 获取网站导航
func (h *MenuEnhancedHandler) GetWebNavs(c *gin.Context) {
	navType := c.Query("type")
	navs, err := h.menuSvc.GetWebNavs(navType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": navs})
}

// CreateWebNav 创建网站导航
func (h *MenuEnhancedHandler) CreateWebNav(c *gin.Context) {
	var nav service.WebNav
	if err := c.ShouldBindJSON(&nav); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.menuSvc.CreateWebNav(&nav); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": nav})
}

// UpdateWebNav 更新网站导航
func (h *MenuEnhancedHandler) UpdateWebNav(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid nav ID"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.menuSvc.UpdateWebNav(uint(id), updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Nav updated"})
}

// DeleteWebNav 删除网站导航
func (h *MenuEnhancedHandler) DeleteWebNav(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid nav ID"})
		return
	}

	if err := h.menuSvc.DeleteWebNav(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Nav deleted"})
}

// GetDefaultSenior 获取默认高级导航
func (h *MenuEnhancedHandler) GetDefaultSenior(c *gin.Context) {
	defaults, err := h.menuSvc.GetDefaultSenior()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": defaults})
}

// AddCustomPage 添加自定义页面
func (h *MenuEnhancedHandler) AddCustomPage(c *gin.Context) {
	var req struct {
		Name    string `json:"name" binding:"required"`
		Slug    string `json:"slug" binding:"required"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.menuSvc.AddCustomPage(req.Name, req.Slug, req.Content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Custom page created"})
}

// AddProductPage 添加产品页面
func (h *MenuEnhancedHandler) AddProductPage(c *gin.Context) {
	var req struct {
		GroupID  uint   `json:"group_id" binding:"required"`
		Name     string `json:"name" binding:"required"`
		Slug     string `json:"slug" binding:"required"`
		Template string `json:"template"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.menuSvc.AddProductPage(req.GroupID, req.Name, req.Slug, req.Template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product page created"})
}

// GetSystemNav 获取系统导航
func (h *MenuEnhancedHandler) GetSystemNav(c *gin.Context) {
	navs, err := h.menuSvc.GetSystemNav()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": navs})
}

// GetProductList 获取产品菜单列表
func (h *MenuEnhancedHandler) GetProductList(c *gin.Context) {
	menus, err := h.menuSvc.GetProductList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": menus})
}

// CreateWebPage 创建网页
func (h *MenuEnhancedHandler) CreateWebPage(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.menuSvc.CreateWebPage(data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Web page created"})
}

// GetMenuType 获取菜单类型
func (h *MenuEnhancedHandler) GetMenuType(c *gin.Context) {
	types, err := h.menuSvc.GetMenuType()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": types})
}

// GetOtherMenu 获取其他菜单
func (h *MenuEnhancedHandler) GetOtherMenu(c *gin.Context) {
	menuType := c.Query("type")
	navs, err := h.menuSvc.GetOtherMenu(menuType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": navs})
}

// DelTwoMenu 删除二级菜单
func (h *MenuEnhancedHandler) DelTwoMenu(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu ID"})
		return
	}

	if err := h.menuSvc.DelTwoMenu(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu deleted"})
}

// GetTypeAllMenu 获取所有类型菜单
func (h *MenuEnhancedHandler) GetTypeAllMenu(c *gin.Context) {
	menus, err := h.menuSvc.GetTypeAllMenu()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": menus})
}

// EditMenuActive 编辑菜单激活状态
func (h *MenuEnhancedHandler) EditMenuActive(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu ID"})
		return
	}

	var req struct {
		Active bool `json:"active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.menuSvc.EditMenuActive(uint(id), req.Active); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu active status updated"})
}

// GetNavType 获取导航类型
func (h *MenuEnhancedHandler) GetNavType(c *gin.Context) {
	types, err := h.menuSvc.GetNavType()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": types})
}

// GetCreateWebData 获取创建网页数据
func (h *MenuEnhancedHandler) GetCreateWebData(c *gin.Context) {
	data, err := h.menuSvc.GetCreateWebData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

// GetLang 获取语言列表
func (h *MenuEnhancedHandler) GetLang(c *gin.Context) {
	langs, err := h.menuSvc.GetLang()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": langs})
}

// DirectDel 直接删除
func (h *MenuEnhancedHandler) DirectDel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu ID"})
		return
	}

	if err := h.menuSvc.DirectDel(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu permanently deleted"})
}

// AddHookMenu 添加Hook菜单
func (h *MenuEnhancedHandler) AddHookMenu(c *gin.Context) {
	var menu service.HookMenu
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.menuSvc.AddHookMenu(&menu); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": menu})
}

// DelHookMenu 删除Hook菜单
func (h *MenuEnhancedHandler) DelHookMenu(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu ID"})
		return
	}

	if err := h.menuSvc.DelHookMenu(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hook menu deleted"})
}

// GetHookMenus 获取Hook菜单列表
func (h *MenuEnhancedHandler) GetHookMenus(c *gin.Context) {
	menus, err := h.menuSvc.GetHookMenus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": menus})
}

// GetOneNavs 获取单个导航
func (h *MenuEnhancedHandler) GetOneNavs(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid nav ID"})
		return
	}

	nav, err := h.menuSvc.GetOneNavs(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": nav})
}

// SaveLinks 保存链接
func (h *MenuEnhancedHandler) SaveLinks(c *gin.Context) {
	navIDStr := c.Param("nav_id")
	navID, err := strconv.ParseUint(navIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid nav ID"})
		return
	}

	var links []map[string]interface{}
	if err := c.ShouldBindJSON(&links); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.menuSvc.SaveLinks(uint(navID), links); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Links saved"})
}

// DeleteLinks 删除链接
func (h *MenuEnhancedHandler) DeleteLinks(c *gin.Context) {
	navIDStr := c.Param("nav_id")
	navID, err := strconv.ParseUint(navIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid nav ID"})
		return
	}

	if err := h.menuSvc.DeleteLinks(uint(navID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Links deleted"})
}

// AllLinks 获取所有链接
func (h *MenuEnhancedHandler) AllLinks(c *gin.Context) {
	links, err := h.menuSvc.AllLinks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": links})
}

// GetWebNavList 获取网站导航列表（管理用）
func (h *MenuEnhancedHandler) GetWebNavList(c *gin.Context) {
	navs, err := h.menuSvc.GetWebNavList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": navs})
}

// SetWebNavList 设置网站导航列表
func (h *MenuEnhancedHandler) SetWebNavList(c *gin.Context) {
	var navs []service.WebNav
	if err := c.ShouldBindJSON(&navs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.menuSvc.SetWebNavList(navs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Web nav list updated"})
}
