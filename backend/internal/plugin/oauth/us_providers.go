package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// 海外OAuth平台实现

func init() {
	// 13. Google
	RegisterProvider("google", func(config map[string]interface{}) (OAuthProvider, error) {
		return &GoogleProvider{
			BaseOAuthProvider: NewBaseOAuthProvider("google", "Google登录"),
		}, nil
	})

	// 14. Facebook
	RegisterProvider("facebook", func(config map[string]interface{}) (OAuthProvider, error) {
		return &FacebookProvider{
			BaseOAuthProvider: NewBaseOAuthProvider("facebook", "Facebook登录"),
		}, nil
	})

	// 15. Twitter/X
	RegisterProvider("twitter", func(config map[string]interface{}) (OAuthProvider, error) {
		return &TwitterProvider{
			BaseOAuthProvider: NewBaseOAuthProvider("twitter", "Twitter/X登录"),
		}, nil
	})

	// 16. GitHub
	RegisterProvider("github", func(config map[string]interface{}) (OAuthProvider, error) {
		return &GitHubProvider{
			BaseOAuthProvider: NewBaseOAuthProvider("github", "GitHub登录"),
		}, nil
	})

	// 17. LinkedIn
	RegisterProvider("linkedin", func(config map[string]interface{}) (OAuthProvider, error) {
		return &LinkedInProvider{
			BaseOAuthProvider: NewBaseOAuthProvider("linkedin", "LinkedIn登录"),
		}, nil
	})

	// 18. Microsoft
	RegisterProvider("microsoft", func(config map[string]interface{}) (OAuthProvider, error) {
		return &MicrosoftProvider{
			BaseOAuthProvider: NewBaseOAuthProvider("microsoft", "Microsoft登录"),
		}, nil
	})

	// 19. Apple
	RegisterProvider("apple", func(config map[string]interface{}) (OAuthProvider, error) {
		return &AppleProvider{
			BaseOAuthProvider: NewBaseOAuthProvider("apple", "Apple登录"),
		}, nil
	})

	// 20. Amazon
	RegisterProvider("amazon", func(config map[string]interface{}) (OAuthProvider, error) {
		return &AmazonProvider{
			BaseOAuthProvider: NewBaseOAuthProvider("amazon", "Amazon登录"),
		}, nil
	})

	// 21. Discord
	RegisterProvider("discord", func(config map[string]interface{}) (OAuthProvider, error) {
		return &DiscordProvider{
			BaseOAuthProvider: NewBaseOAuthProvider("discord", "Discord登录"),
		}, nil
	})

	// 22. Slack
	RegisterProvider("slack", func(config map[string]interface{}) (OAuthProvider, error) {
		return &SlackProvider{
			BaseOAuthProvider: NewBaseOAuthProvider("slack", "Slack登录"),
		}, nil
	})

	// 23. Telegram
	RegisterProvider("telegram", func(config map[string]interface{}) (OAuthProvider, error) {
		return &TelegramProvider{
			BaseOAuthProvider: NewBaseOAuthProvider("telegram", "Telegram登录"),
		}, nil
	})

	// 24. LINE
	RegisterProvider("line", func(config map[string]interface{}) (OAuthProvider, error) {
		return &LineProvider{
			BaseOAuthProvider: NewBaseOAuthProvider("line", "LINE登录"),
		}, nil
	})
}

// GoogleProvider Google登录
type GoogleProvider struct {
	*BaseOAuthProvider
}

func (p *GoogleProvider) GetConfigOptions() []ConfigOption {
	options := p.BaseOAuthProvider.GetConfigOptions()
	options = append(options, ConfigOption{
		Type: "text", Name: "Hosted Domain", Key: "hosted_domain",
		Placeholder: "example.com", Description: "限制特定域名的用户登录（可选）",
	})
	return options
}

