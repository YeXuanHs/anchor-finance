package client

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetPaymentList 获取支付记录列表
// GET /api/client/payments
func GetPaymentList(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetDB()
	query := db.Model(&model.Payment{}).Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var payments []model.Payment
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&payments)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      payments,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetPaymentSummary 获取支付统计
// GET /api/client/payments/summary
func GetPaymentSummary(c *gin.Context) {
	userID, _ := c.Get("user_id")

	db := database.GetDB()

	// 统计总支付金额
	var totalAmount float64
	db.Model(&model.Payment{}).Where("user_id = ? AND status = ?", userID, "success").
		Select("COALESCE(SUM(amount), 0)").Scan(&totalAmount)

	// 统计本月支付金额
	var monthlyAmount float64
	db.Model(&model.Payment{}).Where("user_id = ? AND status = ? AND created_at >= DATE_FORMAT(NOW(), '%Y-%m-01')", userID, "success").
		Select("COALESCE(SUM(amount), 0)").Scan(&monthlyAmount)

	// 统计支付次数
	var totalCount int64
	db.Model(&model.Payment{}).Where("user_id = ? AND status = ?", userID, "success").Count(&totalCount)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"total_amount":   totalAmount,
			"monthly_amount": monthlyAmount,
			"total_count":    totalCount,
		},
	})
}

// GetPaymentDetail 获取支付详情
// GET /api/client/payments/:id
func GetPaymentDetail(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的支付ID",,
			"data": nil})
		return
	}

	db := database.GetDB()
	var payment model.Payment
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&payment).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "支付记录不存在",,
			"data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    payment,
	})
}
