package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/anchorfinance/backend/internal/service"
)

type MenuHandler struct {
	menuService *service.MenuService
}

func NewMenuHandler() *MenuHandler {
	return &MenuHandler{
		menuService: service.NewMenuService(),
	}
}

// GetMenuTree 获取菜单树
func (h *MenuHandler) GetMenuTree(c *gin.Context) {
	// 获取当前管理员ID
	adminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	menuTree, err := h.menuService.GetMenuTree(adminID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menuTree)
}

// GetMenus 获取菜单列表（扁平）
func (h *MenuHandler) GetMenus(c *gin.Context) {
	menus, err := h.menuService.GetAllMenus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menus)
}

// CreateMenu 创建菜单
func (h *MenuHandler) CreateMenu(c *gin.Context) {
	var req struct {
		Name      string `json:"name" binding:"required"`
		ParentID  uint   `json:"parent_id"`
		URL       string `json:"url"`
		Icon      string `json:"icon"`
		SortOrder int    `json:"sort_order"`
		IsVisible bool   `json:"is_visible"`
		IsSystem  bool   `json:"is_system"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	menu, err := h.menuService.CreateMenu(&service.CreateMenuRequest{
		Name:      req.Name,
		ParentID:  req.ParentID,
		URL:       req.URL,
		Icon:      req.Icon,
		SortOrder: req.SortOrder,
		IsVisible: req.IsVisible,
		IsSystem:  req.IsSystem,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, menu)
}

// UpdateMenu 更新菜单
func (h *MenuHandler) UpdateMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var req struct {
		Name      string `json:"name"`
		ParentID  uint   `json:"parent_id"`
		URL       string `json:"url"`
		Icon      string `json:"icon"`
		SortOrder int    `json:"sort_order"`
		IsVisible *bool  `json:"is_visible"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	menu, err := h.menuService.UpdateMenu(uint(id), &service.UpdateMenuRequest{
		Name:      req.Name,
		ParentID:  req.ParentID,
		URL:       req.URL,
		Icon:      req.Icon,
		SortOrder: req.SortOrder,
		IsVisible: req.IsVisible,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menu)
}

// DeleteMenu 删除菜单
func (h *MenuHandler) DeleteMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	if err := h.menuService.DeleteMenu(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// RegisterRoutes 注册路由
func (h *MenuHandler) RegisterRoutes(r *gin.RouterGroup) {
	menu := r.Group("/menus")
	{
		menu.GET("/tree", h.GetMenuTree)
		menu.GET("", h.GetMenus)
		menu.POST("", h.CreateMenu)
		menu.PUT("/:id", h.UpdateMenu)
		menu.DELETE("/:id", h.DeleteMenu)
	}
}