func (p *GoogleProvider) GetLoginURL(ctx context.Context, params *LoginParams) (string, error) {
	clientID, _ := params.Config["app_id"].(string)
	
	loginURL := fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=openid%%20email%%20profile&state=%s",
		clientID,
		url.QueryEscape(params.CallbackURL),
		params.State,
	)
	return loginURL, nil
}

func (p *GoogleProvider) HandleCallback(ctx context.Context, params *CallbackParams) (*UserInfo, error) {
	clientID, _ := params.Config["app_id"].(string)
	clientSecret, _ := params.Config["app_secret"].(string)

	// 获取access_token
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", params.Code)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("redirect_uri", params.Config["callback_url"].(string))

	body, err := p.httpRequest(ctx, "POST", "https://oauth2.googleapis.com/token", data, nil)
	if err != nil {
		return nil, err
	}

	var tokenResult map[string]interface{}
	json.Unmarshal(body, &tokenResult)
	accessToken, _ := tokenResult["access_token"].(string)

	// 获取用户信息
	userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo"
	headers := map[string]string{"Authorization": "Bearer " + accessToken}
	body, err = p.httpRequest(ctx, "GET", userInfoURL, nil, headers)
	if err != nil {
		return nil, err
	}

	var userInfo map[string]interface{}
	json.Unmarshal(body, &userInfo)

	return &UserInfo{
		OpenID:   userInfo["id"].(string),
		Username: userInfo["name"].(string),
		Email:    userInfo["email"].(string),
		Avatar:   userInfo["picture"].(string),
		Data:     userInfo,
	}, nil
}

// FacebookProvider Facebook登录
type FacebookProvider struct {
	*BaseOAuthProvider
}

func (p *FacebookProvider) GetLoginURL(ctx context.Context, params *LoginParams) (string, error) {
	appID, _ := params.Config["app_id"].(string)
	
	loginURL := fmt.Sprintf(
		"https://www.facebook.com/v18.0/dialog/oauth?client_id=%s&redirect_uri=%s&state=%s&scope=email,public_profile",
		appID,
		url.QueryEscape(params.CallbackURL),
		params.State,
	)
	return loginURL, nil
}

func (p *FacebookProvider) HandleCallback(ctx context.Context, params *CallbackParams) (*UserInfo, error) {
	appID, _ := params.Config["app_id"].(string)
	appSecret, _ := params.Config["app_secret"].(string)

	// 获取access_token
	tokenURL := fmt.Sprintf(
		"https://graph.facebook.com/v18.0/oauth/access_token?client_id=%s&client_secret=%s&code=%s&redirect_uri=%s",
		appID, appSecret, params.Code, url.QueryEscape(params.Config["callback_url"].(string)),
	)

	body, err := p.httpRequest(ctx, "GET", tokenURL, nil, nil)
	if err != nil {
		return nil, err
	}

	var tokenResult map[string]interface{}
	json.Unmarshal(body, &tokenResult)
	accessToken, _ := tokenResult["access_token"].(string)

	// 获取用户信息
	userInfoURL := fmt.Sprintf(
		"https://graph.facebook.com/me?fields=id,name,email,picture&access_token=%s",
		accessToken,
	)
	body, err = p.httpRequest(ctx, "GET", userInfoURL, nil, nil)
	if err != nil {
		return nil, err
	}

	var userInfo map[string]interface{}
	json.Unmarshal(body, &userInfo)

	avatar := ""
	if picture, ok := userInfo["picture"].(map[string]interface{}); ok {
		if data, ok := picture["data"].(map[string]interface{}); ok {
			avatar, _ = data["url"].(string)
		}
	}

	return &UserInfo{
		OpenID:   userInfo["id"].(string),
		Username: userInfo["name"].(string),
		Email:    userInfo["email"].(string),
		Avatar:   avatar,
		Data:     userInfo,
	}, nil
}

// TwitterProvider Twitter/X登录
type TwitterProvider struct {
	*BaseOAuthProvider
}

