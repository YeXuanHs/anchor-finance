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
