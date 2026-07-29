package payment

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

// QQPayGateway implements the Gateway interface for QQ Wallet (QQ钱包).
type QQPayGateway struct {
	MchID     string
	APIKey    string
	AppID     string
	NotifyURL string
}

// QQPayConfig is the JSON schema stored in payment_gateways.config.
type QQPayConfig struct {
	MchID     string `json:"mch_id"`
	APIKey    string `json:"api_key"`
	AppID     string `json:"app_id"`
	NotifyURL string `json:"notify_url"`
}

// NewQQPayGateway creates a QQ Pay gateway from a JSON config string.
func NewQQPayGateway(configJSON string) (*QQPayGateway, error) {
	var cfg QQPayConfig
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		return nil, fmt.Errorf("qqpay config parse error: %w", err)
	}
	if cfg.MchID == "" || cfg.APIKey == "" {
		return nil, fmt.Errorf("qqpay mch_id and api_key are required")
	}
	return &QQPayGateway{
		MchID:     cfg.MchID,
		APIKey:    cfg.APIKey,
		AppID:     cfg.AppID,
		NotifyURL: cfg.NotifyURL,
	}, nil
}

func (g *QQPayGateway) Name() string { return "qqpay" }

// qqPaySign generates an MD5 signature for QQ Pay API requests.
// Parameters are sorted by key, concatenated as key=value&key=value, appended with &key=APIKey, then MD5 hashed.
func (g *QQPayGateway) qqPaySign(params map[string]string) string {
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
	buf.WriteString("&key=")
	buf.WriteString(g.APIKey)

	hash := md5.Sum([]byte(buf.String()))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

// buildXMLRequest builds an XML request body from parameters.
func buildXMLRequest(params map[string]string) string {
	var buf strings.Builder
	buf.WriteString("<xml>")
	for k, v := range params {
		buf.WriteString("<" + k + ">" + v + "</" + k + ">")
	}
	buf.WriteString("</xml>")
	return buf.String()
}

// xmlResponse represents a QQ Pay XML response.
type xmlResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`
	PrepayID   string `xml:"prepay_id"`
	CodeURL    string `xml:"code_url"`
	TradeType  string `xml:"trade_type"`
	MchID      string `xml:"mch_id"`
	AppID      string `xml:"appid"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	DeviceInfo string `xml:"device_info"`
}

// CreatePayment creates a QQ Pay unified order and returns a QR code URL.
func (g *QQPayGateway) CreatePayment(ctx context.Context, param *PaymentParam) (*PaymentResult, error) {
	if param.Amount <= 0 {
		return nil, fmt.Errorf("qqpay: amount must be positive")
	}

	totalFee := int(param.Amount * 100) // Convert to fen
	nonceStr := generateNonce(32)

	params := map[string]string{
		"mch_id":       g.MchID,
		"out_trade_no": param.OrderNo,
		"total_fee":    fmt.Sprintf("%d", totalFee),
		"body":         param.Subject,
		"trade_type":   "NATIVE",
		"notify_url":   g.NotifyURL,
		"nonce_str":    nonceStr,
	}
	if g.AppID != "" {
		params["appid"] = g.AppID
	}

	params["sign"] = g.qqPaySign(params)

	xmlBody := buildXMLRequest(params)
	reqURL := "https://qpay.qq.com/cgi-bin/pay/qpay_unified_order.cgi"

	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, bytes.NewBufferString(xmlBody))
	if err != nil {
		return nil, fmt.Errorf("qqpay: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/xml")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("qqpay: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("qqpay: read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("qqpay: error status=%d body=%s", resp.StatusCode, string(body))
	}

	var result xmlResponse
	if err := xml.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("qqpay: parse XML response: %w", err)
	}

	if result.ReturnCode != "SUCCESS" {
		return nil, fmt.Errorf("qqpay: return error: %s", result.ReturnMsg)
	}
	if result.ResultCode != "SUCCESS" {
		return nil, fmt.Errorf("qqpay: result error: %s - %s", result.ErrCode, result.ErrCodeDes)
	}

	codeURL := result.CodeURL
	if codeURL == "" {
		return nil, fmt.Errorf("qqpay: no code_url returned for NATIVE trade type")
	}

	return &PaymentResult{
		TradeNo:   param.OrderNo,
		QrcodeURL: codeURL,
		RawData: map[string]interface{}{
			"prepay_id":  result.PrepayID,
			"trade_type": result.TradeType,
			"total_fee":  totalFee,
			"gateway":    "qqpay",
		},
	}, nil
}

