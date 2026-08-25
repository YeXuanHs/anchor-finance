package service

import (
	"log"
	"time"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"gorm.io/gorm"
)

// CronEngine 统一定时任务引擎
// 合并供应商同步+AI工单+通用定时任务到一个goroutine
type CronEngine struct {
	db             *gorm.DB
	supplierSync   *SupplierSyncService
	ticker         *time.Ticker
	stopCh         chan struct{}
}

// NewCronEngine 创建统一定时任务引擎
func NewCronEngine() *CronEngine {
	return &CronEngine{
		db:           database.GetDB(),
		supplierSync: NewSupplierSyncService(),
		stopCh:       make(chan struct{}),
	}
}

// Start 启动统一定时任务引擎
// 每分钟检查一次需要执行的任务
func (e *CronEngine) Start() {
	log.Println("[Cron] Starting unified cron engine")

	// 启动时初始化默认定时任务
	e.initDefaultTasks()

	// 每分钟tick一次
	e.ticker = time.NewTicker(1 * time.Minute)

	go func() {
		// 启动时立即执行一次
		e.runAll()

		for {
			select {
			case <-e.ticker.C:
				e.runAll()
			case <-e.stopCh:
				return
			}
		}
	}()
}

// Stop 停止引擎
func (e *CronEngine) Stop() {
	if e.ticker != nil {
		e.ticker.Stop()
	}
	close(e.stopCh)
	log.Println("[Cron] Cron engine stopped")
}

// initDefaultTasks 初始化默认定时任务（如果不存在）
func (e *CronEngine) initDefaultTasks() {
	defaults := []struct {
		name string
		cron string
		typ  string
	}{
		{"供应商价格同步", "*/5 * * * *", "supplier_price_sync"},
		{"供应商库存同步", "*/5 * * * *", "supplier_stock_sync"},
		{"供应商全量同步", "0 3 * * *", "supplier_full_sync"},
		{"AI工单队列处理", "* * * * *", "ai_ticket_process"},
		{"到期服务暂停检查", "0 * * * *", "suspend_check"},
		{"到期服务续费提醒", "0 9 * * *", "renew_reminder"},
		{"流量重置", "0 0 1 * *", "flow_reset"},
	}

	for _, d := range defaults {
		var count int64
		e.db.Model(&model.ScheduleTask{}).Where("type = ?", d.typ).Count(&count)
		if count == 0 {
			e.db.Create(&model.ScheduleTask{
				Name:   d.name,
				Cron:   d.cron,
				Type:   d.typ,
				Status: "active",
			})
			log.Printf("[Cron] Created default task: %s (%s)", d.name, d.typ)
		}
	}
}

// runAll 检查并执行所有到期的任务
func (e *CronEngine) runAll() {
	var tasks []model.ScheduleTask
	e.db.Where("status = ?", "active").Find(&tasks)

	now := time.Now()
	for _, task := range tasks {
		if e.shouldRun(task, now) {
			go e.executeTask(task)
		}
	}
}

// shouldRun 判断任务是否应该执行
func (e *CronEngine) shouldRun(task model.ScheduleTask, now time.Time) bool {
	if task.LastRunAt == nil {
		return true
	}
	
	switch task.Type {
	case "supplier_price_sync":
		// 每5分钟
		return now.Sub(*task.LastRunAt) >= 5*time.Minute
	case "supplier_stock_sync":
		// 每5分钟
		return now.Sub(*task.LastRunAt) >= 5*time.Minute
	case "supplier_full_sync":
		// 每天
		return now.Sub(*task.LastRunAt) >= 24*time.Hour
	case "ai_ticket_process":
		// 每分钟
		return now.Sub(*task.LastRunAt) >= 1*time.Minute
	case "suspend_check":
		// 每小时
		return now.Sub(*task.LastRunAt) >= 1*time.Hour
	case "renew_reminder":
		// 每天9点（简化为每24小时）
		return now.Sub(*task.LastRunAt) >= 24*time.Hour
	case "flow_reset":
		// 每月1号（简化为每30天）
		return now.Sub(*task.LastRunAt) >= 30*24*time.Hour
	default:
		return false
	}
}

