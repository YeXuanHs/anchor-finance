package model

import (
	"gorm.io/gorm"
)

// CaptchaConfig 验证码配置模型
type CaptchaConfig struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Key    string `gorm:"uniqueIndex;size:50" json:"key"` // 配置键
	Value  string `gorm:"size:255" json:"value"`           // 配置值
	Status bool   `gorm:"default:true" json:"status"`     // 启用状态
}

func (CaptchaConfig) TableName() string {
	return "captcha_configs"
}

// CaptchaType 验证码类型
const (
	CaptchaTypeImage   = "image"   // 图形验证码
	CaptchaTypeGeetest = "geetest" // 极验验证码
)

// DefaultCaptchaConfigs 默认验证码配置
var DefaultCaptchaConfigs = []CaptchaConfig{
	// 基础配置
	{Key: "is_captcha", Value: "1", Status: true},              // 验证码总开关
	{Key: "captcha_type", Value: "image", Status: true},        // 验证码类型: image/geetest
	{Key: "captcha_length", Value: "4", Status: true},          // 验证码长度（图形验证码用）
	{Key: "captcha_combination", Value: "number", Status: true}, // 验证码组合: number/letter/mixed

	// Geetest 4.0 配置
	{Key: "geetest_captcha_id", Value: "", Status: true},  // Geetest Captcha ID
	{Key: "geetest_captcha_key", Value: "", Status: true}, // Geetest Captcha Key (私钥)

	// 注册场景
	{Key: "allow_register_email_captcha", Value: "1", Status: true}, // 邮件注册显示验证码
	{Key: "allow_register_phone_captcha", Value: "1", Status: true}, // 手机注册显示验证码

	// 登录场景
	{Key: "allow_login_phone_captcha", Value: "1", Status: true},   // 手机登录显示验证码
	{Key: "allow_login_email_captcha", Value: "1", Status: true},   // 邮件登录显示验证码
	{Key: "allow_login_code_captcha", Value: "1", Status: true},    // 验证码登录显示验证码
	{Key: "allow_login_id_captcha", Value: "1", Status: true},      // ID登录显示验证码
	{Key: "allow_login_admin_captcha", Value: "1", Status: true},   // 后台登录显示验证码

	// 密码场景
	{Key: "allow_phone_forgetpwd_captcha", Value: "1", Status: true},  // 手机忘记密码显示验证码
	{Key: "allow_email_forgetpwd_captcha", Value: "1", Status: true},  // 邮件忘记密码显示验证码
	{Key: "allow_resetpwd_captcha", Value: "1", Status: true},         // 重置密码显示验证码
	{Key: "allow_setpwd_captcha", Value: "1", Status: true},           // 设置密码显示验证码

	// 绑定场景
	{Key: "allow_phone_bind_captcha", Value: "1", Status: true},   // 手机绑定显示验证码
	{Key: "allow_email_bind_captcha", Value: "1", Status: true},   // 邮件绑定显示验证码

	// 其他场景
	{Key: "allow_cancel_sms_captcha", Value: "1", Status: true},   // 取消短信提醒显示验证码
	{Key: "allow_cancel_email_captcha", Value: "1", Status: true}, // 取消邮件提醒显示验证码
}

// CaptchaConfigService 验证码配置服务
type CaptchaConfigService struct {
	db *gorm.DB
}

func NewCaptchaConfigService(db *gorm.DB) *CaptchaConfigService {
	return &CaptchaConfigService{db: db}
}

// InitDefaultConfigs 初始化默认配置
func (s *CaptchaConfigService) InitDefaultConfigs() error {
	for _, config := range DefaultCaptchaConfigs {
		var existing CaptchaConfig
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
func (s *CaptchaConfigService) GetAllConfigs() ([]CaptchaConfig, error) {
	var configs []CaptchaConfig
	if err := s.db.Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

// GetConfig 获取单个配置
func (s *CaptchaConfigService) GetConfig(key string) (*CaptchaConfig, error) {
	var config CaptchaConfig
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

// GetCaptchaType 获取验证码类型
func (s *CaptchaConfigService) GetCaptchaType() string {
	captchaType := s.GetValue("captcha_type")
	if captchaType == "" {
		return CaptchaTypeImage
	}
	return captchaType
}

// IsGeetestEnabled 检查是否启用极验
func (s *CaptchaConfigService) IsGeetestEnabled() bool {
	return s.GetCaptchaType() == CaptchaTypeGeetest
}

// GetGeetestConfig 获取极验配置
func (s *CaptchaConfigService) GetGeetestConfig() (captchaID, captchaKey string) {
	captchaID = s.GetValue("geetest_captcha_id")
	captchaKey = s.GetValue("geetest_captcha_key")
	return
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
	return s.db.Model(&CaptchaConfig{}).Where("key = ?", key).Updates(map[string]interface{}{
		"value":  value,
		"status": status,
	}).Error
}

// BatchUpdateConfigs 批量更新配置
func (s *CaptchaConfigService) BatchUpdateConfigs(configs map[string]struct {
	Value  string `json:"value"`
	Status bool   `json:"status"`
}) error {
	tx := s.db.Begin()
	for key, val := range configs {
		if err := tx.Model(&CaptchaConfig{}).Where("key = ?", key).Updates(map[string]interface{}{
			"value":  val.Value,
			"status": val.Status,
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
