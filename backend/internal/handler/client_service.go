package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// ClientServiceHandler handles client service HTTP requests.
type ClientServiceHandler struct {
	svc *service.ClientServiceService
	log *logger.Logger
}

// NewClientServiceHandler creates a new ClientServiceHandler.
func NewClientServiceHandler(svc *service.ClientServiceService, log *logger.Logger) *ClientServiceHandler {
	return &ClientServiceHandler{svc: svc, log: log}
}

// List returns a filtered, paginated list of client services.
// GET /client-services
func (h *ClientServiceHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var userID uint
	if uid := c.Query("user_id"); uid != "" {
		v, _ := strconv.ParseUint(uid, 10, 64)
		userID = uint(v)
	}
	var status int16
	if s := c.Query("status"); s != "" {
		v, _ := strconv.ParseInt(s, 10, 16)
		status = int16(v)
	}

	items, total, err := h.svc.GetList(page, pageSize, userID, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, items, total, page, pageSize)
}

// Get returns a single client service by ID.
// GET /client-services/:id
func (h *ClientServiceHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid service id")
		return
	}

	svc, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "service not found")
		return
	}
	response.Success(c, svc)
}

// Open creates a new client service instance.
// POST /client-services
func (h *ClientServiceHandler) Open(c *gin.Context) {
	var req service.OpenServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	svc, err := h.svc.Open(req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, svc)
}

// Update modifies client service metadata.
// PUT /client-services/:id
func (h *ClientServiceHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid service id")
		return
	}

	var req service.UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.Update(uint(id), req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "service updated")
}

// Suspend pauses an active service.
// POST /client-services/:id/suspend
func (h *ClientServiceHandler) Suspend(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid service id")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)

	if err := h.svc.Suspend(uint(id), req.Reason); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "service suspended")
}

// Terminate permanently terminates a service.
// POST /client-services/:id/terminate
func (h *ClientServiceHandler) Terminate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid service id")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)

	if err := h.svc.Terminate(uint(id), req.Reason); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "service terminated")
}

// Renew extends a service's expiry.
// POST /client-services/:id/renew
func (h *ClientServiceHandler) Renew(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid service id")
		return
	}

	var req struct {
		Months int `json:"months" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.Renew(uint(id), req.Months); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "service renewed")
}

// Resume reactivates a suspended service.
// POST /client-services/:id/resume
func (h *ClientServiceHandler) Resume(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid service id")
		return
	}

	if err := h.svc.Resume(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// return the updated service
	svc, _ := h.svc.GetByID(uint(id))
	response.Success(c, svc)
}

// MyServices returns services for the authenticated user.
// GET /user/services
func (h *ClientServiceHandler) MyServices(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	items, total, err := h.svc.GetList(page, pageSize, userID, 0)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, items, total, page, pageSize)
}

// MyServiceDetail returns a single service for the authenticated user.
// GET /user/services/:id
func (h *ClientServiceHandler) MyServiceDetail(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid service id")
		return
	}

	svc, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "service not found")
		return
	}
	if svc.UserID != userID {
		response.NotFound(c, "service not found")
		return
	}
	response.Success(c, svc)
}

// MyServiceAutoRenew toggles auto-renew for the authenticated user's service.
// POST /user/services/:id/auto-renew
func (h *ClientServiceHandler) MyServiceAutoRenew(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid service id")
		return
	}

	svc, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "service not found")
		return
	}
	if svc.UserID != userID {
		response.NotFound(c, "service not found")
		return
	}

	var req struct {
		AutoRenew bool `json:"auto_renew"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.Update(uint(id), service.UpdateServiceRequest{AutoRenew: &req.AutoRenew}); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "auto-renew updated")
}

// GetByStatus returns service status constants.
// GET /client-services/status-constants
func (h *ClientServiceHandler) GetByStatus(c *gin.Context) {
	response.Success(c, map[string]int16{
		"active":     model.ClientServiceActive,
		"suspended":  model.ClientServiceSuspended,
		"pending":    model.ClientServicePending,
		"terminated": model.ClientServiceTerminated,
		"expired":    model.ClientServiceExpired,
	})
}

// ==================== 新增缺失方法 ====================

