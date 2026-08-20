package handler

import (
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

// List returns all bindings for the authenticated user.
func (h *BindHandler) List(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	items, err := h.svc.ListBindings(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, items)
}

// GetByProvider returns a single binding by provider for the authenticated user.
func (h *BindHandler) GetByProvider(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}
	provider := c.Param("provider")
	if provider == "" {
		response.BadRequest(c, "provider is required")
		return
	}

	item, err := h.svc.GetBinding(userID, provider)
	if err != nil {
		response.NotFound(c, "binding not found")
		return
	}
	response.Success(c, item)
}

// Unbind unbinds a provider for the authenticated user.
func (h *BindHandler) Unbind(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	var req struct {
		Provider string `json:"provider" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.UnbindOAuth(userID, req.Provider); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "binding removed")
}

// GetDetail returns details of a specific binding.
func (h *BindHandler) GetDetail(c *gin.Context) {
	id := c.Param("id")
	binding, err := h.svc.GetDetail(id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, binding)
}

// GetUserBindings returns all bindings for a specific user (admin).
func (h *BindHandler) GetUserBindings(c *gin.Context) {
	userID := c.Param("user_id")
	bindings, err := h.svc.GetUserBindings(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, bindings)
}

// Delete deletes a binding (admin).
func (h *BindHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(id); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "binding deleted")
}

// GetProviders returns available OAuth providers.
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
		IDs []string `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.BatchUnbind(req.IDs); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "bindings removed")
}
