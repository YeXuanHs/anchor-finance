package service

import (
	"anchorfinance/internal/model"

	"gorm.io/gorm"
)

// CaptchaConfigService 验证码配置服务
type CaptchaConfigService struct {
	db *gorm.DB
}

func NewCaptchaConfigService(db *gorm.DB) *CaptchaConfigService {
	return &CaptchaConfigService{db: db}
}

// InitDefaultConfigs 初始化默认配置
func (s *CaptchaConfigService) InitDefaultConfigs() error {
	return s.initDefaultConfigs()
}

func (s *CaptchaConfigService) initDefaultConfigs() error {
	for _, config := range model.DefaultCaptchaConfigs {
		var existing model.CaptchaConfig
		if err := s.db.Where("key = ?", config.Key).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := s.db.Create(&config).Error; err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// GetAllConfigs 获取所有配置
func (s *CaptchaConfigService) GetAllConfigs() ([]model.CaptchaConfig, error) {
	var configs []model.CaptchaConfig
	if err := s.db.Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

// GetConfig 获取单个配置
func (s *CaptchaConfigService) GetConfig(key string) (*model.CaptchaConfig, error) {
	var config model.CaptchaConfig
	if err := s.db.Where("key = ?", key).First(&config).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

// GetValue 获取配置值
func (s *CaptchaConfigService) GetValue(key string) string {
	config, err := s.GetConfig(key)
	if err != nil {
		return ""
	}
	return config.Value
}

// IsEnabled 检查某个配置是否启用
func (s *CaptchaConfigService) IsEnabled(key string) bool {
	config, err := s.GetConfig(key)
	if err != nil {
		return false
	}
	return config.Status && config.Value == "1"
}

// IsCaptchaEnabled 检查验证码是否全局启用
func (s *CaptchaConfigService) IsCaptchaEnabled() bool {
	return s.IsEnabled("is_captcha")
}

// ShouldShowCaptcha 检查某个场景是否应该显示验证码
func (s *CaptchaConfigService) ShouldShowCaptcha(scene string) bool {
	if !s.IsCaptchaEnabled() {
		return false
	}
	return s.IsEnabled(scene)
}

// GetCaptchaLength 获取验证码长度
func (s *CaptchaConfigService) GetCaptchaLength() int {
	length := s.GetValue("captcha_length")
	switch length {
	case "4":
		return 4
	case "5":
		return 5
	case "6":
		return 6
	default:
		return 4
	}
}

// GetCaptchaCombination 获取验证码组合方式
func (s *CaptchaConfigService) GetCaptchaCombination() string {
	combination := s.GetValue("captcha_combination")
	if combination == "" {
		return "number"
	}
	return combination
}

// UpdateConfig 更新配置
func (s *CaptchaConfigService) UpdateConfig(key, value string, status bool) error {
	return s.db.Model(&model.CaptchaConfig{}).Where("key = ?", key).Updates(map[string]interface{}{
		"value":  value,
		"status": status,
	}).Error
}

// BatchUpdateConfigs 批量更新配置
func (s *CaptchaConfigService) BatchUpdateConfigs(configs []struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Status bool   `json:"status"`
}) error {
	tx := s.db.Begin()
	for _, config := range configs {
		if err := tx.Model(&model.CaptchaConfig{}).Where("key = ?", config.Key).Updates(map[string]interface{}{
			"value":  config.Value,
			"status": config.Status,
		}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

// GetSceneConfig 获取场景配置（用于前端）
func (s *CaptchaConfigService) GetSceneConfig() map[string]bool {
	scenes := []string{
		"allow_register_email_captcha",
		"allow_register_phone_captcha",
		"allow_login_phone_captcha",
		"allow_login_email_captcha",
		"allow_login_code_captcha",
		"allow_login_id_captcha",
		"allow_login_admin_captcha",
		"allow_phone_forgetpwd_captcha",
		"allow_email_forgetpwd_captcha",
		"allow_resetpwd_captcha",
		"allow_setpwd_captcha",
		"allow_phone_bind_captcha",
		"allow_email_bind_captcha",
		"allow_cancel_sms_captcha",
		"allow_cancel_email_captcha",
	}

	result := make(map[string]bool)
	for _, scene := range scenes {
		result[scene] = s.ShouldShowCaptcha(scene)
	}
	return result
}

// GetBasicConfig 获取基础配置
func (s *CaptchaConfigService) GetBasicConfig() map[string]interface{} {
	return map[string]interface{}{
		"is_captcha":          s.IsCaptchaEnabled(),
		"captcha_length":      s.GetCaptchaLength(),
		"captcha_combination": s.GetCaptchaCombination(),
	}
}

// GetPublicCaptchaConfig 获取公开的验证码配置（用于前端）
func (s *CaptchaConfigService) GetPublicCaptchaConfig() map[string]interface{} {
	config := map[string]interface{}{
		"enabled": s.IsCaptchaEnabled(),
		"type":    s.GetCaptchaType(),
		"scenes":  s.GetSceneConfig(),
	}

	// 如果是极验，返回极验配置
	if s.IsGeetestEnabled() {
		captchaID, _ := s.GetGeetestConfig()
		config["geetest"] = map[string]interface{}{
			"captcha_id": captchaID,
		}
	}

	return config
}

// GetCaptchaType 获取验证码类型
func (s *CaptchaConfigService) GetCaptchaType() string {
	return s.GetValue("captcha_type")
}

// IsGeetestEnabled 检查极验是否启用
func (s *CaptchaConfigService) IsGeetestEnabled() bool {
	return s.GetCaptchaType() == "geetest"
}

// GetGeetestConfig 获取极验配置
func (s *CaptchaConfigService) GetGeetestConfig() (string, string) {
	return s.GetValue("geetest_captcha_id"), s.GetValue("geetest_captcha_key")
}
