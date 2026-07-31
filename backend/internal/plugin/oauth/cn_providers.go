package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// 国内OAuth平台实现

func init() {
	// 1. 微信
	RegisterProvider("wechat", func(config map[string]interface{}) (OAuthProvider, error) {
		return &WechatProvider{
			BaseOAuthProvider: NewBaseOAuthProvider("wechat", "微信登录"),
		}, nil
	})

	// 2. QQ
	RegisterProvider("qq", func(config map[string]interface{}) (OAuthProvider, error) {
		return &QQProvider{
			BaseOAuthProvider: NewBaseOAuthProvider("qq", "QQ登录"),
		}, nil
	})

	// 3. 微博
	RegisterProvider("weibo", func(config map[string]interface{}) (OAuthProvider, error) {
		return &WeiboProvider{
			BaseOAuthProvider: NewBaseOAuthProvider("weibo", "微博登录"),
		}, nil
	})

	// 4. 支付宝
	RegisterProvider("alipay", func(config map[string]interface{}) (OAuthProvider, error) {
		return &AlipayProvider{
			BaseOAuthProvider: NewBaseOAuthProvider("alipay", "支付宝登录"),
		}, nil
	})

	// 5. 百度
	RegisterProvider("baidu", func(config map[string]interface{}) (OAuthProvider, error) {
		return &BaiduProvider{
			BaseOAuthProvider: NewBaseOAuthProvider("baidu", "百度登录"),
		}, nil
	})

	// 6. 码云
	RegisterProvider("gitee", func(config map[string]interface{}) (OAuthProvider, error) {
		return &GiteeProvider{
			BaseOAuthProvider: NewBaseOAuthProvider("gitee", "码云登录"),
		}, nil
	})

	// 7. 钉钉
	RegisterProvider("dingtalk", func(config map[string]interface{}) (OAuthProvider, error) {
		return &DingTalkProvider{
			BaseOAuthProvider: NewBaseOAuthProvider("dingtalk", "钉钉登录"),
		}, nil
	})

	// 8. 飞书
	RegisterProvider("feishu", func(config map[string]interface{}) (OAuthProvider, error) {
		return &FeishuProvider{
			BaseOAuthProvider: NewBaseOAuthProvider("feishu", "飞书登录"),
		}, nil
	})

	// 9. CSDN
	RegisterProvider("csdn", func(config map[string]interface{}) (OAuthProvider, error) {
		return &CSDNProvider{
			BaseOAuthProvider: NewBaseOAuthProvider("csdn", "CSDN登录"),
		}, nil
	})

	// 10. 开源中国
	RegisterProvider("oschina", func(config map[string]interface{}) (OAuthProvider, error) {
		return &OSChinaProvider{
			BaseOAuthProvider: NewBaseOAuthProvider("oschina", "开源中国登录"),
		}, nil
	})

	// 11. 腾讯云
	RegisterProvider("tencent_cloud", func(config map[string]interface{}) (OAuthProvider, error) {
		return &TencentCloudProvider{
			BaseOAuthProvider: NewBaseOAuthProvider("tencent_cloud", "腾讯云登录"),
		}, nil
	})

	// 12. 阿里云
	RegisterProvider("aliyun", func(config map[string]interface{}) (OAuthProvider, error) {
		return &AliyunProvider{
			BaseOAuthProvider: NewBaseOAuthProvider("aliyun", "阿里云登录"),
		}, nil
	})
}

// WechatProvider 微信登录
type WechatProvider struct {
	*BaseOAuthProvider
}

func (p *WechatProvider) GetLoginURL(ctx context.Context, params *LoginParams) (string, error) {
	appID, _ := params.Config["app_id"].(string)
	state := params.State
	
	loginURL := fmt.Sprintf(
		"https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect",
		appID,
		url.QueryEscape(params.CallbackURL),
		state,
	)
	return loginURL, nil
}

