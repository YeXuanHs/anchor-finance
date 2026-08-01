package service

import (
	"fmt"
	"strings"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type SystemLogService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewSystemLogService(db *gorm.DB, log *logger.Logger) *SystemLogService {
	return &SystemLogService{db: db, log: log}
}

func (s *SystemLogService) Record(level, module, message string, userID uint, ip string, details string) {
	entry := model.SystemLog{
		Level: level, Module: module, Message: message,
		UserID: userID, IP: ip, Details: details,
	}
	s.db.Create(&entry)
}

func (s *SystemLogService) List(page, pageSize int, level, module, keyword string, startTime, endTime *time.Time) ([]model.SystemLog, int64, error) {
	var items []model.SystemLog
	var total int64
	query := s.db.Model(&model.SystemLog{})
	if level != "" {
		query = query.Where("level = ?", level)
	}
	if module != "" {
		query = query.Where("module = ?", module)
	}
	if keyword != "" {
		query = query.Where("message LIKE ?", "%"+keyword+"%")
	}
	if startTime != nil {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime != nil {
		query = query.Where("created_at <= ?", endTime)
	}
	query.Count(&total)
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&items)
	return items, total, nil
}

func (s *SystemLogService) GetByID(id uint) (*model.SystemLog, error) {
	var item model.SystemLog
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *SystemLogService) Delete(id uint) error {
	return s.db.Delete(&model.SystemLog{}, id).Error
}

func (s *SystemLogService) Cleanup(days int) (int64, error) {
	cutoff := time.Now().AddDate(0, 0, -days)
	result := s.db.Where("created_at < ?", cutoff).Delete(&model.SystemLog{})
	return result.RowsAffected, result.Error
}

func (s *SystemLogService) GetStats(days int) (map[string]interface{}, error) {
	cutoff := time.Now().AddDate(0, 0, -days)
	var total int64
	s.db.Model(&model.SystemLog{}).Where("created_at >= ?", cutoff).Count(&total)
	return map[string]interface{}{"total": total, "days": days}, nil
}

func (s *SystemLogService) GetLevelStats(days int) ([]map[string]interface{}, error) {
	cutoff := time.Now().AddDate(0, 0, -days)
	var results []map[string]interface{}
	s.db.Model(&model.SystemLog{}).Where("created_at >= ?", cutoff).
		Select("level, COUNT(*) as count").Group("level").Scan(&results)
	return results, nil
}

func (s *SystemLogService) GetModuleStats(days int) ([]map[string]interface{}, error) {
	cutoff := time.Now().AddDate(0, 0, -days)
	var results []map[string]interface{}
	s.db.Model(&model.SystemLog{}).Where("created_at >= ?", cutoff).
		Select("module, COUNT(*) as count").Group("module").Scan(&results)
	return results, nil
}

func (s *SystemLogService) Export(level, module string, startTime, endTime *time.Time) (string, error) {
	var logs []model.SystemLog
	query := s.db.Model(&model.SystemLog{})
	if level != "" {
		query = query.Where("level = ?", level)
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
	query.Order("created_at DESC").Limit(10000).Find(&logs)

	var buf strings.Builder
	buf.WriteString("ID,Level,Module,Message,IP,User,Data,CreatedAt\n")
	for _, l := range logs {
		msg := strings.ReplaceAll(l.Message, "\"", "\"\"")
		data := strings.ReplaceAll(l.Data, "\"", "\"\"")
		buf.WriteString(fmt.Sprintf("%d,%s,%s,\"%s\",%s,%s,\"%s\",%s\n",
			l.ID, l.Level, l.Module, msg, l.IP, l.User, data, l.CreatedAt.Format("2006-01-02 15:04:05")))
	}
	return buf.String(), nil
}

func (s *SystemLogService) ClearByLevel(level string) (int64, error) {
	result := s.db.Where("level = ?", level).Delete(&model.SystemLog{})
	return result.RowsAffected, result.Error
}

// ─── LogRecord Admin Methods (from zjmf LogRecordController) ───

// AdminGetSystemLog returns paginated activity logs.
func (s *SystemLogService) AdminGetSystemLog(page, pageSize int, searchName, searchDesc, searchIP, searchTime string, isSystem bool) ([]map[string]interface{}, int64, error) {
	query := s.db.Table("activity_log")
	if isSystem {
		query = query.Where("user = ?", "System")
	} else {
		query = query.Where("user != ?", "System")
	}
	if searchName != "" {
		query = query.Where("user LIKE ?", "%"+searchName+"%")
	}
	if searchDesc != "" {
		query = query.Where("description LIKE ?", "%"+searchDesc+"%")
	}
	if searchIP != "" {
		query = query.Where("ipaddr LIKE ?", "%"+searchIP+"%")
	}
	if searchTime != "" {
		// Filter by date
		query = query.Where("create_time >= ? AND create_time < ?",
			searchTime+" 00:00:00", searchTime+" 23:59:59")
	}

	var total int64
	query.Count(&total)

	var items []map[string]interface{}
	offset := (page - 1) * pageSize
	query.Select("create_time, id, ipaddr, description, uid, user, port").
		Order("id DESC").Offset(offset).Limit(pageSize).Find(&items)

	return items, total, nil
}

// AdminGetAdminLog returns paginated admin login logs.
func (s *SystemLogService) AdminGetAdminLog(page, pageSize int, searchName, searchIP, searchTime string) ([]map[string]interface{}, int64, error) {
	query := s.db.Table("admin_log")
	if searchName != "" {
		query = query.Where("admin_username LIKE ?", "%"+searchName+"%")
	}
	if searchIP != "" {
		query = query.Where("ipaddress LIKE ?", "%"+searchIP+"%")
	}
	if searchTime != "" {
		query = query.Where("logintime >= ? AND logintime < ?",
			searchTime+" 00:00:00", searchTime+" 23:59:59")
	}

	var total int64
	query.Count(&total)

	var items []map[string]interface{}
	offset := (page - 1) * pageSize
	query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&items)

	return items, total, nil
}

// AdminGetNotifyLog returns paginated notification logs.
func (s *SystemLogService) AdminGetNotifyLog(page, pageSize int, message, logType, searchTime string) ([]map[string]interface{}, int64, error) {
	query := s.db.Table("notify_log")
	if message != "" {
		query = query.Where("message LIKE ?", "%"+message+"%")
	}
	if logType != "" {
		query = query.Where("type LIKE ?", "%"+logType+"%")
	}
	if searchTime != "" {
		query = query.Where("create_time >= ? AND create_time < ?",
			searchTime+" 00:00:00", searchTime+" 23:59:59")
	}

	var total int64
	query.Count(&total)

	var items []map[string]interface{}
	offset := (page - 1) * pageSize
	query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&items)

	return items, total, nil
}

