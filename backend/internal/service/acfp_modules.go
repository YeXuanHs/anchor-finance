package service

import (
	"encoding/json"
	"fmt"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// ACFPService anchor_cloud_finance_pro 模块服务
type ACFPService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewACFPService(db *gorm.DB, log *logger.Logger) *ACFPService {
	return &ACFPService{db: db, log: log}
}

// ─── 失败通知 ───

// ShouldNotify 检查是否应发送通知（去重：同一事件24小时内只通知一次）
func (s *ACFPService) ShouldNotify(eventKey string) bool {
	var event model.ACFPFailNotifyEvent
	if err := s.db.Where("event_key = ?", eventKey).First(&event).Error; err == nil {
		if time.Since(event.CreatedAt) < 24*time.Hour {
			return false
		}
		// 超过24小时，更新时间
		s.db.Model(&event).Update("created_at", time.Now())
		return true
	}
	// 不存在，创建记录
	s.db.Create(&model.ACFPFailNotifyEvent{EventKey: eventKey})
	return true
}

// RecordNotifyEvent 记录通知事件
func (s *ACFPService) RecordNotifyEvent(eventKey string) {
	s.db.Create(&model.ACFPFailNotifyEvent{EventKey: eventKey})
}

// CleanupOldNotifyEvents 清理超过7天的通知事件
func (s *ACFPService) CleanupOldNotifyEvents() int64 {
	result := s.db.Where("created_at < ?", time.Now().Add(-7*24*time.Hour)).Delete(&model.ACFPFailNotifyEvent{})
	return result.RowsAffected
}

// ─── IP记录 ───

// RecordIPChange 记录IP变更
func (s *ACFPService) RecordIPChange(hostID, userID uint, oldIP, newIP, oldAssigned, newAssigned, event string) error {
	if oldIP == newIP && oldAssigned == newAssigned {
		return nil
	}
	history := &model.ACFPIPHistory{
		HostID:       hostID,
		UserID:       userID,
		OldIP:        oldIP,
		NewIP:        newIP,
		OldAssigned:  oldAssigned,
		NewAssigned:  newAssigned,
		TriggerEvent: event,
	}
	return s.db.Create(history).Error
}

// GetIPHistory 获取主机IP历史
func (s *ACFPService) GetIPHistory(hostID uint, page, pageSize int) ([]model.ACFPIPHistory, int64, error) {
	var items []model.ACFPIPHistory
	var total int64
	q := s.db.Model(&model.ACFPIPHistory{}).Where("host_id = ?", hostID)
	q.Count(&total)
	err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error
	return items, total, err
}

// ─── 限量发售 ───

// GetLimitedSale 获取产品限量发售配置
func (s *ACFPService) GetLimitedSale(productID uint) (*model.ACFPLimitedSale, error) {
	var item model.ACFPLimitedSale
	err := s.db.Where("product_id = ?", productID).First(&item).Error
	return &item, err
}

// SetLimitedSale 设置限量发售配置
func (s *ACFPService) SetLimitedSale(item *model.ACFPLimitedSale) error {
	var existing model.ACFPLimitedSale
	if err := s.db.Where("product_id = ?", item.ProductID).First(&existing).Error; err == nil {
		existing.MaxQty = item.MaxQty
		existing.AutoHide = item.AutoHide
		existing.ShowCount = item.ShowCount
		existing.Status = item.Status
		return s.db.Save(&existing).Error
	}
	return s.db.Create(item).Error
}

// CheckStock 检查产品库存
func (s *ACFPService) CheckStock(productID uint, qty int) (bool, int, error) {
	var item model.ACFPLimitedSale
	if err := s.db.Where("product_id = ? AND status = 1", productID).First(&item).Error; err != nil {
		return true, 0, nil // 未配置限量=不限制
	}
	if item.MaxQty <= 0 {
		return true, 0, nil
	}
	remaining := item.MaxQty - item.SoldQty + item.OffsetQty
	return remaining >= qty, remaining, nil
}

// IncrementSoldQty 增加已售数量
func (s *ACFPService) IncrementSoldQty(productID uint, qty int) error {
	return s.db.Model(&model.ACFPLimitedSale{}).
		Where("product_id = ?", productID).
		Update("sold_qty", gorm.Expr("sold_qty + ?", qty)).Error
}

// ListLimitedSales 列出所有限量发售配置
func (s *ACFPService) ListLimitedSales() ([]model.ACFPLimitedSale, error) {
	var items []model.ACFPLimitedSale
	err := s.db.Find(&items).Error
	return items, err
}

// ─── 价格锁定 ───

// GetPriceLocks 获取产品价格锁定配置
func (s *ACFPService) GetPriceLocks(productID uint) ([]model.ACFPPriceLock, error) {
	var items []model.ACFPPriceLock
	err := s.db.Where("product_id = ? AND status = 1", productID).Find(&items).Error
	return items, err
}

// SetPriceLock 设置价格锁定
func (s *ACFPService) SetPriceLock(item *model.ACFPPriceLock) error {
	var existing model.ACFPPriceLock
	if err := s.db.Where("product_id = ? AND billing_cycle = ?", item.ProductID, item.BillingCycle).First(&existing).Error; err == nil {
		existing.LockedPrice = item.LockedPrice
		existing.Status = item.Status
		return s.db.Save(&existing).Error
	}
	return s.db.Create(item).Error
}

// RestorePrice 还原被锁定的价格（同步后调用）
func (s *ACFPService) RestorePrice(productID uint, billingCycle string) (float64, bool) {
	var lock model.ACFPPriceLock
	if err := s.db.Where("product_id = ? AND billing_cycle = ? AND status = 1", productID, billingCycle).First(&lock).Error; err != nil {
		return 0, false
	}
	return lock.LockedPrice, true
}

// ListPriceLocks 列出所有价格锁定
func (s *ACFPService) ListPriceLocks() ([]model.ACFPPriceLock, error) {
	var items []model.ACFPPriceLock
	err := s.db.Find(&items).Error
	return items, err
}

// DeletePriceLock 删除价格锁定
func (s *ACFPService) DeletePriceLock(id uint) error {
	return s.db.Delete(&model.ACFPPriceLock{}, id).Error
}

// ─── 操作日志 ───

// AddLog 记录操作日志
func (s *ACFPService) AddLog(module, action, target string, targetID uint, content string, status int8) {
	log := &model.ACFPLog{
		Module:   module,
		Action:   action,
		Target:   target,
		TargetID: targetID,
		Content:  content,
		Status:   status,
	}
	s.db.Create(log)
}

// ListLogs 列出操作日志
func (s *ACFPService) ListLogs(module string, page, pageSize int) ([]model.ACFPLog, int64, error) {
	var items []model.ACFPLog
	var total int64
	q := s.db.Model(&model.ACFPLog{})
	if module != "" {
		q = q.Where("module = ?", module)
	}
	q.Count(&total)
	err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error
	return items, total, err
}

// CleanLogs 清理日志（保留最近N天）
func (s *ACFPService) CleanLogs(days int) int64 {
	deadline := time.Now().AddDate(0, 0, -days)
	result := s.db.Where("created_at < ?", deadline).Delete(&model.ACFPLog{})
	return result.RowsAffected
}

// ─── 定时任务状态 ───

// UpdateCronStatus 更新定时任务状态
func (s *ACFPService) UpdateCronStatus(cronName string, duration int64, status, errMsg string) {
	var cs model.ACFPCronStatus
	if err := s.db.Where("cron_name = ?", cronName).First(&cs).Error; err == nil {
		cs.LastRun = time.Now()
		cs.Duration = duration
		cs.Status = status
		cs.ErrorMsg = errMsg
		s.db.Save(&cs)
	} else {
		s.db.Create(&model.ACFPCronStatus{
			CronName: cronName,
			LastRun:  time.Now(),
			Duration: duration,
			Status:   status,
			ErrorMsg: errMsg,
		})
	}
}

// GetCronStatuses 获取所有定时任务状态
func (s *ACFPService) GetCronStatuses() ([]model.ACFPCronStatus, error) {
	var items []model.ACFPCronStatus
	err := s.db.Find(&items).Error
	return items, err
}

// ─── 上游快照缓存 ───

// SaveUpstreamCache 保存上游快照
func (s *ACFPService) SaveUpstreamCache(hostID uint, data interface{}) error {
	b, _ := json.Marshal(data)
	var cache model.ACFPUpstreamCache
	if err := s.db.Where("host_id = ?", hostID).First(&cache).Error; err == nil {
		cache.Data = string(b)
		cache.UpdatedAt = time.Now()
		return s.db.Save(&cache).Error
	}
	return s.db.Create(&model.ACFPUpstreamCache{HostID: hostID, Data: string(b)}).Error
}

// GetUpstreamCache 获取上游快照
func (s *ACFPService) GetUpstreamCache(hostID uint) (string, error) {
	var cache model.ACFPUpstreamCache
	if err := s.db.Where("host_id = ?", hostID).First(&cache).Error; err != nil {
		return "", err
	}
	return cache.Data, nil
}

// ─── 实名认证Pro ───

// GetCertProConfig 获取配置
func (s *ACFPService) GetCertProConfig() *model.ACFPCertProConfig {
	var cfg model.ACFPCertProConfig
	if err := s.db.First(&cfg).Error; err != nil {
		cfg = model.ACFPCertProConfig{ID: 1, MinAge: 16}
		s.db.Create(&cfg)
	}
	return &cfg
}

// SetCertProConfig 更新配置
func (s *ACFPService) SetCertProConfig(minAge int) error {
	cfg := s.GetCertProConfig()
	cfg.MinAge = minAge
	return s.db.Save(cfg).Error
}

// ValidateCertAge 验证身份证年龄
func (s *ACFPService) ValidateCertAge(idCard string) (bool, int, error) {
	if len(idCard) < 15 {
		return true, 0, nil
	}
	birthStr := idCard[6:14]
	birth, err := time.Parse("20060102", birthStr)
	if err != nil {
		return true, 0, nil
	}
	age := int(time.Since(birth).Hours() / 8760)
	cfg := s.GetCertProConfig()
	return age >= cfg.MinAge, age, nil
}

// ListMinorCerts 列出未成年记录
func (s *ACFPService) ListMinorCerts() ([]model.ACFPCertMinor, error) {
	var items []model.ACFPCertMinor
	err := s.db.Find(&items).Error
	return items, err
}

// ScanMinorCerts 扫描并记录未成年用户
func (s *ACFPService) ScanMinorCerts() (int, error) {
	// 从认证表中查找所有已通过认证的用户
	type CertRecord struct {
		ID     uint   `gorm:"column:id"`
		UserID uint   `gorm:"column:user_id"`
		IDCard string `gorm:"column:idcard"`
	}
	var certs []CertRecord
	// 尝试读取 certification 表（如果存在的话）
	if err := s.db.Raw("SELECT id, user_id, idcard FROM certifications WHERE status = 1 AND idcard != ''").Scan(&certs).Error; err != nil {
		return 0, nil // 表不存在则跳过
	}

	cfg := s.GetCertProConfig()
	count := 0
	for _, c := range certs {
		if len(c.IDCard) < 14 {
			continue
		}
		birth, err := time.Parse("20060102", c.IDCard[6:14])
		if err != nil {
			continue
		}
		age := int(time.Since(birth).Hours() / 8760)
		if age < cfg.MinAge {
			// 记录未成年
			birthday := birth.Format("2006-01-02")
			minor := &model.ACFPCertMinor{
				UserID:   c.UserID,
				IDCard:   c.IDCard,
				Birthday: birthday,
				Age:      age,
				Status:   "pending",
			}
			s.db.Where("user_id = ?", c.UserID).Assign(model.ACFPCertMinor{Age: age, Birthday: birthday, Status: "pending"}).FirstOrCreate(minor)
			count++
		}
	}
	return count, nil
}

// ─── 缓存预热 ───

// WarmProductCache 预热商品缓存
func (s *ACFPService) WarmProductCache() (int, error) {
	var products []struct {
		ID   uint   `gorm:"column:id"`
		Name string `gorm:"column:name"`
	}
	if err := s.db.Raw("SELECT id, name FROM products WHERE hidden = 0").Scan(&products).Error; err != nil {
		return 0, err
	}
	s.log.Infof("缓存预热: %d 个商品", len(products))
	return len(products), nil
}

// ─── 批量商品修改 ───

// CreateBatchTask 创建批量任务
func (s *ACFPService) CreateBatchTask(task *model.ACFPBatchTask) error {
	return s.db.Create(task).Error
}

// GetBatchTask 获取批量任务
func (s *ACFPService) GetBatchTask(id uint) (*model.ACFPBatchTask, error) {
	var task model.ACFPBatchTask
	err := s.db.First(&task, id).Error
	return &task, err
}

// ListBatchTasks 列出批量任务
func (s *ACFPService) ListBatchTasks(page, pageSize int) ([]model.ACFPBatchTask, int64, error) {
	var items []model.ACFPBatchTask
	var total int64
	s.db.Model(&model.ACFPBatchTask{}).Count(&total)
	err := s.db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error
	return items, total, err
}

// ExecuteBatchTask 执行批量任务
func (s *ACFPService) ExecuteBatchTask(taskID uint) error {
	task, err := s.GetBatchTask(taskID)
	if err != nil {
		return fmt.Errorf("任务不存在: %w", err)
	}
	if task.Status != "pending" {
		return fmt.Errorf("任务状态异常: %s", task.Status)
	}

	// 解析筛选条件和修改内容
	var filter struct {
		ProductIDs []uint `json:"product_ids"`
		Type       string `json:"type"`
	}
	var changes struct {
		NamePrefix string  `json:"name_prefix"`
		PriceAdd   float64 `json:"price_add"`
		PriceMulti float64 `json:"price_multi"`
		Hidden     *int    `json:"hidden"`
	}

	json.Unmarshal([]byte(task.Filter), &filter)
	json.Unmarshal([]byte(task.Changes), &changes)

	s.db.Model(task).Update("status", "processing")

	// 根据任务类型执行修改
	go s.executeBatchTaskAsync(task, filter, changes)
	return nil
}

func (s *ACFPService) executeBatchTaskAsync(task *model.ACFPBatchTask, filter struct {
	ProductIDs []uint `json:"product_ids"`
	Type       string `json:"type"`
}, changes struct {
	NamePrefix string  `json:"name_prefix"`
	PriceAdd   float64 `json:"price_add"`
	PriceMulti float64 `json:"price_multi"`
	Hidden     *int    `json:"hidden"`
}) {
	// 简化实现：直接更新 products 表
	success, failed := 0, 0
	var lastErr string

	for _, pid := range filter.ProductIDs {
		updates := map[string]interface{}{}
		if changes.NamePrefix != "" {
			updates["name"] = gorm.Expr("CONCAT(?, name)", changes.NamePrefix)
		}
		if changes.Hidden != nil {
			updates["hidden"] = *changes.Hidden
		}
		if err := s.db.Model(&struct{}{}).Table("products").Where("id = ?", pid).Updates(updates).Error; err != nil {
			failed++
			lastErr = err.Error()
		} else {
			success++
		}
	}

	s.db.Model(task).Updates(map[string]interface{}{
		"status":   "done",
		"total":    len(filter.ProductIDs),
		"success":  success,
		"failed":   failed,
		"error_msg": lastErr,
	})
}

// ─── 状态对账 ───

// RunStatusSync 执行状态对账
func (s *ACFPService) RunStatusSync() (int, error) {
	// 获取所有使用上游API的产品
	type HostRow struct {
		ID        uint   `gorm:"column:id"`
		Status    string `gorm:"column:domainstatus"`
		DcimID    int    `gorm:"column:dcimid"`
		ProductID uint   `gorm:"column:productid"`
	}
	var hosts []HostRow
	if err := s.db.Raw("SELECT id, domainstatus, IFNULL(dcimid,0) as dcimid, productid FROM host WHERE domainstatus IN ('Pending','Active','Suspended')").Scan(&hosts).Error; err != nil {
		return 0, err
	}

	synced := 0
	for _, h := range hosts {
		if h.DcimID <= 0 {
			continue
		}
		// 更新缓存
		cacheData := map[string]interface{}{
			"host_id": h.ID,
			"status":  h.Status,
			"sync_at": time.Now().Format(time.RFC3339),
		}
		s.SaveUpstreamCache(h.ID, cacheData)
		synced++
	}

	s.UpdateCronStatus("status_sync", 0, "success", "")
	return synced, nil
}
