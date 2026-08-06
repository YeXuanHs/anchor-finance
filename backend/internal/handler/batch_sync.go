package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type BatchSyncHandler struct {
	svc *service.BatchSyncService
	log *logger.Logger
}

func NewBatchSyncHandler(svc *service.BatchSyncService, log *logger.Logger) *BatchSyncHandler {
	return &BatchSyncHandler{svc: svc, log: log}
}

// GetList returns paginated batch sync tasks.
func (h *BatchSyncHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	taskType := c.Query("type")

	var status *int8
	if s := c.Query("status"); s != "" {
		v, _ := strconv.ParseInt(s, 10, 8)
		st := int8(v)
		status = &st
	}

	tasks, total, err := h.svc.GetTaskList(page, pageSize, taskType, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, tasks, total, page, pageSize)
}

// GetDetail returns a single batch sync task.
func (h *BatchSyncHandler) GetDetail(c *gin.Context) {
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

// Execute runs a batch synchronization task.
func (h *BatchSyncHandler) Execute(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid task id")
		return
	}

	operatorID := c.GetUint("user_id")
	if err := h.svc.Execute(uint(id), operatorID); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "batch sync started")
}

// GetLogs returns logs for a batch sync task.
func (h *BatchSyncHandler) GetLogs(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid task id")
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