// AdminGetEmailLog returns paginated email logs.
func (s *SystemLogService) AdminGetEmailLog(page, pageSize int, subject, username, searchTime string, uid uint) ([]map[string]interface{}, int64, error) {
	query := s.db.Table("email_log a").
		Select("a.id, a.to, a.create_time, a.subject, a.status, a.fail_reason, a.ip, b.username").
		Joins("LEFT JOIN clients b ON b.id = a.uid")
	if uid > 0 {
		query = query.Where("a.uid = ?", uid).Where("a.is_admin = 0")
	}
	if subject != "" {
		query = query.Where("a.subject LIKE ?", "%"+subject+"%")
	}
	if username != "" {
		query = query.Where("a.to LIKE ?", "%"+username+"%")
	}
	if searchTime != "" {
		query = query.Where("a.create_time >= ? AND a.create_time < ?",
			searchTime+" 00:00:00", searchTime+" 23:59:59")
	}

	var total int64
	query.Count(&total)

	var items []map[string]interface{}
	offset := (page - 1) * pageSize
	query.Order("a.id DESC").Offset(offset).Limit(pageSize).Find(&items)

	return items, total, nil
}

// AdminGetEmailDetail returns email detail by ID.
func (s *SystemLogService) AdminGetEmailDetail(id uint) (map[string]interface{}, error) {
	var detail map[string]interface{}
	if err := s.db.Table("email_log").Where("id = ?", id).Find(&detail).Error; err != nil {
		return nil, err
	}
	if detail == nil {
		return nil, fmt.Errorf("email not found")
	}
	return detail, nil
}