// PostTransfer transfers a service to another user (admin).
// POST /admin/client-services/transfer
func (h *ClientServiceHandler) PostTransfer(c *gin.Context) {
	var req struct {
		HostID      uint `json:"hostid" binding:"required"`
		TransferUID uint `json:"transfer_uid" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.TransferService(req.HostID, req.TransferUID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMsg(c, "service transferred successfully")
}

// DeleteHost deletes a host/service (admin).
// DELETE /admin/client-services/host
func (h *ClientServiceHandler) DeleteHost(c *gin.Context) {
	var req struct {
		HostIDs []uint `json:"hostid" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Try single ID
		id, err2 := strconv.ParseUint(c.Param("id"), 10, 64)
		if err2 != nil {
			response.BadRequest(c, err.Error())
			return
		}
		req.HostIDs = []uint{uint(id)}
	}

	if err := h.svc.DeleteHosts(req.HostIDs); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMsg(c, "services deleted")
}

// PostBatchRenewPage returns batch renew page data (admin).
// POST /admin/client-services/batch-renew-page
func (h *ClientServiceHandler) PostBatchRenewPage(c *gin.Context) {
	var req struct {
		UID     uint            `json:"uid" binding:"required"`
		HostIDs []uint          `json:"host_ids" binding:"required"`
		Cycles  map[uint]string `json:"cycles"`
		Amount  map[uint]float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	data, err := h.svc.GetBatchRenewPage(req.UID, req.HostIDs, req.Cycles, req.Amount)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, data)
}

// PostBatchRenew processes batch renewal (admin).
// POST /admin/client-services/batch-renew
func (h *ClientServiceHandler) PostBatchRenew(c *gin.Context) {
	var req struct {
		UID     uint            `json:"uid" binding:"required"`
		HostIDs []uint          `json:"host_ids" binding:"required"`
		Cycles  map[uint]string `json:"cycles"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.svc.ProcessBatchRenew(req.UID, req.HostIDs, req.Cycles)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, result)
}

// GetApplyCreditPage returns credit application page data (admin).
// GET /admin/client-services/apply-credit-page
func (h *ClientServiceHandler) GetApplyCreditPage(c *gin.Context) {
	invoiceID := c.Query("invoiceid")
	uid := c.Query("uid")

	if invoiceID == "" || uid == "" {
		response.BadRequest(c, "invoiceid and uid are required")
		return
	}

	data, err := h.svc.GetApplyCreditPageData(uid, invoiceID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, data)
}

// ApplyCredit applies user credit to an invoice (admin).
// POST /admin/client-services/apply-credit
func (h *ClientServiceHandler) ApplyCredit(c *gin.Context) {
	var req struct {
		UID       uint `json:"uid" binding:"required"`
		InvoiceID uint `json:"invoiceid" binding:"required"`
		Enough    int  `json:"enough"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.svc.ApplyCredit(req.UID, req.InvoiceID, req.Enough)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, result)
}

// PostSearchClient searches clients by keyword (admin).
// POST /admin/client-services/search-client
func (h *ClientServiceHandler) PostSearchClient(c *gin.Context) {
	var req struct {
		ClientID string `json:"client_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	clients, err := h.svc.SearchClients(req.ClientID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"client_list": clients})
}

// GetRefundPage returns refund page data (admin).
// GET /admin/client-services/refund-page
func (h *ClientServiceHandler) GetRefundPage(c *gin.Context) {
	hid := c.Query("hid")
	if hid == "" {
		response.BadRequest(c, "hid is required")
		return
	}

	data, err := h.svc.GetRefundPageData(hid)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, data)
}

// Refund processes a refund for a service (admin).
// POST /admin/client-services/refund
func (h *ClientServiceHandler) Refund(c *gin.Context) {
	var req struct {
		HID          uint    `json:"hid" binding:"required"`
		RefundMethod string  `json:"refund_method" binding:"required"` // day/full/custom
		RefundType   string  `json:"refund_type" binding:"required"`   // addascredit/only
		Amount       float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.ProcessRefund(req.HID, req.RefundMethod, req.RefundType, req.Amount); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMsg(c, "refund processed")
}
