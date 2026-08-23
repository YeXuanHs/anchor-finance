package admin

import (
	crypto_rand "crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// AdminEnableUserAPI 管理员强制开通用户API（无视条件）
// POST /api/admin/users/:id/api-key/enable
func AdminEnableUserAPI(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var user model.User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "用户不存在", "data": nil})
		return
	}

	if user.APIEnabled {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "API已开通", "data": nil})
		return
	}

	apiKey, _ := generateAPIKeyLocal()

	db.Model(&user).Update("api_key", apiKey)
	db.Model(&user).Update("api_enabled", true)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "API已强制开通", "data": gin.H{"api_key": apiKey}})
}

// AdminResetUserAPI 管理员重置用户API密钥
// POST /api/admin/users/:id/api-key/reset
func AdminResetUserAPI(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var user model.User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "用户不存在", "data": nil})
		return
	}

	if !user.APIEnabled {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "用户API未开通", "data": nil})
		return
	}

	apiKey, _ := generateAPIKeyLocal()
	db.Model(&user).Update("api_key", apiKey)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "API密钥已重置", "data": gin.H{"api_key": apiKey}})
}

// AdminDisableUserAPI 管理员关闭用户API
// POST /api/admin/users/:id/api-key/disable
func AdminDisableUserAPI(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var user model.User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "用户不存在", "data": nil})
		return
	}

	db.Model(&user).Updates(map[string]interface{}{
		"api_enabled": false,
		"api_key":     "",
	})

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "API已关闭", "data": nil})
}

// GetAPISettings 获取API开放条件配置
// GET /api/admin/settings/api
func GetAPISettings(c *gin.Context) {
	db := database.GetDB()
	keys := []string{"api_enabled", "api_require_phone", "api_require_realname"}
	settings := make(map[string]string)
	for _, key := range keys {
		var s model.Setting
		if err := db.Where("`key` = ?", key).First(&s).Error; err == nil {
			settings[key] = s.Value
		} else {
			// 默认值
			if key == "api_enabled" {
				settings[key] = "1"
			} else {
				settings[key] = "0"
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": settings})
}

// UpdateAPISettings 更新API开放条件配置
// PUT /api/admin/settings/api
func UpdateAPISettings(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	allowed := map[string]bool{"api_enabled": true, "api_require_phone": true, "api_require_realname": true}
	for key, value := range req {
		if !allowed[key] {
			continue
		}
		var existing model.Setting
		if err := db.Where("`key` = ?", key).First(&existing).Error; err == nil {
			db.Model(&existing).Update("value", value)
		} else {
			db.Create(&model.Setting{Key: key, Value: value, Group: "api"})
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": nil})
}

func generateAPIKeyLocal() (string, error) {
	b := make([]byte, 32)
	if _, err := crypto_rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
