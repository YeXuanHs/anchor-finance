package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ConfigOption 配置选项
type ConfigOption struct {
	gorm.Model
	Group       string         `gorm:"type:varchar(64);not null;index" json:"group"` // 分组：product/server/payment/notification 等
	Name        string         `gorm:"type:varchar(128);not null" json:"name"`
	Code        string         `gorm:"type:varchar(64);not null;uniqueIndex:idx_option_code" json:"code"`
	Type        string         `gorm:"type:varchar(32);not null" json:"type"` // text/textarea/number/select/radio/checkbox/switch/color/date/json
	Value       string         `gorm:"type:text" json:"value"`
	DefaultValue string        `gorm:"type:text" json:"default_value"`
	Options     datatypes.JSON `gorm:"type:jsonb" json:"options"` // 下拉/单选/多选的可选项列表
	Placeholder string         `gorm:"type:varchar(256)" json:"placeholder"`
	Tip         string         `gorm:"type:varchar(512)" json:"tip"`
	Validation  string         `gorm:"type:varchar(256)" json:"validation"` // 验证规则
	IsRequired  bool           `gorm:"default:false" json:"is_required"`
	IsPublic    bool           `gorm:"default:false" json:"is_public"`    // 是否前端可见
	IsReadOnly  bool           `gorm:"default:false" json:"is_read_only"` // 是否只读
	SortOrder   int            `gorm:"default:0;index" json:"sort_order"`
	Status      int16          `gorm:"type:smallint;default:1;not null" json:"status"` // 1=启用 0=禁用
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
