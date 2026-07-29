package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// BindHandler handles account binding requests.
type BindHandler struct {
	svc *service.BindService
	log *logger.Logger
}

// NewBindHandler creates a new BindHandler.
func NewBindHandler(svc *service.BindService, log *logger.Logger) *BindHandler {
	return &BindHandler{svc: svc, log: log}
}

// List returns paginated account bindings.
func (h *BindHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	provider := c.Query("provider")

	var userID *uint
	if uid := c.Query("user_id"); uid != "" {
		v, _ := strconv.ParseUint(uid, 10, 64)
		id := uint(v)
		userID = &id
	}

	items, total, err := h.svc.List(page, pageSize, userID, provider)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, items, total, page, pageSize)
}

// GetDetail returns a single binding by ID.
func (h *BindHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid bind id")
		return
	}

	item, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "binding not found")
		return
	}
	response.Success(c, item)
}

// GetUserBindings returns all bindings for a user.
func (h *BindHandler) GetUserBindings(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	bindings, err := h.svc.GetByUserID(uint(userID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, bindings)
}

// Delete deletes a binding.
func (h *BindHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid bind id")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "binding deleted")
}

// Unbind unbinds a provider for a user.
func (h *BindHandler) Unbind(c *gin.Context) {
	var req struct {
		UserID   uint   `json:"user_id" binding:"required"`
		Provider string `json:"provider" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.Unbind(req.UserID, req.Provider); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "binding removed")
}

// GetProviders returns available binding providers.
func (h *BindHandler) GetProviders(c *gin.Context) {
	providers, err := h.svc.GetProviders()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, providers)
}

// GetStats returns binding statistics.
func (h *BindHandler) GetStats(c *gin.Context) {
	stats, err := h.svc.GetStats()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, stats)
}

// BatchUnbind batch unbinds multiple bindings.
func (h *BindHandler) BatchUnbind(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	count, err := h.svc.BatchDelete(req.IDs)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{
		"deleted_count": count,
	})
}
