package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Order 订单
type Order struct {
	gorm.Model
	OrderNo       string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"order_no"`
	UserID        uint           `gorm:"index;not null" json:"user_id"`
	User          User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ProductID     uint           `gorm:"index" json:"product_id"`
	Product       *Product       `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	UserProductID uint           `gorm:"index" json:"user_product_id"`
	UserProduct   *UserProduct   `gorm:"foreignKey:UserProductID" json:"user_product,omitempty"`
	PromoCodeID   *uint          `gorm:"index" json:"promo_code_id"`
	PromoCode     *PromoCode     `gorm:"foreignKey:PromoCodeID" json:"promo_code,omitempty"`
	Type          string         `gorm:"type:varchar(32);not null;default:'new'" json:"type"` // new/renew/upgrade/transfer
	Amount        datatypes.Decimal `gorm:"type:decimal(20,4);not null" json:"amount"`
	Discount      datatypes.Decimal `gorm:"type:decimal(20,4);default:0" json:"discount"`
	Tax           datatypes.Decimal `gorm:"type:decimal(20,4);default:0" json:"tax"`
	Total         datatypes.Decimal `gorm:"type:decimal(20,4);not null" json:"total"`
	Currency      string         `gorm:"type:varchar(8);default:'CNY';not null" json:"currency"`
	BillingCycle  string         `gorm:"type:varchar(32)" json:"billing_cycle"`
	Quantity      int            `gorm:"default:1" json:"quantity"`
	Description   string         `gorm:"type:varchar(512)" json:"description"`
	PaymentMethod string         `gorm:"type:varchar(64)" json:"payment_method"`
	TransactionID string         `gorm:"type:varchar(256);index" json:"transaction_id"`
	PaidAt        *time.Time     `gorm:"index" json:"paid_at"`
	DueDate       *time.Time     `gorm:"index" json:"due_date"`
	CompletedAt   *time.Time     `json:"completed_at"`
	CancelledAt   *time.Time     `json:"cancelled_at"`
	CancelReason  string         `gorm:"type:varchar(256)" json:"cancel_reason"`
	Status        int16          `gorm:"type:smallint;default:0;not null;index" json:"status"` // 0=待支付 1=已支付 2=处理中 3=已完成 4=已取消 5=已退款 6=部分退款 7=欺诈 8=争议
	PaymentStatus int16          `gorm:"type:smallint;default:0;not null" json:"payment_status"` // 0=未支付 1=已支付 2=部分支付 3=已退款
	AdminNotes    string         `gorm:"type:text" json:"admin_notes"`
	Notes         string         `gorm:"type:text" json:"notes"`
	IPAddress     string         `gorm:"type:varchar(64)" json:"ip_address"`
	Gateway       string         `gorm:"type:varchar(64)" json:"gateway"`
	CommissionPaid bool          `gorm:"default:false" json:"commission_paid"`
	ConfigOptions datatypes.JSON `gorm:"type:json" json:"config_options"`
	CustomFields  datatypes.JSON `gorm:"type:json" json:"custom_fields"`
	Items         []OrderItem    `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

// OrderItem 订单明细
type OrderItem struct {
	gorm.Model
	OrderID   uint           `gorm:"index;not null" json:"order_id"`
	Order     Order          `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	ProductID uint           `gorm:"index" json:"product_id"`
	Product   *Product       `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Name      string         `gorm:"type:varchar(256);not null" json:"name"`
	Quantity  int            `gorm:"default:1" json:"quantity"`
	UnitPrice datatypes.Decimal `gorm:"type:decimal(20,4);not null" json:"unit_price"`
	Amount    datatypes.Decimal `gorm:"type:decimal(20,4);not null" json:"amount"`
	Discount  datatypes.Decimal `gorm:"type:decimal(20,4);default:0" json:"discount"`
	Tax       datatypes.Decimal `gorm:"type:decimal(20,4);default:0" json:"tax"`
	Total     datatypes.Decimal `gorm:"type:decimal(20,4);not null" json:"total"`
	Config    datatypes.JSON   `gorm:"type:json" json:"config"`
}

// OrderNote 订单备注
type OrderNote struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OrderID   uint      `gorm:"index;not null" json:"order_id"`
	Order     Order     `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	AdminID   uint      `gorm:"index" json:"admin_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	IsPrivate bool      `gorm:"default:true" json:"is_private"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName overrides the default table name.
func (OrderNote) TableName() string {
	return "order_notes"
}
