package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type SendMessageHandler struct {
	msgSvc *service.SendMessageService
	log    *logger.Logger
}

func NewSendMessageHandler(msgSvc *service.SendMessageService, log *logger.Logger) *SendMessageHandler {
	return &SendMessageHandler{msgSvc: msgSvc, log: log}
}

// SendEmail sends an email (admin).
func (h *SendMessageHandler) SendEmail(c *gin.Context) {
	var req service.SendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	msg, err := h.msgSvc.SendEmail(req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, msg)
}

// SendSMS sends an SMS (admin).
func (h *SendMessageHandler) SendSMS(c *gin.Context) {
	var req service.SendSMSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	msg, err := h.msgSvc.SendSMS(req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, msg)
}

// SendSiteMessage sends a site message (admin).
func (h *SendMessageHandler) SendSiteMessage(c *gin.Context) {
	var req service.SendSiteMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	msgs, err := h.msgSvc.SendSiteMessage(req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, msgs)
}

// BatchSend sends messages to multiple recipients (admin).
func (h *SendMessageHandler) BatchSend(c *gin.Context) {
	var req service.BatchSendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	batchID, msgs, err := h.msgSvc.BatchSend(req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{
		"batch_id": batchID,
		"count":    len(msgs),
		"messages": msgs,
	})
}

// GetDetail returns a single message (admin).
func (h *SendMessageHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid message id")
		return
	}

	msg, err := h.msgSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "message not found")
		return
	}
	response.Success(c, msg)
}

// GetList returns all messages (admin).
func (h *SendMessageHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var msgType *string
	if t := c.Query("type"); t != "" {
		msgType = &t
	}
	var status *int
	if s := c.Query("status"); s != "" {
		v, _ := strconv.Atoi(s)
		status = &v
	}

	messages, total, err := h.msgSvc.GetList(page, pageSize, msgType, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, messages, total, page, pageSize)
}

// GetByBatchID returns all messages in a batch (admin).
func (h *SendMessageHandler) GetByBatchID(c *gin.Context) {
	batchID := c.Param("batch_id")
	if batchID == "" {
		response.BadRequest(c, "batch_id is required")
		return
	}

	messages, err := h.msgSvc.GetByBatchID(batchID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, messages)
}

// GetUserMessages returns messages for the authenticated user.
func (h *SendMessageHandler) GetUserMessages(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	messages, total, err := h.msgSvc.GetUserMessages(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, messages, total, page, pageSize)
}

// RetryFailed requeues failed messages (admin).
func (h *SendMessageHandler) RetryFailed(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	if err := h.msgSvc.RetryFailed(limit); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "failed messages requeued")
}
