package handler

import (
	"net/http"
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreditHandler handles credit limit HTTP requests.
type CreditHandler struct {
	db      *gorm.DB
	creditSvc *service.CreditService
}

// NewCreditHandler creates a new CreditHandler.
func NewCreditHandler(db *gorm.DB, creditSvc *service.CreditService) *CreditHandler {
	return &CreditHandler{db: db, creditSvc: creditSvc}
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

// applyRequest is the payload for Apply.
type applyRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
	Reason string  `json:"reason"`
}

// Apply creates a credit limit application record.
func (h *CreditHandler) Apply(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Unauthorized(c, "login required")
		return
	}
	var req applyRequest
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

	// Check for pending application
	var pending int64
	h.db.Model(&model.CreditLog{}).
		Where("user_id = ? AND type = ? AND status = ?", userID, "apply", "pending").
		Count(&pending)
	if pending > 0 {
		response.BadRequest(c, "you already have a pending credit application")
		return
	}

	// Create credit application log
	log := model.CreditLog{
		UserID:      userID,
		Type:        "apply",
		Amount:      req.Amount,
		Balance:     0,
		Status:      "pending",
		Remark:      req.Reason,
		RelatedType: "credit_apply",
	}
	if err := h.db.Create(&log).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"id":     log.ID,
		"amount": req.Amount,
		"reason": req.Reason,
		"status": "pending",
	})
}

// repayRequest is the payload for Repay.
type repayRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

// Repay repays used credit balance, reducing the used amount.
func (h *CreditHandler) Repay(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Unauthorized(c, "login required")
		return
	}
	var req repayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Get current credit limit
	var credit model.CreditLimit
	if err := h.db.Where("user_id = ?", userID).First(&credit).Error; err != nil {
		response.BadRequest(c, "no credit limit found")
		return
	}

	if req.Amount > credit.Used {
		response.BadRequest(c, "repay amount exceeds used credit")
		return
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		// Reduce used credit and increase available
		newUsed := credit.Used - req.Amount
		newAvailable := credit.Available + req.Amount

		if err := tx.Model(&credit).Updates(map[string]interface{}{
			"used":      newUsed,
			"available": newAvailable,
		}).Error; err != nil {
			return err
		}

		// Log the repayment
		creditLog := model.CreditLog{
			UserID:      userID,
			Type:        "repay",
			Amount:      req.Amount,
			Balance:     newAvailable,
			Remark:      "credit repayment",
			RelatedType: "credit_repay",
		}
		return tx.Create(&creditLog).Error
	})

	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	h.db.Where("user_id = ?", userID).First(&credit)
	response.Success(c, credit)
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

// ─────────────────────── Billing Cycle Handlers ───────────────────────

// GetBills returns the authenticated user's credit bills.
// GET /credit/bills
func (h *CreditHandler) GetBills(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	bills, total, err := h.creditSvc.GetBills(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessPage(c, bills, total, page, pageSize)
}

// GetBillDetail returns a single bill detail.
// GET /credit/bills/:id
func (h *CreditHandler) GetBillDetail(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	billID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid bill id")
		return
	}

	bill, err := h.creditSvc.GetBillByID(uint(billID))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	// Ensure the bill belongs to the user
	if bill.UserID != userID {
		response.Forbidden(c, "access denied")
		return
	}

	// Calculate current late fee if overdue
	h.creditSvc.CalculateLateFee(bill.ID)
	h.db.First(bill, bill.ID)

	response.Success(c, bill)
}

// payBillRequest is the payload for PayBill.
type payBillRequest struct {
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	PaymentMethod string  `json:"payment_method" binding:"required"`
}

// PayBill pays a credit bill.
// POST /credit/bills/:id/pay
func (h *CreditHandler) PayBill(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	billID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid bill id")
		return
	}

	var req payBillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Verify bill belongs to user
	bill, err := h.creditSvc.GetBillByID(uint(billID))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	if bill.UserID != userID {
		response.Forbidden(c, "access denied")
		return
	}

	bill, err = h.creditSvc.PayBill(uint(billID), req.Amount, req.PaymentMethod)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, bill)
}

// GetCreditConfig returns the user's credit billing configuration.
// GET /credit/config
func (h *CreditHandler) GetCreditConfig(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	config, err := h.creditSvc.GetCreditConfig(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"bill_generation_day": config.BillGenerationDay,
		"repayment_period":    config.RepaymentPeriod,
		"prepayment_balance":  config.PrepaymentBalance,
	})
}

// updateCreditConfigRequest is the payload for UpdateCreditConfig.
type updateCreditConfigRequest struct {
	BillGenerationDay int `json:"bill_generation_day" binding:"required,min=1,max=28"`
	RepaymentPeriod   int `json:"repayment_period" binding:"required,min=1,max=60"`
}

// UpdateCreditConfig updates the user's credit billing configuration.
// PUT /credit/config
func (h *CreditHandler) UpdateCreditConfig(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	var req updateCreditConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.creditSvc.UpdateCreditConfig(userID, req.BillGenerationDay, req.RepaymentPeriod); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMsg(c, "credit config updated")
}

// ─────────────────────── Admin Billing Handlers ───────────────────────

// AdminGenerateBills manually triggers bill generation.
// POST /admin/credit/bills/generate
func (h *CreditHandler) AdminGenerateBills(c *gin.Context) {
	output, err := h.creditSvc.GenerateMonthlyBills()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": output})
}

// AdminGetBills returns all credit bills (admin).
// GET /admin/credit/bills
func (h *CreditHandler) AdminGetBills(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	userIDStr := c.Query("user_id")
	status := c.Query("status")

	var userID uint
	if userIDStr != "" {
		uid, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			response.BadRequest(c, "invalid user_id")
			return
		}
		userID = uint(uid)
	}

	bills, total, err := h.creditSvc.GetAdminBills(page, pageSize, userID, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessPage(c, bills, total, page, pageSize)
}

// waiveLateFeeRequest is the payload for AdminWaiveLateFee.
type waiveLateFeeRequest struct {
	Remark string `json:"remark"`
}

// AdminWaiveLateFee waives the late fee for a bill.
// POST /admin/credit/bills/:id/waive-fee
func (h *CreditHandler) AdminWaiveLateFee(c *gin.Context) {
	adminID := getUserID(c)
	if adminID == 0 {
		return
	}

	billID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid bill id")
		return
	}

	var req waiveLateFeeRequest
	c.ShouldBindJSON(&req)

	bill, err := h.creditSvc.WaiveLateFee(uint(billID), adminID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, bill)
}
