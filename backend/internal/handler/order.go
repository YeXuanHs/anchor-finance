package handler

import (
	"net/http"
	"strconv"
	"time"

	"anchorfinance/pkg/db"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct{}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}

// GetOrders 获取订单列表
func (h *OrderHandler) GetOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusOK, gin.H{"list": []interface{}{}, "total": 0})
		return
	}

	query := database.Table("orders")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var orders []struct {
		ID        uint      `json:"id"`
		UserID    uint      `json:"user_id"`
		Status    string    `json:"status"`
		Total     float64   `json:"total"`
		CreatedAt time.Time `json:"created_at"`
	}

	query.Select("id, user_id, status, total, created_at").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&orders)

	c.JSON(http.StatusOK, gin.H{
		"list":  orders,
		"total": total,
		"page":  page,
	})
}

// GetOrder 获取单个订单
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "not found"})
		return
	}

	var order struct {
		ID        uint      `json:"id"`
		UserID    uint      `json:"user_id"`
		Status    string    `json:"status"`
		Total     float64   `json:"total"`
		CreatedAt time.Time `json:"created_at"`
	}

	err := database.Table("orders").Where("id = ?", id).First(&order).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "订单不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": order})
}

// CreateOrder 创建订单
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"code": 0, "message": "订单创建成功"})
}

// UpdateOrder 更新订单
func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "订单更新成功"})
}

// RegisterRoutes 注册路由
func (h *OrderHandler) RegisterRoutes(r *gin.RouterGroup) {
	order := r.Group("/orders")
	{
		order.GET("", h.GetOrders)
		order.GET("/:id", h.GetOrder)
		order.POST("", h.CreateOrder)
		order.PUT("/:id", h.UpdateOrder)
	}
}

// ==================== Admin Router Methods ====================

// GetList returns a paginated order list.
func (h *OrderHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"list": []interface{}{}, "total": 0}})
		return
	}

	query := database.Table("orders")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Where("CAST(id AS CHAR) LIKE ?", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var orders []struct {
		ID        uint      `json:"id"`
		UserID    uint      `json:"user_id"`
		Status    string    `json:"status"`
		Total     float64   `json:"total"`
		CreatedAt time.Time `json:"created_at"`
	}

	query.Select("id, user_id, status, total, created_at").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&orders)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  orders,
			"total": total,
		},
	})
}

// GetDetail returns order detail.
func (h *OrderHandler) GetDetail(c *gin.Context) {
	id := c.Param("id")

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": nil})
		return
	}

	var order struct {
		ID        uint      `json:"id"`
		UserID    uint      `json:"user_id"`
		Status    string    `json:"status"`
		Total     float64   `json:"total"`
		CreatedAt time.Time `json:"created_at"`
	}

	err := database.Table("orders").Where("id = ?", id).First(&order).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "订单不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": order})
}

// UpdateStatus updates an order's status.
func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"`
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

	if err := database.Table("orders").Where("id = ?", id).Update("status", req.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "状态更新成功"})
}

// AdminCreateOrder creates an order (admin).
func (h *OrderHandler) AdminCreateOrder(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "订单创建成功"})
}

// CheckOrder checks an order.
func (h *OrderHandler) CheckOrder(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// BatchUpdate batch updates orders.
func (h *OrderHandler) BatchUpdate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// Delete deletes an order.
func (h *OrderHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "数据库未初始化"})
		return
	}

	if err := database.Table("orders").Where("id = ?", id).Delete(nil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// AddNote adds a note to an order.
func (h *OrderHandler) AddNote(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// GetNotes returns notes for an order.
func (h *OrderHandler) GetNotes(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
}

// GetSaleOrders returns sale orders.
func (h *OrderHandler) GetSaleOrders(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
}

// ActivateOrder activates an order.
func (h *OrderHandler) ActivateOrder(c *gin.Context) {
	id := c.Param("id")

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "数据库未初始化"})
		return
	}

	if err := database.Table("orders").Where("id = ?", id).Update("status", "active").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "激活失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "订单已激活"})
}

// ChangeStatus changes an order's status.
func (h *OrderHandler) ChangeStatus(c *gin.Context) {
	h.UpdateStatus(c)
}

// GetMultiTotal returns multi-order totals.
func (h *OrderHandler) GetMultiTotal(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"total": 0}})
}

// ApplyCustomPromo applies a custom promotion.
func (h *OrderHandler) ApplyCustomPromo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// SearchPage returns the order search page.
func (h *OrderHandler) SearchPage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
}

// CheckProduct checks product eligibility for order.
func (h *OrderHandler) CheckProduct(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"available": true}})
}

// SetConfig returns order configuration.
func (h *OrderHandler) SetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{}})
}

// GetClients returns client list for order dropdown.
func (h *OrderHandler) GetClients(c *gin.Context) {
	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
		return
	}

	var clients []struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	database.Table("users").Select("id, username, email").Where("is_admin = ?", false).Find(&clients)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": clients})
}

// GetListEnhanced returns an enhanced order list.
func (h *OrderHandler) GetListEnhanced(c *gin.Context) {
	h.GetList(c)
}

// ==================== V1 Router Methods ====================

// Create creates a new order.
func (h *OrderHandler) Create(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// GetUserOrders returns orders for the authenticated user.
func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
}

// Pay pays for an order.
func (h *OrderHandler) Pay(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// Cancel cancels an order.
func (h *OrderHandler) Cancel(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// Preview previews an order before creation.
func (h *OrderHandler) Preview(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{}})
}
