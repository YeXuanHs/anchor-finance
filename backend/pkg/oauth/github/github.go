package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/anchor-finance/backend/pkg/oauth"
)

const (
	authURL         = "https://github.com/login/oauth/authorize"
	tokenURL        = "https://github.com/login/oauth/access_token"
	userInfoURL     = "https://api.github.com/user"
	userEmailURL    = "https://api.github.com/user/emails"
	providerName    = "github"
)

// Provider implements the GitHub OAuth provider.
type Provider struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// New creates a new GitHub OAuth provider.
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

// GetAuthURL returns the GitHub OAuth authorization URL.
func (p *Provider) GetAuthURL(state string) string {
	params := url.Values{}
	params.Set("client_id", p.ClientID)
	params.Set("redirect_uri", p.RedirectURL)
	params.Set("scope", "read:user user:email")
	params.Set("state", state)
	return authURL + "?" + params.Encode()
}

// GetUserInfo exchanges the authorization code for user information.
func (p *Provider) GetUserInfo(code string) (*oauth.UserInfo, error) {
	// Exchange code for access token
	accessToken, err := p.exchangeCode(code)
	if err != nil {
		return nil, fmt.Errorf("github: exchange code: %w", err)
	}

	// Get user profile
	profile, err := p.fetchProfile(accessToken)
	if err != nil {
		return nil, fmt.Errorf("github: fetch profile: %w", err)
	}

	// Get user email (may be private)
	email := p.fetchEmail(accessToken)

	userInfo := &oauth.UserInfo{
		Provider: providerName,
		OpenID:   fmt.Sprintf("%v", profile["id"]),
		Username: getString(profile, "login"),
		Email:    email,
		Avatar:   getString(profile, "avatar_url"),
		RawData:  profile,
	}

	return userInfo, nil
}

func (p *Provider) exchangeCode(code string) (string, error) {
	params := url.Values{}
	params.Set("client_id", p.ClientID)
	params.Set("client_secret", p.ClientSecret)
	params.Set("code", code)

	req, err := http.NewRequest(http.MethodPost, tokenURL, nil)
	if err != nil {
		return "", err
	}
	req.URL.RawQuery = params.Encode()
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
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
		return "", fmt.Errorf("github: no access_token in response")
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

func (p *Provider) fetchEmail(accessToken string) string {
	req, err := http.NewRequest(http.MethodGet, userEmailURL, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	if err := json.Unmarshal(body, &emails); err != nil {
		return ""
	}

	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email
		}
	}

	return ""
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
