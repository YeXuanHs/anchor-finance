package wechat

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/anchor-finance/backend/pkg/oauth"
)

const (
	authURL      = "https://open.weixin.qq.com/connect/qrconnect"
	tokenURL     = "https://api.weixin.qq.com/sns/oauth2/access_token"
	userInfoURL  = "https://api.weixin.qq.com/sns/userinfo"
	providerName = "wechat"
)

// Provider implements the WeChat OAuth provider (PC扫码登录).
type Provider struct {
	AppID     string
	AppSecret string
	RedirectURL string
}

// New creates a new WeChat OAuth provider.
func New(appID, appSecret, redirectURL string) *Provider {
	return &Provider{
		AppID:       appID,
		AppSecret:   appSecret,
		RedirectURL: redirectURL,
	}
}

// Name returns the provider name.
func (p *Provider) Name() string {
	return providerName
}

// GetAuthURL returns the WeChat QR code authorization URL.
func (p *Provider) GetAuthURL(state string) string {
	params := url.Values{}
	params.Set("appid", p.AppID)
	params.Set("redirect_uri", p.RedirectURL)
	params.Set("response_type", "code")
	params.Set("scope", "snsapi_login")
	params.Set("state", state)
	return authURL + "?" + params.Encode() + "#wechat_redirect"
}

// GetUserInfo exchanges the authorization code for user information.
func (p *Provider) GetUserInfo(code string) (*oauth.UserInfo, error) {
	// Exchange code for access_token and openid
	tokenResp, err := p.exchangeCode(code)
	if err != nil {
		return nil, fmt.Errorf("wechat: exchange code: %w", err)
	}

	accessToken, _ := tokenResp["access_token"].(string)
	openID, _ := tokenResp["openid"].(string)
	unionID, _ := tokenResp["unionid"].(string)

	if accessToken == "" || openID == "" {
		return nil, fmt.Errorf("wechat: invalid token response")
	}

	// Get user profile
	profile, err := p.fetchUserInfo(accessToken, openID)
	if err != nil {
		return nil, fmt.Errorf("wechat: fetch user info: %w", err)
	}

	userInfo := &oauth.UserInfo{
		Provider: providerName,
		OpenID:   openID,
		UnionID:  unionID,
		Username: getString(profile, "nickname"),
		Avatar:   getString(profile, "headimgurl"),
		RawData:  profile,
	}

	return userInfo, nil
}

func (p *Provider) exchangeCode(code string) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("appid", p.AppID)
	params.Set("secret", p.AppSecret)
	params.Set("code", code)
	params.Set("grant_type", "authorization_code")

	reqURL := tokenURL + "?" + params.Encode()
	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if errCode, ok := result["errcode"]; ok && errCode != nil && errCode != float64(0) {
		return nil, fmt.Errorf("errcode=%v: %v", errCode, result["errmsg"])
	}

	return result, nil
}

func (p *Provider) fetchUserInfo(accessToken, openID string) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("access_token", accessToken)
	params.Set("openid", openID)
	params.Set("lang", "zh_CN")

	reqURL := userInfoURL + "?" + params.Encode()
	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var profile map[string]interface{}
	if err := json.Unmarshal(body, &profile); err != nil {
		return nil, err
	}

	if errCode, ok := profile["errcode"]; ok && errCode != nil && errCode != float64(0) {
		return nil, fmt.Errorf("errcode=%v: %v", errCode, profile["errmsg"])
	}

	return profile, nil
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
