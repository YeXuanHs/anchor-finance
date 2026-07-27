package service

import (
	"errors"
	"fmt"

	"anchorfinance/internal/model"
	"anchorfinance/internal/util"
	"anchorfinance/pkg/logger"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PaymentGatewayService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewPaymentGatewayService(db *gorm.DB, log *logger.Logger) *PaymentGatewayService {
	return &PaymentGatewayService{db: db, log: log}
}

type CreatePaymentGatewayRequest struct {
	Name             string  `json:"name" binding:"required,max=64"`
	Code             string  `json:"code" binding:"required,max=64"`
	Description      string  `json:"description"`
	Icon             string  `json:"icon" binding:"max=512"`
	Config           string  `json:"config"` // JSON string
	Currencies       string  `json:"currencies"`
	MinAmount        float64 `json:"min_amount"`
	MaxAmount        float64 `json:"max_amount"`
	FixedFee         float64 `json:"fixed_fee"`
	PercentFee       float64 `json:"percent_fee"`
	FeeCurrency      string  `json:"fee_currency" binding:"max=8"`
	IsOnline         bool    `json:"is_online"`
	SortOrder        int     `json:"sort_order"`
	Status           int16   `json:"status"`
	TestMode         bool    `json:"test_mode"`
	SupportRefund    bool    `json:"support_refund"`
	SupportRecurring bool    `json:"support_recurring"`
	WebhookURL       string  `json:"webhook_url" binding:"max=512"`
	ReturnURL        string  `json:"return_url" binding:"max=512"`
	CancelURL        string  `json:"cancel_url" binding:"max=512"`
	NotifyURL        string  `json:"notify_url" binding:"max=512"`
}

type UpdatePaymentGatewayRequest struct {
	Name             *string  `json:"name" binding:"omitempty,max=64"`
	Code             *string  `json:"code" binding:"omitempty,max=64"`
	Description      *string  `json:"description"`
	Icon             *string  `json:"icon" binding:"omitempty,max=512"`
	Config           *string  `json:"config"`
	Currencies       *string  `json:"currencies"`
	MinAmount        *float64 `json:"min_amount"`
	MaxAmount        *float64 `json:"max_amount"`
	FixedFee         *float64 `json:"fixed_fee"`
	PercentFee       *float64 `json:"percent_fee"`
	FeeCurrency      *string  `json:"fee_currency" binding:"omitempty,max=8"`
	IsOnline         *bool    `json:"is_online"`
	SortOrder        *int     `json:"sort_order"`
	Status           *int16   `json:"status"`
	TestMode         *bool    `json:"test_mode"`
	SupportRefund    *bool    `json:"support_refund"`
	SupportRecurring *bool    `json:"support_recurring"`
	WebhookURL       *string  `json:"webhook_url" binding:"omitempty,max=512"`
	ReturnURL        *string  `json:"return_url" binding:"omitempty,max=512"`
	CancelURL        *string  `json:"cancel_url" binding:"omitempty,max=512"`
	NotifyURL        *string  `json:"notify_url" binding:"omitempty,max=512"`
}

func (s *PaymentGatewayService) Create(req CreatePaymentGatewayRequest) (*model.PaymentGateway, error) {
	var count int64
	s.db.Model(&model.PaymentGateway{}).Where("code = ?", req.Code).Count(&count)
	if count > 0 {
		return nil, errors.New("payment gateway code already exists")
	}

	gw := model.PaymentGateway{
		Name:             req.Name,
		Code:             req.Code,
		Description:      req.Description,
		Icon:             req.Icon,
		FeeCurrency:      req.FeeCurrency,
		IsOnline:         req.IsOnline,
		SortOrder:        req.SortOrder,
		Status:           req.Status,
		TestMode:         req.TestMode,
		SupportRefund:    req.SupportRefund,
		SupportRecurring: req.SupportRecurring,
		WebhookURL:       req.WebhookURL,
		ReturnURL:        req.ReturnURL,
		CancelURL:        req.CancelURL,
		NotifyURL:        req.NotifyURL,
	}
	if gw.Status == 0 {
		gw.Status = 1
	}
	if req.Config != "" {
		gw.Config = datatypes.JSON(req.Config)
	}
	if req.Currencies != "" {
		gw.Currencies = datatypes.JSON(req.Currencies)
	}

	// Use raw SQL for decimal fields to avoid type issues
	if err := s.db.Create(&gw).Error; err != nil {
		return nil, err
	}

	// Update decimal fields via raw values
	if req.MinAmount != 0 || req.MaxAmount != 0 || req.FixedFee != 0 || req.PercentFee != 0 {
		updates := map[string]interface{}{}
		if req.MinAmount != 0 {
			updates["min_amount"] = fmt.Sprintf("%.4f", req.MinAmount)
		}
		if req.MaxAmount != 0 {
			updates["max_amount"] = fmt.Sprintf("%.4f", req.MaxAmount)
		}
		if req.FixedFee != 0 {
			updates["fixed_fee"] = fmt.Sprintf("%.4f", req.FixedFee)
		}
		if req.PercentFee != 0 {
			updates["percent_fee"] = fmt.Sprintf("%.4f", req.PercentFee)
		}
		if len(updates) > 0 {
			s.db.Model(&gw).Updates(updates)
		}
		s.db.First(&gw, gw.ID)
	}

	s.log.Infof("payment gateway created: id=%d code=%s", gw.ID, gw.Code)
	return &gw, nil
}

func (s *PaymentGatewayService) Update(id uint, req UpdatePaymentGatewayRequest) (*model.PaymentGateway, error) {
	var gw model.PaymentGateway
	if err := s.db.First(&gw, id).Error; err != nil {
		return nil, err
	}

	if req.Code != nil {
		var count int64
		s.db.Model(&model.PaymentGateway{}).Where("code = ? AND id != ?", *req.Code, id).Count(&count)
		if count > 0 {
			return nil, errors.New("payment gateway code already exists")
		}
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Code != nil {
		updates["code"] = *req.Code
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Icon != nil {
		updates["icon"] = *req.Icon
	}
	if req.Config != nil {
		updates["config"] = datatypes.JSON(*req.Config)
	}
	if req.Currencies != nil {
		updates["currencies"] = datatypes.JSON(*req.Currencies)
	}
	if req.MinAmount != nil {
		updates["min_amount"] = fmt.Sprintf("%.4f", *req.MinAmount)
	}
	if req.MaxAmount != nil {
		updates["max_amount"] = fmt.Sprintf("%.4f", *req.MaxAmount)
	}
	if req.FixedFee != nil {
		updates["fixed_fee"] = fmt.Sprintf("%.4f", *req.FixedFee)
	}
	if req.PercentFee != nil {
		updates["percent_fee"] = fmt.Sprintf("%.4f", *req.PercentFee)
	}
	if req.FeeCurrency != nil {
		updates["fee_currency"] = *req.FeeCurrency
	}
	if req.IsOnline != nil {
		updates["is_online"] = *req.IsOnline
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.TestMode != nil {
		updates["test_mode"] = *req.TestMode
	}
	if req.SupportRefund != nil {
		updates["support_refund"] = *req.SupportRefund
	}
	if req.SupportRecurring != nil {
		updates["support_recurring"] = *req.SupportRecurring
	}
	if req.WebhookURL != nil {
		updates["webhook_url"] = *req.WebhookURL
	}
	if req.ReturnURL != nil {
		updates["return_url"] = *req.ReturnURL
	}
	if req.CancelURL != nil {
		updates["cancel_url"] = *req.CancelURL
	}
	if req.NotifyURL != nil {
		updates["notify_url"] = *req.NotifyURL
	}

	if len(updates) > 0 {
		if err := s.db.Model(&gw).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	if err := s.db.First(&gw, id).Error; err != nil {
		return nil, err
	}
	s.log.Infof("payment gateway updated: id=%d", id)
	return &gw, nil
}

func (s *PaymentGatewayService) Delete(id uint) error {
	result := s.db.Delete(&model.PaymentGateway{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("payment gateway not found")
	}
	s.log.Infof("payment gateway deleted: id=%d", id)
	return nil
}

func (s *PaymentGatewayService) GetByID(id uint) (*model.PaymentGateway, error) {
	var gw model.PaymentGateway
	if err := s.db.First(&gw, id).Error; err != nil {
		return nil, err
	}
	return &gw, nil
}

func (s *PaymentGatewayService) GetList(page, pageSize int, status int) ([]model.PaymentGateway, int64, error) {
	var items []model.PaymentGateway
	var total int64

	query := s.db.Model(&model.PaymentGateway{})
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := util.Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).
		Order("sort_order ASC, id ASC").
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// GetEnabled returns enabled payment gateways for frontend.
func (s *PaymentGatewayService) GetEnabled() ([]model.PaymentGateway, error) {
	var items []model.PaymentGateway
	if err := s.db.Where("status = 1").
		Order("sort_order ASC, id ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (s *PaymentGatewayService) ToggleStatus(id uint) error {
	var gw model.PaymentGateway
	if err := s.db.First(&gw, id).Error; err != nil {
		return err
	}
	newStatus := int16(1)
	if gw.Status == 1 {
		newStatus = 0
	}
	return s.db.Model(&gw).Update("status", newStatus).Error
}

func (s *PaymentGatewayService) UpdateSort(id uint, sortOrder int) error {
	result := s.db.Model(&model.PaymentGateway{}).Where("id = ?", id).Update("sort_order", sortOrder)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("payment gateway not found")
	}
	return nil
}
