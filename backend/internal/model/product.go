package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// JSON 自定义JSON类型，映射PostgreSQL jsonb
type JSON map[string]interface{}

// ProductGroup 商品分组
type ProductGroup struct {
	gorm.Model
	Name        string    `gorm:"type:varchar(128);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Slug        string    `gorm:"type:varchar(128);uniqueIndex" json:"slug"`
	ParentID    *uint     `gorm:"index" json:"parent_id"`
	Parent      *ProductGroup `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []ProductGroup `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Icon        string    `gorm:"type:varchar(512)" json:"icon"`
	SortOrder   int       `gorm:"default:0;index" json:"sort_order"`
	Status      int16     `gorm:"type:smallint;default:1;not null;index" json:"status"` // 1=显示 0=隐藏
	Products    []Product `gorm:"foreignKey:GroupID" json:"products,omitempty"`
}

// Product 商品
type Product struct {
	gorm.Model
	GroupID       uint           `gorm:"index;not null" json:"group_id"`
	Group         ProductGroup   `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	Name          string         `gorm:"type:varchar(256);not null" json:"name"`
	Slug          string         `gorm:"type:varchar(256);uniqueIndex" json:"slug"`
	Description   string         `gorm:"type:text" json:"description"`
	Content       string         `gorm:"type:text" json:"content"`
	Price         datatypes.Decimal `gorm:"type:decimal(20,4);not null" json:"price"`
	OriginalPrice datatypes.Decimal `gorm:"type:decimal(20,4)" json:"original_price"`
	Currency      string         `gorm:"type:varchar(8);default:'CNY';not null" json:"currency"`
	BillingCycle  string         `gorm:"type:varchar(32)" json:"billing_cycle"` // monthly/quarterly/semi-annually/annually/triennially/onetime
	SetupFee      datatypes.Decimal `gorm:"type:decimal(20,4);default:0" json:"setup_fee"`
	Stock         int            `gorm:"default:-1" json:"stock"` // -1=无限库存
	SalesCount    int            `gorm:"default:0" json:"sales_count"`
	Type          string         `gorm:"type:varchar(32);not null;default:'hosting'" json:"type"` // hosting/domain/ssl/vpn/other
	AutoSetup     bool           `gorm:"default:true" json:"auto_setup"`
	StockControl  bool           `gorm:"default:false" json:"stock_control"`
	QuantityMin   int            `gorm:"default:1" json:"quantity_min"`
	QuantityMax   int            `gorm:"default:1" json:"quantity_max"`
	TrialDays     int            `gorm:"default:0" json:"trial_days"`
	SortOrder     int            `gorm:"default:0;index" json:"sort_order"`
	Status        int16          `gorm:"type:smallint;default:1;not null;index" json:"status"` // 1=上架 0=下架 2=缺货
	Hidden        bool           `gorm:"default:false" json:"hidden"`
	Download      bool           `gorm:"default:false" json:"download"`
	Featured      bool           `gorm:"default:false;index" json:"featured"`
	Image         string         `gorm:"type:varchar(512)" json:"image"`
	Images        datatypes.JSON `gorm:"type:jsonb" json:"images"`       // 图片列表
	ConfigOptions datatypes.JSON `gorm:"type:jsonb" json:"config_options"` // 可配置选项
	Metadata      datatypes.JSON `gorm:"type:jsonb" json:"metadata"`
	Tags          datatypes.JSON `gorm:"type:jsonb" json:"tags"`
	UserProducts  []UserProduct  `gorm:"foreignKey:ProductID" json:"user_products,omitempty"`
}

// UserProduct 用户已购商品/服务实例
type UserProduct struct {
	gorm.Model
	UserID        uint           `gorm:"index;not null" json:"user_id"`
	User          User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	OrderID       uint           `gorm:"index" json:"order_id"`
	Order         *Order         `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	ProductID     uint           `gorm:"index;not null" json:"product_id"`
	Product       Product        `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Name          string         `gorm:"type:varchar(256);not null" json:"name"`
	Domain        string         `gorm:"type:varchar(256);index" json:"domain"`
	Username      string         `gorm:"type:varchar(128)" json:"username"`
	Password      string         `gorm:"type:varchar(256)" json:"-"`
	IP            string         `gorm:"type:varchar(64)" json:"ip"`
	DedicatedIP   string         `gorm:"type:varchar(64)" json:"dedicated_ip"`
	Hostname      string         `gorm:"type:varchar(256)" json:"hostname"`
	NS1           string         `gorm:"type:varchar(256)" json:"ns1"`
	NS2           string         `gorm:"type:varchar(256)" json:"ns2"`
	BillingCycle  string         `gorm:"type:varchar(32)" json:"billing_cycle"`
	Amount        datatypes.Decimal `gorm:"type:decimal(20,4);not null" json:"amount"`
	Currency      string         `gorm:"type:varchar(8);default:'CNY'" json:"currency"`
	RegistrationDate *time.Time  `json:"registration_date"`
	NextDueDate   *time.Time     `gorm:"index" json:"next_due_date"`
	TerminationDate *time.Time   `json:"termination_date"`
	SuspendReason string         `gorm:"type:varchar(256)" json:"suspend_reason"`
	Status        int16          `gorm:"type:smallint;default:1;not null;index" json:"status"` // 1=活跃 2=挂起 3=待开通 4=已终止 5=已过期 6=已取消
	ProvisioningStatus int16    `gorm:"type:smallint;default:0" json:"provisioning_status"` // 0=待处理 1=处理中 2=成功 3=失败
	AdminNotes    string         `gorm:"type:text" json:"admin_notes"`
	Notes         string         `gorm:"type:text" json:"notes"`
	ConfigOptions datatypes.JSON `gorm:"type:jsonb" json:"config_options"` // 用户自定义配置
	CustomFields  datatypes.JSON `gorm:"type:jsonb" json:"custom_fields"`
	Metadata      datatypes.JSON `gorm:"type:jsonb" json:"metadata"`
}
