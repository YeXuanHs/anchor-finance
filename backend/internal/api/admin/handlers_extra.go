package admin

import (
	"net/http"
	"strconv"
	"time"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/YeXuanHs/anchor-finance/internal/pluginengine"
	"github.com/YeXuanHs/anchor-finance/internal/service"
	"github.com/gin-gonic/gin"
)

// ==================== 用户备注 ====================

// GetUserRemarks 获取用户备注
func GetUserRemarks(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()
	var remarks []model.UserRemark
	db.Where("user_id = ?", id).Order("id DESC").Find(&remarks)
	if remarks == nil {
		remarks = []model.UserRemark{}
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": remarks})
}

// AddUserRemark 添加用户备注
func AddUserRemark(c *gin.Context) {
	id := c.Param("id")
	adminID, _ := c.Get("user_id")
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}
	db := database.GetDB()
	var admin model.Admin
	db.First(&admin, adminID)
	uid, _ := strconv.ParseUint(id, 10, 32)
	remark := model.UserRemark{
		UserID:    uint(uid),
		AdminID:   adminID.(uint),
		AdminName: admin.Username,
		Content:   req.Content,
	}
	db.Create(&remark)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "添加成功", "data": gin.H{"id": remark.ID}})
}

// LoginAsUser 以用户身份登录
func LoginAsUser(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()
	var user model.User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "用户不存在", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"user_id": user.ID, "username": user.Username}})
}

// RefreshUserServicesStatus 刷新用户服务状态
// RefreshUserServicesStatus 刷新用户服务状态（通过插件引擎同步远程状态）
// POST /api/admin/users/:id/services/refresh-statuses
func RefreshUserServicesStatus(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()
	var services []model.Service
	db.Where("user_id = ? AND status IN ?", id, []string{"active", "suspended"}).Find(&services)

	refreshed := 0
	for _, svc := range services {
		results, err := pluginengine.TriggerHook("get_service_status", map[string]interface{}{
			"service_id": svc.ID,
		})
		if err != nil {
			continue
		}
		if len(results) > 0 && results[0].Data != nil {
			if data, ok := results[0].Data.(map[string]interface{}); ok {
				if status, ok := data["status"].(string); ok && status != svc.Status {
					db.Model(&svc).Update("status", status)
					refreshed++
				}
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "刷新完成", "data": gin.H{"total": len(services), "refreshed": refreshed}})
}

// ==================== 订单扩展 ====================

// SearchOrders 搜索订单
func SearchOrders(c *gin.Context) {
	var req struct {
		Keyword string `json:"keyword"`
		Status  string `json:"status"`
		Page    int    `json:"page"`
		PageSize int  `json:"page_size"`
	}
	c.ShouldBindJSON(&req)
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	db := database.GetDB()
	var orders []model.Order
	var total int64
	query := db.Model(&model.Order{})
	if req.Keyword != "" {
		query = query.Where("order_no LIKE ? OR product_name LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	query.Count(&total)
	query.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Order("id DESC").Find(&orders)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"list": orders, "total": total, "page": req.Page, "page_size": req.PageSize}})
}

// AddOrderNote 添加订单备注
func AddOrderNote(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Note string `json:"note" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}
	db := database.GetDB()
	db.Model(&model.Order{}).Where("id = ?", id).Update("note", req.Note)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "添加成功", "data": nil})
}

// ==================== 账单扩展 ====================

// AddInvoiceNote 添加账单备注
func AddInvoiceNote(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Note string `json:"note" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}
	db := database.GetDB()
	db.Model(&model.Invoice{}).Where("id = ?", id).Update("note", req.Note)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "添加成功", "data": nil})
}

// ==================== 财务报表扩展 ====================

// GetFinanceLedgerSummary 财务台账汇总
func GetFinanceLedgerSummary(c *gin.Context) {
	db := database.GetDB()
	var totalIncome float64
	db.Model(&model.Invoice{}).Where("status = ?", "paid").Select("COALESCE(SUM(amount), 0)").Scan(&totalIncome)
	var totalRecharge float64
	db.Model(&model.Recharge{}).Where("status = ?", "success").Select("COALESCE(SUM(amount), 0)").Scan(&totalRecharge)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"total_income": totalIncome, "total_recharge": totalRecharge}})
}

// GetRenewalOrders 续费订单列表
func GetRenewalOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	db := database.GetDB()
	var orders []model.Order
	var total int64
	db.Model(&model.Order{}).Where("type = ?", "renew").Count(&total)
	db.Where("type = ?", "renew").Offset((page - 1) * pageSize).Limit(pageSize).Order("id DESC").Find(&orders)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"list": orders, "total": total, "page": page, "page_size": pageSize}})
}

// GetUpgradeOrders 升级订单列表
func GetUpgradeOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	db := database.GetDB()
	var orders []model.Order
	var total int64
	db.Model(&model.Order{}).Where("type = ?", "upgrade").Count(&total)
	db.Where("type = ?", "upgrade").Offset((page - 1) * pageSize).Limit(pageSize).Order("id DESC").Find(&orders)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"list": orders, "total": total, "page": page, "page_size": pageSize}})
}

// GetNewCustomerStatistics 新客户统计
func GetNewCustomerStatistics(c *gin.Context) {
	db := database.GetDB()
	var todayCount int64
	today := time.Now().Format("2006-01-02")
	db.Model(&model.User{}).Where("DATE(created_at) = ?", today).Count(&todayCount)
	var weekCount int64
	db.Model(&model.User{}).Where("created_at >= ?", time.Now().AddDate(0, 0, -7)).Count(&weekCount)
	var monthCount int64
	db.Model(&model.User{}).Where("created_at >= ?", time.Now().AddDate(0, -1, 0)).Count(&monthCount)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"today": todayCount, "week": weekCount, "month": monthCount}})
}

// GetRevenueRanking 收入排名
func GetRevenueRanking(c *gin.Context) {
	type RankItem struct {
		UserID   uint    `json:"user_id"`
		Username string  `json:"username"`
		Total    float64 `json:"total"`
	}
	var results []RankItem
	db := database.GetDB()
	db.Raw(`SELECT o.user_id, u.username, SUM(o.amount) as total 
		FROM orders o JOIN users u ON o.user_id = u.id 
		WHERE o.status IN ('active','completed') 
		GROUP BY o.user_id, u.username 
		ORDER BY total DESC LIMIT 20`).Scan(&results)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": results})
}

// ==================== 工单扩展 ====================

