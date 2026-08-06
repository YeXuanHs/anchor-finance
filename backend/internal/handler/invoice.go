package handler

import (
	"fmt"
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type InvoiceHandler struct {
	invoiceSvc  *service.InvoiceService
	enhancedSvc *service.InvoiceEnhancedService
	log         *logger.Logger
}

func NewInvoiceHandler(invoiceSvc *service.InvoiceService, log *logger.Logger) *InvoiceHandler {
	return &InvoiceHandler{invoiceSvc: invoiceSvc, log: log}
}

// SetEnhancedService sets the enhanced invoice service (called from router after init).
func (h *InvoiceHandler) SetEnhancedService(svc *service.InvoiceEnhancedService) {
	h.enhancedSvc = svc
}

// GetDetail returns a single invoice.
func (h *InvoiceHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID错误")
		return
	}

	invoice, err := h.invoiceSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "ID错误")
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
		response.BadRequest(c, "ID错误")
		return
	}

	if err := h.invoiceSvc.PayWithBalance(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "请求成功")
}

// Cancel cancels a pending invoice.
func (h *InvoiceHandler) Cancel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID错误")
		return
	}

	if err := h.invoiceSvc.Cancel(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "取消成功")
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

// CombineInvoices merges multiple unpaid invoices into one (admin).
func (h *InvoiceHandler) CombineInvoices(c *gin.Context) {
	if h.enhancedSvc == nil {
		response.ServerError(c, "enhanced service not available")
		return
	}

	var req struct {
		InvoiceIDs []uint `json:"invoice_ids" binding:"required,min=2"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	combined, err := h.enhancedSvc.CombineInvoices(req.InvoiceIDs)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, combined)
}

// GetCombineInvoices returns invoices eligible for combining for a user (admin).
func (h *InvoiceHandler) GetCombineInvoices(c *gin.Context) {
	if h.enhancedSvc == nil {
		response.ServerError(c, "enhanced service not available")
		return
	}

	userID, err := strconv.ParseUint(c.Query("user_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID错误")
		return
	}

	invoices, err := h.enhancedSvc.GetCombineInvoices(uint(userID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, invoices)
}

// InvoiceLog returns the operation log for an invoice (admin).
func (h *InvoiceHandler) InvoiceLog(c *gin.Context) {
	if h.enhancedSvc == nil {
		response.ServerError(c, "enhanced service not available")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID错误")
		return
	}

	logs, err := h.enhancedSvc.GetInvoiceLog(uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, logs)
}

// Duplicate duplicates an invoice (admin).
func (h *InvoiceHandler) Duplicate(c *gin.Context) {
	if h.enhancedSvc == nil {
		response.ServerError(c, "enhanced service not available")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID错误")
		return
	}

	newInvoice, err := h.enhancedSvc.DuplicateInvoice(uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, newInvoice)
}

// Option updates invoice options like dates, status, payment method, notes (admin).
func (h *InvoiceHandler) Option(c *gin.Context) {
	if h.enhancedSvc == nil {
		response.ServerError(c, "enhanced service not available")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID错误")
		return
	}

	var req struct {
		DueDate   *string `json:"due_date"`
		Status    *int    `json:"status"`
		Payment   *string `json:"payment"`
		Notes     *string `json:"notes"`
		InvoiceNo *string `json:"invoice_no"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	_, err := h.invoiceSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "ID错误")
		return
	}

	adminID := c.GetUint("admin_id")
	updates := map[string]interface{}{}
	detail := ""

	if req.DueDate != nil {
		updates["due_date"] = *req.DueDate
		detail += "due_date updated; "
	}
	if req.Status != nil {
		updates["status"] = *req.Status
		detail += "status updated; "
	}
	if req.Payment != nil {
		updates["payment"] = *req.Payment
		detail += "payment updated; "
	}
	if req.Notes != nil {
		updates["notes"] = *req.Notes
		detail += "notes updated; "
	}
	if req.InvoiceNo != nil {
		updates["invoice_no"] = *req.InvoiceNo
		detail += "invoice_no updated; "
	}

	if len(updates) > 0 {
		if err := h.invoiceSvc.UpdateInvoice(uint(id), updates); err != nil {
			response.BadRequest(c, err.Error())
			return
		}
		h.enhancedSvc.LogAction(uint(id), adminID, "option_updated", detail, c.ClientIP())
	}

	inv, _ := h.invoiceSvc.GetByID(uint(id))
	response.Success(c, inv)
}

