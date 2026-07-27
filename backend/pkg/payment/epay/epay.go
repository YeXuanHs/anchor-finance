package epay

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/anchor-finance/backend/pkg/payment"
)

func init() {
	payment.Register(&Gateway{})
}

// Gateway implements the payment.Gateway interface for Epay (易支付).
//
// Epay is a self-hosted payment aggregation platform. Each instance has its own
// API base URL, so ApiURL is configurable per deployment.
//
// API Reference:
//   - Create payment: POST {api_url}/submit.php
//   - Parameters: pid, type (alipay/wxpay/qq), out_trade_no, notify_url,
//     return_url, name, money, sign (MD5), sign_type ("MD5")
//   - Sign: sort all non-empty params except sign/sign_type by key, join as
//     k1=v1&k2=v2, append KEY, then MD5 (32-char lowercase hex)
//   - Notify callback: GET/POST with trade_status, sign for verification
//   - Success response body is plain text "success"
type Gateway struct {
	ApiURL    string // Base URL of the Epay instance, e.g. "https://pay.example.com"
	PID       string // Merchant ID
	Key       string // Secret key for signing
	NotifyURL string
	ReturnURL string
	Timeout   time.Duration
}

// Name returns the gateway identifier.
func (g *Gateway) Name() string {
	return "epay"
}

// httpClient returns the configured HTTP client.
func (g *Gateway) httpClient() *http.Client {
	timeout := g.Timeout
	if timeout == 0 {
		timeout = 15 * time.Second
	}
	return &http.Client{Timeout: timeout}
}

// CreatePayment submits a payment request to Epay and returns the payment URL.
func (g *Gateway) CreatePayment(ctx context.Context, param *payment.PaymentParam) (*payment.PaymentResult, error) {
	if param.Amount <= 0 {
		return nil, fmt.Errorf("epay: amount must be positive")
	}
	if g.ApiURL == "" || g.PID == "" || g.Key == "" {
		return nil, fmt.Errorf("epay: api_url, pid, and key are required")
	}

	payType := inferPayType(param)
	tradeNo := param.OrderNo

	params := map[string]string{
		"pid":          g.PID,
		"type":         payType,
		"out_trade_no": tradeNo,
		"notify_url":   g.NotifyURL,
		"name":         truncate(param.Subject, 64),
		"money":        formatAmount(param.Amount),
	}
	if g.ReturnURL != "" {
		params["return_url"] = g.ReturnURL
	}

	params["sign"] = g.sign(params)
	params["sign_type"] = "MD5"

	// Build the payment URL with query parameters
	submitURL := strings.TrimRight(g.ApiURL, "/") + "/submit.php"
	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}

	payURL := submitURL + "?" + form.Encode()

	return &payment.PaymentResult{
		TradeNo: tradeNo,
		PayURL:  payURL,
		RawData: map[string]interface{}{
			"gateway":    "epay",
			"pid":        g.PID,
			"type":       payType,
			"order_no":   param.OrderNo,
			"trade_no":   tradeNo,
			"amount":     param.Amount,
			"created_at": time.Now().Format(time.RFC3339),
		},
	}, nil
}

// QueryPayment queries Epay for payment status via GET /api.php?act=order.
func (g *Gateway) QueryPayment(ctx context.Context, tradeNo string) (*payment.QueryResult, error) {
	if tradeNo == "" {
		return nil, fmt.Errorf("epay: trade_no is required")
	}
	if g.ApiURL == "" || g.PID == "" || g.Key == "" {
		return nil, fmt.Errorf("epay: api_url, pid, and key are required")
	}

	params := map[string]string{
		"act":          "order",
		"pid":          g.PID,
		"out_trade_no": tradeNo,
	}
	params["sign"] = g.sign(params)
	params["sign_type"] = "MD5"

	queryURL := strings.TrimRight(g.ApiURL, "/") + "/api.php"
	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, queryURL+"?"+form.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("epay: failed to build query request: %w", err)
	}

	resp, err := g.httpClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("epay: query request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("epay: failed to read query response: %w", err)
	}

	var result epayQueryResponse
	if err := json.Unmarshal(body, &result); err != nil {
		// Epay may return plain text error
		return nil, fmt.Errorf("epay: invalid query response: %s", string(body))
	}
	if result.Code != 1 {
		return nil, fmt.Errorf("epay: query error: %s", result.Msg)
	}

	status := mapTradeStatus(result.TradeStatus)

	qr := &payment.QueryResult{
		TradeNo: tradeNo,
		Status:  status,
		Amount:  parseFloat(result.Money),
	}
	if status == "success" && result.TradeNo != "" {
		now := time.Now()
		qr.PaidAt = &now
	}
	return qr, nil
}