func (p *WechatProvider) HandleCallback(ctx context.Context, params *CallbackParams) (*UserInfo, error) {
	appID, _ := params.Config["app_id"].(string)
	appSecret, _ := params.Config["app_secret"].(string)

	// 获取access_token
	tokenURL := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		appID, appSecret, params.Code,
	)
	
	body, err := p.httpRequest(ctx, "GET", tokenURL, nil, nil)
	if err != nil {
		return nil, err
	}

	var tokenResult map[string]interface{}
	json.Unmarshal(body, &tokenResult)

	accessToken, _ := tokenResult["access_token"].(string)
	openid, _ := tokenResult["openid"].(string)

	if accessToken == "" {
		return nil, fmt.Errorf("获取access_token失败")
	}

	// 获取用户信息
	userInfoURL := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN",
		accessToken, openid,
	)

	body, err = p.httpRequest(ctx, "GET", userInfoURL, nil, nil)
	if err != nil {
		return nil, err
	}

	var userInfo map[string]interface{}
	json.Unmarshal(body, &userInfo)

	return &UserInfo{
		OpenID:   openid,
		Username: userInfo["nickname"].(string),
		Avatar:   userInfo["headimgurl"].(string),
		Gender:   int(userInfo["sex"].(float64)),
		Data:     userInfo,
	}, nil
}

// QQProvider QQ登录
type QQProvider struct {
	*BaseOAuthProvider
}

func (p *QQProvider) GetLoginURL(ctx context.Context, params *LoginParams) (string, error) {
	appID, _ := params.Config["app_id"].(string)
	
	loginURL := fmt.Sprintf(
		"https://graph.qq.com/oauth2.0/authorize?response_type=code&client_id=%s&redirect_uri=%s&state=%s&scope=get_user_info",
		appID,
		url.QueryEscape(params.CallbackURL),
		params.State,
	)
	return loginURL, nil
}

func (p *QQProvider) HandleCallback(ctx context.Context, params *CallbackParams) (*UserInfo, error) {
	appID, _ := params.Config["app_id"].(string)
	appSecret, _ := params.Config["app_secret"].(string)

	// 获取access_token
	tokenURL := fmt.Sprintf(
		"https://graph.qq.com/oauth2.0/token?grant_type=authorization_code&client_id=%s&client_secret=%s&code=%s&redirect_uri=%s",
		appID, appSecret, params.Code, url.QueryEscape(params.Config["callback_url"].(string)),
	)

	body, err := p.httpRequest(ctx, "GET", tokenURL, nil, nil)
	if err != nil {
		return nil, err
	}

	// 解析token (QQ返回的是key=value格式)
	tokenStr := string(body)
	accessToken := parseQQToken(tokenStr)

	// 获取openid
	meURL := fmt.Sprintf("https://graph.qq.com/oauth2.0/me?access_token=%s", accessToken)
	body, err = p.httpRequest(ctx, "GET", meURL, nil, nil)
	if err != nil {
		return nil, err
	}

	var meResult map[string]interface{}
	json.Unmarshal(body, &meResult)
	openid := meResult["openid"].(string)

	// 获取用户信息
	userInfoURL := fmt.Sprintf(
		"https://graph.qq.com/user/get_user_info?access_token=%s&oauth_consumer_key=%s&openid=%s",
		accessToken, appID, openid,
	)

	body, err = p.httpRequest(ctx, "GET", userInfoURL, nil, nil)
	if err != nil {
		return nil, err
	}

	var userInfo map[string]interface{}
	json.Unmarshal(body, &userInfo)

	gender := 0
	if userInfo["gender"].(string) == "男" {
		gender = 1
	} else if userInfo["gender"].(string) == "女" {
		gender = 2
	}

	return &UserInfo{
		OpenID:   openid,
		Username: userInfo["nickname"].(string),
		Avatar:   userInfo["figureurl_qq_2"].(string),
		Gender:   gender,
		Data:     userInfo,
	}, nil
}

// parseQQToken 解析QQ返回的token
func parseQQToken(tokenStr string) string {
	// 格式: access_token=xxx&expires_in=xxx
	parts := strings.Split(tokenStr, "&")
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 && kv[0] == "access_token" {
			return kv[1]
		}
	}
	return ""
}

// WeiboProvider 微博登录
type WeiboProvider struct {
	*BaseOAuthProvider
}

func (p *WeiboProvider) GetLoginURL(ctx context.Context, params *LoginParams) (string, error) {
	appKey, _ := params.Config["app_id"].(string)
	
	loginURL := fmt.Sprintf(
		"https://api.weibo.com/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&state=%s",
		appKey,
		url.QueryEscape(params.CallbackURL),
		params.State,
	)
	return loginURL, nil
}

