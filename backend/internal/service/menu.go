package service

import (
	"errors"

	"anchorfinance/internal/model"
	"anchorfinance/internal/repository"
)

type MenuService struct {
	menuRepo *repository.MenuRepository
}

func NewMenuService() *MenuService {
	return &MenuService{
		menuRepo: repository.NewMenuRepository(),
	}
}

// CreateMenuRequest 创建菜单请求
type CreateMenuRequest struct {
	Name      string `json:"name" binding:"required"`
	ParentID  uint   `json:"parent_id"`
	URL       string `json:"url"`
	Icon      string `json:"icon"`
	SortOrder int    `json:"sort_order"`
	IsVisible bool   `json:"is_visible"`
	IsSystem  bool   `json:"is_system"`
}

// UpdateMenuRequest 更新菜单请求
type UpdateMenuRequest struct {
	Name      string `json:"name"`
	ParentID  uint   `json:"parent_id"`
	URL       string `json:"url"`
	Icon      string `json:"icon"`
	SortOrder int    `json:"sort_order"`
	IsVisible *bool  `json:"is_visible"`
}

// MenuTreeNode 菜单树节点
type MenuTreeNode struct {
	ID        uint            `json:"id"`
	Name      string          `json:"name"`
	URL       string          `json:"url,omitempty"`
	Icon      string          `json:"icon,omitempty"`
	ParentID  uint            `json:"parent_id"`
	IsVisible bool            `json:"is_visible"`
	Children  []MenuTreeNode  `json:"children,omitempty"`
}

// GetMenuTree 获取菜单树
func (s *MenuService) GetMenuTree(adminID uint) ([]MenuTreeNode, error) {
	// 获取所有激活的菜单
	menus, err := s.menuRepo.FindActive()
	if err != nil {
		return nil, err
	}

	// 构建菜单树
	tree := s.buildTree(menus, 0)
	return tree, nil
}

// buildTree 构建菜单树
func (s *MenuService) buildTree(menus []model.Menu, parentID uint) []MenuTreeNode {
	var nodes []MenuTreeNode
	for _, menu := range menus {
		if menu.ParentID == parentID {
			node := MenuTreeNode{
				ID:        menu.ID,
				Name:      menu.Name,
				URL:       menu.URL,
				Icon:      menu.Icon,
				ParentID:  menu.ParentID,
				IsVisible: menu.IsVisible,
				Children:  s.buildTree(menus, menu.ID),
			}
			nodes = append(nodes, node)
		}
	}
	return nodes
}

// GetAllMenus 获取所有菜单（扁平）
func (s *MenuService) GetAllMenus() ([]model.Menu, error) {
	return s.menuRepo.FindAll()
}

// CreateMenu 创建菜单
func (s *MenuService) CreateMenu(req *CreateMenuRequest) (*model.Menu, error) {
	menu := &model.Menu{
		Name:      req.Name,
		ParentID:  req.ParentID,
		URL:       req.URL,
		Icon:      req.Icon,
		SortOrder: req.SortOrder,
		IsVisible: req.IsVisible,
		IsSystem:  req.IsSystem,
		IsActive:  true,
	}

	if err := s.menuRepo.Create(menu); err != nil {
		return nil, err
	}

	return menu, nil
}

// UpdateMenu 更新菜单
func (s *MenuService) UpdateMenu(id uint, req *UpdateMenuRequest) (*model.Menu, error) {
	menu, err := s.menuRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("菜单不存在")
	}

	if req.Name != "" {
		menu.Name = req.Name
	}
	if req.URL != "" {
		menu.URL = req.URL
	}
	if req.Icon != "" {
		menu.Icon = req.Icon
	}
	if req.SortOrder != 0 {
		menu.SortOrder = req.SortOrder
	}
	if req.IsVisible != nil {
		menu.IsVisible = *req.IsVisible
	}
	if req.ParentID != 0 {
		// 检查不能将菜单移到自己的子菜单下
		if req.ParentID == menu.ID {
			return nil, errors.New("不能将菜单移到自己下面")
		}
		menu.ParentID = req.ParentID
	}

	if err := s.menuRepo.Update(menu); err != nil {
		return nil, err
	}

	return menu, nil
}

// DeleteMenu 删除菜单
func (s *MenuService) DeleteMenu(id uint) error {
	menu, err := s.menuRepo.FindByID(id)
	if err != nil {
		return errors.New("菜单不存在")
	}

	if menu.IsSystem {
		return errors.New("系统菜单不能删除")
	}

	// 检查是否有子菜单
	children, err := s.menuRepo.FindByParentID(id)
	if err != nil {
		return err
	}
	if len(children) > 0 {
		return errors.New("请先删除子菜单")
	}

	return s.menuRepo.Delete(id)
}