func (p *TwitterProvider) GetLoginURL(ctx context.Context, params *LoginParams) (string, error) {
	// Twitter使用OAuth 2.0 PKCE
	clientID, _ := params.Config["app_id"].(string)
	
	loginURL := fmt.Sprintf(
		"https://twitter.com/i/oauth2/authorize?response_type=code&client_id=%s&redirect_uri=%s&state=%s&scope=users.read%%20tweet.read",
		clientID,
		url.QueryEscape(params.CallbackURL),
		params.State,
	)
	return loginURL, nil
}

func (p *TwitterProvider) HandleCallback(ctx context.Context, params *CallbackParams) (*UserInfo, error) {
	// Twitter OAuth 2.0需要PKCE，实现较复杂
	return nil, fmt.Errorf("Twitter登录需要PKCE支持，请使用官方SDK")
}

// GitHubProvider GitHub登录
type GitHubProvider struct {
	*BaseOAuthProvider
}

func (p *GitHubProvider) GetLoginURL(ctx context.Context, params *LoginParams) (string, error) {
	clientID, _ := params.Config["app_id"].(string)
	
	loginURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&state=%s&scope=read:user,user:email",
		clientID,
		url.QueryEscape(params.CallbackURL),
		params.State,
	)
	return loginURL, nil
}

func (p *GitHubProvider) HandleCallback(ctx context.Context, params *CallbackParams) (*UserInfo, error) {
	clientID, _ := params.Config["app_id"].(string)
	clientSecret, _ := params.Config["app_secret"].(string)

	// 获取access_token
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", params.Code)

	body, err := p.httpRequest(ctx, "POST", "https://github.com/login/oauth/access_token", data, map[string]string{"Accept": "application/json"})
	if err != nil {
		return nil, err
	}

	var tokenResult map[string]interface{}
	json.Unmarshal(body, &tokenResult)
	accessToken, _ := tokenResult["access_token"].(string)

	// 获取用户信息
	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
		"Accept":        "application/vnd.github.v3+json",
	}
	body, err = p.httpRequest(ctx, "GET", "https://api.github.com/user", nil, headers)
	if err != nil {
		return nil, err
	}

	var userInfo map[string]interface{}
	json.Unmarshal(body, &userInfo)

	return &UserInfo{
		OpenID:   fmt.Sprintf("%.0f", userInfo["id"].(float64)),
		Username: userInfo["login"].(string),
		Email:    userInfo["email"].(string),
		Avatar:   userInfo["avatar_url"].(string),
		Data:     userInfo,
	}, nil
}

// LinkedInProvider LinkedIn登录
type LinkedInProvider struct {
	*BaseOAuthProvider
}

func (p *LinkedInProvider) GetLoginURL(ctx context.Context, params *LoginParams) (string, error) {
	clientID, _ := params.Config["app_id"].(string)
	
	loginURL := fmt.Sprintf(
		"https://www.linkedin.com/oauth/v2/authorization?response_type=code&client_id=%s&redirect_uri=%s&state=%s&scope=r_liteprofile%%20r_emailaddress",
		clientID,
		url.QueryEscape(params.CallbackURL),
		params.State,
	)
	return loginURL, nil
}

func (p *LinkedInProvider) HandleCallback(ctx context.Context, params *CallbackParams) (*UserInfo, error) {
	clientID, _ := params.Config["app_id"].(string)
	clientSecret, _ := params.Config["app_secret"].(string)

	// 获取access_token
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", params.Code)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("redirect_uri", params.Config["callback_url"].(string))

	body, err := p.httpRequest(ctx, "POST", "https://www.linkedin.com/oauth/v2/accessToken", data, nil)
	if err != nil {
		return nil, err
	}

	var tokenResult map[string]interface{}
	json.Unmarshal(body, &tokenResult)
	accessToken, _ := tokenResult["access_token"].(string)

	// 获取用户信息
	headers := map[string]string{"Authorization": "Bearer " + accessToken}
	body, err = p.httpRequest(ctx, "GET", "https://api.linkedin.com/v2/me?projection=(id,firstName,lastName,profilePicture(displayImage~:playableStreams))", nil, headers)
	if err != nil {
		return nil, err
	}

	var userInfo map[string]interface{}
	json.Unmarshal(body, &userInfo)

	return &UserInfo{
		OpenID:   userInfo["id"].(string),
		Username: userInfo["firstName"].(string) + " " + userInfo["lastName"].(string),
		Data:     userInfo,
	}, nil
}

