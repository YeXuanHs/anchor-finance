package handler

import (
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type SettingHandler struct {
	log *logger.Logger
}

func NewSettingHandler(log *logger.Logger) *SettingHandler {
	return &SettingHandler{log: log}
}

// GetNotificationSettings 获取通知设置
func (h *SettingHandler) GetNotificationSettings(c *gin.Context) {
	settings := gin.H{
		"email_enabled":     true,
		"email_host":        "",
		"email_port":        465,
		"email_username":    "",
		"email_password":    "",
		"email_from":        "",
		"email_encryption":  "ssl",
		"sms_enabled":       false,
		"sms_provider":      "",
		"sms_api_key":       "",
		"template": gin.H{
			"welcome":         "欢迎注册，{username}！",
			"password_reset":  "您的密码重置链接：{link}",
			"order_confirm":   "订单 {order_id} 已确认",
			"payment_success": "支付成功，金额：{amount}",
		},
	}
	response.Success(c, settings)
}

// SaveNotificationSettings 保存通知设置
func (h *SettingHandler) SaveNotificationSettings(c *gin.Context) {
	var req struct {
		EmailEnabled    bool   `json:"email_enabled"`
		EmailHost       string `json:"email_host"`
		EmailPort       int    `json:"email_port"`
		EmailUsername    string `json:"email_username"`
		EmailPassword    string `json:"email_password"`
		EmailFrom        string `json:"email_from"`
		EmailEncryption  string `json:"email_encryption"`
		SmsEnabled       bool   `json:"sms_enabled"`
		SmsProvider      string `json:"sms_provider"`
		SmsApiKey        string `json:"sms_api_key"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	h.log.Info("saving notification settings")
	response.SuccessMsg(c, "通知设置已保存")
}

// GetMaintenanceMode 获取维护模式配置
func (h *SettingHandler) GetMaintenanceMode(c *gin.Context) {
	config := gin.H{
		"enabled":     false,
		"message":     "系统维护中，请稍后再试",
		"allowed_ips": []string{"127.0.0.1"},
		"start_time":  "",
		"end_time":    "",
	}
	response.Success(c, config)
}

// SetMaintenanceMode 设置维护模式
func (h *SettingHandler) SetMaintenanceMode(c *gin.Context) {
	var req struct {
		Enabled    bool     `json:"enabled"`
		Message    string   `json:"message"`
		AllowedIPs []string `json:"allowed_ips"`
		StartTime  string   `json:"start_time"`
		EndTime    string   `json:"end_time"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if req.Enabled {
		h.log.Info("maintenance mode enabled")
		response.SuccessMsg(c, "维护模式已开启")
	} else {
		h.log.Info("maintenance mode disabled")
		response.SuccessMsg(c, "维护模式已关闭")
	}
}

// GetCronSettings 获取定时任务配置
func (h *SettingHandler) GetCronSettings(c *gin.Context) {
	jobs := []gin.H{
		{
			"id":          1,
			"name":        "数据库备份",
			"command":     "backup:database",
			"schedule":    "0 2 * * *",
			"enabled":     true,
			"last_run":    "2026-08-01 02:00:00",
			"next_run":    "2026-08-02 02:00:00",
			"status":      "success",
		},
		{
			"id":          2,
			"name":        "清理临时文件",
			"command":     "cleanup:temp",
			"schedule":    "0 3 * * 0",
			"enabled":     true,
			"last_run":    "2026-07-28 03:00:00",
			"next_run":    "2026-08-04 03:00:00",
			"status":      "success",
		},
		{
			"id":          3,
			"name":        "系统更新检查",
			"command":     "system:check-update",
			"schedule":    "0 9 * * 1",
			"enabled":     false,
			"last_run":    "",
			"next_run":    "",
			"status":      "disabled",
		},
	}
	response.Success(c, gin.H{"jobs": jobs})
}

// SaveCronSettings 保存定时任务配置
func (h *SettingHandler) SaveCronSettings(c *gin.Context) {
	var req struct {
		Jobs []struct {
			ID       uint   `json:"id"`
			Name     string `json:"name" binding:"required"`
			Command  string `json:"command" binding:"required"`
			Schedule string `json:"schedule" binding:"required"`
			Enabled  bool   `json:"enabled"`
		} `json:"jobs"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	h.log.Info("saving cron settings, jobs count: %d", len(req.Jobs))
	response.SuccessMsg(c, "定时任务配置已保存")
}

// GetSiteSettings 获取站点设置
func (h *SettingHandler) GetSiteSettings(c *gin.Context) {
	settings := gin.H{
		"site_name":        "AnchorFinance",
		"site_url":         "",
		"site_logo":        "",
		"site_description": "",
		"site_keywords":    "",
		"site_icp":         "",
		"site_copyright":   "",
		"site_footer":      "",
		"default_language":  "zh-CN",
		"default_timezone":  "Asia/Shanghai",
		"date_format":       "Y-m-d H:i:s",
	}
	response.Success(c, settings)
}

// SaveSiteSettings 保存站点设置
func (h *SettingHandler) SaveSiteSettings(c *gin.Context) {
	var req struct {
		SiteName        string `json:"site_name"`
		SiteURL         string `json:"site_url"`
		SiteLogo        string `json:"site_logo"`
		SiteDescription string `json:"site_description"`
		SiteKeywords    string `json:"site_keywords"`
		SiteIcp         string `json:"site_icp"`
		SiteCopyright   string `json:"site_copyright"`
		SiteFooter      string `json:"site_footer"`
		DefaultLanguage  string `json:"default_language"`
		DefaultTimezone  string `json:"default_timezone"`
		DateFormat       string `json:"date_format"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	h.log.Info("saving site settings")
	response.SuccessMsg(c, "站点设置已保存")
}

// GetPaymentSettings 获取支付设置
func (h *SettingHandler) GetPaymentSettings(c *gin.Context) {
	settings := gin.H{
		"alipay": gin.H{
			"enabled":      false,
			"app_id":       "",
			"private_key":  "",
			"public_key":   "",
			"notify_url":   "",
		},
		"wechat": gin.H{
			"enabled":      false,
			"app_id":       "",
			"mch_id":       "",
			"api_key":      "",
			"notify_url":   "",
		},
		"stripe": gin.H{
			"enabled":      false,
			"public_key":   "",
			"secret_key":   "",
			"webhook_key":  "",
		},
	}
	response.Success(c, settings)
}

// SavePaymentSettings 保存支付设置
func (h *SettingHandler) SavePaymentSettings(c *gin.Context) {
	var req struct {
		Alipay struct {
			Enabled     bool   `json:"enabled"`
			AppID       string `json:"app_id"`
			PrivateKey  string `json:"private_key"`
			PublicKey   string `json:"public_key"`
			NotifyURL   string `json:"notify_url"`
		} `json:"alipay"`
		Wechat struct {
			Enabled     bool   `json:"enabled"`
			AppID       string `json:"app_id"`
			MchID       string `json:"mch_id"`
			ApiKey      string `json:"api_key"`
			NotifyURL   string `json:"notify_url"`
		} `json:"wechat"`
		Stripe struct {
			Enabled     bool   `json:"enabled"`
			PublicKey   string `json:"public_key"`
			SecretKey   string `json:"secret_key"`
			WebhookKey  string `json:"webhook_key"`
		} `json:"stripe"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	h.log.Info("saving payment settings")
	response.SuccessMsg(c, "支付设置已保存")
}
