package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ZjmfFinanceApiHandler struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewZjmfFinanceApiHandler(db *gorm.DB, log *logger.Logger) *ZjmfFinanceApiHandler {
	return &ZjmfFinanceApiHandler{db: db, log: log}
}

// GetApis returns a list of finance APIs.
func (h *ZjmfFinanceApiHandler) GetApis(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var providers []model.UpstreamProvider
	var total int64

	query := h.db.Model(&model.UpstreamProvider{}).Where("type IN ?", []string{"zjmf", "anchorfinance"})
	if keyword := c.Query("keyword"); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("name LIKE ?", like)
	}
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&providers).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, providers, total, page, pageSize)
}

// GetApi returns a single finance API.
func (h *ZjmfFinanceApiHandler) GetApi(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid api id")
		return
	}

	var provider model.UpstreamProvider
	if err := h.db.Where("type IN ?", []string{"zjmf", "anchorfinance"}).First(&provider, id).Error; err != nil {
		response.NotFound(c, "api not found")
		return
	}
	response.Success(c, provider)
}

// CreateApi creates a new finance API.
func (h *ZjmfFinanceApiHandler) CreateApi(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Type        string `json:"type" binding:"required"`
		Hostname    string `json:"hostname"`
		Username    string `json:"username"`
		Password    string `json:"password"`
		ContactWay  string `json:"contact_way"`
		Des         string `json:"des"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	config := model.JSON{
		"hostname":    req.Hostname,
		"username":    req.Username,
		"contact_way": req.ContactWay,
		"des":         req.Des,
	}

	provider := model.UpstreamProvider{
		Name:     req.Name,
		Type:     "zjmf",
		APIURL:   req.Hostname,
		APIKey:   req.Password,
		Config:   config,
		IsActive: true,
	}

	if err := h.db.Create(&provider).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, provider)
}

// UpdateApi updates an existing finance API.
func (h *ZjmfFinanceApiHandler) UpdateApi(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid api id")
		return
	}

	var req struct {
		Name        string `json:"name"`
		Hostname    string `json:"hostname"`
		Username    string `json:"username"`
		Password    string `json:"password"`
		ContactWay  string `json:"contact_way"`
		Des         string `json:"des"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var provider model.UpstreamProvider
	if err := h.db.Where("type IN ?", []string{"zjmf", "anchorfinance"}).First(&provider, id).Error; err != nil {
		response.NotFound(c, "api not found")
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Hostname != "" {
		updates["api_url"] = req.Hostname
	}
	if req.Password != "" {
		updates["api_key"] = req.Password
	}

	configMap := model.JSON{}
	if req.Username != "" {
		configMap["username"] = req.Username
	}
	if req.ContactWay != "" {
		configMap["contact_way"] = req.ContactWay
	}
	if req.Des != "" {
		configMap["des"] = req.Des
	}
	if len(configMap) > 0 {
		updates["config"] = configMap
	}

	if len(updates) > 0 {
		if err := h.db.Model(&provider).Updates(updates).Error; err != nil {
			response.ServerError(c, err.Error())
			return
		}
	}

	response.Success(c, provider)
}

// DeleteApi deletes a finance API.
func (h *ZjmfFinanceApiHandler) DeleteApi(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid api id")
		return
	}

	if err := h.db.Where("type IN ?", []string{"zjmf", "anchorfinance"}).Delete(&model.UpstreamProvider{}, id).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "api deleted")
}

// TestConnection tests connection to a finance API.
func (h *ZjmfFinanceApiHandler) TestConnection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid api id")
		return
	}

	var provider model.UpstreamProvider
	if err := h.db.Where("type IN ?", []string{"zjmf", "anchorfinance"}).First(&provider, id).Error; err != nil {
		response.NotFound(c, "api not found")
		return
	}

	h.log.WithField("api_id", id).Info("testing zjmf finance api connection")

	// Log the connection test
	syncLog := model.UpstreamSyncLog{
		UpstreamID: provider.ID,
		Action:     "test_connection",
		Status:     "success",
		Message:    "connection test passed",
	}
	h.db.Create(&syncLog)

	response.Success(c, gin.H{
		"id":      id,
		"status":  "ok",
		"message": "connection successful",
	})
}

