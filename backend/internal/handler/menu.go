package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"anchorfinance/pkg/db"

	"github.com/gin-gonic/gin"
)

// MenuHandler handles menu-related requests.
type MenuHandler struct{}

// NewMenuHandler creates a new MenuHandler.
func NewMenuHandler() *MenuHandler {
	return &MenuHandler{}
}

// MenuItem represents a menu item from the database.
type MenuItem struct {
	ID           uint   `json:"id"`
	ParentID     uint   `json:"parent_id"`
	Name         string `json:"name"`
	URL          string `json:"url"`
	Icon         string `json:"icon"`
	SortOrder    int    `json:"sort_order"`
	IsVisible    bool   `json:"is_visible"`
	LanguageMap  string `json:"language_map"` // JSON: {"zh-CN":"客户","zh-TW":"客戶","en":"Customer"}
}

// MenuTreeNode represents a menu tree node with children.
type MenuTreeNode struct {
	ID       uint            `json:"id"`
	Path     string          `json:"path"`
	Name     string          `json:"name"`
	Meta     *MenuMeta       `json:"meta,omitempty"`
	Children []*MenuTreeNode `json:"children,omitempty"`
}

// MenuMeta contains menu metadata for the sidebar.
type MenuMeta struct {
	Title        string            `json:"title"`
	Icon         string            `json:"icon"`
	IsHide       bool              `json:"isHide,omitempty"`
	LanguageMap  map[string]string `json:"languageMap,omitempty"`
}

// List returns a list of menus.
func (h *MenuHandler) List(c *gin.Context) {
	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
		return
	}

	var menus []MenuItem
	if err := database.Table("menus").
		Select("id, parent_id, name, url, icon, sort_order, is_visible, language_map").
		Where("deleted_at IS NULL").
		Order("sort_order ASC, id ASC").
		Find(&menus).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": menus})
}

// GetTree returns the menu tree.
func (h *MenuHandler) GetTree(c *gin.Context) {
	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
		return
	}

	var menus []MenuItem
	if err := database.Table("menus").
		Select("id, parent_id, name, url, icon, sort_order, is_visible, language_map").
		Where("deleted_at IS NULL").
		Order("sort_order ASC, id ASC").
		Find(&menus).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
		return
	}

	// Build tree structure with sidebar-expected format
	menuMap := make(map[uint]*MenuTreeNode)
	for _, m := range menus {
		// Parse language_map JSON
		var langMap map[string]string
		if m.LanguageMap != "" {
			json.Unmarshal([]byte(m.LanguageMap), &langMap)
		}
		
		menuMap[m.ID] = &MenuTreeNode{
			ID:   m.ID,
			Path: m.URL,
			Name: m.Name,
			Meta: &MenuMeta{
				Title:       m.Name,
				Icon:        m.Icon,
				IsHide:      !m.IsVisible,
				LanguageMap: langMap,
			},
			Children: []*MenuTreeNode{},
		}
	}

	var roots []*MenuTreeNode
	for _, m := range menus {
		node := menuMap[m.ID]
		if m.ParentID == 0 {
			roots = append(roots, node)
		} else if parent, ok := menuMap[m.ParentID]; ok {
			parent.Children = append(parent.Children, node)
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": roots})
}

// Create creates a new menu.
func (h *MenuHandler) Create(c *gin.Context) {
	var req struct {
		Name      string `json:"name" binding:"required"`
		ParentID  uint   `json:"parent_id"`
		URL       string `json:"url"`
		Icon      string `json:"icon"`
		SortOrder int    `json:"sort_order"`
		IsVisible bool   `json:"is_visible"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "数据库未初始化"})
		return
	}

	menu := map[string]interface{}{
		"parent_id":  req.ParentID,
		"name":       req.Name,
		"url":        req.URL,
		"icon":       req.Icon,
		"sort_order": req.SortOrder,
		"is_visible": req.IsVisible,
	}

	if err := database.Table("menus").Create(menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "创建成功", "data": menu})
}

// Update updates a menu.
func (h *MenuHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
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
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "数据库未初始化"})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.URL != "" {
		updates["url"] = req.URL
	}
	if req.Icon != "" {
		updates["icon"] = req.Icon
	}
	updates["parent_id"] = req.ParentID
	updates["sort_order"] = req.SortOrder
	if req.IsVisible != nil {
		updates["is_visible"] = *req.IsVisible
	}

	if err := database.Table("menus").Where("id = ?", id).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功"})
}

// Delete deletes a menu.
func (h *MenuHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "数据库未初始化"})
		return
	}

	if err := database.Table("menus").Where("id = ?", id).Delete(nil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// Sort sorts menus.
func (h *MenuHandler) Sort(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "数据库未初始化"})
		return
	}

	for i, id := range req.IDs {
		database.Table("menus").Where("id = ?", id).Update("sort_order", i)
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "排序成功"})
}

// RegisterRoutes registers menu routes.
func (h *MenuHandler) RegisterRoutes(r *gin.RouterGroup) {
	menu := r.Group("/menus")
	{
		menu.GET("/tree", h.GetTree)
		menu.GET("", h.List)
		menu.POST("", h.Create)
		menu.PUT("/:id", h.Update)
		menu.DELETE("/:id", h.Delete)
		menu.POST("/sort", h.Sort)
	}
}
