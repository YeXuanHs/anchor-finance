package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// CronService 定时任务业务逻辑
type CronService struct {
	db       *gorm.DB
	log      *logger.Logger
	creditSvc *CreditService
	provSvc  *ProvisionService
}

func NewCronService(db *gorm.DB, log *logger.Logger, provSvc *ProvisionService) *CronService {
	return &CronService{
		db:       db,
		log:      log,
		creditSvc: NewCreditService(db, log),
		provSvc:  provSvc,
	}
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
	var output string
	var taskErr error
	status := int8(2) // success

	timeout := time.Duration(task.Timeout) * time.Second
	if timeout <= 0 {
		timeout = 60 * time.Second
	}

	switch task.Type {
	case "command":
		taskErr, output = s.execCommandTask(task, timeout)
	case "http":
		taskErr, output = s.execHTTPTask(task, timeout)
	case "script":
		taskErr, output = s.execScriptTask(task, timeout)
	case "database":
		taskErr, output = s.execDatabaseTask(task)
	case "cleanup":
		taskErr, output = s.execCleanupTask()
	case "auto_suspend":
		taskErr, output = s.execAutoSuspendTask()
	case "auto_terminate":
		taskErr, output = s.execAutoTerminateTask()
	case "invoice_reminder":
		taskErr, output = s.execInvoiceReminderTask()
	case "renewal_reminder":
		taskErr, output = s.execRenewalReminderTask()
	case "auto_renew":
		taskErr, output = s.execAutoRenewTask()
	case "credit_bill_generation":
		output, taskErr = s.creditSvc.GenerateMonthlyBills()
	case "credit_late_fee":
		output, taskErr = s.creditSvc.ApplyLateFees()
	default:
		taskErr = fmt.Errorf("unknown task type: %s", task.Type)
		output = "Unknown task type"
	}

	if taskErr != nil {
		status = 3 // failed
		s.log.Errorf("cron task failed: id=%d type=%s err=%v", task.ID, task.Type, taskErr)
	}

	finishTime := time.Now()
	duration := int(finishTime.Sub(now).Milliseconds())

	logUpdates := map[string]interface{}{
		"status":      status,
		"output":      output,
		"finished_at": &finishTime,
		"duration":    duration,
	}
	if taskErr != nil {
		logUpdates["error_msg"] = taskErr.Error()
		logUpdates["status"] = 3
	}
	s.db.Model(taskLog).Updates(logUpdates)

	lastResult := "success"
	lastError := ""
	if taskErr != nil {
		lastResult = "failed"
		lastError = taskErr.Error()
	}
	taskUpdates := map[string]interface{}{
		"last_run_at": &now,
		"last_result": lastResult,
		"last_error":  lastError,
		"run_count":   gorm.Expr("run_count + 1"),
	}
	s.db.Model(task).Where("id = ?", task.ID).Updates(taskUpdates)

	s.log.Infof("cron task executed: id=%d type=%s duration=%dms result=%s", task.ID, task.Type, duration, lastResult)
}

// execCommandTask 执行系统命令
func (s *CronService) execCommandTask(task *model.CronTask, timeout time.Duration) (error, string) {
	if task.Command == "" {
		return fmt.Errorf("command is empty"), ""
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "cmd", "/C", task.Command)
	} else {
		cmd = exec.CommandContext(ctx, "/bin/sh", "-c", task.Command)
	}

	outputBytes, err := cmd.CombinedOutput()
	output := string(outputBytes)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("command timed out after %s", timeout), output
		}
		return fmt.Errorf("command failed: %w", err), output
	}
	return nil, fmt.Sprintf("Command executed successfully\n%s", output)
}

