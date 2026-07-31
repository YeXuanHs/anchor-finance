package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Invoice 发票/账单
type Invoice struct {
	gorm.Model
	InvoiceNo    string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"invoice_no"`
	UserID       uint           `gorm:"index;not null" json:"user_id"`
	User         User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	OrderID      *uint          `gorm:"index" json:"order_id"`
	Order        *Order         `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Type         string         `gorm:"type:varchar(32);not null;default:'invoice'" json:"type"` // invoice/proforma/credit/debit
	Currency     string         `gorm:"type:varchar(8);default:'CNY';not null" json:"currency"`
	SubTotal     datatypes.Decimal `gorm:"type:decimal(20,4);not null" json:"sub_total"`
	Tax          datatypes.Decimal `gorm:"type:decimal(20,4);default:0" json:"tax"`
	TaxRate      datatypes.Decimal `gorm:"type:decimal(5,4);default:0" json:"tax_rate"`
	Discount     datatypes.Decimal `gorm:"type:decimal(20,4);default:0" json:"discount"`
	Credit       datatypes.Decimal `gorm:"type:decimal(20,4);default:0" json:"credit"`
	Total        datatypes.Decimal `gorm:"type:decimal(20,4);not null" json:"total"`
	PaidAmount   datatypes.Decimal `gorm:"type:decimal(20,4);default:0" json:"paid_amount"`
	Balance      datatypes.Decimal `gorm:"type:decimal(20,4);default:0" json:"balance"`
	Status       int16          `gorm:"type:smallint;default:0;not null;index" json:"status"` // 0=待支付 1=已支付 2=部分支付 3=已取消 4=已退款 5=逾期 6=已催付
	PaymentMethod string        `gorm:"type:varchar(64)" json:"payment_method"`
	TransactionID string        `gorm:"type:varchar(256);index" json:"transaction_id"`
	DueDate      *time.Time     `gorm:"index" json:"due_date"`
	PaidAt       *time.Time     `json:"paid_at"`
	CancelledAt  *time.Time     `json:"cancelled_at"`
	// 发票抬头信息
	BillingName    string         `gorm:"type:varchar(128)" json:"billing_name"`
	BillingAddress string         `gorm:"type:varchar(512)" json:"billing_address"`
	BillingCity    string         `gorm:"type:varchar(128)" json:"billing_city"`
	BillingState   string         `gorm:"type:varchar(128)" json:"billing_state"`
	BillingZip     string         `gorm:"type:varchar(16)" json:"billing_zip"`
	BillingCountry string         `gorm:"type:varchar(64)" json:"billing_country"`
	TaxID          string         `gorm:"type:varchar(64)" json:"tax_id"`
	Notes          string         `gorm:"type:text" json:"notes"`
	AdminNotes     string         `gorm:"type:text" json:"admin_notes"`
	// 自动开票
	AutoBilling    bool           `gorm:"default:false" json:"auto_billing"`
	LastAutoAttempt *time.Time    `json:"last_auto_attempt"`
	Attempts       int            `gorm:"default:0" json:"attempts"`
	Items          []InvoiceItem  `gorm:"foreignKey:InvoiceID" json:"items,omitempty"`
	Transactions   []Transaction  `gorm:"foreignKey:InvoiceID" json:"transactions,omitempty"`
	Metadata       datatypes.JSON `gorm:"type:json" json:"metadata"`
}

// InvoiceItem 发票明细
type InvoiceItem struct {
	gorm.Model
	InvoiceID uint           `gorm:"index;not null" json:"invoice_id"`
	Invoice   Invoice        `gorm:"foreignKey:InvoiceID" json:"invoice,omitempty"`
	Type      string         `gorm:"type:varchar(32);not null;default:'product'" json:"type"` // product/fee/discount/credit/custom
	RelID     uint           `gorm:"index" json:"rel_id"` // 关联ID（商品ID等）
	RelType   string         `gorm:"type:varchar(32)" json:"rel_type"` // 关联类型
	Description string       `gorm:"type:varchar(512);not null" json:"description"`
	Quantity  int            `gorm:"default:1" json:"quantity"`
	UnitPrice datatypes.Decimal `gorm:"type:decimal(20,4);not null" json:"unit_price"`
	Discount  datatypes.Decimal `gorm:"type:decimal(20,4);default:0" json:"discount"`
	Tax       datatypes.Decimal `gorm:"type:decimal(20,4);default:0" json:"tax"`
	Total     datatypes.Decimal `gorm:"type:decimal(20,4);not null" json:"total"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
}

// Transaction 支付交易记录
type Transaction struct {
	gorm.Model
	TransactionNo string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"transaction_no"`
	UserID        uint           `gorm:"index;not null" json:"user_id"`
	User          User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	InvoiceID     *uint          `gorm:"index" json:"invoice_id"`
	Invoice       *Invoice       `gorm:"foreignKey:InvoiceID" json:"invoice,omitempty"`
	OrderID       *uint          `gorm:"index" json:"order_id"`
	Order         *Order         `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Gateway       string         `gorm:"type:varchar(64);not null" json:"gateway"`
	GatewayTransID string        `gorm:"type:varchar(256);index" json:"gateway_trans_id"`
	Amount        datatypes.Decimal `gorm:"type:decimal(20,4);not null" json:"amount"`
	Fee           datatypes.Decimal `gorm:"type:decimal(20,4);default:0" json:"fee"`
	Currency      string         `gorm:"type:varchar(8);default:'CNY';not null" json:"currency"`
	ExchangeRate  datatypes.Decimal `gorm:"type:decimal(16,8);default:1" json:"exchange_rate"`
	Type          string         `gorm:"type:varchar(32);not null;default:'payment'" json:"type"` // payment/refund/credit/debit/withdrawal
	Status        int16          `gorm:"type:smallint;default:0;not null;index" json:"status"` // 0=待处理 1=成功 2=失败 3=已取消 4=已退款 5=争议中
	CompletedAt   *time.Time     `json:"completed_at"`
	Notes         string         `gorm:"type:text" json:"notes"`
	AdminNotes    string         `gorm:"type:text" json:"admin_notes"`
	IPAddress     string         `gorm:"type:varchar(64)" json:"ip_address"`
	CallbackData  datatypes.JSON `gorm:"type:json" json:"callback_data"` // 网关回调原始数据
	Metadata      datatypes.JSON `gorm:"type:json" json:"metadata"`
}
