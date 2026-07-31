package handler

import (
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type ConfigGeneralHandler struct {
	svc *service.ConfigGeneralService
	log *logger.Logger
}

func NewConfigGeneralHandler(svc *service.ConfigGeneralService, log *logger.Logger) *ConfigGeneralHandler {
	return &ConfigGeneralHandler{svc: svc, log: log}
}

// Get is an alias for GetConfig.
func (h *ConfigGeneralHandler) Get(c *gin.Context) { h.GetConfig(c) }

// Update is an alias for UpdateConfig.
func (h *ConfigGeneralHandler) Update(c *gin.Context) { h.UpdateConfig(c) }

// ==================== General Config ====================

// GetConfig returns the general site configuration.
func (h *ConfigGeneralHandler) GetConfig(c *gin.Context) {
	cfg, err := h.svc.GetConfig()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cfg)
}

// UpdateConfig updates the general site configuration.
func (h *ConfigGeneralHandler) UpdateConfig(c *gin.Context) {
	var req service.GeneralConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.UpdateConfig(req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "general config updated")
}

// ==================== Email Config ====================

func (h *ConfigGeneralHandler) GetEmailConfig(c *gin.Context) {
	cfg, err := h.svc.GetEmailConfig()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cfg)
}

func (h *ConfigGeneralHandler) UpdateEmailConfig(c *gin.Context) {
	var req service.EmailConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateEmailConfig(req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "email config updated")
}

// ==================== Email Support ====================

func (h *ConfigGeneralHandler) GetEmailSupport(c *gin.Context) {
	cfg, err := h.svc.GetEmailSupport()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cfg)
}

func (h *ConfigGeneralHandler) UpdateEmailSupport(c *gin.Context) {
	var req service.EmailSupportConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateEmailSupport(req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "email support config updated")
}

// ==================== Affiliate Ladders ====================

func (h *ConfigGeneralHandler) GetAffiliateLadders(c *gin.Context) {
	cfg, err := h.svc.GetAffiliateLadders()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cfg)
}

func (h *ConfigGeneralHandler) UpdateAffiliateLadders(c *gin.Context) {
	var req service.AffiliateLaddersConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateAffiliateLadders(req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "affiliate ladders updated")
}

// ==================== Safe Config ====================

func (h *ConfigGeneralHandler) GetSafeConfig(c *gin.Context) {
	cfg, err := h.svc.GetSafeConfig()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cfg)
}

func (h *ConfigGeneralHandler) UpdateSafeConfig(c *gin.Context) {
	var req service.SafeConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateSafeConfig(req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "safe config updated")
}

// ==================== Recharge Config ====================

func (h *ConfigGeneralHandler) GetRechargeConfig(c *gin.Context) {
	cfg, err := h.svc.GetRechargeConfig()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cfg)
}

func (h *ConfigGeneralHandler) UpdateRechargeConfig(c *gin.Context) {
	var req service.RechargeConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateRechargeConfig(req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "recharge config updated")
}

// ==================== Invoice Config ====================

func (h *ConfigGeneralHandler) GetInvoiceConfig(c *gin.Context) {
	cfg, err := h.svc.GetInvoiceConfig()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cfg)
}

func (h *ConfigGeneralHandler) UpdateInvoiceConfig(c *gin.Context) {
	var req service.InvoiceConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateInvoiceConfig(req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "invoice config updated")
}

// ==================== Register Config ====================

func (h *ConfigGeneralHandler) GetRegisterConfig(c *gin.Context) {
	cfg, err := h.svc.GetRegisterConfig()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cfg)
}

func (h *ConfigGeneralHandler) UpdateRegisterConfig(c *gin.Context) {
	var req service.RegisterConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateRegisterConfig(req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "register config updated")
}

// ==================== Login Config ====================

func (h *ConfigGeneralHandler) GetLoginConfig(c *gin.Context) {
	cfg, err := h.svc.GetLoginConfig()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cfg)
}

func (h *ConfigGeneralHandler) UpdateLoginConfig(c *gin.Context) {
	var req service.LoginConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateLoginConfig(req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "login config updated")
}

// ==================== API Config ====================

func (h *ConfigGeneralHandler) GetAPIConfig(c *gin.Context) {
	cfg, err := h.svc.GetAPIConfig()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cfg)
}

func (h *ConfigGeneralHandler) UpdateAPIConfig(c *gin.Context) {
	var req service.APIConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateAPIConfig(req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "API config updated")
}

// ==================== 2FA Config ====================

func (h *ConfigGeneralHandler) GetTwoFactorConfig(c *gin.Context) {
	cfg, err := h.svc.GetTwoFactorConfig()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cfg)
}

func (h *ConfigGeneralHandler) UpdateTwoFactorConfig(c *gin.Context) {
	var req service.TwoFactorConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateTwoFactorConfig(req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "2FA config updated")
}

// ==================== Debug Mode ====================

func (h *ConfigGeneralHandler) GetDebugMode(c *gin.Context) {
	enabled, err := h.svc.GetDebugMode()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"enabled": enabled})
}

func (h *ConfigGeneralHandler) SetDebugMode(c *gin.Context) {
	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.SetDebugMode(req.Enabled); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "debug mode updated")
}

// ==================== SMTP Test ====================

func (h *ConfigGeneralHandler) TestSMTP(c *gin.Context) {
	var req struct {
		ToEmail string `json:"to_email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.TestSMTP(req.ToEmail); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "test email sent")
}

// ==================== SMS Test ====================

func (h *ConfigGeneralHandler) TestSMS(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.TestSMS(req.Phone); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "test SMS sent")
}
