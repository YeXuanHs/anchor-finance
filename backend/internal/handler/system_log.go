package handler

import (
	"strconv"
	"time"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// SystemLogHandler handles system log management requests.
type SystemLogHandler struct {
	svc *service.SystemLogService
	log *logger.Logger
}

// NewSystemLogHandler creates a new SystemLogHandler.
func NewSystemLogHandler(svc *service.SystemLogService, log *logger.Logger) *SystemLogHandler {
	return &SystemLogHandler{svc: svc, log: log}
}

// List returns paginated system logs with filters.
func (h *SystemLogHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	level := c.Query("level")
	module := c.Query("module")
	keyword := c.Query("keyword")

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

	items, total, err := h.svc.List(page, pageSize, level, module, keyword, startTime, endTime)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, items, total, page, pageSize)
}

// GetDetail returns a single system log by ID.
func (h *SystemLogHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid system log id")
		return
	}

	item, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "system log not found")
		return
	}
	response.Success(c, item)
}

// Delete deletes a system log.
func (h *SystemLogHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid system log id")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "system log deleted")
}

// Cleanup deletes old system logs.
func (h *SystemLogHandler) Cleanup(c *gin.Context) {
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

// GetStats returns system log statistics.
func (h *SystemLogHandler) GetStats(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))

	stats, err := h.svc.GetStats(days)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, stats)
}

// GetLevelStats returns log counts grouped by level.
func (h *SystemLogHandler) GetLevelStats(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))

	stats, err := h.svc.GetLevelStats(days)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, stats)
}

// GetModuleStats returns log counts grouped by module.
func (h *SystemLogHandler) GetModuleStats(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))

	stats, err := h.svc.GetModuleStats(days)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, stats)
}

// Export exports system logs to a CSV file.
func (h *SystemLogHandler) Export(c *gin.Context) {
	level := c.Query("level")
	module := c.Query("module")

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

	filepath, err := h.svc.Export(level, module, startTime, endTime)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filepath)
	c.Header("Content-Type", "text/csv")
	c.File(filepath)
}

// ClearByLevel deletes all logs of a specific level.
func (h *SystemLogHandler) ClearByLevel(c *gin.Context) {
	level := c.Query("level")
	if level == "" {
		response.BadRequest(c, "level is required")
		return
	}

	count, err := h.svc.ClearByLevel(level)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{
		"deleted_count": count,
		"level":         level,
	})
}
