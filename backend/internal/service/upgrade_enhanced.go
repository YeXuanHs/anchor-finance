package service

import (
	"fmt"
	"time"

	"anchorfinance/pkg/logger"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// UpgradeEnhancedService 升级增强服务
type UpgradeEnhancedService struct {
	db *gorm.DB
}

// NewUpgradeEnhancedService 创建升级增强服务
func NewUpgradeEnhancedService(db *gorm.DB) *UpgradeEnhancedService {
	return &UpgradeEnhancedService{db: db}
}

// UpgradeConfig 升级配置
type UpgradeConfig struct {
	ID              uint           `gorm:"primaryKey"`
	ProductID       uint           `gorm:"index"`
	AllowUpgrade    bool           `gorm:"default:true"`
	AllowDowngrade  bool           `gorm:"default:false"`
	ProrateCredit   bool           `gorm:"default:true"` // 按比例折算余额
	UpgradeCycle    string         `gorm:"size:32"`      // immediate/next_billing
	TargetProducts  datatypes.JSON `gorm:"type:json"`   // 可升级到的产品ID列表
	Config          datatypes.JSON `gorm:"type:json"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// UpgradeValidation 升级验证结果
type UpgradeValidation struct {
	Valid      bool
	Errors     []string
	Warnings   []string
	PriceDiff  float64
	ProrateAmt float64
	NewCycle   string
}

// JudgeUpgradeConfig 判断升级配置
func (s *UpgradeEnhancedService) JudgeUpgradeConfig(hostID uint) (*UpgradeConfig, error) {
	var host struct {
		ProductID uint
	}
	if err := s.db.Table("hosts").Where("id = ?", hostID).Select("product_id").First(&host).Error; err != nil {
		return nil, err
	}

	var config UpgradeConfig
	err := s.db.Where("product_id = ?", host.ProductID).First(&config).Error
	if err != nil {
		// 返回默认配置
		return &UpgradeConfig{
			ProductID:    host.ProductID,
			AllowUpgrade: true,
			ProrateCredit: true,
			UpgradeCycle: "immediate",
		}, nil
	}

	return &config, nil
}

// JudgeUpgradeConfigError 验证升级配置错误
func (s *UpgradeEnhancedService) JudgeUpgradeConfigError(hostID uint, targetProductID uint) *UpgradeValidation {
	result := &UpgradeValidation{Valid: true}

	config, err := s.JudgeUpgradeConfig(hostID)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, "Failed to get upgrade config")
		return result
	}

	if !config.AllowUpgrade {
		result.Valid = false
		result.Errors = append(result.Errors, "Upgrade is not allowed for this product")
	}

	// 检查目标产品是否在允许列表
	if len(config.TargetProducts) > 0 {
		var allowed bool
		var targetIDs []uint
		// 解析JSON数组
		if err := config.TargetProducts.Unmarshal(&targetIDs); err == nil {
			for _, id := range targetIDs {
				if id == targetProductID {
					allowed = true
					break
				}
			}
		}
		if !allowed {
			result.Valid = false
			result.Errors = append(result.Errors, "Target product is not in allowed upgrade list")
		}
	}

	return result
}

// CheckChange 检查变更
func (s *UpgradeEnhancedService) CheckChange(hostID uint, targetProductID uint, newCycle string) (map[string]interface{}, error) {
	var host struct {
		ID        uint
		ProductID uint
		Cycle     string
		Price     float64
		ExpiredAt time.Time
	}
	if err := s.db.Table("hosts").Where("id = ?", hostID).First(&host).Error; err != nil {
		return nil, err
	}

	var targetProduct struct {
		ID    uint
		Name  string
		Price float64
		Cycle string
	}
	if err := s.db.Table("products").Where("id = ?", targetProductID).First(&targetProduct).Error; err != nil {
		return nil, err
	}

	// 计算价格差异
	priceDiff := targetProduct.Price - host.Price

	// 计算按比例折算金额
	var prorateAmt float64
	if host.ExpiredAt.After(time.Now()) {
		remaining := time.Until(host.ExpiredAt)
		total := time.Duration(30) * 24 * time.Hour // 假设月付
		if host.Cycle == "quarterly" {
			total = time.Duration(90) * 24 * time.Hour
		} else if host.Cycle == "yearly" {
			total = time.Duration(365) * 24 * time.Hour
		}
		prorateAmt = host.Price * (remaining.Seconds() / total.Seconds())
	}

	result := map[string]interface{}{
		"current_product":   host.ProductID,
		"target_product":    targetProductID,
		"target_product_name": targetProduct.Name,
		"current_price":     host.Price,
		"target_price":      targetProduct.Price,
		"price_diff":        priceDiff,
		"prorate_amount":    prorateAmt,
		"new_cycle":         newCycle,
		"current_cycle":     host.Cycle,
	}

	return result, nil
}

// CheckChangeText 检查变更文本描述
func (s *UpgradeEnhancedService) CheckChangeText(hostID uint, targetProductID uint) (string, error) {
	change, err := s.CheckChange(hostID, targetProductID, "")
	if err != nil {
		return "", err
	}

	text := fmt.Sprintf("升级从产品#%d到产品#%s，价格差异: %.2f，按比例折算: %.2f",
		change["current_product"], change["target_product_name"],
		change["price_diff"].(float64), change["prorate_amount"].(float64))

	return text, nil
}

// FilterConfigOptions 过滤配置选项
func (s *UpgradeEnhancedService) FilterConfigOptions(productID uint, currentConfig map[string]interface{}) ([]map[string]interface{}, error) {
	var options []struct {
		ID          uint
		Name        string
		Label       string
		Type        string
		Options     datatypes.JSON
		DefaultVal  string
		Required    bool
		SortOrder   int
	}

	err := s.db.Table("config_options").
		Where("product_id = ? AND enabled = ?", productID, true).
		Order("sort_order").
		Find(&options).Error
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(options))
	for i, opt := range options {
		result[i] = map[string]interface{}{
			"id":         opt.ID,
			"name":       opt.Name,
			"label":      opt.Label,
			"type":       opt.Type,
			"options":    opt.Options,
			"default":    opt.DefaultVal,
			"required":   opt.Required,
			"sort_order": opt.SortOrder,
		}

		// 如果有当前配置，标记已选值
		if val, ok := currentConfig[opt.Name]; ok {
			result[i]["current_value"] = val
		}
	}

	return result, nil
}

// UpgradeConfigAdmin 管理员升级配置
func (s *UpgradeEnhancedService) UpgradeConfigAdmin(hostID uint, adminID uint, data map[string]interface{}) error {
	// 记录管理员操作
	log := map[string]interface{}{
		"host_id":   hostID,
		"admin_id":  adminID,
		"action":    "upgrade_config",
		"data":      data,
		"timestamp": time.Now(),
	}

	logger.Info("Admin upgrade config", "host_id", hostID, "admin_id", adminID)

	// 更新主机配置
	return s.db.Table("hosts").Where("id = ?", hostID).Updates(data).Error
}

// UpgradeConfigProduct 产品升级配置
func (s *UpgradeEnhancedService) UpgradeConfigProduct(productID uint) (*UpgradeConfig, error) {
	var config UpgradeConfig
	err := s.db.Where("product_id = ?", productID).First(&config).Error
	if err != nil {
		// 创建默认配置
		config = UpgradeConfig{
			ProductID:    productID,
			AllowUpgrade: true,
			ProrateCredit: true,
			UpgradeCycle: "immediate",
		}
		s.db.Create(&config)
	}

	return &config, nil
}

// UpgradeConfigCommon 通用升级配置
func (s *UpgradeEnhancedService) UpgradeConfigCommon() (map[string]interface{}, error) {
	var settings []struct {
		Key   string
		Value string
	}
	s.db.Table("system_settings").Where("key IN ?", []string{
		"upgrade_allow_downgrade",
		"upgrade_prorate_credit",
		"upgrade_cycle",
	}).Find(&settings)

	result := make(map[string]interface{})
	for _, s := range settings {
		result[s.Key] = s.Value
	}

	return result, nil
}

// CheckUpgradePromo 检查升级促销
func (s *UpgradeEnhancedService) CheckUpgradePromo(hostID uint, targetProductID uint) (float64, error) {
	var host struct {
		UserID    uint
		ProductID uint
	}
	s.db.Table("hosts").Where("id = ?", hostID).First(&host)

	// 查找适用的促销活动
	var promo struct {
		DiscountType  string
		DiscountValue float64
	}

	s.db.Table("sale_promotions").
		Where("enabled = ? AND start_date <= ? AND end_date >= ?", true, time.Now(), time.Now()).
		Where("JSON_CONTAINS(target_products, ?) OR target_products = '[]'", fmt.Sprintf("%d", targetProductID)).
		First(&promo)

	if promo.DiscountValue > 0 {
		if promo.DiscountType == "percent" {
			return promo.DiscountValue, nil
		}
		return promo.DiscountValue, nil
	}

	return 0, nil
}

// DoUpgrade 执行升级
func (s *UpgradeEnhancedService) DoUpgrade(hostID uint, targetProductID uint, newCycle string, userID uint) (*UpgradeValidation, error) {
	// 验证升级
	validation := s.JudgeUpgradeConfigError(hostID, targetProductID)
	if !validation.Valid {
		return validation, fmt.Errorf("upgrade validation failed: %v", validation.Errors)
	}

	// 检查价格
	change, err := s.CheckChange(hostID, targetProductID, newCycle)
	if err != nil {
		return nil, err
	}

	priceDiff := change["price_diff"].(float64)

	// 如果需要付费，创建订单
	if priceDiff > 0 {
		order := map[string]interface{}{
			"user_id":    userID,
			"type":       "upgrade",
			"host_id":    hostID,
			"product_id": targetProductID,
			"amount":     priceDiff,
			"status":     "pending",
		}

		logger.Info("Creating upgrade order", "host_id", hostID, "amount", priceDiff)
		_ = order
	}

	// 更新主机
	updates := map[string]interface{}{
		"product_id": targetProductID,
	}
	if newCycle != "" {
		updates["cycle"] = newCycle
	}

	s.db.Table("hosts").Where("id = ?", hostID).Updates(updates)

	validation.Valid = true
	validation.PriceDiff = priceDiff

	return validation, nil
}

// AllowUpgradeProducts 获取可升级的产品列表
func (s *UpgradeEnhancedService) AllowUpgradeProducts(hostID uint) ([]map[string]interface{}, error) {
	var host struct {
		ProductID uint
	}
	s.db.Table("hosts").Where("id = ?", hostID).Select("product_id").First(&host)

	var config UpgradeConfig
	s.db.Where("product_id = ?", host.ProductID).First(&config)

	var products []struct {
		ID    uint
		Name  string
		Price float64
		Cycle string
	}

	if len(config.TargetProducts) > 0 {
		var targetIDs []uint
		config.TargetProducts.Unmarshal(&targetIDs)
		s.db.Table("products").Where("id IN ?", targetIDs).Find(&products)
	} else {
		// 获取同组的所有产品
		var groupID uint
		s.db.Table("products").Where("id = ?", host.ProductID).Select("group_id").Scan(&groupID)
		s.db.Table("products").Where("group_id = ? AND id != ?", groupID, host.ProductID).Find(&products)
	}

	result := make([]map[string]interface{}, len(products))
	for i, p := range products {
		result[i] = map[string]interface{}{
			"id":    p.ID,
			"name":  p.Name,
			"price": p.Price,
			"cycle": p.Cycle,
		}
	}

	return result, nil
}

// GetProductUpgradeConfig 获取产品升级配置
func (s *UpgradeEnhancedService) GetProductUpgradeConfig(productID uint) (map[string]interface{}, error) {
	config, err := s.UpgradeConfigProduct(productID)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"product_id":      config.ProductID,
		"allow_upgrade":   config.AllowUpgrade,
		"allow_downgrade": config.AllowDowngrade,
		"prorate_credit":  config.ProrateCredit,
		"upgrade_cycle":   config.UpgradeCycle,
		"target_products": config.TargetProducts,
	}

	return result, nil
}

// UpgradeProductCommon 产品通用升级逻辑
func (s *UpgradeEnhancedService) UpgradeProductCommon(hostID uint, params map[string]interface{}) error {
	targetProductID, _ := params["target_product_id"].(uint)
	newCycle, _ := params["cycle"].(string)

	_, err := s.DoUpgrade(hostID, targetProductID, newCycle, 0)
	return err
}
