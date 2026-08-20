package handler

import (
	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CommonHandler handles common admin API requests.
type CommonHandler struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewCommonHandler creates a new CommonHandler.
func NewCommonHandler(db *gorm.DB, log *logger.Logger) *CommonHandler {
	return &CommonHandler{db: db, log: log}
}

// Common returns common admin dashboard data.
// GET /admin/common
func (h *CommonHandler) Common(c *gin.Context) {
	data := map[string]interface{}{}

	// Company info
	var companyConfig model.SystemConfig
	h.db.Where("key = ?", "company_name").First(&companyConfig)
	data["company_name"] = companyConfig.Value

	// Domain
	var domainConfig model.SystemConfig
	h.db.Where("key = ?", "domain").First(&domainConfig)
	data["domain"] = domainConfig.Value

	// System URL
	var urlConfig model.SystemConfig
	h.db.Where("key = ?", "system_url").First(&urlConfig)
	data["system_url"] = urlConfig.Value

	// Gateways
	var gateways []model.PaymentGateway
	h.db.Where("is_active = ?", true).Find(&gateways)
	data["gateway"] = gateways

	// Current admin
	adminID, _ := c.Get("admin_id")
	data["admin"] = adminID

	// Language
	var langConfig model.SystemConfig
	h.db.Where("key = ?", "language").First(&langConfig)
	data["system_language"] = langConfig.Value

	// System title
	var titleConfig model.SystemConfig
	h.db.Where("key = ?", "title").First(&titleConfig)
	data["system_title"] = titleConfig.Value

	// Affiliate enabled
	var affConfig model.SystemConfig
	h.db.Where("key = ?", "affiliate_enabled").First(&affConfig)
	data["is_aff"] = affConfig.Value == "1"

	// Per page limit
	var limitConfig model.SystemConfig
	h.db.Where("key = ?", "per_page_limit").First(&limitConfig)
	if limitConfig.Value != "" {
		data["per_page_limit"] = limitConfig.Value
	} else {
		data["per_page_limit"] = 50
	}

	response.Success(c, data)
}

// InfoNotice returns admin notification messages.
// GET /admin/info-notice
func (h *CommonHandler) InfoNotice(c *gin.Context) {
	var notices []struct {
		Info string `json:"info"`
	}
	h.db.Table("info_notices").Where("info <> ?", "").Select("info").Find(&notices)

	info := ""
	for _, n := range notices {
		info += n.Info + "\n"
	}

	response.Success(c, gin.H{"info": info})
}

// GetGetways returns all available payment gateways.
// GET /admin/gateways
func (h *CommonHandler) GetGetways(c *gin.Context) {
	var gateways []model.PaymentGateway
	if err := h.db.Find(&gateways).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"gateway": gateways})
}

// GetEmailTem returns email templates filtered by type.
// GET /admin/email-templates
func (h *CommonHandler) GetEmailTem(c *gin.Context) {
	templateType := c.Query("type")

	query := h.db.Model(&model.EmailTemplate{})
	if templateType != "" {
		query = query.Where("type = ?", templateType)
	}

	var templates []model.EmailTemplate
	if err := query.Find(&templates).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"email": templates})
}

// SmsCountry model for SMS country codes.
type SmsCountry struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Nicename string `gorm:"type:varchar(128)" json:"nicename"`
	NameZh   string `gorm:"type:varchar(128)" json:"name_zh"`
	PhoneCode string `gorm:"type:varchar(16)" json:"phone_code"`
}

func (SmsCountry) TableName() string {
	return "sms_country"
}

// GetSmsCountry returns SMS country list with phone codes.
// GET /admin/sms-countries
func (h *CommonHandler) GetSmsCountry(c *gin.Context) {
	var countries []SmsCountry
	if err := h.db.Find(&countries).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"sms_country": countries})
}
