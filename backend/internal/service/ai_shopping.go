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

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AIShoppingService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewAIShoppingService(db *gorm.DB, log *logger.Logger) *AIShoppingService {
	return &AIShoppingService{db: db, log: log}
}

// ─── 配置管理 ───

func (s *AIShoppingService) GetConfig() *model.AIShoppingAssistantConfig {
	var config model.AIShoppingAssistantConfig
	if err := s.db.First(&config).Error; err != nil {
		return &model.AIShoppingAssistantConfig{
			Enabled:             false,
			WelcomeMessage:      "您好！我是AI购物助手，可以帮您推荐合适的产品和服务。请问您有什么需求？",
			MaxRecommendations:  5,
			IncludePricing:      true,
			ShowOnAllPages:      true,
		}
	}
	return &config
}

func (s *AIShoppingService) SaveConfig(config *model.AIShoppingAssistantConfig) error {
	var existing model.AIShoppingAssistantConfig
	if err := s.db.First(&existing).Error; err != nil {
		return s.db.Create(config).Error
	}
	config.ID = existing.ID
	config.UpdatedAt = time.Now()
	return s.db.Save(config).Error
}

// ─── 商品目录配置 ───

func (s *AIShoppingService) GetCatalogConfig() *model.ProductCatalogConfig {
	var config model.ProductCatalogConfig
	if err := s.db.First(&config).Error; err != nil {
		return &model.ProductCatalogConfig{
			LayoutStyle:      "grid",
			ShowFilters:      true,
			ShowComparison:   true,
			ShowReviews:      true,
			ShowTechSpecs:    true,
			EnableSort:       true,
			DefaultSort:      "recommend",
			ProductsPerPage:  12,
			ShowCategoryTree: true,
		}
	}
	return &config
}

func (s *AIShoppingService) SaveCatalogConfig(config *model.ProductCatalogConfig) error {
	var existing model.ProductCatalogConfig
	if err := s.db.First(&existing).Error; err != nil {
		return s.db.Create(config).Error
	}
	config.ID = existing.ID
	return s.db.Save(config).Error
}

// ─── AI 购物助手对话 ───

// StartSession 开始新的购物助手会话
func (s *AIShoppingService) StartSession(userID uint) (string, error) {
	sessionID := uuid.New().String()
	session := model.AIShoppingSession{
		UserID:    userID,
		SessionID: sessionID,
		Status:    "active",
	}
	if err := s.db.Create(&session).Error; err != nil {
		return "", err
	}

	// 发送欢迎消息
	config := s.GetConfig()
	welcomeMsg := model.AIShoppingMessage{
		SessionID: sessionID,
		Role:      "assistant",
		Content:   config.WelcomeMessage,
	}
	s.db.Create(&welcomeMsg)

	return sessionID, nil
}

// SendMessage 发送消息并获取 AI 回复
func (s *AIShoppingService) SendMessage(sessionID string, userMessage string) (string, []ProductRecommendation, error) {
	// 验证会话
	var session model.AIShoppingSession
	if err := s.db.Where("session_id = ? AND status = ?", sessionID, "active").First(&session).Error; err != nil {
		return "", nil, fmt.Errorf("会话不存在或已关闭")
	}

	// 保存用户消息
	userMsg := model.AIShoppingMessage{
		SessionID: sessionID,
		Role:      "user",
		Content:   userMessage,
	}
	s.db.Create(&userMsg)

	// 获取历史消息（最近10条）
	var history []model.AIShoppingMessage
	s.db.Where("session_id = ?", sessionID).Order("id DESC").Limit(10).Find(&history)

	// 获取产品信息用于推荐
	products := s.getAvailableProducts()

	// 获取 AI 配置
	config := s.GetConfig()
	var aiConfig model.AIConfig
	if config.AIConfigID > 0 {
		s.db.First(&aiConfig, config.AIConfigID)
	}

	if !config.Enabled || aiConfig.ID == 0 || !aiConfig.IsActive {
		// AI 未启用，返回默认回复
		return "抱歉，AI 购物助手暂时不可用。请浏览我们的产品页面或联系客服获取帮助。", nil, nil
	}

	// 构建 prompt
	systemPrompt := config.SystemPrompt
	if systemPrompt == "" {
		systemPrompt = fmt.Sprintf(`你是一个专业的购物助手。你的职责是：
1. 了解客户需求
2. 推荐合适的产品
3. 解答产品相关问题
4. 帮助客户做出购买决策

可用产品列表：
%s

请根据客户需求推荐产品。推荐时请说明推荐理由。`, products)
	}

	// 构建消息历史
	var messages []openaiMessage
	messages = append(messages, openaiMessage{Role: "system", Content: systemPrompt})
	for i := len(history) - 1; i >= 0; i-- {
		messages = append(messages, openaiMessage{
			Role:    history[i].Role,
			Content: history[i].Content,
		})
	}

	// 调用 AI
	reply, recommendations, err := s.callShoppingAI(&aiConfig, messages, config.MaxRecommendations)
	if err != nil {
		s.log.Warnf("AI 购物助手调用失败: %v", err)
		return "抱歉，AI 助手暂时无法响应。请稍后重试或联系客服。", nil, nil
	}

	// 保存 AI 回复
	productsJSON := ""
	if len(recommendations) > 0 {
		data, _ := json.Marshal(recommendations)
		productsJSON = string(data)
	}
	aiMsg := model.AIShoppingMessage{
		SessionID: sessionID,
		Role:      "assistant",
		Content:   reply,
		Products:  productsJSON,
	}
	s.db.Create(&aiMsg)

	return reply, recommendations, nil
}

