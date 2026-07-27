package model

import "time"

// UpstreamProvider 上游供应商
type UpstreamProvider struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Type      string    `gorm:"type:varchar(50);not null" json:"type"` // manual, v10, zjmfv3, anchor, whmcs
	APIURL    string    `gorm:"type:varchar(500)" json:"api_url"`
	APIKey    string    `gorm:"type:varchar(255)" json:"-"`
	Config    JSON      `gorm:"type:jsonb" json:"config"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpstreamProduct 上游产品映射
type UpstreamProduct struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	LocalProductID  uint      `gorm:"index;not null" json:"local_product_id"`
	UpstreamID      uint      `gorm:"index;not null" json:"upstream_id"`
	RemoteProductID string    `gorm:"type:varchar(100);not null" json:"remote_product_id"`
	Config          JSON      `gorm:"type:jsonb" json:"config"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// UpstreamSyncLog 上游同步日志
type UpstreamSyncLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UpstreamID uint      `gorm:"index" json:"upstream_id"`
	Action     string    `gorm:"type:varchar(50)" json:"action"`
	Status     string    `gorm:"type:varchar(20)" json:"status"` // success, failed
	Message    string    `gorm:"type:text" json:"message"`
	CreatedAt  time.Time `json:"created_at"`
}
