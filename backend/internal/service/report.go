package service

import (
	"time"

	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// DailyReport 每日统计报表
type DailyReport struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Date           string    `gorm:"type:varchar(10);uniqueIndex" json:"date"`
	NewUsers       int       `gorm:"default:0" json:"new_users"`
	NewOrders      int       `gorm:"default:0" json:"new_orders"`
	PaidOrders     int       `gorm:"default:0" json:"paid_orders"`
	Revenue        float64   `gorm:"type:decimal(12,2);default:0" json:"revenue"`
	Refunds        float64   `gorm:"type:decimal(12,2);default:0" json:"refunds"`
	NewTickets     int       `gorm:"default:0" json:"new_tickets"`
	ActiveProducts int       `gorm:"default:0" json:"active_products"`
	CreatedAt      time.Time `json:"created_at"`
}

// MonthlyReport 月度统计
type MonthlyReport struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Month      string    `gorm:"type:varchar(7);uniqueIndex" json:"month"`
	NewUsers   int       `gorm:"default:0" json:"new_users"`
	TotalUsers int       `gorm:"default:0" json:"total_users"`
	NewOrders  int       `gorm:"default:0" json:"new_orders"`
	PaidOrders int       `gorm:"default:0" json:"paid_orders"`
	Revenue    float64   `gorm:"type:decimal(12,2);default:0" json:"revenue"`
	Expenses   float64   `gorm:"type:decimal(12,2);default:0" json:"expenses"`
	Profit     float64   `gorm:"type:decimal(12,2);default:0" json:"profit"`
	CreatedAt  time.Time `json:"created_at"`
}

type DashboardSummary struct {
	TotalUsers   int64   `json:"total_users"`
	TotalOrders  int64   `json:"total_orders"`
	TotalRevenue float64 `json:"total_revenue"`
	TodayUsers   int64   `json:"today_users"`
	TodayOrders  int64   `json:"today_orders"`
	TodayRevenue float64 `json:"today_revenue"`
	ActiveTickets int64  `json:"active_tickets"`
}

type ChartPoint struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
	Count  int     `json:"count"`
}

type UserStatsPoint struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type OrderStatsPoint struct {
	Date       string  `json:"date"`
	Count      int     `json:"count"`
	Revenue    float64 `json:"revenue"`
	PaidCount  int     `json:"paid_count"`
}

type ReportService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewReportService(db *gorm.DB, log *logger.Logger) *ReportService {
	return &ReportService{db: db, log: log}
}

// GenerateDailyReport generates a daily statistics report for the given date.
func (s *ReportService) GenerateDailyReport(date string) (*DailyReport, error) {
	startDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}
	endDate := startDate.AddDate(0, 0, 1)

	var newUsers int64
	s.db.Model(&User{}).Where("created_at >= ? AND created_at < ?", startDate, endDate).Count(&newUsers)

	var newOrders int64
	s.db.Model(&Order{}).Where("created_at >= ? AND created_at < ?", startDate, endDate).Count(&newOrders)

	var paidOrders int64
	s.db.Model(&Order{}).Where("status = 1 AND updated_at >= ? AND updated_at < ?", startDate, endDate).Count(&paidOrders)

	var revenue float64
	s.db.Model(&Order{}).Where("status = 1 AND updated_at >= ? AND updated_at < ?", startDate, endDate).
		Select("COALESCE(SUM(total_price), 0)").Scan(&revenue)

	var activeProducts int64
	s.db.Model(&Product{}).Where("status = 1").Count(&activeProducts)

	report := &DailyReport{
		Date:           date,
		NewUsers:       int(newUsers),
		NewOrders:      int(newOrders),
		PaidOrders:     int(paidOrders),
		Revenue:        revenue,
		ActiveProducts: int(activeProducts),
	}

	if err := s.db.Where("date = ?", date).Assign(report).FirstOrCreate(report).Error; err != nil {
		return nil, err
	}

	s.log.Infof("daily report generated: %s", date)
	return report, nil
}

