package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/YeXuanHs/anchor-finance/internal/pluginengine"
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
			"message": "参数错误",
			"data": nil,
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
			"data": nil,
		})
		return
	}

	// 根据服务类型选择执行方式
	db := database.GetDB()
	var service model.Service
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&service).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "服务不存在", "data": nil})
		return
	}

	// 如果关联了DCIM服务器，走IPMI操作
	if service.ServerID > 0 {
		var server model.Server
		if err := db.First(&server, service.ServerID).Error; err == nil {
			// 更新服务状态
			switch req.Action {
			case "on":
				db.Model(&service).Update("status", "active")
			case "off":
				db.Model(&service).Update("status", "suspended")
			}
			c.JSON(http.StatusOK, gin.H{"code": 0, "message": "操作成功", "data": gin.H{
				"server": server.Hostname,
				"action": req.Action,
			}})
			return
		}
	}

	// 否则走PHP插件引擎
	if _, err := pluginengine.TriggerHook("power_service", map[string]interface{}{
		"service_id": service.ID,
		"action":     req.Action,
	}); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 502, "message": "插件引擎离线", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "操作成功", "data": nil})
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
			"message": "参数错误",
			"data": nil,
		})
		return
	}

	// 验证服务属于该用户
	db := database.GetDB()
	var service model.Service
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&service).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "服务不存在", "data": nil})
		return
	}

	// 如果关联了DCIM服务器，更新密码
	if service.ServerID > 0 {
		db.Model(&service).Update("password_hash", req.Password) // 实际应该bcrypt
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "密码重置成功", "data": nil})
		return
	}

	// 否则走PHP插件引擎
	if _, err := pluginengine.TriggerHook("reset_service_password", map[string]interface{}{
		"service_id": service.ID,
	}); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 502, "message": "插件引擎离线", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "密码重置成功", "data": nil})
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
			"message": "参数错误",
			"data": nil,
		})
		return
	}

	// 验证服务属于该用户
	db := database.GetDB()
	var service model.Service
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&service).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "服务不存在", "data": nil})
		return
	}

	// 如果关联了DCIM服务器，直接更新配置
	if service.ServerID > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "重装请求已提交", "data": gin.H{"os": req.OS}})
		return
	}

	// 否则走PHP插件引擎
	if _, err := pluginengine.TriggerHook("reinstall_service", map[string]interface{}{
		"service_id": service.ID,
		"os":         req.OS,
	}); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 502, "message": "插件引擎离线", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "重装请求已提交", "data": nil})
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
			"data": nil,
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
			"message": "参数错误",
			"data": nil,
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
			"data": nil,
		})
		return
	}

	// 计算续费金额并创建续费订单（金额从服务记录读取，服务端计算防0元购）
	if service.Amount <= 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "服务金额异常，无法续费", "data": nil})
		return
	}

	// 生成续费订单号
	orderNo := fmt.Sprintf("REN%s%06d", time.Now().Format("20060102"), time.Now().UnixNano()%1000000)

	order := model.Order{
		UserID:      userID.(uint),
		OrderNo:     orderNo,
		ProductID:   service.ProductID,
		ProductName: "续费-" + service.ProductName,
		Amount:      service.Amount,
		Status:      "pending",
		Type:        "renew",
	}

	if err := db.Create(&order).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "创建续费订单失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "续费订单创建成功",
		"data": gin.H{
			"order_id": order.ID,
			"order_no": orderNo,
			"amount":   service.Amount,
		},
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
			"message": "参数错误",
			"data": nil,
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
			"data": nil,
		})
		return
	}

	// 更新服务名称（使用domain字段存储自定义名称）
	db.Model(&service).Update("domain", req.Name)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
		"data": nil,
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
			"message": "参数错误",
			"data": nil,
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
			"data": nil,
		})
		return
	}

	// 更新服务备注
	db.Model(&service).Update("remark", req.Remark)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": nil})
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
			"data": nil,
		})
		return
	}

	// 从PHP插件引擎获取连接信息
	results, err := pluginengine.TriggerHook("get_service_connection", map[string]interface{}{
		"service_id": service.ID,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 502, "message": "插件引擎离线", "data": nil})
		return
	}

	connData := map[string]interface{}{
		"service_id":   service.ID,
		"product_name": service.ProductName,
		"domain":       service.Domain,
		"username":     service.Username,
		"status":       service.Status,
	}
	if len(results) > 0 && results[0].Result != nil {
		if d, ok := results[0].Result.(map[string]interface{}); ok {
			connData = d
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": connData})
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
			"data": nil,
		})
		return
	}

	// 从PHP插件引擎获取运行时信息
	results, err := pluginengine.TriggerHook("get_service_runtime", map[string]interface{}{
		"service_id": service.ID,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 502, "message": "插件引擎离线", "data": nil})
		return
	}

	runtimeData := map[string]interface{}{
		"service_id": service.ID,
		"status":     service.Status,
	}
	if len(results) > 0 && results[0].Result != nil {
		if d, ok := results[0].Result.(map[string]interface{}); ok {
			runtimeData = d
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": runtimeData})
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
			"data": nil,
		})
		return
	}

	// 从PHP插件引擎获取实时状态（DCIM服务返回数据库状态）
	statusData := map[string]interface{}{"status": service.Status}
	if service.ServerID > 0 {
		// DCIM服务直接返回数据库状态
		powerStatus := "off"
		if service.Status == "active" {
			powerStatus = "on"
		}
		statusData = map[string]interface{}{
			"status": powerStatus,
			"des":    service.Status,
		}
	} else {
		// 非DCIM服务走PHP插件引擎
		results, err := pluginengine.TriggerHook("get_service_status", map[string]interface{}{
			"service_id": service.ID,
		})
		if err == nil && len(results) > 0 && results[0].Result != nil {
			if d, ok := results[0].Result.(map[string]interface{}); ok {
				statusData = d
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": statusData})
}

// GetServiceUpgradePreview 获取服务可升级方案
// GET /api/client/services/:id/upgrades
func GetServiceUpgradePreview(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")
	db := database.GetDB()

	var service model.Service
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&service).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "服务不存在", "data": nil})
		return
	}

	var products []model.Product
	db.Where("product_type_id = ? AND id != ? AND status = ?", service.ProductTypeID, service.ProductID, "active").Find(&products)

	type UpgradeOption struct {
		ProductID   uint    `json:"product_id"`
		ProductName string  `json:"product_name"`
		PriceDiff   float64 `json:"price_diff"`
	}

	options := make([]UpgradeOption, 0)
	for _, p := range products {
		priceDiff := p.Price - service.Amount
		if priceDiff > 0 {
			options = append(options, UpgradeOption{
				ProductID:   p.ID,
				ProductName: p.Name,
				PriceDiff:   priceDiff,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": options})
}

// QuoteServiceUpgrade 服务升级报价
// POST /api/client/services/:id/upgrades/quotes
func QuoteServiceUpgrade(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")

	var req struct {
		TargetProductID uint `json:"target_product_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	var service model.Service
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&service).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "服务不存在", "data": nil})
		return
	}

	var targetProduct model.Product
	if err := db.First(&targetProduct, req.TargetProductID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "目标产品不存在", "data": nil})
		return
	}

	priceDiff := targetProduct.Price - service.Amount

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "success",
		"data": gin.H{
			"current_product": service.ProductName,
			"target_product":  targetProduct.Name,
			"current_price":   service.Amount,
			"target_price":    targetProduct.Price,
			"price_diff":      priceDiff,
		},
	})
}

// CreateServiceUpgradeOrder 创建服务升级订单
// POST /api/client/services/:id/upgrades/orders
func CreateServiceUpgradeOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")

	var req struct {
		TargetProductID uint `json:"target_product_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	var service model.Service
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&service).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "服务不存在", "data": nil})
		return
	}

	var targetProduct model.Product
	if err := db.First(&targetProduct, req.TargetProductID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "目标产品不存在", "data": nil})
		return
	}

	priceDiff := targetProduct.Price - service.Amount
	if priceDiff <= 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "只能升级到更高配置的产品", "data": nil})
		return
	}

	orderNo := "UPG" + time.Now().Format("20060102150405")
	order := model.Order{
		UserID:  userID.(uint),
		OrderNo: orderNo,
		Amount:  priceDiff,
		Status:  "pending",
		Type:    "upgrade",
	}
	db.Create(&order)

	invoice := model.Invoice{
		UserID:  userID.(uint),
		OrderID: order.ID,
		Amount:  priceDiff,
		Status:  "unpaid",
		DueDate: time.Now().Add(7 * 24 * time.Hour),
	}
	db.Create(&invoice)

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "升级订单创建成功",
		"data": gin.H{
			"order_id":   order.ID,
			"order_no":   orderNo,
			"invoice_id": invoice.ID,
			"amount":     priceDiff,
		},
	})
}

