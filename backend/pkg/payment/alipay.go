package payment

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// AliPayConfig 支付宝官方配置
type AliPayConfig struct {
	AppID           string `json:"app_id"`            // 应用ID
	PrivateKey      string `json:"private_key"`       // 应用私钥
	AlipayPublicKey string `json:"alipay_public_key"` // 支付宝公钥
	ProductPC       bool   `json:"product_pc"`        // 是否支持PC支付
	ProductWAP      bool   `json:"product_wap"`       // 是否支持WAP支付
	ProductQR       bool   `json:"product_qr"`        // 是否支持扫码支付
}

// AliPayGateway 支付宝官方接口
type AliPayGateway struct {
	config *AliPayConfig
}

// NewAliPayGateway 创建支付宝官方实例
func NewAliPayGateway(configJSON string) (*AliPayGateway, error) {
	var config AliPayConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, fmt.Errorf("invalid alipay config: %w", err)
	}
	if config.AppID == "" || config.PrivateKey == "" || config.AlipayPublicKey == "" {
		return nil, fmt.Errorf("alipay config missing required fields (app_id, private_key, alipay_public_key)")
	}
	return &AliPayGateway{config: &config}, nil
}

func (g *AliPayGateway) Name() string { return GatewayAliPay }

func (g *AliPayGateway) CreatePayment(ctx context.Context, param *PaymentParam) (*PaymentResult, error) {
	// 根据配置选择支付产品
	if g.config.ProductPC {
		return g.pagePay(param)
	} else if g.config.ProductQR {
		return g.qrPay(param)
	} else if g.config.ProductWAP {
		return g.wapPay(param)
	}
	return nil, fmt.Errorf("no payment product configured")
}

// pagePay PC网页支付
func (g *AliPayGateway) pagePay(param *PaymentParam) (*PaymentResult, error) {
	bizContent := map[string]interface{}{
		"out_trade_no": param.OrderNo,
		"total_amount": fmt.Sprintf("%.2f", param.Amount),
		"subject":      param.Subject,
		"product_code": "FAST_INSTANT_TRADE_PAY",
	}

	params := g.buildCommonParams("alipay.trade.page.pay")
	params["biz_content"] = marshalJSON(bizContent)
	params["return_url"] = param.ReturnURL
	params["notify_url"] = param.NotifyURL

	// 生成签名
	params["sign"] = g.sign(params)

	// 构建跳转URL
	gatewayURL := "https://openapi.alipay.com/gateway.do"
	query := url.Values{}
	for k, v := range params {
		query.Set(k, v)
	}
	payURL := fmt.Sprintf("%s?%s", gatewayURL, query.Encode())

	return &PaymentResult{
		Type:    "jump",
		Data:    payURL,
		OrderNo: param.OrderNo,
	}, nil
}

// wapPay 手机网页支付
func (g *AliPayGateway) wapPay(param *PaymentParam) (*PaymentResult, error) {
	bizContent := map[string]interface{}{
		"out_trade_no": param.OrderNo,
		"total_amount": fmt.Sprintf("%.2f", param.Amount),
		"subject":      param.Subject,
		"product_code": "QUICK_WAP_WAY",
	}

	params := g.buildCommonParams("alipay.trade.wap.pay")
	params["biz_content"] = marshalJSON(bizContent)
	params["return_url"] = param.ReturnURL
	params["notify_url"] = param.NotifyURL

	params["sign"] = g.sign(params)

	gatewayURL := "https://openapi.alipay.com/gateway.do"
	query := url.Values{}
	for k, v := range params {
		query.Set(k, v)
	}
	payURL := fmt.Sprintf("%s?%s", gatewayURL, query.Encode())

	return &PaymentResult{
		Type:    "jump",
		Data:    payURL,
		OrderNo: param.OrderNo,
	}, nil
}

