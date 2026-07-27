package handler

import (
	"strconv"

	"github.com/anchor-finance/backend/internal/service"
	"github.com/anchor-finance/backend/pkg/logger"
	"github.com/anchor-finance/backend/pkg/response"

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
