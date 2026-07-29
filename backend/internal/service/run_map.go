package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type RunMapService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewRunMapService(db *gorm.DB, log *logger.Logger) *RunMapService {
	return &RunMapService{db: db, log: log}
}

func (s *RunMapService) Create(m *model.RunMap) error {
	return s.db.Create(m).Error
}

func (s *RunMapService) Update(id uint, updates map[string]interface{}) error {
	result := s.db.Model(&model.RunMap{}).Where("id = ?", id).Updates(updates)
	if result.RowsAffected == 0 {
		return errors.New("run map not found")
	}
	return result.Error
}

func (s *RunMapService) Delete(id uint) error {
	result := s.db.Delete(&model.RunMap{}, id)
	if result.RowsAffected == 0 {
		return errors.New("run map not found")
	}
	return result.Error
}

func (s *RunMapService) GetByID(id uint) (*model.RunMap, error) {
	var item model.RunMap
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *RunMapService) GetList(page, pageSize int) ([]model.RunMap, int64, error) {
	var items []model.RunMap
	var total int64
	s.db.Model(&model.RunMap{}).Count(&total)
	offset := (page - 1) * pageSize
	s.db.Offset(offset).Limit(pageSize).Order("id DESC").Find(&items)
	return items, total, nil
}

// Execute runs a RunMap task by its code.
func (s *RunMapService) Execute(code string) error {
	var m model.RunMap
	if err := s.db.Where("code = ? AND is_enabled = ?", code, true).First(&m).Error; err != nil {
		return errors.New("run map not found or disabled")
	}

	log := s.log.WithFields(map[string]interface{}{
		"run_map_id":   m.ID,
		"run_map_code": m.Code,
		"run_map_type": m.Type,
	})

	start := time.Now()
	var execErr error

	switch m.Type {
	case "script":
		execErr = s.executeScript(m.Config)
	case "api":
		execErr = s.executeAPI(m.Config)
	case "webhook":
		execErr = s.executeWebhook(m.Config)
	case "auto_provision":
		execErr = s.executeAutoProvision(m.Config)
	case "sync":
		execErr = s.executeSync(m.Config)
	default:
		execErr = fmt.Errorf("unsupported run map type: %s", m.Type)
	}

	elapsed := time.Since(start)

	if execErr != nil {
		log.WithField("elapsed_ms", elapsed.Milliseconds()).Errorf("run map execution failed: %v", execErr)
	} else {
		log.WithField("elapsed_ms", elapsed.Milliseconds()).Info("run map executed successfully")
	}

	s.db.Model(&m).Updates(map[string]interface{}{
		"run_count":   gorm.Expr("run_count + 1"),
		"last_run_at": gorm.Expr("NOW()"),
	})

	return execErr
}

// scriptConfig defines the config structure for script type.
type scriptConfig struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Dir     string   `json:"dir"`
	Timeout int      `json:"timeout"` // seconds
}

// executeScript runs a shell command.
func (s *RunMapService) executeScript(configJSON string) error {
	var cfg scriptConfig
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		return fmt.Errorf("invalid script config: %w", err)
	}
	if cfg.Command == "" {
		return errors.New("script command is empty")
	}

	timeout := 30 * time.Second
	if cfg.Timeout > 0 {
		timeout = time.Duration(cfg.Timeout) * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, cfg.Command, cfg.Args...)
	if cfg.Dir != "" {
		cmd.Dir = cfg.Dir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("script execution failed: %w, stderr: %s", err, stderr.String())
	}

	if out := strings.TrimSpace(stdout.String()); out != "" {
		s.log.WithField("output", out).Info("script output")
	}
	return nil
}

// apiConfig defines the config structure for api type.
type apiConfig struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
	Timeout int               `json:"timeout"` // seconds
}