// qrPay 扫码支付
func (g *AliPayGateway) qrPay(param *PaymentParam) (*PaymentResult, error) {
	bizContent := map[string]interface{}{
		"out_trade_no": param.OrderNo,
		"total_amount": fmt.Sprintf("%.2f", param.Amount),
		"subject":      param.Subject,
		"product_code": "FACE_TO_FACE_PAYMENT",
	}

	params := g.buildCommonParams("alipay.trade.precreate")
	params["biz_content"] = marshalJSON(bizContent)
	params["notify_url"] = param.NotifyURL

	params["sign"] = g.sign(params)

	// 调用支付宝API
	gatewayURL := "https://openapi.alipay.com/gateway.do"
	result, err := g.doRequest(gatewayURL, params)
	if err != nil {
		return nil, fmt.Errorf("alipay precreate failed: %w", err)
	}

	// 解析响应
	response, ok := result["alipay_trade_precreate_response"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid alipay response")
	}

	code, _ := response["code"].(string)
	if code != "10000" {
		msg, _ := response["sub_msg"].(string)
		return nil, fmt.Errorf("alipay error: %s", msg)
	}

	qrCode, _ := response["qr_code"].(string)
	if qrCode == "" {
		return nil, fmt.Errorf("alipay returned empty qr code")
	}

	return &PaymentResult{
		Type:    "url",
		Data:    qrCode,
		OrderNo: param.OrderNo,
	}, nil
}

func (g *AliPayGateway) VerifyNotification(ctx context.Context, data map[string]string) (*NotificationResult, error) {
	// 验证签名
	sign := data["sign"]
	signType := data["sign_type"]
	delete(data, "sign")
	delete(data, "sign_type")

	if !g.verifySign(data, sign, signType) {
		return nil, fmt.Errorf("invalid alipay signature")
	}

	// 解析业务参数
	tradeStatus := data["trade_status"]
	status := "pending"
	if tradeStatus == "TRADE_SUCCESS" || tradeStatus == "TRADE_FINISHED" {
		status = "success"
	}

	var amount float64
	fmt.Sscanf(data["total_amount"], "%f", &amount)

	return &NotificationResult{
		OrderNo: data["out_trade_no"],
		TradeNo: data["trade_no"],
		Amount:  amount,
		Status:  status,
	}, nil
}

// buildCommonParams 构建公共参数
func (g *AliPayGateway) buildCommonParams(method string) map[string]string {
	return map[string]string{
		"app_id":      g.config.AppID,
		"method":      method,
		"format":      "JSON",
		"charset":     "utf-8",
		"sign_type":   "RSA2",
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		"version":     "1.0",
	}
}

// sign RSA2签名
func (g *AliPayGateway) sign(params map[string]string) string {
	// 按key排序
	keys := make([]string, 0, len(params))
	for k, v := range params {
		if v != "" && k != "sign" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	// 拼接待签名字符串
	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, params[k]))
	}
	signContent := strings.Join(parts, "&")

	// 解析私钥
	block, _ := pem.Decode([]byte(g.config.PrivateKey))
	if block == nil {
		return ""
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return ""
	}

	// SHA256WithRSA签名
	hash := sha256.Sum256([]byte(signContent))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), crypto.SHA256, hash[:])
	if err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(signature)
}

// verifySign 验证RSA2签名
func (g *AliPayGateway) verifySign(params map[string]string, sign string, signType string) bool {
	// 按key排序
	keys := make([]string, 0, len(params))
	for k, v := range params {
		if v != "" && k != "sign" && k != "sign_type" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	// 拼接待验签字符串
	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, params[k]))
	}
	signContent := strings.Join(parts, "&")

	// 解析支付宝公钥
	block, _ := pem.Decode([]byte(g.config.AlipayPublicKey))
	if block == nil {
		return false
	}
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false
	}

	// 验证签名
	signBytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false
	}

	hash := sha256.Sum256([]byte(signContent))
	err = rsa.VerifyPKCS1v15(publicKey.(*rsa.PublicKey), crypto.SHA256, hash[:], signBytes)
	return err == nil
}

// doRequest 调用支付宝API
func (g *AliPayGateway) doRequest(gatewayURL string, params map[string]string) (map[string]interface{}, error) {
	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}

	resp, err := http.Post(gatewayURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("invalid response: %s", string(body))
	}

	return result, nil
}

func marshalJSON(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}
