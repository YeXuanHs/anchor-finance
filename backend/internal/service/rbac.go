package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"anchorfinance/pkg/logger"

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

// ─── RBAC Admin Methods (from zjmf RbacController) ───

// AdminGetRoles returns all roles with admin users.
func (s *RbacService) AdminGetRoles(order, sort string) ([]map[string]interface{}, error) {
	if order == "" {
		order = "a.id"
	}
	if sort == "" {
		sort = "DESC"
	}

	var roles []map[string]interface{}
	if err := s.db.Table("role a").
		Select("a.id, a.name, a.status, a.remark").
		Order(order + " " + sort).
		Find(&roles).Error; err != nil {
		return nil, err
	}

	// Attach admin users to each role
	for i, role := range roles {
		roleID, _ := role["id"].(uint)
		var users []string
		s.db.Table("role_user b").
			Select("c.user_login").
			Joins("LEFT JOIN user c ON c.id = b.user_id").
			Where("b.role_id = ?", roleID).
			Pluck("c.user_login", &users)
		role["user_login"] = users
		roles[i] = role
	}

	return roles, nil
}

// AdminGetAuthTree returns auth rules as a tree.
func (s *RbacService) AdminGetAuthTree() ([]map[string]interface{}, error) {
	var rules []struct {
		ID    uint   `json:"id"`
		PID   uint   `json:"pid"`
		Title string `json:"title"`
	}
	if err := s.db.Table("auth_rule").Where("status = 1").Select("id, pid, title").Find(&rules).Error; err != nil {
		return nil, err
	}

	// Build tree
	ruleMap := make(map[uint]*map[string]interface{})
	var result []map[string]interface{}
	for _, r := range rules {
		node := map[string]interface{}{
			"id":    r.ID,
			"pid":   r.PID,
			"title": r.Title,
		}
		nodePtr := &node
		ruleMap[r.ID] = nodePtr
		if r.PID == 0 {
			result = append(result, node)
		}
	}
	for _, r := range rules {
		if r.PID != 0 {
			if parent, ok := ruleMap[r.PID]; ok {
				children, _ := (*parent)["sublevel"].([]map[string]interface{})
				children = append(children, *ruleMap[r.ID])
				(*parent)["sublevel"] = children
			}
		}
	}

	return result, nil
}

// AdminAddRole creates a role with auth rules.
func (s *RbacService) AdminAddRole(name, remark string, status int, authIDs []uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Exec("INSERT INTO role (name, remark, status, auth_role) VALUES (?, ?, ?, ?)",
			name, remark, status, joinUint(authIDs))
		if res.Error != nil {
			return res.Error
		}
		var lid struct{ ID uint }
		tx.Raw("SELECT LAST_INSERT_ID() as id").Scan(&lid)

		// Insert auth_access
		if len(authIDs) > 0 {
			var rules []struct {
				ID   uint
				Name string
			}
			tx.Table("auth_rule").Where("id IN ?", authIDs).Select("id, name").Find(&rules)
			for _, r := range rules {
				tx.Exec("INSERT INTO auth_access (role_id, rule_name, rule_id, type) VALUES (?, ?, ?, 'admin_url')",
					lid.ID, r.Name, r.ID)
			}
		}
		return nil
	})
}

// AdminEditRolePage returns data for editing a role.
func (s *RbacService) AdminEditRolePage(roleID uint) (map[string]interface{}, error) {
	if roleID == 1 {
		return nil, fmt.Errorf("不允许的操作")
	}

	var role map[string]interface{}
	if err := s.db.Table("role").Where("id = ?", roleID).Find(&role).Error; err != nil {
		return nil, err
	}
	if role == nil {
		return nil, fmt.Errorf("角色不存在")
	}

	// Get auth rules assigned to this role
	var authRole []struct{ RuleID uint }
	s.db.Table("auth_access").Where("role_id = ?", roleID).Select("rule_id").Find(&authRole)
	authSelect := make([]uint, len(authRole))
	for i, a := range authRole {
		authSelect[i] = a.RuleID
	}

	// Get all auth rules
	var rules []struct {
		ID         uint   `json:"id"`
		PID        uint   `json:"pid"`
		IsDisplay  int    `json:"is_display"`
		Name       string `json:"name"`
		Title      string `json:"title"`
	}
	s.db.Table("auth_rule").Select("id, pid, is_display, name, title").Find(&rules)

	// Build tree
	auths := buildAuthTree(rules)

	// Get users in this role
	var users []map[string]interface{}
	s.db.Table("role_user a").
		Select("b.id, b.user_login, b.user_nickname").
		Joins("LEFT JOIN user b ON a.user_id = b.id").
		Where("a.role_id = ?", roleID).
		Find(&users)

	return map[string]interface{}{
		"role":        role,
		"auths":       auths,
		"auth_select": authSelect,
		"user":        users,
	}, nil
}