func (p *WeiboProvider) HandleCallback(ctx context.Context, params *CallbackParams) (*UserInfo, error) {
	appKey, _ := params.Config["app_id"].(string)
	appSecret, _ := params.Config["app_secret"].(string)

	// 获取access_token
	data := url.Values{}
	data.Set("client_id", appKey)
	data.Set("client_secret", appSecret)
	data.Set("grant_type", "authorization_code")
	data.Set("code", params.Code)
	data.Set("redirect_uri", params.Config["callback_url"].(string))

	body, err := p.httpRequest(ctx, "POST", "https://api.weibo.com/oauth2/access_token", data, nil)
	if err != nil {
		return nil, err
	}

	var tokenResult map[string]interface{}
	json.Unmarshal(body, &tokenResult)

	accessToken, _ := tokenResult["access_token"].(string)
	uid, _ := tokenResult["uid"].(string)

	// 获取用户信息
	userInfoURL := fmt.Sprintf("https://api.weibo.com/2/users/show.json?access_token=%s&uid=%s", accessToken, uid)
	body, err = p.httpRequest(ctx, "GET", userInfoURL, nil, nil)
	if err != nil {
		return nil, err
	}

	var userInfo map[string]interface{}
	json.Unmarshal(body, &userInfo)

	gender := 0
	if userInfo["gender"].(string) == "m" {
		gender = 1
	} else if userInfo["gender"].(string) == "f" {
		gender = 2
	}

	return &UserInfo{
		OpenID:   uid,
		Username: userInfo["screen_name"].(string),
		Avatar:   userInfo["avatar_large"].(string),
		Gender:   gender,
		Data:     userInfo,
	}, nil
}

// AlipayProvider 支付宝登录
type AlipayProvider struct {
	*BaseOAuthProvider
}

func (p *AlipayProvider) GetLoginURL(ctx context.Context, params *LoginParams) (string, error) {
	appID, _ := params.Config["app_id"].(string)
	
	loginURL := fmt.Sprintf(
		"https://openauth.alipay.com/oauth2/publicAppAuthorize.htm?app_id=%s&scope=auth_user&redirect_uri=%s&state=%s",
		appID,
		url.QueryEscape(params.CallbackURL),
		params.State,
	)
	return loginURL, nil
}

func (p *AlipayProvider) HandleCallback(ctx context.Context, params *CallbackParams) (*UserInfo, error) {
	// 支付宝需要RSA签名，这里简化处理
	// 实际实现需要使用支付宝SDK
	return nil, fmt.Errorf("支付宝登录需要RSA签名，请使用官方SDK")
}

// BaiduProvider 百度登录
type BaiduProvider struct {
	*BaseOAuthProvider
}

func (p *BaiduProvider) GetLoginURL(ctx context.Context, params *LoginParams) (string, error) {
	appKey, _ := params.Config["app_id"].(string)
	
	loginURL := fmt.Sprintf(
		"https://openapi.baidu.com/oauth/2.0/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=basic&display=page&state=%s",
		appKey,
		url.QueryEscape(params.CallbackURL),
		params.State,
	)
	return loginURL, nil
}

func (p *BaiduProvider) HandleCallback(ctx context.Context, params *CallbackParams) (*UserInfo, error) {
	appKey, _ := params.Config["app_id"].(string)
	appSecret, _ := params.Config["app_secret"].(string)

	// 获取access_token
	tokenURL := fmt.Sprintf(
		"https://openapi.baidu.com/oauth/2.0/token?grant_type=authorization_code&code=%s&client_id=%s&client_secret=%s&redirect_uri=%s",
		params.Code, appKey, appSecret, url.QueryEscape(params.Config["callback_url"].(string)),
	)

	body, err := p.httpRequest(ctx, "GET", tokenURL, nil, nil)
	if err != nil {
		return nil, err
	}

	var tokenResult map[string]interface{}
	json.Unmarshal(body, &tokenResult)
	accessToken, _ := tokenResult["access_token"].(string)

	// 获取用户信息
	userInfoURL := fmt.Sprintf("https://openapi.baidu.com/rest/2.0/passport/users/getInfo?access_token=%s", accessToken)
	body, err = p.httpRequest(ctx, "GET", userInfoURL, nil, nil)
	if err != nil {
		return nil, err
	}

	var userInfo map[string]interface{}
	json.Unmarshal(body, &userInfo)

	return &UserInfo{
		OpenID:   userInfo["userid"].(string),
		Username: userInfo["username"].(string),
		Avatar:   userInfo["portrait"].(string),
		Data:     userInfo,
	}, nil
}

