package service

import (
	"encoding/json"
	"errors"
	"math"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type MenuService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewMenuService(db *gorm.DB, log *logger.Logger) *MenuService {
	return &MenuService{db: db, log: log}
}

type CreateMenuRequest struct {
	Name       string             `json:"name" binding:"required,max=64"`
	Icon       string             `json:"icon"`
	URL        string             `json:"url"`
	ParentID   *uint              `json:"parent_id"`
	SortOrder  int                `json:"sort_order"`
	Type       string             `json:"type" binding:"required,oneof=admin top side bottom"`
	IsVisible  *bool              `json:"is_visible"`
	Permission string             `json:"permission"`
	Target     string             `json:"target" binding:"omitempty,oneof=_self _blank"`
	Badge      string             `json:"badge"`
	BadgeType  string             `json:"badge_type" binding:"omitempty,oneof=dot number text"`
	Extra      map[string]interface{} `json:"extra"`
}

type UpdateMenuRequest struct {
	Name       string             `json:"name" binding:"omitempty,max=64"`
	Icon       string             `json:"icon"`
	URL        string             `json:"url"`
	ParentID   *uint              `json:"parent_id"`
	SortOrder  *int               `json:"sort_order"`
	Type       string             `json:"type" binding:"omitempty,oneof=admin top side bottom"`
	IsVisible  *bool              `json:"is_visible"`
	Permission string             `json:"permission"`
	Target     string             `json:"target" binding:"omitempty,oneof=_self _blank"`
	Badge      string             `json:"badge"`
	BadgeType  string             `json:"badge_type" binding:"omitempty,oneof=dot number text"`
	Extra      map[string]interface{} `json:"extra"`
	IsActive   *bool              `json:"is_active"`
}

type SortMenuRequest struct {
	IDs []uint `json:"ids" binding:"required"`
}

// List returns all menus (flat list).
func (s *MenuService) List(menuType string) ([]model.Menu, error) {
	var menus []model.Menu
	query := s.db.Order("sort_order ASC, id ASC")
	if menuType != "" {
		query = query.Where("type = ?", menuType)
	}
	if err := query.Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

// GetTree returns menus in tree structure.
func (s *MenuService) GetTree(menuType string) ([]model.Menu, error) {
	var menus []model.Menu
	query := s.db.Order("sort_order ASC, id ASC")
	if menuType != "" {
		query = query.Where("type = ?", menuType)
	}
	if err := query.Find(&menus).Error; err != nil {
		return nil, err
	}
	s.log.Infof("GetTree: loaded %d menus", len(menus))
	if len(menus) > 0 {
		s.log.Infof("GetTree: first menu id=%d parent_id=%v name=%s", menus[0].ID, menus[0].ParentID, menus[0].Name)
	}
	tree := buildMenuTree(menus, nil)
	s.log.Infof("GetTree: tree has %d top-level items", len(tree))
	return tree, nil
}

// GetVisibleTree returns only visible menus in tree structure.
func (s *MenuService) GetVisibleTree(menuType string) ([]model.Menu, error) {
	var menus []model.Menu
	query := s.db.Where("is_visible = ? AND is_active = ?", true, true).Order("sort_order ASC, id ASC")
	if menuType != "" {
		query = query.Where("type = ?", menuType)
	}
	if err := query.Find(&menus).Error; err != nil {
		return nil, err
	}
	return buildMenuTree(menus, nil), nil
}

// GetByID returns a single menu by ID.
func (s *MenuService) GetByID(id uint) (*model.Menu, error) {
	var menu model.Menu
	if err := s.db.First(&menu, id).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

// Create creates a new menu item.
func (s *MenuService) Create(req CreateMenuRequest) (*model.Menu, error) {
	if req.ParentID != nil && *req.ParentID > 0 {
		var parent model.Menu
		if err := s.db.First(&parent, *req.ParentID).Error; err != nil {
			return nil, errors.New("parent menu not found")
		}
		if parent.Type != req.Type {
			return nil, errors.New("child menu type must match parent")
		}
	}

	isVisible := true
	if req.IsVisible != nil {
		isVisible = *req.IsVisible
	}
	target := req.Target
	if target == "" {
		target = "_self"
	}

	menu := &model.Menu{
		Name:       req.Name,
		Icon:       req.Icon,
		URL:        req.URL,
		ParentID:   req.ParentID,
		SortOrder:  req.SortOrder,
		Type:       req.Type,
		IsVisible:  isVisible,
		Permission: req.Permission,
		Target:     target,
		Badge:      req.Badge,
		BadgeType:  req.BadgeType,
		Extra:      toDatatypesJSON(req.Extra),
		IsActive:   true,
	}

	if err := s.db.Create(menu).Error; err != nil {
		return nil, err
	}

	s.log.Infof("menu created: %s (id=%d)", menu.Name, menu.ID)
	return menu, nil
}

// Update updates a menu item.
func (s *MenuService) Update(id uint, req UpdateMenuRequest) (*model.Menu, error) {
	var menu model.Menu
	if err := s.db.First(&menu, id).Error; err != nil {
		return nil, err
	}

	// Prevent circular reference
	if req.ParentID != nil && *req.ParentID == id {
		return nil, errors.New("menu cannot be its own parent")
	}
	if req.ParentID != nil && *req.ParentID > 0 {
		var parent model.Menu
		if err := s.db.First(&parent, *req.ParentID).Error; err != nil {
			return nil, errors.New("parent menu not found")
		}
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	updates["icon"] = req.Icon
	updates["url"] = req.URL
	if req.ParentID != nil {
		updates["parent_id"] = req.ParentID
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.IsVisible != nil {
		updates["is_visible"] = *req.IsVisible
	}
	updates["permission"] = req.Permission
	if req.Target != "" {
		updates["target"] = req.Target
	}
	updates["badge"] = req.Badge
	if req.BadgeType != "" {
		updates["badge_type"] = req.BadgeType
	}
	if req.Extra != nil {
		updates["extra"] = toDatatypesJSON(req.Extra)
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(updates) > 0 {
		if err := s.db.Model(&menu).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	s.log.Infof("menu updated: %s (id=%d)", menu.Name, menu.ID)
	return s.GetByID(id)
}

// Delete deletes a menu and its children.
func (s *MenuService) Delete(id uint) error {
	var menu model.Menu
	if err := s.db.First(&menu, id).Error; err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Delete children first
		if err := tx.Where("parent_id = ?", id).Delete(&model.Menu{}).Error; err != nil {
			return err
		}
		return tx.Delete(&menu).Error
	})
}

// Sort updates sort order for multiple menus.
func (s *MenuService) Sort(req SortMenuRequest) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for i, id := range req.IDs {
			if err := tx.Model(&model.Menu{}).Where("id = ?", id).Update("sort_order", i).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// buildMenuTree recursively builds a tree from a flat list.
func buildMenuTree(menus []model.Menu, parentID *uint) []model.Menu {
	var tree []model.Menu
	for _, menu := range menus {
		menuParentID := uint(0)
		if menu.ParentID != nil {
			menuParentID = *menu.ParentID
		}

		isRoot := parentID == nil && menuParentID == 0
		isChild := parentID != nil && menuParentID == *parentID
		if isRoot || isChild {
			children := buildMenuTree(menus, &menu.ID)
			menu.Children = children
			tree = append(tree, menu)
		}
	}
	if parentID == nil {
		logger.Default().Infof("buildMenuTree(root): scanned %d items, found %d roots", len(menus), len(tree))
	}
	return tree
}

func toDatatypesJSON(v map[string]interface{}) datatypes.JSON {
	if v == nil {
		return nil
	}
	data, _ := json.Marshal(v)
	return datatypes.JSON(data)
}

func menuPaginate(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	pageSize = int(math.Min(float64(pageSize), 100))
	return (page - 1) * pageSize, pageSize
}
