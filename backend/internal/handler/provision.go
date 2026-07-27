package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type ProvisionHandler struct {
	provSvc *service.ProvisionService
	log     *logger.Logger
}

func NewProvisionHandler(provSvc *service.ProvisionService, log *logger.Logger) *ProvisionHandler {
	return &ProvisionHandler{provSvc: provSvc, log: log}
}

// GetList returns all provision modules.
// GET /provisions
func (h *ProvisionHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	moduleType := c.Query("type")

	var active *bool
	if v := c.Query("active"); v != "" {
		b := v == "true" || v == "1"
		active = &b
	}

	modules, total, err := h.provSvc.GetList(page, pageSize, moduleType, active)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, modules, total, page, pageSize)
}

// GetDetail returns a single provision module by ID.
// GET /provisions/:id
func (h *ProvisionHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid module id")
		return
	}

	module, err := h.provSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "module not found")
		return
	}
	response.Success(c, module)
}

// Create creates a new provision module.
// POST /provisions
func (h *ProvisionHandler) Create(c *gin.Context) {
	var req service.CreateProvisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	module, err := h.provSvc.Create(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, module)
}

// Update updates a provision module.
// PUT /provisions/:id
func (h *ProvisionHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid module id")
		return
	}

	var req service.UpdateProvisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	module, err := h.provSvc.Update(uint(id), req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, module)
}

// Delete deletes a provision module.
// DELETE /provisions/:id
func (h *ProvisionHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid module id")
		return
	}

	if err := h.provSvc.Delete(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "module deleted")
}

// TestConnection tests connectivity to a provision module's server.
// POST /provisions/:id/test
func (h *ProvisionHandler) TestConnection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid module id")
		return
	}

	module, err := h.provSvc.TestConnection(uint(id))
	if err != nil {
		response.Success(c, gin.H{
			"connected": false,
			"error":     err.Error(),
			"module":    module,
		})
		return
	}
	response.Success(c, gin.H{
		"connected": true,
		"module":    module,
	})
}

// GetLogs returns provision module operation logs.
// GET /provisions/logs
func (h *ProvisionHandler) GetLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	action := c.Query("action")

	var moduleID *uint
	if mid := c.Query("module_id"); mid != "" {
		v, _ := strconv.ParseUint(mid, 10, 64)
		id := uint(v)
		moduleID = &id
	}

	logs, total, err := h.provSvc.GetLogs(page, pageSize, moduleID, action)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, logs, total, page, pageSize)
}
