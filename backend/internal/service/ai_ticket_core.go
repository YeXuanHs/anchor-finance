package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// AITicketCoreService AI工单核心服务
// 移植自 mianyu_ai_ticket 的核心逻辑
type AITicketCoreService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewAITicketCoreService(db *gorm.DB, log *logger.Logger) *AITicketCoreService {
	return &AITicketCoreService{db: db, log: log}
}

// ─── 全局配置 ───

func (s *AITicketCoreService) GetConfig(key string) string {
	var cfg model.AITicketConfig
	if err := s.db.Where("`key` = ?", key).First(&cfg).Error; err != nil {
		return ""
	}
	return cfg.Value
}

func (s *AITicketCoreService) SetConfig(key, value string) error {
	var cfg model.AITicketConfig
	if err := s.db.Where("`key` = ?", key).First(&cfg).Error; err != nil {
		cfg = model.AITicketConfig{Key: key, Value: value}
		return s.db.Create(&cfg).Error
	}
	cfg.Value = value
	return s.db.Save(&cfg).Error
}

// GetDashboardConfig 获取控制台配置
func (s *AITicketCoreService) GetDashboardConfig() map[string]interface{} {
	return map[string]interface{}{
		"ai_enabled":  s.GetConfig("ai_enabled"),
		"ai_base_url": s.GetConfig("ai_base_url"),
		"ai_models":   s.GetConfig("ai_models"),
		"max_turns":   s.GetConfig("max_turns"),
		"timeout":     s.GetConfig("timeout"),
		"concurrency": s.GetConfig("concurrency"),
		"mode":        s.GetConfig("mode"),
	}
}

// SaveDashboardConfig 保存控制台配置
func (s *AITicketCoreService) SaveDashboardConfig(data map[string]interface{}) error {
	for k, v := range data {
		if k == "ai_api_key" {
			if str, ok := v.(string); ok && str != "" {
				s.SetConfig(k, str)
			}
		} else {
			if str, ok := v.(string); ok {
				s.SetConfig(k, str)
			}
		}
	}
	return nil
}

// ─── 知识库 CRUD ───