// Refund is not supported by Epay.
func (g *Gateway) Refund(ctx context.Context, param *payment.RefundParam) (*payment.RefundResult, error) {
	return nil, fmt.Errorf("epay: refund is not supported")
}

// ParseNotify parses and verifies an Epay async notification callback.
//
// Epay sends notifications as GET/POST form data with trade_status, sign, etc.
// The handler should respond with plain text "success" if verification passes.
func (g *Gateway) ParseNotify(ctx context.Context, data []byte) (*payment.NotifyResult, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("epay: empty notification data")
	}

	values, err := url.ParseQuery(string(data))
	if err != nil {
		return nil, fmt.Errorf("epay: failed to parse notification: %w", err)
	}

	params := make(map[string]string)
	for k, v := range values {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}

	receivedSign := params["sign"]
	if receivedSign == "" {
		return nil, fmt.Errorf("epay: notification missing sign")
	}

	// Verify signature
	expectedSign := g.sign(params)
	if receivedSign != expectedSign {
		return nil, fmt.Errorf("epay: signature verification failed (expected %s, got %s)", expectedSign, receivedSign)
	}

	status := mapTradeStatus(values.Get("trade_status"))

	return &payment.NotifyResult{
		TradeNo: values.Get("trade_no"),
		OrderNo: values.Get("out_trade_no"),
		Amount:  parseFloat(values.Get("money")),
		Status:  status,
		Sign:    receivedSign,
	}, nil
}

// NotifySuccessBody returns the plain text body that Epay expects in response
// to a successful notification callback.
func NotifySuccessBody() string {
	return "success"
}

// sign computes the Epay MD5 signature.
//
// Algorithm:
//  1. Collect all non-empty params where key != "sign" and key != "sign_type"
//  2. Sort keys alphabetically
//  3. Join as key1=value1&key2=value2
//  4. Append the secret KEY directly (no separator)
//  5. MD5 hash, 32-char lowercase hex
func (g *Gateway) sign(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k, v := range params {
		if v != "" && k != "sign" && k != "sign_type" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var buf strings.Builder
	for i, k := range keys {
		if i > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(params[k])
	}
	buf.WriteString(g.Key)

	h := md5.Sum([]byte(buf.String()))
	return fmt.Sprintf("%x", h)
}

// inferPayType determines the Epay payment type from the payment parameters.
// Defaults to "alipay" if no specific type can be determined.
func inferPayType(param *payment.PaymentParam) string {
	// If the caller specifies a pay type via description, use it.
	// Supported: alipay, wxpay, qq
	lower := strings.ToLower(param.Description)
	for _, t := range []string{"alipay", "wxpay", "qq"} {
		if strings.Contains(lower, t) {
			return t
		}
	}
	return "alipay"
}

// mapTradeStatus maps Epay trade_status to the internal status enum.
func mapTradeStatus(tradeStatus string) string {
	switch tradeStatus {
	case "TRADE_SUCCESS":
		return "success"
	case "TRADE_FINISHED":
		return "success"
	case "TRADE_CLOSED":
		return "closed"
	default:
		return "pending"
	}
}

// epayQueryResponse represents the JSON response from Epay order query API.
type epayQueryResponse struct {
	Code         int    `json:"code"`
	Msg          string `json:"msg"`
	TradeNo      string `json:"trade_no"`
	OutTradeNo   string `json:"out_trade_no"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	Money        string `json:"money"`
	TradeStatus  string `json:"trade_status"`
}

// formatAmount formats a float64 amount to 2 decimal places.
func formatAmount(amount float64) string {
	return strconv.FormatFloat(amount, 'f', 2, 64)
}

// parseFloat safely parses a string to float64, returning 0 on error.
func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

// truncate truncates a string to the given max byte length.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}

// generateNonceStr creates a random alphanumeric string of the given length.
// Exported for use in testing.
func generateNonceStr(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}
