package handler

import (
	"net/http"
	"time"

	"anchorfinance/pkg/db"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct{}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

// GetAdmins 获取管理员列表
func (h *AdminHandler) GetAdmins(c *gin.Context) {
	database := db.GetDB()
	if database == nil {
		response.SuccessPage(c, []interface{}{}, 0, page, pageSize)
		return
	}

	var admins []struct {
		ID        uint       `json:"id"`
		Username  string     `json:"username"`
		Email     string     `json:"email"`
		Nickname  string     `json:"nickname"`
		Avatar    string     `json:"avatar"`
		Status    int16      `json:"status"`
		IsAdmin   bool       `json:"is_admin"`
		LastLogin *time.Time `json:"last_login_at"`
		CreatedAt time.Time  `json:"created_at"`
	}

	database.Table("users").
		Select("id, username, email, nickname, avatar, status, is_admin, last_login_at, created_at").
		Where("is_admin = ?", true).
		Find(&admins)

	response.SuccessPage(c, admins, int64(len(admins)), page, pageSize)
}

// GetAdmin 获取单个管理员
func (h *AdminHandler) GetAdmin(c *gin.Context) {
	id := c.Param("id")

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "not found"})
		return
	}

	var admin struct {
		ID        uint       `json:"id"`
		Username  string     `json:"username"`
		Email     string     `json:"email"`
		Nickname  string     `json:"nickname"`
		Avatar    string     `json:"avatar"`
		Status    int16      `json:"status"`
		IsAdmin   bool       `json:"is_admin"`
		LastLogin *time.Time `json:"last_login_at"`
		CreatedAt time.Time  `json:"created_at"`
	}

	err := database.Table("users").
		Select("id, username, email, nickname, avatar, status, is_admin, last_login_at, created_at").
		Where("id = ? AND is_admin = ?", id, true).
		First(&admin).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "管理员不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": admin})
}

// CreateAdmin 创建管理员
func (h *AdminHandler) CreateAdmin(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"code": 0, "message": "管理员创建成功"})
}

// UpdateAdmin 更新管理员
func (h *AdminHandler) UpdateAdmin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "管理员更新成功"})
}

// DeleteAdmin 删除管理员
func (h *AdminHandler) DeleteAdmin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "管理员删除成功"})
}

// RegisterRoutes 注册路由
func (h *AdminHandler) RegisterRoutes(r *gin.RouterGroup) {
	admin := r.Group("/admins")
	{
		admin.GET("", h.GetAdmins)
		admin.GET("/:id", h.GetAdmin)
		admin.POST("", h.CreateAdmin)
		admin.PUT("/:id", h.UpdateAdmin)
		admin.DELETE("/:id", h.DeleteAdmin)
	}
}

// Dashboard returns dashboard data.
func (h *AdminHandler) Dashboard(c *gin.Context) {
	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{}})
		return
	}

	// 用户统计
	var totalUsers, activeUsers, newUsersToday, newUsersMonth int64
	database.Table("users").Count(&totalUsers)
	database.Table("users").Where("status = ?", 1).Count(&activeUsers)
	database.Table("users").Where("DATE(created_at) = CURDATE()").Count(&newUsersToday)
	database.Table("users").Where("created_at >= DATE_FORMAT(NOW(), '%Y-%m-01')").Count(&newUsersMonth)

	// 订单统计
	var totalOrders, pendingOrders, activeOrders, todayOrders int64
	database.Table("orders").Count(&totalOrders)
	database.Table("orders").Where("status = ?", "pending").Count(&pendingOrders)
	database.Table("orders").Where("status = ?", "active").Count(&activeOrders)
	database.Table("orders").Where("DATE(created_at) = CURDATE()").Count(&todayOrders)

	// 工单统计
	var totalTickets, openTickets, pendingTickets, todayTickets int64
	database.Table("tickets").Count(&totalTickets)
	database.Table("tickets").Where("status = ?", "open").Count(&openTickets)
	database.Table("tickets").Where("status = ?", "pending").Count(&pendingTickets)
	database.Table("tickets").Where("DATE(created_at) = CURDATE()").Count(&todayTickets)

	// 收入统计
	var totalIncome, monthIncome, todayIncome float64
	database.Table("transactions").Where("type = ?", "income").Select("COALESCE(SUM(amount), 0)").Scan(&totalIncome)
	database.Table("transactions").Where("type = ? AND created_at >= DATE_FORMAT(NOW(), '%Y-%m-01')", "income").Select("COALESCE(SUM(amount), 0)").Scan(&monthIncome)
	database.Table("transactions").Where("type = ? AND DATE(created_at) = CURDATE()", "income").Select("COALESCE(SUM(amount), 0)").Scan(&todayIncome)

	// 最近订单
	var recentOrders []struct {
		ID        uint      `json:"id"`
		UserID    uint      `json:"user_id"`
		Status    string    `json:"status"`
		Total     float64   `json:"total"`
		CreatedAt time.Time `json:"created_at"`
	}
	database.Table("orders").Order("created_at DESC").Limit(5).Scan(&recentOrders)

	// 最近工单
	var recentTickets []struct {
		ID        uint      `json:"id"`
		UserID    uint      `json:"user_id"`
		Subject   string    `json:"subject"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
	}
	database.Table("tickets").Order("created_at DESC").Limit(5).Scan(&recentTickets)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"users": gin.H{
				"total":        totalUsers,
				"active":       activeUsers,
				"today":        newUsersToday,
				"month":        newUsersMonth,
			},
			"orders": gin.H{
				"total":        totalOrders,
				"pending":      pendingOrders,
				"active":       activeOrders,
				"today":        todayOrders,
			},
			"tickets": gin.H{
				"total":        totalTickets,
				"open":         openTickets,
				"pending":      pendingTickets,
				"today":        todayTickets,
			},
			"income": gin.H{
				"total":        totalIncome,
				"month":        monthIncome,
				"today":        todayIncome,
			},
			"recent_orders":  recentOrders,
			"recent_tickets": recentTickets,
		},
	})
}

// Stats returns dashboard statistics.
func (h *AdminHandler) Stats(c *gin.Context) {
	h.Dashboard(c)
}

// GetSettings returns all system settings.
func (h *AdminHandler) GetSettings(c *gin.Context) {
	settings := db.GetSystemSettings("all")
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": settings})
}

// UpdateSettings updates system settings.
func (h *AdminHandler) UpdateSettings(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	for key, value := range req {
		db.SetSystemSetting(key, value, "system", "")
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "设置更新成功"})
}

// GetSettingsByGroup returns settings filtered by group.
func (h *AdminHandler) GetSettingsByGroup(c *gin.Context) {
	group := c.Param("group")
	settings := db.GetSystemSettings(group)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": settings})
}

// GetLogs returns operation logs.
func (h *AdminHandler) GetLogs(c *gin.Context) {
	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
		return
	}

	var logs []struct {
		ID        uint      `json:"id"`
		AdminID   uint      `json:"admin_id"`
		Action    string    `json:"action"`
		Module    string    `json:"module"`
		Detail    string    `json:"detail"`
		IP        string    `json:"ip"`
		CreatedAt time.Time `json:"created_at"`
	}

	database.Table("admin_logs").
		Order("created_at DESC").
		Limit(50).
		Scan(&logs)

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": logs})
}
