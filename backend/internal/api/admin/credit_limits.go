package admin

import (
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

// GetCreditLimitConfig 获取信用额配置
// GET /api/admin/credit-limits/config
func GetCreditLimitConfig(c *gin.Context) {
	// 返回信用额配置
	config := gin.H{
		"enabled":        true,
		"default_amount": 0,
		"max_amount":     100000,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    config,
	})
}

// SaveCreditLimitConfig 保存信用额配置
// POST /api/admin/credit-limits/config
func SaveCreditLimitConfig(c *gin.Context) {
	var req struct {
		Enabled       bool    `json:"enabled"`
		DefaultAmount float64 `json:"default_amount"`
		MaxAmount     float64 `json:"max_amount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// TODO: 保存到settings表

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "保存成功",
	})
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
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	db := database.GetDB()

	// 检查用户是否存在
	var user model.User
	if err := db.First(&user, req.UserID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "用户不存在",
		})
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
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	db := database.GetDB()
	var limit model.CreditLimit
	if err := db.First(&limit, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "信用额度不存在",
		})
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
			"message": "信用额度不存在",
		})
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