// GiteeProvider 码云登录
type GiteeProvider struct {
	*BaseOAuthProvider
}

func (p *GiteeProvider) GetLoginURL(ctx context.Context, params *LoginParams) (string, error) {
	appID, _ := params.Config["app_id"].(string)
	
	loginURL := fmt.Sprintf(
		"https://gitee.com/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code&state=%s",
		appID,
		url.QueryEscape(params.CallbackURL),
		params.State,
	)
	return loginURL, nil
}

func (p *GiteeProvider) HandleCallback(ctx context.Context, params *CallbackParams) (*UserInfo, error) {
	appID, _ := params.Config["app_id"].(string)
	appSecret, _ := params.Config["app_secret"].(string)

	// 获取access_token
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", params.Code)
	data.Set("client_id", appID)
	data.Set("client_secret", appSecret)
	data.Set("redirect_uri", params.Config["callback_url"].(string))

	body, err := p.httpRequest(ctx, "POST", "https://gitee.com/oauth/token", data, nil)
	if err != nil {
		return nil, err
	}

	var tokenResult map[string]interface{}
	json.Unmarshal(body, &tokenResult)
	accessToken, _ := tokenResult["access_token"].(string)

	// 获取用户信息
	userInfoURL := fmt.Sprintf("https://gitee.com/api/v5/user?access_token=%s", accessToken)
	body, err = p.httpRequest(ctx, "GET", userInfoURL, nil, nil)
	if err != nil {
		return nil, err
	}

	var userInfo map[string]interface{}
	json.Unmarshal(body, &userInfo)

	return &UserInfo{
		OpenID:   fmt.Sprintf("%.0f", userInfo["id"].(float64)),
		Username: userInfo["login"].(string),
		Avatar:   userInfo["avatar_url"].(string),
		Data:     userInfo,
	}, nil
}

// DingTalkProvider 钉钉登录
type DingTalkProvider struct {
	*BaseOAuthProvider
}

func (p *DingTalkProvider) GetLoginURL(ctx context.Context, params *LoginParams) (string, error) {
	appID, _ := params.Config["app_id"].(string)
	
	loginURL := fmt.Sprintf(
		"https://login.dingtalk.com/oauth2/auth?response_type=code&client_id=%s&redirect_uri=%s&state=%s&scope=openid",
		appID,
		url.QueryEscape(params.CallbackURL),
		params.State,
	)
	return loginURL, nil
}

