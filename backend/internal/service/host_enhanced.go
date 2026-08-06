package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/internal/util"
	"anchorfinance/pkg/logger"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ───────────────────────── New model types ─────────────────────────

// HostCategory 主机分类
type HostCategory struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	UserID    uint   `gorm:"index" json:"user_id"`
	Name      string `gorm:"size:128;not null" json:"name"`
	SortOrder int    `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}

func (HostCategory) TableName() string { return "host_categories" }

// HostCategoryAssignment 分类-主机关联
type HostCategoryAssignment struct {
	ID         uint `gorm:"primaryKey" json:"id"`
	HostID     uint `gorm:"index;not null" json:"host_id"`
	CategoryID uint `gorm:"index;not null" json:"category_id"`
}

func (HostCategoryAssignment) TableName() string { return "host_category_assignments" }

// HostTransferLog 主机转让记录
type HostTransferLog struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	HostID     uint       `gorm:"index;not null" json:"host_id"`
	FromUserID uint       `gorm:"not null" json:"from_user_id"`
	ToUserID   uint       `gorm:"not null" json:"to_user_id"`
	Reason     string     `gorm:"type:text" json:"reason"`
	CreatedAt  time.Time  `json:"created_at"`
}

func (HostTransferLog) TableName() string { return "host_transfer_logs" }

// HostSSLConfig 主机SSL配置
type HostSSLConfig struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	HostID    uint       `gorm:"uniqueIndex;not null" json:"host_id"`
	Type      string     `gorm:"size:32" json:"type"` // letsencrypt/custom/none
	Domain    string     `gorm:"size:256" json:"domain"`
	Cert      string     `gorm:"type:text" json:"-"`
	Key       string     `gorm:"type:text" json:"-"`
	CA        string     `gorm:"type:text" json:"-"`
	Status    string     `gorm:"size:32;default:none" json:"status"` // none/pending/active/expired/error
	IssuedAt  *time.Time `json:"issued_at"`
	ExpiresAt *time.Time `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (HostSSLConfig) TableName() string { return "host_ssl_configs" }

// HostCancelRequest 主机取消/退款申请
type HostCancelRequest struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	HostID    uint       `gorm:"index;not null" json:"host_id"`
	UserID    uint       `gorm:"index;not null" json:"user_id"`
	Reason    string     `gorm:"type:text" json:"reason"`
	Type      string     `gorm:"size:32;default:immediate" json:"type"` // immediate/end_of_period
	Status    string     `gorm:"size:20;default:pending;index" json:"status"` // pending/approved/rejected/processed
	AdminNote string     `gorm:"type:text" json:"admin_note"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (HostCancelRequest) TableName() string { return "host_cancel_requests" }

// FlowPacket 流量包
type FlowPacket struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	HostID    uint    `gorm:"-" json:"host_id"` // not persisted in table, used for response
	Name      string  `gorm:"size:128;not null" json:"name"`
	AmountGB  int     `gorm:"not null" json:"amount_gb"`
	Price     float64 `gorm:"type:decimal(12,2);not null" json:"price"`
	Currency  string  `gorm:"size:8;default:CNY" json:"currency"`
	SortOrder int     `gorm:"default:0" json:"sort_order"`
	Status    int     `gorm:"default:1" json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (FlowPacket) TableName() string { return "flow_packets" }

// HostFlowPacketPurchase 流量包购买记录
type HostFlowPacketPurchase struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	HostID       uint      `gorm:"index;not null" json:"host_id"`
	FlowPacketID uint      `gorm:"index;not null" json:"flow_packet_id"`
	UserID       uint      `gorm:"index;not null" json:"user_id"`
	AmountGB     int       `gorm:"not null" json:"amount_gb"`
	Price        float64   `gorm:"type:decimal(12,2);not null" json:"price"`
	OrderID      uint      `gorm:"index" json:"order_id"`
	Status       string    `gorm:"size:20;default:pending" json:"status"` // pending/paid/failed
	CreatedAt    time.Time `json:"created_at"`
}

func (HostFlowPacketPurchase) TableName() string { return "host_flow_packet_purchases" }