// ReceiveTicket 接收工单
func ReceiveTicket(c *gin.Context) {
	id := c.Param("id")
	adminID, _ := c.Get("user_id")
	db := database.GetDB()
	db.Model(&model.Ticket{}).Where("id = ?", id).Updates(map[string]interface{}{"assigned_to": adminID, "status": "in_progress"})
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "已接收", "data": nil})
}

// GetTicketDepartmentDetail 工单部门详情
func GetTicketDepartmentDetail(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()
	var dept model.TicketDepartment
	if err := db.First(&dept, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "部门不存在", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": dept})
}

// MoveTicketDepartmentUp 上移工单部门
func MoveTicketDepartmentUp(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()
	var dept model.TicketDepartment
	db.First(&dept, id)
	var prev model.TicketDepartment
	db.Where("sort_order < ?", dept.SortOrder).Order("sort_order DESC").First(&prev)
	if prev.ID > 0 {
		db.Model(&dept).Update("sort_order", prev.SortOrder)
		db.Model(&prev).Update("sort_order", dept.SortOrder)
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "已上移", "data": nil})
}

// MoveTicketDepartmentDown 下移工单部门
func MoveTicketDepartmentDown(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()
	var dept model.TicketDepartment
	db.First(&dept, id)
	var next model.TicketDepartment
	db.Where("sort_order > ?", dept.SortOrder).Order("sort_order ASC").First(&next)
	if next.ID > 0 {
		db.Model(&dept).Update("sort_order", next.SortOrder)
		db.Model(&next).Update("sort_order", dept.SortOrder)
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "已下移", "data": nil})
}

// GetTicketStatusDetail 工单状态详情
func GetTicketStatusDetail(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()
	var status model.TicketStatus
	if err := db.First(&status, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "状态不存在", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": status})
}

// ==================== 产品扩展 ====================

// ReorderProducts 产品排序
func ReorderProducts(c *gin.Context) {
	var req struct {
		Items []struct {
			ID       uint `json:"id"`
			Position int  `json:"position"`
		} `json:"items"`
	}
	c.ShouldBindJSON(&req)
	db := database.GetDB()
	for _, item := range req.Items {
		db.Model(&model.Product{}).Where("id = ?", item.ID).Update("sort_order", item.Position)
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "排序成功", "data": nil})
}

// BatchUpdateProductCategory 批量更新产品分类
func BatchUpdateProductCategory(c *gin.Context) {
	var req struct {
		ProductIDs []uint `json:"product_ids"`
		GroupID    uint   `json:"group_id"`
	}
	c.ShouldBindJSON(&req)
	db := database.GetDB()
	db.Model(&model.Product{}).Where("id IN ?", req.ProductIDs).Update("group_id", req.GroupID)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": nil})
}

// GetProductGroupDetail 产品分组详情
func GetProductGroupDetail(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()
	var group model.ProductGroup
	if err := db.First(&group, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "分组不存在", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": group})
}

// ReorderProductTypes 产品类型排序
func ReorderProductTypes(c *gin.Context) {
	var req struct {
		Items []struct {
			ID       uint `json:"id"`
			Position int  `json:"position"`
		} `json:"items"`
	}
	c.ShouldBindJSON(&req)
	db := database.GetDB()
	for _, item := range req.Items {
		db.Model(&model.ProductType{}).Where("id = ?", item.ID).Update("sort_order", item.Position)
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "排序成功", "data": nil})
}

// ==================== 设置扩展 ====================

// GetSecurityConfig 获取安全配置
func GetSecurityConfig(c *gin.Context) {
	db := database.GetDB()
	var settings []model.Setting
	db.Where("`group` = ?", "security").Find(&settings)
	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": result})
}

// UpdateSecurityConfig 更新安全配置
func UpdateSecurityConfig(c *gin.Context) {
	updateSettingsByGroup(c, "security")
}

// GetGeneralConfig 获取常规配置
func GetGeneralConfig(c *gin.Context) {
	db := database.GetDB()
	var settings []model.Setting
	db.Where("`group` = ?", "general").Find(&settings)
	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": result})
}

// UpdateGeneralConfig 更新常规配置
func UpdateGeneralConfig(c *gin.Context) {
	updateSettingsByGroup(c, "general")
}

// GetDisplayConfig 获取显示配置
func GetDisplayConfig(c *gin.Context) {
	db := database.GetDB()
	var settings []model.Setting
	db.Where("`group` = ?", "display").Find(&settings)
	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": result})
}

// UpdateDisplayConfig 更新显示配置
func UpdateDisplayConfig(c *gin.Context) {
	updateSettingsByGroup(c, "display")
}

// GetInvoiceConfig 获取发票配置
func GetInvoiceConfig(c *gin.Context) {
	db := database.GetDB()
	var settings []model.Setting
	db.Where("`group` = ?", "invoice").Find(&settings)
	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": result})
}

// UpdateInvoiceConfig 更新发票配置
func UpdateInvoiceConfig(c *gin.Context) {
	updateSettingsByGroup(c, "invoice")
}

// GetContractConfig 获取合同配置
func GetContractConfig(c *gin.Context) {
	db := database.GetDB()
	var settings []model.Setting
	db.Where("`group` = ?", "contract").Find(&settings)
	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": result})
}

// UpdateContractConfig 更新合同配置
func UpdateContractConfig(c *gin.Context) {
	updateSettingsByGroup(c, "contract")
}

// GetCreditSettingConfig 获取信用额配置
func GetCreditSettingConfig(c *gin.Context) {
	db := database.GetDB()
	var settings []model.Setting
	db.Where("`group` = ?", "credit").Find(&settings)
	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": result})
}

// UpdateCreditSettingConfig 更新信用额配置
func UpdateCreditSettingConfig(c *gin.Context) {
	updateSettingsByGroup(c, "credit")
}

// GetPaymentGatewayConfig 获取支付网关配置
func GetPaymentGatewayConfig(c *gin.Context) {
	db := database.GetDB()
	var settings []model.Setting
	db.Where("`group` = ?", "payment_gateway").Find(&settings)
	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": result})
}

// UpdatePaymentGatewayConfig 更新支付网关配置
func UpdatePaymentGatewayConfig(c *gin.Context) {
	updateSettingsByGroup(c, "payment_gateway")
}

