package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type ConfigOptionHandler struct {
	svc *service.ConfigOptionService
	log *logger.Logger
}

func NewConfigOptionHandler(svc *service.ConfigOptionService, log *logger.Logger) *ConfigOptionHandler {
	return &ConfigOptionHandler{svc: svc, log: log}
}

// GetList returns paginated config options.
func (h *ConfigOptionHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	group := c.Query("group")
	keyword := c.Query("keyword")

	items, total, err := h.svc.GetList(page, pageSize, group, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, items, total, page, pageSize)
}

// GetByGroup returns config options filtered by group.
func (h *ConfigOptionHandler) GetByGroup(c *gin.Context) {
	group := c.Param("group")
	if group == "" {
		response.BadRequest(c, "group is required")
		return
	}

	items, err := h.svc.GetByGroup(group)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, items)
}

// GetGroups returns all distinct option groups.
func (h *ConfigOptionHandler) GetGroups(c *gin.Context) {
	groups, err := h.svc.GetGroups()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, groups)
}

// GetDetail returns a single config option by ID.
func (h *ConfigOptionHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid option id")
		return
	}

	item, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "config option not found")
		return
	}
	response.Success(c, item)
}

// Create creates a config option.
func (h *ConfigOptionHandler) Create(c *gin.Context) {
	var req service.CreateConfigOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	item, err := h.svc.Create(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, item)
}

// Update updates a config option.
func (h *ConfigOptionHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid option id")
		return
	}

	var req service.UpdateConfigOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	item, err := h.svc.Update(uint(id), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, item)
}

// Delete deletes a config option.
func (h *ConfigOptionHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid option id")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "config option deleted")
}

// BatchUpdateSort batch-updates sort order for config options.
func (h *ConfigOptionHandler) BatchUpdateSort(c *gin.Context) {
	var req []service.SortItem
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.BatchUpdateSort(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "sort order updated")
}

// UpdateSort updates sort order for a single config option.
func (h *ConfigOptionHandler) UpdateSort(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid option id")
		return
	}

	var req struct {
		SortOrder int `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.UpdateSort(uint(id), req.SortOrder); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "sort order updated")
}

// UpdateValue updates the value of a config option by code.
func (h *ConfigOptionHandler) UpdateValue(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		response.BadRequest(c, "option code is required")
		return
	}

	var req struct {
		Value string `json:"value"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.UpdateValue(code, req.Value); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "option value updated")
}

// BatchUpdateValue batch-updates option values by code.
func (h *ConfigOptionHandler) BatchUpdateValue(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.BatchUpdateValue(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "option values updated")
}
