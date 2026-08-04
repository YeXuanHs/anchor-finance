package service

import (
	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/email"
	"encoding/json"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type ConfigGeneralService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewConfigGeneralService(db *gorm.DB, log *logger.Logger) *ConfigGeneralService {
	return &ConfigGeneralService{db: db, log: log}
}

// ==================== General Config ====================

type GeneralConfig struct {
	SiteName       string `json:"site_name"`
	SiteURL        string `json:"site_url"`
	Logo           string `json:"logo"`
	Favicon        string `json:"favicon"`
	Description    string `json:"description"`
	Keywords       string `json:"keywords"`
	ICP            string `json:"icp"`
	PSB            string `json:"psb"`
	Copyright      string `json:"copyright"`
	ContactEmail   string `json:"contact_email"`
	ContactPhone   string `json:"contact_phone"`
	ContactAddress string `json:"contact_address"`
	CompanyName    string `json:"company_name"`
	CompanyLogo    string `json:"company_logo"`
	TermsURL       string `json:"terms_url"`
	PrivacyURL     string `json:"privacy_url"`
	HomepageTitle  string `json:"homepage_title"`
	HomepageDesc   string `json:"homepage_desc"`
	OpenRegister   bool   `json:"open_register"`
	VerifyEmail    bool   `json:"verify_email"`
	DefaultLang    string `json:"default_lang"`
	DefaultTheme   string `json:"default_theme"`
	CustomCSS      string `json:"custom_css"`
	CustomJS       string `json:"custom_js"`
	FooterHTML     string `json:"footer_html"`
}

var generalConfigKeys = []string{
	"site_name", "site_url", "logo", "favicon", "description", "keywords",
	"icp", "psb", "copyright", "contact_email", "contact_phone", "contact_address",
	"company_name", "company_logo", "terms_url", "privacy_url",
	"homepage_title", "homepage_desc", "open_register", "verify_email",
	"default_lang", "default_theme", "custom_css", "custom_js", "footer_html",
}

func (s *ConfigGeneralService) GetConfig() (*GeneralConfig, error) {
	var configs []model.SystemConfig
	if err := s.db.Where("`key` IN ?", generalConfigKeys).Find(&configs).Error; err != nil {
		return nil, err
	}

	configMap := make(map[string]string, len(configs))
	for _, c := range configs {
		configMap[c.Key] = c.Value
	}

	return &GeneralConfig{
		SiteName:       configMap["site_name"],
		SiteURL:        configMap["site_url"],
		Logo:           configMap["logo"],
		Favicon:        configMap["favicon"],
		Description:    configMap["description"],
		Keywords:       configMap["keywords"],
		ICP:            configMap["icp"],
		PSB:            configMap["psb"],
		Copyright:      configMap["copyright"],
		ContactEmail:   configMap["contact_email"],
		ContactPhone:   configMap["contact_phone"],
		ContactAddress: configMap["contact_address"],
		CompanyName:    configMap["company_name"],
		CompanyLogo:    configMap["company_logo"],
		TermsURL:       configMap["terms_url"],
		PrivacyURL:     configMap["privacy_url"],
		HomepageTitle:  configMap["homepage_title"],
		HomepageDesc:   configMap["homepage_desc"],
		OpenRegister:   configMap["open_register"] == "true",
		VerifyEmail:    configMap["verify_email"] == "true",
		DefaultLang:    configMap["default_lang"],
		DefaultTheme:   configMap["default_theme"],
		CustomCSS:      configMap["custom_css"],
		CustomJS:       configMap["custom_js"],
		FooterHTML:     configMap["footer_html"],
	}, nil
}

func (s *ConfigGeneralService) UpdateConfig(req GeneralConfig) error {
	configs := map[string]string{
		"site_name":        req.SiteName,
		"site_url":         req.SiteURL,
		"logo":             req.Logo,
		"favicon":          req.Favicon,
		"description":      req.Description,
		"keywords":         req.Keywords,
		"icp":              req.ICP,
		"psb":              req.PSB,
		"copyright":        req.Copyright,
		"contact_email":    req.ContactEmail,
		"contact_phone":    req.ContactPhone,
		"contact_address":  req.ContactAddress,
		"company_name":     req.CompanyName,
		"company_logo":     req.CompanyLogo,
		"terms_url":        req.TermsURL,
		"privacy_url":      req.PrivacyURL,
		"homepage_title":   req.HomepageTitle,
		"homepage_desc":    req.HomepageDesc,
		"open_register":    BoolStr(req.OpenRegister),
		"verify_email":     BoolStr(req.VerifyEmail),
		"default_lang":     req.DefaultLang,
		"default_theme":    req.DefaultTheme,
		"custom_css":       req.CustomCSS,
		"custom_js":        req.CustomJS,
		"footer_html":      req.FooterHTML,
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		for key, value := range configs {
			result := tx.Where("key = ?", key).
				Assign(model.SystemConfig{Value: value}).
				FirstOrCreate(&model.SystemConfig{
					Key:   key,
					Value: value,
					Group: "general",
					Name:  key,
					Type:  "string",
				})
			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})
}

func BoolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// ==================== Email Config ====================

type EmailConfig struct {
	SMTPHost     string `json:"smtp_host"`
	SMTPPort     int    `json:"smtp_port"`
	SMTPUser     string `json:"smtp_user"`
	SMTPPassword string `json:"smtp_password"`
	FromName     string `json:"from_name"`
	FromAddress  string `json:"from_address"`
	Encryption   string `json:"encryption"` // none/ssl/tls
}

var emailConfigKeys = []string{
	"email_smtp_host", "email_smtp_port", "email_smtp_user",
	"email_smtp_password", "email_from_name", "email_from_address",
	"email_encryption",
}

func (s *ConfigGeneralService) GetEmailConfig() (*EmailConfig, error) {
	var configs []model.SystemConfig
	if err := s.db.Where("`key` IN ?", emailConfigKeys).Find(&configs).Error; err != nil {
		return nil, err
	}

	m := make(map[string]string, len(configs))
	for _, c := range configs {
		m[c.Key] = c.Value
	}

	return &EmailConfig{
		SMTPHost:     m["email_smtp_host"],
		SMTPPort:     parseInt(m["email_smtp_port"]),
		SMTPUser:     m["email_smtp_user"],
		SMTPPassword: m["email_smtp_password"],
		FromName:     m["email_from_name"],
		FromAddress:  m["email_from_address"],
		Encryption:   m["email_encryption"],
	}, nil
}

func (s *ConfigGeneralService) UpdateEmailConfig(req EmailConfig) error {
	configs := map[string]string{
		"email_smtp_host":     req.SMTPHost,
		"email_smtp_port":     intStr(req.SMTPPort),
		"email_smtp_user":     req.SMTPUser,
		"email_smtp_password": req.SMTPPassword,
		"email_from_name":     req.FromName,
		"email_from_address":  req.FromAddress,
		"email_encryption":    req.Encryption,
	}
	return s.saveConfigMap(configs, "email")
}

// ==================== Email Support ====================

type EmailSupportConfig struct {
	SupportEmail    string            `json:"support_email"`
	DepartmentEmails map[string]string `json:"department_emails"` // dept_name -> email
}

var emailSupportKeys = []string{
	"email_support_address", "email_support_departments",
}

func (s *ConfigGeneralService) GetEmailSupport() (*EmailSupportConfig, error) {
	var configs []model.SystemConfig
	if err := s.db.Where("`key` IN ?", emailSupportKeys).Find(&configs).Error; err != nil {
		return nil, err
	}

	m := make(map[string]string, len(configs))
	for _, c := range configs {
		m[c.Key] = c.Value
	}

	depts := make(map[string]string)
	if v := m["email_support_departments"]; v != "" {
		_ = json.Unmarshal([]byte(v), &depts)
	}

	return &EmailSupportConfig{
		SupportEmail:     m["email_support_address"],
		DepartmentEmails: depts,
	}, nil
}

func (s *ConfigGeneralService) UpdateEmailSupport(req EmailSupportConfig) error {
	deptsJSON, _ := json.Marshal(req.DepartmentEmails)
	configs := map[string]string{
		"email_support_address":     req.SupportEmail,
		"email_support_departments": string(deptsJSON),
	}
	return s.saveConfigMap(configs, "email")
}

// ==================== Affiliate Ladders ====================

type AffiliateLadder struct {
	MinAmount float64 `json:"min_amount"`
	MaxAmount float64 `json:"max_amount"`
	Rate      float64 `json:"rate"` // percent
}

type AffiliateLaddersConfig struct {
	Ladders []AffiliateLadder `json:"ladders"`
}

func (s *ConfigGeneralService) GetAffiliateLadders() (*AffiliateLaddersConfig, error) {
	var cfg model.SystemConfig
	if err := s.db.Where("`key` = ?", "affiliate_ladders").First(&cfg).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &AffiliateLaddersConfig{Ladders: []AffiliateLadder{}}, nil
		}
		return nil, err
	}

	var ladders []AffiliateLadder
	if cfg.Value != "" {
		_ = json.Unmarshal([]byte(cfg.Value), &ladders)
	}
	return &AffiliateLaddersConfig{Ladders: ladders}, nil
}

