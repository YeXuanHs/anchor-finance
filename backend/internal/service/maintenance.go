package service

import (
	"encoding/json"
	"errors"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/internal/util"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// MaintenanceService 维护模式服务
type MaintenanceService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewMaintenanceService(db *gorm.DB, log *logger.Logger) *MaintenanceService {
	return &MaintenanceService{db: db, log: log}
}

// MaintenanceStatus 维护状态
type MaintenanceStatus struct {
	Enabled   bool       `json:"enabled"`
	Message   string     `json:"message"`
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
}

// ScheduleRequest 维护计划请求
type ScheduleRequest struct {
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
	Message   string `json:"message"`
}

// getOrCreate 获取或创建维护配置
func (s *MaintenanceService) getOrCreate() (*model.ConfigOption, error) {
	var config model.ConfigOption
	err := s.db.Where("code = ?", "maintenance_mode").First(&config).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		config = model.ConfigOption{
			Code:  "maintenance_mode",
			Group: "system",
			Name:  "维护模式",
			Type:  "json",
			Value: `{"enabled":false,"message":"","start_time":null,"end_time":null}`,
		}
		if err := s.db.Create(&config).Error; err != nil {
			return nil, err
		}
		return &config, nil
	}
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// Enable 启用维护模式
func (s *MaintenanceService) Enable(message, allowedIPs, estimatedAt string) error {
	config, err := s.getOrCreate()
	if err != nil {
		return err
	}
	value := map[string]interface{}{
		"enabled":       true,
		"message":       message,
		"allowed_ips":   allowedIPs,
		"estimated_at":  estimatedAt,
		"start_time":    nil,
		"end_time":      nil,
	}
	jsonBytes, _ := json.Marshal(value)
	return s.db.Model(config).Update("value", string(jsonBytes)).Error
}

// Disable 禁用维护模式
func (s *MaintenanceService) Disable() error {
	config, err := s.getOrCreate()
	if err != nil {
		return err
	}
	return s.db.Model(config).Update("value", `{"enabled":false,"message":"","start_time":null,"end_time":null}`).Error
}

// GetStatus 获取维护状态
func (s *MaintenanceService) GetStatus() (*MaintenanceStatus, error) {
	config, err := s.getOrCreate()
	if err != nil {
		return nil, err
	}
	var cfg struct {
		Enabled   bool    `json:"enabled"`
		Message   string  `json:"message"`
		StartTime *string `json:"start_time"`
		EndTime   *string `json:"end_time"`
	}
	if err := json.Unmarshal([]byte(config.Value), &cfg); err != nil {
		return &MaintenanceStatus{Enabled: false}, nil
	}
	status := &MaintenanceStatus{
		Enabled: cfg.Enabled,
		Message: cfg.Message,
	}
	if cfg.StartTime != nil {
		t, _ := time.Parse(time.RFC3339, *cfg.StartTime)
		status.StartTime = &t
	}
	if cfg.EndTime != nil {
		t, _ := time.Parse(time.RFC3339, *cfg.EndTime)
		status.EndTime = &t
	}
	return status, nil
}

// SetMessage 设置维护消息
func (s *MaintenanceService) SetMessage(message string) error {
	config, err := s.getOrCreate()
	if err != nil {
		return err
	}
	return s.db.Model(config).Update("value", `{"enabled":true,"message":"`+message+`","start_time":null,"end_time":null}`).Error
}

// ScheduleMaintenance 计划维护
func (s *MaintenanceService) ScheduleMaintenance(req ScheduleRequest) error {
	startTime, err := util.ParseTime(req.StartTime)
	if err != nil {
		return errors.New("invalid start_time")
	}
	endTime, err := util.ParseTime(req.EndTime)
	if err != nil {
		return errors.New("invalid end_time")
	}
	if endTime.Before(startTime) {
		return errors.New("end_time must be after start_time")
	}

	config, err := s.getOrCreate()
	if err != nil {
		return err
	}

	message := req.Message
	if message == "" {
		message = "System maintenance in progress"
	}

	return s.db.Model(config).Update("value",
		`{"enabled":false,"message":"`+message+`","start_time":"`+startTime.Format(time.RFC3339)+`","end_time":"`+endTime.Format(time.RFC3339)+`"}`,
	).Error
}

// Update updates maintenance mode settings.
func (s *MaintenanceService) Update(updates map[string]interface{}) error {
	config, err := s.getOrCreate()
	if err != nil {
		return err
	}

	var cfg map[string]interface{}
	json.Unmarshal([]byte(config.Value), &cfg)

	for k, v := range updates {
		cfg[k] = v
	}

	jsonBytes, _ := json.Marshal(cfg)
	return s.db.Model(config).Update("value", string(jsonBytes)).Error
}

// GetHistory returns maintenance mode history.
func (s *MaintenanceService) GetHistory(page, pageSize int) ([]map[string]interface{}, int64, error) {
	var items []map[string]interface{}
	var total int64

	// For now, return empty history as we don't have a history table
	return items, total, nil
}

// AddAllowedIP adds an IP to the allowed list.
func (s *MaintenanceService) AddAllowedIP(ip string) error {
	config, err := s.getOrCreate()
	if err != nil {
		return err
	}

	var cfg map[string]interface{}
	json.Unmarshal([]byte(config.Value), &cfg)

	allowedIPs := ""
	if v, ok := cfg["allowed_ips"].(string); ok {
		allowedIPs = v
	}
	if allowedIPs != "" {
		allowedIPs += ","
	}
	allowedIPs += ip
	cfg["allowed_ips"] = allowedIPs

	jsonBytes, _ := json.Marshal(cfg)
	return s.db.Model(config).Update("value", string(jsonBytes)).Error
}

// RemoveAllowedIP removes an IP from the allowed list.
func (s *MaintenanceService) RemoveAllowedIP(ip string) error {
	config, err := s.getOrCreate()
	if err != nil {
		return err
	}

	var cfg map[string]interface{}
	json.Unmarshal([]byte(config.Value), &cfg)

	allowedIPs := ""
	if v, ok := cfg["allowed_ips"].(string); ok {
		allowedIPs = v
	}

	// Remove the IP from the list
	ips := make([]string, 0)
	for _, existingIP := range splitIPs(allowedIPs) {
		if existingIP != ip {
			ips = append(ips, existingIP)
		}
	}
	cfg["allowed_ips"] = joinIPs(ips)

	jsonBytes, _ := json.Marshal(cfg)
	return s.db.Model(config).Update("value", string(jsonBytes)).Error
}

// GetAllowedIPs returns the list of allowed IPs.
func (s *MaintenanceService) GetAllowedIPs() ([]string, error) {
	config, err := s.getOrCreate()
	if err != nil {
		return nil, err
	}

	var cfg map[string]interface{}
	json.Unmarshal([]byte(config.Value), &cfg)

	allowedIPs := ""
	if v, ok := cfg["allowed_ips"].(string); ok {
		allowedIPs = v
	}

	return splitIPs(allowedIPs), nil
}

// TestMode tests maintenance mode display.
func (s *MaintenanceService) TestMode(message string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"enabled": true,
		"message": message,
		"test":    true,
	}, nil
}

// splitIPs splits comma-separated IPs into a slice.
func splitIPs(ips string) []string {
	if ips == "" {
		return []string{}
	}
	result := make([]string, 0)
	for _, ip := range splitString(ips, ",") {
		ip = trimSpace(ip)
		if ip != "" {
			result = append(result, ip)
		}
	}
	return result
}

// joinIPs joins IPs with comma.
func joinIPs(ips []string) string {
	result := ""
	for i, ip := range ips {
		if i > 0 {
			result += ","
		}
		result += ip
	}
	return result
}

func splitString(s, sep string) []string {
	// Simple split implementation
	result := make([]string, 0)
	start := 0
	for i := 0; i < len(s); i++ {
		if string(s[i]) == sep {
			result = append(result, s[start:i])
			start = i + 1
		}
	}
	result = append(result, s[start:])
	return result
}

func trimSpace(s string) string {
	// Simple trim implementation
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}
