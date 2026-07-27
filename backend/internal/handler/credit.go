package handler

import (
	"net/http"
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreditHandler handles credit limit HTTP requests.
type CreditHandler struct {
	db *gorm.DB
}

// NewCreditHandler creates a new CreditHandler.
func NewCreditHandler(db *gorm.DB) *CreditHandler {
	return &CreditHandler{db: db}
}

// GetInfo returns the authenticated user's credit limit info.
// GET /user/credit
func (h *CreditHandler) GetInfo(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	var credit model.CreditLimit
	if err := h.db.Where("user_id = ?", userID).First(&credit).Error; err != nil {
		// Return zeroed credit if not yet created
		response.Success(c, gin.H{
			"limit":     0,
			"used":      0,
			"available": 0,
		})
		return
	}

	response.Success(c, credit)
}

// GetLogs returns the authenticated user's credit change logs.
// GET /user/credit/logs
func (h *CreditHandler) GetLogs(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var logs []model.CreditLog
	var total int64

	query := h.db.Model(&model.CreditLog{}).Where("user_id = ?", userID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs)

	response.SuccessPage(c, logs, total, page, pageSize)
}

// adminAdjustCreditRequest is the payload for AdminAdjust.
type adminAdjustCreditRequest struct {
	Limit       float64 `json:"limit" binding:"required,gte=0"`
	Description string  `json:"description"`
}

// AdminAdjust adjusts a user's credit limit (admin).
// POST /admin/users/:id/credit
func (h *CreditHandler) AdminAdjust(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req adminAdjustCreditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Verify user exists
	var user model.User
	if err := h.db.First(&user, userID).Error; err != nil {
		response.NotFound(c, "user not found")
		return
	}

	adminID := getUserID(c)
	if adminID == 0 {
		return
	}

	var credit model.CreditLimit
	err = h.db.Transaction(func(tx *gorm.DB) error {
		// Find or create credit limit
		if err := tx.Where("user_id = ?", userID).FirstOrCreate(&credit, model.CreditLimit{
			UserID: uint(userID),
		}).Error; err != nil {
			return err
		}

		oldLimit := credit.Limit
		newAvailable := credit.Available + (req.Limit - oldLimit)

		updates := map[string]interface{}{
			"limit":     req.Limit,
			"available": newAvailable,
		}
		if req.Description != "" {
			updates["description"] = req.Description
		}

		if err := tx.Model(&credit).Updates(updates).Error; err != nil {
			return err
		}

		// Log the adjustment
		log := model.CreditLog{
			UserID:      uint(userID),
			Type:        "adjust",
			Amount:      req.Limit - oldLimit,
			Balance:     newAvailable,
			AdminID:     &adminID,
			Remark:      req.Description,
			RelatedType: "admin_adjust",
		}
		return tx.Create(&log).Error
	})

	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	// Reload with updated values
	h.db.Where("user_id = ?", userID).First(&credit)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "credit limit adjusted",
		"data":    credit,
	})
}

// AdminGetLogs returns all credit logs with optional filters (admin).
// GET /admin/credit/logs
func (h *CreditHandler) AdminGetLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	userID := c.Query("user_id")
	logType := c.Query("type")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var logs []model.CreditLog
	var total int64

	query := h.db.Model(&model.CreditLog{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if logType != "" {
		query = query.Where("type = ?", logType)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs)

	response.SuccessPage(c, logs, total, page, pageSize)
}