func (s *ConfigGeneralService) UpdateAffiliateLadders(req AffiliateLaddersConfig) error {
	data, _ := json.Marshal(req.Ladders)
	return s.saveSingleKey("affiliate_ladders", string(data), "affiliate")
}

// ==================== Safe Config ====================

type SafeConfig struct {
	MaxLoginAttempts int      `json:"max_login_attempts"`
	LockMinutes      int      `json:"lock_minutes"`
	IPWhitelist      []string `json:"ip_whitelist"`
	Force2FA         bool     `json:"force_2fa"`
	CaptchaLogin     bool     `json:"captcha_login"`
	CaptchaRegister  bool     `json:"captcha_register"`
}

var safeConfigKeys = []string{
	"safe_max_login_attempts", "safe_lock_minutes",
	"safe_ip_whitelist", "safe_force_2fa",
	"safe_captcha_login", "safe_captcha_register",
}

func (s *ConfigGeneralService) GetSafeConfig() (*SafeConfig, error) {
	var configs []model.SystemConfig
	if err := s.db.Where("`key` IN ?", safeConfigKeys).Find(&configs).Error; err != nil {
		return nil, err
	}

	m := make(map[string]string, len(configs))
	for _, c := range configs {
		m[c.Key] = c.Value
	}

	var whitelist []string
	if v := m["safe_ip_whitelist"]; v != "" {
		_ = json.Unmarshal([]byte(v), &whitelist)
	}

	return &SafeConfig{
		MaxLoginAttempts: parseInt(m["safe_max_login_attempts"]),
		LockMinutes:      parseInt(m["safe_lock_minutes"]),
		IPWhitelist:      whitelist,
		Force2FA:         m["safe_force_2fa"] == "true",
		CaptchaLogin:     m["safe_captcha_login"] == "true",
		CaptchaRegister:  m["safe_captcha_register"] == "true",
	}, nil
}

func (s *ConfigGeneralService) UpdateSafeConfig(req SafeConfig) error {
	wlJSON, _ := json.Marshal(req.IPWhitelist)
	configs := map[string]string{
		"safe_max_login_attempts": intStr(req.MaxLoginAttempts),
		"safe_lock_minutes":       intStr(req.LockMinutes),
		"safe_ip_whitelist":       string(wlJSON),
		"safe_force_2fa":          BoolStr(req.Force2FA),
		"safe_captcha_login":      BoolStr(req.CaptchaLogin),
		"safe_captcha_register":   BoolStr(req.CaptchaRegister),
	}
	return s.saveConfigMap(configs, "safe")
}

// ==================== Recharge Config ====================

type RechargeConfig struct {
	MinAmount    float64  `json:"min_amount"`
	MaxAmount    float64  `json:"max_amount"`
	FixedAmounts []float64 `json:"fixed_amounts"`
	AutoApprove  bool     `json:"auto_approve"`
}

var rechargeConfigKeys = []string{
	"recharge_min_amount", "recharge_max_amount",
	"recharge_fixed_amounts", "recharge_auto_approve",
}

