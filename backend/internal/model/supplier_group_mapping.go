package model

import "time"

// SupplierGroupMapping 供应商分组映射（MD 7.2.5：上游分组→本地分组）
type SupplierGroupMapping struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	SupplierID     uint   `gorm:"index;not null" json:"supplier_id"`
	RemoteGroupID  string `gorm:"size:100;not null" json:"remote_group_id"`  // 上游分组ID
	RemoteGroupName string `gorm:"size:200" json:"remote_group_name"`        // 上游分组名称（显示用）
	LocalGroupID   uint   `gorm:"index;not null" json:"local_group_id"`      // 本地分组ID
	ProfitRate     float64 `gorm:"type:decimal(5,2);default:25" json:"profit_rate"` // 该映射的利润率
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
