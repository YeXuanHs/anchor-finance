package juhe

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/anchor-finance/backend/pkg/oauth"
)

// Provider implements the aggregated login (聚合登录) OAuth provider.
// It wraps third-party social login APIs (QQ, WeChat, Alipay, etc.)
// through a unified aggregation service.
type Provider struct {
	apiURL     string // e.g. https://login.blogcloud.cn
	appID      string
	appKey     string
	loginType  string // qq, wx, alipay, sina, github, gitee, dingtalk
	redirectURL string
}

// New creates a new aggregated login provider.
func New(apiURL, appID, appKey, loginType, redirectURL string) *Provider {
	return &Provider{
		apiURL:      apiURL,
		appID:       appID,
		appKey:      appKey,
		loginType:   loginType,
		redirectURL: redirectURL,
	}
}

// Name returns the provider name in the format "juhe_{type}".
func (p *Provider) Name() string {
	return "juhe_" + p.loginType
}

// LoginType returns the underlying social login type (qq, wx, etc.).
func (p *Provider) LoginType() string {
	return p.loginType
}

// GetAuthURL returns the aggregated login authorization URL.
func (p *Provider) GetAuthURL(state string) string {
	redirectURI := url.QueryEscape(p.redirectURL)
	return fmt.Sprintf("%s/connect.php?act=login&appid=%s&appkey=%s&type=%s&redirect_uri=%s&state=%s",
		p.apiURL, p.appID, p.appKey, p.loginType, redirectURI, state)
}

// GetUserInfo exchanges the authorization code for user information
// by calling the aggregated login callback API.
func (p *Provider) GetUserInfo(code string) (*oauth.UserInfo, error) {
	apiURL := fmt.Sprintf("%s/connect.php?act=callback&appid=%s&appkey=%s&type=%s&code=%s",
		p.apiURL, p.appID, p.appKey, p.loginType, code)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("juhe: callback request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("juhe: read response failed: %w", err)
	}

	var result callbackResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("juhe: parse response failed: %w", err)
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("juhe: api error [%d]: %s", result.Code, result.Msg)
	}

	return &oauth.UserInfo{
		Provider: p.Name(),
		OpenID:   result.Data.SocialUID,
		Username: result.Data.Nickname,
		Avatar:   result.Data.Avatar,
		Email:    result.Data.Email,
		RawData: map[string]interface{}{
			"social_uid": result.Data.SocialUID,
			"nickname":   result.Data.Nickname,
			"avatar":     result.Data.Avatar,
			"sex":        result.Data.Sex,
			"email":      result.Data.Email,
		},
	}, nil
}

// QueryUser queries user information by social_uid through the aggregation API.
func (p *Provider) QueryUser(socialUID string) (*oauth.UserInfo, error) {
	apiURL := fmt.Sprintf("%s/connect.php?act=query&appid=%s&appkey=%s&type=%s&social_uid=%s",
		p.apiURL, p.appID, p.appKey, p.loginType, socialUID)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("juhe: query request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("juhe: read response failed: %w", err)
	}

	var result callbackResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("juhe: parse response failed: %w", err)
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("juhe: query error [%d]: %s", result.Code, result.Msg)
	}

	return &oauth.UserInfo{
		Provider: p.Name(),
		OpenID:   result.Data.SocialUID,
		Username: result.Data.Nickname,
		Avatar:   result.Data.Avatar,
		Email:    result.Data.Email,
		RawData: map[string]interface{}{
			"social_uid": result.Data.SocialUID,
			"nickname":   result.Data.Nickname,
			"avatar":     result.Data.Avatar,
			"sex":        result.Data.Sex,
			"email":      result.Data.Email,
		},
	}, nil
}

// callbackResponse is the aggregated login API response structure.
type callbackResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		SocialUID string `json:"social_uid"`
		Nickname  string `json:"nickname"`
		Avatar    string `json:"avatar"`
		Sex       string `json:"sex"`
		Email     string `json:"email"`
	} `json:"data"`
}
