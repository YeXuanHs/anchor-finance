package model

import (
	"time"

	"gorm.io/gorm"
)

// ─── 交易市场模型 ───

// MarketplaceListing 市场挂售列表
type MarketplaceListing struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	UserID          uint           `gorm:"index" json:"user_id"`
	User            User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	HostID          uint           `gorm:"index" json:"host_id"`
	Host            Host           `gorm:"foreignKey:HostID" json:"host,omitempty"`
	ProductName     string         `gorm:"type:varchar(256)" json:"product_name"`
	Description     string         `gorm:"type:text" json:"description"`
	OriginalPrice   float64        `gorm:"type:decimal(10,2)" json:"original_price"` // 原价（月付）
	SellPrice       float64        `gorm:"type:decimal(10,2)" json:"sell_price"`     // 售价
	Currency        string         `gorm:"type:varchar(10);default:CNY" json:"currency"`
	Status          int8           `gorm:"default:1;comment:1=在售 2=已售 3=下架" json:"status"`
	ViewCount       int            `gorm:"default:0" json:"view_count"`
	ExpiresAt       *time.Time     `json:"expires_at"` // 到期时间
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// MarketplaceOrder 市场订单
type MarketplaceOrder struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	OrderNo         string         `gorm:"type:varchar(64);uniqueIndex" json:"order_no"`
	ListingID       uint           `gorm:"index" json:"listing_id"`
	Listing         MarketplaceListing `gorm:"foreignKey:ListingID" json:"listing,omitempty"`
	BuyerID         uint           `gorm:"index" json:"buyer_id"`
	Buyer           User           `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`
	SellerID        uint           `gorm:"index" json:"seller_id"`
	Seller          User           `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	HostID          uint           `json:"host_id"`
	Amount          float64        `gorm:"type:decimal(10,2)" json:"amount"`       // 订单金额
	Fee             float64        `gorm:"type:decimal(10,2)" json:"fee"`          // 手续费
	TotalAmount     float64        `gorm:"type:decimal(10,2)" json:"total_amount"` // 总支付金额
	PaymentMethod   string         `gorm:"type:varchar(32)" json:"payment_method"` // full=全额 fee_only=仅手续费
	Status          int8           `gorm:"default:0;comment:0=待支付 1=已支付 2=已转移 3=已完成 4=已取消 5=已退款" json:"status"`
	TransferStatus  int8           `gorm:"default:0;comment:0=未转移 1=转移中 2=转移成功 3=转移失败" json:"transfer_status"`
	TransferID      uint           `json:"transfer_id"` // 关联的产品转移ID
	PaidAt          *time.Time     `json:"paid_at"`
	CompletedAt     *time.Time     `json:"completed_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// MarketplaceChat 市场私聊
type MarketplaceChat struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	ListingID       uint           `gorm:"index" json:"listing_id"`
	Listing         MarketplaceListing `gorm:"foreignKey:ListingID" json:"listing,omitempty"`
	SenderID        uint           `gorm:"index" json:"sender_id"`
	Sender          User           `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
	ReceiverID      uint           `gorm:"index" json:"receiver_id"`
	Receiver        User           `gorm:"foreignKey:ReceiverID" json:"receiver,omitempty"`
	Content         string         `gorm:"type:text" json:"content"`
	IsRead          bool           `gorm:"default:false" json:"is_read"`
	CreatedAt       time.Time      `json:"created_at"`
}

// MarketplaceChatSession 聊天会话
type MarketplaceChatSession struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	ListingID       uint           `gorm:"uniqueIndex:idx_chat_session" json:"listing_id"`
	User1ID         uint           `gorm:"uniqueIndex:idx_chat_session" json:"user1_id"`
	User2ID         uint           `gorm:"uniqueIndex:idx_chat_session" json:"user2_id"`
	LastMessageID   uint           `json:"last_message_id"`
	LastMessage     string         `gorm:"type:varchar(500)" json:"last_message"`
	LastMessageAt   time.Time      `json:"last_message_at"`
	UnreadCount     int            `gorm:"default:0" json:"unread_count"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

// MarketplaceConfig 市场配置
type MarketplaceConfig struct {
	ID                uint    `gorm:"primaryKey" json:"id"`
	Enabled           bool    `gorm:"default:true" json:"enabled"`                       // 是否启用市场
	FeeRate           float64 `gorm:"type:decimal(5,2);default:5" json:"fee_rate"`       // 手续费率（%）
	MinFee            float64 `gorm:"type:decimal(10,2);default:1" json:"min_fee"`       // 最低手续费
	MaxListingDays    int     `gorm:"default:30" json:"max_listing_days"`                 // 挂售最长天数
	MinHoldDays       int     `gorm:"default:7" json:"min_hold_days"`                    // 最少持有天数才能挂售
	RequireRealName   bool    `gorm:"default:false" json:"require_real_name"`             // 买家需要实名
	AllowFeeOnly      bool    `gorm:"default:true" json:"allow_fee_only"`                 // 允许仅付手续费模式
	AutoTransfer      bool    `gorm:"default:true" json:"auto_transfer"`                  // 自动转移
	NotifyEmail       bool    `gorm:"default:true" json:"notify_email"`                   // 邮件通知
	UpdatedAt         time.Time `json:"updated_at"`
}
