package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/internal/plugin/mail"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// NotifyService AI工单通知服务 - 支持多渠道
type NotifyService struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewNotifyService 创建通知服务
func NewNotifyService(db *gorm.DB, log *logger.Logger) *NotifyService {
	return &NotifyService{db: db, log: log}
}

// NotifyConfig 通知配置
type NotifyConfig struct {
	Enabled         bool   `json:"enabled"`
	EmailEnabled    bool   `json:"email_enabled"`
	NotifyEmail     string `json:"notify_email"`
	SMTPHost        string `json:"smtp_host"`
	SMTPPort        int    `json:"smtp_port"`
	SMTPUser        string `json:"smtp_user"`
	SMTPPass        string `json:"smtp_pass"`
	SMTPFrom        string `json:"smtp_from"`
	WechatEnabled   bool   `json:"wechat_enabled"`
	WechatWebhook   string `json:"wechat_webhook"`
	DingtalkEnabled bool   `json:"dingtalk_enabled"`
	DingtalkWebhook string `json:"dingtalk_webhook"`
	FeishuEnabled   bool   `json:"feishu_enabled"`
	FeishuWebhook   string `json:"feishu_webhook"`
	TelegramEnabled bool   `json:"telegram_enabled"`
	TelegramToken   string `json:"telegram_bot_token"`
	TelegramChatID  string `json:"telegram_chat_id"`
	PushplusEnabled bool   `json:"pushplus_enabled"`
	PushplusToken   string `json:"pushplus_token"`
	WebhookEnabled  bool   `json:"webhook_enabled"`
	WebhookURL      string `json:"webhook_url"`
	NotifyCooldown  int    `json:"notify_cooldown"` // 冷却时间（分钟）
	TicketURLTpl    string `json:"ticket_url_template"`
	SiteAdminURL    string `json:"site_admin_url"`
}

// SendHumanHandoffNotify 发送转人工通知
func (s *NotifyService) SendHumanHandoffNotify(ticketID uint, reason string) error {
	// 读取配置
	cfg := s.loadConfig()
	if !cfg.Enabled {
		return nil
	}

	// 检查冷却时间
	if !s.checkCooldown(ticketID, cfg.NotifyCooldown) {
		s.log.Debugf("通知冷却中，跳过: ticket_id=%d", ticketID)
		return nil
	}

	// 获取工单信息
	ctxSvc := NewTicketContextService(s.db, s.log)
	ticketInfo := ctxSvc.BuildTicketInfoForNotify(ticketID)
	if ticketInfo == nil {
		return fmt.Errorf("工单不存在: %d", ticketID)
	}

	subject := fmt.Sprintf("【AI工单提醒】有工单需要人工处理 #%s", ticketInfo["tid"])
	textContent := s.buildTextContent(ticketInfo, reason, cfg)

	anySuccess := false

	// 邮件
	if cfg.EmailEnabled && cfg.NotifyEmail != "" {
		sent := s.sendEmail(cfg, subject, textContent)
		s.logNotification(ticketID, ticketInfo["tid"].(string), "mail", subject, sent)
		if sent {
			anySuccess = true
		}
	}

	// 企业微信
	if cfg.WechatEnabled && cfg.WechatWebhook != "" {
		sent := s.sendWechat(cfg.WechatWebhook, subject, textContent)
		s.logNotification(ticketID, ticketInfo["tid"].(string), "wechat", subject, sent)
		if sent {
			anySuccess = true
		}
	}

	// 钉钉
	if cfg.DingtalkEnabled && cfg.DingtalkWebhook != "" {
		sent := s.sendDingtalk(cfg.DingtalkWebhook, subject, textContent)
		s.logNotification(ticketID, ticketInfo["tid"].(string), "dingtalk", subject, sent)
		if sent {
			anySuccess = true
		}
	}

	// 飞书
	if cfg.FeishuEnabled && cfg.FeishuWebhook != "" {
		sent := s.sendFeishu(cfg.FeishuWebhook, subject, textContent)
		s.logNotification(ticketID, ticketInfo["tid"].(string), "feishu", subject, sent)
		if sent {
			anySuccess = true
		}
	}

	// Telegram
	if cfg.TelegramEnabled && cfg.TelegramToken != "" && cfg.TelegramChatID != "" {
		sent := s.sendTelegram(cfg, subject, textContent)
		s.logNotification(ticketID, ticketInfo["tid"].(string), "telegram", subject, sent)
		if sent {
			anySuccess = true
		}
	}

	// PushPlus
	if cfg.PushplusEnabled && cfg.PushplusToken != "" {
		sent := s.sendPushplus(cfg.PushplusToken, subject, textContent)
		s.logNotification(ticketID, ticketInfo["tid"].(string), "pushplus", subject, sent)
		if sent {
			anySuccess = true
		}
	}

	// 自定义Webhook
	if cfg.WebhookEnabled && cfg.WebhookURL != "" {
		sent := s.sendCustomWebhook(cfg.WebhookURL, ticketInfo, reason, subject, textContent)
		s.logNotification(ticketID, ticketInfo["tid"].(string), "webhook", subject, sent)
		if sent {
			anySuccess = true
		}
	}

	if !anySuccess {
		s.log.Warnf("所有通知渠道均失败: ticket_id=%d", ticketID)
	}

	return nil
}

