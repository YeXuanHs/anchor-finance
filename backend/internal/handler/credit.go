package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
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

// ─────────────────────── Admin Client List ───────────────────────

// AdminGetClientList returns all credit-enabled users with search/filter/sort.
// GET /admin/credit/clients
func (h *CreditHandler) AdminGetClientList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	paymentStatus := c.Query("payment_status")
	order := c.DefaultQuery("order", "id")
	sort := c.DefaultQuery("sort", "desc")

	items, total, err := h.creditSvc.GetClientList(page, pageSize, keyword, paymentStatus, order, sort)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessPage(c, items, total, page, pageSize)
}

// ─────────────────────── Admin Enable/Disable Credit ───────────────────────

// adminEnableCreditRequest is the payload for AdminEnableCredit.
type adminEnableCreditRequest struct {
	UserID           uint    `json:"user_id" binding:"required"`
	Limit            float64 `json:"limit" binding:"required,gte=0"`
	BillGenerationDay int    `json:"bill_generation_day" binding:"omitempty,min=1,max=28"`
	RepaymentPeriod  int     `json:"repayment_period" binding:"omitempty,min=1,max=60"`
}

// AdminEnableCredit enables credit for a user.
// POST /admin/credit/users/enable
func (h *CreditHandler) AdminEnableCredit(c *gin.Context) {
	var req adminEnableCreditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Verify user exists
	var user model.User
	if err := h.db.First(&user, req.UserID).Error; err != nil {
		response.NotFound(c, "user not found")
		return
	}

	billDay := req.BillGenerationDay
	if billDay == 0 {
		billDay = 1
	}
	repayPeriod := req.RepaymentPeriod
	if repayPeriod == 0 {
		repayPeriod = 15
	}

	if err := h.creditSvc.EnableCredit(req.UserID, req.Limit, billDay, repayPeriod); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	// Log the action
	adminID := getUserID(c)
	credit, _ := h.creditSvc.GetByUserID(req.UserID)
	if credit != nil {
		log := model.CreditLog{
			UserID:      req.UserID,
			Type:        "adjust",
			Amount:      req.Limit,
			Balance:     credit.Available,
			AdminID:     &adminID,
			Remark:      fmt.Sprintf("Credit enabled with limit %.2f", req.Limit),
			RelatedType: "admin_enable",
		}
		h.db.Create(&log)
	}

	response.SuccessMsg(c, "credit enabled")
}

