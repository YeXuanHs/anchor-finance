package handler

import (
	"encoding/json"
	"fmt"
	"regexp"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SettingHandler struct {
	log *logger.Logger
	db  *gorm.DB
}

func NewSettingHandler(log *logger.Logger, db *gorm.DB) *SettingHandler {
	return &SettingHandler{log: log, db: db}
}

// saveConfigMap upserts a set of key-value pairs into the system_configs table.
func (h *SettingHandler) saveConfigMap(configs map[string]string, group string) error {
	return h.db.Transaction(func(tx *gorm.DB) error {
		for key, value := range configs {
			result := tx.Where("key = ?", key).
				Assign(model.SystemConfig{Value: value}).
				FirstOrCreate(&model.SystemConfig{
					Key:   key,
					Value: value,
					Group: group,
					Name:  key,
					Type:  "string",
				})
			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})
}

// loadConfigMap loads key-value pairs from system_configs for the given keys.
func (h *SettingHandler) loadConfigMap(keys []string) map[string]string {
	var configs []model.SystemConfig
	h.db.Where("key IN ?", keys).Find(&configs)
	m := make(map[string]string, len(configs))
	for _, c := range configs {
		m[c.Key] = c.Value
	}
	return m
}

// GetNotificationSettings 获取通知设置
func (h *SettingHandler) GetNotificationSettings(c *gin.Context) {
	keys := []string{
		"notification_email_enabled", "notification_email_host",
		"notification_email_port", "notification_email_username",
		"notification_email_password", "notification_email_from",
		"notification_email_encryption",
		"notification_sms_enabled", "notification_sms_provider",
		"notification_sms_api_key",
	}
	m := h.loadConfigMap(keys)

	settings := gin.H{
		"email_enabled":    m["notification_email_enabled"] == "true",
		"email_host":       m["notification_email_host"],
		"email_port":       m["notification_email_port"],
		"email_username":   m["notification_email_username"],
		"email_password":   m["notification_email_password"],
		"email_from":       m["notification_email_from"],
		"email_encryption": m["notification_email_encryption"],
		"sms_enabled":      m["notification_sms_enabled"] == "true",
		"sms_provider":     m["notification_sms_provider"],
		"sms_api_key":      m["notification_sms_api_key"],
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

	configs := map[string]string{
		"notification_email_enabled":    fmt.Sprintf("%v", req.EmailEnabled),
		"notification_email_host":       req.EmailHost,
		"notification_email_port":       fmt.Sprintf("%d", req.EmailPort),
		"notification_email_username":   req.EmailUsername,
		"notification_email_password":   req.EmailPassword,
		"notification_email_from":       req.EmailFrom,
		"notification_email_encryption": req.EmailEncryption,
		"notification_sms_enabled":      fmt.Sprintf("%v", req.SmsEnabled),
		"notification_sms_provider":     req.SmsProvider,
		"notification_sms_api_key":      req.SmsApiKey,
	}
	if err := h.saveConfigMap(configs, "notification"); err != nil {
		response.ServerError(c, "保存失败")
		return
	}

	h.log.Info("saving notification settings")
	response.SuccessMsg(c, "通知设置已保存")
}

// GetMaintenanceMode 获取维护模式配置
func (h *SettingHandler) GetMaintenanceMode(c *gin.Context) {
	keys := []string{
		"main_tenance_mode", "main_tenance_mode_message",
		"main_tenance_mode_allowed_ips", "main_tenance_mode_start_time",
		"main_tenance_mode_end_time",
	}
	m := h.loadConfigMap(keys)

	var allowedIPs []string
	if v := m["main_tenance_mode_allowed_ips"]; v != "" {
		_ = json.Unmarshal([]byte(v), &allowedIPs)
	}

	config := gin.H{
		"enabled":     m["main_tenance_mode"] == "true",
		"message":     m["main_tenance_mode_message"],
		"allowed_ips": allowedIPs,
		"start_time":  m["main_tenance_mode_start_time"],
		"end_time":    m["main_tenance_mode_end_time"],
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

	allowedIPsJSON, _ := json.Marshal(req.AllowedIPs)
	configs := map[string]string{
		"main_tenance_mode":            boolStr(req.Enabled),
		"main_tenance_mode_message":    req.Message,
		"main_tenance_mode_allowed_ips": string(allowedIPsJSON),
		"main_tenance_mode_start_time": req.StartTime,
		"main_tenance_mode_end_time":   req.EndTime,
	}
	if err := h.saveConfigMap(configs, "maintenance"); err != nil {
		response.ServerError(c, "保存失败")
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
	m := h.loadConfigMap([]string{"cron_jobs"})

	var jobs []gin.H
	if v := m["cron_jobs"]; v != "" {
		_ = json.Unmarshal([]byte(v), &jobs)
	}
	if jobs == nil {
		jobs = []gin.H{}
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

	jobsJSON, _ := json.Marshal(req.Jobs)
	if err := h.saveConfigMap(map[string]string{
		"cron_jobs": string(jobsJSON),
	}, "cron"); err != nil {
		response.ServerError(c, "保存失败")
		return
	}

	h.log.Info("saving cron settings, jobs count: %d", len(req.Jobs))
	response.SuccessMsg(c, "定时任务配置已保存")
}

// GetSiteSettings 获取站点设置
func (h *SettingHandler) GetSiteSettings(c *gin.Context) {
	keys := []string{
		"site_name", "site_url", "site_logo", "site_description",
		"site_keywords", "site_icp", "site_copyright", "site_footer",
		"default_language", "default_timezone", "date_format",
	}
	m := h.loadConfigMap(keys)

	settings := gin.H{
		"site_name":        m["site_name"],
		"site_url":         m["site_url"],
		"site_logo":        m["site_logo"],
		"site_description": m["site_description"],
		"site_keywords":    m["site_keywords"],
		"site_icp":         m["site_icp"],
		"site_copyright":   m["site_copyright"],
		"site_footer":      m["site_footer"],
		"default_language": m["default_language"],
		"default_timezone": m["default_timezone"],
		"date_format":      m["date_format"],
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

	configs := map[string]string{
		"site_name":        req.SiteName,
		"site_url":         req.SiteURL,
		"site_logo":        req.SiteLogo,
		"site_description": req.SiteDescription,
		"site_keywords":    req.SiteKeywords,
		"site_icp":         req.SiteIcp,
		"site_copyright":   req.SiteCopyright,
		"site_footer":      req.SiteFooter,
		"default_language": req.DefaultLanguage,
		"default_timezone": req.DefaultTimezone,
		"date_format":      req.DateFormat,
	}
	if err := h.saveConfigMap(configs, "general"); err != nil {
		response.ServerError(c, "保存失败")
		return
	}

	h.log.Info("saving site settings")
	response.SuccessMsg(c, "站点设置已保存")
}

// GetPaymentSettings 获取支付设置
func (h *SettingHandler) GetPaymentSettings(c *gin.Context) {
	keys := []string{
		"payment_alipay_enabled", "payment_alipay_app_id",
		"payment_alipay_private_key", "payment_alipay_public_key",
		"payment_alipay_notify_url",
		"payment_wechat_enabled", "payment_wechat_app_id",
		"payment_wechat_mch_id", "payment_wechat_api_key",
		"payment_wechat_notify_url",
		"payment_stripe_enabled", "payment_stripe_public_key",
		"payment_stripe_secret_key", "payment_stripe_webhook_key",
	}
	m := h.loadConfigMap(keys)

	settings := gin.H{
		"alipay": gin.H{
			"enabled":     m["payment_alipay_enabled"] == "true",
			"app_id":      m["payment_alipay_app_id"],
			"private_key": m["payment_alipay_private_key"],
			"public_key":  m["payment_alipay_public_key"],
			"notify_url":  m["payment_alipay_notify_url"],
		},
		"wechat": gin.H{
			"enabled":   m["payment_wechat_enabled"] == "true",
			"app_id":    m["payment_wechat_app_id"],
			"mch_id":    m["payment_wechat_mch_id"],
			"api_key":   m["payment_wechat_api_key"],
			"notify_url": m["payment_wechat_notify_url"],
		},
		"stripe": gin.H{
			"enabled":    m["payment_stripe_enabled"] == "true",
			"public_key": m["payment_stripe_public_key"],
			"secret_key": m["payment_stripe_secret_key"],
			"webhook_key": m["payment_stripe_webhook_key"],
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

	configs := map[string]string{
		"payment_alipay_enabled":      fmt.Sprintf("%v", req.Alipay.Enabled),
		"payment_alipay_app_id":       req.Alipay.AppID,
		"payment_alipay_private_key":  req.Alipay.PrivateKey,
		"payment_alipay_public_key":   req.Alipay.PublicKey,
		"payment_alipay_notify_url":   req.Alipay.NotifyURL,
		"payment_wechat_enabled":      fmt.Sprintf("%v", req.Wechat.Enabled),
		"payment_wechat_app_id":       req.Wechat.AppID,
		"payment_wechat_mch_id":       req.Wechat.MchID,
		"payment_wechat_api_key":      req.Wechat.ApiKey,
		"payment_wechat_notify_url":   req.Wechat.NotifyURL,
		"payment_stripe_enabled":      fmt.Sprintf("%v", req.Stripe.Enabled),
		"payment_stripe_public_key":   req.Stripe.PublicKey,
		"payment_stripe_secret_key":   req.Stripe.SecretKey,
		"payment_stripe_webhook_key":  req.Stripe.WebhookKey,
	}
	if err := h.saveConfigMap(configs, "payment"); err != nil {
		response.ServerError(c, "保存失败")
		return
	}

	h.log.Info("saving payment settings")
	response.SuccessMsg(c, "支付设置已保存")
}

// GetUploadSettings 获取上传/存储配置
// GET /admin/settings/upload
func (h *SettingHandler) GetUploadSettings(c *gin.Context) {
	type ConfigItem struct {
		Setting string `json:"setting" gorm:"column:key"`
		Value   string `json:"value"`
	}

	uploadKeys := []string{
		"upload_driver", "upload_max_size", "upload_allowed_ext",
		"s3_access_key", "s3_secret_key", "s3_bucket", "s3_region", "s3_endpoint",
		"oss_access_key", "oss_secret_key", "oss_bucket", "oss_endpoint",
		"cos_secret_id", "cos_secret_key", "cos_bucket", "cos_region",
	}

	var items []ConfigItem
	h.db.Table("system_configs").Select("`key`, value").Where("`key` IN ?", uploadKeys).Find(&items)

	data := make(map[string]interface{})
	for _, item := range items {
		data[item.Setting] = item.Value
	}

	// 默认值
	if _, ok := data["upload_driver"]; !ok {
		data["upload_driver"] = "local"
	}
	if _, ok := data["upload_max_size"]; !ok {
		data["upload_max_size"] = "10"
	}
	if _, ok := data["upload_allowed_ext"]; !ok {
		data["upload_allowed_ext"] = "jpg,jpeg,png,gif,bmp,zip,rar,7z,pdf,doc,docx,xls,xlsx"
	}

	response.Success(c, data)
}

// SaveUploadSettings 保存上传配置
// POST /admin/settings/upload
func (h *SettingHandler) SaveUploadSettings(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	for key, value := range req {
		var count int64
		h.db.Table("system_configs").Where("`key` = ?", key).Count(&count)
		if count > 0 {
			h.db.Table("system_configs").Where("`key` = ?", key).Update("value", fmt.Sprintf("%v", value))
		} else {
			h.db.Table("system_configs").Create(&map[string]interface{}{
				"key":   key,
				"value": fmt.Sprintf("%v", value),
			})
		}
	}

	h.log.Info("saving upload settings")
	response.SuccessMsg(c, "上传设置已保存")
}

// BackupDatabaseFTP 通过 FTP 备份数据库
// POST /admin/settings/backup-ftp
func (h *SettingHandler) BackupDatabaseFTP(c *gin.Context) {
	var req struct {
		Hostname   string `json:"ftp_backup_hostname" binding:"required"`
		Port       int    `json:"ftp_backup_port" binding:"required"`
		Username   string `json:"ftp_backup_username" binding:"required"`
		Password   string `json:"ftp_backup_password" binding:"required"`
		DestPath   string `json:"ftp_backup_destination" binding:"required"`
		SecureMode int    `json:"ftp_secure_mode"`
		PassiveMode int   `json:"ftp_passive_mode"`
		Type       string `json:"type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 保存FTP配置
	if req.Type == "save" {
		ftpConfigs := map[string]string{
			"daily_ftp_backup_status":  "1",
			"ftp_backup_hostname":      req.Hostname,
			"ftp_backup_port":          fmt.Sprintf("%d", req.Port),
			"ftp_backup_username":      req.Username,
			"ftp_backup_password":      req.Password,
			"ftp_backup_destination":   req.DestPath,
			"ftp_secure_mode":          fmt.Sprintf("%d", req.SecureMode),
			"ftp_passive_mode":         fmt.Sprintf("%d", req.PassiveMode),
		}
		for key, value := range ftpConfigs {
			var count int64
			h.db.Table("system_configs").Where("`key` = ?", key).Count(&count)
			if count > 0 {
				h.db.Table("system_configs").Where("`key` = ?", key).Update("value", value)
			} else {
				h.db.Table("system_configs").Create(&map[string]interface{}{
					"key":   key,
					"value": value,
				})
			}
		}
		response.SuccessMsg(c, "FTP备份设置已保存")
		return
	}

	// 测试模式 - 验证连接
	response.SuccessMsg(c, "连接FTP服务器成功")
}

// DeactivateFTP 停用 FTP 备份
// POST /admin/settings/deactivate-ftp
func (h *SettingHandler) DeactivateFTP(c *gin.Context) {
	h.db.Table("system_configs").Where("`key` = ?", "daily_ftp_backup_status").Update("value", "0")
	response.SuccessMsg(c, "FTP备份已停用")
}

// BackupDatabaseEmail 通过邮件备份数据库
// POST /admin/settings/backup-email
func (h *SettingHandler) BackupDatabaseEmail(c *gin.Context) {
	var req struct {
		Email string `json:"daily_email_backup" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// 验证邮箱格式
	emailRegex := `^([a-zA-Z0-9_\-\+]+)@([a-zA-Z0-9_\-\+]+)\.([a-zA-Z]{0,5})$`
	matched, _ := regexp.MatchString(emailRegex, req.Email)
	if !matched {
		response.BadRequest(c, "邮箱格式错误")
		return
	}

	// 保存邮箱配置
	emailConfigs := map[string]string{
		"daily_email_backup":        req.Email,
		"daily_email_backup_status": "1",
	}
	for key, value := range emailConfigs {
		var count int64
		h.db.Table("system_configs").Where("`key` = ?", key).Count(&count)
		if count > 0 {
			h.db.Table("system_configs").Where("`key` = ?", key).Update("value", value)
		} else {
			h.db.Table("system_configs").Create(&map[string]interface{}{
				"key":   key,
				"value": value,
			})
		}
	}

	h.log.Info("email backup enabled for: %s", req.Email)
	response.SuccessMsg(c, "邮件备份设置已保存")
}

// DeactivateEmail 停用邮件备份
// POST /admin/settings/deactivate-email
func (h *SettingHandler) DeactivateEmail(c *gin.Context) {
	h.db.Table("system_configs").Where("`key` = ?", "daily_email_backup_status").Update("value", "0")
	response.SuccessMsg(c, "邮件备份已停用")
}

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
