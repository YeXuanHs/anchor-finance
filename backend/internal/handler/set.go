package handler

import (
	"encoding/json"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SetHandler struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewSetHandler(db *gorm.DB, log *logger.Logger) *SetHandler {
	return &SetHandler{db: db, log: log}
}

// GetSiteSettings returns site settings.
func (h *SetHandler) GetSiteSettings(c *gin.Context) {
	var options []model.ConfigOption
	if err := h.db.Where("`group` IN ?", []string{"cdn", "cmf", "admin"}).Find(&options).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}

	grouped := map[string]gin.H{}
	for _, opt := range options {
		if _, ok := grouped[opt.Group]; !ok {
			grouped[opt.Group] = gin.H{}
		}
		grouped[opt.Group][opt.Code] = opt.Value
	}

	themes := []string{"default", "dark", "blue", "green", "purple"}
	var themeOpt model.ConfigOption
	if err := h.db.Where("code = ?", "admin_theme").First(&themeOpt).Error; err == nil && themeOpt.Value != "" {
		response.Success(c, gin.H{
			"cdn_settings":   grouped["cdn"],
			"cmf_settings":   grouped["cmf"],
			"admin_settings": grouped["admin"],
			"admin_themes":   themes,
			"current_theme":  themeOpt.Value,
		})
		return
	}

	response.Success(c, gin.H{
		"cdn_settings":   grouped["cdn"],
		"cmf_settings":   grouped["cmf"],
		"admin_settings": grouped["admin"],
		"admin_themes":   themes,
	})
}

// UpdateSiteSettings updates site settings.
func (h *SetHandler) UpdateSiteSettings(c *gin.Context) {
	var req struct {
		CDNSettings   interface{} `json:"cdn_settings"`
		CMFSettings   interface{} `json:"cmf_settings"`
		AdminSettings interface{} `json:"admin_settings"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	groups := map[string]interface{}{
		"cdn":   req.CDNSettings,
		"cmf":   req.CMFSettings,
		"admin": req.AdminSettings,
	}

	for group, data := range groups {
		if data == nil {
			continue
		}
		settingsMap, ok := data.(map[string]interface{})
		if !ok {
			continue
		}
		for code, val := range settingsMap {
			valueBytes, _ := json.Marshal(val)
			result := h.db.Model(&model.ConfigOption{}).
				Where("`group` = ? AND code = ?", group, code).
				Update("value", string(valueBytes))
			if result.RowsAffected == 0 {
				h.db.Create(&model.ConfigOption{
					Group: group,
					Code:  code,
					Name:  code,
					Type:  "json",
					Value: string(valueBytes),
				})
			}
		}
	}

	response.Success(c, gin.H{
		"message": "site settings updated",
	})
}

// GetAdminThemes returns available admin themes.
func (h *SetHandler) GetAdminThemes(c *gin.Context) {
	themes := []gin.H{
		{"name": "default", "label": "默认主题"},
		{"name": "dark", "label": "暗色主题"},
		{"name": "blue", "label": "蓝色主题"},
		{"name": "green", "label": "绿色主题"},
		{"name": "purple", "label": "紫色主题"},
	}

	current := "default"
	var opt model.ConfigOption
	if err := h.db.Where("code = ?", "admin_theme").First(&opt).Error; err == nil && opt.Value != "" {
		current = opt.Value
	}

	response.Success(c, gin.H{
		"themes":         themes,
		"current_theme":  current,
	})
}

// SetAdminTheme sets the active admin theme.
func (h *SetHandler) SetAdminTheme(c *gin.Context) {
	var req struct {
		Theme string `json:"theme" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result := h.db.Model(&model.ConfigOption{}).
		Where("code = ?", "admin_theme").
		Update("value", req.Theme)
	if result.RowsAffected == 0 {
		h.db.Create(&model.ConfigOption{
			Group: "admin",
			Code:  "admin_theme",
			Name:  "后台主题",
			Type:  "text",
			Value: req.Theme,
		})
	}

	response.Success(c, gin.H{
		"message": "admin theme set",
		"theme":   req.Theme,
	})
}