// loadConfig 从数据库加载通知配置
func (s *NotifyService) loadConfig() NotifyConfig {
	cfg := NotifyConfig{
		NotifyCooldown: 60, // 默认60分钟冷却
	}

	var configs []model.AITicketConfig
	s.db.Find(&configs)

	configMap := make(map[string]string)
	for _, c := range configs {
		configMap[c.Key] = c.Value
	}

	cfg.Enabled = configMap["notify_enabled"] == "1"
	cfg.EmailEnabled = configMap["email_enabled"] == "1"
	cfg.NotifyEmail = configMap["notify_email"]
	cfg.SMTPHost = configMap["smtp_host"]
	fmt.Sscanf(configMap["smtp_port"], "%d", &cfg.SMTPPort)
	if cfg.SMTPPort == 0 {
		cfg.SMTPPort = 465
	}
	cfg.SMTPUser = configMap["smtp_user"]
	cfg.SMTPPass = configMap["smtp_pass"]
	cfg.SMTPFrom = configMap["smtp_from"]
	cfg.WechatEnabled = configMap["wechat_enabled"] == "1"
	cfg.WechatWebhook = configMap["wechat_webhook"]
	cfg.DingtalkEnabled = configMap["dingtalk_enabled"] == "1"
	cfg.DingtalkWebhook = configMap["dingtalk_webhook"]
	cfg.FeishuEnabled = configMap["feishu_enabled"] == "1"
	cfg.FeishuWebhook = configMap["feishu_webhook"]
	cfg.TelegramEnabled = configMap["telegram_enabled"] == "1"
	cfg.TelegramToken = configMap["telegram_bot_token"]
	cfg.TelegramChatID = configMap["telegram_chat_id"]
	cfg.PushplusEnabled = configMap["pushplus_enabled"] == "1"
	cfg.PushplusToken = configMap["pushplus_token"]
	cfg.WebhookEnabled = configMap["webhook_enabled"] == "1"
	cfg.WebhookURL = configMap["webhook_url"]
	fmt.Sscanf(configMap["notify_cooldown"], "%d", &cfg.NotifyCooldown)
	if cfg.NotifyCooldown <= 0 {
		cfg.NotifyCooldown = 60
	}
	cfg.TicketURLTpl = configMap["ticket_url_template"]
	cfg.SiteAdminURL = configMap["site_admin_url"]

	return cfg
}

// buildTextContent 构建通知文本内容
func (s *NotifyService) buildTextContent(ticketInfo map[string]interface{}, reason string, cfg NotifyConfig) string {
	priority, _ := ticketInfo["priority"].(string)
	tid, _ := ticketInfo["tid"].(string)
	title, _ := ticketInfo["title"].(string)
	department, _ := ticketInfo["department"].(string)
	content, _ := ticketInfo["content"].(string)
	ticketTime, _ := ticketInfo["time"].(string)
	ticketID, _ := ticketInfo["ticket_id"].(uint)

	if priority == "" {
		priority = "中"
	}
	if department == "" {
		department = "未分配"
	}

	text := "有一个工单需要人工处理，请尽快上线处理。\n\n"
	text += "工单号：" + tid + "\n"
	text += fmt.Sprintf("工单ID：%d\n", ticketID)
	text += "标题：" + title + "\n"
	text += "部门：" + department + "\n"
	text += "优先级：" + priority + "\n"
	text += "转人工原因：" + reason + "\n\n"
	text += "工单内容：" + content + "\n\n"
	text += "发送时间：" + ticketTime + "\n"

	ticketURL := s.getTicketURL(ticketID, cfg)
	if ticketURL != "" {
		text += "工单链接：" + ticketURL + "\n"
	}

	return text
}

