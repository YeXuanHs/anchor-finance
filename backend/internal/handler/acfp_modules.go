package handler

import (
	"encoding/json"
	"strconv"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type ACFPHandler struct {
	svc *service.ACFPService
	log *logger.Logger
}

func NewACFPHandler(svc *service.ACFPService, log *logger.Logger) *ACFPHandler {
	return &ACFPHandler{svc: svc, log: log}
}

// ─── IP历史 ───

func (h *ACFPHandler) GetIPHistory(c *gin.Context) {
	hostID, _ := strconv.ParseUint(c.Param("host_id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 20 }
	items, total, err := h.svc.GetIPHistory(uint(hostID), page, pageSize)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}
	response.Success(c, gin.H{"items": items, "total": total, "page": page})
}

// ─── 限量发售 ───

func (h *ACFPHandler) ListLimitedSales(c *gin.Context) {
	items, err := h.svc.ListLimitedSales()
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}
	response.Success(c, items)
}

func (h *ACFPHandler) SetLimitedSale(c *gin.Context) {
	var item model.ACFPLimitedSale
	if err := c.ShouldBindJSON(&item); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.svc.SetLimitedSale(&item); err != nil {
		response.ServerError(c, "保存失败")
		return
	}
	response.SuccessMsg(c, "保存成功")
}

func (h *ACFPHandler) CheckStock(c *gin.Context) {
	productID, _ := strconv.ParseUint(c.Param("product_id"), 10, 64)
	qty, _ := strconv.Atoi(c.DefaultQuery("qty", "1"))
	ok, remaining, _ := h.svc.CheckStock(uint(productID), qty)
	response.Success(c, gin.H{"available": ok, "remaining": remaining})
}

// ─── 价格锁定 ───

func (h *ACFPHandler) ListPriceLocks(c *gin.Context) {
	items, err := h.svc.ListPriceLocks()
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}
	response.Success(c, items)
}

func (h *ACFPHandler) SetPriceLock(c *gin.Context) {
	var item model.ACFPPriceLock
	if err := c.ShouldBindJSON(&item); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.svc.SetPriceLock(&item); err != nil {
		response.ServerError(c, "保存失败")
		return
	}
	response.SuccessMsg(c, "保存成功")
}

func (h *ACFPHandler) DeletePriceLock(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeletePriceLock(uint(id)); err != nil {
		response.ServerError(c, "删除失败")
		return
	}
	response.SuccessMsg(c, "删除成功")
}

// ─── 操作日志 ───

func (h *ACFPHandler) ListLogs(c *gin.Context) {
	module := c.Query("module")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 20 }
	items, total, err := h.svc.ListLogs(module, page, pageSize)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}
	response.Success(c, gin.H{"items": items, "total": total, "page": page})
}

func (h *ACFPHandler) CleanLogs(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "90"))
	if days < 1 { days = 90 }
	count := h.svc.CleanLogs(days)
	response.Success(c, gin.H{"cleaned": count})
}

// ─── 定时任务状态 ───

func (h *ACFPHandler) GetCronStatuses(c *gin.Context) {
	items, err := h.svc.GetCronStatuses()
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}
	response.Success(c, items)
}

// ─── 实名认证Pro ───

func (h *ACFPHandler) GetCertProConfig(c *gin.Context) {
	cfg := h.svc.GetCertProConfig()
	response.Success(c, cfg)
}

