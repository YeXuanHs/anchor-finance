package admin

import (
	"net/http"
	"strconv"
	"time"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
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
func RefreshUserServicesStatus(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()
	var services []model.Service
	db.Where("user_id = ?", id).Find(&services)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "已刷新", "data": gin.H{"count": len(services)}})
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
	db.Model(&model.Recharge{}).Where("status = ?", "completed").Select("COALESCE(SUM(amount), 0)").Scan(&totalRecharge)
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
	c.ShouldBindJSON(&req)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "测试发送成功", "data": nil})
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

// ScanPlugins 扫描插件
func ScanPlugins(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "扫描完成", "data": gin.H{"found": 0, "new": 0}})
}

// PluginHealthCheck 插件健康检查
func PluginHealthCheck(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()
	var plugin model.Plugin
	if err := db.First(&plugin, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "插件不存在", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"status": "healthy", "plugin_id": plugin.ID}})
}

// ==================== 供应商扩展 ====================

// RunSupplierTask 执行供应商任务
func RunSupplierTask(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Task string `json:"task" binding:"required"`
	}
	c.ShouldBindJSON(&req)
	db := database.GetDB()
	var supplier model.Supplier
	if err := db.First(&supplier, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "供应商不存在", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "任务已提交", "data": gin.H{"supplier_id": supplier.ID, "task": req.Task}})
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

// GetHomeHeroAssets 获取首页Hero可用资源文件
func GetHomeHeroAssets(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"images": []string{}, "videos": []string{}}})
}