// AdminEditRole updates a role with auth rules.
func (s *RbacService) AdminEditRole(roleID uint, name, remark string, status int, authIDs []uint) error {
	if roleID == 1 {
		return fmt.Errorf("不允许的操作")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		tx.Exec("UPDATE role SET name = ?, remark = ?, status = ?, auth_role = ?, update_time = ? WHERE id = ?",
			name, remark, status, joinUint(authIDs), time.Now().Unix(), roleID)

		// Replace auth_access
		tx.Exec("DELETE FROM auth_access WHERE role_id = ?", roleID)
		if len(authIDs) > 0 {
			var rules []struct {
				ID   uint
				Name string
			}
			tx.Table("auth_rule").Where("id IN ?", authIDs).Select("id, name").Find(&rules)
			for _, r := range rules {
				tx.Exec("INSERT INTO auth_access (role_id, rule_name, rule_id, type) VALUES (?, ?, ?, 'admin_url')",
					roleID, r.Name, r.ID)
			}
		}
		return nil
	})
}

// AdminDeleteRole deletes a role (system role 1 excluded).
func (s *RbacService) AdminDeleteRole(roleID uint) error {
	if roleID == 1 {
		return fmt.Errorf("不允许删除系统角色")
	}
	var count int64
	s.db.Table("RoleUser").Where("role_id = ?", roleID).Count(&count)
	if count > 0 {
		return fmt.Errorf("该角色下存在管理员，不可删除")
	}
	return s.db.Delete(&Role{}, roleID).Error
}

// AdminCopyRole duplicates a role.
func (s *RbacService) AdminCopyRole(roleID uint, newName, newRemark string) error {
	var role Role
	if err := s.db.First(&role, roleID).Error; err != nil {
		return fmt.Errorf("要复制的分组不存在")
	}
	// Check name uniqueness
	var count int64
	s.db.Model(&Role{}).Where("name = ?", newName).Count(&count)
	if count > 0 {
		return fmt.Errorf("分组名称已存在")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		newRole := Role{
			Name:        newName,
			Description: newRemark,
		}
		if err := tx.Create(&newRole).Error; err != nil {
			return err
		}

		// Copy auth access
		var access []struct {
			RuleName string
			RuleID   uint
		}
		tx.Table("auth_access").Where("role_id = ?", roleID).Select("rule_name, rule_id").Find(&access)
		for _, a := range access {
			tx.Exec("INSERT INTO auth_access (role_id, rule_name, rule_id, type) VALUES (?, ?, ?, 'admin_url')",
				newRole.ID, a.RuleName, a.RuleID)
		}
		return nil
	})
}

// Helper functions

func joinUint(ids []uint) string {
	if len(ids) == 0 {
		return ""
	}
	parts := make([]string, len(ids))
	for i, id := range ids {
		parts[i] = strconv.FormatUint(uint64(id), 10)
	}
	return strings.Join(parts, ",")
}

func buildAuthTree(rules []struct {
	ID        uint   `json:"id"`
	PID       uint   `json:"pid"`
	IsDisplay int    `json:"is_display"`
	Name      string `json:"name"`
	Title     string `json:"title"`
}) []map[string]interface{} {
	nodeMap := make(map[uint]*map[string]interface{})
	var result []map[string]interface{}

	for _, r := range rules {
		node := map[string]interface{}{
			"id":          r.ID,
			"pid":         r.PID,
			"is_display":  r.IsDisplay,
			"name":        r.Name,
			"title":       r.Title,
		}
		nodeMap[r.ID] = &node
		if r.PID == 0 {
			result = append(result, node)
		}
	}

	for _, r := range rules {
		if r.PID != 0 {
			if parent, ok := nodeMap[r.PID]; ok {
				children, _ := (*parent)["sublevel"].([]map[string]interface{})
				children = append(children, *nodeMap[r.ID])
				(*parent)["sublevel"] = children
			}
		}
	}
	return result
}
