package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/anchor-finance/backend/internal/model"
	"github.com/anchor-finance/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// CouponHandler handles coupon-related HTTP requests.
type CouponHandler struct {
	couponService *service.CouponService
}

// NewCouponHandler creates a new CouponHandler.
func NewCouponHandler(couponService *service.CouponService) *CouponHandler {
	return &CouponHandler{couponService: couponService}
}

// validateCouponRequest is the payload for ValidateCoupon.
type validateCouponRequest struct {
	Code      string  `json:"code" binding:"required"`
	ProductID uint    `json:"product_id"`
	Amount    float64 `json:"amount" binding:"required"`
}

// createCouponRequest is the payload for admin coupon creation.
type createCouponRequest struct {
	Code         string   `json:"code" binding:"required"`
	Type         string   `json:"type" binding:"required,oneof=fixed percent"`
	Value        float64  `json:"value" binding:"required,gt=0"`
	MinAmount    float64  `json:"min_amount"`
	MaxDiscount  float64  `json:"max_discount"`
	ProductID    *uint    `json:"product_id"`
	UserID       *uint    `json:"user_id"`
	StartDate    string   `json:"start_date" binding:"required"`
	EndDate      string   `json:"end_date" binding:"required"`
	TotalUses    int      `json:"total_uses"`
	PerUserLimit int      `json:"per_user_limit"`
	Enabled      *bool    `json:"enabled"`
	Description  string   `json:"description"`
}

// updateCouponRequest is the payload for admin coupon update.
type updateCouponRequest struct {
	Code         string   `json:"code"`
	Type         string   `json:"type" binding:"omitempty,oneof=fixed percent"`
	Value        float64  `json:"value"`
	MinAmount    float64  `json:"min_amount"`
	MaxDiscount  float64  `json:"max_discount"`
	ProductID    *uint    `json:"product_id"`
	UserID       *uint    `json:"user_id"`
	StartDate    string   `json:"start_date"`
	EndDate      string   `json:"end_date"`
	TotalUses    int      `json:"total_uses"`
	PerUserLimit int      `json:"per_user_limit"`
	Enabled      *bool    `json:"enabled"`
	Description  string   `json:"description"`
}

// ValidateCoupon checks if a coupon is valid for the given order.
// POST /coupons/validate
func (h *CouponHandler) ValidateCoupon(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	var req validateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	discount, coupon, err := h.couponService.Validate(req.Code, userID, req.ProductID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"valid": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":    true,
		"discount": discount,
		"coupon":   coupon,
	})
}

// GetUserCoupons returns the authenticated user's available coupons.
// GET /user/coupons
func (h *CouponHandler) GetUserCoupons(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	coupons, err := h.couponService.GetUserCoupons(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get coupons"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": coupons})
}

// CreateCoupon creates a new coupon (admin).
// POST /admin/coupons
func (h *CouponHandler) CreateCoupon(c *gin.Context) {
	var req createCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate, err := parseTime(req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format, use RFC3339 or 2006-01-02"})
		return
	}
	endDate, err := parseTime(req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format, use RFC3339 or 2006-01-02"})
		return
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	coupon := &model.Coupon{
		Code:         req.Code,
		Type:         req.Type,
		Value:        req.Value,
		MinAmount:    req.MinAmount,
		MaxDiscount:  req.MaxDiscount,
		ProductID:    req.ProductID,
		UserID:       req.UserID,
		StartDate:    startDate,
		EndDate:      endDate,
		TotalUses:    req.TotalUses,
		PerUserLimit: req.PerUserLimit,
		Enabled:      enabled,
		Description:  req.Description,
	}

	if err := h.couponService.Create(coupon); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": coupon})
}

// UpdateCoupon updates an existing coupon (admin).
// PUT /admin/coupons/:id
func (h *CouponHandler) UpdateCoupon(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid coupon id"})
		return
	}

	coupon, err := h.couponService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "coupon not found"})
		return
	}

	var req updateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Code != "" {
		coupon.Code = req.Code
	}
	if req.Type != "" {
		coupon.Type = req.Type
	}
	if req.Value > 0 {
		coupon.Value = req.Value
	}
	if req.MinAmount > 0 {
		coupon.MinAmount = req.MinAmount
	}
	if req.MaxDiscount > 0 {
		coupon.MaxDiscount = req.MaxDiscount
	}
	if req.ProductID != nil {
		coupon.ProductID = req.ProductID
	}
	if req.UserID != nil {
		coupon.UserID = req.UserID
	}
	if req.StartDate != "" {
		t, err := parseTime(req.StartDate)
		if err == nil {
			coupon.StartDate = t
		}
	}
	if req.EndDate != "" {
		t, err := parseTime(req.EndDate)
		if err == nil {
			coupon.EndDate = t
		}
	}
	if req.TotalUses > 0 {
		coupon.TotalUses = req.TotalUses
	}
	if req.PerUserLimit > 0 {
		coupon.PerUserLimit = req.PerUserLimit
	}
	if req.Enabled != nil {
		coupon.Enabled = *req.Enabled
	}
	if req.Description != "" {
		coupon.Description = req.Description
	}

	if err := h.couponService.Update(coupon); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update coupon"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": coupon})
}

// DeleteCoupon soft-deletes a coupon (admin).
// DELETE /admin/coupons/:id
func (h *CouponHandler) DeleteCoupon(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid coupon id"})
		return
	}

	if err := h.couponService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "coupon deleted"})
}

// ListCoupons returns all coupons with pagination (admin).
// GET /admin/coupons
func (h *CouponHandler) ListCoupons(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	coupons, total, err := h.couponService.GetAll(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list coupons"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      coupons,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// parseTime attempts to parse a time string in RFC3339 or date-only format.
func parseTime(s string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot parse time: %s", s)
}
