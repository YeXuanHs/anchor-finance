package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type AffiliateHandler struct {
	affSvc *service.AffiliateService
	log    *logger.Logger
}

func NewAffiliateHandler(affSvc *service.AffiliateService, log *logger.Logger) *AffiliateHandler {
	return &AffiliateHandler{affSvc: affSvc, log: log}
}

// GetInfo returns the current user's affiliate info.
func (h *AffiliateHandler) GetInfo(c *gin.Context) {
	userID := c.GetUint("user_id")
	aff, err := h.affSvc.GetByUserID(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, aff)
}

// GetRecords returns paginated commission records.
func (h *AffiliateHandler) GetRecords(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	aff, err := h.affSvc.GetByUserID(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	records, total, err := h.affSvc.GetRecords(int(aff.ID), page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, records, total, page, pageSize)
}

type applyWithdrawRequest struct {
	Amount  float64 `json:"amount" binding:"required,gt=0"`
	Method  string  `json:"method" binding:"required,oneof=alipay bank balance"`
	Account string  `json:"account" binding:"required,max=100"`
}

// ApplyWithdraw creates a withdrawal request.
func (h *AffiliateHandler) ApplyWithdraw(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req applyWithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	aff, err := h.affSvc.GetByUserID(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	withdraw, err := h.affSvc.ApplyWithdraw(aff.ID, req.Amount, req.Method, req.Account)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, withdraw)
}

// GetWithdraws returns paginated withdrawal records.
func (h *AffiliateHandler) GetWithdraws(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	aff, err := h.affSvc.GetByUserID(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	withdraws, total, err := h.affSvc.GetWithdraws(int(aff.ID), page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, withdraws, total, page, pageSize)
}

// AdminGetList returns a paginated affiliate list.
func (h *AffiliateHandler) AdminGetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	affiliates, total, err := h.affSvc.GetList(page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, affiliates, total, page, pageSize)
}

type confirmRecordRequest struct {
	Status int8 `json:"status" binding:"required,oneof=2 3"`
}

// AdminConfirmRecord confirms or rejects a commission record.
func (h *AffiliateHandler) AdminConfirmRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid record id")
		return
	}

	var req confirmRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if req.Status == 2 {
		if err := h.affSvc.ConfirmRecord(uint(id)); err != nil {
			response.ServerError(c, err.Error())
			return
		}
		response.SuccessMsg(c, "record confirmed")
	} else {
		response.BadRequest(c, "only confirmation (status=2) is supported via this endpoint")
	}
}

type processWithdrawRequest struct {
	Approve   bool   `json:"approve" binding:"required"`
	AdminNote string `json:"admin_note" binding:"omitempty"`
}

// AdminProcessWithdraw approves or rejects a withdrawal.
func (h *AffiliateHandler) AdminProcessWithdraw(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid withdraw id")
		return
	}

	var req processWithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.affSvc.ProcessWithdraw(uint(id), req.Approve, req.AdminNote); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if req.Approve {
		response.SuccessMsg(c, "withdrawal approved")
	} else {
		response.SuccessMsg(c, "withdrawal rejected")
	}
}

// GatewayList returns available payment gateways for affiliate.
// GET /admin/affiliate/gateways
func (h *AffiliateHandler) GatewayList(c *gin.Context) {
	gateways, err := h.affSvc.GetGatewayList()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"gateway": gateways})
}

// ProductAffiPage returns affiliate settings for a product.
// GET /admin/affiliate/product/:pid
func (h *AffiliateHandler) ProductAffiPage(c *gin.Context) {
	pid, err := strconv.ParseUint(c.Param("pid"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product id")
		return
	}

	setting, err := h.affSvc.GetProductAffiSetting(uint(pid))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, setting)
}

// ProductAffiPost creates or updates affiliate settings for a product.
// POST /admin/affiliate/product
func (h *AffiliateHandler) ProductAffiPost(c *gin.Context) {
	var req service.ProductAffiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.affSvc.SaveProductAffiSetting(&req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "product affiliate setting saved")
}

// UserAffiBuyRecord returns affiliate purchase records for a user.
// GET /admin/affiliate/buy-records
func (h *AffiliateHandler) UserAffiBuyRecord(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	uid := c.Query("uid")
	if uid == "" {
		response.BadRequest(c, "uid is required")
		return
	}

	records, total, err := h.affSvc.GetUserBuyRecords(uid, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, records, total, page, pageSize)
}

// ==================== 新增缺失方法 ====================

