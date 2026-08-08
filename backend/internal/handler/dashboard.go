package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetDashboardStats 获取仪表盘统计
func GetDashboardStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"total_clients":     0,
		"total_orders":      0,
		"total_tickets":     0,
		"total_revenue":     0,
		"pending_tickets":   0,
		"unpaid_orders":     0,
		"new_clients_today": 0,
		"revenue_today":     0,
	})
}

// GetIncomeData 获取收入数据
func GetIncomeData(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"months":  []string{},
		"income":  []float64{},
		"expense": []float64{},
	})
}

// GetRecentOrders 获取最近订单
func GetRecentOrders(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"list": []interface{}{}})
}

// GetRecentTickets 获取最近工单
func GetRecentTickets(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"list": []interface{}{}})
}
