package payment

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// AlipayGateway implements the Gateway interface for Alipay (支付宝).
type AlipayGateway struct {
	AppID      string
	PrivateKey string
	PublicKey  string
	NotifyURL  string
	ReturnURL  string
	GatewayURL string
	SignType   string // RSA2
	httpClient *http.Client
}

// AlipayConfig is the JSON schema stored in payment_gateways.config.
type AlipayConfig struct {
	AppID      string `json:"app_id"`
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
	NotifyURL  string `json:"notify_url"`
	ReturnURL  string `json:"return_url"`
	GatewayURL string `json:"gateway_url"`
	SignType   string `json:"sign_type"`
}

// NewAlipayGateway creates an Alipay gateway from a JSON config string.
func NewAlipayGateway(configJSON string) (*AlipayGateway, error) {
	var cfg AlipayConfig
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		return nil, fmt.Errorf("alipay config parse error: %w", err)
	}
	if cfg.AppID == "" {
		return nil, fmt.Errorf("alipay app_id is required")
	}
	if cfg.GatewayURL == "" {
		cfg.GatewayURL = "https://openapi.alipay.com/gateway.do"
	}
	if cfg.SignType == "" {
		cfg.SignType = "RSA2"
	}
	return &AlipayGateway{
		AppID:      cfg.AppID,
		PrivateKey: cfg.PrivateKey,
		PublicKey:  cfg.PublicKey,
		NotifyURL:  cfg.NotifyURL,
		ReturnURL:  cfg.ReturnURL,
		GatewayURL: cfg.GatewayURL,
		SignType:   cfg.SignType,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}, nil
}

func (g *AlipayGateway) Name() string { return "alipay" }

func (g *AlipayGateway) CreatePayment(ctx context.Context, param *PaymentParam) (*PaymentResult, error) {
	bizContent := map[string]interface{}{
		"out_trade_no":    param.OrderNo,
		"total_amount":    fmt.Sprintf("%.2f", param.Amount),
		"subject":         param.Subject,
		"product_code":    "QUICK_MSECURITY_PAY",
		"timeout_express": "30m",
	}

	bizJSON, err := json.Marshal(bizContent)
	if err != nil {
		return nil, fmt.Errorf("marshal biz_content: %w", err)
	}

	params := map[string]string{
		"app_id":      g.AppID,
		"method":      "alipay.trade.page.pay",
		"charset":     "utf-8",
		"sign_type":   g.SignType,
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		"version":     "1.0",
		"biz_content": string(bizJSON),
		"notify_url":  g.NotifyURL,
		"return_url":  g.ReturnURL,
	}

	sign, err := g.sign(params)
	if err != nil {
		return nil, fmt.Errorf("sign error: %w", err)
	}
	params["sign"] = sign

	values := url.Values{}
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		values.Set(k, params[k])
	}
	payURL := g.GatewayURL + "?" + values.Encode()

	return &PaymentResult{
		TradeNo: param.OrderNo,
		PayURL:  payURL,
	}, nil
}