// execHTTPTask 执行HTTP请求
func (s *CronService) execHTTPTask(task *model.CronTask, timeout time.Duration) (error, string) {
	var params struct {
		Method  string            `json:"method"`
		URL     string            `json:"url"`
		Headers map[string]string `json:"headers"`
		Body    string            `json:"body"`
	}
	if task.Params != "" {
		if err := json.Unmarshal([]byte(task.Params), &params); err != nil {
			return fmt.Errorf("invalid params JSON: %w", err), ""
		}
	}
	if params.URL == "" {
		params.URL = task.Command
	}
	if params.URL == "" {
		return fmt.Errorf("URL is empty"), ""
	}
	if params.Method == "" {
		params.Method = "GET"
	}

	var bodyReader io.Reader
	if params.Body != "" {
		bodyReader = bytes.NewBufferString(params.Body)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, strings.ToUpper(params.Method), params.URL, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err), ""
	}
	for k, v := range params.Headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("HTTP request timed out after %s", timeout), ""
		}
		return fmt.Errorf("HTTP request failed: %w", err), ""
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	output := fmt.Sprintf("HTTP %s %s\nStatus: %d\nBody: %s", params.Method, params.URL, resp.StatusCode, string(respBody))

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status), output
	}
	return nil, output
}

// execScriptTask 执行脚本文件
func (s *CronService) execScriptTask(task *model.CronTask, timeout time.Duration) (error, string) {
	scriptPath := task.Command
	if scriptPath == "" {
		var params struct {
			Path string `json:"path"`
		}
		if task.Params != "" {
			json.Unmarshal([]byte(task.Params), &params)
		}
		scriptPath = params.Path
	}
	if scriptPath == "" {
		return fmt.Errorf("script path is empty"), ""
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "cmd", "/C", scriptPath)
	} else {
		cmd = exec.CommandContext(ctx, "/bin/sh", scriptPath)
	}

	outputBytes, err := cmd.CombinedOutput()
	output := string(outputBytes)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("script timed out after %s", timeout), output
		}
		return fmt.Errorf("script failed: %w", err), output
	}
	return nil, fmt.Sprintf("Script executed successfully\n%s", output)
}

// execDatabaseTask 执行数据库维护
func (s *CronService) execDatabaseTask(task *model.CronTask) (error, string) {
	var results []string
	var hasErr bool

	if err := s.db.Exec("VACUUM").Error; err != nil {
		results = append(results, fmt.Sprintf("VACUUM failed: %v", err))
		s.log.Warnf("VACUUM failed: %v", err)
		hasErr = true
	} else {
		results = append(results, "VACUUM completed")
	}

	if err := s.db.Exec("ANALYZE").Error; err != nil {
		results = append(results, fmt.Sprintf("ANALYZE failed: %v", err))
		s.log.Warnf("ANALYZE failed: %v", err)
		hasErr = true
	} else {
		results = append(results, "ANALYZE completed")
	}

	output := strings.Join(results, "\n")
	if hasErr {
		return fmt.Errorf("database maintenance had errors"), output
	}
	return nil, output
}

// execCleanupTask 清理过期日志
func (s *CronService) execCleanupTask() (error, string) {
	var results []string
	var hasErr bool

	loginCutoff := time.Now().AddDate(0, 0, -90)
	result := s.db.Where("created_at < ?", loginCutoff).Delete(&model.LoginLog{})
	if result.Error != nil {
		results = append(results, fmt.Sprintf("login_logs cleanup failed: %v", result.Error))
		s.log.Warnf("login_logs cleanup failed: %v", result.Error)
		hasErr = true
	} else {
		results = append(results, fmt.Sprintf("Deleted %d login_logs older than 90 days", result.RowsAffected))
	}

	logCutoff := time.Now().AddDate(0, 0, -30)
	result = s.db.Where("created_at < ?", logCutoff).Delete(&model.SystemLog{})
	if result.Error != nil {
		results = append(results, fmt.Sprintf("system_logs cleanup failed: %v", result.Error))
		s.log.Warnf("system_logs cleanup failed: %v", result.Error)
		hasErr = true
	} else {
		results = append(results, fmt.Sprintf("Deleted %d system_logs older than 30 days", result.RowsAffected))
	}

	output := strings.Join(results, "\n")
	if hasErr {
		return fmt.Errorf("cleanup had errors"), output
	}
	return nil, output
}