// getTicketURL 获取工单链接
func (s *NotifyService) getTicketURL(ticketID uint, cfg NotifyConfig) string {
	template := strings.TrimSpace(cfg.TicketURLTpl)
	if template == "" {
		siteURL := strings.TrimRight(cfg.SiteAdminURL, "/")
		if siteURL == "" {
			siteURL = "/admin"
		}
		template = siteURL + "/ticket/detail/id/{ticket_id}"
	}
	return strings.ReplaceAll(template, "{ticket_id}", fmt.Sprintf("%d", ticketID))
}

// checkCooldown 检查通知冷却时间
func (s *NotifyService) checkCooldown(ticketID uint, cooldownMinutes int) bool {
	if cooldownMinutes <= 0 {
		cooldownMinutes = 60
	}

	var count int64
	s.db.Model(&model.AITicketNotifyLog{}).
		Where("ticket_id = ? AND created_at > ?", ticketID, time.Now().Add(-time.Duration(cooldownMinutes)*time.Minute)).
		Count(&count)

	return count == 0
}

// logNotification 记录通知日志
func (s *NotifyService) logNotification(ticketID uint, tid, channel, content string, success bool) {
	status := "failed"
	if success {
		status = "success"
	}

	if len(content) > 500 {
		content = content[:500]
	}

	s.db.Create(&model.AITicketNotifyLog{
		TicketID: ticketID,
		Channel:  channel,
		Status:   status,
		Content:  content,
	})
}

// ─── 渠道发送方法 ───

// sendEmail 发送邮件通知
func (s *NotifyService) sendEmail(cfg NotifyConfig, subject, content string) bool {
	smtpCfg := mail.SmtpConfig{
		Host:        cfg.SMTPHost,
		Port:        cfg.SMTPPort,
		Username:    cfg.SMTPUser,
		Password:    cfg.SMTPPass,
		FromName:    "",
		SystemEmail: cfg.SMTPFrom,
		SmtpSecure:  "ssl",
	}
	sender := mail.NewSmtpSender(smtpCfg)
	if err := sender.Send(cfg.NotifyEmail, subject, content, nil); err != nil {
		s.log.Errorf("邮件发送失败: %v", err)
		return false
	}
	return true
}

// sendWechat 发送企业微信 Webhook
func (s *NotifyService) sendWechat(webhook, subject, content string) bool {
	data := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": subject + "\n\n" + content,
		},
	}
	return s.httpPost(webhook, data)
}

// sendDingtalk 发送钉钉 Webhook
func (s *NotifyService) sendDingtalk(webhook, subject, content string) bool {
	data := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": subject + "\n\n" + content,
		},
	}
	return s.httpPost(webhook, data)
}

// sendFeishu 发送飞书 Webhook
func (s *NotifyService) sendFeishu(webhook, subject, content string) bool {
	data := map[string]interface{}{
		"msg_type": "text",
		"content": map[string]string{
			"text": subject + "\n\n" + content,
		},
	}
	return s.httpPost(webhook, data)
}

// sendTelegram 发送 Telegram Bot 消息
func (s *NotifyService) sendTelegram(cfg NotifyConfig, subject, content string) bool {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", cfg.TelegramToken)
	data := map[string]interface{}{
		"chat_id":    cfg.TelegramChatID,
		"text":       subject + "\n\n" + content,
		"parse_mode": "HTML",
	}
	return s.httpPost(url, data)
}

// sendPushplus 发送 PushPlus 消息
func (s *NotifyService) sendPushplus(token, subject, content string) bool {
	url := "http://www.pushplus.plus/send"
	data := map[string]interface{}{
		"token":    token,
		"title":    subject,
		"content":  strings.ReplaceAll(content, "\n", "<br>"),
		"template": "txt",
	}
	return s.httpPost(url, data)
}

