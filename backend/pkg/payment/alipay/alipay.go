package alipay

import (
	"context"
	"fmt"
	"time"

	"github.com/anchor-finance/backend/pkg/payment"
)

func init() {
	payment.Register(&Gateway{})
}

// Gateway implements the payment.Gateway interface for Alipay.
type Gateway struct {
	AppID      string
	PrivateKey string
	PublicKey  string
	NotifyURL  string
	ReturnURL  string
}

// Name returns the gateway identifier.
func (g *Gateway) Name() string {
	return "alipay"
}

// CreatePayment creates an Alipay payment and returns a pay URL.
func (g *Gateway) CreatePayment(ctx context.Context, param *payment.PaymentParam) (*payment.PaymentResult, error) {
	if param.Amount <= 0 {
		return nil, fmt.Errorf("alipay: amount must be positive")
	}

	tradeNo := param.TradeNo
	if tradeNo == "" {
		tradeNo = fmt.Sprintf("ALI_%s_%d", param.OrderNo, time.Now().UnixNano())
	}

	// In production, this would construct a real Alipay payment URL
	// using RSA-signed parameters and the Alipay gateway API.
	payURL := fmt.Sprintf(
		"https://openapi.alipay.com/gateway.do?trade_no=%s&total_amount=%.2f&subject=%s",
		tradeNo, param.Amount, param.Subject,
	)

	return &payment.PaymentResult{
		TradeNo: tradeNo,
		PayURL:  payURL,
		RawData: map[string]interface{}{
			"gateway":    "alipay",
			"order_no":   param.OrderNo,
			"trade_no":   tradeNo,
			"amount":     param.Amount,
			"subject":    param.Subject,
			"created_at": time.Now().Format(time.RFC3339),
		},
	}, nil
}

// QueryPayment queries the status of an Alipay payment.
func (g *Gateway) QueryPayment(ctx context.Context, tradeNo string) (*payment.QueryResult, error) {
	if tradeNo == "" {
		return nil, fmt.Errorf("alipay: trade_no is required")
	}

	// Mock implementation: in production, call Alipay's alipay.trade.query API.
	return &payment.QueryResult{
		TradeNo: tradeNo,
		Status:  "pending",
		Amount:  0,
		PaidAt:  nil,
	}, nil
}

// Refund processes a refund through Alipay.
func (g *Gateway) Refund(ctx context.Context, param *payment.RefundParam) (*payment.RefundResult, error) {
	if param.Amount <= 0 {
		return nil, fmt.Errorf("alipay: refund amount must be positive")
	}

	// Mock implementation: in production, call alipay.trade.refund API.
	return &payment.RefundResult{
		RefundNo: param.RefundNo,
		Status:   "processing",
	}, nil
}

// ParseNotify parses an Alipay async notification callback.
func (g *Gateway) ParseNotify(ctx context.Context, data []byte) (*payment.NotifyResult, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("alipay: empty notification data")
	}

	// Mock implementation: in production, parse the form-encoded notification,
	// verify the RSA signature, and extract trade status.
	// The real implementation should:
	// 1. Parse form values from data
	// 2. Verify sign using Alipay public key
	// 3. Check trade_status is TRADE_SUCCESS or TRADE_FINISHED
	// 4. Return parsed result
	return &payment.NotifyResult{
		Status: "success",
		Sign:   "mock_signature",
	}, nil
}