// MicrosoftProvider Microsoft登录
type MicrosoftProvider struct {
	*BaseOAuthProvider
}

func (p *MicrosoftProvider) GetLoginURL(ctx context.Context, params *LoginParams) (string, error) {
	clientID, _ := params.Config["app_id"].(string)
	
	loginURL := fmt.Sprintf(
		"https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=%s&response_type=code&redirect_uri=%s&state=%s&scope=openid%%20email%%20profile",
		clientID,
		url.QueryEscape(params.CallbackURL),
		params.State,
	)
	return loginURL, nil
}

func (p *MicrosoftProvider) HandleCallback(ctx context.Context, params *CallbackParams) (*UserInfo, error) {
	clientID, _ := params.Config["app_id"].(string)
	clientSecret, _ := params.Config["app_secret"].(string)

	// 获取access_token
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", params.Code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", params.Config["callback_url"].(string))

	body, err := p.httpRequest(ctx, "POST", "https://login.microsoftonline.com/common/oauth2/v2.0/token", data, nil)
	if err != nil {
		return nil, err
	}

	var tokenResult map[string]interface{}
	json.Unmarshal(body, &tokenResult)
	accessToken, _ := tokenResult["access_token"].(string)

	// 获取用户信息
	headers := map[string]string{"Authorization": "Bearer " + accessToken}
	body, err = p.httpRequest(ctx, "GET", "https://graph.microsoft.com/v1.0/me", nil, headers)
	if err != nil {
		return nil, err
	}

	var userInfo map[string]interface{}
	json.Unmarshal(body, &userInfo)

	return &UserInfo{
		OpenID:   userInfo["id"].(string),
		Username: userInfo["displayName"].(string),
		Email:    userInfo["mail"].(string),
		Data:     userInfo,
	}, nil
}

// AppleProvider Apple登录
type AppleProvider struct {
	*BaseOAuthProvider
}

func (p *AppleProvider) GetLoginURL(ctx context.Context, params *LoginParams) (string, error) {
	clientID, _ := params.Config["app_id"].(string)
	
	loginURL := fmt.Sprintf(
		"https://appleid.apple.com/auth/authorize?client_id=%s&redirect_uri=%s&response_type=code&state=%s&scope=name%%20email&response_mode=form_post",
		clientID,
		url.QueryEscape(params.CallbackURL),
		params.State,
	)
	return loginURL, nil
}

func (p *AppleProvider) HandleCallback(ctx context.Context, params *CallbackParams) (*UserInfo, error) {
	// Apple登录需要JWT验证，实现较复杂
	return nil, fmt.Errorf("Apple登录需要JWT支持，请使用官方SDK")
}

// AmazonProvider Amazon登录
type AmazonProvider struct {
	*BaseOAuthProvider
}

func (p *AmazonProvider) GetLoginURL(ctx context.Context, params *LoginParams) (string, error) {
	clientID, _ := params.Config["app_id"].(string)
	
	loginURL := fmt.Sprintf(
		"https://www.amazon.com/ap/oa?client_id=%s&response_type=code&redirect_uri=%s&state=%s&scope=profile",
		clientID,
		url.QueryEscape(params.CallbackURL),
		params.State,
	)
	return loginURL, nil
}

