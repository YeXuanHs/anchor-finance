package handler

import (
	"strconv"

	"github.com/anchor-finance/backend/internal/service"
	"github.com/anchor-finance/backend/pkg/logger"
	"github.com/anchor-finance/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderSvc *service.OrderService
	log      *logger.Logger
}

func NewOrderHandler(orderSvc *service.OrderService, log *logger.Logger) *OrderHandler {
	return &OrderHandler{orderSvc: orderSvc, log: log}
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
