package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

var _ = service.HostActionRequest{}

type HostHandler struct {
	hostSvc *service.HostService
	log     *logger.Logger
}

func NewHostHandler(hostSvc *service.HostService, log *logger.Logger) *HostHandler {
	return &HostHandler{hostSvc: hostSvc, log: log}
}

// GetDetail returns a single host.
func (h *HostHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	host, err := h.hostSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "host not found")
		return
	}
	response.Success(c, host)
}

// GetList returns all hosts (admin).
func (h *HostHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	var status *int
	if s := c.Query("status"); s != "" {
		v, _ := strconv.Atoi(s)
		status = &v
	}
	var ownerID *uint
	if o := c.Query("owner_id"); o != "" {
		v, _ := strconv.ParseUint(o, 10, 64)
		uid := uint(v)
		ownerID = &uid
	}

	hosts, total, err := h.hostSvc.GetList(page, pageSize, status, keyword, ownerID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, hosts, total, page, pageSize)
}

// GetUserHosts returns hosts for the authenticated user.
func (h *HostHandler) GetUserHosts(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	hosts, total, err := h.hostSvc.GetUserHosts(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, hosts, total, page, pageSize)
}

// PerformAction executes an action on a host (admin).
func (h *HostHandler) PerformAction(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	var req service.HostActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	operatorID := c.GetUint("user_id")
	op, err := h.hostSvc.PerformAction(uint(id), operatorID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, op)
}

// GetOperations returns operations for a host (admin).
func (h *HostHandler) GetOperations(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	ops, total, err := h.hostSvc.GetHostOperations(uint(id), page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, ops, total, page, pageSize)
}

// Boot powers on a host.
func (h *HostHandler) Boot(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}
	operatorID := c.GetUint("user_id")
	op, err := h.hostSvc.PerformAction(uint(id), operatorID, service.HostActionRequest{Action: "boot"})
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, op)
}

// Shutdown powers off a host.
func (h *HostHandler) Shutdown(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}
	operatorID := c.GetUint("user_id")
	op, err := h.hostSvc.PerformAction(uint(id), operatorID, service.HostActionRequest{Action: "shutdown"})
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, op)
}

// Reboot restarts a host.
func (h *HostHandler) Reboot(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}
	operatorID := c.GetUint("user_id")
	op, err := h.hostSvc.PerformAction(uint(id), operatorID, service.HostActionRequest{Action: "reboot"})
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, op)
}

// GetExpiringHosts returns hosts expiring within N days (admin).
func (h *HostHandler) GetExpiringHosts(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))

	hosts, err := h.hostSvc.GetExpiringHosts(days)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, hosts)
}

// ==================== P3-18: GetTimetype ====================

// GetTimetype 获取时间类型选项
func (h *HostHandler) GetTimetype(c *gin.Context) {
	timeTypes := []map[string]interface{}{
		{"value": "monthly", "label": "月付"},
		{"value": "quarterly", "label": "季付"},
		{"value": "semi-annually", "label": "半年付"},
		{"value": "annually", "label": "年付"},
		{"value": "biennially", "label": "两年付"},
		{"value": "triennially", "label": "三年付"},
		{"value": "free", "label": "免费"},
		{"value": "onetime", "label": "一次性"},
	}
	response.Success(c, gin.H{"data": timeTypes})
}

// GetBilling returns billing information for a host.
// GET /hosts/:id/billing
func (h *HostHandler) GetBilling(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	// Verify host belongs to user
	host, err := h.hostSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "host not found")
		return
	}
	if host.OwnerID == nil || *host.OwnerID != userID {
		response.Forbidden(c, "access denied")
		return
	}

	// Get billing info from service
	billing, err := h.hostSvc.GetBillingInfo(uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, billing)
}

// GetLog returns operation logs for a host.
// GET /hosts/:id/log
func (h *HostHandler) GetLog(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	// Verify host belongs to user
	host, err := h.hostSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "host not found")
		return
	}
	if host.OwnerID == nil || *host.OwnerID != userID {
		response.Forbidden(c, "access denied")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	logs, total, err := h.hostSvc.GetHostOperations(uint(id), page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, logs, total, page, pageSize)
}

// GetDownload returns download files for a host.
// GET /hosts/:id/download
func (h *HostHandler) GetDownload(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	// Verify host belongs to user
	host, err := h.hostSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "host not found")
		return
	}
	if host.OwnerID == nil || *host.OwnerID != userID {
		response.Forbidden(c, "access denied")
		return
	}

	// Get download files
	downloads, err := h.hostSvc.GetDownloadFiles(uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, downloads)
}

// PostRemark adds or updates a remark for a host.
// POST /hosts/:id/remark
func (h *HostHandler) PostRemark(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host id")
		return
	}

	// Verify host belongs to user
	host, err := h.hostSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "host not found")
		return
	}
	if host.OwnerID == nil || *host.OwnerID != userID {
		response.Forbidden(c, "access denied")
		return
	}

	var req struct {
		Remark string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.hostSvc.UpdateRemark(uint(id), req.Remark); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessMsg(c, "remark updated")
}
