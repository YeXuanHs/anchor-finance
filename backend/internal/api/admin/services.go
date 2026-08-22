package admin

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/YeXuanHs/anchor-finance/internal/pluginengine"
	"github.com/gin-gonic/gin"
)

// GetServiceList 获取服务列表
// GET /api/admin/services
func GetServiceList(c *gin.Context) {
	// 1. 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	status := c.Query("status")
	userID := c.Query("user_id")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 2. 构建查询
	db := database.GetDB()
	query := db.Model(&model.Service{})

	// 关键词搜索（产品名称、域名、用户名）
	if keyword != "" {
		query = query.Where("product_name LIKE ? OR domain LIKE ? OR username LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 用户筛选
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// 3. 获取总数
	var total int64
	query.Count(&total)

	// 4. 分页查询
	var services []model.Service
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&services)

	// 5. 返回统一格式
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

// GetService 获取服务详情
// GET /api/admin/services/:id
func GetService(c *gin.Context) {
	// 1. 获取服务ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的服务ID",
			"data":    nil,
		})
		return
	}

	// 2. 查询服务
	db := database.GetDB()
	var service model.Service
	if err := db.First(&service, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "服务不存在",
			"data":    nil,
		})
		return
	}

	// 3. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    service,
	})
}

// UpdateService 更新服务
// PUT /api/admin/services/:id
func UpdateService(c *gin.Context) {
	// 1. 获取服务ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
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
		ProductName  string  `json:"product_name"`
		Domain       string  `json:"domain"`
		Username     string  `json:"username"`
		Status       string  `json:"status"`
		BillingCycle string  `json:"billing_cycle"`
		Amount       float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 3. 查询服务
	db := database.GetDB()
	var service model.Service
	if err := db.First(&service, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "服务不存在",
			"data":    nil,
		})
		return
	}

	// 4. 更新服务
	updates := map[string]interface{}{}
	if req.ProductName != "" {
		updates["product_name"] = req.ProductName
	}
	if req.Domain != "" {
		updates["domain"] = req.Domain
	}
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.BillingCycle != "" {
		updates["billing_cycle"] = req.BillingCycle
	}
	if req.Amount > 0 {
		updates["amount"] = req.Amount
	}

	if err := db.Model(&service).Updates(updates).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "更新服务失败: " + err.Error(),
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

// SuspendService 暂停服务
// POST /api/admin/services/:id/suspend
func SuspendService(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "无效的服务ID"})
		return
	}

	db := database.GetDB()
	var service model.Service
	if err := db.First(&service, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "服务不存在"})
		return
	}

	if service.Status != "active" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "只有活跃状态的服务才能暂停"})
		return
	}

	// 调用PHP插件引擎执行上游暂停操作
	if _, err := pluginengine.TriggerHook("suspend_service", map[string]interface{}{
		"service_id": service.ID,
		"product_id": service.ProductID,
		"user_id":    service.UserID,
	}); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 502, "message": "插件引擎离线: " + err.Error()})
		return
	}

	// 上游暂停成功后，更新本地状态
	db.Model(&service).Update("status", "suspended")

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "暂停成功"})
}

// UnsuspendService 取消暂停
// POST /api/admin/services/:id/unsuspend
func UnsuspendService(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "无效的服务ID"})
		return
	}

	db := database.GetDB()
	var service model.Service
	if err := db.First(&service, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "服务不存在"})
		return
	}

	if service.Status != "suspended" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "只有暂停状态的服务才能取消暂停"})
		return
	}

	// 调用PHP插件引擎执行上游取消暂停操作
	if _, err := pluginengine.TriggerHook("unsuspend_service", map[string]interface{}{
		"service_id": service.ID,
		"product_id": service.ProductID,
		"user_id":    service.UserID,
	}); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 502, "message": "插件引擎离线: " + err.Error()})
		return
	}

	db.Model(&service).Update("status", "active")

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "取消暂停成功"})
}

// TerminateService 终止服务
// POST /api/admin/services/:id/terminate
func TerminateService(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "无效的服务ID"})
		return
	}

	db := database.GetDB()
	var service model.Service
	if err := db.First(&service, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "服务不存在"})
		return
	}

	if service.Status == "terminated" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "服务已终止"})
		return
	}

	// 调用PHP插件引擎执行上游终止操作
	if _, err := pluginengine.TriggerHook("terminate_service", map[string]interface{}{
		"service_id": service.ID,
		"product_id": service.ProductID,
		"user_id":    service.UserID,
	}); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 502, "message": "插件引擎离线: " + err.Error()})
		return
	}

	db.Model(&service).Update("status", "terminated")

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "终止成功"})
}
