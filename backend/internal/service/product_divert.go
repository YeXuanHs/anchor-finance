package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// ProductTransfer 产品转移记录
type ProductTransfer struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	TransferNo     string     `gorm:"type:varchar(64);uniqueIndex;not null" json:"transfer_no"`
	TransferCode   string     `gorm:"type:varchar(32);index" json:"transfer_code"` // 转移码（类似域名EPP码）
	FromUserID     uint       `gorm:"index;not null" json:"from_user_id"`
	FromUser       User       `gorm:"foreignKey:FromUserID" json:"from_user,omitempty"`
	ToUserID       uint       `gorm:"index;not null" json:"to_user_id"`
	ToUser         User       `gorm:"foreignKey:ToUserID" json:"to_user,omitempty"`
	UserProductID  uint       `gorm:"index;not null" json:"user_product_id"`
	ProductID      uint       `gorm:"index;not null" json:"product_id"`
	ProductName    string     `gorm:"type:varchar(256)" json:"product_name"`
	TransferFee    float64    `gorm:"type:decimal(10,2);default:0" json:"transfer_fee"` // 转移费用
	Reason         string     `gorm:"type:text" json:"reason"`
	Status         int16      `gorm:"type:smallint;default:0;not null;index" json:"status"` // 0=待确认 1=已接受 2=已拒绝 3=已取消 4=已过期
	ProcessedAt    *time.Time `json:"processed_at"`
	ExpiresAt      time.Time  `json:"expires_at"` // 转移码过期时间
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// TransferConfig 转移配置
type TransferConfig struct {
	ID                 uint    `gorm:"primaryKey" json:"id"`
	Enabled            bool    `gorm:"default:true" json:"enabled"`                        // 是否启用转移功能
	TransferFee        float64 `gorm:"type:decimal(10,2);default:0" json:"transfer_fee"`   // 默认转移费用
	MinHoldDays        int     `gorm:"default:30" json:"min_hold_days"`                    // 最少持有天数
	RequireRealName    bool    `gorm:"default:true" json:"require_real_name"`               // 转入方需要实名
	RequirePhone       bool    `gorm:"default:false" json:"require_phone"`                  // 转入方需要绑定手机
	CodeExpiryHours    int     `gorm:"default:72" json:"code_expiry_hours"`                // 转移码有效期（小时）
	MaxActiveTransfers int     `gorm:"default:3" json:"max_active_transfers"`               // 最大进行中的转移数
	UpdatedAt          time.Time `json:"updated_at"`
}

type ProductDivertService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewProductDivertService(db *gorm.DB, log *logger.Logger) *ProductDivertService {
	return &ProductDivertService{db: db, log: log}
}

// GetDB 获取数据库连接
func (s *ProductDivertService) GetDB() *gorm.DB {
	return s.db
}

// GetConfig 获取转移配置
func (s *ProductDivertService) GetConfig() *TransferConfig {
	var config TransferConfig
	if err := s.db.First(&config).Error; err != nil {
		config = TransferConfig{
			ID:                 1,
			Enabled:            true,
			TransferFee:        0,
			MinHoldDays:        30,
			RequireRealName:    true,
			RequirePhone:       false,
			CodeExpiryHours:    72,
			MaxActiveTransfers: 3,
		}
		s.db.Create(&config)
	}
	return &config
}

// SaveConfig 保存转移配置
func (s *ProductDivertService) SaveConfig(config *TransferConfig) error {
	config.ID = 1
	return s.db.Save(config).Error
}

// generateTransferCode 生成转移码
func generateTransferCode() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

type CreateTransferRequest struct {
	UserProductID uint   `json:"user_product_id" binding:"required"`
	ToUserEmail   string `json:"to_user_email" binding:"required,email"`
	Reason        string `json:"reason" binding:"omitempty,max=500"`
}

