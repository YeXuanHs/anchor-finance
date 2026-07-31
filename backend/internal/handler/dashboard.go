package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DashboardHandler 仪表盘数据
type DashboardHandler struct {
	db *gorm.DB
}

func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

// GetStats 获取统计数据
// GET /admin/dashboard/stats
func (h *DashboardHandler) GetStats(c *gin.Context) {
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	// 今日收入
	var todayIncome float64
	h.db.Raw(`
		SELECT COALESCE(SUM(amount), 0) FROM orders 
		WHERE DATE(created_at) = ? AND status IN (3, 5)
	`, today).Scan(&todayIncome)

	// 昨日收入
	var yesterdayIncome float64
	h.db.Raw(`
		SELECT COALESCE(SUM(amount), 0) FROM orders 
		WHERE DATE(created_at) = ? AND status IN (3, 5)
	`, yesterday).Scan(&yesterdayIncome)

	// 收入变化百分比
	incomeChange := 0.0
	if yesterdayIncome > 0 {
		incomeChange = ((todayIncome - yesterdayIncome) / yesterdayIncome) * 100
	}

	// 今日新增客户
	var newClients int64
	h.db.Raw(`SELECT COUNT(*) FROM users WHERE DATE(created_at) = ?`, today).Scan(&newClients)

	// 昨日新增客户
	var yesterdayClients int64
	h.db.Raw(`SELECT COUNT(*) FROM users WHERE DATE(created_at) = ?`, yesterday).Scan(&yesterdayClients)

	// 客户变化百分比
	clientChange := 0.0
	if yesterdayClients > 0 {
		clientChange = (float64(newClients-yesterdayClients) / float64(yesterdayClients)) * 100
	}

	// 待处理工单
	var pendingTickets int64
	h.db.Raw(`SELECT COUNT(*) FROM tickets WHERE status IN (0, 2)`).Scan(&pendingTickets)

	// 紧急工单
	var urgentTickets int64
	h.db.Raw(`SELECT COUNT(*) FROM tickets WHERE priority = 4 AND status IN (0, 2)`).Scan(&urgentTickets)

	// 待处理订单
	var pendingOrders int64
	h.db.Raw(`SELECT COUNT(*) FROM orders WHERE status IN (0, 1)`).Scan(&pendingOrders)

	// 待审核订单
	var reviewOrders int64
	h.db.Raw(`SELECT COUNT(*) FROM orders WHERE status = 1`).Scan(&reviewOrders)

	c.JSON(200, gin.H{
		"code": 0,
		"data": gin.H{
			"today_income":   todayIncome,
			"income_change":  incomeChange,
			"new_clients":    newClients,
			"client_change":  clientChange,
			"pending_tickets": pendingTickets,
			"urgent_tickets": urgentTickets,
			"pending_orders": pendingOrders,
			"review_orders":  reviewOrders,
		},
	})
}

// GetIncomeTrend 获取收入趋势
// GET /admin/dashboard/income-trend
func (h *DashboardHandler) GetIncomeTrend(c *gin.Context) {
	period := c.DefaultQuery("period", "week")

	var days int
	switch period {
	case "month":
		days = 30
	case "year":
		days = 12
	default:
		days = 7
	}

	type IncomeData struct {
		Dates  []string  `json:"dates"`
		Values []float64 `json:"values"`
	}

	result := IncomeData{
		Dates:  make([]string, days),
		Values: make([]float64, days),
	}

	if period == "year" {
		// 按月查询
		rows, _ := h.db.Raw(`
			SELECT DATE_FORMAT(created_at, '%Y-%m') as month, COALESCE(SUM(amount), 0) as total
			FROM orders 
			WHERE created_at >= DATE_SUB(NOW(), INTERVAL ? MONTH) AND status IN (3, 5)
			GROUP BY DATE_FORMAT(created_at, '%Y-%m')
			ORDER BY month
		`, days).Rows()

		defer rows.Close()

		monthData := make(map[string]float64)
		for rows.Next() {
			var month string
			var total float64
			rows.Scan(&month, &total)
			monthData[month] = total
		}

		for i := 0; i < days; i++ {
			date := time.Now().AddDate(0, -(days - 1 - i), 0)
			monthStr := date.Format("2006-01")
			result.Dates[i] = date.Format("1月")
			result.Values[i] = monthData[monthStr]
		}
	} else {
		// 按天查询
		rows, _ := h.db.Raw(`
			SELECT DATE(created_at) as date, COALESCE(SUM(amount), 0) as total
			FROM orders 
			WHERE created_at >= DATE_SUB(NOW(), INTERVAL ? DAY) AND status IN (3, 5)
			GROUP BY DATE(created_at)
			ORDER BY date
		`, days).Rows()

		defer rows.Close()

		dateData := make(map[string]float64)
		for rows.Next() {
			var date string
			var total float64
			rows.Scan(&date, &total)
			dateData[date] = total
		}

		for i := 0; i < days; i++ {
			date := time.Now().AddDate(0, 0, -(days - 1 - i))
			dateStr := date.Format("2006-01-02")
			result.Dates[i] = date.Format("01/02")
			result.Values[i] = dateData[dateStr]
		}
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": result,
	})
}

// GetProductDistribution 获取产品分布
// GET /admin/dashboard/product-distribution
func (h *DashboardHandler) GetProductDistribution(c *gin.Context) {
	type ProductStat struct {
		Name  string `json:"name"`
		Value int64  `json:"value"`
	}

	var stats []ProductStat
	h.db.Raw(`
		SELECT p.name, COUNT(*) as value
		FROM orders o
		JOIN products p ON o.product_id = p.id
		WHERE o.status IN (3, 5)
		GROUP BY p.name
		ORDER BY value DESC
		LIMIT 7
	`).Scan(&stats)

	c.JSON(200, gin.H{
		"code": 0,
		"data": stats,
	})
}
