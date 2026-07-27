package service

import (
	"errors"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type ConfigServerService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewConfigServerService(db *gorm.DB, log *logger.Logger) *ConfigServerService {
	return &ConfigServerService{db: db, log: log}
}

// ---------- Server Config ----------

type CreateServerConfigRequest struct {
	Name            string  `json:"name" binding:"required,max=128"`
	Code            string  `json:"code" binding:"required,max=64"`
	Type            string  `json:"type" binding:"required,oneof=vps dedicated cloud reseller"`
	Provider        string  `json:"provider" binding:"max=128"`
	TemplateID      *uint   `json:"template_id"`
	CPU             string  `json:"cpu" binding:"max=64"`
	Memory          int     `json:"memory"`
	Disk            int     `json:"disk"`
	Bandwidth       int     `json:"bandwidth"`
	TrafficLimit    int64   `json:"traffic_limit"`
	IPCount         int     `json:"ip_count"`
	Location        string  `json:"location" binding:"max=128"`
	Datacenter      string  `json:"datacenter" binding:"max=128"`
	OS              string  `json:"os"`
	Features        string  `json:"features"`
	PriceMonthly    float64 `json:"price_monthly"`
	PriceQuarter    float64 `json:"price_quarter"`
	PriceSemiAnn    float64 `json:"price_semi_annual"`
	PriceAnnual     float64 `json:"price_annual"`
	PriceBiennial   float64 `json:"price_biennial"`
	PriceTriennial  float64 `json:"price_triennial"`
	PricingStrategy string  `json:"pricing_strategy" binding:"omitempty,oneof=fixed graduated promotional"`
	StockTotal      int     `json:"stock_total"`
	MaxPerUser      int     `json:"max_per_user"`
	SortOrder       int     `json:"sort_order"`
	Status          int16   `json:"status"`
	Remark          string  `json:"remark"`
}

type UpdateServerConfigRequest struct {
	Name            *string  `json:"name"`
	Code            *string  `json:"code"`
	Type            *string  `json:"type"`
	Provider        *string  `json:"provider"`
	TemplateID      *uint    `json:"template_id"`
	CPU             *string  `json:"cpu"`
	Memory          *int     `json:"memory"`
	Disk            *int     `json:"disk"`
	Bandwidth       *int     `json:"bandwidth"`
	TrafficLimit    *int64   `json:"traffic_limit"`
	IPCount         *int     `json:"ip_count"`
	Location        *string  `json:"location"`
	Datacenter      *string  `json:"datacenter"`
	OS              *string  `json:"os"`
	Features        *string  `json:"features"`
	PriceMonthly    *float64 `json:"price_monthly"`
	PriceQuarter    *float64 `json:"price_quarter"`
	PriceSemiAnn    *float64 `json:"price_semi_annual"`
	PriceAnnual     *float64 `json:"price_annual"`
	PriceBiennial   *float64 `json:"price_biennial"`
	PriceTriennial  *float64 `json:"price_triennial"`
	PricingStrategy *string  `json:"pricing_strategy"`
	StockTotal      *int     `json:"stock_total"`
	MaxPerUser      *int     `json:"max_per_user"`
	SortOrder       *int     `json:"sort_order"`
	Status          *int16   `json:"status"`
	Remark          *string  `json:"remark"`
}

func (s *ConfigServerService) GetList(page, pageSize int, serverType string, keyword string) ([]model.ServerConfig, int64, error) {
	var items []model.ServerConfig
	var total int64

	query := s.db.Model(&model.ServerConfig{})
	if serverType != "" {
		query = query.Where("type = ?", serverType)
	}
	if keyword != "" {
		q := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR code LIKE ?", q, q)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).
		Order("sort_order ASC, id ASC").
		Preload("Template").
		Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (s *ConfigServerService) GetByID(id uint) (*model.ServerConfig, error) {
	var item model.ServerConfig
	if err := s.db.Preload("Template").First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *ConfigServerService) Create(req CreateServerConfigRequest) (*model.ServerConfig, error) {
	var count int64
	s.db.Model(&model.ServerConfig{}).Where("code = ?", req.Code).Count(&count)
	if count > 0 {
		return nil, errors.New("server config code already exists")
	}

	item := &model.ServerConfig{
		Name:            req.Name,
		Code:            req.Code,
		Type:            req.Type,
		Provider:        req.Provider,
		TemplateID:      req.TemplateID,
		CPU:             req.CPU,
		Memory:          req.Memory,
		Disk:            req.Disk,
		Bandwidth:       req.Bandwidth,
		TrafficLimit:    req.TrafficLimit,
		IPCount:         req.IPCount,
		Location:        req.Location,
		Datacenter:      req.Datacenter,
		PriceMonthly:    req.PriceMonthly,
		PriceQuarter:    req.PriceQuarter,
		PriceSemiAnn:    req.PriceSemiAnn,
		PriceAnnual:     req.PriceAnnual,
		PriceBiennial:   req.PriceBiennial,
		PriceTriennial:  req.PriceTriennial,
		PricingStrategy: req.PricingStrategy,
		StockTotal:      req.StockTotal,
		StockUsed:       0,
		MaxPerUser:      req.MaxPerUser,
		SortOrder:       req.SortOrder,
		Status:          req.Status,
		Remark:          req.Remark,
	}
	if item.PricingStrategy == "" {
		item.PricingStrategy = "fixed"
	}
	if item.IPCount == 0 {
		item.IPCount = 1
	}

	if err := s.db.Create(item).Error; err != nil {
		return nil, err
	}
	s.log.Infof("server config created: id=%d code=%s", item.ID, item.Code)
	return item, nil
}

func (s *ConfigServerService) Update(id uint, req UpdateServerConfigRequest) (*model.ServerConfig, error) {
	var item model.ServerConfig
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}

	if req.Code != nil {
		var count int64
		s.db.Model(&model.ServerConfig{}).Where("code = ? AND id != ?", *req.Code, id).Count(&count)
		if count > 0 {
			return nil, errors.New("server config code already exists")
		}
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Code != nil {
		updates["code"] = *req.Code
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Provider != nil {
		updates["provider"] = *req.Provider
	}
	if req.TemplateID != nil {
		updates["template_id"] = req.TemplateID
	}
	if req.CPU != nil {
		updates["cpu"] = *req.CPU
	}
	if req.Memory != nil {
		updates["memory"] = *req.Memory
	}
	if req.Disk != nil {
		updates["disk"] = *req.Disk
	}
	if req.Bandwidth != nil {
		updates["bandwidth"] = *req.Bandwidth
	}
	if req.TrafficLimit != nil {
		updates["traffic_limit"] = *req.TrafficLimit
	}
	if req.IPCount != nil {
		updates["ip_count"] = *req.IPCount
	}
	if req.Location != nil {
		updates["location"] = *req.Location
	}
	if req.Datacenter != nil {
		updates["datacenter"] = *req.Datacenter
	}
	if req.OS != nil {
		updates["os"] = *req.OS
	}
	if req.Features != nil {
		updates["features"] = *req.Features
	}
	if req.PriceMonthly != nil {
		updates["price_monthly"] = *req.PriceMonthly
	}
	if req.PriceQuarter != nil {
		updates["price_quarter"] = *req.PriceQuarter
	}
	if req.PriceSemiAnn != nil {
		updates["price_semi_annual"] = *req.PriceSemiAnn
	}
	if req.PriceAnnual != nil {
		updates["price_annual"] = *req.PriceAnnual
	}
	if req.PriceBiennial != nil {
		updates["price_biennial"] = *req.PriceBiennial
	}
	if req.PriceTriennial != nil {
		updates["price_triennial"] = *req.PriceTriennial
	}
	if req.PricingStrategy != nil {
		updates["pricing_strategy"] = *req.PricingStrategy
	}
	if req.StockTotal != nil {
		updates["stock_total"] = *req.StockTotal
	}
	if req.MaxPerUser != nil {
		updates["max_per_user"] = *req.MaxPerUser
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Remark != nil {
		updates["remark"] = *req.Remark
	}

	if len(updates) > 0 {
		if err := s.db.Model(&item).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	if err := s.db.Preload("Template").First(&item, id).Error; err != nil {
		return nil, err
	}
	s.log.Infof("server config updated: id=%d", id)
	return &item, nil
}

func (s *ConfigServerService) Delete(id uint) error {
	var count int64
	s.db.Model(&model.ServerProduct{}).Where("server_config_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("server config has associated products, remove them first")
	}
	if err := s.db.Delete(&model.ServerConfig{}, id).Error; err != nil {
		return err
	}
	s.log.Infof("server config deleted: id=%d", id)
	return nil
}

func (s *ConfigServerService) BatchUpdateStatus(ids []uint, status int16) error {
	if len(ids) == 0 {
		return errors.New("ids is empty")
	}
	result := s.db.Model(&model.ServerConfig{}).Where("id IN ?", ids).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	s.log.Infof("server config batch status updated: ids=%v status=%d", ids, status)
	return nil
}

func (s *ConfigServerService) BatchDelete(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("ids is empty")
	}
	var count int64
	s.db.Model(&model.ServerProduct{}).Where("server_config_id IN ?", ids).Count(&count)
	if count > 0 {
		return errors.New("some server configs have associated products")
	}
	result := s.db.Delete(&model.ServerConfig{}, ids)
	if result.Error != nil {
		return result.Error
	}
	s.log.Infof("server config batch deleted: ids=%v", ids)
	return nil
}

func (s *ConfigServerService) UpdateSort(id uint, sortOrder int) error {
	result := s.db.Model(&model.ServerConfig{}).Where("id = ?", id).Update("sort_order", sortOrder)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("server config not found")
	}
	return nil
}

// ---------- Server Template ----------

type CreateServerTemplateRequest struct {
	Name        string `json:"name" binding:"required,max=128"`
	Code        string `json:"code" binding:"required,max=64"`
	Type        string `json:"type" binding:"required"`
	Description string `json:"description"`
	Config      string `json:"config" binding:"required"`
	SortOrder   int    `json:"sort_order"`
	Status      int16  `json:"status"`
}

type UpdateServerTemplateRequest struct {
	Name        *string `json:"name"`
	Code        *string `json:"code"`
	Type        *string `json:"type"`
	Description *string `json:"description"`
	Config      *string `json:"config"`
	SortOrder   *int    `json:"sort_order"`
	Status      *int16  `json:"status"`
}

func (s *ConfigServerService) GetTemplateList(page, pageSize int, templateType string) ([]model.ServerTemplate, int64, error) {
	var items []model.ServerTemplate
	var total int64

	query := s.db.Model(&model.ServerTemplate{})
	if templateType != "" {
		query = query.Where("type = ?", templateType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("sort_order ASC, id ASC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (s *ConfigServerService) CreateTemplate(req CreateServerTemplateRequest) (*model.ServerTemplate, error) {
	var count int64
	s.db.Model(&model.ServerTemplate{}).Where("code = ?", req.Code).Count(&count)
	if count > 0 {
		return nil, errors.New("template code already exists")
	}

	item := &model.ServerTemplate{
		Name:        req.Name,
		Code:        req.Code,
		Type:        req.Type,
		Description: req.Description,
		Config:      []byte(req.Config),
		SortOrder:   req.SortOrder,
		Status:      req.Status,
	}
	if item.Status == 0 {
		item.Status = 1
	}

	if err := s.db.Create(item).Error; err != nil {
		return nil, err
	}
	s.log.Infof("server template created: id=%d code=%s", item.ID, item.Code)
	return item, nil
}

func (s *ConfigServerService) UpdateTemplate(id uint, req UpdateServerTemplateRequest) (*model.ServerTemplate, error) {
	var item model.ServerTemplate
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}

	if req.Code != nil {
		var count int64
		s.db.Model(&model.ServerTemplate{}).Where("code = ? AND id != ?", *req.Code, id).Count(&count)
		if count > 0 {
			return nil, errors.New("template code already exists")
		}
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Code != nil {
		updates["code"] = *req.Code
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Config != nil {
		updates["config"] = *req.Config
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) > 0 {
		if err := s.db.Model(&item).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	s.log.Infof("server template updated: id=%d", id)
	return &item, nil
}

func (s *ConfigServerService) DeleteTemplate(id uint) error {
	var count int64
	s.db.Model(&model.ServerConfig{}).Where("template_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("template is in use by server configs")
	}
	if err := s.db.Delete(&model.ServerTemplate{}, id).Error; err != nil {
		return err
	}
	s.log.Infof("server template deleted: id=%d", id)
	return nil
}