// GetDashboard returns summary stats for the admin dashboard.
func (s *ReportService) GetDashboard() (*DashboardSummary, error) {
	summary := &DashboardSummary{}

	s.db.Model(&User{}).Count(&summary.TotalUsers)
	s.db.Model(&Order{}).Count(&summary.TotalOrders)
	s.db.Model(&Order{}).Where("status = 1").Select("COALESCE(SUM(total_price), 0)").Scan(&summary.TotalRevenue)

	today := time.Now().Format("2006-01-02")
	startToday, _ := time.Parse("2006-01-02", today)
	endToday := startToday.AddDate(0, 0, 1)

	s.db.Model(&User{}).Where("created_at >= ? AND created_at < ?", startToday, endToday).Count(&summary.TodayUsers)
	s.db.Model(&Order{}).Where("created_at >= ? AND created_at < ?", startToday, endToday).Count(&summary.TodayOrders)
	s.db.Model(&Order{}).Where("status = 1 AND updated_at >= ? AND updated_at < ?", startToday, endToday).
		Select("COALESCE(SUM(total_price), 0)").Scan(&summary.TodayRevenue)

	// Active tickets (status = 0 = open)
	s.db.Model(&struct{ ID uint }{}).Table("tickets").Where("status = 0").Count(&summary.ActiveTickets)

	return summary, nil
}

// GetRevenueChart returns revenue chart data for a date range.
func (s *ReportService) GetRevenueChart(startDate, endDate string) ([]ChartPoint, error) {
	var reports []DailyReport
	query := s.db.Order("date ASC")
	if startDate != "" {
		query = query.Where("date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("date <= ?", endDate)
	}
	if err := query.Find(&reports).Error; err != nil {
		return nil, err
	}

	points := make([]ChartPoint, len(reports))
	for i, r := range reports {
		points[i] = ChartPoint{
			Date:   r.Date,
			Amount: r.Revenue,
			Count:  r.PaidOrders,
		}
	}
	return points, nil
}

// GetUserStats returns user registration trends.
func (s *ReportService) GetUserStats(days int) ([]UserStatsPoint, error) {
	if days <= 0 {
		days = 30
	}

	var results []UserStatsPoint
	startDate := time.Now().AddDate(0, 0, -days)

	rows, err := s.db.Raw(`
		SELECT DATE(created_at) AS date, COUNT(*) AS count
		FROM users
		WHERE created_at >= ? AND deleted_at IS NULL
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`, startDate).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p UserStatsPoint
		if err := rows.Scan(&p.Date, &p.Count); err != nil {
			return nil, err
		}
		results = append(results, p)
	}
	return results, nil
}

// GetMonthlyReport aggregates daily reports for a given month.
func (s *ReportService) GetMonthlyReport(month string) (*MonthlyReport, error) {
	var dailyReports []DailyReport
	if err := s.db.Where("date LIKE ?", month+"%").Order("date ASC").Find(&dailyReports).Error; err != nil {
		return nil, err
	}

	report := &MonthlyReport{Month: month}
	for _, r := range dailyReports {
		report.NewUsers += r.NewUsers
		report.NewOrders += r.NewOrders
		report.PaidOrders += r.PaidOrders
		report.Revenue += r.Revenue
	}
	return report, nil
}

// GetOrderStats returns order statistics for a date range.
func (s *ReportService) GetOrderStats(days int) ([]OrderStatsPoint, error) {
	if days <= 0 {
		days = 30
	}

	var results []OrderStatsPoint
	startDate := time.Now().AddDate(0, 0, -days)

	rows, err := s.db.Raw(`
		SELECT DATE(created_at) AS date,
		       COUNT(*) AS count,
		       COALESCE(SUM(CASE WHEN status = 1 THEN total_price ELSE 0 END), 0) AS revenue,
		       SUM(CASE WHEN status = 1 THEN 1 ELSE 0 END) AS paid_count
		FROM orders
		WHERE created_at >= ? AND deleted_at IS NULL
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`, startDate).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p OrderStatsPoint
		if err := rows.Scan(&p.Date, &p.Count, &p.Revenue, &p.PaidCount); err != nil {
			return nil, err
		}
		results = append(results, p)
	}
	return results, nil
}
