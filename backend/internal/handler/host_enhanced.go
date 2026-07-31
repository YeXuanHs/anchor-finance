package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// HostEnhancedHandler handles enhanced host management endpoints.
type HostEnhancedHandler struct {
	svc *service.HostEnhancedService
	log *logger.Logger
}

func NewHostEnhancedHandler(svc *service.HostEnhancedService, log *logger.Logger) *HostEnhancedHandler {
	return &HostEnhancedHandler{svc: svc, log: log}
}

// ═══════════════════ Renewal ═══════════════════

// GetRenewalPage returns renewal options for a host.
func (h *HostEnhancedHandler) GetRenewalPage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	data, err := h.svc.GetRenewalPage(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, data)
}

// GetRenewalPrice calculates renewal price.
func (h *HostEnhancedHandler) GetRenewalPrice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	cycle := c.Query("cycle")
	if cycle == "" {
		response.BadRequest(c, "cycle is required")
		return
	}

	data, err := h.svc.GetRenewalPrice(uint(id), cycle)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

// SubmitRenewal creates a renewal order.
func (h *HostEnhancedHandler) SubmitRenewal(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	var req struct {
		Cycle         string `json:"cycle" binding:"required"`
		PaymentMethod string `json:"payment_method"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetUint("user_id")
	result, err := h.svc.SubmitRenewal(userID, uint(id), req.Cycle, req.PaymentMethod)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// SetAutoRenew toggles auto-renewal.
func (h *HostEnhancedHandler) SetAutoRenew(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.SetAutoRenew(uint(id), req.Enabled); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "auto-renew updated")
}

// BatchRenew renews multiple hosts.
func (h *HostEnhancedHandler) BatchRenew(c *gin.Context) {
	var req struct {
		HostIDs []uint `json:"host_ids" binding:"required,min=1"`
		Cycle   string `json:"cycle" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetUint("user_id")
	summary, err := h.svc.BatchRenew(userID, req.HostIDs, req.Cycle)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, summary)
}

// ═══════════════════ Transfer ═══════════════════

// TransferHost transfers a host to another user.
func (h *HostEnhancedHandler) TransferHost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	var req struct {
		ToUserID uint   `json:"to_user_id" binding:"required"`
		Reason   string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetUint("user_id")
	if err := h.svc.TransferHost(uint(id), userID, req.ToUserID, req.Reason); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "host transferred")
}

// GetTransferHistory returns transfer history for a host.
func (h *HostEnhancedHandler) GetTransferHistory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	logs, err := h.svc.GetTransferHistory(uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, logs)
}

// ═══════════════════ Categories ═══════════════════

// GetHostCategories returns user's host categories.
func (h *HostEnhancedHandler) GetHostCategories(c *gin.Context) {
	userID := c.GetUint("user_id")
	cats, err := h.svc.GetHostCategories(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cats)
}

// CreateHostCategory creates a new category.
func (h *HostEnhancedHandler) CreateHostCategory(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetUint("user_id")
	cat, err := h.svc.CreateHostCategory(userID, req.Name)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cat)
}

// AssignCategory assigns a host to a category.
func (h *HostEnhancedHandler) AssignCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	var req struct {
		CategoryID uint `json:"category_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.AssignCategory(uint(id), req.CategoryID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "category assigned")
}

// DeleteCategory deletes a category.
func (h *HostEnhancedHandler) DeleteCategory(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("category_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid category id")
		return
	}

	userID := c.GetUint("user_id")
	if err := h.svc.DeleteCategory(userID, uint(categoryID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "category deleted")
}

// ═══════════════════ SSL Config ═══════════════════

// GetSSLConfig returns SSL config for a host.
func (h *HostEnhancedHandler) GetSSLConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	cfg, err := h.svc.GetSSLConfig(uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cfg)
}

// SetSSLConfig updates SSL config.
func (h *HostEnhancedHandler) SetSSLConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	var req struct {
		Type   string `json:"type" binding:"required,oneof=letsencrypt custom"`
		Domain string `json:"domain" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	cfg, err := h.svc.SetSSLConfig(uint(id), req.Type, req.Domain)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cfg)
}

// InstallSSL installs an SSL certificate.
func (h *HostEnhancedHandler) InstallSSL(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	var req struct {
		Cert string `json:"cert" binding:"required"`
		Key  string `json:"key" binding:"required"`
		CA   string `json:"ca"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.InstallSSL(uint(id), req.Cert, req.Key, req.CA); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "SSL certificate installed")
}

// ═══════════════════ Downstream / Second Verify ═══════════════════

// SetDownstream sets a downstream provider.
func (h *HostEnhancedHandler) SetDownstream(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	var req struct {
		DownstreamID string `json:"downstream_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.SetDownstream(uint(id), req.DownstreamID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "downstream set")
}

// SetSecondVerify toggles second verification.
func (h *HostEnhancedHandler) SetSecondVerify(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.SetSecondVerify(uint(id), req.Enabled); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "second verify updated")
}

