package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UpstreamHandler struct {
	db  *gorm.DB
	log *logger.Logger
	svc *service.UpstreamService
}

func NewUpstreamHandler(db *gorm.DB, log *logger.Logger) *UpstreamHandler {
	return &UpstreamHandler{
		db:  db,
		log: log,
		svc: service.NewUpstreamService(db, log),
	}
}

// GetProviders returns all upstream providers (admin).
func (h *UpstreamHandler) GetProviders(c *gin.Context) {
	var providers []model.UpstreamProvider
	if err := h.db.Order("id DESC").Find(&providers).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, providers)
}

type CreateProviderRequest struct {
	Name   string      `json:"name" binding:"required,max=100"`
	Type   string      `json:"type" binding:"required,max=50"`
	APIURL string      `json:"api_url" binding:"max=500"`
	APIKey string      `json:"api_key" binding:"max=255"`
	Config model.JSON  `json:"config"`
}

// CreateProvider creates a new upstream provider (admin).
func (h *UpstreamHandler) CreateProvider(c *gin.Context) {
	var req CreateProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	provider := model.UpstreamProvider{
		Name:   req.Name,
		Type:   req.Type,
		APIURL: req.APIURL,
		APIKey: req.APIKey,
		Config: req.Config,
	}

	if err := h.db.Create(&provider).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, provider)
}

type UpdateProviderRequest struct {
	Name      *string    `json:"name"`
	Type      *string    `json:"type"`
	APIURL    *string    `json:"api_url"`
	APIKey    *string    `json:"api_key"`
	Config    *model.JSON `json:"config"`
	IsActive  *bool      `json:"is_active"`
}

// UpdateProvider updates an upstream provider (admin).
func (h *UpstreamHandler) UpdateProvider(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid provider id")
		return
	}

	var req UpdateProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.APIURL != nil {
		updates["api_url"] = *req.APIURL
	}
	if req.APIKey != nil {
		updates["api_key"] = *req.APIKey
	}
	if req.Config != nil {
		updates["config"] = *req.Config
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(updates) > 0 {
		if err := h.db.Model(&model.UpstreamProvider{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			response.ServerError(c, err.Error())
			return
		}
	}
	response.SuccessMsg(c, "provider updated")
}

// DeleteProvider deletes an upstream provider (admin).
func (h *UpstreamHandler) DeleteProvider(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid provider id")
		return
	}

	if err := h.db.Delete(&model.UpstreamProvider{}, id).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "provider deleted")
}

// TestConnection tests the connection to an upstream provider (admin).
func (h *UpstreamHandler) TestConnection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid provider id")
		return
	}

	var provider model.UpstreamProvider
	if err := h.db.First(&provider, id).Error; err != nil {
		response.NotFound(c, "provider not found")
		return
	}

	result, err := h.svc.TestConnection(uint(id))
	if err != nil {
		msg := err.Error()
		var latency int64
		if result != nil {
			latency = result.Latency
		}
		response.Success(c, gin.H{
			"provider_id": id,
			"status":      "failed",
			"message":     msg,
			"latency_ms":  latency,
		})
		return
	}

	response.Success(c, gin.H{
		"provider_id": id,
		"status":      "ok",
		"message":     result.Message,
		"latency_ms":  result.Latency,
	})
}

// SyncProducts triggers product sync from an upstream provider (admin).
func (h *UpstreamHandler) SyncProducts(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid provider id")
		return
	}

	var provider model.UpstreamProvider
	if err := h.db.First(&provider, id).Error; err != nil {
		response.NotFound(c, "provider not found")
		return
	}

	synced, err := h.svc.SyncProducts(uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	h.log.WithField("provider_id", id).Infof("upstream product sync completed: %d products", synced)

	response.Success(c, gin.H{
		"provider_id": id,
		"synced":      synced,
		"message":     "sync completed",
	})
}

// GetSyncLogs returns sync logs for an upstream provider (admin).
func (h *UpstreamHandler) GetSyncLogs(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid provider id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var logs []model.UpstreamSyncLog
	var total int64

	query := h.db.Model(&model.UpstreamSyncLog{}).Where("upstream_id = ?", id)
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, logs, total, page, pageSize)
}