// AdminDisableCredit disables credit for a user.
// POST /admin/credit/users/:id/disable
func (h *CreditHandler) AdminDisableCredit(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	// Verify user exists
	var user model.User
	if err := h.db.First(&user, userID).Error; err != nil {
		response.NotFound(c, "user not found")
		return
	}

	if err := h.creditSvc.DisableCredit(uint(userID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Log the action
	adminID := getUserID(c)
	log := model.CreditLog{
		UserID:      uint(userID),
		Type:        "adjust",
		Amount:      0,
		Balance:     0,
		AdminID:     &adminID,
		Remark:      "Credit disabled",
		RelatedType: "admin_disable",
	}
	h.db.Create(&log)

	response.SuccessMsg(c, "credit disabled")
}

// adminUpdateCreditSettingsRequest is the payload for AdminUpdateCreditSettings.
type adminUpdateCreditSettingsRequest struct {
	Limit            *float64 `json:"limit" binding:"omitempty,gte=0"`
	BillGenerationDay *int    `json:"bill_generation_day" binding:"omitempty,min=1,max=28"`
	RepaymentPeriod  *int     `json:"repayment_period" binding:"omitempty,min=1,max=60"`
}

// AdminUpdateCreditSettings updates credit settings for a user.
// PUT /admin/credit/users/:id/settings
func (h *CreditHandler) AdminUpdateCreditSettings(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req adminUpdateCreditSettingsRequest
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

	if err := h.creditSvc.UpdateUserCreditSettings(uint(userID), req.Limit, req.BillGenerationDay, req.RepaymentPeriod); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Log the action
	adminID := getUserID(c)
	log := model.CreditLog{
		UserID:      uint(userID),
		Type:        "adjust",
		Amount:      0,
		Balance:     0,
		AdminID:     &adminID,
		Remark:      "Credit settings updated",
		RelatedType: "admin_update_settings",
	}
	h.db.Create(&log)

	response.SuccessMsg(c, "credit settings updated")
}

// ─────────────────────── Admin User Credit Invoices ───────────────────────

// AdminGetUserCreditInvoices returns credit invoices for a specific user.
// GET /admin/credit/users/:id/invoices
func (h *CreditHandler) AdminGetUserCreditInvoices(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	items, total, err := h.creditSvc.GetUserCreditInvoices(uint(userID), page, pageSize, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessPage(c, items, total, page, pageSize)
}

// AdminGetCreditInvoiceSubItems returns sub-items under a credit invoice.
// GET /admin/credit/invoices/:id/items
func (h *CreditHandler) AdminGetCreditInvoiceSubItems(c *gin.Context) {
	billID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid bill id")
		return
	}

	items, err := h.creditSvc.GetCreditInvoiceSubItems(uint(billID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, items)
}

// ─────────────────────── Admin Global Credit Config ───────────────────────

// AdminGetGlobalCreditConfig returns the global credit limit configuration.
// GET /admin/credit/config
func (h *CreditHandler) AdminGetGlobalCreditConfig(c *gin.Context) {
	config, err := h.creditSvc.GetGlobalCreditConfig()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, config)
}

// AdminUpdateGlobalCreditConfig updates the global credit limit configuration.
// PUT /admin/credit/config
func (h *CreditHandler) AdminUpdateGlobalCreditConfig(c *gin.Context) {
	var config service.GlobalCreditConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.creditSvc.UpdateGlobalCreditConfig(&config); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessMsg(c, "global credit config updated")
}

// AdminIndex returns detailed credit limit info for a specific user (admin).
// GET /admin/credit/index
func (h *CreditHandler) AdminIndex(c *gin.Context) {
	uid := c.Query("uid")
	if uid == "" {
		response.BadRequest(c, "uid is required")
		return
	}

	var user struct {
		Username              string  `json:"username"`
		PhoneNumber           string  `json:"phonenumber"`
		Email                 string  `json:"email"`
		IsOpenCreditLimit     int     `json:"is_open_credit_limit"`
		CreditLimit           float64 `json:"credit_limit"`
		RepaymentDate         int     `json:"repayment_date"`
		BillGenerationDate    int     `json:"bill_generation_date"`
		BillRepaymentPeriod   int     `json:"bill_repayment_period"`
		CreditLimitCreateTime int64   `json:"credit_limit_create_time"`
		Prefix                string  `json:"prefix"`
		Suffix                string  `json:"suffix"`
	}
	h.db.Table("users u").
		Select("u.username, u.phone_number, u.email, u.is_open_credit_limit, u.credit_limit, u.repayment_date, u.bill_generation_date, u.bill_repayment_period, u.credit_limit_create_time, cu.prefix, cu.suffix").
		Joins("LEFT JOIN currencies cu ON u.currency = cu.id").
		Where("u.id = ?", uid).
		Scan(&user)

	// Calculate amount to be settled
	var amountToBeSettled float64
	h.db.Model(&model.Invoice{}).
		Where("status = ? AND use_credit_limit = ? AND invoice_id = ? AND uid = ?", "Paid", 1, 0, uid).
		Select("COALESCE(SUM(total), 0)").
		Scan(&amountToBeSettled)

	// Calculate unpaid
	var unpaid float64
	h.db.Model(&model.Invoice{}).
		Where("type = ? AND status = ? AND uid = ?", "credit_limit", "Unpaid", uid).
		Select("COALESCE(SUM(total), 0)").
		Scan(&unpaid)

	creditLimitUsed := amountToBeSettled + unpaid
	creditLimitBalance := user.CreditLimit - creditLimitUsed
	if creditLimitBalance < 0 {
		creditLimitBalance = 0
	}

	creditUsedPercent := 0.0
	if user.CreditLimit > 0 && user.CreditLimit > creditLimitUsed {
		creditUsedPercent = creditLimitUsed / user.CreditLimit * 100
	} else if user.CreditLimit <= creditLimitUsed {
		creditUsedPercent = 100
	}

	// Get this month's bill
	var thisMonthBill map[string]interface{}
	h.db.Table("invoices").
		Where("type = ? AND uid = ? AND credit_limit_prepayment = 0", "credit_limit", uid).
		Where("created_at >= ?", time.Now().Format("2006-01")).
		Order("created_at DESC").
		Limit(1).
		Find(&thisMonthBill)

	// Get credit config
	creditConfig, _ := h.creditSvc.GetGlobalCreditConfig()

	response.Success(c, gin.H{
		"user": gin.H{
			"username":                   user.Username,
			"phonenumber":                user.PhoneNumber,
			"email":                      user.Email,
			"is_open_credit_limit":       user.IsOpenCreditLimit,
			"credit_limit":               user.CreditLimit,
			"repayment_date":             user.RepaymentDate,
			"bill_generation_date":       user.BillGenerationDate,
			"bill_repayment_period":      user.BillRepaymentPeriod,
			"amount_to_be_settled":       amountToBeSettled,
			"credit_limit_used":          creditLimitUsed,
			"credit_limit_balance":       creditLimitBalance,
			"credit_limit_used_percent":  creditUsedPercent,
			"this_month_bill":            thisMonthBill,
		},
		"credit_limit_config": creditConfig,
	})
}

// AdminCreditLog returns credit change logs for a user (admin).
// GET /admin/credit/log
func (h *CreditHandler) AdminCreditLog(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	uid := c.Query("uid")
	keywords := c.Query("keywords")
	logType := c.Query("type")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	query := h.db.Table("credit_logs cl").
		Joins("LEFT JOIN users u ON cl.admin_id = u.id").
		Select("cl.id, cl.user_id, cl.created_at, cl.description, cl.type, cl.ip, u.username as admin_name")

	if uid != "" {
		query = query.Where("cl.user_id = ?", uid)
	}
	if keywords != "" {
		query = query.Where("cl.description LIKE ?", "%"+keywords+"%")
	}
	if logType != "" {
		query = query.Where("cl.type = ?", logType)
	}
	if startTime != "" && endTime != "" {
		query = query.Where("cl.created_at >= ? AND cl.created_at <= ?", startTime, endTime)
	}

	var total int64
	query.Count(&total)

	var logs []map[string]interface{}
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("cl.id DESC").Find(&logs)

	response.SuccessPage(c, logs, total, page, pageSize)
}

// AdminCreditInvoiceList returns credit invoices list (admin).
// GET /admin/credit/invoices
func (h *CreditHandler) AdminCreditInvoiceList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	order := c.DefaultQuery("order", "id")
	sort := c.DefaultQuery("sort", "DESC")
	uid := c.Query("uid")
	paymentStatus := c.Query("payment_status")

	query := h.db.Model(&model.Invoice{}).Where("use_credit_limit = ?", 1)

	if uid != "" {
		query = query.Where("uid = ?", uid)
	}
	if paymentStatus != "" {
		switch paymentStatus {
		case "Paid":
			query = query.Where("status = ?", "Paid")
		case "Unpaid":
			query = query.Where("status = ? AND due_time > ?", "Unpaid", time.Now().Unix())
		case "Overdue":
			query = query.Where("status = ? AND due_time <= ?", "Unpaid", time.Now().Unix())
		}
	}

	var total int64
	query.Count(&total)

	var invoices []model.Invoice
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order(order + " " + sort).Find(&invoices)

	// Get invoice status config
	statusConfig := map[string]string{
		"Paid":      "已支付",
		"Unpaid":    "未支付",
		"Cancelled": "已取消",
	}

	response.SuccessPage(c, gin.H{
		"invoices":              invoices,
		"invoice_status":        statusConfig,
		"credit_limit_invoice_status": map[string]string{
			"Prepayment": "预付款",
			"Paid":       "已支付",
			"Unpaid":     "未支付",
			"Overdue":    "已逾期",
		},
	}, total, page, pageSize)
}

// Delete disables credit for a user (admin).
// DELETE /admin/credit/users/:uid
func (h *CreditHandler) Delete(c *gin.Context) {
	uid, err := strconv.ParseUint(c.Param("uid"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	// Check user exists
	var user model.User
	if err := h.db.First(&user, uid).Error; err != nil {
		response.NotFound(c, "user not found")
		return
	}

	// Disable credit
	if err := h.db.Model(&model.User{}).Where("id = ?", uid).
		Update("is_open_credit_limit", 0).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessMsg(c, "credit disabled successfully")
}

// GetSearch performs advanced search on credit invoices (admin).
// GET /admin/credit/search
func (h *CreditHandler) GetSearch(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := h.db.Model(&model.Invoice{}).Where("use_credit_limit = ?", 1)

	// Apply filters
	if uid := c.Query("uid"); uid != "" {
		query = query.Where("uid = ?", uid)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if paymentStatus := c.Query("payment_status"); paymentStatus != "" {
		switch paymentStatus {
		case "Paid":
			query = query.Where("status = ?", "Paid")
		case "Unpaid":
			query = query.Where("status = ? AND due_time > ?", "Unpaid", time.Now().Unix())
		case "Overdue":
			query = query.Where("status = ? AND due_time <= ?", "Unpaid", time.Now().Unix())
		}
	}
	if startTime := c.Query("start_time"); startTime != "" {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime := c.Query("end_time"); endTime != "" {
		query = query.Where("created_at <= ?", endTime)
	}

	var total int64
	query.Count(&total)

	var invoices []model.Invoice
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&invoices).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessPage(c, invoices, total, page, pageSize)
}

// Prepayment allows a user to repay credit limit bills before due date.
// POST /credit/prepayment
func (h *CreditHandler) Prepayment(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	// Check if user has unpaid credit limit bills
	var unpaidCount int64
	h.db.Model(&model.Invoice{}).
		Where("user_id = ? AND type = ? AND status = ?", userID, "credit_limit", 0).
		Count(&unpaidCount)
	if unpaidCount > 0 {
		response.BadRequest(c, "当前有未还款信用额账单，不可提前还款")
		return
	}

	// Calculate amount from paid invoices using credit limit that have no associated prepayment invoice
	var totalAmount float64
	h.db.Model(&model.Invoice{}).
		Where("user_id = ? AND status = ? AND use_credit_limit = ? AND linked_invoice_id = ?", userID, 1, 1, 0).
		Select("COALESCE(SUM(total), 0)").
		Scan(&totalAmount)

	if totalAmount <= 0 {
		response.BadRequest(c, "无需提前还款的金额")
		return
	}

	// Create prepayment invoice
	invoice := model.Invoice{
		UserID:                userID,
		Type:                  "credit_limit",
		CreditLimitPrepayment: 1,
		Total:                 datatypes.DecimalFromString(fmt.Sprintf("%.4f", totalAmount)),
		Status:                0, // Unpaid
	}
	if err := h.db.Create(&invoice).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}

	// Link original invoices to this prepayment invoice
	h.db.Model(&model.Invoice{}).
		Where("user_id = ? AND status = ? AND use_credit_limit = ? AND linked_invoice_id = ?", userID, 1, 1, 0).
		Update("linked_invoice_id", invoice.ID)

	response.Success(c, gin.H{
		"invoice_id": invoice.ID,
		"amount":     totalAmount,
		"message":    "提前还款账单已生成",
	})
}
