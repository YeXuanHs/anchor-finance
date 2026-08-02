package handler

import (
	"strconv"
	"time"

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

// GetYearIncomeStatistics 年度收入统计（按货币、按月分组）
// GET /admin/reports/year-income-statistics?page=1&page_size=20
func (h *ReportHandler) GetYearIncomeStatistics(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	type MonthIncome struct {
		Date     string  `json:"date"`
		Income   float64 `json:"income"`
		Expenses float64 `json:"expenses"`
		Balance  float64 `json:"last"`
	}
	type CurrencyIncome struct {
		Currency        string        `json:"currency"`
		IsDefault       bool          `json:"is_default_currency"`
		YearCount       []MonthIncome `json:"year_count"`
		YearCountTotal  int64         `json:"year_count_num"`
	}

	var currencies []CurrencyIncome
	h.db.Raw(`SELECT code as currency, is_default as is_default FROM currencies`).Scan(&currencies)

	offset := (page - 1) * pageSize
	for i := range currencies {
		var rows []MonthIncome
		h.db.Raw(`
			SELECT DATE_FORMAT(FROM_UNIXTIME(pay_time), '%Y年%m月') as date,
			       SUM(amount_in) as income,
			       SUM(amount_out) as expenses,
			       SUM(amount_in - amount_out) as last
			FROM accounts
			WHERE delete_time = 0 AND currency = ?
			GROUP BY date
			ORDER BY pay_time DESC
			LIMIT ? OFFSET ?
		`, currencies[i].Currency, pageSize, offset).Scan(&rows)

		var total int64
		h.db.Raw(`
			SELECT COUNT(*) FROM (
				SELECT DATE_FORMAT(FROM_UNIXTIME(pay_time), '%Y年%m月') as date
				FROM accounts
				WHERE delete_time = 0 AND currency = ?
				GROUP BY date
			) t
		`, currencies[i].Currency).Scan(&total)

		currencies[i].YearCount = rows
		currencies[i].YearCountTotal = total
	}

	// 返回默认货币的数据
	var result interface{}
	for _, c := range currencies {
		if c.IsDefault {
			result = gin.H{
				"year_count":     c.YearCount,
				"year_count_num": c.YearCountTotal,
			}
			break
		}
	}
	if result == nil {
		result = gin.H{}
	}

	response.Success(c, result)
}

// GetYearIncomeStatisticsForChart 年度收入统计图表数据（多年趋势对比）
// GET /admin/reports/year-income-statistics-chart
func (h *ReportHandler) GetYearIncomeStatisticsForChart(c *gin.Context) {
	type CurrencyAll struct {
		Income   float64 `json:"income"`
		Expenses float64 `json:"expenses"`
		Balance  float64 `json:"last"`
		Code     string  `json:"currency_code"`
		Prefix   string  `json:"currency_prefix"`
	}
	type YearData struct {
		Year   string    `json:"year"`
		Data   []float64 `json:"data"`
		Prefix string    `json:"prefix"`
	}
	type ChartResult struct {
		Years []string   `json:"years"`
		List  []YearData `json:"list"`
	}

	var currencies []struct {
		Code      string `json:"code"`
		Prefix    string `json:"prefix"`
		IsDefault bool   `json:"is_default"`
	}
	h.db.Raw(`SELECT code, prefix, is_default FROM currencies`).Scan(&currencies)

	var result interface{}
	for _, currency := range currencies {
		// 获取汇总数据
		var all CurrencyAll
		h.db.Raw(`
			SELECT SUM(amount_in) as income, SUM(amount_out) as expenses,
			       SUM(amount_in - amount_out) as last
			FROM accounts
			WHERE delete_time = 0 AND currency = ?
		`, currency.Code).Scan(&all)
		all.Code = currency.Code
		all.Prefix = currency.Prefix

		// 获取多年趋势
		var years []string
		h.db.Raw(`
			SELECT DISTINCT DATE_FORMAT(FROM_UNIXTIME(pay_time), '%Y') as year
			FROM accounts
			WHERE delete_time = 0 AND refund = 0
			ORDER BY year ASC
		`).Scan(&years)

		var yearList []YearData
		for _, year := range years {
			var data []float64
			// 预填充12个月为0
			monthData := make(map[int]float64)
			type MonthSum struct {
				Month   int     `json:"month"`
				Balance float64 `json:"last"`
			}
			var months []MonthSum
			h.db.Raw(`
				SELECT MONTH(FROM_UNIXTIME(pay_time)) as month,
				       SUM(amount_in - amount_out) as last
				FROM accounts
				WHERE delete_time = 0 AND currency = ?
				  AND YEAR(FROM_UNIXTIME(pay_time)) = ?
				GROUP BY month
				ORDER BY month ASC
			`, currency.Code, year).Scan(&months)

			for _, m := range months {
				monthData[m.Month] = m.Balance
			}
			for i := 1; i <= 12; i++ {
				data = append(data, monthData[i])
			}

			yearList = append(yearList, YearData{
				Year:   year + "年",
				Data:   data,
				Prefix: currency.Prefix,
			})
		}

		if currency.IsDefault {
			result = gin.H{
				"all":   all,
				"chart": ChartResult{Years: years, List: yearList},
			}
			break
		}
	}

	response.Success(c, result)
}

// GetNewClientStatistics 新客户统计
// GET /admin/reports/new-client-statistics?year=2026&month=8
func (h *ReportHandler) GetNewClientStatistics(c *gin.Context) {
	year := c.DefaultQuery("year", "")
	month := c.DefaultQuery("month", "")
	if year == "" {
		year = time.Now().Format("2006")
	}
	if month == "" {
		month = time.Now().Format("01")
	}

	type DayStats struct {
		Day               string `json:"day_string"`
		NewClients        int64  `json:"new_clients_count"`
		NewOrders         int64  `json:"new_order_count"`
		CompleteOrders    int64  `json:"complete_order_count"`
		NewTickets        int64  `json:"new_ticket_count"`
		ReplyTickets      int64  `json:"reply_ticket_count"`
		CancelRequests    int64  `json:"cancel_requests_count"`
	}
	type YearOption struct {
		Label string `json:"label"`
		Value string `json:"value"`
	}

	// 获取年份列表
	var minYear string
	h.db.Raw(`SELECT MIN(DATE_FORMAT(FROM_UNIXTIME(created_at), '%Y')) FROM users WHERE status = 1`).Scan(&minYear)
	if minYear == "" {
		minYear = time.Now().Format("2006")
	}
	var yearOptions []YearOption
	minY, _ := strconv.Atoi(minYear)
	nowY := time.Now().Year()
	for y := minY; y <= nowY; y++ {
		yearOptions = append(yearOptions, YearOption{
			Label: strconv.Itoa(y) + "年",
			Value: strconv.Itoa(y),
		})
	}

	// 计算月份天数
	yearInt, _ := strconv.Atoi(year)
	monthInt, _ := strconv.Atoi(month)
	daysInMonth := time.Date(yearInt, time.Month(monthInt+1), 0, 0, 0, 0, 0, time.UTC).Day()

	var dayStats []DayStats
	for i := 1; i <= daysInMonth; i++ {
		startTime := time.Date(yearInt, time.Month(monthInt), i, 0, 0, 0, 0, time.UTC)
		endTime := time.Date(yearInt, time.Month(monthInt), i, 23, 59, 59, 0, time.UTC)

		startTS := startTime.Unix()
		endTS := endTime.Unix()

		var day DayStats
		day.Day = month + "月" + strconv.Itoa(i) + "日"

		// 新客户数
		h.db.Raw(`SELECT COUNT(*) FROM users WHERE created_at BETWEEN ? AND ? AND status = 1`, startTS, endTS).Scan(&day.NewClients)

		// 新订单数
		h.db.Raw(`SELECT COUNT(*) FROM orders WHERE created_at BETWEEN ? AND ?`, startTS, endTS).Scan(&day.NewOrders)

		// 完成订单数（已支付）
		h.db.Raw(`SELECT COUNT(*) FROM orders WHERE paid_at IS NOT NULL AND paid_at BETWEEN ? AND ?`, startTS, endTS).Scan(&day.CompleteOrders)

		// 新工单数
		h.db.Raw(`SELECT COUNT(*) FROM tickets WHERE created_at BETWEEN ? AND ?`, startTS, endTS).Scan(&day.NewTickets)

		// 回复工单数（非初始状态的工单）
		h.db.Raw(`SELECT COUNT(DISTINCT ticket_id) FROM ticket_replies WHERE created_at BETWEEN ? AND ?`, startTS, endTS).Scan(&day.ReplyTickets)

		// 取消请求数
		h.db.Raw(`SELECT COUNT(*) FROM cancel_requests WHERE created_at BETWEEN ? AND ?`, startTS, endTS).Scan(&day.CancelRequests)

		dayStats = append(dayStats, day)
	}

	response.Success(c, gin.H{
		"client_data":            dayStats,
		"new_client_years_group": yearOptions,
	})
}

// GetRevenueRanking 客户收入排名（前10名客户，按收入减支出排序）
// GET /admin/reports/revenue-ranking
func (h *ReportHandler) GetRevenueRanking(c *gin.Context) {
	type ClientRanking struct {
		ID           uint    `json:"id"`
		Username     string  `json:"username"`
		CompanyName  string  `json:"companyname"`
		IncomeSum    float64 `json:"income_sum"`
		ExpenseSum   float64 `json:"expense_sum"`
		Balance      float64 `json:"last"`
		CurrencyCode string  `json:"currency"`
		Prefix       string  `json:"prefix"`
		Suffix       string  `json:"suffix"`
	}

	var rankings []ClientRanking
	h.db.Raw(`
		SELECT u.id, u.username, u.nickname as companyname,
		       COALESCE(SUM(a.amount_in), 0) as income_sum,
		       COALESCE(SUM(a.amount_out), 0) as expense_sum,
		       COALESCE(SUM(a.amount_in - a.amount_out), 0) as balance,
		       a.currency as currency_code,
		       IFNULL(c.symbol, '¥') as prefix,
		       '' as suffix
		FROM users u
		JOIN accounts a ON a.uid = u.id
		LEFT JOIN currencies c ON c.code = a.currency
		WHERE u.status = 1 AND a.delete_time = 0
		GROUP BY u.id, a.currency
		HAVING balance > 0
		ORDER BY balance DESC
		LIMIT 30
	`).Scan(&rankings)

	// 按用户ID汇总（取默认货币）
	userBalanceMap := make(map[uint]*ClientRanking)
	for _, r := range rankings {
		if existing, ok := userBalanceMap[r.ID]; ok {
			existing.IncomeSum += r.IncomeSum
			existing.ExpenseSum += r.ExpenseSum
			existing.Balance += r.Balance
		} else {
			copy := r
			userBalanceMap[r.ID] = &copy
		}
	}

	var result []ClientRanking
	for _, v := range userBalanceMap {
		result = append(result, *v)
	}

	// 按balance降序排序
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j].Balance > result[i].Balance {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	// 只取前10
	if len(result) > 10 {
		result = result[:10]
	}

	// 格式化金额
	for i := range result {
		result[i].IncomeSum = float64(int(result[i].IncomeSum*100)) / 100
		result[i].ExpenseSum = float64(int(result[i].ExpenseSum*100)) / 100
		result[i].Balance = float64(int(result[i].Balance*100)) / 100
	}

	response.Success(c, result)
}
