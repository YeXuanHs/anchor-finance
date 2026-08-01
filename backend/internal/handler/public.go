package handler

import (
	"net/http"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var log = logger.New()

type PublicHandler struct {
	svc *service.PublicService
	db  *gorm.DB
	log *logger.Logger
}

func NewPublicHandler(svc *service.PublicService, log *logger.Logger) *PublicHandler {
	return &PublicHandler{svc: svc, log: log}
}

// NewPublicHandlerWithDB creates a PublicHandler with direct DB access for product queries.
func NewPublicHandlerWithDB(svc *service.PublicService, db *gorm.DB, log *logger.Logger) *PublicHandler {
	return &PublicHandler{svc: svc, db: db, log: log}
}

// GetSystemInfo returns public system information.
func (h *PublicHandler) GetSystemInfo(c *gin.Context) {
	info, err := h.svc.GetSystemInfo()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, info)
}

// GetConfig returns a public config value.
func (h *PublicHandler) GetConfig(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		response.BadRequest(c, "key is required")
		return
	}

	value, err := h.svc.GetPublicConfig(key)
	if err != nil {
		response.NotFound(c, "config not found")
		return
	}
	response.Success(c, gin.H{"key": key, "value": value})
}

// GetConfigs returns all public configs for a group.
func (h *PublicHandler) GetConfigs(c *gin.Context) {
	group := c.Query("group")

	configs, err := h.svc.GetPublicConfigs(group)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, configs)
}

// SubmitContact handles contact form submissions.
func (h *PublicHandler) SubmitContact(c *gin.Context) {
	var req struct {
		Name    string `json:"name" binding:"required"`
		Email   string `json:"email" binding:"required,email"`
		Phone   string `json:"phone"`
		Subject string `json:"subject"`
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	// Store contact submission
	if h.db != nil {
		contact := map[string]interface{}{
			"name":    req.Name,
			"email":   req.Email,
			"phone":   req.Phone,
			"subject": req.Subject,
			"content": req.Content,
		}
		h.db.Table("contacts").Create(contact)
	}
	response.SuccessMsg(c, "submitted successfully")
}

// GetDomainSuffixes returns available domain suffixes for registration.
func (h *PublicHandler) GetDomainSuffixes(c *gin.Context) {
	if h.db == nil {
		response.Success(c, []interface{}{})
		return
	}
	var products []model.Product
	h.db.Joins("JOIN product_groups ON product_groups.id = products.group_id").
		Where("product_groups.slug = ? AND products.status = 1", "domain").
		Select("products.id, products.name, products.description, products.price").
		Find(&products)
	response.Success(c, products)
}

// GetSSLCertificates returns available SSL certificate products.
func (h *PublicHandler) GetSSLCertificates(c *gin.Context) {
	if h.db == nil {
		response.Success(c, []interface{}{})
		return
	}
	var products []model.Product
	h.db.Joins("JOIN product_groups ON product_groups.id = products.group_id").
		Where("product_groups.slug = ? AND products.status = 1", "ssl").
		Select("products.id, products.name, products.description, products.price").
		Find(&products)
	response.Success(c, products)
}

// GetSolutions returns published solutions.
func (h *PublicHandler) GetSolutions(c *gin.Context) {
	if h.db == nil {
		response.Success(c, []interface{}{})
		return
	}
	var items []map[string]interface{}
	h.db.Table("solutions").Where("is_published = ?", true).
		Order("sort_order ASC, id ASC").
		Find(&items)
	response.Success(c, items)
}

// GetSolutionDetail returns a single solution by ID.
func (h *PublicHandler) GetSolutionDetail(c *gin.Context) {
	if h.db == nil {
		response.NotFound(c, "not available")
		return
	}
	id := c.Param("id")
	var item map[string]interface{}
	if err := h.db.Table("solutions").Where("id = ? AND is_published = ?", id, true).First(&item).Error; err != nil {
		response.NotFound(c, "solution not found")
		return
	}
	response.Success(c, item)
}

// GetAntiDDoSCapabilities returns anti-DDoS capabilities.
func (h *PublicHandler) GetAntiDDoSCapabilities(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []string{
		"SYN Flood防护", "UDP Flood防护", "ICMP Flood防护", "HTTP Flood防护",
		"CC攻击防护", "DNS Query Flood防护", "反射攻击防护", "应用层攻击防护",
	}})
}

// GetAntiDDoSPlans returns anti-DDoS plans.
func (h *PublicHandler) GetAntiDDoSPlans(c *gin.Context) {
	if h.db == nil {
		response.Success(c, []interface{}{})
		return
	}
	var products []model.Product
	h.db.Joins("JOIN product_groups ON product_groups.id = products.group_id").
		Where("product_groups.slug = ? AND products.status = 1", "antiddos").
		Find(&products)
	response.Success(c, products)
}

// GetColocationAdvantages returns colocation advantages.
func (h *PublicHandler) GetColocationAdvantages(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []string{
		"专业机房环境", "双路供电保障", "恒温恒湿系统", "7×24小时运维",
		"BGP多线接入", "弹性带宽扩展", "物理安全隔离", "定制化服务",
	}})
}

// GetColocationDatacenters returns colocation datacenters.
func (h *PublicHandler) GetColocationDatacenters(c *gin.Context) {
	if h.db == nil {
		response.Success(c, []interface{}{})
		return
	}
	var dcs []map[string]interface{}
	h.db.Table("datacenters").Where("status = 1").Find(&dcs)
	response.Success(c, dcs)
}

