package admin

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetOrderList 获取订单列表
// GET /api/admin/orders
func GetOrderList(c *gin.Context) {
	// 1. 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	status := c.Query("status")
	userID := c.Query("user_id")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 2. 构建查询
	db := database.GetDB()
	query := db.Model(&model.Order{})

	// 关键词搜索（订单号、产品名称）
	if keyword != "" {
		query = query.Where("order_no LIKE ? OR product_name LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 用户筛选
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// 3. 获取总数
	var total int64
	query.Count(&total)

	// 4. 分页查询
	var orders []model.Order
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&orders)

	// 5. 返回统一格式
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

// GetOrder 获取订单详情
// GET /api/admin/orders/:id
func GetOrder(c *gin.Context) {
	// 1. 获取订单ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的订单ID",
			"data":    nil,
		})
		return
	}

	// 2. 查询订单
	db := database.GetDB()
	var order model.Order
	if err := db.First(&order, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "订单不存在",
			"data":    nil,
		})
		return
	}

	// 3. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    order,
	})
}

// CreateOrder 创建订单
// POST /api/admin/orders
func CreateOrder(c *gin.Context) {
	// 1. 解析请求参数
	var req struct {
		UserID      uint    `json:"user_id" binding:"required"`
		ProductID   uint    `json:"product_id" binding:"required"`
		ProductName string  `json:"product_name" binding:"required"`
		Quantity    int     `json:"quantity"`
		Amount      float64 `json:"amount" binding:"required"`
		Remark      string  `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	// 0元购防护：订单金额必须大于0
	if req.Amount <= 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "订单金额必须大于0", "data": nil})
		return
	}

	// 2. 验证用户是否存在
	db := database.GetDB()
	var user model.User
	if err := db.First(&user, req.UserID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "用户不存在",
			"data":    nil,
		})
		return
	}

	// 3. 生成订单号
	orderNo := fmt.Sprintf("ORD%s%06d", time.Now().Format("20060102"), time.Now().UnixNano()%1000000)

	// 4. 创建订单
	if req.Quantity == 0 {
		req.Quantity = 1
	}

	order := model.Order{
		UserID:      req.UserID,
		OrderNo:     orderNo,
		ProductID:   req.ProductID,
		ProductName: req.ProductName,
		Quantity:    req.Quantity,
		Amount:      req.Amount,
		Status:      "pending",
		Remark:      req.Remark,
	}

	if err := db.Create(&order).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建订单失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 5. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id":      order.ID,
			"order_no": order.OrderNo,
		},
	})
}

// UpdateOrder 更新订单
// PUT /api/admin/orders/:id
func UpdateOrder(c *gin.Context) {
	// 1. 获取订单ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的订单ID",
			"data":    nil,
		})
		return
	}

	// 2. 解析请求参数
	var req struct {
		ProductName string  `json:"product_name"`
		Quantity    int     `json:"quantity"`
		Amount      float64 `json:"amount"`
		Status      string  `json:"status"`
		Remark      string  `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 3. 查询订单
	db := database.GetDB()
	var order model.Order
	if err := db.First(&order, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "订单不存在",
			"data":    nil,
		})
		return
	}

	// 4. 更新订单
	updates := map[string]interface{}{}
	if req.ProductName != "" {
		updates["product_name"] = req.ProductName
	}
	if req.Quantity > 0 {
		updates["quantity"] = req.Quantity
	}
	if req.Amount > 0 {
		updates["amount"] = req.Amount
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}

	if err := db.Model(&order).Updates(updates).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "更新订单失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 5. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
		"data":    nil,
	})
}

// ActivateOrder 激活订单
// POST /api/admin/orders/:id/activate
func ActivateOrder(c *gin.Context) {
	// 1. 获取订单ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的订单ID",
			"data":    nil,
		})
		return
	}

	// 2. 查询订单
	db := database.GetDB()
	var order model.Order
	if err := db.First(&order, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "订单不存在",
			"data":    nil,
		})
		return
	}

	// 3. 检查订单状态
	if order.Status != "pending" && order.Status != "paid" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "订单状态不允许激活",
			"data":    nil,
		})
		return
	}

	// 4. 激活订单（乐观锁防重复激活）
	now := time.Now()
	result := db.Model(&model.Order{}).
		Where("id = ? AND status IN ?", id, []string{"pending", "paid"}).
		Updates(map[string]interface{}{
			"status":  "active",
			"paid_at": &now,
		})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "订单状态已变更，无法激活", "data": nil})
		return
	}

	// 5. 返回统一格式
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "激活成功", "data": nil})
}

// CancelOrder 取消订单
// POST /api/admin/orders/:id/cancel
func CancelOrder(c *gin.Context) {
	// 1. 获取订单ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的订单ID",
			"data":    nil,
		})
		return
	}

	// 2. 查询订单
	db := database.GetDB()
	var order model.Order
	if err := db.First(&order, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "订单不存在",
			"data":    nil,
		})
		return
	}

	// 3. 检查订单状态
	if order.Status == "cancelled" || order.Status == "refunded" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "订单已取消或退款",
			"data":    nil,
		})
		return
	}

	// 4. 取消订单
	if err := db.Model(&order).Update("status", "cancelled").Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "取消订单失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 5. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "取消成功",
		"data":    nil,
	})
}
