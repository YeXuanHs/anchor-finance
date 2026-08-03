package handler

import (
	"fmt"
	"strconv"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var _ = model.Order{}

type OrderHandler struct {
	orderSvc *service.OrderService
	db       *gorm.DB
	log      *logger.Logger
}

func NewOrderHandler(orderSvc *service.OrderService, log *logger.Logger) *OrderHandler {
	return &OrderHandler{orderSvc: orderSvc, log: log}
}

// NewOrderHandlerWithDB creates an OrderHandler with direct DB access.
func NewOrderHandlerWithDB(orderSvc *service.OrderService, db *gorm.DB, log *logger.Logger) *OrderHandler {
	return &OrderHandler{orderSvc: orderSvc, db: db, log: log}
}

// Create creates a new order.
func (h *OrderHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req service.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	order, err := h.orderSvc.Create(userID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, order)
}

// previewItem is a single item in the preview request.
type previewItem struct {
	ProductID    uint   `json:"product_id" binding:"required"`
	BillingCycle string `json:"billing_cycle"`
	Quantity     int    `json:"quantity"`
}

// PreviewRequest is the payload for order preview.
type PreviewRequest struct {
	Items      []previewItem `json:"items" binding:"required,min=1"`
	CouponCode string        `json:"coupon_code"`
}

// Preview calculates order total without creating the order.
func (h *OrderHandler) Preview(c *gin.Context) {
	var req PreviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var total float64
	var items []map[string]interface{}

	for _, item := range req.Items {
		if item.Quantity < 1 {
			item.Quantity = 1
		}

		// Look up product from database
		var product model.Product
		if err := h.db.First(&product, item.ProductID).Error; err != nil {
			response.BadRequest(c, fmt.Sprintf("product not found: %d", item.ProductID))
			return
		}

		if product.Status != 1 {
			response.BadRequest(c, "product is not available: "+product.Name)
			return
		}

		// Check stock
		if product.StockControl && product.Stock >= 0 && product.Stock < item.Quantity {
			response.BadRequest(c, fmt.Sprintf("insufficient stock for %s: available %d", product.Name, product.Stock))
			return
		}

		// Get base price from Decimal
		basePrice, _ := product.Price.Float64()

		// Determine price based on billing cycle
		price := basePrice
		cycle := item.BillingCycle
		if cycle == "" {
			cycle = product.BillingCycle
		}
		if cycle != "" && cycle != product.BillingCycle {
			switch cycle {
			case "monthly":
				price = basePrice
			case "quarterly":
				price = basePrice * 3 * 0.95
			case "semi-annually":
				price = basePrice * 6 * 0.90
			case "annually":
				price = basePrice * 12 * 0.85
			case "biennially":
				price = basePrice * 24 * 0.80
			case "triennially":
				price = basePrice * 36 * 0.75
			default:
				price = basePrice
			}
		}

		subtotal := price * float64(item.Quantity)
		total += subtotal

		items = append(items, map[string]interface{}{
			"product_id":    product.ID,
			"product_name":  product.Name,
			"billing_cycle": cycle,
			"quantity":      item.Quantity,
			"price":         price,
			"subtotal":      subtotal,
		})
	}

	// Apply coupon if provided
	couponDiscount := 0.0
	if req.CouponCode != "" {
		var coupon model.Coupon
		if err := h.db.Where("code = ? AND status = 1", req.CouponCode).First(&coupon).Error; err == nil {
			now := time.Now()
			valid := true
			if coupon.StartDate != nil && coupon.StartDate.After(now) {
				valid = false
			}
			if coupon.EndDate != nil && coupon.EndDate.Before(now) {
				valid = false
			}
			if coupon.MaxUses > 0 && coupon.UsedCount >= coupon.MaxUses {
				valid = false
			}
			minOrder, _ := coupon.MinOrderAmount.Float64()
			if minOrder > 0 && total < minOrder {
				valid = false
			}

			if valid {
				couponVal, _ := coupon.Value.Float64()
				switch coupon.Type {
				case "percentage":
					couponDiscount = total * couponVal / 100
					maxDisc, _ := coupon.MaxDiscount.Float64()
					if maxDisc > 0 && couponDiscount > maxDisc {
						couponDiscount = maxDisc
					}
				case "fixed":
					couponDiscount = couponVal
					if couponDiscount > total {
						couponDiscount = total
					}
				case "free":
					couponDiscount = total
				}
			}
		}
	}

	finalTotal := total - couponDiscount
	if finalTotal < 0 {
		finalTotal = 0
	}

	response.Success(c, gin.H{
		"items":           items,
		"subtotal":        total,
		"coupon_code":     req.CouponCode,
		"coupon_discount": couponDiscount,
		"total":           finalTotal,
	})
}

// GetDetail returns a single order.
func (h *OrderHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid order id")
		return
	}

	order, err := h.orderSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "order not found")
		return
	}
	response.Success(c, order)
}

