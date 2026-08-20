package service

import (
	"errors"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type ServiceDetailService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewServiceDetailService(db *gorm.DB, log *logger.Logger) *ServiceDetailService {
	return &ServiceDetailService{db: db, log: log}
}

// ServiceDetail represents a detailed service record.
type ServiceDetail struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	UserID       uint       `gorm:"index" json:"user_id"`
	ProductID    uint       `json:"product_id"`
	Domain       string     `gorm:"type:varchar(255)" json:"domain"`
	Username     string     `gorm:"type:varchar(128)" json:"username"`
	Password     string     `gorm:"type:varchar(255)" json:"-"`
	StartDate    *time.Time `json:"start_date"`
	DueDate      *time.Time `json:"due_date"`
	Price        float64    `gorm:"type:decimal(12,2)" json:"price"`
	BillingCycle string     `gorm:"type:varchar(32)" json:"billing_cycle"`
	Status       int16      `gorm:"default:1" json:"status"`
	Note         string     `gorm:"type:text" json:"note"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (s *ServiceDetailService) GetServiceDetail(serviceID, userID uint) (map[string]interface{}, error) {
	var service map[string]interface{}
	err := s.db.Table("client_services").
		Where("id = ? AND user_id = ?", serviceID, userID).
		Scan(&service).Error
	if err != nil {
		return nil, err
	}
	return service, nil
}

func (s *ServiceDetailService) GetServiceStats(userID uint) (map[string]interface{}, error) {
	var total, active, suspended, terminated int64
	s.db.Table("client_services").Where("user_id = ?", userID).Count(&total)
	s.db.Table("client_services").Where("user_id = ? AND status = ?", userID, 1).Count(&active)
	s.db.Table("client_services").Where("user_id = ? AND status = ?", userID, 2).Count(&suspended)
	s.db.Table("client_services").Where("user_id = ? AND status = ?", userID, 3).Count(&terminated)
	return map[string]interface{}{
		"total": total, "active": active, "suspended": suspended, "terminated": terminated,
	}, nil
}

// List returns paginated service details.
func (s *ServiceDetailService) List(page, pageSize int, userID, productID *uint, keyword string, status *int16) ([]ServiceDetail, int64, error) {
	var items []ServiceDetail
	var total int64

	query := s.db.Model(&ServiceDetail{})
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if productID != nil {
		query = query.Where("product_id = ?", *productID)
	}
	if keyword != "" {
		query = query.Where("domain LIKE ? OR username LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&items)
	return items, total, nil
}

// GetByID returns a service detail by ID.
func (s *ServiceDetailService) GetByID(id uint) (*ServiceDetail, error) {
	var item ServiceDetail
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// GetByUserID returns services for a user.
func (s *ServiceDetailService) GetByUserID(userID uint, page, pageSize int) ([]ServiceDetail, int64, error) {
	var items []ServiceDetail
	var total int64

	query := s.db.Model(&ServiceDetail{}).Where("user_id = ?", userID)
	query.Count(&total)
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&items)
	return items, total, nil
}

// Create creates a new service detail.
func (s *ServiceDetailService) Create(userID, productID uint, domain, username, password, startDate, dueDate string, price float64, billingCycle, note string) (*ServiceDetail, error) {
	item := &ServiceDetail{
		UserID:       userID,
		ProductID:    productID,
		Domain:       domain,
		Username:     username,
		Password:     password,
		Price:        price,
		BillingCycle: billingCycle,
		Status:       1,
		Note:         note,
	}

	if startDate != "" {
		t, _ := time.Parse("2006-01-02", startDate)
		item.StartDate = &t
	}
	if dueDate != "" {
		t, _ := time.Parse("2006-01-02", dueDate)
		item.DueDate = &t
	}

	if err := s.db.Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

// Update updates a service detail.
func (s *ServiceDetailService) Update(id uint, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}
	result := s.db.Model(&ServiceDetail{}).Where("id = ?", id).Updates(updates)
	if result.RowsAffected == 0 {
		return errors.New("service not found")
	}
	return result.Error
}

// Delete deletes a service detail.
func (s *ServiceDetailService) Delete(id uint) error {
	result := s.db.Delete(&ServiceDetail{}, id)
	if result.RowsAffected == 0 {
		return errors.New("service not found")
	}
	return result.Error
}

// Suspend suspends a service.
func (s *ServiceDetailService) Suspend(id uint, reason string) error {
	return s.db.Model(&ServiceDetail{}).Where("id = ?", id).Update("status", 2).Error
}

// Unsuspend unsuspends a service.
func (s *ServiceDetailService) Unsuspend(id uint) error {
	return s.db.Model(&ServiceDetail{}).Where("id = ?", id).Update("status", 1).Error
}

// Terminate terminates a service.
func (s *ServiceDetailService) Terminate(id uint, reason string) error {
	return s.db.Model(&ServiceDetail{}).Where("id = ?", id).Update("status", 3).Error
}

// Renew renews a service.
func (s *ServiceDetailService) Renew(id uint, period int, periodUnit string) (map[string]interface{}, error) {
	var item ServiceDetail
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, errors.New("service not found")
	}

	var months int
	switch periodUnit {
	case "month":
		months = period
	case "year":
		months = period * 12
	default:
		months = period
	}

	var base time.Time
	if item.DueDate != nil && item.DueDate.After(time.Now()) {
		base = *item.DueDate
	} else {
		base = time.Now()
	}
	newDue := base.AddDate(0, months, 0)

	s.db.Model(&item).Updates(map[string]interface{}{
		"due_date": &newDue,
		"status":   1,
	})

	return map[string]interface{}{
		"service_id": id,
		"new_due":    newDue,
		"period":     period,
		"period_unit": periodUnit,
	}, nil
}

// GetStats returns overall service statistics.
func (s *ServiceDetailService) GetStats() (map[string]interface{}, error) {
	var total, active, suspended, terminated int64
	s.db.Model(&model.ClientService{}).Count(&total)
	s.db.Model(&model.ClientService{}).Where("status = ?", 1).Count(&active)
	s.db.Model(&model.ClientService{}).Where("status = ?", 2).Count(&suspended)
	s.db.Model(&model.ClientService{}).Where("status = ?", 3).Count(&terminated)

	return map[string]interface{}{
		"total":      total,
		"active":     active,
		"suspended":  suspended,
		"terminated": terminated,
	}, nil
}

// GetServiceLogs returns logs for a service.
func (s *ServiceDetailService) GetServiceLogs(serviceID uint, page, pageSize int) ([]map[string]interface{}, int64, error) {
	var logs []map[string]interface{}
	var total int64

	query := s.db.Table("system_logs").Where("target_id = ? AND target_type = ?", serviceID, "service")
	query.Count(&total)
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs)

	return logs, total, nil
}
