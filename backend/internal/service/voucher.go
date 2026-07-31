package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"anchorfinance/internal/util"
	"anchorfinance/pkg/logger"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Voucher 代金券
type Voucher struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Name       string         `gorm:"type:varchar(100);not null" json:"name"`
	Code       string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Type       int8           `gorm:"type:smallint;not null" json:"type"` // 1固定金额 2百分比
	Value      float64        `gorm:"type:decimal(10,2);not null" json:"value"`
	MaxAmount  float64        `gorm:"type:decimal(10,2);default:0" json:"max_amount"`
	MinOrder   float64        `gorm:"type:decimal(10,2);default:0" json:"min_order"`
	TotalCount int            `gorm:"default:0" json:"total_count"` // 0不限
	UsedCount  int            `gorm:"default:0" json:"used_count"`
	UserCount  int            `gorm:"default:1" json:"user_count"`
	ProductIDs datatypes.JSON `gorm:"type:json" json:"product_ids"`
	StartDate  time.Time      `json:"start_date"`
	EndDate    time.Time      `json:"end_date"`
	IsActive   bool           `gorm:"default:true" json:"is_active"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// VoucherRecord 代金券使用记录
type VoucherRecord struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	VoucherID uint      `gorm:"index;not null" json:"voucher_id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	OrderID   uint      `gorm:"index;not null" json:"order_id"`
	Amount    float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