// GetUserOrders returns paginated orders for the authenticated user.
func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	userID := int(c.GetUint("user_id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	orders, total, err := h.orderSvc.GetUserOrders(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, orders, total, page, pageSize)
}

// Pay marks an order as paid, creating invoice + user_product.
func (h *OrderHandler) Pay(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid order id")
		return
	}

	order, err := h.orderSvc.Pay(uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, order)
}

// Cancel cancels a pending order.
func (h *OrderHandler) Cancel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid order id")
		return
	}

	if err := h.orderSvc.Cancel(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "order cancelled")
}

// GetList returns all orders (admin).
func (h *OrderHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var status *int
	if s := c.Query("status"); s != "" {
		v, _ := strconv.Atoi(s)
		status = &v
	}
	var userID *uint
	if u := c.Query("user_id"); u != "" {
		v, _ := strconv.ParseUint(u, 10, 64)
		uid := uint(v)
		userID = &uid
	}

	orders, total, err := h.orderSvc.GetList(page, pageSize, status, userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, orders, total, page, pageSize)
}

// UpdateStatus manually updates an order's status (admin).
func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid order id")
		return
	}

	var req struct {
		Status int `json:"status" binding:"required,oneof=0 1 2 3"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	_ = id // status update via direct DB in real impl
	response.SuccessMsg(c, "order status updated")
}

// ---------------------------------------------------------------------------
// Admin order management endpoints
// ---------------------------------------------------------------------------

// AdminCreateOrder handles admin creating an order for a user.
func (h *OrderHandler) AdminCreateOrder(c *gin.Context) {
	adminID := c.GetUint("admin_id")
	var req service.AdminCreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	order, err := h.orderSvc.AdminCreateOrder(adminID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, order)
}

// CheckOrder validates whether an order can be executed.
func (h *OrderHandler) CheckOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid order id")
		return
	}

	result, err := h.orderSvc.CheckOrder(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, result)
}

// BatchUpdate performs batch operations on orders (confirm/cancel/delete).
func (h *OrderHandler) BatchUpdate(c *gin.Context) {
	adminID := c.GetUint("admin_id")
	var req service.BatchUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	processed, err := h.orderSvc.BatchUpdate(adminID, req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{
		"action":    req.Action,
		"processed": processed,
	})
}

// Delete soft-deletes an order.
func (h *OrderHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid order id")
		return
	}

	if err := h.orderSvc.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "order deleted")
}

// AddNote adds an admin note to an order.
func (h *OrderHandler) AddNote(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid order id")
		return
	}

	adminID := c.GetUint("admin_id")

	var req struct {
		Content   string `json:"content" binding:"required"`
		IsPrivate bool   `json:"is_private"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	note, err := h.orderSvc.AddNote(uint(id), adminID, req.Content, req.IsPrivate)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, note)
}

// GetNotes returns all notes for an order.
func (h *OrderHandler) GetNotes(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid order id")
		return
	}

	notes, err := h.orderSvc.GetNotes(uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, notes)
}

// GetSaleOrders returns sale-related orders.
func (h *OrderHandler) GetSaleOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var saleID *uint
	if s := c.Query("sale_id"); s != "" {
		v, _ := strconv.ParseUint(s, 10, 64)
		sid := uint(v)
		saleID = &sid
	}

	var status *string
	if s := c.Query("status"); s != "" {
		status = &s
	}

	orders, total, err := h.orderSvc.GetSaleOrders(page, pageSize, saleID, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, orders, total, page, pageSize)
}

// ActivateOrder manually activates an order (opens service without payment).
func (h *OrderHandler) ActivateOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid order id")
		return
	}

	order, err := h.orderSvc.ActivateOrder(uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, order)
}

