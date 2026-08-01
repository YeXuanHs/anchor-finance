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

// AIShoppingCoreService AI购物助手核心服务
// 移植自 mahiru_ai_shopping
type AIShoppingCoreService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewAIShoppingCoreService(db *gorm.DB, log *logger.Logger) *AIShoppingCoreService {
	return &AIShoppingCoreService{db: db, log: log}
}

// ─── 配置管理 ───

func (s *AIShoppingCoreService) GetConfig(key string) string {
	var cfg model.AIShoppingConfig
	if err := s.db.Where("`key` = ?", key).First(&cfg).Error; err != nil {
		return ""
	}
	return cfg.Value
}

func (s *AIShoppingCoreService) SetConfig(key, value string) error {
	var cfg model.AIShoppingConfig
	if err := s.db.Where("`key` = ?", key).First(&cfg).Error; err != nil {
		cfg = model.AIShoppingConfig{Key: key, Value: value}
		return s.db.Create(&cfg).Error
	}
	cfg.Value = value
	return s.db.Save(&cfg).Error
}

// GetAllConfig 获取所有配置
func (s *AIShoppingCoreService) GetAllConfig() map[string]interface{} {
	enabled := s.GetConfig("ai_enabled")
	if enabled == "" {
		enabled = "0"
	}
	baseURL := s.GetConfig("ai_base_url")
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	modelName := s.GetConfig("ai_model")
	if modelName == "" {
		modelName = "gpt-4o-mini"
	}
	widgetTitle := s.GetConfig("widget_title")
	if widgetTitle == "" {
		widgetTitle = "AI导购"
	}
	welcome := s.GetConfig("welcome_message")
	if welcome == "" {
		welcome = "你好，我是 AI 导购，可以帮你对比商品、选配置、估算价格。直接告诉我你的需求即可。"
	}
	return map[string]interface{}{
		"ai_enabled":      enabled,
		"ai_base_url":     baseURL,
		"ai_model":        modelName,
		"widget_title":    widgetTitle,
		"welcome_message": welcome,
		"has_api_key":     s.GetConfig("ai_api_key") != "",
	}
}

// SaveAllConfig 保存所有配置
func (s *AIShoppingCoreService) SaveAllConfig(data map[string]interface{}) error {
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

// ─── 商品搜索（工具调用） ───

// SearchProducts 搜索商品（供 AI 工具调用）
func (s *AIShoppingCoreService) SearchProducts(keyword string, limit int) []map[string]interface{} {
	type product struct {
		ID          uint    `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}
	var products []product
	q := s.db.Table("products").Select("id, name, description, price").Where("status = ?", 1)
	if keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("name LIKE ? OR description LIKE ?", like, like)
	}
	q.Limit(limit).Find(&products)
	var result []map[string]interface{}
	for _, p := range products {
		result = append(result, map[string]interface{}{
			"id":          p.ID,
			"name":        p.Name,
			"description": p.Description,
			"price":       p.Price,
		})
	}
	return result
}

// GetAllProducts 获取所有商品
func (s *AIShoppingCoreService) GetAllProducts(limit int) []map[string]interface{} {
	return s.SearchProducts("", limit)
}

// ─── 聊天会话 ───

func (s *AIShoppingCoreService) Chat(sessionID string, userID uint, message string) (string, error) {
	// 保存用户消息
	s.db.Create(&model.AIShoppingChatLog{
		SessionID: sessionID,
		UserID:    userID,
		Role:      "user",
		Content:   message,
	})

	// 获取历史消息
	var history []model.AIShoppingChatLog
	s.db.Where("session_id = ?", sessionID).Order("id DESC").Limit(10).Find(&history)

	// 获取配置
	apiKey := s.GetConfig("ai_api_key")
	baseURL := s.GetConfig("ai_base_url")
	modelName := s.GetConfig("ai_model")
	systemPrompt := s.GetConfig("system_prompt")

	if apiKey == "" || baseURL == "" {
		return "抱歉，AI 购物助手暂未配置。", nil
	}

	if systemPrompt == "" {
		systemPrompt = `你是一个专业的AI购物助手。你的职责是：
1. 了解客户需求
2. 根据商品数据推荐合适的产品
3. 解答产品相关问题
4. 帮助客户做出购买决策
回复要简洁友好，重点突出产品优势和价格。`
	}

	// 构建消息
	messages := []map[string]string{
		{"role": "system", "content": systemPrompt},
	}
	for i := len(history) - 1; i >= 0; i-- {
		messages = append(messages, map[string]string{
			"role":    history[i].Role,
			"content": history[i].Content,
		})
	}

	reply, err := s.callAI(baseURL, apiKey, modelName, messages)
	if err != nil {
		return "抱歉，AI 暂时无法响应，请稍后重试。", nil
	}

	// 保存 AI 回复
	s.db.Create(&model.AIShoppingChatLog{
		SessionID: sessionID,
		UserID:    userID,
		Role:      "assistant",
		Content:   reply,
	})

	return reply, nil
}

func (s *AIShoppingCoreService) GetChatHistory(sessionID string) []model.AIShoppingChatLog {
	var messages []model.AIShoppingChatLog
	s.db.Where("session_id = ?", sessionID).Order("id ASC").Find(&messages)
	return messages
}

func (s *AIShoppingCoreService) callAI(baseURL, apiKey, modelName string, messages []map[string]string) (string, error) {
	endpoint := strings.TrimRight(baseURL, "/") + "/chat/completions"
	reqBody := map[string]interface{}{
		"model":       modelName,
		"messages":    messages,
		"max_tokens":  2000,
		"temperature": 0.7,
	}
	jsonData, _ := json.Marshal(reqBody)
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API 错误 %d: %s", resp.StatusCode, string(body))
	}
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("无回复")
	}
	return result.Choices[0].Message.Content, nil
}