// Create creates a product transfer request.
func (s *ProductDivertService) Create(fromUserID uint, req CreateTransferRequest) (*ProductTransfer, error) {
	if fromUserID == 0 {
		return nil, errors.New("invalid sender")
	}

	config := s.GetConfig()
	if !config.Enabled {
		return nil, errors.New("transfer feature is disabled")
	}

	// 检查进行中的转移数量
	var activeCount int64
	s.db.Model(&ProductTransfer{}).Where("from_user_id = ? AND status = 0", fromUserID).Count(&activeCount)
	if int(activeCount) >= config.MaxActiveTransfers {
		return nil, fmt.Errorf("进行中的转移数量已达上限(%d)", config.MaxActiveTransfers)
	}

	// Verify the product belongs to the sender
	var userProduct struct {
		ID        uint
		UserID    uint
		ProductID uint
		Name      string
		CreatedAt time.Time
	}
	if err := s.db.Table("user_products").Where("id = ? AND user_id = ?", req.UserProductID, fromUserID).
		First(&userProduct).Error; err != nil {
		return nil, errors.New("product not found or not owned by you")
	}

	// 检查持有天数
	holdDays := int(time.Since(userProduct.CreatedAt).Hours() / 24)
	if holdDays < config.MinHoldDays {
		return nil, fmt.Errorf("产品持有天数不足，需要%d天，当前%d天", config.MinHoldDays, holdDays)
	}

	// Find recipient by email
	var toUser struct {
		ID       uint
		Email    string
		RealName string
		Phone    string
	}
	if err := s.db.Table("users").Where("email = ?", req.ToUserEmail).First(&toUser).Error; err != nil {
		return nil, errors.New("recipient user not found")
	}

	if toUser.ID == fromUserID {
		return nil, errors.New("cannot transfer to yourself")
	}

	// 检查目标用户是否满足产品购买条件
	var product struct {
		RequireRealName int
		RequirePhone    int
	}
	s.db.Table("products").Where("id = ?", userProduct.ProductID).First(&product)

	if config.RequireRealName || product.RequireRealName == 1 {
		if toUser.RealName == "" {
			return nil, errors.New("recipient must complete real name verification")
		}
	}

	if config.RequirePhone || product.RequirePhone == 1 {
		if toUser.Phone == "" {
			return nil, errors.New("recipient must bind phone number")
		}
	}

	// 生成转移码
	transferCode := generateTransferCode()
	codeExpiry := time.Now().Add(time.Duration(config.CodeExpiryHours) * time.Hour)

	transfer := &ProductTransfer{
		TransferNo:    fmt.Sprintf("TRF%s%04d", time.Now().Format("20060102150405"), time.Now().UnixNano()%10000),
		TransferCode:  transferCode,
		FromUserID:    fromUserID,
		ToUserID:      toUser.ID,
		UserProductID: req.UserProductID,
		ProductID:     userProduct.ProductID,
		ProductName:   userProduct.Name,
		TransferFee:   config.TransferFee,
		Reason:        req.Reason,
		Status:        0,
		ExpiresAt:     codeExpiry,
	}

	if err := s.db.Create(transfer).Error; err != nil {
		return nil, err
	}

	s.log.Infof("product transfer created: %s from=%d to=%d", transfer.TransferNo, fromUserID, toUser.ID)
	return transfer, nil
}

// GetTransferCode 获取转移码（只有发起人可以查看）
func (s *ProductDivertService) GetTransferCode(userID, transferID uint) (string, error) {
	var transfer ProductTransfer
	if err := s.db.First(&transfer, transferID).Error; err != nil {
		return "", errors.New("transfer not found")
	}

	if transfer.FromUserID != userID {
		return "", errors.New("not authorized")
	}

	if transfer.Status != 0 {
		return "", errors.New("transfer is not pending")
	}

	if time.Now().After(transfer.ExpiresAt) {
		return "", errors.New("transfer code has expired")
	}

	return transfer.TransferCode, nil
}

// RegenerateCode 重新生成转移码
func (s *ProductDivertService) RegenerateCode(userID, transferID uint) (string, error) {
	var transfer ProductTransfer
	if err := s.db.First(&transfer, transferID).Error; err != nil {
		return "", errors.New("transfer not found")
	}

	if transfer.FromUserID != userID {
		return "", errors.New("not authorized")
	}

	if transfer.Status != 0 {
		return "", errors.New("transfer is not pending")
	}

	config := s.GetConfig()
	newCode := generateTransferCode()
	newExpiry := time.Now().Add(time.Duration(config.CodeExpiryHours) * time.Hour)

	if err := s.db.Model(&transfer).Updates(map[string]interface{}{
		"transfer_code": newCode,
		"expires_at":    newExpiry,
	}).Error; err != nil {
		return "", err
	}

	return newCode, nil
}

// AcceptByCode 通过转移码接受转移
func (s *ProductDivertService) AcceptByCode(userID uint, transferCode string) error {
	var transfer ProductTransfer
	if err := s.db.Where("transfer_code = ? AND status = 0", transferCode).First(&transfer).Error; err != nil {
		return errors.New("invalid transfer code")
	}

	if transfer.ToUserID != userID {
		return errors.New("not authorized to accept this transfer")
	}

	if time.Now().After(transfer.ExpiresAt) {
		return errors.New("transfer code has expired")
	}

	return s.executeTransfer(&transfer)
}

