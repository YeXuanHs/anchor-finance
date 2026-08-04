package handler

import (
	"strconv"

	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserServicesHandler handles user-facing service management requests.
type UserServicesHandler struct {
	db *gorm.DB
}

// NewUserServicesHandler creates a new UserServicesHandler.
func NewUserServicesHandler(db *gorm.DB) *UserServicesHandler {
	return &UserServicesHandler{db: db}
}

// GetSMSServices returns the user's SMS services.
// GET /user/services/sms
func (h *UserServicesHandler) GetSMSServices(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var items []map[string]interface{}
	var total int64

	query := h.db.Table("client_services").
		Where("user_id = ? AND type = ?", userID, "sms")
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&items)

	response.SuccessPage(c, items, total, page, pageSize)
}

// GetSMSRecords returns the user's SMS sending records.
// GET /user/services/sms/records
func (h *UserServicesHandler) GetSMSRecords(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var items []map[string]interface{}
	var total int64

	query := h.db.Table("sms_logs").
		Where("user_id = ?", userID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&items)

	response.SuccessPage(c, items, total, page, pageSize)
}

// GetSoftwareServices returns the user's software services.
// GET /user/services/software
func (h *UserServicesHandler) GetSoftwareServices(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var items []map[string]interface{}
	var total int64

	query := h.db.Table("client_services").
		Where("user_id = ? AND type = ?", userID, "software")
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&items)

	response.SuccessPage(c, items, total, page, pageSize)
}

// PostSoftwareResetKey resets the API key for a software service.
// POST /user/services/software/reset-key
func (h *UserServicesHandler) PostSoftwareResetKey(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		ServiceID uint `json:"service_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Verify service belongs to user
	var service map[string]interface{}
	if err := h.db.Table("client_services").
		Where("id = ? AND user_id = ? AND type = ?", req.ServiceID, userID, "software").
		First(&service).Error; err != nil {
		response.NotFound(c, "service not found")
		return
	}

	// Generate new API key
	newKey := "sk-" + strconv.FormatUint(uint64(req.ServiceID), 10) + "-" + strconv.FormatInt(int64(userID), 10)

	h.db.Table("client_services").Where("id = ?", req.ServiceID).Update("api_key", newKey)

	response.Success(c, gin.H{
		"service_id": req.ServiceID,
		"api_key":    newKey,
	})
}