func (s *AITicketCoreService) ListKnowledge(status *int, keyword string) ([]model.AITicketKnowledge, error) {
	var items []model.AITicketKnowledge
	q := s.db.Order("sort ASC, id DESC")
	if status != nil {
		q = q.Where("status = ?", *status)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("title LIKE ? OR keywords LIKE ?", like, like)
	}
	if err := q.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (s *AITicketCoreService) CreateKnowledge(item *model.AITicketKnowledge) error {
	return s.db.Create(item).Error
}

func (s *AITicketCoreService) UpdateKnowledge(id uint, updates map[string]interface{}) error {
	updates["updated_at"] = gorm.Expr("NOW()")
	return s.db.Model(&model.AITicketKnowledge{}).Where("id = ?", id).Updates(updates).Error
}

func (s *AITicketCoreService) DeleteKnowledge(id uint) error {
	return s.db.Delete(&model.AITicketKnowledge{}, id).Error
}

// SearchKnowledge 搜索知识库（用于AI回复）
func (s *AITicketCoreService) SearchKnowledge(query string, limit int) []model.AITicketKnowledge {
	var items []model.AITicketKnowledge
	keywords := extractKeywordsSimple(query)
	if len(keywords) == 0 {
		return items
	}
	q := s.db.Model(&model.AITicketKnowledge{}).Where("status = ?", 1)
	conditions := []string{}
	args := []interface{}{}
	for _, kw := range keywords {
		like := "%" + kw + "%"
		conditions = append(conditions, "(title LIKE ? OR keywords LIKE ? OR question LIKE ?)")
		args = append(args, like, like, like)
	}
	if len(conditions) > 0 {
		q = q.Where(strings.Join(conditions, " OR "), args...)
	}
	q.Order("sort ASC").Limit(limit).Find(&items)
	return items
}

func extractKeywordsSimple(text string) []string {
	words := strings.FieldsFunc(text, func(r rune) bool {
		return r == ' ' || r == ',' || r == '。' || r == '，' || r == '！' || r == '？' || r == '\n'
	})
	var result []string
	seen := map[string]bool{}
	for _, w := range words {
		w = strings.TrimSpace(strings.ToLower(w))
		if len(w) >= 2 && !seen[w] {
			seen[w] = true
			result = append(result, w)
		}
	}
	return result
}

// ImportDefaultKnowledge 导入默认知识库
func (s *AITicketCoreService) ImportDefaultKnowledge() int {
	defaults := []model.AITicketKnowledge{
		{Title: "如何重置密码", Keywords: "密码,重置,忘记,登录", Question: "我忘记了密码怎么办？", Answer: "请在登录页面点击「忘记密码」，输入注册邮箱后系统会发送重置链接。如无法收到邮件，请联系客服。", Sort: 1, Status: 1},
		{Title: "服务器开通时间", Keywords: "开通,时间,多久,等待", Question: "购买后多久能开通？", Answer: "自动开通的产品在付款后即时开通。需要人工审核的产品通常在1-24小时内完成。", Sort: 2, Status: 1},
		{Title: "退款政策", Keywords: "退款,退钱,不想要", Question: "可以退款吗？", Answer: "请参考我们的退款政策页面。一般情况下，未使用的服务可在7天内申请退款。", Sort: 3, Status: 1},
		{Title: "如何提交工单", Keywords: "工单,提交,问题,帮助", Question: "怎么提交工单？", Answer: "登录后进入「工单」页面，点击「新建工单」，选择部门并描述您的问题即可。", Sort: 4, Status: 1},
		{Title: "服务器无法连接", Keywords: "服务器,无法连接,连不上,SSH,RDP", Question: "服务器连不上怎么办？", Answer: "请先检查服务器状态是否正常。如显示运行中但无法连接，请提交工单并附上IP地址，技术会为您排查。", Sort: 5, Status: 1},
	}
	count := 0
	for _, d := range defaults {
		var exists int64
		s.db.Model(&model.AITicketKnowledge{}).Where("title = ?", d.Title).Count(&exists)
		if exists == 0 {
			s.db.Create(&d)
			count++
		}
	}
	return count
}

// ─── 规则管理 ───

func (s *AITicketCoreService) ListRules(status *int) ([]model.AITicketRule, error) {
	var items []model.AITicketRule
	q := s.db.Order("priority ASC, id DESC")
	if status != nil {
		q = q.Where("status = ?", *status)
	}
	if err := q.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (s *AITicketCoreService) CreateRule(item *model.AITicketRule) error {
	return s.db.Create(item).Error
}

func (s *AITicketCoreService) UpdateRule(id uint, updates map[string]interface{}) error {
	updates["updated_at"] = gorm.Expr("NOW()")
	return s.db.Model(&model.AITicketRule{}).Where("id = ?", id).Updates(updates).Error
}

func (s *AITicketCoreService) DeleteRule(id uint) error {
	return s.db.Delete(&model.AITicketRule{}, id).Error
}

// MatchRule 匹配规则
func (s *AITicketCoreService) MatchRule(content string, deptID uint) *model.AITicketRule {
	var rules []model.AITicketRule
	s.db.Where("status = ?", 1).Order("priority ASC").Find(&rules)
	lower := strings.ToLower(content)
	for _, rule := range rules {
		// 检查部门过滤
		if rule.DeptFilter != "" {
			deptIDs := strings.Split(rule.DeptFilter, ",")
			found := false
			for _, id := range deptIDs {
				if strings.TrimSpace(id) == fmt.Sprintf("%d", deptID) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		// 检查关键词
		if rule.Keywords != "" {
			keywords := strings.Split(rule.Keywords, ",")
			matched := false
			for _, kw := range keywords {
				if strings.Contains(lower, strings.ToLower(strings.TrimSpace(kw))) {
					matched = true
					break
				}
			}
			if matched {
				return &rule
			}
		}
	}
	return nil
}

// ─── 队列管理 ───

func (s *AITicketCoreService) ListQueue(status string, page, pageSize int) ([]model.AITicketQueue, int64, error) {
	var items []model.AITicketQueue
	var total int64
	q := s.db.Model(&model.AITicketQueue{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	q.Count(&total)
	offset := (page - 1) * pageSize
	if err := q.Order("id DESC").Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (s *AITicketCoreService) GetQueueStats() map[string]interface{} {
	var pending, processing, completed, failed int64
	s.db.Model(&model.AITicketQueue{}).Where("status = ?", "pending").Count(&pending)
	s.db.Model(&model.AITicketQueue{}).Where("status = ?", "processing").Count(&processing)
	s.db.Model(&model.AITicketQueue{}).Where("status = ?", "completed").Count(&completed)
	s.db.Model(&model.AITicketQueue{}).Where("status = ?", "failed").Count(&failed)
	return map[string]interface{}{
		"pending":    pending,
		"processing": processing,
		"completed":  completed,
		"failed":     failed,
	}
}

// ─── 日志 ───

func (s *AITicketCoreService) ListProcessLogs(ticketID uint, page, pageSize int) ([]model.AITicketProcessLog, int64, error) {
	var items []model.AITicketProcessLog
	var total int64
	q := s.db.Model(&model.AITicketProcessLog{})
	if ticketID > 0 {
		q = q.Where("ticket_id = ?", ticketID)
	}
	q.Count(&total)
	offset := (page - 1) * pageSize
	if err := q.Order("id DESC").Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (s *AITicketCoreService) ListNotifyLogs(ticketID uint, page, pageSize int) ([]model.AITicketNotifyLog, int64, error) {
	var items []model.AITicketNotifyLog
	var total int64
	q := s.db.Model(&model.AITicketNotifyLog{})
	if ticketID > 0 {
		q = q.Where("ticket_id = ?", ticketID)
	}
	q.Count(&total)
	offset := (page - 1) * pageSize
	if err := q.Order("id DESC").Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// ─── 工单模式控制 ───

func (s *AITicketCoreService) GetTicketMode(ticketID uint) string {
	var mode model.AITicketMode
	if err := s.db.First(&mode, ticketID).Error; err != nil {
		return "ai"
	}
	return mode.Mode
}

func (s *AITicketCoreService) SetTicketMode(ticketID uint, mode string) error {
	var existing model.AITicketMode
	if err := s.db.First(&existing, ticketID).Error; err != nil {
		return s.db.Create(&model.AITicketMode{TicketID: ticketID, Mode: mode}).Error
	}
	existing.Mode = mode
	existing.UpdatedAt = time.Now()
	return s.db.Save(&existing).Error
}

// ─── AI 调用（带工具调用支持） ───

type TicketAIRequest struct {
	TicketID uint   `json:"ticket_id"`
	Subject  string `json:"subject"`
	Content  string `json:"content"`
	DeptID   uint   `json:"dept_id"`
}

type TicketAIResult struct {
	ShouldReply bool    `json:"should_reply"`
	Reply       string  `json:"reply"`
	Confidence  float64 `json:"confidence"`
	Decision    string  `json:"decision"`
	Error       string  `json:"error,omitempty"`
}

// ProcessTicket 处理工单（核心逻辑）
func (s *AITicketCoreService) ProcessTicket(req TicketAIRequest) TicketAIResult {
	// 检查是否启用
	if s.GetConfig("ai_enabled") != "1" {
		return TicketAIResult{ShouldReply: false}
	}

	// 检查工单模式
	mode := s.GetTicketMode(req.TicketID)
	if mode == "human" {
		return TicketAIResult{ShouldReply: false}
	}

	// 匹配规则
	rule := s.MatchRule(req.Content, req.DeptID)

	// 搜索知识库
	kbItems := s.SearchKnowledge(req.Subject+" "+req.Content, 5)

	// 构建 prompt
	systemPrompt := "你是一个专业的客服助手。根据用户的问题和知识库内容给出准确回答。"
	prompt := fmt.Sprintf("用户工单主题：%s\n用户工单内容：%s", req.Subject, req.Content)
	if len(kbItems) > 0 {
		kbContext := "\n\n参考知识库："
		for _, item := range kbItems {
			kbContext += fmt.Sprintf("\n【%s】问：%s 答：%s", item.Title, item.Question, item.Answer)
		}
		prompt += kbContext
	}
	if rule != nil && rule.PromptExtra != "" {
		systemPrompt += "\n" + rule.PromptExtra
	}
	if rule != nil && rule.SampleReply != "" {
		prompt += "\n\n参考回复格式：" + rule.SampleReply
	}

	// 调用 AI
	apiKey := s.GetConfig("ai_api_key")
	baseURL := s.GetConfig("ai_base_url")
	models := s.GetConfig("ai_models")
	if apiKey == "" || baseURL == "" || models == "" {
		return TicketAIResult{ShouldReply: false, Error: "AI 配置不完整"}
	}

	reply, confidence, err := s.callAI(baseURL, apiKey, models, systemPrompt, prompt)
	if err != nil {
		s.log.Warnf("AI 调用失败: %v", err)
		// 记录日志
		s.db.Create(&model.AITicketProcessLog{
			TicketID:  req.TicketID,
			Event:     "auto_reply",
			Decision:  "error",
			Status:    "failed",
			Message:   err.Error(),
			ModelUsed: models,
		})
		return TicketAIResult{ShouldReply: false, Error: err.Error()}
	}

	decision := "reply"
	if rule != nil {
		decision = rule.Action
	}

	// 记录日志
	s.db.Create(&model.AITicketProcessLog{
		TicketID:   req.TicketID,
		Event:      "auto_reply",
		Decision:   decision,
		Confidence: confidence,
		Status:     "success",
		Message:    reply,
		ModelUsed:  models,
	})

	// 加入队列
	s.db.Create(&model.AITicketQueue{
		TicketID:  req.TicketID,
		EventType: "auto_reply",
		Status:    "completed",
		ModelUsed: models,
	})

	return TicketAIResult{
		ShouldReply: true,
		Reply:       reply,
		Confidence:  confidence,
		Decision:    decision,
	}
}

func (s *AITicketCoreService) callAI(baseURL, apiKey, modelName, systemPrompt, userPrompt string) (string, float64, error) {
	endpoint := strings.TrimRight(baseURL, "/") + "/chat/completions"
	reqBody := map[string]interface{}{
		"model": modelName,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userPrompt},
		},
		"max_tokens":  2000,
		"temperature": 0.7,
	}
	jsonData, _ := json.Marshal(reqBody)
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", 0, fmt.Errorf("API 错误 %d: %s", resp.StatusCode, string(body))
	}
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", 0, err
	}
	if len(result.Choices) == 0 {
		return "", 0, fmt.Errorf("无回复")
	}
	return result.Choices[0].Message.Content, 0.85, nil
}