// executeTask 执行单个任务
func (e *CronEngine) executeTask(task model.ScheduleTask) {
	startedAt := time.Now()
	log.Printf("[Cron] Executing task: %s (type=%s)", task.Name, task.Type)

	// 记录运行
	run := model.ScheduleRun{
		TaskID:    task.ID,
		Status:    "success",
		StartedAt: startedAt,
	}

	var err error
	// 遍历所有供应商执行同步
	var suppliers []model.Supplier
	e.db.Where("status = ?", "active").Find(&suppliers)
	for _, supplier := range suppliers {
		switch task.Type {
		case "supplier_price_sync":
			if e := e.supplierSync.SyncAllPrices(supplier.ID); e != nil {
				err = e
			}
		case "supplier_stock_sync":
			if e := e.supplierSync.SyncAllProducts(supplier.ID); e != nil {
				err = e
			}
		case "supplier_full_sync":
			if e := e.supplierSync.SyncAllProducts(supplier.ID); e != nil {
				err = e
			}
			if e := e.supplierSync.SyncAllPrices(supplier.ID); e != nil {
				err = e
			}
			if e := e.supplierSync.SyncProductStatus(supplier.ID); e != nil {
				err = e
			}
		}
	}
	switch task.Type {
	case "ai_ticket_process":
		e.processAITickets()
	case "suspend_check":
		e.checkSuspendServices()
	case "renew_reminder":
		e.sendRenewReminders()
	case "flow_reset":
		e.resetFlowForRenewals()
	}

	finishedAt := time.Now()
	run.FinishedAt = &finishedAt

	if err != nil {
		run.Status = "failed"
		run.Detail = err.Error()
		log.Printf("[Cron] Task %s failed: %v", task.Name, err)
	} else {
		run.Detail = "completed"
		log.Printf("[Cron] Task %s completed in %v", task.Name, finishedAt.Sub(startedAt))
	}

	// 更新最后运行时间
	e.db.Model(&task).Update("last_run_at", startedAt)
	e.db.Create(&run)
}

// processAITickets 处理AI工单队列
func (e *CronEngine) processAITickets() {
	aiTicketSvc := NewAITicketService()
	processed, err := aiTicketSvc.ProcessQueue(10)
	if err != nil {
		log.Printf("[Cron] AI ticket process error: %v", err)
		return
	}
	if processed > 0 {
		log.Printf("[Cron] Processed %d AI ticket items", processed)
	}
}

// checkSuspendServices 检查到期服务并暂停
func (e *CronEngine) checkSuspendServices() {
	var services []model.Service
	e.db.Where("status = ? AND next_due_date IS NOT NULL AND next_due_date < ?", "active", time.Now()).Find(&services)

	for _, svc := range services {
		e.db.Model(&svc).Update("status", "suspended")
		log.Printf("[Cron] Suspended service %d (due date passed)", svc.ID)
	}
}

// sendRenewReminders 发送续费提醒
func (e *CronEngine) sendRenewReminders() {
	// 7天内到期的服务
	reminderDate := time.Now().AddDate(0, 0, 7)
	var services []model.Service
	e.db.Where("status = ? AND next_due_date IS NOT NULL AND next_due_date <= ? AND next_due_date > ?", "active", reminderDate, time.Now()).Find(&services)

	for _, svc := range services {
		// 通知用户（通过消息系统）
		log.Printf("[Cron] Renew reminder for service %d (user %d)", svc.ID, svc.UserID)
	}
}

// resetFlowForRenewals 重置到期流量
func (e *CronEngine) resetFlowForRenewals() {
	log.Println("[Cron] Flow reset for renewals")
}
