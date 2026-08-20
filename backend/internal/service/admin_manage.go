package service

import (
	"errors"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AdminManageService provides admin management operations.
type AdminManageService struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewAdminManageService creates a new AdminManageService.
func NewAdminManageService(db *gorm.DB, log *logger.Logger) *AdminManageService {
	return &AdminManageService{db: db, log: log}
}

// CreateAdminRequest is the payload for creating an admin.
type CreateAdminRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=6,max=128"`
	Email    string `json:"email" binding:"omitempty,email"`
	RealName string `json:"real_name" binding:"omitempty,max=64"`
	RoleID   uint   `json:"role_id"`
	Avatar   string `json:"avatar"`
}

// UpdateAdminRequest is the payload for updating an admin.
type UpdateAdminRequest struct {
	Email    *string `json:"email"`
	RealName *string `json:"real_name"`
	RoleID   *uint   `json:"role_id"`
	Avatar   *string `json:"avatar"`
}

// Create creates a new admin user.
func (s *AdminManageService) Create(req CreateAdminRequest) (*model.Admin, error) {
	// Check username uniqueness
	var count int64
	s.db.Model(&model.Admin{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		return nil, errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	admin := &model.Admin{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		RealName: req.RealName,
		RoleID:   req.RoleID,
		Status:   1,
		Avatar:   req.Avatar,
	}

	if err := s.db.Create(admin).Error; err != nil {
		return nil, err
	}

	s.log.Infof("admin created: username=%s", req.Username)
	return admin, nil
}

// GetByID fetches an admin by ID.
func (s *AdminManageService) GetByID(id uint) (*model.Admin, error) {
	var admin model.Admin
	if err := s.db.First(&admin, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("admin not found")
		}
		return nil, err
	}
	return &admin, nil
}

// GetList returns a paginated admin list.
func (s *AdminManageService) GetList(page, pageSize int, keyword string) ([]model.Admin, int64, error) {
	var admins []model.Admin
	var total int64

	query := s.db.Model(&model.Admin{})
	if keyword != "" {
		q := "%" + keyword + "%"
		query = query.Where("username LIKE ? OR email LIKE ? OR real_name LIKE ?", q, q, q)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("id ASC").Find(&admins).Error; err != nil {
		return nil, 0, err
	}
	return admins, total, nil
}

// Update modifies an existing admin.
func (s *AdminManageService) Update(id uint, req UpdateAdminRequest) error {
	admin, err := s.GetByID(id)
	if err != nil {
		return err
	}

	updates := map[string]interface{}{}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.RealName != nil {
		updates["real_name"] = *req.RealName
	}
	if req.RoleID != nil {
		updates["role_id"] = *req.RoleID
	}
	if req.Avatar != nil {
		updates["avatar"] = *req.Avatar
	}

	if len(updates) == 0 {
		return nil
	}

	return s.db.Model(admin).Updates(updates).Error
}

// Delete removes an admin by ID.
func (s *AdminManageService) Delete(id uint) error {
	result := s.db.Delete(&model.Admin{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("admin not found")
	}
	return nil
}

// SetStatus enables or disables an admin.
func (s *AdminManageService) SetStatus(id uint, status int) error {
	if status != 0 && status != 1 {
		return errors.New("invalid status value")
	}
	result := s.db.Model(&model.Admin{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("admin not found")
	}
	return nil
}

// ResetPassword sets a new password for an admin.
func (s *AdminManageService) ResetPassword(id uint, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	result := s.db.Model(&model.Admin{}).Where("id = ?", id).Update("password", string(hashedPassword))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("admin not found")
	}
	return nil
}

// UpdateLastLogin updates the last login time for an admin.
func (s *AdminManageService) UpdateLastLogin(id uint) error {
	now := time.Now()
	return s.db.Model(&model.Admin{}).Where("id = ?", id).Update("last_login", now).Error
}

// GetOperationLogs returns operation logs with optional admin filter.
func (s *AdminManageService) GetOperationLogs(page, pageSize int, adminID uint, module string) ([]model.AdminLog, int64, error) {
	var logs []model.AdminLog
	var total int64

	query := s.db.Model(&model.AdminLog{})
	if adminID > 0 {
		query = query.Where("admin_id = ?", adminID)
	}
	if module != "" {
		query = query.Where("module = ?", module)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	if err := query.Preload("Admin").Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

// CreateOperationLog creates an operation log entry.
func (s *AdminManageService) CreateOperationLog(log *model.AdminLog) error {
	return s.db.Create(log).Error
}
