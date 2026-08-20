package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// AITicketCoreService AI工单核心服务
// 移植自 mianyu_ai_ticket 的核心逻辑，支持 Agent（Function Calling）模式
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

// ─── 工具管理 ───

// ListTools 获取所有工具列表（按分类）
func (s *AITicketCoreService) ListTools() []ToolCategory {
	categories := AllToolCategories()

	// 读取启用/禁用配置
	var configs []model.AIToolConfig
	s.db.Find(&configs)
	configMap := make(map[string]int8, len(configs))
	for _, c := range configs {
		configMap[c.Name] = c.Enabled
	}

	// 如果没有配置记录，所有工具默认启用
	for i, cat := range categories {
		for j, tool := range cat.Tools {
			if enabled, ok := configMap[tool.Name]; ok {
				if enabled == 0 {
					categories[i].Tools[j].RiskLevel = tool.RiskLevel + " (已禁用)"
				}
			}
		}
	}

	return categories
}

// GetEnabledTools 获取已启用的工具名列表
func (s *AITicketCoreService) GetEnabledTools() []string {
	var configs []model.AIToolConfig
	s.db.Find(&configs)

	if len(configs) == 0 {
		return AllToolNames()
	}

	var enabled []string
	for _, c := range configs {
		if c.Enabled == 1 {
			enabled = append(enabled, c.Name)
		}
	}
	return enabled
}

// SetToolEnabled 启用/禁用工具
func (s *AITicketCoreService) SetToolEnabled(name string, enabled bool) error {
	var cfg model.AIToolConfig
	enabledVal := int8(1)
	if !enabled {
		enabledVal = 0
	}

	if err := s.db.Where("name = ?", name).First(&cfg).Error; err != nil {
		// 找到对应的 risk level
		riskLevel := "low"
		for _, t := range allToolsMap() {
			if t.Name == name {
				riskLevel = t.RiskLevel
				break
			}
		}
		cfg = model.AIToolConfig{Name: name, Enabled: enabledVal, RiskLevel: riskLevel}
		return s.db.Create(&cfg).Error
	}

	cfg.Enabled = enabledVal
	return s.db.Save(&cfg).Error
}

