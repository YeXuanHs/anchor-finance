package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// BaseOAuthProvider 基础OAuth提供商
type BaseOAuthProvider struct {
	name         string
	title        string
	config       map[string]interface{}
	client       *http.Client
	authURL      string
	tokenURL     string
	userInfoURL  string
	scope        string
}

// NewBaseOAuthProvider 创建基础提供商
func NewBaseOAuthProvider(name, title string) *BaseOAuthProvider {
	return &BaseOAuthProvider{
		name:   name,
		title:  title,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *BaseOAuthProvider) Name() string  { return p.name }
func (p *BaseOAuthProvider) Title() string { return p.title }

// GetConfigOptions 获取配置选项
func (p *BaseOAuthProvider) GetConfigOptions() []ConfigOption {
	return []ConfigOption{
		{Type: "text", Name: "App ID", Key: "app_id", Placeholder: "应用ID", Description: "应用唯一标识", Required: true},
		{Type: "text", Name: "App Secret", Key: "app_secret", Placeholder: "应用密钥", Description: "应用密钥", Required: true},
		{Type: "text", Name: "回调地址", Key: "callback_url", Placeholder: "https://example.com/oauth/callback", Description: "OAuth回调地址", Required: true},
	}
}

// getAccessToken 获取访问令牌
func (p *BaseOAuthProvider) getAccessToken(ctx context.Context, code string, config map[string]interface{}) (string, error) {
	appID, _ := config["app_id"].(string)
	appSecret, _ := config["app_secret"].(string)
	callbackURL, _ := config["callback_url"].(string)

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("client_id", appID)
	data.Set("client_secret", appSecret)
	data.Set("redirect_uri", callbackURL)

	req, err := http.NewRequestWithContext(ctx, "POST", p.tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if token, ok := result["access_token"].(string); ok {
		return token, nil
	}

	return "", fmt.Errorf("获取access_token失败: %s", string(body))
}

// httpRequest 发送HTTP请求
func (p *BaseOAuthProvider) httpRequest(ctx context.Context, method, reqURL string, data url.Values, headers map[string]string) ([]byte, error) {
	var req *http.Request
	var err error

	if data != nil {
		req, err = http.NewRequestWithContext(ctx, method, reqURL, strings.NewReader(data.Encode()))
	} else {
		req, err = http.NewRequestWithContext(ctx, method, reqURL, nil)
	}
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
