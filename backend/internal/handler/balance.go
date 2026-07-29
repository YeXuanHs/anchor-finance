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
	"github.com/shopspring/decimal"
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
	Gateway string  `json:"gateway" binding:"required"` // payment gateway code: alipay/wechat/qqpay/usdt/xunhupay/epay/stripe/paypal
}

// rechargeResponse is the response for a successful recharge initiation.
type rechargeResponse struct {
	InvoiceNo string  `json:"invoice_no"`
	TradeNo   string  `json:"trade_no"`
	Amount    float64 `json:"amount"`
	Gateway   string  `json:"gateway"`
	PayURL    string  `json:"pay_url,omitempty"`
	QrcodeURL string  `json:"qrcode_url,omitempty"`
	Status    string  `json:"status"`
}

// GetBalance returns the authenticated user's current balance.
// GET /balances
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
// GET /balances/logs
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
// POST /balances/recharge
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
	amountDecimal := decimal.NewFromFloat(req.Amount)
	invoice := model.Invoice{
		InvoiceNo:     invoiceNo,
		UserID:        userID,
		Type:          "recharge",
		Currency:      "CNY",
		SubTotal:      amountDecimal,
		Total:         amountDecimal,
		Status:        0, // pending
		PaymentMethod: req.Gateway,
	}
	if err := h.db.Create(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create invoice"})
		return
	}

	// 4. Create the payment via the gateway
	paymentGW, err := payment.Factory(gw.Code, gw.Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "payment gateway initialization failed"})
		return
	}

	orderNo := util.GenerateTransactionNo()
	param := &payment.PaymentParam{
		OrderNo:  orderNo,
		Amount:   req.Amount,
		Subject:  fmt.Sprintf("Balance recharge %.2f", req.Amount),
		ClientIP: c.ClientIP(),
	}

	result, err := paymentGW.CreatePayment(c.Request.Context(), param)
	if err != nil {
		// Update invoice status to cancelled
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
		Amount:        amountDecimal,
		Currency:      "CNY",
		Type:          "payment",
		Status:        0, // pending
		IPAddress:     c.ClientIP(),
	}
	if err := h.db.Create(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create transaction"})
		return
	}

	// 6. Return payment details to the user
	c.JSON(http.StatusOK, gin.H{
		"data": rechargeResponse{
			InvoiceNo: invoiceNo,
			TradeNo:   result.TradeNo,
			Amount:    req.Amount,
			Gateway:   req.Gateway,
			PayURL:    result.PayURL,
			QrcodeURL: result.QrcodeURL,
			Status:    "pending",
		},
	})
}

// RechargeNotify handles payment gateway callbacks for balance recharge.
// POST /balances/recharge/notify/:gateway
func (h *BalanceHandler) RechargeNotify(c *gin.Context) {
	gatewayCode := c.Param("gateway")
	if gatewayCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gateway parameter required"})
		return
	}

	// Read the raw body for gateway signature verification
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}

	// Look up the gateway config
	var gw model.PaymentGateway
	if err := h.db.Where("code = ? AND is_enabled = true", gatewayCode).First(&gw).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gateway not found or disabled"})
		return
	}

	// Create gateway instance and parse the notification
	paymentGW, err := payment.Factory(gw.Code, gw.Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gateway initialization failed"})
		return
	}

	// Prefer verified webhook parsing when the gateway supports it.
	var notify *payment.NotifyResult
	if verifier, ok := paymentGW.(payment.WebhookVerifier); ok {
		headers := make(map[string]string)
		for _, key := range []string{
			"Stripe-Signature",
			"paypal-auth-algo", "paypal-cert-url", "paypal-transmission-id",
			"paypal-transmission-sig", "paypal-transmission-time",
		} {
			if v := c.GetHeader(key); v != "" {
				headers[key] = v
			}
		}
		notify, err = verifier.VerifyAndParseNotify(c.Request.Context(), body, headers)
	} else {
		notify, err = paymentGW.ParseNotify(c.Request.Context(), body)
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid notification: %s", err.Error())})
		return
	}

	if notify.Status != "success" {
		c.JSON(http.StatusOK, gin.H{"status": "ignored"})
		return
	}

	// Find the transaction by trade/order number
	var transaction model.Transaction
	if err := h.db.Where("transaction_no = ? AND status = 0", notify.OrderNo).First(&transaction).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "transaction not found or already processed"})
		return
	}

	// Process in a database transaction
	err = h.db.Transaction(func(tx *gorm.DB) error {
		// Update transaction status to success
		if err := tx.Model(&transaction).Updates(map[string]interface{}{
			"status":       1, // success
			"completed_at": gorm.Expr("NOW()"),
		}).Error; err != nil {
			return fmt.Errorf("update transaction: %w", err)
		}

		// Update invoice status to paid
		if transaction.InvoiceID != nil {
			if err := tx.Model(&model.Invoice{}).Where("id = ?", *transaction.InvoiceID).Updates(map[string]interface{}{
				"status":         1, // paid
				"transaction_id": notify.OrderNo,
				"paid_at":        gorm.Expr("NOW()"),
			}).Error; err != nil {
				return fmt.Errorf("update invoice: %w", err)
			}
		}

		// Credit the user's balance
		amount := notify.Amount
		if amount <= 0 {
			// Fallback: read amount from transaction
			f64, _ := transaction.Amount.Float64Value()
			amount = f64
		}
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
// GET /balances/recharge/status/:invoice_no
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
// GET /balances/gateways
func (h *BalanceHandler) GetEnabledGateways(c *gin.Context) {
	var gateways []model.PaymentGateway
	if err := h.db.Where("is_enabled = true").
		Order("sort_order ASC, id ASC").
		Find(&gateways).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get gateways"})
		return
	}

	type gatewayInfo struct {
		ID                  uint    `json:"id"`
		Name                string  `json:"name"`
		Code                string  `json:"code"`
		Description         string  `json:"description"`
		Icon                string  `json:"icon"`
		FeeRate             float64 `json:"fee_rate"`
		MinAmount           float64 `json:"min_amount"`
		MaxAmount           float64 `json:"max_amount"`
		SupportedCurrencies string  `json:"supported_currencies"`
	}

	result := make([]gatewayInfo, len(gateways))
	for i, gw := range gateways {
		result[i] = gatewayInfo{
			ID:                  gw.ID,
			Name:                gw.Name,
			Code:                gw.Code,
			Description:         gw.Description,
			Icon:                gw.Icon,
			FeeRate:             gw.FeeRate,
			MinAmount:           gw.MinAmount,
			MaxAmount:           gw.MaxAmount,
			SupportedCurrencies: gw.SupportedCurrencies,
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
		Amount float64 `json:"amount" binding:"required,gt=0"`
		Method string  `json:"method"` // bank/alipay/wechat
		Account string `json:"account"` // withdrawal account
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check balance
	balance, err := h.balanceService.GetBalance(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get balance"})
		return
	}
	if balance < req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient balance"})
		return
	}

	// Create withdrawal request (stored as a balance log with negative amount)
	err = h.db.Transaction(func(tx *gorm.DB) error {
		// Deduct balance
		result := tx.Model(&model.User{}).Where("id = ?", userID).
			UpdateColumn("balance", gorm.Expr("balance - ?", req.Amount))
		if result.Error != nil {
			return result.Error
		}

		// Get new balance
		var user model.User
		tx.Select("balance").First(&user, userID)

		// Log withdrawal
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