// ListToolExecutionLogs 获取工具执行日志
func (s *AITicketCoreService) ListToolExecutionLogs(ticketID uint, page, pageSize int) ([]model.AIToolExecutionLog, int64, error) {
	var items []model.AIToolExecutionLog
	var total int64
	q := s.db.Model(&model.AIToolExecutionLog{})
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

// ─── AI 调用（带 Function Calling 支持） ───

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

// AIChatMessage AI 聊天消息
type AIChatMessage struct {
	Role       string           `json:"role"`
	Content    string           `json:"content,omitempty"`
	ToolCalls  []AIToolCall     `json:"tool_calls,omitempty"`
	ToolCallID string           `json:"tool_call_id,omitempty"`
}

// AIToolCall AI 工具调用
type AIToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

// AIChatResponse AI 聊天响应
type AIChatResponse struct {
	Choices []struct {
		Message struct {
			Content   string       `json:"content"`
			ToolCalls []AIToolCall `json:"tool_calls"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

// ProcessTicket 处理工单（核心逻辑，支持 Agent 模式）
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

	// 构建上下文（使用 TicketContextService）
	ctxSvc := NewTicketContextService(s.db, s.log)
	contextJson := ctxSvc.BuildContextJson(req.TicketID, 0)

	// 系统提示词
	systemPrompt := s.GetConfig("system_prompt")
	if systemPrompt == "" {
		systemPrompt = DefaultSystemPrompt()
	}
	// 注入上下文
	if strings.Contains(systemPrompt, "{context}") {
		systemPrompt = strings.Replace(systemPrompt, "{context}", contextJson, 1)
	} else {
		systemPrompt += "\n\n## 上下文数据（JSON）\n" + contextJson
	}

	// 规则额外提示词
	if rule != nil && rule.PromptExtra != "" {
		systemPrompt += "\n" + rule.PromptExtra
	}

	// 构建消息列表（从工单对话历史）
	messages := ctxSvc.BuildMessages(req.TicketID)
	if len(messages) == 0 {
		return TicketAIResult{ShouldReply: false, Error: "工单无有效对话内容"}
	}

	// 追加知识库参考到用户最后一条消息
	if len(kbItems) > 0 {
		kbContext := "\n\n参考知识库："
		for _, item := range kbItems {
			kbContext += fmt.Sprintf("\n【%s】问：%s 答：%s", item.Title, item.Question, item.Answer)
		}
		lastIdx := len(messages) - 1
		messages[lastIdx].Content += kbContext
	}

	// 规则示例回复
	if rule != nil && rule.SampleReply != "" {
		lastIdx := len(messages) - 1
		messages[lastIdx].Content += "\n\n参考回复格式：" + rule.SampleReply
	}

	// 在最前面插入系统提示词
	messages = append([]AIChatMessage{{Role: "system", Content: systemPrompt}}, messages...)

	// 调用 AI（带 Function Calling）
	apiKey := s.GetConfig("ai_api_key")
	baseURL := s.GetConfig("ai_base_url")
	modelName := s.GetConfig("ai_models")
	if apiKey == "" || baseURL == "" || modelName == "" {
		return TicketAIResult{ShouldReply: false, Error: "AI 配置不完整"}
	}

	// 获取已启用的工具
	enabledTools := s.GetEnabledTools()
	openAITools := GetOpenAITools(enabledTools)

	// Function Calling 循环
	reply, confidence, err := s.callAIWithTools(baseURL, apiKey, modelName, messages, openAITools, req.TicketID, 0)
	if err != nil {
		s.log.Warnf("AI 调用失败: %v", err)
		s.db.Create(&model.AITicketProcessLog{
			TicketID:  req.TicketID,
			Event:     "auto_reply",
			Decision:  "error",
			Status:    "failed",
			Message:   err.Error(),
			ModelUsed: modelName,
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
		ModelUsed:  modelName,
	})

	s.db.Create(&model.AITicketQueue{
		TicketID:  req.TicketID,
		EventType: "auto_reply",
		Status:    "completed",
		ModelUsed: modelName,
	})

	// 使用回复服务写入工单
	adminID, _ := strconv.Atoi(s.GetConfig("reply_admin_id"))
	if adminID == 0 {
		adminID = 1
	}
	replySvc := NewTicketReplyService(s.db, s.log)
	if err := replySvc.PostAdminReply(req.TicketID, reply, uint(adminID)); err != nil {
		s.log.Warnf("AI 回复写入失败: %v", err)
	}

	return TicketAIResult{
		ShouldReply: true,
		Reply:       reply,
		Confidence:  confidence,
		Decision:    decision,
	}
}

// DefaultSystemPrompt 默认系统提示词
func DefaultSystemPrompt() string {
	return `你是锚点财务 IDC 客户工单 AI 客服助手，为客户解答问题。

## 身份与保密
- 你代表技术支持团队回复工单
- 对话仅客户、管理员与你可见，可基于下方用户与产品数据回答问题
- 语气专业、简洁、友好，使用中文

## 对话风格
- 你可以像真人客服一样自然对话，不必拘泥于一问一答
- 如果需要多条信息要告知客户，可以分多条消息发送
- 主动关怀：如果发现客户有问题但没明确提问，可以主动询问
- 全程禁止使用固定模板话术，所有提示语由你自主智能生成

## 能力范围
- 根据工单标题、内容、部门、优先级理解客户诉求
- 结合用户账户信息与产品列表（状态、IP、端口、登录用户名/密码、到期时间、specs 配置参数等）回答问题
- 客户询问如何连接服务器时，必须直接写出 IP、端口、用户名、密码明文及连接方式（SSH/RDP）
- 常见问题：续费、开机/关机、网络、配置、账单、密码重置指引等
- 无法确定或需人工操作时，明确说明已转交人工

## 工具使用规则
- 优先使用工具获取数据，不要编造数据
- 需要操作服务器前，先确认客户意图
- 高危操作（重装系统等）必须先获得客户明确确认
- 不要在回复中提及工具名称或技术细节，用自然语言描述操作
- 首次服务时，调用 list_available_tools 查询当前可用工具

## 跨部门分流规则
- 售前部门：产品咨询、套餐询价、下单购机
- 技术部门：服务器运维、故障、网络、重装
- 财务部门：退款、账单、费用、发票
- 识别到部门错配后使用 transfer_ticket_department 自动迁移

## 上游透传规则
适用场景：服务器宕机、机房网络中断、硬件故障、新机器未开通、上游卡单
执行规则：
1. 自动识别上游渠道
2. 整理客户诉求，生成规范工单
3. 上游工单仅提交问题/退款诉求，绝不填写任何金额
4. 不擅自承诺处理结果

## 上游回复处理规则
- 禁止直接转发上游原始回复
- 需对上游回复进行脱敏、简化、通俗化、润色改写
- 隐藏上游内部工单号和系统名称

## 转人工规则
禁止转接：开关机、重启、状态查询、产品咨询、基础故障排查
允许转接：机房深层硬件故障、退款争议、复杂纠纷
转接后 AI 停止主动应答，静默监听；若客户问题降级为简单问题，自动切回 AI 接管

## 退款规则（铁律）
1. 仅新开通订单 + 24小时以内允许退款
2. 续费订单、超时订单直接驳回
3. 上游工单只提交退款诉求，不填写金额
4. 金额核算由本机处理，以客户实付金额为准

## 限制
- 不要编造不存在的产品或操作结果
- 登录信息仅来自上下文，没有则说明需人工核实
- 不要承诺具体完成时间
- 回复使用 Markdown 排版
- 不要直接输出 HTML 标签

## 上下文数据（JSON）
{context}`
}

// callAIWithTools 带 Function Calling 的 AI 调用（含多轮循环）
func (s *AITicketCoreService) callAIWithTools(baseURL, apiKey, modelName string, messages []AIChatMessage, tools []map[string]interface{}, ticketID uint, round int) (string, float64, error) {
	if round >= MaxToolRounds {
		return "", 0, fmt.Errorf("工具调用轮次超限 (%d)", MaxToolRounds)
	}

	// 调用 AI API
	aiResp, err := s.callAIChatCompletion(baseURL, apiKey, modelName, messages, tools)
	if err != nil {
		return "", 0, err
	}

	if len(aiResp.Choices) == 0 {
		return "", 0, fmt.Errorf("AI 无回复")
	}

	choice := aiResp.Choices[0]

	// 如果 AI 返回 tool_calls
	if choice.FinishReason == "tool_calls" && len(choice.Message.ToolCalls) > 0 {
		// 追加 assistant message（含 tool_calls）到消息列表
		assistantMsg := AIChatMessage{
			Role:      "assistant",
			ToolCalls: choice.Message.ToolCalls,
		}
		messages = append(messages, assistantMsg)

		// 初始化工具执行器
		executor := NewToolExecutor(s.db, s.log, ticketID, 0, 0)

		// 执行每个工具调用
		for _, tc := range choice.Message.ToolCalls {
			fnName := tc.Function.Name
			fnArgs := make(map[string]interface{})
			if tc.Function.Arguments != "" {
				json.Unmarshal([]byte(tc.Function.Arguments), &fnArgs)
			}

			// 检查工具是否启用
			toolResult := ""
			if !isToolEnabled(fnName, s.GetEnabledTools()) {
				b, _ := json.Marshal(map[string]interface{}{"error": "工具未启用: " + fnName})
				toolResult = string(b)
			} else {
				toolResult = executor.Execute(fnName, fnArgs)
			}

			// 追加 tool result 消息
			messages = append(messages, AIChatMessage{
				Role:       "tool",
				ToolCallID: tc.ID,
				Content:    toolResult,
			})
		}

		// 保存工具调用日志
		executor.SaveCallLog()

		// 递归调用，让 AI 基于工具结果生成最终回复
		return s.callAIWithTools(baseURL, apiKey, modelName, messages, tools, ticketID, round+1)
	}

	// AI 返回普通文本
	reply := strings.TrimSpace(choice.Message.Content)
	if reply == "" {
		return "", 0, fmt.Errorf("AI 返回空内容")
	}

	return reply, 0.85, nil
}

// callAIChatCompletion 调用 AI Chat Completion API
func (s *AITicketCoreService) callAIChatCompletion(baseURL, apiKey, modelName string, messages []AIChatMessage, tools []map[string]interface{}) (*AIChatResponse, error) {
	endpoint := strings.TrimRight(baseURL, "/") + "/chat/completions"

	payload := map[string]interface{}{
		"model":       modelName,
		"messages":    messages,
		"temperature": 0.5,
	}

	// 注入工具定义
	if len(tools) > 0 {
		payload["tools"] = tools
		payload["tool_choice"] = "auto"
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("请求构建失败: %v", err)
	}

	client := &http.Client{Timeout: 120 * time.Second}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("连接失败: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API 错误 %d: %s", resp.StatusCode, string(body))
	}

	var result AIChatResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &result, nil
}

// isToolEnabled 检查工具是否在已启用列表中
func isToolEnabled(toolName string, enabledTools []string) bool {
	for _, name := range enabledTools {
		if name == toolName {
			return true
		}
	}
	return false
}

// allToolsMap 返回所有工具名 → 定义的映射（内部辅助）
func allToolsMapLocal() map[string]ToolDefinition {
	return allToolsMap()
}

// ─── 监控模式（人工工单 AI 静默监听 + 自动接管） ───

// ProcessMonitoring 处理监控模式：AI 评估人工工单中客户新消息是否简单到可以自动接管
// 移植自 mianyu_ai_ticket 的 enqueueMonitoring 逻辑
func (s *AITicketCoreService) ProcessMonitoring(ticketID uint, content string) MonitoringResult {
	if s.GetConfig("ai_enabled") != "1" {
		return MonitoringResult{Action: "skip"}
	}

	// 检查工单是否处于人工模式
	mode := s.GetTicketMode(ticketID)
	if mode != "human" {
		return MonitoringResult{Action: "skip"}
	}

	// 构建监控消息
	ctxSvc := NewTicketContextService(s.db, s.log)
	contextJson := ctxSvc.BuildContextJson(ticketID, 0)
	messages := ctxSvc.BuildMessages(ticketID)
	if len(messages) == 0 {
		return MonitoringResult{Action: "skip"}
	}

	// 追加监控指令
	messages = append(messages, AIChatMessage{
		Role:    "system",
		Content: MonitoringPrompt(),
	})

	// 系统提示词
	systemPrompt := s.GetConfig("system_prompt")
	if systemPrompt == "" {
		systemPrompt = DefaultSystemPrompt()
	}
	if strings.Contains(systemPrompt, "{context}") {
		systemPrompt = strings.Replace(systemPrompt, "{context}", contextJson, 1)
	} else {
		systemPrompt += "\n\n## 上下文数据（JSON）\n" + contextJson
	}
	messages = append([]AIChatMessage{{Role: "system", Content: systemPrompt}}, messages...)

	// 调用 AI
	apiKey := s.GetConfig("ai_api_key")
	baseURL := s.GetConfig("ai_base_url")
	modelName := s.GetConfig("ai_models")
	if apiKey == "" || baseURL == "" || modelName == "" {
		return MonitoringResult{Action: "skip"}
	}

	enabledTools := s.GetEnabledTools()
	openAITools := GetOpenAITools(enabledTools)

	reply, _, err := s.callAIWithTools(baseURL, apiKey, modelName, messages, openAITools, ticketID, 0)
	if err != nil {
		s.log.Warnf("监控模式 AI 调用失败: %v", err)
		return MonitoringResult{Action: "skip"}
	}

	// 判断 AI 是否决定沉默
	if isStaySilent(reply) {
		return MonitoringResult{Action: "silent"}
	}

	// AI 决定接管：切回 AI 模式
	s.SetTicketMode(ticketID, "ai")

	// 发送接管通知
	s.SendAITakeoverNotification(ticketID)

	// 写入回复
	adminID, _ := strconv.Atoi(s.GetConfig("reply_admin_id"))
	if adminID == 0 {
		adminID = 1
	}
	replySvc := NewTicketReplyService(s.db, s.log)
	replySvc.PostAdminReply(ticketID, reply, uint(adminID))

	return MonitoringResult{Action: "takeover", Reply: reply}
}

// MonitoringResult 监控模式结果
type MonitoringResult struct {
	Action string `json:"action"` // skip / silent / takeover
	Reply  string `json:"reply,omitempty"`
}

// MonitoringPrompt 监控模式指令
func MonitoringPrompt() string {
	return `---
【监控模式指令】
当前工单处于人工模式，由人工客服负责回复。你现在作为后台监控 AI，需要评估客户最新一条消息：

**判断标准：**
- 如果这是一个简单的、你有能力独立解决的常规问题（如产品信息查询、续费指引、常见技术问题等），你可以选择接管回复
- 如果这是一个复杂问题、涉及争议、需要人工判断、或涉及上游/退款等你无法独立处理的问题，请保持沉默

**输出规则：**
- 如果你决定保持沉默（不接管），只回复 [STAY_SILENT]，不输出任何其他内容
- 如果你决定接管，直接输出你对客户的回复内容（正常回答即可，系统会自动将工单切回AI模式）
- 不要解释你的判断理由`
}

// isStaySilent 判断 AI 回复是否表示"保持沉默"
func isStaySilent(text string) bool {
	text = strings.TrimSpace(text)
	return strings.ToUpper(text) == "[STAY_SILENT]" || strings.HasPrefix(strings.ToUpper(text), "[STAY_SILENT]")
}

// ─── AI 接管通知（动态开场白） ───

// SendAITakeoverNotification AI 首次回复工单前，动态生成转接提示消息
// 移植自 mianyu_ai_ticket 的 QueueService.sendAiTakeoverNotification
func (s *AITicketCoreService) SendAITakeoverNotification(ticketID uint) {
	// 检查是否已经发送过 AI 接管通知
	var count int64
	s.db.Table("ticket_reply").Where("tid = ? AND editor = ?", ticketID, "ai").Count(&count)
	if count > 0 {
		return // 已有 AI 回复，不需要再发通知
	}

	// 获取工单标题
	var ticket struct {
		Title string `gorm:"column:title"`
	}
	s.db.Table("ticket").Where("id = ?", ticketID).Select("title").First(&ticket)
	ticketTitle := ticket.Title
	if ticketTitle == "" {
		ticketTitle = "用户咨询"
	}

	// 调用 AI 生成动态开场白
	apiKey := s.GetConfig("ai_api_key")
	baseURL := s.GetConfig("ai_base_url")
	modelName := s.GetConfig("ai_models")
	if apiKey == "" || baseURL == "" || modelName == "" {
		return
	}

	messages := []AIChatMessage{
		{Role: "system", Content: "你是一个AI智能客服。请生成一段简短友好的开场白（1-2句话），告诉用户你将为他处理工单。语气自然亲切，不要用固定模板，每次都要不同。直接输出消息内容，不要加任何前缀或说明。"},
		{Role: "user", Content: "工单标题: " + ticketTitle},
	}

	reply, _, err := s.callAIWithTools(baseURL, apiKey, modelName, messages, nil, ticketID, 0)
	if err != nil || strings.TrimSpace(reply) == "" {
		return
	}

	// 获取管理员信息
	adminID, _ := strconv.Atoi(s.GetConfig("reply_admin_id"))
	if adminID == 0 {
		adminID = 1
	}

	// 写入开场白
	s.db.Table("ticket_reply").Create(map[string]interface{}{
		"tid":         ticketID,
		"uid":         0,
		"create_time": time.Now().Unix(),
		"content":     "<p>" + reply + "</p>",
		"admin_id":    adminID,
		"admin":       "AI客服",
		"attachment":  "",
		"editor":      "ai",
	})

	// 更新工单状态
	s.db.Table("ticket").Where("id = ?", ticketID).Updates(map[string]interface{}{
		"last_reply_time": time.Now().Unix(),
		"client_unread":   1,
	})
}

// ─── 工具调用日志写入工单内部备注 ───

// SaveToolCallLogAsNote 将工具调用日志保存为工单内部备注（仅管理员可见）
// 移植自 mianyu_ai_ticket 的 QueueService.saveToolCallLog
func (s *AITicketCoreService) SaveToolCallLogAsNote(ticketID uint, callLog []ToolCallEntry) {
	if len(callLog) == 0 {
		return
	}

	var lines []string
	for _, entry := range callLog {
		status := "OK"
		if !entry.Success {
			status = "FAIL"
		}
		argsJSON, _ := json.Marshal(entry.Args)
		argsStr := string(argsJSON)
		if len(argsStr) > 200 {
			argsStr = argsStr[:200]
		}
		lines = append(lines, fmt.Sprintf("[%s] %s (%dms) args=%s", status, entry.Tool, entry.Elapsed, argsStr))
	}

	adminID, _ := strconv.Atoi(s.GetConfig("reply_admin_id"))
	if adminID == 0 {
		adminID = 1
	}

	content := "<p><em>[AI工具调用日志]</em></p><pre>" + strings.Join(lines, "\n") + "</pre>"

	s.db.Table("ticket_reply").Create(map[string]interface{}{
		"tid":         ticketID,
		"uid":         0,
		"create_time": time.Now().Unix(),
		"content":     content,
		"admin_id":    adminID,
		"admin":       "系统",
		"attachment":  "",
		"editor":      "plain",
	})
}

// ─── 自动关单 ───

// AutoCloseTicket AI 自动关闭已解决的工单
// 移植自 mianyu_ai_ticket 的 auto_close_ticket 工具逻辑
func (s *AITicketCoreService) AutoCloseTicket(ticketID uint, reason string) error {
	// 更新工单状态为已关闭
	err := s.db.Table("ticket").Where("id = ?", ticketID).Updates(map[string]interface{}{
		"status":       "closed",
		"close_time":   time.Now().Unix(),
		"last_reply_time": time.Now().Unix(),
	}).Error
	if err != nil {
		return err
	}

	// 记录关闭原因
	if reason != "" {
		adminID, _ := strconv.Atoi(s.GetConfig("reply_admin_id"))
		if adminID == 0 {
			adminID = 1
		}
		s.db.Table("ticket_reply").Create(map[string]interface{}{
			"tid":         ticketID,
			"uid":         0,
			"create_time": time.Now().Unix(),
			"content":     "<p><em>" + reason + "</em></p>",
			"admin_id":    adminID,
			"admin":       "AI客服",
			"attachment":  "",
			"editor":      "ai",
		})
	}

	// 记录日志
	s.db.Create(&model.AITicketProcessLog{
		TicketID: ticketID,
		Event:    "auto_close",
		Decision: "close",
		Status:   "success",
		Message:  reason,
	})

	return nil
}

// ─── 跨部门分流 ───

// TransferTicketDepartment AI 将工单转移到其他部门
// 移植自 mianyu_ai_ticket 的 transfer_ticket_department 工具逻辑
func (s *AITicketCoreService) TransferTicketDepartment(ticketID uint, targetDeptID uint, reason string) error {
	// 获取原工单信息
	var ticket struct {
		DeptID  uint   `gorm:"column:dptid"`
		Title   string `gorm:"column:title"`
		Content string `gorm:"column:content"`
		UID     uint   `gorm:"column:uid"`
	}
	s.db.Table("ticket").Where("id = ?", ticketID).Select("dptid, title, content, uid").First(&ticket)

	if ticket.DeptID == targetDeptID {
		return nil // 已在目标部门
	}

	// 记录转移日志
	s.db.Create(&model.AITicketProcessLog{
		TicketID: ticketID,
		Event:    "transfer_dept",
		Decision: fmt.Sprintf("%d->%d", ticket.DeptID, targetDeptID),
		Status:   "success",
		Message:  reason,
	})

	// 更新工单部门
	err := s.db.Table("ticket").Where("id = ?", ticketID).Updates(map[string]interface{}{
		"dptid":          targetDeptID,
		"last_reply_time": time.Now().Unix(),
	}).Error
	if err != nil {
		return err
	}

	// 写入转移备注
	adminID, _ := strconv.Atoi(s.GetConfig("reply_admin_id"))
	if adminID == 0 {
		adminID = 1
	}

	transferNote := fmt.Sprintf("<p><em>[AI部门转移] 从部门 %d 转移到部门 %d</em></p><p>原因：%s</p>", ticket.DeptID, targetDeptID, reason)
	s.db.Table("ticket_reply").Create(map[string]interface{}{
		"tid":         ticketID,
		"uid":         0,
		"create_time": time.Now().Unix(),
		"content":     transferNote,
		"admin_id":    adminID,
		"admin":       "AI客服",
		"attachment":  "",
		"editor":      "ai",
	})

	return nil
}

// ─── 队列恢复（处理超时任务） ───

// RecoverStuckJobs 恢复卡住的队列任务
// 移植自 mianyu_ai_ticket 的 QueueService.recoverStuckJobs
func (s *AITicketCoreService) RecoverStuckJobs(staleSeconds int) {
	if staleSeconds <= 0 {
		staleSeconds = 90
	}
	threshold := time.Now().Add(-time.Duration(staleSeconds) * time.Second)
	s.db.Model(&model.AITicketQueue{}).
		Where("status = ? AND updated_at < ?", "processing", threshold).
		Updates(map[string]interface{}{
			"status":     "pending",
			"error_msg":  "处理超时已自动恢复",
			"updated_at": time.Now(),
		})
}