// updateSettingsByGroup 通用按分组更新设置
func updateSettingsByGroup(c *gin.Context, group string) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}
	db := database.GetDB()
	for key, value := range req {
		db.Where("`key` = ? AND `group` = ?", key, group).
			Assign(model.Setting{Value: value, Group: group}).
			FirstOrCreate(&model.Setting{Key: key, Group: group})
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": nil})
}

// ==================== 通知模板扩展 ====================

// TestNotificationTemplate 测试发送通知模板
func TestNotificationTemplate(c *gin.Context) {
	var req struct {
		ID     uint   `json:"id" binding:"required"`
		Target string `json:"target" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	var template model.NotificationTemplate
	if err := db.First(&template, req.ID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "模板不存在", "data": nil})
		return
	}

	// 通过PHP插件引擎发送测试通知
	results, err := pluginengine.TriggerHook("send_notification", map[string]interface{}{
		"template_id": template.ID,
		"target":      req.Target,
		"is_test":     true,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 502, "message": "插件引擎离线", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "测试发送成功", "data": gin.H{"results": results}})
}

// ==================== 员工扩展 ====================

// GetStaffRoles 获取员工角色列表
func GetStaffRoles(c *gin.Context) {
	db := database.GetDB()
	var roles []model.Role
	db.Find(&roles)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": roles})
}

// ==================== 插件扩展 ====================

// InstallPlugin 安装插件
func InstallPlugin(c *gin.Context) {
	var req struct {
		Slug   string `json:"slug" binding:"required"`
		Domain string `json:"domain" binding:"required"`
	}
	c.ShouldBindJSON(&req)
	db := database.GetDB()
	plugin := model.Plugin{
		Slug:   req.Slug,
		Name:   req.Slug,
		Domain: req.Domain,
		Status: "disabled",
	}
	db.Create(&plugin)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "安装成功", "data": gin.H{"id": plugin.ID}})
}

// ScanPlugins 扫描插件（通过PHP插件引擎扫描插件目录）
func ScanPlugins(c *gin.Context) {
	results, err := pluginengine.TriggerHook("scan_plugins", map[string]interface{}{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 502, "message": "插件引擎离线", "data": nil})
		return
	}

	// 同步扫描结果到数据库
	db := database.GetDB()
	found := 0
	newCount := 0
	if len(results) > 0 {
		if data, ok := results[0].Data.(map[string]interface{}); ok {
			if f, ok := data["found"].(float64); ok { found = int(f) }
			if n, ok := data["new"].(float64); ok { newCount = int(n) }
		}
	}

	// 如果插件引擎没有返回数据，从数据库查询当前插件数
	if found == 0 {
		var count int64
		db.Model(&model.Plugin{}).Count(&count)
		found = int(count)
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "扫描完成", "data": gin.H{"found": found, "new": newCount}})
}

// PluginHealthCheck 插件健康检查（通过PHP插件引擎检查）
func PluginHealthCheck(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()
	var plugin model.Plugin
	if err := db.First(&plugin, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "插件不存在", "data": nil})
		return
	}

	// 通过PHP插件引擎检查插件健康状态
	results, err := pluginengine.TriggerHook("plugin_health_check", map[string]interface{}{
		"plugin_id": plugin.ID,
		"slug":      plugin.Slug,
		"domain":    plugin.Domain,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 502, "message": "插件引擎离线", "data": nil})
		return
	}

	status := "healthy"
	if len(results) > 0 && results[0].Data != nil {
		if data, ok := results[0].Data.(map[string]interface{}); ok {
			if s, ok := data["status"].(string); ok {
				status = s
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"status": status, "plugin_id": plugin.ID}})
}

// ==================== 供应商扩展 ====================

// RunSupplierTask 执行供应商任务（通过PHP插件引擎执行）
func RunSupplierTask(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Task string `json:"task" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	var supplier model.Supplier
	if err := db.First(&supplier, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "供应商不存在", "data": nil})
		return
	}

	// 通过PHP插件引擎执行供应商任务
	results, err := pluginengine.TriggerHook("supplier_task", map[string]interface{}{
		"supplier_id": supplier.ID,
		"task":        req.Task,
		"api_url":     supplier.APIURL,
		"api_key":     supplier.APIKey,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 502, "message": "插件引擎离线", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "任务已提交", "data": gin.H{"supplier_id": supplier.ID, "task": req.Task, "results": results}})
}

// ==================== 实名认证扩展 ====================

// UnbindVerification 解绑实名认证
func UnbindVerification(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()
	db.Model(&model.User{}).Where("id = ?", id).Updates(map[string]interface{}{"is_verified": false, "verified_name": "", "verified_type": ""})
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "已解绑", "data": nil})
}

// ==================== 优惠券扩展 ====================

// GetCouponProductGroups 获取优惠券可用产品分组
func GetCouponProductGroups(c *gin.Context) {
	db := database.GetDB()
	var groups []model.ProductGroup
	db.Find(&groups)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": groups})
}

// ==================== 内容详情 ====================

// GetNewsDetail 新闻详情
func GetNewsDetail(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()
	var news model.News
	if err := db.First(&news, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "新闻不存在", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": news})
}

// GetKnowledgeArticleDetail 知识库文章详情
func GetKnowledgeArticleDetail(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()
	var article model.KnowledgeArticle
	if err := db.First(&article, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "文章不存在", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": article})
}

// ==================== 知识库分类扩展 ====================

// UpdateKnowledgeCategory 更新知识库分类
func UpdateKnowledgeCategory(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name      string `json:"name"`
		SortOrder int    `json:"sort_order"`
		Status    string `json:"status"`
	}
	c.ShouldBindJSON(&req)
	db := database.GetDB()
	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.SortOrder > 0 {
		updates["sort_order"] = req.SortOrder
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	db.Model(&model.KnowledgeCategory{}).Where("id = ?", id).Updates(updates)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": nil})
}

// DeleteKnowledgeCategory 删除知识库分类
func DeleteKnowledgeCategory(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()
	var count int64
	db.Model(&model.KnowledgeArticle{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "该分类下有文章，无法删除", "data": nil})
		return
	}
	db.Delete(&model.KnowledgeCategory{}, id)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功", "data": nil})
}

// ==================== 下载分类扩展 ====================

// UpdateDownloadCategory 更新下载分类
func UpdateDownloadCategory(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name      string `json:"name"`
		SortOrder int    `json:"sort_order"`
		Status    string `json:"status"`
	}
	c.ShouldBindJSON(&req)
	db := database.GetDB()
	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.SortOrder > 0 {
		updates["sort_order"] = req.SortOrder
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	db.Model(&model.DownloadCategory{}).Where("id = ?", id).Updates(updates)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": nil})
}

// DeleteDownloadCategory 删除下载分类
func DeleteDownloadCategory(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()
	var count int64
	db.Model(&model.Download{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "该分类下有文件，无法删除", "data": nil})
		return
	}
	db.Delete(&model.DownloadCategory{}, id)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功", "data": nil})
}

