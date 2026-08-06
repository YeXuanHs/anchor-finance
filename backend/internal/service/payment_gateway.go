package service

import (
	"errors"
	"fmt"

	"anchorfinance/internal/model"
	"anchorfinance/internal/util"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/payment"

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
	Name      string  `json:"name" binding:"required,max=64"`
	Title     string  `json:"title" binding:"required,max=64"`
	Gateway   string  `json:"gateway" binding:"required,max=32"`  // epay, xunhupay, alipay, wxpay, balance, bank
	Code      string  `json:"code" binding:"required,max=32"`     // alipay, wechat, qqpay, usdt, bank
	Config    string  `json:"config"`                              // JSON配置
	FeeRate   float64 `json:"fee_rate"`
	MinAmount float64 `json:"min_amount"`
	MaxAmount float64 `json:"max_amount"`
	SortOrder int     `json:"sort_order"`
	IsEnabled bool    `json:"is_enabled"`
}

type UpdatePaymentGatewayRequest struct {
	Title     *string  `json:"title" binding:"omitempty,max=64"`
	Gateway   *string  `json:"gateway" binding:"omitempty,max=32"`
	Code      *string  `json:"code" binding:"omitempty,max=32"`
	Config    *string  `json:"config"`
	FeeRate   *float64 `json:"fee_rate"`
	MinAmount *float64 `json:"min_amount"`
	MaxAmount *float64 `json:"max_amount"`
	SortOrder *int     `json:"sort_order"`
	IsEnabled *bool    `json:"is_enabled"`
}

func (s *PaymentGatewayService) Create(req CreatePaymentGatewayRequest) (*model.PaymentGateway, error) {
	// 检查name唯一性
	var count int64
	s.db.Model(&model.PaymentGateway{}).Where("name = ?", req.Name).Count(&count)
	if count > 0 {
		return nil, errors.New("支付方式标识已存在")
	}

	// 验证接口类型
	if _, ok := payment.GatewayLabels[req.Gateway]; !ok {
		return nil, fmt.Errorf("不支持的支付接口: %s", req.Gateway)
	}

	gw := model.PaymentGateway{
		Name:      req.Name,
		Title:     req.Title,
		Gateway:   req.Gateway,
		Code:      req.Code,
		Config:    req.Config,
		FeeRate:   req.FeeRate,
		MinAmount: req.MinAmount,
		MaxAmount: req.MaxAmount,
		SortOrder: req.SortOrder,
		IsEnabled: req.IsEnabled,
	}

	if err := s.db.Create(&gw).Error; err != nil {
		return nil, err
	}

	s.log.Infof("支付方式创建成功: id=%d name=%s title=%s", gw.ID, gw.Name, gw.Title)
	return &gw, nil
}