// QueryPayment queries the status of a QQ Pay order.
func (g *QQPayGateway) QueryPayment(ctx context.Context, tradeNo string) (*QueryResult, error) {
	if tradeNo == "" {
		return nil, fmt.Errorf("qqpay: out_trade_no is required")
	}

	nonceStr := generateNonce(32)
	params := map[string]string{
		"mch_id":       g.MchID,
		"out_trade_no": tradeNo,
		"nonce_str":    nonceStr,
	}
	params["sign"] = g.qqPaySign(params)

	xmlBody := buildXMLRequest(params)
	reqURL := "https://qpay.qq.com/cgi-bin/pay/qpay_order_query.cgi"

	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, bytes.NewBufferString(xmlBody))
	if err != nil {
		return nil, fmt.Errorf("qqpay: create query request: %w", err)
	}
	req.Header.Set("Content-Type", "application/xml")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("qqpay: query request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("qqpay: read query response: %w", err)
	}

	var result struct {
		ReturnCode  string `xml:"return_code"`
		ReturnMsg   string `xml:"return_msg"`
		ResultCode  string `xml:"result_code"`
		TradeState  string `xml:"trade_state"`
		TotalFee    int    `xml:"total_fee"`
		TransactionID string `xml:"transaction_id"`
		TimeEnd     string `xml:"time_end"`
	}
	if err := xml.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("qqpay: parse query response: %w", err)
	}

	if result.ReturnCode != "SUCCESS" {
		return nil, fmt.Errorf("qqpay: query return error: %s", result.ReturnMsg)
	}

	status := "pending"
	switch result.TradeState {
	case "SUCCESS":
		status = "success"
	case "CLOSED", "REVOKED", "PAYERROR":
		status = "failed"
	}

	var paidAt *time.Time
	if status == "success" && result.TimeEnd != "" {
		t, err := time.Parse("20060102150405", result.TimeEnd)
		if err == nil {
			paidAt = &t
		}
	}

	return &QueryResult{
		TradeNo: tradeNo,
		Status:  status,
		Amount:  float64(result.TotalFee) / 100,
		PaidAt:  paidAt,
	}, nil
}

// Refund processes a refund through QQ Pay.
func (g *QQPayGateway) Refund(ctx context.Context, param *RefundParam) (*RefundResult, error) {
	if param.Amount <= 0 {
		return nil, fmt.Errorf("qqpay: refund amount must be positive")
	}

	totalFee := int(param.Amount * 100)
	nonceStr := generateNonce(32)

	params := map[string]string{
		"mch_id":        g.MchID,
		"out_trade_no":  param.TradeNo,
		"out_refund_no": param.RefundNo,
		"total_fee":     fmt.Sprintf("%d", totalFee),
		"refund_fee":    fmt.Sprintf("%d", totalFee),
		"nonce_str":     nonceStr,
	}
	params["sign"] = g.qqPaySign(params)

	xmlBody := buildXMLRequest(params)
	reqURL := "https://qpay.qq.com/cgi-bin/pay/qpay_refund.cgi"

	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, bytes.NewBufferString(xmlBody))
	if err != nil {
		return nil, fmt.Errorf("qqpay: create refund request: %w", err)
	}
	req.Header.Set("Content-Type", "application/xml")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("qqpay: refund request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("qqpay: read refund response: %w", err)
	}

	var result struct {
		ReturnCode string `xml:"return_code"`
		ReturnMsg  string `xml:"return_msg"`
		ResultCode string `xml:"result_code"`
		RefundID   string `xml:"refund_id"`
	}
	if err := xml.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("qqpay: parse refund response: %w", err)
	}

	if result.ReturnCode != "SUCCESS" {
		return nil, fmt.Errorf("qqpay: refund return error: %s", result.ReturnMsg)
	}

	status := "processing"
	if result.ResultCode == "SUCCESS" {
		status = "success"
	}

	return &RefundResult{
		RefundNo: result.RefundID,
		Status:   status,
	}, nil
}

// ParseNotify parses a QQ Pay payment callback notification.
func (g *QQPayGateway) ParseNotify(ctx context.Context, data []byte) (*NotifyResult, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("qqpay: empty notification data")
	}

	var notify struct {
		ReturnCode    string `xml:"return_code"`
		ResultCode    string `xml:"result_code"`
		OutTradeNo    string `xml:"out_trade_no"`
		TotalFee      int    `xml:"total_fee"`
		TransactionID string `xml:"transaction_id"`
		TimeEnd       string `xml:"time_end"`
		Sign          string `xml:"sign"`
	}
	if err := xml.Unmarshal(data, &notify); err != nil {
		return nil, fmt.Errorf("qqpay: parse notify: %w", err)
	}

	if notify.ReturnCode != "SUCCESS" {
		return nil, fmt.Errorf("qqpay: notify return error")
	}

	status := "pending"
	if notify.ResultCode == "SUCCESS" {
		status = "success"
	}

	return &NotifyResult{
		TradeNo: notify.TransactionID,
		OrderNo: notify.OutTradeNo,
		Amount:  float64(notify.TotalFee) / 100,
		Status:  status,
	}, nil
}