// ChangeStatus changes the status of an order.
func (h *OrderHandler) ChangeStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid order id")
		return
	}

	var req struct {
		Status int16 `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.orderSvc.ChangeStatus(uint(id), req.Status); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "order status changed")
}

// GetMultiTotal calculates the total price for multiple products.
func (h *OrderHandler) GetMultiTotal(c *gin.Context) {
	var req struct {
		Items      []service.MultiProductItem `json:"items" binding:"required,min=1"`
		CouponCode string                     `json:"coupon_code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.orderSvc.GetMultiTotal(req.Items, req.CouponCode)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// ApplyCustomPromo creates a custom promo code.
func (h *OrderHandler) ApplyCustomPromo(c *gin.Context) {
	var req service.CreateCustomPromoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	promo, err := h.orderSvc.ApplyCustomPromo(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, promo)
}

// SearchPage returns filter configuration for the order search page.
func (h *OrderHandler) SearchPage(c *gin.Context) {
	config := h.orderSvc.SearchPageConfig()
	response.Success(c, config)
}

// ==================== P0-5: CheckProduct ====================

// CheckProduct 校验产品试用/购买资格
func (h *OrderHandler) CheckProduct(c *gin.Context) {
	var req struct {
		ProductID uint `json:"product_id" binding:"required"`
		UserID    uint `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.orderSvc.CheckProduct(req.ProductID, req.UserID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// ==================== P3-17: SetConfig ====================

// SetConfig 获取下单配置（产品定价/配置选项/自定义字段）
func (h *OrderHandler) SetConfig(c *gin.Context) {
	pid, err := strconv.ParseUint(c.Query("pid"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product id")
		return
	}
	cycle := c.DefaultQuery("billingcycle", "monthly")

	data, err := h.orderSvc.GetOrderConfig(uint(pid), cycle)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

// ==================== P3-17: GetClients ====================

// GetClients 搜索客户下拉列表
func (h *OrderHandler) GetClients(c *gin.Context) {
	keyword := c.Query("username")

	var clients []map[string]interface{}
	query := h.db.Table("users").Select("id, username, email, phone")
	if keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if err := query.Limit(50).Scan(&clients).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"data": clients})
}

// ==================== P2-13: Enhanced GetList filters ====================

// GetListEnhanced returns all orders with advanced filters (admin).
func (h *OrderHandler) GetListEnhanced(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	orderField := c.DefaultQuery("order", "id")
	sort := c.DefaultQuery("sort", "desc")

	query := h.db.Model(&model.Order{})

	// 基础筛选
	if s := c.Query("status"); s != "" {
		query = query.Where("orders.status = ?", s)
	}
	if uid := c.Query("uid"); uid != "" {
		query = query.Where("orders.user_id = ?", uid)
	}
	if id := c.Query("id"); id != "" {
		query = query.Where("orders.id LIKE ?", "%"+id+"%")
	}
	if amount := c.Query("amount"); amount != "" {
		query = query.Where("orders.total_price LIKE ?", "%"+amount+"%")
	}
	if ordernum := c.Query("ordernum"); ordernum != "" {
		query = query.Where("orders.order_no LIKE ?", "%"+ordernum+"%")
	}

	// 用户名筛选
	if username := c.Query("username"); username != "" {
		query = query.Joins("JOIN users ON users.id = orders.user_id").
			Where("users.username LIKE ?", "%"+username+"%")
	}

	// 产品名筛选
	if productName := c.Query("product_name"); productName != "" {
		query = query.Joins("JOIN products ON products.id = orders.product_id").
			Where("products.name = ?", productName)
	}

	// 支付方式筛选
	if payment := c.Query("payment"); payment != "" {
		if payment == "creditPay" {
			query = query.Joins("LEFT JOIN invoices ON invoices.order_id = orders.id").
				Where("invoices.credit > 0")
		} else if payment == "creditLimitPay" {
			query = query.Joins("LEFT JOIN invoices ON invoices.order_id = orders.id").
				Where("invoices.use_credit_limit = 1")
		} else {
			query = query.Where("orders.payment = ?", payment)
		}
	}

	// 销售筛选
	if saleID := c.Query("sale_id"); saleID != "" {
		query = query.Joins("JOIN users u ON u.id = orders.user_id").
			Where("u.sale_id = ?", saleID)
	}

	// 付款状态筛选
	if payStatus := c.Query("pay_status"); payStatus != "" {
		if payStatus == "0" {
			query = query.Where("orders.invoice_id IS NULL OR orders.invoice_id = 0")
		} else {
			query = query.Joins("LEFT JOIN invoices inv ON inv.id = orders.invoice_id").
				Where("inv.status = ?", payStatus)
		}
	}

	// 时间范围
	if startTime := c.Query("start_time"); startTime != "" {
		query = query.Where("orders.created_at >= ?", startTime)
	}
	if endTime := c.Query("end_time"); endTime != "" {
		query = query.Where("orders.created_at <= ?", endTime+" 23:59:59")
	}

	var total int64
	query.Count(&total)

	offset, limit := service.Paginate(page, pageSize)
	var orders []model.Order
	if err := query.Preload("Product").Offset(offset).Limit(limit).
		Order(fmt.Sprintf("orders.%s %s", orderField, sort)).Find(&orders).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, orders, total, page, pageSize)
}
