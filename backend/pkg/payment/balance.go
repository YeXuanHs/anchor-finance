package payment

import "context"

// BalanceGateway 余额支付接口
type BalanceGateway struct{}

func NewBalanceGateway() (*BalanceGateway, error) {
	return &BalanceGateway{}, nil
}

func (g *BalanceGateway) Name() string { return GatewayBalance }

func (g *BalanceGateway) CreatePayment(ctx context.Context, param *PaymentParam) (*PaymentResult, error) {
	// 余额支付直接返回成功，由业务层扣减余额
	return &PaymentResult{
		Type:    "balance",
		Data:    "balance_payment",
		OrderNo: param.OrderNo,
	}, nil
}

func (g *BalanceGateway) VerifyNotification(ctx context.Context, data map[string]string) (*NotificationResult, error) {
	return &NotificationResult{
		OrderNo: data["order_no"],
		Status:  "success",
	}, nil
}
