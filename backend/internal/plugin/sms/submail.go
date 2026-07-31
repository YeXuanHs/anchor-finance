package sms

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const submailBaseURL = "http://api.mysubmail.com/"

// SubmailConfig holds configuration for the Submail SMS provider.
type SubmailConfig struct {
	AppID              string `json:"app_id"`
	AppKey             string `json:"app_key"`
	AppSign            string `json:"app_sign"`
	InternationalAppID string `json:"international_app_id,omitempty"`
	InternationalAppKey string `json:"international_app_key,omitempty"`
	InternationalAppSign string `json:"international_app_sign,omitempty"`
}

// SubmailPlugin implements SmsSender for the Submail (赛邮) SMS service.
type SubmailPlugin struct {
	config SubmailConfig
	client *http.Client
}

// NewSubmailPlugin creates a new Submail SMS plugin instance.
func NewSubmailPlugin(cfg SubmailConfig) *SubmailPlugin {
	return &SubmailPlugin{
		config: cfg,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// Info returns plugin metadata.
func (p *SubmailPlugin) Info() PluginInfo {
	return PluginInfo{
		Name:        "submail",
		Title:       "赛邮",
		Description: "赛邮短信服务",
		Author:      "智简魔方",
		Version:     "1.0",
		HelpURL:     "https://www.mysubmail.com/",
	}
}

// Send sends a plain text SMS message.
func (p *SubmailPlugin) Send(mobile, content string) (*SendResult, error) {
	if mobile == "" {
		return nil, ErrInvalidPhone
	}
	if content == "" {
		return nil, ErrEmptyContent
	}

	sign := formatSign(p.config.AppSign)
	params := map[string]string{
		"to":      mobile,
		"content": sign + content,
	}

	resp, err := p.apiRequest("cn", "message/send.json", params, http.MethodPost)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSendFailed, err)
	}

	if resp["status"] == "success" {
		return &SendResult{Status: "success", Content: content}, nil
	}

	return &SendResult{
		Status: "error",
		Content: content,
		Msg:    fmt.Sprintf("%v", resp["msg"]),
	}, fmt.Errorf("%w: %v", ErrSendFailed, resp["msg"])
}

// SendTemplate sends an SMS using a Submail template with variables.
func (p *SubmailPlugin) SendTemplate(mobile, templateID string, vars map[string]string) (*SendResult, error) {
	if mobile == "" {
		return nil, ErrInvalidPhone
	}
	if templateID == "" {
		return nil, ErrTemplateNotFound
	}

	params := map[string]string{
		"to":          mobile,
		"template_id": templateID,
	}

	// Submail uses vars as part of the request
	for k, v := range vars {
		params[k] = v
	}

	resp, err := p.apiRequest("cn", "message/xsend.json", params, http.MethodPost)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSendFailed, err)
	}

	if resp["status"] == "success" {
		return &SendResult{Status: "success"}, nil
	}

	return &SendResult{
		Status: "error",
		Msg:    fmt.Sprintf("%v", resp["msg"]),
	}, fmt.Errorf("%w: %v", ErrSendFailed, resp["msg"])
}

// GetTemplate queries the status of a template.
func (p *SubmailPlugin) GetTemplate(templateID string) (*TemplateResult, error) {
	if templateID == "" {
		return nil, ErrTemplateNotFound
	}

	params := map[string]string{
		"template_id": templateID,
	}

	resp, err := p.apiRequest("cn", "message/template.json", params, http.MethodGet)
	if err != nil {
		return nil, err
	}

	result := &TemplateResult{}
	if resp["status"] == "success" {
		result.Status = "success"
		if tpl, ok := resp["template"].(map[string]interface{}); ok {
			if id, ok := tpl["template_id"].(string); ok {
				result.TemplateID = id
			}
			if status, ok := tpl["template_status"].(float64); ok {
				result.TemplateStatus = int(status)
			}
		}
	} else {
		result.Status = "error"
		result.Msg = fmt.Sprintf("%v", resp["msg"])
	}

	return result, nil
}

