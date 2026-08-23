package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/YeXuanHs/anchor-finance/internal/pluginengine"
)

// AIService AI服务（可选启用，配置存数据库）
type AIService struct {
	enabled  bool
	provider string // openai, deepseek, etc.
	apiKey   string
	baseURL  string
	model    string
	client   *http.Client
}

// NewAIService 从数据库读取AI配置并初始化
func NewAIService() *AIService {
	s := &AIService{
		client: &http.Client{Timeout: 60 * time.Second},
	}
	s.loadConfig()
	return s
}

// loadConfig 从settings表加载AI配置（按key查询，兼容所有group）
func (s *AIService) loadConfig() {
	db := database.GetDB()

	keys := []string{"ai_enabled", "ai_provider", "ai_api_key", "ai_base_url", "ai_model"}
	configMap := make(map[string]string)
	for _, key := range keys {
		var setting model.Setting
		if err := db.Where("`key` = ?", key).First(&setting).Error; err == nil {
			configMap[key] = setting.Value
		}
	}

	if enabled, ok := configMap["ai_enabled"]; ok && enabled == "1" {
		s.enabled = true
	} else {
		s.enabled = false
		return
	}

	s.provider = configMap["ai_provider"]
	if s.provider == "" {
		s.provider = "openai"
	}
	s.apiKey = configMap["ai_api_key"]
	s.baseURL = configMap["ai_base_url"]
	if s.baseURL == "" {
		s.baseURL = "https://api.openai.com/v1"
	}
	s.model = configMap["ai_model"]
	if s.model == "" {
		s.model = "gpt-3.5-turbo"
	}
}

// IsEnabled 是否启用AI
func (s *AIService) IsEnabled() bool {
	return s.enabled
}

// GetConfig 获取AI配置（脱敏）
func (s *AIService) GetConfig() map[string]interface{} {
	maskedKey := ""
	if len(s.apiKey) > 8 {
		maskedKey = s.apiKey[:4] + "****" + s.apiKey[len(s.apiKey)-4:]
	}
	return map[string]interface{}{
		"enabled":   s.enabled,
		"provider":  s.provider,
		"api_key":   maskedKey,
		"base_url":  s.baseURL,
		"model":     s.model,
	}
}

