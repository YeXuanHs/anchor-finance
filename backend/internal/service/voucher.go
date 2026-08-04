package service

import (
	"errors"
	"fmt"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// CreateVoucherRequest 创建发票申请请求
type CreateVoucherRequest struct {
	InvoiceID uint    `json:"invoice_id" binding:"required"`
	PostID    uint    `json:"post_id" binding:"required"`
	TypeID    uint    `json:"type_id" binding:"required"`
	ExpressID uint    `json:"express_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"required,gt=0"`
	Notes     string  `json:"notes" binding:"max=500"`
}

// CreateVoucherTypeRequest 创建发票抬头请求
type CreateVoucherTypeRequest struct {
	Title       string `json:"title" binding:"required,max=50"`
	IssueType   string `json:"issue_type" binding:"required,oneof=person company"`
	InvoiceType string `json:"invoice_type" binding:"required,oneof=common dedicated"`
	TaxID       string `json:"tax_id" binding:"max=100"`
	Bank        string `json:"bank" binding:"max=100"`
	Account     string `json:"account" binding:"max=100"`
	Address     string `json:"address" binding:"max=100"`
	Phone       string `json:"phone" binding:"max=100"`
}

// CreateVoucherPostRequest 创建收件地址请求
type CreateVoucherPostRequest struct {
	Username string `json:"username" binding:"required,max=50"`
	Phone    string `json:"phone" binding:"required,max=50"`
	Province string `json:"province" binding:"required,max=100"`
	City     string `json:"city" binding:"required,max=100"`
	Region   string `json:"region" binding:"required,max=100"`
	Detail   string `json:"detail" binding:"required,max=500"`
	Post     string `json:"post" binding:"max=50"`
}

type VoucherService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewVoucherService(db *gorm.DB, log *logger.Logger) *VoucherService {
	return &VoucherService{db: db, log: log}
}

// ==================== 管理员接口 ====================

// GetRateConfig 获取费率配置
func (s *VoucherService) GetRateConfig() (map[string]interface{}, error) {
	manager := s.getConfigValue("voucher_manager", "0")
	rate := s.getConfigValue("voucher_rate", "0")
	return map[string]interface{}{
		"voucher_manager": manager,
		"voucher_rate":    rate,
	}, nil
}

// UpdateRateConfig 更新费率配置
func (s *VoucherService) UpdateRateConfig(manager int, rate float64) error {
	err := s.setConfigValue("voucher_manager", fmt.Sprintf("%d", manager))
	if err != nil {
		return err
	}
	return s.setConfigValue("voucher_rate", fmt.Sprintf("%.2f", rate))
}

// GetVoucherList 获取发票申请列表
func (s *VoucherService) GetVoucherList(page, pageSize int, status, order, sort string) ([]map[string]interface{}, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	var total int64
	query := s.db.Model(&model.Voucher{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	query.Count(&total)

	var vouchers []model.Voucher
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order(order + " " + sort).Find(&vouchers).Error
	if err != nil {
		return nil, 0, err
	}

	var results []map[string]interface{}
	for _, v := range vouchers {
		item := map[string]interface{}{
			"id":          v.ID,
			"uid":         v.UID,
			"invoice_id":  v.InvoiceID,
			"post_id":     v.PostID,
			"type_id":     v.TypeID,
			"express_id":  v.ExpressID,
			"amount":      v.Amount,
			"status":      v.Status,
			"notes":       v.Notes,
			"check_time":  v.CheckTime,
			"create_time": v.CreateTime,
			"update_time": v.UpdateTime,
		}

		// 关联用户名
		if v.User != nil {
			item["username"] = v.User.Username
		} else {
			var user model.User
			if err := s.db.First(&user, v.UID).Error; err == nil {
				item["username"] = user.Username
			}
		}

		// 关联发票抬头
		if v.VoucherType != nil {
			item["voucher_type"] = v.VoucherType
		} else if v.TypeID > 0 {
			var vt model.VoucherType
			if err := s.db.First(&vt, v.TypeID).Error; err == nil {
				item["voucher_type"] = vt
			}
		}

		// 关联收件地址
		if v.Post != nil {
			item["post"] = v.Post
		} else if v.PostID > 0 {
			var vp model.VoucherPost
			if err := s.db.First(&vp, v.PostID).Error; err == nil {
				item["post"] = vp
			}
		}

		// 关联快递
		if v.Express != nil {
			item["express"] = v.Express
		} else if v.ExpressID > 0 {
			var ex model.Express
			if err := s.db.First(&ex, v.ExpressID).Error; err == nil {
				item["express"] = ex
			}
		}

		// 关联发票
		if v.Invoice != nil {
			item["invoice"] = v.Invoice
		} else if v.InvoiceID > 0 {
			var inv model.Invoice
			if err := s.db.First(&inv, v.InvoiceID).Error; err == nil {
				item["invoice"] = inv
			}
		}

		results = append(results, item)
	}

	return results, total, nil
}

// GetVoucherDetail 获取发票申请详情
func (s *VoucherService) GetVoucherDetail(id uint) (map[string]interface{}, error) {
	var voucher model.Voucher
	if err := s.db.First(&voucher, id).Error; err != nil {
		return nil, errors.New("voucher not found")
	}

	result := map[string]interface{}{
		"id":          voucher.ID,
		"uid":         voucher.UID,
		"invoice_id":  voucher.InvoiceID,
		"post_id":     voucher.PostID,
		"type_id":     voucher.TypeID,
		"express_id":  voucher.ExpressID,
		"amount":      voucher.Amount,
		"status":      voucher.Status,
		"notes":       voucher.Notes,
		"check_time":  voucher.CheckTime,
		"create_time": voucher.CreateTime,
		"update_time": voucher.UpdateTime,
	}

	// 关联用户名
	var user model.User
	if err := s.db.First(&user, voucher.UID).Error; err == nil {
		result["username"] = user.Username
		result["user"] = map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		}
	}

	// 关联发票抬头
	if voucher.TypeID > 0 {
		var vt model.VoucherType
		if err := s.db.First(&vt, voucher.TypeID).Error; err == nil {
			result["voucher_type"] = vt
		}
	}

	// 关联收件地址
	if voucher.PostID > 0 {
		var vp model.VoucherPost
		if err := s.db.First(&vp, voucher.PostID).Error; err == nil {
			result["post"] = vp
		}
	}

	// 关联快递
	if voucher.ExpressID > 0 {
		var ex model.Express
		if err := s.db.First(&ex, voucher.ExpressID).Error; err == nil {
			result["express"] = ex
		}
	}

	// 关联发票账单详情
	if voucher.InvoiceID > 0 {
		var inv model.Invoice
		if err := s.db.Preload("Items").First(&inv, voucher.InvoiceID).Error; err == nil {
			result["invoice"] = inv
			// 计算税额
			rate := s.getRateValue()
			taxAmount := inv.SubTotal * rate / 100
			result["tax_amount"] = taxAmount
			result["tax_rate"] = rate
			result["total_with_tax"] = inv.SubTotal + taxAmount
		}
	}

	return result, nil
}

// UpdateVoucherStatus 更新发票申请状态
func (s *VoucherService) UpdateVoucherStatus(id uint, status, notes string) error {
	var voucher model.Voucher
	if err := s.db.First(&voucher, id).Error; err != nil {
		return errors.New("voucher not found")
	}

	if voucher.Status != "Pending" && voucher.Status != "Unpaid" {
		return errors.New("only Pending or Unpaid vouchers can be updated")
	}

	updates := map[string]interface{}{
		"status":      status,
		"check_time":  time.Now().Unix(),
		"update_time": time.Now().Unix(),
	}
	if notes != "" {
		updates["notes"] = notes
	}

	return s.db.Model(&voucher).Updates(updates).Error
}

// ==================== 用户接口 ====================

// GetUserVoucherList 获取用户发票申请列表
func (s *VoucherService) GetUserVoucherList(userID uint, page, pageSize int) ([]model.Voucher, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	var total int64
	s.db.Model(&model.Voucher{}).Where("uid = ?", userID).Count(&total)

	var vouchers []model.Voucher
	offset := (page - 1) * pageSize
	err := s.db.Where("uid = ?", userID).
		Preload("VoucherType").
		Preload("Post").
		Preload("Express").
		Preload("Invoice").
		Offset(offset).Limit(pageSize).
		Order("create_time DESC").
		Find(&vouchers).Error

	return vouchers, total, err
}

// CreateUserVoucher 创建发票申请
func (s *VoucherService) CreateUserVoucher(userID uint, req CreateVoucherRequest) (*model.Voucher, error) {
	now := time.Now().Unix()
	voucher := &model.Voucher{
		UID:        userID,
		InvoiceID:  req.InvoiceID,
		PostID:     req.PostID,
		TypeID:     req.TypeID,
		ExpressID:  req.ExpressID,
		Amount:     req.Amount,
		Status:     "Pending",
		Notes:      req.Notes,
		CreateTime: now,
		UpdateTime: now,
	}

	if err := s.db.Create(voucher).Error; err != nil {
		return nil, err
	}

	// 创建发票关联
	if req.InvoiceID > 0 {
		vi := &model.VoucherInvoice{
			VoucherID: voucher.ID,
			InvoiceID: req.InvoiceID,
			CreatedAt: time.Now(),
		}
		if err := s.db.Create(vi).Error; err != nil {
			return nil, fmt.Errorf("link voucher to invoice: %w", err)
		}
	}

	return voucher, nil
}

// ==================== 发票抬头 CRUD ====================

// GetVoucherTypes 获取用户发票抬头列表
func (s *VoucherService) GetVoucherTypes(userID uint) ([]model.VoucherType, error) {
	var types []model.VoucherType
	err := s.db.Where("uid = ?", userID).Order("create_time DESC").Find(&types).Error
	return types, err
}

// CreateVoucherType 创建发票抬头
func (s *VoucherService) CreateVoucherType(userID uint, req CreateVoucherTypeRequest) (*model.VoucherType, error) {
	now := time.Now().Unix()
	vt := &model.VoucherType{
		UID:         userID,
		Title:       req.Title,
		IssueType:   req.IssueType,
		InvoiceType: req.InvoiceType,
		TaxID:       req.TaxID,
		Bank:        req.Bank,
		Account:     req.Account,
		Address:     req.Address,
		Phone:       req.Phone,
		CreateTime:  now,
		UpdateTime:  now,
	}

	if err := s.db.Create(vt).Error; err != nil {
		return nil, err
	}
	return vt, nil
}

// UpdateVoucherType 更新发票抬头
func (s *VoucherService) UpdateVoucherType(userID, id uint, req CreateVoucherTypeRequest) (*model.VoucherType, error) {
	var vt model.VoucherType
	if err := s.db.Where("id = ? AND uid = ?", id, userID).First(&vt).Error; err != nil {
		return nil, errors.New("voucher type not found")
	}

	vt.Title = req.Title
	vt.IssueType = req.IssueType
	vt.InvoiceType = req.InvoiceType
	vt.TaxID = req.TaxID
	vt.Bank = req.Bank
	vt.Account = req.Account
	vt.Address = req.Address
	vt.Phone = req.Phone
	vt.UpdateTime = time.Now().Unix()

	if err := s.db.Save(&vt).Error; err != nil {
		return nil, err
	}
	return &vt, nil
}

// DeleteVoucherType 删除发票抬头
func (s *VoucherService) DeleteVoucherType(userID, id uint) error {
	result := s.db.Where("id = ? AND uid = ?", id, userID).Delete(&model.VoucherType{})
	if result.RowsAffected == 0 {
		return errors.New("voucher type not found")
	}
	return result.Error
}

// ==================== 收件地址 CRUD ====================

// GetVoucherPosts 获取用户收件地址列表
func (s *VoucherService) GetVoucherPosts(userID uint) ([]model.VoucherPost, error) {
	var posts []model.VoucherPost
	err := s.db.Where("uid = ?", userID).Order("is_default DESC, create_time DESC").Find(&posts).Error
	return posts, err
}

// CreateVoucherPost 创建收件地址
func (s *VoucherService) CreateVoucherPost(userID uint, req CreateVoucherPostRequest) (*model.VoucherPost, error) {
	now := time.Now().Unix()
	vp := &model.VoucherPost{
		UID:        userID,
		Username:   req.Username,
		Phone:      req.Phone,
		Province:   req.Province,
		City:       req.City,
		Region:     req.Region,
		Detail:     req.Detail,
		Post:       req.Post,
		CreateTime: now,
		UpdateTime: now,
	}

	if err := s.db.Create(vp).Error; err != nil {
		return nil, err
	}
	return vp, nil
}

// UpdateVoucherPost 更新收件地址
func (s *VoucherService) UpdateVoucherPost(userID, id uint, req CreateVoucherPostRequest) (*model.VoucherPost, error) {
	var vp model.VoucherPost
	if err := s.db.Where("id = ? AND uid = ?", id, userID).First(&vp).Error; err != nil {
		return nil, errors.New("voucher post not found")
	}

	vp.Username = req.Username
	vp.Phone = req.Phone
	vp.Province = req.Province
	vp.City = req.City
	vp.Region = req.Region
	vp.Detail = req.Detail
	vp.Post = req.Post
	vp.UpdateTime = time.Now().Unix()

	if err := s.db.Save(&vp).Error; err != nil {
		return nil, err
	}
	return &vp, nil
}

// DeleteVoucherPost 删除收件地址
func (s *VoucherService) DeleteVoucherPost(userID, id uint) error {
	result := s.db.Where("id = ? AND uid = ?", id, userID).Delete(&model.VoucherPost{})
	if result.RowsAffected == 0 {
		return errors.New("voucher post not found")
	}
	return result.Error
}

// ==================== 内部方法 ====================

func (s *VoucherService) getConfigValue(key, defaultValue string) string {
	var config model.SystemConfig
	if err := s.db.Where("`key` = ?", key).First(&config).Error; err != nil {
		return defaultValue
	}
	if config.Value == "" {
		return defaultValue
	}
	return config.Value
}

func (s *VoucherService) setConfigValue(key, value string) error {
	result := s.db.Model(&model.SystemConfig{}).Where("`key` = ?", key).Update("value", value)
	if result.RowsAffected == 0 {
		return s.db.Create(&model.SystemConfig{
			Key:   key,
			Value: value,
			Group: "voucher",
			Name:  key,
			Type:  "string",
		}).Error
	}
	return result.Error
}

func (s *VoucherService) getRateValue() float64 {
	val := s.getConfigValue("voucher_rate", "0")
	var rate float64
	fmt.Sscanf(val, "%f", &rate)
	return rate
}
