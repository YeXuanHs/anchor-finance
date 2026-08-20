package model

import (
	"time"
)

// UserDownload 用户专属下载（管理员为用户上传的文件）
type UserDownload struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UID       uint      `gorm:"not null;index" json:"uid"`       // 用户ID
	Name      string    `gorm:"size:255;not null" json:"name"`   // 附件显示名
	URL       string    `gorm:"size:255;not null" json:"url"`    // 文件地址
	DowName   string    `gorm:"size:255" json:"downame"`         // 原始文件名
	Remarks   string    `gorm:"size:255" json:"remarks"`         // 备注
	AdminID   uint      `gorm:"default:0" json:"admin_id"`       // 上传管理员ID
	CreatedAt time.Time `json:"created_at"`
}

func (UserDownload) TableName() string {
	return "user_downloads"
}
