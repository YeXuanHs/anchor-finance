package handler

import (
	"net/http"

	"anchorfinance/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ConfigHandler 系统配置管理
type ConfigHandler struct {
	configService *model.ConfigService
}

func NewConfigHandler(db *gorm.DB) *ConfigHandler {
	return &ConfigHandler{
		configService: model.NewConfigService(db),
	}
}

// GetGroups 获取配置分组列表
// GET /admin/config/groups
func (h *ConfigHandler) GetGroups(c *gin.Context) {
	groups := h.configService.GetGroups()
	c.JSON(http.StatusOK, gin.H{"data": groups})
}

// GetByGroup 按分组获取配置
// GET /admin/config/group/:group
func (h *ConfigHandler) GetByGroup(c *gin.Context) {
	group := c.Param("group")
	configs, err := h.configService.GetByGroup(group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取配置失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": configs})
}

// GetAll 获取所有配置
// GET /admin/config/all
func (h *ConfigHandler) GetAll(c *gin.Context) {
	configs, err := h.configService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取配置失败"})
		return
	}

	// 按分组整理
	grouped := make(map[string][]model.SystemConfig)
	for _, config := range configs {
		grouped[config.Group] = append(grouped[config.Group], config)
	}

	c.JSON(http.StatusOK, gin.H{"data": grouped})
}

// UpdateConfig 更新单个配置
// PUT /admin/config/:key
func (h *ConfigHandler) UpdateConfig(c *gin.Context) {
	key := c.Param("key")

	var req struct {
		Value string `json:"value"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := h.configService.Set(key, req.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// BatchUpdateConfig 批量更新配置
// PUT /admin/config/batch
func (h *ConfigHandler) BatchUpdateConfig(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := h.configService.SetBatch(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// InitDefaultConfigs 初始化默认配置
// POST /admin/config/init
func (h *ConfigHandler) InitDefaultConfigs(c *gin.Context) {
	if err := h.configService.InitDefaultConfigs(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "初始化失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "初始化成功"})
}

// GetPublicConfig 获取公开配置（前端用）
// GET /config/public
func (h *ConfigHandler) GetPublicConfig(c *gin.Context) {
	config := h.configService.GetPublicConfig()
	c.JSON(http.StatusOK, gin.H{"data": config})
}

// GetLoginConfig 获取登录配置
// GET /config/login
func (h *ConfigHandler) GetLoginConfig(c *gin.Context) {
	config := map[string]interface{}{
		"login_methods":    h.configService.GetLoginMethods(),
		"register_methods": h.configService.GetRegisterMethods(),
		"captcha_enabled":  h.configService.GetBool("is_captcha"),
	}
	c.JSON(http.StatusOK, gin.H{"data": config})
}

// GetMaintenanceStatus 获取维护模式状态
// GET /config/maintenance
func (h *ConfigHandler) GetMaintenanceStatus(c *gin.Context) {
	enabled, url, message := h.configService.GetMaintenanceConfig()
	c.JSON(http.StatusOK, gin.H{
		"data": map[string]interface{}{
			"enabled": enabled,
			"url":     url,
			"message": message,
		},
	})
}
