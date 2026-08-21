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
