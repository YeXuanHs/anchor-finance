package validator

import (
	"regexp"
	"unicode"
)

// ────────────────────────────────────────────────────────────
// 验证规则（移植自 zjmf）
// ────────────────────────────────────────────────────────────

var (
	// 用户名：4-20位字母数字下划线
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{4,20}$`)
	
	// 密码：6-32位任意字符
	passwordMinLen = 6
	passwordMaxLen = 32
	
	// 手机号：4-11位数字
	phoneRegex = regexp.MustCompile(`^\d{4,11}$`)
	
	// 邮箱格式
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	
	// API密钥：只能包含字母、数字、下划线、破折号，6-32位
	apiKeyRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-]{6,32}$`)
	
	// QQ号：最多20位
	qqMaxLen = 20
	
	// 公司名：最多50字符
	companyMaxLen = 50
	
	// 地址：最多100字符
	addressMaxLen = 100
	
	// 个性签名：最多20字符
	signatureMaxLen = 20
)

// ValidateUsername 验证用户名（移植自 zjmf）
// 规则：必填，4-20位字母数字下划线
func ValidateUsername(username string) (bool, string) {
	if username == "" {
		return false, "用户名不能为空"
	}
	if len(username) < 4 || len(username) > 20 {
		return false, "用户名长度必须在4-20位之间"
	}
	if !usernameRegex.MatchString(username) {
		return false, "用户名只能包含字母、数字和下划线"
	}
	return true, ""
}

// ValidatePassword 验证密码（移植自 zjmf）
// 规则：必填，6-32位
func ValidatePassword(password string) (bool, string) {
	if password == "" {
		return false, "密码不能为空"
	}
	if len(password) < passwordMinLen {
		return false, "密码长度至少6位"
	}
	if len(password) > passwordMaxLen {
		return false, "密码长度最多32位"
	}
	return true, ""
}

// ValidatePasswordMatch 验证两次密码是否一致
func ValidatePasswordMatch(password, confirmPassword string) (bool, string) {
	if password != confirmPassword {
		return false, "两次输入密码不一致"
	}
	return true, ""
}

// ValidatePasswordNotSame 验证新密码不能和旧密码相同
func ValidatePasswordNotSame(oldPassword, newPassword string) (bool, string) {
	if oldPassword == newPassword {
		return false, "新密码不能与旧密码相同"
	}
	return true, ""
}

// ValidatePhone 验证手机号（移植自 zjmf）
// 规则：必填，4-11位数字
func ValidatePhone(phone string) (bool, string) {
	if phone == "" {
		return false, "手机号不能为空"
	}
	if !phoneRegex.MatchString(phone) {
		return false, "手机号格式错误，应为4-11位数字"
	}
	return true, ""
}

// ValidateEmail 验证邮箱（移植自 zjmf）
// 规则：必填，邮箱格式
func ValidateEmail(email string) (bool, string) {
	if email == "" {
		return false, "邮箱不能为空"
	}
	if !emailRegex.MatchString(email) {
		return false, "邮箱格式错误"
	}
	return true, ""
}

// ValidateEmailOptional 验证邮箱（可选）
// 规则：如果填写了必须是邮箱格式
func ValidateEmailOptional(email string) (bool, string) {
	if email == "" {
		return true, ""
	}
	if !emailRegex.MatchString(email) {
		return false, "邮箱格式错误"
	}
	return true, ""
}

// ValidateAPIKey 验证 API 密钥（移植自 zjmf）
// 规则：不能包含中文，必须包含大小写字母、数字、特殊字符，8-32位
func ValidateAPIKey(apiKey string) (bool, string) {
	if apiKey == "" {
		return false, "API密钥不能为空"
	}
	
	// 检查是否包含中文
	for _, r := range apiKey {
		if unicode.Is(unicode.Han, r) {
			return false, "API密钥不能包含中文"
		}
	}
	
	// 检查长度
	if len(apiKey) < 8 || len(apiKey) > 32 {
		return false, "API密钥长度必须在8-32位之间"
	}
	
	// 检查是否包含大小写字母、数字、特殊字符
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false
	
	for _, r := range apiKey {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r) || r == '_':
			hasSpecial = true
		}
	}
	
	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return false, "API密钥必须包含大小写字母、数字和特殊字符"
	}
	
	return true, ""
}

// ValidateQQ 验证 QQ 号（移植自 zjmf）
// 规则：最多20字符
func ValidateQQ(qq string) (bool, string) {
	if len(qq) > qqMaxLen {
		return false, "QQ号不能超过20个字符"
	}
	return true, ""
}

