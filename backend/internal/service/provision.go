package service

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type ProvisionService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewProvisionService(db *gorm.DB, log *logger.Logger) *ProvisionService {
	return &ProvisionService{db: db, log: log}
}

// GetList returns all provision modules with optional type and active filter.
func (s *ProvisionService) GetList(page, pageSize int, moduleType string, active *bool) ([]model.ProvisionModule, int64, error) {
	var modules []model.ProvisionModule
	var total int64

	query := s.db.Model(&model.ProvisionModule{})
	if moduleType != "" {
		query = query.Where("type = ?", moduleType)
	}
	if active != nil {
		query = query.Where("active = ?", *active)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).
		Order("priority DESC, weight DESC, id ASC").
		Find(&modules).Error; err != nil {
		return nil, 0, err
	}
	return modules, total, nil
}

// GetByID returns a single provision module by ID.
func (s *ProvisionService) GetByID(id uint) (*model.ProvisionModule, error) {
	var module model.ProvisionModule
	if err := s.db.First(&module, id).Error; err != nil {
		return nil, err
	}
	return &module, nil
}

// Create creates a new provision module.
func (s *ProvisionService) Create(req CreateProvisionRequest) (*model.ProvisionModule, error) {
	var existing model.ProvisionModule
	if err := s.db.Where("slug = ?", req.Slug).First(&existing).Error; err == nil {
		return nil, errors.New("module with this slug already exists")
	}

	module := &model.ProvisionModule{
		Name:              req.Name,
		Slug:              req.Slug,
		Description:       req.Description,
		Type:              req.Type,
		SupportedProducts: req.SupportedProducts,
		Config:            req.Config,
		ServerURL:         req.ServerURL,
		ServerIP:          req.ServerIP,
		APIKey:            req.APIKey,
		APISecret:         req.APISecret,
		Username:          req.Username,
		Password:          req.Password,
		Hash:              req.Hash,
		Active:            true,
		Priority:          req.Priority,
		Weight:            req.Weight,
		MaxRetries:        req.MaxRetries,
		Timeout:           req.Timeout,
		Metadata:          req.Metadata,
	}
	if module.Weight == 0 {
		module.Weight = 1
	}
	if module.MaxRetries == 0 {
		module.MaxRetries = 3
	}
	if module.Timeout == 0 {
		module.Timeout = 30
	}
	if err := s.db.Create(module).Error; err != nil {
		return nil, err
	}
	return module, nil
}

