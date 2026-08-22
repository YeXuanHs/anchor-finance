package admin

import (
	"fmt"
	"net/http"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetEmailConfig 获取邮件配置
// GET /api/admin/settings/email
func GetEmailConfig(c *gin.Context) {
	db := database.GetDB()

	// 从settings表获取邮件配置
	var settings []model.Setting
	db.Where("`group` = ?", "email").Find(&settings)

	// 转换为map
	config := make(map[string]interface{})
	for _, s := range settings {
		config[s.Key] = s.Value
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    config,
	})
}

// UpdateEmailConfig 更新邮件配置
// PUT /api/admin/settings/email
func UpdateEmailConfig(c *gin.Context) {
	var req struct {
		SMTPHost     string `json:"smtp_host"`
		SMTPPort     int    `json:"smtp_port"`
		SMTPUser     string `json:"smtp_user"`
		SMTPPassword string `json:"smtp_password"`
		SMTPFrom     string `json:"smtp_from"`
		SMTPFromName string `json:"smtp_from_name"`
		EnableSSL    bool   `json:"enable_ssl"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	db := database.GetDB()

	// 保存配置到settings表
	settings := map[string]string{
		"smtp_host":      req.SMTPHost,
		"smtp_port":      string(rune(req.SMTPPort)),
		"smtp_user":      req.SMTPUser,
		"smtp_password":  req.SMTPPassword,
		"smtp_from":      req.SMTPFrom,
		"smtp_from_name": req.SMTPFromName,
		"enable_ssl":     "false",
	}

	if req.EnableSSL {
		settings["enable_ssl"] = "true"
	}

	for key, value := range settings {
		var setting model.Setting
		result := db.Where("`key` = ? AND `group` = ?", key, "email").First(&setting)
		if result.Error != nil {
			// 不存在则创建
			setting = model.Setting{
				Key:   key,
				Value: value,
				Group: "email",
			}
			db.Create(&setting)
		} else {
			// 存在则更新
			db.Model(&setting).Update("value", value)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "保存成功",
	})
}

// GetSMSConfig 获取短信配置
// GET /api/admin/settings/sms
func GetSMSConfig(c *gin.Context) {
	db := database.GetDB()

	var settings []model.Setting
	db.Where("`group` = ?", "sms").Find(&settings)

	config := make(map[string]interface{})
	for _, s := range settings {
		config[s.Key] = s.Value
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    config,
	})
}

// UpdateSMSConfig 更新短信配置
// PUT /api/admin/settings/sms
func UpdateSMSConfig(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	db := database.GetDB()

	for key, value := range req {
		var setting model.Setting
		result := db.Where("`key` = ? AND `group` = ?", key, "sms").First(&setting)
		if result.Error != nil {
			setting = model.Setting{
				Key:   key,
				Value: fmt.Sprintf("%v", value),
				Group: "sms",
			}
			db.Create(&setting)
		} else {
			db.Model(&setting).Update("value", fmt.Sprintf("%v", value))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "保存成功",
	})
}

// GetRegisterLoginConfig 获取注册登录配置
// GET /api/admin/settings/register-login
func GetRegisterLoginConfig(c *gin.Context) {
	db := database.GetDB()

	var settings []model.Setting
	db.Where("`group` = ?", "register_login").Find(&settings)

	config := make(map[string]interface{})
	for _, s := range settings {
		config[s.Key] = s.Value
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    config,
	})
}

// UpdateRegisterLoginConfig 更新注册登录配置
// PUT /api/admin/settings/register-login
func UpdateRegisterLoginConfig(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	db := database.GetDB()

	for key, value := range req {
		var setting model.Setting
		result := db.Where("`key` = ? AND `group` = ?", key, "register_login").First(&setting)
		if result.Error != nil {
			setting = model.Setting{
				Key:   key,
				Value: fmt.Sprintf("%v", value),
				Group: "register_login",
			}
			db.Create(&setting)
		} else {
			db.Model(&setting).Update("value", fmt.Sprintf("%v", value))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "保存成功",
	})
}

// GetCaptchaConfig 获取验证码配置
// GET /api/admin/settings/captcha
func GetCaptchaConfig(c *gin.Context) {
	db := database.GetDB()

	var settings []model.Setting
	db.Where("`group` = ?", "captcha").Find(&settings)

	config := make(map[string]interface{})
	for _, s := range settings {
		config[s.Key] = s.Value
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    config,
	})
}

// UpdateCaptchaConfig 更新验证码配置
// PUT /api/admin/settings/captcha
func UpdateCaptchaConfig(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	db := database.GetDB()

	for key, value := range req {
		var setting model.Setting
		result := db.Where("`key` = ? AND `group` = ?", key, "captcha").First(&setting)
		if result.Error != nil {
			setting = model.Setting{
				Key:   key,
				Value: fmt.Sprintf("%v", value),
				Group: "captcha",
			}
			db.Create(&setting)
		} else {
			db.Model(&setting).Update("value", fmt.Sprintf("%v", value))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "保存成功",
	})
}
