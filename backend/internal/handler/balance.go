package handler

import (
	"net/http"
	"strconv"

	"github.com/anchor-finance/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// BalanceHandler handles user balance HTTP requests.
type BalanceHandler struct {
	balanceService *service.BalanceLogService
}

// NewBalanceHandler creates a new BalanceHandler.
func NewBalanceHandler(balanceService *service.BalanceLogService) *BalanceHandler {
	return &BalanceHandler{balanceService: balanceService}
}

// rechargeRequest is the payload for Recharge.
type rechargeRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

// GetBalance returns the authenticated user's current balance.
// GET /user/balance
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
// GET /user/balance/logs
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

// Recharge creates a recharge invoice for the user.
// POST /user/balance/recharge
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

	// TODO: Integrate with payment gateway to create a payment invoice.
	// For now, return a placeholder invoice response.
	// In a real implementation, this would:
	// 1. Create a pending invoice in the database
	// 2. Return payment URL or QR code for the user to complete payment
	// 3. Payment gateway callback would call BalanceLogService.AddBalance

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"amount":  req.Amount,
			"status":  "pending",
			"message": "recharge invoice created, please complete payment",
		},
	})
}