// Update updates a provision module.
func (s *ProvisionService) Update(id uint, req UpdateProvisionRequest) (*model.ProvisionModule, error) {
	var module model.ProvisionModule
	if err := s.db.First(&module, id).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Slug != nil {
		updates["slug"] = *req.Slug
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.SupportedProducts != nil {
		updates["supported_products"] = *req.SupportedProducts
	}
	if req.Config != nil {
		updates["config"] = *req.Config
	}
	if req.ServerURL != nil {
		updates["server_url"] = *req.ServerURL
	}
	if req.ServerIP != nil {
		updates["server_ip"] = *req.ServerIP
	}
	if req.APIKey != nil {
		updates["api_key"] = *req.APIKey
	}
	if req.APISecret != nil {
		updates["api_secret"] = *req.APISecret
	}
	if req.Username != nil {
		updates["username"] = *req.Username
	}
	if req.Password != nil {
		updates["password"] = *req.Password
	}
	if req.Hash != nil {
		updates["hash"] = *req.Hash
	}
	if req.Active != nil {
		updates["active"] = *req.Active
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}
	if req.Weight != nil {
		updates["weight"] = *req.Weight
	}
	if req.MaxRetries != nil {
		updates["max_retries"] = *req.MaxRetries
	}
	if req.Timeout != nil {
		updates["timeout"] = *req.Timeout
	}

	if len(updates) > 0 {
		if err := s.db.Model(&module).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	if err := s.db.First(&module, id).Error; err != nil {
		return nil, err
	}
	return &module, nil
}

// Delete soft-deletes a provision module.
func (s *ProvisionService) Delete(id uint) error {
	return s.db.Delete(&model.ProvisionModule{}, id).Error
}

// TestConnection tests the connectivity to a provision module's server.
func (s *ProvisionService) TestConnection(id uint) (*model.ProvisionModule, error) {
	module, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	module.LastTestAt = &now

	testURL := module.ServerURL
	if testURL == "" {
		module.LastTestOK = false
		module.LastError = "server_url is empty"
		s.db.Model(module).Updates(map[string]interface{}{
			"last_test_at": now,
			"last_test_ok": false,
			"last_error":   "server_url is empty",
		})
		return module, errors.New("server_url is empty")
	}

	parsed, err := url.Parse(testURL)
	if err != nil {
		module.LastTestOK = false
		module.LastError = "invalid server_url: " + err.Error()
		s.db.Model(module).Updates(map[string]interface{}{
			"last_test_at": now,
			"last_test_ok": false,
			"last_error":   module.LastError,
		})
		return module, err
	}

	client := &http.Client{
		Timeout: time.Duration(module.Timeout) * time.Second,
	}

	var lastErr error
	for attempt := 0; attempt < module.MaxRetries; attempt++ {
		resp, err := client.Get(parsed.String())
		if err != nil {
			lastErr = err
			continue
		}
		resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 500 {
			module.LastTestOK = true
			module.LastError = ""
			s.db.Model(module).Updates(map[string]interface{}{
				"last_test_at": now,
				"last_test_ok": true,
				"last_error":   "",
			})
			return module, nil
		}
		lastErr = errors.New("server returned status " + resp.Status)
	}

	module.LastTestOK = false
	if lastErr != nil {
		module.LastError = lastErr.Error()
	} else {
		module.LastError = "connection failed after retries"
	}
	s.db.Model(module).Updates(map[string]interface{}{
		"last_test_at": now,
		"last_test_ok": false,
		"last_error":   module.LastError,
	})
	return module, lastErr
}

// GetLogs returns provision module operation logs with pagination.
func (s *ProvisionService) GetLogs(page, pageSize int, moduleID *uint, action string) ([]model.ProvisionLog, int64, error) {
	var logs []model.ProvisionLog
	var total int64

	query := s.db.Model(&model.ProvisionLog{})
	if moduleID != nil {
		query = query.Where("module_id = ?", *moduleID)
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).
		Order("id DESC").
		Preload("Module").
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

// IncrProvisionCount increments provision counters.
func (s *ProvisionService) IncrProvisionCount(moduleID uint, success bool) {
	updates := map[string]interface{}{
		"provision_count": gorm.Expr("provision_count + 1"),
	}
	if success {
		updates["success_count"] = gorm.Expr("success_count + 1")
	} else {
		updates["fail_count"] = gorm.Expr("fail_count + 1")
	}
	s.db.Model(&model.ProvisionModule{}).Where("id = ?", moduleID).Updates(updates)
}

// ---------- Request DTOs ----------

type CreateProvisionRequest struct {
	Name              string `json:"name" binding:"required,max=128"`
	Slug              string `json:"slug" binding:"required,max=128"`
	Description       string `json:"description"`
	Type              string `json:"type" binding:"required"`
	SupportedProducts string `json:"supported_products"`
	Config            string `json:"config"`
	ServerURL         string `json:"server_url"`
	ServerIP          string `json:"server_ip"`
	APIKey            string `json:"api_key"`
	APISecret         string `json:"api_secret"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	Hash              string `json:"hash"`
	Priority          int    `json:"priority"`
	Weight            int    `json:"weight"`
	MaxRetries        int    `json:"max_retries"`
	Timeout           int    `json:"timeout"`
	Metadata          string `json:"metadata"`
}

type UpdateProvisionRequest struct {
	Name              *string `json:"name"`
	Slug              *string `json:"slug"`
	Description       *string `json:"description"`
	Type              *string `json:"type"`
	SupportedProducts *string `json:"supported_products"`
	Config            *string `json:"config"`
	ServerURL         *string `json:"server_url"`
	ServerIP          *string `json:"server_ip"`
	APIKey            *string `json:"api_key"`
	APISecret         *string `json:"api_secret"`
	Username          *string `json:"username"`
	Password          *string `json:"password"`
	Hash              *string `json:"hash"`
	Active            *bool   `json:"active"`
	Priority          *int    `json:"priority"`
	Weight            *int    `json:"weight"`
	MaxRetries        *int    `json:"max_retries"`
	Timeout           *int    `json:"timeout"`
}
