package admin

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/YeXuanHs/anchor-finance/internal/service"
	"github.com/gin-gonic/gin"
)

// GetAdminList 获取管理员列表
// GET /api/admin/admins
func GetAdminList(c *gin.Context) {
	// 1. 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 2. 查询管理员
	db := database.GetDB()
	var total int64
	db.Model(&model.Admin{}).Count(&total)

	var admins []model.Admin
	offset := (page - 1) * pageSize
	db.Offset(offset).Limit(pageSize).Order("id DESC").Find(&admins)

	// 3. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      admins,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetRoleList 获取角色列表
// GET /api/admin/roles
func GetRoleList(c *gin.Context) {
	db := database.GetDB()
	var roles []model.Role
	db.Find(&roles)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    roles,
	})
}

// CreateAdmin 创建管理员
// POST /api/admin/admins
func CreateAdmin(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		RealName string `json:"real_name"`
		RoleID   uint   `json:"role_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	// 检查用户名是否已存在
	db := database.GetDB()
	var count int64
	db.Model(&model.Admin{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "用户名已存在",
			"data":    nil,
		})
		return
	}

	// 创建管理员
	hashedPassword, err := service.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "密码加密失败",
			"data":    nil,
		})
		return
	}

	admin := model.Admin{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		RealName:     req.RealName,
		RoleID:       req.RoleID,
		Status:       "active",
	}

	if err := db.Create(&admin).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建管理员失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": admin.ID,
		},
	})
}

// UpdateAdmin 更新管理员
// PUT /api/admin/admins/:id
func UpdateAdmin(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Email    string `json:"email"`
		RealName string `json:"real_name"`
		RoleID   uint   `json:"role_id"`
		Status   string `json:"status"`
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
	var admin model.Admin
	if err := db.First(&admin, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "管理员不存在",
			"data":    nil,
		})
		return
	}

	updates := map[string]interface{}{}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}
	if req.RoleID > 0 {
		updates["role_id"] = req.RoleID
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	db.Model(&admin).Updates(updates)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
		"data":    nil,
	})
}

// GetCronTasks 获取定时任务列表
// GET /api/admin/cron-tasks
func GetCronTasks(c *gin.Context) {
	db := database.GetDB()
	var tasks []model.ScheduleTask
	db.Order("id ASC").Find(&tasks)

	if tasks == nil {
		tasks = []model.ScheduleTask{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      tasks,
			"total":     len(tasks),
			"page":      1,
			"page_size": len(tasks),
		},
	})
}

// GetRoleDetail 获取角色详情
// GET /api/admin/roles/:id
func GetRoleDetail(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var role model.Role
	if err := db.First(&role, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "角色不存在",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    role,
	})
}

// CreateRole 创建角色
// POST /api/admin/roles
func CreateRole(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Permissions string `json:"permissions"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
			"data": nil,
		})
		return
	}

	db := database.GetDB()
	role := model.Role{
		Name:        req.Name,
		Description: req.Description,
		Permissions: req.Permissions,
	}

	if err := db.Create(&role).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建角色失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": role.ID,
		},
	})
}

// UpdateRole 更新角色
// PUT /api/admin/roles/:id
func UpdateRole(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Permissions string `json:"permissions"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
			"data": nil,
		})
		return
	}

	db := database.GetDB()
	var role model.Role
	if err := db.First(&role, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "角色不存在",
			"data": nil,
		})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Permissions != "" {
		updates["permissions"] = req.Permissions
	}

	db.Model(&role).Updates(updates)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
	})
}

// DeleteRole 删除角色
// DELETE /api/admin/roles/:id
func DeleteRole(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var role model.Role
	if err := db.First(&role, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "角色不存在",
			"data": nil,
		})
		return
	}

	if role.IsSuper {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "不能删除超级管理员角色",
			"data": nil,
		})
		return
	}

	db.Delete(&role)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}

// CopyRole 复制角色
// POST /api/admin/roles/:id/copies
func CopyRole(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var role model.Role
	if err := db.First(&role, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "角色不存在",
			"data": nil,
		})
		return
	}

	// 创建副本
	newRole := model.Role{
		Name:        role.Name + " (副本)",
		Description: role.Description,
		Permissions: role.Permissions,
		IsSuper:     false,
	}

	if err := db.Create(&newRole).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "复制角色失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "复制成功",
		"data": gin.H{
			"id": newRole.ID,
		},
	})
}

// GetPermissions 获取权限列表
// GET /api/admin/permissions
func GetPermissions(c *gin.Context) {
	// 返回系统所有权限
	permissions := []gin.H{
		{"key": "dashboard_view", "name": "查看仪表盘"},
		{"key": "user_list", "name": "用户列表"},
		{"key": "user_detail", "name": "用户详情"},
		{"key": "user_manage", "name": "用户管理"},
		{"key": "order_list", "name": "订单列表"},
		{"key": "order_detail", "name": "订单详情"},
		{"key": "invoice_list", "name": "账单列表"},
		{"key": "invoice_manage", "name": "账单管理"},
		{"key": "ticket_list", "name": "工单列表"},
		{"key": "ticket_manage", "name": "工单管理"},
		{"key": "product_list", "name": "产品列表"},
		{"key": "product_manage", "name": "产品管理"},
		{"key": "supplier_list", "name": "供应商列表"},
		{"key": "supplier_manage", "name": "供应商管理"},
		{"key": "plugin_list", "name": "插件列表"},
		{"key": "plugin_manage", "name": "插件管理"},
		{"key": "settings_view", "name": "查看设置"},
		{"key": "settings_manage", "name": "设置管理"},
		{"key": "log_list", "name": "查看日志"},
		{"key": "staff_list", "name": "员工列表"},
		{"key": "staff_manage", "name": "员工管理"},
		{"key": "role_list", "name": "角色列表"},
		{"key": "role_manage", "name": "角色管理"},
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    permissions,
	})
}