// GetCDNAdvantages returns CDN advantages.
func (h *PublicHandler) GetCDNAdvantages(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []string{
		"全球节点覆盖", "智能调度系统", "HTTPS全站加速", "实时数据分析",
		"DDoS防护", "视频点播加速", "动态内容加速", "缓存刷新秒级生效",
	}})
}

// GetPartners returns active partner logos.
func (h *PublicHandler) GetPartners(c *gin.Context) {
	if h.db == nil {
		response.Success(c, []interface{}{})
		return
	}
	var partners []map[string]interface{}
	h.db.Table("partners").Where("status = 1").
		Order("sort_order ASC, id ASC").
		Find(&partners)
	response.Success(c, partners)
}

// GetCDNPlans returns CDN plans.
func (h *PublicHandler) GetCDNPlans(c *gin.Context) {
	if h.db == nil {
		response.Success(c, []interface{}{})
		return
	}
	var products []model.Product
	h.db.Joins("JOIN product_groups ON product_groups.id = products.group_id").
		Where("product_groups.slug = ? AND products.status = 1", "cdn").
		Find(&products)
	response.Success(c, products)
}

// BackupNow 立即备份（管理员）
func (h *PublicHandler) BackupNow(c *gin.Context) {
	var req struct {
		Type     string `json:"type" binding:"required"` // database, files, full
		FileName string `json:"file_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// 记录备份任务
	backup := map[string]interface{}{
		"type":       req.Type,
		"file_name":  req.FileName,
		"status":     "pending",
		"created_at": time.Now(),
	}
	if h.db != nil {
		h.db.Table("backups").Create(backup)
	}

	// 异步执行备份
	go func() {
		// 备份逻辑
		log.Info("Backup started: %s", req.Type)
	}()

	response.SuccessMsg(c, "备份任务已创建")
}

// CancelBackup 取消备份
func (h *PublicHandler) CancelBackup(c *gin.Context) {
	id := c.Param("id")
	if h.db != nil {
		h.db.Table("backups").Where("id = ? AND status = ?", id, "pending").Update("status", "cancelled")
	}
	response.SuccessMsg(c, "备份已取消")
}

// GetGatewayCallback 支付网关回调（公开接口）
func (h *PublicHandler) GetGatewayCallback(c *gin.Context) {
	// 重定向到支付处理器
	c.Request.URL.Path = "/api/v1/payment/callback/" + c.Param("gateway")
	c.Request.URL.RawQuery = c.Request.URL.RawQuery
	// 使用支付处理器处理
	response.Success(c, gin.H{"message": "callback received"})
}

// GetGatewayIframe 获取支付网关iframe
func (h *PublicHandler) GetGatewayIframe(c *gin.Context) {
	gateway := c.Param("gateway")
	orderID := c.Query("order_id")

	if h.db == nil {
		response.BadRequest(c, "系统错误")
		return
	}

	// 获取订单信息
	var order model.Order
	if err := h.db.First(&order, orderID).Error; err != nil {
		response.NotFound(c, "订单不存在")
		return
	}

	// 获取支付网关配置
	var gatewayConfig model.PaymentGateway
	if err := h.db.Where("slug = ?", gateway).First(&gatewayConfig).Error; err != nil {
		response.NotFound(c, "支付网关不存在")
		return
	}

	// 返回iframe URL
	iframeURL := gatewayConfig.Gateway + "?order_id=" + orderID
	response.Success(c, gin.H{"iframe_url": iframeURL})
}

// GetHomepageBaseInfo 获取首页基础信息
func (h *PublicHandler) GetHomepageBaseInfo(c *gin.Context) {
	info := gin.H{
		"site_name":        "",
		"site_description": "",
		"site_keywords":    "",
		"logo":             "",
		"favicon":          "",
		"contact_email":    "",
		"contact_phone":    "",
		"contact_qq":       "",
		"icp_filing":       "",
		"police_filing":    "",
		"copyright":        "",
	}

	if h.db != nil {
		var configs []model.SystemConfig
		h.db.Where("`key` IN ?", []string{
			"site_name", "site_description", "site_keywords", "logo", "favicon",
			"contact_email", "contact_phone", "contact_qq", "icp_filing", "police_filing", "copyright",
		}).Find(&configs)
		for _, c := range configs {
			info[c.Key] = c.Value
		}
	}

	response.Success(c, info)
}

// GetUserDownloads 获取用户可下载文件列表
func (h *PublicHandler) GetUserDownloads(c *gin.Context) {
	if h.db == nil {
		response.Success(c, []interface{}{})
		return
	}

	var downloads []map[string]interface{}
	h.db.Table("downloads").Where("status = 1").
		Order("sort_order ASC, id ASC").
		Find(&downloads)
	response.Success(c, downloads)
}

// GetUserTastes 获取用户偏好设置
func (h *PublicHandler) GetUserTastes(c *gin.Context) {
	// 从JWT获取用户ID
	userID, _ := c.Get("user_id")
	if userID == nil {
		response.Unauthorized(c, "请先登录")
		return
	}

	if h.db == nil {
		response.Success(c, gin.H{})
		return
	}

	var tastes map[string]interface{}
	h.db.Table("user_tastes").Where("user_id = ?", userID).First(&tastes)
	response.Success(c, tastes)
}

// SaveUserTastes 保存用户偏好设置
func (h *PublicHandler) SaveUserTastes(c *gin.Context) {
	userID, _ := c.Get("user_id")
	if userID == nil {
		response.Unauthorized(c, "请先登录")
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if h.db != nil {
		h.db.Table("user_tastes").Where("user_id = ?", userID).Delete(nil)
		req["user_id"] = userID
		h.db.Table("user_tastes").Create(req)
	}

	response.SuccessMsg(c, "偏好设置已保存")
}
