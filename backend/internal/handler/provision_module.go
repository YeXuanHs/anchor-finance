package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type ProvisionModuleHandler struct {
	provSvc *service.ProvisionService
	log     *logger.Logger
}

func NewProvisionModuleHandler(provSvc *service.ProvisionService, log *logger.Logger) *ProvisionModuleHandler {
	return &ProvisionModuleHandler{provSvc: provSvc, log: log}
}

// ==================== Module Management ====================

// GetModules returns all provision modules.
// GET /provision-modules
func (h *ProvisionModuleHandler) GetModules(c *gin.Context) {
	moduleType := c.Query("type")
	var enabled *bool
	if v := c.Query("enabled"); v != "" {
		b := v == "true" || v == "1"
		enabled = &b
	}

	modules, err := h.provSvc.GetModules(moduleType, enabled)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, modules)
}

// GetModule returns a single module by ID.
// GET /provision-modules/:id
func (h *ProvisionModuleHandler) GetModule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid module id")
		return
	}

	module, err := h.provSvc.GetModuleByID(uint(id))
	if err != nil {
		response.NotFound(c, "module not found")
		return
	}
	response.Success(c, module)
}

// CreateModule creates a new provision module.
// POST /provision-modules
func (h *ProvisionModuleHandler) CreateModule(c *gin.Context) {
	var req service.CreateModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	module, err := h.provSvc.CreateModule(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, module)
}

// UpdateModule updates a provision module.
// PUT /provision-modules/:id
func (h *ProvisionModuleHandler) UpdateModule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid module id")
		return
	}

	var req service.UpdateModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	module, err := h.provSvc.UpdateModule(uint(id), req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, module)
}

// DeleteModule deletes a provision module.
// DELETE /provision-modules/:id
func (h *ProvisionModuleHandler) DeleteModule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid module id")
		return
	}

	if err := h.provSvc.DeleteModule(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "module deleted")
}

// TestModule tests connectivity to a module.
// POST /provision-modules/:id/test
func (h *ProvisionModuleHandler) TestModule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid module id")
		return
	}

	module, err := h.provSvc.TestModule(uint(id))
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

// ==================== Client Area ====================

