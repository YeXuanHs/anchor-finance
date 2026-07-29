package handler

import (
	"strconv"
	"time"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// LoginLogHandler handles login log management requests.
type LoginLogHandler struct {
	svc *service.LoginLogService
	log *logger.Logger
}

// NewLoginLogHandler creates a new LoginLogHandler.
func NewLoginLogHandler(svc *service.LoginLogService, log *logger.Logger) *LoginLogHandler {
	return &LoginLogHandler{svc: svc, log: log}
}

// List returns paginated login logs with filters.
func (h *LoginLogHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	username := c.Query("username")
	ip := c.Query("ip")
	status := c.Query("status")

	var userID *uint
	if uid := c.Query("user_id"); uid != "" {
		v, _ := strconv.ParseUint(uid, 10, 64)
		id := uint(v)
		userID = &id
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

	items, total, err := h.svc.List(page, pageSize, userID, username, ip, status, startTime, endTime)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, items, total, page, pageSize)
}

// GetDetail returns a single login log by ID.
func (h *LoginLogHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid login log id")
		return
	}

	item, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "login log not found")
		return
	}
	response.Success(c, item)
}

// Delete deletes a login log.
func (h *LoginLogHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid login log id")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "login log deleted")
}

// Cleanup deletes old login logs.
func (h *LoginLogHandler) Cleanup(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "90"))
	if days < 7 {
		response.BadRequest(c, "cleanup days must be at least 7")
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

// GetStats returns login statistics.
func (h *LoginLogHandler) GetStats(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))

	stats, err := h.svc.GetStats(days)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, stats)
}

// GetFailedAttempts returns failed login attempts for a user.
func (h *LoginLogHandler) GetFailedAttempts(c *gin.Context) {
	var userID *uint
	if uid := c.Query("user_id"); uid != "" {
		v, _ := strconv.ParseUint(uid, 10, 64)
		id := uint(v)
		userID = &id
	}
	ip := c.Query("ip")
	minutes, _ := strconv.Atoi(c.DefaultQuery("minutes", "30"))

	count, err := h.svc.GetFailedAttempts(userID, ip, minutes)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{
		"failed_count": count,
		"minutes":      minutes,
	})
}

// Export exports login logs to a CSV file.
func (h *LoginLogHandler) Export(c *gin.Context) {
	username := c.Query("username")
	ip := c.Query("ip")

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

	filepath, err := h.svc.Export(username, ip, startTime, endTime)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filepath)
	c.Header("Content-Type", "text/csv")
	c.File(filepath)
}