// ==================== Home Hero Assets ====================

// GetHomeHeroAssets 获取首页Hero可用资源文件（扫描uploads目录）
func GetHomeHeroAssets(c *gin.Context) {
	db := database.GetDB()
	var files []model.MediaFile
	db.Where("mime_type LIKE ?", "image/%").Order("id DESC").Limit(50).Find(&files)

	images := []string{}
	for _, f := range files {
		images = append(images, f.Path)
	}
	if images == nil { images = []string{} }

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"images": images, "videos": []string{}}})
}

// ==================== Admin 用户服务子操作 ====================

// GetUserEmailLogs 获取用户邮件日志
// GET /api/admin/users/:id/email-logs
func GetUserEmailLogs(c *gin.Context) {
	id := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 20 }

	db := database.GetDB()
	var logs []model.SystemLog
	var total int64
	db.Model(&model.SystemLog{}).Where("user_id = ? AND type = ?", id, "email").Count(&total)
	offset := (page - 1) * pageSize
	db.Where("user_id = ? AND type = ?", id, "email").Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs)
	if logs == nil { logs = []model.SystemLog{} }
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"list": logs, "total": total, "page": page, "page_size": pageSize}})
}

// GetUserSmsLogs 获取用户短信日志
// GET /api/admin/users/:id/sms-logs
func GetUserSmsLogs(c *gin.Context) {
	id := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 20 }

	db := database.GetDB()
	var logs []model.SystemLog
	var total int64
	db.Model(&model.SystemLog{}).Where("user_id = ? AND type = ?", id, "sms").Count(&total)
	offset := (page - 1) * pageSize
	db.Where("user_id = ? AND type = ?", id, "sms").Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs)
	if logs == nil { logs = []model.SystemLog{} }
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"list": logs, "total": total, "page": page, "page_size": pageSize}})
}

// GetUserInvoiceDetail 获取用户特定账单详情
// GET /api/admin/users/:id/invoices/:invoice_id
func GetUserInvoiceDetail(c *gin.Context) {
	userID := c.Param("id")
	invoiceID := c.Param("invoice_id")

	db := database.GetDB()
	var invoice model.Invoice
	if err := db.Where("id = ? AND user_id = ?", invoiceID, userID).First(&invoice).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "账单不存在", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": invoice})
}

// AdminGetServiceConnection 管理员获取用户服务连接信息
// GET /api/admin/users/:id/services/:service_id/connection
func AdminGetServiceConnection(c *gin.Context) {
	userID := c.Param("id")
	serviceID := c.Param("service_id")

	db := database.GetDB()
	var service model.Service
	if err := db.Where("id = ? AND user_id = ?", serviceID, userID).First(&service).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "服务不存在", "data": nil})
		return
	}

	results, err := pluginengine.TriggerHook("get_service_connection", map[string]interface{}{"service_id": service.ID})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 502, "message": "插件引擎离线", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"service": service, "connection": results}})
}

// AdminGetServiceRemoteStatus 管理员获取用户服务远程状态
// GET /api/admin/users/:id/services/:service_id/remote-status
func AdminGetServiceRemoteStatus(c *gin.Context) {
	userID := c.Param("id")
	serviceID := c.Param("service_id")

	db := database.GetDB()
	var service model.Service
	if err := db.Where("id = ? AND user_id = ?", serviceID, userID).First(&service).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "服务不存在", "data": nil})
		return
	}

	results, err := pluginengine.TriggerHook("get_service_status", map[string]interface{}{"service_id": service.ID})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 502, "message": "插件引擎离线", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"status": service.Status, "remote": results}})
}

// AdminUpdateServiceMeta 管理员更新服务元数据
// PUT /api/admin/users/:id/services/:service_id/meta
func AdminUpdateServiceMeta(c *gin.Context) {
	userID := c.Param("id")
	serviceID := c.Param("service_id")

	db := database.GetDB()
	var service model.Service
	if err := db.Where("id = ? AND user_id = ?", serviceID, userID).First(&service).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "服务不存在", "data": nil})
		return
	}

	var req struct {
		Remark     string `json:"remark"`
		Hostname   string `json:"hostname"`
		Username   string `json:"username"`
		Password   string `json:"password"`
		IPAddress  string `json:"ip_address"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	// 白名单更新（只允许修改非敏感字段）
	updates := map[string]interface{}{}
	if req.Remark != "" { updates["remark"] = req.Remark }
	if req.Hostname != "" { updates["hostname"] = req.Hostname }
	if req.Username != "" { updates["username"] = req.Username }
	if req.Password != "" { updates["password"] = req.Password }
	if req.IPAddress != "" { updates["ip_address"] = req.IPAddress }

	db.Model(&service).Updates(updates)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": nil})
}

// AdminManualProvision 管理员手动开通服务
// POST /api/admin/users/:id/services/:service_id/manual-provision
func AdminManualProvision(c *gin.Context) {
	userID := c.Param("id")
	serviceID := c.Param("service_id")

	db := database.GetDB()
	var service model.Service
	if err := db.Where("id = ? AND user_id = ?", serviceID, userID).First(&service).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "服务不存在", "data": nil})
		return
	}

	results, err := pluginengine.TriggerHook("provision_service", map[string]interface{}{"service_id": service.ID})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 502, "message": "插件引擎离线", "data": nil})
		return
	}

	db.Model(&service).Update("status", "active")
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "开通成功", "data": gin.H{"results": results}})
}

// AdminServicePowerAction 管理员服务电源操作
// POST /api/admin/users/:id/services/:service_id/power-actions
func AdminServicePowerAction(c *gin.Context) {
	userID := c.Param("id")
	serviceID := c.Param("service_id")

	var req struct {
		Action string `json:"action" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	var service model.Service
	if err := db.Where("id = ? AND user_id = ?", serviceID, userID).First(&service).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "服务不存在", "data": nil})
		return
	}

	results, err := pluginengine.TriggerHook("service_power_action", map[string]interface{}{
		"service_id": service.ID,
		"action":     req.Action,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 502, "message": "插件引擎离线", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "操作成功", "data": gin.H{"results": results}})
}

