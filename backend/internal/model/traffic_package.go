package model

import (
	"time"

	"gorm.io/gorm"
)

// TrafficPackage 流量包模型
type TrafficPackage struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Volume    int64          `gorm:"not null" json:"volume"` // 流量大小，单位MB
	Price     float64        `gorm:"not null" json:"price"`
	Unit      string         `gorm:"size:10;default:GB" json:"unit"`
	Status    string         `gorm:"size:20;default:active" json:"status"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (TrafficPackage) TableName() string {
	return "traffic_packages"
}

// TrafficLog 流量使用记录模型
type TrafficLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	ServiceID uint      `gorm:"index;not null" json:"service_id"`
	PackageID uint      `gorm:"index" json:"package_id"`
	Volume    int64     `gorm:"not null" json:"volume"` // 使用流量，单位MB
	Direction string    `gorm:"size:10;not null" json:"direction"` // in, out
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (TrafficLog) TableName() string {
	return "traffic_logs"
}
