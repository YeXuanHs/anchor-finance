package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// WorkTimeService 工作时段判断服务
// 移植自 anchor_cloud_finance_pro 的 WorkTimeService
// 用于客服聊天系统判断当前是否在工作时间内
type WorkTimeService struct{}

func NewWorkTimeService() *WorkTimeService {
	return &WorkTimeService{}
}

// WorkTimeResult 工作时段检查结果
type WorkTimeResult struct {
	State  string `json:"state"`  // online / offhours / rest
	Notice string `json:"notice"` // 非工作时间提示文案
}

const OffhoursText = "当前为非工作时间，您可以先留言，我们会在工作时间尽快回复。"

// WeekdayConfig 每天的工作时间配置
type WeekdayConfig struct {
	Enabled bool   `json:"enabled"`
	Start   string `json:"start"` // HH:mm
	End     string `json:"end"`   // HH:mm
}

// ExceptionConfig 例外日期配置
type ExceptionConfig struct {
	Date  string `json:"date"`  // YYYY-MM-DD 或 MM-DD
	Mode  string `json:"mode"`  // work / rest
	Start string `json:"start"` // HH:mm（mode=work时有效）
	End   string `json:"end"`   // HH:mm（mode=work时有效）
}

// WorkTimeConfig 完整工作时段配置
type WorkTimeConfig struct {
	Weekdays   map[int]WeekdayConfig `json:"weekdays"`
	Exceptions []ExceptionConfig     `json:"exceptions"`
}

// Check 检查当前是否在工作时间
// cfg 可以是 JSON 字符串或 WorkTimeConfig 结构
func (s *WorkTimeService) Check(cfg interface{}) WorkTimeResult {
	var config WorkTimeConfig

	switch v := cfg.(type) {
	case string:
		if err := json.Unmarshal([]byte(v), &config); err != nil {
			config = s.DefaultConfig()
		}
	case []byte:
		if err := json.Unmarshal(v, &config); err != nil {
			config = s.DefaultConfig()
		}
	case WorkTimeConfig:
		config = v
	default:
		// 尝试 JSON 序列化再反序列化
		if data, err := json.Marshal(cfg); err == nil {
			json.Unmarshal(data, &config)
		} else {
			config = s.DefaultConfig()
		}
	}

	now := time.Now()
	ymd := now.Format("2006-01-02")
	md := now.Format("01-02")
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // 周日=7
	}
	hmStr := now.Format("15:04")

	// 1. 例外日期优先（精确到某年 > 每年固定日）
	if ex := s.matchException(config.Exceptions, ymd, md); ex != nil {
		return s.evalPeriod(ex.Mode, ex.Start, ex.End, hmStr)
	}

	// 2. 周几时段
	day, ok := config.Weekdays[weekday]
	if !ok {
		day = s.defaultWeekday(weekday)
	}
	if !day.Enabled {
		return WorkTimeResult{State: "rest", Notice: OffhoursText}
	}
	return s.evalPeriod("work", day.Start, day.End, hmStr)
}

// IsOnline 快速判断是否在工作时间
func (s *WorkTimeService) IsOnline(cfg interface{}) bool {
	return s.Check(cfg).State == "online"
}

// matchException 匹配例外日期
func (s *WorkTimeService) matchException(exceptions []ExceptionConfig, ymd, md string) *ExceptionConfig {
	var yearlyHit *ExceptionConfig
	for _, ex := range exceptions {
		d := strings.TrimSpace(ex.Date)
		if d == "" {
			continue
		}
		// 带年份：精确匹配
		if len(d) == 10 && d[4] == '-' && d[7] == '-' {
			if s.normalizeYmd(d) == ymd {
				return &ex
			}
		} else if len(d) == 5 && d[2] == '-' {
			// 不带年份：每年该日
			if s.normalizeMd(d) == md {
				yearlyHit = &ex
			}
		}
	}
	return yearlyHit
}

// evalPeriod 评估时间段
func (s *WorkTimeService) evalPeriod(mode, start, end, hmStr string) WorkTimeResult {
	if mode == "rest" {
		return WorkTimeResult{State: "rest", Notice: OffhoursText}
	}
	start = s.normalizeHm(start)
	end = s.normalizeHm(end)
	if hmStr < start || hmStr > end {
		return WorkTimeResult{State: "offhours", Notice: OffhoursText}
	}
	return WorkTimeResult{State: "online", Notice: ""}
}

// DefaultConfig 默认配置（周一到周五 09:00-18:00）
func (s *WorkTimeService) DefaultConfig() WorkTimeConfig {
	weekdays := make(map[int]WeekdayConfig)
	for i := 1; i <= 7; i++ {
		weekdays[i] = WeekdayConfig{
			Enabled: i <= 5,
			Start:   "09:00",
			End:     "18:00",
		}
	}
	return WorkTimeConfig{
		Weekdays:   weekdays,
		Exceptions: []ExceptionConfig{},
	}
}

// defaultWeekday 默认周几配置
func (s *WorkTimeService) defaultWeekday(weekday int) WeekdayConfig {
	return WeekdayConfig{
		Enabled: weekday <= 5,
		Start:   "09:00",
		End:     "18:00",
	}
}

// normalizeHm 规范化时间格式为 HH:mm
func (s *WorkTimeService) normalizeHm(hm string) string {
	hm = strings.TrimSpace(hm)
	parts := strings.Split(hm, ":")
	if len(parts) == 2 {
		h, m := parts[0], parts[1]
		if len(h) == 1 {
			h = "0" + h
		}
		return fmt.Sprintf("%s:%s", h, m)
	}
	return "09:00"
}

// normalizeYmd 规范化日期格式为 YYYY-MM-DD
func (s *WorkTimeService) normalizeYmd(d string) string {
	parts := strings.Split(d, "-")
	if len(parts) != 3 {
		return d
	}
	return fmt.Sprintf("%s-%02s-%02s", parts[0], parts[1], parts[2])
}

// normalizeMd 规范化日期格式为 MM-DD
func (s *WorkTimeService) normalizeMd(d string) string {
	parts := strings.Split(d, "-")
	if len(parts) != 2 {
		return d
	}
	return fmt.Sprintf("%02s-%02s", parts[0], parts[1])
}
