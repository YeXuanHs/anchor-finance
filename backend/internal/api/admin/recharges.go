package admin

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetRechargeList 获取充值记录列表
// GET /api/admin/finance/recharges
func GetRechargeList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")
	userID := c.Query("user_id")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetDB()
	query := db.Model(&model.Recharge{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	var total int64
	query.Count(&total)

	var recharges []model.Recharge
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&recharges)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      recharges,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetRechargeSummary 获取充值统计
// GET /api/admin/finance/recharges/summary
func GetRechargeSummary(c *gin.Context) {
	db := database.GetDB()

	// 统计总充值金额
	var totalAmount float64
	db.Model(&model.Recharge{}).Where("status = ?", "success").
		Select("COALESCE(SUM(amount), 0)").Scan(&totalAmount)

	// 统计本月充值金额
	var monthlyAmount float64
	db.Model(&model.Recharge{}).Where("status = ? AND created_at >= DATE_FORMAT(NOW(), '%Y-%m-01')", "success").
		Select("COALESCE(SUM(amount), 0)").Scan(&monthlyAmount)

	// 统计充值次数
	var totalCount int64
	db.Model(&model.Recharge{}).Where("status = ?", "success").Count(&totalCount)

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
