package admin

import (
	"net/http"
	"time"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetNewCustomerDailySummary 新客户每日统计
// GET /api/admin/finance/new-customer-daily-summary
func GetNewCustomerDailySummary(c *gin.Context) {
	db := database.GetDB()

	type DailyData struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}

	var results []DailyData
	now := time.Now()

	// 最近30天
	for i := 29; i >= 0; i-- {
		day := now.AddDate(0, 0, -i)
		startOfDay := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
		endOfDay := startOfDay.AddDate(0, 0, 1)

		var count int64
		db.Model(&model.User{}).
			Where("created_at >= ? AND created_at < ?", startOfDay, endOfDay).
			Count(&count)

		results = append(results, DailyData{
			Date:  startOfDay.Format("2006-01-02"),
			Count: count,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    results,
	})
}

// GetProductIncomeSummary 产品收入统计
// GET /api/admin/finance/product-income-summary
func GetProductIncomeSummary(c *gin.Context) {
	db := database.GetDB()

	type ProductIncome struct {
		ProductName string  `json:"product_name"`
		TotalAmount float64 `json:"total_amount"`
		OrderCount  int64   `json:"order_count"`
	}

	var results []ProductIncome
	db.Model(&model.Order{}).
		Select("product_name, SUM(amount) as total_amount, COUNT(*) as order_count").
		Where("status = ?", "paid").
		Group("product_name").
		Order("total_amount DESC").
		Limit(20).
		Scan(&results)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    results,
	})
}

// GetFinanceLedger 财务账本
// GET /api/admin/finance/ledger
func GetFinanceLedger(c *gin.Context) {
	db := database.GetDB()

	// 统计总收入
	var totalIncome float64
	db.Model(&model.Invoice{}).
		Where("status = ?", "paid").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalIncome)

	// 统计本月收入
	now := time.Now()
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	var monthlyIncome float64
	db.Model(&model.Invoice{}).
		Where("status = ? AND created_at >= ?", "paid", firstOfMonth).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&monthlyIncome)

	// 统计待支付金额
	var unpaidAmount float64
	db.Model(&model.Invoice{}).
		Where("status = ?", "unpaid").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&unpaidAmount)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"total_income":   totalIncome,
			"monthly_income": monthlyIncome,
			"unpaid_amount":  unpaidAmount,
		},
	})
}
