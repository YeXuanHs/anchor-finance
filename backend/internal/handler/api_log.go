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

	var userID uint
	if uid := c.Query("user_id"); uid != "" {
		v, _ := strconv.ParseUint(uid, 10, 64)
		userID = uint(v)
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

	items, total, err := h.svc.List(page, pageSize, userID, method, startTime, endTime)
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
