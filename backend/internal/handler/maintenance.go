package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// MaintenanceHandler handles maintenance mode requests.
type MaintenanceHandler struct {
	svc *service.MaintenanceService
	log *logger.Logger
}

// NewMaintenanceHandler creates a new MaintenanceHandler.
func NewMaintenanceHandler(svc *service.MaintenanceService, log *logger.Logger) *MaintenanceHandler {
	return &MaintenanceHandler{svc: svc, log: log}
}

// GetStatus returns the current maintenance mode status.
func (h *MaintenanceHandler) GetStatus(c *gin.Context) {
	status, err := h.svc.GetStatus()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, status)
}

// Enable enables maintenance mode.
func (h *MaintenanceHandler) Enable(c *gin.Context) {
	var req struct {
		Message     string `json:"message"`
		AllowedIPs  string `json:"allowed_ips"`
		EstimatedAt string `json:"estimated_at"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.Enable(req.Message, req.AllowedIPs, req.EstimatedAt); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "maintenance mode enabled")
}

// Disable disables maintenance mode.
func (h *MaintenanceHandler) Disable(c *gin.Context) {
	if err := h.svc.Disable(); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "maintenance mode disabled")
}

// Update updates maintenance mode settings.
func (h *MaintenanceHandler) Update(c *gin.Context) {
	var req struct {
		Message     *string `json:"message"`
		AllowedIPs  *string `json:"allowed_ips"`
		EstimatedAt *string `json:"estimated_at"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Message != nil {
		updates["message"] = *req.Message
	}
	if req.AllowedIPs != nil {
		updates["allowed_ips"] = *req.AllowedIPs
	}
	if req.EstimatedAt != nil {
		updates["estimated_at"] = *req.EstimatedAt
	}

	if err := h.svc.Update(updates); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "maintenance settings updated")
}

// GetHistory returns maintenance mode history.
func (h *MaintenanceHandler) GetHistory(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	items, total, err := h.svc.GetHistory(page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, items, total, page, pageSize)
}

// AddAllowedIP adds an IP to the allowed list.
func (h *MaintenanceHandler) AddAllowedIP(c *gin.Context) {
	var req struct {
		IP string `json:"ip" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.AddAllowedIP(req.IP); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "ip added to allowed list")
}

// RemoveAllowedIP removes an IP from the allowed list.
func (h *MaintenanceHandler) RemoveAllowedIP(c *gin.Context) {
	ip := c.Query("ip")
	if ip == "" {
		response.BadRequest(c, "ip is required")
		return
	}

	if err := h.svc.RemoveAllowedIP(ip); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "ip removed from allowed list")
}

// GetAllowedIPs returns the list of allowed IPs.
func (h *MaintenanceHandler) GetAllowedIPs(c *gin.Context) {
	ips, err := h.svc.GetAllowedIPs()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, ips)
}

// TestMode tests maintenance mode display without enabling it.
func (h *MaintenanceHandler) TestMode(c *gin.Context) {
	var req struct {
		Message string `json:"message"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.svc.TestMode(req.Message)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, result)
}
