package admin

import (
	"net/http"
	"time"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetDashboardStats 获取仪表盘统计
// GET /api/admin/dashboard/stats
func GetDashboardStats(c *gin.Context) {
	db := database.GetDB()

	// 统计用户数
	var userCount int64
	db.Model(&model.User{}).Count(&userCount)

	// 统计订单数
	var orderCount int64
	db.Model(&model.Order{}).Count(&orderCount)

	// 统计服务数
	var serviceCount int64
	db.Model(&model.Service{}).Where("status = ?", "active").Count(&serviceCount)

	// 统计工单数
	var ticketCount int64
	db.Model(&model.Ticket{}).Where("status != ?", "closed").Count(&ticketCount)

	// 统计本月收入
	var monthlyIncome float64
	now := time.Now()
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	db.Model(&model.Order{}).
		Where("status = ? AND created_at >= ?", "paid", firstOfMonth).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&monthlyIncome)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"user_count":     userCount,
			"order_count":    orderCount,
			"service_count":  serviceCount,
			"ticket_count":   ticketCount,
			"monthly_income": monthlyIncome,
		},
	})
}

// GetIncomeTrend 获取收入趋势
// GET /api/admin/dashboard/income-trend
func GetIncomeTrend(c *gin.Context) {
	db := database.GetDB()

	// 获取最近12个月的收入
	type MonthlyData struct {
		Month  string  `json:"month"`
		Amount float64 `json:"amount"`
	}

	var results []MonthlyData
	now := time.Now()

	for i := 11; i >= 0; i-- {
		month := now.AddDate(0, -i, 0)
		firstOfMonth := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

		var amount float64
		db.Model(&model.Order{}).
			Where("status = ? AND created_at >= ? AND created_at <= ?", "paid", firstOfMonth, lastOfMonth).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&amount)

		results = append(results, MonthlyData{
			Month:  firstOfMonth.Format("2006-01"),
			Amount: amount,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    results,
	})
}

// GetOnlineAdmins 获取在线管理员
// GET /api/admin/dashboard/online-admins
func GetOnlineAdmins(c *gin.Context) {
	// 暂时返回最近登录的管理员
	db := database.GetDB()
	var admins []model.Admin
	db.Where("last_login_at IS NOT NULL").
		Order("last_login_at DESC").
		Limit(10).
		Find(&admins)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    admins,
	})
}
