package service

import (
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type BatchSyncService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewBatchSyncService(db *gorm.DB, log *logger.Logger) *BatchSyncService {
	return &BatchSyncService{db: db, log: log}
}

// GetTaskList returns paginated batch sync tasks.
func (s *BatchSyncService) GetTaskList(page, pageSize int, taskType string, status *int8) ([]model.BatchSyncTask, int64, error) {
	var tasks []model.BatchSyncTask
	var total int64

	query := s.db.Model(&model.BatchSyncTask{})
	if taskType != "" {
		query = query.Where("type = ?", taskType)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&tasks).Error; err != nil {
		return nil, 0, err
	}
	return tasks, total, nil
}

// GetTaskByID returns a single task by ID.
func (s *BatchSyncService) GetTaskByID(id uint) (*model.BatchSyncTask, error) {
	var task model.BatchSyncTask
	if err := s.db.First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

// CreateTask creates a new batch sync task.
func (s *BatchSyncService) CreateTask(task *model.BatchSyncTask) error {
	return s.db.Create(task).Error
}

// Execute starts a batch sync task execution.
func (s *BatchSyncService) Execute(taskID uint, operatorID uint) error {
	now := time.Now()
	return s.db.Model(&model.BatchSyncTask{}).Where("id = ?", taskID).Updates(map[string]interface{}{
		"status":     1,
		"started_at": &now,
	}).Error
}

// Complete marks a batch sync task as completed.
func (s *BatchSyncService) Complete(taskID uint, success, failed, skipped int) error {
	now := time.Now()
	return s.db.Model(&model.BatchSyncTask{}).Where("id = ?", taskID).Updates(map[string]interface{}{
		"status":      2,
		"success":     success,
		"failed":      failed,
		"skipped":     skipped,
		"finished_at": &now,
	}).Error
}

// GetLogs returns logs for a batch sync task.
func (s *BatchSyncService) GetLogs(taskID uint, page, pageSize int) ([]model.BatchSyncLog, int64, error) {
	var logs []model.BatchSyncLog
	var total int64

	query := s.db.Model(&model.BatchSyncLog{}).Where("task_id = ?", taskID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}
