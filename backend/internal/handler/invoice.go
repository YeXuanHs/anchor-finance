package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type InvoiceHandler struct {
	invoiceSvc *service.InvoiceService
	log        *logger.Logger
}

func NewInvoiceHandler(invoiceSvc *service.InvoiceService, log *logger.Logger) *InvoiceHandler {
	return &InvoiceHandler{invoiceSvc: invoiceSvc, log: log}
}

// GetDetail returns a single invoice.
func (h *InvoiceHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid invoice id")
		return
	}

	invoice, err := h.invoiceSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "invoice not found")
		return
	}
	response.Success(c, invoice)
}

// GetUserInvoices returns paginated invoices for the authenticated user.
func (h *InvoiceHandler) GetUserInvoices(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	invoices, total, err := h.invoiceSvc.GetUserInvoices(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, invoices, total, page, pageSize)
}

// Pay pays an invoice using the user's balance.
func (h *InvoiceHandler) Pay(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid invoice id")
		return
	}

	if err := h.invoiceSvc.PayWithBalance(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "invoice paid")
}

// Cancel cancels a pending invoice.
func (h *InvoiceHandler) Cancel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid invoice id")
		return
	}

	if err := h.invoiceSvc.Cancel(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "invoice cancelled")
}

// GetList returns all invoices (admin).
func (h *InvoiceHandler) GetList(c *gin.Context) {
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

	invoices, total, err := h.invoiceSvc.GetList(page, pageSize, status, userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, invoices, total, page, pageSize)
}
