package handler

import (
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type ConfigMessageHandler struct {
	svc *service.ConfigMessageService
	log *logger.Logger
}

func NewConfigMessageHandler(svc *service.ConfigMessageService, log *logger.Logger) *ConfigMessageHandler {
	return &ConfigMessageHandler{svc: svc, log: log}
}

// GetAll returns all message channel configs.
func (h *ConfigMessageHandler) GetAll(c *gin.Context) {
	items, err := h.svc.GetAll()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, items)
}

// GetByChannel returns a single message channel config.
func (h *ConfigMessageHandler) GetByChannel(c *gin.Context) {
	channel := c.Param("channel")
	if channel == "" {
		response.BadRequest(c, "channel is required")
		return
	}

	item, err := h.svc.GetByChannel(channel)
	if err != nil {
		response.NotFound(c, "message channel config not found")
		return
	}
	response.Success(c, item)
}

// Update updates a message channel config.
func (h *ConfigMessageHandler) Update(c *gin.Context) {
	channel := c.Param("channel")
	if channel == "" {
		response.BadRequest(c, "channel is required")
		return
	}

	var req service.UpdateMessageConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	item, err := h.svc.Update(channel, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, item)
}

// TestSend tests sending via a message channel.
func (h *ConfigMessageHandler) TestSend(c *gin.Context) {
	channel := c.Param("channel")
	if channel == "" {
		response.BadRequest(c, "channel is required")
		return
	}

	if err := h.svc.TestSend(channel); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "test message sent successfully")
}

// GetEnabled returns all enabled message channels.
func (h *ConfigMessageHandler) GetEnabled(c *gin.Context) {
	items, err := h.svc.GetEnabledChannels()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, items)
}
