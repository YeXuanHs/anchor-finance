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

type TicketDeliverService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewTicketDeliverService(db *gorm.DB, log *logger.Logger) *TicketDeliverService {
	return &TicketDeliverService{db: db, log: log}
}

// GetAddPage returns data needed for adding a deliver rule.
func (s *TicketDeliverService) GetAddPage() (map[string]interface{}, error) {
	return map[string]interface{}{
		"departments": []interface{}{},
		"products":    []interface{}{},
	}, nil
}

// GetRules returns all deliver rules.
func (s *TicketDeliverService) GetRules(page, pageSize int, keyword string) ([]model.TicketDeliverRule, int64, error) {
	var rules []model.TicketDeliverRule
	var total int64

	query := s.db.Model(&model.TicketDeliverRule{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("priority DESC, id ASC").Find(&rules).Error; err != nil {
		return nil, 0, err
	}
	return rules, total, nil
}

// GetRuleByID returns a single rule by ID.
func (s *TicketDeliverService) GetRuleByID(id uint) (*model.TicketDeliverRule, error) {
	var rule model.TicketDeliverRule
	if err := s.db.First(&rule, id).Error; err != nil {
		return nil, err
	}
	return &rule, nil
}

// CreateRule creates a new deliver rule.
func (s *TicketDeliverService) CreateRule(rule *model.TicketDeliverRule) error {
	return s.db.Create(rule).Error
}

// UpdateRule updates a deliver rule.
func (s *TicketDeliverService) UpdateRule(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.TicketDeliverRule{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteRule deletes a deliver rule.
func (s *TicketDeliverService) DeleteRule(id uint) error {
	return s.db.Delete(&model.TicketDeliverRule{}, id).Error
}

// GetLogs returns deliver logs for a ticket.
func (s *TicketDeliverService) GetLogs(ticketID uint, page, pageSize int) ([]model.TicketDeliverLog, int64, error) {
	var logs []model.TicketDeliverLog
	var total int64

	query := s.db.Model(&model.TicketDeliverLog{})
	if ticketID > 0 {
		query = query.Where("ticket_id = ?", ticketID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

// ────────────────────────────────────────────────────────────
// 上游透传功能
// ────────────────────────────────────────────────────────────

// GetUpstreams returns all upstream configurations.
func (s *TicketDeliverService) GetUpstreams() ([]model.TicketUpstream, error) {
	var upstreams []model.TicketUpstream
	if err := s.db.Where("status = 1").Find(&upstreams).Error; err != nil {
		return nil, err
	}
	return upstreams, nil
}

// GetUpstreamByID returns a single upstream by ID.
func (s *TicketDeliverService) GetUpstreamByID(id uint) (*model.TicketUpstream, error) {
	var upstream model.TicketUpstream
	if err := s.db.First(&upstream, id).Error; err != nil {
		return nil, err
	}
	return &upstream, nil
}

// CreateUpstream creates a new upstream configuration.
func (s *TicketDeliverService) CreateUpstream(upstream *model.TicketUpstream) error {
	return s.db.Create(upstream).Error
}

// UpdateUpstream updates an upstream configuration.
func (s *TicketDeliverService) UpdateUpstream(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.TicketUpstream{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteUpstream deletes an upstream configuration.
func (s *TicketDeliverService) DeleteUpstream(id uint) error {
	return s.db.Delete(&model.TicketUpstream{}, id).Error
}

// TestUpstreamConnection tests connection to an upstream system.
func (s *TicketDeliverService) TestUpstreamConnection(id uint) (bool, string, error) {
	upstream, err := s.GetUpstreamByID(id)
	if err != nil {
		return false, "上游系统不存在", err
	}

	// 根据类型测试连接
	switch upstream.Type {
	case "anchorfinance":
		return s.testAnchorFinanceConnection(upstream)
	case "zjmf":
		return s.testZjmfConnection(upstream)
	case "v10":
		return s.testV10Connection(upstream)
	default:
		return false, "不支持的上游类型", nil
	}
}

// testAnchorFinanceConnection tests connection to an AnchorFinance instance.
func (s *TicketDeliverService) testAnchorFinanceConnection(upstream *model.TicketUpstream) (bool, string, error) {
	url := strings.TrimRight(upstream.URL, "/") + "/api/v1/health"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, "创建请求失败", err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, "连接失败: " + err.Error(), err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, "连接成功", nil
	}
	return false, fmt.Sprintf("连接失败: HTTP %d", resp.StatusCode), nil
}

// testZjmfConnection tests connection to a zjmf instance.
func (s *TicketDeliverService) testZjmfConnection(upstream *model.TicketUpstream) (bool, string, error) {
	url := strings.TrimRight(upstream.URL, "/") + "/api/v1/system/info"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, "创建请求失败", err
	}
	if upstream.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+upstream.APIKey)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, "连接失败: " + err.Error(), err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, "连接成功", nil
	}
	return false, fmt.Sprintf("连接失败: HTTP %d", resp.StatusCode), nil
}

// testV10Connection tests connection to a v10 instance.
func (s *TicketDeliverService) testV10Connection(upstream *model.TicketUpstream) (bool, string, error) {
	url := strings.TrimRight(upstream.URL, "/") + "/api/v1/system/info"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, "创建请求失败", err
	}
	if upstream.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+upstream.APIKey)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, "连接失败: " + err.Error(), err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, "连接成功", nil
	}
	return false, fmt.Sprintf("连接失败: HTTP %d", resp.StatusCode), nil
}

// ForwardTicketToUpstream forwards a ticket to an upstream system.
func (s *TicketDeliverService) ForwardTicketToUpstream(ticketID uint, upstreamID uint, deptID uint, productID uint, maskKeywords string) error {
	// 获取上游配置
	upstream, err := s.GetUpstreamByID(upstreamID)
	if err != nil {
		return fmt.Errorf("上游系统不存在: %v", err)
	}

	// 获取工单信息
	var ticket struct {
		ID       uint   `json:"id"`
		Title    string `json:"title"`
		Content  string `json:"content"`
		Priority string `json:"priority"`
		UID      uint   `json:"uid"`
	}
	if err := s.db.Table("tickets").Select("id, title, content, priority, user_id as uid").Where("id = ?", ticketID).First(&ticket).Error; err != nil {
		return fmt.Errorf("工单不存在: %v", err)
	}

	// 敏感关键词脱敏
	title := ticket.Title
	content := ticket.Content
	if maskKeywords != "" {
		keywords := strings.Split(maskKeywords, "\n")
		for _, keyword := range keywords {
			keyword = strings.TrimSpace(keyword)
			if keyword != "" {
				mask := strings.Repeat("*", len([]rune(keyword)))
				title = strings.ReplaceAll(title, keyword, mask)
				content = strings.ReplaceAll(content, keyword, mask)
			}
		}
	}

	// 根据上游类型转发
	switch upstream.Type {
	case "anchorfinance":
		return s.forwardToAnchorFinance(upstream, ticketID, title, content, deptID, productID)
	case "zjmf":
		return s.forwardToZjmf(upstream, ticketID, title, content, deptID, productID)
	case "v10":
		return s.forwardToV10(upstream, ticketID, title, content, deptID, productID)
	default:
		return fmt.Errorf("不支持的上游类型: %s", upstream.Type)
	}
}

// forwardToAnchorFinance forwards a ticket to another AnchorFinance instance.
func (s *TicketDeliverService) forwardToAnchorFinance(upstream *model.TicketUpstream, ticketID uint, title, content string, deptID, productID uint) error {
	url := strings.TrimRight(upstream.URL, "/") + "/api/v1/tickets/create"

	payload := map[string]interface{}{
		"title":       title,
		"content":     content,
		"department":  deptID,
		"product_id": productID,
		"is_api":      true,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if upstream.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+upstream.APIKey)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("上游返回错误: %s", string(body))
	}

	// 解析响应获取上游工单ID
	var result struct {
		TicketID string `json:"ticket_id"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	// 保存映射关系
	mapping := model.TicketUpstreamMapping{
		LocalTicketID:   ticketID,
		UpstreamTicketID: result.TicketID,
		UpstreamID:      upstream.ID,
		Direction:       "upstream",
	}
	return s.db.Create(&mapping).Error
}

// forwardToZjmf forwards a ticket to a zjmf instance.
func (s *TicketDeliverService) forwardToZjmf(upstream *model.TicketUpstream, ticketID uint, title, content string, deptID, productID uint) error {
	url := strings.TrimRight(upstream.URL, "/") + "/api/v1/ticket/create"

	payload := map[string]interface{}{
		"title":      title,
		"content":    content,
		"dptid":      deptID,
		"hostid":     productID,
		"is_api":     1,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if upstream.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+upstream.APIKey)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("上游返回错误: %s", string(body))
	}

	// 解析响应获取上游工单ID
	var result struct {
		Tid string `json:"tid"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	// 保存映射关系
	mapping := model.TicketUpstreamMapping{
		LocalTicketID:   ticketID,
		UpstreamTicketID: result.Tid,
		UpstreamID:      upstream.ID,
		Direction:       "upstream",
	}
	return s.db.Create(&mapping).Error
}

// forwardToV10 forwards a ticket to a v10 instance.
func (s *TicketDeliverService) forwardToV10(upstream *model.TicketUpstream, ticketID uint, title, content string, deptID, productID uint) error {
	url := strings.TrimRight(upstream.URL, "/") + "/api/v1/tickets/create"

	payload := map[string]interface{}{
		"title":      title,
		"content":    content,
		"department": deptID,
		"product_id": productID,
		"is_api":     true,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if upstream.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+upstream.APIKey)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("上游返回错误: %s", string(body))
	}

	// 解析响应获取上游工单ID
	var result struct {
		TicketID string `json:"ticket_id"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	// 保存映射关系
	mapping := model.TicketUpstreamMapping{
		LocalTicketID:   ticketID,
		UpstreamTicketID: result.TicketID,
		UpstreamID:      upstream.ID,
		Direction:       "upstream",
	}
	return s.db.Create(&mapping).Error
}

// GetTicketUpstreamMapping returns the upstream mapping for a ticket.
func (s *TicketDeliverService) GetTicketUpstreamMapping(ticketID uint) (*model.TicketUpstreamMapping, error) {
	var mapping model.TicketUpstreamMapping
	if err := s.db.Where("local_ticket_id = ?", ticketID).First(&mapping).Error; err != nil {
		return nil, err
	}
	return &mapping, nil
}

// SyncUpstreamReply syncs a reply from upstream back to the local ticket.
func (s *TicketDeliverService) SyncUpstreamReply(upstreamID uint, upstreamTicketID string, content string) error {
	// 查找本地工单映射
	var mapping model.TicketUpstreamMapping
	if err := s.db.Where("upstream_id = ? AND upstream_ticket_id = ?", upstreamID, upstreamTicketID).First(&mapping).Error; err != nil {
		return fmt.Errorf("未找到工单映射: %v", err)
	}

	// 创建本地回复
	reply := struct {
		TicketID  uint   `json:"ticket_id"`
		Content   string `json:"content"`
		IsAdmin   bool   `json:"is_admin"`
		Source    string `json:"source"`
		CreatedAt time.Time `json:"created_at"`
	}{
		TicketID:  mapping.LocalTicketID,
		Content:   content,
		IsAdmin:   true,
		Source:    "upstream",
		CreatedAt: time.Now(),
	}

	return s.db.Table("ticket_replies").Create(&reply).Error
}
