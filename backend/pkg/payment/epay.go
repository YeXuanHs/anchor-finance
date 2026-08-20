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
	"time"
)

// EpayConfig 易支付配置
type EpayConfig struct {
	PID      string `json:"pid"`       // 商户ID
	Key      string `json:"key"`       // 商户密钥
	APIURL   string `json:"api_url"`   // 接口地址，如 https://pay.example.com
}

// EpayGateway 易支付接口
// 支持 type: alipay, wxpay/wechat, qqpay, usdt, bank
type EpayGateway struct {
	config *EpayConfig
	code   string // alipay, wechat, qqpay, usdt, bank
}

// NewEpayGateway 创建易支付实例
func NewEpayGateway(configJSON string) (*EpayGateway, error) {
	var config EpayConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, fmt.Errorf("invalid epay config: %w", err)
	}
	if config.PID == "" || config.Key == "" || config.APIURL == "" {
		return nil, fmt.Errorf("epay config missing required fields (pid, key, api_url)")
	}
	return &EpayGateway{config: &config}, nil
}

// SetCode 设置支付类型
func (g *EpayGateway) SetCode(code string) {
	g.code = code
}

func (g *EpayGateway) Name() string { return GatewayEpay }

func (g *EpayGateway) CreatePayment(ctx context.Context, param *PaymentParam) (*PaymentResult, error) {
	if g.code == "" {
		return nil, fmt.Errorf("epay code not set")
	}

	// 易支付 type 映射：wechat -> wxpay（易支付接口用 wxpay 不是 wechat）
	epayType := g.code
	if epayType == "wechat" {
		epayType = "wxpay"
	}

	params := map[string]string{
		"pid":          g.config.PID,
		"type":         epayType, // alipay, wxpay, qqpay, usdt, bank
		"out_trade_no": param.OrderNo,
		"name":         param.Subject,
		"money":        fmt.Sprintf("%.2f", param.Amount),
		"notify_url":   param.NotifyURL,
		"return_url":   param.ReturnURL,
	}

	// 生成签名
	params["sign"] = g.sign(params)
	params["sign_type"] = "MD5"

	// 构建跳转URL（参考 zjmf EpayCore::getPayLink）
	query := url.Values{}
	for k, v := range params {
		query.Set(k, v)
	}
	payURL := fmt.Sprintf("%s/submit.php?%s", strings.TrimRight(g.config.APIURL, "/"), query.Encode())

	return &PaymentResult{
		Type:    "jump",
		Data:    payURL,
		OrderNo: param.OrderNo,
	}, nil
}

func (g *EpayGateway) VerifyNotification(ctx context.Context, data map[string]string) (*NotificationResult, error) {
	// 验证签名
	sign := data["sign"]
	delete(data, "sign")
	delete(data, "sign_type")

	expectedSign := g.sign(data)
	if sign != expectedSign {
		return nil, fmt.Errorf("invalid signature")
	}

	tradeStatus := data["trade_status"]
	status := "pending"
	if tradeStatus == "TRADE_SUCCESS" || tradeStatus == "TRADE_FINISHED" {
		status = "success"
	}

	var amount float64
	fmt.Sscanf(data["money"], "%f", &amount)

	return &NotificationResult{
		OrderNo: data["out_trade_no"],
		TradeNo: data["trade_no"],
		Amount:  amount,
		Status:  status,
	}, nil
}

func (g *EpayGateway) sign(params map[string]string) string {
	// 按key排序（参考 zjmf EpayCore::getSign）
	keys := make([]string, 0, len(params))
	for k, v := range params {
		if v != "" && k != "sign" && k != "sign_type" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	// 拼接
	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, params[k]))
	}
	str := strings.Join(parts, "&") + g.config.Key

	// MD5
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// QueryOrder 查询订单状态（参考 zjmf EpayCore::queryOrder）
func (g *EpayGateway) QueryOrder(orderNo string) (string, error) {
	params := map[string]string{
		"act":          "order",
		"pid":          g.config.PID,
		"out_trade_no": orderNo,
	}
	params["sign"] = g.sign(params)
	params["sign_type"] = "MD5"

	query := url.Values{}
	for k, v := range params {
		query.Set(k, v)
	}
	apiURL := fmt.Sprintf("%s/api.php?%s", strings.TrimRight(g.config.APIURL, "/"), query.Encode())

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
}
