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

// XunhuPayGateway implements the Gateway interface for 虎皮虾 (XunhuPay).
// 虎皮虾 is a Chinese payment aggregator supporting multiple payment methods.
type XunhuPayGateway struct {
	AppID     string
	AppSecret string
	NotifyURL string
	ReturnURL string
	GatewayURL string
}

// XunhuPayConfig is the JSON schema stored in payment_gateways.config.
type XunhuPayConfig struct {
	AppID      string `json:"app_id"`
	AppSecret  string `json:"app_secret"`
	NotifyURL  string `json:"notify_url"`
	ReturnURL  string `json:"return_url"`
	GatewayURL string `json:"gateway_url"`
}

// NewXunhuPayGateway creates a XunhuPay gateway from a JSON config string.
func NewXunhuPayGateway(configJSON string) (*XunhuPayGateway, error) {
	var cfg XunhuPayConfig
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		return nil, fmt.Errorf("xunhupay config parse error: %w", err)
	}
	if cfg.AppID == "" || cfg.AppSecret == "" {
		return nil, fmt.Errorf("xunhupay app_id and app_secret are required")
	}
	if cfg.GatewayURL == "" {
		cfg.GatewayURL = "https://api.xunhupay.com/payment/do"
	}
	return &XunhuPayGateway{
		AppID:      cfg.AppID,
		AppSecret:  cfg.AppSecret,
		NotifyURL:  cfg.NotifyURL,
		ReturnURL:  cfg.ReturnURL,
		GatewayURL: cfg.GatewayURL,
	}, nil
}

func (g *XunhuPayGateway) Name() string { return "xunhupay" }

func (g *XunhuPayGateway) CreatePayment(ctx context.Context, param *PaymentParam) (*PaymentResult, error) {
	params := map[string]string{
		"hash":         g.AppID,
		"trade_order_id": param.OrderNo,
		"total_fee":    fmt.Sprintf("%.2f", param.Amount),
		"title":        param.Subject,
		"notify_url":   g.NotifyURL,
		"return_url":   g.ReturnURL,
		"nonce":        generateNonce(16),
	}

	params["hash"] = g.sign(params)

	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}

	resp, err := http.PostForm(g.GatewayURL, values)
	if err != nil {
		return nil, fmt.Errorf("xunhupay request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			URL     string `json:"url"`
			Qrcode  string `json:"qrcode"`
			OrderID string `json:"order_id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse xunhupay response: %w", err)
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("xunhupay error: %s", result.Msg)
	}

	return &PaymentResult{
		TradeNo:   param.OrderNo,
		PayURL:    result.Data.URL,
		QrcodeURL: result.Data.Qrcode,
	}, nil
}

func (g *XunhuPayGateway) QueryPayment(ctx context.Context, tradeNo string) (*QueryResult, error) {
	params := map[string]string{
		"hash":           g.AppID,
		"trade_order_id": tradeNo,
		"nonce":          generateNonce(16),
	}
	params["hash"] = g.sign(params)

	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}

	resp, err := http.PostForm("https://api.xunhupay.com/payment/query", values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Code int `json:"code"`
		Data struct {
			Status string `json:"status"`
		} `json:"data"`
	}
	json.Unmarshal(body, &result)

	status := "pending"
	switch result.Data.Status {
	case "1", "TRADE_SUCCESS":
		status = "success"
	case "2", "TRADE_CLOSED":
		status = "closed"
	}

	return &QueryResult{TradeNo: tradeNo, Status: status}, nil
}

func (g *XunhuPayGateway) Refund(ctx context.Context, param *RefundParam) (*RefundResult, error) {
	return &RefundResult{
		RefundNo: param.RefundNo,
		Status:   "pending", // XunhuPay typically requires manual refund
	}, nil
}

func (g *XunhuPayGateway) ParseNotify(ctx context.Context, data []byte) (*NotifyResult, error) {
	form, err := url.ParseQuery(string(data))
	if err != nil {
		return nil, fmt.Errorf("parse xunhupay notify: %w", err)
	}

	orderID := form.Get("trade_order_id")
	statusStr := form.Get("trade_status")

	status := "pending"
	if statusStr == "TRADE_SUCCESS" {
		status = "success"
	} else if statusStr == "TRADE_CLOSED" {
		status = "closed"
	}

	return &NotifyResult{
		TradeNo: orderID,
		OrderNo: orderID,
		Status:  status,
		Sign:    form.Get("hash"),
	}, nil
}

func (g *XunhuPayGateway) sign(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
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
	buf.WriteString(g.AppSecret)

	hash := md5.Sum([]byte(buf.String()))
	return hex.EncodeToString(hash[:])
}
