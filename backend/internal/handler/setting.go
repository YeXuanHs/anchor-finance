package handler

import (
	"net/http"

	"anchorfinance/pkg/db"

	"github.com/gin-gonic/gin"
)

// SettingHandler handles settings-related requests.
type SettingHandler struct{}

// NewSettingHandler creates a new SettingHandler.
func NewSettingHandler() *SettingHandler {
	return &SettingHandler{}
}

// GetSettings gets settings by group.
func (h *SettingHandler) GetSettings(c *gin.Context) {
	group := c.Param("group")

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{}})
		return
	}

	settings := db.GetSystemSettings(group)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": settings})
}

// UpdateSettings updates settings by group.
func (h *SettingHandler) UpdateSettings(c *gin.Context) {
	group := c.Param("group")

	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	for key, value := range req {
		db.SetSystemSetting(key, value, group, "")
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "设置更新成功"})
}

// RegisterRoutes registers setting routes.
func (h *SettingHandler) RegisterRoutes(r *gin.RouterGroup) {
	setting := r.Group("/settings")
	{
		setting.GET("/:group", h.GetSettings)
		setting.POST("/:group", h.UpdateSettings)
	}
}

// GetNotificationSettings returns notification settings.
func (h *SettingHandler) GetNotificationSettings(c *gin.Context) {
	settings := db.GetSystemSettings("notification")
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": settings})
}

// SaveNotificationSettings saves notification settings.
func (h *SettingHandler) SaveNotificationSettings(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	for key, value := range req {
		db.SetSystemSetting(key, value, "notification", "")
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "保存成功"})
}

// GetMaintenanceMode returns maintenance mode status.
func (h *SettingHandler) GetMaintenanceMode(c *gin.Context) {
	enabled := db.GetSystemSetting("maintenance_mode")
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"enabled": enabled == "true",
		},
	})
}

// SetMaintenanceMode sets maintenance mode.
func (h *SettingHandler) SetMaintenanceMode(c *gin.Context) {
	var req struct {
		Enabled bool `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	value := "false"
	if req.Enabled {
		value = "true"
	}

	db.SetSystemSetting("maintenance_mode", value, "system", "维护模式")
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "设置成功"})
}

// GetCronSettings returns cron settings.
func (h *SettingHandler) GetCronSettings(c *gin.Context) {
	settings := db.GetSystemSettings("cron")
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": settings})
}

// SaveCronSettings saves cron settings.
func (h *SettingHandler) SaveCronSettings(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	for key, value := range req {
		db.SetSystemSetting(key, value, "cron", "")
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "保存成功"})
}

// GetSiteSettings returns site settings.
func (h *SettingHandler) GetSiteSettings(c *gin.Context) {
	settings := db.GetSystemSettings("site")
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": settings})
}

// SaveSiteSettings saves site settings.
func (h *SettingHandler) SaveSiteSettings(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	for key, value := range req {
		db.SetSystemSetting(key, value, "site", "")
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "保存成功"})
}

// GetPaymentSettings returns payment settings.
func (h *SettingHandler) GetPaymentSettings(c *gin.Context) {
	settings := db.GetSystemSettings("payment")
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": settings})
}

// SavePaymentSettings saves payment settings.
func (h *SettingHandler) SavePaymentSettings(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	for key, value := range req {
		db.SetSystemSetting(key, value, "payment", "")
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "保存成功"})
}
