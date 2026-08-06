package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/internal/util"
	"anchorfinance/pkg/payment"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BalanceHandler handles user balance HTTP requests.
type BalanceHandler struct {
	balanceService *service.BalanceLogService
	db             *gorm.DB
}

// NewBalanceHandler creates a new BalanceHandler.
func NewBalanceHandler(balanceService *service.BalanceLogService, db *gorm.DB) *BalanceHandler {
	return &BalanceHandler{
		balanceService: balanceService,
		db:             db,
	}
}

// rechargeRequest is the payload for Recharge.
type rechargeRequest struct {
	Amount  float64 `json:"amount" binding:"required,gt=0"`
	Gateway string  `json:"gateway" binding:"required"`
}

// GetBalance returns the authenticated user's current balance.
func (h *BalanceHandler) GetBalance(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	balance, err := h.balanceService.GetBalance(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"balance": balance}})
}

// GetBalanceLogs returns the user's balance change history.
func (h *BalanceHandler) GetBalanceLogs(c *gin.Context) {
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

	logs, total, err := h.balanceService.GetLogs(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get balance logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      logs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// Recharge creates a recharge invoice and initiates payment via the selected gateway.
func (h *BalanceHandler) Recharge(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	var req rechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Look up the payment gateway from the database
	var gw model.PaymentGateway
	if err := h.db.Where("code = ?", req.Gateway).First(&gw).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment gateway not found"})
		return
	}
	if !gw.IsEnabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment gateway is disabled"})
		return
	}

	// 2. Validate amount against gateway limits
	if gw.MinAmount > 0 && req.Amount < gw.MinAmount {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("minimum recharge amount is %.2f", gw.MinAmount),
		})
		return
	}
	if gw.MaxAmount > 0 && req.Amount > gw.MaxAmount {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("maximum recharge amount is %.2f", gw.MaxAmount),
		})
		return
	}

	// 3. Create a recharge invoice record
	invoiceNo := util.GenerateInvoiceNo()
	invoice := model.Invoice{
		InvoiceNo:     invoiceNo,
		UserID:        userID,
		Type:          "recharge",
		Currency:      "CNY",
		SubTotal:      req.Amount,
		Total:         req.Amount,
		Status:        0,
		PaymentMethod: req.Gateway,
	}
	if err := h.db.Create(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create invoice"})
		return
	}

	// 4. Create the payment via the gateway
	paymentGW, err := payment.Factory(gw.Gateway, gw.Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "payment gateway initialization failed"})
		return
	}

	// 如果是易支付或迅虎支付，需要设置支付类型
	if epayGW, ok := paymentGW.(*payment.EpayGateway); ok {
		epayGW.SetCode(gw.Code)
	}
	if xunhuGW, ok := paymentGW.(*payment.XunhuPayGateway); ok {
		xunhuGW.SetCode(gw.Code)
	}

	orderNo := util.GenerateTransactionNo()
	param := &payment.PaymentParam{
		OrderNo:   orderNo,
		Amount:    req.Amount,
		Subject:   fmt.Sprintf("余额充值 %.2f元", req.Amount),
		ReturnURL: fmt.Sprintf("%s/user/balance", getDomain(c)),
		NotifyURL: fmt.Sprintf("%s/api/v1/payments/notify/%s", getDomain(c), gw.Name),
		ClientIP:  c.ClientIP(),
	}

	result, err := paymentGW.CreatePayment(c.Request.Context(), param)
	if err != nil {
		h.db.Model(&invoice).Update("status", 3)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("payment creation failed: %s", err.Error())})
		return
	}

	// 5. Create a transaction record
	transaction := model.Transaction{
		TransactionNo: orderNo,
		UserID:        userID,
		InvoiceID:     &invoice.ID,
		Gateway:       req.Gateway,
		Amount:        req.Amount,
		Currency:      "CNY",
		Type:          "payment",
		Status:        0,
		IPAddress:     c.ClientIP(),
	}
	if err := h.db.Create(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create transaction"})
		return
	}

	// 6. Return payment details to the user
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"invoice_no": invoiceNo,
			"trade_no":   result.OrderNo,
			"amount":     req.Amount,
			"gateway":    req.Gateway,
			"pay_url":    result.Data,
			"type":       result.Type,
			"status":     "pending",
		},
	})
}

