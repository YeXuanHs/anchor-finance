package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// ────────────────────────────────────────────────────────────
// 敏感数据加密（移植自 zjmf cmf_encrypt/cmf_decrypt）
// ────────────────────────────────────────────────────────────

// Encryptor 数据加密器
type Encryptor struct {
	key []byte
}

// NewEncryptor 创建加密器
func NewEncryptor(secretKey string) *Encryptor {
	// 使用 SHA256 生成 32 字节密钥
	hash := sha256.Sum256([]byte(secretKey))
	return &Encryptor{key: hash[:]}
}

// Encrypt 加密数据（AES-256-GCM）
func (e *Encryptor) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密数据（AES-256-GCM）
func (e *Encryptor) Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// ────────────────────────────────────────────────────────────
// 数据脱敏（敏感信息隐藏）
// ────────────────────────────────────────────────────────────

// MaskPhone 手机号脱敏（保留前3后4）
// 例: 13812345678 -> 138****5678
func MaskPhone(phone string) string {
	if len(phone) < 7 {
		return phone
	}
	return phone[:3] + "****" + phone[len(phone)-4:]
}

// MaskEmail 邮箱脱敏（保留前3和@后内容）
// 例: test@example.com -> t***@example.com
func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}
	username := parts[0]
	if len(username) <= 3 {
		return username + "@" + parts[1]
	}
	return username[:3] + "***@" + parts[1]
}

// MaskName 姓名脱敏（保留姓）
// 例: 张三 -> 张*, 李四 -> 李*
func MaskName(name string) string {
	if len(name) == 0 {
		return name
	}
	runes := []rune(name)
	if len(runes) == 1 {
		return "*"
	}
	return string(runes[0]) + strings.Repeat("*", len(runes)-1)
}

// MaskIDCard 身份证号脱敏（保留前4后4）
// 例: 110101199001011234 -> 1101**********1234
func MaskIDCard(idCard string) string {
	if len(idCard) < 8 {
		return idCard
	}
	return idCard[:4] + strings.Repeat("*", len(idCard)-8) + idCard[len(idCard)-4:]
}

// MaskBankCard 银行卡号脱敏（保留后4位）
// 例: 6222021234567890123 -> ***********0123
func MaskBankCard(cardNo string) string {
	if len(cardNo) < 4 {
		return cardNo
	}
	return strings.Repeat("*", len(cardNo)-4) + cardNo[len(cardNo)-4:]
}

// MaskPassword 密码脱敏（全部隐藏）
func MaskPassword(password string) string {
	return "********"
}

// MaskIP IP地址脱敏（保留前两段）
// 例: 192.168.1.100 -> 192.168.*.*
func MaskIP(ip string) string {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return ip
	}
	return parts[0] + "." + parts[1] + ".*.*"
}

// MaskAPIKey API密钥脱敏（保留前6后4）
// 例: abcdefghijklmnop -> abcdef****mnop
func MaskAPIKey(apiKey string) string {
	if len(apiKey) <= 10 {
		return strings.Repeat("*", len(apiKey))
	}
	return apiKey[:6] + "****" + apiKey[len(apiKey)-4:]
}

// ────────────────────────────────────────────────────────────
// 用户信息脱敏（前台显示用）
// ────────────────────────────────────────────────────────────

// MaskedUserInfo 脱敏后的用户信息（前台显示）
type MaskedUserInfo struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Phone     string `json:"phone,omitempty"`
	Email     string `json:"email,omitempty"`
	RealName  string `json:"real_name,omitempty"`
	IDCard    string `json:"id_card,omitempty"`
	BankCard  string `json:"bank_card,omitempty"`
	IP        string `json:"ip,omitempty"`
	CreatedAt string `json:"created_at"`
}

// MaskUserInfoForFrontend 前台用户信息脱敏
// 只在前台（用户端）显示时使用，后台管理员看到完整信息
func MaskUserInfoForFrontend(user interface{}) interface{} {
	// 根据具体结构体类型进行脱敏
	// 这里只是示例，实际使用时需要根据 User 模型调整
	return user
}