// AddPayInvoicePage returns page data for adding credit payment to an invoice (admin).
func (h *InvoiceHandler) AddPayInvoicePage(c *gin.Context) {
	if h.enhancedSvc == nil {
		response.ServerError(c, "enhanced service not available")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID错误")
		return
	}

	inv, err := h.invoiceSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "ID错误")
		return
	}

	user, err := h.invoiceSvc.GetUser(inv.UserID)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	response.Success(c, gin.H{
		"invoice":       inv,
		"user_credit":   user.Balance,
		"surplus":       inv.Amount,
		"invoice_credit": 0.0,
	})
}

// DeletePayInvoice removes credit payment from an invoice (admin).
func (h *InvoiceHandler) DeletePayInvoice(c *gin.Context) {
	if h.enhancedSvc == nil {
		response.ServerError(c, "enhanced service not available")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID错误")
		return
	}

	if err := h.enhancedSvc.DeletePayInvoice(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	adminID := c.GetUint("admin_id")
	h.enhancedSvc.LogAction(uint(id), adminID, "payment_deleted", "payment record deleted", c.ClientIP())
	response.SuccessMsg(c, "删除成功")
}

// InvoicePayAfterHandle handles post-payment processing (admin).
func (h *InvoiceHandler) InvoicePayAfterHandle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID错误")
		return
	}

	inv, err := h.invoiceSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "ID错误")
		return
	}

	if inv.Status != 1 {
		response.BadRequest(c, "invoice is not paid")
		return
	}

	if h.enhancedSvc != nil {
		adminID := c.GetUint("admin_id")
		h.enhancedSvc.LogAction(uint(id), adminID, "pay_after_handle", "post-payment processing triggered", c.ClientIP())
	}

	response.Success(c, gin.H{
		"invoice_id": id,
		"status":     "processed",
	})
}

// BatchStatus updates status for multiple invoices (admin).
func (h *InvoiceHandler) BatchStatus(c *gin.Context) {
	var req struct {
		InvoiceIDs []uint `json:"invoice_ids" binding:"required"`
		Status     int    `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	for _, id := range req.InvoiceIDs {
		if err := h.invoiceSvc.UpdateInvoice(id, map[string]interface{}{"status": req.Status}); err != nil {
			h.log.Errorf("failed to update invoice %d status: %v", id, err)
		}
	}

	response.SuccessMsg(c, "请求成功")
}

// MarkPaid marks invoices as paid (admin).
func (h *InvoiceHandler) MarkPaid(c *gin.Context) {
	var req struct {
		InvoiceIDs []uint `json:"invoice_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	for _, id := range req.InvoiceIDs {
		if err := h.invoiceSvc.MarkPaid(id); err != nil {
			h.log.Errorf("failed to mark invoice %d as paid: %v", id, err)
		}
	}

	response.SuccessMsg(c, "请求成功")
}

// MarkUnpaid marks invoices as unpaid (admin).
func (h *InvoiceHandler) MarkUnpaid(c *gin.Context) {
	var req struct {
		InvoiceIDs []uint `json:"invoice_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	for _, id := range req.InvoiceIDs {
		if err := h.invoiceSvc.UpdateInvoice(id, map[string]interface{}{"status": 0, "paid_at": nil}); err != nil {
			h.log.Errorf("failed to mark invoice %d as unpaid: %v", id, err)
		}
	}

	response.SuccessMsg(c, "请求成功")
}

// BatchDelete deletes multiple invoices (admin).
func (h *InvoiceHandler) BatchDelete(c *gin.Context) {
	var req struct {
		InvoiceIDs []uint `json:"invoice_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	for _, id := range req.InvoiceIDs {
		if err := h.invoiceSvc.DeleteInvoice(id); err != nil {
			h.log.Errorf("failed to delete invoice %d: %v", id, err)
		}
	}

	response.SuccessMsg(c, "删除成功")
}

// SendInvoiceEmail sends invoice email (admin).
func (h *InvoiceHandler) SendInvoiceEmail(c *gin.Context) {
	if h.enhancedSvc == nil {
		response.ServerError(c, "enhanced service not available")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID错误")
		return
	}

	email := c.Query("email")
	if err := h.enhancedSvc.SendInvoiceEmail(uint(id), email); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "发送成功")
}