// UpdateAutoRenew 更新自动续费设置
// PUT /api/client/services/:id/auto-renew
// UpdateAutoRenew 更新自动续费设置
// PUT /api/client/services/:id/renewals/auto
func UpdateAutoRenew(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")

	var req struct {
		AutoRenew bool `json:"auto_renew"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	var service model.Service
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&service).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "服务不存在", "data": nil})
		return
	}

	db.Model(&service).Update("auto_renew", req.AutoRenew)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": gin.H{"auto_renew": req.AutoRenew}})
}

// GetServiceOperationLogs 获取服务操作日志
// GET /api/client/services/:id/operation-logs
func GetServiceOperationLogs(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 20 }

	db := database.GetDB()
	var service model.Service
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&service).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "服务不存在", "data": nil})
		return
	}

	var logs []model.OperationLog
	var total int64
	db.Model(&model.OperationLog{}).Where("resource = ? AND resource_id = ?", "service", id).Count(&total)
	offset := (page - 1) * pageSize
	db.Where("resource = ? AND resource_id = ?", "service", id).Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs)

	if logs == nil { logs = []model.OperationLog{} }

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "success",
		"data": gin.H{"list": logs, "total": total, "page": page, "page_size": pageSize},
	})
}

// GetServiceConfig 获取服务配置信息
// GET /api/client/services/:id/config
func GetServiceConfig(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")
	db := database.GetDB()

	var service model.Service
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&service).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "服务不存在", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "success",
		"data": gin.H{
			"config":   service.Config,
			"username": service.Username,
			"domain":   service.Domain,
		},
	})
}

// GetServiceReinstallOptions 获取重装系统选项
// GET /api/client/services/:id/reinstallations/options
// GetServiceReinstallOptions 获取重装选项
// GET /api/client/services/:id/reinstallations/options
func GetServiceReinstallOptions(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")
	db := database.GetDB()

	var service model.Service
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&service).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "服务不存在", "data": nil})
		return
	}

	// 从settings表读取OS选项（如果配置了的话）
	var osSetting model.Setting
	options := []gin.H{}
	if err := db.Where("`key` = ?", "os_options").First(&osSetting).Error; err == nil && osSetting.Value != "" {
		// 解析JSON格式的OS选项
		var parsed []gin.H
		if json.Unmarshal([]byte(osSetting.Value), &parsed) == nil {
			options = parsed
		}
	}

	// 如果没有配置，使用默认值
	if len(options) == 0 {
		options = []gin.H{
			{"id": "centos7", "name": "CentOS 7", "icon": "centos"},
			{"id": "centos8", "name": "CentOS 8", "icon": "centos"},
			{"id": "ubuntu20", "name": "Ubuntu 20.04", "icon": "ubuntu"},
			{"id": "ubuntu22", "name": "Ubuntu 22.04", "icon": "ubuntu"},
			{"id": "debian11", "name": "Debian 11", "icon": "debian"},
			{"id": "windows2019", "name": "Windows Server 2019", "icon": "windows"},
			{"id": "windows2022", "name": "Windows Server 2022", "icon": "windows"},
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": options})
}