// UserAffiPage returns affiliate settings page data for a user.
// GET /admin/affiliate/user-affi-page
func (h *AffiliateHandler) UserAffiPage(c *gin.Context) {
	uid := c.Query("uid")
	if uid == "" {
		response.BadRequest(c, "uid is required")
		return
	}

	data, datauser, err := h.affSvc.GetUserAffiPage(uid)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"data":     data,
		"datauser": datauser,
	})
}

// UserAffiBalance updates affiliate balance for a user (admin).
// POST /admin/affiliate/user-affi-balance
func (h *AffiliateHandler) UserAffiBalance(c *gin.Context) {
	var req struct {
		UID       uint    `json:"uid" binding:"required"`
		Withdrawn float64 `json:"withdrawn"`
		Balance   float64 `json:"balance"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.affSvc.UpdateUserAffiBalance(req.UID, req.Withdrawn, req.Balance); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessMsg(c, "affiliate balance updated")
}

// UserAffiPost creates or updates affiliate settings for a user (admin).
// POST /admin/affiliate/user-affi-post
func (h *AffiliateHandler) UserAffiPost(c *gin.Context) {
	var req service.UserAffiSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.affSvc.SaveUserAffiSetting(&req); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessMsg(c, "user affiliate setting saved")
}

// UserAffiList returns referred users list for an affiliate.
// GET /admin/affiliate/user-affi-list
func (h *AffiliateHandler) UserAffiList(c *gin.Context) {
	uid := c.Query("uid")
	if uid == "" {
		response.BadRequest(c, "uid is required")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	users, total, err := h.affSvc.GetUserAffiList(uid, page, pageSize, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessPage(c, users, total, page, pageSize)
}

// UserAffiRecord returns withdrawal records for a specific user affiliate.
// GET /admin/affiliate/user-affi-record
func (h *AffiliateHandler) UserAffiRecord(c *gin.Context) {
	uid := c.Query("uid")
	if uid == "" {
		response.BadRequest(c, "uid is required")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	records, total, err := h.affSvc.GetUserAffiRecord(uid, page, pageSize, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessPage(c, records, total, page, pageSize)
}

// GetTimeType returns available time type options.
// GET /admin/affiliate/time-type
func (h *AffiliateHandler) GetTimeType(c *gin.Context) {
	timeTypes := []map[string]interface{}{
		{"label": "最近7天", "value": "7day"},
		{"label": "最近30天", "value": "30day"},
		{"label": "本月", "value": "this_month"},
		{"label": "上月", "value": "last_month"},
		{"label": "本年", "value": "this_year"},
	}
	response.Success(c, gin.H{"time_type": timeTypes})
}

// GetIDs returns affiliate user IDs.
// GET /admin/affiliate/get-ids
func (h *AffiliateHandler) GetIDs(c *gin.Context) {
	uid := c.Query("uid")
	if uid == "" {
		response.BadRequest(c, "uid is required")
		return
	}

	uids, err := h.affSvc.GetAffiliateUserIDs(uid)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"uids": uids})
}

// AffiWithdrawRecord returns affiliate withdrawal records (admin).
// GET /admin/affiliate/withdraw-records
func (h *AffiliateHandler) AffiWithdrawRecord(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	status := c.Query("status")

	records, total, err := h.affSvc.GetAffiWithdrawRecords(page, pageSize, keyword, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessPage(c, records, total, page, pageSize)
}

// AffiWithdrawSH approves or rejects an affiliate withdrawal (admin).
// POST /admin/affiliate/withdraw-sh
func (h *AffiliateHandler) AffiWithdrawSH(c *gin.Context) {
	var req struct {
		ID      uint    `json:"id" binding:"required"`
		Type    int     `json:"type"`    // 1=balance, 3=account
		Status  int     `json:"status"`  // 2=approve, 3=reject
		Payment string  `json:"payment"`
		TransID string  `json:"trans_id"`
		Reason  string  `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	adminID := c.GetUint("user_id")
	if err := h.affSvc.ProcessAffiWithdrawSH(req.ID, req.Type, req.Status, req.Payment, req.TransID, req.Reason, adminID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if req.Status == 2 {
		response.SuccessMsg(c, "withdrawal approved")
	} else {
		response.SuccessMsg(c, "withdrawal rejected")
	}
}

// GetR is a helper that enriches commission records with status info.
// This is used internally but exposed as endpoint for frontend processing.
func (h *AffiliateHandler) GetR(c *gin.Context) {
	uid := c.Query("uid")
	if uid == "" {
		response.BadRequest(c, "uid is required")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	records, total, err := h.affSvc.GetCommissionRecords(uid, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessPage(c, records, total, page, pageSize)
}