func (s *PaymentGatewayService) Update(id uint, req UpdatePaymentGatewayRequest) (*model.PaymentGateway, error) {
	var gw model.PaymentGateway
	if err := s.db.First(&gw, id).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Gateway != nil {
		if _, ok := payment.GatewayLabels[*req.Gateway]; !ok {
			return nil, fmt.Errorf("不支持的支付接口: %s", *req.Gateway)
		}
		updates["gateway"] = *req.Gateway
	}
	if req.Code != nil {
		updates["code"] = *req.Code
	}
	if req.Config != nil {
		updates["config"] = *req.Config
	}
	if req.FeeRate != nil {
		updates["fee_rate"] = *req.FeeRate
	}
	if req.MinAmount != nil {
		updates["min_amount"] = *req.MinAmount
	}
	if req.MaxAmount != nil {
		updates["max_amount"] = *req.MaxAmount
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.IsEnabled != nil {
		updates["is_enabled"] = *req.IsEnabled
	}

	if len(updates) > 0 {
		if err := s.db.Model(&gw).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	if err := s.db.First(&gw, id).Error; err != nil {
		return nil, err
	}
	s.log.Infof("支付方式更新成功: id=%d", id)
	return &gw, nil
}

func (s *PaymentGatewayService) Delete(id uint) error {
	result := s.db.Delete(&model.PaymentGateway{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("支付方式不存在")
	}
	s.log.Infof("支付方式删除成功: id=%d", id)
	return nil
}

func (s *PaymentGatewayService) GetByID(id uint) (*model.PaymentGateway, error) {
	var gw model.PaymentGateway
	if err := s.db.First(&gw, id).Error; err != nil {
		return nil, err
	}
	return &gw, nil
}

func (s *PaymentGatewayService) GetList(page, pageSize int, isEnabled *bool) ([]model.PaymentGateway, int64, error) {
	var items []model.PaymentGateway
	var total int64

	query := s.db.Model(&model.PaymentGateway{})
	if isEnabled != nil {
		query = query.Where("is_enabled = ?", *isEnabled)
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

// GetEnabled 返回启用的支付方式（用户前台）
func (s *PaymentGatewayService) GetEnabled() ([]model.PaymentGateway, error) {
	var items []model.PaymentGateway
	if err := s.db.Where("is_enabled = ?", true).
		Order("sort_order ASC, id ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// ToggleStatus 切换启用状态
func (s *PaymentGatewayService) ToggleStatus(id uint) error {
	var gw model.PaymentGateway
	if err := s.db.First(&gw, id).Error; err != nil {
		return err
	}
	return s.db.Model(&gw).Update("is_enabled", !gw.IsEnabled).Error
}

// TestConnection 测试支付接口配置
func (s *PaymentGatewayService) TestConnection(id uint) error {
	var gw model.PaymentGateway
	if err := s.db.First(&gw, id).Error; err != nil {
		return errors.New("支付方式不存在")
	}

	if gw.Config == "" {
		return errors.New("支付接口未配置")
	}

	// 尝试创建支付接口实例
	_, err := payment.Factory(gw.Gateway, gw.Config)
	if err != nil {
		return fmt.Errorf("支付接口配置错误: %w", err)
	}

	s.log.Infof("支付接口测试通过: id=%d gateway=%s", gw.ID, gw.Gateway)
	return nil
}

// UpdateSort 更新排序
func (s *PaymentGatewayService) UpdateSort(id uint, sortOrder int) error {
	result := s.db.Model(&model.PaymentGateway{}).Where("id = ?", id).Update("sort_order", sortOrder)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("支付方式不存在")
	}
	return nil
}

// GetSupportedInfo 返回支持的接口和支付类型信息
func (s *PaymentGatewayService) GetSupportedInfo() map[string]interface{} {
	return map[string]interface{}{
		"gateways": payment.GatewayLabels,
		"codes":    payment.CodeLabels,
	}
}

// InitDefaults 初始化默认支付方式（首次启动）
func (s *PaymentGatewayService) InitDefaults() {
	var count int64
	s.db.Model(&model.PaymentGateway{}).Count(&count)
	if count > 0 {
		return
	}

	defaults := []model.PaymentGateway{
		{Name: "EpayAlipay", Title: "支付宝-易支付", Gateway: "epay", Code: "alipay", IsEnabled: false, SortOrder: 1},
		{Name: "EpayWechat", Title: "微信支付-易支付", Gateway: "epay", Code: "wechat", IsEnabled: false, SortOrder: 2},
		{Name: "EpayQQPay", Title: "QQ支付-易支付", Gateway: "epay", Code: "qqpay", IsEnabled: false, SortOrder: 3},
		{Name: "EpayUsdt", Title: "USDT-易支付", Gateway: "epay", Code: "usdt", IsEnabled: false, SortOrder: 4},
		{Name: "EpayBank", Title: "银联-易支付", Gateway: "epay", Code: "bank", IsEnabled: false, SortOrder: 5},
		{Name: "HpjAlipay", Title: "支付宝-虎皮椒", Gateway: "xunhupay", Code: "alipay", IsEnabled: false, SortOrder: 6},
		{Name: "HpjWechat", Title: "微信支付-虎皮椒", Gateway: "xunhupay", Code: "wechat", IsEnabled: false, SortOrder: 7},
		{Name: "AliPay", Title: "支付宝-官方", Gateway: "alipay", Code: "alipay", IsEnabled: false, SortOrder: 8},
		{Name: "WxPay", Title: "微信支付-官方", Gateway: "wxpay", Code: "wechat", IsEnabled: false, SortOrder: 9},
		{Name: "Balance", Title: "余额支付", Gateway: "balance", Code: "balance", IsEnabled: true, SortOrder: 10},
	}

	for _, gw := range defaults {
		s.db.Create(&gw)
	}

	s.log.Info("默认支付方式初始化完成")
}
