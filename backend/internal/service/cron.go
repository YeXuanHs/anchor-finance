package service

import (
	"errors"
	"fmt"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// CronService 定时任务业务逻辑
type CronService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewCronService(db *gorm.DB, log *logger.Logger) *CronService {
	return &CronService{db: db, log: log}
}

// CreateTask 创建定时任务
func (s *CronService) CreateTask(task *model.CronTask) error {
	if err := s.db.Create(task).Error; err != nil {
		return fmt.Errorf("create cron task: %w", err)
	}
	s.log.Infof("cron task created: id=%d name=%s", task.ID, task.Name)
	return nil
}

// UpdateTask 更新定时任务
func (s *CronService) UpdateTask(id uint, updates map[string]interface{}) error {
	result := s.db.Model(&model.CronTask{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("task not found")
	}
	s.log.Infof("cron task updated: id=%d", id)
	return nil
}

// DeleteTask 删除定时任务（软删除）
func (s *CronService) DeleteTask(id uint) error {
	result := s.db.Where("id = ?", id).Delete(&model.CronTask{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("task not found")
	}
	s.log.Infof("cron task deleted: id=%d", id)
	return nil
}

// GetTaskByID 根据ID获取任务
func (s *CronService) GetTaskByID(id uint) (*model.CronTask, error) {
	var task model.CronTask
	if err := s.db.First(&task, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("task not found")
		}
		return nil, err
	}
	return &task, nil
}

// GetTaskList 获取任务列表（分页）
func (s *CronService) GetTaskList(page, pageSize int, keyword string, status *int8) ([]model.CronTask, int64, error) {
	var tasks []model.CronTask
	var total int64

	query := s.db.Model(&model.CronTask{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("priority DESC, id DESC").Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// SetStatus 启用/禁用任务
func (s *CronService) SetStatus(id uint, status int8) error {
	result := s.db.Model(&model.CronTask{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("task not found")
	}
	s.log.Infof("cron task status changed: id=%d status=%d", id, status)
	return nil
}

// RunTask 手动执行任务
func (s *CronService) RunTask(taskID uint, operatorID uint) (*model.CronTaskLog, error) {
	task, err := s.GetTaskByID(taskID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	taskLog := &model.CronTaskLog{
		TaskID:    task.ID,
		TaskName:  task.Name,
		Status:    1,
		StartedAt: now,
		Trigger:   "manual",
	}

	if err := s.db.Create(taskLog).Error; err != nil {
		return nil, fmt.Errorf("create task log: %w", err)
	}

	// 模拟执行任务并记录结果
	go s.executeTask(task, taskLog)

	s.log.Infof("cron task manually triggered: task_id=%d operator=%d", taskID, operatorID)
	return taskLog, nil
}

// executeTask 执行任务并记录结果
func (s *CronService) executeTask(task *model.CronTask, taskLog *model.CronTaskLog) {
	now := time.Now()
	output := fmt.Sprintf("Task [%s] executed at %s", task.Name, now.Format("2006-01-02 15:04:05"))

	// 根据任务类型执行不同的逻辑
	switch task.Type {
	case "custom":
		output += "\nExecuting custom command: " + task.Command
	case "system":
		output += "\nExecuting system maintenance task"
	case "plugin":
		output += "\nExecuting plugin: " + task.Command
	default:
		output += "\nUnknown task type"
	}

	finishTime := time.Now()
	duration := int(finishTime.Sub(now).Milliseconds())

	// 更新日志
	logUpdates := map[string]interface{}{
		"status":      2,
		"output":      output,
		"finished_at": &finishTime,
		"duration":    duration,
	}
	s.db.Model(taskLog).Updates(logUpdates)

	// 更新任务状态
	taskUpdates := map[string]interface{}{
		"last_run_at":  &now,
		"last_result":  "success",
		"last_error":   "",
		"run_count":    gorm.Expr("run_count + 1"),
	}
	s.db.Model(task).Where("id = ?", task.ID).Updates(taskUpdates)

	s.log.Infof("cron task executed: id=%d duration=%dms", task.ID, duration)
}

// GetTaskLogs 获取任务执行日志
func (s *CronService) GetTaskLogs(taskID uint, page, pageSize int) ([]model.CronTaskLog, int64, error) {
	var logs []model.CronTaskLog
	var total int64

	query := s.db.Model(&model.CronTaskLog{})
	if taskID > 0 {
		query = query.Where("task_id = ?", taskID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetLogByID 根据ID获取日志详情
func (s *CronService) GetLogByID(logID uint) (*model.CronTaskLog, error) {
	var taskLog model.CronTaskLog
	if err := s.db.First(&taskLog, logID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("log not found")
		}
		return nil, err
	}
	return &taskLog, nil
}

// GetEnabledTasks 获取所有已启用的任务
func (s *CronService) GetEnabledTasks() ([]model.CronTask, error) {
	var tasks []model.CronTask
	if err := s.db.Where("status = 1").Order("priority DESC").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}
