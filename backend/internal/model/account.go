package model

// Account 交易流水记录
type Account struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	UID         uint    `gorm:"index;not null" json:"uid"`
	Currency    string  `gorm:"type:varchar(16);default:'CNY'" json:"currency"`
	Gateway     string  `gorm:"type:varchar(64)" json:"gateway"`
	PayTime     int64   `gorm:"index" json:"pay_time"`
	UpdateTime  int64   `json:"update_time"`
	CreateTime  int64   `json:"create_time"`
	DeleteTime  int64   `gorm:"index;default:0" json:"delete_time"`
	Description string  `gorm:"type:text" json:"description"`
	TransID     string  `gorm:"type:varchar(128);index" json:"trans_id"`
	InvoiceID   uint    `gorm:"index" json:"invoice_id"`
	AmountIn    float64 `gorm:"type:decimal(20,4);default:0" json:"amount_in"`
	Fees        float64 `gorm:"type:decimal(20,4);default:0" json:"fees"`
	AmountOut   float64 `gorm:"type:decimal(20,4);default:0" json:"amount_out"`
	Rate        float64 `gorm:"type:decimal(20,8);default:1" json:"rate"`
	Refund      int     `gorm:"default:0" json:"refund"`
}

func (Account) TableName() string {
	return "accounts"
}
