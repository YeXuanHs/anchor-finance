package handler

import (
	"net/http"

	"anchorfinance/internal/service"

	"github.com/gin-gonic/gin"
)

// CaptchaConfigHandler 验证码配置管理
type CaptchaConfigHandler struct {
	captchaService *service.CaptchaService
}

func NewCaptchaConfigHandler(captchaService *service.CaptchaService) *CaptchaConfigHandler {
	return &CaptchaConfigHandler{captchaService: captchaService}
}

// GetConfigs 获取所有验证码配置
// GET /admin/captcha-config
func (h *CaptchaConfigHandler) GetConfigs(c *gin.Context) {
	configService := h.captchaService.GetCaptchaConfigService()

	configs, err := configService.GetAllConfigs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取配置失败"})
		return
	}

	// 分类整理配置
	basic := make(map[string]interface{})
	scenes := make(map[string]bool)

	for _, config := range configs {
		switch config.Key {
		case "is_captcha", "captcha_type", "captcha_length", "captcha_combination",
			"geetest_captcha_id", "geetest_captcha_key":
			basic[config.Key] = map[string]interface{}{
				"value":  config.Value,
				"status": config.Status,
			}
		default:
			scenes[config.Key] = config.Status && config.Value == "1"
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"basic":  basic,
			"scenes": scenes,
		},
	})
}

// UpdateBasicConfig 更新基础配置
// PUT /admin/captcha-config/basic
func (h *CaptchaConfigHandler) UpdateBasicConfig(c *gin.Context) {
	var req struct {
		IsCaptcha          *bool  `json:"is_captcha"`
		CaptchaType        string `json:"captcha_type"`
		CaptchaLength      *int   `json:"captcha_length"`
		CaptchaCombination string `json:"captcha_combination"`
		GeetestCaptchaID   string `json:"geetest_captcha_id"`
		GeetestCaptchaKey  string `json:"geetest_captcha_key"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	configService := h.captchaService.GetCaptchaConfigService()

	if req.IsCaptcha != nil {
		value := "0"
		if *req.IsCaptcha {
			value = "1"
		}
		if err := configService.UpdateConfig("is_captcha", value, *req.IsCaptcha); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
			return
		}
	}

	if req.CaptchaType != "" {
		if err := configService.UpdateConfig("captcha_type", req.CaptchaType, true); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
			return
		}
	}

	if req.CaptchaLength != nil {
		value := "4"
		switch *req.CaptchaLength {
		case 4:
			value = "4"
		case 5:
			value = "5"
		case 6:
			value = "6"
		}
		if err := configService.UpdateConfig("captcha_length", value, true); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
			return
		}
	}

	if req.CaptchaCombination != "" {
		if err := configService.UpdateConfig("captcha_combination", req.CaptchaCombination, true); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
			return
		}
	}

	// 更新极验配置
	if req.GeetestCaptchaID != "" {
		if err := configService.UpdateConfig("geetest_captcha_id", req.GeetestCaptchaID, true); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
			return
		}
	}

	if req.GeetestCaptchaKey != "" {
		if err := configService.UpdateConfig("geetest_captcha_key", req.GeetestCaptchaKey, true); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// UpdateSceneConfig 更新场景配置
// PUT /admin/captcha-config/scenes
func (h *CaptchaConfigHandler) UpdateSceneConfig(c *gin.Context) {
	var req map[string]bool

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	configService := h.captchaService.GetCaptchaConfigService()

	for key, enabled := range req {
		value := "0"
		if enabled {
			value = "1"
		}
		if err := configService.UpdateConfig(key, value, enabled); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// GetSceneStatus 获取场景状态（前端用）
// GET /captcha/status
func (h *CaptchaConfigHandler) GetSceneStatus(c *gin.Context) {
	sceneConfig := h.captchaService.GetSceneConfig()
	c.JSON(http.StatusOK, gin.H{"data": sceneConfig})
}

// InitDefaultConfigs 初始化默认配置
// POST /admin/captcha-config/init
func (h *CaptchaConfigHandler) InitDefaultConfigs(c *gin.Context) {
	if err := h.captchaService.InitDefaultConfigs(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "初始化失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "初始化成功"})
}
