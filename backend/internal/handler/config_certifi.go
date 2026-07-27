package handler

import (
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type ConfigCertifiHandler struct {
	svc *service.ConfigCertifiService
	log *logger.Logger
}

func NewConfigCertifiHandler(svc *service.ConfigCertifiService, log *logger.Logger) *ConfigCertifiHandler {
	return &ConfigCertifiHandler{svc: svc, log: log}
}

// GetConfig returns the certification configuration.
func (h *ConfigCertifiHandler) GetConfig(c *gin.Context) {
	cfg, err := h.svc.GetConfig()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cfg)
}

// UpdateConfig updates the certification configuration.
func (h *ConfigCertifiHandler) UpdateConfig(c *gin.Context) {
	var req service.CertificationConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.UpdateConfig(req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "certification config updated")
}
