package service

import (
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
	s.db.Table("client_services").Where("user_id = ? AND status = ?", userID, "active").Count(&active)
	s.db.Table("client_services").Where("user_id = ? AND status = ?", userID, "suspended").Count(&suspended)
	s.db.Table("client_services").Where("user_id = ? AND status = ?", userID, "terminated").Count(&terminated)
	return map[string]interface{}{
		"total": total, "active": active, "suspended": suspended, "terminated": terminated,
	}, nil
}