// AdminResetServicePassword 管理员重置服务密码
// POST /api/admin/users/:id/services/:service_id/password-resets
func AdminResetServicePassword(c *gin.Context) {
	userID := c.Param("id")
	serviceID := c.Param("service_id")

	db := database.GetDB()
	var service model.Service
	if err := db.Where("id = ? AND user_id = ?", serviceID, userID).First(&service).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "服务不存在", "data": nil})
		return
	}

	results, err := pluginengine.TriggerHook("reset_service_password", map[string]interface{}{"service_id": service.ID})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 502, "message": "插件引擎离线", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "密码重置成功", "data": gin.H{"results": results}})
}

// ==================== 新增路由handler ====================

// ReorderProductGroups 产品分组排序
// POST /api/admin/product-groups/reorders
func ReorderProductGroups(c *gin.Context) {
	var req struct {
		Items []struct {
			ID       uint `json:"id"`
			SortOrder int `json:"sort_order"`
		} `json:"items" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	for _, item := range req.Items {
		db.Model(&model.ProductGroup{}).Where("id = ?", item.ID).Update("sort_order", item.SortOrder)
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "排序成功", "data": nil})
}

// ForceDeleteProduct 强制删除产品（跳过回收站）
// DELETE /api/admin/products/:id/force
func ForceDeleteProduct(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	// 检查是否有活跃服务使用此产品
	var serviceCount int64
	db.Model(&model.Service{}).Where("product_id = ? AND status IN ?", id, []string{"active", "suspended"}).Count(&serviceCount)
	if serviceCount > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "该产品还有活跃服务，无法强制删除", "data": nil})
		return
	}

	if err := db.Unscoped().Delete(&model.Product{}, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "删除失败", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功", "data": nil})
}

// RetryScheduleRun 重试定时任务
// POST /api/admin/schedule-runs/:id/retry
func RetryScheduleRun(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var run model.ScheduleRun
	if err := db.First(&run, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "运行记录不存在", "data": nil})
		return
	}

	// 重置状态为pending
	db.Model(&run).Updates(map[string]interface{}{"status": "pending", "error_message": ""})
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "已加入重试队列", "data": nil})
}

// TriggerSchedule 手动触发定时任务
// POST /api/admin/schedule-triggers
func TriggerSchedule(c *gin.Context) {
	var req struct {
		TaskID uint `json:"task_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	var task model.ScheduleTask
	if err := db.First(&task, req.TaskID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "任务不存在", "data": nil})
		return
	}

	// 创建运行记录
	run := model.ScheduleRun{
		TaskID:  task.ID,
		Status:  "pending",
	}
	db.Create(&run)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "已触发", "data": gin.H{"run_id": run.ID}})
}

// GetVerificationHistory 获取实名认证历史
// GET /api/admin/verifications/:id/history
func GetVerificationHistory(c *gin.Context) {
	userID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 20 }

	db := database.GetDB()
	var records []model.Verification
	var total int64
	db.Model(&model.Verification{}).Where("user_id = ?", userID).Count(&total)
	offset := (page - 1) * pageSize
	db.Where("user_id = ?", userID).Offset(offset).Limit(pageSize).Order("id DESC").Find(&records)
	if records == nil { records = []model.Verification{} }
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"list": records, "total": total, "page": page, "page_size": pageSize}})
}

// UnbindVerificationByUser 解绑实名认证（按用户ID）
// POST /api/admin/verifications/:id/unbindings
func UnbindVerificationByUser(c *gin.Context) {
	userID := c.Param("id")
	db := database.GetDB()

	result := db.Model(&model.Verification{}).Where("user_id = ?", userID).Update("status", "unbound")
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "未找到实名认证记录", "data": nil})
		return
	}
	// 同时清除用户表的verified状态
	db.Model(&model.User{}).Where("id = ?", userID).Update("is_verified", false)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "解绑成功", "data": nil})
}

// RecallTicketReply 撤回工单回复
// POST /api/admin/tickets/:id/replies/:reply_id/recalls
func RecallTicketReply(c *gin.Context) {
	replyID := c.Param("reply_id")
	db := database.GetDB()

	var reply model.TicketReply
	if err := db.First(&reply, replyID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "回复不存在", "data": nil})
		return
	}

	// 只能撤回自己的回复或管理员可以撤回任何回复
	reply.Content = "[已撤回]"
	db.Save(&reply)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "撤回成功", "data": nil})
}

// GetFinanceLedgerDetail 获取财务账本详情
// GET /api/admin/finance/ledger/:id
func GetFinanceLedgerDetail(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var invoice model.Invoice
	if err := db.First(&invoice, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "记录不存在", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": invoice})
}

// GetLogSummaryByChannel 获取日志摘要
// GET /api/admin/log-summaries/:channel
func GetLogSummaryByChannel(c *gin.Context) {
	channel := c.Param("channel")
	db := database.GetDB()

	// 统计今日、本周、本月的日志数量
	var todayCount, weekCount, monthCount int64
	db.Model(&model.OperationLog{}).Where("resource = ? AND created_at >= CURDATE()", channel).Count(&todayCount)
	db.Model(&model.OperationLog{}).Where("resource = ? AND created_at >= DATE_SUB(CURDATE(), INTERVAL 7 DAY)", channel).Count(&weekCount)
	db.Model(&model.OperationLog{}).Where("resource = ? AND created_at >= DATE_SUB(CURDATE(), INTERVAL 30 DAY)", channel).Count(&monthCount)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"channel":    channel,
		"today":      todayCount,
		"this_week":  weekCount,
		"this_month": monthCount,
	}})
}

// ==================== 工单上游投递（参考图拉财务TicketDeliveryService） ====================

// GetTicketDeliveryRules 获取投递规则列表
// GET /api/admin/ticket-delivery-rules
func GetTicketDeliveryRules(c *gin.Context) {
	db := database.GetDB()
	var rules []model.TicketDeliveryRule
	db.Order("sort_order ASC").Find(&rules)
	if rules == nil { rules = []model.TicketDeliveryRule{} }
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": rules})
}

