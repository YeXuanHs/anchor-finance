package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type CronURLService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewCronURLService(db *gorm.DB, log *logger.Logger) *CronURLService {
	return &CronURLService{db: db, log: log}
}

// GetTaskList returns paginated URL cron tasks.
func (s *CronURLService) GetTaskList(page, pageSize int, keyword string, status *int8) ([]model.CronURLTask, int64, error) {
	var tasks []model.CronURLTask
	var total int64

	query := s.db.Model(&model.CronURLTask{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR url LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
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
func (s *CronURLService) GetTaskByID(id uint) (*model.CronURLTask, error) {
	var task model.CronURLTask
	if err := s.db.First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

// CreateTask creates a new URL cron task.
func (s *CronURLService) CreateTask(task *model.CronURLTask) error {
	return s.db.Create(task).Error
}

// UpdateTask updates a URL cron task.
func (s *CronURLService) UpdateTask(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.CronURLTask{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteTask deletes a URL cron task.
func (s *CronURLService) DeleteTask(id uint) error {
	return s.db.Delete(&model.CronURLTask{}, id).Error
}

// SetStatus sets the status of a URL cron task.
func (s *CronURLService) SetStatus(id uint, status int8) error {
	return s.db.Model(&model.CronURLTask{}).Where("id = ?", id).Update("status", status).Error
}

// RunTask manually executes a URL cron task.
func (s *CronURLService) RunTask(taskID uint) (*model.CronURLTaskLog, error) {
	task, err := s.GetTaskByID(taskID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	logEntry := &model.CronURLTaskLog{
		TaskID:    taskID,
		TaskName:  task.Name,
		Status:    1,
		StartedAt: now,
	}

	if err := s.db.Create(logEntry).Error; err != nil {
		return nil, err
	}

	go s.executeURLRequest(task, logEntry)

	return logEntry, nil
}

// executeURLRequest 执行HTTP请求并记录结果
func (s *CronURLService) executeURLRequest(task *model.CronURLTask, logEntry *model.CronURLTaskLog) {
	timeout := time.Duration(task.Timeout) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	method := strings.ToUpper(task.Method)
	if method == "" {
		method = "GET"
	}

	var bodyReader io.Reader
	if task.Body != "" && (method == "POST" || method == "PUT" || method == "PATCH") {
		bodyReader = bytes.NewBufferString(task.Body)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, task.URL, bodyReader)
	if err != nil {
		s.finishURLLog(logEntry, 2, 0, "", err)
		s.updateURLTaskStatus(task, "failed", err.Error())
		return
	}

	if task.Headers != "" {
		var headers map[string]string
		if err := json.Unmarshal([]byte(task.Headers), &headers); err == nil {
			for k, v := range headers {
				req.Header.Set(k, v)
			}
		}
	}

	if method == "POST" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			s.finishURLLog(logEntry, 3, 0, "", fmt.Errorf("request timed out after %s", timeout))
			s.updateURLTaskStatus(task, "timeout", fmt.Sprintf("timed out after %s", timeout))
			return
		}
		s.finishURLLog(logEntry, 2, 0, "", err)
		s.updateURLTaskStatus(task, "failed", err.Error())
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 8192))

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		s.finishURLLog(logEntry, 1, resp.StatusCode, string(respBody), nil)
		s.updateURLTaskStatus(task, "success", "")
	} else {
		s.finishURLLog(logEntry, 2, resp.StatusCode, string(respBody), fmt.Errorf("HTTP %d", resp.StatusCode))
		s.updateURLTaskStatus(task, "failed", fmt.Sprintf("HTTP %d", resp.StatusCode))
	}
}

// finishURLLog 更新URL任务日志
func (s *CronURLService) finishURLLog(logEntry *model.CronURLTaskLog, status int8, statusCode int, response string, err error) {
	now := time.Now()
	duration := int(now.Sub(logEntry.StartedAt).Milliseconds())
	updates := map[string]interface{}{
		"status":      status,
		"status_code": statusCode,
		"response":    response,
		"duration":    duration,
		"finished_at": &now,
	}
	if err != nil {
		updates["error_msg"] = err.Error()
	}
	s.db.Model(logEntry).Updates(updates)
}

// updateURLTaskStatus 更新URL任务状态
func (s *CronURLService) updateURLTaskStatus(task *model.CronURLTask, result string, errMsg string) {
	now := time.Now()
	updates := map[string]interface{}{
		"last_run_at": &now,
		"run_count":   gorm.Expr("run_count + 1"),
		"last_result": result,
		"last_error":  errMsg,
	}
	if result == "failed" || result == "timeout" {
		updates["fail_count"] = gorm.Expr("fail_count + 1")
	}
	s.db.Model(task).Updates(updates)
}

// GetLogs returns paginated logs for URL cron tasks.
func (s *CronURLService) GetLogs(taskID uint, page, pageSize int) ([]model.CronURLTaskLog, int64, error) {
	var logs []model.CronURLTaskLog
	var total int64

	query := s.db.Model(&model.CronURLTaskLog{})
	if taskID > 0 {
		query = query.Where("task_id = ?", taskID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}
