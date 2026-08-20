package service

import (
	"errors"
	"math"

	"anchorfinance/internal/model"
	"anchorfinance/internal/util"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// CurrencyService 货币管理服务
type CurrencyService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewCurrencyService(db *gorm.DB, log *logger.Logger) *CurrencyService {
	return &CurrencyService{db: db, log: log}
}

// CreateCurrencyRequest 创建货币请求
type CreateCurrencyRequest struct {
	Code      string  `json:"code" binding:"required,max=10"`
	Name      string  `json:"name" binding:"required,max=50"`
	Symbol    string  `json:"symbol" binding:"required,max=10"`
	Rate      float64 `json:"rate"`
	IsActive  bool    `json:"is_active"`
	Precision int     `json:"precision"`
}

// Create 创建货币
func (s *CurrencyService) Create(req CreateCurrencyRequest) (*model.Currency, error) {
	var count int64
	s.db.Model(&model.Currency{}).Where("code = ?", req.Code).Count(&count)
	if count > 0 {
		return nil, errors.New("currency code already exists")
	}

	currency := model.Currency{
		Code:      req.Code,
		Name:      req.Name,
		Symbol:    req.Symbol,
		Rate:      req.Rate,
		IsActive:  req.IsActive,
		Precision: req.Precision,
	}
	if currency.Rate == 0 {
		currency.Rate = 1
	}
	if currency.Precision == 0 {
		currency.Precision = 2
	}

	if err := s.db.Create(&currency).Error; err != nil {
		return nil, err
	}
	s.log.Infof("currency created: code=%s", currency.Code)
	return &currency, nil
}

// Get 获取货币
func (s *CurrencyService) Get(id uint) (*model.Currency, error) {
	var currency model.Currency
	if err := s.db.First(&currency, id).Error; err != nil {
		return nil, err
	}
	return &currency, nil
}

// GetByCode 根据代码获取货币
func (s *CurrencyService) GetByCode(code string) (*model.Currency, error) {
	var currency model.Currency
	if err := s.db.Where("code = ?", code).First(&currency).Error; err != nil {
		return nil, err
	}
	return &currency, nil
}

// List 获取货币列表
func (s *CurrencyService) List(page, pageSize int, activeOnly bool) ([]model.Currency, int64, error) {
	var items []model.Currency
	var total int64

	query := s.db.Model(&model.Currency{})
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	query.Count(&total)
	offset, limit := util.Paginate(page, pageSize)
	query.Offset(offset).Limit(limit).Order("is_default DESC, code ASC").Find(&items)
	return items, total, nil
}

// Update 更新货币
func (s *CurrencyService) Update(id uint, updates map[string]interface{}) (*model.Currency, error) {
	var currency model.Currency
	if err := s.db.First(&currency, id).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&currency).Updates(updates).Error; err != nil {
		return nil, err
	}
	s.db.First(&currency, id)
	return &currency, nil
}

// Delete 删除货币
func (s *CurrencyService) Delete(id uint) error {
	var currency model.Currency
	if err := s.db.First(&currency, id).Error; err != nil {
		return err
	}
	if currency.IsDefault {
		return errors.New("cannot delete default currency")
	}
	return s.db.Delete(&model.Currency{}, id).Error
}

// SetDefault 设置默认货币
func (s *CurrencyService) SetDefault(id uint) error {
	var currency model.Currency
	if err := s.db.First(&currency, id).Error; err != nil {
		return err
	}
	if !currency.IsActive {
		return errors.New("cannot set inactive currency as default")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Currency{}).Where("is_default = ?", true).Update("is_default", false).Error; err != nil {
			return err
		}
		return tx.Model(&currency).Update("is_default", true).Error
	})
}

// Convert 货币转换
func (s *CurrencyService) Convert(amount float64, fromCode, toCode string) (float64, error) {
	if fromCode == toCode {
		return amount, nil
	}

	var from, to model.Currency
	if err := s.db.Where("code = ? AND is_active = ?", fromCode, true).First(&from).Error; err != nil {
		return 0, errors.New("source currency not found or inactive")
	}
	if err := s.db.Where("code = ? AND is_active = ?", toCode, true).First(&to).Error; err != nil {
		return 0, errors.New("target currency not found or inactive")
	}

	converted := amount * to.Rate / from.Rate
	precision := math.Pow(10, float64(to.Precision))
	return math.Round(converted*precision) / precision, nil
}
