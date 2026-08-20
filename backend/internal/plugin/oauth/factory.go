package oauth

import (
	"fmt"
	"sync"
)

// registry 提供商注册表
var (
	registry   = make(map[string]ProviderConstructor)
	registryMu sync.RWMutex
)

// RegisterProvider 注册OAuth提供商
func RegisterProvider(name string, constructor ProviderConstructor) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[name] = constructor
}

// CreateProvider 创建OAuth提供商实例
func CreateProvider(name string, config map[string]interface{}) (OAuthProvider, error) {
	registryMu.RLock()
	constructor, ok := registry[name]
	registryMu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("oauth provider %q not registered", name)
	}

	return constructor(config)
}

// ListProviders 返回所有已注册的提供商名称
func ListProviders() []string {
	registryMu.RLock()
	defer registryMu.RUnlock()

	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}

// GetAllProviderInfo 获取所有提供商信息
func GetAllProviderInfo() []map[string]interface{} {
	registryMu.RLock()
	defer registryMu.RUnlock()

	var providers []map[string]interface{}
	for name, constructor := range registry {
		provider, err := constructor(nil)
		if err != nil {
			continue
		}
		providers = append(providers, map[string]interface{}{
			"name":    name,
			"title":   provider.Title(),
			"options": provider.GetConfigOptions(),
		})
	}
	return providers
}

// GetSupportedProviders 获取支持的提供商列表（24个）
func GetSupportedProviders() []map[string]string {
	return []map[string]string{
		// 国内平台（12个）
		{"name": "wechat", "title": "微信登录", "region": "cn"},
		{"name": "qq", "title": "QQ登录", "region": "cn"},
		{"name": "weibo", "title": "微博登录", "region": "cn"},
		{"name": "alipay", "title": "支付宝登录", "region": "cn"},
		{"name": "baidu", "title": "百度登录", "region": "cn"},
		{"name": "gitee", "title": "码云登录", "region": "cn"},
		{"name": "dingtalk", "title": "钉钉登录", "region": "cn"},
		{"name": "feishu", "title": "飞书登录", "region": "cn"},
		{"name": "csdn", "title": "CSDN登录", "region": "cn"},
		{"name": "oschina", "title": "开源中国登录", "region": "cn"},
		{"name": "tencent_cloud", "title": "腾讯云登录", "region": "cn"},
		{"name": "aliyun", "title": "阿里云登录", "region": "cn"},

		// 海外平台（12个）
		{"name": "google", "title": "Google登录", "region": "us"},
		{"name": "facebook", "title": "Facebook登录", "region": "us"},
		{"name": "twitter", "title": "Twitter/X登录", "region": "us"},
		{"name": "github", "title": "GitHub登录", "region": "us"},
		{"name": "linkedin", "title": "LinkedIn登录", "region": "us"},
		{"name": "microsoft", "title": "Microsoft登录", "region": "us"},
		{"name": "apple", "title": "Apple登录", "region": "us"},
		{"name": "amazon", "title": "Amazon登录", "region": "us"},
		{"name": "discord", "title": "Discord登录", "region": "us"},
		{"name": "slack", "title": "Slack登录", "region": "us"},
		{"name": "telegram", "title": "Telegram登录", "region": "us"},
		{"name": "line", "title": "LINE登录", "region": "us"},
	}
}
