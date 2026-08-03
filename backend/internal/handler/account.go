package handler

import (
	"strconv"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AccountHandler handles transaction record (交易流水) HTTP requests.
type AccountHandler struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewAccountHandler creates a new AccountHandler.
func NewAccountHandler(db *gorm.DB, log *logger.Logger) *AccountHandler {
	return &AccountHandler{db: db, log: log}
}

// Index returns paginated transaction records with filters (admin).
// GET /admin/accounts
func (h *AccountHandler) Index(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	order := c.DefaultQuery("order", "id")
	sort := c.DefaultQuery("sort", "DESC")

	allowedOrders := map[string]bool{
		"id": true, "username": true, "create_time": true,
		"gateway": true, "description": true, "amount_in": true,
		"fees": true, "amount_out": true,
	}
	if !allowedOrders[order] {
		order = "id"
	}
	if sort != "ASC" && sort != "DESC" {
		sort = "DESC"
	}

	query := h.db.Model(&model.Account{}).Where("delete_time = ?", 0)

	if uid := c.Query("uid"); uid != "" {
		query = query.Where("uid = ?", uid)
	}
	if gateway := c.Query("gateway"); gateway != "" {
		query = query.Where("gateway = ?", gateway)
	}
	if transID := c.Query("trans_id"); transID != "" {
		query = query.Where("trans_id LIKE ?", "%"+transID+"%")
	}
	if desc := c.Query("description"); desc != "" {
		query = query.Where("description LIKE ?", "%"+desc+"%")
	}
	if startTime := c.Query("start_time"); startTime != "" {
		query = query.Where("pay_time >= ?", startTime)
	}
	if endTime := c.Query("end_time"); endTime != "" {
		query = query.Where("pay_time <= ?", endTime)
	}

	var total int64
	query.Count(&total)

	var accounts []model.Account
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order(order + " " + sort).Find(&accounts).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessPage(c, accounts, total, page, pageSize)
}

// Create returns available data for creating a transaction record.
// GET /admin/accounts/create
func (h *AccountHandler) Create(c *gin.Context) {
	uid := c.Query("uid")

	type UserOption struct {
		ID       uint   `json:"value"`
		Nickname string `json:"label"`
	}
	var users []UserOption
	h.db.Model(&model.User{}).Where("is_sale = ?", 1).Select("id as value, user_nickname as label").Find(&users)

	type GatewayOption struct {
		Name  string `json:"name"`
		Title string `json:"title"`
	}
	var gateways []GatewayOption
	h.db.Model(&model.PaymentGateway{}).Where("status = ?", 1).Select("name, title").Find(&gateways)
	otherPay := []GatewayOption{
		{Name: "creditPay", Title: "余额支付"},
		{Name: "creditLimitPay", Title: "信用额支付"},
	}
	gateways = append(gateways, otherPay...)

	var invoices []uint
	if uid != "" {
		h.db.Model(&model.Invoice{}).Where("user_id = ?", uid).Pluck("id", &invoices)
	}

	response.Success(c, gin.H{
		"users":    users,
		"gateway":  gateways,
		"invoices": invoices,
	})
}

// CreateInvoice returns invoice list for a user.
// GET /admin/accounts/create-invoice
func (h *AccountHandler) CreateInvoice(c *gin.Context) {
	uid := c.Query("uid")
	if uid == "" {
		response.BadRequest(c, "uid is required")
		return
	}

	var invoices []uint
	h.db.Model(&model.Invoice{}).Where("user_id = ?", uid).Pluck("id", &invoices)

	response.Success(c, gin.H{"invoices": invoices})
}

// Save creates a new transaction record (admin).
// POST /admin/accounts
func (h *AccountHandler) Save(c *gin.Context) {
	var req struct {
		UID         uint    `json:"uid" binding:"required"`
		Currency    string  `json:"currency" binding:"required"`
		PayTime     int64   `json:"pay_time"`
		Description string  `json:"description"`
		TransID     string  `json:"trans_id"`
		InvoiceID   uint    `json:"invoice_id"`
		Gateway     string  `json:"gateway"`
		AmountIn    float64 `json:"amount_in"`
		Fees        float64 `json:"fees"`
		AmountOut   float64 `json:"amount_out"`
		Refund      int     `json:"refund"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if req.Description == "" && req.TransID == "" {
		response.BadRequest(c, "description or trans_id is required")
		return
	}
	if req.AmountIn == 0 && req.AmountOut == 0 && req.Fees == 0 {
		response.BadRequest(c, "at least one amount field must be non-zero")
		return
	}

	// Verify currency exists
	var currencyCount int64
	h.db.Model(&model.Currency{}).Where("code = ?", req.Currency).Count(&currencyCount)
	if currencyCount == 0 {
		response.BadRequest(c, "invalid currency")
		return
	}

	// Verify invoice if provided
	if req.InvoiceID > 0 {
		var invoice model.Invoice
		if err := h.db.First(&invoice, req.InvoiceID).Error; err != nil {
			response.BadRequest(c, "invoice not found")
			return
		}
	}

	if req.Refund > 0 && req.InvoiceID > 0 {
		response.BadRequest(c, "refund to balance cannot be linked to an invoice")
		return
	}

	account := model.Account{
		UID:         req.UID,
		Currency:    req.Currency,
		PayTime:     time.Now().Unix(),
		Description: req.Description,
		TransID:     req.TransID,
		InvoiceID:   req.InvoiceID,
		Gateway:     req.Gateway,
		AmountIn:    req.AmountIn,
		Fees:        req.Fees,
		AmountOut:   req.AmountOut,
		CreateTime:  time.Now().Unix(),
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&account).Error; err != nil {
			return err
		}
		// If refund to balance, add credit
		if req.Refund > 0 && req.AmountIn > 0 {
			if err := tx.Model(&model.User{}).Where("id = ?", req.UID).
				UpdateColumn("credit", gorm.Expr("credit + ?", req.AmountIn)).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		response.BadRequest(c, "failed to create account record")
		return
	}

	response.SuccessMsg(c, "请求成功")
}

// Read returns a single transaction record (admin).
// GET /admin/accounts/:id
func (h *AccountHandler) Read(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid account id")
		return
	}

	var account model.Account
	if err := h.db.First(&account, id).Error; err != nil {
		response.NotFound(c, "account not found")
		return
	}

	var invoices []uint
	h.db.Model(&model.Invoice{}).Where("user_id = ?", account.UID).Pluck("id", &invoices)

	var gateways []model.PaymentGateway
	h.db.Where("status = ?", 1).Find(&gateways)

	response.Success(c, gin.H{
		"list":     account,
		"gateway":  gateways,
		"invoices": invoices,
	})
}

// Update updates a transaction record (admin).
// PUT /admin/accounts/:id
func (h *AccountHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid account id")
		return
	}

	var req struct {
		UID         uint    `json:"uid"`
		Gateway     string  `json:"gateway"`
		PayTime     int64   `json:"pay_time"`
		Description string  `json:"description"`
		AmountIn    float64 `json:"amount_in"`
		Fees        float64 `json:"fees"`
		AmountOut   float64 `json:"amount_out"`
		InvoiceID   uint    `json:"invoice_id"`
		TransID     string  `json:"trans_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := map[string]interface{}{
		"uid":          req.UID,
		"gateway":      req.Gateway,
		"pay_time":     req.PayTime,
		"description":  req.Description,
		"amount_in":    req.AmountIn,
		"fees":         req.Fees,
		"amount_out":   req.AmountOut,
		"invoice_id":   req.InvoiceID,
		"trans_id":     req.TransID,
		"update_time":  time.Now().Unix(),
	}

	if err := h.db.Model(&model.Account{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMsg(c, "请求成功")
}

// Delete deletes transaction records (admin).
// DELETE /admin/accounts/:id
func (h *AccountHandler) Delete(c *gin.Context) {
	ids := c.QueryArray("ids")
	if len(ids) == 0 {
		if id := c.Param("id"); id != "" {
			ids = []string{id}
		}
	}

	if len(ids) == 0 {
		response.BadRequest(c, "id is required")
		return
	}

	for _, id := range ids {
		if err := h.db.Delete(&model.Account{}, id).Error; err != nil {
			h.log.Errorf("failed to delete account %s: %v", id, err)
		}
	}

	response.SuccessMsg(c, "删除成功")
}

// SearchPage returns search page data (gateways, sale list).
// GET /admin/accounts/search-page
func (h *AccountHandler) SearchPage(c *gin.Context) {
	type UserOption struct {
		ID       uint   `json:"value"`
		Nickname string `json:"label"`
	}
	var saleList []UserOption
	h.db.Model(&model.User{}).Where("is_sale = ?", 1).Select("id as value, user_nickname as label").Find(&saleList)

	type GatewayOption struct {
		Name  string `json:"name"`
		Title string `json:"title"`
	}
	var gateways []GatewayOption
	h.db.Model(&model.PaymentGateway{}).Where("status = ?", 1).Select("name, title").Find(&gateways)
	otherPay := []GatewayOption{
		{Name: "creditPay", Title: "余额支付"},
		{Name: "creditLimitPay", Title: "信用额支付"},
	}
	gateways = append(gateways, otherPay...)

	response.Success(c, gin.H{
		"gateway":  gateways,
		"salelist": saleList,
	})
}
