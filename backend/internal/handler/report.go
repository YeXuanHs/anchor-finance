package handler

import (
	"strconv"

	"github.com/anchor-finance/backend/internal/service"
	"github.com/anchor-finance/backend/pkg/logger"
	"github.com/anchor-finance/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	reportSvc *service.ReportService
	log       *logger.Logger
}

func NewReportHandler(reportSvc *service.ReportService, log *logger.Logger) *ReportHandler {
	return &ReportHandler{reportSvc: reportSvc, log: log}
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