// Refund processes a refund for an invoice (admin).
func (h *InvoiceHandler) Refund(c *gin.Context) {
	if h.enhancedSvc == nil {
		response.ServerError(c, "enhanced service not available")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID错误")
		return
	}

	var req struct {
		Amount float64 `json:"amount" binding:"required,gt=0"`
		Reason string  `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	refund, err := h.enhancedSvc.RefundInvoice(uint(id), req.Amount, req.Reason)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, refund)
}

// GetSummary returns invoice summary statistics (admin).
func (h *InvoiceHandler) GetSummary(c *gin.Context) {
	if h.enhancedSvc == nil {
		response.ServerError(c, "enhanced service not available")
		return
	}

	summary, err := h.enhancedSvc.GetInvoiceSummary()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, summary)
}

// GetRenewInvoices returns renewal invoices (admin).
func (h *InvoiceHandler) GetRenewInvoices(c *gin.Context) {
	if h.enhancedSvc == nil {
		response.ServerError(c, "enhanced service not available")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	invoices, total, err := h.enhancedSvc.GetRenewInvoices(page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, invoices, total, page, pageSize)
}

// ==================== P1-12: Notes ====================

// NotesPage 获取发票备注页面数据
func (h *InvoiceHandler) NotesPage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID错误")
		return
	}

	inv, err := h.invoiceSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "ID错误")
		return
	}

	response.Success(c, gin.H{
		"id":    id,
		"notes": inv.Notes,
	})
}

// Notes 更新发票备注
func (h *InvoiceHandler) Notes(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID错误")
		return
	}

	var req struct {
		Notes string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.invoiceSvc.UpdateInvoice(uint(id), map[string]interface{}{"notes": req.Notes}); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMsg(c, "备注已更新")
}

// ==================== P2-14: Enhanced GetList ====================

// GetListEnhanced 发票列表（支持复杂筛选：金额区间/时间范围/行项描述/付款方式细分）
func (h *InvoiceHandler) GetListEnhanced(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	orderField := c.DefaultQuery("order", "id")
	sort := c.DefaultQuery("sort", "desc")

	var status *int
	if s := c.Query("status"); s != "" {
		v, _ := strconv.Atoi(s)
		status = &v
	}
	var userID *uint
	if u := c.Query("uid"); u != "" {
		v, _ := strconv.ParseUint(u, 10, 64)
		uid := uint(v)
		userID = &uid
	}

	// 金额区间
	totalSmall := c.Query("total_small")
	totalBig := c.Query("total_big")
	// 时间范围
	createTimeStart := c.Query("create_time_start")
	createTimeEnd := c.Query("create_time_end")
	dueTimeStart := c.Query("due_time_start")
	dueTimeEnd := c.Query("due_time_end")
	paidTimeStart := c.Query("paid_time_start")
	paidTimeEnd := c.Query("paid_time_end")
	// 付款方式细分
	payment := c.Query("payment")
	// 行项描述
	lineItemDesc := c.Query("lineitem_desc")
	// 发票类型
	invType := c.Query("type")
	// 销售
	saleID := c.Query("sale_id")
	// 发票号
	invoiceID := c.Query("invoice_id")

	invoices, total, err := h.invoiceSvc.GetListEnhanced(page, pageSize, status, userID,
		totalSmall, totalBig, createTimeStart, createTimeEnd,
		dueTimeStart, dueTimeEnd, paidTimeStart, paidTimeEnd,
		payment, lineItemDesc, invType, saleID, invoiceID,
		orderField, sort)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, invoices, total, page, pageSize)
}

// ==================== 新增缺失方法 ====================

// EditItem updates an invoice item (admin).
// PUT /admin/invoices/items/:id
func (h *InvoiceHandler) EditItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID错误")
		return
	}

	var req struct {
		Description *string  `json:"description"`
		Amount      *float64 `json:"amount"`
		Taxed       *int     `json:"taxed"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Amount != nil {
		updates["amount"] = *req.Amount
	}
	if req.Taxed != nil {
		updates["taxed"] = *req.Taxed
	}

	if len(updates) == 0 {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.invoiceSvc.UpdateInvoiceItem(uint(id), updates); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMsg(c, "请求成功")
}

// DeleteItems deletes multiple invoice items (admin).
// DELETE /admin/invoices/items
func (h *InvoiceHandler) DeleteItems(c *gin.Context) {
	var req struct {
		ItemIDs []uint `json:"item_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	for _, id := range req.ItemIDs {
		if err := h.invoiceSvc.DeleteInvoiceItem(id); err != nil {
			h.log.Errorf("failed to delete invoice item %d: %v", id, err)
		}
	}

	response.SuccessMsg(c, "删除成功")
}

// DelAccount removes a payment record from an invoice (admin).
// DELETE /admin/invoices/:id/account
func (h *InvoiceHandler) DelAccount(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID错误")
		return
	}

	if err := h.invoiceSvc.DeleteInvoiceAccount(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	adminID := c.GetUint("admin_id")
	if h.enhancedSvc != nil {
		h.enhancedSvc.LogAction(uint(id), adminID, "account_deleted", "payment account deleted", c.ClientIP())
	}

	response.SuccessMsg(c, "删除成功")
}
