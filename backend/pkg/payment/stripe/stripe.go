package stripe

import (
	"context"
	"fmt"
	"time"

	"github.com/anchor-finance/backend/pkg/payment"
)

func init() {
	payment.Register(&Gateway{})
}

// Gateway implements the payment.Gateway interface for Stripe.
type Gateway struct {
	SecretKey      string
	WebhookSecret  string
	Currency       string // default: "usd"
}

// Name returns the gateway identifier.
func (g *Gateway) Name() string {
	return "stripe"
}

// currency returns the configured currency or defaults to "usd".
func (g *Gateway) currency() string {
	if g.Currency != "" {
		return g.Currency
	}
	return "usd"
}

// CreatePayment creates a Stripe Checkout Session and returns the session URL.
func (g *Gateway) CreatePayment(ctx context.Context, param *payment.PaymentParam) (*payment.PaymentResult, error) {
	if param.Amount <= 0 {
		return nil, fmt.Errorf("stripe: amount must be positive")
	}

	tradeNo := param.TradeNo
	if tradeNo == "" {
		tradeNo = fmt.Sprintf("STR_%s_%d", param.OrderNo, time.Now().UnixNano())
	}

	// In production, this would create a Stripe Checkout Session via
	// POST /v1/checkout/sessions with line_items and success/cancel URLs.
	checkoutURL := fmt.Sprintf(
		"https://checkout.stripe.com/pay/%s",
		tradeNo,
	)

	return &payment.PaymentResult{
		TradeNo: tradeNo,
		PayURL:  checkoutURL,
		RawData: map[string]interface{}{
			"gateway":     "stripe",
			"order_no":    param.OrderNo,
			"trade_no":    tradeNo,
			"amount":      param.Amount,
			"currency":    g.currency(),
			"created_at":  time.Now().Format(time.RFC3339),
		},
	}, nil
}

// QueryPayment queries the status of a Stripe payment intent.
func (g *Gateway) QueryPayment(ctx context.Context, tradeNo string) (*payment.QueryResult, error) {
	if tradeNo == "" {
		return nil, fmt.Errorf("stripe: trade_no is required")
	}

	// Mock: in production, retrieve PaymentIntent via /v1/payment_intents/{id}.
	return &payment.QueryResult{
		TradeNo: tradeNo,
		Status:  "pending",
		Amount:  0,
		PaidAt:  nil,
	}, nil
}

// Refund processes a refund through Stripe.
func (g *Gateway) Refund(ctx context.Context, param *payment.RefundParam) (*payment.RefundResult, error) {
	if param.Amount <= 0 {
		return nil, fmt.Errorf("stripe: refund amount must be positive")
	}

	// Mock: in production, call /v1/refunds with payment_intent and amount.
	return &payment.RefundResult{
		RefundNo: param.RefundNo,
		Status:   "processing",
	}, nil
}

// ParseNotify parses a Stripe webhook event.
func (g *Gateway) ParseNotify(ctx context.Context, data []byte) (*payment.NotifyResult, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("stripe: empty notification data")
	}

	// Mock: in production, verify the webhook signature using
	// stripe.ConstructEvent with the endpoint signing secret,
	// then handle checkout.session.completed or payment_intent.succeeded events.
	return &payment.NotifyResult{
		Status: "success",
		Sign:   "mock_signature",
	}, nil
}
