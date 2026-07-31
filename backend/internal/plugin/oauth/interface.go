package oauth

import "context"

// UserInfo OAuth用户信息
type UserInfo struct {
	OpenID   string                 `json:"openid"`    // 平台唯一标识
	UnionID  string                 `json:"unionid"`   // 统一标识（微信等）
	Username string                 `json:"username"`  // 用户名
	Email    string                 `json:"email"`     // 邮箱
	Phone    string                 `json:"phone"`     // 手机号
	Avatar   string                 `json:"avatar"`    // 头像URL
	Gender   int                    `json:"gender"`    // 性别 0未知 1男 2女
	Data     map[string]interface{} `json:"data"`      // 原始数据
}

// OAuthProvider OAuth提供商接口
type OAuthProvider interface {
	// Name 提供商标识
	Name() string

	// Title 显示名称
	Title() string

	// GetLoginURL 获取登录URL
	GetLoginURL(ctx context.Context, params *LoginParams) (string, error)

	// HandleCallback 处理回调，获取用户信息
	HandleCallback(ctx context.Context, params *CallbackParams) (*UserInfo, error)

	// GetConfigOptions 获取配置选项
	GetConfigOptions() []ConfigOption
}

// LoginParams 登录参数
type LoginParams struct {
	CallbackURL string                 `json:"callback_url"`
	State       string                 `json:"state"`
	Config      map[string]interface{} `json:"config"`
}

// CallbackParams 回调参数
type CallbackParams struct {
	Code  string `json:"code"`
	State string `json:"state"`
	// Config 提供商配置
	Config map[string]interface{} `json:"config"`
}

// ConfigOption 配置选项
type ConfigOption struct {
	Type        string `json:"type"`        // text/password
	Name        string `json:"name"`        // 显示名称
	Key         string `json:"key"`         // 配置键名
	Placeholder string `json:"placeholder"` // 占位符
	Description string `json:"description"` // 描述
	Required    bool   `json:"required"`    // 是否必填
}

// ProviderConstructor 提供商构造函数
type ProviderConstructor func(config map[string]interface{}) (OAuthProvider, error)
