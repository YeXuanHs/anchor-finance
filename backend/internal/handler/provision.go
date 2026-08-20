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

// List is an alias for GetList.
func (h *ProvisionHandler) List(c *gin.Context) {
	h.GetList(c)
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

// Execute manually triggers provisioning for an order or product.
// POST /provision/execute
func (h *ProvisionHandler) Execute(c *gin.Context) {
	var req struct {
		OrderID       *uint  `json:"order_id"`
		UserProductID *uint  `json:"user_product_id"`
		Action        string `json:"action" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	validActions := map[string]bool{
		"create": true, "suspend": true, "terminate": true,
		"unsuspend": true, "rebuild": true, "renew": true,
	}
	if !validActions[req.Action] {
		response.BadRequest(c, "invalid action")
		return
	}

	var err error
	if req.OrderID != nil {
		err = h.provSvc.ProvisionOrder(*req.OrderID)
	} else if req.UserProductID != nil {
		err = h.provSvc.ProvisionProduct(*req.UserProductID, req.Action)
	} else {
		response.BadRequest(c, "order_id or user_product_id required")
		return
	}

	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "provisioning triggered")
}

// Suspend suspends a service.
// POST /provision/:id/suspend
func (h *ProvisionHandler) Suspend(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)
	if req.Reason == "" {
		req.Reason = "Admin suspended"
	}

	if err := h.provSvc.SuspendService(uint(id), req.Reason); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "service suspended")
}

// Terminate terminates a service.
// POST /provision/:id/terminate
func (h *ProvisionHandler) Terminate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)
	if req.Reason == "" {
		req.Reason = "Admin terminated"
	}

	if err := h.provSvc.TerminateService(uint(id), req.Reason); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "service terminated")
}

// Unsuspend reactivates a suspended service.
// POST /provision/:id/unsuspend
func (h *ProvisionHandler) Unsuspend(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.provSvc.UnsuspendService(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "service unsuspended")
}
