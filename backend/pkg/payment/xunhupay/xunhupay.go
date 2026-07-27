package xunhupay

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

	"anchorfinance/pkg/payment"
)

const (
	paymentURL = "https://api.xunhupay.com/payment/do.html"
	apiVersion = "1.1"
)

func init() {
	payment.Register(&Gateway{})
}

// Gateway implements the payment.Gateway interface for XunhuPay (虎皮椒).
//
// API Reference:
//   - POST https://api.xunhupay.com/payment/do.html
//   - Sign: sort all non-empty params alphabetically, join as k1=v1&k2=v2,
//     append APPSECRET, then MD5 (32-char lowercase hex).
//   - Response JSON contains orderid, url_qrcode (PC QR), and url (deprecated).
type Gateway struct {
	AppID     string
	AppSecret string
	NotifyURL string
	ReturnURL string
	Timeout   time.Duration // HTTP request timeout; default 15s
}

// Name returns the gateway identifier.
func (g *Gateway) Name() string {
	return "xunhupay"
}

// httpClient returns the configured HTTP client.
func (g *Gateway) httpClient() *http.Client {
	timeout := g.Timeout
	if timeout == 0 {
		timeout = 15 * time.Second
	}
	return &http.Client{Timeout: timeout}
}

// CreatePayment submits a payment request to XunhuPay and returns the payment URL / QR code.
func (g *Gateway) CreatePayment(ctx context.Context, param *payment.PaymentParam) (*payment.PaymentResult, error) {
	if param.Amount <= 0 {
		return nil, fmt.Errorf("xunhupay: amount must be positive")
	}
	if g.AppID == "" || g.AppSecret == "" {
		return nil, fmt.Errorf("xunhupay: appid and appsecret are required")
	}

	params := map[string]string{
		"version":        apiVersion,
		"appid":          g.AppID,
		"trade_order_id": param.OrderNo,
		"total_fee":      formatAmount(param.Amount),
		"title":          truncate(param.Subject, 32),
		"time":           strconv.FormatInt(time.Now().Unix(), 10),
		"notify_url":     g.NotifyURL,
		"nonce_str":      generateNonceStr(32),
	}

	if g.ReturnURL != "" {
		params["return_url"] = g.ReturnURL
	}
	if param.Description != "" {
		params["attach"] = truncate(param.Description, 255)
	}

	params["hash"] = g.sign(params)

	// Build form data
	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, paymentURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("xunhupay: failed to build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := g.httpClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("xunhupay: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("xunhupay: failed to read response: %w", err)
	}

	var result xunhuResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("xunhupay: invalid response JSON: %w (body: %s)", err, string(body))
	}
	if result.ErrCode != 0 {
		return nil, fmt.Errorf("xunhupay: error %d: %s", result.ErrCode, result.ErrMsg)
	}
	if result.URLQrcode == "" && result.URL == "" {
		return nil, fmt.Errorf("xunhupay: response missing payment URL (orderid=%s)", result.OrderID)
	}

	return &payment.PaymentResult{
		TradeNo:   result.OrderID,
		PayURL:    result.URL,
		QrcodeURL: result.URLQrcode,
		RawData: map[string]interface{}{
			"gateway":      "xunhupay",
			"orderid":      result.OrderID,
			"url_qrcode":   result.URLQrcode,
			"url":          result.URL,
			"order_no":     param.OrderNo,
			"amount":       param.Amount,
			"created_at":   time.Now().Format(time.RFC3339),
		},
	}, nil
}

// QueryPayment queries XunhuPay for payment status.
// XunhuPay does not provide a standard query API, so this relies on the notify callback.
func (g *Gateway) QueryPayment(ctx context.Context, tradeNo string) (*payment.QueryResult, error) {
	if tradeNo == "" {
		return nil, fmt.Errorf("xunhupay: trade_no is required")
	}
	return &payment.QueryResult{
		TradeNo: tradeNo,
		Status:  "pending",
	}, nil
}

// Refund is not supported by XunhuPay.
func (g *Gateway) Refund(ctx context.Context, param *payment.RefundParam) (*payment.RefundResult, error) {
	return nil, fmt.Errorf("xunhupay: refund is not supported")
}

// ParseNotify parses and verifies an XunhuPay async notification callback.
//
// The callback is POSTed as form data containing the same parameters used to
// create the payment, plus the hash field for verification.
func (g *Gateway) ParseNotify(ctx context.Context, data []byte) (*payment.NotifyResult, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("xunhupay: empty notification data")
	}

	values, err := url.ParseQuery(string(data))
	if err != nil {
		return nil, fmt.Errorf("xunhupay: failed to parse notification: %w", err)
	}

	params := make(map[string]string)
	for k, v := range values {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}

	receivedHash := params["hash"]
	if receivedHash == "" {
		return nil, fmt.Errorf("xunhupay: notification missing hash")
	}

	// Verify signature: recompute and compare
	delete(params, "hash")
	expectedHash := g.sign(params)
	if receivedHash != expectedHash {
		return nil, fmt.Errorf("xunhupay: signature verification failed (expected %s, got %s)", expectedHash, receivedHash)
	}

	status := "pending"
	switch values.Get("trade_status") {
	case "TRADE_SUCCESS":
		status = "success"
	case "TRADE_CLOSED":
		status = "closed"
	}

	return &payment.NotifyResult{
		TradeNo: values.Get("trade_order_id"),
		OrderNo: values.Get("trade_order_id"),
		Amount:  parseFloat(values.Get("total_fee")),
		Status:  status,
		Sign:    receivedHash,
	}, nil
}

// sign computes the XunhuPay MD5 signature.
//
// Algorithm:
//  1. Collect all non-empty params where key != "hash"
//  2. Sort keys alphabetically
//  3. Join as key1=value1&key2=value2
//  4. Append APPSECRET directly (no separator)
//  5. MD5 hash, 32-char lowercase hex
func (g *Gateway) sign(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k, v := range params {
		if v != "" && k != "hash" {
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
	buf.WriteString(g.AppSecret)

	h := md5.Sum([]byte(buf.String()))
	return fmt.Sprintf("%x", h)
}

// xunhuResponse represents the JSON response from XunhuPay API.
type xunhuResponse struct {
	OrderID   string `json:"orderid"`
	URL       string `json:"url"`
	URLQrcode string `json:"url_qrcode"`
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
}

// generateNonceStr creates a random alphanumeric string of the given length.
func generateNonceStr(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
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

// truncate truncates a string to the given max length.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
