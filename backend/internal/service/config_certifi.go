package service

import (
	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type ConfigCertifiService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewConfigCertifiService(db *gorm.DB, log *logger.Logger) *ConfigCertifiService {
	return &ConfigCertifiService{db: db, log: log}
}

type CertificationConfig struct {
	ForceCertify     bool   `json:"force_certify"`      // 是否强制实名认证
	AllowType        string `json:"allow_type"`         // individual/enterprise/both
	EnableIndividual bool   `json:"enable_individual"`  // 是否开启个人认证
	EnableEnterprise bool   `json:"enable_enterprise"`  // 是否开启企业认证
	RequireIDCard    bool   `json:"require_id_card"`    // 是否要求上传身份证
	RequireHandImage bool   `json:"require_hand_image"` // 是否要求上传手持身份证
	MaxReviewDays    int    `json:"max_review_days"`    // 最大审核天数
	AutoApprove      bool   `json:"auto_approve"`       // 是否自动通过
	CertifyNotice    string `json:"certify_notice"`     // 认证提示文案
	ReviewTemplate   string `json:"review_template"`    // 审核结果通知模板
	AllowResubmit    bool   `json:"allow_resubmit"`     // 是否允许重新提交
	ResubmitLimit    int    `json:"resubmit_limit"`     // 重新提交次数限制，0=不限
}

var certifiConfigKeys = []string{
	"certifi_force_certify", "certifi_allow_type",
	"certifi_enable_individual", "certifi_enable_enterprise",
	"certifi_require_id_card", "certifi_require_hand_image",
	"certifi_max_review_days", "certifi_auto_approve",
	"certifi_certify_notice", "certifi_review_template",
	"certifi_allow_resubmit", "certifi_resubmit_limit",
}

func (s *ConfigCertifiService) GetConfig() (*CertificationConfig, error) {
	var configs []model.SystemConfig
	if err := s.db.Where("`key` IN ?", certifiConfigKeys).Find(&configs).Error; err != nil {
		return nil, err
	}

	configMap := make(map[string]string, len(configs))
	for _, c := range configs {
		configMap[c.Key] = c.Value
	}

	return &CertificationConfig{
		ForceCertify:     configMap["certifi_force_certify"] == "true",
		AllowType:        configMap["certifi_allow_type"],
		EnableIndividual: configMap["certifi_enable_individual"] == "true",
		EnableEnterprise: configMap["certifi_enable_enterprise"] == "true",
		RequireIDCard:    configMap["certifi_require_id_card"] == "true",
		RequireHandImage: configMap["certifi_require_hand_image"] == "true",
		MaxReviewDays:    parseInt(configMap["certifi_max_review_days"]),
		AutoApprove:      configMap["certifi_auto_approve"] == "true",
		CertifyNotice:    configMap["certifi_certify_notice"],
		ReviewTemplate:   configMap["certifi_review_template"],
		AllowResubmit:    configMap["certifi_allow_resubmit"] == "true",
		ResubmitLimit:    parseInt(configMap["certifi_resubmit_limit"]),
	}, nil
}

func (s *ConfigCertifiService) UpdateConfig(req CertificationConfig) error {
	configs := map[string]string{
		"certifi_force_certify":      boolStr(req.ForceCertify),
		"certifi_allow_type":         req.AllowType,
		"certifi_enable_individual":  boolStr(req.EnableIndividual),
		"certifi_enable_enterprise":  boolStr(req.EnableEnterprise),
		"certifi_require_id_card":    boolStr(req.RequireIDCard),
		"certifi_require_hand_image": boolStr(req.RequireHandImage),
		"certifi_max_review_days":    intStr(req.MaxReviewDays),
		"certifi_auto_approve":       boolStr(req.AutoApprove),
		"certifi_certify_notice":     req.CertifyNotice,
		"certifi_review_template":    req.ReviewTemplate,
		"certifi_allow_resubmit":     boolStr(req.AllowResubmit),
		"certifi_resubmit_limit":     intStr(req.ResubmitLimit),
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		for key, value := range configs {
			result := tx.Where("key = ?", key).
				Assign(model.SystemConfig{Value: value}).
				FirstOrCreate(&model.SystemConfig{
					Key:   key,
					Value: value,
					Group: "certification",
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

func parseInt(s string) int {
	if s == "" {
		return 0
	}
	n := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
}

func intStr(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}
