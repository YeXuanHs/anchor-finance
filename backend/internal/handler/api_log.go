package handler

import (
	"strconv"
	"time"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// APILogHandler handles API call log requests.
type APILogHandler struct {
	svc *service.APILogService
	log *logger.Logger
}

// NewAPILogHandler creates a new APILogHandler.
func NewAPILogHandler(svc *service.APILogService, log *logger.Logger) *APILogHandler {
	return &APILogHandler{svc: svc, log: log}
}

// List returns paginated API call logs with filters.
func (h *APILogHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	method := c.Query("method")
	path := c.Query("path")
	statusCode := c.Query("status_code")

	var apiKeyID *uint
	if akid := c.Query("api_key_id"); akid != "" {
		v, _ := strconv.ParseUint(akid, 10, 64)
		id := uint(v)
		apiKeyID = &id
	}

	var startTime, endTime *time.Time
	if st := c.Query("start_time"); st != "" {
		t, err := time.Parse("2006-01-02", st)
		if err == nil {
			startTime = &t
		}
	}
	if et := c.Query("end_time"); et != "" {
		t, err := time.Parse("2006-01-02", et)
		if err == nil {
			end := t.Add(24*time.Hour - time.Second)
			endTime = &end
		}
	}

	items, total, err := h.svc.List(page, pageSize, apiKeyID, method, path, statusCode, startTime, endTime)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, items, total, page, pageSize)
}

// GetDetail returns a single API log by ID.
func (h *APILogHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid api log id")
		return
	}

	item, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "api log not found")
		return
	}
	response.Success(c, item)
}

// Delete deletes an API log.
func (h *APILogHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid api log id")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "api log deleted")
}

// Cleanup deletes old API logs.
func (h *APILogHandler) Cleanup(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	if days < 1 {
		response.BadRequest(c, "cleanup days must be at least 1")
		return
	}

	count, err := h.svc.Cleanup(days)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{
		"deleted_count": count,
		"older_than":    days,
	})
}

// GetStats returns API call statistics.
func (h *APILogHandler) GetStats(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))

	stats, err := h.svc.GetStats(days)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, stats)
}

// GetTopEndpoints returns most frequently called endpoints.
func (h *APILogHandler) GetTopEndpoints(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	endpoints, err := h.svc.GetTopEndpoints(days, limit)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, endpoints)
}

// GetSlowRequests returns slowest API requests.
func (h *APILogHandler) GetSlowRequests(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	threshold, _ := strconv.Atoi(c.DefaultQuery("threshold", "1000"))

	requests, err := h.svc.GetSlowRequests(days, limit, threshold)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, requests)
}

// GetErrorRate returns API error rate statistics.
func (h *APILogHandler) GetErrorRate(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))

	rate, err := h.svc.GetErrorRate(days)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, rate)
}

// Export exports API logs to a CSV file.
func (h *APILogHandler) Export(c *gin.Context) {
	method := c.Query("method")
	path := c.Query("path")

	var startTime, endTime *time.Time
	if st := c.Query("start_time"); st != "" {
		t, err := time.Parse("2006-01-02", st)
		if err == nil {
			startTime = &t
		}
	}
	if et := c.Query("end_time"); et != "" {
		t, err := time.Parse("2006-01-02", et)
		if err == nil {
			end := t.Add(24*time.Hour - time.Second)
			endTime = &end
		}
	}

	filepath, err := h.svc.Export(method, path, startTime, endTime)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filepath)
	c.Header("Content-Type", "text/csv")
	c.File(filepath)
}
