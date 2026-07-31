package model

import (
	"time"

	"gorm.io/gorm"
)

// ClientTrackRecord 客户跟踪记录
type ClientTrackRecord struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UID       uint      `gorm:"not null;index" json:"uid"`     // 客户ID
	AdminID   uint      `gorm:"index" json:"admin_id"`         // 操作管理员ID
	Des       string    `gorm:"type:text" json:"des"`          // 记录内容
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联
	Admin     *Admin    `gorm:"foreignKey:AdminID" json:"admin,omitempty"`
	Remarks   []ClientTrackRemark `gorm:"foreignKey:TrackID" json:"remarks,omitempty"`
}

func (ClientTrackRecord) TableName() string {
	return "client_track_records"
}

// ClientTrackRemark 客户跟踪备注
type ClientTrackRemark struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TrackID   uint      `gorm:"not null;index" json:"track_id"` // 记录ID
	AdminID   uint      `gorm:"index" json:"admin_id"`          // 操作管理员ID
	Remark    string    `gorm:"type:text" json:"remark"`        // 补充说明
	CreatedAt time.Time `json:"created_at"`

	// 关联
	Admin     *Admin    `gorm:"foreignKey:AdminID" json:"admin,omitempty"`
}

func (ClientTrackRemark) TableName() string {
	return "client_track_remarks"
}
