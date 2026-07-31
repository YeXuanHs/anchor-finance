package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/sms"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// SMSService handles SMS operations including templates, logging, batch sending.
type SMSService struct {
	db     *gorm.DB
	log    *logger.Logger
	smsSnd *sms.Sender
}

func NewSMSService(db *gorm.DB, log *logger.Logger) *SMSService {
	return &SMSService{
		db:     db,
		log:    log,
		smsSnd: sms.NewSender(db),
	}
}

// ─── Operator Detection ───

// DetectOperator detects the mobile operator from a phone number prefix.
func (s *SMSService) DetectOperator(phone string) string {
	if len(phone) < 7 {
		return "unknown"
	}
	prefix := phone[:7]

	// Mobile: 134-139, 150-152, 157-159, 188
	mobilePrefixes := []string{
		"134", "135", "136", "137", "138", "139",
		"150", "151", "152",
		"157", "158", "159",
		"188",
	}
	for _, p := range mobilePrefixes {
		if strings.HasPrefix(prefix, p) {
			return "mobile"
		}
	}

	// Unicom: 130-132, 155-156, 186
	unicomPrefixes := []string{
		"130", "131", "132",
		"155", "156",
		"186",
	}
	for _, p := range unicomPrefixes {
		if strings.HasPrefix(prefix, p) {
			return "unicom"
		}
	}

	// Telecom: 133, 153, 189
	telecomPrefixes := []string{
		"133", "153", "189",
	}
	for _, p := range telecomPrefixes {
		if strings.HasPrefix(prefix, p) {
			return "telecom"
		}
	}

	return "unknown"
}

// ValidatePhone validates a Chinese phone number format (11 digits starting with 1).
func (s *SMSService) ValidatePhone(phone string) bool {
	re := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return re.MatchString(phone)
}

// ─── Template Management ───

