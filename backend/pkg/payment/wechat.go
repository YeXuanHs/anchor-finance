package payment

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

// WechatGateway implements the Gateway interface for WeChat Pay (微信支付).
type WechatGateway struct {
	AppID      string
	MchID      string
	APIKey     string
	NotifyURL  string
	GatewayURL string
	httpClient *http.Client
	tlsClient  *http.Client // For mTLS (required by secapi endpoints like refund)
}

// WechatConfig is the JSON schema stored in payment_gateways.config.
type WechatConfig struct {
	AppID      string `json:"app_id"`
	MchID      string `json:"mch_id"`
	APIKey     string `json:"api_key"`
	NotifyURL  string `json:"notify_url"`
	GatewayURL string `json:"gateway_url"`
	CertPath   string `json:"cert_path"`    // Client certificate file for mTLS
	KeyPath    string `json:"key_path"`     // Client private key file for mTLS
	RootCAPath string `json:"root_ca_path"` // Root CA certificate (optional)
}

// NewWechatGateway creates a WeChat Pay gateway from a JSON config string.
func NewWechatGateway(configJSON string) (*WechatGateway, error) {
	var cfg WechatConfig
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		return nil, fmt.Errorf("wechat config parse error: %w", err)
	}
	if cfg.AppID == "" || cfg.MchID == "" || cfg.APIKey == "" {
		return nil, fmt.Errorf("wechat app_id, mch_id, and api_key are required")
	}
	if cfg.GatewayURL == "" {
		cfg.GatewayURL = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	}

	gw := &WechatGateway{
		AppID:      cfg.AppID,
		MchID:      cfg.MchID,
		APIKey:     cfg.APIKey,
		NotifyURL:  cfg.NotifyURL,
		GatewayURL: cfg.GatewayURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	// Initialize mTLS client if cert paths are provided.
	if cfg.CertPath != "" && cfg.KeyPath != "" {
		tlsClient, err := newTLSClient(cfg.CertPath, cfg.KeyPath, cfg.RootCAPath)
		if err != nil {
			return nil, fmt.Errorf("create TLS client: %w", err)
		}
		gw.tlsClient = tlsClient
	}

	return gw, nil
}

