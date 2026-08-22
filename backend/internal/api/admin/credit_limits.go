package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetCreditLimitList 获取信用额度列表
// GET /api/admin/credit-limits
func GetCreditLimitList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetDB()
	var total int64
	db.Model(&model.CreditLimit{}).Count(&total)

	var limits []model.CreditLimit
	offset := (page - 1) * pageSize
	db.Offset(offset).Limit(pageSize).Order("id DESC").Find(&limits)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      limits,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetCreditLimitConfig 获取信用额配置（从settings表读）
// GET /api/admin/credit-limits/config
func GetCreditLimitConfig(c *gin.Context) {
	db := database.GetDB()
	var enabledSetting, defaultSetting, maxSetting model.Setting

	enabled := false
	db.Where("`group` = ? AND `key` = ?", "credit", "enabled").First(&enabledSetting)
	if enabledSetting.Value == "1" || enabledSetting.Value == "true" {
		enabled = true
	}

	defaultAmount := 0.0
	db.Where("`group` = ? AND `key` = ?", "credit", "default_amount").First(&defaultSetting)
	if v, err := strconv.ParseFloat(defaultSetting.Value, 64); err == nil {
		defaultAmount = v
	}

	maxAmount := 0.0
	db.Where("`group` = ? AND `key` = ?", "credit", "max_amount").First(&maxSetting)
	if v, err := strconv.ParseFloat(maxSetting.Value, 64); err == nil {
		maxAmount = v
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"enabled":        enabled,
			"default_amount": defaultAmount,
			"max_amount":     maxAmount,
		},
	})
}

// SaveCreditLimitConfig 保存信用额配置（存settings表）
// POST /api/admin/credit-limits/config
func SaveCreditLimitConfig(c *gin.Context) {
	var req struct {
		Enabled       bool    `json:"enabled"`
		DefaultAmount float64 `json:"default_amount"`
		MaxAmount     float64 `json:"max_amount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	// 0元购防护：最大额度必须>=0
	if req.MaxAmount < 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "最大额度不能为负数", "data": nil})
		return
	}

	db := database.GetDB()
	saveSetting := func(key, value string) {
		db.Where("key = ?", key).Assign(map[string]interface{}{"value": value, "group": "credit"}).FirstOrCreate(&model.Setting{})
	}

	saveSetting("credit_enabled", fmt.Sprintf("%v", req.Enabled))
	saveSetting("credit_default_amount", fmt.Sprintf("%v", req.DefaultAmount))
	saveSetting("credit_max_amount", fmt.Sprintf("%v", req.MaxAmount))

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "保存成功", "data": nil})
}

// SaveCreditLimit 设置用户信用额度
// POST /api/admin/credit-limits
func SaveCreditLimit(c *gin.Context) {
	var req struct {
		UserID uint    `json:"user_id" binding:"required"`
		Amount float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),,
			"data": nil})
		return
	}

	db := database.GetDB()

	// 检查用户是否存在
	var user model.User
	if err := db.First(&user, req.UserID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "用户不存在",,
			"data": nil})
		return
	}

	// 查找或创建信用额度
	var limit model.CreditLimit
	result := db.Where("user_id = ?", req.UserID).First(&limit)
	if result.Error != nil {
		// 创建新的
		limit = model.CreditLimit{
			UserID: req.UserID,
			Amount: req.Amount,
			Status: "active",
		}
		db.Create(&limit)
	} else {
		// 更新现有的
		db.Model(&limit).Update("amount", req.Amount)
	}

	// 记录日志
	log := model.CreditLimitLog{
		UserID: req.UserID,
		Type:   "add",
		Amount: req.Amount,
		Remark: "管理员设置信用额度",
	}
	db.Create(&log)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "设置成功",
	})
}

// UpdateCreditLimit 更新用户信用额度
// PUT /api/admin/credit-limits/:id
func UpdateCreditLimit(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Amount float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),,
			"data": nil})
		return
	}

	db := database.GetDB()
	var limit model.CreditLimit
	if err := db.First(&limit, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "信用额度不存在",,
			"data": nil})
		return
	}

	// 记录日志
	log := model.CreditLimitLog{
		UserID: limit.UserID,
		Type:   "add",
		Amount: req.Amount - limit.Amount,
		Remark: "管理员修改信用额度",
	}
	db.Create(&log)

	db.Model(&limit).Update("amount", req.Amount)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
	})
}

// DeleteCreditLimit 删除用户信用额度
// DELETE /api/admin/credit-limits/:id
func DeleteCreditLimit(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var limit model.CreditLimit
	if err := db.First(&limit, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "信用额度不存在",,
			"data": nil})
		return
	}

	db.Delete(&limit)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}

// GetCreditLimitLogs 获取信用额日志
// GET /api/admin/credit-limits/logs
func GetCreditLimitLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	userID := c.Query("user_id")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetDB()
	query := db.Model(&model.CreditLimitLog{})

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	var total int64
	query.Count(&total)

	var logs []model.CreditLimitLog
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      logs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
