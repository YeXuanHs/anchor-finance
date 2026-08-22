package model

import "strings"

// BillingCycle 计费周期常量（参考图拉财务 BillingCycle.php）
const (
	BillingCycleMonthly      = "monthly"
	BillingCycleQuarterly    = "quarterly"
	BillingCycleSemiannually = "semiannually"
	BillingCycleAnnually     = "annually"
	BillingCycleBiennially   = "biennially"
	BillingCycleTriennially  = "triennially"
	BillingCycleOneTime      = "one_time"
	BillingCycleFree         = "free"
)

// BillingCycleLabels 计费周期中文标签
var BillingCycleLabels = map[string]string{
	BillingCycleMonthly:      "月付",
	BillingCycleQuarterly:    "季付",
	BillingCycleSemiannually: "半年付",
	BillingCycleAnnually:     "年付",
	BillingCycleBiennially:   "两年付",
	BillingCycleTriennially:  "三年付",
	BillingCycleOneTime:      "一次性",
	BillingCycleFree:         "免费",
}

// BillingCycleMonths 计费周期对应的月数
var BillingCycleMonths = map[string]int{
	BillingCycleMonthly:      1,
	BillingCycleQuarterly:    3,
	BillingCycleSemiannually: 6,
	BillingCycleAnnually:     12,
	BillingCycleBiennially:   24,
	BillingCycleTriennially:  36,
	BillingCycleOneTime:      0,
	BillingCycleFree:         0,
}

// BillingCycleAliases 历史别名 → 规范值
var BillingCycleAliases = map[string]string{
	"onetime":      BillingCycleOneTime,
	"one-time":     BillingCycleOneTime,
	"yearly":       BillingCycleAnnually,
	"biannually":   BillingCycleSemiannually,
	"semi-annually": BillingCycleSemiannually,
}

// RenewableCycles 支持下单与续费的周期白名单
var RenewableCycles = []string{
	BillingCycleMonthly,
	BillingCycleQuarterly,
	BillingCycleSemiannually,
	BillingCycleAnnually,
}

// NormalizeBillingCycle 别名归一 + 去空白 + 转小写
func NormalizeBillingCycle(cycle string) string {
	if cycle == "" {
		return ""
	}
	normalized := strings.ToLower(strings.TrimSpace(cycle))
	if alias, ok := BillingCycleAliases[normalized]; ok {
		return alias
	}
	return normalized
}

// BillingCycleLabel 周期 → 中文标签
func BillingCycleLabel(cycle string) string {
	normalized := NormalizeBillingCycle(cycle)
	if label, ok := BillingCycleLabels[normalized]; ok {
		return label
	}
	return normalized
}

// BillingCycleMonths 周期 → 月数，未知返回0
func BillingCycleMonthsValue(cycle string) int {
	normalized := NormalizeBillingCycle(cycle)
	if months, ok := BillingCycleMonths[normalized]; ok {
		return months
	}
	return 0
}

// IsRenewableCycle 是否是支持续费的周期
func IsRenewableCycle(cycle string) string {
	normalized := NormalizeBillingCycle(cycle)
	for _, c := range RenewableCycles {
		if c == normalized {
			return normalized
		}
	}
	return ""
}