// CloseSession 关闭会话
func (s *AIShoppingService) CloseSession(sessionID string) error {
	now := time.Now()
	return s.db.Model(&model.AIShoppingSession{}).
		Where("session_id = ?", sessionID).
		Updates(map[string]interface{}{
			"status":    "closed",
			"closed_at": &now,
		}).Error
}

// GetSessionMessages 获取会话消息
func (s *AIShoppingService) GetSessionMessages(sessionID string) ([]model.AIShoppingMessage, error) {
	var messages []model.AIShoppingMessage
	if err := s.db.Where("session_id = ?", sessionID).Order("id ASC").Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

// ─── 产品推荐 ───

type ProductRecommendation struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Reason      string  `json:"reason"`
}

func (s *AIShoppingService) getAvailableProducts() string {
	type productInfo struct {
		Name        string
		Description string
		Price       float64
	}
	var products []productInfo
	s.db.Table("products").Select("name, description, price").Where("status = ?", 1).Limit(50).Find(&products)

	var parts []string
	for _, p := range products {
		parts = append(parts, fmt.Sprintf("- %s: %s (¥%.2f)", p.Name, p.Description, p.Price))
	}
	return strings.Join(parts, "\n")
}

type shoppingAIRequest struct {
	Model       string          `json:"model"`
	Messages    []openaiMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float64         `json:"temperature"`
}

type shoppingAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (s *AIShoppingService) callShoppingAI(config *model.AIConfig, messages []openaiMessage, maxRecs int) (string, []ProductRecommendation, error) {
	endpoint := config.APIEndpoint
	if endpoint == "" {
		endpoint = "https://api.openai.com/v1/chat/completions"
	}

	reqBody := shoppingAIRequest{
		Model:       config.Model,
		Messages:    messages,
		MaxTokens:   config.MaxTokens,
		Temperature: config.Temperature,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", nil, err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("AI API 请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("AI API 返回错误 %d: %s", resp.StatusCode, string(body))
	}

	var result shoppingAIResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", nil, fmt.Errorf("解析 AI 响应失败: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", nil, fmt.Errorf("AI 未返回任何回复")
	}

	reply := result.Choices[0].Message.Content

	// 尝试从回复中提取产品推荐（简单实现）
	recommendations := s.extractRecommendations(reply)

	return reply, recommendations, nil
}

// extractRecommendations 从 AI 回复中提取产品推荐
func (s *AIShoppingService) extractRecommendations(reply string) []ProductRecommendation {
	var recs []ProductRecommendation

	// 查询所有产品
	type product struct {
		ID          uint
		Name        string
		Description string
		Price       float64
	}
	var products []product
	s.db.Table("products").Select("id, name, description, price").Where("status = ?", 1).Find(&products)

	// 匹配回复中提到的产品名
	lowerReply := strings.ToLower(reply)
	for _, p := range products {
		if strings.Contains(lowerReply, strings.ToLower(p.Name)) {
			recs = append(recs, ProductRecommendation{
				ID:          p.ID,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
				Reason:      "AI 推荐",
			})
		}
	}

	return recs
}