// Accept accepts a product transfer request.
func (s *ProductDivertService) Accept(userID, transferID uint) error {
	var transfer ProductTransfer
	if err := s.db.First(&transfer, transferID).Error; err != nil {
		return errors.New("transfer request not found")
	}

	if transfer.ToUserID != userID {
		return errors.New("not authorized to accept this transfer")
	}

	if transfer.Status != 0 {
		return errors.New("transfer is not pending")
	}

	if time.Now().After(transfer.ExpiresAt) {
		return errors.New("transfer has expired")
	}

	return s.executeTransfer(&transfer)
}

// executeTransfer 执行转移
func (s *ProductDivertService) executeTransfer(transfer *ProductTransfer) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Transfer product ownership
		if err := tx.Table("user_products").Where("id = ?", transfer.UserProductID).
			Update("user_id", transfer.ToUserID).Error; err != nil {
			return err
		}

		// 扣除转移费用（如果有的话）
		if transfer.TransferFee > 0 {
			// 从接收方扣除余额
			if err := tx.Exec("UPDATE users SET balance = balance - ? WHERE id = ? AND balance >= ?",
				transfer.TransferFee, transfer.ToUserID, transfer.TransferFee).Error; err != nil {
				return errors.New("insufficient balance for transfer fee")
			}
		}

		now := time.Now()
		return tx.Model(transfer).Updates(map[string]interface{}{
			"status":       1,
			"processed_at": &now,
		}).Error
	})
}

// Reject rejects a product transfer request.
func (s *ProductDivertService) Reject(userID, transferID uint) error {
	var transfer ProductTransfer
	if err := s.db.First(&transfer, transferID).Error; err != nil {
		return errors.New("transfer request not found")
	}

	if transfer.ToUserID != userID {
		return errors.New("not authorized to reject this transfer")
	}

	if transfer.Status != 0 {
		return errors.New("transfer is not pending")
	}

	now := time.Now()
	return s.db.Model(&transfer).Updates(map[string]interface{}{
		"status":       2,
		"processed_at": &now,
	}).Error
}

// Cancel cancels a product transfer (by sender).
func (s *ProductDivertService) Cancel(userID, transferID uint) error {
	var transfer ProductTransfer
	if err := s.db.First(&transfer, transferID).Error; err != nil {
		return errors.New("transfer request not found")
	}

	if transfer.FromUserID != userID {
		return errors.New("not authorized to cancel this transfer")
	}

	if transfer.Status != 0 {
		return errors.New("transfer is not pending")
	}

	now := time.Now()
	return s.db.Model(&transfer).Updates(map[string]interface{}{
		"status":       3,
		"processed_at": &now,
	}).Error
}

// GetSent returns paginated transfers sent by a user.
func (s *ProductDivertService) GetSent(userID uint, page, pageSize int) ([]ProductTransfer, int64, error) {
	var transfers []ProductTransfer
	var total int64

	query := s.db.Model(&ProductTransfer{}).Where("from_user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").
		Preload("ToUser").Find(&transfers).Error; err != nil {
		return nil, 0, err
	}
	return transfers, total, nil
}

// GetReceived returns paginated transfers received by a user.
func (s *ProductDivertService) GetReceived(userID uint, page, pageSize int) ([]ProductTransfer, int64, error) {
	var transfers []ProductTransfer
	var total int64

	query := s.db.Model(&ProductTransfer{}).Where("to_user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").
		Preload("FromUser").Find(&transfers).Error; err != nil {
		return nil, 0, err
	}
	return transfers, total, nil
}

// GetByID returns a transfer request by ID.
func (s *ProductDivertService) GetByID(userID, transferID uint) (*ProductTransfer, error) {
	var transfer ProductTransfer
	if err := s.db.Where("id = ?", transferID).First(&transfer).Error; err != nil {
		return nil, err
	}

	if transfer.FromUserID != userID && transfer.ToUserID != userID {
		return nil, errors.New("not authorized")
	}

	return &transfer, nil
}

// CleanupExpiredTransfers 清理过期的转移请求
func (s *ProductDivertService) CleanupExpiredTransfers() int64 {
	result := s.db.Model(&ProductTransfer{}).
		Where("status = 0 AND expires_at < ?", time.Now()).
		Update("status", 4) // 4=已过期
	return result.RowsAffected
}
