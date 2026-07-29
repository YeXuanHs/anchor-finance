package model

import (
	"time"

	"gorm.io/gorm"
)

// WechatUser 微信用户
type WechatUser struct {
	gorm.Model
	UserID      uint       `gorm:"index" json:"user_id"`
	OpenID      string     `gorm:"type:varchar(64);uniqueIndex" json:"open_id"`
	UnionID     string     `gorm:"type:varchar(64);index" json:"union_id"`
	Nickname    string     `gorm:"type:varchar(64)" json:"nickname"`
	Avatar      string     `gorm:"type:varchar(256)" json:"avatar"`
	Sex         int        `gorm:"default:0" json:"sex"`
	Province    string     `gorm:"type:varchar(32)" json:"province"`
	City        string     `gorm:"type:varchar(32)" json:"city"`
	Country     string     `gorm:"type:varchar(32)" json:"country"`
	SubscribeAt *time.Time `json:"subscribe_at"`
	IsSubscribe bool       `gorm:"default:true" json:"is_subscribe"`
	Extra       string     `gorm:"type:jsonb" json:"extra"`
}
