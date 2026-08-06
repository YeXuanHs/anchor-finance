package handler

import (
	"net/http"
	"strconv"

	"anchorfinance/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

// CartHandler handles shopping cart HTTP requests.
type CartHandler struct {
	cartService *service.CartService
}

// NewCartHandler creates a new CartHandler.
func NewCartHandler(cartService *service.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

// addToCartRequest is the payload for AddToCart.
type addToCartRequest struct {
	ProductID    uint           `json:"product_id" binding:"required"`
	BillingCycle string        `json:"billing_cycle" binding:"required"`
	Quantity     int           `json:"quantity"`
	Config       datatypes.JSON `json:"config"`
	Domain       string         `json:"domain"`
}

// updateCartRequest is the payload for UpdateCart.
type updateCartRequest struct {
	Quantity int            `json:"quantity"`
	Config   datatypes.JSON `json:"config"`
}

// checkoutRequest is the payload for Checkout.
type checkoutRequest struct {
	CouponCode string `json:"coupon_code"`
}

// GetCart returns all items in the authenticated user's cart.
// GET /cart
func (h *CartHandler) GetCart(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	items, err := h.cartService.GetCart(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}

// AddToCart adds a product to the user's cart.
// POST /cart
func (h *CartHandler) AddToCart(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	var req addToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	item, err := h.cartService.AddItem(userID, req.ProductID, req.BillingCycle, req.Config, req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

// UpdateCart updates a cart item's quantity or config.
// PUT /cart/:id
func (h *CartHandler) UpdateCart(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cart item id"})
		return
	}

	var req updateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.cartService.UpdateItem(uint(id), userID, req.Quantity, req.Config)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cart item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

// RemoveFromCart removes an item from the cart.
// DELETE /cart/:id
func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cart item id"})
		return
	}

	if err := h.cartService.RemoveItem(uint(id), userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "removed"})
}

// ClearCart clears all items from the user's cart.
// DELETE /cart/clear
func (h *CartHandler) ClearCart(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	if err := h.cartService.ClearCart(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to clear cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cart cleared"})
}

// Checkout creates an order from the cart items.
// POST /cart/checkout
func (h *CartHandler) Checkout(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	var req checkoutRequest
	// Coupon code is optional, so we don't require binding errors
	_ = c.ShouldBindJSON(&req)

	order, err := h.cartService.Checkout(userID, req.CouponCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

// BatchDelete removes multiple items from the cart.
// POST /cart/batch-delete
func (h *CartHandler) BatchDelete(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	var req struct {
		IDs []string `json:"ids" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, idStr := range req.IDs {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			continue
		}
		h.cartService.RemoveItem(uint(id), userID)
	}

	c.JSON(http.StatusOK, gin.H{"message": "批量删除成功"})
}

// getUserID extracts the authenticated user ID from the gin context.
// This assumes auth middleware sets "user_id" in the context.
func getUserID(c *gin.Context) uint {
	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return 0
	}
	uid, ok := val.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user context"})
		return 0
	}
	return uid
}