// VerifySecond verifies a second factor code.
func (h *HostEnhancedHandler) VerifySecond(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	ok, err := h.svc.VerifySecond(uint(id), req.Code)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	if !ok {
		response.BadRequest(c, "invalid verification code")
		return
	}
	response.SuccessMsg(c, "verified")
}

// ═══════════════════ Traffic & Power ═══════════════════

// GetTrafficUsage returns traffic usage for a host.
func (h *HostEnhancedHandler) GetTrafficUsage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	data, err := h.svc.GetTrafficUsage(uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

// GetTrafficChart returns traffic chart data.
func (h *HostEnhancedHandler) GetTrafficChart(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	period := c.DefaultQuery("period", "30d")
	data, err := h.svc.GetTrafficChart(uint(id), period)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

// RefreshPowerStatus refreshes power status from upstream.
func (h *HostEnhancedHandler) RefreshPowerStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	data, err := h.svc.RefreshPowerStatus(uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

// ═══════════════════ Host Info ═══════════════════

// GetHostDetail returns full host detail.
func (h *HostEnhancedHandler) GetHostDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	data, err := h.svc.GetHostDetail(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, data)
}

// GetHostStatus returns real-time status.
func (h *HostEnhancedHandler) GetHostStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	data, err := h.svc.GetHostStatus(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, data)
}

// PostRemark updates host remark.
func (h *HostEnhancedHandler) PostRemark(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	var req struct {
		Remark string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetUint("user_id")
	if err := h.svc.PostRemark(uint(id), userID, req.Remark); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "remark updated")
}

// HideHost hides a host from user list.
func (h *HostEnhancedHandler) HideHost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	userID := c.GetUint("user_id")
	if err := h.svc.HideHost(uint(id), userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "host hidden")
}

// ═══════════════════ Cancel / Terminate ═══════════════════

// GetCancelPage returns cancellation options.
func (h *HostEnhancedHandler) GetCancelPage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	userID := c.GetUint("user_id")
	data, err := h.svc.GetCancelPage(uint(id), userID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

// SubmitCancel submits a cancellation request.
func (h *HostEnhancedHandler) SubmitCancel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	var req struct {
		Reason     string `json:"reason"`
		CancelType string `json:"cancel_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetUint("user_id")
	result, err := h.svc.SubmitCancel(uint(id), userID, req.Reason, req.CancelType)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// DeleteCancel cancels a pending cancellation request.
func (h *HostEnhancedHandler) DeleteCancel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	userID := c.GetUint("user_id")
	if err := h.svc.DeleteCancel(uint(id), userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "cancellation request removed")
}

// ═══════════════════ Dedicated Server ═══════════════════

// GetDedicatedServer returns dedicated server info.
func (h *HostEnhancedHandler) GetDedicatedServer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	data, err := h.svc.GetDedicatedServer(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, data)
}

// ═══════════════════ Flow Packets ═══════════════════

// GetFlowPackets lists available flow packets.
func (h *HostEnhancedHandler) GetFlowPackets(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	data, err := h.svc.GetFlowPackets(uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

// BuyFlowPacket purchases a flow packet.
func (h *HostEnhancedHandler) BuyFlowPacket(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	packetID, err := strconv.ParseUint(c.Param("packet_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid packet id")
		return
	}

	userID := c.GetUint("user_id")
	result, err := h.svc.BuyFlowPacket(userID, uint(id), uint(packetID))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// ═══════════════════ Reinstall (enhanced) ═══════════════════

// CheckReinstall checks if reinstall is allowed.
func (h *HostEnhancedHandler) CheckReinstall(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	data, err := h.svc.CheckReinstall(uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

// GetReinstallStatus returns reinstall progress.
func (h *HostEnhancedHandler) GetReinstallStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	data, err := h.svc.GetReinstallStatus(uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

// CancelReinstall cancels an in-progress reinstall.
func (h *HostEnhancedHandler) CancelReinstall(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	if err := h.svc.CancelReinstall(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "reinstall cancelled")
}

// ═══════════════════ Host Recharge ═══════════════════

// GetHostRecharge returns recharge options.
func (h *HostEnhancedHandler) GetHostRecharge(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	userID := c.GetUint("user_id")
	data, err := h.svc.GetHostRecharge(uint(id), userID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}
