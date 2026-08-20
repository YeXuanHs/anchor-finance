package model

import "gorm.io/gorm"

// EmailSuffixWhitelist 邮箱后缀白名单
// 只允许使用白名单内的邮箱后缀注册
type EmailSuffixWhitelist struct {
	gorm.Model
	Suffix    string `gorm:"type:varchar(64);uniqueIndex;not null;comment:邮箱后缀如 gmail.com" json:"suffix"`
	IsDefault bool   `gorm:"default:false;comment:是否为默认后缀（安装时导入的）" json:"is_default"`
	IsActive  bool   `gorm:"default:true;comment:是否启用" json:"is_active"`
	Remark    string `gorm:"type:varchar(255);comment:备注" json:"remark"`
}
