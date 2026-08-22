package client

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetUserOrders 获取用户订单列表
// GET /api/client/orders
func GetUserOrders(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetDB()
	query := db.Model(&model.Order{}).Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var orders []model.Order
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&orders)

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

// GetUserOrder 获取用户订单详情
// GET /api/client/orders/:id
func GetUserOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的订单ID",,
			"data": nil})
		return
	}

	db := database.GetDB()
	var order model.Order
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&order).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "订单不存在",,
			"data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    order,
	})
}

// CancelUserOrder 取消用户订单
// POST /api/client/orders/:id/cancel
func CancelUserOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的订单ID",,
			"data": nil})
		return
	}

	db := database.GetDB()
	var order model.Order
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&order).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "订单不存在",,
			"data": nil})
		return
	}

	if order.Status != "pending" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "只有待支付的订单才能取消",,
			"data": nil})
		return
	}

	db.Model(&order).Update("status", "cancelled")

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "取消成功",
	})
}
