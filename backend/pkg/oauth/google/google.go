package google

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/anchor-finance/backend/pkg/oauth"
)

const (
	authURL      = "https://accounts.google.com/o/oauth2/v2/auth"
	tokenURL     = "https://oauth2.googleapis.com/token"
	userInfoURL  = "https://www.googleapis.com/oauth2/v2/userinfo"
	providerName = "google"
)

// Provider implements the Google OAuth provider.
type Provider struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// New creates a new Google OAuth provider.
func New(clientID, clientSecret, redirectURL string) *Provider {
	return &Provider{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
	}
}

// Name returns the provider name.
func (p *Provider) Name() string {
	return providerName
}

// GetAuthURL returns the Google OAuth authorization URL.
func (p *Provider) GetAuthURL(state string) string {
	params := url.Values{}
	params.Set("client_id", p.ClientID)
	params.Set("redirect_uri", p.RedirectURL)
	params.Set("response_type", "code")
	params.Set("scope", "openid email profile")
	params.Set("state", state)
	return authURL + "?" + params.Encode()
}

// GetUserInfo exchanges the authorization code for user information.
func (p *Provider) GetUserInfo(code string) (*oauth.UserInfo, error) {
	// Exchange code for access token
	accessToken, err := p.exchangeCode(code)
	if err != nil {
		return nil, fmt.Errorf("google: exchange code: %w", err)
	}

	// Get user profile
	profile, err := p.fetchProfile(accessToken)
	if err != nil {
		return nil, fmt.Errorf("google: fetch profile: %w", err)
	}

	userInfo := &oauth.UserInfo{
		Provider: providerName,
		OpenID:   getString(profile, "id"),
		Username: getString(profile, "name"),
		Email:    getString(profile, "email"),
		Avatar:   getString(profile, "picture"),
		RawData:  profile,
	}

	return userInfo, nil
}

func (p *Provider) exchangeCode(code string) (string, error) {
	data := url.Values{}
	data.Set("client_id", p.ClientID)
	data.Set("client_secret", p.ClientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", p.RedirectURL)

	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
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

	if errMsg, ok := result["error"]; ok {
		return "", fmt.Errorf("%v: %v", errMsg, result["error_description"])
	}

	token, ok := result["access_token"].(string)
	if !ok || token == "" {
		return "", fmt.Errorf("google: no access_token in response")
	}

	return token, nil
}

func (p *Provider) fetchProfile(accessToken string) (map[string]interface{}, error) {
	req, err := http.NewRequest(http.MethodGet, userInfoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
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
