package admin

import (
	"net/http"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetOrderConfig 获取商品订购设置
// GET /api/admin/order-config
func GetOrderConfig(c *gin.Context) {
	db := database.GetDB()
	var configs []model.OrderConfig
	db.Where("`group` = ?", "order").Find(&configs)

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

// UpdateOrderConfig 更新商品订购设置
// PUT /api/admin/order-config
func UpdateOrderConfig(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	for key, value := range req {
		db.Where("key = ?", key).Assign(map[string]interface{}{"value": value}).FirstOrCreate(&model.OrderConfig{})
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": nil})
}

// GetFinanceConfig 获取财务配置
// GET /api/admin/finance-config
func GetFinanceConfig(c *gin.Context) {
	db := database.GetDB()
	var configs []model.OrderConfig
	db.Where("`group` = ?", "finance").Find(&configs)

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

// UpdateFinanceConfig 更新财务配置
// PUT /api/admin/finance-config
func UpdateFinanceConfig(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	for key, value := range req {
		db.Where("key = ?", key).Assign(map[string]interface{}{"value": value}).FirstOrCreate(&model.OrderConfig{})
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": nil})
}
