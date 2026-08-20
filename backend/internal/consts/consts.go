package consts

// 错误码
const (
	ErrCodeSuccess       = 0
	ErrCodeBadRequest    = 400
	ErrCodeUnauthorized  = 401
	ErrCodeForbidden     = 403
	ErrCodeNotFound      = 404
	ErrCodeServer        = 500
	ErrCodeTooMany       = 429
	
	// 业务错误码
	ErrCodeUserExists     = 1001
	ErrCodeUserNotFound   = 1002
	ErrCodePasswordWrong  = 1003
	ErrCodeAccountBanned  = 1004
	ErrCodeTokenExpired   = 1005
	ErrCodeTokenInvalid   = 1006
	ErrCodeCaptchaWrong   = 1007
	
	ErrCodeProductNotFound = 2001
	ErrCodeProductOffSale  = 2002
	ErrCodeStockNotEnough  = 2003
	
	ErrCodeOrderNotFound   = 3001
	ErrCodeOrderPaid       = 3002
	ErrCodeOrderCancelled  = 3003
	
	ErrCodeInvoiceNotFound = 4001
	ErrCodeInvoicePaid     = 4002
	ErrCodeBalanceNotEnough = 4003
	
	ErrCodeTicketNotFound  = 5001
	ErrCodeTicketClosed    = 5002
)

// 用户状态
const (
	UserStatusActive   = 1
	UserStatusInactive = 2
	UserStatusBanned   = 3
	UserStatusClosed   = 4
)

// 订单状态
const (
	OrderStatusPending   = 1
	OrderStatusPaid      = 2
	OrderStatusActive    = 3
	OrderStatusCancelled = 4
	OrderStatusFraud     = 5
	OrderStatusRefunded  = 6
)

// 账单状态
const (
	InvoiceStatusPending   = 1
	InvoiceStatusPaid      = 2
	InvoiceStatusCancelled = 3
	InvoiceStatusRefunded  = 4
	InvoiceStatusOverdue   = 5
)

// 工单状态
const (
	TicketStatusOpen     = 1
	TicketStatusReplied  = 2
	TicketStatusPending  = 3
	TicketStatusClosed   = 4
	TicketStatusResolved = 5
)

// 产品状态
const (
	ProductStatusActive   = 1
	ProductStatusInactive = 2
	ProductStatusHidden   = 3
)

// 支付方式
const (
	PayMethodBalance = "balance"
	PayMethodAlipay  = "alipay"
	PayMethodWechat  = "wechat"
	PayMethodCredit  = "credit"
)

// 分页默认值
const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)