// execAutoSuspendTask 自动挂起逾期产品
func (s *CronService) execAutoSuspendTask() (error, string) {
	now := time.Now()
	graceDays := 3

	var invoices []model.Invoice
	if err := s.db.Where("status IN (0, 5) AND due_date IS NOT NULL AND due_date < ?",
		now.AddDate(0, 0, -graceDays)).Find(&invoices).Error; err != nil {
		return fmt.Errorf("query overdue invoices: %w", err), ""
	}

	suspended := 0
	for _, inv := range invoices {
		var userProducts []model.UserProduct
		if err := s.db.Where("user_id = ? AND status = 1", inv.UserID).Find(&userProducts).Error; err != nil {
			s.log.Warnf("query user products for user %d: %v", inv.UserID, err)
			continue
		}

		for _, up := range userProducts {
			reason := fmt.Sprintf("Auto-suspended: overdue invoice %s", inv.InvoiceNo)
			if s.provSvc != nil {
				if err := s.provSvc.SuspendService(up.ID, reason); err != nil {
					s.log.Warnf("provision suspend user_product %d: %v", up.ID, err)
					// Fallback: update DB directly
					s.db.Model(&up).Updates(map[string]interface{}{
						"status":         2,
						"suspend_reason": reason,
					})
				}
			} else {
				s.db.Model(&up).Updates(map[string]interface{}{
					"status":         2,
					"suspend_reason": reason,
				})
			}
			suspended++
		}

		if inv.Status == 0 {
			s.db.Model(&inv).Update("status", 5) // mark overdue
		}
	}

	output := fmt.Sprintf("Processed %d overdue invoices, suspended %d products", len(invoices), suspended)
	return nil, output
}

// execAutoTerminateTask 自动终止长期挂起产品
func (s *CronService) execAutoTerminateTask() (error, string) {
	cutoff := time.Now().AddDate(0, 0, -30)

	var products []model.UserProduct
	if err := s.db.Where("status = 2 AND updated_at < ?", cutoff).Find(&products).Error; err != nil {
		return fmt.Errorf("query suspended products: %w", err), ""
	}

	terminated := 0
	for _, up := range products {
		reason := "Auto-terminated: suspended >30 days"
		if s.provSvc != nil {
			if err := s.provSvc.TerminateService(up.ID, reason); err != nil {
				s.log.Warnf("provision terminate user_product %d: %v", up.ID, err)
				// Fallback: update DB directly
				now := time.Now()
				s.db.Model(&up).Updates(map[string]interface{}{
					"status":           4,
					"termination_date": &now,
				})
			}
		} else {
			now := time.Now()
			s.db.Model(&up).Updates(map[string]interface{}{
				"status":           4,
				"termination_date": &now,
			})
		}
		terminated++
	}

	output := fmt.Sprintf("Checked %d suspended products, terminated %d products suspended >30 days", len(products), terminated)
	return nil, output
}

// execInvoiceReminderTask 发送账单到期提醒
func (s *CronService) execInvoiceReminderTask() (error, string) {
	now := time.Now()
	reminderDate := now.AddDate(0, 0, 3)

	var invoices []model.Invoice
	if err := s.db.Where("status IN (0, 5, 6) AND due_date IS NOT NULL AND due_date <= ? AND due_date >= ?",
		reminderDate, now).Find(&invoices).Error; err != nil {
		return fmt.Errorf("query upcoming invoices: %w", err), ""
	}

	reminded := 0
	for _, inv := range invoices {
		var existing model.SystemMessage
		if err := s.db.Where("user_id = ? AND type = 'system' AND title LIKE ? AND created_at > ?",
			inv.UserID, fmt.Sprintf("%%%s%%", inv.InvoiceNo), now.AddDate(0, 0, -1)).First(&existing).Error; err == nil {
			continue
		}

		msg := &model.SystemMessage{
			UserID:  inv.UserID,
			Title:   fmt.Sprintf("账单到期提醒: %s", inv.InvoiceNo),
			Content: fmt.Sprintf("您的账单 %s 将于 %s 到期，请及时支付。金额: %.2f %s", inv.InvoiceNo, inv.DueDate.Format("2006-01-02"), inv.Total, inv.Currency),
			Type:    "system",
		}
		if err := s.db.Create(msg).Error; err != nil {
			s.log.Warnf("create reminder for invoice %s: %v", inv.InvoiceNo, err)
			continue
		}

		if inv.Status == 0 || inv.Status == 5 {
			s.db.Model(&inv).Update("status", 6) // marked as reminded
		}
		reminded++
	}

	output := fmt.Sprintf("Checked %d upcoming invoices, sent %d reminders", len(invoices), reminded)
	return nil, output
}

