package service

import (
	"errors"
	"time"

	"github.com/anchor-finance/backend/internal/model"
	"github.com/anchor-finance/backend/pkg/logger"

	"gorm.io/gorm"
)

type CertificationService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewCertificationService(db *gorm.DB, log *logger.Logger) *CertificationService {
	return &CertificationService{db: db, log: log}
}

type SubmitCertificationRequest struct {
	Type            string `json:"type" binding:"required,oneof=individual enterprise"`
	RealName        string `json:"real_name" binding:"required,max=50"`
	IDCard          string `json:"id_card" binding:"omitempty,max=50"`
	FrontImage      string `json:"front_image" binding:"omitempty,max=255"`
	BackImage       string `json:"back_image" binding:"omitempty,max=255"`
	HandImage       string `json:"hand_image" binding:"omitempty,max=255"`
	EnterpriseName  string `json:"enterprise_name" binding:"omitempty,max=100"`
	BusinessLicense string `json:"business_license" binding:"omitempty,max=255"`
}

type ReviewCertificationRequest struct {
	Status       int8   `json:"status" binding:"required,oneof=2 3"` // 2通过 3拒绝
	RejectReason string `json:"reject_reason" binding:"omitempty"`
}

// Submit creates or updates a certification request.
func (s *CertificationService) Submit(userID uint, req SubmitCertificationRequest) (*model.Certification, error) {
	var cert model.Certification
	err := s.db.Where("user_id = ?", userID).First(&cert).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		cert = model.Certification{
			UserID:          userID,
			Type:            req.Type,
			RealName:        req.RealName,
			IDCard:          req.IDCard,
			FrontImage:      req.FrontImage,
			BackImage:       req.BackImage,
			HandImage:       req.HandImage,
			EnterpriseName:  req.EnterpriseName,
			BusinessLicense: req.BusinessLicense,
			Status:          1,
		}
		if err := s.db.Create(&cert).Error; err != nil {
			return nil, err
		}
		s.log.Infof("certification submitted: user=%d type=%s", userID, req.Type)
		return &cert, nil
	}
	if err != nil {
		return nil, err
	}

	// Only allow re-submission if previously rejected
	if cert.Status != 3 {
		return nil, errors.New("certification is pending or already approved")
	}

	updates := map[string]interface{}{
		"type":             req.Type,
		"real_name":        req.RealName,
		"id_card":          req.IDCard,
		"front_image":      req.FrontImage,
		"back_image":       req.BackImage,
		"hand_image":       req.HandImage,
		"enterprise_name":  req.EnterpriseName,
		"business_license": req.BusinessLicense,
		"status":           1,
		"reject_reason":    "",
		"reviewed_by":      nil,
		"reviewed_at":      nil,
	}
	if err := s.db.Model(&cert).Updates(updates).Error; err != nil {
		return nil, err
	}

	s.log.Infof("certification re-submitted: user=%d type=%s", userID, req.Type)
	return s.GetByUserID(userID)
}

// GetByUserID returns the certification for a user.
func (s *CertificationService) GetByUserID(userID uint) (*model.Certification, error) {
	var cert model.Certification
	if err := s.db.Where("user_id = ?", userID).First(&cert).Error; err != nil {
		return nil, err
	}
	return &cert, nil
}

// GetList returns a paginated certification list (admin).
func (s *CertificationService) GetList(page, pageSize int, status int) ([]model.Certification, int64, error) {
	var certs []model.Certification
	var total int64

	query := s.db.Model(&model.Certification{})
	if status > 0 {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Preload("User").Find(&certs).Error; err != nil {
		return nil, 0, err
	}

	return certs, total, nil
}

// Review approves or rejects a certification (admin).
func (s *CertificationService) Review(certID, reviewerID uint, req ReviewCertificationRequest) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":      req.Status,
		"reviewed_by": reviewerID,
		"reviewed_at": &now,
	}
	if req.Status == 3 {
		updates["reject_reason"] = req.RejectReason
	}

	result := s.db.Model(&model.Certification{}).
		Where("id = ? AND status = 1", certID).
		Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("certification not found or already reviewed")
	}

	s.log.Infof("certification reviewed: id=%d status=%d reviewer=%d", certID, req.Status, reviewerID)
	return nil
}
