package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// MultiRenewHandler handles batch renewal requests.
type MultiRenewHandler struct {
	svc *service.MultiRenewService
	log *logger.Logger
}

// NewMultiRenewHandler creates a new MultiRenewHandler.
func NewMultiRenewHandler(svc *service.MultiRenewService, log *logger.Logger) *MultiRenewHandler {
	return &MultiRenewHandler{svc: svc, log: log}
}

// List returns paginated batch renewal tasks.
func (h *MultiRenewHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var status *int16
	if s := c.Query("status"); s != "" {
		v, _ := strconv.Atoi(s)
		sv := int16(v)
		status = &sv
	}

	items, total, err := h.svc.List(page, pageSize, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, items, total, page, pageSize)
}

// GetDetail returns a single batch renewal task by ID.
func (h *MultiRenewHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid renew task id")
		return
	}

	item, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "renew task not found")
		return
	}
	response.Success(c, item)
}

// Create creates a new batch renewal task.
func (h *MultiRenewHandler) Create(c *gin.Context) {
	var req struct {
		Name       string `json:"name" binding:"required"`
		ServiceIDs []uint `json:"service_ids" binding:"required,min=1"`
		Period     int    `json:"period" binding:"required,min=1"`
		PeriodUnit string `json:"period_unit" binding:"required,oneof=month year"`
		AutoPay    bool   `json:"auto_pay"`
		Note       string `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	task, err := h.svc.Create(req.Name, req.ServiceIDs, req.Period, req.PeriodUnit, req.AutoPay, req.Note)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, task)
}

// Execute executes a batch renewal task.
func (h *MultiRenewHandler) Execute(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid renew task id")
		return
	}

	result, err := h.svc.Execute(uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, result)
}

// Cancel cancels a batch renewal task.
func (h *MultiRenewHandler) Cancel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid renew task id")
		return
	}

	if err := h.svc.Cancel(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "renew task cancelled")
}

// Delete deletes a batch renewal task.
func (h *MultiRenewHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid renew task id")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "renew task deleted")
}

// GetLogs returns logs for a batch renewal task.
func (h *MultiRenewHandler) GetLogs(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid renew task id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	logs, total, err := h.svc.GetLogs(uint(id), page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, logs, total, page, pageSize)
}

// GetStats returns batch renewal statistics.
func (h *MultiRenewHandler) GetStats(c *gin.Context) {
	stats, err := h.svc.GetStats()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, stats)
}

// Preview previews the services that would be renewed.
func (h *MultiRenewHandler) Preview(c *gin.Context) {
	var req struct {
		ServiceIDs []uint `json:"service_ids" binding:"required,min=1"`
		Period     int    `json:"period" binding:"required,min=1"`
		PeriodUnit string `json:"period_unit" binding:"required,oneof=month year"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	preview, err := h.svc.Preview(req.ServiceIDs, req.Period, req.PeriodUnit)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, preview)
}
