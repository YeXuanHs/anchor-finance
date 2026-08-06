package handler

import (
	"strconv"
	"time"

	"anchorfinance/internal/model"
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

// GetStatistics 获取销售统计
func (h *SaleHandler) GetStatistics(c *gin.Context) {
	saleID, _ := strconv.ParseUint(c.DefaultQuery("sale_id", "0"), 10, 32)
	timeRange := c.DefaultQuery("time", "month")

	now := time.Now()
	var startTime time.Time

	switch timeRange {
	case "today":
		startTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	case "week":
		startTime = now.AddDate(0, 0, -7)
	case "month":
		startTime = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	case "last_month":
		startTime = time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location())
	case "year":
		startTime = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	default:
		startTime = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	}

	query := h.saleSvc.GetDB().Model(&model.Invoice{}).Where("status = ? AND created_at >= ?", 1, startTime)
	if saleID > 0 {
		query = query.Where("sale_id = ?", saleID)
	}

	var totalAmount float64
	var totalCount int64
	query.Select("COALESCE(SUM(total), 0)").Scan(&totalAmount)
	query.Count(&totalCount)

	// 获取本月和上月对比
	lastMonthStart := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location())
	lastMonthEnd := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	var lastMonthAmount float64
	h.saleSvc.GetDB().Model(&model.Invoice{}).Where("status = ? AND created_at >= ? AND created_at < ?", 1, lastMonthStart, lastMonthEnd).
		Select("COALESCE(SUM(total), 0)").Scan(&lastMonthAmount)

	rate := float64(0)
	if lastMonthAmount > 0 {
		rate = (totalAmount - lastMonthAmount) / lastMonthAmount * 100
	}

	response.Success(c, gin.H{
		"total_amount":   totalAmount,
		"total_count":    totalCount,
		"last_month":     lastMonthAmount,
		"comparison":     rate,
		"time_range":     timeRange,
	})
}

// GetSaleRecords 获取销售记录
func (h *SaleHandler) GetSaleRecords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	saleID, _ := strconv.ParseUint(c.DefaultQuery("sale_id", "0"), 10, 32)
	username := c.Query("username")
	productName := c.Query("product_name")

	query := h.saleSvc.GetDB().Model(&model.Invoice{}).Where("sale_id > 0")
	if saleID > 0 {
		query = query.Where("sale_id = ?", saleID)
	}
	if username != "" {
		query = query.Joins("LEFT JOIN users ON users.id = invoices.user_id").Where("users.username LIKE ?", "%"+username+"%")
	}
	if productName != "" {
		query = query.Joins("LEFT JOIN products ON products.id = invoices.product_id").Where("products.name LIKE ?", "%"+productName+"%")
	}

	var total int64
	query.Count(&total)

	type SaleRecord struct {
		ID          uint    `json:"id"`
		InvoiceID   uint    `json:"invoice_id"`
		UserID      uint    `json:"user_id"`
		Username    string  `json:"username"`
		ProductName string  `json:"product_name"`
		Amount      float64 `json:"amount"`
		SaleID      uint    `json:"sale_id"`
		CreatedAt   string  `json:"created_at"`
	}

	var records []SaleRecord
	query.Select("invoices.id, invoices.id as invoice_id, invoices.user_id, users.username, '' as product_name, invoices.total as amount, invoices.sale_id, invoices.created_at").
		Joins("LEFT JOIN users ON users.id = invoices.user_id").
		Order("invoices.id DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&records)

	response.Success(c, gin.H{"list": records, "total": total})
}

// GetSaleUsers 获取销售关联的用户列表
func (h *SaleHandler) GetSaleUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	saleID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var total int64
	h.saleSvc.GetDB().Model(&model.User{}).Where("sale_id = ?", saleID).Count(&total)

	var users []model.User
	h.saleSvc.GetDB().Where("sale_id = ?", saleID).Select("id, username, email, phone, created_at").
		Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users)

	response.Success(c, gin.H{"list": users, "total": total})
}

// GetAdminList 获取销售管理员列表
func (h *SaleHandler) GetAdminList(c *gin.Context) {
	var admins []struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Status   int8   `json:"status"`
	}
	h.saleSvc.GetDB().Model(&model.User{}).Where("is_sale = 1").
		Select("id, username, email, status").
		Find(&admins)
	response.Success(c, admins)
}

// SetSaleStatus 设置销售启用状态
func (h *SaleHandler) SetSaleStatus(c *gin.Context) {
	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// 保存到系统配置
	value := "0"
	if req.Enabled {
		value = "1"
	}
	h.saleSvc.GetDB().Where("key = ?", "sale_enabled").Assign(map[string]interface{}{"key": "sale_enabled", "value": value}).
		FirstOrCreate(&model.SystemConfig{})

	response.Success(c, gin.H{"enabled": req.Enabled})
}

// GetSaleStatus 获取销售启用状态
func (h *SaleHandler) GetSaleStatus(c *gin.Context) {
	var config model.SystemConfig
	result := h.saleSvc.GetDB().Where("key = ?", "sale_enabled").First(&config)
	enabled := result.Error == nil && config.Value == "1"
	response.Success(c, gin.H{"enabled": enabled})
}

// GetTimetype returns time type options for sale promotions.
// GET /admin/sale/timetypes
func (h *SaleHandler) GetTimetype(c *gin.Context) {
	timeTypes := []map[string]interface{}{
		{"id": "month", "name": "按月"},
		{"id": "quarter", "name": "按季度"},
		{"id": "half_year", "name": "半年"},
		{"id": "year", "name": "按年"},
		{"id": "biennial", "name": "两年"},
		{"id": "triennial", "name": "三年"},
		{"id": "onetime", "name": "一次性"},
		{"id": "free", "name": "免费"},
	}
	response.Success(c, gin.H{"time_type": timeTypes})
}

// DelSaleLadder deletes a sale ladder entry.
// DELETE /admin/sale/ladder/:id
func (h *SaleHandler) DelSaleLadder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid ladder id")
		return
	}

	if err := h.saleSvc.DeleteLadder(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "sale ladder deleted")
}