func (p *AmazonProvider) HandleCallback(ctx context.Context, params *CallbackParams) (*UserInfo, error) {
	clientID, _ := params.Config["app_id"].(string)
	clientSecret, _ := params.Config["app_secret"].(string)

	// 获取access_token
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", params.Code)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("redirect_uri", params.Config["callback_url"].(string))

	body, err := p.httpRequest(ctx, "POST", "https://api.amazon.com/auth/o2/token", data, nil)
	if err != nil {
		return nil, err
	}

	var tokenResult map[string]interface{}
	json.Unmarshal(body, &tokenResult)
	accessToken, _ := tokenResult["access_token"].(string)

	// 获取用户信息
	headers := map[string]string{"Authorization": "Bearer " + accessToken}
	body, err = p.httpRequest(ctx, "GET", "https://api.amazon.com/user/profile", nil, headers)
	if err != nil {
		return nil, err
	}

	var userInfo map[string]interface{}
	json.Unmarshal(body, &userInfo)

	return &UserInfo{
		OpenID:   userInfo["user_id"].(string),
		Username: userInfo["name"].(string),
		Email:    userInfo["email"].(string),
		Data:     userInfo,
	}, nil
}

// DiscordProvider Discord登录
type DiscordProvider struct {
	*BaseOAuthProvider
}

func (p *DiscordProvider) GetLoginURL(ctx context.Context, params *LoginParams) (string, error) {
	clientID, _ := params.Config["app_id"].(string)
	
	loginURL := fmt.Sprintf(
		"https://discord.com/api/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&state=%s&scope=identify%%20email",
		clientID,
		url.QueryEscape(params.CallbackURL),
		params.State,
	)
	return loginURL, nil
}

func (p *DiscordProvider) HandleCallback(ctx context.Context, params *CallbackParams) (*UserInfo, error) {
	clientID, _ := params.Config["app_id"].(string)
	clientSecret, _ := params.Config["app_secret"].(string)

	// 获取access_token
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", params.Code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", params.Config["callback_url"].(string))

	body, err := p.httpRequest(ctx, "POST", "https://discord.com/api/oauth2/token", data, nil)
	if err != nil {
		return nil, err
	}

	var tokenResult map[string]interface{}
	json.Unmarshal(body, &tokenResult)
	accessToken, _ := tokenResult["access_token"].(string)

	// 获取用户信息
	headers := map[string]string{"Authorization": "Bearer " + accessToken}
	body, err = p.httpRequest(ctx, "GET", "https://discord.com/api/users/@me", nil, headers)
	if err != nil {
		return nil, err
	}

	var userInfo map[string]interface{}
	json.Unmarshal(body, &userInfo)

	avatar := fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", userInfo["id"], userInfo["avatar"])

	return &UserInfo{
		OpenID:   userInfo["id"].(string),
		Username: userInfo["username"].(string),
		Email:    userInfo["email"].(string),
		Avatar:   avatar,
		Data:     userInfo,
	}, nil
}

// SlackProvider Slack登录
type SlackProvider struct {
	*BaseOAuthProvider
}

func (p *SlackProvider) GetLoginURL(ctx context.Context, params *LoginParams) (string, error) {
	clientID, _ := params.Config["app_id"].(string)
	
	loginURL := fmt.Sprintf(
		"https://slack.com/oauth/v2/authorize?client_id=%s&redirect_uri=%s&state=%s&scope=openid%%20profile%%20email",
		clientID,
		url.QueryEscape(params.CallbackURL),
		params.State,
	)
	return loginURL, nil
}

