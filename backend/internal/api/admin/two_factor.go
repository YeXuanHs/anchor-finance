package admin

import (
	"net/http"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetTwoFactorConfig 获取二次验证配置
// GET /api/admin/two-factor-config
func GetTwoFactorConfig(c *gin.Context) {
	db := database.GetDB()
	var configs []model.TwoFactorConfig
	db.Where("`group` = ?", "two_factor").Find(&configs)

	result := make(map[string]interface{})
	for _, config := range configs {
		result[config.Key] = gin.H{
			"value":   config.Value,
			"name":    config.Name,
			"type":    config.Type,
			"options": config.Options,
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": result})
}

// UpdateTwoFactorConfig 更新二次验证配置
// PUT /api/admin/two-factor-config
func UpdateTwoFactorConfig(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	for key, value := range req {
		db.Where("key = ?", key).Assign(map[string]interface{}{"value": value}).FirstOrCreate(&model.TwoFactorConfig{})
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": nil})
}
