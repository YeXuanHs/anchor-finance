package handler

import (
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type ConfigGeneralHandler struct {
	svc *service.ConfigGeneralService
	log *logger.Logger
}

func NewConfigGeneralHandler(svc *service.ConfigGeneralService, log *logger.Logger) *ConfigGeneralHandler {
	return &ConfigGeneralHandler{svc: svc, log: log}
}

// Get is an alias for GetConfig.
func (h *ConfigGeneralHandler) Get(c *gin.Context) { h.GetConfig(c) }

// Update is an alias for UpdateConfig.
func (h *ConfigGeneralHandler) Update(c *gin.Context) { h.UpdateConfig(c) }

// GetConfig returns the general site configuration.
func (h *ConfigGeneralHandler) GetConfig(c *gin.Context) {
	cfg, err := h.svc.GetConfig()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cfg)
}

// UpdateConfig updates the general site configuration.
func (h *ConfigGeneralHandler) UpdateConfig(c *gin.Context) {
	var req service.GeneralConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.UpdateConfig(req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "general config updated")
}
