package admin

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/YeXuanHs/anchor-finance/internal/service"
	"github.com/gin-gonic/gin"
)

// GetStaffList 获取员工列表
// GET /api/admin/staff
func GetStaffList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetDB()
	query := db.Model(&model.Staff{})

	if keyword != "" {
		query = query.Where("username LIKE ? OR real_name LIKE ? OR email LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var staff []model.Staff
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&staff)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      staff,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetStaffDetail 获取员工详情
// GET /api/admin/staff/:id
func GetStaffDetail(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var staff model.Staff
	if err := db.First(&staff, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "员工不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    staff,
	})
}

// CreateStaff 创建员工
// POST /api/admin/staff
func CreateStaff(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		RealName string `json:"real_name"`
		Phone    string `json:"phone"`
		RoleID   uint   `json:"role_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 检查用户名是否已存在
	db := database.GetDB()
	var count int64
	db.Model(&model.Staff{}).Where("username = ? OR email = ?", req.Username, req.Email).Count(&count)
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "用户名或邮箱已存在",
		})
		return
	}

	// 创建员工
	hashedPassword, err := service.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "密码加密失败",
		})
		return
	}

	staff := model.Staff{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		RealName:     req.RealName,
		Phone:        req.Phone,
		RoleID:       req.RoleID,
		Status:       "active",
	}

	if err := db.Create(&staff).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建员工失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": staff.ID,
		},
	})
}

// UpdateStaff 更新员工
// PUT /api/admin/staff/:id
func UpdateStaff(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Email    string `json:"email"`
		RealName string `json:"real_name"`
		Phone    string `json:"phone"`
		RoleID   uint   `json:"role_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	db := database.GetDB()
	var staff model.Staff
	if err := db.First(&staff, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "员工不存在",
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
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.RoleID > 0 {
		updates["role_id"] = req.RoleID
	}

	db.Model(&staff).Updates(updates)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
	})
}

// DeleteStaff 删除员工
// DELETE /api/admin/staff/:id
func DeleteStaff(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var staff model.Staff
	if err := db.First(&staff, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "员工不存在",
		})
		return
	}

	db.Delete(&staff)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}

// UpdateStaffStatus 更新员工状态
// PATCH /api/admin/staff/:id/status
func UpdateStaffStatus(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 验证状态值
	validStatuses := map[string]bool{
		"active":   true,
		"disabled": true,
	}
	if !validStatuses[req.Status] {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的状态值",
		})
		return
	}

	db := database.GetDB()
	var staff model.Staff
	if err := db.First(&staff, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "员工不存在",
		})
		return
	}

	db.Model(&staff).Update("status", req.Status)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "状态更新成功",
	})
}

// ResetStaffPassword 重置员工密码
// POST /api/admin/staff/:id/password-resets
func ResetStaffPassword(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	db := database.GetDB()
	var staff model.Staff
	if err := db.First(&staff, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "员工不存在",
		})
		return
	}

	// 生成新密码hash
	hashedPassword, err := service.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "密码加密失败",
		})
		return
	}

	db.Model(&staff).Update("password_hash", hashedPassword)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "密码重置成功",
	})
}
