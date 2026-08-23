package admin

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetCurrencyList 获取货币列表
// GET /api/admin/currencies
func GetCurrencyList(c *gin.Context) {
	db := database.GetDB()
	var currencies []model.Currency
	db.Order("is_default DESC, code ASC").Find(&currencies)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    currencies,
	})
}

// CreateCurrency 创建货币
// POST /api/admin/currencies
func CreateCurrency(c *gin.Context) {
	var req struct {
		Code      string  `json:"code" binding:"required"`
		Name      string  `json:"name" binding:"required"`
		Symbol    string  `json:"symbol"`
		Rate      float64 `json:"rate"`
		IsDefault bool    `json:"is_default"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	db := database.GetDB()

	// 如果设为默认，取消其他默认
	if req.IsDefault {
		db.Model(&model.Currency{}).Where("is_default = ?", true).Update("is_default", false)
	}

	if req.Rate == 0 {
		req.Rate = 1
	}

	currency := model.Currency{
		Code:      req.Code,
		Name:      req.Name,
		Symbol:    req.Symbol,
		Rate:      req.Rate,
		IsDefault: req.IsDefault,
		Status:    "active",
	}

	if err := db.Create(&currency).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建货币失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": currency.ID,
		},
	})
}

// UpdateCurrency 更新货币
// PUT /api/admin/currencies/:id
func UpdateCurrency(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的货币ID",
			"data":    nil,
		})
		return
	}

	var req struct {
		Name      string  `json:"name"`
		Symbol    string  `json:"symbol"`
		Rate      float64 `json:"rate"`
		IsDefault *bool   `json:"is_default"`
		Status    string  `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	db := database.GetDB()
	var currency model.Currency
	if err := db.First(&currency, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "货币不存在",
			"data":    nil,
		})
		return
	}

	// 如果设为默认，取消其他默认
	if req.IsDefault != nil && *req.IsDefault {
		db.Model(&model.Currency{}).Where("is_default = ? AND id != ?", true, id).Update("is_default", false)
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Symbol != "" {
		updates["symbol"] = req.Symbol
	}
	if req.Rate > 0 {
		updates["rate"] = req.Rate
	}
	if req.IsDefault != nil {
		updates["is_default"] = *req.IsDefault
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	db.Model(&currency).Updates(updates)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
		"data":    nil,
	})
}
