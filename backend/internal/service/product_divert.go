package service

import (
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
	FromUserID     uint       `gorm:"index;not null" json:"from_user_id"`
	FromUser       User       `gorm:"foreignKey:FromUserID" json:"from_user,omitempty"`
	ToUserID       uint       `gorm:"index;not null" json:"to_user_id"`
	ToUser         User       `gorm:"foreignKey:ToUserID" json:"to_user,omitempty"`
	UserProductID  uint       `gorm:"index;not null" json:"user_product_id"`
	ProductID      uint       `gorm:"index;not null" json:"product_id"`
	ProductName    string     `gorm:"type:varchar(256)" json:"product_name"`
	Reason         string     `gorm:"type:text" json:"reason"`
	Status         int16      `gorm:"type:smallint;default:0;not null;index" json:"status"` // 0=待确认 1=已接受 2=已拒绝 3=已取消
	ProcessedAt    *time.Time `json:"processed_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type ProductDivertService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewProductDivertService(db *gorm.DB, log *logger.Logger) *ProductDivertService {
	return &ProductDivertService{db: db, log: log}
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

	// Verify the product belongs to the sender
	var userProduct struct {
		ID        uint
		UserID    uint
		ProductID uint
		Name      string
	}
	if err := s.db.Table("user_products").Where("id = ? AND user_id = ?", req.UserProductID, fromUserID).
		First(&userProduct).Error; err != nil {
		return nil, errors.New("product not found or not owned by you")
	}

	// Find recipient by email
	var toUser struct {
		ID    uint
		Email string
	}
	if err := s.db.Table("users").Where("email = ?", req.ToUserEmail).First(&toUser).Error; err != nil {
		return nil, errors.New("recipient user not found")
	}

	if toUser.ID == fromUserID {
		return nil, errors.New("cannot transfer to yourself")
	}

	transfer := &ProductTransfer{
		TransferNo:    fmt.Sprintf("TRF%s%04d", time.Now().Format("20060102150405"), time.Now().UnixNano()%10000),
		FromUserID:    fromUserID,
		ToUserID:      toUser.ID,
		UserProductID: req.UserProductID,
		ProductID:     userProduct.ProductID,
		ProductName:   userProduct.Name,
		Reason:        req.Reason,
		Status:        0,
	}

	if err := s.db.Create(transfer).Error; err != nil {
		return nil, err
	}

	s.log.Infof("product transfer created: %s from=%d to=%d", transfer.TransferNo, fromUserID, toUser.ID)
	return transfer, nil
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

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Transfer product ownership
		if err := tx.Table("user_products").Where("id = ?", transfer.UserProductID).
			Update("user_id", transfer.ToUserID).Error; err != nil {
			return err
		}

		now := time.Now()
		return tx.Model(&transfer).Updates(map[string]interface{}{
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

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").
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

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").
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
