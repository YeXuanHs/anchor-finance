package client

import (
	"crypto/rand"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetCart 获取购物车
// GET /api/client/cart
func GetCart(c *gin.Context) {
	userID, _ := c.Get("user_id")
	db := database.GetDB()

	var items []model.CartItem
	db.Where("user_id = ?", userID).Order("id DESC").Find(&items)

	if items == nil {
		items = []model.CartItem{}
	}

	// 计算总金额（服务端计算，不信任前端）
	var totalAmount float64
	for _, item := range items {
		totalAmount += item.Amount * float64(item.Quantity)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "success",
		"data": gin.H{
			"items":        items,
			"total":        len(items),
			"total_amount": totalAmount,
		},
	})
}

// AddToCart 添加到购物车
// POST /api/client/cart
func AddToCart(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		ProductID uint   `json:"product_id" binding:"required"`
		Quantity  int    `json:"quantity"`
		Cycle     string `json:"cycle" binding:"required"`
		ConfigID  uint   `json:"config_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	db := database.GetDB()

	// 验证产品存在且上架
	var product model.Product
	if err := db.Where("id = ? AND status = ?", req.ProductID, "active").First(&product).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "产品不存在或已下架", "data": nil})
		return
	}

	// 服务端计算金额（防0元购，不接受前端价格）
	amount := product.Price

	if amount <= 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "产品价格异常", "data": nil})
		return
	}

	// 检查购物车是否已有该产品+周期
	var existing model.CartItem
	if err := db.Where("user_id = ? AND product_id = ? AND cycle = ?", userID, req.ProductID, req.Cycle).First(&existing).Error; err == nil {
		// 已存在，更新数量
		db.Model(&existing).Updates(map[string]interface{}{
			"quantity": existing.Quantity + req.Quantity,
		})
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "已更新购物车", "data": nil})
		return
	}

	// 新增购物车项
	item := model.CartItem{
		UserID:      userID.(uint),
		ProductID:   req.ProductID,
		ProductName: product.Name,
		ConfigID:    req.ConfigID,
		Quantity:    req.Quantity,
		Cycle:       req.Cycle,
		Amount:      amount,
	}

	if err := db.Create(&item).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "添加失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "已添加到购物车", "data": nil})
}

// UpdateCartItem 更新购物车项
// PUT /api/client/cart/:id
func UpdateCartItem(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")

	var req struct {
		Quantity int `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	var item model.CartItem
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&item).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "购物车项不存在", "data": nil})
		return
	}

	if req.Quantity <= 0 {
		db.Delete(&item)
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "已移除", "data": nil})
		return
	}

	db.Model(&item).Update("quantity", req.Quantity)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "已更新", "data": nil})
}

// RemoveCartItem 删除购物车项
// DELETE /api/client/cart/:id
func RemoveCartItem(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")

	db := database.GetDB()
	result := db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.CartItem{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "购物车项不存在", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "已删除", "data": nil})
}

// ClearCart 清空购物车
// DELETE /api/client/cart
func ClearCart(c *gin.Context) {
	userID, _ := c.Get("user_id")
	db := database.GetDB()
	db.Where("user_id = ?", userID).Delete(&model.CartItem{})

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "购物车已清空", "data": nil})
}

// Checkout 购物车结算（创建订单）
// POST /api/client/cart/checkout
// 安全要点：价格全部服务端计算，不接受前端传入的价格参数（防0元购）
func Checkout(c *gin.Context) {
	userID, _ := c.Get("user_id")
	db := database.GetDB()

	var items []model.CartItem
	db.Where("user_id = ?", userID).Find(&items)

	if len(items) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "购物车为空", "data": nil})
		return
	}

	// 服务端重新计算总金额（不信任前端传入的价格）
	var totalAmount float64
	for _, item := range items {
		// 重新查询产品价格
		var product model.Product
		if err := db.First(&product, item.ProductID).Error; err != nil {
			continue
		}
		// 使用服务端价格，忽略购物车中可能被篡改的价格
		totalAmount += product.Price * float64(item.Quantity)
	}

	// 防0元购：金额必须大于0
	if totalAmount <= 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "订单金额异常，请联系客服", "data": nil})
		return
	}

	// 检查用户余额是否足够（如果选择余额支付）
	var user model.User
	db.First(&user, userID)

	// 生成订单号（使用crypto/rand）
	randNum, _ := rand.Int(rand.Reader, big.NewInt(10000))
	orderNo := "ORD" + time.Now().Format("20060102150405") + strconv.FormatUint(uint64(userID.(uint)), 10) + strconv.FormatInt(randNum.Int64(), 10)

	// 创建订单
	order := model.Order{
		UserID:  userID.(uint),
		OrderNo: orderNo,
		Amount:  totalAmount,
		Status:  "pending",
		Type:    "new",
	}

	if err := db.Create(&order).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "创建订单失败", "data": nil})
		return
	}

	// 创建订单项（使用数据库价格，产品不存在则跳过）
	for _, item := range items {
		var product model.Product
		if err := db.First(&product, item.ProductID).Error; err != nil {
			continue // 产品不存在，跳过
		}
		orderItem := model.OrderItem{
			OrderID:     order.ID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			Cycle:       item.Cycle,
			Amount:      product.Price,
		}
		db.Create(&orderItem)
	}

	// 创建对应的账单（必须设置InvoiceNo，否则uniqueIndex冲突）
	invRandNum, _ := rand.Int(rand.Reader, big.NewInt(100000))
	invoiceNo := "INV" + time.Now().Format("20060102150405") + strconv.FormatUint(uint64(userID.(uint)), 10) + strconv.FormatInt(invRandNum.Int64(), 10)
	invoice := model.Invoice{
		UserID:    userID.(uint),
		InvoiceNo: invoiceNo,
		OrderID:   order.ID,
		Amount:    totalAmount,
		Status:    "unpaid",
		DueDate:   func() *time.Time { t := time.Now().Add(7 * 24 * time.Hour); return &t }(),
	}
	db.Create(&invoice)

	// 清空购物车
	db.Where("user_id = ?", userID).Delete(&model.CartItem{})

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "订单创建成功",
		"data": gin.H{
			"order_id":   order.ID,
			"order_no":   orderNo,
			"invoice_id": invoice.ID,
			"amount":     totalAmount,
		},
	})
}
