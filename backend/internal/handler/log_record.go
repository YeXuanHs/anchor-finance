package handler

import (
	"strconv"
	"time"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type LogRecordHandler struct {
	logRecordSvc *service.LogRecordService
	log          *logger.Logger
}

func NewLogRecordHandler(logRecordSvc *service.LogRecordService, log *logger.Logger) *LogRecordHandler {
	return &LogRecordHandler{logRecordSvc: logRecordSvc, log: log}
}

// List returns paginated log records with filters.
func (h *LogRecordHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	action := c.Query("action")
	module := c.Query("module")
	keyword := c.Query("keyword")

	var adminID *uint
	if aid := c.Query("admin_id"); aid != "" {
		v, _ := strconv.ParseUint(aid, 10, 64)
		id := uint(v)
		adminID = &id
	}

	var status *int16
	if s := c.Query("status"); s != "" {
		v, _ := strconv.Atoi(s)
		sv := int16(v)
		status = &sv
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

	records, total, err := h.logRecordSvc.List(page, pageSize, adminID, action, module, keyword, status, startTime, endTime)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, records, total, page, pageSize)
}

// GetDetail returns a single log record by ID.
func (h *LogRecordHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid log id")
		return
	}

	record, err := h.logRecordSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "log record not found")
		return
	}
	response.Success(c, record)
}

// Search searches log records by keyword.
func (h *LogRecordHandler) Search(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	if keyword == "" {
		response.BadRequest(c, "keyword is required")
		return
	}

	records, total, err := h.logRecordSvc.Search(page, pageSize, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, records, total, page, pageSize)
}

// Export exports log records to a CSV file.
func (h *LogRecordHandler) Export(c *gin.Context) {
	action := c.Query("action")
	module := c.Query("module")

	var adminID *uint
	if aid := c.Query("admin_id"); aid != "" {
		v, _ := strconv.ParseUint(aid, 10, 64)
		id := uint(v)
		adminID = &id
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

	filepath, err := h.logRecordSvc.Export(adminID, action, module, startTime, endTime)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filepath)
	c.Header("Content-Type", "text/csv")
	c.File(filepath)
}

// Stats returns log statistics.
func (h *LogRecordHandler) Stats(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))

	stats, err := h.logRecordSvc.Stats(days)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, stats)
}

// ModuleStats returns log counts grouped by module.
func (h *LogRecordHandler) ModuleStats(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))

	stats, err := h.logRecordSvc.ModuleStats(days)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, stats)
}

// ActionStats returns log counts grouped by action.
func (h *LogRecordHandler) ActionStats(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))

	stats, err := h.logRecordSvc.ActionStats(days)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, stats)
}

// Cleanup deletes old log records.
func (h *LogRecordHandler) Cleanup(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "90"))
	if days < 7 {
		response.BadRequest(c, "cleanup days must be at least 7")
		return
	}

	count, err := h.logRecordSvc.Cleanup(days)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{
		"deleted_count": count,
		"older_than":    days,
	})
}
