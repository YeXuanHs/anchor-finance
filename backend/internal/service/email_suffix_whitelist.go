package service

import (
	"strings"

	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type EmailSuffixWhitelistService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewEmailSuffixWhitelistService(db *gorm.DB, log *logger.Logger) *EmailSuffixWhitelistService {
	return &EmailSuffixWhitelistService{db: db, log: log}
}

type EmailSuffixWhitelistItem struct {
	ID        uint   `json:"id"`
	Suffix    string `json:"suffix"`
	IsDefault bool   `json:"is_default"`
	IsActive  bool   `json:"is_active"`
	Remark    string `json:"remark"`
}

// List 获取所有邮箱后缀白名单
func (s *EmailSuffixWhitelistService) List(showInactive bool) ([]EmailSuffixWhitelistItem, error) {
	var items []EmailSuffixWhitelistItem
	q := s.db.Table("email_suffix_whitelists").Select("id, suffix, is_default, is_active, remark").Order("suffix ASC")
	if !showInactive {
		q = q.Where("is_active = ?", true)
	}
	if err := q.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// Add 添加邮箱后缀
func (s *EmailSuffixWhitelistService) Add(suffix, remark string) error {
	suffix = normalizeSuffix(suffix)
	return s.db.Exec("INSERT INTO email_suffix_whitelists (suffix, is_default, is_active, remark, created_at, updated_at) VALUES (?, false, true, ?, NOW(), NOW())", suffix, remark).Error
}

// Update 更新邮箱后缀
func (s *EmailSuffixWhitelistService) Update(id uint, isActive *bool, remark *string) error {
	updates := map[string]interface{}{}
	if isActive != nil {
		updates["is_active"] = *isActive
	}
	if remark != nil {
		updates["remark"] = *remark
	}
	if len(updates) == 0 {
		return nil
	}
	updates["updated_at"] = gorm.Expr("NOW()")
	return s.db.Table("email_suffix_whitelists").Where("id = ?", id).Updates(updates).Error
}

// Delete 删除邮箱后缀
func (s *EmailSuffixWhitelistService) Delete(id uint) error {
	return s.db.Exec("DELETE FROM email_suffix_whitelists WHERE id = ?", id).Error
}

// BatchDelete 批量删除
func (s *EmailSuffixWhitelistService) BatchDelete(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return s.db.Exec("DELETE FROM email_suffix_whitelists WHERE id IN ?", ids).Error
}

// IsAllowed 检查邮箱后缀是否允许注册
func (s *EmailSuffixWhitelistService) IsAllowed(email string) bool {
	email = strings.TrimSpace(strings.ToLower(email))
	atIdx := strings.LastIndex(email, "@")
	if atIdx < 0 || atIdx == len(email)-1 {
		return false
	}
	suffix := email[atIdx+1:]

	var count int64
	s.db.Table("email_suffix_whitelists").Where("suffix = ? AND is_active = ?", suffix, true).Count(&count)
	return count > 0
}

// ImportDefaults 导入默认邮箱后缀
func (s *EmailSuffixWhitelistService) ImportDefaults() {
	defaults := getDefaultSuffixes()
	for _, d := range defaults {
		s.db.Exec("INSERT IGNORE INTO email_suffix_whitelists (suffix, is_default, is_active, remark, created_at, updated_at) VALUES (?, true, true, ?, NOW(), NOW())", d.suffix, d.remark)
	}
}

type suffixEntry struct {
	suffix string
	remark string
}

func normalizeSuffix(suffix string) string {
	suffix = strings.TrimSpace(strings.ToLower(suffix))
	suffix = strings.TrimPrefix(suffix, "@")
	return suffix
}

func getDefaultSuffixes() []suffixEntry {
	return []suffixEntry{
		// 国际主流
		{"gmail.com", "Google 邮箱"},
		{"outlook.com", "Microsoft Outlook"},
		{"hotmail.com", "Microsoft Hotmail"},
		{"live.com", "Microsoft Live"},
		{"yahoo.com", "Yahoo 邮箱"},
		{"yahoo.co.jp", "Yahoo Japan"},
		{"yahoo.co.uk", "Yahoo UK"},
		{"icloud.com", "Apple iCloud"},
		{"me.com", "Apple Me"},
		{"mac.com", "Apple Mac"},
		{"aol.com", "AOL 邮箱"},
		{"protonmail.com", "ProtonMail 加密邮箱"},
		{"proton.me", "Proton 新域名"},
		{"zoho.com", "Zoho 邮箱"},
		{"yandex.com", "Yandex 邮箱"},
		{"mail.ru", "Mail.ru"},
		{"gmx.com", "GMX 邮箱"},
		{"gmx.net", "GMX 德国"},
		{"fastmail.com", "Fastmail"},
		{"tutanota.com", "Tutanota 加密邮箱"},
		{"tuta.io", "Tuta 新域名"},

		// 国内主流
		{"qq.com", "QQ 邮箱"},
		{"vip.qq.com", "QQ 会员邮箱"},
		{"foxmail.com", "Foxmail 邮箱"},
		{"163.com", "网易 163 邮箱"},
		{"126.com", "网易 126 邮箱"},
		{"yeah.net", "网易 Yeah 邮箱"},
		{"sina.com", "新浪邮箱"},
		{"sina.cn", "新浪 CN 邮箱"},
		{"sohu.com", "搜狐邮箱"},
		{"aliyun.com", "阿里云邮箱"},
		{"139.com", "中国移动 139 邮箱"},
		{"189.cn", "中国电信 189 邮箱"},
		{"wo.cn", "中国联通沃邮箱"},
		{"tom.com", "TOM 邮箱"},
		{"21cn.com", "21CN 邮箱"},

		// 教育/组织
		{"edu.cn", "中国教育邮箱"},
		{"ac.cn", "中国科学院"},
		{"edu", "国际教育域名"},
		{"org", "组织域名"},
		{"gov.cn", "中国政府域名"},

		// 企业常用
		{"microsoft.com", "Microsoft"},
		{"apple.com", "Apple"},
		{"google.com", "Google"},
		{"amazon.com", "Amazon"},
		{"meta.com", "Meta"},
		{"twitter.com", "Twitter/X"},
		{"github.com", "GitHub"},
	}
}
