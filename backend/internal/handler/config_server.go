package handler

import (
	"strconv"
	"strings"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type ConfigServerHandler struct {
	svc *service.ConfigServerService
	log *logger.Logger
}

func NewConfigServerHandler(svc *service.ConfigServerService, log *logger.Logger) *ConfigServerHandler {
	return &ConfigServerHandler{svc: svc, log: log}
}

// ---------- Server Config ----------

// GetList returns paginated server configs.
func (h *ConfigServerHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	serverType := c.Query("type")
	keyword := c.Query("keyword")

	items, total, err := h.svc.GetList(page, pageSize, serverType, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, items, total, page, pageSize)
}

// GetDetail returns a single server config by ID.
func (h *ConfigServerHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server config id")
		return
	}

	item, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "server config not found")
		return
	}
	response.Success(c, item)
}

// Create creates a server config.
func (h *ConfigServerHandler) Create(c *gin.Context) {
	var req service.CreateServerConfigRequest
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

// Update updates a server config.
func (h *ConfigServerHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server config id")
		return
	}

	var req service.UpdateServerConfigRequest
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

// Delete deletes a server config.
func (h *ConfigServerHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server config id")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "server config deleted")
}

// BatchUpdateStatus batch-updates status for server configs.
func (h *ConfigServerHandler) BatchUpdateStatus(c *gin.Context) {
	var req struct {
		IDs    []uint `json:"ids" binding:"required,min=1"`
		Status int16  `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.BatchUpdateStatus(req.IDs, req.Status); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "server configs status updated")
}

// BatchDelete batch-deletes server configs.
func (h *ConfigServerHandler) BatchDelete(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.BatchDelete(req.IDs); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "server configs deleted")
}

// UpdateSort updates sort order for a server config.
func (h *ConfigServerHandler) UpdateSort(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server config id")
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

// ---------- Server Template ----------

// GetTemplateList returns paginated server templates.
func (h *ConfigServerHandler) GetTemplateList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	templateType := c.Query("type")

	items, total, err := h.svc.GetTemplateList(page, pageSize, templateType)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, items, total, page, pageSize)
}

// CreateTemplate creates a server template.
func (h *ConfigServerHandler) CreateTemplate(c *gin.Context) {
	var req service.CreateServerTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	item, err := h.svc.CreateTemplate(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, item)
}

// UpdateTemplate updates a server template.
func (h *ConfigServerHandler) UpdateTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid template id")
		return
	}

	var req service.UpdateServerTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	item, err := h.svc.UpdateTemplate(uint(id), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, item)
}

// DeleteTemplate deletes a server template.
func (h *ConfigServerHandler) DeleteTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid template id")
		return
	}

	if err := h.svc.DeleteTemplate(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "template deleted")
}

// helper: convert comma-separated string to uint slice
func parseIDs(s string) []uint {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	ids := make([]uint, 0, len(parts))
	for _, p := range parts {
		if v, err := strconv.ParseUint(strings.TrimSpace(p), 10, 64); err == nil {
			ids = append(ids, uint(v))
		}
	}
	return ids
}
