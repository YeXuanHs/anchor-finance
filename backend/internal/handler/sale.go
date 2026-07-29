package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type SaleHandler struct {
	saleSvc *service.SaleService
	log     *logger.Logger
}

func NewSaleHandler(saleSvc *service.SaleService, log *logger.Logger) *SaleHandler {
	return &SaleHandler{saleSvc: saleSvc, log: log}
}

// Create creates a new sale promotion (admin).
func (h *SaleHandler) Create(c *gin.Context) {
	var req service.CreateSalePromotionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	promo, err := h.saleSvc.Create(req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, promo)
}

// GetDetail returns a single promotion.
func (h *SaleHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid promotion id")
		return
	}

	promo, err := h.saleSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "promotion not found")
		return
	}
	response.Success(c, promo)
}

// GetList returns all promotions (admin).
func (h *SaleHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var promoType *string
	if t := c.Query("type"); t != "" {
		promoType = &t
	}
	var status *int
	if s := c.Query("status"); s != "" {
		v, _ := strconv.Atoi(s)
		status = &v
	}

	promos, total, err := h.saleSvc.GetList(page, pageSize, promoType, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, promos, total, page, pageSize)
}

// Update modifies a promotion (admin).
func (h *SaleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid promotion id")
		return
	}

	var req service.UpdateSalePromotionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	promo, err := h.saleSvc.Update(uint(id), req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, promo)
}

// Delete soft-deletes a promotion (admin).
func (h *SaleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid promotion id")
		return
	}

	if err := h.saleSvc.Delete(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "promotion deleted")
}

// Enable activates a promotion (admin).
func (h *SaleHandler) Enable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid promotion id")
		return
	}

	if err := h.saleSvc.Enable(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "promotion enabled")
}

// Disable deactivates a promotion (admin).
func (h *SaleHandler) Disable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid promotion id")
		return
	}

	if err := h.saleSvc.Disable(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "promotion disabled")
}

// SetStatus enables or disables a sale promotion (admin).
func (h *SaleHandler) SetStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid promotion id")
		return
	}

	var req struct {
		Status int `json:"status" binding:"required,oneof=0 1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.saleSvc.SetStatus(uint(id), req.Status); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "promotion status updated")
}

// GetUsageStats returns usage statistics for a promotion (admin).
func (h *SaleHandler) GetUsageStats(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid promotion id")
		return
	}

	stats, err := h.saleSvc.GetUsageStats(uint(id))
	if err != nil {
		response.NotFound(c, "promotion not found")
		return
	}
	response.Success(c, stats)
}

// GetUsageLogs returns usage logs for a promotion (admin).
func (h *SaleHandler) GetUsageLogs(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid promotion id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	logs, total, err := h.saleSvc.GetUsageLogs(uint(id), page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, logs, total, page, pageSize)
}

// GetCommissionLadder returns commission tiers for a promotion (admin).
func (h *SaleHandler) GetCommissionLadder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid promotion id")
		return
	}

	tiers, err := h.saleSvc.GetCommissionLadder(uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, tiers)
}

// SetCommissionLadder replaces commission tiers for a promotion (admin).
func (h *SaleHandler) SetCommissionLadder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid promotion id")
		return
	}

	var req []service.SaleCommission
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.saleSvc.SetCommissionLadder(uint(id), req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "commission ladder updated")
}

// CalculateCommission calculates commission for a given amount (admin).
func (h *SaleHandler) CalculateCommission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid promotion id")
		return
	}

	var req struct {
		Amount float64 `json:"amount" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	commission, err := h.saleSvc.CalculateCommission(uint(id), req.Amount)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"commission": commission})
}

// ValidateUsage checks if a user can use a promotion (admin).
func (h *SaleHandler) ValidateUsage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid promotion id")
		return
	}

	var req struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	valid, reason, err := h.saleSvc.ValidateUsage(uint(id), req.UserID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"valid": valid, "reason": reason})
}
