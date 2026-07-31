package service

import (
	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"
	"encoding/json"
	"fmt"

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
		"open_register":    boolStr(req.OpenRegister),
		"verify_email":     boolStr(req.VerifyEmail),
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

func boolStr(b bool) string {
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
		"safe_force_2fa":          boolStr(req.Force2FA),
		"safe_captcha_login":      boolStr(req.CaptchaLogin),
		"safe_captcha_register":   boolStr(req.CaptchaRegister),
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
		"recharge_auto_approve":  boolStr(req.AutoApprove),
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
		"invoice_auto_generate":  boolStr(req.AutoGenerate),
		"invoice_tax_rate":       fmt.Sprintf("%.2f", req.TaxRate),
		"invoice_title_required": boolStr(req.TitleRequired),
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
		"register_enable":       boolStr(req.EnableRegister),
		"register_email_verify": boolStr(req.EmailVerify),
		"register_phone_verify": boolStr(req.PhoneVerify),
		"register_captcha":      boolStr(req.Captcha),
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
		"login_captcha":     boolStr(req.Captcha),
		"login_remember_me": boolStr(req.RememberMe),
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
		"api_enable":      boolStr(req.EnableAPI),
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
		"twofactor_enable":       boolStr(req.Enable),
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
	return s.saveSingleKey("debug_mode", boolStr(enabled), "system")
}

// ==================== SMTP Test ====================

func (s *ConfigGeneralService) TestSMTP(toEmail string) error {
	// Load SMTP config
	cfg, err := s.GetEmailConfig()
	if err != nil {
		return fmt.Errorf("failed to load email config: %w", err)
	}
	if cfg.SMTPHost == "" {
		return fmt.Errorf("SMTP host not configured")
	}

	// Build connection address
	addr := fmt.Sprintf("%s:%d", cfg.SMTPHost, cfg.SMTPPort)

	// TODO: Integrate with actual SMTP sending library (e.g., net/smtp or gomail).
	// For now, validate configuration is complete and return success.
	s.log.Info("SMTP test", "addr", addr, "from", cfg.FromAddress, "to", toEmail)
	return nil
}

// ==================== SMS Test ====================

func (s *ConfigGeneralService) TestSMS(phone string) error {
	if phone == "" {
		return fmt.Errorf("phone number is required")
	}

	// TODO: Integrate with actual SMS provider (e.g., Aliyun SMS, Tencent SMS).
	// For now, validate input and return success.
	s.log.Info("SMS test", "phone", phone)
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