// UserVoucher 用户领取的代金券
type UserVoucher struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"index;not null" json:"user_id"`
	VoucherID uint       `gorm:"index;not null" json:"voucher_id"`
	UsedAt    *time.Time `json:"used_at"`
	OrderID   *uint      `gorm:"index" json:"order_id"`
	CreatedAt time.Time  `json:"created_at"`
	Voucher   *Voucher   `gorm:"foreignKey:VoucherID" json:"voucher,omitempty"`
}

// CreateVoucherRequest 创建代金券请求
type CreateVoucherRequest struct {
	Name       string  `json:"name" binding:"required,max=100"`
	Code       string  `json:"code" binding:"required,max=50"`
	Type       int8    `json:"type" binding:"required,oneof=1 2"`
	Value      float64 `json:"value" binding:"required,gt=0"`
	MaxAmount  float64 `json:"max_amount"`
	MinOrder   float64 `json:"min_order"`
	TotalCount int     `json:"total_count"`
	UserCount  int     `json:"user_count"`
	ProductIDs []uint  `json:"product_ids"`
	StartDate  string  `json:"start_date" binding:"required"`
	EndDate    string  `json:"end_date" binding:"required"`
}

// UpdateVoucherRequest 更新代金券请求
type UpdateVoucherRequest struct {
	Name       *string  `json:"name"`
	Code       *string  `json:"code"`
	Type       *int8    `json:"type"`
	Value      *float64 `json:"value"`
	MaxAmount  *float64 `json:"max_amount"`
	MinOrder   *float64 `json:"min_order"`
	TotalCount *int     `json:"total_count"`
	UserCount  *int     `json:"user_count"`
	ProductIDs []uint   `json:"product_ids"`
	StartDate  *string  `json:"start_date"`
	EndDate    *string  `json:"end_date"`
	IsActive   *bool    `json:"is_active"`
}

// ValidateVoucherRequest 验证代金券请求
type ValidateVoucherRequest struct {
	Code      string  `json:"code" binding:"required"`
	ProductID uint    `json:"product_id"`
	Amount    float64 `json:"amount" binding:"required"`
}

type VoucherService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewVoucherService(db *gorm.DB, log *logger.Logger) *VoucherService {
	return &VoucherService{db: db, log: log}
}

// GetUserVouchers 获取用户可用的代金券
func (s *VoucherService) GetUserVouchers(userID uint) ([]UserVoucher, error) {
	var vouchers []UserVoucher
	now := time.Now()

	err := s.db.Preload("Voucher").
		Where("user_id = ?", userID).
		Where("used_at IS NULL").
		Joins("JOIN vouchers ON vouchers.id = user_vouchers.voucher_id").
		Where("vouchers.is_active = ? AND vouchers.start_date <= ? AND vouchers.end_date >= ?", true, now, now).
		Where("vouchers.deleted_at IS NULL").
		Order("user_vouchers.created_at DESC").
		Find(&vouchers).Error
	if err != nil {
		return nil, err
	}

	// 过滤已用完的代金券
	var result []UserVoucher
	for _, uv := range vouchers {
		if uv.Voucher != nil && uv.Voucher.TotalCount > 0 && uv.Voucher.UsedCount >= uv.Voucher.TotalCount {
			continue
		}
		result = append(result, uv)
	}
	return result, nil
}

// ClaimVoucher 用户领取代金券
func (s *VoucherService) ClaimVoucher(userID, voucherID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var voucher Voucher
		if err := tx.First(&voucher, voucherID).Error; err != nil {
			return errors.New("voucher not found")
		}

		if !voucher.IsActive {
			return errors.New("voucher is not active")
		}

		now := time.Now()
		if now.Before(voucher.StartDate) || now.After(voucher.EndDate) {
			return errors.New("voucher is not in valid date range")
		}

		if voucher.TotalCount > 0 && voucher.UsedCount >= voucher.TotalCount {
			return errors.New("voucher is fully claimed")
		}

		// 检查用户领取次数
		var count int64
		tx.Model(&UserVoucher{}).Where("user_id = ? AND voucher_id = ?", userID, voucherID).Count(&count)
		if voucher.UserCount > 0 && int(count) >= voucher.UserCount {
			return errors.New("user claim limit reached")
		}

		uv := &UserVoucher{
			UserID:    userID,
			VoucherID: voucherID,
		}
		return tx.Create(uv).Error
	})
}

// ValidateVoucher 验证代金券是否可用于特定订单
func (s *VoucherService) ValidateVoucher(userID uint, code string, productID uint, amount float64) (float64, *Voucher, error) {
	var voucher Voucher
	if err := s.db.Where("code = ? AND is_active = ?", code, true).First(&voucher).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil, errors.New("voucher not found")
		}
		return 0, nil, err
	}

	now := time.Now()
	if now.Before(voucher.StartDate) || now.After(voucher.EndDate) {
		return 0, nil, errors.New("voucher is not in valid date range")
	}

	if voucher.TotalCount > 0 && voucher.UsedCount >= voucher.TotalCount {
		return 0, nil, errors.New("voucher is fully used")
	}

	// 检查产品限制
	if len(voucher.ProductIDs) > 0 {
		var productIDs []uint
		if err := json.Unmarshal(voucher.ProductIDs, &productIDs); err == nil && len(productIDs) > 0 {
			found := false
			for _, pid := range productIDs {
				if pid == productID {
					found = true
					break
				}
			}
			if !found {
				return 0, nil, errors.New("voucher is not applicable to this product")
			}
		}
	}

	// 检查最低消费
	if voucher.MinOrder > 0 && amount < voucher.MinOrder {
		return 0, nil, fmt.Errorf("order amount %.2f is below minimum %.2f", amount, voucher.MinOrder)
	}

	// 检查用户领取情况
	var userVoucher UserVoucher
	err := s.db.Where("user_id = ? AND voucher_id = ? AND used_at IS NULL", userID, voucher.ID).
		First(&userVoucher).Error
	if err != nil {
		return 0, nil, errors.New("user does not have this voucher")
	}

	// 计算抵扣金额
	var discount float64
	switch voucher.Type {
	case 1: // 固定金额
		discount = voucher.Value
	case 2: // 百分比
		discount = amount * voucher.Value / 100
		if voucher.MaxAmount > 0 && discount > voucher.MaxAmount {
			discount = voucher.MaxAmount
		}
	default:
		return 0, nil, errors.New("unknown voucher type")
	}

	if discount > amount {
		discount = amount
	}
	if discount < 0 {
		discount = 0
	}

	return discount, &voucher, nil
}

// UseVoucher 标记代金券已使用
func (s *VoucherService) UseVoucher(userID, voucherID, orderID uint, amount float64) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		// 更新用户代金券记录
		result := tx.Model(&UserVoucher{}).
			Where("user_id = ? AND voucher_id = ? AND used_at IS NULL", userID, voucherID).
			Updates(map[string]interface{}{
				"used_at":  now,
				"order_id": orderID,
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("no available voucher to use")
		}

		// 增加使用次数
		if err := tx.Model(&Voucher{}).Where("id = ?", voucherID).
			UpdateColumn("used_count", gorm.Expr("used_count + 1")).Error; err != nil {
			return err
		}

		// 创建使用记录
		record := &VoucherRecord{
			VoucherID: voucherID,
			UserID:    userID,
			OrderID:   orderID,
			Amount:    amount,
		}
		return tx.Create(record).Error
	})
}

// Create 创建代金券 (admin)
func (s *VoucherService) Create(req CreateVoucherRequest) (*Voucher, error) {
	var count int64
	s.db.Model(&Voucher{}).Where("code = ?", req.Code).Count(&count)
	if count > 0 {
		return nil, errors.New("voucher code already exists")
	}

	startDate, err := util.ParseTime(req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date: %w", err)
	}
	endDate, err := util.ParseTime(req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end_date: %w", err)
	}

	var productIDs datatypes.JSON
	if len(req.ProductIDs) > 0 {
		data, err := json.Marshal(req.ProductIDs)
		if err != nil {
			return nil, err
		}
		productIDs = datatypes.JSON(data)
	}

	userCount := req.UserCount
	if userCount < 1 {
		userCount = 1
	}

	voucher := &Voucher{
		Name:       req.Name,
		Code:       req.Code,
		Type:       req.Type,
		Value:      req.Value,
		MaxAmount:  req.MaxAmount,
		MinOrder:   req.MinOrder,
		TotalCount: req.TotalCount,
		UserCount:  userCount,
		ProductIDs: productIDs,
		StartDate:  startDate,
		EndDate:    endDate,
		IsActive:   true,
	}

	if err := s.db.Create(voucher).Error; err != nil {
		return nil, err
	}
	return voucher, nil
}

// GetByID 根据ID获取代金券
func (s *VoucherService) GetByID(id uint) (*Voucher, error) {
	var voucher Voucher
	if err := s.db.First(&voucher, id).Error; err != nil {
		return nil, err
	}
	return &voucher, nil
}

// Update 更新代金券 (admin)
func (s *VoucherService) Update(id uint, req UpdateVoucherRequest) (*Voucher, error) {
	voucher, err := s.GetByID(id)
	if err != nil {
		return nil, err
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
	if req.Value != nil {
		updates["value"] = *req.Value
	}
	if req.MaxAmount != nil {
		updates["max_amount"] = *req.MaxAmount
	}
	if req.MinOrder != nil {
		updates["min_order"] = *req.MinOrder
	}
	if req.TotalCount != nil {
		updates["total_count"] = *req.TotalCount
	}
	if req.UserCount != nil {
		updates["user_count"] = *req.UserCount
	}
	if req.ProductIDs != nil {
		data, err := json.Marshal(req.ProductIDs)
		if err == nil {
			updates["product_ids"] = datatypes.JSON(data)
		}
	}
	if req.StartDate != nil {
		t, err := util.ParseTime(*req.StartDate)
		if err == nil {
			updates["start_date"] = t
		}
	}
	if req.EndDate != nil {
		t, err := util.ParseTime(*req.EndDate)
		if err == nil {
			updates["end_date"] = t
		}
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(updates) > 0 {
		if err := s.db.Model(voucher).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	return s.GetByID(id)
}

// Delete 删除代金券 (admin)
func (s *VoucherService) Delete(id uint) error {
	result := s.db.Delete(&Voucher{}, id)
	if result.RowsAffected == 0 {
		return errors.New("voucher not found")
	}
	return result.Error
}

// GetAll 获取所有代金券 (admin, 分页)
func (s *VoucherService) GetAll(page, pageSize int) ([]Voucher, int64, error) {
	var vouchers []Voucher
	var total int64

	s.db.Model(&Voucher{}).Count(&total)
	offset, limit := Paginate(page, pageSize)
	err := s.db.Offset(offset).Limit(limit).Order("id DESC").Find(&vouchers).Error
	return vouchers, total, err
}
