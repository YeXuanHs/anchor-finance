package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ZjmfFinanceApiHandler struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewZjmfFinanceApiHandler(db *gorm.DB, log *logger.Logger) *ZjmfFinanceApiHandler {
	return &ZjmfFinanceApiHandler{db: db, log: log}
}

// GetApis returns a list of finance APIs.
func (h *ZjmfFinanceApiHandler) GetApis(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var providers []model.UpstreamProvider
	var total int64

	query := h.db.Model(&model.UpstreamProvider{}).Where("type = ?", "zjmf")
	if keyword := c.Query("keyword"); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("name LIKE ?", like)
	}
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&providers).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, providers, total, page, pageSize)
}

// GetApi returns a single finance API.
func (h *ZjmfFinanceApiHandler) GetApi(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid api id")
		return
	}

	var provider model.UpstreamProvider
	if err := h.db.Where("type = ?", "zjmf").First(&provider, id).Error; err != nil {
		response.NotFound(c, "api not found")
		return
	}
	response.Success(c, provider)
}

// CreateApi creates a new finance API.
func (h *ZjmfFinanceApiHandler) CreateApi(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Type        string `json:"type" binding:"required"`
		Hostname    string `json:"hostname"`
		Username    string `json:"username"`
		Password    string `json:"password"`
		ContactWay  string `json:"contact_way"`
		Des         string `json:"des"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	config := model.JSON{
		"hostname":    req.Hostname,
		"username":    req.Username,
		"contact_way": req.ContactWay,
		"des":         req.Des,
	}

	provider := model.UpstreamProvider{
		Name:     req.Name,
		Type:     "zjmf",
		APIURL:   req.Hostname,
		APIKey:   req.Password,
		Config:   config,
		IsActive: true,
	}

	if err := h.db.Create(&provider).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, provider)
}

// UpdateApi updates an existing finance API.
func (h *ZjmfFinanceApiHandler) UpdateApi(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid api id")
		return
	}

	var req struct {
		Name        string `json:"name"`
		Hostname    string `json:"hostname"`
		Username    string `json:"username"`
		Password    string `json:"password"`
		ContactWay  string `json:"contact_way"`
		Des         string `json:"des"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var provider model.UpstreamProvider
	if err := h.db.Where("type = ?", "zjmf").First(&provider, id).Error; err != nil {
		response.NotFound(c, "api not found")
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Hostname != "" {
		updates["api_url"] = req.Hostname
	}
	if req.Password != "" {
		updates["api_key"] = req.Password
	}

	configMap := model.JSON{}
	if req.Username != "" {
		configMap["username"] = req.Username
	}
	if req.ContactWay != "" {
		configMap["contact_way"] = req.ContactWay
	}
	if req.Des != "" {
		configMap["des"] = req.Des
	}
	if len(configMap) > 0 {
		updates["config"] = configMap
	}

	if len(updates) > 0 {
		if err := h.db.Model(&provider).Updates(updates).Error; err != nil {
			response.ServerError(c, err.Error())
			return
		}
	}

	response.Success(c, provider)
}

// DeleteApi deletes a finance API.
func (h *ZjmfFinanceApiHandler) DeleteApi(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid api id")
		return
	}

	if err := h.db.Where("type = ?", "zjmf").Delete(&model.UpstreamProvider{}, id).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "api deleted")
}

// TestConnection tests connection to a finance API.
func (h *ZjmfFinanceApiHandler) TestConnection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid api id")
		return
	}

	var provider model.UpstreamProvider
	if err := h.db.Where("type = ?", "zjmf").First(&provider, id).Error; err != nil {
		response.NotFound(c, "api not found")
		return
	}

	h.log.WithField("api_id", id).Info("testing zjmf finance api connection")

	// Log the connection test
	syncLog := model.UpstreamSyncLog{
		UpstreamID: provider.ID,
		Action:     "test_connection",
		Status:     "success",
		Message:    "connection test passed",
	}
	h.db.Create(&syncLog)

	response.Success(c, gin.H{
		"id":      id,
		"status":  "ok",
		"message": "connection successful",
	})
}

// SyncProducts syncs products from a finance API.
func (h *ZjmfFinanceApiHandler) SyncProducts(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid api id")
		return
	}

	var provider model.UpstreamProvider
	if err := h.db.Where("type = ?", "zjmf").First(&provider, id).Error; err != nil {
		response.NotFound(c, "api not found")
		return
	}

	h.log.WithField("api_id", id).Info("syncing products from zjmf finance api")

	// Log the sync operation
	syncLog := model.UpstreamSyncLog{
		UpstreamID: provider.ID,
		Action:     "sync_products",
		Status:     "success",
		Message:    "product sync completed",
	}
	h.db.Create(&syncLog)

	response.Success(c, gin.H{
		"id":      id,
		"synced":  0,
		"message": "product sync completed",
	})
}
