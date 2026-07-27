package service

import (
	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type ConfigGeneralService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewConfigGeneralService(db *gorm.DB, log *logger.Logger) *ConfigGeneralService {
	return &ConfigGeneralService{db: db, log: log}
}

type GeneralConfig struct {
	SiteName       string `json:"site_name"`
	SiteURL        string `json:"site_url"`
	Logo           string `json:"logo"`
	Favicon        string `json:"favicon"`
	Description    string `json:"description"`
	Keywords       string `json:"keywords"`
	ICP            string `json:"icp"`             // 备案号
	PSB            string `json:"psb"`             // 公安备案号
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
