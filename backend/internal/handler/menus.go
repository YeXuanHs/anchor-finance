package handler

import (
	"encoding/json"
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type MenusHandler struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewMenusHandler(db *gorm.DB, log *logger.Logger) *MenusHandler {
	return &MenusHandler{db: db, log: log}
}

// GetMenu returns a menu by ID.
func (h *MenusHandler) GetMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid menu id")
		return
	}

	var menu model.Menu
	if err := h.db.First(&menu, id).Error; err != nil {
		response.NotFound(c, "menu not found")
		return
	}
	response.Success(c, menu)
}

// GetMenuList returns a list of menus.
func (h *MenusHandler) GetMenuList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	menuType := c.Query("type")

	var menus []model.Menu
	var total int64

	query := h.db.Model(&model.Menu{})
	if menuType != "" {
		query = query.Where("type = ?", menuType)
	}
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("sort_order ASC, id ASC").Find(&menus).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, menus, total, page, pageSize)
}

// SetNavList sets the navigation list for a menu.
func (h *MenusHandler) SetNavList(c *gin.Context) {
	var req struct {
		MenuID  uint        `json:"menu_id" binding:"required"`
		NavList interface{} `json:"nav_list" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var menu model.Menu
	if err := h.db.First(&menu, req.MenuID).Error; err != nil {
		response.NotFound(c, "menu not found")
		return
	}

	extraJSON, _ := json.Marshal(gin.H{"nav_list": req.NavList})
	if err := h.db.Model(&menu).Update("extra", datatypes.JSON(extraJSON)).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "nav list set",
	})
}

// CreateMenu creates a new menu.
func (h *MenusHandler) CreateMenu(c *gin.Context) {
	var req struct {
		Name    string `json:"name" binding:"required"`
		Type    int    `json:"type"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	menu := model.Menu{
		Name: req.Name,
		Type: strconv.Itoa(req.Type),
	}
	if req.Content != "" {
		menu.Extra = datatypes.JSON([]byte(`{"content":"` + req.Content + `"}`))
	}

	if err := h.db.Create(&menu).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, menu)
}

// UpdateMenu updates an existing menu.
func (h *MenusHandler) UpdateMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid menu id")
		return
	}

	var req struct {
		Name    string `json:"name"`
		Type    int    `json:"type"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var menu model.Menu
	if err := h.db.First(&menu, id).Error; err != nil {
		response.NotFound(c, "menu not found")
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Type != 0 {
		updates["type"] = strconv.Itoa(req.Type)
	}
	if req.Content != "" {
		updates["extra"] = datatypes.JSON([]byte(`{"content":"` + req.Content + `"}`))
	}

	if len(updates) > 0 {
		if err := h.db.Model(&menu).Updates(updates).Error; err != nil {
			response.ServerError(c, err.Error())
			return
		}
	}

	response.Success(c, menu)
}

// DeleteMenu deletes a menu.
func (h *MenusHandler) DeleteMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid menu id")
		return
	}

	if err := h.db.Delete(&model.Menu{}, id).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "menu deleted")
}

// SetActive sets a menu as active.
func (h *MenusHandler) SetActive(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid menu id")
		return
	}

	var menu model.Menu
	if err := h.db.First(&menu, id).Error; err != nil {
		response.NotFound(c, "menu not found")
		return
	}

	if err := h.db.Model(&menu).Update("is_active", true).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"id":         id,
		"is_active":  true,
		"message":    "menu set as active",
	})
}
