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
