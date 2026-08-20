package model

import "time"

// Language 语言
type Language struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Code      string    `gorm:"type:varchar(10);uniqueIndex;not null" json:"code"`      // zh-CN, en-US, zh-TW
	Name      string    `gorm:"type:varchar(50);not null" json:"name"`                  // 中文简体
	Flag      string    `gorm:"type:varchar(10)" json:"flag"`                           // CN, US, TW
	IsDefault bool      `gorm:"default:false" json:"is_default"`                        // 是否默认语言
	Status    int16     `gorm:"default:1" json:"status"`                                // 1=启用 0=禁用
	CreatedAt time.Time `json:"created_at"`
}

// LangKey 语言键值
type LangKey struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Key       string    `gorm:"type:varchar(100);uniqueIndex:idx_lang_key;not null" json:"key"` // 键名
	Module    string    `gorm:"type:varchar(50);index" json:"module"`                           // 模块
	CreatedAt time.Time `json:"created_at"`
}

// LangTranslation 语言翻译
type LangTranslation struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	KeyID      uint   `gorm:"uniqueIndex:idx_lang_trans;not null" json:"key_id"`
	LangCode   string `gorm:"type:varchar(10);uniqueIndex:idx_lang_trans;not null" json:"lang_code"`
	Value      string `gorm:"type:text" json:"value"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TableName 表名
func (Language) TableName() string         { return "languages" }
func (LangKey) TableName() string          { return "lang_keys" }
func (LangTranslation) TableName() string  { return "lang_translations" }