// newTLSClient creates an HTTP client with mTLS certificate support.
func newTLSClient(certPath, keyPath, rootCAPath string) (*http.Client, error) {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("load client cert: %w", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	if rootCAPath != "" {
		caCert, err := os.ReadFile(rootCAPath)
		if err != nil {
			return nil, fmt.Errorf("read root CA: %w", err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
	}

	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}, nil
}

func (g *WechatGateway) Name() string { return "wechat" }

// WechatXMLRequest is the XML structure for unified order request.
type WechatXMLRequest struct {
	XMLName        xml.Name `xml:"xml"`
	AppID          string   `xml:"appid"`
	MchID          string   `xml:"mch_id"`
	NonceStr       string   `xml:"nonce_str"`
	Sign           string   `xml:"sign"`
	Body           string   `xml:"body"`
	OutTradeNo     string   `xml:"out_trade_no"`
	TotalFee       int      `xml:"total_fee"`
	SpbillCreateIP string   `xml:"spbill_create_ip"`
	NotifyURL      string   `xml:"notify_url"`
	TradeType      string   `xml:"trade_type"`
}

// WechatXMLResponse is the XML structure for unified order response.
type WechatXMLResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	ResultCode string `xml:"result_code"`
	PrepayID   string `xml:"prepay_id"`
	CodeURL    string   `xml:"code_url"`
}

func (g *WechatGateway) CreatePayment(ctx context.Context, param *PaymentParam) (*PaymentResult, error) {
	nonceStr := generateNonce(32)
	totalFee := int(param.Amount * 100) // Convert to fen

	params := map[string]string{
		"appid":            g.AppID,
		"mch_id":           g.MchID,
		"nonce_str":        nonceStr,
		"body":             param.Subject,
		"out_trade_no":     param.OrderNo,
		"total_fee":        fmt.Sprintf("%d", totalFee),
		"spbill_create_ip": param.ClientIP,
		"notify_url":       g.NotifyURL,
		"trade_type":       "NATIVE",
	}

	sign := g.sign(params)
	params["sign"] = sign

	req := WechatXMLRequest{
		AppID:          g.AppID,
		MchID:          g.MchID,
		NonceStr:       nonceStr,
		Sign:           sign,
		Body:           param.Subject,
		OutTradeNo:     param.OrderNo,
		TotalFee:       totalFee,
		SpbillCreateIP: param.ClientIP,
		NotifyURL:      g.NotifyURL,
		TradeType:      "NATIVE",
	}

	xmlData, err := xml.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal xml: %w", err)
	}

	resp, err := g.httpClient.Post(g.GatewayURL, "application/xml", strings.NewReader(string(xmlData)))
	if err != nil {
		return nil, fmt.Errorf("wechat api request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var wxResp WechatXMLResponse
	if err := xml.Unmarshal(body, &wxResp); err != nil {
		return nil, fmt.Errorf("parse wechat response: %w", err)
	}

	if wxResp.ReturnCode != "SUCCESS" || wxResp.ResultCode != "SUCCESS" {
		return nil, fmt.Errorf("wechat error: %s - %s", wxResp.ReturnCode, wxResp.ReturnMsg)
	}

	return &PaymentResult{
		TradeNo:   param.OrderNo,
		QrcodeURL: wxResp.CodeURL, // QR code URL for NATIVE pay
	}, nil
}

// QueryPayment queries the payment status via pay/orderquery API.
func (g *WechatGateway) QueryPayment(ctx context.Context, tradeNo string) (*QueryResult, error) {
	params := map[string]string{
		"appid":        g.AppID,
		"mch_id":       g.MchID,
		"out_trade_no": tradeNo,
		"nonce_str":    generateNonce(32),
	}
	params["sign"] = g.sign(params)

	type queryReq struct {
		XMLName    xml.Name `xml:"xml"`
		AppID      string   `xml:"appid"`
		MchID      string   `xml:"mch_id"`
		OutTradeNo string   `xml:"out_trade_no"`
		NonceStr   string   `xml:"nonce_str"`
		Sign       string   `xml:"sign"`
	}

	xmlData, err := xml.Marshal(queryReq{
		AppID: g.AppID, MchID: g.MchID, OutTradeNo: tradeNo,
		NonceStr: params["nonce_str"], Sign: params["sign"],
	})
	if err != nil {
		return nil, fmt.Errorf("marshal xml: %w", err)
	}

	resp, err := g.httpClient.Post(
		"https://api.mch.weixin.qq.com/pay/orderquery",
		"application/xml",
		strings.NewReader(string(xmlData)),
	)
	if err != nil {
		return nil, fmt.Errorf("wechat query request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	// Parse response XML into a map for signature verification.
	respParams, err := parseXMLToMap(body)
	if err != nil {
		return nil, fmt.Errorf("parse response xml: %w", err)
	}

	// Verify response signature.
	if err := g.verifySign(respParams); err != nil {
		log.Printf("wechat query response signature verification failed: %v", err)
		return nil, fmt.Errorf("response signature verification failed: %w", err)
	}

	if respParams["return_code"] != "SUCCESS" {
		return nil, fmt.Errorf("wechat query error: %s", respParams["return_msg"])
	}
	if respParams["result_code"] != "SUCCESS" {
		return nil, fmt.Errorf("wechat query error: %s - %s", respParams["err_code"], respParams["err_code_des"])
	}

	status := "pending"
	switch respParams["trade_state"] {
	case "SUCCESS":
		status = "success"
	case "CLOSED", "REVOKED", "PAYERROR":
		status = "failed"
	case "NOTPAY":
		status = "pending"
	}

	var amount float64
	if totalFee, ok := respParams["total_fee"]; ok {
		var fee int
		fmt.Sscanf(totalFee, "%d", &fee)
		amount = float64(fee) / 100.0
	}

	return &QueryResult{
		TradeNo: tradeNo,
		Status:  status,
		Amount:  amount,
	}, nil
}

// Refund processes a refund via secapi/pay/refund API (requires mTLS client certificate).
func (g *WechatGateway) Refund(ctx context.Context, param *RefundParam) (*RefundResult, error) {
	params := map[string]string{
		"appid":         g.AppID,
		"mch_id":        g.MchID,
		"nonce_str":     generateNonce(32),
		"out_trade_no":  param.TradeNo,
		"out_refund_no": param.RefundNo,
		"total_fee":     fmt.Sprintf("%d", int(param.Amount*100)),
		"refund_fee":    fmt.Sprintf("%d", int(param.Amount*100)),
	}
	if param.Reason != "" {
		params["refund_desc"] = param.Reason
	}
	params["sign"] = g.sign(params)

	type refundReq struct {
		XMLName     xml.Name `xml:"xml"`
		AppID       string   `xml:"appid"`
		MchID       string   `xml:"mch_id"`
		NonceStr    string   `xml:"nonce_str"`
		Sign        string   `xml:"sign"`
		OutTradeNo  string   `xml:"out_trade_no"`
		OutRefundNo string   `xml:"out_refund_no"`
		TotalFee    string   `xml:"total_fee"`
		RefundFee   string   `xml:"refund_fee"`
		RefundDesc  string   `xml:"refund_desc,omitempty"`
	}

	xmlData, err := xml.Marshal(refundReq{
		AppID: g.AppID, MchID: g.MchID, NonceStr: params["nonce_str"],
		Sign: params["sign"], OutTradeNo: param.TradeNo,
		OutRefundNo: param.RefundNo, TotalFee: params["total_fee"],
		RefundFee: params["refund_fee"], RefundDesc: params["refund_desc"],
	})
	if err != nil {
		return nil, fmt.Errorf("marshal xml: %w", err)
	}

	// WeChat secapi endpoints require mTLS with client certificate.
	client := g.httpClient
	if g.tlsClient != nil {
		client = g.tlsClient
	} else {
		log.Println("WARNING: WeChat refund API requires mTLS client certificate; request may fail without cert_path/key_path configured")
	}

	resp, err := client.Post(
		"https://api.mch.weixin.qq.com/secapi/pay/refund",
		"application/xml",
		strings.NewReader(string(xmlData)),
	)
	if err != nil {
		return nil, fmt.Errorf("wechat refund request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	// Parse response XML into a map for signature verification.
	respParams, err := parseXMLToMap(body)
	if err != nil {
		return nil, fmt.Errorf("parse response xml: %w", err)
	}

	// Verify response signature.
	if err := g.verifySign(respParams); err != nil {
		log.Printf("wechat refund response signature verification failed: %v", err)
		return nil, fmt.Errorf("response signature verification failed: %w", err)
	}

	if respParams["return_code"] != "SUCCESS" {
		return nil, fmt.Errorf("wechat refund error: %s", respParams["return_msg"])
	}
	if respParams["result_code"] != "SUCCESS" {
		return nil, fmt.Errorf("wechat refund error: %s - %s", respParams["err_code"], respParams["err_code_des"])
	}

	return &RefundResult{
		RefundNo: param.RefundNo,
		Status:   "success",
	}, nil
}

// ParseNotify verifies and parses a WeChat Pay async payment notification.
// It extracts all XML fields, validates the MD5 signature, and returns a NotifyResult.
func (g *WechatGateway) ParseNotify(ctx context.Context, data []byte) (*NotifyResult, error) {
	// Parse XML into a map for signature verification.
	params, err := parseXMLToMap(data)
	if err != nil {
		return nil, fmt.Errorf("parse notify xml: %w", err)
	}

	// Verify MD5 signature.
	if err := g.verifySign(params); err != nil {
		log.Printf("wechat notify signature verification failed: %v", err)
		return nil, fmt.Errorf("notification signature verification failed: %w", err)
	}

	status := "pending"
	if params["return_code"] == "SUCCESS" && params["result_code"] == "SUCCESS" {
		status = "success"
	}

	var amount float64
	if totalFee, ok := params["total_fee"]; ok {
		var fee int
		fmt.Sscanf(totalFee, "%d", &fee)
		amount = float64(fee) / 100.0
	}

	return &NotifyResult{
		TradeNo: params["transaction_id"],
		OrderNo: params["out_trade_no"],
		Amount:  amount,
		Status:  status,
		Sign:    params["sign"],
	}, nil
}

// sign computes an MD5 signature for the given params.
// All non-empty params are sorted alphabetically, joined with & and &key=<api_key>,
// then MD5 hashed and returned as uppercase hex.
func (g *WechatGateway) sign(params map[string]string) string {
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
	buf.WriteString("&key=")
	buf.WriteString(g.APIKey)

	hash := md5.Sum([]byte(buf.String()))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

// verifySign verifies the MD5 signature contained in the response params.
// It recomputes the signature from all params except "sign" and compares.
func (g *WechatGateway) verifySign(params map[string]string) error {
	receivedSign := params["sign"]
	if receivedSign == "" {
		return fmt.Errorf("missing sign in response")
	}

	// Build a copy without the sign field for verification.
	verifyParams := make(map[string]string)
	for k, v := range params {
		if k != "sign" {
			verifyParams[k] = v
		}
	}

	computedSign := g.sign(verifyParams)
	if computedSign != receivedSign {
		return fmt.Errorf("sign mismatch: expected %s, got %s", computedSign, receivedSign)
	}

	return nil
}

// parseXMLToMap parses generic XML key-value pairs into a map[string]string.
// This is used to extract all fields from WeChat API responses for signature verification.
func parseXMLToMap(data []byte) (map[string]string, error) {
	type xmlNode struct {
		XMLName xml.Name
		Content string    `xml:",chardata"`
		Nodes   []xmlNode `xml:",any"`
	}

	var node xmlNode
	if err := xml.Unmarshal(data, &node); err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, child := range node.Nodes {
		result[child.XMLName.Local] = strings.TrimSpace(child.Content)
	}

	return result, nil
}

// generateNonce generates a random hex string of the given length.
func generateNonce(length int) string {
	b := make([]byte, (length+1)/2)
	if _, err := rand.Read(b); err != nil {
		// Fallback to time-based (not crypto-safe but functional).
		for i := range b {
			b[i] = byte(time.Now().UnixNano() >> (uint(i) * 8))
		}
	}
	return hex.EncodeToString(b)[:length]
}

// WechatGatewayH5 implements H5 payment (for mobile browser redirect).
type WechatGatewayH5 struct {
	*WechatGateway
}

func (g *WechatGatewayH5) Name() string { return "wechat_h5" }

func (g *WechatGatewayH5) CreatePayment(ctx context.Context, param *PaymentParam) (*PaymentResult, error) {
	nonceStr := generateNonce(32)
	totalFee := int(param.Amount * 100)

	params := map[string]string{
		"appid":            g.AppID,
		"mch_id":           g.MchID,
		"nonce_str":        nonceStr,
		"body":             param.Subject,
		"out_trade_no":     param.OrderNo,
		"total_fee":        fmt.Sprintf("%d", totalFee),
		"spbill_create_ip": param.ClientIP,
		"notify_url":       g.NotifyURL,
		"trade_type":       "MWEB",
	}
	params["sign"] = g.sign(params)

	req := WechatXMLRequest{
		AppID: g.AppID, MchID: g.MchID, NonceStr: nonceStr,
		Sign: params["sign"], Body: param.Subject, OutTradeNo: param.OrderNo,
		TotalFee: totalFee, SpbillCreateIP: param.ClientIP,
		NotifyURL: g.NotifyURL, TradeType: "MWEB",
	}

	xmlData, _ := xml.Marshal(req)
	resp, err := g.httpClient.Post(g.GatewayURL, "application/xml", strings.NewReader(string(xmlData)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var wxResp WechatXMLResponse
	xml.Unmarshal(body, &wxResp)

	if wxResp.ReturnCode != "SUCCESS" || wxResp.ResultCode != "SUCCESS" {
		return nil, fmt.Errorf("wechat h5 error: %s", wxResp.ReturnMsg)
	}

	return &PaymentResult{
		TradeNo: param.OrderNo,
		PayURL:  wxResp.PrepayID, // For H5, this contains the MWEB URL
	}, nil
}