func (s *ConfigGeneralService) GetRechargeConfig() (*RechargeConfig, error) {
	var configs []model.SystemConfig
	if err := s.db.Where("`key` IN ?", rechargeConfigKeys).Find(&configs).Error; err != nil {
		return nil, err
	}

	m := make(map[string]string, len(configs))
	for _, c := range configs {
		m[c.Key] = c.Value
	}

	var fixed []float64
	if v := m["recharge_fixed_amounts"]; v != "" {
		_ = json.Unmarshal([]byte(v), &fixed)
	}

	return &RechargeConfig{
		MinAmount:    parseFloat(m["recharge_min_amount"]),
		MaxAmount:    parseFloat(m["recharge_max_amount"]),
		FixedAmounts: fixed,
		AutoApprove:  m["recharge_auto_approve"] == "true",
	}, nil
}

func (s *ConfigGeneralService) UpdateRechargeConfig(req RechargeConfig) error {
	fixedJSON, _ := json.Marshal(req.FixedAmounts)
	configs := map[string]string{
		"recharge_min_amount":    fmt.Sprintf("%.2f", req.MinAmount),
		"recharge_max_amount":    fmt.Sprintf("%.2f", req.MaxAmount),
		"recharge_fixed_amounts": string(fixedJSON),
		"recharge_auto_approve":  BoolStr(req.AutoApprove),
	}
	return s.saveConfigMap(configs, "recharge")
}

// ==================== Invoice Config ====================

type InvoiceConfig struct {
	InvoicePrefix  string  `json:"invoice_prefix"`
	AutoGenerate   bool    `json:"auto_generate"`
	TaxRate        float64 `json:"tax_rate"`
	TitleRequired  bool    `json:"title_required"`
}

var invoiceConfigKeys = []string{
	"invoice_prefix", "invoice_auto_generate",
	"invoice_tax_rate", "invoice_title_required",
}

func (s *ConfigGeneralService) GetInvoiceConfig() (*InvoiceConfig, error) {
	var configs []model.SystemConfig
	if err := s.db.Where("`key` IN ?", invoiceConfigKeys).Find(&configs).Error; err != nil {
		return nil, err
	}

	m := make(map[string]string, len(configs))
	for _, c := range configs {
		m[c.Key] = c.Value
	}

	return &InvoiceConfig{
		InvoicePrefix: m["invoice_prefix"],
		AutoGenerate:  m["invoice_auto_generate"] == "true",
		TaxRate:       parseFloat(m["invoice_tax_rate"]),
		TitleRequired: m["invoice_title_required"] == "true",
	}, nil
}

func (s *ConfigGeneralService) UpdateInvoiceConfig(req InvoiceConfig) error {
	configs := map[string]string{
		"invoice_prefix":         req.InvoicePrefix,
		"invoice_auto_generate":  BoolStr(req.AutoGenerate),
		"invoice_tax_rate":       fmt.Sprintf("%.2f", req.TaxRate),
		"invoice_title_required": BoolStr(req.TitleRequired),
	}
	return s.saveConfigMap(configs, "invoice")
}

// ==================== Register Config ====================

type RegisterConfig struct {
	EnableRegister bool     `json:"enable_register"`
	EmailVerify    bool     `json:"email_verify"`
	PhoneVerify    bool     `json:"phone_verify"`
	Captcha        bool     `json:"captcha"`
	ShowFields     []string `json:"show_fields"` // username/email/phone/real_name/company
}

var registerConfigKeys = []string{
	"register_enable", "register_email_verify",
	"register_phone_verify", "register_captcha",
	"register_show_fields",
}

func (s *ConfigGeneralService) GetRegisterConfig() (*RegisterConfig, error) {
	var configs []model.SystemConfig
	if err := s.db.Where("`key` IN ?", registerConfigKeys).Find(&configs).Error; err != nil {
		return nil, err
	}

	m := make(map[string]string, len(configs))
	for _, c := range configs {
		m[c.Key] = c.Value
	}

	var fields []string
	if v := m["register_show_fields"]; v != "" {
		_ = json.Unmarshal([]byte(v), &fields)
	}

	return &RegisterConfig{
		EnableRegister: m["register_enable"] == "true",
		EmailVerify:    m["register_email_verify"] == "true",
		PhoneVerify:    m["register_phone_verify"] == "true",
		Captcha:        m["register_captcha"] == "true",
		ShowFields:     fields,
	}, nil
}

func (s *ConfigGeneralService) UpdateRegisterConfig(req RegisterConfig) error {
	fieldsJSON, _ := json.Marshal(req.ShowFields)
	configs := map[string]string{
		"register_enable":       BoolStr(req.EnableRegister),
		"register_email_verify": BoolStr(req.EmailVerify),
		"register_phone_verify": BoolStr(req.PhoneVerify),
		"register_captcha":      BoolStr(req.Captcha),
		"register_show_fields":  string(fieldsJSON),
	}
	return s.saveConfigMap(configs, "register")
}

// ==================== Login Config ====================

type LoginConfig struct {
	Methods    []string `json:"methods"`     // email/phone/username
	Captcha    bool     `json:"captcha"`
	RememberMe bool     `json:"remember_me"`
}

var loginConfigKeys = []string{
	"login_methods", "login_captcha", "login_remember_me",
}

func (s *ConfigGeneralService) GetLoginConfig() (*LoginConfig, error) {
	var configs []model.SystemConfig
	if err := s.db.Where("`key` IN ?", loginConfigKeys).Find(&configs).Error; err != nil {
		return nil, err
	}

	m := make(map[string]string, len(configs))
	for _, c := range configs {
		m[c.Key] = c.Value
	}

	var methods []string
	if v := m["login_methods"]; v != "" {
		_ = json.Unmarshal([]byte(v), &methods)
	}

	return &LoginConfig{
		Methods:    methods,
		Captcha:    m["login_captcha"] == "true",
		RememberMe: m["login_remember_me"] == "true",
	}, nil
}

