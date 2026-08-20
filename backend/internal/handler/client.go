package handler

import (
	"anchorfinance/pkg/response"
	"net/http"
	"strconv"
	"time"

	"anchorfinance/pkg/db"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type ClientHandler struct{}

func NewClientHandler() *ClientHandler {
	return &ClientHandler{}
}

// ClientResponse represents a client in API responses.
type ClientResponse struct {
	ID          uint       `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	Phone       string     `json:"phone"`
	Nickname    string     `json:"nickname"`
	Avatar      string     `json:"avatar"`
	Status      int16      `json:"status"`
	Balance     float64    `json:"balance"`
	GroupID     uint       `json:"group_id"`
	CompanyName string     `json:"company"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

// GetClients 获取客户列表
func (h *ClientHandler) GetClients(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	status := c.Query("status")
	groupID := c.Query("group_id")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	database := db.GetDB()
	if database == nil {
		response.SuccessPage(c, []interface{}{}, 0, page, pageSize)
		return
	}

	query := database.Table("users").Where("is_admin = ?", false)

	if keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ? OR phone LIKE ? OR nickname LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if groupID != "" {
		query = query.Where("group_id = ?", groupID)
	}

	var total int64
	query.Count(&total)

	var clients []ClientResponse
	query.Select("id, username, email, phone, nickname, avatar, status, balance, group_id, company, last_login_at, created_at").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&clients)

	c.JSON(http.StatusOK, gin.H{
		"list":      clients,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetClient 获取单个客户
func (h *ClientHandler) GetClient(c *gin.Context) {
	id := c.Param("id")

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "not found"})
		return
	}

	var client struct {
		ID          uint       `json:"id"`
		Username    string     `json:"username"`
		Email       string     `json:"email"`
		Phone       string     `json:"phone"`
		Nickname    string     `json:"nickname"`
		Avatar      string     `json:"avatar"`
		Status      int16      `json:"status"`
		Balance     float64    `json:"balance"`
		CreditLimit float64    `json:"credit_limit"`
		GroupID     uint       `json:"group_id"`
		CompanyName string     `json:"company"`
		Address     string     `json:"address"`
		City        string     `json:"city"`
		Country     string     `json:"country"`
		Language    string     `json:"language"`
		Currency    string     `json:"currency"`
		LastLoginAt *time.Time `json:"last_login_at"`
		LastLoginIP string     `json:"last_login_ip"`
		CreatedAt   time.Time  `json:"created_at"`
	}

	err := database.Table("users").
		Where("id = ?", id).
		First(&client).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "客户不存在"})
		return
	}

	// 获取客户的产品数量
	var productCount int64
	database.Table("user_products").Where("user_id = ?", id).Count(&productCount)

	// 获取客户的订单数量
	var orderCount int64
	database.Table("orders").Where("user_id = ?", id).Count(&orderCount)

	// 获取客户的工单数量
	var ticketCount int64
	database.Table("tickets").Where("user_id = ?", id).Count(&ticketCount)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"client":        client,
			"product_count": productCount,
			"order_count":   orderCount,
			"ticket_count":  ticketCount,
		},
	})
}

// CreateClient 创建客户
func (h *ClientHandler) CreateClient(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Phone    string `json:"phone"`
		Nickname string `json:"nickname"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "数据库未初始化"})
		return
	}

	// 检查用户名是否已存在
	var count int64
	database.Table("users").Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "用户名已存在"})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "密码加密失败"})
		return
	}

	user := map[string]interface{}{
		"username": req.Username,
		"email":    req.Email,
		"password": string(hashedPassword),
		"phone":    req.Phone,
		"nickname": req.Nickname,
		"status":   1,
		"is_admin": false,
	}

	if err := database.Table("users").Create(user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "客户创建成功"})
}

// UpdateClient 更新客户
func (h *ClientHandler) UpdateClient(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Status   *int16 `json:"status"`
		GroupID  *uint  `json:"group_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "数据库未初始化"})
		return
	}

	updates := map[string]interface{}{}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.GroupID != nil {
		updates["group_id"] = *req.GroupID
	}

	if err := database.Table("users").Where("id = ?", id).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "客户更新成功"})
}

// DeleteClient 删除客户
func (h *ClientHandler) DeleteClient(c *gin.Context) {
	id := c.Param("id")

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "数据库未初始化"})
		return
	}

	if err := database.Table("users").Where("id = ?", id).Delete(nil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "客户删除成功"})
}

// RegisterRoutes 注册路由
func (h *ClientHandler) RegisterRoutes(r *gin.RouterGroup) {
	client := r.Group("/clients")
	{
		client.GET("", h.GetClients)
		client.GET("/:id", h.GetClient)
		client.POST("", h.CreateClient)
		client.PUT("/:id", h.UpdateClient)
		client.DELETE("/:id", h.DeleteClient)
	}
}