// ValidateCompany 验证公司名（移植自 zjmf）
// 规则：最多50字符
func ValidateCompany(company string) (bool, string) {
	if len(company) > companyMaxLen {
		return false, "公司名不能超过50个字符"
	}
	return true, ""
}

// ValidateAddress 验证地址（移植自 zjmf）
// 规则：最多100字符
func ValidateAddress(address string) (bool, string) {
	if len(address) > addressMaxLen {
		return false, "地址不能超过100个字符"
	}
	return true, ""
}

// ValidateSignature 验证个性签名（移植自 zjmf）
// 规则：最多20字符
func ValidateSignature(signature string) (bool, string) {
	if len(signature) > signatureMaxLen {
		return false, "个性签名不能超过20个字符"
	}
	return true, ""
}

// ValidateNickname 验证昵称（移植自 zjmf）
// 规则：最多64字符
func ValidateNickname(nickname string) (bool, string) {
	if len(nickname) > 64 {
		return false, "昵称不能超过64个字符"
	}
	return true, ""
}

// ValidateSiteName 验证站点名称（移植自 zjmf）
// 规则：不能为空
func ValidateSiteName(name string) (bool, string) {
	if name == "" {
		return false, "站点名称不能为空"
	}
	return true, ""
}

// ValidateSiteURL 验证站点URL（移植自 zjmf）
// 规则：必须以 http:// 或 https:// 开头，不能以 / 结尾
func ValidateSiteURL(url string) (bool, string) {
	if url == "" {
		return false, "站点URL不能为空"
	}
	if len(url) < 8 {
		return false, "站点URL格式错误"
	}
	if url[:7] != "http://" && url[:8] != "https://" {
		return false, "站点URL必须以 http:// 或 https:// 开头"
	}
	if url[len(url)-1] == '/' {
		return false, "站点URL不能以 / 结尾"
	}
	return true, ""
}

// ValidateAdminPath 验证后台路径（移植自 zjmf）
// 规则：必须同时包含数字和字母，2-32位
func ValidateAdminPath(path string) (bool, string) {
	if path == "" {
		return false, "后台路径不能为空"
	}
	if len(path) < 2 || len(path) > 32 {
		return false, "后台路径长度必须在2-32位之间"
	}
	
	hasLetter := false
	hasDigit := false
	
	for _, r := range path {
		if unicode.IsLetter(r) {
			hasLetter = true
		} else if unicode.IsDigit(r) {
			hasDigit = true
		} else {
			return false, "后台路径只能包含字母和数字"
		}
	}
	
	if !hasLetter {
		return false, "后台路径不能全是数字"
	}
	if !hasDigit {
		return false, "后台路径不能全是字母，必须包含数字"
	}
	
	return true, ""
}

// ValidateTicketSubject 验证工单主题（移植自 zjmf）
// 规则：不能为空，最多200字符
func ValidateTicketSubject(subject string) (bool, string) {
	if subject == "" {
		return false, "工单主题不能为空"
	}
	if len(subject) > 200 {
		return false, "工单主题不能超过200个字符"
	}
	return true, ""
}

// ValidateTicketContent 验证工单内容（移植自 zjmf）
// 规则：不能为空
func ValidateTicketContent(content string) (bool, string) {
	if content == "" {
		return false, "工单内容不能为空"
	}
	return true, ""
}

// ValidateOrderAmount 验证订单金额
// 规则：必须大于0
func ValidateOrderAmount(amount float64) (bool, string) {
	if amount <= 0 {
		return false, "订单金额必须大于0"
	}
	return true, ""
}

// ValidateInvoiceAmount 验证发票金额
// 规则：必须大于0
func ValidateInvoiceAmount(amount float64) (bool, string) {
	if amount <= 0 {
		return false, "发票金额必须大于0"
	}
	return true, ""
}

// ValidatePort 验证端口号
// 规则：1-65535
func ValidatePort(port int) (bool, string) {
	if port < 1 || port > 65535 {
		return false, "端口号必须在1-65535之间"
	}
	return true, ""
}

// ValidateIP 验证IP地址
// 规则：IPv4格式
func ValidateIP(ip string) (bool, string) {
	if ip == "" {
		return false, "IP地址不能为空"
	}
	// 简单的IPv4验证
	parts := regexp.MustCompile(`^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})$`).FindStringSubmatch(ip)
	if parts == nil {
		return false, "IP地址格式错误"
	}
	for i := 1; i <= 4; i++ {
		num := 0
		for _, c := range parts[i] {
			num = num*10 + int(c-'0')
		}
		if num > 255 {
			return false, "IP地址格式错误"
		}
	}
	return true, ""
}
