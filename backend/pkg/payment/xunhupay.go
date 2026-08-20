package payment

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

// XunhuPayConfig 迅虎支付配置
type XunhuPayConfig struct {
	APIKey    string `json:"api_key"`    // 应用ID (zjmf: api_key)
	SecretKey string `json:"secret_key"` // 应用密钥 (zjmf: secret_key)
	Gateway   string `json:"gateway"`    // 接口地址，默认 https://api.xunhupay.com
}

// XunhuPayGateway 迅虎支付接口
type XunhuPayGateway struct {
	config *XunhuPayConfig
	code   string // alipay, wechat, qqpay
}

// NewXunhuPayGateway 创建迅虎支付实例
func NewXunhuPayGateway(configJSON string) (*XunhuPayGateway, error) {
	var config XunhuPayConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, fmt.Errorf("invalid xunhupay config: %w", err)
	}
	if config.APIKey == "" || config.SecretKey == "" {
		return nil, fmt.Errorf("xunhupay config missing required fields (api_key, secret_key)")
	}
	if config.Gateway == "" {
		config.Gateway = "https://api.xunhupay.com"
	}
	return &XunhuPayGateway{config: &config}, nil
}

// SetCode 设置支付类型
func (g *XunhuPayGateway) SetCode(code string) {
	g.code = code
}

func (g *XunhuPayGateway) Name() string { return GatewayXunhuPay }

func (g *XunhuPayGateway) CreatePayment(ctx context.Context, param *PaymentParam) (*PaymentResult, error) {
	if g.code == "" {
		return nil, fmt.Errorf("xunhupay code not set")
	}

	// 迅虎支付参数（参考 zjmf HpjAlipayPayPlugin/HpjWechatPayPlugin）
	params := map[string]interface{}{
		"version":       "1.1",
		"trade_order_id": param.OrderNo,
		"payment":       g.code, // alipay or wechat
		"total_fee":     fmt.Sprintf("%.2f", param.Amount),
		"title":         param.Subject,
		"notify_url":    param.NotifyURL,
		"return_url":    param.ReturnURL,
	}

	// 调用迅虎支付接口
	result, err := g.execute("/payment/do.html", params)
	if err != nil {
		return nil, fmt.Errorf("xunhupay payment creation failed: %w", err)
	}

	url, _ := result["url"].(string)
	if url == "" {
		return nil, fmt.Errorf("xunhupay returned empty payment url")
	}

	return &PaymentResult{
		Type:    "jump",
		Data:    url,
		OrderNo: param.OrderNo,
	}, nil
}

func (g *XunhuPayGateway) VerifyNotification(ctx context.Context, data map[string]string) (*NotificationResult, error) {
	// 验证签名
	hash := data["hash"]
	delete(data, "hash")

	expectedHash := g.generateHash(data)
	if hash != expectedHash {
		return nil, fmt.Errorf("invalid signature")
	}

	status := "pending"
	if data["trade_status"] == "TRADE_SUCCESS" {
		status = "success"
	}

	var amount float64
	fmt.Sscanf(data["total_fee"], "%f", &amount)

	return &NotificationResult{
		OrderNo: data["trade_order_id"],
		TradeNo: data["transaction_id"],
		Amount:  amount,
		Status:  status,
	}, nil
}

// execute 调用迅虎支付API（参考 zjmf XunhupayClient::execute）
func (g *XunhuPayGateway) execute(path string, params map[string]interface{}) (map[string]interface{}, error) {
	url := strings.TrimRight(g.config.Gateway, "/") + path

	// 添加公共参数
	params["appid"] = g.config.APIKey
	params["time"] = time.Now().Unix()
	params["nonce_str"] = fmt.Sprintf("%d", time.Now().UnixNano())

	// 生成签名
	params["hash"] = g.generateHash(params)

	// 发送请求
	jsonData, _ := json.Marshal(params)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("invalid response: %s", string(body))
	}

	// 检查错误码
	if errcode, ok := result["errcode"].(float64); ok && errcode != 0 {
		errmsg, _ := result["errmsg"].(string)
		return nil, fmt.Errorf("xunhupay error: %s", errmsg)
	}

	return result, nil
}

// generateHash 生成签名（参考 zjmf XunhupayClient::generate_hash）
func (g *XunhuPayGateway) generateHash(params interface{}) string {
	// 将参数转为 map
	var paramMap map[string]interface{}
	switch v := params.(type) {
	case map[string]interface{}:
		paramMap = v
	case map[string]string:
		paramMap = make(map[string]interface{})
		for k, val := range v {
			paramMap[k] = val
		}
	default:
		return ""
	}

	// 按key排序
	keys := make([]string, 0, len(paramMap))
	for k := range paramMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 拼接
	var parts []string
	for _, k := range keys {
		if k == "hash" {
			continue
		}
		v := paramMap[k]
		if v == nil || v == "" {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}
	str := strings.Join(parts, "&") + g.config.SecretKey

	// MD5
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
