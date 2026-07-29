package upstream

import "strings"

// normalizeBillingCycle maps common billing cycle names to a canonical form.
func normalizeBillingCycle(s string) string {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "monthly", "month", "m", "月付":
		return "monthly"
	case "quarterly", "quarter", "q", "季付":
		return "quarterly"
	case "semi-annually", "semiannual", "semi_annually", "半年付":
		return "semi-annually"
	case "annually", "annual", "yearly", "y", "年付":
		return "annually"
	case "biennially", "biennial", "两年付":
		return "biennially"
	case "triennially", "triennial", "三年付":
		return "triennially"
	case "onetime", "one_time", "one-time", "一次性":
		return "onetime"
	case "free", "免费":
		return "free"
	default:
		return s
	}
}
