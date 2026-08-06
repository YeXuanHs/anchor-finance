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
func (s *MaintenanceService) Enable(message string) error {
	config, err := s.getOrCreate()
	if err != nil {
		return err
	}
	return s.db.Model(config).Update("value", `{"enabled":true,"message":"`+message+`","start_time":null,"end_time":null}`).Error
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