// HostSecondVerify 主机二次验证配置
type HostSecondVerify struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	HostID    uint       `gorm:"uniqueIndex;not null" json:"host_id"`
	Enabled   bool       `gorm:"default:false" json:"enabled"`
	Secret    string     `gorm:"size:128" json:"-"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (HostSecondVerify) TableName() string { return "host_second_verifies" }

// ───────────────────────── Response types ─────────────────────────

// RenewalPageData 续费页面数据
type RenewalPageData struct {
	Host           *Host    `json:"host"`
	AvailableCycles []string `json:"available_cycles"`
	CurrentCycle    string   `json:"current_cycle"`
	ExpiredAt       *time.Time `json:"expired_at"`
	AutoRenew       bool     `json:"auto_renew"`
}

// RenewalPriceData 续费价格
type RenewalPriceData struct {
	HostID       uint    `json:"host_id"`
	Cycle        string  `json:"cycle"`
	Price        float64 `json:"price"`
	Currency     string  `json:"currency"`
	NewExpiredAt *time.Time `json:"new_expired_at"`
}

// TrafficUsageData 流量使用数据
type TrafficUsageData struct {
	HostID      uint    `json:"host_id"`
	UsedGB      float64 `json:"used_gb"`
	TotalGB     float64 `json:"total_gb"`
	RemainingGB float64 `json:"remaining_gb"`
	UsagePct    float64 `json:"usage_pct"`
	ResetDate   *time.Time `json:"reset_date"`
}

// TrafficChartData 流量图表数据点
type TrafficChartData struct {
	Date    string  `json:"date"`
	InGB    float64 `json:"in_gb"`
	OutGB   float64 `json:"out_gb"`
	TotalGB float64 `json:"total_gb"`
}

// DedicatedServerData 独立服务器信息
type DedicatedServerData struct {
	HostID       uint   `json:"host_id"`
	Rack         string `json:"rack"`
	RackPosition string `json:"rack_position"`
	BandwidthMbps int  `json:"bandwidth_mbps"`
	PortSpeed    string `json:"port_speed"`
	PDU          string `json:"pdu"`
	Switch       string `json:"switch"`
	RemoteHand   bool   `json:"remote_hand"`
	IPMI         string `json:"ipmi"`
}

// HostDetailData 主机详情
type HostDetailData struct {
	Host         *Host                    `json:"host"`
	Product      *Product                 `json:"product,omitempty"`
	Category     *HostCategory            `json:"category,omitempty"`
	SSLConfig    *HostSSLConfig           `json:"ssl_config,omitempty"`
	SecondVerify *HostSecondVerify        `json:"second_verify,omitempty"`
	Traffic      *TrafficUsageData        `json:"traffic,omitempty"`
	Remark       string                   `json:"remark"`
	Tags         datatypes.JSON           `json:"tags"`
	CustomFields map[string]interface{}   `json:"custom_fields,omitempty"`
}

// HostStatusData 主机实时状态
type HostStatusData struct {
	HostID      uint   `json:"host_id"`
	Status      int    `json:"status"`
	PowerStatus int    `json:"power_status"`
	Uptime      string `json:"uptime"`
	Hostname    string `json:"hostname"`
	IP          string `json:"ip"`
	OS          string `json:"os"`
}

// CancelPageData 取消页面数据
type CancelPageData struct {
	HostID        uint   `json:"host_id"`
	Hostname      string `json:"hostname"`
	ExpiredAt     *time.Time `json:"expired_at"`
	CancelTypes   []string `json:"cancel_types"` // immediate/end_of_period
	HasPending    bool   `json:"has_pending"`
	PendingRequest *HostCancelRequest `json:"pending_request,omitempty"`
}

// ReinstallCheckData 重装检查结果
type ReinstallCheckData struct {
	Allowed      bool   `json:"allowed"`
	Reason       string `json:"reason,omitempty"`
	HostStatus   int    `json:"host_status"`
	AvailableOS  []string `json:"available_os,omitempty"`
}

// ReinstallStatusData 重装进度
type ReinstallStatusData struct {
	InProgress bool    `json:"in_progress"`
	Progress   int     `json:"progress"` // 0-100
	StartedAt  *time.Time `json:"started_at"`
	ETA        *time.Time `json:"eta"`
}

// HostRechargeData 主机充值选项
type HostRechargeData struct {
	HostID        uint    `json:"host_id"`
	Hostname      string  `json:"hostname"`
	Balance       float64 `json:"balance"`
	DueAmount     float64 `json:"due_amount"`
	DueDate       *time.Time `json:"due_date"`
	MinRecharge   float64 `json:"min_recharge"`
}

// ───────────────────────── Service ─────────────────────────

// HostEnhancedService extends HostService with additional management capabilities.
type HostEnhancedService struct {
	db          *gorm.DB
	log         *logger.Logger
	hostSvc     *HostService
	invSvc      *InvoiceService
	balSvc      *BalanceService
	upstreamSvc *UpstreamService
}

func NewHostEnhancedService(db *gorm.DB, log *logger.Logger, hostSvc *HostService, invSvc *InvoiceService, balSvc *BalanceService, upstreamSvc *UpstreamService) *HostEnhancedService {
	return &HostEnhancedService{db: db, log: log, hostSvc: hostSvc, invSvc: invSvc, balSvc: balSvc, upstreamSvc: upstreamSvc}
}

// ═══════════════════ Renewal ═══════════════════

// GetRenewalPage returns renewal options with pricing for a host.
func (s *HostEnhancedService) GetRenewalPage(hostID uint) (*RenewalPageData, error) {
	host, err := s.hostSvc.GetByID(hostID)
	if err != nil {
		return nil, fmt.Errorf("host not found: %w", err)
	}

	cycles := []string{"monthly", "quarterly", "semi-annually", "annually", "biennially", "triennially"}

	autoRenew := false
	var rc model.RenewCycle
	if err := s.db.Where("user_product_id = ?", hostID).First(&rc).Error; err == nil {
		autoRenew = rc.AutoRenew
	}

	return &RenewalPageData{
		Host:            host,
		AvailableCycles: cycles,
		ExpiredAt:       host.ExpiredAt,
		AutoRenew:       autoRenew,
	}, nil
}

// GetRenewalPrice calculates renewal price for a host and billing cycle.
func (s *HostEnhancedService) GetRenewalPrice(hostID uint, cycle string) (*RenewalPriceData, error) {
	host, err := s.hostSvc.GetByID(hostID)
	if err != nil {
		return nil, fmt.Errorf("host not found: %w", err)
	}

	price, currency := s.calculateRenewalPrice(host, cycle)

	var newExpiredAt *time.Time
	base := time.Now()
	if host.ExpiredAt != nil && host.ExpiredAt.After(time.Now()) {
		base = *host.ExpiredAt
	}
	t := s.extendDate(base, cycle)
	newExpiredAt = &t

	return &RenewalPriceData{
		HostID:       hostID,
		Cycle:        cycle,
		Price:        price,
		Currency:     currency,
		NewExpiredAt: newExpiredAt,
	}, nil
}

// SubmitRenewal creates a renewal order and invoice for a host.
func (s *HostEnhancedService) SubmitRenewal(userID, hostID uint, cycle, paymentMethod string) (*RenewResult, error) {
	host, err := s.hostSvc.GetByID(hostID)
	if err != nil {
		return nil, fmt.Errorf("host not found: %w", err)
	}

	price, _ := s.calculateRenewalPrice(host, cycle)
	if price <= 0 {
		return nil, errors.New("invalid renewal price")
	}

	result := &RenewResult{ServiceID: hostID, Amount: price}

	order := &model.Order{
		OrderNo:       util.GenerateOrderNo(),
		UserID:        userID,
		Type:          "renew",
		Amount:        price,
		Total:         price,
		Currency:      "CNY",
		BillingCycle:  cycle,
		Quantity:      1,
		Description:   fmt.Sprintf("Renewal for host %s (%s)", host.Hostname, cycle),
		Status:        0,
		PaymentStatus: 0,
	}
	if host.ProductID != nil {
		order.ProductID = *host.ProductID
	}

	if err := s.db.Create(order).Error; err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}
	result.OrderID = order.ID

	invoice, err := s.invSvc.Create(userID, order.ID, order.OrderNo, price)
	if err != nil {
		return nil, fmt.Errorf("create invoice: %w", err)
	}
	result.InvoiceID = invoice.ID

	if paymentMethod == "balance" || paymentMethod == "" {
		if err := s.balSvc.Deduct(userID, price, fmt.Sprintf("Renewal order %s", order.OrderNo)); err != nil {
			result.Error = fmt.Sprintf("balance payment failed: %v", err)
			return result, nil
		}
		now := time.Now()
		s.db.Model(order).Updates(map[string]interface{}{
			"status":         1,
			"payment_status": 1,
			"payment_method": "balance",
			"paid_at":        &now,
		})
		s.db.Model(invoice).Updates(map[string]interface{}{
			"status": 1,
			"paid_at": &now,
		})
		result.PaymentOK = true
	}

	base := time.Now()
	if host.ExpiredAt != nil && host.ExpiredAt.After(time.Now()) {
		base = *host.ExpiredAt
	}
	newDueDate := s.extendDate(base, cycle)

	s.db.Model(host).Updates(map[string]interface{}{
		"expired_at": &newDueDate,
		"status":     1,
	})
	result.NewDueDate = &newDueDate

	s.log.Infof("host renewed: host=%d user=%d cycle=%s new_due=%s price=%.2f",
		hostID, userID, cycle, newDueDate.Format("2006-01-02"), price)

	return result, nil
}

// SetAutoRenew toggles auto-renewal for a host.
func (s *HostEnhancedService) SetAutoRenew(hostID uint, enabled bool) error {
	var rc model.RenewCycle
	err := s.db.Where("user_product_id = ?", hostID).First(&rc).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		rc = model.RenewCycle{
			UserProductID: hostID,
			AutoRenew:     enabled,
			Status:        "active",
		}
		return s.db.Create(&rc).Error
	}
	if err != nil {
		return err
	}
	return s.db.Model(&rc).Update("auto_renew", enabled).Error
}

// BatchRenew renews multiple hosts with the same billing cycle.
func (s *HostEnhancedService) BatchRenew(userID uint, hostIDs []uint, cycle string) (*BatchRenewSummary, error) {
	summary := &BatchRenewSummary{
		Results: make([]RenewResult, 0, len(hostIDs)),
	}

	for _, hostID := range hostIDs {
		result := s.renewSingleHost(userID, hostID, cycle)
		summary.Results = append(summary.Results, result)

		if result.Error != "" {
			summary.Failed++
		} else {
			summary.Renewed++
			summary.TotalAmount += result.Amount
		}
	}

	s.log.WithFields(map[string]interface{}{
		"user_id":   userID,
		"cycle":     cycle,
		"renewed":   summary.Renewed,
		"failed":    summary.Failed,
		"total_amt": summary.TotalAmount,
	}).Info("batch host renew completed")

	return summary, nil
}

// ═══════════════════ Transfer ═══════════════════

// TransferHost transfers a host to another user.
func (s *HostEnhancedService) TransferHost(hostID, fromUserID, toUserID uint, reason string) error {
	host, err := s.hostSvc.GetByID(hostID)
	if err != nil {
		return fmt.Errorf("host not found: %w", err)
	}
	if host.OwnerID == nil || *host.OwnerID != fromUserID {
		return errors.New("host ownership mismatch")
	}

	var toUserCount int64
	if err := s.db.Table("users").Where("id = ?", toUserID).Count(&toUserCount).Error; err != nil || toUserCount == 0 {
		return errors.New("target user not found")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(host).Update("owner_id", toUserID).Error; err != nil {
			return err
		}

		log := HostTransferLog{
			HostID:     hostID,
			FromUserID: fromUserID,
			ToUserID:   toUserID,
			Reason:     reason,
		}
		return tx.Create(&log).Error
	})
}

// GetTransferHistory returns transfer history for a host.
func (s *HostEnhancedService) GetTransferHistory(hostID uint) ([]HostTransferLog, error) {
	var logs []HostTransferLog
	if err := s.db.Where("host_id = ?", hostID).Order("id DESC").Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

// ═══════════════════ Categories ═══════════════════

// GetHostCategories returns categories for a user.
func (s *HostEnhancedService) GetHostCategories(userID uint) ([]HostCategory, error) {
	var cats []HostCategory
	if err := s.db.Where("user_id = ?", userID).Order("sort_order ASC, id ASC").Find(&cats).Error; err != nil {
		return nil, err
	}
	return cats, nil
}

// CreateHostCategory creates a new host category.
func (s *HostEnhancedService) CreateHostCategory(userID uint, name string) (*HostCategory, error) {
	cat := &HostCategory{
		UserID: userID,
		Name:   name,
	}
	if err := s.db.Create(cat).Error; err != nil {
		return nil, err
	}
	return cat, nil
}

// AssignCategory assigns a host to a category.
func (s *HostEnhancedService) AssignCategory(hostID, categoryID uint) error {
	var cat HostCategory
	if err := s.db.First(&cat, categoryID).Error; err != nil {
		return errors.New("category not found")
	}

	var existing HostCategoryAssignment
	err := s.db.Where("host_id = ?", hostID).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return s.db.Create(&HostCategoryAssignment{HostID: hostID, CategoryID: categoryID}).Error
	}
	if err != nil {
		return err
	}
	return s.db.Model(&existing).Update("category_id", categoryID).Error
}

// DeleteCategory deletes a category and its assignments.
func (s *HostEnhancedService) DeleteCategory(userID, categoryID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND user_id = ?", categoryID, userID).Delete(&HostCategory{}).Error; err != nil {
			return err
		}
		return tx.Where("category_id = ?", categoryID).Delete(&HostCategoryAssignment{}).Error
	})
}

// ═══════════════════ SSL Config ═══════════════════

// GetSSLConfig returns SSL configuration for a host.
func (s *HostEnhancedService) GetSSLConfig(hostID uint) (*HostSSLConfig, error) {
	var cfg HostSSLConfig
	err := s.db.Where("host_id = ?", hostID).First(&cfg).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &HostSSLConfig{HostID: hostID, Status: "none"}, nil
	}
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// SetSSLConfig updates SSL configuration for a host.
func (s *HostEnhancedService) SetSSLConfig(hostID uint, cfgType, domain string) (*HostSSLConfig, error) {
	var cfg HostSSLConfig
	err := s.db.Where("host_id = ?", hostID).First(&cfg).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		cfg = HostSSLConfig{
			HostID: hostID,
			Type:   cfgType,
			Domain: domain,
			Status: "pending",
		}
		if err := s.db.Create(&cfg).Error; err != nil {
			return nil, err
		}
		return &cfg, nil
	}
	if err != nil {
		return nil, err
	}

	cfg.Type = cfgType
	cfg.Domain = domain
	cfg.Status = "pending"
	if err := s.db.Save(&cfg).Error; err != nil {
		return nil, err
	}
	return &cfg, nil
}

// InstallSSL installs an SSL certificate for a host.
func (s *HostEnhancedService) InstallSSL(hostID uint, cert, key, ca string) error {
	var cfg HostSSLConfig
	err := s.db.Where("host_id = ?", hostID).First(&cfg).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		cfg = HostSSLConfig{
			HostID: hostID,
			Type:   "custom",
		}
	}
	cfg.Cert = cert
	cfg.Key = key
	cfg.CA = ca
	cfg.Status = "active"
	now := time.Now()
	cfg.IssuedAt = &now

	if cfg.ID == 0 {
		return s.db.Create(&cfg).Error
	}
	return s.db.Save(&cfg).Error
}

// ═══════════════════ Downstream / Second Verify ═══════════════════

// SetDownstream sets a downstream provider for a host.
func (s *HostEnhancedService) SetDownstream(hostID uint, downstreamID string) error {
	host, err := s.hostSvc.GetByID(hostID)
	if err != nil {
		return fmt.Errorf("host not found: %w", err)
	}

	metadata := map[string]interface{}{}
	if len(host.Metadata) > 0 {
		_ = json.Unmarshal(host.Metadata, &metadata)
	}
	metadata["downstream_id"] = downstreamID
	data, _ := json.Marshal(metadata)
	return s.db.Model(host).Update("metadata", datatypes.JSON(data)).Error
}

// SetSecondVerify toggles second verification for a host.
func (s *HostEnhancedService) SetSecondVerify(hostID uint, enabled bool) error {
	var sv HostSecondVerify
	err := s.db.Where("host_id = ?", hostID).First(&sv).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		sv = HostSecondVerify{
			HostID:  hostID,
			Enabled: enabled,
		}
		return s.db.Create(&sv).Error
	}
	if err != nil {
		return err
	}
	return s.db.Model(&sv).Update("enabled", enabled).Error
}

// VerifySecond verifies a second factor code for a host operation.
func (s *HostEnhancedService) VerifySecond(hostID uint, code string) (bool, error) {
	var sv HostSecondVerify
	if err := s.db.Where("host_id = ? AND enabled = true", hostID).First(&sv).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true, nil // no second verify configured, pass through
		}
		return false, err
	}

	// TOTP verification: for now, accept any 6-digit code when secret is set.
	// Real implementation would validate against sv.Secret using TOTP algorithm.
	if len(code) == 6 {
		return true, nil
	}
	return false, nil
}

// ═══════════════════ Traffic & Power ═══════════════════

// GetTrafficUsage returns traffic usage for a host.
func (s *HostEnhancedService) GetTrafficUsage(hostID uint) (*TrafficUsageData, error) {
	host, err := s.hostSvc.GetByID(hostID)
	if err != nil {
		return nil, fmt.Errorf("host not found: %w", err)
	}

	totalGB := float64(host.TrafficGB)
	if totalGB <= 0 {
		totalGB = 0 // unlimited
	}

	// Traffic data from host metadata
	var usedGB float64
	if len(host.Metadata) > 0 {
		var meta map[string]interface{}
		if err := json.Unmarshal(host.Metadata, &meta); err == nil {
			if v, ok := meta["traffic_used_gb"].(float64); ok {
				usedGB = v
			}
		}
	}

	var usagePct float64
	if totalGB > 0 {
		usagePct = usedGB / totalGB * 100
	}

	var resetDate *time.Time
	if len(host.Metadata) > 0 {
		var meta map[string]interface{}
		if err := json.Unmarshal(host.Metadata, &meta); err == nil {
			if v, ok := meta["traffic_reset_date"].(string); ok && v != "" {
				if t, err := time.Parse("2006-01-02", v); err == nil {
					resetDate = &t
				}
			}
		}
	}

	return &TrafficUsageData{
		HostID:      hostID,
		UsedGB:      usedGB,
		TotalGB:     totalGB,
		RemainingGB: totalGB - usedGB,
		UsagePct:    usagePct,
		ResetDate:   resetDate,
	}, nil
}

// GetTrafficChart returns traffic chart data for a host.
func (s *HostEnhancedService) GetTrafficChart(hostID uint, period string) ([]TrafficChartData, error) {
	_, err := s.hostSvc.GetByID(hostID)
	if err != nil {
		return nil, fmt.Errorf("host not found: %w", err)
	}

	days := 30
	switch period {
	case "7d":
		days = 7
	case "90d":
		days = 90
	}

	since := time.Now().AddDate(0, 0, -days)
	var logs []model.DcimTrafficLog
	if err := s.db.Where("server_id = ? AND recorded_at >= ?", hostID, since).
		Order("recorded_at ASC").Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("query traffic logs: %w", err)
	}

	// Aggregate by date
	type daily struct {
		inBytes  int64
		outBytes int64
	}
	dateMap := make(map[string]*daily)
	for _, l := range logs {
		key := l.RecordedAt.Format("2006-01-02")
		if _, ok := dateMap[key]; !ok {
			dateMap[key] = &daily{}
		}
		dateMap[key].inBytes += l.InBytes
		dateMap[key].outBytes += l.OutBytes
	}

	charts := make([]TrafficChartData, 0, days)
	now := time.Now()
	for i := days - 1; i >= 0; i-- {
		d := now.AddDate(0, 0, -i)
		key := d.Format("2006-01-02")
		inGB, outGB := 0.0, 0.0
		if dd, ok := dateMap[key]; ok {
			inGB = float64(dd.inBytes) / (1024 * 1024 * 1024)
			outGB = float64(dd.outBytes) / (1024 * 1024 * 1024)
		}
		charts = append(charts, TrafficChartData{
			Date:    key,
			InGB:    math.Round(inGB*100) / 100,
			OutGB:   math.Round(outGB*100) / 100,
			TotalGB: math.Round((inGB+outGB)*100) / 100,
		})
	}

	return charts, nil
}

// RefreshPowerStatus refreshes power status from upstream.
func (s *HostEnhancedService) RefreshPowerStatus(hostID uint) (*HostStatusData, error) {
	host, err := s.hostSvc.GetByID(hostID)
	if err != nil {
		return nil, fmt.Errorf("host not found: %w", err)
	}

	upstreamErr := s.hostSvc.callUpstreamAction(host, "status", "")
	if upstreamErr != nil {
		s.log.Warnf("upstream status refresh failed for host %d: %v", hostID, upstreamErr)
	}

	return &HostStatusData{
		HostID:      hostID,
		Status:      int(host.Status),
		PowerStatus: int(host.PowerStatus),
		Hostname:    host.Hostname,
		IP:          host.IP,
		OS:          host.OS,
	}, nil
}

// ═══════════════════ Host Info ═══════════════════

// GetHostDetail returns full host detail with all extended info.
func (s *HostEnhancedService) GetHostDetail(hostID uint) (*HostDetailData, error) {
	host, err := s.hostSvc.GetByID(hostID)
	if err != nil {
		return nil, fmt.Errorf("host not found: %w", err)
	}

	detail := &HostDetailData{
		Host:   host,
		Remark: host.Remark,
		Tags:   host.Tags,
	}

	if host.ProductID != nil {
		var product Product
		if err := s.db.First(&product, *host.ProductID).Error; err == nil {
			detail.Product = &product
		}
	}

	var assign HostCategoryAssignment
	if err := s.db.Where("host_id = ?", hostID).First(&assign).Error; err == nil {
		var cat HostCategory
		if err := s.db.First(&cat, assign.CategoryID).Error; err == nil {
			detail.Category = &cat
		}
	}

	ssl, _ := s.GetSSLConfig(hostID)
	if ssl != nil {
		detail.SSLConfig = ssl
	}

	var sv HostSecondVerify
	if err := s.db.Where("host_id = ?", hostID).First(&sv).Error; err == nil {
		detail.SecondVerify = &sv
	}

	traffic, _ := s.GetTrafficUsage(hostID)
	if traffic != nil {
		detail.Traffic = traffic
	}

	return detail, nil
}

// GetHostStatus returns real-time status for a host.
func (s *HostEnhancedService) GetHostStatus(hostID uint) (*HostStatusData, error) {
	host, err := s.hostSvc.GetByID(hostID)
	if err != nil {
		return nil, fmt.Errorf("host not found: %w", err)
	}

	return &HostStatusData{
		HostID:      hostID,
		Status:      int(host.Status),
		PowerStatus: int(host.PowerStatus),
		Hostname:    host.Hostname,
		IP:          host.IP,
		OS:          host.OS,
	}, nil
}

// PostRemark updates the remark for a host.
func (s *HostEnhancedService) PostRemark(hostID, userID uint, remark string) error {
	host, err := s.hostSvc.GetByID(hostID)
	if err != nil {
		return fmt.Errorf("host not found: %w", err)
	}
	if host.OwnerID == nil || *host.OwnerID != userID {
		return errors.New("access denied")
	}
	return s.db.Model(host).Update("remark", remark).Error
}

// HideHost hides a host from the user's list.
func (s *HostEnhancedService) HideHost(hostID, userID uint) error {
	host, err := s.hostSvc.GetByID(hostID)
	if err != nil {
		return fmt.Errorf("host not found: %w", err)
	}
	if host.OwnerID == nil || *host.OwnerID != userID {
		return errors.New("access denied")
	}

	tags := map[string]interface{}{}
	if len(host.Tags) > 0 {
		_ = json.Unmarshal(host.Tags, &tags)
	}
	tags["hidden"] = true
	data, _ := json.Marshal(tags)
	return s.db.Model(host).Update("tags", datatypes.JSON(data)).Error
}

// ═══════════════════ Cancel / Terminate ═══════════════════

// GetCancelPage returns cancellation options for a host.
func (s *HostEnhancedService) GetCancelPage(hostID, userID uint) (*CancelPageData, error) {
	host, err := s.hostSvc.GetByID(hostID)
	if err != nil {
		return nil, fmt.Errorf("host not found: %w", err)
	}
	if host.OwnerID == nil || *host.OwnerID != userID {
		return nil, errors.New("access denied")
	}

	data := &CancelPageData{
		HostID:      hostID,
		Hostname:    host.Hostname,
		ExpiredAt:   host.ExpiredAt,
		CancelTypes: []string{"immediate", "end_of_period"},
	}

	var pending HostCancelRequest
	if err := s.db.Where("host_id = ? AND status = ?", hostID, "pending").First(&pending).Error; err == nil {
		data.HasPending = true
		data.PendingRequest = &pending
	}

	return data, nil
}

// SubmitCancel creates a cancellation request for a host.
func (s *HostEnhancedService) SubmitCancel(hostID, userID uint, reason, cancelType string) (*HostCancelRequest, error) {
	host, err := s.hostSvc.GetByID(hostID)
	if err != nil {
		return nil, fmt.Errorf("host not found: %w", err)
	}
	if host.OwnerID == nil || *host.OwnerID != userID {
		return nil, errors.New("access denied")
	}

	var existing HostCancelRequest
	if err := s.db.Where("host_id = ? AND status = ?", hostID, "pending").First(&existing).Error; err == nil {
		return nil, errors.New("cancellation request already pending")
	}

	if cancelType == "" {
		cancelType = "immediate"
	}

	req := &HostCancelRequest{
		HostID: hostID,
		UserID: userID,
		Reason: reason,
		Type:   cancelType,
		Status: "pending",
	}
	if err := s.db.Create(req).Error; err != nil {
		return nil, err
	}

	s.log.Infof("host cancel request: host=%d user=%d type=%s", hostID, userID, cancelType)
	return req, nil
}

// DeleteCancel cancels a pending cancellation request.
func (s *HostEnhancedService) DeleteCancel(hostID, userID uint) error {
	result := s.db.Where("host_id = ? AND user_id = ? AND status = ?", hostID, userID, "pending").
		Delete(&HostCancelRequest{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no pending cancellation request found")
	}
	return nil
}

// ═══════════════════ Dedicated Server ═══════════════════

// GetDedicatedServer returns dedicated server specific info for a host.
func (s *HostEnhancedService) GetDedicatedServer(hostID uint) (*DedicatedServerData, error) {
	host, err := s.hostSvc.GetByID(hostID)
	if err != nil {
		return nil, fmt.Errorf("host not found: %w", err)
	}

	data := &DedicatedServerData{
		HostID:        hostID,
		Rack:          host.Rack,
		RackPosition:  host.RackPosition,
		BandwidthMbps: host.BandwidthMbps,
	}

	if len(host.Metadata) > 0 {
		var meta map[string]interface{}
		if err := json.Unmarshal(host.Metadata, &meta); err == nil {
			if v, ok := meta["port_speed"].(string); ok {
				data.PortSpeed = v
			}
			if v, ok := meta["pdu"].(string); ok {
				data.PDU = v
			}
			if v, ok := meta["switch"].(string); ok {
				data.Switch = v
			}
			if v, ok := meta["remote_hand"].(bool); ok {
				data.RemoteHand = v
			}
			if v, ok := meta["ipmi"].(string); ok {
				data.IPMI = v
			}
		}
	}

	return data, nil
}

// ═══════════════════ Flow Packets ═══════════════════

// GetFlowPackets lists available flow packets for a host.
func (s *HostEnhancedService) GetFlowPackets(hostID uint) ([]FlowPacket, error) {
	host, err := s.hostSvc.GetByID(hostID)
	if err != nil {
		return nil, fmt.Errorf("host not found: %w", err)
	}

	var packets []FlowPacket
	query := s.db.Model(&FlowPacket{}).Where("status = 1")
	if host.ProductID != nil {
		// Filter packets applicable to this product if needed
		_ = *host.ProductID
	}
	if err := query.Order("sort_order ASC, id ASC").Find(&packets).Error; err != nil {
		return nil, err
	}

	for i := range packets {
		packets[i].HostID = hostID
	}
	return packets, nil
}

// BuyFlowPacket purchases a flow packet for a host.
func (s *HostEnhancedService) BuyFlowPacket(userID, hostID, packetID uint) (*HostFlowPacketPurchase, error) {
	host, err := s.hostSvc.GetByID(hostID)
	if err != nil {
		return nil, fmt.Errorf("host not found: %w", err)
	}
	if host.OwnerID == nil || *host.OwnerID != userID {
		return nil, errors.New("access denied")
	}

	var packet FlowPacket
	if err := s.db.First(&packet, packetID).Error; err != nil {
		return nil, errors.New("flow packet not found")
	}

	purchase := &HostFlowPacketPurchase{
		HostID:       hostID,
		FlowPacketID: packetID,
		UserID:       userID,
		AmountGB:     packet.AmountGB,
		Price:        packet.Price,
		Status:       "pending",
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(purchase).Error; err != nil {
			return err
		}

		order := &model.Order{
			OrderNo:       util.GenerateOrderNo(),
			UserID:        userID,
			Type:          "new",
			Amount:        packet.Price,
			Total:         packet.Price,
			Currency:      "CNY",
			BillingCycle:  "onetime",
			Quantity:      1,
			Description:   fmt.Sprintf("Flow packet %s (%dGB) for host %s", packet.Name, packet.AmountGB, host.Hostname),
			Status:        0,
			PaymentStatus: 0,
		}
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		purchase.OrderID = order.ID

		if err := tx.Model(purchase).Update("order_id", order.ID).Error; err != nil {
			return err
		}

		if err := s.balSvc.Deduct(userID, packet.Price, fmt.Sprintf("Flow packet order %s", order.OrderNo)); err != nil {
			return fmt.Errorf("balance payment failed: %w", err)
		}

		now := time.Now()
		tx.Model(order).Updates(map[string]interface{}{
			"status":         1,
			"payment_status": 1,
			"payment_method": "balance",
			"paid_at":        &now,
		})

		// Add traffic to host
		currentUsed := 0.0
		metadata := map[string]interface{}{}
		if len(host.Metadata) > 0 {
			if err := json.Unmarshal(host.Metadata, &metadata); err == nil {
				if v, ok := metadata["traffic_used_gb"].(float64); ok {
					currentUsed = v
				}
			}
		}
		metadata["flow_packet_added_gb"] = float64(packet.AmountGB)
		data, _ := json.Marshal(metadata)
		tx.Model(host).Update("metadata", datatypes.JSON(data))

		tx.Model(purchase).Update("status", "paid")
		return nil
	})

	if err != nil {
		return nil, err
	}

	s.log.Infof("flow packet purchased: host=%d packet=%d amount=%dGB", hostID, packetID, packet.AmountGB)
	return purchase, nil
}

// ═══════════════════ Reinstall (enhanced) ═══════════════════

// CheckReinstall checks if reinstall is allowed for a host.
func (s *HostEnhancedService) CheckReinstall(hostID uint) (*ReinstallCheckData, error) {
	host, err := s.hostSvc.GetByID(hostID)
	if err != nil {
		return nil, fmt.Errorf("host not found: %w", err)
	}

	data := &ReinstallCheckData{
		HostStatus: int(host.Status),
	}

	if host.Status == 3 {
		data.Allowed = false
		data.Reason = "host is in maintenance/reinstalling"
		return data, nil
	}

	if host.Status == 1 {
		data.Allowed = false
		data.Reason = "host must be stopped before reinstall"
		return data, nil
	}

	data.Allowed = true
	data.AvailableOS = []string{"ubuntu-22.04", "ubuntu-20.04", "debian-12", "debian-11", "centos-7", "almalinux-9"}

	return data, nil
}

// GetReinstallStatus returns reinstall progress for a host.
func (s *HostEnhancedService) GetReinstallStatus(hostID uint) (*ReinstallStatusData, error) {
	host, err := s.hostSvc.GetByID(hostID)
	if err != nil {
		return nil, fmt.Errorf("host not found: %w", err)
	}

	data := &ReinstallStatusData{
		InProgress: host.Status == 3,
	}

	if host.Status == 3 {
		var lastOp HostOperation
		if err := s.db.Where("host_id = ? AND action = ?", hostID, "reinstall").
			Order("id DESC").First(&lastOp).Error; err == nil {
			data.StartedAt = &lastOp.StartedAt
			if lastOp.Status == 1 {
				data.Progress = 50
			} else if lastOp.Status == 2 {
				data.Progress = 100
				data.InProgress = false
			}
		}
	}

	return data, nil
}

// CancelReinstall cancels an in-progress reinstall.
func (s *HostEnhancedService) CancelReinstall(hostID uint) error {
	host, err := s.hostSvc.GetByID(hostID)
	if err != nil {
		return fmt.Errorf("host not found: %w", err)
	}
	if host.Status != 3 {
		return errors.New("host is not in reinstall state")
	}

	var lastOp HostOperation
	if err := s.db.Where("host_id = ? AND action = ? AND status = ?", hostID, "reinstall", 1).
		Order("id DESC").First(&lastOp).Error; err != nil {
		return errors.New("no active reinstall operation found")
	}

	now := time.Now()
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&lastOp).Updates(map[string]interface{}{
			"status":      3,
			"error_msg":   "cancelled by user",
			"finished_at": &now,
		}).Error; err != nil {
			return err
		}
		return tx.Model(host).Update("status", 0).Error
	})
}

// ═══════════════════ Host Recharge ═══════════════════

// GetHostRecharge returns recharge options for a host.
func (s *HostEnhancedService) GetHostRecharge(hostID, userID uint) (*HostRechargeData, error) {
	host, err := s.hostSvc.GetByID(hostID)
	if err != nil {
		return nil, fmt.Errorf("host not found: %w", err)
	}
	if host.OwnerID == nil || *host.OwnerID != userID {
		return nil, errors.New("access denied")
	}

	var balance float64
	s.db.Table("users").Where("id = ?", userID).Select("balance").Scan(&balance)

	data := &HostRechargeData{
		HostID:      hostID,
		Hostname:    host.Hostname,
		Balance:     balance,
		ExpiredAt:   host.ExpiredAt,
		MinRecharge: 10,
	}

	if host.ProductID != nil {
		var product Product
		if err := s.db.First(&product, *host.ProductID).Error; err == nil {
			data.DueAmount = product.Price
		}
	}

	return data, nil
}

// ═══════════════════ Internal helpers ═══════════════════

func (s *HostEnhancedService) calculateRenewalPrice(host *Host, cycle string) (float64, string) {
	if host.ProductID != nil {
		var product Product
		if err := s.db.First(&product, *host.ProductID).Error; err == nil {
			return product.Price, "CNY"
		}
	}

	switch cycle {
	case "monthly":
		return 99.0, "CNY"
	case "quarterly":
		return 269.0, "CNY"
	case "semi-annually":
		return 499.0, "CNY"
	case "annually":
		return 899.0, "CNY"
	case "biennially":
		return 1599.0, "CNY"
	case "triennially":
		return 2199.0, "CNY"
	default:
		return 99.0, "CNY"
	}
}

func (s *HostEnhancedService) extendDate(from time.Time, cycle string) time.Time {
	switch cycle {
	case "monthly":
		return from.AddDate(0, 1, 0)
	case "quarterly":
		return from.AddDate(0, 3, 0)
	case "semi-annually":
		return from.AddDate(0, 6, 0)
	case "annually":
		return from.AddDate(1, 0, 0)
	case "biennially":
		return from.AddDate(2, 0, 0)
	case "triennially":
		return from.AddDate(3, 0, 0)
	default:
		return from.AddDate(0, 1, 0)
	}
}

func (s *HostEnhancedService) renewSingleHost(userID, hostID uint, cycle string) RenewResult {
	result := RenewResult{ServiceID: hostID}

	host, err := s.hostSvc.GetByID(hostID)
	if err != nil {
		result.Error = "host not found"
		return result
	}
	if host.OwnerID == nil || *host.OwnerID != userID {
		result.Error = "access denied"
		return result
	}

	price, _ := s.calculateRenewalPrice(host, cycle)
	if price <= 0 {
		result.Error = "invalid price"
		return result
	}
	result.Amount = price

	order := &model.Order{
		OrderNo:       util.GenerateOrderNo(),
		UserID:        userID,
		Type:          "renew",
		Amount:        price,
		Total:         price,
		Currency:      "CNY",
		BillingCycle:  cycle,
		Quantity:      1,
		Description:   fmt.Sprintf("Renewal for host %s (%s)", host.Hostname, cycle),
		Status:        0,
		PaymentStatus: 0,
	}
	if host.ProductID != nil {
		order.ProductID = *host.ProductID
	}
	if err := s.db.Create(order).Error; err != nil {
		result.Error = fmt.Sprintf("create order: %v", err)
		return result
	}
	result.OrderID = order.ID

	invoice, err := s.invSvc.Create(userID, order.ID, order.OrderNo, price)
	if err != nil {
		result.Error = fmt.Sprintf("create invoice: %v", err)
		return result
	}
	result.InvoiceID = invoice.ID

	if err := s.balSvc.Deduct(userID, price, fmt.Sprintf("Renewal order %s", order.OrderNo)); err != nil {
		result.Error = fmt.Sprintf("balance payment failed: %v", err)
		return result
	}

	now := time.Now()
	s.db.Model(order).Updates(map[string]interface{}{
		"status":         1,
		"payment_status": 1,
		"payment_method": "balance",
		"paid_at":        &now,
	})
	s.db.Model(invoice).Updates(map[string]interface{}{
		"status": 1,
		"paid_at": &now,
	})
	result.PaymentOK = true

	base := time.Now()
	if host.ExpiredAt != nil && host.ExpiredAt.After(time.Now()) {
		base = *host.ExpiredAt
	}
	newDueDate := s.extendDate(base, cycle)

	s.db.Model(host).Updates(map[string]interface{}{
		"expired_at": &newDueDate,
		"status":     1,
	})
	result.NewDueDate = &newDueDate

	return result
}
