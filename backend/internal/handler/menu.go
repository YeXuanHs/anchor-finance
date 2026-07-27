package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	menuSvc *service.MenuService
	log     *logger.Logger
}

func NewMenuHandler(menuSvc *service.MenuService, log *logger.Logger) *MenuHandler {
	return &MenuHandler{menuSvc: menuSvc, log: log}
}

// List returns all menus as a flat list.
func (h *MenuHandler) List(c *gin.Context) {
	menuType := c.Query("type")

	menus, err := h.menuSvc.List(menuType)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, menus)
}

// GetTree returns menus in tree structure.
func (h *MenuHandler) GetTree(c *gin.Context) {
	menuType := c.Query("type")

	tree, err := h.menuSvc.GetTree(menuType)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, tree)
}

// GetVisibleTree returns only visible menus in tree structure (public).
func (h *MenuHandler) GetVisibleTree(c *gin.Context) {
	menuType := c.Query("type")

	tree, err := h.menuSvc.GetVisibleTree(menuType)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, tree)
}

// GetDetail returns a single menu by ID.
func (h *MenuHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid menu id")
		return
	}

	menu, err := h.menuSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "menu not found")
		return
	}
	response.Success(c, menu)
}

// Create creates a new menu item.
func (h *MenuHandler) Create(c *gin.Context) {
	var req service.CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	menu, err := h.menuSvc.Create(req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, menu)
}

// Update updates a menu item.
func (h *MenuHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid menu id")
		return
	}

	var req service.UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	menu, err := h.menuSvc.Update(uint(id), req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, menu)
}

// Delete deletes a menu and its children.
func (h *MenuHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid menu id")
		return
	}

	if err := h.menuSvc.Delete(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "menu deleted")
}

// Sort updates sort order for multiple menus.
func (h *MenuHandler) Sort(c *gin.Context) {
	var req service.SortMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.menuSvc.Sort(req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "menu sorted")
}
