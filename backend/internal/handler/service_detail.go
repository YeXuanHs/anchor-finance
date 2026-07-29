package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// ServiceDetailHandler handles service detail requests.
type ServiceDetailHandler struct {
	svc *service.ServiceDetailService
	log *logger.Logger
}

// NewServiceDetailHandler creates a new ServiceDetailHandler.
func NewServiceDetailHandler(svc *service.ServiceDetailService, log *logger.Logger) *ServiceDetailHandler {
	return &ServiceDetailHandler{svc: svc, log: log}
}

// List returns paginated service details.
func (h *ServiceDetailHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	var userID *uint
	if uid := c.Query("user_id"); uid != "" {
		v, _ := strconv.ParseUint(uid, 10, 64)
		id := uint(v)
		userID = &id
	}

	var productID *uint
	if pid := c.Query("product_id"); pid != "" {
		v, _ := strconv.ParseUint(pid, 10, 64)
		id := uint(v)
		productID = &id
	}

	var status *int16
	if s := c.Query("status"); s != "" {
		v, _ := strconv.Atoi(s)
		sv := int16(v)
		status = &sv
	}

	items, total, err := h.svc.List(page, pageSize, userID, productID, keyword, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, items, total, page, pageSize)
}

// GetDetail returns a single service detail by ID.
func (h *ServiceDetailHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid service id")
		return
	}

	item, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "service not found")
		return
	}
	response.Success(c, item)
}

// GetByUser returns all services for a user.
func (h *ServiceDetailHandler) GetByUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	items, total, err := h.svc.GetByUserID(uint(userID), page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, items, total, page, pageSize)
}

// Create creates a new service detail.
func (h *ServiceDetailHandler) Create(c *gin.Context) {
	var req struct {
		UserID      uint   `json:"user_id" binding:"required"`
		ProductID   uint   `json:"product_id" binding:"required"`
		Domain      string `json:"domain"`
		Username    string `json:"username"`
		Password    string `json:"password"`
		StartDate   string `json:"start_date"`
		DueDate     string `json:"due_date"`
		Price       float64 `json:"price"`
		BillingCycle string `json:"billing_cycle"`
		Note        string `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	item, err := h.svc.Create(req.UserID, req.ProductID, req.Domain, req.Username, req.Password, req.StartDate, req.DueDate, req.Price, req.BillingCycle, req.Note)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, item)
}

// Update updates a service detail.
func (h *ServiceDetailHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid service id")
		return
	}

	var req struct {
		Domain       *string  `json:"domain"`
		Username     *string  `json:"username"`
		Password     *string  `json:"password"`
		StartDate    *string  `json:"start_date"`
		DueDate      *string  `json:"due_date"`
		Price        *float64 `json:"price"`
		BillingCycle *string  `json:"billing_cycle"`
		Note         *string  `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Domain != nil {
		updates["domain"] = *req.Domain
	}
	if req.Username != nil {
		updates["username"] = *req.Username
	}
	if req.Password != nil {
		updates["password"] = *req.Password
	}
	if req.StartDate != nil {
		updates["start_date"] = *req.StartDate
	}
	if req.DueDate != nil {
		updates["due_date"] = *req.DueDate
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.BillingCycle != nil {
		updates["billing_cycle"] = *req.BillingCycle
	}
	if req.Note != nil {
		updates["note"] = *req.Note
	}

	if err := h.svc.Update(uint(id), updates); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "service updated")
}

// Delete deletes a service detail.
func (h *ServiceDetailHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid service id")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "service deleted")
}

// Suspend suspends a service.
func (h *ServiceDetailHandler) Suspend(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid service id")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.Suspend(uint(id), req.Reason); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "service suspended")
}

// Unsuspend unsuspends a service.
func (h *ServiceDetailHandler) Unsuspend(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid service id")
		return
	}

	if err := h.svc.Unsuspend(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "service unsuspended")
}

// Terminate terminates a service.
func (h *ServiceDetailHandler) Terminate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid service id")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.Terminate(uint(id), req.Reason); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "service terminated")
}

// Renew renews a service.
func (h *ServiceDetailHandler) Renew(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid service id")
		return
	}

	var req struct {
		Period     int    `json:"period" binding:"required,min=1"`
		PeriodUnit string `json:"period_unit" binding:"required,oneof=month year"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.svc.Renew(uint(id), req.Period, req.PeriodUnit)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, result)
}

// GetStats returns service statistics.
func (h *ServiceDetailHandler) GetStats(c *gin.Context) {
	stats, err := h.svc.GetStats()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, stats)
}

// GetServiceLogs returns logs for a service.
func (h *ServiceDetailHandler) GetServiceLogs(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid service id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	logs, total, err := h.svc.GetServiceLogs(uint(id), page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, logs, total, page, pageSize)
}

// GetUserDashboard returns dashboard summary for the authenticated user.
func (h *ServiceDetailHandler) GetUserDashboard(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Unauthorized(c, "login required")
		return
	}

	stats, err := h.svc.GetServiceStats(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, stats)
}
