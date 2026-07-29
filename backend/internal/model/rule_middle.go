package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// RuleMiddle 规则中间件
type RuleMiddle struct {
	gorm.Model
	Name         string         `gorm:"type:varchar(128);not null;uniqueIndex" json:"name"`
	AddRole      int            `gorm:"default:0" json:"add_role"`
	AddRoleMenu  string         `gorm:"type:varchar(256)" json:"add_role_menu"`
	CatRole      int            `gorm:"default:0" json:"cat_role"`
	CatRoleMenu  string         `gorm:"type:varchar(256)" json:"cat_role_menu"`
	DelRole      int            `gorm:"default:0" json:"del_role"`
	DelRoleMenu  string         `gorm:"type:varchar(256)" json:"del_role_menu"`
	EditRole     int            `gorm:"default:0" json:"edit_role"`
	EditRoleMenu string         `gorm:"type:varchar(256)" json:"edit_role_menu"`
	Status       int16          `gorm:"type:smallint;default:1;not null;index" json:"status"`
	SortOrder    int            `gorm:"default:0" json:"sort_order"`
	Extra        datatypes.JSON `gorm:"type:jsonb" json:"extra"`
}