// ChatCompletion 调用AI对话接口（支持function calling）
func (s *AIService) ChatCompletion(messages []map[string]string, tools []map[string]interface{}) (string, error) {
	if !s.enabled {
		return "", fmt.Errorf("AI服务未启用")
	}

	payload := map[string]interface{}{
		"model":    s.model,
		"messages": messages,
	}
	if len(tools) > 0 {
		payload["tools"] = tools
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", s.baseURL+"/chat/completions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("AI请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("AI返回错误: %d %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content   string                   `json:"content"`
				ToolCalls []map[string]interface{} `json:"tool_calls,omitempty"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("解析AI响应失败: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("AI无响应")
	}

	// 如果AI返回tool_calls，执行工具调用后继续对话
	msg := result.Choices[0].Message
	if len(msg.ToolCalls) > 0 {
		return s.handleToolCalls(messages, msg.ToolCalls)
	}

	return msg.Content, nil
}

// handleToolCalls 执行AI工具调用后继续对话
func (s *AIService) handleToolCalls(messages []map[string]string, toolCalls []map[string]interface{}) (string, error) {
	// 执行每个tool并收集结果文本
	var toolContext string
	for _, tc := range toolCalls {
		fn, _ := tc["function"].(map[string]interface{})
		funcName, _ := fn["name"].(string)
		funcArgs, _ := fn["arguments"].(string)
		result := ExecuteAITool(funcName, funcArgs)
		toolContext += fmt.Sprintf("工具[%s]结果: %s\n", funcName, result)
	}

	// 把工具结果作为assistant消息追加，再让AI总结
	messages = append(messages, map[string]string{
		"role":    "assistant",
		"content": toolContext,
	})

	return s.ChatCompletion(messages, nil)
}

// GenerateProductDescription AI生成商品简介
func (s *AIService) GenerateProductDescription(productName string, config map[string]interface{}, template string) (string, error) {
	messages := []map[string]string{
		{
			"role": "system",
			"content": `你是一个专业的IDC商品文案撰写专家。根据提供的商品信息和模板，生成专业的商品简介。
要求：
1. 严格遵循模板的格式和风格
2. 用实际商品信息替换占位符
3. 商品没有的信息不能编造
4. 输出纯HTML，可直接用于商品简介`,
		},
		{
			"role": "user",
			"content": fmt.Sprintf("模板：\n%s\n\n商品名称：%s\n商品配置：%v", template, productName, config),
		},
	}

	return s.ChatCompletion(messages, nil)
}

// TicketAutoReply AI工单自动回复
func (s *AIService) TicketAutoReply(ticketSubject string, ticketContent string, context map[string]interface{}) (string, error) {
	systemPrompt := `你是一个专业的IDC客服AI。根据客户的工单内容，生成专业、友好的回复。
要求：
1. 回复要专业、简洁、有帮助
2. 如果是技术问题，给出具体的解决步骤
3. 如果无法解决，建议转接人工客服
4. 不要编造不确定的信息`

	userPrompt := fmt.Sprintf("工单主题：%s\n工单内容：%s", ticketSubject, ticketContent)
	if ctx, ok := context["customer_info"]; ok {
		userPrompt += fmt.Sprintf("\n客户信息：%v", ctx)
	}
	if ctx, ok := context["service_info"]; ok {
		userPrompt += fmt.Sprintf("\n服务信息：%v", ctx)
	}

	messages := []map[string]string{
		{"role": "system", "content": systemPrompt},
		{"role": "user", "content": userPrompt},
	}

	return s.ChatCompletion(messages, nil)
}

// GetAITools 获取AI可用工具列表（function calling定义）
func GetAITools() []map[string]interface{} {
	return []map[string]interface{}{
		{"type": "function", "function": map[string]interface{}{
			"name": "get_user_info", "description": "获取客户信息",
			"parameters": map[string]interface{}{"type": "object", "properties": map[string]interface{}{
				"user_id": map[string]interface{}{"type": "integer", "description": "用户ID"},
			}, "required": []string{"user_id"}},
		}},
		{"type": "function", "function": map[string]interface{}{
			"name": "get_service_status", "description": "查询服务状态",
			"parameters": map[string]interface{}{"type": "object", "properties": map[string]interface{}{
				"service_id": map[string]interface{}{"type": "integer", "description": "服务ID"},
			}, "required": []string{"service_id"}},
		}},
		{"type": "function", "function": map[string]interface{}{
			"name": "reboot_service", "description": "重启客户服务器",
			"parameters": map[string]interface{}{"type": "object", "properties": map[string]interface{}{
				"service_id": map[string]interface{}{"type": "integer", "description": "服务ID"},
			}, "required": []string{"service_id"}},
		}},
		{"type": "function", "function": map[string]interface{}{
			"name": "get_user_orders", "description": "获取客户订单列表",
			"parameters": map[string]interface{}{"type": "object", "properties": map[string]interface{}{
				"user_id": map[string]interface{}{"type": "integer", "description": "用户ID"},
			}, "required": []string{"user_id"}},
		}},
		{"type": "function", "function": map[string]interface{}{
			"name": "get_user_invoices", "description": "获取客户账单列表",
			"parameters": map[string]interface{}{"type": "object", "properties": map[string]interface{}{
				"user_id": map[string]interface{}{"type": "integer", "description": "用户ID"},
			}, "required": []string{"user_id"}},
		}},
	}
}

// ExecuteAITool 执行AI工具调用
func ExecuteAITool(funcName string, argsJSON string) string {
	db := database.GetDB()
	var args map[string]interface{}
	json.Unmarshal([]byte(argsJSON), &args)

	switch funcName {
	case "get_user_info":
		userID := int(args["user_id"].(float64))
		var user model.User
		if err := db.First(&user, userID).Error; err != nil {
			return `{"error":"用户不存在"}`
		}
		r, _ := json.Marshal(map[string]interface{}{"id": user.ID, "username": user.Username, "email": user.Email, "balance": user.Balance, "status": user.Status})
		return string(r)

	case "get_service_status":
		serviceID := int(args["service_id"].(float64))
		var svc model.Service
		if err := db.First(&svc, serviceID).Error; err != nil {
			return `{"error":"服务不存在"}`
		}
		r, _ := json.Marshal(map[string]interface{}{"id": svc.ID, "status": svc.Status, "product_name": svc.ProductName, "username": svc.Username, "domain": svc.Domain})
		return string(r)

	case "reboot_service":
		serviceID := int(args["service_id"].(float64))
		var svc model.Service
		if err := db.First(&svc, serviceID).Error; err != nil {
			return `{"error":"服务不存在"}`
		}
		_, err := pluginengine.TriggerHook("service_reboot", map[string]interface{}{"service_id": svc.ID, "user_id": svc.UserID})
		if err != nil {
			return fmt.Sprintf(`{"error":"重启失败: %s"}`, err.Error())
		}
		return `{"success":true,"message":"重启命令已发送"}`

	case "get_user_orders":
		userID := int(args["user_id"].(float64))
		var orders []model.Order
		db.Where("user_id = ?", userID).Order("id DESC").Limit(10).Find(&orders)
		r, _ := json.Marshal(orders)
		return string(r)

	case "get_user_invoices":
		userID := int(args["user_id"].(float64))
		var invoices []model.Invoice
		db.Where("user_id = ?", userID).Order("id DESC").Limit(10).Find(&invoices)
		r, _ := json.Marshal(invoices)
		return string(r)

	default:
		return `{"error":"未知工具"}`
	}
}
