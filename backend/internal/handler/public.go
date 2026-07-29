package handler

import (
	"net/http"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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
