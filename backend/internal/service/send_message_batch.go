package service

import (
	"fmt"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/email"
	"anchorfinance/pkg/sms"

	"gorm.io/gorm"
)

type SendMessageBatchService struct {
	db        *gorm.DB
	log       *logger.Logger
	emailSnd  *email.Sender
	smsSnd    *sms.Sender
}

func NewSendMessageBatchService(db *gorm.DB, log *logger.Logger) *SendMessageBatchService {
	return &SendMessageBatchService{
		db:       db,
		log:      log,
		emailSnd: email.NewSender(db),
		smsSnd:   sms.NewSender(db),
	}
}

// GetSearchParams returns available search parameters for batch messaging.
func (s *SendMessageBatchService) GetSearchParams() (map[string]interface{}, error) {
	return map[string]interface{}{
		"send_methods":  []string{"email", "mobile", "system"},
		"client_status": []string{"active", "suspended", "terminated"},
	}, nil
}

// SendBatch creates a batch send operation and dispatches messages to target users.
func (s *SendMessageBatchService) SendBatch(batch *model.SendMessageBatch) error {
	batch.BatchNo = fmt.Sprintf("BATCH_%d", time.Now().UnixNano())
	batch.Status = 1 // sending
	now := time.Now()
	batch.StartedAt = &now

	if err := s.db.Create(batch).Error; err != nil {
		return err
	}

	// Run sending in background
	go s.executeBatch(batch)

	return nil
}

// executeBatch queries target users and sends messages to each.
func (s *SendMessageBatchService) executeBatch(batch *model.SendMessageBatch) {
	// Query target users based on batch filter criteria
	var users []model.User
	query := s.db.Model(&model.User{}).Where("status = 1")
	if err := query.Find(&users).Error; err != nil {
		s.log.Errorf("batch %s: failed to query users: %v", batch.BatchNo, err)
		s.markBatchFailed(batch.ID, fmt.Sprintf("query users: %v", err))
		return
	}

	total := len(users)
	s.db.Model(batch).Update("total", total)

	var successCount, failCount int
	for _, user := range users {
		record := model.SendMessageBatchRecord{
			BatchID: batch.ID,
			UserID:  user.ID,
			Status:  3, // pending
		}

		var err error
		switch batch.SendMethod {
		case "email":
			record.Target = user.Email
			err = s.sendEmailToUser(user.Email, batch.Subject, batch.Content)
		case "mobile":
			record.Target = user.Phone
			err = s.sendSMSToUser(user.Phone, batch.Content)
		case "system":
			record.Target = fmt.Sprintf("user:%d", user.ID)
			err = nil // site messages are stored directly in DB
		default:
			err = fmt.Errorf("unsupported send method: %s", batch.SendMethod)
		}

		if err != nil {
			record.Status = 2 // failed
			record.ErrorMsg = err.Error()
			failCount++
		} else {
			now := time.Now()
			record.Status = 1 // success
			record.SentAt = &now
			successCount++
		}

		s.db.Create(&record)

		// Update batch progress periodically
		s.db.Model(batch).Updates(map[string]interface{}{
			"success": successCount,
			"failed":  failCount,
		})
	}

	// Mark batch as completed
	finishTime := time.Now()
	s.db.Model(batch).Updates(map[string]interface{}{
		"status":      2, // completed
		"success":     successCount,
		"failed":      failCount,
		"finished_at": &finishTime,
	})

	s.log.Infof("batch %s completed: total=%d success=%d failed=%d",
		batch.BatchNo, total, successCount, failCount)
}

func (s *SendMessageBatchService) sendEmailToUser(to, subject, content string) error {
	if to == "" {
		return fmt.Errorf("email address is empty")
	}
	if s.emailSnd == nil {
		return fmt.Errorf("email sender not configured")
	}
	return s.emailSnd.Send(to, subject, content)
}

func (s *SendMessageBatchService) sendSMSToUser(phone, content string) error {
	if phone == "" {
		return fmt.Errorf("phone number is empty")
	}
	if s.smsSnd == nil {
		return fmt.Errorf("sms sender not configured")
	}
	return s.smsSnd.Send(phone, content)
}

func (s *SendMessageBatchService) markBatchFailed(batchID uint, errMsg string) {
	finishTime := time.Now()
	s.db.Model(&model.SendMessageBatch{}).Where("id = ?", batchID).Updates(map[string]interface{}{
		"status":      3, // failed
		"error_msg":   errMsg,
		"finished_at": &finishTime,
	})
}

// GetBatchByID returns a single batch by ID.
func (s *SendMessageBatchService) GetBatchByID(id uint) (*model.SendMessageBatch, error) {
	var batch model.SendMessageBatch
	if err := s.db.First(&batch, id).Error; err != nil {
		return nil, err
	}
	return &batch, nil
}

// GetProgress returns the progress of a batch send operation.
func (s *SendMessageBatchService) GetProgress(batchID uint) (*model.SendMessageBatch, error) {
	return s.GetBatchByID(batchID)
}

// GetRecords returns paginated batch send records.
func (s *SendMessageBatchService) GetRecords(page, pageSize int, batchID *uint) ([]model.SendMessageBatchRecord, int64, error) {
	var records []model.SendMessageBatchRecord
	var total int64

	query := s.db.Model(&model.SendMessageBatchRecord{})
	if batchID != nil {
		query = query.Where("batch_id = ?", *batchID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&records).Error; err != nil {
		return nil, 0, err
	}
	return records, total, nil
}

// GetBatches returns paginated batch send tasks.
func (s *SendMessageBatchService) GetBatches(page, pageSize int, sendMethod string, status *int8) ([]model.SendMessageBatch, int64, error) {
	var batches []model.SendMessageBatch
	var total int64

	query := s.db.Model(&model.SendMessageBatch{})
	if sendMethod != "" {
		query = query.Where("send_method = ?", sendMethod)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&batches).Error; err != nil {
		return nil, 0, err
	}
	return batches, total, nil
}
