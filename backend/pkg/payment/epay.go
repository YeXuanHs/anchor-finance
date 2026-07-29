package payment

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

// EpayGateway implements the Gateway interface for 易支付 (Epay).
// Epay is an open-source payment gateway widely used in Chinese hosting/IDC industry.
type EpayGateway struct {
	PID       string
	Key       string
	APIURL    string
	NotifyURL string
	ReturnURL string
}

// EpayConfig is the JSON schema stored in payment_gateways.config.
type EpayConfig struct {
	PID       string `json:"pid"`
	Key       string `json:"key"`
	APIURL    string `json:"api_url"`
	NotifyURL string `json:"notify_url"`
	ReturnURL string `json:"return_url"`
}

// NewEpayGateway creates an Epay gateway from a JSON config string.
func NewEpayGateway(configJSON string) (*EpayGateway, error) {
	var cfg EpayConfig
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		return nil, fmt.Errorf("epay config parse error: %w", err)
	}
	if cfg.PID == "" || cfg.Key == "" {
		return nil, fmt.Errorf("epay pid and key are required")
	}
	if cfg.APIURL == "" {
		return nil, fmt.Errorf("epay api_url is required")
	}
	// Remove trailing slash
	cfg.APIURL = strings.TrimRight(cfg.APIURL, "/")
	return &EpayGateway{
		PID:       cfg.PID,
		Key:       cfg.Key,
		APIURL:    cfg.APIURL,
		NotifyURL: cfg.NotifyURL,
		ReturnURL: cfg.ReturnURL,
	}, nil
}

func (g *EpayGateway) Name() string { return "epay" }

func (g *EpayGateway) CreatePayment(ctx context.Context, param *PaymentParam) (*PaymentResult, error) {
	params := map[string]string{
		"pid":          g.PID,
		"type":         "alipay", // Default to alipay; can be overridden via Extra
		"out_trade_no": param.OrderNo,
		"notify_url":   g.NotifyURL,
		"return_url":   g.ReturnURL,
		"name":         param.Subject,
		"money":        fmt.Sprintf("%.2f", param.Amount),
	}

	// Allow selecting specific payment method via Extra
	if param.ClientIP != "" {
		params["clientip"] = param.ClientIP
	}

	params["sign"] = g.sign(params)
	params["sign_type"] = "MD5"

	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}

	payURL := fmt.Sprintf("%s/submit.php?%s", g.APIURL, values.Encode())

	return &PaymentResult{
		TradeNo: param.OrderNo,
		PayURL:  payURL,
	}, nil
}

func (g *EpayGateway) QueryPayment(ctx context.Context, tradeNo string) (*QueryResult, error) {
	params := map[string]string{
		"act":          "order",
		"pid":          g.PID,
		"out_trade_no": tradeNo,
	}
	params["sign"] = g.sign(params)
	params["sign_type"] = "MD5"

	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}

	queryURL := fmt.Sprintf("%s/api.php?%s", g.APIURL, values.Encode())
	resp, err := http.Get(queryURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Code    int    `json:"code"`
		Status  string `json:"status"`
		TradeNo string `json:"trade_no"`
	}
	json.Unmarshal(body, &result)

	status := "pending"
	switch result.Status {
	case "TRADE_SUCCESS":
		status = "success"
	case "TRADE_CLOSED":
		status = "closed"
	}

	return &QueryResult{TradeNo: tradeNo, Status: status}, nil
}

func (g *EpayGateway) Refund(ctx context.Context, param *RefundParam) (*RefundResult, error) {
	return &RefundResult{
		RefundNo: param.RefundNo,
		Status:   "pending", // Epay typically requires manual refund via admin panel
	}, nil
}

func (g *EpayGateway) ParseNotify(ctx context.Context, data []byte) (*NotifyResult, error) {
	form, err := url.ParseQuery(string(data))
	if err != nil {
		return nil, fmt.Errorf("parse epay notify: %w", err)
	}

	tradeNo := form.Get("out_trade_no")
	tradeStatus := form.Get("trade_status")

	status := "pending"
	switch tradeStatus {
	case "TRADE_SUCCESS":
		status = "success"
	case "TRADE_CLOSED":
		status = "closed"
	}

	return &NotifyResult{
		TradeNo: tradeNo,
		OrderNo: tradeNo,
		Status:  status,
		Sign:    form.Get("sign"),
	}, nil
}

func (g *EpayGateway) sign(params map[string]string) string {
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
			buf.WriteString("&")
		}
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(params[k])
	}
	buf.WriteString(g.Key)

	hash := md5.Sum([]byte(buf.String()))
	return strings.ToLower(hex.EncodeToString(hash[:]))
}
