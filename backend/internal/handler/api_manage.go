package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// APIManageHandler handles API management requests.
type APIManageHandler struct {
	svc *service.APIManageService
	log *logger.Logger
}

// NewAPIManageHandler creates a new APIManageHandler.
func NewAPIManageHandler(svc *service.APIManageService, log *logger.Logger) *APIManageHandler {
	return &APIManageHandler{svc: svc, log: log}
}

// List returns paginated API keys.
func (h *APIManageHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	var status *int16
	if s := c.Query("status"); s != "" {
		v, _ := strconv.Atoi(s)
		sv := int16(v)
		status = &sv
	}

	items, total, err := h.svc.List(page, pageSize, keyword, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, items, total, page, pageSize)
}

// GetDetail returns a single API key by ID.
func (h *APIManageHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid api key id")
		return
	}

	item, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "api key not found")
		return
	}
	response.Success(c, item)
}

// Create creates a new API key.
func (h *APIManageHandler) Create(c *gin.Context) {
	var req struct {
		Name        string   `json:"name" binding:"required"`
		Description string   `json:"description"`
		Permissions []string `json:"permissions"`
		RateLimit   int      `json:"rate_limit"`
		ExpiresAt   string   `json:"expires_at"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	apiKey := &model.APIKey{
		Name:        req.Name,
		Description: req.Description,
		Status:      1,
		RateLimit:   req.RateLimit,
	}

	if err := h.svc.Create(apiKey); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, apiKey)
}

// Update updates an API key.
func (h *APIManageHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid api key id")
		return
	}

	var req struct {
		Name        *string  `json:"name"`
		Description *string  `json:"description"`
		Permissions []string `json:"permissions"`
		RateLimit   *int     `json:"rate_limit"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.RateLimit != nil {
		updates["rate_limit"] = *req.RateLimit
	}

	if err := h.svc.Update(uint(id), updates); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "api key updated")
}

// Delete deletes an API key.
func (h *APIManageHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid api key id")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "api key deleted")
}

// SetStatus sets the status of an API key.
func (h *APIManageHandler) SetStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid api key id")
		return
	}

	var req struct {
		Status int16 `json:"status" binding:"required,oneof=0 1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.SetStatus(uint(id), req.Status); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "api key status updated")
}

// Regenerate regenerates the API key secret.
func (h *APIManageHandler) Regenerate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid api key id")
		return
	}

	newKey, err := h.svc.Regenerate(uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{
		"api_key": newKey,
	})
}

// GetPermissions returns available API permissions.
func (h *APIManageHandler) GetPermissions(c *gin.Context) {
	permissions, err := h.svc.GetPermissions()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, permissions)
}
