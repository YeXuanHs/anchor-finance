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
	CouponID      *uint          `gorm:"index" json:"coupon_id"`
	Coupon        *Coupon        `gorm:"foreignKey:CouponID" json:"coupon,omitempty"`
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

// Coupon 优惠券
type Coupon struct {
	gorm.Model
	Code              string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"code"`
	Name              string         `gorm:"type:varchar(128);not null" json:"name"`
	Description       string         `gorm:"type:text" json:"description"`
	Type              string         `gorm:"type:varchar(16);not null;default:'percentage'" json:"type"` // percentage/fixed/override/free
	Value             datatypes.Decimal `gorm:"type:decimal(20,4);not null" json:"value"`
	MaxDiscount       datatypes.Decimal `gorm:"type:decimal(20,4)" json:"max_discount"`
	MinOrderAmount    datatypes.Decimal `gorm:"type:decimal(20,4);default:0" json:"min_order_amount"`
	MaxUses           int            `gorm:"default:0" json:"max_uses"` // 0=不限
	UsedCount         int            `gorm:"default:0" json:"used_count"`
	MaxUsesPerUser    int            `gorm:"default:1" json:"max_uses_per_user"`
	StartDate         *time.Time     `gorm:"index" json:"start_date"`
	EndDate           *time.Time     `gorm:"index" json:"end_date"`
	ProductIDs        datatypes.JSON `gorm:"type:json" json:"product_ids"`        // 适用商品ID列表
	Cycles            datatypes.JSON `gorm:"type:json" json:"cycles"`             // 适用计费周期列表
	GroupIDs          datatypes.JSON `gorm:"type:json" json:"group_ids"`          // 适用用户组ID列表
	ExcludeProductIDs datatypes.JSON `gorm:"type:json" json:"exclude_product_ids"`
	OnlyNewClient     bool           `gorm:"default:false" json:"only_new_client"`  // 仅新客户可用
	OnlyOldClient     bool           `gorm:"default:false" json:"only_old_client"`  // 仅老客户可用
	OncePerClient     bool           `gorm:"default:false" json:"once_per_client"`  // 每客户仅限一次
	RequiresExist     *uint          `json:"requires_exist"`                        // 需要已购买的商品ID
	ApplyOnce         bool           `gorm:"default:false" json:"apply_once"`       // 仅首单可用
	Recurring         bool           `gorm:"default:false" json:"recurring"`        // 续费时也享受
	Status            int16          `gorm:"type:smallint;default:1;not null;index" json:"status"` // 1=启用 0=禁用
	Orders            []Order        `gorm:"foreignKey:CouponID" json:"orders,omitempty"`
}

// CouponUsageLog 优惠券使用记录
type CouponUsageLog struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	CouponID uint    `gorm:"index;not null" json:"coupon_id"`
	UserID   uint    `gorm:"index;not null" json:"user_id"`
	OrderID  uint    `gorm:"index" json:"order_id"`
	Discount float64 `gorm:"type:decimal(20,4);not null" json:"discount"`
}