func (s *ConfigGeneralService) UpdateLoginConfig(req LoginConfig) error {
	methodsJSON, _ := json.Marshal(req.Methods)
	configs := map[string]string{
		"login_methods":     string(methodsJSON),
		"login_captcha":     BoolStr(req.Captcha),
		"login_remember_me": BoolStr(req.RememberMe),
	}
	return s.saveConfigMap(configs, "login")
}

// ==================== API Config ====================

type APIConfig struct {
	EnableAPI  bool     `json:"enable_api"`
	RateLimit  int      `json:"rate_limit"` // requests per minute
	AllowedIPs []string `json:"allowed_ips"`
}

var apiConfigKeys = []string{
	"api_enable", "api_rate_limit", "api_allowed_ips",
}

func (s *ConfigGeneralService) GetAPIConfig() (*APIConfig, error) {
	var configs []model.SystemConfig
	if err := s.db.Where("`key` IN ?", apiConfigKeys).Find(&configs).Error; err != nil {
		return nil, err
	}

	m := make(map[string]string, len(configs))
	for _, c := range configs {
		m[c.Key] = c.Value
	}

	var ips []string
	if v := m["api_allowed_ips"]; v != "" {
		_ = json.Unmarshal([]byte(v), &ips)
	}

	return &APIConfig{
		EnableAPI:  m["api_enable"] == "true",
		RateLimit:  parseInt(m["api_rate_limit"]),
		AllowedIPs: ips,
	}, nil
}

func (s *ConfigGeneralService) UpdateAPIConfig(req APIConfig) error {
	ipsJSON, _ := json.Marshal(req.AllowedIPs)
	configs := map[string]string{
		"api_enable":      BoolStr(req.EnableAPI),
		"api_rate_limit":  intStr(req.RateLimit),
		"api_allowed_ips": string(ipsJSON),
	}
	return s.saveConfigMap(configs, "api")
}

// ==================== 2FA Config ====================

type TwoFactorConfig struct {
	Enable      bool     `json:"enable"`
	Methods     []string `json:"methods"`       // totp/sms/email
	ForcedRoles []string `json:"forced_roles"`   // role names that must use 2FA
}

var twoFactorConfigKeys = []string{
	"twofactor_enable", "twofactor_methods", "twofactor_forced_roles",
}

func (s *ConfigGeneralService) GetTwoFactorConfig() (*TwoFactorConfig, error) {
	var configs []model.SystemConfig
	if err := s.db.Where("`key` IN ?", twoFactorConfigKeys).Find(&configs).Error; err != nil {
		return nil, err
	}

	m := make(map[string]string, len(configs))
	for _, c := range configs {
		m[c.Key] = c.Value
	}

	var methods []string
	if v := m["twofactor_methods"]; v != "" {
		_ = json.Unmarshal([]byte(v), &methods)
	}
	var roles []string
	if v := m["twofactor_forced_roles"]; v != "" {
		_ = json.Unmarshal([]byte(v), &roles)
	}

	return &TwoFactorConfig{
		Enable:      m["twofactor_enable"] == "true",
		Methods:     methods,
		ForcedRoles: roles,
	}, nil
}

func (s *ConfigGeneralService) UpdateTwoFactorConfig(req TwoFactorConfig) error {
	methodsJSON, _ := json.Marshal(req.Methods)
	rolesJSON, _ := json.Marshal(req.ForcedRoles)
	configs := map[string]string{
		"twofactor_enable":       BoolStr(req.Enable),
		"twofactor_methods":      string(methodsJSON),
		"twofactor_forced_roles": string(rolesJSON),
	}
	return s.saveConfigMap(configs, "twofactor")
}

// ==================== Debug Mode ====================

func (s *ConfigGeneralService) GetDebugMode() (bool, error) {
	var cfg model.SystemConfig
	if err := s.db.Where("`key` = ?", "debug_mode").First(&cfg).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return cfg.Value == "true", nil
}

func (s *ConfigGeneralService) SetDebugMode(enabled bool) error {
	return s.saveSingleKey("debug_mode", BoolStr(enabled), "system")
}

// ==================== SMTP Test ====================

func (s *ConfigGeneralService) TestSMTP(toEmail string) error {
	if toEmail == "" {
		return fmt.Errorf("recipient email is required")
	}

	sender := email.NewSender(s.db)
	return sender.Send(toEmail, "AnchorFinance SMTP Test", "<p>This is a test email from AnchorFinance. SMTP configuration is working correctly.</p>")
}

// ==================== SMS Test ====================

func (s *ConfigGeneralService) TestSMS(phone string) error {
	if phone == "" {
		return fmt.Errorf("phone number is required")
	}

	// Load SMS config to verify it's configured
	cfg, err := s.GetSmsConfig()
	if err != nil {
		return fmt.Errorf("failed to load SMS config: %w", err)
	}
	if cfg.Provider == "" {
		return fmt.Errorf("SMS provider not configured")
	}
	if cfg.AccessKeyID == "" || cfg.AccessSecret == "" {
		return fmt.Errorf("SMS access credentials not configured")
	}
	if cfg.SignName == "" {
		return fmt.Errorf("SMS signature not configured")
	}
	if cfg.TemplateCode == "" {
		return fmt.Errorf("SMS template code not configured")
	}

	s.log.Info("SMS test: configuration validated for provider=%s, phone=%s", cfg.Provider, phone)
	return nil
}

// ==================== Helper methods ====================

func (s *ConfigGeneralService) saveConfigMap(configs map[string]string, group string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for key, value := range configs {
			result := tx.Where("key = ?", key).
				Assign(model.SystemConfig{Value: value}).
				FirstOrCreate(&model.SystemConfig{
					Key:   key,
					Value: value,
					Group: group,
					Name:  key,
					Type:  "string",
				})
			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})
}

func (s *ConfigGeneralService) saveSingleKey(key, value, group string) error {
	return s.db.Where("key = ?", key).
		Assign(model.SystemConfig{Value: value}).
		FirstOrCreate(&model.SystemConfig{
			Key:   key,
			Value: value,
			Group: group,
			Name:  key,
			Type:  "string",
		}).Error
}

