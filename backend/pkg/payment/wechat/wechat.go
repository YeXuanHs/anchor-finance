package wechat

import (
	"context"
	"fmt"
	"time"

	"anchorfinance/pkg/payment"
)

func init() {
	payment.Register(&Gateway{})
}

// Gateway implements the payment.Gateway interface for WeChat Pay.
type Gateway struct {
	AppID     string
	MchID     string
	APIKey    string
	NotifyURL string
}

// Name returns the gateway identifier.
func (g *Gateway) Name() string {
	return "wechat"
}

// CreatePayment creates a WeChat Pay native payment and returns a QR code URL.
func (g *Gateway) CreatePayment(ctx context.Context, param *payment.PaymentParam) (*payment.PaymentResult, error) {
	if param.Amount <= 0 {
		return nil, fmt.Errorf("wechat: amount must be positive")
	}

	tradeNo := param.TradeNo
	if tradeNo == "" {
		tradeNo = fmt.Sprintf("WX_%s_%d", param.OrderNo, time.Now().UnixNano())
	}

	// In production, this would call WeChat Pay's unified order API
	// to obtain a code_url for native QR code payment.
	qrcodeURL := fmt.Sprintf(
		"weixin://wxpay/bizpayurl?trade_no=%s&amount=%d",
		tradeNo, int(param.Amount*100),
	)

	return &payment.PaymentResult{
		TradeNo:   tradeNo,
		QrcodeURL: qrcodeURL,
		RawData: map[string]interface{}{
			"gateway":    "wechat",
			"order_no":   param.OrderNo,
			"trade_no":   tradeNo,
			"amount":     param.Amount,
			"created_at": time.Now().Format(time.RFC3339),
		},
	}, nil
}

// QueryPayment queries the status of a WeChat Pay payment.
func (g *Gateway) QueryPayment(ctx context.Context, tradeNo string) (*payment.QueryResult, error) {
	if tradeNo == "" {
		return nil, fmt.Errorf("wechat: trade_no is required")
	}

	// Mock: in production, call /pay/orderquery API.
	return &payment.QueryResult{
		TradeNo: tradeNo,
		Status:  "pending",
		Amount:  0,
		PaidAt:  nil,
	}, nil
}

// Refund processes a refund through WeChat Pay.
func (g *Gateway) Refund(ctx context.Context, param *payment.RefundParam) (*payment.RefundResult, error) {
	if param.Amount <= 0 {
		return nil, fmt.Errorf("wechat: refund amount must be positive")
	}

	// Mock: in production, call /secapi/pay/refund API.
	return &payment.RefundResult{
		RefundNo: param.RefundNo,
		Status:   "processing",
	}, nil
}

// ParseNotify parses a WeChat Pay async notification callback.
func (g *Gateway) ParseNotify(ctx context.Context, data []byte) (*payment.NotifyResult, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("wechat: empty notification data")
	}

	// Mock implementation: in production, parse XML notification,
	// verify HMAC-SHA256 signature, and check result_code.
	return &payment.NotifyResult{
		Status: "success",
		Sign:   "mock_signature",
	}, nil
}
