package admin

import (
	"fmt"
	"net/http"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetSettings 获取所有设置
// GET /api/admin/settings
func GetSettings(c *gin.Context) {
	db := database.GetDB()
	var settings []model.Setting
	db.Find(&settings)

	// 转换为key-value格式
	result := make(map[string]interface{})
	for _, s := range settings {
		result[s.Key] = s.Value
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    result,
	})
}

// GetSettingsByGroup 获取分组设置
// GET /api/admin/settings/:group
func GetSettingsByGroup(c *gin.Context) {
	group := c.Param("group")

	db := database.GetDB()
	var settings []model.Setting
	db.Where("group = ?", group).Find(&settings)

	// 转换为key-value格式
	result := make(map[string]interface{})
	for _, s := range settings {
		result[s.Key] = s.Value
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    result,
	})
}

// UpdateSettings 更新设置
// PUT /api/admin/settings
func UpdateSettings(c *gin.Context) {
	// 解析请求参数
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	db := database.GetDB()

	// 逐个更新设置
	for key, value := range req {
		// 转换value为字符串
		var strValue string
		switch v := value.(type) {
		case string:
			strValue = v
		default:
			strValue = fmt.Sprintf("%v", v)
		}

		// 查找或创建设置
		var setting model.Setting
		result := db.Where("key = ?", key).First(&setting)
		if result.Error != nil {
			// 不存在则创建
			setting = model.Setting{
				Key:   key,
				Value: strValue,
				Group: "general",
			}
			db.Create(&setting)
		} else {
			// 存在则更新
			db.Model(&setting).Update("value", strValue)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
		"data":    nil,
	})
}

// GetMenus 获取菜单列表
// GET /api/admin/menus
func GetMenus(c *gin.Context) {
	db := database.GetDB()
	var menus []model.Menu
	db.Where("is_visible = ?", true).Order("sort_order ASC").Find(&menus)

	// 构建树形结构
	tree := buildMenuTree(menus, 0)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    tree,
	})
}

// buildMenuTree 构建菜单树
func buildMenuTree(menus []model.Menu, parentID uint) []gin.H {
	var result []gin.H
	for _, menu := range menus {
		if menu.ParentID == parentID {
			children := buildMenuTree(menus, menu.ID)
			item := gin.H{
				"id":         menu.ID,
				"name":       menu.Name,
				"path":       menu.Path,
				"icon":       menu.Icon,
				"sort_order": menu.SortOrder,
			}
			if len(children) > 0 {
				item["children"] = children
			}
			result = append(result, item)
		}
	}
	return result
}

// CreateMenu 创建菜单
// POST /api/admin/menus
func CreateMenu(c *gin.Context) {
	var req struct {
		ParentID  uint   `json:"parent_id"`
		Name      string `json:"name" binding:"required"`
		Path      string `json:"path"`
		Icon      string `json:"icon"`
		SortOrder int    `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	db := database.GetDB()
	menu := model.Menu{
		ParentID:  req.ParentID,
		Name:      req.Name,
		Path:      req.Path,
		Icon:      req.Icon,
		SortOrder: req.SortOrder,
		IsVisible: true,
	}

	if err := db.Create(&menu).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建菜单失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": menu.ID,
		},
	})
}

// UpdateMenu 更新菜单
// PUT /api/admin/menus/:id
func UpdateMenu(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name      string `json:"name"`
		Path      string `json:"path"`
		Icon      string `json:"icon"`
		SortOrder int    `json:"sort_order"`
		IsVisible *bool  `json:"is_visible"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	db := database.GetDB()
	var menu model.Menu
	if err := db.First(&menu, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "菜单不存在",
			"data":    nil,
		})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Path != "" {
		updates["path"] = req.Path
	}
	if req.Icon != "" {
		updates["icon"] = req.Icon
	}
	if req.SortOrder > 0 {
		updates["sort_order"] = req.SortOrder
	}
	if req.IsVisible != nil {
		updates["is_visible"] = *req.IsVisible
	}

	db.Model(&menu).Updates(updates)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
		"data":    nil,
	})
}

// GetMenuTypeList 获取菜单类型列表
// GET /api/admin/menu-types
func GetMenuTypeList(c *gin.Context) {
	types := []gin.H{
		{"id": "system", "name": "系统菜单", "description": "系统内置菜单"},
		{"id": "custom", "name": "自定义菜单", "description": "用户自定义菜单"},
		{"id": "product", "name": "产品菜单", "description": "产品相关菜单"},
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    types,
	})
}

// DeleteMenu 删除菜单
// DELETE /api/admin/menus/:id
func DeleteMenu(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var menu model.Menu
	if err := db.First(&menu, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "菜单不存在",
			"data":    nil,
		})
		return
	}

	db.Delete(&menu)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
		"data":    nil,
	})
}
