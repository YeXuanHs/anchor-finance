package mail

import "fmt"

// SenderFactory creates a MailSender instance based on the plugin name.
type SenderFactory struct{}

// NewSenderFactory creates a new SenderFactory.
func NewSenderFactory() *SenderFactory {
	return &SenderFactory{}
}

// Create creates a MailSender by plugin name with the provided configuration.
// Supported names: smtp, submail, alimail.
func (f *SenderFactory) Create(name string, config map[string]interface{}) (MailSender, error) {
	switch name {
	case "smtp":
		return f.createSmtp(config)
	case "submail":
		return f.createSubmail(config)
	case "alimail":
		return f.createAlimail(config)
	default:
		return nil, fmt.Errorf("unsupported mail plugin: %s", name)
	}
}

func (f *SenderFactory) createSmtp(config map[string]interface{}) (MailSender, error) {
	cfg := SmtpConfig{
		Host:        getString(config, "host"),
		Port:        getInt(config, "port"),
		Username:    getString(config, "username"),
		Password:    getString(config, "password"),
		FromName:    getString(config, "from_name"),
		SystemEmail: getString(config, "system_email"),
		SmtpSecure:  getString(config, "smtp_secure"),
		Charset:     getString(config, "charset"),
	}
	if cfg.Host == "" {
		return nil, fmt.Errorf("smtp: host is required")
	}
	if cfg.Port == 0 {
		cfg.Port = 25
	}
	return NewSmtpSender(cfg), nil
}

func (f *SenderFactory) createSubmail(config map[string]interface{}) (MailSender, error) {
	cfg := SubmailConfig{
		AppID:     getString(config, "app_id"),
		AppKey:    getString(config, "app_key"),
		SignType:  getString(config, "sign_type"),
		FromName:  getString(config, "from_name"),
		FromEmail: getString(config, "from_email"),
	}
	if cfg.AppID == "" || cfg.AppKey == "" {
		return nil, fmt.Errorf("submail: app_id and app_key are required")
	}
	return NewSubmailSender(cfg), nil
}

func (f *SenderFactory) createAlimail(config map[string]interface{}) (MailSender, error) {
	cfg := AlimailConfig{
		AccessKeyID:     getString(config, "access_key_id"),
		AccessKeySecret: getString(config, "access_key_secret"),
		AccountName:     getString(config, "account_name"),
		FromAlias:       getString(config, "from_alias"),
	}
	if cfg.AccessKeyID == "" || cfg.AccessKeySecret == "" {
		return nil, fmt.Errorf("alimail: access_key_id and access_key_secret are required")
	}
	return NewAlimailSender(cfg), nil
}

// getString extracts a string value from a config map.
func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// getInt extracts an int value from a config map.
func getInt(m map[string]interface{}, key string) int {
	if v, ok := m[key]; ok {
		switch n := v.(type) {
		case int:
			return n
		case float64:
			return int(n)
		}
	}
	return 0
}

// SubmailConfig holds the configuration for Submail mail sender.
type SubmailConfig struct {
	AppID     string `json:"app_id"`
	AppKey    string `json:"app_key"`
	SignType  string `json:"sign_type"`
	FromName  string `json:"from_name"`
	FromEmail string `json:"from_email"`
}

// SubmailSender implements MailSender for Submail API.
type SubmailSender struct {
	config SubmailConfig
}

// NewSubmailSender creates a new SubmailSender with the given configuration.
func NewSubmailSender(config SubmailConfig) *SubmailSender {
	return &SubmailSender{config: config}
}

// Send sends an email via Submail API.
// TODO: implement Submail REST API integration
func (s *SubmailSender) Send(to, subject, content string, attachments []string) error {
	return fmt.Errorf("submail: not yet implemented")
}

// AlimailConfig holds the configuration for Aliyun DirectMail sender.
type AlimailConfig struct {
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	AccountName     string `json:"account_name"`
	FromAlias       string `json:"from_alias"`
}

// AlimailSender implements MailSender for Aliyun DirectMail API.
type AlimailSender struct {
	config AlimailConfig
}

// NewAlimailSender creates a new AlimailSender with the given configuration.
func NewAlimailSender(config AlimailConfig) *AlimailSender {
	return &AlimailSender{config: config}
}

// Send sends an email via Aliyun DirectMail API.
// TODO: implement Aliyun DirectMail SDK integration
func (s *AlimailSender) Send(to, subject, content string, attachments []string) error {
	return fmt.Errorf("alimail: not yet implemented")
}
