package sms

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Config holds SMS provider configuration read from message_configs table.
type Config struct {
	Provider    string `json:"provider"`     // aliyun_sms / tencent_sms / huyi / custom
	AccessKey   string `json:"access_key"`   // AccessKey ID or API account
	AccessSecret string `json:"access_secret"` // AccessKey secret or API password
	SignName    string `json:"sign_name"`    // SMS signature
	TemplateID  string `json:"template_id"`  // SMS template code
	Endpoint    string `json:"endpoint"`     // Custom API endpoint (for custom provider)
	AppID       string `json:"app_id"`       // Tencent SMS app ID (SdkAppId)
}

// Sender sends SMS messages via configured providers.
type Sender struct {
	db *gorm.DB
}

// NewSender creates a new SMS Sender.
func NewSender(db *gorm.DB) *Sender {
	return &Sender{db: db}
}

// loadConfig reads SMS config from the message_configs table (channel = "sms").
func (s *Sender) loadConfig() (*Config, error) {
	var row struct {
		Provider     string `gorm:"column:provider"`
		Config       []byte `gorm:"column:config"`
		Signature    string `gorm:"column:signature"`
		IsEnabled    bool   `gorm:"column:is_enabled"`
		Status       int16  `gorm:"column:status"`
	}

	err := s.db.Table("message_configs").
		Select("provider, config, signature, is_enabled, status").
		Where("channel = ?", "sms").
		First(&row).Error
	if err != nil {
		return nil, fmt.Errorf("sms config not found: %w", err)
	}

	if !row.IsEnabled || row.Status != 1 {
		return nil, errors.New("sms channel is disabled")
	}

	cfg := &Config{
		Provider:  row.Provider,
		SignName:  row.Signature,
	}

	if len(row.Config) > 0 {
		if err := json.Unmarshal(row.Config, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse sms config: %w", err)
		}
	}

	// Signature from column takes precedence
	if row.Signature != "" {
		cfg.SignName = row.Signature
	}

	return cfg, nil
}

// Send sends an SMS verification code to the given phone number.
func (s *Sender) Send(phone, code string) error {
	cfg, err := s.loadConfig()
	if err != nil {
		return err
	}

	switch cfg.Provider {
	case "aliyun_sms":
		return s.sendAliyun(cfg, phone, code)
	case "tencent_sms":
		return s.sendTencent(cfg, phone, code)
	case "huyi":
		return s.sendHuyi(cfg, phone, code)
	case "custom":
		return s.sendCustom(cfg, phone, code)
	default:
		return fmt.Errorf("unsupported sms provider: %s", cfg.Provider)
	}
}

// sendAliyun sends SMS via Alibaba Cloud SMS API.
func (s *Sender) sendAliyun(cfg *Config, phone, code string) error {
	if cfg.AccessKey == "" || cfg.AccessSecret == "" {
		return errors.New("aliyun sms: access_key and access_secret are required")
	}

	params := map[string]string{
		"AccessKeyId":      cfg.AccessKey,
		"Action":           "SendSms",
		"Format":           "JSON",
		"PhoneNumbers":     phone,
		"RegionId":         "cn-hangzhou",
		"SignName":         cfg.SignName,
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureNonce":   fmt.Sprintf("%d", time.Now().UnixNano()),
		"SignatureVersion": "1.0",
		"TemplateCode":     cfg.TemplateID,
		"TemplateParam":    fmt.Sprintf(`{"code":"%s"}`, code),
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"Version":          "2017-05-25",
	}

	endpoint := cfg.Endpoint
	if endpoint == "" {
		endpoint = "https://dysmsapi.aliyuncs.com"
	}

	return s.sendHTTPRequest(endpoint, params, cfg.AccessSecret)
}

// sendTencent sends SMS via Tencent Cloud SMS API.
func (s *Sender) sendTencent(cfg *Config, phone, code string) error {
	if cfg.AccessKey == "" || cfg.AccessSecret == "" {
		return errors.New("tencent sms: access_key and access_secret are required")
	}

	endpoint := cfg.Endpoint
	if endpoint == "" {
		endpoint = "https://sms.tencentcloudapi.com"
	}

	payload := map[string]interface{}{
		"SmsSdkAppId": cfg.AppID,
		"Sign":        cfg.SignName,
		"TemplateId":  cfg.TemplateID,
		"TemplateParamSet": []string{code},
		"PhoneNumberSet":   []string{phone},
	}

	return s.sendTencentRequest(endpoint, payload, cfg.AccessKey, cfg.AccessSecret)
}

