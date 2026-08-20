package sms

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// SmsbaoConfig holds configuration for the Smsbao SMS provider.
type SmsbaoConfig struct {
	User string `json:"user"` // 账号
	Pass string `json:"pass"` // 密码 (明文，内部自动MD5)
	Sign string `json:"sign"` // 短信签名
}

// SmsbaoPlugin implements SmsSender for the Smsbao (短信宝) SMS service.
type SmsbaoPlugin struct {
	config SmsbaoConfig
	client *http.Client
}

// NewSmsbaoPlugin creates a new Smsbao SMS plugin instance.
func NewSmsbaoPlugin(cfg SmsbaoConfig) *SmsbaoPlugin {
	return &SmsbaoPlugin{
		config: cfg,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// Info returns plugin metadata.
func (p *SmsbaoPlugin) Info() PluginInfo {
	return PluginInfo{
		Name:        "smsbao",
		Title:       "短信宝",
		Description: "短信宝短信服务",
		Author:      "智简魔方",
		Version:     "1.0",
		HelpURL:     "http://www.smsbao.com/",
	}
}

// smsbaoStatusMap maps Smsbao error codes to human-readable messages.
var smsbaoStatusMap = map[string]string{
	"0":                     "短信发送成功",
	"-1":                    "参数不全",
	"-2":                    "服务器空间不支持，请确认支持curl或者fsocket",
	"30":                    "密码错误",
	"40":                    "账号不存在",
	"41":                    "余额不足",
	"42":                    "帐户已过期",
	"43":                    "IP地址限制",
	"50":                    "内容含有敏感词",
	"51":                    "手机号码不正确",
}

// Send sends a plain text SMS message.
func (p *SmsbaoPlugin) Send(mobile, content string) (*SendResult, error) {
	if mobile == "" {
		return nil, ErrInvalidPhone
	}
	if content == "" {
		return nil, ErrEmptyContent
	}

	sign := formatSign(p.config.Sign)
	fullContent := sign + content

	result, err := p.sendRequest("cn", mobile, fullContent)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// SendTemplate sends an SMS using template parameters.
// Smsbao doesn't have native template support, so we replace {key} placeholders in the content.
func (p *SmsbaoPlugin) SendTemplate(mobile, templateContent string, params map[string]string) (*SendResult, error) {
	if mobile == "" {
		return nil, ErrInvalidPhone
	}

	content := templateContent
	for k, v := range params {
		content = replacePlaceholder(content, k, v)
	}

	return p.Send(mobile, content)
}

// GetTemplate returns success with status=2 (Smsbao doesn't support remote template management).
func (p *SmsbaoPlugin) GetTemplate(templateID string) (*TemplateResult, error) {
	return &TemplateResult{
		Status:         "success",
		TemplateStatus: 2,
	}, nil
}

// CreateTemplate returns success (Smsbao doesn't support remote template management).
func (p *SmsbaoPlugin) CreateTemplate(title, content string) (*TemplateResult, error) {
	return &TemplateResult{
		Status:         "success",
		TemplateStatus: 2,
	}, nil
}

// UpdateTemplate returns success (Smsbao doesn't support remote template management).
func (p *SmsbaoPlugin) UpdateTemplate(templateID, title, content string) (*TemplateResult, error) {
	return &TemplateResult{
		Status:         "success",
		TemplateStatus: 2,
	}, nil
}

// DeleteTemplate returns success (Smsbao doesn't support remote template management).
func (p *SmsbaoPlugin) DeleteTemplate(templateID string) (*TemplateResult, error) {
	return &TemplateResult{
		Status: "success",
	}, nil
}

// sendRequest makes an HTTP request to the Smsbao API.
func (p *SmsbaoPlugin) sendRequest(smsType, mobile, content string) (*SendResult, error) {
	var apiURL string
	if smsType == "global" {
		apiURL = "http://api.smsbao.com/wsms"
	} else {
		apiURL = "http://api.smsbao.com/sms"
	}

	passMD5 := fmt.Sprintf("%x", md5.Sum([]byte(p.config.Pass)))

	params := url.Values{}
	params.Set("u", p.config.User)
	params.Set("p", passMD5)
	params.Set("m", mobile)
	params.Set("c", content)

	fullURL := apiURL + "?" + params.Encode()

	resp, err := p.client.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("%w: http request failed: %v", ErrSendFailed, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: read response failed: %v", ErrSendFailed, err)
	}

	code := string(body)

	if code == "0" {
		return &SendResult{
			Status:  "success",
			Content: content,
		}, nil
	}

	errMsg := smsbaoStatusMap[code]
	if errMsg == "" {
		errMsg = "未知错误"
	}

	return &SendResult{
		Status:  "error",
		Content: content,
		Msg:     fmt.Sprintf("%s (code: %s)", errMsg, code),
	}, fmt.Errorf("%w: %s (code: %s)", ErrSendFailed, errMsg, code)
}

// replacePlaceholder replaces {key} placeholders in content.
func replacePlaceholder(content, key, value string) string {
	// Simple replacement of {key} with value
	result := ""
	search := "{" + key + "}"
	for {
		idx := indexOf(content, search)
		if idx < 0 {
			result += content
			break
		}
		result += content[:idx] + value
		content = content[idx+len(search):]
	}
	return result
}

// indexOf returns the index of substr in s, or -1 if not found.
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// Ensure interface compliance at compile time.
var _ SmsSender = (*SubmailPlugin)(nil)
var _ SmsSender = (*SmsbaoPlugin)(nil)

// init registers the plugins with the factory.
func init() {
	RegisterPlugin("submail", func(config map[string]interface{}) (SmsSender, error) {
		cfg := SubmailConfig{}
		cfgBytes, _ := json.Marshal(config)
		if err := json.Unmarshal(cfgBytes, &cfg); err != nil {
			return nil, fmt.Errorf("parse submail config: %w", err)
		}
		if cfg.AppID == "" || cfg.AppKey == "" {
			return nil, fmt.Errorf("%w: submail requires app_id and app_key", ErrMissingConfig)
		}
		return NewSubmailPlugin(cfg), nil
	})

	RegisterPlugin("smsbao", func(config map[string]interface{}) (SmsSender, error) {
		cfg := SmsbaoConfig{}
		cfgBytes, _ := json.Marshal(config)
		if err := json.Unmarshal(cfgBytes, &cfg); err != nil {
			return nil, fmt.Errorf("parse smsbao config: %w", err)
		}
		if cfg.User == "" || cfg.Pass == "" {
			return nil, fmt.Errorf("%w: smsbao requires user and pass", ErrMissingConfig)
		}
		return NewSmsbaoPlugin(cfg), nil
	})
}
