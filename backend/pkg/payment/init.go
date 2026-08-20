package payment

import (
	"context"
	"encoding/json"
	"fmt"
)

// Gateway 支付接口抽象
type Gateway interface {
	// Name 接口名称
	Name() string
	// CreatePayment 创建支付
	CreatePayment(ctx context.Context, param *PaymentParam) (*PaymentResult, error)
	// VerifyNotification 验证回调
	VerifyNotification(ctx context.Context, data map[string]string) (*NotificationResult, error)
}

// PaymentParam 支付参数
type PaymentParam struct {
	OrderNo    string  // 订单号
	Amount     float64 // 金额（元）
	Subject    string  // 商品名称
	ReturnURL  string  // 同步回调URL
	NotifyURL  string  // 异步回调URL
	ClientIP   string  // 客户端IP
}

// PaymentResult 支付结果
type PaymentResult struct {
	Type     string // 跳转类型：jump=跳转, url=二维码链接, balance=余额支付, bank_transfer=银行转账
	Data     string // 跳转数据：URL或HTML
	OrderNo  string // 订单号
}

// NotificationResult 回调结果
type NotificationResult struct {
	OrderNo    string  // 订单号
	TradeNo    string  // 交易平台流水号
	Amount     float64 // 金额
	Status     string  // 状态：success, failed, pending
	RawData    string  // 原始数据
}

// 支付接口类型常量
const (
	GatewayEpay      = "epay"       // 易支付（支持 alipay/wechat/qqpay/usdt/bank）
	GatewayXunhuPay  = "xunhupay"   // 迅虎支付（支持 alipay/wechat）
	GatewayAliPay    = "alipay"     // 支付宝官方
	GatewayWxPay     = "wxpay"      // 微信支付官方
	GatewayBalance   = "balance"    // 余额支付
)

// 支付类型常量（用户可见的支付方式）
const (
	CodeAlipay = "alipay"  // 支付宝
	CodeWechat = "wechat"  // 微信支付
	CodeQQPay  = "qqpay"   // QQ支付
	CodeUsdt   = "usdt"    // USDT
	CodeBank   = "bank"    // 银联（通过易支付接口）
)

// GatewayLabels 接口显示名称
var GatewayLabels = map[string]string{
	GatewayEpay:     "易支付",
	GatewayXunhuPay: "迅虎支付",
	GatewayAliPay:   "支付宝官方",
	GatewayWxPay:    "微信支付官方",
	GatewayBalance:  "余额支付",
}

// CodeLabels 支付类型显示名称
var CodeLabels = map[string]string{
	CodeAlipay: "支付宝",
	CodeWechat: "微信支付",
	CodeQQPay:  "QQ支付",
	CodeUsdt:   "USDT",
	CodeBank:   "银联",
}

// Factory 根据网关类型创建支付接口实例
func Factory(gateway, configJSON string) (Gateway, error) {
	switch gateway {
	case GatewayEpay:
		return NewEpayGateway(configJSON)
	case GatewayXunhuPay:
		return NewXunhuPayGateway(configJSON)
	case GatewayAliPay:
		return NewAliPayGateway(configJSON)
	case GatewayWxPay:
		return NewWxPayGateway(configJSON)
	case GatewayBalance:
		return NewBalanceGateway()
	default:
		return nil, fmt.Errorf("unsupported gateway: %s", gateway)
	}
}

// GetGatewayIcon 获取支付方式图标路径
// 前端根据 code 显示对应图标：alipay.png, wechat.png, qqpay.png 等
func GetGatewayIcon(code string) string {
	icons := map[string]string{
		CodeAlipay: "/assets/payment/alipay.png",
		CodeWechat: "/assets/payment/wechat.png",
		CodeQQPay:  "/assets/payment/qqpay.png",
		CodeUsdt:   "/assets/payment/usdt.png",
		CodeBank:   "/assets/payment/bank.png",
	}
	if icon, ok := icons[code]; ok {
		return icon
	}
	return "/assets/payment/default.png"
}

// ConfigToString 将配置对象序列化为JSON字符串
func ConfigToString(config interface{}) (string, error) {
	data, err := json.Marshal(config)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// GetAll returns all registered payment gateways.
func GetAll() map[string]Gateway {
	return map[string]Gateway{
		GatewayBalance: &BalanceGateway{},
	}
}
