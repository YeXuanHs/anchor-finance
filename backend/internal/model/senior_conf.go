package model

import "gorm.io/gorm"

// SeniorConf 高级配置
type SeniorConf struct {
	gorm.Model
	Name        string `gorm:"type:varchar(64);not null" json:"name"`
	Code        string `gorm:"type:varchar(32);uniqueIndex;not null" json:"code"`
	Type        string `gorm:"type:varchar(16);not null" json:"type"` // text/number/boolean/json
	Value       string `gorm:"type:text" json:"value"`
	Default     string `gorm:"type:text" json:"default"`
	Group       string `gorm:"type:varchar(32);index" json:"group"`
	Description string `gorm:"type:varchar(256)" json:"description"`
	IsPublic    bool   `gorm:"default:false" json:"is_public"`
}