func (p *SlackProvider) HandleCallback(ctx context.Context, params *CallbackParams) (*UserInfo, error) {
	clientID, _ := params.Config["app_id"].(string)
	clientSecret, _ := params.Config["app_secret"].(string)

	// 获取access_token
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", params.Code)

	body, err := p.httpRequest(ctx, "POST", "https://slack.com/api/oauth.v2.access", data, nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if !result["ok"].(bool) {
		return nil, fmt.Errorf("Slack认证失败: %s", result["error"].(string))
	}

	authedUser := result["authed_user"].(map[string]interface{})
	userID, _ := authedUser["id"].(string)

	// 获取用户信息
	headers := map[string]string{"Authorization": "Bearer " + authedUser["access_token"].(string)}
	body, err = p.httpRequest(ctx, "GET", "https://slack.com/api/users.info?user="+userID, nil, headers)
	if err != nil {
		return nil, err
	}

	var userInfoResult map[string]interface{}
	json.Unmarshal(body, &userInfoResult)
	user := userInfoResult["user"].(map[string]interface{})

	return &UserInfo{
		OpenID:   userID,
		Username: user["name"].(string),
		Email:    user["email"].(string),
		Avatar:   user["image_192"].(string),
		Data:     user,
	}, nil
}

// TelegramProvider Telegram登录
type TelegramProvider struct {
	*BaseOAuthProvider
}

func (p *TelegramProvider) GetConfigOptions() []ConfigOption {
	return []ConfigOption{
		{Type: "text", Name: "Bot Token", Key: "bot_token", Placeholder: "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11", Description: "Telegram Bot Token", Required: true},
		{Type: "text", Name: "Bot Username", Key: "bot_username", Placeholder: "MyBot", Description: "Bot用户名", Required: true},
	}
}

func (p *TelegramProvider) GetLoginURL(ctx context.Context, params *LoginParams) (string, error) {
	botUsername, _ := params.Config["bot_username"].(string)
	
	// Telegram使用Login Widget，需要前端集成
	return fmt.Sprintf("https://oauth.telegram.org/auth?bot_id=%s&origin=%s&request_access=write", 
		botUsername, url.QueryEscape(params.CallbackURL)), nil
}

func (p *TelegramProvider) HandleCallback(ctx context.Context, params *CallbackParams) (*UserInfo, error) {
	// Telegram登录需要验证hash，实现参考官方文档
	return nil, fmt.Errorf("Telegram登录需要前端Widget集成，请参考官方文档")
}

// LineProvider LINE登录
type LineProvider struct {
	*BaseOAuthProvider
}

func (p *LineProvider) GetLoginURL(ctx context.Context, params *LoginParams) (string, error) {
	clientID, _ := params.Config["app_id"].(string)
	
	loginURL := fmt.Sprintf(
		"https://access.line.me/oauth2/v2.1/authorize?client_id=%s&redirect_uri=%s&response_type=code&state=%s&scope=profile%%20openid%%20email",
		clientID,
		url.QueryEscape(params.CallbackURL),
		params.State,
	)
	return loginURL, nil
}

func (p *LineProvider) HandleCallback(ctx context.Context, params *CallbackParams) (*UserInfo, error) {
	clientID, _ := params.Config["app_id"].(string)
	clientSecret, _ := params.Config["app_secret"].(string)

	// 获取access_token
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", params.Code)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("redirect_uri", params.Config["callback_url"].(string))

	body, err := p.httpRequest(ctx, "POST", "https://api.line.me/oauth2/v2.1/token", data, nil)
	if err != nil {
		return nil, err
	}

	var tokenResult map[string]interface{}
	json.Unmarshal(body, &tokenResult)
	accessToken, _ := tokenResult["access_token"].(string)
	idToken, _ := tokenResult["id_token"].(string)

	// 解析id_token获取用户信息（简化处理）
	// 实际应该验证JWT签名
	headers := map[string]string{"Authorization": "Bearer " + accessToken}
	body, err = p.httpRequest(ctx, "GET", "https://api.line.me/v2/profile", nil, headers)
	if err != nil {
		return nil, err
	}

	var userInfo map[string]interface{}
	json.Unmarshal(body, &userInfo)

	return &UserInfo{
		OpenID:   userInfo["userId"].(string),
		Username: userInfo["displayName"].(string),
		Avatar:   userInfo["pictureUrl"].(string),
		Data:     userInfo,
	}, nil
}
