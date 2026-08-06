package model

import (
	"time"

	"gorm.io/gorm"
)

// ProductTransfer 产品转移记录
type ProductTransfer struct {
	gorm.Model
	TransferNo    string     `gorm:"type:varchar(64);uniqueIndex;not null" json:"transfer_no"`
	FromUserID    uint       `gorm:"index;not null" json:"from_user_id"`
	ToUserID      uint       `gorm:"index;not null" json:"to_user_id"`
	UserProductID uint       `gorm:"index;not null" json:"user_product_id"`
	Reason        string     `gorm:"type:text" json:"reason"`
	Status        int8       `gorm:"type:smallint;default:1" json:"status"` // 1=处理中 2=成功 3=失败
	ProcessedAt   *time.Time `json:"processed_at"`
	Remark        string     `gorm:"type:text" json:"remark"`
}