// CreateTemplate creates a new template on Submail.
func (p *SubmailPlugin) CreateTemplate(title, content string) (*TemplateResult, error) {
	if content == "" {
		return nil, ErrEmptyContent
	}

	params := map[string]string{
		"sms_title":     title,
		"sms_signature": p.config.AppSign,
		"sms_content":   content,
	}

	resp, err := p.apiRequest("cn", "message/template.json", params, http.MethodPost)
	if err != nil {
		return nil, err
	}

	result := &TemplateResult{}
	if resp["status"] == "success" {
		result.Status = "success"
		if id, ok := resp["template_id"].(string); ok {
			result.TemplateID = id
		}
		result.TemplateStatus = 1
	} else {
		result.Status = "error"
		result.Msg = fmt.Sprintf("%v", resp["msg"])
	}

	return result, nil
}

// UpdateTemplate updates an existing template on Submail.
func (p *SubmailPlugin) UpdateTemplate(templateID, title, content string) (*TemplateResult, error) {
	if templateID == "" {
		return nil, ErrTemplateNotFound
	}

	params := map[string]string{
		"template_id":   templateID,
		"sms_signature": p.config.AppSign,
		"sms_content":   content,
	}
	if title != "" {
		params["sms_title"] = title
	}

	resp, err := p.apiRequest("cn", "message/template.json", params, http.MethodPut)
	if err != nil {
		return nil, err
	}

	result := &TemplateResult{}
	if resp["status"] == "success" {
		result.Status = "success"
		result.TemplateStatus = 1
	} else {
		result.Status = "error"
		result.Msg = fmt.Sprintf("%v", resp["msg"])
	}

	return result, nil
}

// DeleteTemplate deletes a template from Submail.
func (p *SubmailPlugin) DeleteTemplate(templateID string) (*TemplateResult, error) {
	if templateID == "" {
		return nil, ErrTemplateNotFound
	}

	params := map[string]string{
		"template_id": templateID,
	}

	resp, err := p.apiRequest("cn", "message/template.json", params, http.MethodDelete)
	if err != nil {
		return nil, err
	}

	result := &TemplateResult{}
	if resp["status"] == "success" {
		result.Status = "success"
	} else {
		result.Status = "error"
		result.Msg = fmt.Sprintf("%v", resp["msg"])
	}

	return result, nil
}

// apiRequest makes an HTTP request to the Submail API.
func (p *SubmailPlugin) apiRequest(smsType, apiPath string, data map[string]string, method string) (map[string]interface{}, error) {
	var appID, appKey string
	if smsType == "global" {
		appID = p.config.InternationalAppID
		appKey = p.config.InternationalAppKey
	} else {
		appID = p.config.AppID
		appKey = p.config.AppKey
	}

	// Build request params
	reqParams := url.Values{}
	reqParams.Set("appid", appID)
	reqParams.Set("appkey", appKey)
	reqParams.Set("timestamp", p.getTimestamp())
	reqParams.Set("signature", appKey)

	for k, v := range data {
		reqParams.Set(k, v)
	}

	apiURL := submailBaseURL + apiPath

	var resp *http.Response
	var err error

	switch strings.ToUpper(method) {
	case http.MethodGet:
		resp, err = p.client.Get(apiURL + "?" + reqParams.Encode())
	case http.MethodDelete:
		req, reqErr := http.NewRequest(http.MethodDelete, apiURL, strings.NewReader(reqParams.Encode()))
		if reqErr != nil {
			return nil, fmt.Errorf("create request failed: %w", reqErr)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp, err = p.client.Do(req)
	default:
		resp, err = p.client.PostForm(apiURL, reqParams)
	}

	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}

	// Remove BOM if present
	bodyStr := strings.TrimPrefix(string(body), "\ufeff")

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(bodyStr), &result); err != nil {
		return nil, fmt.Errorf("parse response failed: %w", err)
	}

	return result, nil
}

// getTimestamp fetches the server timestamp from Submail.
func (p *SubmailPlugin) getTimestamp() string {
	resp, err := p.client.Get(submailBaseURL + "service/timestamp.json")
	if err != nil {
		return fmt.Sprintf("%d", time.Now().Unix())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("%d", time.Now().Unix())
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Sprintf("%d", time.Now().Unix())
	}

	if ts, ok := result["timestamp"].(string); ok {
		return ts
	}

	return fmt.Sprintf("%d", time.Now().Unix())
}

// formatSign wraps the sign with 【】brackets.
func formatSign(sign string) string {
	sign = strings.ReplaceAll(sign, "【", "")
	sign = strings.ReplaceAll(sign, "】", "")
	return "【" + sign + "】"
}
