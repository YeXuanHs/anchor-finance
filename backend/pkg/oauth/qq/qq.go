package qq

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"anchorfinance/pkg/oauth"
)

const (
	authURL         = "https://graph.qq.com/oauth2.0/authorize"
	tokenURL        = "https://graph.qq.com/oauth2.0/token"
	openIDURL       = "https://graph.qq.com/oauth2.0/me"
	userInfoURL     = "https://graph.qq.com/user/get_user_info"
	providerName    = "qq"
)

// Provider implements the QQ OAuth provider.
type Provider struct {
	AppID     string
	AppKey    string
	RedirectURL string
}

// New creates a new QQ OAuth provider.
func New(appID, appKey, redirectURL string) *Provider {
	return &Provider{
		AppID:       appID,
		AppKey:      appKey,
		RedirectURL: redirectURL,
	}
}

// Name returns the provider name.
func (p *Provider) Name() string {
	return providerName
}

// GetAuthURL returns the QQ OAuth authorization URL.
func (p *Provider) GetAuthURL(state string) string {
	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("client_id", p.AppID)
	params.Set("redirect_uri", p.RedirectURL)
	params.Set("state", state)
	params.Set("scope", "get_user_info")
	return authURL + "?" + params.Encode()
}

// GetUserInfo exchanges the authorization code for user information.
func (p *Provider) GetUserInfo(code string) (*oauth.UserInfo, error) {
	// Exchange code for access token
	accessToken, err := p.exchangeCode(code)
	if err != nil {
		return nil, fmt.Errorf("qq: exchange code: %w", err)
	}

	// Get OpenID
	openID, err := p.getOpenID(accessToken)
	if err != nil {
		return nil, fmt.Errorf("qq: get openid: %w", err)
	}

	// Get user profile
	profile, err := p.fetchUserInfo(accessToken, openID)
	if err != nil {
		return nil, fmt.Errorf("qq: fetch user info: %w", err)
	}

	userInfo := &oauth.UserInfo{
		Provider: providerName,
		OpenID:   openID,
		Username: getString(profile, "nickname"),
		Avatar:   getString(profile, "figureurl_qq_2"),
		RawData:  profile,
	}

	return userInfo, nil
}

func (p *Provider) exchangeCode(code string) (string, error) {
	params := url.Values{}
	params.Set("grant_type", "authorization_code")
	params.Set("client_id", p.AppID)
	params.Set("client_secret", p.AppKey)
	params.Set("code", code)
	params.Set("redirect_uri", p.RedirectURL)
	params.Set("fmt", "json")

	reqURL := tokenURL + "?" + params.Encode()
	resp, err := http.Get(reqURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if errMsg, ok := result["error_description"]; ok {
		return "", fmt.Errorf("%v", errMsg)
	}

	token, ok := result["access_token"].(string)
	if !ok || token == "" {
		return "", fmt.Errorf("qq: no access_token in response")
	}

	return token, nil
}

func (p *Provider) getOpenID(accessToken string) (string, error) {
	params := url.Values{}
	params.Set("access_token", accessToken)
	params.Set("fmt", "json")

	reqURL := openIDURL + "?" + params.Encode()
	resp, err := http.Get(reqURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	openID, ok := result["openid"].(string)
	if !ok || openID == "" {
		return "", fmt.Errorf("qq: no openid in response")
	}

	return openID, nil
}

func (p *Provider) fetchUserInfo(accessToken, openID string) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("access_token", accessToken)
	params.Set("oauth_consumer_key", p.AppID)
	params.Set("openid", openID)
	params.Set("format", "json")

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

	ret, _ := profile["ret"].(float64)
	if ret != 0 {
		return nil, fmt.Errorf("qq: get_user_info ret=%v: %v", ret, profile["msg"])
	}

	return profile, nil
}

// callbackResponse is used to parse QQ callback which returns as key=value pairs.
func parseTokenResponse(raw string) (string, error) {
	re := regexp.MustCompile(`access_token=([^&]+)`)
	matches := re.FindStringSubmatch(raw)
	if len(matches) < 2 {
		return "", fmt.Errorf("qq: no access_token in response")
	}
	return strings.TrimSpace(matches[1]), nil
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
