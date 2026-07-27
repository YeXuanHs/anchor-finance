package balance

import (
	"context"
	"fmt"
	"time"

	"github.com/anchor-finance/backend/pkg/payment"
)

func init() {
	payment.Register(&Gateway{})
}

// BalanceStore defines the interface for accessing user balance.
// Implement this in your repository layer.
type BalanceStore interface {
	GetBalance(userID string) (float64, error)
	DeductBalance(userID string, amount float64) error
	RefundBalance(userID string, amount float64) error
}

// Gateway implements the payment.Gateway interface for balance payments.
type Gateway struct {
	Store BalanceStore
}

// Name returns the gateway identifier.
func (g *Gateway) Name() string {
	return "balance"
}

// CreatePayment deducts the amount from user balance immediately.
func (g *Gateway) CreatePayment(ctx context.Context, param *payment.PaymentParam) (*payment.PaymentResult, error) {
	if param.Amount <= 0 {
		return nil, fmt.Errorf("balance: amount must be positive")
	}

	// Extract userID from context; set by middleware or service layer.
	userID, _ := ctx.Value("userID").(string)
	if userID == "" {
		return nil, fmt.Errorf("balance: user_id is required in context")
	}

	if g.Store == nil {
		return nil, fmt.Errorf("balance: store not configured")
	}

	balance, err := g.Store.GetBalance(userID)
	if err != nil {
		return nil, fmt.Errorf("balance: failed to get balance: %w", err)
	}
	if balance < param.Amount {
		return nil, fmt.Errorf("balance: insufficient balance (have %.2f, need %.2f)", balance, param.Amount)
	}

	if err := g.Store.DeductBalance(userID, param.Amount); err != nil {
		return nil, fmt.Errorf("balance: failed to deduct balance: %w", err)
	}

	tradeNo := param.TradeNo
	if tradeNo == "" {
		tradeNo = fmt.Sprintf("BAL_%s_%d", param.OrderNo, time.Now().UnixNano())
	}

	return &payment.PaymentResult{
		TradeNo: tradeNo,
		RawData: map[string]interface{}{
			"gateway":    "balance",
			"order_no":   param.OrderNo,
			"trade_no":   tradeNo,
			"amount":     param.Amount,
			"user_id":    userID,
			"created_at": time.Now().Format(time.RFC3339),
		},
	}, nil
}

// QueryPayment always returns success for balance payments (instant settlement).
func (g *Gateway) QueryPayment(ctx context.Context, tradeNo string) (*payment.QueryResult, error) {
	if tradeNo == "" {
		return nil, fmt.Errorf("balance: trade_no is required")
	}

	now := time.Now()
	return &payment.QueryResult{
		TradeNo: tradeNo,
		Status:  "success",
		PaidAt:  &now,
	}, nil
}

// Refund refunds the amount back to user balance.
func (g *Gateway) Refund(ctx context.Context, param *payment.RefundParam) (*payment.RefundResult, error) {
	if param.Amount <= 0 {
		return nil, fmt.Errorf("balance: refund amount must be positive")
	}

	userID, _ := ctx.Value("userID").(string)
	if userID == "" {
		return nil, fmt.Errorf("balance: user_id is required in context")
	}

	if g.Store == nil {
		return nil, fmt.Errorf("balance: store not configured")
	}

	if err := g.Store.RefundBalance(userID, param.Amount); err != nil {
		return nil, fmt.Errorf("balance: failed to refund balance: %w", err)
	}

	return &payment.RefundResult{
		RefundNo: param.RefundNo,
		Status:   "success",
	}, nil
}

// ParseNotify is not applicable for balance payments (no async callback).
func (g *Gateway) ParseNotify(ctx context.Context, data []byte) (*payment.NotifyResult, error) {
	return nil, fmt.Errorf("balance: notify not supported for balance payments")
}
