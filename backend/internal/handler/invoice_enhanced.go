package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// InvoiceEnhancedHandler 账单增强处理器
type InvoiceEnhancedHandler struct {
	svc *service.InvoiceEnhancedService
	log *logger.Logger
}

// NewInvoiceEnhancedHandler creates a new InvoiceEnhancedHandler.
func NewInvoiceEnhancedHandler(svc *service.InvoiceEnhancedService, log *logger.Logger) *InvoiceEnhancedHandler {
	return &InvoiceEnhancedHandler{svc: svc, log: log}
}

// ==================== Status Filters ====================

// GetByStatus 按状态获取账单列表
func (h *InvoiceEnhancedHandler) GetByStatus(c *gin.Context) {
	status, err := strconv.Atoi(c.Param("status"))
	if err != nil {
		response.BadRequest(c, "invalid status")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	invoices, total, err := h.svc.GetInvoicesByStatus(status, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, invoices, total, page, pageSize)
}

// GetPaidInvoices 获取已支付账单
func (h *InvoiceEnhancedHandler) GetPaidInvoices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	invoices, total, err := h.svc.GetPaidInvoices(page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, invoices, total, page, pageSize)
}

// GetUnpaidInvoices 获取未支付账单
func (h *InvoiceEnhancedHandler) GetUnpaidInvoices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	invoices, total, err := h.svc.GetUnpaidInvoices(page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, invoices, total, page, pageSize)
}

// GetCancelledInvoices 获取已取消账单
func (h *InvoiceEnhancedHandler) GetCancelledInvoices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	invoices, total, err := h.svc.GetCancelledInvoices(page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, invoices, total, page, pageSize)
}

// GetOverdueInvoices 获取逾期账单
func (h *InvoiceEnhancedHandler) GetOverdueInvoices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	invoices, total, err := h.svc.GetOverdueInvoicesPage(page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, invoices, total, page, pageSize)
}

// GetInvoiceSummary 获取账单统计汇总
func (h *InvoiceEnhancedHandler) GetInvoiceSummary(c *gin.Context) {
	summary, err := h.svc.GetInvoiceSummary()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, summary)
}

// ==================== Invoice Operations ====================

// AddPayInvoice 添加支付记录
func (h *InvoiceEnhancedHandler) AddPayInvoice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid invoice id")
		return
	}

	var req struct {
		Amount float64 `json:"amount" binding:"required,gt=0"`
		Method string  `json:"method" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	invoice, err := h.svc.AddPayInvoice(uint(id), req.Amount, req.Method)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 记录操作日志
	ip := c.ClientIP()
	h.svc.LogAction(uint(id), c.GetUint("user_id"), "payment_added", "手动添加支付记录", ip)

	response.Success(c, invoice)
}

// DeletePayInvoice 删除支付记录
func (h *InvoiceEnhancedHandler) DeletePayInvoice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid invoice id")
		return
	}

	if err := h.svc.DeletePayInvoice(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	ip := c.ClientIP()
	h.svc.LogAction(uint(id), c.GetUint("user_id"), "payment_deleted", "删除支付记录", ip)

	response.SuccessMsg(c, "payment record deleted")
}

// RefundInvoice 处理退款
func (h *InvoiceEnhancedHandler) RefundInvoice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid invoice id")
		return
	}

	var req struct {
		Amount float64 `json:"amount" binding:"required,gt=0"`
		Reason string  `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	refund, err := h.svc.RefundInvoice(uint(id), req.Amount, req.Reason)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	ip := c.ClientIP()
	h.svc.LogAction(uint(id), c.GetUint("user_id"), "refund", req.Reason, ip)

	response.Success(c, refund)
}

// GetRefundPage 获取退款选项
func (h *InvoiceEnhancedHandler) GetRefundPage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid invoice id")
		return
	}

	data, err := h.svc.GetRefundPage(uint(id))
	if err != nil {
		response.NotFound(c, "invoice not found")
		return
	}
	response.Success(c, data)
}

// ==================== Invoice Notes ====================

// AddNote 添加管理员备注
func (h *InvoiceEnhancedHandler) AddNote(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid invoice id")
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	adminID := c.GetUint("user_id")
	note, err := h.svc.AddNote(uint(id), adminID, req.Content)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	ip := c.ClientIP()
	h.svc.LogAction(uint(id), adminID, "note_added", "添加管理员备注", ip)

	response.Success(c, note)
}