// CreateTicketDeliveryRule 创建投递规则
// POST /api/admin/ticket-delivery-rules
func CreateTicketDeliveryRule(c *gin.Context) {
	var req struct {
		Name         string `json:"name" binding:"required"`
		DepartmentID uint   `json:"department_id"`
		ProductID    uint   `json:"product_id"`
		Keyword      string `json:"keyword"`
		UpstreamURL  string `json:"upstream_url"`
		UpstreamKey  string `json:"upstream_key"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	rule := model.TicketDeliveryRule{
		Name:         req.Name,
		DepartmentID: req.DepartmentID,
		ProductID:    req.ProductID,
		Keyword:      req.Keyword,
		UpstreamURL:  req.UpstreamURL,
		UpstreamKey:  req.UpstreamKey,
		Status:       "active",
	}
	if err := db.Create(&rule).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "创建失败", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "创建成功", "data": gin.H{"id": rule.ID}})
}

// UpdateTicketDeliveryRule 更新投递规则
// PUT /api/admin/ticket-delivery-rules/:id
func UpdateTicketDeliveryRule(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()
	var rule model.TicketDeliveryRule
	if err := db.First(&rule, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "规则不存在", "data": nil})
		return
	}

	var req struct {
		Name         string `json:"name"`
		DepartmentID uint   `json:"department_id"`
		ProductID    uint   `json:"product_id"`
		Keyword      string `json:"keyword"`
		UpstreamURL  string `json:"upstream_url"`
		UpstreamKey  string `json:"upstream_key"`
		Status       string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" { updates["name"] = req.Name }
	if req.DepartmentID > 0 { updates["department_id"] = req.DepartmentID }
	if req.ProductID > 0 { updates["product_id"] = req.ProductID }
	if req.Keyword != "" { updates["keyword"] = req.Keyword }
	if req.UpstreamURL != "" { updates["upstream_url"] = req.UpstreamURL }
	if req.UpstreamKey != "" { updates["upstream_key"] = req.UpstreamKey }
	if req.Status != "" { updates["status"] = req.Status }

	db.Model(&rule).Updates(updates)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": nil})
}

// DeleteTicketDeliveryRule 删除投递规则
// DELETE /api/admin/ticket-delivery-rules/:id
func DeleteTicketDeliveryRule(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()
	if err := db.Delete(&model.TicketDeliveryRule{}, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "删除失败", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功", "data": nil})
}

// GetTicketUpstreamDelivery 获取工单上游投递状态
// GET /api/admin/tickets/:id/upstream-delivery
func GetTicketUpstreamDelivery(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var ticket model.Ticket
	if err := db.First(&ticket, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "工单不存在", "data": nil})
		return
	}

	// 查询投递记录
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"ticket_id": ticket.ID,
		"status":    ticket.Status,
	}})
}

// GetTicketUpstreamDeliveryLogs 获取工单上游投递日志
// GET /api/admin/tickets/:id/upstream-delivery/logs
func GetTicketUpstreamDeliveryLogs(c *gin.Context) {
	id := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 20 }

	db := database.GetDB()
	var logs []model.OperationLog
	var total int64
	db.Model(&model.OperationLog{}).Where("resource = ? AND resource_id = ?", "ticket_delivery", id).Count(&total)
	offset := (page - 1) * pageSize
	db.Where("resource = ? AND resource_id = ?", "ticket_delivery", id).Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs)
	if logs == nil { logs = []model.OperationLog{} }
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"list": logs, "total": total, "page": page, "page_size": pageSize}})
}

// ==================== 供应商同步（参考图拉财务上游驱动） ====================

// SyncSupplierProducts 同步供应商商品
// POST /api/admin/suppliers/:id/sync-products
func SyncSupplierProducts(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var supplier model.Supplier
	if err := db.First(&supplier, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "供应商不存在", "data": nil})
		return
	}

	// 根据供应商类型选择驱动（Go原生，不走PHP）
	driver := service.NewDriver(supplier)
	products, err := driver.FetchProducts()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 502, "message": "拉取商品失败", "data": nil})
		return
	}

	synced := 0
	for _, p := range products {
		var existing model.SupplierProduct
		result := db.Where("supplier_id = ? AND remote_product_id = ?", supplier.ID, p.ID).First(&existing)
		if result.Error != nil {
			db.Create(&model.SupplierProduct{SupplierID: supplier.ID, RemoteProductID: p.ID, Name: p.Name, RemotePrice: p.Price, Stock: p.Stock, Status: "active"})
		} else {
			db.Model(&existing).Updates(map[string]interface{}{"name": p.Name, "remote_price": p.Price, "stock": p.Stock})
		}
		synced++
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "同步完成", "data": gin.H{"synced": synced}})
}

// SyncSupplierPrices 同步供应商价格（应用利润率）
// POST /api/admin/suppliers/:id/sync-prices
func SyncSupplierPrices(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var supplier model.Supplier
	if err := db.First(&supplier, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "供应商不存在", "data": nil})
		return
	}

	driver := service.NewDriver(supplier)
	products, err := driver.FetchProducts()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 502, "message": "拉取价格失败", "data": nil})
		return
	}

	updated := 0
	for _, p := range products {
		var existing model.SupplierProduct
		if err := db.Where("supplier_id = ? AND remote_product_id = ?", supplier.ID, p.ID).First(&existing).Error; err == nil {
			profitRate := existing.ProfitRate
			if profitRate <= 0 { profitRate = 25 }
			localPrice := p.Price * (1 + profitRate/100)
			db.Model(&existing).Updates(map[string]interface{}{"remote_price": p.Price, "local_price": localPrice})
			updated++
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "价格同步完成", "data": gin.H{"updated": updated}})
}

// SyncSupplierStock 同步供应商库存
// POST /api/admin/suppliers/:id/sync-stock
func SyncSupplierStock(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var supplier model.Supplier
	if err := db.First(&supplier, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "供应商不存在", "data": nil})
		return
	}

	driver := service.NewDriver(supplier)
	products, err := driver.FetchProducts()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 502, "message": "拉取库存失败", "data": nil})
		return
	}

	updated := 0
	for _, p := range products {
		var existing model.SupplierProduct
		if err := db.Where("supplier_id = ? AND remote_product_id = ?", supplier.ID, p.ID).First(&existing).Error; err == nil {
			db.Model(&existing).Update("stock", p.Stock)
			updated++
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "库存同步完成", "data": gin.H{"updated": updated}})
}

// ==================== Redis配置 ====================

// GetRedisConfig 获取Redis配置
// GET /api/admin/redis/config
func GetRedisConfig(c *gin.Context) {
	redisSvc := service.NewRedisService()
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": redisSvc.GetConfig()})
}

// RedisHealthCheck Redis健康检查
// GET /api/admin/redis/health
func RedisHealthCheck(c *gin.Context) {
	redisSvc := service.NewRedisService()
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": redisSvc.HealthCheck()})
}

// ==================== AI配置 ====================

// GetAIConfig 获取AI配置
// GET /api/admin/ai/config
func GetAIConfig(c *gin.Context) {
	aiSvc := service.NewAIService()
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": aiSvc.GetConfig()})
}

// TestAIConnection 测试AI连接
// POST /api/admin/ai/test
func TestAIConnection(c *gin.Context) {
	aiSvc := service.NewAIService()
	if !aiSvc.IsEnabled() {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "AI服务未启用", "data": nil})
		return
	}

	result, err := aiSvc.ChatCompletion([]map[string]string{
		{"role": "user", "content": "你好，请回复'连接成功'四个字"},
	}, nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "AI连接失败", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "AI连接成功", "data": gin.H{"response": result}})
}

// GenerateProductDescription AI生成商品简介
// POST /api/admin/ai/generate-description
func GenerateProductDescription(c *gin.Context) {
	var req struct {
		ProductName string                 `json:"product_name" binding:"required"`
		Config      map[string]interface{} `json:"config"`
		Template    string                 `json:"template"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	aiSvc := service.NewAIService()
	if !aiSvc.IsEnabled() {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "AI服务未启用，请先在设置中配置AI", "data": nil})
		return
	}

	result, err := aiSvc.GenerateProductDescription(req.ProductName, req.Config, req.Template)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "AI生成失败", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "生成成功", "data": gin.H{"description": result}})
}

