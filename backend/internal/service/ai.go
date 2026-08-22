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

// loadConfig 从settings表加载AI配置
func (s *AIService) loadConfig() {
	db := database.GetDB()

	var settings []model.Setting
	db.Where("`group` = ?", "ai").Find(&settings)

	configMap := make(map[string]string)
	for _, setting := range settings {
		configMap[setting.Key] = setting.Value
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

// ChatCompletion 调用AI对话接口
func (s *AIService) ChatCompletion(messages []map[string]string) (string, error) {
	if !s.enabled {
		return "", fmt.Errorf("AI服务未启用")
	}

	payload := map[string]interface{}{
		"model":    s.model,
		"messages": messages,
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
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("解析AI响应失败: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("AI无响应")
	}

	return result.Choices[0].Message.Content, nil
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

	return s.ChatCompletion(messages)
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

	return s.ChatCompletion(messages)
}
