package service

import (
	"encoding/json"
	"errors"

	"github.com/anchorfinance/backend/internal/model"
	"github.com/anchorfinance/backend/internal/repository"
)

type PluginService struct {
	pluginRepo *repository.PluginRepository
}

func NewPluginService() *PluginService {
	return &PluginService{
		pluginRepo: repository.NewPluginRepository(),
	}
}

// GetPlugins 获取插件列表
func (s *PluginService) GetPlugins(params *model.PluginQueryParams) (*model.PluginListResponse, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 100
	}

	plugins, total, err := s.pluginRepo.List(params)
	if err != nil {
		return nil, err
	}

	return &model.PluginListResponse{
		List:  plugins,
		Total: total,
	}, nil
}

// GetPluginByID 根据ID获取插件
func (s *PluginService) GetPluginByID(id uint) (*model.Plugin, error) {
	return s.pluginRepo.FindByID(id)
}

// TogglePlugin 切换插件状态
func (s *PluginService) TogglePlugin(id uint, enabled bool) error {
	plugin, err := s.pluginRepo.FindByID(id)
	if err != nil {
		return errors.New("插件不存在")
	}

	plugin.Enabled = enabled
	return s.pluginRepo.Update(plugin)
}

// GetPluginConfig 获取插件配置
func (s *PluginService) GetPluginConfig(id uint) (map[string]interface{}, error) {
	plugin, err := s.pluginRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("插件不存在")
	}

	// 解析当前配置
	config := make(map[string]interface{})
	if plugin.Config != "" {
		json.Unmarshal([]byte(plugin.Config), &config)
	}

	// 解析配置字段定义
	var configFields []model.PluginConfig
	if plugin.ConfigFields != "" {
		json.Unmarshal([]byte(plugin.ConfigFields), &configFields)
	}

	return map[string]interface{}{
		"config":        config,
		"config_fields": configFields,
	}, nil
}

// UpdatePluginConfig 更新插件配置
func (s *PluginService) UpdatePluginConfig(id uint, config map[string]interface{}) error {
	plugin, err := s.pluginRepo.FindByID(id)
	if err != nil {
		return errors.New("插件不存在")
	}

	// 验证必填字段
	var configFields []model.PluginConfig
	if plugin.ConfigFields != "" {
		json.Unmarshal([]byte(plugin.ConfigFields), &configFields)
	}

	for _, field := range configFields {
		if field.Required {
			if val, ok := config[field.Key]; !ok || val == nil || val == "" {
				return errors.New(field.Label + "不能为空")
			}
		}
	}

	// 保存配置
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	plugin.Config = string(configJSON)
	return s.pluginRepo.Update(plugin)
}

// InitDefaultPlugins 初始化默认插件
func (s *PluginService) InitDefaultPlugins() error {
	defaultPlugins := []model.Plugin{
		// OAuth插件
		{
			Name:        "微信登录",
			Code:        "wechat_login",
			Description: "支持微信扫码登录",
			Version:     "1.0.0",
			Author:      "AnchorFinance",
			Category:    "oauth",
			ConfigFields: `[{"key":"app_id","label":"App ID","type":"text","required":true},{"key":"app_secret","label":"App Secret","type":"password","required":true}]`,
		},
		{
			Name:        "QQ登录",
			Code:        "qq_login",
			Description: "支持QQ账号登录",
			Version:     "1.0.0",
			Author:      "AnchorFinance",
			Category:    "oauth",
			ConfigFields: `[{"key":"app_id","label":"App ID","type":"text","required":true},{"key":"app_secret","label":"App Secret","type":"password","required":true}]`,
		},
		// 支付插件
		{
			Name:        "支付宝",
			Code:        "alipay",
			Description: "支付宝在线支付",
			Version:     "1.0.0",
			Author:      "AnchorFinance",
			Category:    "payment",
			ConfigFields: `[{"key":"app_id","label":"应用ID","type":"text","required":true},{"key":"private_key","label":"应用私钥","type":"textarea","required":true},{"key":"public_key","label":"支付宝公钥","type":"textarea","required":true}]`,
		},
		{
			Name:        "微信支付",
			Code:        "wechat_pay",
			Description: "微信在线支付",
			Version:     "1.0.0",
			Author:      "AnchorFinance",
			Category:    "payment",
			ConfigFields: `[{"key":"mch_id","label":"商户号","type":"text","required":true},{"key":"api_key","label":"API密钥","type":"password","required":true},{"key":"cert_path","label":"证书路径","type":"text"}]`,
		},
		// 实名认证插件
		{
			Name:        "阿里云实名认证",
			Code:        "aliyun_verification",
			Description: "阿里云实名认证服务",
			Version:     "1.0.0",
			Author:      "AnchorFinance",
			Category:    "verification",
			ConfigFields: `[{"key":"access_key","label":"AccessKey","type":"text","required":true},{"key":"access_secret","label":"AccessSecret","type":"password","required":true}]`,
		},
		// 其他插件
		{
			Name:        "阿里云短信",
			Code:        "aliyun_sms",
			Description: "阿里云短信服务",
			Version:     "1.0.0",
			Author:      "AnchorFinance",
			Category:    "other",
			ConfigFields: `[{"key":"access_key","label":"AccessKey","type":"text","required":true},{"key":"access_secret","label":"AccessSecret","type":"password","required":true},{"key":"sign_name","label":"签名","type":"text","required":true}]`,
		},
	}

	for _, plugin := range defaultPlugins {
		existing, _ := s.pluginRepo.FindByCode(plugin.Code)
		if existing == nil {
			s.pluginRepo.Create(&plugin)
		}
	}

	return nil
}