// RenderClientArea renders the client area for a host.
// GET /provision-modules/client-area/:host_id
func (h *ProvisionModuleHandler) RenderClientArea(c *gin.Context) {
	hostID, err := strconv.ParseUint(c.Param("host_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	result, err := h.provSvc.RenderClientArea(uint(hostID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, result)
}

// RenderClientAreaDetail renders the detailed client area.
// GET /provision-modules/client-area/:host_id/detail
func (h *ProvisionModuleHandler) RenderClientAreaDetail(c *gin.Context) {
	hostID, err := strconv.ParseUint(c.Param("host_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	result, err := h.provSvc.RenderClientAreaDetail(uint(hostID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, result)
}

// GetClientButtons returns available client buttons for a host.
// GET /provision-modules/client-area/:host_id/buttons
func (h *ProvisionModuleHandler) GetClientButtons(c *gin.Context) {
	hostID, err := strconv.ParseUint(c.Param("host_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	buttons, err := h.provSvc.GetClientButtons(uint(hostID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, buttons)
}

// ExecuteClientButton executes a client button action.
// POST /provision-modules/client-area/:host_id/execute
func (h *ProvisionModuleHandler) ExecuteClientButton(c *gin.Context) {
	hostID, err := strconv.ParseUint(c.Param("host_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	var req struct {
		Action string `json:"action" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.provSvc.ExecuteClientButton(uint(hostID), req.Action)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// ==================== Admin Area ====================

// RenderAdminArea renders the admin area for a host.
// GET /provision-modules/admin-area/:host_id
func (h *ProvisionModuleHandler) RenderAdminArea(c *gin.Context) {
	hostID, err := strconv.ParseUint(c.Param("host_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	result, err := h.provSvc.RenderAdminArea(uint(hostID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, result)
}

// GetAdminButtons returns admin buttons for a host.
// GET /provision-modules/admin-area/:host_id/buttons
func (h *ProvisionModuleHandler) GetAdminButtons(c *gin.Context) {
	hostID, err := strconv.ParseUint(c.Param("host_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	buttons, err := h.provSvc.GetAdminButtons(uint(hostID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, buttons)
}

// ExecuteAdminButton executes an admin button action.
// POST /provision-modules/admin-area/:host_id/execute
func (h *ProvisionModuleHandler) ExecuteAdminButton(c *gin.Context) {
	hostID, err := strconv.ParseUint(c.Param("host_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	var req struct {
		Action string `json:"action" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.provSvc.ExecuteAdminButton(uint(hostID), req.Action)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// GetDefaultButtons returns the default provision buttons.
// GET /provision-modules/default-buttons
func (h *ProvisionModuleHandler) GetDefaultButtons(c *gin.Context) {
	buttons := h.provSvc.GetDefaultButtons()
	response.Success(c, buttons)
}

// ==================== Charts ====================

// GetCharts returns charts for a host.
// GET /provision-modules/charts/:host_id
func (h *ProvisionModuleHandler) GetCharts(c *gin.Context) {
	hostID, err := strconv.ParseUint(c.Param("host_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	charts, err := h.provSvc.GetCharts(uint(hostID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, charts)
}

// GetChartData returns data for a specific chart.
// GET /provision-modules/charts/:host_id/:chart_id
func (h *ProvisionModuleHandler) GetChartData(c *gin.Context) {
	hostID, err := strconv.ParseUint(c.Param("host_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}
	chartID, err := strconv.ParseUint(c.Param("chart_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid chart id")
		return
	}
	period := c.DefaultQuery("period", "7d")

	data, err := h.provSvc.GetChartData(uint(hostID), uint(chartID), period)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, data)
}

// CreateChart creates a chart configuration.
// POST /provision-modules/charts
func (h *ProvisionModuleHandler) CreateChart(c *gin.Context) {
	var req service.CreateChartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	chart, err := h.provSvc.CreateChart(req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, chart)
}

// UpdateChart updates a chart configuration.
// PUT /provision-modules/charts/:id
func (h *ProvisionModuleHandler) UpdateChart(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid chart id")
		return
	}

	var req service.UpdateChartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	chart, err := h.provSvc.UpdateChart(uint(id), req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, chart)
}

// ==================== Usage Tracking ====================

// GetUsage returns usage data for a host.
// GET /provision-modules/usage/:host_id
func (h *ProvisionModuleHandler) GetUsage(c *gin.Context) {
	hostID, err := strconv.ParseUint(c.Param("host_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	usage, err := h.provSvc.GetUsage(uint(hostID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, usage)
}

// UpdateUsage updates usage data for a host.
// PUT /provision-modules/usage/:host_id
func (h *ProvisionModuleHandler) UpdateUsage(c *gin.Context) {
	hostID, err := strconv.ParseUint(c.Param("host_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	var req service.UpdateUsageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	usage, err := h.provSvc.UpdateUsage(uint(hostID), req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, usage)
}

// CheckDefineUsage checks custom usage limits.
// GET /provision-modules/usage/:host_id/check
func (h *ProvisionModuleHandler) CheckDefineUsage(c *gin.Context) {
	hostID, err := strconv.ParseUint(c.Param("host_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	result, err := h.provSvc.CheckDefineUsage(uint(hostID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, result)
}

// TrafficUsage returns traffic usage for a host.
// GET /provision-modules/usage/:host_id/traffic
func (h *ProvisionModuleHandler) TrafficUsage(c *gin.Context) {
	hostID, err := strconv.ParseUint(c.Param("host_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	result, err := h.provSvc.TrafficUsage(uint(hostID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, result)
}

// ==================== Custom Functions ====================

// GetCustomFunctions returns custom functions for a module.
// GET /provision-modules/:id/functions
func (h *ProvisionModuleHandler) GetCustomFunctions(c *gin.Context) {
	moduleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid module id")
		return
	}

	functions, err := h.provSvc.GetCustomFunctions(uint(moduleID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, functions)
}

// CreateCustomFunction creates a custom function.
// POST /provision-modules/functions
func (h *ProvisionModuleHandler) CreateCustomFunction(c *gin.Context) {
	var req service.CreateFunctionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	fn, err := h.provSvc.CreateCustomFunction(req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, fn)
}

// ExecuteCustomFunction executes a custom function.
// POST /provision-modules/functions/:id/execute
func (h *ProvisionModuleHandler) ExecuteCustomFunction(c *gin.Context) {
	fnID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid function id")
		return
	}

	var req struct {
		HostID uint `json:"host_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.provSvc.ExecuteCustomFunction(uint(fnID), req.HostID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// DeleteCustomFunction deletes a custom function.
// DELETE /provision-modules/functions/:id
func (h *ProvisionModuleHandler) DeleteCustomFunction(c *gin.Context) {
	fnID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid function id")
		return
	}

	if err := h.provSvc.DeleteCustomFunction(uint(fnID)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "function deleted")
}

// ==================== Buttons CRUD ====================

// GetButtons returns buttons for a module.
// GET /provision-modules/:id/buttons
func (h *ProvisionModuleHandler) GetButtons(c *gin.Context) {
	moduleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid module id")
		return
	}
	buttonType := c.Query("type")

	buttons, err := h.provSvc.GetButtons(uint(moduleID), buttonType)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, buttons)
}

// CreateButton creates a button.
// POST /provision-modules/buttons
func (h *ProvisionModuleHandler) CreateButton(c *gin.Context) {
	var req service.CreateButtonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	button, err := h.provSvc.CreateButton(req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, button)
}

// UpdateButton updates a button.
// PUT /provision-modules/buttons/:id
func (h *ProvisionModuleHandler) UpdateButton(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid button id")
		return
	}

	var req service.UpdateButtonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	button, err := h.provSvc.UpdateButton(uint(id), req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, button)
}

// DeleteButton deletes a button.
// DELETE /provision-modules/buttons/:id
func (h *ProvisionModuleHandler) DeleteButton(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid button id")
		return
	}

	if err := h.provSvc.DeleteButton(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "button deleted")
}

// ==================== SSL/Download ====================

// SSLButton returns SSL button configuration.
// GET /provision-modules/ssl/:host_id
func (h *ProvisionModuleHandler) SSLButton(c *gin.Context) {
	hostID, err := strconv.ParseUint(c.Param("host_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	result, err := h.provSvc.SSLButton(uint(hostID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, result)
}

// DownloadResource downloads a module resource.
// GET /provision-modules/download/:host_id/:resource_id
func (h *ProvisionModuleHandler) DownloadResource(c *gin.Context) {
	hostID, err := strconv.ParseUint(c.Param("host_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}
	resourceID, err := strconv.ParseUint(c.Param("resource_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid resource id")
		return
	}

	result, err := h.provSvc.DownloadResource(uint(hostID), uint(resourceID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, result)
}

// ==================== Flow Packets ====================

// AfterFlowPacketPaid handles post-payment hook for flow packets.
// POST /provision-modules/flow-packet/paid
func (h *ProvisionModuleHandler) AfterFlowPacketPaid(c *gin.Context) {
	var req struct {
		HostID   uint `json:"host_id" binding:"required"`
		PacketID uint `json:"packet_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.provSvc.AfterFlowPacketPaid(req.HostID, req.PacketID); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "flow packet processed")
}
