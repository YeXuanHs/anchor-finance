package service

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type LogRecordService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewLogRecordService(db *gorm.DB, log *logger.Logger) *LogRecordService {
	return &LogRecordService{db: db, log: log}
}

type CreateLogRecordRequest struct {
	AdminID       uint   `json:"admin_id" binding:"required"`
	Action        string `json:"action" binding:"required"`
	Module        string `json:"module" binding:"required"`
	TargetID      uint   `json:"target_id"`
	TargetType    string `json:"target_type"`
	Title         string `json:"title"`
	IPAddress     string `json:"ip_address" binding:"required"`
	UserAgent     string `json:"user_agent"`
	RequestMethod string `json:"request_method"`
	RequestPath   string `json:"request_path"`
	Duration      int64  `json:"duration"`
	Remark        string `json:"remark"`
	Status        int16  `json:"status"`
	ErrorMsg      string `json:"error_msg"`
}

// List returns paginated log records with filters.
func (s *LogRecordService) List(page, pageSize int, adminID *uint, action, module, keyword string, status *int16, startTime, endTime *time.Time) ([]model.LogRecord, int64, error) {
	var records []model.LogRecord
	var total int64

	query := s.db.Model(&model.LogRecord{})
	if adminID != nil {
		query = query.Where("admin_id = ?", *adminID)
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if module != "" {
		query = query.Where("module = ?", module)
	}
	if keyword != "" {
		query = query.Where("title LIKE ? OR remark LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if startTime != nil {
		query = query.Where("created_at >= ?", *startTime)
	}
	if endTime != nil {
		query = query.Where("created_at <= ?", *endTime)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := logRecordPaginate(page, pageSize)
	if err := query.Preload("Admin").Offset(offset).Limit(limit).Order("id DESC").Find(&records).Error; err != nil {
		return nil, 0, err
	}
	return records, total, nil
}

// GetByID returns a single log record by ID.
func (s *LogRecordService) GetByID(id uint) (*model.LogRecord, error) {
	var record model.LogRecord
	if err := s.db.Preload("Admin").First(&record, id).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// Create creates a new log record.
func (s *LogRecordService) Create(req CreateLogRecordRequest) (*model.LogRecord, error) {
	record := &model.LogRecord{
		AdminID:       req.AdminID,
		Action:        req.Action,
		Module:        req.Module,
		TargetID:      req.TargetID,
		TargetType:    req.TargetType,
		Title:         req.Title,
		IPAddress:     req.IPAddress,
		UserAgent:     req.UserAgent,
		RequestMethod: req.RequestMethod,
		RequestPath:   req.RequestPath,
		Duration:      req.Duration,
		Remark:        req.Remark,
		Status:        req.Status,
		ErrorMsg:      req.ErrorMsg,
	}

	if record.Status == 0 {
		record.Status = 1
	}

	if err := s.db.Create(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

// QuickLog is a convenience method for quick log creation.
func (s *LogRecordService) QuickLog(adminID uint, action, module, title, ip string) {
	record := &model.LogRecord{
		AdminID:   adminID,
		Action:    action,
		Module:    module,
		Title:     title,
		IPAddress: ip,
		Status:    1,
	}
	if err := s.db.Create(record).Error; err != nil {
		s.log.Errorf("quick log failed: %v", err)
	}
}

// Search searches log records by keyword across multiple fields.
func (s *LogRecordService) Search(page, pageSize int, keyword string) ([]model.LogRecord, int64, error) {
	var records []model.LogRecord
	var total int64

	query := s.db.Model(&model.LogRecord{}).Where(
		"title LIKE ? OR remark LIKE ? OR ip_address LIKE ? OR error_msg LIKE ?",
		"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%",
	)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := logRecordPaginate(page, pageSize)
	if err := query.Preload("Admin").Offset(offset).Limit(limit).Order("id DESC").Find(&records).Error; err != nil {
		return nil, 0, err
	}
	return records, total, nil
}

// Export exports log records to a CSV file.
func (s *LogRecordService) Export(adminID *uint, action, module string, startTime, endTime *time.Time) (string, error) {
	var records []model.LogRecord

	query := s.db.Model(&model.LogRecord{})
	if adminID != nil {
		query = query.Where("admin_id = ?", *adminID)
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if module != "" {
		query = query.Where("module = ?", module)
	}
	if startTime != nil {
		query = query.Where("created_at >= ?", *startTime)
	}
	if endTime != nil {
		query = query.Where("created_at <= ?", *endTime)
	}

	if err := query.Preload("Admin").Order("id DESC").Find(&records).Error; err != nil {
		return "", err
	}

	filename := fmt.Sprintf("log_export_%s.csv", time.Now().Format("20060102150405"))
	filepath := fmt.Sprintf("tmp/%s", filename)

	// Ensure tmp directory exists
	if err := os.MkdirAll("tmp", 0755); err != nil {
		return "", err
	}

	file, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Write BOM for Excel UTF-8 compatibility
	file.WriteString("\xEF\xBB\xBF")

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"ID", "管理员", "操作类型", "模块", "标题", "IP地址", "状态", "备注", "创建时间"}
	if err := writer.Write(header); err != nil {
		return "", err
	}

	// Write data
	for _, r := range records {
		adminName := ""
		if r.Admin.Username != "" {
			adminName = r.Admin.Username
		}
		statusStr := "成功"
		if r.Status == 0 {
			statusStr = "失败"
		}
		row := []string{
			fmt.Sprintf("%d", r.ID),
			adminName,
			r.Action,
			r.Module,
			r.Title,
			r.IPAddress,
			statusStr,
			r.Remark,
			r.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(row); err != nil {
			return "", err
		}
	}

	return filepath, nil
}

// Stats returns log statistics grouped by module and action.
func (s *LogRecordService) Stats(days int) ([]model.LogRecordStat, error) {
	var stats []model.LogRecordStat

	startDate := time.Now().AddDate(0, 0, -days)

	rows, err := s.db.Model(&model.LogRecord{}).
		Select("DATE(created_at) as date, module, action, COUNT(*) as count").
		Where("created_at >= ?", startDate).
		Group("DATE(created_at), module, action").
		Order("date DESC").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var stat model.LogRecordStat
		if err := rows.Scan(&stat.Date, &stat.Module, &stat.Action, &stat.Count); err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}
	return stats, nil
}

// ModuleStats returns log counts grouped by module.
func (s *LogRecordService) ModuleStats(days int) (map[string]int64, error) {
	type result struct {
		Module string
		Count  int64
	}
	var results []result

	startDate := time.Now().AddDate(0, 0, -days)
	if err := s.db.Model(&model.LogRecord{}).
		Select("module, COUNT(*) as count").
		Where("created_at >= ?", startDate).
		Group("module").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, r := range results {
		stats[r.Module] = r.Count
	}
	return stats, nil
}

// ActionStats returns log counts grouped by action.
func (s *LogRecordService) ActionStats(days int) (map[string]int64, error) {
	type result struct {
		Action string
		Count  int64
	}
	var results []result

	startDate := time.Now().AddDate(0, 0, -days)
	if err := s.db.Model(&model.LogRecord{}).
		Select("action, COUNT(*) as count").
		Where("created_at >= ?", startDate).
		Group("action").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, r := range results {
		stats[r.Action] = r.Count
	}
	return stats, nil
}

// Cleanup deletes log records older than the specified days.
func (s *LogRecordService) Cleanup(days int) (int64, error) {
	cutoff := time.Now().AddDate(0, 0, -days)
	result := s.db.Where("created_at < ?", cutoff).Delete(&model.LogRecord{})
	if result.Error != nil {
		return 0, result.Error
	}
	s.log.Infof("cleaned up %d log records older than %d days", result.RowsAffected, days)
	return result.RowsAffected, nil
}

func logRecordPaginate(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	pageSize = int(math.Min(float64(pageSize), 100))
	return (page - 1) * pageSize, pageSize
}