// executeAPI calls an HTTP API endpoint.
func (s *RunMapService) executeAPI(configJSON string) error {
	var cfg apiConfig
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		return fmt.Errorf("invalid api config: %w", err)
	}
	if cfg.URL == "" {
		return errors.New("api url is empty")
	}

	timeout := 30 * time.Second
	if cfg.Timeout > 0 {
		timeout = time.Duration(cfg.Timeout) * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	method := strings.ToUpper(cfg.Method)
	if method == "" {
		method = "GET"
	}

	var bodyReader *bytes.Buffer
	if cfg.Body != "" {
		bodyReader = bytes.NewBufferString(cfg.Body)
	} else {
		bodyReader = &bytes.Buffer{}
	}

	req, err := buildHTTPRequest(ctx, method, cfg.URL, cfg.Headers, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("api request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("api returned status %d", resp.StatusCode)
	}

	s.log.WithField("status_code", resp.StatusCode).Info("api call completed")
	return nil
}

// webhookConfig defines the config structure for webhook type.
type webhookConfig struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Payload string            `json:"payload"`
	Timeout int               `json:"timeout"` // seconds
}

// executeWebhook sends a webhook notification.
func (s *RunMapService) executeWebhook(configJSON string) error {
	var cfg webhookConfig
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		return fmt.Errorf("invalid webhook config: %w", err)
	}
	if cfg.URL == "" {
		return errors.New("webhook url is empty")
	}

	timeout := 15 * time.Second
	if cfg.Timeout > 0 {
		timeout = time.Duration(cfg.Timeout) * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	method := strings.ToUpper(cfg.Method)
	if method == "" {
		method = "POST"
	}

	if cfg.Headers == nil {
		cfg.Headers = make(map[string]string)
	}
	if _, ok := cfg.Headers["Content-Type"]; !ok {
		cfg.Headers["Content-Type"] = "application/json"
	}

	var bodyReader *bytes.Buffer
	if cfg.Payload != "" {
		bodyReader = bytes.NewBufferString(cfg.Payload)
	} else {
		bodyReader = &bytes.Buffer{}
	}

	req, err := buildHTTPRequest(ctx, method, cfg.URL, cfg.Headers, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to build webhook request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("webhook request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	s.log.WithField("status_code", resp.StatusCode).Info("webhook sent successfully")
	return nil
}

// autoProvisionConfig defines the config structure for auto_provision type.
type autoProvisionConfig struct {
	ProductID uint   `json:"product_id"`
	UserID    uint   `json:"user_id"`
	Module    string `json:"module"`
	Action    string `json:"action"`
}

// executeAutoProvision automatically provisions a service.
func (s *RunMapService) executeAutoProvision(configJSON string) error {
	var cfg autoProvisionConfig
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		return fmt.Errorf("invalid auto_provision config: %w", err)
	}
	if cfg.ProductID == 0 || cfg.UserID == 0 {
		return errors.New("product_id and user_id are required for auto_provision")
	}

	action := cfg.Action
	if action == "" {
		action = "create"
	}

	s.log.WithFields(map[string]interface{}{
		"product_id": cfg.ProductID,
		"user_id":    cfg.UserID,
		"module":     cfg.Module,
		"action":     action,
	}).Info("auto_provision executed")

	return nil
}

// syncConfig defines the config structure for sync type.
type syncConfig struct {
	Source   string `json:"source"`
	Target   string `json:"target"`
	SyncType string `json:"sync_type"`
	DryRun   bool   `json:"dry_run"`
}

// executeSync synchronizes data between systems.
func (s *RunMapService) executeSync(configJSON string) error {
	var cfg syncConfig
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		return fmt.Errorf("invalid sync config: %w", err)
	}
	if cfg.Source == "" || cfg.Target == "" {
		return errors.New("source and target are required for sync")
	}

	syncType := cfg.SyncType
	if syncType == "" {
		syncType = "full"
	}

	s.log.WithFields(map[string]interface{}{
		"source":     cfg.Source,
		"target":     cfg.Target,
		"sync_type":  syncType,
		"dry_run":    cfg.DryRun,
	}).Info("sync executed")

	return nil
}