// sendCustomWebhook 发送自定义 Webhook
func (s *NotifyService) sendCustomWebhook(url string, ticketInfo map[string]interface{}, reason, subject, content string) bool {
	ticketID, _ := ticketInfo["ticket_id"].(uint)
	tid, _ := ticketInfo["tid"].(string)
	title, _ := ticketInfo["title"].(string)
	department, _ := ticketInfo["department"].(string)
	priority, _ := ticketInfo["priority"].(string)
	ticketContent, _ := ticketInfo["content"].(string)

	data := map[string]interface{}{
		"event":      "ai_handoff",
		"ticket_id":  ticketID,
		"tid":        tid,
		"title":      title,
		"department": department,
		"priority":   priority,
		"reason":     reason,
		"content":    ticketContent,
		"subject":    subject,
		"text":       content,
		"timestamp":  time.Now().Unix(),
	}
	return s.httpPost(url, data)
}

// httpPost 发送 HTTP POST 请求
func (s *NotifyService) httpPost(url string, data interface{}) bool {
	jsonData, err := json.Marshal(data)
	if err != nil {
		s.log.Errorf("JSON序列化失败: %v", err)
		return false
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Post(url, "application/json; charset=utf-8", bytes.NewBuffer(jsonData))
	if err != nil {
		s.log.Errorf("HTTP请求失败: %v", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

// SendTestNotify 发送测试通知
func (s *NotifyService) SendTestNotify() map[string]interface{} {
	cfg := s.loadConfig()

	var results []string
	anySuccess := false

	testSubject := "【AI工单测试】通知测试"
	testContent := "这是一条测试消息。\n发送时间：" + time.Now().Format("2006-01-02 15:04:05")

	if cfg.EmailEnabled && cfg.NotifyEmail != "" {
		sent := s.sendEmail(cfg, testSubject, testContent)
		results = append(results, "邮件: "+boolStatus(sent))
		if sent {
			anySuccess = true
		}
	}

	if cfg.WechatEnabled && cfg.WechatWebhook != "" {
		sent := s.sendWechat(cfg.WechatWebhook, testSubject, testContent)
		results = append(results, "企业微信: "+boolStatus(sent))
		if sent {
			anySuccess = true
		}
	}

	if cfg.DingtalkEnabled && cfg.DingtalkWebhook != "" {
		sent := s.sendDingtalk(cfg.DingtalkWebhook, testSubject, testContent)
		results = append(results, "钉钉: "+boolStatus(sent))
		if sent {
			anySuccess = true
		}
	}

	if cfg.FeishuEnabled && cfg.FeishuWebhook != "" {
		sent := s.sendFeishu(cfg.FeishuWebhook, testSubject, testContent)
		results = append(results, "飞书: "+boolStatus(sent))
		if sent {
			anySuccess = true
		}
	}

	if cfg.TelegramEnabled && cfg.TelegramToken != "" {
		sent := s.sendTelegram(cfg, testSubject, testContent)
		results = append(results, "Telegram: "+boolStatus(sent))
		if sent {
			anySuccess = true
		}
	}

	if cfg.PushplusEnabled && cfg.PushplusToken != "" {
		sent := s.sendPushplus(cfg.PushplusToken, testSubject, testContent)
		results = append(results, "PushPlus: "+boolStatus(sent))
		if sent {
			anySuccess = true
		}
	}

	if cfg.WebhookEnabled && cfg.WebhookURL != "" {
		sent := s.httpPost(cfg.WebhookURL, map[string]interface{}{
			"event": "test",
			"text":  testSubject,
			"time":  time.Now().Format("2006-01-02 15:04:05"),
		})
		results = append(results, "Webhook: "+boolStatus(sent))
		if sent {
			anySuccess = true
		}
	}

	if len(results) == 0 {
		return map[string]interface{}{
			"success": false,
			"msg":     "请先配置至少一个通知渠道",
		}
	}

	return map[string]interface{}{
		"success": anySuccess,
		"msg":     strings.Join(results, "；"),
	}
}

func boolStatus(b bool) string {
	if b {
		return "成功"
	}
	return "失败"
}
