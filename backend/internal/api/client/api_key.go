package client

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/YeXuanHs/anchor-finance/internal/util"
	"github.com/gin-gonic/gin"
)

// GetAPIKeyStatus 查看API状态
// GET /api/client/api-key/status
func GetAPIKeyStatus(c *gin.Context) {
	userID, _ := c.Get("user_id")
	db := database.GetDB()

	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "用户不存在", "data": nil})
		return
	}

	// 读取开放条件配置
	requirePhone := false
	var phoneSetting model.Setting
	if err := db.Where("`key` = ?", "api_require_phone").First(&phoneSetting).Error; err == nil && phoneSetting.Value == "1" {
		requirePhone = true
	}
	requireRealName := false
	var realNameSetting model.Setting
	if err := db.Where("`key` = ?", "api_require_realname").First(&realNameSetting).Error; err == nil && realNameSetting.Value == "1" {
		requireRealName = true
	}
	apiEnabled := true
	var apiSwitch model.Setting
	if err := db.Where("`key` = ?", "api_enabled").First(&apiSwitch).Error; err == nil && apiSwitch.Value == "0" {
		apiEnabled = false
	}

	// 解密API密钥（如果已开通）
	apiKey := ""
	if user.APIEnabled && user.APIKey != "" {
		apiKey, _ = util.DecryptAES(user.APIKey)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "success",
		"data": gin.H{
			"api_enabled":       user.APIEnabled,
			"api_key":           apiKey,
			"system_api_open":   apiEnabled,
			"require_phone":     requirePhone,
			"require_realname":  requireRealName,
			"has_phone":         user.Phone != "",
			"is_verified":       user.IsVerified,
		},
	})
}

// EnableAPIKey 自助开通API
// POST /api/client/api-key/enable
func EnableAPIKey(c *gin.Context) {
	userID, _ := c.Get("user_id")
	db := database.GetDB()

	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "用户不存在", "data": nil})
		return
	}

	if user.APIEnabled {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "API已开通", "data": nil})
		return
	}

	// 检查总开关
	var apiSwitch model.Setting
	if err := db.Where("`key` = ?", "api_enabled").First(&apiSwitch).Error; err == nil && apiSwitch.Value == "0" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "API功能已关闭", "data": nil})
		return
	}

	// 检查条件
	var phoneSetting model.Setting
	if err := db.Where("`key` = ?", "api_require_phone").First(&phoneSetting).Error; err == nil && phoneSetting.Value == "1" {
		if user.Phone == "" {
			c.JSON(http.StatusOK, gin.H{"code": 403, "message": "请先绑定手机号", "data": nil})
			return
		}
	}
	var realNameSetting model.Setting
	if err := db.Where("`key` = ?", "api_require_realname").First(&realNameSetting).Error; err == nil && realNameSetting.Value == "1" {
		if !user.IsVerified {
			c.JSON(http.StatusOK, gin.H{"code": 403, "message": "请先完成实名认证", "data": nil})
			return
		}
	}

	// 生成API密钥
	apiKey, _ := generateAPIKey()
	encrypted, encErr := util.EncryptAES(apiKey)
	if encErr != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "密钥加密失败", "data": nil})
		return
	}

	db.Model(&user).Update("api_key", encrypted)
	db.Model(&user).Update("api_enabled", true)

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "API开通成功",
		"data": gin.H{"api_key": apiKey},
	})
}

// ResetAPIKey 重置API密钥
// POST /api/client/api-key/reset
func ResetAPIKey(c *gin.Context) {
	userID, _ := c.Get("user_id")
	db := database.GetDB()

	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "用户不存在", "data": nil})
		return
	}

	if !user.APIEnabled {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "API未开通，请先开通", "data": nil})
		return
	}

	// 生成新密钥
	apiKey, _ := generateAPIKey()
	encrypted, encErr := util.EncryptAES(apiKey)
	if encErr != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "密钥加密失败", "data": nil})
		return
	}

	db.Model(&user).Update("api_key", encrypted)

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "API密钥已重置",
		"data": gin.H{"api_key": apiKey},
	})
}

// DisableAPIKey 关闭API功能
// POST /api/client/api-key/disable
func DisableAPIKey(c *gin.Context) {
	userID, _ := c.Get("user_id")
	db := database.GetDB()

	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "用户不存在", "data": nil})
		return
	}

	db.Model(&user).Update("api_enabled", false)
	db.Model(&user).Update("api_key", "")

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "API已关闭", "data": nil})
}

// generateAPIKey 生成64位随机API密钥
func generateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
