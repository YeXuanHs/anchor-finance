package payment

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// WxPayConfig 微信支付官方配置
type WxPayConfig struct {
	AppID     string `json:"app_id"`     // 公众号/小程序ID
	MchID     string `json:"mch_id"`     // 商户号
	APIKey    string `json:"api_key"`    // API密钥
	AppSecret string `json:"app_secret"` // 应用密钥
	ProductNative bool `json:"product_native"` // 是否支持扫码支付
	ProductJSAPI  bool `json:"product_jsapi"`  // 是否支持JSAPI支付
	ProductWAP    bool `json:"product_wap"`    // 是否支持H5支付
}

// WxPayGateway 微信支付官方接口
type WxPayGateway struct {
	config *WxPayConfig
}

// WxPayRequest 微信支付请求
type WxPayRequest struct {
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
	ProductID      string   `xml:"product_id,omitempty"`
	SceneInfo      string   `xml:"scene_info,omitempty"`
}

// WxPayResponse 微信支付响应
type WxPayResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	ResultCode string `xml:"result_code"`
	PrepayID   string `xml:"prepay_id"`
	CodeURL    string `xml:"code_url"`
	MwebURL    string `xml:"mweb_url"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`
}

// WxPayNotifyResult 微信支付回调结果
type WxPayNotifyResult struct {
	ReturnCode    string `xml:"return_code"`
	ResultCode    string `xml:"result_code"`
	OutTradeNo    string `xml:"out_trade_no"`
	TransactionID string `xml:"transaction_id"`
	TotalFee      int    `xml:"total_fee"`
	TimeEnd       string `xml:"time_end"`
}

// NewWxPayGateway 创建微信支付官方实例
func NewWxPayGateway(configJSON string) (*WxPayGateway, error) {
	var config WxPayConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, fmt.Errorf("invalid wxpay config: %w", err)
	}
	if config.AppID == "" || config.MchID == "" || config.APIKey == "" {
		return nil, fmt.Errorf("wxpay config missing required fields (app_id, mch_id, api_key)")
	}
	return &WxPayGateway{config: &config}, nil
}

func (g *WxPayGateway) Name() string { return GatewayWxPay }

func (g *WxPayGateway) CreatePayment(ctx context.Context, param *PaymentParam) (*PaymentResult, error) {
	// 根据配置选择支付产品
	if g.config.ProductNative {
		return g.nativePay(param)
	} else if g.config.ProductWAP {
		return g.wapPay(param)
	} else if g.config.ProductJSAPI {
		return g.jsapiPay(param)
	}
	return nil, fmt.Errorf("no payment product configured")
}

// nativePay 扫码支付
func (g *WxPayGateway) nativePay(param *PaymentParam) (*PaymentResult, error) {
	totalFee := int(param.Amount * 100) // 转换为分

	req := &WxPayRequest{
		AppID:          g.config.AppID,
		MchID:          g.config.MchID,
		NonceStr:       generateNonceStr(),
		Body:           param.Subject,
		OutTradeNo:     param.OrderNo,
		TotalFee:       totalFee,
		SpbillCreateIP: param.ClientIP,
		NotifyURL:      param.NotifyURL,
		TradeType:      "NATIVE",
		ProductID:      "01001",
	}
	req.Sign = g.sign(req)

	result, err := g.unifiedOrder(req)
	if err != nil {
		return nil, err
	}

	if result.CodeURL == "" {
		return nil, fmt.Errorf("wxpay returned empty code url")
	}

	return &PaymentResult{
		Type:    "url",
		Data:    result.CodeURL,
		OrderNo: param.OrderNo,
	}, nil
}

// wapPay H5支付
func (g *WxPayGateway) wapPay(param *PaymentParam) (*PaymentResult, error) {
	totalFee := int(param.Amount * 100)

	sceneInfo := map[string]interface{}{
		"h5_info": map[string]string{
			"type":    "Wap",
			"wap_url": param.ReturnURL,
			"wap_name": "智简魔方",
		},
	}
	sceneInfoJSON, _ := json.Marshal(sceneInfo)

	req := &WxPayRequest{
		AppID:          g.config.AppID,
		MchID:          g.config.MchID,
		NonceStr:       generateNonceStr(),
		Body:           param.Subject,
		OutTradeNo:     param.OrderNo,
		TotalFee:       totalFee,
		SpbillCreateIP: param.ClientIP,
		NotifyURL:      param.NotifyURL,
		TradeType:      "MWEB",
		SceneInfo:      string(sceneInfoJSON),
	}
	req.Sign = g.sign(req)

	result, err := g.unifiedOrder(req)
	if err != nil {
		return nil, err
	}

	if result.MwebURL == "" {
		return nil, fmt.Errorf("wxpay returned empty mweb url")
	}

	// H5支付需要拼接redirect_url
	mwebURL := result.MwebURL
	if param.ReturnURL != "" {
		if strings.Contains(mwebURL, "?") {
			mwebURL += "&redirect_url=" + url.QueryEscape(param.ReturnURL)
		} else {
			mwebURL += "?redirect_url=" + url.QueryEscape(param.ReturnURL)
		}
	}

	return &PaymentResult{
		Type:    "jump",
		Data:    mwebURL,
		OrderNo: param.OrderNo,
	}, nil
}

