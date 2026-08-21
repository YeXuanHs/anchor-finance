package client

import (
	"net/http"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// PowerService 服务电源操作
// POST /api/client/services/:id/power-actions
func PowerService(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")

	var req struct {
		Action string `json:"action" binding:"required"` // on, off, reboot, hard_off, hard_reboot
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 验证操作类型
	validActions := map[string]bool{
		"on":           true,
		"off":          true,
		"reboot":       true,
		"hard_off":     true,
		"hard_reboot":  true,
	}
	if !validActions[req.Action] {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的操作类型",
		})
		return
	}

	// 验证服务属于该用户
	db := database.GetDB()
	var service model.Service
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&service).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "服务不存在",
		})
		return
	}

	// TODO: 调用上游模块执行电源操作
	// 这里暂时只返回成功

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "操作成功",
	})
}

// ResetServicePassword 重置服务密码
// POST /api/client/services/:id/password-resets
func ResetServicePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")
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

	// 验证服务属于该用户
	db := database.GetDB()
	var service model.Service
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&service).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "服务不存在",
		})
		return
	}

	// TODO: 调用上游模块重置密码
	// 这里暂时只返回成功

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "密码重置成功",
	})
}

// ReinstallService 重装服务系统
// POST /api/client/services/:id/reinstallations
func ReinstallService(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")

	var req struct {
		OS string `json:"os" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 验证服务属于该用户
	db := database.GetDB()
	var service model.Service
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&service).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "服务不存在",
		})
		return
	}

	// TODO: 调用上游模块重装系统
	// 这里暂时只返回成功

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "重装请求已提交",
	})
}

// GetServiceRenewPreview 获取服务续费预览
// GET /api/client/services/:id/renewals
func GetServiceRenewPreview(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")

	// 验证服务属于该用户
	db := database.GetDB()
	var service model.Service
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&service).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "服务不存在",
		})
		return
	}

	// 返回续费信息
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"service_id":    service.ID,
			"product_name":  service.ProductName,
			"billing_cycle": service.BillingCycle,
			"amount":        service.Amount,
			"next_due_date": service.NextDueDate,
		},
	})
}

// CreateRenewOrder 创建续费订单
// POST /api/client/services/:id/renewals
func CreateRenewOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")

	var req struct {
		Cycle string `json:"cycle" binding:"required"` // monthly, quarterly, yearly
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 验证服务属于该用户
	db := database.GetDB()
	var service model.Service
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&service).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "服务不存在",
		})
		return
	}

	// TODO: 计算续费金额，创建续费订单
	// 这里暂时返回成功

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "续费订单创建成功",
	})
}

// GetUserServices 获取用户服务列表（带分组概览）
// GET /api/client/services/grouped-overview
func GetUserServicesGroupedOverview(c *gin.Context) {
	userID, _ := c.Get("user_id")

	db := database.GetDB()

	// 按产品名称分组统计
	type GroupInfo struct {
		ProductName string `json:"product_name"`
		Count       int64  `json:"count"`
	}

	var groups []GroupInfo
	db.Model(&model.Service{}).
		Select("product_name, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("product_name").
		Scan(&groups)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    groups,
	})
}

// UpdateServiceName 更新服务名称
// PUT /api/client/services/:id/name
func UpdateServiceName(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 验证服务属于该用户
	db := database.GetDB()
	var service model.Service
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&service).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "服务不存在",
		})
		return
	}

	// 更新服务名称（使用domain字段存储自定义名称）
	db.Model(&service).Update("domain", req.Name)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
	})
}

// UpdateServiceRemark 更新服务备注
// PUT /api/client/services/:id/remark
func UpdateServiceRemark(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")

	var req struct {
		Remark string `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 验证服务属于该用户
	db := database.GetDB()
	var service model.Service
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&service).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "服务不存在",
		})
		return
	}

	// TODO: 更新服务备注（需要添加remark字段到service表）

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
	})
}

// GetServiceConnection 获取服务连接信息
// GET /api/client/services/:id/connection
func GetServiceConnection(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")

	// 验证服务属于该用户
	db := database.GetDB()
	var service model.Service
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&service).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "服务不存在",
		})
		return
	}

	// TODO: 从上游模块获取连接信息
	// 这里暂时返回基本信息

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"service_id": service.ID,
			"product_name": service.ProductName,
			"domain": service.Domain,
			"username": service.Username,
			"status": service.Status,
		},
	})
}

// GetServiceRuntime 获取服务运行时信息
// GET /api/client/services/:id/runtime
func GetServiceRuntime(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")

	// 验证服务属于该用户
	db := database.GetDB()
	var service model.Service
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&service).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "服务不存在",
		})
		return
	}

	// TODO: 从上游模块获取运行时信息（CPU、内存、磁盘等）

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"service_id": service.ID,
			"status": service.Status,
		},
	})
}

// GetServiceStatus 获取服务状态
// GET /api/client/services/:id/module-status
func GetServiceStatus(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")

	// 验证服务属于该用户
	db := database.GetDB()
	var service model.Service
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&service).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "服务不存在",
		})
		return
	}

	// TODO: 从上游模块获取实时状态
	// 这里暂时返回数据库中的状态

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"status": service.Status,
		},
	})
}
