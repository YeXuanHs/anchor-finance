package payment

import (
	"context"
	"time"
)

// Gateway defines the interface that all payment gateways must implement.
type Gateway interface {
	Name() string
	CreatePayment(ctx context.Context, param *PaymentParam) (*PaymentResult, error)
	QueryPayment(ctx context.Context, tradeNo string) (*QueryResult, error)
	Refund(ctx context.Context, param *RefundParam) (*RefundResult, error)
	ParseNotify(ctx context.Context, data []byte) (*NotifyResult, error)
}

// WebhookVerifier is an optional interface for gateways that support
// cryptographic webhook signature verification. Handlers should prefer
// VerifyAndParseNotify over ParseNotify when the gateway implements this.
type WebhookVerifier interface {
	VerifyAndParseNotify(ctx context.Context, data []byte, headers map[string]string) (*NotifyResult, error)
}

// PaymentParam contains the parameters needed to create a payment.
type PaymentParam struct {
	OrderNo     string
	TradeNo     string
	Amount      float64
	Subject     string
	Description string
	ClientIP    string
}

// PaymentResult contains the result of creating a payment.
type PaymentResult struct {
	TradeNo   string
	PayURL    string
	QrcodeURL string
	RawData   map[string]interface{}
}

// QueryResult contains the status of a payment query.
type QueryResult struct {
	TradeNo string
	Status  string // pending, success, failed, closed
	Amount  float64
	PaidAt  *time.Time
}

// RefundParam contains the parameters needed to process a refund.
type RefundParam struct {
	TradeNo  string
	RefundNo string
	Amount   float64
	Reason   string
}

// RefundResult contains the result of a refund operation.
type RefundResult struct {
	RefundNo string
	Status   string
}

// NotifyResult contains the parsed result of an async payment notification.
type NotifyResult struct {
	TradeNo string
	OrderNo string
	Amount  float64
	Status  string
	Sign    string
}

// Global gateway registry.
var gateways = make(map[string]Gateway)

// Register adds a payment gateway to the registry.
func Register(gw Gateway) {
	gateways[gw.Name()] = gw
}

// Get retrieves a payment gateway by name.
func Get(name string) (Gateway, bool) {
	gw, ok := gateways[name]
	return gw, ok
}

// GetAll returns all registered payment gateways.
func GetAll() map[string]Gateway {
	return gateways
}
