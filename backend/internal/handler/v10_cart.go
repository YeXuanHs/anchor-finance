package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// V10CartHandler handles V10 shopping cart HTTP requests.
type V10CartHandler struct {
	cartSvc *service.V10CartService
	log     *logger.Logger
}

// NewV10CartHandler creates a new V10CartHandler.
func NewV10CartHandler(cartSvc *service.V10CartService, log *logger.Logger) *V10CartHandler {
	return &V10CartHandler{cartSvc: cartSvc, log: log}
}

// GetCart returns all items in the user's cart with summary.
// GET /v10/cart
func (h *V10CartHandler) GetCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	summary, err := h.cartSvc.GetCart(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, summary)
}

// AddItem adds a product to the cart.
// POST /v10/cart
func (h *V10CartHandler) AddItem(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req service.V10AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	item, err := h.cartSvc.AddItem(userID, req.ProductID, req.BillingCycle, req.Config, req.Quantity, req.Domain)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, item)
}

// UpdateItem updates a cart item.
// PUT /v10/cart/:id
func (h *V10CartHandler) UpdateItem(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID错误")
		return
	}

	var req service.V10UpdateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	item, err := h.cartSvc.UpdateItem(userID, uint(id), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, item)
}

// RemoveItem removes an item from the cart.
// DELETE /v10/cart/:id
func (h *V10CartHandler) RemoveItem(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID错误")
		return
	}

	if err := h.cartSvc.RemoveItem(userID, uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "删除成功")
}

// ClearCart clears all items from the cart.
// DELETE /v10/cart/clear
func (h *V10CartHandler) ClearCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	if err := h.cartSvc.ClearCart(userID); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "请求成功")
}

// ApplyCoupon applies a coupon code to the cart.
// POST /v10/cart/coupon
func (h *V10CartHandler) ApplyCoupon(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req service.V10ApplyCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.cartSvc.ApplyCoupon(userID, req.CouponCode); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "优惠码应用成功")
}

// RemoveCoupon removes the coupon from the cart.
// DELETE /v10/cart/coupon
func (h *V10CartHandler) RemoveCoupon(c *gin.Context) {
	userID := c.GetUint("user_id")

	if err := h.cartSvc.RemoveCoupon(userID); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "优惠码移除成功")
}

// Checkout creates orders from cart items.
// POST /v10/cart/checkout
func (h *V10CartHandler) Checkout(c *gin.Context) {
	userID := c.GetUint("user_id")

	orders, err := h.cartSvc.Checkout(userID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"orders": orders})
}

// GetItemCount returns the number of items in the cart.
// GET /v10/cart/count
func (h *V10CartHandler) GetItemCount(c *gin.Context) {
	userID := c.GetUint("user_id")

	count, err := h.cartSvc.GetItemCount(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"count": count})
}
