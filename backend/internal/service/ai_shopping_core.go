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

// ChatRequest 聊天请求
type ChatRequest struct {
	SessionID   string `json:"session_id"`
	UserID      uint   `json:"user_id"`
	Message     string `json:"message"`
	PageContext string `json:"page_context"` // 页面上下文（当前浏览的分组/商品）
}

// Chat 聊天（支持页面上下文和多消息分段）
func (s *AIShoppingCoreService) Chat(sessionID string, userID uint, message string, pageContext string) (string, error) {
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
		systemPrompt = DefaultShoppingPrompt()
	}

	// 注入页面上下文
	if pageContext != "" {
		systemPrompt += "\n\n## 用户当前浏览上下文\n" + pageContext
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

	// 处理多消息分段（[[[MSG]]] 分隔符）
	reply = s.formatMultiMessage(reply)

	// 保存 AI 回复
	s.db.Create(&model.AIShoppingChatLog{
		SessionID: sessionID,
		UserID:    userID,
		Role:      "assistant",
		Content:   reply,
	})

	return reply, nil
}

// formatMultiMessage 处理多消息分段格式
// 移植自 mahiru_ai_shopping 的多消息格式：用 [[[MSG]]] 分隔多条消息
func (s *AIShoppingCoreService) formatMultiMessage(reply string) string {
	// 如果包含分隔符，保持原样（前端会解析）
	if strings.Contains(reply, "[[[MSG]]]") {
		return reply
	}
	return reply
}

// DefaultShoppingPrompt 默认导购提示词
// 移植自 mahiru_ai_shopping 的 PromptDefaults
func DefaultShoppingPrompt() string {
	return `你是购物车里的 AI 导购助手，帮助用户选商品、对比配置、了解价格。
回答使用简洁、友好的中文，可用 Markdown（列表、加粗、链接、表格）。
像真人导购一样分多条消息回复。**必须**用单独一行 [[[MSG]]] 分隔每条消息，前后各空一行。示例格式：

好的，帮您查一下～

[[[MSG]]]

找到了！这款很适合您：

[商品名](链接)

[[[MSG]]]

还有其他需求随时告诉我哦

每段 1-3 句，不要把所有内容挤在一大段里。不要用列表把多条消息写在一起。不要用 --- 分隔。

## 核心约束（必须遵守）
- 面向用户的回复中严禁出现 product_id、商品 ID、内部编号等标识符
- 推荐商品必须使用 Markdown 可点击链接：[短名称](order_url)
- 禁止在 [ ] 里换行；禁止在链接标签里写 |（竖线）
- 禁止只写商品名称而不给链接；禁止编造不存在的商品、价格或配置

## 商品工具（按需调用，禁止臆造）
- 当你需要查询商品、对比配置或价格时，必须先用工具，不要编造不存在的套餐/价格
- 搜索关键词规则：只传短词，如「香港」「HK」「轻量」「CN2」
- 不要一次性索取全站商品；只在用户明确需要时搜索/取详情
- 工具调用节制：同一需求最多搜索 1～2 次

## 其它
- 若有「用户当前浏览上下文」，可优先结合当前分组作推荐线索，但仍须用工具核实商品与价格
- 充值、购买、找回密码等平台操作不确定时，建议用户提交工单或联系人工`
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
