package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type SystemHandler struct {
	svc *service.SystemService
	log *logger.Logger
}

func NewSystemHandler(svc *service.SystemService, log *logger.Logger) *SystemHandler {
	return &SystemHandler{svc: svc, log: log}
}

// GetCommonInfo returns common system information.
func (h *SystemHandler) GetCommonInfo(c *gin.Context) {
	info, err := h.svc.GetCommonInfo()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, info)
}

// GetUpdateContent returns update content.
func (h *SystemHandler) GetUpdateContent(c *gin.Context) {
	update, err := h.svc.GetUpdateContent()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	if update == nil {
		response.Success(c, gin.H{"message": "no updates available"})
		return
	}
	response.Success(c, update)
}

// CheckUpdate checks for system updates.
func (h *SystemHandler) CheckUpdate(c *gin.Context) {
	update, err := h.svc.CheckUpdate()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	if update == nil {
		response.Success(c, gin.H{"has_update": false})
		return
	}
	response.Success(c, gin.H{"has_update": true, "update": update})
}

// GetUpdateList returns paginated system updates.
func (h *SystemHandler) GetUpdateList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	updates, total, err := h.svc.GetUpdateList(page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, updates, total, page, pageSize)
}

// InstallUpdate installs a system update.
func (h *SystemHandler) InstallUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid update id")
		return
	}

	if err := h.svc.InstallUpdate(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "update installed")
}

// GetSystemLog returns system logs.
func (h *SystemHandler) GetSystemLog(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	level := c.Query("level")
	module := c.Query("module")

	logs, total, err := h.svc.GetSystemLog(page, pageSize, level, module)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, logs, total, page, pageSize)
}

// ClearCache clears system cache.
func (h *SystemHandler) ClearCache(c *gin.Context) {
	if err := h.svc.ClearCache(); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "cache cleared")
}
