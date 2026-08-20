package job

import (
	"time"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type CronJobs struct {
	db        *gorm.DB
	log       *logger.Logger
	invSvc    *service.InvoiceService
	orderSvc  *service.OrderService
}

func NewCronJobs(db *gorm.DB, log *logger.Logger, invSvc *service.InvoiceService, orderSvc *service.OrderService) *CronJobs {
	return &CronJobs{
		db:       db,
		log:      log,
		invSvc:   invSvc,
		orderSvc: orderSvc,
	}
}

// Start registers and starts all cron jobs.
func (j *CronJobs) Start() *cron.Cron {
	c := cron.New(cron.WithSeconds())

	// Auto-renew check every hour
	c.AddFunc("0 0 * * * *", j.AutoRenew)

	// Invoice overdue check every 30 minutes
	c.AddFunc("0 */30 * * * *", j.InvoiceCheck)

	// Product status check every 6 hours
	c.AddFunc("0 0 */6 * * *", j.ProductStatusCheck)

	c.Start()
	j.log.Info("cron jobs started")
	return c
}

// AutoRenew processes user products that are set to auto-renew and are near expiry.
func (j *CronJobs) AutoRenew() {
	j.log.Info("running AutoRenew job")

	var userProducts []service.UserProduct
	threshold := time.Now().AddDate(0, 0, 3) // 3 days before expiry

	if err := j.db.Preload("Product").
		Where("auto_renew = ? AND status = 1 AND expire_at <= ? AND expire_at > ?",
			true, threshold, time.Now()).
		Find(&userProducts).Error; err != nil {
		j.log.Errorf("AutoRenew: failed to query user products: %v", err)
		return
	}

	for _, up := range userProducts {
		// Create renewal invoice
		_, err := j.invSvc.CreateRenew(up.UserID, up.Product.Price, "auto-renew for product #"+up.Product.Name)
		if err != nil {
			j.log.Errorf("AutoRenew: failed to create invoice for user %d: %v", up.UserID, err)
			continue
		}
		j.log.Infof("AutoRenew: renewal invoice created for user %d, product %s", up.UserID, up.Product.Name)
	}
}

// InvoiceCheck marks overdue invoices.
func (j *CronJobs) InvoiceCheck() {
	j.log.Info("running InvoiceCheck job")

	invoices, err := j.invSvc.GetOverdueInvoices()
	if err != nil {
		j.log.Errorf("InvoiceCheck: failed to get overdue invoices: %v", err)
		return
	}

	for _, inv := range invoices {
		if err := j.db.Model(&inv).Update("status", 3).Error; err != nil {
			j.log.Errorf("InvoiceCheck: failed to mark invoice %s overdue: %v", inv.InvoiceNo, err)
			continue
		}
		j.log.Infof("InvoiceCheck: invoice %s marked overdue", inv.InvoiceNo)
	}
}

// ProductStatusCheck deactivates expired user products.
func (j *CronJobs) ProductStatusCheck() {
	j.log.Info("running ProductStatusCheck job")

	result := j.db.Model(&service.UserProduct{}).
		Where("status = 1 AND expire_at < ?", time.Now()).
		Update("status", 2)

	if result.Error != nil {
		j.log.Errorf("ProductStatusCheck: failed to update: %v", result.Error)
		return
	}

	j.log.Infof("ProductStatusCheck: %d products marked as expired", result.RowsAffected)
}

// Package-level cron instance for Start/StopAll
var defaultCron *cron.Cron

// Start starts the default cron jobs (package-level convenience).
func Start() {
	if defaultCron == nil {
		defaultCron = cron.New()
	}
	defaultCron.Start()
}

// StopAll stops all cron jobs.
func StopAll() {
	if defaultCron != nil {
		ctx := defaultCron.Stop()
		<-ctx.Done()
	}
}