// AdminGetWechatLog returns paginated WeChat logs.
func (s *SystemLogService) AdminGetWechatLog(page, pageSize int) ([]map[string]interface{}, int64, error) {
	var total int64
	s.db.Table("notify_log").Where("type = ?", "wechat").Count(&total)

	var items []map[string]interface{}
	offset := (page - 1) * pageSize
	s.db.Table("notify_log").Where("type = ?", "wechat").
		Order("id DESC").Offset(offset).Limit(pageSize).Find(&items)

	return items, total, nil
}

// AdminGetSmsLog returns paginated SMS logs.
func (s *SystemLogService) AdminGetSmsLog(page, pageSize int, phone, username, searchTime string, uid uint) ([]map[string]interface{}, int64, error) {
	query := s.db.Table("message_log a").
		Select("a.id, a.uid, a.create_time, a.content, a.fail_reason, a.status, a.phone, a.phone_code, b.username, a.ip").
		Joins("LEFT JOIN clients b ON b.id = a.uid")
	if uid > 0 {
		query = query.Where("a.uid = ?", uid)
	}
	if phone != "" {
		query = query.Where("a.phone LIKE ?", "%"+phone+"%")
	}
	if username != "" {
		query = query.Where("b.username LIKE ?", "%"+username+"%")
	}
	if searchTime != "" {
		query = query.Where("a.create_time >= ? AND a.create_time < ?",
			searchTime+" 00:00:00", searchTime+" 23:59:59")
	}

	var total int64
	query.Count(&total)

	var items []map[string]interface{}
	offset := (page - 1) * pageSize
	query.Order("a.id DESC").Offset(offset).Limit(pageSize).Find(&items)

	return items, total, nil
}

// AdminGetSystemMessageLog returns paginated system message logs.
func (s *SystemLogService) AdminGetSystemMessageLog(page, pageSize int, uid uint, keywords, username, readType, searchTimeStart, searchTimeEnd string) ([]map[string]interface{}, int64, error) {
	query := s.db.Table("system_message sm").
		Select("sm.*, c.username, c.phonenumber, c.email").
		Joins("JOIN clients c ON c.id = sm.uid")
	if uid > 0 {
		query = query.Where("sm.uid = ?", uid)
	}
	if keywords != "" {
		query = query.Where("sm.title LIKE ?", "%"+keywords+"%")
	}
	if username != "" {
		query = query.Where("c.username LIKE ?", "%"+username+"%")
	}
	if readType == "0" {
		query = query.Where("sm.read_time = 0")
	} else if readType == "1" {
		query = query.Where("sm.read_time > 0")
	}
	if searchTimeStart != "" {
		query = query.Where("sm.create_time >= ?", searchTimeStart)
	}
	if searchTimeEnd != "" {
		query = query.Where("sm.create_time <= ?", searchTimeEnd)
	}

	var total int64
	query.Count(&total)

	var items []map[string]interface{}
	offset := (page - 1) * pageSize
	query.Order("sm.id DESC").Offset(offset).Limit(pageSize).Find(&items)

	return items, total, nil
}

