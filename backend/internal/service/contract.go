package service

import (
	"errors"
	"fmt"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/internal/util"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// ContractService 合同管理服务
type ContractService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewContractService(db *gorm.DB, log *logger.Logger) *ContractService {
	return &ContractService{db: db, log: log}
}

// CreateContractRequest 创建合同请求
type CreateContractRequest struct {
	UserID    uint    `json:"user_id" binding:"required"`
	Title     string  `json:"title" binding:"required,max=255"`
	Content   string  `json:"content"`
	Type      string  `json:"type" binding:"omitempty,oneof=service sales custom"`
	Amount    float64 `json:"amount"`
	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`
	AdminID   uint    `json:"admin_id"`
}

// generateContractNo 生成合同编号
func (s *ContractService) generateContractNo() string {
	return fmt.Sprintf("CON%s%04d", time.Now().Format("20060102150405"), time.Now().UnixNano()%10000)
}

// Create 创建合同
func (s *ContractService) Create(req CreateContractRequest) (*model.Contract, error) {
	contract := model.Contract{
		ContractNo: s.generateContractNo(),
		UserID:     req.UserID,
		Title:      req.Title,
		Content:    req.Content,
		Type:       req.Type,
		Status:     1,
		Amount:     req.Amount,
		AdminID:    req.AdminID,
	}
	if contract.Type == "" {
		contract.Type = "service"
	}

	if req.StartDate != nil && *req.StartDate != "" {
		t, err := util.ParseTime(*req.StartDate)
		if err != nil {
			return nil, errors.New("invalid start_date")
		}
		contract.StartDate = &t
	}
	if req.EndDate != nil && *req.EndDate != "" {
		t, err := util.ParseTime(*req.EndDate)
		if err != nil {
			return nil, errors.New("invalid end_date")
		}
		contract.EndDate = &t
	}

	if err := s.db.Create(&contract).Error; err != nil {
		return nil, err
	}
	s.log.Infof("contract created: id=%d no=%s", contract.ID, contract.ContractNo)
	return &contract, nil
}

// Get 获取合同
func (s *ContractService) Get(id uint) (*model.Contract, error) {
	var contract model.Contract
	if err := s.db.Preload("User").First(&contract, id).Error; err != nil {
		return nil, err
	}
	return &contract, nil
}

// List 获取合同列表
func (s *ContractService) List(page, pageSize int, userID uint, status int) ([]model.Contract, int64, error) {
	var items []model.Contract
	var total int64

	query := s.db.Model(&model.Contract{})
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if status > 0 {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	offset, limit := util.Paginate(page, pageSize)
	query.Offset(offset).Limit(limit).Order("id DESC").Preload("User").Find(&items)
	return items, total, nil
}

// Update 更新合同
func (s *ContractService) Update(id uint, updates map[string]interface{}) (*model.Contract, error) {
	var contract model.Contract
	if err := s.db.First(&contract, id).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&contract).Updates(updates).Error; err != nil {
		return nil, err
	}
	s.db.First(&contract, id)
	return &contract, nil
}

// Delete 删除合同
func (s *ContractService) Delete(id uint) error {
	result := s.db.Delete(&model.Contract{}, id)
	if result.RowsAffected == 0 {
		return errors.New("contract not found")
	}
	return result.Error
}

// Sign 签署合同
func (s *ContractService) Sign(id uint) (*model.Contract, error) {
	var contract model.Contract
	if err := s.db.First(&contract, id).Error; err != nil {
		return nil, err
	}
	if contract.Status != 2 {
		return nil, errors.New("contract is not in pending sign status")
	}
	now := time.Now()
	updates := map[string]interface{}{
		"status":     3,
		"signed_at":  &now,
	}
	if err := s.db.Model(&contract).Updates(updates).Error; err != nil {
		return nil, err
	}
	s.db.First(&contract, id)
	s.log.Infof("contract signed: id=%d", id)
	return &contract, nil
}

// Expire 过期合同
func (s *ContractService) Expire(id uint) (*model.Contract, error) {
	var contract model.Contract
	if err := s.db.First(&contract, id).Error; err != nil {
		return nil, err
	}
	if contract.Status != 3 {
		return nil, errors.New("only signed contracts can be expired")
	}
	if err := s.db.Model(&contract).Update("status", 4).Error; err != nil {
		return nil, err
	}
	s.db.First(&contract, id)
	s.log.Infof("contract expired: id=%d", id)
	return &contract, nil
}
