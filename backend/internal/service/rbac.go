package service

import (
	"errors"
	"time"

	"github.com/anchor-finance/backend/pkg/logger"

	"gorm.io/gorm"
)

// Role 角色（service 层副本）
type Role struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	Name        string       `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Description string       `gorm:"type:text" json:"description"`
	IsSystem    bool         `gorm:"default:false" json:"is_system"`
	SortOrder   int          `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}

// Permission 权限
type Permission struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Code      string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"code"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Module    string    `gorm:"type:varchar(50);index" json:"module"`
	Type      string    `gorm:"type:varchar(20)" json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

// UserRole 用户角色关联
type UserRole struct {
	UserID    uint      `gorm:"primaryKey" json:"user_id"`
	RoleID    uint      `gorm:"primaryKey" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

// RolePermission 角色权限关联
type RolePermission struct {
	RoleID       uint `gorm:"primaryKey" json:"role_id"`
	PermissionID uint `gorm:"primaryKey" json:"permission_id"`
}

type RbacService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewRbacService(db *gorm.DB, log *logger.Logger) *RbacService {
	return &RbacService{db: db, log: log}
}

type CreateRoleRequest struct {
	Name          string `json:"name" binding:"required,max=50"`
	Description   string `json:"description"`
	PermissionIDs []uint `json:"permission_ids"`
}

type UpdateRoleRequest struct {
	Name          string `json:"name" binding:"omitempty,max=50"`
	Description   string `json:"description"`
	PermissionIDs []uint `json:"permission_ids"`
}

type AssignRoleRequest struct {
	RoleIDs []uint `json:"role_ids" binding:"required"`
}

// GetRoles returns all roles with permissions.
func (s *RbacService) GetRoles() ([]Role, error) {
	var roles []Role
	if err := s.db.Preload("Permissions").Order("sort_order ASC, id ASC").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// CreateRole creates a new role with optional permissions.
func (s *RbacService) CreateRole(req CreateRoleRequest) (*Role, error) {
	role := &Role{
		Name:        req.Name,
		Description: req.Description,
		IsSystem:    false,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(role).Error; err != nil {
			return err
		}
		if len(req.PermissionIDs) > 0 {
			var perms []Permission
			if err := tx.Where("id IN ?", req.PermissionIDs).Find(&perms).Error; err != nil {
				return err
			}
			if err := tx.Model(role).Association("Permissions").Replace(perms); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	s.log.Infof("role created: %s (id=%d)", role.Name, role.ID)
	return s.GetRoleByID(role.ID)
}

// UpdateRole updates role info and permissions.
func (s *RbacService) UpdateRole(roleID uint, req UpdateRoleRequest) (*Role, error) {
	var role Role
	if err := s.db.First(&role, roleID).Error; err != nil {
		return nil, err
	}
	if role.IsSystem {
		return nil, errors.New("system role cannot be modified")
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if len(updates) > 0 {
			if err := tx.Model(&role).Updates(updates).Error; err != nil {
				return err
			}
		}
		if req.PermissionIDs != nil {
			var perms []Permission
			if err := tx.Where("id IN ?", req.PermissionIDs).Find(&perms).Error; err != nil {
				return err
			}
			if err := tx.Model(&role).Association("Permissions").Replace(perms); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	s.log.Infof("role updated: %s (id=%d)", role.Name, role.ID)
	return s.GetRoleByID(roleID)
}

// DeleteRole soft-deletes a role (system roles excluded).
func (s *RbacService) DeleteRole(roleID uint) error {
	var role Role
	if err := s.db.First(&role, roleID).Error; err != nil {
		return err
	}
	if role.IsSystem {
		return errors.New("system role cannot be deleted")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&role).Association("Permissions").Clear(); err != nil {
			return err
		}
		if err := tx.Model(&role).Association("Users").Clear(); err != nil {
			return err
		}
		return tx.Delete(&role).Error
	})
}

// GetRoleByID fetches a role by ID with permissions.
func (s *RbacService) GetRoleByID(roleID uint) (*Role, error) {
	var role Role
	if err := s.db.Preload("Permissions").First(&role, roleID).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// GetPermissions returns all permissions grouped by module.
func (s *RbacService) GetPermissions() (map[string][]Permission, error) {
	var perms []Permission
	if err := s.db.Order("module ASC, code ASC").Find(&perms).Error; err != nil {
		return nil, err
	}

	grouped := make(map[string][]Permission)
	for _, p := range perms {
		grouped[p.Module] = append(grouped[p.Module], p)
	}
	return grouped, nil
}

// AssignRole assigns roles to a user (replaces existing).
func (s *RbacService) AssignRole(userID uint, roleIDs []uint) error {
	var user User
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	var roles []Role
	if err := s.db.Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
		return err
	}
	if len(roles) != len(roleIDs) {
		return errors.New("some roles not found")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// 清除旧关联
		if err := tx.Where("user_id = ?", userID).Delete(&UserRole{}).Error; err != nil {
			return err
		}
		// 写入新关联
		for _, rid := range roleIDs {
			ur := UserRole{UserID: userID, RoleID: rid, CreatedAt: time.Now()}
			if err := tx.Create(&ur).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetUserRoles returns roles assigned to a user.
func (s *RbacService) GetUserRoles(userID uint) ([]Role, error) {
	var userRoles []UserRole
	if err := s.db.Where("user_id = ?", userID).Find(&userRoles).Error; err != nil {
		return nil, err
	}

	if len(userRoles) == 0 {
		return []Role{}, nil
	}

	roleIDs := make([]uint, len(userRoles))
	for i, ur := range userRoles {
		roleIDs[i] = ur.RoleID
	}

	var roles []Role
	if err := s.db.Preload("Permissions").Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// HasPermission checks if a user has a specific permission (by code).
func (s *RbacService) HasPermission(userID uint, permCode string) bool {
	roles, err := s.GetUserRoles(userID)
	if err != nil {
		return false
	}

	for _, role := range roles {
		for _, perm := range role.Permissions {
			if perm.Code == permCode {
				return true
			}
		}
	}
	return false
}

// GetUserPermissions returns all unique permissions for a user across all roles.
func (s *RbacService) GetUserPermissions(userID uint) ([]Permission, error) {
	roles, err := s.GetUserRoles(userID)
	if err != nil {
		return nil, err
	}

	seen := make(map[uint]bool)
	var perms []Permission
	for _, role := range roles {
		for _, p := range role.Permissions {
			if !seen[p.ID] {
				seen[p.ID] = true
				perms = append(perms, p)
			}
		}
	}
	return perms, nil
}
