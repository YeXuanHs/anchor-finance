package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type CronHandler struct {
	cronSvc *service.CronService
	log     *logger.Logger
}

func NewCronHandler(cronSvc *service.CronService, log *logger.Logger) *CronHandler {
	return &CronHandler{cronSvc: cronSvc, log: log}
}

// ==================== 请求结构体 ====================

type createCronTaskRequest struct {
	Name        string `json:"name" binding:"required,max=128"`
	Type        string `json:"type" binding:"required,oneof=custom system plugin"`
	CronExpr    string `json:"cron_expr" binding:"required,max=64"`
	Command     string `json:"command" binding:"omitempty"`
	Params      string `json:"params" binding:"omitempty"`
	Timeout     int    `json:"timeout" binding:"omitempty,gte=1,lte=3600"`
	Priority    int    `json:"priority" binding:"omitempty,gte=0,lte=99"`
	Description string `json:"description" binding:"omitempty"`
}

type updateCronTaskRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=128"`
	Type        *string `json:"type" binding:"omitempty,oneof=custom system plugin"`
	CronExpr    *string `json:"cron_expr" binding:"omitempty,max=64"`
	Command     *string `json:"command"`
	Params      *string `json:"params"`
	Timeout     *int    `json:"timeout" binding:"omitempty,gte=1,lte=3600"`
	Priority    *int    `json:"priority" binding:"omitempty,gte=0,lte=99"`
	Description *string `json:"description"`
}

type setCronStatusRequest struct {
	Status int8 `json:"status" binding:"required,oneof=0 1"`
}

type runCronTaskRequest struct {
	Comment string `json:"comment" binding:"omitempty"`
}

// ==================== 任务管理 ====================

// GetList 获取定时任务列表
func (h *CronHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	var status *int8
	if s := c.Query("status"); s != "" {
		v, err := strconv.ParseInt(s, 10, 8)
		if err == nil {
			st := int8(v)
			status = &st
		}
	}

	tasks, total, err := h.cronSvc.GetTaskList(page, pageSize, keyword, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, tasks, total, page, pageSize)
}

// GetDetail 获取任务详情
func (h *CronHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid task id")
		return
	}

	task, err := h.cronSvc.GetTaskByID(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, task)
}

// Create 创建任务
func (h *CronHandler) Create(c *gin.Context) {
	var req createCronTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	task := &model.CronTask{
		Name:        req.Name,
		Type:        req.Type,
		CronExpr:    req.CronExpr,
		Command:     req.Command,
		Params:      req.Params,
		Status:      1,
		Timeout:     req.Timeout,
		Priority:    req.Priority,
		Description: req.Description,
		CreatedBy:   c.GetUint("user_id"),
	}
	if task.Timeout == 0 {
		task.Timeout = 60
	}

	if err := h.cronSvc.CreateTask(task); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, task)
}

// Update 编辑任务
func (h *CronHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid task id")
		return
	}

	var req updateCronTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.CronExpr != nil {
		updates["cron_expr"] = *req.CronExpr
	}
	if req.Command != nil {
		updates["command"] = *req.Command
	}
	if req.Params != nil {
		updates["params"] = *req.Params
	}
	if req.Timeout != nil {
		updates["timeout"] = *req.Timeout
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}

	if len(updates) == 0 {
		response.BadRequest(c, "no fields to update")
		return
	}

	if err := h.cronSvc.UpdateTask(uint(id), updates); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "task updated")
}

// Delete 删除任务
func (h *CronHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid task id")
		return
	}

	if err := h.cronSvc.DeleteTask(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "task deleted")
}

// SetStatus 启用/禁用任务
func (h *CronHandler) SetStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid task id")
		return
	}

	var req setCronStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.cronSvc.SetStatus(uint(id), req.Status); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if req.Status == 1 {
		response.SuccessMsg(c, "task enabled")
	} else {
		response.SuccessMsg(c, "task disabled")
	}
}

// RunTask 手动执行任务
func (h *CronHandler) RunTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid task id")
		return
	}

	operatorID := c.GetUint("user_id")
	taskLog, err := h.cronSvc.RunTask(uint(id), operatorID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, taskLog)
}

// ==================== 执行日志 ====================

// GetLogs 获取执行日志列表
func (h *CronHandler) GetLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var taskID uint
	if tid := c.Query("task_id"); tid != "" {
		v, err := strconv.ParseUint(tid, 10, 64)
		if err == nil {
			taskID = uint(v)
		}
	}

	logs, total, err := h.cronSvc.GetTaskLogs(taskID, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, logs, total, page, pageSize)
}

// GetLogDetail 获取日志详情
func (h *CronHandler) GetLogDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid log id")
		return
	}

	taskLog, err := h.cronSvc.GetLogByID(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, taskLog)
}
