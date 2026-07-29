package payment

import "fmt"

// Factory creates a Gateway from a code and config JSON.
// This is used by the balance/order handler to instantiate gateways
// from the payment_gateways table.
func Factory(code, configJSON string) (Gateway, error) {
	switch code {
	case "alipay":
		return NewAlipayGateway(configJSON)
	case "wechat":
		return NewWechatGateway(configJSON)
	case "wechat_h5":
		gw, err := NewWechatGateway(configJSON)
		if err != nil {
			return nil, err
		}
		return &WechatGatewayH5{gw}, nil
	case "qqpay":
		return NewQQPayGateway(configJSON)
	case "usdt":
		return NewUSDTGateway(configJSON)
	case "xunhupay":
		return NewXunhuPayGateway(configJSON)
	case "epay":
		return NewEpayGateway(configJSON)
	case "stripe":
		return NewStripeGateway(configJSON)
	case "paypal":
		return NewPayPalGateway(configJSON)
	default:
		return nil, fmt.Errorf("unsupported payment gateway: %s", code)
	}
}

// SupportedGateways returns the list of supported gateway codes.
func SupportedGateways() []string {
	return []string{
		"alipay",
		"wechat",
		"wechat_h5",
		"qqpay",
		"usdt",
		"xunhupay",
		"epay",
		"stripe",
		"paypal",
	}
}
