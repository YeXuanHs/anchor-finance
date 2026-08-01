package service

import (
	"strings"
	"unicode/utf8"
)

// KbService 客服内置知识库
// 移植自 anchor_cloud_finance_pro 的 KbService
// 基于常见问题提供快捷回答
type KbService struct{}

func NewKbService() *KbService {
	return &KbService{}
}

// Match 匹配用户问题，返回回答（无匹配返回空字符串）
func (s *KbService) Match(question string) string {
	q := s.normalize(question)
	if q == "" {
		return ""
	}

	// 精确匹配
	exactMap := map[string]string{
		"如何充值账户余额":   "recharge",
		"怎么购买和续费产品": "buy_renew",
		"忘记密码怎么找回":   "pwreset",
	}
	for k, code := range exactMap {
		if q == s.normalize(k) || q == s.normalize(k+"?") || q == s.normalize(k+"？") {
			return s.answer(code)
		}
	}

	// 模糊匹配
	if s.hasAll(q, []string{"充值"}) && s.hasAny(q, []string{"余额", "账户", "帐户", "怎么", "如何", "怎样"}) {
		return s.answer("recharge")
	}
	if s.hasAny(q, []string{"续费"}) || (s.hasAny(q, []string{"购买", "订购"}) && s.hasAny(q, []string{"产品", "怎么", "如何"})) {
		return s.answer("buy_renew")
	}
	if (s.hasAny(q, []string{"忘记"}) && s.hasAny(q, []string{"密码"})) ||
		(s.hasAny(q, []string{"找回"}) && s.hasAny(q, []string{"密码"})) {
		return s.answer("pwreset")
	}

	return ""
}

// answer 根据问题代码返回回答
func (s *KbService) answer(code string) string {
	switch code {
	case "recharge":
		return strings.Join([]string{
			"账户充值可以直接点下面链接进入，不必一层层找菜单：",
			"",
			"- **一键直达**：[账户充值](/addfunds)",
			"- 相关入口：[账单列表](/billing) · [交易记录](/transaction)",
			"",
			"**如需按侧栏操作：**",
			"1. 登录后，打开左侧 **财务管理**",
			"2. 点击 **账户充值**（或直接用上面的链接）",
			"3. 填写 **充值金额**，选择 **付款方式**，点 **立刻充值**",
			"",
			"支付完成后余额会到账，可用于购买/续费。若有疑问可切到「人工客服」，或 [提交工单](/submitticket)。",
		}, "\n")

	case "buy_renew":
		return strings.Join([]string{
			"购买与续费可以点链接直达，省得在侧栏里翻：",
			"",
			"**购买新产品**",
			"- **一键订购**：[订购产品](/cart)",
			"- 登录后也可打开侧栏 **产品与服务 → 订购产品**",
			"",
			"**续费已开通产品**",
			"- **产品列表**：[我的产品](/service)",
			"- 进入对应产品详情页后，点击 **立即续费**，选择周期并支付",
			"",
			"**账单支付**",
			"- 若已有待支付账单：[账单列表](/billing)",
			"",
			"卡住了可切「人工客服」，或 [提交工单](/submitticket)。",
		}, "\n")

	case "pwreset":
		return strings.Join([]string{
			"忘记密码可以直接进入找回页：",
			"",
			"- **一键找回**：[忘记密码 / 重置密码](/pwreset)",
			"- 也可在 [登录页](/login) 点击「忘记密码?」",
			"",
			"**找回步骤：**",
			"1. 打开上面的找回链接",
			"2. 选择 **手机找回** 或 **邮箱找回**",
			"3. 填写绑定手机号/注册邮箱，获取验证码并设置新密码",
			"4. 返回 [登录](/login) 使用新密码登录",
			"",
			"收不到验证码时，可切「人工客服」协助，或 [提交工单](/submitticket)。",
		}, "\n")
	}

	return ""
}

// normalize 规范化文本用于匹配
func (s *KbService) normalize(text string) string {
	text = strings.TrimSpace(strings.ToLower(text))
	// 移除标点符号
	replacer := strings.NewReplacer(
		"?", "", "？", "", "。", "", ".", "",
		"！", "", "!", "", "，", "", ",", "",
		"、", "", "：", "", ":", "",
	)
	text = replacer.Replace(text)
	// 移除空白
	text = strings.ReplaceAll(text, " ", "")
	return text
}

// hasAll 检查文本是否包含所有关键词
func (s *KbService) hasAll(text string, words []string) bool {
	for _, w := range words {
		if !strings.Contains(text, s.normalize(w)) {
			return false
		}
	}
	return true
}

// hasAny 检查文本是否包含任一关键词
func (s *KbService) hasAny(text string, words []string) bool {
	for _, w := range words {
		if strings.Contains(text, s.normalize(w)) {
			return true
		}
	}
	return false
}

// HasCJK 检查文本是否包含中日韩字符（辅助判断是否中文问题）
func (s *KbService) HasCJK(text string) bool {
	for _, r := range text {
		if r >= 0x4E00 && r <= 0x9FFF {
			return true
		}
		if r >= 0x3400 && r <= 0x4DBF {
			return true
		}
	}
	return false
}

// Truncate 截断文本到指定 rune 长度
func (s *KbService) Truncate(text string, maxLen int) string {
	if utf8.RuneCountInString(text) <= maxLen {
		return text
	}
	runes := []rune(text)
	return string(runes[:maxLen]) + "..."
}
