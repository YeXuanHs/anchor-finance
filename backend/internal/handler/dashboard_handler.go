package handler

import (
	"net/http"
	"time"

	"anchorfinance/pkg/db"

	"github.com/gin-gonic/gin"
)

// DashboardHandler handles dashboard-related requests.
type DashboardHandler struct{}

// GetIncomeTrend returns income trend data for the last 12 months.
func (h *DashboardHandler) GetIncomeTrend(c *gin.Context) {
	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
		return
	}

	type MonthlyIncome struct {
		Month  string  `json:"month"`
		Amount float64 `json:"amount"`
	}

	var results []MonthlyIncome
	now := time.Now()

	for i := 11; i >= 0; i-- {
		month := now.AddDate(0, -i, 0)
		startOfMonth := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
		endOfMonth := startOfMonth.AddDate(0, 1, 0)

		var total float64
		database.Model(&struct {
			Amount float64
		}{}).Table("transactions").
			Where("created_at >= ? AND created_at < ? AND type = ?", startOfMonth, endOfMonth, "income").
			Select("COALESCE(SUM(amount), 0)").
			Scan(&total)

		results = append(results, MonthlyIncome{
			Month:  month.Format("2006-01"),
			Amount: total,
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": results})
}

// GetProductDistribution returns product distribution data.
func (h *DashboardHandler) GetProductDistribution(c *gin.Context) {
	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
		return
	}

	type ProductCount struct {
		Name  string `json:"name"`
		Count int64  `json:"count"`
	}

	var results []ProductCount
	database.Table("products").
		Select("name, COUNT(*) as count").
		Group("name").
		Order("count DESC").
		Limit(10).
		Scan(&results)

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": results})
}

// GlobalSearch performs global search across users, orders, tickets.
func (h *DashboardHandler) GlobalSearch(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{}})
		return
	}

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{}})
		return
	}

	type SearchResult struct {
		Type string      `json:"type"`
		ID   uint        `json:"id"`
		Name string      `json:"name"`
		Data interface{} `json:"data"`
	}

	var results []SearchResult

	// Search users
	var users []struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	database.Table("users").
		Where("username LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Limit(5).
		Scan(&users)
	for _, u := range users {
		results = append(results, SearchResult{
			Type: "user",
			ID:   u.ID,
			Name: u.Username,
			Data: u,
		})
	}

	// Search orders
	var orders []struct {
		ID     uint   `json:"id"`
		UserID uint   `json:"user_id"`
		Status string `json:"status"`
	}
	database.Table("orders").
		Where("CAST(id AS CHAR) LIKE ?", "%"+keyword+"%").
		Limit(5).
		Scan(&orders)
	for _, o := range orders {
		results = append(results, SearchResult{
			Type: "order",
			ID:   o.ID,
			Name: "Order #" + string(rune(o.ID)),
			Data: o,
		})
	}

	// Search tickets
	var tickets []struct {
		ID     uint   `json:"id"`
		UserID uint   `json:"user_id"`
		Subject string `json:"subject"`
	}
	database.Table("tickets").
		Where("subject LIKE ?", "%"+keyword+"%").
		Limit(5).
		Scan(&tickets)
	for _, t := range tickets {
		results = append(results, SearchResult{
			Type: "ticket",
			ID:   t.ID,
			Name: t.Subject,
			Data: t,
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": results})
}

// GetAdminIndex returns admin dashboard overview data.
func (h *DashboardHandler) GetAdminIndex(c *gin.Context) {
	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{}})
		return
	}

	// Count users
	var userCount int64
	database.Table("users").Count(&userCount)

	// Count orders
	var orderCount int64
	database.Table("orders").Count(&orderCount)

	// Count tickets
	var ticketCount int64
	database.Table("tickets").Count(&ticketCount)

	// Count products
	var productCount int64
	database.Table("products").Count(&productCount)

	// Get total income
	var totalIncome float64
	database.Table("transactions").
		Where("type = ?", "income").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalIncome)

	// Get recent orders
	var recentOrders []struct {
		ID        uint      `json:"id"`
		UserID    uint      `json:"user_id"`
		Status    string    `json:"status"`
		Amount    float64   `json:"amount"`
		CreatedAt time.Time `json:"created_at"`
	}
	database.Table("orders").
		Order("created_at DESC").
		Limit(5).
		Scan(&recentOrders)

	// Get recent tickets
	var recentTickets []struct {
		ID        uint      `json:"id"`
		UserID    uint      `json:"user_id"`
		Subject   string    `json:"subject"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
	}
	database.Table("tickets").
		Order("created_at DESC").
		Limit(5).
		Scan(&recentTickets)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"user_count":     userCount,
			"order_count":    orderCount,
			"ticket_count":   ticketCount,
			"product_count":  productCount,
			"total_income":   totalIncome,
			"recent_orders":  recentOrders,
			"recent_tickets": recentTickets,
		},
	})
}

// ExpiringProduct represents a product expiring soon.
type ExpiringProduct struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	ProductName string  `json:"product_name"`
	ExpiresAt time.Time `json:"expires_at"`
}

// GetExpiringProducts returns products expiring within 7 days.
func (h *DashboardHandler) GetExpiringProducts(c *gin.Context) {
	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
		return
	}

	var results []ExpiringProduct
	sevenDaysLater := time.Now().AddDate(0, 0, 7)

	database.Raw(`
		SELECT o.id, o.user_id, COALESCE(u.username, '') as username,
		       COALESCE(p.name, '') as product_name, o.expires_at
		FROM orders o
		LEFT JOIN users u ON o.user_id = u.id
		LEFT JOIN products p ON o.product_id = p.id
		WHERE o.status = 'active' AND o.expires_at IS NOT NULL
		AND o.expires_at <= ? AND o.expires_at >= ?
		ORDER BY o.expires_at ASC
		LIMIT 20
	`, sevenDaysLater, time.Now()).
		Scan(&results)

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": results})
}

// OnlineClient represents an online client.
type OnlineClient struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	IP       string `json:"ip"`
	LastSeen string `json:"last_seen"`
}

// GetOnlineClients returns currently online clients.
func (h *DashboardHandler) GetOnlineClients(c *gin.Context) {
	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
		return
	}

	var results []OnlineClient

	// Try to get from a user_activity or sessions table if it exists
	type tableInfo struct {
		Name string
	}
	var tables []tableInfo
	database.Raw("SHOW TABLES LIKE 'user_sessions'").Scan(&tables)

	if len(tables) > 0 {
		database.Raw(`
			SELECT s.user_id, COALESCE(u.username, '') as username,
			       COALESCE(s.ip, '') as ip, COALESCE(s.last_active, '') as last_seen
			FROM user_sessions s
			LEFT JOIN users u ON s.user_id = u.id
			WHERE s.last_active > DATE_SUB(NOW(), INTERVAL 30 MINUTE)
			ORDER BY s.last_active DESC
			LIMIT 50
		`).Scan(&results)
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": results})
}
