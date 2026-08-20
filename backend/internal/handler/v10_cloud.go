package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// V10CloudHandler handles V10 cloud server HTTP requests.
type V10CloudHandler struct {
	cloudSvc *service.V10CloudService
	log      *logger.Logger
}

// NewV10CloudHandler creates a new V10CloudHandler.
func NewV10CloudHandler(cloudSvc *service.V10CloudService, log *logger.Logger) *V10CloudHandler {
	return &V10CloudHandler{cloudSvc: cloudSvc, log: log}
}

// ═══════════════════ Product Browsing ═══════════════════

// GetProductList lists cloud products by group.
// GET /v10/cloud/products
func (h *V10CloudHandler) GetProductList(c *gin.Context) {
	var req service.CloudProductListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 20
	}

	products, total, err := h.cloudSvc.GetProductList(req.GroupID, req.Page, req.PageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, products, total, req.Page, req.PageSize)
}

// GetProductDetail returns a single cloud product.
// GET /v10/cloud/products/:id
func (h *V10CloudHandler) GetProductDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product id")
		return
	}

	product, err := h.cloudSvc.GetProductDetail(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, product)
}

// GetRegions lists available regions.
// GET /v10/cloud/regions
func (h *V10CloudHandler) GetRegions(c *gin.Context) {
	regions, err := h.cloudSvc.GetRegions()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, regions)
}

// GetOSTypes lists available OS types.
// GET /v10/cloud/os-types
func (h *V10CloudHandler) GetOSTypes(c *gin.Context) {
	osTypes, err := h.cloudSvc.GetOSTypes()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, osTypes)
}

// ═══════════════════ Configuration ═══════════════════

// GetConfigOptions returns configurable options for a product.
// GET /v10/cloud/products/:id/config
func (h *V10CloudHandler) GetConfigOptions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product id")
		return
	}

	opts, err := h.cloudSvc.GetConfigOptions(uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, opts)
}

// CalculatePrice calculates price based on configuration.
// POST /v10/cloud/calculate-price
func (h *V10CloudHandler) CalculatePrice(c *gin.Context) {
	var req service.CalculatePriceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.Quantity < 1 {
		req.Quantity = 1
	}

	breakdown, err := h.cloudSvc.CalculatePrice(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, breakdown)
}

// GetLinkAgeList returns cascading config options.
// GET /v10/cloud/products/:id/linkage
func (h *V10CloudHandler) GetLinkAgeList(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product id")
		return
	}

	var parentID *uint
	if pid := c.Query("parent_id"); pid != "" {
		p, err := strconv.ParseUint(pid, 10, 64)
		if err == nil {
			uid := uint(p)
			parentID = &uid
		}
	}

	opts, err := h.cloudSvc.GetLinkAgeList(uint(id), parentID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, opts)
}

// FilterConfigOptions filters available options.
// GET /v10/cloud/products/:id/config/filter
func (h *V10CloudHandler) FilterConfigOptions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product id")
		return
	}

	var filters service.ConfigFilter
	if err := c.ShouldBindQuery(&filters); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	opts, err := h.cloudSvc.FilterConfigOptions(uint(id), filters)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, opts)
}

// ═══════════════════ Cart Operations ═══════════════════

// AddToCart adds a configured cloud product to cart.
// POST /v10/cloud/cart
func (h *V10CloudHandler) AddToCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req service.CartItemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.Quantity < 1 {
		req.Quantity = 1
	}

	item, err := h.cloudSvc.AddToCart(userID, req.ProductID, req.Config, req.Cycle, req.Quantity)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, item)
}

// UpdateCartItem updates a cart item.
// PUT /v10/cloud/cart/:id
func (h *V10CloudHandler) UpdateCartItem(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cart item id")
		return
	}

	var req service.UpdateCartItemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	item, err := h.cloudSvc.UpdateCartItem(userID, uint(id), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, item)
}

// GetCartSummary returns cart summary with totals.
// GET /v10/cloud/cart
func (h *V10CloudHandler) GetCartSummary(c *gin.Context) {
	userID := c.GetUint("user_id")

	summary, err := h.cloudSvc.GetCartSummary(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, summary)
}

// GetCartItems lists all cart items.
// GET /v10/cloud/cart/items
func (h *V10CloudHandler) GetCartItems(c *gin.Context) {
	userID := c.GetUint("user_id")

	items, err := h.cloudSvc.GetCartItems(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, items)
}

// SettleCart validates and prepares cart for checkout.
// POST /v10/cloud/cart/settle
func (h *V10CloudHandler) SettleCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	items, err := h.cloudSvc.SettleCart(userID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"items": items, "item_count": len(items)})
}

// ═══════════════════ Order Flow ═══════════════════

// CreateOrder creates orders from cart items.
// POST /v10/cloud/orders
func (h *V10CloudHandler) CreateOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req service.CreateOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	orders, err := h.cloudSvc.CreateOrder(userID, req.CartItemIDs, req.CouponCode)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"orders": orders})
}

// GetOrderDetail returns order details.
// GET /v10/cloud/orders/:id
func (h *V10CloudHandler) GetOrderDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid order id")
		return
	}

	order, err := h.cloudSvc.GetOrderDetail(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, order)
}

// PayOrder processes payment for an order.
// POST /v10/cloud/orders/:id/pay
func (h *V10CloudHandler) PayOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid order id")
		return
	}

	var req service.PayOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	order, err := h.cloudSvc.PayOrder(uint(id), req.PaymentMethod)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, order)
}

// ═══════════════════ Host Management ═══════════════════

// GetHostInfo returns cloud host information.
// GET /v10/cloud/hosts/:id
func (h *V10CloudHandler) GetHostInfo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	info, err := h.cloudSvc.GetHostInfo(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, info)
}

// GetHostConfig returns host configuration.
// GET /v10/cloud/hosts/:id/config
func (h *V10CloudHandler) GetHostConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	config, err := h.cloudSvc.GetHostConfig(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, gin.H{"config": config})
}

// GetTrafficUsage returns traffic usage stats.
// GET /v10/cloud/hosts/:id/traffic
func (h *V10CloudHandler) GetTrafficUsage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	usage, err := h.cloudSvc.GetTrafficUsage(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, usage)
}

// GetOSList returns available OS for reinstall.
// GET /v10/cloud/hosts/:id/os
func (h *V10CloudHandler) GetOSList(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	osList, err := h.cloudSvc.GetOSList(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, osList)
}
