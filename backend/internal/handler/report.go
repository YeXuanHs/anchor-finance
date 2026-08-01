package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReportHandler struct {
	db        *gorm.DB
	reportSvc *service.ReportService
	log       *logger.Logger
}

func NewReportHandler(db *gorm.DB, reportSvc *service.ReportService, log *logger.Logger) *ReportHandler {
	return &ReportHandler{db: db, reportSvc: reportSvc, log: log}
}

// GetDashboard returns dashboard summary stats.
// GET /admin/reports/dashboard
func (h *ReportHandler) GetDashboard(c *gin.Context) {
	summary, err := h.reportSvc.GetDashboard()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, summary)
}

// GetDailyReport returns or generates a daily report.
// GET /admin/reports/daily?date=2024-01-01
func (h *ReportHandler) GetDailyReport(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		response.BadRequest(c, "date is required (format: 2006-01-02)")
		return
	}

	report, err := h.reportSvc.GenerateDailyReport(date)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, report)
}

// GetMonthlyReport returns monthly report data aggregated from daily reports.
// GET /admin/reports/monthly?month=2024-01
func (h *ReportHandler) GetMonthlyReport(c *gin.Context) {
	month := c.Query("month")
	if month == "" {
		response.BadRequest(c, "month is required (format: 2006-01)")
		return
	}

	report, err := h.reportSvc.GetMonthlyReport(month)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, report)
}

// GetRevenueChart returns revenue chart data for a date range.
// GET /admin/reports/revenue?start=2024-01-01&end=2024-01-31
func (h *ReportHandler) GetRevenueChart(c *gin.Context) {
	start := c.Query("start")
	end := c.Query("end")

	points, err := h.reportSvc.GetRevenueChart(start, end)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, points)
}

// GetUserStats returns user registration trends.
// GET /admin/reports/users?days=30
func (h *ReportHandler) GetUserStats(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	stats, err := h.reportSvc.GetUserStats(days)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, stats)
}

// GetOrderStats returns order statistics.
// GET /admin/reports/orders?days=30
func (h *ReportHandler) GetOrderStats(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	stats, err := h.reportSvc.GetOrderStats(days)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, stats)
}

// GetTopClients returns top clients ranked by spending.
// GET /admin/reports/top-clients?limit=10&period=month
func (h *ReportHandler) GetTopClients(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	period := c.DefaultQuery("period", "month")

	data, err := h.reportSvc.GetTopClients(limit, period)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, data)
}

// GetProductIncome returns income breakdown by product.
// GET /admin/reports/product-income?period=month
func (h *ReportHandler) GetProductIncome(c *gin.Context) {
	period := c.DefaultQuery("period", "month")

	data, err := h.reportSvc.GetProductIncome(period)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, data)
}

// BaseInfoModule represents a system info module.
type BaseInfoModule struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Name   string `gorm:"type:varchar(128);not null" json:"name"`
	Desc   string `gorm:"type:varchar(255)" json:"desc"`
	Enable int    `gorm:"default:1" json:"enable"`
	Sort   int    `gorm:"default:0" json:"sort"`
}

func (BaseInfoModule) TableName() string {
	return "base_info"
}

// GetSystemInfoModulesList returns all system info modules sorted by sort order.
// GET /admin/reports/modules
func (h *ReportHandler) GetSystemInfoModulesList(c *gin.Context) {
	var modules []BaseInfoModule
	if err := h.db.Where("delete_time = ?", 0).
		Order("sort ASC, id ASC").
		Find(&modules).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, modules)
}

// UpdateSystemInfoModulesSort updates the sort order of system info modules.
// POST /admin/reports/modules/sort
func (h *ReportHandler) UpdateSystemInfoModulesSort(c *gin.Context) {
	var req struct {
		Modules []struct {
			ID   uint `json:"id"`
			Sort int  `json:"sort"`
		} `json:"modules" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	for _, m := range req.Modules {
		if err := h.db.Model(&BaseInfoModule{}).Where("id = ?", m.ID).
			Update("sort", m.Sort).Error; err != nil {
			response.ServerError(c, err.Error())
			return
		}
	}
	response.SuccessMsg(c, "module sort order updated")
}