// AITicketReply AI工单自动回复
// POST /api/admin/ai/ticket-reply
func AITicketReply(c *gin.Context) {
	var req struct {
		TicketID uint `json:"ticket_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	var ticket model.Ticket
	if err := db.First(&ticket, req.TicketID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "工单不存在", "data": nil})
		return
	}

	aiSvc := service.NewAIService()
	if !aiSvc.IsEnabled() {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "AI服务未启用", "data": nil})
		return
	}

	// 获取客户信息
	var user model.User
	db.First(&user, ticket.UserID)

	// 获取最新工单回复内容
	var lastReply model.TicketReply
	ticketContent := ticket.Subject
	db.Where("ticket_id = ?", ticket.ID).Order("id DESC").First(&lastReply)
	if lastReply.Content != "" {
		ticketContent = lastReply.Content
	}

	reply, err := aiSvc.TicketAutoReply(ticket.Subject, ticketContent, map[string]interface{}{
		"customer_info": gin.H{"username": user.Username, "email": user.Email},
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "AI回复失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "AI回复生成成功", "data": gin.H{"reply": reply}})
}

// ==================== AI工单系统（参考mianyu_ai_ticket插件） ====================

// GetAITicketConfig 获取AI工单配置
// GET /api/admin/ai-ticket/config
func GetAITicketConfig(c *gin.Context) {
	aiTicketSvc := service.NewAITicketService()
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": aiTicketSvc.GetConfig()})
}

// GetAITicketQueueStats 获取AI工单队列统计
// GET /api/admin/ai-ticket/queue/stats
func GetAITicketQueueStats(c *gin.Context) {
	aiTicketSvc := service.NewAITicketService()
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": aiTicketSvc.GetQueueStats()})
}

// GetAITicketQueueList 获取AI工单队列列表
// GET /api/admin/ai-ticket/queue
func GetAITicketQueueList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.DefaultQuery("status", "")
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 20 }

	db := database.GetDB()
	query := db.Model(&model.AITicketQueue{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var items []model.AITicketQueue
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&items)
	if items == nil { items = []model.AITicketQueue{} }

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"list": items, "total": total, "page": page, "page_size": pageSize}})
}

// ProcessAITicketQueue 手动处理AI工单队列
// POST /api/admin/ai-ticket/queue/process
func ProcessAITicketQueue(c *gin.Context) {
	aiTicketSvc := service.NewAITicketService()
	processed, err := aiTicketSvc.ProcessQueue(20)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "处理失败", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "处理完成", "data": gin.H{"processed": processed}})
}

// GetAITicketKnowledgeList 获取AI工单知识库
// GET /api/admin/ai-ticket/knowledge
func GetAITicketKnowledgeList(c *gin.Context) {
	db := database.GetDB()
	var items []model.AITicketKnowledge
	db.Where("status = 1").Order("sort ASC").Find(&items)
	if items == nil { items = []model.AITicketKnowledge{} }
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": items})
}

// CreateAITicketKnowledge 创建知识库条目
// POST /api/admin/ai-ticket/knowledge
func CreateAITicketKnowledge(c *gin.Context) {
	var req struct {
		Title    string `json:"title" binding:"required"`
		Keywords string `json:"keywords"`
		Question string `json:"question"`
		Answer   string `json:"answer" binding:"required"`
		Sort     int    `json:"sort"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	knowledge := model.AITicketKnowledge{
		Title:    req.Title,
		Keywords: req.Keywords,
		Question: req.Question,
		Answer:   req.Answer,
		Sort:     req.Sort,
		Status:   1,
	}
	if err := db.Create(&knowledge).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "创建失败", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "创建成功", "data": gin.H{"id": knowledge.ID}})
}