func (h *ACFPHandler) SetCertProConfig(c *gin.Context) {
	var req struct {
		MinAge int `json:"min_age"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if req.MinAge < 1 || req.MinAge > 100 {
		response.BadRequest(c, "年龄范围无效")
		return
	}
	if err := h.svc.SetCertProConfig(req.MinAge); err != nil {
		response.ServerError(c, "保存失败")
		return
	}
	response.SuccessMsg(c, "保存成功")
}

// SaveCertProConfig 保存配置（别名）
func (h *ACFPHandler) SaveCertProConfig(c *gin.Context) {
	h.SetCertProConfig(c)
}

// GetCertReviewList 获取认证审核列表
func (h *ACFPHandler) GetCertReviewList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.DefaultQuery("status", "")

	type CertReview struct {
		ID        uint   `json:"id"`
		UserID    uint   `json:"user_id"`
		Username  string `json:"username"`
		RealName  string `json:"real_name"`
		IDCard    string `json:"id_card"`
		Status    int    `json:"status"`
		CreatedAt string `json:"created_at"`
	}

	var items []CertReview
	var total int64

	q := h.svc.GetDB().Table("certifications c").
		Select("c.id, c.user_id, u.username, c.real_name, c.idcard as id_card, c.status, c.created_at").
		Joins("LEFT JOIN users u ON c.user_id = u.id")

	if status != "" {
		q = q.Where("c.status = ?", status)
	}

	q.Count(&total)
	q.Order("c.id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Scan(&items)

	response.Success(c, gin.H{
		"list":  items,
		"total": total,
		"page":  page,
	})
}

// ReviewCert 审核认证
func (h *ACFPHandler) ReviewCert(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Status int    `json:"status"` // 1=通过 2=拒绝
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if req.Status != 1 && req.Status != 2 {
		response.BadRequest(c, "状态无效")
		return
	}

	if err := h.svc.GetDB().Exec("UPDATE certifications SET status = ?, reject_reason = ? WHERE id = ?", req.Status, req.Reason, id).Error; err != nil {
		response.ServerError(c, "审核失败")
		return
	}
	response.SuccessMsg(c, "审核成功")
}

func (h *ACFPHandler) ListMinorCerts(c *gin.Context) {
	items, err := h.svc.ListMinorCerts()
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}
	response.Success(c, items)
}

func (h *ACFPHandler) ScanMinorCerts(c *gin.Context) {
	count, err := h.svc.ScanMinorCerts()
	if err != nil {
		response.ServerError(c, "扫描失败")
		return
	}
	response.Success(c, gin.H{"found": count})
}

func (h *ACFPHandler) RejectUnderageSubmissions(c *gin.Context) {
	count, err := h.svc.RejectUnderageSubmissions()
	if err != nil {
		response.ServerError(c, "拒绝失败")
		return
	}
	response.Success(c, gin.H{"rejected": count})
}

// ─── 缓存预热 ───

func (h *ACFPHandler) WarmCache(c *gin.Context) {
	count, err := h.svc.WarmProductCache()
	if err != nil {
		response.ServerError(c, "预热失败")
		return
	}
	response.Success(c, gin.H{"warmed": count})
}

// ─── 批量商品修改 ───

func (h *ACFPHandler) CreateBatchTask(c *gin.Context) {
	var task model.ACFPBatchTask
	if err := c.ShouldBindJSON(&task); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.svc.CreateBatchTask(&task); err != nil {
		response.ServerError(c, "创建失败")
		return
	}
	response.Success(c, task)
}

func (h *ACFPHandler) ListBatchTasks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 20 }
	items, total, err := h.svc.ListBatchTasks(page, pageSize)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}
	response.Success(c, gin.H{"items": items, "total": total, "page": page})
}

func (h *ACFPHandler) ExecuteBatchTask(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.ExecuteBatchTask(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "执行中")
}

// ─── 状态对账 ───

func (h *ACFPHandler) RunStatusSync(c *gin.Context) {
	count, err := h.svc.RunStatusSync()
	if err != nil {
		response.ServerError(c, "对账失败")
		return
	}
	response.Success(c, gin.H{"synced": count})
}

// ─── 通知去重 ───

func (h *ACFPHandler) GetNotifyStats(c *gin.Context) {
	response.Success(c, gin.H{"message": "通知去重运行中"})
}

func (h *ACFPHandler) CleanNotifyEvents(c *gin.Context) {
	count := h.svc.CleanupOldNotifyEvents()
	response.Success(c, gin.H{"cleaned": count})
}

// ─── 通用模块配置 ───

// GetModuleConfig 获取模块配置
func (h *ACFPHandler) GetModuleConfig(c *gin.Context) {
	key := c.Param("key")
	// 从数据库读取模块配置
	var config struct {
		ID        uint   `json:"id"`
		Key       string `json:"key"`
		Value     string `json:"value"`
		UpdatedAt string `json:"updated_at"`
	}
	if err := h.svc.GetDB().Table("acfp_module_configs").Where("module_key = ?", key).First(&config).Error; err != nil {
		response.NotFound(c, "模块配置不存在")
		return
	}
	response.Success(c, config)
}

// ToggleModule 切换模块启用/禁用
func (h *ACFPHandler) ToggleModule(c *gin.Context) {
	key := c.Param("key")
	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	enabled := 0
	if req.Enabled {
		enabled = 1
	}
	if err := h.svc.GetDB().Exec("UPDATE acfp_module_configs SET enabled = ?, updated_at = NOW() WHERE module_key = ?", enabled, key).Error; err != nil {
		response.ServerError(c, "更新失败")
		return
	}
	response.SuccessMsg(c, "模块状态已更新")
}

// ─── 失败通知配置 ───

// GetFailNotifyConfig 获取失败通知配置
func (h *ACFPHandler) GetFailNotifyConfig(c *gin.Context) {
	var config struct {
		Enabled    bool   `json:"enabled"`
		Email      string `json:"email"`
		Webhook    string `json:"webhook"`
		RetryCount int    `json:"retry_count"`
	}
	// 从配置表读取
	h.svc.GetDB().Table("acfp_configs").Where("config_key = ?", "fail_notify").Scan(&config)
	response.Success(c, config)
}

// SaveFailNotifyConfig 保存失败通知配置
func (h *ACFPHandler) SaveFailNotifyConfig(c *gin.Context) {
	var req struct {
		Enabled    bool   `json:"enabled"`
		Email      string `json:"email"`
		Webhook    string `json:"webhook"`
		RetryCount int    `json:"retry_count"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	// 保存到配置表
	data, _ := json.Marshal(req)
	h.svc.GetDB().Exec("INSERT INTO acfp_configs (config_key, config_value, updated_at) VALUES (?, ?, NOW()) ON DUPLICATE KEY UPDATE config_value = ?, updated_at = NOW()", "fail_notify", string(data), string(data))
	response.SuccessMsg(c, "配置已保存")
}

// ─── 状态对账配置 ───

// GetStatusSyncConfig 获取状态对账配置
func (h *ACFPHandler) GetStatusSyncConfig(c *gin.Context) {
	var config struct {
		Enabled    bool   `json:"enabled"`
		Interval   int    `json:"interval"`
		NotifyFail bool   `json:"notify_fail"`
	}
	h.svc.GetDB().Table("acfp_configs").Where("config_key = ?", "status_sync").Scan(&config)
	response.Success(c, config)
}

// SaveStatusSyncConfig 保存状态对账配置
func (h *ACFPHandler) SaveStatusSyncConfig(c *gin.Context) {
	var req struct {
		Enabled    bool `json:"enabled"`
		Interval   int  `json:"interval"`
		NotifyFail bool `json:"notify_fail"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	data, _ := json.Marshal(req)
	h.svc.GetDB().Exec("INSERT INTO acfp_configs (config_key, config_value, updated_at) VALUES (?, ?, NOW()) ON DUPLICATE KEY UPDATE config_value = ?, updated_at = NOW()", "status_sync", string(data), string(data))
	response.SuccessMsg(c, "配置已保存")
}

// GetUpstreamCache 获取上游缓存
func (h *ACFPHandler) GetUpstreamCache(c *gin.Context) {
	hostID, _ := strconv.ParseUint(c.Param("host_id"), 10, 64)
	data, err := h.svc.GetUpstreamCache(uint(hostID))
	if err != nil {
		response.NotFound(c, "缓存不存在")
		return
	}
	response.Success(c, gin.H{"host_id": hostID, "data": data})
}

// ─── 限量发售管理 ───

// UpdateLimitedSale 更新限量发售
func (h *ACFPHandler) UpdateLimitedSale(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var item model.ACFPLimitedSale
	if err := c.ShouldBindJSON(&item); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	item.ID = uint(id)
	if err := h.svc.SetLimitedSale(&item); err != nil {
		response.ServerError(c, "更新失败")
		return
	}
	response.SuccessMsg(c, "更新成功")
}

// DeleteLimitedSale 删除限量发售
func (h *ACFPHandler) DeleteLimitedSale(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.GetDB().Delete(&model.ACFPLimitedSale{}, id).Error; err != nil {
		response.ServerError(c, "删除失败")
		return
	}
	response.SuccessMsg(c, "删除成功")
}

// ResetLimitedSaleQuota 重置限量发售配额
func (h *ACFPHandler) ResetLimitedSaleQuota(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.GetDB().Model(&model.ACFPLimitedSale{}).Where("id = ?", id).Update("sold_qty", 0).Error; err != nil {
		response.ServerError(c, "重置失败")
		return
	}
	response.SuccessMsg(c, "配额已重置")
}

// ─── 批量商品修改 ───

// BatchUpdateProducts 批量更新商品
func (h *ACFPHandler) BatchUpdateProducts(c *gin.Context) {
	var req struct {
		ProductIDs []uint `json:"product_ids"`
		Action     string `json:"action"`
		Value      string `json:"value"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	// 创建批量任务
	task := &model.ACFPBatchTask{
		Status: "pending",
		Total:  len(req.ProductIDs),
	}
	data, _ := json.Marshal(req)
	task.Filter = string(data)
	task.Changes = string(data)
	if err := h.svc.CreateBatchTask(task); err != nil {
		response.ServerError(c, "创建任务失败")
		return
	}
	// 执行任务
	h.svc.ExecuteBatchTask(task.ID)
	response.Success(c, task)
}

// GetCacheWarmStatus 获取缓存预热状态
func (h *ACFPHandler) GetCacheWarmStatus(c *gin.Context) {
	// 返回缓存预热状态
	response.Success(c, gin.H{
		"status":     "ready",
		"last_warm":  nil,
		"total_keys": 0,
	})
}

// ─── 业务列表 Pro ───

// GetBusinessList 获取业务列表
func (h *ACFPHandler) GetBusinessList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	type BusinessRow struct {
		ID          uint   `json:"id"`
		UserID      uint   `json:"user_id"`
		Username    string `json:"username"`
		ProductID   uint   `json:"product_id"`
		ProductName string `json:"product_name"`
		Domain      string `json:"domain"`
		Status      string `json:"status"`
		IP          string `json:"ip"`
		CreatedAt   string `json:"created_at"`
	}

	var items []BusinessRow
	var total int64

	h.svc.GetDB().Raw(`
		SELECT h.id, h.userid as user_id, u.username, h.productid as product_id, 
		       p.name as product_name, h.domain, h.domainstatus as status, 
		       h.dedicatedip as ip, h.created_at
		FROM host h
		LEFT JOIN users u ON h.userid = u.id
		LEFT JOIN products p ON h.productid = p.id
	`).Count(&total).Offset((page - 1) * pageSize).Limit(pageSize).Scan(&items)

	response.Success(c, gin.H{"items": items, "total": total, "page": page})
}

// GetBusinessRow 获取单条业务详情
func (h *ACFPHandler) GetBusinessRow(c *gin.Context) {
	hostID, _ := strconv.ParseUint(c.Param("host_id"), 10, 64)
	type BusinessDetail struct {
		ID          uint   `json:"id"`
		UserID      uint   `json:"user_id"`
		Username    string `json:"username"`
		ProductID   uint   `json:"product_id"`
		ProductName string `json:"product_name"`
		Domain      string `json:"domain"`
		Status      string `json:"status"`
		IP          string `json:"ip"`
		BillingCycle string `json:"billing_cycle"`
		Amount      float64 `json:"amount"`
		CreatedAt   string `json:"created_at"`
	}
	var detail BusinessDetail
	if err := h.svc.GetDB().Raw(`
		SELECT h.id, h.userid as user_id, u.username, h.productid as product_id,
		       p.name as product_name, h.domain, h.domainstatus as status,
		       h.dedicatedip as ip, h.billingcycle as billing_cycle,
		       h.amount, h.created_at
		FROM host h
		LEFT JOIN users u ON h.userid = u.id
		LEFT JOIN products p ON h.productid = p.id
		WHERE h.id = ?
	`, hostID).Scan(&detail).Error; err != nil {
		response.NotFound(c, "业务不存在")
		return
	}
	response.Success(c, detail)
}

// GetBusinessSnapshot 获取业务快照
func (h *ACFPHandler) GetBusinessSnapshot(c *gin.Context) {
	hostID, _ := strconv.ParseUint(c.Param("host_id"), 10, 64)
	data, err := h.svc.GetUpstreamCache(uint(hostID))
	if err != nil {
		response.NotFound(c, "快照不存在")
		return
	}
	response.Success(c, gin.H{"host_id": hostID, "snapshot": data})
}

// SyncOneBusiness 同步单条业务
func (h *ACFPHandler) SyncOneBusiness(c *gin.Context) {
	hostID, _ := strconv.ParseUint(c.Param("host_id"), 10, 64)
	// 获取业务信息
	type HostInfo struct {
		Status string `gorm:"column:domainstatus"`
		DcimID int    `gorm:"column:dcimid"`
	}
	var host HostInfo
	if err := h.svc.GetDB().Raw("SELECT domainstatus, IFNULL(dcimid,0) as dcimid FROM host WHERE id = ?", hostID).Scan(&host).Error; err != nil {
		response.NotFound(c, "业务不存在")
		return
	}
	// 更新缓存
	cacheData := map[string]interface{}{
		"host_id": hostID,
		"status":  host.Status,
		"sync_at": time.Now().Format(time.RFC3339),
	}
	h.svc.SaveUpstreamCache(uint(hostID), cacheData)
	response.SuccessMsg(c, "同步成功")
}

// SuspendOneBusiness 暂停单条业务
func (h *ACFPHandler) SuspendOneBusiness(c *gin.Context) {
	hostID, _ := strconv.ParseUint(c.Param("host_id"), 10, 64)
	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)
	if err := h.svc.GetDB().Exec("UPDATE host SET domainstatus = 'Suspended' WHERE id = ?", hostID).Error; err != nil {
		response.ServerError(c, "暂停失败")
		return
	}
	response.SuccessMsg(c, "业务已暂停")
}