// ShouldMask 判断是否需要脱敏
// isAdmin: 是否是管理员
// 返回 true 表示需要脱敏（前台），false 表示不脱敏（后台）
func ShouldMask(isAdmin bool) bool {
	return !isAdmin
}

// MaskString 根据是否管理员决定是否脱敏
func MaskString(value string, isAdmin bool, maskFunc func(string) string) string {
	if isAdmin {
		return value
	}
	return maskFunc(value)
}

// MaskPhoneByRole 根据角色脱敏手机号
func MaskPhoneByRole(phone string, isAdmin bool) string {
	return MaskString(phone, isAdmin, MaskPhone)
}

// MaskEmailByRole 根据角色脱敏邮箱
func MaskEmailByRole(email string, isAdmin bool) string {
	return MaskString(email, isAdmin, MaskEmail)
}

// MaskIDCardByRole 根据角色脱敏身份证
func MaskIDCardByRole(idCard string, isAdmin bool) string {
	return MaskString(idCard, isAdmin, MaskIDCard)
}

// MaskBankCardByRole 根据角色脱敏银行卡
func MaskBankCardByRole(cardNo string, isAdmin bool) string {
	return MaskString(cardNo, isAdmin, MaskBankCard)
}

// MaskIPByRole 根据角色脱敏IP
func MaskIPByRole(ip string, isAdmin bool) string {
	return MaskString(ip, isAdmin, MaskIP)
}

// ────────────────────────────────────────────────────────────
// 敏感关键词过滤（移植自 zjmf mask_keywords）
// ────────────────────────────────────────────────────────────

// KeywordFilter 关键词过滤器
type KeywordFilter struct {
	keywords []string
}

// NewKeywordFilter 创建关键词过滤器
func NewKeywordFilter(keywords []string) *KeywordFilter {
	return &KeywordFilter{keywords: keywords}
}

// ContainsSensitive 检查文本是否包含敏感关键词
func (f *KeywordFilter) ContainsSensitive(text string) (bool, string) {
	for _, keyword := range f.keywords {
		if keyword != "" && strings.Contains(text, keyword) {
			return true, keyword
		}
	}
	return false, ""
}

// Filter 过滤敏感关键词（替换为*）
func (f *KeywordFilter) Filter(text string) string {
	result := text
	for _, keyword := range f.keywords {
		if keyword != "" {
			replacement := strings.Repeat("*", len([]rune(keyword)))
			result = strings.ReplaceAll(result, keyword, replacement)
		}
	}
	return result
}

// ────────────────────────────────────────────────────────────
// XSS 防护
// ────────────────────────────────────────────────────────────

var xssPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`),
	regexp.MustCompile(`(?i)<iframe[^>]*>.*?</iframe>`),
	regexp.MustCompile(`(?i)<object[^>]*>.*?</object>`),
	regexp.MustCompile(`(?i)<embed[^>]*>`),
	regexp.MustCompile(`(?i)<applet[^>]*>.*?</applet>`),
	regexp.MustCompile(`(?i)javascript:`),
	regexp.MustCompile(`(?i)vbscript:`),
	regexp.MustCompile(`(?i)on\w+\s*=`),
}

// SanitizeHTML 清理 HTML（防止 XSS）
func SanitizeHTML(html string) string {
	result := html
	for _, pattern := range xssPatterns {
		result = pattern.ReplaceAllString(result, "")
	}
	return result
}

// EscapeHTML 转义 HTML 特殊字符
func EscapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}

// ────────────────────────────────────────────────────────────
// SQL 注入防护
// ────────────────────────────────────────────────────────────

var sqlInjectionPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(union\s+select)`),
	regexp.MustCompile(`(?i)(select\s+.*\s+from)`),
	regexp.MustCompile(`(?i)(insert\s+into)`),
	regexp.MustCompile(`(?i)(delete\s+from)`),
	regexp.MustCompile(`(?i)(drop\s+(table|database))`),
	regexp.MustCompile(`(?i)(update\s+.*\s+set)`),
	regexp.MustCompile(`(?i)(--|;|/\*|\*/)`),
	regexp.MustCompile(`(?i)(or\s+1\s*=\s*1)`),
	regexp.MustCompile(`(?i)(or\s+'1'\s*=\s*'1')`),
	regexp.MustCompile(`(?i)(exec\s*\()`),
}

