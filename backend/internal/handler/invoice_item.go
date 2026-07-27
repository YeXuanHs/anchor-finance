package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type InvoiceItemHandler struct {
	itemSvc *service.InvoiceItemService
	log     *logger.Logger
}

func NewInvoiceItemHandler(itemSvc *service.InvoiceItemService, log *logger.Logger) *InvoiceItemHandler {
	return &InvoiceItemHandler{itemSvc: itemSvc, log: log}
}

// Create creates a new invoice item (admin).
func (h *InvoiceItemHandler) Create(c *gin.Context) {
	var req service.CreateInvoiceItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	item, err := h.itemSvc.Create(req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, item)
}

// GetDetail returns a single invoice item.
func (h *InvoiceItemHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid item id")
		return
	}

	item, err := h.itemSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "item not found")
		return
	}
	response.Success(c, item)
}

// GetByInvoiceID returns all items for an invoice.
func (h *InvoiceItemHandler) GetByInvoiceID(c *gin.Context) {
	invoiceID, err := strconv.ParseUint(c.Param("invoice_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid invoice id")
		return
	}

	items, err := h.itemSvc.GetByInvoiceID(uint(invoiceID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, items)
}

// Update modifies an invoice item (admin).
func (h *InvoiceItemHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid item id")
		return
	}

	var req service.UpdateInvoiceItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	item, err := h.itemSvc.Update(uint(id), req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, item)
}

// Delete soft-deletes an invoice item (admin).
func (h *InvoiceItemHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid item id")
		return
	}

	if err := h.itemSvc.Delete(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "item deleted")
}

// BatchCreate creates multiple invoice items (admin).
func (h *InvoiceItemHandler) BatchCreate(c *gin.Context) {
	var req service.BatchCreateInvoiceItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	items, err := h.itemSvc.BatchCreate(req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, items)
}

// BatchDelete deletes multiple invoice items (admin).
func (h *InvoiceItemHandler) BatchDelete(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.itemSvc.BatchDelete(req.IDs); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "items deleted")
}

// BatchUpdate updates multiple invoice items (admin).
func (h *InvoiceItemHandler) BatchUpdate(c *gin.Context) {
	var req struct {
		IDs  []uint                       `json:"ids" binding:"required,min=1"`
		Data service.UpdateInvoiceItemRequest `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.itemSvc.BatchUpdate(req.IDs, req.Data); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "items updated")
}

// CalculateInvoiceTotal calculates the total for an invoice.
func (h *InvoiceItemHandler) CalculateInvoiceTotal(c *gin.Context) {
	invoiceID, err := strconv.ParseUint(c.Param("invoice_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid invoice id")
		return
	}

	totals, err := h.itemSvc.CalculateInvoiceTotal(uint(invoiceID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, totals)
}

// AddDiscount adds a discount to an invoice item (admin).
func (h *InvoiceItemHandler) AddDiscount(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid item id")
		return
	}

	var req struct {
		DiscountType string  `json:"discount_type" binding:"required,oneof=amount percent"`
		Value        float64 `json:"value" binding:"required,gt=0"`
		Description  string  `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.itemSvc.AddDiscount(uint(id), req.DiscountType, req.Value, req.Description); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "discount added")
}

// AddTax adds a tax to an invoice item (admin).
func (h *InvoiceItemHandler) AddTax(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid item id")
		return
	}

	var req struct {
		TaxName     string  `json:"tax_name" binding:"required,max=64"`
		TaxRate     float64 `json:"tax_rate" binding:"required,gt=0,lte=1"`
		Description string  `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.itemSvc.AddTax(uint(id), req.TaxName, req.TaxRate, req.Description); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "tax added")
}
