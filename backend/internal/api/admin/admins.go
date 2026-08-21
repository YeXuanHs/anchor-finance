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
			"message": "参数错误: " + err.Error(),
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
			"message": "创建管理员失败: " + err.Error(),
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
			"message": "参数错误: " + err.Error(),
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
	// 暂时返回空列表，后续实现定时任务
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      []interface{}{},
			"total":     0,
			"page":      1,
			"page_size": 20,
		},
	})
}