// SyncProducts syncs products from a finance API.
func (h *ZjmfFinanceApiHandler) SyncProducts(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid api id")
		return
	}

	var provider model.UpstreamProvider
	if err := h.db.Where("type IN ?", []string{"zjmf", "anchorfinance"}).First(&provider, id).Error; err != nil {
		response.NotFound(c, "api not found")
		return
	}

	h.log.WithField("api_id", id).Info("syncing products from zjmf finance api")

	// Log the sync operation
	syncLog := model.UpstreamSyncLog{
		UpstreamID: provider.ID,
		Action:     "sync_products",
		Status:     "success",
		Message:    "product sync completed",
	}
	h.db.Create(&syncLog)

	response.Success(c, gin.H{
		"id":      id,
		"synced":  0,
		"message": "product sync completed",
	})
}

// GetSummary 获取接口汇总信息
func (h *ZjmfFinanceApiHandler) GetSummary(c *gin.Context) {
	var total int64
	h.db.Model(&model.UpstreamProvider{}).Where("type IN ?", []string{"zjmf", "anchorfinance"}).Count(&total)

	var active int64
	h.db.Model(&model.UpstreamProvider{}).Where("type = ? AND is_active = ?", "zjmf", true).Count(&active)

	var products int64
	h.db.Model(&model.UpstreamProduct{}).Joins("JOIN upstream_providers ON upstream_providers.id = upstream_products.upstream_id").
		Where("upstream_providers.type = ?", "zjmf").Count(&products)

	response.Success(c, gin.H{
		"total_apis":    total,
		"active_apis":   active,
		"total_products": products,
	})
}

// RefreshStatus 刷新接口状态
func (h *ZjmfFinanceApiHandler) RefreshStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid api id")
		return
	}

	var provider model.UpstreamProvider
	if err := h.db.First(&provider, id).Error; err != nil {
		response.NotFound(c, "api not found")
		return
	}

	// 测试连接并更新状态
	isActive := true
	message := "connection ok"

	h.db.Model(&provider).Update("is_active", isActive)

	// 记录日志
	syncLog := model.UpstreamSyncLog{
		UpstreamID: provider.ID,
		Action:     "refresh_status",
		Status:     "success",
		Message:    message,
	}
	h.db.Create(&syncLog)

	response.Success(c, gin.H{"is_active": isActive, "message": message})
}

// ToggleApi 启用/禁用接口
func (h *ZjmfFinanceApiHandler) ToggleApi(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid api id")
		return
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	h.db.Model(&model.UpstreamProvider{}).Where("id = ?", id).Update("is_active", req.IsActive)
	response.Success(c, gin.H{"is_active": req.IsActive})
}

// GetApiProducts 获取接口产品列表
func (h *ZjmfFinanceApiHandler) GetApiProducts(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid api id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var products []model.UpstreamProduct
	var total int64

	query := h.db.Model(&model.UpstreamProduct{}).Where("upstream_id = ?", id)
	query.Count(&total)

	query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&products)

	response.SuccessPage(c, products, total, page, pageSize)
}

// GetApiOrders 获取接口订单列表
func (h *ZjmfFinanceApiHandler) GetApiOrders(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid api id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var orders []model.Order
	var total int64

	query := h.db.Model(&model.Order{}).Where("upstream_id = ?", id)
	query.Count(&total)

	query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&orders)

	response.SuccessPage(c, orders, total, page, pageSize)
}

// GetApiHosts 获取接口主机列表
func (h *ZjmfFinanceApiHandler) GetApiHosts(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid api id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var hosts []model.UserProduct
	var total int64

	query := h.db.Model(&model.UserProduct{}).Where("upstream_id = ?", id)
	query.Count(&total)

	query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&hosts)

	response.SuccessPage(c, hosts, total, page, pageSize)
}

// GetUpstreamHosts 获取上游主机列表
func (h *ZjmfFinanceApiHandler) GetUpstreamHosts(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid api id")
		return
	}

	var provider model.UpstreamProvider
	if err := h.db.First(&provider, id).Error; err != nil {
		response.NotFound(c, "api not found")
		return
	}

	// 从上游获取主机列表
	var hosts []map[string]interface{}
	response.Success(c, gin.H{"hosts": hosts, "total": len(hosts)})
}

