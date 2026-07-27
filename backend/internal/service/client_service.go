package service

import (
	"errors"
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
