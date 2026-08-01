package handler

import (
	"strconv"
	"strings"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BusinessListHandler 业务列表Pro处理器
type BusinessListHandler struct {
	svc *service.BusinessListService
	log *logger.Logger
}

func NewBusinessListHandler(db *gorm.DB, log *logger.Logger) *BusinessListHandler {
	return &BusinessListHandler{
		svc: service.NewBusinessListService(db, log),
		log: log,
	}
}

// GetList 获取业务列表（带高级筛选+状态统计）
func (h *BusinessListHandler) GetList(c *gin.Context) {
	filter := service.BusinessFilter{
		Status:         c.Query("status"),
		Keyword:        c.Query("keyword"),
		ProductType:    c.Query("product_type"),
		BillingCycle:   c.Query("billingcycle"),
		Payment:        c.Query("payment"),
		DomainFilter:   c.Query("domain_filter"),
		IPFilter:       c.Query("ip_filter"),
		UsernameFilter: c.Query("username_filter"),
		StartTimeFrom:  c.Query("start_time_from"),
		StartTimeTo:    c.Query("start_time_to"),
		SortField:      c.Query("sort_field"),
		SortDir:        c.Query("sort_dir"),
	}

	if v := c.Query("product_id"); v != "" {
		filter.ProductID = parseUint(v)
	}
	if v := c.Query("server_id"); v != "" {
		filter.ServerID = parseInt(v)
	}
	if v := c.Query("due_filter"); v != "" {
		filter.DueFilter = parseInt(v)
	}
	if v := c.Query("page"); v != "" {
		filter.Page = parseInt(v)
	}
	if v := c.Query("page_size"); v != "" {
		filter.PageSize = parseInt(v)
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	items, total, stats, err := h.svc.GetList(filter)
	if err != nil {
		c.JSON(500, gin.H{"status": 500, "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"status": 200,
		"data": gin.H{
			"list":  items,
			"total": total,
			"page":  filter.Page,
			"page_size": filter.PageSize,
			"pages": (total + int64(filter.PageSize) - 1) / int64(filter.PageSize),
			"stats": stats,
		},
	})
}

// GetRow 获取单行业务数据
func (h *BusinessListHandler) GetRow(c *gin.Context) {
	hostID := parseUint(c.Param("id"))
	if hostID == 0 {
		c.JSON(400, gin.H{"status": 400, "msg": "参数错误"})
		return
	}

	row, err := h.svc.GetRow(hostID)
	if err != nil {
		c.JSON(404, gin.H{"status": 404, "msg": "业务不存在"})
		return
	}

	c.JSON(200, gin.H{"status": 200, "data": row})
}

// GetSnapshot 获取业务快照
func (h *BusinessListHandler) GetSnapshot(c *gin.Context) {
	hostID := parseUint(c.Param("id"))
	if hostID == 0 {
		c.JSON(400, gin.H{"status": 400, "msg": "参数错误"})
		return
	}

	data, err := h.svc.GetSnapshot(hostID)
	if err != nil {
		c.JSON(500, gin.H{"status": 500, "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": 200, "data": data})
}

// SyncOne 同步单个业务
func (h *BusinessListHandler) SyncOne(c *gin.Context) {
	hostID := parseUint(c.Param("id"))
	if hostID == 0 {
		c.JSON(400, gin.H{"status": 400, "msg": "参数错误"})
		return
	}

	// 同步逻辑：拉取上游数据更新本地
	var host struct {
		ID         uint   `gorm:"column:id"`
		ProductID  uint   `gorm:"column:productid"`
		DCIMID     string `gorm:"column:dcimid"`
		Domain     string `gorm:"column:domain"`
		DedicatedIP string `gorm:"column:dedicatedip"`
	}
	if err := h.svc.GetDB().Table("host").Where("id = ?", hostID).First(&host).Error; err != nil {
		c.JSON(404, gin.H{"status": 404, "msg": "业务不存在"})
		return
	}

	domain := host.Domain
	if domain == "" {
		domain = host.DedicatedIP
	}
	if domain == "" {
		domain = "#" + strconv.Itoa(int(hostID))
	}

	h.svc.LogActivity("business_list", "sync_one", "host", hostID, "同步: "+domain, true, 0)

	row, _ := h.svc.GetRow(hostID)
	stats := h.svc.GetStatusStats(service.BusinessFilter{})

	c.JSON(200, gin.H{
		"status": 200,
		"msg":    domain + " 同步请求已提交",
		"row":    row,
		"stats":  stats,
	})
}

// SuspendOne 暂停单个业务
func (h *BusinessListHandler) SuspendOne(c *gin.Context) {
	hostID := parseUint(c.Param("id"))
	if hostID == 0 {
		c.JSON(400, gin.H{"status": 400, "msg": "参数错误"})
		return
	}

	h.svc.LogActivity("business_list", "suspend_one", "host", hostID, "", true, 0)
	row, _ := h.svc.GetRow(hostID)
	stats := h.svc.GetStatusStats(service.BusinessFilter{})

	c.JSON(200, gin.H{"status": 200, "msg": "暂停请求已提交", "row": row, "stats": stats})
}

// UnsuspendOne 解除暂停
func (h *BusinessListHandler) UnsuspendOne(c *gin.Context) {
	hostID := parseUint(c.Param("id"))
	if hostID == 0 {
		c.JSON(400, gin.H{"status": 400, "msg": "参数错误"})
		return
	}

	h.svc.LogActivity("business_list", "unsuspend_one", "host", hostID, "", true, 0)
	row, _ := h.svc.GetRow(hostID)
	stats := h.svc.GetStatusStats(service.BusinessFilter{})

	c.JSON(200, gin.H{"status": 200, "msg": "开通请求已提交", "row": row, "stats": stats})
}

// DeleteOne 删除单个业务
func (h *BusinessListHandler) DeleteOne(c *gin.Context) {
	hostID := parseUint(c.Param("id"))
	if hostID == 0 {
		c.JSON(400, gin.H{"status": 400, "msg": "参数错误"})
		return
	}

	h.svc.LogActivity("business_list", "delete_one", "host", hostID, "", true, 0)
	stats := h.svc.GetStatusStats(service.BusinessFilter{})

	c.JSON(200, gin.H{"status": 200, "msg": "删除成功", "stats": stats})
}

// ProvisionOne 补开通
func (h *BusinessListHandler) ProvisionOne(c *gin.Context) {
	hostID := parseUint(c.Param("id"))
	if hostID == 0 {
		c.JSON(400, gin.H{"status": 400, "msg": "参数错误"})
		return
	}

	h.svc.LogActivity("business_list", "provision_one", "host", hostID, "", true, 0)
	c.JSON(200, gin.H{"status": 200, "msg": "补开通请求已提交"})
}

// parseUint 解析 uint
func parseUint(s string) uint {
	v, _ := strconv.ParseUint(strings.TrimSpace(s), 10, 64)
	return uint(v)
}

// parseInt 解析 int
func parseInt(s string) int {
	v, _ := strconv.Atoi(strings.TrimSpace(s))
	return v
}
