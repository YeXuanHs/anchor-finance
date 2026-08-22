package admin

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/YeXuanHs/anchor-finance/internal/service"
	"github.com/gin-gonic/gin"
)

// GetUserList 获取用户列表
// GET /api/admin/users
func GetUserList(c *gin.Context) {
	// 1. 解析分页参数
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

	// 2. 构建查询
	db := database.GetDB()
	query := db.Model(&model.User{})

	// 关键词搜索
	if keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ? OR phone LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 3. 获取总数
	var total int64
	query.Count(&total)

	// 4. 分页查询
	var users []model.User
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&users)

	// 5. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      users,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetUser 获取用户详情
// GET /api/admin/users/:id
func GetUser(c *gin.Context) {
	// 1. 获取用户ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的用户ID",
			"data":    nil,
		})
		return
	}

	// 2. 查询用户
	db := database.GetDB()
	var user model.User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "用户不存在",
			"data":    nil,
		})
		return
	}

	// 3. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    user,
	})
}

// CreateUser 创建用户
// POST /api/admin/users
func CreateUser(c *gin.Context) {
	// 1. 解析请求参数
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Phone    string `json:"phone"`
		Company  string `json:"company"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 2. 检查用户名是否已存在
	db := database.GetDB()
	var count int64
	db.Model(&model.User{}).Where("username = ? OR email = ?", req.Username, req.Email).Count(&count)
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "用户名或邮箱已存在",
			"data":    nil,
		})
		return
	}

	// 3. 创建用户
	hashedPassword, err := service.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "密码加密失败",
			"data":    nil,
		})
		return
	}

	user := model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Phone:        req.Phone,
		Company:      req.Company,
		Status:       "active",
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建用户失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 4. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": user.ID,
		},
	})
}

// UpdateUser 更新用户
// PUT /api/admin/users/:id
func UpdateUser(c *gin.Context) {
	// 1. 获取用户ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的用户ID",
			"data":    nil,
		})
		return
	}

	// 2. 解析请求参数
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Company  string `json:"company"`
		Status   string `json:"status"`
		Balance  float64 `json:"balance"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 3. 查询用户
	db := database.GetDB()
	var user model.User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "用户不存在",
			"data":    nil,
		})
		return
	}

	// 4. 更新用户
	updates := map[string]interface{}{}
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Company != "" {
		updates["company"] = req.Company
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.Balance != 0 {
		updates["balance"] = req.Balance
	}

	if err := db.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "更新用户失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 5. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
		"data":    nil,
	})
}

// DeleteUser 删除用户（软删除）
// DELETE /api/admin/users/:id
func DeleteUser(c *gin.Context) {
	// 1. 获取用户ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的用户ID",
			"data":    nil,
		})
		return
	}

	// 2. 查询用户
	db := database.GetDB()
	var user model.User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "用户不存在",
			"data":    nil,
		})
		return
	}

	// 3. 软删除用户
	if err := db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "删除用户失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 4. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
		"data":    nil,
	})
}

// GetUserOrders 获取用户订单
// GET /api/admin/users/:id/orders
func GetUserOrders(c *gin.Context) {
	// 1. 获取用户ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的用户ID",
			"data":    nil,
		})
		return
	}

	// 2. 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 3. 查询用户订单
	db := database.GetDB()
	var orders []model.Order
	var total int64

	db.Model(&model.Order{}).Where("user_id = ?", id).Count(&total)
	db.Where("user_id = ?", id).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("id DESC").
		Find(&orders)

	// 4. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      orders,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetUserInvoices 获取用户账单
// GET /api/admin/users/:id/invoices
func GetUserInvoices(c *gin.Context) {
	// 1. 获取用户ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的用户ID",
			"data":    nil,
		})
		return
	}

	// 2. 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 3. 查询用户账单
	db := database.GetDB()
	var invoices []model.Invoice
	var total int64

	db.Model(&model.Invoice{}).Where("user_id = ?", id).Count(&total)
	db.Where("user_id = ?", id).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("id DESC").
		Find(&invoices)

	// 4. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      invoices,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetUserTickets 获取用户工单
// GET /api/admin/users/:id/tickets
func GetUserTickets(c *gin.Context) {
	// 1. 获取用户ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的用户ID",
			"data":    nil,
		})
		return
	}

	// 2. 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 3. 查询用户工单
	db := database.GetDB()
	var tickets []model.Ticket
	var total int64

	db.Model(&model.Ticket{}).Where("user_id = ?", id).Count(&total)
	db.Where("user_id = ?", id).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("id DESC").
		Find(&tickets)

	// 4. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      tickets,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetUserServices 获取用户服务
