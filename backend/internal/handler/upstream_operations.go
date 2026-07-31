package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type UpstreamOpsHandler struct {
	svc *service.UpstreamOperations
	log *logger.Logger
}

func NewUpstreamOpsHandler(upstreamSvc *service.UpstreamService, log *logger.Logger) *UpstreamOpsHandler {
	return &UpstreamOpsHandler{
		svc: service.NewUpstreamOperations(upstreamSvc),
		log: log,
	}
}

// ==================== Power Operations ====================

// Boot powers on a server via the upstream provider.
func (h *UpstreamOpsHandler) Boot(c *gin.Context) {
	providerID, hostID, ok := parseProviderHostIDs(c)
	if !ok {
		return
	}
	operatorID := c.GetUint("user_id")

	result, err := h.svc.UpstreamBoot(providerID, hostID, operatorID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// Shutdown powers off a server via the upstream provider.
func (h *UpstreamOpsHandler) Shutdown(c *gin.Context) {
	providerID, hostID, ok := parseProviderHostIDs(c)
	if !ok {
		return
	}
	operatorID := c.GetUint("user_id")

	result, err := h.svc.UpstreamShutdown(providerID, hostID, operatorID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// Reboot reboots a server via the upstream provider.
func (h *UpstreamOpsHandler) Reboot(c *gin.Context) {
	providerID, hostID, ok := parseProviderHostIDs(c)
	if !ok {
		return
	}
	operatorID := c.GetUint("user_id")

	result, err := h.svc.UpstreamReboot(providerID, hostID, operatorID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// GetStatus gets server status from the upstream provider.
func (h *UpstreamOpsHandler) GetStatus(c *gin.Context) {
	providerID, hostID, ok := parseProviderHostIDs(c)
	if !ok {
		return
	}
	operatorID := c.GetUint("user_id")

	result, err := h.svc.UpstreamGetStatus(providerID, hostID, operatorID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// ==================== Console Access ====================

// VNC gets the VNC console URL from the upstream provider.
func (h *UpstreamOpsHandler) VNC(c *gin.Context) {
	providerID, hostID, ok := parseProviderHostIDs(c)
	if !ok {
		return
	}
	operatorID := c.GetUint("user_id")

	result, err := h.svc.UpstreamVNC(providerID, hostID, operatorID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// KVM gets the KVM console URL from the upstream provider.
func (h *UpstreamOpsHandler) KVM(c *gin.Context) {
	providerID, hostID, ok := parseProviderHostIDs(c)
	if !ok {
		return
	}
	operatorID := c.GetUint("user_id")

	result, err := h.svc.UpstreamKVM(providerID, hostID, operatorID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// IPMIStatus gets IPMI status from the upstream provider.
func (h *UpstreamOpsHandler) IPMIStatus(c *gin.Context) {
	providerID, hostID, ok := parseProviderHostIDs(c)
	if !ok {
		return
	}
	operatorID := c.GetUint("user_id")

	result, err := h.svc.UpstreamIPMIStatus(providerID, hostID, operatorID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// IPMIOn powers on via IPMI through the upstream provider.
func (h *UpstreamOpsHandler) IPMIOn(c *gin.Context) {
	providerID, hostID, ok := parseProviderHostIDs(c)
	if !ok {
		return
	}
	operatorID := c.GetUint("user_id")

	result, err := h.svc.UpstreamIPMIOn(providerID, hostID, operatorID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// IPMIOff powers off via IPMI through the upstream provider.
func (h *UpstreamOpsHandler) IPMIOff(c *gin.Context) {
	providerID, hostID, ok := parseProviderHostIDs(c)
	if !ok {
		return
	}
	operatorID := c.GetUint("user_id")

	result, err := h.svc.UpstreamIPMIOff(providerID, hostID, operatorID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// IPMIReboot reboots via IPMI through the upstream provider.
func (h *UpstreamOpsHandler) IPMIReboot(c *gin.Context) {
	providerID, hostID, ok := parseProviderHostIDs(c)
	if !ok {
		return
	}
	operatorID := c.GetUint("user_id")

	result, err := h.svc.UpstreamIPMIReboot(providerID, hostID, operatorID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// IPMIVNC gets the IPMI VNC console from the upstream provider.
func (h *UpstreamOpsHandler) IPMIVNC(c *gin.Context) {
	providerID, hostID, ok := parseProviderHostIDs(c)
	if !ok {
		return
	}
	operatorID := c.GetUint("user_id")

	result, err := h.svc.UpstreamIPMIVNC(providerID, hostID, operatorID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// ==================== Reinstall ====================

// Reinstall initiates OS reinstall via the upstream provider.
func (h *UpstreamOpsHandler) Reinstall(c *gin.Context) {
	providerID, hostID, ok := parseProviderHostIDs(c)
	if !ok {
		return
	}

	var req struct {
		OS       string `json:"os" binding:"required"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	operatorID := c.GetUint("user_id")

	result, err := h.svc.UpstreamReinstall(providerID, hostID, operatorID, req.OS, req.Password)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// GetReinstallStatus checks the reinstall progress from the upstream.
func (h *UpstreamOpsHandler) GetReinstallStatus(c *gin.Context) {
	providerID, hostID, ok := parseProviderHostIDs(c)
	if !ok {
		return
	}
	operatorID := c.GetUint("user_id")

	result, err := h.svc.UpstreamGetReinstallStatus(providerID, hostID, operatorID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// CancelReinstall cancels an in-progress reinstall via the upstream.
func (h *UpstreamOpsHandler) CancelReinstall(c *gin.Context) {
	providerID, hostID, ok := parseProviderHostIDs(c)
	if !ok {
		return
	}
	operatorID := c.GetUint("user_id")

	result, err := h.svc.UpstreamCancelReinstall(providerID, hostID, operatorID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// GetOSList gets the available OS list from the upstream provider.
func (h *UpstreamOpsHandler) GetOSList(c *gin.Context) {
	providerID, err := strconv.ParseUint(c.Param("provider_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid provider id")
		return
	}
	operatorID := c.GetUint("user_id")

	result, err := h.svc.UpstreamGetOSList(uint(providerID), operatorID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// CrackPassword resets a server's password via the upstream provider.
func (h *UpstreamOpsHandler) CrackPassword(c *gin.Context) {
	providerID, hostID, ok := parseProviderHostIDs(c)
	if !ok {
		return
	}
	operatorID := c.GetUint("user_id")

	result, err := h.svc.UpstreamCrackPassword(providerID, hostID, operatorID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// ==================== DCIM Client Operations ====================

// DcimClientStatus gets the DCIM client status from the upstream.
func (h *UpstreamOpsHandler) DcimClientStatus(c *gin.Context) {
	providerID, hostID, ok := parseProviderHostIDs(c)
	if !ok {
		return
	}
	operatorID := c.GetUint("user_id")

	result, err := h.svc.DcimClientStatus(providerID, hostID, operatorID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// DcimClientOn powers on via DCIM client through the upstream.
func (h *UpstreamOpsHandler) DcimClientOn(c *gin.Context) {
	providerID, hostID, ok := parseProviderHostIDs(c)
	if !ok {
		return
	}
	operatorID := c.GetUint("user_id")

	result, err := h.svc.DcimClientOn(providerID, hostID, operatorID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// DcimClientOff powers off via DCIM client through the upstream.
func (h *UpstreamOpsHandler) DcimClientOff(c *gin.Context) {
	providerID, hostID, ok := parseProviderHostIDs(c)
	if !ok {
		return
	}
	operatorID := c.GetUint("user_id")

	result, err := h.svc.DcimClientOff(providerID, hostID, operatorID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// DcimClientReboot reboots via DCIM client through the upstream.
func (h *UpstreamOpsHandler) DcimClientReboot(c *gin.Context) {
	providerID, hostID, ok := parseProviderHostIDs(c)
	if !ok {
		return
	}
	operatorID := c.GetUint("user_id")

	result, err := h.svc.DcimClientReboot(providerID, hostID, operatorID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// DcimClientVNC gets the DCIM client VNC console from the upstream.
func (h *UpstreamOpsHandler) DcimClientVNC(c *gin.Context) {
	providerID, hostID, ok := parseProviderHostIDs(c)
	if !ok {
		return
	}
	operatorID := c.GetUint("user_id")

	result, err := h.svc.DcimClientVNC(providerID, hostID, operatorID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// DcimClientReinstall reinstalls via DCIM client through the upstream.
func (h *UpstreamOpsHandler) DcimClientReinstall(c *gin.Context) {
	providerID, hostID, ok := parseProviderHostIDs(c)
	if !ok {
		return
	}

	var req struct {
		OS string `json:"os" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	operatorID := c.GetUint("user_id")

	result, err := h.svc.DcimClientReinstall(providerID, hostID, operatorID, req.OS)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// DcimClientCrackPass resets password via DCIM client through the upstream.
func (h *UpstreamOpsHandler) DcimClientCrackPass(c *gin.Context) {
	providerID, hostID, ok := parseProviderHostIDs(c)
	if !ok {
		return
	}
	operatorID := c.GetUint("user_id")

	result, err := h.svc.DcimClientCrackPass(providerID, hostID, operatorID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// DcimClientCancelReinstall cancels reinstall via DCIM client through the upstream.
func (h *UpstreamOpsHandler) DcimClientCancelReinstall(c *gin.Context) {
	providerID, hostID, ok := parseProviderHostIDs(c)
	if !ok {
		return
	}
	operatorID := c.GetUint("user_id")

	result, err := h.svc.DcimClientCancelReinstall(providerID, hostID, operatorID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// DcimClientReinstallStatus gets the DCIM client reinstall status from the upstream.
func (h *UpstreamOpsHandler) DcimClientReinstallStatus(c *gin.Context) {
	providerID, hostID, ok := parseProviderHostIDs(c)
	if !ok {
		return
	}
	operatorID := c.GetUint("user_id")

	result, err := h.svc.DcimClientReinstallStatus(providerID, hostID, operatorID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// DcimClientGetOS gets the OS list from the DCIM client through the upstream.
func (h *UpstreamOpsHandler) DcimClientGetOS(c *gin.Context) {
	providerID, err := strconv.ParseUint(c.Param("provider_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid provider id")
		return
	}
	operatorID := c.GetUint("user_id")

	result, err := h.svc.DcimClientGetOS(uint(providerID), operatorID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// ==================== Module Buttons ====================

// ModuleClientButton executes a client-side module button via the upstream.
func (h *UpstreamOpsHandler) ModuleClientButton(c *gin.Context) {
	providerID, hostID, ok := parseProviderHostIDs(c)
	if !ok {
		return
	}

	var req struct {
		Button string `json:"button" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	operatorID := c.GetUint("user_id")

	result, err := h.svc.ModuleClientButton(providerID, hostID, operatorID, req.Button)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// ModuleAdminButton executes an admin-side module button via the upstream.
func (h *UpstreamOpsHandler) ModuleAdminButton(c *gin.Context) {
	providerID, hostID, ok := parseProviderHostIDs(c)
	if !ok {
		return
	}

	var req struct {
		Button string `json:"button" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	operatorID := c.GetUint("user_id")

	result, err := h.svc.ModuleAdminButton(providerID, hostID, operatorID, req.Button)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// ModulePowerStatus gets the power status from the upstream module.
func (h *UpstreamOpsHandler) ModulePowerStatus(c *gin.Context) {
	providerID, hostID, ok := parseProviderHostIDs(c)
	if !ok {
		return
	}
	operatorID := c.GetUint("user_id")

	result, err := h.svc.ModulePowerStatus(providerID, hostID, operatorID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// ==================== Helper ====================

// parseProviderHostIDs extracts provider_id and host_id from the URL params.
func parseProviderHostIDs(c *gin.Context) (uint, uint, bool) {
	providerID, err := strconv.ParseUint(c.Param("provider_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid provider id")
		return 0, 0, false
	}
	hostID, err := strconv.ParseUint(c.Param("host_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return 0, 0, false
	}
	return uint(providerID), uint(hostID), true
}