// DetectSQLInjection 检测 SQL 注入
func DetectSQLInjection(input string) bool {
	for _, pattern := range sqlInjectionPatterns {
		if pattern.MatchString(input) {
			return true
		}
	}
	return false
}

// ────────────────────────────────────────────────────────────
// 密码强度检查
// ────────────────────────────────────────────────────────────

// PasswordStrength 密码强度
type PasswordStrength int

const (
	PasswordWeak   PasswordStrength = iota // 弱
	PasswordMedium                         // 中
	PasswordStrong                         // 强
	PasswordVeryStrong                     // 非常强
)

// CheckPasswordStrength 检查密码强度
func CheckPasswordStrength(password string) PasswordStrength {
	score := 0

	// 长度检查
	if len(password) >= 8 {
		score++
	}
	if len(password) >= 12 {
		score++
	}
	if len(password) >= 16 {
		score++
	}

	// 包含小写字母
	if matched, _ := regexp.MatchString(`[a-z]`, password); matched {
		score++
	}

	// 包含大写字母
	if matched, _ := regexp.MatchString(`[A-Z]`, password); matched {
		score++
	}

	// 包含数字
	if matched, _ := regexp.MatchString(`[0-9]`, password); matched {
		score++
	}

	// 包含特殊字符
	if matched, _ := regexp.MatchString(`[^a-zA-Z0-9]`, password); matched {
		score++
	}

	// 计算强度
	switch {
	case score >= 6:
		return PasswordVeryStrong
	case score >= 4:
		return PasswordStrong
	case score >= 2:
		return PasswordMedium
	default:
		return PasswordWeak
	}
}

// GetPasswordStrengthDesc 获取密码强度描述
func GetPasswordStrengthDesc(strength PasswordStrength) string {
	switch strength {
	case PasswordVeryStrong:
		return "非常强"
	case PasswordStrong:
		return "强"
	case PasswordMedium:
		return "中"
	default:
		return "弱"
	}
}

// ────────────────────────────────────────────────────────────
// 安全配置检查
// ────────────────────────────────────────────────────────────

// SecurityConfig 安全配置
type SecurityConfig struct {
	JWTSecret          string `json:"jwt_secret"`
	PasswordMinLength  int    `json:"password_min_length"`
	PasswordMaxLength  int    `json:"password_max_length"`
	LoginMaxAttempts   int    `json:"login_max_attempts"`
	LoginLockMinutes   int    `json:"login_lock_minutes"`
	SessionTimeout     int    `json:"session_timeout"`
	EnableHTTPS        bool   `json:"enable_https"`
	EnableCSRF         bool   `json:"enable_csrf"`
	EnableXSSFilter    bool   `json:"enable_xss_filter"`
	EnableSQLFilter    bool   `json:"enable_sql_filter"`
	AllowedOrigins     string `json:"allowed_origins"`
}

// DefaultSecurityConfig 默认安全配置
func DefaultSecurityConfig() *SecurityConfig {
	return &SecurityConfig{
		PasswordMinLength: 6,
		PasswordMaxLength: 32,
		LoginMaxAttempts:  5,
		LoginLockMinutes:  15,
		SessionTimeout:    72,
		EnableHTTPS:       false,
		EnableCSRF:        true,
		EnableXSSFilter:   true,
		EnableSQLFilter:   true,
		AllowedOrigins:    "*",
	}
}

// Validate 验证安全配置
func (c *SecurityConfig) Validate() []string {
	var errors []string

	if len(c.JWTSecret) < 32 {
		errors = append(errors, "JWT密钥长度至少32位")
	}

	if c.PasswordMinLength < 6 {
		errors = append(errors, "密码最小长度不能小于6")
	}

	if c.PasswordMaxLength > 32 {
		errors = append(errors, "密码最大长度不能超过32")
	}

	if c.LoginMaxAttempts < 3 || c.LoginMaxAttempts > 20 {
		errors = append(errors, "登录尝试次数应在3-20之间")
	}

	if c.LoginLockMinutes < 5 || c.LoginLockMinutes > 60 {
		errors = append(errors, "登录锁定时间应在5-60分钟之间")
	}

	return errors
}
