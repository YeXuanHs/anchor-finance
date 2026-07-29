package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type CronURLHandler struct {
	svc *service.CronURLService
	log *logger.Logger
}

func NewCronURLHandler(svc *service.CronURLService, log *logger.Logger) *CronURLHandler {
	return &CronURLHandler{svc: svc, log: log}
}

// GetList returns paginated URL cron tasks.
func (h *CronURLHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	var status *int8
	if s := c.Query("status"); s != "" {
		v, _ := strconv.ParseInt(s, 10, 8)
		st := int8(v)
		status = &st
	}

	tasks, total, err := h.svc.GetTaskList(page, pageSize, keyword, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, tasks, total, page, pageSize)
}

// GetDetail returns a single URL cron task.
func (h *CronURLHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid task id")
		return
	}

	task, err := h.svc.GetTaskByID(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, task)
}

// Create creates a new URL cron task.
func (h *CronURLHandler) Create(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		URL         string `json:"url" binding:"required"`
		Method      string `json:"method" binding:"omitempty,oneof=GET POST"`
		Headers     string `json:"headers"`
		Body        string `json:"body"`
		CronExpr    string `json:"cron_expr" binding:"required"`
		Timeout     int    `json:"timeout"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if req.Method == "" {
		req.Method = "GET"
	}
	if req.Timeout == 0 {
		req.Timeout = 30
	}

	task := &model.CronURLTask{
		Name:        req.Name,
		URL:         req.URL,
		Method:      req.Method,
		Headers:     req.Headers,
		Body:        req.Body,
		CronExpr:    req.CronExpr,
		Status:      1,
		Timeout:     req.Timeout,
		Description: req.Description,
		CreatedBy:   c.GetUint("user_id"),
	}

	if err := h.svc.CreateTask(task); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, task)
}

// Update updates a URL cron task.
func (h *CronURLHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid task id")
		return
	}

	var req struct {
		Name        *string `json:"name"`
		URL         *string `json:"url"`
		Method      *string `json:"method"`
		Headers     *string `json:"headers"`
		Body        *string `json:"body"`
		CronExpr    *string `json:"cron_expr"`
		Timeout     *int    `json:"timeout"`
		Description *string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.URL != nil {
		updates["url"] = *req.URL
	}
	if req.Method != nil {
		updates["method"] = *req.Method
	}
	if req.Headers != nil {
		updates["headers"] = *req.Headers
	}
	if req.Body != nil {
		updates["body"] = *req.Body
	}
	if req.CronExpr != nil {
		updates["cron_expr"] = *req.CronExpr
	}
	if req.Timeout != nil {
		updates["timeout"] = *req.Timeout
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}

	if err := h.svc.UpdateTask(uint(id), updates); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "task updated")
}

// Delete deletes a URL cron task.
func (h *CronURLHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid task id")
		return
	}

	if err := h.svc.DeleteTask(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "task deleted")
}

// SetStatus sets the status of a URL cron task.
func (h *CronURLHandler) SetStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid task id")
		return
	}

	var req struct {
		Status int8 `json:"status" binding:"required,oneof=0 1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.SetStatus(uint(id), req.Status); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "status updated")
}

// RunCron manually runs a URL cron task.
func (h *CronURLHandler) RunCron(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid task id")
		return
	}

	log, err := h.svc.RunTask(uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, log)
}

// GetLogs returns logs for URL cron tasks.
func (h *CronURLHandler) GetLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var taskID uint
	if tid := c.Query("task_id"); tid != "" {
		v, _ := strconv.ParseUint(tid, 10, 64)
		taskID = uint(v)
	}

	logs, total, err := h.svc.GetLogs(taskID, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, logs, total, page, pageSize)
}