func parseFloat(s string) float64 {
	if s == "" {
		return 0
	}
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

func parseInt(s string) int {
	if s == "" {
		return 0
	}
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}

func intStr(i int) string {
	return fmt.Sprintf("%d", i)
}

// ==================== Payment Config ====================

type PaymentConfig struct {
	AlipayEnabled    bool   `json:"alipay_enabled"`
	WechatEnabled    bool   `json:"wechat_enabled"`
	AlipayAppID      string `json:"alipay_app_id"`
	AlipayPrivateKey string `json:"alipay_private_key"`
	WechatAppID      string `json:"wechat_app_id"`
	WechatMchID      string `json:"wechat_mch_id"`
	WechatAPIKey     string `json:"wechat_api_key"`
	AutoInvoice      bool   `json:"auto_invoice"`
}

var paymentConfigKeys = []string{
	"payment_alipay_enabled", "payment_wechat_enabled",
	"payment_alipay_app_id", "payment_alipay_private_key",
	"payment_wechat_app_id", "payment_wechat_mch_id", "payment_wechat_api_key",
	"payment_auto_invoice",
}

func (s *ConfigGeneralService) GetPaymentConfig() (*PaymentConfig, error) {
	var configs []model.SystemConfig
	if err := s.db.Where("`key` IN ?", paymentConfigKeys).Find(&configs).Error; err != nil {
		return nil, err
	}
	m := make(map[string]string, len(configs))
	for _, c := range configs {
		m[c.Key] = c.Value
	}
	return &PaymentConfig{
		AlipayEnabled:    m["payment_alipay_enabled"] == "true",
		WechatEnabled:    m["payment_wechat_enabled"] == "true",
		AlipayAppID:      m["payment_alipay_app_id"],
		AlipayPrivateKey: m["payment_alipay_private_key"],
		WechatAppID:      m["payment_wechat_app_id"],
		WechatMchID:      m["payment_wechat_mch_id"],
		WechatAPIKey:     m["payment_wechat_api_key"],
		AutoInvoice:      m["payment_auto_invoice"] == "true",
	}, nil
}

func (s *ConfigGeneralService) UpdatePaymentConfig(req PaymentConfig) error {
	configs := map[string]string{
		"payment_alipay_enabled":     BoolStr(req.AlipayEnabled),
		"payment_wechat_enabled":     BoolStr(req.WechatEnabled),
		"payment_alipay_app_id":      req.AlipayAppID,
		"payment_alipay_private_key": req.AlipayPrivateKey,
		"payment_wechat_app_id":      req.WechatAppID,
		"payment_wechat_mch_id":      req.WechatMchID,
		"payment_wechat_api_key":     req.WechatAPIKey,
		"payment_auto_invoice":       BoolStr(req.AutoInvoice),
	}
	return s.saveConfigMap(configs, "payment")
}

// ==================== SMS Config ====================

type SmsConfig struct {
	Provider       string `json:"provider"`
	AccessKeyID    string `json:"access_key_id"`
	AccessSecret   string `json:"access_secret"`
	SignName       string `json:"sign_name"`
	TemplateCode   string `json:"template_code"`
	DailyLimit     int    `json:"daily_limit"`
	PhoneLimit     int    `json:"phone_limit"`
	Enabled        bool   `json:"enabled"`
}

var smsConfigKeys = []string{
	"sms_provider", "sms_access_key_id", "sms_access_secret",
	"sms_sign_name", "sms_template_code", "sms_daily_limit", "sms_phone_limit",
	"sms_enabled",
}

func (s *ConfigGeneralService) GetSmsConfig() (*SmsConfig, error) {
	var configs []model.SystemConfig
	if err := s.db.Where("`key` IN ?", smsConfigKeys).Find(&configs).Error; err != nil {
		return nil, err
	}
	m := make(map[string]string, len(configs))
	for _, c := range configs {
		m[c.Key] = c.Value
	}
	return &SmsConfig{
		Provider:     m["sms_provider"],
		AccessKeyID:  m["sms_access_key_id"],
		AccessSecret: m["sms_access_secret"],
		SignName:     m["sms_sign_name"],
		TemplateCode: m["sms_template_code"],
		DailyLimit:   parseInt(m["sms_daily_limit"]),
		PhoneLimit:   parseInt(m["sms_phone_limit"]),
		Enabled:      m["sms_enabled"] == "true",
	}, nil
}

func (s *ConfigGeneralService) UpdateSmsConfig(req SmsConfig) error {
	configs := map[string]string{
		"sms_provider":       req.Provider,
		"sms_access_key_id":  req.AccessKeyID,
		"sms_access_secret":  req.AccessSecret,
		"sms_sign_name":      req.SignName,
		"sms_template_code":  req.TemplateCode,
		"sms_daily_limit":    intStr(req.DailyLimit),
		"sms_phone_limit":    intStr(req.PhoneLimit),
		"sms_enabled":        BoolStr(req.Enabled),
	}
	return s.saveConfigMap(configs, "sms")
}

// ==================== Security Config ====================

type SecurityConfig struct {
	RequiredPasswordStrength string `json:"required_password_strength"`
	InvalidLoginsBanLength   int    `json:"invalid_logins_ban_length"`
	IPCheckFrontend          bool   `json:"ip_check_frontend"`
	IPCheckAdmin             bool   `json:"ip_check_admin"`
	LoginErrorMaxNum         int    `json:"login_error_max_num"`
}

var securityConfigKeys = []string{
	"required_pwstrength", "invalid_logins_banlength",
	"home_ip_check", "admin_ip_check", "login_error_max_num",
}

func (s *ConfigGeneralService) GetSecurityConfig() (*SecurityConfig, error) {
	var configs []model.SystemConfig
	if err := s.db.Where("`key` IN ?", securityConfigKeys).Find(&configs).Error; err != nil {
		return nil, err
	}
	m := make(map[string]string, len(configs))
	for _, c := range configs {
		m[c.Key] = c.Value
	}
	return &SecurityConfig{
		RequiredPasswordStrength: m["required_pwstrength"],
		InvalidLoginsBanLength:   parseInt(m["invalid_logins_banlength"]),
		IPCheckFrontend:          m["home_ip_check"] == "1",
		IPCheckAdmin:             m["admin_ip_check"] == "1",
		LoginErrorMaxNum:         parseInt(m["login_error_max_num"]),
	}, nil
}

func (s *ConfigGeneralService) UpdateSecurityConfig(req SecurityConfig) error {
	configs := map[string]string{
		"required_pwstrength":   req.RequiredPasswordStrength,
		"invalid_logins_banlength": intStr(req.InvalidLoginsBanLength),
		"home_ip_check":         boolToIntStr(req.IPCheckFrontend),
		"admin_ip_check":        boolToIntStr(req.IPCheckAdmin),
		"login_error_max_num":   intStr(req.LoginErrorMaxNum),
	}
	return s.saveConfigMap(configs, "security")
}

func boolToIntStr(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

// ==================== Local Config ====================

type LocalConfig struct {
	Charset           string `json:"charset"`
	DateFormat        string `json:"date_format"`
	ClientDateFormat  string `json:"client_date_format"`
	DefaultCountry    string `json:"default_country"`
	Language          string `json:"language"`
	AllowUserLanguage bool   `json:"allow_user_language"`
	TelCCInput        bool   `json:"tel_cc_input"`
}

var localConfigKeys = []string{
	"charset", "date_format", "client_date_format",
	"default_country", "language", "allow_user_language", "tel_cc_input",
}

func (s *ConfigGeneralService) GetLocalConfig() (*LocalConfig, error) {
	var configs []model.SystemConfig
	if err := s.db.Where("`key` IN ?", localConfigKeys).Find(&configs).Error; err != nil {
		return nil, err
	}
	m := make(map[string]string, len(configs))
	for _, c := range configs {
		m[c.Key] = c.Value
	}
	return &LocalConfig{
		Charset:           m["charset"],
		DateFormat:        m["date_format"],
		ClientDateFormat:  m["client_date_format"],
		DefaultCountry:    m["default_country"],
		Language:          m["language"],
		AllowUserLanguage: m["allow_user_language"] == "true",
		TelCCInput:        m["tel_cc_input"] == "true",
	}, nil
}

func (s *ConfigGeneralService) UpdateLocalConfig(req LocalConfig) error {
	configs := map[string]string{
		"charset":            req.Charset,
		"date_format":        req.DateFormat,
		"client_date_format": req.ClientDateFormat,
		"default_country":    req.DefaultCountry,
		"language":           req.Language,
		"allow_user_language": BoolStr(req.AllowUserLanguage),
		"tel_cc_input":       BoolStr(req.TelCCInput),
	}
	return s.saveConfigMap(configs, "local")
}

// ==================== Affiliate Config ====================

type AffiliateConfig struct {
	Enabled             bool    `json:"enabled"`
	BonusDeposit        float64 `json:"bonus_deposit"`
	Rate                float64 `json:"rate"`
	Type                int     `json:"type"`
	CookieDays          int     `json:"cookie_days"`
	MinWithdraw         float64 `json:"min_withdraw"`
	RequireAuth         bool    `json:"require_auth"`
	DelayCommission     int     `json:"delay_commission"`
	IsReorder           bool    `json:"is_reorder"`
	ReorderRate         float64 `json:"reorder_rate"`
	ReorderType         int     `json:"reorder_type"`
	IsRenew             bool    `json:"is_renew"`
	RenewRate           float64 `json:"renew_rate"`
	RenewType           int     `json:"renew_type"`
	Invited             bool    `json:"invited"`
	InvitedMoney        float64 `json:"invited_money"`
	InvitedType         int     `json:"invited_type"`
}

var affiliateConfigKeys = []string{
	"affiliate_enabled", "affiliate_bonusde_posit", "affiliate_bates", "affiliate_type",
	"affiliate_cookie", "affiliate_withdraw", "affiliate_is_authentication",
	"affiliate_delay_commission", "affiliate_is_reorder", "affiliate_reorder", "affiliate_reorder_type",
	"affiliate_is_renew", "affiliate_renew", "affiliate_renew_type",
	"affiliate_invited", "affiliate_invited_money", "affiliate_invited_type",
}

func (s *ConfigGeneralService) GetAffiliateConfig() (*AffiliateConfig, error) {
	var configs []model.SystemConfig
	if err := s.db.Where("`key` IN ?", affiliateConfigKeys).Find(&configs).Error; err != nil {
		return nil, err
	}
	m := make(map[string]string, len(configs))
	for _, c := range configs {
		m[c.Key] = c.Value
	}
	return &AffiliateConfig{
		Enabled:         m["affiliate_enabled"] == "1",
		BonusDeposit:    parseFloat(m["affiliate_bonusde_posit"]),
		Rate:            parseFloat(m["affiliate_bates"]),
		Type:            parseInt(m["affiliate_type"]),
		CookieDays:      parseInt(m["affiliate_cookie"]),
		MinWithdraw:     parseFloat(m["affiliate_withdraw"]),
		RequireAuth:     m["affiliate_is_authentication"] == "1",
		DelayCommission: parseInt(m["affiliate_delay_commission"]),
		IsReorder:       m["affiliate_is_reorder"] == "1",
		ReorderRate:     parseFloat(m["affiliate_reorder"]),
		ReorderType:     parseInt(m["affiliate_reorder_type"]),
		IsRenew:         m["affiliate_is_renew"] == "1",
		RenewRate:       parseFloat(m["affiliate_renew"]),
		RenewType:       parseInt(m["affiliate_renew_type"]),
		Invited:         m["affiliate_invited"] == "1",
		InvitedMoney:    parseFloat(m["affiliate_invited_money"]),
		InvitedType:     parseInt(m["affiliate_invited_type"]),
	}, nil
}

func (s *ConfigGeneralService) UpdateAffiliateConfig(req AffiliateConfig) error {
	configs := map[string]string{
		"affiliate_enabled":           boolOneStr(req.Enabled),
		"affiliate_bonusde_posit":     fmt.Sprintf("%.2f", req.BonusDeposit),
		"affiliate_bates":             fmt.Sprintf("%.2f", req.Rate),
		"affiliate_type":              intStr(req.Type),
		"affiliate_cookie":            intStr(req.CookieDays),
		"affiliate_withdraw":          fmt.Sprintf("%.2f", req.MinWithdraw),
		"affiliate_is_authentication": boolOneStr(req.RequireAuth),
		"affiliate_delay_commission":  intStr(req.DelayCommission),
		"affiliate_is_reorder":        boolOneStr(req.IsReorder),
		"affiliate_reorder":           fmt.Sprintf("%.2f", req.ReorderRate),
		"affiliate_reorder_type":      intStr(req.ReorderType),
		"affiliate_is_renew":          boolOneStr(req.IsRenew),
		"affiliate_renew":             fmt.Sprintf("%.2f", req.RenewRate),
		"affiliate_renew_type":        intStr(req.RenewType),
		"affiliate_invited":           boolOneStr(req.Invited),
		"affiliate_invited_money":     fmt.Sprintf("%.2f", req.InvitedMoney),
		"affiliate_invited_type":      intStr(req.InvitedType),
	}
	return s.saveConfigMap(configs, "affiliate")
}

func boolOneStr(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

// ==================== Captcha Config ====================

type CaptchaConfigData struct {
	Enabled               bool   `json:"enabled"`
	Length                int    `json:"length"`
	Combination           int    `json:"combination"`
	RegisterEmailCaptcha  bool   `json:"register_email_captcha"`
	RegisterPhoneCaptcha  bool   `json:"register_phone_captcha"`
	LoginPhoneCaptcha     bool   `json:"login_phone_captcha"`
	LoginEmailCaptcha     bool   `json:"login_email_captcha"`
	LoginCodeCaptcha      bool   `json:"login_code_captcha"`
	LoginIDCaptcha        bool   `json:"login_id_captcha"`
	LoginAdminCaptcha     bool   `json:"login_admin_captcha"`
}

var captchaConfigKeys = []string{
	"is_captcha", "captcha_length", "captcha_combination",
	"allow_register_email_captcha", "allow_register_phone_captcha",
	"allow_login_phone_captcha", "allow_login_email_captcha",
	"allow_login_code_captcha", "allow_login_id_captcha",
	"allow_login_admin_captcha",
}

func (s *ConfigGeneralService) GetCaptchaConfig() (*CaptchaConfigData, error) {
	var configs []model.SystemConfig
	if err := s.db.Where("`key` IN ?", captchaConfigKeys).Find(&configs).Error; err != nil {
		return nil, err
	}
	m := make(map[string]string, len(configs))
	for _, c := range configs {
		m[c.Key] = c.Value
	}
	return &CaptchaConfigData{
		Enabled:              m["is_captcha"] != "0",
		Length:               parseInt(m["captcha_length"]),
		Combination:          parseInt(m["captcha_combination"]),
		RegisterEmailCaptcha: m["allow_register_email_captcha"] != "0",
		RegisterPhoneCaptcha: m["allow_register_phone_captcha"] != "0",
		LoginPhoneCaptcha:    m["allow_login_phone_captcha"] != "0",
		LoginEmailCaptcha:    m["allow_login_email_captcha"] != "0",
		LoginCodeCaptcha:     m["allow_login_code_captcha"] != "0",
		LoginIDCaptcha:       m["allow_login_id_captcha"] != "0",
		LoginAdminCaptcha:    m["allow_login_admin_captcha"] != "0",
	}, nil
}

func (s *ConfigGeneralService) UpdateCaptchaConfig(req CaptchaConfigData) error {
	configs := map[string]string{
		"is_captcha":                  boolOneStr(req.Enabled),
		"captcha_length":              intStr(req.Length),
		"captcha_combination":         intStr(req.Combination),
		"allow_register_email_captcha": boolOneStr(req.RegisterEmailCaptcha),
		"allow_register_phone_captcha": boolOneStr(req.RegisterPhoneCaptcha),
		"allow_login_phone_captcha":   boolOneStr(req.LoginPhoneCaptcha),
		"allow_login_email_captcha":   boolOneStr(req.LoginEmailCaptcha),
		"allow_login_code_captcha":    boolOneStr(req.LoginCodeCaptcha),
		"allow_login_id_captcha":      boolOneStr(req.LoginIDCaptcha),
		"allow_login_admin_captcha":   boolOneStr(req.LoginAdminCaptcha),
	}
	return s.saveConfigMap(configs, "captcha")
}

// ==================== Buy Product Config ====================

type BuyProductConfig struct {
	MustBindPhone     bool   `json:"must_bind_phone"`
	RequireRealName   bool   `json:"require_real_name"`
	OrderPageStyle    string `json:"order_page_style"`
}

var buyProductConfigKeys = []string{
	"buy_product_must_bind_phone", "certifi_isrealname", "order_page_style",
}

func (s *ConfigGeneralService) GetBuyProductConfig() (*BuyProductConfig, error) {
	var configs []model.SystemConfig
	if err := s.db.Where("`key` IN ?", buyProductConfigKeys).Find(&configs).Error; err != nil {
		return nil, err
	}
	m := make(map[string]string, len(configs))
	for _, c := range configs {
		m[c.Key] = c.Value
	}
	return &BuyProductConfig{
		MustBindPhone:   m["buy_product_must_bind_phone"] == "1",
		RequireRealName: m["certifi_isrealname"] == "1",
		OrderPageStyle:  m["order_page_style"],
	}, nil
}

func (s *ConfigGeneralService) UpdateBuyProductConfig(req BuyProductConfig) error {
	configs := map[string]string{
		"buy_product_must_bind_phone": boolOneStr(req.MustBindPhone),
		"certifi_isrealname":         boolOneStr(req.RequireRealName),
		"order_page_style":           req.OrderPageStyle,
	}
	return s.saveConfigMap(configs, "buy_product")
}

// ==================== Second Verify Config ====================

type SecondVerifyConfig struct {
	HomeEnabled      bool     `json:"home_enabled"`
	HomeActions      []string `json:"home_actions"`
	HomeActionTypes  []string `json:"home_action_types"`
	AdminEnabled     bool     `json:"admin_enabled"`
	AdminActions     []string `json:"admin_actions"`
}

var secondVerifyConfigKeys = []string{
	"second_verify_home", "second_verify_action_home", "second_verify_action_home_type",
	"second_verify_admin", "second_verify_action_admin",
}

func (s *ConfigGeneralService) GetSecondVerifyConfig() (*SecondVerifyConfig, error) {
	var configs []model.SystemConfig
	if err := s.db.Where("`key` IN ?", secondVerifyConfigKeys).Find(&configs).Error; err != nil {
		return nil, err
	}
	m := make(map[string]string, len(configs))
	for _, c := range configs {
		m[c.Key] = c.Value
	}

	homeActions := splitCSV(m["second_verify_action_home"])
	homeActionTypes := splitCSV(m["second_verify_action_home_type"])
	adminActions := splitCSV(m["second_verify_action_admin"])

	return &SecondVerifyConfig{
		HomeEnabled:     m["second_verify_home"] != "0",
		HomeActions:     homeActions,
		HomeActionTypes: homeActionTypes,
		AdminEnabled:    m["second_verify_admin"] != "0",
		AdminActions:    adminActions,
	}, nil
}

func (s *ConfigGeneralService) UpdateSecondVerifyConfig(req SecondVerifyConfig) error {
	configs := map[string]string{
		"second_verify_home":            boolOneStr(req.HomeEnabled),
		"second_verify_action_home":     joinCSV(req.HomeActions),
		"second_verify_action_home_type": joinCSV(req.HomeActionTypes),
		"second_verify_admin":           boolOneStr(req.AdminEnabled),
		"second_verify_action_admin":    joinCSV(req.AdminActions),
	}
	return s.saveConfigMap(configs, "second_verify")
}

func splitCSV(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

func joinCSV(items []string) string {
	return strings.Join(items, ",")
}

// ==================== Nav Group ====================

type NavGroupReq struct {
	ID       uint   `json:"id"`
	Groupname string `json:"groupname"`
	FaIcon   string `json:"fa_icon"`
}

type NavGroup struct {
	ID        uint   `json:"id"`
	Groupname string `json:"groupname"`
	FaIcon    string `json:"fa_icon"`
	Order     int    `json:"order"`
}

func (s *ConfigGeneralService) GetNavGroups() ([]NavGroup, error) {
	var groups []NavGroup
	if err := s.db.Table("nav_groups").Order("`order` ASC").Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (s *ConfigGeneralService) CreateNavGroup(req NavGroupReq) error {
	var maxOrder int
	s.db.Table("nav_groups").Select("COALESCE(MAX(`order`), 0)").Scan(&maxOrder)
	return s.db.Table("nav_groups").Create(map[string]interface{}{
		"groupname": req.Groupname,
		"fa_icon":   req.FaIcon,
		"order":     maxOrder + 1,
	}).Error
}

func (s *ConfigGeneralService) UpdateNavGroup(req NavGroupReq) error {
	return s.db.Table("nav_groups").Where("id = ?", req.ID).Updates(map[string]interface{}{
		"groupname": req.Groupname,
		"fa_icon":   req.FaIcon,
	}).Error
}

func (s *ConfigGeneralService) DeleteNavGroup(id, toID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("nav_groups").Where("id = ?", id).Delete(nil).Error; err != nil {
			return err
		}
		if toID > 0 {
			if err := tx.Table("nav_group_user").Where("groupid = ?", id).Update("groupid", toID).Error; err != nil {
				return err
			}
			if err := tx.Table("products").Where("groupid = ?", id).Update("groupid", toID).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ==================== Language Config ====================

type LanguageConfig struct {
	Languages       []LanguageItem `json:"languages"`
	CurrentLanguage string         `json:"current_language"`
}

type LanguageItem struct {
	Name   string `json:"name"`
	NameZH string `json:"name_zh"`
}

func (s *ConfigGeneralService) GetLanguageConfig() (*LanguageConfig, error) {
	var cfg model.SystemConfig
	if err := s.db.Where("`key` = ?", "language").First(&cfg).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	langs := []LanguageItem{
		{Name: "zh-cn", NameZH: "简体中文"},
		{Name: "en", NameZH: "English"},
		{Name: "ja", NameZH: "日本語"},
	}
	return &LanguageConfig{
		Languages:       langs,
		CurrentLanguage: cfg.Value,
	}, nil
}

func (s *ConfigGeneralService) SetAdminLanguage(lang string) error {
	return s.saveSingleKey("language", lang, "local")
}

// ==================== Header/Footer Config ====================

type HeaderFooterConfig struct {
	Header string `json:"header"`
	Footer string `json:"footer"`
}

func (s *ConfigGeneralService) GetHeaderFooter() (*HeaderFooterConfig, error) {
	var configs []model.SystemConfig
	if err := s.db.Where("`key` IN ?", []string{"header", "footer"}).Find(&configs).Error; err != nil {
		return nil, err
	}
	m := make(map[string]string, len(configs))
	for _, c := range configs {
		m[c.Key] = c.Value
	}
	return &HeaderFooterConfig{
		Header: m["header"],
		Footer: m["footer"],
	}, nil
}

func (s *ConfigGeneralService) UpdateHeaderFooter(req HeaderFooterConfig) error {
	configs := map[string]string{
		"header": req.Header,
		"footer": req.Footer,
	}
	return s.saveConfigMap(configs, "general")
}

// ==================== New Login Page Config ====================

type NewLoginPageConfig struct {
	AllowNewLoginTemplate []string `json:"allow_new_login_template"`
}

func (s *ConfigGeneralService) GetNewLoginPageConfig() (*NewLoginPageConfig, error) {
	var cfg model.SystemConfig
	if err := s.db.Where("`key` = ?", "allow_new_login_template").First(&cfg).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &NewLoginPageConfig{AllowNewLoginTemplate: []string{"default"}}, nil
		}
		return nil, err
	}
	templates := splitCSV(cfg.Value)
	if len(templates) == 0 {
		templates = []string{"default"}
	}
	return &NewLoginPageConfig{AllowNewLoginTemplate: templates}, nil
}

func (s *ConfigGeneralService) UpdateNewLoginPageConfig(req NewLoginPageConfig) error {
	return s.saveSingleKey("allow_new_login_template", joinCSV(req.AllowNewLoginTemplate), "login")
}