// RechargeNotify handles payment gateway callbacks for balance recharge.
func (h *BalanceHandler) RechargeNotify(c *gin.Context) {
	gatewayCode := c.Param("gateway")
	if gatewayCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gateway parameter required"})
		return
	}

	// Look up the gateway config
	var gw model.PaymentGateway
	if err := h.db.Where("code = ? AND is_enabled = true", gatewayCode).First(&gw).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gateway not found or disabled"})
		return
	}

	// Create gateway instance
	paymentGW, err := payment.Factory(gw.Code, gw.Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gateway initialization failed"})
		return
	}

	// Parse form/query params for notification verification
	notifyData := make(map[string]string)
	c.Request.ParseForm()
	for k, v := range c.Request.PostForm {
		if len(v) > 0 {
			notifyData[k] = v[0]
		}
	}
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			notifyData[k] = v[0]
		}
	}

	notify, err := paymentGW.VerifyNotification(c.Request.Context(), notifyData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid notification: %s", err.Error())})
		return
	}

	if !notify.Success {
		c.JSON(http.StatusOK, gin.H{"status": "ignored"})
		return
	}

	// Find the transaction by order number
	var transaction model.Transaction
	if err := h.db.Where("transaction_no = ? AND status = 0", notify.OrderNo).First(&transaction).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "transaction not found or already processed"})
		return
	}

	// Process in a database transaction
	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&transaction).Updates(map[string]interface{}{
			"status":       1,
			"completed_at": gorm.Expr("NOW()"),
		}).Error; err != nil {
			return fmt.Errorf("update transaction: %w", err)
		}

		if transaction.InvoiceID != nil {
			if err := tx.Model(&model.Invoice{}).Where("id = ?", *transaction.InvoiceID).Updates(map[string]interface{}{
				"status":         1,
				"transaction_id": notify.OrderNo,
				"paid_at":        gorm.Expr("NOW()"),
			}).Error; err != nil {
				return fmt.Errorf("update invoice: %w", err)
			}
		}

		amount := transaction.Amount
		if amount > 0 {
			desc := fmt.Sprintf("Balance recharge via %s, order: %s", gatewayCode, notify.OrderNo)
			if err := h.balanceService.AddBalance(transaction.UserID, amount, transaction.ID, "recharge", desc); err != nil {
				return fmt.Errorf("add balance: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// GetRechargeStatus checks the status of a pending recharge invoice.
func (h *BalanceHandler) GetRechargeStatus(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	invoiceNo := c.Param("invoice_no")
	if invoiceNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invoice_no required"})
		return
	}

	var invoice model.Invoice
	if err := h.db.Where("invoice_no = ? AND user_id = ?", invoiceNo, userID).First(&invoice).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invoice not found"})
		return
	}

	statusText := "pending"
	switch invoice.Status {
	case 1:
		statusText = "paid"
	case 3:
		statusText = "cancelled"
	case 4:
		statusText = "refunded"
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"invoice_no": invoice.InvoiceNo,
			"amount":     invoice.Total,
			"gateway":    invoice.PaymentMethod,
			"status":     statusText,
			"paid_at":    invoice.PaidAt,
		},
	})
}

// GetEnabledGateways returns the list of enabled payment gateways for recharge.
func (h *BalanceHandler) GetEnabledGateways(c *gin.Context) {
	var gateways []model.PaymentGateway
	if err := h.db.Where("is_enabled = true").
		Order("sort_order ASC, id ASC").
		Find(&gateways).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get gateways"})
		return
	}

	type gatewayInfo struct {
		ID        uint    `json:"id"`
		Name      string  `json:"name"`
		Title     string  `json:"title"`
		Code      string  `json:"code"`
		Icon      string  `json:"icon"`
		FeeRate   float64 `json:"fee_rate"`
		MinAmount float64 `json:"min_amount"`
		MaxAmount float64 `json:"max_amount"`
	}

	result := make([]gatewayInfo, len(gateways))
	for i, gw := range gateways {
		icon := payment.GetGatewayIcon(gw.Code)
		result[i] = gatewayInfo{
			ID:        gw.ID,
			Name:      gw.Name,
			Title:     gw.Title,
			Code:      gw.Code,
			Icon:      icon,
			FeeRate:   gw.FeeRate,
			MinAmount: gw.MinAmount,
			MaxAmount: gw.MaxAmount,
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// Withdraw processes a balance withdrawal request.
func (h *BalanceHandler) Withdraw(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "login required"})
		return
	}

	var req struct {
		Amount  float64 `json:"amount" binding:"required,gt=0"`
		Method  string  `json:"method"`
		Account string  `json:"account"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	balance, err := h.balanceService.GetBalance(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get balance"})
		return
	}
	if balance < req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient balance"})
		return
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&model.User{}).Where("id = ?", userID).
			UpdateColumn("balance", gorm.Expr("balance - ?", req.Amount))
		if result.Error != nil {
			return result.Error
		}

		var user model.User
		tx.Select("balance").First(&user, userID)

		log := &model.BalanceLog{
			UserID:      userID,
			Amount:      -req.Amount,
			Balance:     user.Balance,
			RelatedType: "withdrawal",
			Description: fmt.Sprintf("Withdraw %.2f via %s to %s", req.Amount, req.Method, req.Account),
		}
		return tx.Create(log).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "withdrawal submitted", "data": gin.H{
		"amount": req.Amount, "method": req.Method, "status": "pending",
	}})
}

// getDomain gets the domain from the request
func getDomain(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	host := c.Request.Host
	if h := c.GetHeader("X-Forwarded-Host"); h != "" {
		host = h
	}
	return fmt.Sprintf("%s://%s", scheme, host)
}
