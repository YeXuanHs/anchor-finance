package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// AdvancedOption 高级配置选项
type AdvancedOption struct {
	gorm.Model
	ProductID   uint           `gorm:"index;not null" json:"product_id"`
	Name        string         `gorm:"type:varchar(128);not null" json:"name"`
	Type        string         `gorm:"type:varchar(32);not null" json:"type"` // select/radio/checkbox/text/textarea
	Value       string         `gorm:"type:text" json:"value"`
	Options     datatypes.JSON `gorm:"type:json" json:"options"` // 可选值列表
	Required    bool           `gorm:"default:false" json:"required"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	Status      int16          `gorm:"type:smallint;default:1;not null;index" json:"status"`
	Description string         `gorm:"type:text" json:"description"`
	Metadata    datatypes.JSON `gorm:"type:json" json:"metadata"`
}

// AdvancedOptionLink 配置选项联动条件
type AdvancedOptionLink struct {
	gorm.Model
	ProductID   uint           `gorm:"index;not null" json:"product_id"`
	ConfigID    uint           `gorm:"index;not null" json:"config_id"` // 关联的配置选项ID
	Relation    string         `gorm:"type:varchar(32);not null" json:"relation"` // show/hide/enable/disable
	SubID       datatypes.JSON `gorm:"type:json;not null" json:"sub_id"` // 触发条件 {option_id: value}
	Result      datatypes.JSON `gorm:"type:json" json:"result"` // 执行结果
	Status      int16          `gorm:"type:smallint;default:1;not null" json:"status"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	LastSyncAt  *time.Time     `json:"last_sync_at"`
}
