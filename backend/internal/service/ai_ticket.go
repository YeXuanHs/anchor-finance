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

type AITicketService struct {
	db  *gorm.DB
	log *logger.Logger
	kb  *KnowledgeBaseService
}

func NewAITicketService(db *gorm.DB, log *logger.Logger, kb *KnowledgeBaseService) *AITicketService {
	return &AITicketService{db: db, log: log, kb: kb}
}

// ─── AI 配置管理 ───

func (s *AITicketService) ListAIConfigs() ([]model.AIConfig, error) {
	var configs []model.AIConfig
	if err := s.db.Order("id ASC").Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

func (s *AITicketService) GetAIConfig(id uint) (*model.AIConfig, error) {
	var config model.AIConfig
	if err := s.db.First(&config, id).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

func (s *AITicketService) SaveAIConfig(config *model.AIConfig) error {
	if config.ID == 0 {
		return s.db.Create(config).Error
	}
	config.UpdatedAt = time.Now()
	return s.db.Save(config).Error
}

func (s *AITicketService) DeleteAIConfig(id uint) error {
	return s.db.Delete(&model.AIConfig{}, id).Error
}

// ─── 自动回复配置管理 ───

func (s *AITicketService) GetAutoReplyConfig() *model.AITicketAutoReplyConfig {
	var config model.AITicketAutoReplyConfig
	if err := s.db.First(&config).Error; err != nil {
		return &model.AITicketAutoReplyConfig{
			Enabled:           false,
			ConfidenceThreshold: 0.7,
			MaxAutoReplies:    3,
			IncludeKBContent:  true,
			KBSearchLimit:     5,
			AddDisclaimer:     true,
			DisclaimerText:    "此回复由AI生成，仅供参考。如需人工帮助请回复「转人工」",
		}
	}
	return &config
}

func (s *AITicketService) SaveAutoReplyConfig(config *model.AITicketAutoReplyConfig) error {
	var existing model.AITicketAutoReplyConfig
	if err := s.db.First(&existing).Error; err != nil {
		return s.db.Create(config).Error
	}
	config.ID = existing.ID
	config.UpdatedAt = time.Now()
	return s.db.Save(config).Error
}

// ─── AI 自动回复核心逻辑 ───

type AutoReplyRequest struct {
	TicketID    uint   `json:"ticket_id"`
	Subject     string `json:"subject"`
	Content     string `json:"content"`
	DeptID      uint   `json:"dept_id"`
}

type AutoReplyResult struct {
	ShouldReply bool   `json:"should_reply"`
	Reply       string `json:"reply"`
	Confidence  float64 `json:"confidence"`
	KBMatchIDs  []uint `json:"kb_match_ids"`
	Error       string `json:"error,omitempty"`
}

// GenerateAutoReply 生成工单自动回复
func (s *AITicketService) GenerateAutoReply(req AutoReplyRequest) AutoReplyResult {
	config := s.GetAutoReplyConfig()
	if !config.Enabled {
		return AutoReplyResult{ShouldReply: false}
	}

	// 检查部门是否适用
	if config.DeptIDs != "" {
		deptIDs := strings.Split(config.DeptIDs, ",")
		found := false
		for _, id := range deptIDs {
			if strings.TrimSpace(id) == fmt.Sprintf("%d", req.DeptID) {
				found = true
				break
			}
		}
		if !found {
			return AutoReplyResult{ShouldReply: false}
		}
	}

	// 检查排除关键词
	if config.ExcludeKeywords != "" {
		keywords := strings.Split(config.ExcludeKeywords, ",")
		content := strings.ToLower(req.Subject + " " + req.Content)
		for _, kw := range keywords {
			if strings.Contains(content, strings.ToLower(strings.TrimSpace(kw))) {
				return AutoReplyResult{ShouldReply: false}
			}
		}
	}

	// 检查同一工单自动回复次数
	var replyCount int64
	s.db.Model(&model.AITicketLog{}).Where("ticket_id = ?", req.TicketID).Count(&replyCount)
	if int(replyCount) >= config.MaxAutoReplies {
		return AutoReplyResult{ShouldReply: false}
	}

	// 获取 AI 配置
	aiConfig, err := s.GetAIConfig(config.AIConfigID)
	if err != nil || !aiConfig.IsActive {
		return AutoReplyResult{ShouldReply: false, Error: "AI 配置未找到或未启用"}
	}

	// 搜索知识库
	var kbContext string
	var kbMatchIDs []uint
	if config.IncludeKBContent && s.kb != nil {
		query := req.Subject + " " + req.Content
		articles, err := s.kb.SearchForAI(query, config.KBSearchLimit)
		if err == nil && len(articles) > 0 {
			var parts []string
			for _, a := range articles {
				parts = append(parts, fmt.Sprintf("【%s】%s", a.Title, a.Summary))
				kbMatchIDs = append(kbMatchIDs, a.ID)
			}
			kbContext = "\n\n参考知识库内容：\n" + strings.Join(parts, "\n")
		}
	}

	// 构建 prompt
	prompt := fmt.Sprintf("用户工单主题：%s\n用户工单内容：%s%s", req.Subject, req.Content, kbContext)

	// 调用 AI
	reply, confidence, tokensUsed, err := s.callAI(aiConfig, prompt)
	if err != nil {
		s.log.Warnf("AI 自动回复失败: %v", err)
		return AutoReplyResult{ShouldReply: false, Error: err.Error()}
	}

	// 检查置信度
	if confidence < config.ConfidenceThreshold {
		return AutoReplyResult{ShouldReply: false, Confidence: confidence}
	}

	// 添加声明
	if config.AddDisclaimer {
		reply = reply + "\n\n---\n" + config.DisclaimerText
	}

	// 记录日志
	kbMatchIDsStr := ""
	if len(kbMatchIDs) > 0 {
		parts := make([]string, len(kbMatchIDs))
		for i, id := range kbMatchIDs {
			parts[i] = fmt.Sprintf("%d", id)
		}
		kbMatchIDsStr = strings.Join(parts, ",")
	}

	log := model.AITicketLog{
		TicketID:   req.TicketID,
		Question:   prompt,
		Answer:     reply,
		Confidence: confidence,
		KBMatchIDs: kbMatchIDsStr,
		TokensUsed: tokensUsed,
	}
	s.db.Create(&log)

	return AutoReplyResult{
		ShouldReply: true,
		Reply:       reply,
		Confidence:  confidence,
		KBMatchIDs:  kbMatchIDs,
	}
}

// GetAutoReplyLogs 获取自动回复日志
func (s *AITicketService) GetAutoReplyLogs(ticketID uint, page, pageSize int) ([]model.AITicketLog, int64, error) {
	var logs []model.AITicketLog
	var total int64
	q := s.db.Model(&model.AITicketLog{})
	if ticketID > 0 {
		q = q.Where("ticket_id = ?", ticketID)
	}
	q.Count(&total)
	offset := (page - 1) * pageSize
	if err := q.Order("id DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

// MarkReplyAccepted 标记回复是否被接受
func (s *AITicketService) MarkReplyAccepted(logID uint, accepted bool) error {
	return s.db.Model(&model.AITicketLog{}).Where("id = ?", logID).Update("accepted", accepted).Error
}

// ─── AI 调用 ───

type openaiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openaiRequest struct {
	Model       string          `json:"model"`
	Messages    []openaiMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float64         `json:"temperature"`
}

type openaiResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
}

func (s *AITicketService) callAI(config *model.AIConfig, userPrompt string) (string, float64, int, error) {
	endpoint := config.APIEndpoint
	if endpoint == "" {
		endpoint = "https://api.openai.com/v1/chat/completions"
	}

	systemPrompt := config.SystemPrompt
	if systemPrompt == "" {
		systemPrompt = "你是一个专业的客服助手。请根据用户的问题给出准确、有帮助的回答。如果不确定，请说明并建议用户等待人工客服。"
	}

	reqBody := openaiRequest{
		Model: config.Model,
		Messages: []openaiMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		MaxTokens:   config.MaxTokens,
		Temperature: config.Temperature,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", 0, 0, err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", 0, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", 0, 0, fmt.Errorf("AI API 请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return "", 0, 0, fmt.Errorf("AI API 返回错误 %d: %s", resp.StatusCode, string(body))
	}

	var result openaiResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", 0, 0, fmt.Errorf("解析 AI 响应失败: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", 0, 0, fmt.Errorf("AI 未返回任何回复")
	}

	// 简单的置信度估算（基于回复长度和是否包含不确定性词汇）
	reply := result.Choices[0].Message.Content
	confidence := estimateConfidence(reply)

	return reply, confidence, result.Usage.TotalTokens, nil
}

// estimateConfidence 简单估算 AI 回复的置信度
func estimateConfidence(reply string) float64 {
	confidence := 0.85

	// 不确定性词汇降低置信度
	uncertainPhrases := []string{
		"不确定", "不清楚", "可能", "也许", "大概", "建议咨询",
		"请联系", "我无法", "抱歉", "不太了解", "需要确认",
		"not sure", "maybe", "perhaps", "I'm not certain",
		"sorry", "I cannot", "I'm unable",
	}

	lower := strings.ToLower(reply)
	for _, phrase := range uncertainPhrases {
		if strings.Contains(lower, strings.ToLower(phrase)) {
			confidence -= 0.15
			break
		}
	}

	// 回复太短也降低置信度
	if len([]rune(reply)) < 20 {
		confidence -= 0.2
	}

	if confidence < 0.3 {
		confidence = 0.3
	}
	return confidence
}