// UnsuspendOneBusiness 恢复单条业务
func (h *ACFPHandler) UnsuspendOneBusiness(c *gin.Context) {
	hostID, _ := strconv.ParseUint(c.Param("host_id"), 10, 64)
	if err := h.svc.GetDB().Exec("UPDATE host SET domainstatus = 'Active' WHERE id = ?", hostID).Error; err != nil {
		response.ServerError(c, "恢复失败")
		return
	}
	response.SuccessMsg(c, "业务已恢复")
}

// DeleteOneBusiness 删除单条业务
func (h *ACFPHandler) DeleteOneBusiness(c *gin.Context) {
	hostID, _ := strconv.ParseUint(c.Param("host_id"), 10, 64)
	if err := h.svc.GetDB().Exec("UPDATE host SET domainstatus = 'Terminated' WHERE id = ?", hostID).Error; err != nil {
		response.ServerError(c, "删除失败")
		return
	}
	response.SuccessMsg(c, "业务已删除")
}

// ProvisionOneBusiness 开通单条业务
func (h *ACFPHandler) ProvisionOneBusiness(c *gin.Context) {
	hostID, _ := strconv.ParseUint(c.Param("host_id"), 10, 64)
	if err := h.svc.GetDB().Exec("UPDATE host SET domainstatus = 'Active' WHERE id = ?", hostID).Error; err != nil {
		response.ServerError(c, "开通失败")
		return
	}
	response.SuccessMsg(c, "业务已开通")
}

// GetBusinessStats 获取业务统计
func (h *ACFPHandler) GetBusinessStats(c *gin.Context) {
	type Stats struct {
		Total     int64 `json:"total"`
		Active    int64 `json:"active"`
		Suspended int64 `json:"suspended"`
		Pending   int64 `json:"pending"`
		Terminated int64 `json:"terminated"`
	}
	var stats Stats
	h.svc.GetDB().Raw("SELECT COUNT(*) as total FROM host").Scan(&stats.Total)
	h.svc.GetDB().Raw("SELECT COUNT(*) as total FROM host WHERE domainstatus = 'Active'").Scan(&stats.Active)
	h.svc.GetDB().Raw("SELECT COUNT(*) as total FROM host WHERE domainstatus = 'Suspended'").Scan(&stats.Suspended)
	h.svc.GetDB().Raw("SELECT COUNT(*) as total FROM host WHERE domainstatus = 'Pending'").Scan(&stats.Pending)
	h.svc.GetDB().Raw("SELECT COUNT(*) as total FROM host WHERE domainstatus = 'Terminated'").Scan(&stats.Terminated)
	response.Success(c, stats)
}
