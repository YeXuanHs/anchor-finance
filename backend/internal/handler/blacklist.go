package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// BlacklistHandler handles blacklist requests.
type BlacklistHandler struct {
	blacklistSvc *service.BlacklistService
	log          *logger.Logger
}

// NewBlacklistHandler creates a new BlacklistHandler.
func NewBlacklistHandler(blacklistSvc *service.BlacklistService, log *logger.Logger) *BlacklistHandler {
	return &BlacklistHandler{blacklistSvc: blacklistSvc, log: log}
}

// List returns all blacklist entries.
// GET /admin/blacklist?page=1&page_size=20&type=ip
func (h *BlacklistHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	blacklistType := c.Query("type")

	items, total, err := h.blacklistSvc.List(page, pageSize, blacklistType)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      items,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// Create adds a new blacklist entry.
// POST /admin/blacklist
func (h *BlacklistHandler) Create(c *gin.Context) {
	var req struct {
		Type   string `json:"type" binding:"required"`
		Value  string `json:"value" binding:"required"`
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	adminID, _ := c.Get("admin_id")

	item := &model.Blacklist{
		Type:    req.Type,
		Value:   req.Value,
		Reason:  req.Reason,
		AdminID: adminID.(uint),
	}

	if err := h.blacklistSvc.Create(item); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, item)
}

// Update updates a blacklist entry.
// PUT /admin/blacklist/:id
func (h *BlacklistHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := map[string]interface{}{
		"reason": req.Reason,
	}

	if err := h.blacklistSvc.Update(uint(id), updates); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessMsg(c, "updated")
}

// Delete removes a blacklist entry.
// DELETE /admin/blacklist/:id
func (h *BlacklistHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.blacklistSvc.Delete(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessMsg(c, "deleted")
}
