package certification

import "fmt"

// Factory 根据插件名称创建认证实例
func Factory(name, configJSON string) (Certification, error) {
	switch name {
	case "Wechat":
		return NewWechatPlugin(configJSON)
	case "Idcsmartali":
		return NewIdcsmartaliPlugin(configJSON)
	default:
		return nil, fmt.Errorf("unsupported certification plugin: %s", name)
	}
}