func (p *DingTalkProvider) HandleCallback(ctx context.Context, params *CallbackParams) (*UserInfo, error) {
	appID, _ := params.Config["app_id"].(string)
	appSecret, _ := params.Config["app_secret"].(string)

	// 获取access_token
	data := map[string]string{
		"client_id":     appID,
		"client_secret": appSecret,
		"code":          params.Code,
		"grant_type":    "authorization_code",
	}
	
	dataBytes, _ := json.Marshal(data)
	req, _ := http.NewRequestWithContext(ctx, "POST", "https://api.dingtalk.com/v1.0/oauth2/userAccessToken", strings.NewReader(string(dataBytes)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var tokenResult map[string]interface{}
	json.Unmarshal(body, &tokenResult)
	accessToken, _ := tokenResult["accessToken"].(string)

	// 获取用户信息
	req, _ = http.NewRequestWithContext(ctx, "GET", "https://api.dingtalk.com/v1.0/contact/users/me", nil)
	req.Header.Set("x-acs-action", "GetUserByMobile")
	req.Header.Set("x-acs-action", "GetUserByAuthCode")
	req.Header.Set("x-acs-action", "GetUserInfo")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err = p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	var userInfo map[string]interface{}
	json.Unmarshal(body, &userInfo)

	return &UserInfo{
		OpenID:   userInfo["unionId"].(string),
		Username: userInfo["nick"].(string),
		Avatar:   userInfo["avatarUrl"].(string),
		Data:     userInfo,
	}, nil
}

// FeishuProvider 飞书登录
type FeishuProvider struct {
	*BaseOAuthProvider
}

func (p *FeishuProvider) GetLoginURL(ctx context.Context, params *LoginParams) (string, error) {
	appID, _ := params.Config["app_id"].(string)
	
	loginURL := fmt.Sprintf(
		"https://open.feishu.cn/open-apis/authen/v1/authorize?app_id=%s&redirect_uri=%s&state=%s",
		appID,
		url.QueryEscape(params.CallbackURL),
		params.State,
	)
	return loginURL, nil
}

func (p *FeishuProvider) HandleCallback(ctx context.Context, params *CallbackParams) (*UserInfo, error) {
	// 飞书需要先获取app_access_token，再获取user_access_token
	// 实现略复杂，这里提供框架
	return nil, fmt.Errorf("飞书登录需要完整实现，请参考飞书开放平台文档")
}

// CSDNProvider CSDN登录
type CSDNProvider struct {
	*BaseOAuthProvider
}

func (p *CSDNProvider) GetLoginURL(ctx context.Context, params *LoginParams) (string, error) {
	appID, _ := params.Config["app_id"].(string)
	
	loginURL := fmt.Sprintf(
		"https://oauth.csdn.net/authorize?response_type=code&client_id=%s&redirect_uri=%s&state=%s",
		appID,
		url.QueryEscape(params.CallbackURL),
		params.State,
	)
	return loginURL, nil
}

func (p *CSDNProvider) HandleCallback(ctx context.Context, params *CallbackParams) (*UserInfo, error) {
	// CSDN OAuth实现
	return nil, fmt.Errorf("CSDN登录暂未实现")
}

// OSChinaProvider 开源中国登录
type OSChinaProvider struct {
	*BaseOAuthProvider
}

func (p *OSChinaProvider) GetLoginURL(ctx context.Context, params *LoginParams) (string, error) {
	appID, _ := params.Config["app_id"].(string)
	
	loginURL := fmt.Sprintf(
		"https://www.oschina.net/action/oauth2/authorize?response_type=code&client_id=%s&redirect_uri=%s&state=%s",
		appID,
		url.QueryEscape(params.CallbackURL),
		params.State,
	)
	return loginURL, nil
}

func (p *OSChinaProvider) HandleCallback(ctx context.Context, params *CallbackParams) (*UserInfo, error) {
	// 开源中国OAuth实现
	return nil, fmt.Errorf("开源中国登录暂未实现")
}

// TencentCloudProvider 腾讯云登录
type TencentCloudProvider struct {
	*BaseOAuthProvider
}

func (p *TencentCloudProvider) GetLoginURL(ctx context.Context, params *LoginParams) (string, error) {
	appID, _ := params.Config["app_id"].(string)
	
	loginURL := fmt.Sprintf(
		"https://cloud.tencent.com/login/qq?state=%s&redirect_uri=%s",
		params.State,
		url.QueryEscape(params.CallbackURL),
	)
	return loginURL, nil
}

func (p *TencentCloudProvider) HandleCallback(ctx context.Context, params *CallbackParams) (*UserInfo, error) {
	// 腾讯云登录实现
	return nil, fmt.Errorf("腾讯云登录暂未实现")
}

// AliyunProvider 阿里云登录
type AliyunProvider struct {
	*BaseOAuthProvider
}

func (p *AliyunProvider) GetLoginURL(ctx context.Context, params *LoginParams) (string, error) {
	appID, _ := params.Config["app_id"].(string)
	
	loginURL := fmt.Sprintf(
		"https://passport.aliyun.com/mini_login.htm?oauthType=authorize&clientId=%s&redirectUri=%s&state=%s",
		appID,
		url.QueryEscape(params.CallbackURL),
		params.State,
	)
	return loginURL, nil
}

func (p *AliyunProvider) HandleCallback(ctx context.Context, params *CallbackParams) (*UserInfo, error) {
	// 阿里云登录实现
	return nil, fmt.Errorf("阿里云登录暂未实现")
}