// execRenewalReminderTask 发送到期续费提醒
func (s *CronService) execRenewalReminderTask() (error, string) {
	now := time.Now()
	reminderDate := now.AddDate(0, 0, 7)

	var products []model.UserProduct
	if err := s.db.Where("status = 1 AND next_due_date IS NOT NULL AND next_due_date <= ? AND next_due_date >= ?",
		reminderDate, now).Find(&products).Error; err != nil {
		return fmt.Errorf("query expiring products: %w", err), ""
	}

	reminded := 0
	for _, up := range products {
		var existing model.SystemMessage
		if err := s.db.Where("user_id = ? AND type = 'system' AND title LIKE ? AND created_at > ?",
			up.UserID, "%续费提醒%", now.AddDate(0, 0, -1)).First(&existing).Error; err == nil {
			continue
		}

		msg := &model.SystemMessage{
			UserID:  up.UserID,
			Title:   fmt.Sprintf("续费提醒: %s", up.Name),
			Content: fmt.Sprintf("您的产品 %s 将于 %s 到期，请及时续费。", up.Name, up.NextDueDate.Format("2006-01-02")),
			Type:    "system",
		}
		if err := s.db.Create(msg).Error; err != nil {
			s.log.Warnf("create renewal reminder for product %d: %v", up.ID, err)
			continue
		}
		reminded++
	}

	output := fmt.Sprintf("Checked %d expiring products, sent %d renewal reminders", len(products), reminded)
	return nil, output
}

// execAutoRenewTask 自动续费
func (s *CronService) execAutoRenewTask() (error, string) {
	now := time.Now()
	renewWindow := now.AddDate(0, 0, 7)

	var renewCycles []model.RenewCycle
	if err := s.db.Where("auto_renew = true AND status = 'active' AND next_due_date IS NOT NULL AND next_due_date <= ?",
		renewWindow).Find(&renewCycles).Error; err != nil {
		return fmt.Errorf("query auto-renew cycles: %w", err), ""
	}

	created := 0
	for _, rc := range renewCycles {
		var up model.UserProduct
		if err := s.db.First(&up, rc.UserProductID).Error; err != nil {
			s.log.Warnf("user_product %d not found for renew cycle %d: %v", rc.UserProductID, rc.ID, err)
			continue
		}
		if up.Status != 1 {
			continue
		}

		var product model.Product
		s.db.First(&product, up.ProductID)

		order := &model.Order{
			OrderNo:      fmt.Sprintf("RN%d%d", now.Unix(), rc.ID),
			UserID:       up.UserID,
			ProductID:    up.ProductID,
			UserProductID: up.ID,
			Type:         "renew",
			Amount:       rc.Amount,
			Total:        rc.Amount,
			Currency:     up.Currency,
			BillingCycle: rc.Cycle,
			Status:       0,
		}
		if err := s.db.Create(order).Error; err != nil {
			s.log.Warnf("create renewal order for product %d: %v", up.ID, err)
			continue
		}

		msg := &model.SystemMessage{
			UserID:  up.UserID,
			Title:   fmt.Sprintf("自动续费订单已创建: %s", up.Name),
			Content: fmt.Sprintf("您的产品 %s 自动续费订单 %s 已创建，金额: %.2f %s，请及时支付。", up.Name, order.OrderNo, rc.Amount, up.Currency),
			Type:    "order",
			Link:    fmt.Sprintf("/order/%d", order.ID),
		}
		s.db.Create(msg)
		created++
	}

	output := fmt.Sprintf("Checked %d auto-renew cycles, created %d renewal orders", len(renewCycles), created)
	return nil, output
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