// AdminGetApiLog returns paginated API resource logs.
func (s *SystemLogService) AdminGetApiLog(page, pageSize int, keywords string, uid uint, searchTime string) ([]map[string]interface{}, int64, error) {
	query := s.db.Table("api_resource_log a").
		Select("a.id, a.create_time, a.description, a.ip, b.username, a.uid").
		Joins("LEFT JOIN clients b ON a.uid = b.id")
	if keywords != "" {
		query = query.Where("a.description LIKE ? OR a.ip LIKE ?", "%"+keywords+"%", "%"+keywords+"%")
	}
	if uid > 0 {
		query = query.Where("a.uid = ?", uid)
	}
	if searchTime != "" {
		query = query.Where("a.create_time >= ? AND a.create_time < ?",
			searchTime+" 00:00:00", searchTime+" 23:59:59")
	}

	var total int64
	query.Count(&total)

	var items []map[string]interface{}
	offset := (page - 1) * pageSize
	query.Order("a.id DESC").Offset(offset).Limit(pageSize).Find(&items)

	return items, total, nil
}

// AdminGetLogCount returns log count for a given type.
func (s *SystemLogService) AdminGetLogCount(logType string) (int64, error) {
	tableMap := map[string]string{
		"system_log":        "activity_log",
		"admin_log":         "admin_log",
		"email_log":         "email_log",
		"sms_log":           "message_log",
		"system_message_log": "system_message",
		"cron_system_log":   "activity_log",
		"api_log":           "api_resource_log",
	}
	table, ok := tableMap[logType]
	if !ok {
		return 0, nil
	}

	var count int64
	query := s.db.Table(table)
	if logType == "system_log" {
		query = query.Where("user != ?", "System")
	} else if logType == "cron_system_log" {
		query = query.Where("user = ?", "System")
	}
	query.Count(&count)
	return count, nil
}

// AdminGetLogCountBefore returns log count before a given time.
func (s *SystemLogService) AdminGetLogCountBefore(logType, timeStr string) (int64, error) {
	tableMap := map[string]struct{ Table, TimeCol string }{
		"system_log":         {"activity_log", "create_time"},
		"admin_log":          {"admin_log", "logintime"},
		"email_log":          {"email_log", "create_time"},
		"sms_log":            {"message_log", "create_time"},
		"system_message_log": {"system_message", "create_time"},
		"cron_system_log":    {"activity_log", "create_time"},
		"api_log":            {"api_resource_log", "create_time"},
	}
	info, ok := tableMap[logType]
	if !ok {
		return 0, fmt.Errorf("invalid log type")
	}

	var count int64
	query := s.db.Table(info.Table)
	if timeStr != "" {
		query = query.Where(info.TimeCol+" <= ?", timeStr)
	}
	if logType == "system_log" {
		query = query.Where("user != ?", "System")
	} else if logType == "cron_system_log" {
		query = query.Where("user = ?", "System")
	}
	query.Count(&count)
	return count, nil
}

// AdminDeleteLog deletes logs by type and optional time filter.
func (s *SystemLogService) AdminDeleteLog(logType, timeStr string) (int64, error) {
	tableMap := map[string]struct{ Table, TimeCol string }{
		"system_log":         {"activity_log", "create_time"},
		"admin_log":          {"admin_log", "logintime"},
		"email_log":          {"email_log", "create_time"},
		"sms_log":            {"message_log", "create_time"},
		"system_message_log": {"system_message", "create_time"},
		"cron_system_log":    {"activity_log", "create_time"},
		"api_log":            {"api_resource_log", "create_time"},
	}
	info, ok := tableMap[logType]
	if !ok {
		return 0, fmt.Errorf("invalid log type")
	}

	query := s.db.Table(info.Table)
	if timeStr != "" {
		query = query.Where(info.TimeCol+" <= ?", timeStr)
	}
	if logType == "system_log" {
		query = query.Where("user != ?", "System")
	} else if logType == "cron_system_log" {
		query = query.Where("user = ?", "System")
	}

	result := query.Delete(nil)
	return result.RowsAffected, result.Error
}
