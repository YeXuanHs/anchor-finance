package admin

import (
	"net/http"
	"strconv"
	"time"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetSalesStatistics 获取销售统计
// GET /api/admin/sales/statistics
func GetSalesStatistics(c *gin.Context) {
	db := database.GetDB()

	// 统计本月销售额
	now := time.Now()
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	var monthlySales float64
	db.Model(&model.Order{}).
		Where("status = ? AND created_at >= ?", "paid", firstOfMonth).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&monthlySales)

	// 统计本月订单数
	var monthlyOrderCount int64
	db.Model(&model.Order{}).
		Where("created_at >= ?", firstOfMonth).
		Count(&monthlyOrderCount)

	// 统计本月新客户数
	var monthlyCustomerCount int64
	db.Model(&model.User{}).
		Where("created_at >= ?", firstOfMonth).
		Count(&monthlyCustomerCount)

	// 历年总销售额
	var totalSales float64
	db.Model(&model.Order{}).
		Where("status = ?", "paid").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalSales)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"monthly_sales":     monthlySales,
			"monthly_orders":    monthlyOrderCount,
			"monthly_customers": monthlyCustomerCount,
			"total_sales":       totalSales,
		},
	})
}

// GetSalesRecords 获取销售记录
// GET /api/admin/sales/records
func GetSalesRecords(c *gin.Context) {
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
	db.Model(&model.Order{}).Where("status = ?", "paid").Count(&total)

	var orders []model.Order
	offset := (page - 1) * pageSize
	db.Where("status = ?", "paid").
		Offset(offset).
		Limit(pageSize).
		Order("id DESC").
		Find(&orders)

	if orders == nil {
		orders = []model.Order{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      orders,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
