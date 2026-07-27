package paypal

import (
	"context"
	"fmt"
	"time"

	"github.com/anchor-finance/backend/pkg/payment"
)

func init() {
	payment.Register(&Gateway{})
}

// Gateway implements the payment.Gateway interface for PayPal.
type Gateway struct {
	ClientID     string
	ClientSecret string
	Mode         string // "sandbox" or "live"
	ReturnURL    string
	CancelURL    string
}

// Name returns the gateway identifier.
func (g *Gateway) Name() string {
	return "paypal"
}

// baseURL returns the PayPal API base URL based on mode.
func (g *Gateway) baseURL() string {
	if g.Mode == "live" {
		return "https://api-m.paypal.com"
	}
	return "https://api-m.sandbox.paypal.com"
}

// CreatePayment creates a PayPal payment and returns an approval URL.
func (g *Gateway) CreatePayment(ctx context.Context, param *payment.PaymentParam) (*payment.PaymentResult, error) {
	if param.Amount <= 0 {
		return nil, fmt.Errorf("paypal: amount must be positive")
	}

	tradeNo := param.TradeNo
	if tradeNo == "" {
		tradeNo = fmt.Sprintf("PP_%s_%d", param.OrderNo, time.Now().UnixNano())
	}

	// In production, this would call PayPal's /v2/checkout/orders API
	// to create an order and obtain the approval link.
	approvalURL := fmt.Sprintf(
		"%s/checkout/orders/%s?approve",
		g.baseURL(), tradeNo,
	)

	return &payment.PaymentResult{
		TradeNo: tradeNo,
		PayURL:  approvalURL,
		RawData: map[string]interface{}{
			"gateway":    "paypal",
			"order_no":   param.OrderNo,
			"trade_no":   tradeNo,
			"amount":     param.Amount,
			"currency":   "USD",
			"created_at": time.Now().Format(time.RFC3339),
		},
	}, nil
}

// QueryPayment queries the status of a PayPal order.
func (g *Gateway) QueryPayment(ctx context.Context, tradeNo string) (*payment.QueryResult, error) {
	if tradeNo == "" {
		return nil, fmt.Errorf("paypal: trade_no is required")
	}

	// Mock: in production, call /v2/checkout/orders/{order_id}.
	return &payment.QueryResult{
		TradeNo: tradeNo,
		Status:  "pending",
		Amount:  0,
		PaidAt:  nil,
	}, nil
}

// Refund processes a refund through PayPal.
func (g *Gateway) Refund(ctx context.Context, param *payment.RefundParam) (*payment.RefundResult, error) {
	if param.Amount <= 0 {
		return nil, fmt.Errorf("paypal: refund amount must be positive")
	}

	// Mock: in production, call /v2/payments/captures/{capture_id}/refund.
	return &payment.RefundResult{
		RefundNo: param.RefundNo,
		Status:   "processing",
	}, nil
}

// ParseNotify parses a PayPal IPN (Instant Payment Notification) webhook.
func (g *Gateway) ParseNotify(ctx context.Context, data []byte) (*payment.NotifyResult, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("paypal: empty notification data")
	}

	// Mock: in production, verify webhook signature via PayPal's
	// /v1/notifications/verify-webhook-signature API.
	return &payment.NotifyResult{
		Status: "success",
		Sign:   "mock_signature",
	}, nil
}