// QueryPayment queries the payment status via alipay.trade.query API.
func (g *AlipayGateway) QueryPayment(ctx context.Context, tradeNo string) (*QueryResult, error) {
	params := map[string]string{
		"app_id":      g.AppID,
		"method":      "alipay.trade.query",
		"charset":     "utf-8",
		"sign_type":   g.SignType,
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		"version":     "1.0",
		"biz_content": fmt.Sprintf(`{"out_trade_no":"%s"}`, tradeNo),
	}

	sign, err := g.sign(params)
	if err != nil {
		return nil, fmt.Errorf("sign error: %w", err)
	}
	params["sign"] = sign

	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}

	resp, err := g.httpClient.PostForm(g.GatewayURL, values)
	if err != nil {
		return nil, fmt.Errorf("alipay query request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	// Parse with RawMessage to extract the response content for signature verification.
	var rawResult struct {
		AlipayTradeQueryResponse json.RawMessage `json:"alipay_trade_query_response"`
		Sign                     string          `json:"sign"`
	}
	if err := json.Unmarshal(body, &rawResult); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	// Verify response signature against the raw response JSON.
	if rawResult.Sign != "" {
		if err := g.verifySign(string(rawResult.AlipayTradeQueryResponse), rawResult.Sign); err != nil {
			log.Printf("alipay query response signature verification failed: %v", err)
			return nil, fmt.Errorf("response signature verification failed: %w", err)
		}
	}

	var respData struct {
		Code        string `json:"code"`
		Msg         string `json:"msg"`
		SubCode     string `json:"sub_code"`
		SubMsg      string `json:"sub_msg"`
		TradeNo     string `json:"trade_no"`
		TradeStatus string `json:"trade_status"`
		TotalAmount string `json:"total_amount"`
		SendPayDate string `json:"send_pay_date"`
	}
	if err := json.Unmarshal(rawResult.AlipayTradeQueryResponse, &respData); err != nil {
		return nil, fmt.Errorf("parse query response: %w", err)
	}

	if respData.Code != "10000" {
		return nil, fmt.Errorf("alipay query error: %s - %s (%s)", respData.Code, respData.Msg, respData.SubMsg)
	}

	status := "pending"
	switch respData.TradeStatus {
	case "TRADE_SUCCESS", "TRADE_FINISHED":
		status = "success"
	case "TRADE_CLOSED":
		status = "closed"
	case "WAIT_BUYER_PAY":
		status = "pending"
	}

	var amount float64
	fmt.Sscanf(respData.TotalAmount, "%f", &amount)

	return &QueryResult{
		TradeNo: tradeNo,
		Status:  status,
		Amount:  amount,
	}, nil
}

// Refund processes a refund via alipay.trade.refund API.
func (g *AlipayGateway) Refund(ctx context.Context, param *RefundParam) (*RefundResult, error) {
	bizContent := map[string]interface{}{
		"out_trade_no":   param.TradeNo,
		"out_request_no": param.RefundNo,
		"refund_amount":  fmt.Sprintf("%.2f", param.Amount),
	}
	if param.Reason != "" {
		bizContent["refund_reason"] = param.Reason
	}

	bizJSON, err := json.Marshal(bizContent)
	if err != nil {
		return nil, fmt.Errorf("marshal biz_content: %w", err)
	}

	params := map[string]string{
		"app_id":      g.AppID,
		"method":      "alipay.trade.refund",
		"charset":     "utf-8",
		"sign_type":   g.SignType,
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		"version":     "1.0",
		"biz_content": string(bizJSON),
	}

	sign, err := g.sign(params)
	if err != nil {
		return nil, fmt.Errorf("sign error: %w", err)
	}
	params["sign"] = sign

	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}

	resp, err := g.httpClient.PostForm(g.GatewayURL, values)
	if err != nil {
		return nil, fmt.Errorf("alipay refund request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	// Parse with RawMessage to extract the response content for signature verification.
	var rawResult struct {
		AlipayTradeRefundResponse json.RawMessage `json:"alipay_trade_refund_response"`
		Sign                      string          `json:"sign"`
	}
	if err := json.Unmarshal(body, &rawResult); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	// Verify response signature against the raw response JSON.
	if rawResult.Sign != "" {
		if err := g.verifySign(string(rawResult.AlipayTradeRefundResponse), rawResult.Sign); err != nil {
			log.Printf("alipay refund response signature verification failed: %v", err)
			return nil, fmt.Errorf("response signature verification failed: %w", err)
		}
	}

	var respData struct {
		Code    string `json:"code"`
		Msg     string `json:"msg"`
		SubCode string `json:"sub_code"`
		SubMsg  string `json:"sub_msg"`
	}
	if err := json.Unmarshal(rawResult.AlipayTradeRefundResponse, &respData); err != nil {
		return nil, fmt.Errorf("parse refund response: %w", err)
	}

	if respData.Code != "10000" {
		return nil, fmt.Errorf("alipay refund error: %s - %s (%s)", respData.Code, respData.Msg, respData.SubMsg)
	}

	return &RefundResult{
		RefundNo: param.RefundNo,
		Status:   "success",
	}, nil
}

// ParseNotify verifies and parses an Alipay async payment notification.
// It verifies the RSA2 signature against all form params (excluding sign, sign_type, and empty values).
func (g *AlipayGateway) ParseNotify(ctx context.Context, data []byte) (*NotifyResult, error) {
	form, err := url.ParseQuery(string(data))
	if err != nil {
		return nil, fmt.Errorf("parse notify form: %w", err)
	}

	sign := form.Get("sign")
	if sign == "" {
		return nil, fmt.Errorf("missing sign in notification")
	}

	// Build verification string: all params except sign, sign_type, and empty values.
	params := make(map[string]string)
	for k, v := range form {
		if k == "sign" || k == "sign_type" {
			continue
		}
		if len(v) > 0 && v[0] != "" {
			params[k] = v[0]
		}
	}

	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
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

	// Verify RSA2 signature using the Alipay public key.
	if err := g.verifySign(buf.String(), sign); err != nil {
		log.Printf("alipay notify signature verification failed: %v", err)
		return nil, fmt.Errorf("notification signature verification failed: %w", err)
	}

	// Extract notification fields.
	tradeNo := form.Get("out_trade_no")
	alipayTradeNo := form.Get("trade_no")
	tradeStatus := form.Get("trade_status")
	totalAmount := form.Get("total_amount")

	status := "pending"
	switch tradeStatus {
	case "TRADE_SUCCESS", "TRADE_FINISHED":
		status = "success"
	case "TRADE_CLOSED":
		status = "closed"
	}

	var amount float64
	fmt.Sscanf(totalAmount, "%f", &amount)

	return &NotifyResult{
		TradeNo: alipayTradeNo,
		OrderNo: tradeNo,
		Amount:  amount,
		Status:  status,
		Sign:    sign,
	}, nil
}

// verifySign verifies an RSA2 (SHA256WithRSA) signature against content using the Alipay public key.
func (g *AlipayGateway) verifySign(content string, sign string) error {
	block, _ := pem.Decode([]byte(g.PublicKey))
	if block == nil {
		return fmt.Errorf("failed to parse public key PEM")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("parse public key: %w", err)
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("not an RSA public key")
	}

	sigBytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return fmt.Errorf("decode signature: %w", err)
	}

	hashed := sha256.Sum256([]byte(content))
	if err := rsa.VerifyPKCS1v15(rsaPub, crypto.SHA256, hashed[:], sigBytes); err != nil {
		return fmt.Errorf("verify signature failed: %w", err)
	}

	return nil
}

// sign signs the params using RSA2 (SHA256WithRSA) with the app private key.
// All non-empty params are sorted alphabetically, joined with &, then signed.
func (g *AlipayGateway) sign(params map[string]string) (string, error) {
	keys := make([]string, 0, len(params))
	for k, v := range params {
		if v != "" {
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

	block, _ := pem.Decode([]byte(g.PrivateKey))
	if block == nil {
		return "", fmt.Errorf("failed to parse private key PEM")
	}

	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		parsedKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return "", fmt.Errorf("parse private key: %w", err)
		}
	}

	privKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("not an RSA private key")
	}

	hashed := sha256.Sum256([]byte(buf.String()))
	sig, err := rsa.SignPKCS1v15(nil, privKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", fmt.Errorf("sign error: %w", err)
	}

	return base64.StdEncoding.EncodeToString(sig), nil
}
