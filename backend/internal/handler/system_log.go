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
	response.SuccessMsg(c, "删除成功")
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

// ─── LogRecord Admin Methods (from zjmf LogRecordController) ───

// GetSystemLog returns paginated activity logs (non-System user).
func (h *SystemLogHandler) GetSystemLog(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	searchName := c.Query("search_name")
	searchDesc := c.Query("search_desc")
	searchIP := c.Query("search_ip")
	searchTime := c.Query("search_time")

	items, total, err := h.svc.AdminGetSystemLog(page, pageSize, searchName, searchDesc, searchIP, searchTime, false)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"count": total, "log_list": items})
}

// GetCronSystemLog returns paginated activity logs for System user.
func (h *SystemLogHandler) GetCronSystemLog(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	searchName := c.Query("search_name")
	searchDesc := c.Query("search_desc")
	searchIP := c.Query("search_ip")
	searchTime := c.Query("search_time")

	items, total, err := h.svc.AdminGetSystemLog(page, pageSize, searchName, searchDesc, searchIP, searchTime, true)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"count": total, "log_list": items})
}

// GetAdminLog returns paginated admin login logs.
func (h *SystemLogHandler) GetAdminLog(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	searchName := c.Query("search_name")
	searchIP := c.Query("search_ip")
	searchTime := c.Query("search_time")

	items, total, err := h.svc.AdminGetAdminLog(page, pageSize, searchName, searchIP, searchTime)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"count": total, "data": items})
}

// GetNotifyLog returns paginated notification logs.
func (h *SystemLogHandler) GetNotifyLog(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	message := c.Query("message")
	logType := c.Query("type")
	searchTime := c.Query("search_time")

	items, total, err := h.svc.AdminGetNotifyLog(page, pageSize, message, logType, searchTime)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"count": total, "data": items, "type": []string{"email", "sms", "wechat"}})
}

// GetEmailLog returns paginated email logs.
func (h *SystemLogHandler) GetEmailLog(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	subject := c.Query("subject")
	username := c.Query("username")
	searchTime := c.Query("search_time")
	uid, _ := strconv.ParseUint(c.Query("uid"), 10, 64)

	items, total, err := h.svc.AdminGetEmailLog(page, pageSize, subject, username, searchTime, uint(uid))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"count": total, "data": items})
}

// GetEmailDetail returns email detail by ID.
func (h *SystemLogHandler) GetEmailDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	if id == 0 {
		response.BadRequest(c, "id is required")
		return
	}

	detail, err := h.svc.AdminGetEmailDetail(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, gin.H{"detail": detail})
}

// GetWechatLog returns paginated WeChat logs.
func (h *SystemLogHandler) GetWechatLog(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	items, total, err := h.svc.AdminGetWechatLog(page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"count": total, "data": items})
}

// GetSmsLog returns paginated SMS logs.
func (h *SystemLogHandler) GetSmsLog(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	phone := c.Query("phone")
	username := c.Query("username")
	searchTime := c.Query("search_time")
	uid, _ := strconv.ParseUint(c.Query("uid"), 10, 64)

	items, total, err := h.svc.AdminGetSmsLog(page, pageSize, phone, username, searchTime, uint(uid))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"count": total, "data": items})
}

// GetSmsLogM returns paginated SMS logs for a specific user.
func (h *SystemLogHandler) GetSmsLogM(c *gin.Context) {
	uid, _ := strconv.ParseUint(c.Query("uid"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	phone := c.Query("phone")
	searchTime := c.Query("search_time")

	items, total, err := h.svc.AdminGetSmsLog(page, pageSize, phone, "", searchTime, uint(uid))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"count": total, "data": items})
}

// GetSystemMessageLog returns paginated system message logs.
func (h *SystemLogHandler) GetSystemMessageLog(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	uid, _ := strconv.ParseUint(c.Query("uid"), 10, 64)
	keywords := c.Query("keywords")
	username := c.Query("username")
	readType := c.Query("read_type")
	searchTimeStart := c.Query("search_time_start")
	searchTimeEnd := c.Query("search_time_end")

	items, total, err := h.svc.AdminGetSystemMessageLog(page, pageSize, uint(uid), keywords, username, readType, searchTimeStart, searchTimeEnd)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"count": total, "list": items})
}

// GetApiLog returns paginated API resource logs.
func (h *SystemLogHandler) GetApiLog(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	keywords := c.Query("keywords")
	uid, _ := strconv.ParseUint(c.Query("uid"), 10, 64)
	searchTime := c.Query("time")

	items, total, err := h.svc.AdminGetApiLog(page, pageSize, keywords, uint(uid), searchTime)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"count": total, "list": items})
}

// GetDeleteLogPage returns the count of logs for a given type.
func (h *SystemLogHandler) GetDeleteLogPage(c *gin.Context) {
	logType := c.Query("type")
	if logType == "" {
		logType = "system_log"
	}

	count, err := h.svc.AdminGetLogCount(logType)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"count": count})
}

// GetAffirmDeleteLogPage returns the count of logs for deletion confirmation.
func (h *SystemLogHandler) GetAffirmDeleteLogPage(c *gin.Context) {
	logType := c.Query("type")
	timeStr := c.Query("time")

	count, err := h.svc.AdminGetLogCountBefore(logType, timeStr)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"count": count})
}

// DeleteLog deletes logs by type and optional time filter.
func (h *SystemLogHandler) DeleteLog(c *gin.Context) {
	logType := c.Query("type")
	timeStr := c.Query("time")

	count, err := h.svc.AdminDeleteLog(logType, timeStr)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"deleted_count": count})
}
