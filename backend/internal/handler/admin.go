package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminHandler struct {
	db        *gorm.DB
	userSvc   *service.UserService
	orderSvc  *service.OrderService
	invSvc    *service.InvoiceService
	log       *logger.Logger
}

func NewAdminHandler(db *gorm.DB, userSvc *service.UserService, orderSvc *service.OrderService, invSvc *service.InvoiceService, log *logger.Logger) *AdminHandler {
	return &AdminHandler{
		db:       db,
		userSvc:  userSvc,
		orderSvc: orderSvc,
		invSvc:   invSvc,
		log:      log,
	}
}

// Dashboard returns summary stats for the admin dashboard.
func (h *AdminHandler) Dashboard(c *gin.Context) {
	var userCount, orderCount, invoiceCount int64
	h.db.Model(&service.User{}).Count(&userCount)
	h.db.Model(&service.Order{}).Count(&orderCount)
	h.db.Model(&service.Invoice{}).Count(&invoiceCount)

	var totalRevenue float64
	h.db.Model(&service.Invoice{}).Where("status = 1").Select("COALESCE(SUM(amount),0)").Scan(&totalRevenue)

	response.Success(c, gin.H{
		"user_count":    userCount,
		"order_count":   orderCount,
		"invoice_count": invoiceCount,
		"total_revenue": totalRevenue,
	})
}

// Stats returns period-based statistics.
func (h *AdminHandler) Stats(c *gin.Context) {
	period := c.DefaultQuery("period", "today")

	var startClause string
	switch period {
	case "today":
		startClause = "CURRENT_DATE"
	case "week":
		startClause = "CURRENT_DATE - INTERVAL '7 days'"
	case "month":
		startClause = "CURRENT_DATE - INTERVAL '30 days'"
	case "year":
		startClause = "CURRENT_DATE - INTERVAL '1 year'"
	default:
		startClause = "CURRENT_DATE"
	}

	var newUsers, newOrders int64
	h.db.Model(&service.User{}).Where("created_at >= "+startClause).Count(&newUsers)
	h.db.Model(&service.Order{}).Where("created_at >= "+startClause).Count(&newOrders)

	var revenue float64
	h.db.Model(&service.Invoice{}).Where("status = 1 AND paid_at >= "+startClause).
		Select("COALESCE(SUM(amount),0)").Scan(&revenue)

	response.Success(c, gin.H{
		"new_users":  newUsers,
		"new_orders": newOrders,
		"revenue":    revenue,
		"period":     period,
	})
}

type Announcement struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Title   string `gorm:"size:256;not null" json:"title"`
	Content string `gorm:"type:text" json:"content"`
	Status  int    `gorm:"default:1" json:"status"`
}

// GetAnnouncements returns all announcements.
func (h *AdminHandler) GetAnnouncements(c *gin.Context) {
	var items []Announcement
	h.db.Order("id DESC").Find(&items)
	response.Success(c, items)
}

// CreateAnnouncement adds an announcement.
func (h *AdminHandler) CreateAnnouncement(c *gin.Context) {
	var req Announcement
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.db.Create(&req).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, req)
}

// UpdateAnnouncement modifies an announcement.
func (h *AdminHandler) UpdateAnnouncement(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req Announcement
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.db.Model(&Announcement{}).Where("id = ?", id).Updates(&req).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "announcement updated")
}

// DeleteAnnouncement removes an announcement.
func (h *AdminHandler) DeleteAnnouncement(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.db.Delete(&Announcement{}, id).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "announcement deleted")
}

type Setting struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Key   string `gorm:"uniqueIndex;size:128;not null" json:"key"`
	Value string `gorm:"type:text" json:"value"`
	Group string `gorm:"size:64;default:general" json:"group"`
}

// GetSettings returns all settings.
func (h *AdminHandler) GetSettings(c *gin.Context) {
	var items []Setting
	h.db.Find(&items)
	response.Success(c, items)
}

// UpdateSettings batch-updates settings.
func (h *AdminHandler) UpdateSettings(c *gin.Context) {
	var req []Setting
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	for _, s := range req {
		h.db.Where(Setting{Key: s.Key}).Assign(Setting{Value: s.Value, Group: s.Group}).FirstOrCreate(&s)
	}
	response.SuccessMsg(c, "settings updated")
}

// GetSettingsByGroup returns settings filtered by group.
func (h *AdminHandler) GetSettingsByGroup(c *gin.Context) {
	group := c.Param("group")
	var items []Setting
	h.db.Where("`group` = ?", group).Find(&items)
	response.Success(c, items)
}

type LogEntry struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Action    string `gorm:"size:128" json:"action"`
	Detail    string `gorm:"type:text" json:"detail"`
	UserID    uint   `gorm:"index" json:"user_id"`
	IP        string `gorm:"size:64" json:"ip"`
	CreatedAt string `gorm:"autoCreateTime" json:"created_at"`
}

// GetLogs returns system logs with pagination.
func (h *AdminHandler) GetLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var logs []LogEntry
	var total int64

	query := h.db.Model(&LogEntry{})
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs)

	response.SuccessPage(c, logs, total, page, pageSize)
}