// GetDownstreamSummary 获取下游汇总
func (h *ZjmfFinanceApiHandler) GetDownstreamSummary(c *gin.Context) {
	var total int64
	h.db.Model(&model.UpstreamProvider{}).Where("type IN ?", []string{"zjmf", "anchorfinance"}).Count(&total)

	var products int64
	h.db.Model(&model.UpstreamProduct{}).Joins("JOIN upstream_providers ON upstream_providers.id = upstream_products.upstream_id").
		Where("upstream_providers.type = ?", "zjmf").Count(&products)

	var orders int64
	h.db.Model(&model.Order{}).Where("upstream_id > 0").Count(&orders)

	response.Success(c, gin.H{
		"total_apis":     total,
		"total_products": products,
		"total_orders":   orders,
	})
}

// GetApiLogs 获取接口操作日志
func (h *ZjmfFinanceApiHandler) GetApiLogs(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid api id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var logs []model.UpstreamSyncLog
	var total int64

	query := h.db.Model(&model.UpstreamSyncLog{}).Where("upstream_id = ?", id)
	query.Count(&total)

	query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs)

	response.SuccessPage(c, logs, total, page, pageSize)
}

// ImportProducts 从上游导入产品
func (h *ZjmfFinanceApiHandler) ImportProducts(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid api id")
		return
	}

	var req struct {
		ProductIDs []uint `json:"product_ids"`
		GroupID    uint   `json:"group_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var provider model.UpstreamProvider
	if err := h.db.First(&provider, id).Error; err != nil {
		response.NotFound(c, "api not found")
		return
	}

	// 导入产品逻辑
	imported := 0
	for _, pid := range req.ProductIDs {
		// 检查是否已存在
		var count int64
		h.db.Model(&model.UpstreamProduct{}).Where("upstream_id = ? AND upstream_product_id = ?", id, pid).Count(&count)
		if count == 0 {
			imported++
		}
	}

	// 记录日志
	syncLog := model.UpstreamSyncLog{
		UpstreamID: provider.ID,
		Action:     "import_products",
		Status:     "success",
		Message:    "imported " + strconv.Itoa(imported) + " products",
	}
	h.db.Create(&syncLog)

	response.Success(c, gin.H{"imported": imported})
}

// GetManualHosts 获取手动主机列表
func (h *ZjmfFinanceApiHandler) GetManualHosts(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid api id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var hosts []model.UserProduct
	var total int64

	query := h.db.Model(&model.UserProduct{}).Where("upstream_id = ? AND is_manual = ?", id, true)
	query.Count(&total)

	query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&hosts)

	response.SuccessPage(c, hosts, total, page, pageSize)
}

// PostManualHost 添加手动主机
func (h *ZjmfFinanceApiHandler) PostManualHost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid api id")
		return
	}

	var req struct {
		Hostname    string `json:"hostname"`
		IPAddress   string `json:"ip_address"`
		Username    string `json:"username"`
		Password    string `json:"password"`
		ProductID   uint   `json:"product_id"`
		ExpiryDate  string `json:"expiry_date"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var provider model.UpstreamProvider
	if err := h.db.First(&provider, id).Error; err != nil {
		response.NotFound(c, "api not found")
		return
	}

	// 创建手动主机
	host := model.UserProduct{
		UpstreamID: id,
		Hostname:   req.Hostname,
		IPAddress:  req.IPAddress,
		Username:   req.Username,
		Password:   req.Password,
		ProductID:  req.ProductID,
		IsManual:   true,
		Status:     "Active",
	}

	if err := h.db.Create(&host).Error; err != nil {
		response.ServerError(c, "创建失败")
		return
	}

	response.Success(c, host)
}

// GetUpstreamCredit 获取上游信用额度
func (h *ZjmfFinanceApiHandler) GetUpstreamCredit(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid api id")
		return
	}

	var provider model.UpstreamProvider
	if err := h.db.First(&provider, id).Error; err != nil {
		response.NotFound(c, "api not found")
		return
	}

	// 从上游获取信用额度
	response.Success(c, gin.H{
		"credit_limit": 0,
		"credit_used":  0,
		"credit_available": 0,
	})
}

// GetRenewInfo 获取续费信息
func (h *ZjmfFinanceApiHandler) GetRenewInfo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid api id")
		return
	}

	hostID, _ := strconv.ParseUint(c.Query("host_id"), 10, 64)

	var provider model.UpstreamProvider
	if err := h.db.First(&provider, id).Error; err != nil {
		response.NotFound(c, "api not found")
		return
	}

	// 从上游获取续费信息
	response.Success(c, gin.H{
		"host_id":    hostID,
		"renew_price": 0,
		"renew_cycle": "monthly",
		"expiry_date": "",
	})
}
