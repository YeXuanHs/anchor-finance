package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type PaymentGatewayHandler struct {
	svc *service.PaymentGatewayService
	log *logger.Logger
}

func NewPaymentGatewayHandler(svc *service.PaymentGatewayService, log *logger.Logger) *PaymentGatewayHandler {
	return &PaymentGatewayHandler{svc: svc, log: log}
}

// AdminGetList returns paginated payment gateway list.
func (h *PaymentGatewayHandler) AdminGetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status, _ := strconv.Atoi(c.DefaultQuery("status", "-1"))

	items, total, err := h.svc.GetList(page, pageSize, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, items, total, page, pageSize)
}

// AdminGetDetail returns a single payment gateway by ID.
func (h *PaymentGatewayHandler) AdminGetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid gateway id")
		return
	}

	item, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "payment gateway not found")
		return
	}
	response.Success(c, item)
}

// AdminCreate creates a new payment gateway.
func (h *PaymentGatewayHandler) AdminCreate(c *gin.Context) {
	var req service.CreatePaymentGatewayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	item, err := h.svc.Create(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, item)
}

// AdminUpdate updates a payment gateway.
func (h *PaymentGatewayHandler) AdminUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid gateway id")
		return
	}

	var req service.UpdatePaymentGatewayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	item, err := h.svc.Update(uint(id), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, item)
}

// AdminDelete deletes a payment gateway.
func (h *PaymentGatewayHandler) AdminDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid gateway id")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "payment gateway deleted")
}

// AdminToggleStatus toggles a gateway's enabled/disabled status.
func (h *PaymentGatewayHandler) AdminToggleStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid gateway id")
		return
	}

	if err := h.svc.ToggleStatus(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "payment gateway status toggled")
}

// AdminUpdateSort updates sort order for a payment gateway.
func (h *PaymentGatewayHandler) AdminUpdateSort(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid gateway id")
		return
	}

	var req struct {
		SortOrder int `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.UpdateSort(uint(id), req.SortOrder); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "sort order updated")
}

// GetEnabled returns enabled payment gateways for frontend.
func (h *PaymentGatewayHandler) GetEnabled(c *gin.Context) {
	items, err := h.svc.GetEnabled()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, items)
}