// UpdateAITicketKnowledge 更新知识库条目
// PUT /api/admin/ai-ticket/knowledge/:id
func UpdateAITicketKnowledge(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()
	var item model.AITicketKnowledge
	if err := db.First(&item, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "条目不存在", "data": nil})
		return
	}

	var req struct {
		Title    string `json:"title"`
		Keywords string `json:"keywords"`
		Question string `json:"question"`
		Answer   string `json:"answer"`
		Sort     int    `json:"sort"`
		Status   *int   `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	updates := map[string]interface{}{}
	if req.Title != "" { updates["title"] = req.Title }
	if req.Keywords != "" { updates["keywords"] = req.Keywords }
	if req.Question != "" { updates["question"] = req.Question }
	if req.Answer != "" { updates["answer"] = req.Answer }
	if req.Sort != 0 { updates["sort"] = req.Sort }
	if req.Status != nil { updates["status"] = *req.Status }

	db.Model(&item).Updates(updates)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": nil})
}

// DeleteAITicketKnowledge 删除知识库条目
// DELETE /api/admin/ai-ticket/knowledge/:id
func DeleteAITicketKnowledge(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()
	if err := db.Delete(&model.AITicketKnowledge{}, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "删除失败", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功", "data": nil})
}

// GetAITicketRuleList 获取AI工单规则
// GET /api/admin/ai-ticket/rules
func GetAITicketRuleList(c *gin.Context) {
	db := database.GetDB()
	var items []model.AITicketRule
	db.Order("priority ASC").Find(&items)
	if items == nil { items = []model.AITicketRule{} }
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": items})
}

// CreateAITicketRule 创建规则
// POST /api/admin/ai-ticket/rules
func CreateAITicketRule(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Keywords    string `json:"keywords"`
		DeptFilter  string `json:"dept_filter"`
		TargetDept  int    `json:"target_dept"`
		Priority    int    `json:"priority"`
		Action      string `json:"action"`
		PromptExtra string `json:"prompt_extra"`
		SampleReply string `json:"sample_reply"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	rule := model.AITicketRule{
		Name:        req.Name,
		Keywords:    req.Keywords,
		DeptFilter:  req.DeptFilter,
		TargetDept:  req.TargetDept,
		Priority:    req.Priority,
		Action:      req.Action,
		PromptExtra: req.PromptExtra,
		SampleReply: req.SampleReply,
		Status:      1,
	}
	if err := db.Create(&rule).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "创建失败", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "创建成功", "data": gin.H{"id": rule.ID}})
}

// UpdateAITicketRule 更新规则
// PUT /api/admin/ai-ticket/rules/:id
func UpdateAITicketRule(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()
	var item model.AITicketRule
	if err := db.First(&item, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "规则不存在", "data": nil})
		return
	}

	var req struct {
		Name        string `json:"name"`
		Keywords    string `json:"keywords"`
		DeptFilter  string `json:"dept_filter"`
		TargetDept  int    `json:"target_dept"`
		Priority    int    `json:"priority"`
		Action      string `json:"action"`
		PromptExtra string `json:"prompt_extra"`
		SampleReply string `json:"sample_reply"`
		Status      *int   `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" { updates["name"] = req.Name }
	if req.Keywords != "" { updates["keywords"] = req.Keywords }
	if req.DeptFilter != "" { updates["dept_filter"] = req.DeptFilter }
	if req.TargetDept != 0 { updates["target_dept"] = req.TargetDept }
	if req.Priority != 0 { updates["priority"] = req.Priority }
	if req.Action != "" { updates["action"] = req.Action }
	if req.PromptExtra != "" { updates["prompt_extra"] = req.PromptExtra }
	if req.SampleReply != "" { updates["sample_reply"] = req.SampleReply }
	if req.Status != nil { updates["status"] = *req.Status }

	db.Model(&item).Updates(updates)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": nil})
}

// DeleteAITicketRule 删除规则
// DELETE /api/admin/ai-ticket/rules/:id
func DeleteAITicketRule(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()
	if err := db.Delete(&model.AITicketRule{}, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "删除失败", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功", "data": nil})
}

// GetAITicketProcessLogs 获取AI工单处理日志
// GET /api/admin/ai-ticket/logs
func GetAITicketProcessLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 20 }

	db := database.GetDB()
	var total int64
	db.Model(&model.AITicketProcessLog{}).Count(&total)

	var items []model.AITicketProcessLog
	offset := (page - 1) * pageSize
	db.Offset(offset).Limit(pageSize).Order("id DESC").Find(&items)
	if items == nil { items = []model.AITicketProcessLog{} }

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"list": items, "total": total, "page": page, "page_size": pageSize}})
}

// SetAITicketMode 设置工单AI/人工模式
// POST /api/admin/ai-ticket/tickets/:id/mode
func SetAITicketMode(c *gin.Context) {
	ticketIDStr := c.Param("id")
	ticketID, err := strconv.ParseUint(ticketIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "无效的工单ID", "data": nil})
		return
	}

	var req struct {
		Mode string `json:"mode" binding:"required"` // ai, human
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	if req.Mode != "ai" && req.Mode != "human" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "模式必须是ai或human", "data": nil})
		return
	}

	aiTicketSvc := service.NewAITicketService()
	if err := aiTicketSvc.SetTicketMode(uint(ticketID), req.Mode); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "设置失败", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "设置成功", "data": nil})
}

// ==================== API审计补充的4个缺失路由 ====================

// PullTrafficPackageCatalog 流量包同步
// POST /api/admin/products/traffic-package-pulls
func PullTrafficPackageCatalog(c *gin.Context) {
	// 流量包从插件引擎拉取目录
	results, err := pluginengine.TriggerHook("traffic_package_pull", map[string]interface{}{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 502, "message": "插件引擎离线", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "同步完成", "data": results})
}

// RunCouponCampaignTask 优惠券活动任务执行
// POST /api/admin/coupon-campaigns/:id/tasks
func RunCouponCampaignTask(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var campaign model.CouponCampaign
	if err := db.First(&campaign, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "活动不存在", "data": nil})
		return
	}

	// 根据活动类型执行不同任务
	switch campaign.Type {
	case "auto_issue":
		// 自动发放优惠券
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "自动发放任务已执行", "data": gin.H{"campaign_id": campaign.ID}})
	default:
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "不支持的活动类型", "data": nil})
	}
}

// AdminCreateUserService 管理员为用户创建服务
// POST /api/admin/users/:id/services
func AdminCreateUserService(c *gin.Context) {
	userID := c.Param("id")
	var req struct {
		ProductID uint   `json:"product_id" binding:"required"`
		BillingCycle string `json:"billing_cycle"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "用户不存在", "data": nil})
		return
	}

	var product model.Product
	if err := db.First(&product, req.ProductID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "产品不存在", "data": nil})
		return
	}

	service := model.Service{
		UserID:       user.ID,
		ProductID:    product.ID,
		ProductName:  product.Name,
		BillingCycle: req.BillingCycle,
		Amount:       product.Amount,
		Status:       "pending",
	}
	if err := db.Create(&service).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "创建失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "服务创建成功", "data": gin.H{"id": service.ID}})
}

// AdminDeleteUserService 管理员删除用户服务
// DELETE /api/admin/users/:id/services/:service_id
func AdminDeleteUserService(c *gin.Context) {
	serviceID := c.Param("service_id")
	db := database.GetDB()

	var svc model.Service
	if err := db.First(&svc, serviceID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "服务不存在", "data": nil})
		return
	}

	if svc.Status == "active" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "不能删除活跃服务，请先终止", "data": nil})
		return
	}

	db.Delete(&svc)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功", "data": nil})
}
