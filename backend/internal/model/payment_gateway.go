package model

import "gorm.io/gorm"

// PaymentGateway 支付方式
// 每条记录代表一种用户可见的支付方式，如"支付宝-易支付"、"微信支付-虎皮椒"
type PaymentGateway struct {
	gorm.Model
	Name        string `gorm:"type:varchar(64);uniqueIndex;not null" json:"name"`        // 唯一标识，如 EpayAlipay、HpjWechatPay
	Title       string `gorm:"type:varchar(64);not null" json:"title"`                   // 显示名称，如"支付宝-易支付"
	Description string `gorm:"type:varchar(256)" json:"description"`                     // 描述
	Gateway     string `gorm:"type:varchar(32);not null;index" json:"gateway"`           // 支付接口：epay、xunhupay、alipay、wxpay、balance
	Code        string `gorm:"type:varchar(32);not null" json:"code"`                    // 支付类型：alipay、wechat、qqpay、usdt、bank
	Config      string `gorm:"type:json" json:"config"`                                 // 接口配置（商户ID、密钥等）
	FeeRate     float64 `gorm:"type:decimal(5,4);default:0" json:"fee_rate"`             // 手续费率
	MinAmount   float64 `gorm:"type:decimal(12,2);default:0" json:"min_amount"`          // 最低金额
	MaxAmount   float64 `gorm:"type:decimal(12,2);default:0" json:"max_amount"`          // 最高金额
	SortOrder   int    `gorm:"default:0;index" json:"sort_order"`                        // 排序
	IsEnabled   bool   `gorm:"default:true;index" json:"is_enabled"`                     // 是否启用
}
