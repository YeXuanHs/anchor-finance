package service

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type HookService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewHookService(db *gorm.DB, log *logger.Logger) *HookService {
	return &HookService{db: db, log: log}
}

// GetList returns paginated hooks.
func (s *HookService) GetList(page, pageSize int, event string, status *int16) ([]model.Hook, int64, error) {
	var hooks []model.Hook
	var total int64

	query := s.db.Model(&model.Hook{})
	if event != "" {
		query = query.Where("event = ?", event)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&hooks).Error; err != nil {
		return nil, 0, err
	}
	return hooks, total, nil
}

// GetByID returns a single hook by ID.
func (s *HookService) GetByID(id uint) (*model.Hook, error) {
	var hook model.Hook
	if err := s.db.First(&hook, id).Error; err != nil {
		return nil, err
	}
	return &hook, nil
}

// GetByEvent returns all active hooks for an event.
func (s *HookService) GetByEvent(event string) ([]model.Hook, error) {
	var hooks []model.Hook
	if err := s.db.Where("event = ? AND status = 1", event).Order("id ASC").Find(&hooks).Error; err != nil {
		return nil, err
	}
	return hooks, nil
}

// Create creates a new hook.
func (s *HookService) Create(hook *model.Hook) error {
	return s.db.Create(hook).Error
}

// Update updates a hook.
func (s *HookService) Update(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.Hook{}).Where("id = ?", id).Updates(updates).Error
}

// Delete deletes a hook.
func (s *HookService) Delete(id uint) error {
	var hook model.Hook
	if err := s.db.First(&hook, id).Error; err != nil {
		return err
	}
	if hook.IsSystem {
		return gorm.ErrInvalidData
	}
	return s.db.Delete(&hook).Error
}

// SetStatus sets the status of a hook.
func (s *HookService) SetStatus(id uint, status int16) error {
	return s.db.Model(&model.Hook{}).Where("id = ?", id).Update("status", status).Error
}

// Trigger triggers hooks for an event.
func (s *HookService) Trigger(event string, data interface{}) error {
	hooks, err := s.GetByEvent(event)
	if err != nil {
		return err
	}

	for _, hook := range hooks {
		go s.executeHook(hook, event, data)
	}
	return nil
}

// executeHook 执行单个钩子
func (s *HookService) executeHook(hook model.Hook, event string, data interface{}) {
	timeout := time.Duration(hook.Timeout) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	startTime := time.Now()

	payload := map[string]interface{}{
		"event":     event,
		"data":      data,
		"timestamp": startTime.Unix(),
		"hook_id":   hook.ID,
		"hook_name": hook.Name,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		s.recordHookLog(hook, event, nil, "", 0, 0, err)
		s.updateHookStatus(hook, false)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", hook.URL, bytes.NewReader(bodyBytes))
	if err != nil {
		s.recordHookLog(hook, event, payload, "", 0, 0, err)
		s.updateHookStatus(hook, false)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Hook-Event", event)

	if len(hook.Headers) > 0 {
		var headers map[string]string
		if err := json.Unmarshal(hook.Headers, &headers); err == nil {
			for k, v := range headers {
				req.Header.Set(k, v)
			}
		}
	}

	if len(hook.Params) > 0 {
		var params map[string]interface{}
		if err := json.Unmarshal(hook.Params, &params); err == nil {
			if secret, ok := params["secret"].(string); ok && secret != "" {
				sig := computeHMAC(bodyBytes, secret)
				req.Header.Set("X-Hook-Signature", sig)
			}
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			s.recordHookLog(hook, event, payload, "", 0, 0, fmt.Errorf("timed out after %s", timeout))
		} else {
			s.recordHookLog(hook, event, payload, "", 0, 0, err)
		}
		s.updateHookStatus(hook, false)
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	duration := int(time.Since(startTime).Milliseconds())

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		s.recordHookLog(hook, event, payload, string(respBody), resp.StatusCode, duration, nil)
		s.updateHookStatus(hook, true)
	} else {
		s.recordHookLog(hook, event, payload, string(respBody), resp.StatusCode, duration, fmt.Errorf("HTTP %d", resp.StatusCode))
		s.updateHookStatus(hook, false)
	}
}

// computeHMAC 计算HMAC-SHA256签名
func computeHMAC(data []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(data)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

// recordHookLog 记录钩子执行日志
func (s *HookService) recordHookLog(hook model.Hook, event string, request interface{}, response string, statusCode int, duration int, err error) {
	status := int8(1)
	errMsg := ""
	if err != nil {
		status = 2
		errMsg = err.Error()
	}

	var requestJSON []byte
	if request != nil {
		requestJSON, _ = json.Marshal(request)
	}

	log := &model.HookLog{
		HookID:     hook.ID,
		Event:      event,
		Request:    requestJSON,
		Response:   response,
		StatusCode: statusCode,
		Status:     status,
		ErrorMsg:   errMsg,
		Duration:   duration,
	}
	if err := s.db.Create(log).Error; err != nil {
		s.log.Warnf("failed to create hook log: %v", err)
	}
}

// updateHookStatus 更新钩子运行状态
func (s *HookService) updateHookStatus(hook model.Hook, success bool) {
	now := time.Now()
	updates := map[string]interface{}{
		"last_run_at": &now,
		"run_count":   gorm.Expr("run_count + 1"),
	}
	if !success {
		updates["fail_count"] = gorm.Expr("fail_count + 1")
	}
	s.db.Model(&hook).Updates(updates)
}

// GetLogs returns paginated hook execution logs.
func (s *HookService) GetLogs(hookID uint, page, pageSize int) ([]model.HookLog, int64, error) {
	var logs []model.HookLog
	var total int64

	query := s.db.Model(&model.HookLog{})
	if hookID > 0 {
		query = query.Where("hook_id = ?", hookID)
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
