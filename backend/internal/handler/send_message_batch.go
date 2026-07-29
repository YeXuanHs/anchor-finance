package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type SendMessageBatchHandler struct {
	svc *service.SendMessageBatchService
	log *logger.Logger
}

func NewSendMessageBatchHandler(svc *service.SendMessageBatchService, log *logger.Logger) *SendMessageBatchHandler {
	return &SendMessageBatchHandler{svc: svc, log: log}
}

// GetSearchParams returns search parameters for batch messaging.
func (h *SendMessageBatchHandler) GetSearchParams(c *gin.Context) {
	params, err := h.svc.GetSearchParams()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, params)
}

// GetBatches returns paginated batch send tasks.
func (h *SendMessageBatchHandler) GetBatches(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	sendMethod := c.Query("send_method")

	var status *int8
	if s := c.Query("status"); s != "" {
		v, _ := strconv.ParseInt(s, 10, 8)
		st := int8(v)
		status = &st
	}

	batches, total, err := h.svc.GetBatches(page, pageSize, sendMethod, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, batches, total, page, pageSize)
}

// SendBatch sends messages in batch.
func (h *SendMessageBatchHandler) SendBatch(c *gin.Context) {
	var req struct {
		SendMethod string `json:"send_method" binding:"required"`
		Subject    string `json:"subject"`
		Content    string `json:"content" binding:"required"`
		UserIDs    []uint `json:"user_ids"`
		GroupIDs   []uint `json:"group_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	batch := &model.SendMessageBatch{
		SendMethod: req.SendMethod,
		Subject:    req.Subject,
		Content:    req.Content,
		Status:     0,
		CreatedBy:  c.GetUint("user_id"),
	}

	if err := h.svc.SendBatch(batch); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, batch)
}

// GetProgress returns the progress of a batch send operation.
func (h *SendMessageBatchHandler) GetProgress(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid batch id")
		return
	}

	batch, err := h.svc.GetProgress(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, batch)
}

// GetRecords returns batch send records.
func (h *SendMessageBatchHandler) GetRecords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var batchID *uint
	if bid := c.Query("batch_id"); bid != "" {
		v, _ := strconv.ParseUint(bid, 10, 64)
		id := uint(v)
		batchID = &id
	}

	records, total, err := h.svc.GetRecords(page, pageSize, batchID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, records, total, page, pageSize)
}