// GET /api/admin/users/:id/services
func GetUserServices(c *gin.Context) {
	// 1. 获取用户ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的用户ID",
			"data":    nil,
		})
		return
	}

	// 2. 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 3. 查询用户服务
	db := database.GetDB()
	var services []model.Service
	var total int64

	db.Model(&model.Service{}).Where("user_id = ?", id).Count(&total)
	db.Where("user_id = ?", id).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("id DESC").
		Find(&services)

	// 4. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      services,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// UpdateUserStatus 更新用户状态
// PATCH /api/admin/users/:id/status
func UpdateUserStatus(c *gin.Context) {
	// 1. 获取用户ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的用户ID",
			"data":    nil,
		})
		return
	}

	// 2. 解析请求参数
	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 3. 验证状态值
	validStatuses := map[string]bool{
		"active":    true,
		"suspended": true,
		"closed":    true,
	}
	if !validStatuses[req.Status] {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的状态值",
			"data":    nil,
		})
		return
	}

	// 4. 查询用户
	db := database.GetDB()
	var user model.User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "用户不存在",
			"data":    nil,
		})
		return
	}

	// 5. 更新状态
	if err := db.Model(&user).Update("status", req.Status).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "更新状态失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 6. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "状态更新成功",
		"data":    nil,
	})
}

// GetUserBalanceLogs 获取用户余额日志
// GET /api/admin/users/:id/balance-logs
func GetUserBalanceLogs(c *gin.Context) {
	// 1. 获取用户ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的用户ID",
			"data":    nil,
		})
		return
	}

	// 2. 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 3. 查询余额变动（充值记录）
	db := database.GetDB()
	var total int64
	db.Model(&model.Recharge{}).Where("user_id = ?", id).Count(&total)

	var logs []model.Recharge
	offset := (page - 1) * pageSize
	db.Where("user_id = ?", id).
		Offset(offset).
		Limit(pageSize).
		Order("id DESC").
		Find(&logs)

	// 4. 返回统一格式
	if logs == nil {
		logs = []model.Recharge{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      logs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// RefundUserService 退款用户服务
// POST /api/admin/users/:id/services/:service_id/refunds
func RefundUserService(c *gin.Context) {
	// 1. 获取用户ID和服务ID
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的用户ID",
			"data":    nil,
		})
		return
	}

	serviceID, err := strconv.ParseUint(c.Param("service_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的服务ID",
			"data":    nil,
		})
		return
	}

	// 2. 解析请求参数
	var req struct {
		Amount float64 `json:"amount" binding:"required"`
		Reason string  `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 3. 查询用户
	db := database.GetDB()
	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "用户不存在",
			"data":    nil,
		})
		return
	}

	// 4. 查询服务
	var service model.Service
	if err := db.Where("id = ? AND user_id = ?", serviceID, userID).First(&service).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "服务不存在",
			"data":    nil,
		})
		return
	}

	// 5. 退款到用户余额
	newBalance := user.Balance + req.Amount
	if err := db.Model(&user).Update("balance", newBalance).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "退款失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 6. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "退款成功",
		"data": gin.H{
			"balance": newBalance,
		},
	})
}

// GetUserOperationLogs 获取用户操作日志
// GET /api/admin/users/:id/operation-logs
func GetUserOperationLogs(c *gin.Context) {
	// 1. 获取用户ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的用户ID",
			"data":    nil,
		})
		return
	}

	// 2. 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 3. 查询用户操作日志
	db := database.GetDB()
	var logs []model.OperationLog
	var total int64

	db.Model(&model.OperationLog{}).Where("user_id = ?", id).Count(&total)
	db.Where("user_id = ?", id).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("id DESC").
		Find(&logs)

	// 4. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      logs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// RechargeUser 用户充值
// POST /api/admin/users/:id/recharges
func RechargeUser(c *gin.Context) {
	// 1. 获取用户ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的用户ID",
			"data":    nil,
		})
		return
	}

	// 2. 解析请求参数
	var req struct {
		Amount float64 `json:"amount" binding:"required"`
		Remark string  `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 3. 查询用户
	db := database.GetDB()
	var user model.User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "用户不存在",
			"data":    nil,
		})
		return
	}

	// 4. 更新余额
	newBalance := user.Balance + req.Amount
	if err := db.Model(&user).Update("balance", newBalance).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "充值失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 5. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "充值成功",
		"data": gin.H{
			"balance": newBalance,
		},
	})
}
