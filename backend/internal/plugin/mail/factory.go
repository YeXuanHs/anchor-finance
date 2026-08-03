package mail

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// SenderFactory creates a MailSender instance based on the plugin name.
type SenderFactory struct{}

// NewSenderFactory creates a new SenderFactory.
func NewSenderFactory() *SenderFactory {
	return &SenderFactory{}
}

// Create creates a MailSender by plugin name with the provided configuration.
// Supported names: smtp, submail, alimail.
func (f *SenderFactory) Create(name string, config map[string]interface{}) (MailSender, error) {
	switch name {
	case "smtp":
		return f.createSmtp(config)
	case "submail":
		return f.createSubmail(config)
	case "alimail":
		return f.createAlimail(config)
	default:
		return nil, fmt.Errorf("unsupported mail plugin: %s", name)
	}
}

func (f *SenderFactory) createSmtp(config map[string]interface{}) (MailSender, error) {
	cfg := SmtpConfig{
		Host:        getString(config, "host"),
		Port:        getInt(config, "port"),
		Username:    getString(config, "username"),
		Password:    getString(config, "password"),
		FromName:    getString(config, "from_name"),
		SystemEmail: getString(config, "system_email"),
		SmtpSecure:  getString(config, "smtp_secure"),
		Charset:     getString(config, "charset"),
	}
	if cfg.Host == "" {
		return nil, fmt.Errorf("smtp: host is required")
	}
	if cfg.Port == 0 {
		cfg.Port = 25
	}
	return NewSmtpSender(cfg), nil
}

func (f *SenderFactory) createSubmail(config map[string]interface{}) (MailSender, error) {
	cfg := SubmailConfig{
		AppID:     getString(config, "app_id"),
		AppKey:    getString(config, "app_key"),
		SignType:  getString(config, "sign_type"),
		FromName:  getString(config, "from_name"),
		FromEmail: getString(config, "from_email"),
	}
	if cfg.AppID == "" || cfg.AppKey == "" {
		return nil, fmt.Errorf("submail: app_id and app_key are required")
	}
	return NewSubmailSender(cfg), nil
}

func (f *SenderFactory) createAlimail(config map[string]interface{}) (MailSender, error) {
	cfg := AlimailConfig{
		AccessKeyID:     getString(config, "access_key_id"),
		AccessKeySecret: getString(config, "access_key_secret"),
		AccountName:     getString(config, "account_name"),
		FromAlias:       getString(config, "from_alias"),
	}
	if cfg.AccessKeyID == "" || cfg.AccessKeySecret == "" {
		return nil, fmt.Errorf("alimail: access_key_id and access_key_secret are required")
	}
	return NewAlimailSender(cfg), nil
}

// getString extracts a string value from a config map.
func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// getInt extracts an int value from a config map.
func getInt(m map[string]interface{}, key string) int {
	if v, ok := m[key]; ok {
		switch n := v.(type) {
		case int:
			return n
		case float64:
			return int(n)
		}
	}
	return 0
}

// SubmailConfig holds the configuration for Submail mail sender.
type SubmailConfig struct {
	AppID     string `json:"app_id"`
	AppKey    string `json:"app_key"`
	SignType  string `json:"sign_type"`
	FromName  string `json:"from_name"`
	FromEmail string `json:"from_email"`
}

// SubmailSender implements MailSender for Submail API.
type SubmailSender struct {
	config SubmailConfig
}

// NewSubmailSender creates a new SubmailSender with the given configuration.
func NewSubmailSender(config SubmailConfig) *SubmailSender {
	return &SubmailSender{config: config}
}

// Send sends an email via Submail API.
func (s *SubmailSender) Send(to, subject, content string, attachments []string) error {
	if s.config.AppID == "" || s.config.AppKey == "" {
		return fmt.Errorf("submail: app_id and app_key are required")
	}

	apiURL := "https://api.submail.cn/mail.send"
	form := url.Values{}
	form.Set("appid", s.config.AppID)
	form.Set("to", to)
	form.Set("subject", subject)
	form.Set("html", content)
	form.Set("signature", s.config.AppKey)
	if s.config.SignType != "" {
		form.Set("sign_type", s.config.SignType)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post(apiURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("submail: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("submail: invalid response: %s", string(body))
	}

	if status, ok := result["status"].(string); ok && status != "success" {
		msg, _ := result["msg"].(string)
		return fmt.Errorf("submail: %s", msg)
	}

	return nil
}

// AlimailConfig holds the configuration for Aliyun DirectMail sender.
type AlimailConfig struct {
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	AccountName     string `json:"account_name"`
	FromAlias       string `json:"from_alias"`
}

// AlimailSender implements MailSender for Aliyun DirectMail API.
type AlimailSender struct {
	config AlimailConfig
}

// NewAlimailSender creates a new AlimailSender with the given configuration.
func NewAlimailSender(config AlimailConfig) *AlimailSender {
	return &AlimailSender{config: config}
}

// Send sends an email via Aliyun DirectMail API.
func (s *AlimailSender) Send(to, subject, content string, attachments []string) error {
	if s.config.AccessKeyID == "" || s.config.AccessKeySecret == "" {
		return fmt.Errorf("alimail: access_key_id and access_key_secret are required")
	}

	accountName := s.config.AccountName
	if accountName == "" {
		return fmt.Errorf("alimail: account_name is required")
	}

	fromAlias := s.config.FromAlias

	now := time.Now().UTC().Format("2006-01-02T15:04:05Z")

	params := map[string]string{
		"Action":           "SingleSendMail",
		"AccountName":      accountName,
		"AddressType":      "1",
		"ToAddress":        to,
		"Subject":          subject,
		"HtmlBody":         content,
		"Format":           "JSON",
		"Version":          "2015-11-23",
		"AccessKeyId":      s.config.AccessKeyID,
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureVersion": "1.0",
		"SignatureNonce":   fmt.Sprintf("%d", time.Now().UnixNano()),
		"Timestamp":        now,
	}
	if fromAlias != "" {
		params["FromAlias"] = fromAlias
	}

	// Sort parameters lexicographically
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var canonicalQuery strings.Builder
	for i, k := range keys {
		if i > 0 {
			canonicalQuery.WriteString("&")
		}
		canonicalQuery.WriteString(url.QueryEscape(k))
		canonicalQuery.WriteString("=")
		canonicalQuery.WriteString(url.QueryEscape(params[k]))
	}

	stringToSign := fmt.Sprintf("POST&%%2F&%s", url.QueryEscape(canonicalQuery.String()))

	mac := hmac.New(sha1.New, []byte(s.config.AccessKeySecret+"&"))
	mac.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	params["Signature"] = signature

	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post("https://dm.aliyuncs.com/", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("alimail: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("alimail: invalid response: %s", string(body))
	}

	if resp.StatusCode >= 400 {
		errMsg, _ := result["Message"].(string)
		if errMsg == "" {
			errMsg = string(body)
		}
		return fmt.Errorf("alimail: %s", errMsg)
	}

	return nil
}
