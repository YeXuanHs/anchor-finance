package service

import (
	"errors"
	"fmt"
	"time"

	"anchorfinance/internal/model"

	"gorm.io/gorm"
)

// ClientServiceService manages client service (product instance) operations.
type ClientServiceService struct {
	db *gorm.DB
}

// NewClientServiceService creates a new ClientServiceService.
func NewClientServiceService(db *gorm.DB) *ClientServiceService {
	return &ClientServiceService{db: db}
}

// OpenServiceRequest is the payload for opening a client service.
type OpenServiceRequest struct {
	UserID    uint   `json:"user_id" binding:"required"`
	ProductID uint   `json:"product_id" binding:"required"`
	Name      string `json:"name" binding:"required,max=256"`
	Remark    string `json:"remark"`
}

// UpdateServiceRequest is the payload for updating service metadata.
type UpdateServiceRequest struct {
	Name      *string `json:"name"`
	Remark    *string `json:"remark"`
	AutoRenew *bool   `json:"auto_renew"`
}

// Open creates a new service instance for a user.
func (s *ClientServiceService) Open(req OpenServiceRequest) (*model.ClientService, error) {
	now := time.Now()
	svc := &model.ClientService{
		UserID:    req.UserID,
		ProductID: req.ProductID,
		Name:      req.Name,
		Status:    model.ClientServiceActive,
		OpenedAt:  &now,
		AutoRenew: false,
		Remark:    req.Remark,
	}
	if err := s.db.Create(svc).Error; err != nil {
		return nil, err
	}
	return svc, nil
}

// GetByID fetches a client service by ID.
func (s *ClientServiceService) GetByID(id uint) (*model.ClientService, error) {
	var svc model.ClientService
	if err := s.db.Preload("User").Preload("Product").First(&svc, id).Error; err != nil {
		return nil, err
	}
	return &svc, nil
}