// GetNotes 获取账单备注
func (h *InvoiceEnhancedHandler) GetNotes(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid invoice id")
		return
	}

	notes, err := h.svc.GetNotes(uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, notes)
}

// DeleteNote 删除备注
func (h *InvoiceEnhancedHandler) DeleteNote(c *gin.Context) {
	noteID, err := strconv.ParseUint(c.Param("note_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid note id")
		return
	}

	if err := h.svc.DeleteNote(uint(noteID)); err != nil {
		response.NotFound(c, "note not found")
		return
	}
	response.SuccessMsg(c, "note deleted")
}

// ==================== Combine Invoices ====================

// GetCombineInvoices 获取可合并账单
func (h *InvoiceEnhancedHandler) GetCombineInvoices(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	invoices, err := h.svc.GetCombineInvoices(uint(userID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, invoices)
}

// CombineInvoices 合并账单
func (h *InvoiceEnhancedHandler) CombineInvoices(c *gin.Context) {
	var req struct {
		InvoiceIDs []uint `json:"invoice_ids" binding:"required,min=2"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	combined, err := h.svc.CombineInvoices(req.InvoiceIDs)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	ip := c.ClientIP()
	h.svc.LogAction(combined.ID, c.GetUint("user_id"), "combined", "合并账单", ip)

	response.Success(c, combined)
}

// ==================== Invoice Email ====================

// SendInvoiceEmail 发送账单邮件
func (h *InvoiceEnhancedHandler) SendInvoiceEmail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid invoice id")
		return
	}

	var req struct {
		Email string `json:"email"`
	}
	c.ShouldBindJSON(&req)

	if err := h.svc.SendInvoiceEmail(uint(id), req.Email); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "invoice email sent")
}

// InvoiceEmail 获取账单邮件模板
func (h *InvoiceEnhancedHandler) InvoiceEmail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid invoice id")
		return
	}

	data, err := h.svc.InvoiceEmail(uint(id))
	if err != nil {
		response.NotFound(c, "invoice not found")
		return
	}
	response.Success(c, data)
}

// ==================== Renew Invoices ====================

// CreateRenewInvoice 创建续费账单
func (h *InvoiceEnhancedHandler) CreateRenewInvoice(c *gin.Context) {
	var req struct {
		HostID uint   `json:"host_id" binding:"required"`
		Cycle  string `json:"cycle" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	invoice, err := h.svc.CreateRenewInvoice(req.HostID, req.Cycle)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	ip := c.ClientIP()
	h.svc.LogAction(invoice.ID, c.GetUint("user_id"), "renew_created", "创建续费账单", ip)

	response.Success(c, invoice)
}

// GetRenewInvoices 续费账单列表
func (h *InvoiceEnhancedHandler) GetRenewInvoices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	invoices, total, err := h.svc.GetRenewInvoices(page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, invoices, total, page, pageSize)
}

// ==================== Invoice Log ====================

// GetInvoiceLog 获取账单操作日志
func (h *InvoiceEnhancedHandler) GetInvoiceLog(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid invoice id")
		return
	}

	logs, err := h.svc.GetInvoiceLog(uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, logs)
}

// ==================== Search ====================

// SearchInvoices 搜索账单
func (h *InvoiceEnhancedHandler) SearchInvoices(c *gin.Context) {
	query := c.Query("q")

	invoices, err := h.svc.SearchInvoices(query)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, invoices)
}

// SearchPage 高级搜索
func (h *InvoiceEnhancedHandler) SearchPage(c *gin.Context) {
	var filters service.InvoiceSearchFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	invoices, total, err := h.svc.SearchPage(filters)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, invoices, total, filters.Page, filters.PageSize)
}

// ==================== Duplicate ====================

// DuplicateInvoice 复制账单
func (h *InvoiceEnhancedHandler) DuplicateInvoice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid invoice id")
		return
	}

	invoice, err := h.svc.DuplicateInvoice(uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	ip := c.ClientIP()
	h.svc.LogAction(invoice.ID, c.GetUint("user_id"), "duplicated", "复制账单", ip)

	response.Success(c, invoice)
}
