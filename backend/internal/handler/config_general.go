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

// ==================== Payment Config ====================

func (h *ConfigGeneralHandler) GetPaymentConfig(c *gin.Context) {
	cfg, err := h.svc.GetPaymentConfig()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cfg)
}

func (h *ConfigGeneralHandler) UpdatePaymentConfig(c *gin.Context) {
	var req service.PaymentConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdatePaymentConfig(req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "payment config updated")
}

// ==================== SMS Config ====================

func (h *ConfigGeneralHandler) GetSmsConfig(c *gin.Context) {
	cfg, err := h.svc.GetSmsConfig()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cfg)
}

func (h *ConfigGeneralHandler) UpdateSmsConfig(c *gin.Context) {
	var req service.SmsConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateSmsConfig(req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "SMS config updated")
}

// ==================== Security Config ====================

func (h *ConfigGeneralHandler) GetSecurityConfig(c *gin.Context) {
	cfg, err := h.svc.GetSecurityConfig()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cfg)
}

func (h *ConfigGeneralHandler) UpdateSecurityConfig(c *gin.Context) {
	var req service.SecurityConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateSecurityConfig(req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "security config updated")
}

// ==================== Local Config ====================

func (h *ConfigGeneralHandler) GetLocalConfig(c *gin.Context) {
	cfg, err := h.svc.GetLocalConfig()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cfg)
}

func (h *ConfigGeneralHandler) UpdateLocalConfig(c *gin.Context) {
	var req service.LocalConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateLocalConfig(req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "local config updated")
}

// ==================== Affiliate Config ====================

func (h *ConfigGeneralHandler) GetAffiliateConfig(c *gin.Context) {
	cfg, err := h.svc.GetAffiliateConfig()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cfg)
}

func (h *ConfigGeneralHandler) UpdateAffiliateConfig(c *gin.Context) {
	var req service.AffiliateConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateAffiliateConfig(req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "affiliate config updated")
}

// ==================== Captcha Config ====================

func (h *ConfigGeneralHandler) GetCaptchaConfig(c *gin.Context) {
	cfg, err := h.svc.GetCaptchaConfig()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cfg)
}

func (h *ConfigGeneralHandler) UpdateCaptchaConfig(c *gin.Context) {
	var req service.CaptchaConfigData
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateCaptchaConfig(req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "captcha config updated")
}

// ==================== Buy Product Config ====================

func (h *ConfigGeneralHandler) GetBuyProductConfig(c *gin.Context) {
	cfg, err := h.svc.GetBuyProductConfig()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cfg)
}

func (h *ConfigGeneralHandler) UpdateBuyProductConfig(c *gin.Context) {
	var req service.BuyProductConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateBuyProductConfig(req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "buy product config updated")
}

// ==================== Second Verify Config ====================

func (h *ConfigGeneralHandler) GetSecondVerifyConfig(c *gin.Context) {
	cfg, err := h.svc.GetSecondVerifyConfig()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cfg)
}

func (h *ConfigGeneralHandler) UpdateSecondVerifyConfig(c *gin.Context) {
	var req service.SecondVerifyConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateSecondVerifyConfig(req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "second verify config updated")
}

// ==================== Nav Group ====================

func (h *ConfigGeneralHandler) GetNavGroups(c *gin.Context) {
	groups, err := h.svc.GetNavGroups()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, groups)
}

func (h *ConfigGeneralHandler) CreateNavGroup(c *gin.Context) {
	var req service.NavGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.CreateNavGroup(req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "nav group created")
}

func (h *ConfigGeneralHandler) UpdateNavGroup(c *gin.Context) {
	var req service.NavGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateNavGroup(req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "nav group updated")
}

func (h *ConfigGeneralHandler) DeleteNavGroup(c *gin.Context) {
	var req struct {
		ID     uint `json:"id" binding:"required"`
		ToID   uint `json:"to_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.DeleteNavGroup(req.ID, req.ToID); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "nav group deleted")
}

// ==================== Language Config ====================

func (h *ConfigGeneralHandler) GetLanguageConfig(c *gin.Context) {
	cfg, err := h.svc.GetLanguageConfig()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cfg)
}

func (h *ConfigGeneralHandler) SetAdminLanguage(c *gin.Context) {
	var req struct {
		Lang string `json:"lang" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.SetAdminLanguage(req.Lang); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "language set")
}

// ==================== Header/Footer Config ====================

func (h *ConfigGeneralHandler) GetHeaderFooter(c *gin.Context) {
	cfg, err := h.svc.GetHeaderFooter()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cfg)
}

func (h *ConfigGeneralHandler) UpdateHeaderFooter(c *gin.Context) {
	var req service.HeaderFooterConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateHeaderFooter(req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "header/footer config updated")
}

// ==================== New Login Page Config ====================

func (h *ConfigGeneralHandler) GetNewLoginPageConfig(c *gin.Context) {
	cfg, err := h.svc.GetNewLoginPageConfig()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cfg)
}

func (h *ConfigGeneralHandler) UpdateNewLoginPageConfig(c *gin.Context) {
	var req service.NewLoginPageConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateNewLoginPageConfig(req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "new login page config updated")
}