// sendHuyi sends SMS via 互亿无线 API.
func (s *Sender) sendHuyi(cfg *Config, phone, code string) error {
	if cfg.AccessKey == "" || cfg.AccessSecret == "" {
		return errors.New("huyi sms: account and password are required")
	}

	endpoint := cfg.Endpoint
	if endpoint == "" {
		endpoint = "http://106.ihuyi.com/webservice/sms.php?method=Submit"
	}

	content := fmt.Sprintf("您的验证码是：%s。请在5分钟内完成验证。", code)

	data := url.Values{}
	data.Set("account", cfg.AccessKey)
	data.Set("password", cfg.AccessSecret)
	data.Set("mobile", phone)
	data.Set("content", content)

	resp, err := http.PostForm(endpoint, data)
	if err != nil {
		return fmt.Errorf("huyi sms request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("huyi sms read response failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("huyi sms error: status=%d body=%s", resp.StatusCode, string(body))
	}

	// Huyi returns XML-like response, check for success
	respStr := string(body)
	if strings.Contains(respStr, "<code>2</code>") {
		return nil
	}

	return fmt.Errorf("huyi sms failed: %s", respStr)
}

// sendCustom sends SMS via a custom HTTP API.
func (s *Sender) sendCustom(cfg *Config, phone, code string) error {
	if cfg.Endpoint == "" {
		return errors.New("custom sms: endpoint is required")
	}

	payload := map[string]string{
		"phone":      phone,
		"code":       code,
		"sign_name":  cfg.SignName,
		"template":   cfg.TemplateID,
		"access_key": cfg.AccessKey,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("custom sms marshal failed: %w", err)
	}

	resp, err := http.Post(cfg.Endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("custom sms request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("custom sms error: status=%d body=%s", resp.StatusCode, string(body))
	}

	return nil
}

// sendHTTPRequest sends an Alibaba Cloud style signed HTTP request using HMAC-SHA1.
func (s *Sender) sendHTTPRequest(endpoint string, params map[string]string, accessSecret string) error {
	// Sort parameters
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build canonical query string
	var canonicalQuery strings.Builder
	for i, k := range keys {
		if i > 0 {
			canonicalQuery.WriteString("&")
		}
		canonicalQuery.WriteString(url.QueryEscape(k))
		canonicalQuery.WriteString("=")
		canonicalQuery.WriteString(url.QueryEscape(params[k]))
	}

	// Build string to sign (Alibaba Cloud uses HMAC-SHA1)
	stringToSign := "GET&" + url.QueryEscape("/") + "&" + url.QueryEscape(canonicalQuery.String())

	// Calculate HMAC-SHA1 signature
	mac := hmac.New(sha1.New, []byte(accessSecret+"&"))
	mac.Write([]byte(stringToSign))
	signature := hex.EncodeToString(mac.Sum(nil))

	// Build final URL
	finalURL := endpoint + "?" + canonicalQuery.String() + "&Signature=" + url.QueryEscape(signature)

	resp, err := http.Get(finalURL)
	if err != nil {
		return fmt.Errorf("aliyun sms request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("aliyun sms read response failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("aliyun sms error: status=%d body=%s", resp.StatusCode, string(body))
	}

	// Parse response to check for errors
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err == nil {
		if code, ok := result["Code"].(string); ok && code != "OK" {
			return fmt.Errorf("aliyun sms error: %s - %s", code, result["Message"])
		}
	}

	return nil
}

// sendTencentRequest sends a Tencent Cloud style signed HTTP request.
func (s *Sender) sendTencentRequest(endpoint string, payload map[string]interface{}, secretID, secretKey string) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// For simplicity, use TC3-HMAC-SHA256 signing
	now := time.Now().Unix()
	date := time.Unix(now, 0).UTC().Format("2006-01-02")

	// Build canonical request
	httpRequestMethod := "POST"
	canonicalURI := "/"
	canonicalQueryString := ""
	canonicalHeaders := fmt.Sprintf("content-type:application/json\nhost:sms.tencentcloudapi.com\nx-tc-action:sendsms\n")
	signedHeaders := "content-type;host;x-tc-action"
	hashedPayload := sha256Hex(string(jsonData))
	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		httpRequestMethod, canonicalURI, canonicalQueryString,
		canonicalHeaders, signedHeaders, hashedPayload)

	// Build string to sign
	algorithm := "TC3-HMAC-SHA256"
	credentialScope := fmt.Sprintf("%s/sms/tc3_request", date)
	hashedCanonicalRequest := sha256Hex(canonicalRequest)
	stringToSign := fmt.Sprintf("%s\n%d\n%s\n%s", algorithm, now, credentialScope, hashedCanonicalRequest)

	// Calculate signature
	secretDate := hmacSHA256("TC3"+secretKey, date)
	secretService := hmacSHA256(secretDate, "sms")
	secretSigning := hmacSHA256(secretService, "tc3_request")
	signature := hex.EncodeToString(hmacSHA256(secretSigning, stringToSign))

	// Build authorization header
	authorization := fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		algorithm, secretID, credentialScope, signedHeaders, signature)

	// Make request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", "sms.tencentcloudapi.com")
	req.Header.Set("X-TC-Action", "SendSms")
	req.Header.Set("X-TC-Version", "2021-01-11")
	req.Header.Set("X-TC-Timestamp", fmt.Sprintf("%d", now))
	req.Header.Set("X-TC-Region", "ap-guangzhou")
	req.Header.Set("Authorization", authorization)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("tencent sms request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("tencent sms read response failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("tencent sms error: status=%d body=%s", resp.StatusCode, string(body))
	}

	// Parse response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err == nil {
		if response, ok := result["Response"].(map[string]interface{}); ok {
			if errObj, ok := response["Error"].(map[string]interface{}); ok {
				return fmt.Errorf("tencent sms error: %v - %v", errObj["Code"], errObj["Message"])
			}
		}
	}

	return nil
}

// sha256Hex returns hex-encoded SHA256 hash.
func sha256Hex(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

// hmacSHA256 returns HMAC-SHA256 bytes.
func hmacSHA256(key interface{}, data string) []byte {
	var keyBytes []byte
	switch v := key.(type) {
	case string:
		keyBytes = []byte(v)
	case []byte:
		keyBytes = v
	default:
		keyBytes = []byte(fmt.Sprintf("%v", key))
	}
	h := hmac.New(sha256.New, keyBytes)
	h.Write([]byte(data))
	return h.Sum(nil)
}