// GetList returns a filtered, paginated list of client services.
func (s *ClientServiceService) GetList(page, pageSize int, userID uint, status int16) ([]model.ClientService, int64, error) {
	var items []model.ClientService
	var total int64

	query := s.db.Model(&model.ClientService{})
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if status > 0 {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("User").Preload("Product").
		Offset(offset).Limit(pageSize).Order("id DESC").
		Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// Update modifies service metadata (name, remark, auto_renew).
func (s *ClientServiceService) Update(id uint, req UpdateServiceRequest) error {
	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Remark != nil {
		updates["remark"] = *req.Remark
	}
	if req.AutoRenew != nil {
		updates["auto_renew"] = *req.AutoRenew
	}
	if len(updates) == 0 {
		return nil
	}
	return s.db.Model(&model.ClientService{}).Where("id = ?", id).Updates(updates).Error
}

// Suspend pauses an active service.
func (s *ClientServiceService) Suspend(id uint, reason string) error {
	svc, err := s.GetByID(id)
	if err != nil {
		return err
	}
	if svc.Status != model.ClientServiceActive {
		return errors.New("only active services can be suspended")
	}
	remark := svc.Remark
	if reason != "" {
		remark = reason
	}
	return s.db.Model(&model.ClientService{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status": model.ClientServiceSuspended,
			"remark": remark,
		}).Error
}

// Terminate permanently terminates a service.
func (s *ClientServiceService) Terminate(id uint, reason string) error {
	svc, err := s.GetByID(id)
	if err != nil {
		return err
	}
	if svc.Status == model.ClientServiceTerminated {
		return errors.New("service is already terminated")
	}
	remark := svc.Remark
	if reason != "" {
		remark = reason
	}
	return s.db.Model(&model.ClientService{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status": model.ClientServiceTerminated,
			"remark": remark,
		}).Error
}

// Renew extends a service's expiry date.
func (s *ClientServiceService) Renew(id uint, months int) error {
	if months <= 0 {
		return errors.New("months must be positive")
	}
	svc, err := s.GetByID(id)
	if err != nil {
		return err
	}
	if svc.Status == model.ClientServiceTerminated {
		return errors.New("cannot renew a terminated service")
	}

	var base time.Time
	if svc.ExpiredAt != nil && svc.ExpiredAt.After(time.Now()) {
		base = *svc.ExpiredAt
	} else {
		base = time.Now()
	}
	newExpiry := base.AddDate(0, months, 0)

	return s.db.Model(&model.ClientService{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"expired_at": newExpiry,
			"status":     model.ClientServiceActive,
		}).Error
}

// Resume reactivates a suspended service.
func (s *ClientServiceService) Resume(id uint) error {
	svc, err := s.GetByID(id)
	if err != nil {
		return err
	}
	if svc.Status != model.ClientServiceSuspended {
		return errors.New("only suspended services can be resumed")
	}
	return s.db.Model(&model.ClientService{}).Where("id = ?", id).
		Update("status", model.ClientServiceActive).Error
}

// GetDB returns the database instance.
func (s *ClientServiceService) GetDB() *gorm.DB {
	return s.db
}

// ==================== P1-10: HostRenew ====================

// HostRenew 单台主机续费（创建续费发票）
func (s *ClientServiceService) HostRenew(hostID uint, cycle string) (map[string]interface{}, error) {
	var host struct {
		ID            uint
		UserID        uint
		ProductID     uint
		BillingCycle  string
		NextDueDate   *time.Time
		FirstPayment  float64
	}
	if err := s.db.Table("host").Where("id = ?", hostID).
		Select("id, userid as user_id, productid as product_id, billingcycle as billing_cycle, nextduedate as next_due_date, firstpaymentamount as first_payment").
		Scan(&host).Error; err != nil {
		return nil, errors.New("host not found")
	}

	// 获取产品价格
	var product struct {
		ID    uint
		Name  string
		Price float64
	}
	s.db.Table("products").Where("id = ?", host.ProductID).Select("id, name, price").Scan(&product)

	// 计算续费金额
	basePrice := product.Price
	months := 1
	switch cycle {
	case "monthly":
		months = 1
	case "quarterly":
		months = 3
	case "semi-annually":
		months = 6
	case "annually":
		months = 12
	case "biennially":
		months = 24
	case "triennially":
		months = 36
	}
	amount := basePrice * float64(months)

	// 创建续费发票
	invoice := map[string]interface{}{
		"user_id":    host.UserID,
		"type":       "renew",
		"amount":     amount,
		"status":     0,
		"host_id":    hostID,
		"remark":     fmt.Sprintf("主机续费 - %s (%s)", product.Name, cycle),
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}
	if err := s.db.Table("invoices").Create(&invoice).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"host_id":   hostID,
		"cycle":     cycle,
		"amount":    amount,
		"invoice":   invoice,
	}, nil
}

// ==================== P1-11: UpgradeConfig ====================

// UpgradeConfig 升级配置（创建升级发票）
func (s *ClientServiceService) UpgradeConfig(hostID, newProductID uint, cycle string) (map[string]interface{}, error) {
	var host struct {
		ID           uint
		UserID       uint
		ProductID    uint
		BillingCycle string
	}
	if err := s.db.Table("host").Where("id = ?", hostID).
		Select("id, userid as user_id, productid as product_id, billingcycle as billing_cycle").
		Scan(&host).Error; err != nil {
		return nil, errors.New("host not found")
	}

	// 获取当前产品价格
	var oldProduct struct {
		Price float64
		Name  string
	}
	s.db.Table("products").Where("id = ?", host.ProductID).Select("price, name").Scan(&oldProduct)

	// 获取新产品价格
	var newProduct struct {
		Price float64
		Name  string
	}
	if err := s.db.Table("products").Where("id = ?", newProductID).Select("price, name").Scan(&newProduct).Error; err != nil {
		return nil, errors.New("new product not found")
	}

	// 计算升级差价
	months := 1
	switch cycle {
	case "quarterly":
		months = 3
	case "semi-annually":
		months = 6
	case "annually":
		months = 12
	case "biennially":
		months = 24
	case "triennially":
		months = 36
	}
	diff := (newProduct.Price - oldProduct.Price) * float64(months)
	if diff < 0 {
		diff = 0
	}

	// 创建升级发票
	invoice := map[string]interface{}{
		"user_id":    host.UserID,
		"type":       "upgrade",
		"amount":     diff,
		"status":     0,
		"host_id":    hostID,
		"remark":     fmt.Sprintf("配置升级 %s -> %s", oldProduct.Name, newProduct.Name),
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}
	if err := s.db.Table("invoices").Create(&invoice).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"host_id":       hostID,
		"old_product":   oldProduct.Name,
		"new_product":   newProduct.Name,
		"upgrade_price": diff,
		"invoice":       invoice,
	}, nil
}
