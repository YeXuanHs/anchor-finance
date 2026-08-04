package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// DeveloperHandler 开发者管理处理器
type DeveloperHandler struct {
	svc *service.DeveloperService
	log *logger.Logger
}

// NewDeveloperHandler creates a new DeveloperHandler.
func NewDeveloperHandler(svc *service.DeveloperService, log *logger.Logger) *DeveloperHandler {
	return &DeveloperHandler{svc: svc, log: log}
}

// List returns paginated developers.
// GET /developers
func (h *DeveloperHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	items, total, err := h.svc.List(page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, items, total, page, pageSize)
}

// GetDetail returns developer detail.
// GET /developers/:id
func (h *DeveloperHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	detail, err := h.svc.GetDetail(uint(id))
	if err != nil {
		response.NotFound(c, "developer not found")
		return
	}
	response.Success(c, detail)
}

// Create creates a new developer.
// POST /developers
func (h *DeveloperHandler) Create(c *gin.Context) {
	var dev service.Developer
	if err := c.ShouldBindJSON(&dev); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.Create(&dev); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, dev)
}

// Update updates a developer.
// PUT /developers/:id
func (h *DeveloperHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.Update(uint(id), updates); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "updated")
}

// Delete deletes a developer.
// DELETE /developers/:id
func (h *DeveloperHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "deleted")
}

// Approve approves a developer.
// POST /developers/:id/approve
func (h *DeveloperHandler) Approve(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.svc.Approve(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "approved")
}

// Reject rejects a developer.
// POST /developers/:id/reject
func (h *DeveloperHandler) Reject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)

	if err := h.svc.Reject(uint(id), req.Reason); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "rejected")
}

// GetAPIKeys returns API keys for a developer.
// GET /developers/:id/api-keys
func (h *DeveloperHandler) GetAPIKeys(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	items, err := h.svc.GetAPIKeys(uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, items)
}

// GenerateAPIKey generates a new API key.
// POST /developers/:id/api-keys
func (h *DeveloperHandler) GenerateAPIKey(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req struct {
		Name   string `json:"name" binding:"required"`
		Scopes string `json:"scopes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	apiKey, err := h.svc.GenerateAPIKey(uint(id), req.Name, req.Scopes)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, apiKey)
}

// RevokeAPIKey revokes an API key.
// DELETE /developers/:id/api-keys/:key_id
func (h *DeveloperHandler) RevokeAPIKey(c *gin.Context) {
	devID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid developer id")
		return
	}
	keyID, err := strconv.ParseUint(c.Param("key_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid key id")
		return
	}

	if err := h.svc.RevokeAPIKey(uint(devID), uint(keyID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "revoked")
}

// GetDocs returns developer documentation.
// GET /developers/docs
func (h *DeveloperHandler) GetDocs(c *gin.Context) {
	items, err := h.svc.GetDocs()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, items)
}

// UpdateDocs updates developer documentation.
// PUT /developers/docs
func (h *DeveloperHandler) UpdateDocs(c *gin.Context) {
	var docs []service.DeveloperDoc
	if err := c.ShouldBindJSON(&docs); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.UpdateDocs(docs); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "updated")
}

// GetBilling returns billing records.
// GET /developers/billing
func (h *DeveloperHandler) GetBilling(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	items, total, err := h.svc.GetBilling(page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, items, total, page, pageSize)
}

// SettleBilling settles pending billing.
// POST /developers/billing/settle
func (h *DeveloperHandler) SettleBilling(c *gin.Context) {
	var req struct {
		DeveloperID uint   `json:"developer_id" binding:"required"`
		Period      string `json:"period"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.SettleBilling(req.DeveloperID, req.Period); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "settled")
}
