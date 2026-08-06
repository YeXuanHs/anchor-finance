package service

import (
	"encoding/json"
	"fmt"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// CSChatService 客服聊天服务
// 移植自 anchor_cloud_finance_pro 的 AdminCsChatController + 业务逻辑
type CSChatService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewCSChatService(db *gorm.DB, log *logger.Logger) *CSChatService {
	return &CSChatService{db: db, log: log}
}

// ─── 配置管理 ───

// GetConfig 获取配置值
func (s *CSChatService) GetConfig(key string) string {
	var cfg model.CSChatConfig
	if err := s.db.Where("`key` = ?", key).First(&cfg).Error; err != nil {
		return ""
	}
	return cfg.Value
}

// SetConfig 设置配置值
func (s *CSChatService) SetConfig(key, value, remark string) error {
	var cfg model.CSChatConfig
	if err := s.db.Where("`key` = ?", key).First(&cfg).Error; err != nil {
		cfg = model.CSChatConfig{Key: key, Value: value, Remark: remark}
		return s.db.Create(&cfg).Error
	}
	cfg.Value = value
	if remark != "" {
		cfg.Remark = remark
	}
	return s.db.Save(&cfg).Error
}

// GetAIConfig 获取 AI 配置
func (s *CSChatService) GetAIConfig() map[string]interface{} {
	raw := s.GetConfig("ai_config")
	if raw == "" {
		return map[string]interface{}{
			"api_endpoint": "",
			"api_key":      "",
			"model":        "",
		}
	}
	var result map[string]interface{}
	json.Unmarshal([]byte(raw), &result)
	return result
}

// SaveAIConfig 保存 AI 配置
func (s *CSChatService) SaveAIConfig(endpoint, key, model_name string) error {
	data := map[string]interface{}{
		"api_endpoint": endpoint,
		"model":        model_name,
	}
	if key != "" {
		data["api_key"] = key
	} else {
		// 保留旧密钥
		old := s.GetAIConfig()
		if v, ok := old["api_key"]; ok {
			data["api_key"] = v
		}
	}
	jsonData, _ := json.Marshal(data)
	return s.SetConfig("ai_config", string(jsonData), "AI配置")
}

// GetAppearanceConfig 获取外观配置
func (s *CSChatService) GetAppearanceConfig() map[string]interface{} {
	raw := s.GetConfig("appearance")
	if raw == "" {
		return map[string]interface{}{
			"avatar_url":  "",
			"theme_color": "#1890ff",
		}
	}
	var result map[string]interface{}
	json.Unmarshal([]byte(raw), &result)
	return result
}

// SaveAppearanceConfig 保存外观配置
func (s *CSChatService) SaveAppearanceConfig(data map[string]interface{}) error {
	jsonData, _ := json.Marshal(data)
	return s.SetConfig("appearance", string(jsonData), "外观配置")
}

// GetWorkingHoursConfig 获取工作时间配置
func (s *CSChatService) GetWorkingHoursConfig() map[string]interface{} {
	raw := s.GetConfig("working_hours")
	if raw == "" {
		return map[string]interface{}{
			"weekdays": map[string]interface{}{
				"1": map[string]interface{}{"enabled": true, "start": "09:00", "end": "18:00"},
				"2": map[string]interface{}{"enabled": true, "start": "09:00", "end": "18:00"},
				"3": map[string]interface{}{"enabled": true, "start": "09:00", "end": "18:00"},
				"4": map[string]interface{}{"enabled": true, "start": "09:00", "end": "18:00"},
				"5": map[string]interface{}{"enabled": true, "start": "09:00", "end": "18:00"},
				"6": map[string]interface{}{"enabled": false, "start": "09:00", "end": "18:00"},
				"7": map[string]interface{}{"enabled": false, "start": "09:00", "end": "18:00"},
			},
		}
	}
	var result map[string]interface{}
	json.Unmarshal([]byte(raw), &result)
	return result
}

// SaveWorkingHoursConfig 保存工作时间配置
func (s *CSChatService) SaveWorkingHoursConfig(data map[string]interface{}) error {
	jsonData, _ := json.Marshal(data)
	return s.SetConfig("working_hours", string(jsonData), "工作时间配置")
}

// IsWithinWorkingHours 检查当前是否在工作时间内
func (s *CSChatService) IsWithinWorkingHours() bool {
	cfg := s.GetWorkingHoursConfig()
	weekdays, ok := cfg["weekdays"].(map[string]interface{})
	if !ok {
		return false
	}
	now := time.Now()
	weekday := fmt.Sprintf("%d", int(now.Weekday()))
	if weekday == "0" {
		weekday = "7"
	}
	dayCfg, ok := weekdays[weekday].(map[string]interface{})
	if !ok {
		return false
	}
	enabled, _ := dayCfg["enabled"].(bool)
	if !enabled {
		return false
	}
	start, _ := dayCfg["start"].(string)
	end, _ := dayCfg["end"].(string)
	if start == "" || end == "" {
		return true
	}
	current := now.Format("15:04")
	return current >= start && current <= end
}

// ─── 会话管理 ───

type CSChatSessionItem struct {
	ID              uint       `json:"id"`
	ClientName      string     `json:"client_name"`
	ClientEmail     string     `json:"client_email"`
	ClientPhone     string     `json:"client_phone"`
	Mode            string     `json:"mode"`
	Status          string     `json:"status"`
	Rating          int        `json:"rating"`
	LastMessage     string     `json:"last_message"`
	ThreadCount     int        `json:"thread_count"`
	CreatedAt       time.Time  `json:"created_at"`
}

// ListSessions 获取会话列表
func (s *CSChatService) ListSessions(status string, page, pageSize int) ([]CSChatSessionItem, int64, error) {
	var sessions []model.CSChatSession
	var total int64
	q := s.db.Model(&model.CSChatSession{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	q.Count(&total)
	offset := (page - 1) * pageSize
	if err := q.Order("id DESC").Offset(offset).Limit(pageSize).Find(&sessions).Error; err != nil {
		return nil, 0, err
	}
	var items []CSChatSessionItem
	for _, sess := range sessions {
		items = append(items, CSChatSessionItem{
			ID:          sess.ID,
			ClientName:  sess.ClientName,
			ClientEmail: sess.ClientEmail,
			ClientPhone: sess.ClientPhone,
			Mode:        sess.Mode,
			Status:      sess.Status,
			Rating:      sess.Rating,
			LastMessage: sess.LastMessage,
			CreatedAt:   sess.CreatedAt,
		})
	}
	return items, total, nil
}

// GetSession 获取单个会话
func (s *CSChatService) GetSession(id uint) (*model.CSChatSession, error) {
	var session model.CSChatSession
	if err := s.db.First(&session, id).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

// GetSessionMessages 获取会话消息
func (s *CSChatService) GetSessionMessages(sessionID uint) ([]model.CSChatMessage, error) {
	var messages []model.CSChatMessage
	if err := s.db.Where("session_id = ?", sessionID).Order("id ASC").Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

// SendReply 发送回复（人工客服）
func (s *CSChatService) SendReply(sessionID uint, content string) error {
	msg := model.CSChatMessage{
		SessionID: sessionID,
		Role:      "human",
		Content:   content,
	}
	if err := s.db.Create(&msg).Error; err != nil {
		return err
	}
	// 更新会话最后消息
	s.db.Model(&model.CSChatSession{}).Where("id = ?", sessionID).Updates(map[string]interface{}{
		"last_message": content,
		"updated_at":   gorm.Expr("NOW()"),
	})
	return nil
}

// TransferToHuman 转人工
func (s *CSChatService) TransferToHuman(sessionID uint, staffID uint) error {
	return s.db.Model(&model.CSChatSession{}).Where("id = ?", sessionID).Updates(map[string]interface{}{
		"mode":      "human",
		"staff_id":  staffID,
		"updated_at": gorm.Expr("NOW()"),
	}).Error
}

// CloseSession 关闭会话
func (s *CSChatService) CloseSession(sessionID uint) error {
	return s.db.Model(&model.CSChatSession{}).Where("id = ?", sessionID).Updates(map[string]interface{}{
		"status":     "closed",
		"closed_at":  gorm.Expr("NOW()"),
		"updated_at": gorm.Expr("NOW()"),
	}).Error
}

// RateSession 评价会话
func (s *CSChatService) RateSession(sessionID uint, rating int, comment string) error {
	return s.db.Model(&model.CSChatSession{}).Where("id = ?", sessionID).Updates(map[string]interface{}{
		"rating":         rating,
		"rating_comment": comment,
	}).Error
}

// GetStats 获取统计
func (s *CSChatService) GetStats() map[string]interface{} {
	var totalOpen, totalClosed, totalToday int64
	s.db.Model(&model.CSChatSession{}).Where("status = ?", "open").Count(&totalOpen)
	s.db.Model(&model.CSChatSession{}).Where("status = ?", "closed").Count(&totalClosed)
	s.db.Model(&model.CSChatSession{}).Where("DATE(created_at) = CURDATE()").Count(&totalToday)

	var avgRating float64
	s.db.Model(&model.CSChatSession{}).Where("rating > 0").Select("AVG(rating)").Scan(&avgRating)

	return map[string]interface{}{
		"open":        totalOpen,
		"closed":      totalClosed,
		"today":       totalToday,
		"avg_rating":  avgRating,
	}
}

// GetOrCreateSession gets or creates a chat session for a user.
func (s *CSChatService) GetOrCreateSession(userID uint) (*model.CSChatSession, error) {
	var session model.CSChatSession
	err := s.db.Where("user_id = ? AND status = ?", userID, "open").First(&session).Error
	if err == nil {
		return &session, nil
	}

	session = model.CSChatSession{
		UserID: userID,
		Status: "open",
	}
	if err := s.db.Create(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

// SendMessage sends a message in a chat session.
func (s *CSChatService) SendMessage(sessionID, senderID uint, senderType, content string) (*model.CSChatMessage, error) {
	msg := &model.CSChatMessage{
		SessionID: sessionID,
		Role:      senderType,
		Content:   content,
	}
	if err := s.db.Create(msg).Error; err != nil {
		return nil, err
	}
	return msg, nil
}

// GetOrCreateSessionByVisitorID gets or creates a session by visitor ID (string).
func (s *CSChatService) GetOrCreateSessionByVisitorID(visitorID string) (*model.CSChatSession, error) {
	var session model.CSChatSession
	err := s.db.Where("visitor_id = ? AND status = ?", visitorID, "open").First(&session).Error
	if err == nil {
		return &session, nil
	}

	session = model.CSChatSession{
		VisitorID: visitorID,
		Status:    "open",
	}
	if err := s.db.Create(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

// AIReply generates an AI reply for a message.
func (s *CSChatService) AIReply(sessionID uint, userMessage string) (string, error) {
	// Simple auto-reply implementation
	reply := "您好，感谢您的咨询！客服正在处理中，请稍候。"

	msg := &model.CSChatMessage{
		SessionID: sessionID,
		Role:      "assistant",
		Content:   reply,
	}
	if err := s.db.Create(msg).Error; err != nil {
		return "", err
	}
	return reply, nil
}