// GetTemplates returns all SMS templates.
func (s *SMSService) GetTemplates() ([]model.SMSTemplate, error) {
	var templates []model.SMSTemplate
	if err := s.db.Order("id DESC").Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

// GetTemplateByID returns a single template by ID.
func (s *SMSService) GetTemplateByID(id uint) (*model.SMSTemplate, error) {
	var tpl model.SMSTemplate
	if err := s.db.First(&tpl, id).Error; err != nil {
		return nil, err
	}
	return &tpl, nil
}

// GetTemplateByCode returns a template by its code.
func (s *SMSService) GetTemplateByCode(code string) (*model.SMSTemplate, error) {
	var tpl model.SMSTemplate
	if err := s.db.Where("code = ?", code).First(&tpl).Error; err != nil {
		return nil, err
	}
	return &tpl, nil
}

// CreateTemplate creates a new SMS template.
func (s *SMSService) CreateTemplate(tpl *model.SMSTemplate) error {
	return s.db.Create(tpl).Error
}

// UpdateTemplate updates an existing SMS template.
func (s *SMSService) UpdateTemplate(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.SMSTemplate{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteTemplate soft-deletes an SMS template.
func (s *SMSService) DeleteTemplate(id uint) error {
	return s.db.Delete(&model.SMSTemplate{}, id).Error
}

// ReplaceTemplateParams replaces {1}, {2} or {param_name} placeholders in the template content.
func (s *SMSService) ReplaceTemplateParams(tplCode string, params map[string]string) (string, error) {
	tpl, err := s.GetTemplateByCode(tplCode)
	if err != nil {
		return "", fmt.Errorf("template not found: %s", tplCode)
	}

	if !tpl.Enabled {
		return "", errors.New("template is disabled")
	}

	content := tpl.Content
	for k, v := range params {
		content = strings.ReplaceAll(content, "{"+k+"}", v)
	}

	return content, nil
}

// ─── Enhanced Sending ───

// SendSMSSend sends an SMS using a template with parameters.
func (s *SMSService) SendSMSSend(phone, tplCode string, params map[string]string) (*model.SMSLog, error) {
	if !s.ValidatePhone(phone) {
		return nil, errors.New("invalid phone number format")
	}

	content, err := s.ReplaceTemplateParams(tplCode, params)
	if err != nil {
		return nil, err
	}

	operator := s.DetectOperator(phone)
	logEntry := &model.SMSLog{
		Phone:      phone,
		Content:    content,
		TemplateID: tplCode,
		Params:     marshalParams(params),
		Status:     "pending",
		Operator:   operator,
	}

	if err := s.db.Create(logEntry).Error; err != nil {
		return nil, err
	}

	// Send via provider
	sendErr := s.smsSnd.Send(phone, content)
	now := time.Now()
	logEntry.SentAt = &now

	if sendErr != nil {
		logEntry.Status = "failed"
		logEntry.Response = sendErr.Error()
		s.db.Save(logEntry)
		s.log.Errorf("sms send failed: phone=%s tpl=%s err=%v", phone, tplCode, sendErr)
		return logEntry, sendErr
	}

	logEntry.Status = "sent"
	logEntry.Response = "ok"
	s.db.Save(logEntry)
	s.log.Infof("sms sent: phone=%s tpl=%s operator=%s", phone, tplCode, operator)
	return logEntry, nil
}

// SendBatchSMS sends SMS to multiple phones using a template (async via goroutine).
func (s *SMSService) SendBatchSMS(phones []string, tplCode string, params map[string]string, batchID *uint) {
	go func() {
		for _, phone := range phones {
			if !s.ValidatePhone(phone) {
				s.log.Warnf("skipping invalid phone: %s", phone)
				continue
			}

			content, err := s.ReplaceTemplateParams(tplCode, params)
			if err != nil {
				s.log.Errorf("template replace failed for phone %s: %v", phone, err)
				continue
			}

			operator := s.DetectOperator(phone)
			logEntry := &model.SMSLog{
				Phone:      phone,
				Content:    content,
				TemplateID: tplCode,
				Params:     marshalParams(params),
				Status:     "pending",
				Operator:   operator,
				BatchID:    batchID,
			}
			s.db.Create(logEntry)

			sendErr := s.smsSnd.Send(phone, content)
			now := time.Now()
			logEntry.SentAt = &now

			if sendErr != nil {
				logEntry.Status = "failed"
				logEntry.Response = sendErr.Error()
				s.log.Errorf("batch sms failed: phone=%s err=%v", phone, sendErr)
			} else {
				logEntry.Status = "sent"
				logEntry.Response = "ok"
			}
			s.db.Save(logEntry)
		}
	}()
}

// SendMarketingSMS sends marketing SMS to a target group.
func (s *SMSService) SendMarketingSMS(targetGroup, tplCode string, params map[string]string) (*model.SMSBatch, error) {
	tpl, err := s.GetTemplateByCode(tplCode)
	if err != nil {
		return nil, fmt.Errorf("template not found: %s", tplCode)
	}

	if tpl.Type != "marketing" {
		return nil, errors.New("template is not a marketing template")
	}

	// Create batch
	batch := &model.SMSBatch{
		Name:        fmt.Sprintf("marketing_%s_%d", targetGroup, time.Now().Unix()),
		TemplateID:  tpl.ID,
		TargetGroup: targetGroup,
		Status:      "sending",
	}
	if err := s.db.Create(batch).Error; err != nil {
		return nil, err
	}

	// Get target phones
	phones, err := s.getPhonesByGroup(targetGroup)
	if err != nil {
		s.markBatchFailed(batch.ID, err.Error())
		return nil, err
	}

	s.db.Model(batch).Update("total_count", len(phones))

	// Execute async
	go s.executeBatch(batch, phones, tplCode, params)

	return batch, nil
}

// SendVerifyCode sends a verification code SMS using the verify template.
func (s *SMSService) SendVerifyCode(phone string) (*model.SMSLog, error) {
	if !s.ValidatePhone(phone) {
		return nil, errors.New("invalid phone number format")
	}

	code := generateNumericCode(6)
	params := map[string]string{"code": code}
	return s.SendSMSSend(phone, "verify_code", params)
}

// SendSmsBefore performs pre-send validation checks.
func (s *SMSService) SendSmsBefore(phone, tplCode string) error {
	if !s.ValidatePhone(phone) {
		return errors.New("invalid phone number format")
	}

	tpl, err := s.GetTemplateByCode(tplCode)
	if err != nil {
		return fmt.Errorf("template not found: %s", tplCode)
	}

	if !tpl.Enabled {
		return errors.New("template is disabled")
	}

	return nil
}

// ─── Logging ───

// GetSMSLogs returns paginated SMS logs.
func (s *SMSService) GetSMSLogs(page, pageSize int) ([]model.SMSLog, int64, error) {
	var logs []model.SMSLog
	var total int64

	query := s.db.Model(&model.SMSLog{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

// GetSMSLogByID returns a single SMS log by ID.
func (s *SMSService) GetSMSLogByID(id uint) (*model.SMSLog, error) {
	var logEntry model.SMSLog
	if err := s.db.First(&logEntry, id).Error; err != nil {
		return nil, err
	}
	return &logEntry, nil
}

// GetSMSLogByPhone returns SMS logs for a specific phone number.
func (s *SMSService) GetSMSLogByPhone(phone string, page, pageSize int) ([]model.SMSLog, int64, error) {
	var logs []model.SMSLog
	var total int64

	query := s.db.Model(&model.SMSLog{}).Where("phone = ?", phone)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

// GetSMSLogByUser returns SMS logs for a specific user.
func (s *SMSService) GetSMSLogByUser(userID uint, page, pageSize int) ([]model.SMSLog, int64, error) {
	var logs []model.SMSLog
	var total int64

	query := s.db.Model(&model.SMSLog{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

// LogSMS creates a new SMS log entry.
func (s *SMSService) LogSMS(phone, content, templateID string, params map[string]string, status, response string) (*model.SMSLog, error) {
	logEntry := &model.SMSLog{
		Phone:      phone,
		Content:    content,
		TemplateID: templateID,
		Params:     marshalParams(params),
		Status:     status,
		Response:   response,
		Operator:   s.DetectOperator(phone),
	}
	if status == "sent" {
		now := time.Now()
		logEntry.SentAt = &now
	}
	if err := s.db.Create(logEntry).Error; err != nil {
		return nil, err
	}
	return logEntry, nil
}

// ─── Statistics ───

// SMSStats represents SMS sending statistics.
type SMSStats struct {
	Total    int64 `json:"total"`
	Sent     int64 `json:"sent"`
	Failed   int64 `json:"failed"`
	Pending  int64 `json:"pending"`
}

// OperatorStats represents statistics grouped by operator.
type OperatorStats struct {
	Operator string `json:"operator"`
	Total    int64  `json:"total"`
	Sent     int64  `json:"sent"`
	Failed   int64  `json:"failed"`
}

// GetSMSStats returns SMS sending statistics for a given period.
func (s *SMSService) GetSMSStats(period string) (*SMSStats, error) {
	stats := &SMSStats{}

	var since time.Time
	now := time.Now()
	switch period {
	case "today":
		since = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	case "week":
		since = now.AddDate(0, 0, -7)
	case "month":
		since = now.AddDate(0, -1, 0)
	default:
		since = time.Time{} // all time
	}

	query := s.db.Model(&model.SMSLog{})
	if !since.IsZero() {
		query = query.Where("created_at >= ?", since)
	}

	query.Count(&stats.Total)
	query.Where("status = ?", "sent").Count(&stats.Sent)
	query.Where("status = ?", "failed").Count(&stats.Failed)
	query.Where("status = ?", "pending").Count(&stats.Pending)

	return stats, nil
}

// GetOperatorStats returns SMS statistics grouped by operator.
func (s *SMSService) GetOperatorStats() ([]OperatorStats, error) {
	var results []OperatorStats

	rows, err := s.db.Model(&model.SMSLog{}).
		Select("operator, count(*) as total, sum(case when status = 'sent' then 1 else 0 end) as sent, sum(case when status = 'failed' then 1 else 0 end) as failed").
		Group("operator").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var stat OperatorStats
		if err := rows.Scan(&stat.Operator, &stat.Total, &stat.Sent, &stat.Failed); err != nil {
			return nil, err
		}
		results = append(results, stat)
	}

	return results, nil
}

// ─── Batch Operations ───

// CreateBatch creates a new SMS batch job.
func (s *SMSService) CreateBatch(name string, templateID uint, targetGroup string) (*model.SMSBatch, error) {
	batch := &model.SMSBatch{
		Name:        name,
		TemplateID:  templateID,
		TargetGroup: targetGroup,
		Status:      "pending",
	}
	if err := s.db.Create(batch).Error; err != nil {
		return nil, err
	}
	return batch, nil
}

// GetBatches returns all SMS batches.
func (s *SMSService) GetBatches() ([]model.SMSBatch, error) {
	var batches []model.SMSBatch
	if err := s.db.Order("id DESC").Find(&batches).Error; err != nil {
		return nil, err
	}
	return batches, nil
}

// ExecuteBatch executes an SMS batch job asynchronously.
func (s *SMSService) ExecuteBatch(batchID uint) error {
	var batch model.SMSBatch
	if err := s.db.First(&batch, batchID).Error; err != nil {
		return err
	}

	if batch.Status != "pending" {
		return errors.New("batch is not in pending status")
	}

	tpl, err := s.GetTemplateByID(batch.TemplateID)
	if err != nil {
		return fmt.Errorf("template not found: %d", batch.TemplateID)
	}

	phones, err := s.getPhonesByGroup(batch.TargetGroup)
	if err != nil {
		s.markBatchFailed(batch.ID, err.Error())
		return err
	}

	s.db.Model(&batch).Updates(map[string]interface{}{
		"status":      "sending",
		"total_count": len(phones),
	})

	go s.executeBatch(&batch, phones, tpl.Code, nil)

	return nil
}

// ─── Internal Helpers ───

func (s *SMSService) executeBatch(batch *model.SMSBatch, phones []string, tplCode string, params map[string]string) {
	var sentCount, failedCount int

	for _, phone := range phones {
		if !s.ValidatePhone(phone) {
			failedCount++
			continue
		}

		content, err := s.ReplaceTemplateParams(tplCode, params)
		if err != nil {
			failedCount++
			continue
		}

		operator := s.DetectOperator(phone)
		logEntry := &model.SMSLog{
			Phone:      phone,
			Content:    content,
			TemplateID: tplCode,
			Params:     marshalParams(params),
			Status:     "pending",
			Operator:   operator,
			BatchID:    &batch.ID,
		}
		s.db.Create(logEntry)

		sendErr := s.smsSnd.Send(phone, content)
		now := time.Now()
		logEntry.SentAt = &now

		if sendErr != nil {
			logEntry.Status = "failed"
			logEntry.Response = sendErr.Error()
			failedCount++
		} else {
			logEntry.Status = "sent"
			logEntry.Response = "ok"
			sentCount++
		}
		s.db.Save(logEntry)

		// Update batch progress
		s.db.Model(batch).Updates(map[string]interface{}{
			"sent_count":   sentCount,
			"failed_count": failedCount,
		})
	}

	completedAt := time.Now()
	s.db.Model(batch).Updates(map[string]interface{}{
		"status":       "completed",
		"sent_count":   sentCount,
		"failed_count": failedCount,
		"completed_at": &completedAt,
	})

	s.log.Infof("batch %d completed: total=%d sent=%d failed=%d", batch.ID, len(phones), sentCount, failedCount)
}

func (s *SMSService) markBatchFailed(batchID uint, errMsg string) {
	completedAt := time.Now()
	s.db.Model(&model.SMSBatch{}).Where("id = ?", batchID).Updates(map[string]interface{}{
		"status":       "failed",
		"completed_at": &completedAt,
	})
}

func (s *SMSService) getPhonesByGroup(group string) ([]string, error) {
	var users []model.User
	query := s.db.Model(&model.User{}).Where("status = 1 AND phone != ''")

	switch group {
	case "all":
		// no filter
	case "new":
		query = query.Where("created_at >= ?", time.Now().AddDate(0, 0, -30))
	case "active":
		query = query.Where("updated_at >= ?", time.Now().AddDate(0, 0, -90))
	case "vip":
		query = query.Where("group_id > 1")
	default:
		return nil, fmt.Errorf("unknown target group: %s", group)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	phones := make([]string, 0, len(users))
	for _, u := range users {
		if u.Phone != "" {
			phones = append(phones, u.Phone)
		}
	}
	return phones, nil
}

func marshalParams(params map[string]string) datatypes.JSON {
	if params == nil {
		return nil
	}
	data, _ := json.Marshal(params)
	return datatypes.JSON(data)
}