// jsapiPay JSAPI支付
func (g *WxPayGateway) jsapiPay(param *PaymentParam) (*PaymentResult, error) {
	totalFee := int(param.Amount * 100)

	req := &WxPayRequest{
		AppID:          g.config.AppID,
		MchID:          g.config.MchID,
		NonceStr:       generateNonceStr(),
		Body:           param.Subject,
		OutTradeNo:     param.OrderNo,
		TotalFee:       totalFee,
		SpbillCreateIP: param.ClientIP,
		NotifyURL:      param.NotifyURL,
		TradeType:      "JSAPI",
	}
	req.Sign = g.sign(req)

	result, err := g.unifiedOrder(req)
	if err != nil {
		return nil, err
	}

	if result.PrepayID == "" {
		return nil, fmt.Errorf("wxpay returned empty prepay id")
	}

	// 返回JSAPI调用参数
	jsapiParams := map[string]string{
		"appId":     g.config.AppID,
		"timeStamp": fmt.Sprintf("%d", time.Now().Unix()),
		"nonceStr":  generateNonceStr(),
		"package":   fmt.Sprintf("prepay_id=%s", result.PrepayID),
		"signType":  "MD5",
	}
	jsapiParams["paySign"] = g.signMap(jsapiParams)

	paramsJSON, _ := json.Marshal(jsapiParams)

	return &PaymentResult{
		Type:    "jsapi",
		Data:    string(paramsJSON),
		OrderNo: param.OrderNo,
	}, nil
}

func (g *WxPayGateway) VerifyNotification(ctx context.Context, data map[string]string) (*NotificationResult, error) {
	// 从map中提取数据
	returnCode := data["return_code"]
	resultCode := data["result_code"]
	if returnCode != "SUCCESS" || resultCode != "SUCCESS" {
		return nil, fmt.Errorf("wxpay notify failed")
	}

	// 验证签名
	sign := data["sign"]
	delete(data, "sign")
	expectedSign := g.signMap(data)
	if sign != expectedSign {
		return nil, fmt.Errorf("invalid wxpay signature")
	}

	// 解析金额
	totalFee := 0
	fmt.Sscanf(data["total_fee"], "%d", &totalFee)
	amount := float64(totalFee) / 100 // 转换为元

	return &NotificationResult{
		OrderNo: data["out_trade_no"],
		TradeNo: data["transaction_id"],
		Amount:  amount,
		Status:  "success",
	}, nil
}

// unifiedOrder 统一下单
func (g *WxPayGateway) unifiedOrder(req *WxPayRequest) (*WxPayResponse, error) {
	xmlData, err := xml.Marshal(req)
	if err != nil {
		return nil, err
	}

	apiURL := "https://api.mch.weixin.qq.com/pay/unifiedorder"
	resp, err := http.Post(apiURL, "application/xml", strings.NewReader(string(xmlData)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result WxPayResponse
	if err := xml.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("invalid wxpay response: %s", string(body))
	}

	if result.ReturnCode != "SUCCESS" {
		return nil, fmt.Errorf("wxpay error: %s", result.ReturnMsg)
	}
	if result.ResultCode != "SUCCESS" {
		return nil, fmt.Errorf("wxpay error: %s - %s", result.ErrCode, result.ErrCodeDes)
	}

	return &result, nil
}

// sign 签名（结构体）
func (g *WxPayGateway) sign(req *WxPayRequest) string {
	params := map[string]string{
		"appid":            req.AppID,
		"mch_id":           req.MchID,
		"nonce_str":        req.NonceStr,
		"body":             req.Body,
		"out_trade_no":     req.OutTradeNo,
		"total_fee":        fmt.Sprintf("%d", req.TotalFee),
		"spbill_create_ip": req.SpbillCreateIP,
		"notify_url":       req.NotifyURL,
		"trade_type":       req.TradeType,
	}
	if req.ProductID != "" {
		params["product_id"] = req.ProductID
	}
	if req.SceneInfo != "" {
		params["scene_info"] = req.SceneInfo
	}
	return g.signMap(params)
}

// signMap 签名（map）
func (g *WxPayGateway) signMap(params map[string]string) string {
	// 按key排序
	keys := make([]string, 0, len(params))
	for k, v := range params {
		if v != "" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	// 拼接
	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, params[k]))
	}
	str := strings.Join(parts, "&") + "&key=" + g.config.APIKey

	// MD5
	h := md5.New()
	h.Write([]byte(str))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

// generateNonceStr 生成随机字符串
func generateNonceStr() string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, 32)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}
