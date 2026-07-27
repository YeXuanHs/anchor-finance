package service

import (
	"errors"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"anchorfinance/pkg/logger"
)

// Host 主机
type Host struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Hostname       string         `gorm:"size:256;not null;index" json:"hostname"`
	IP             string         `gorm:"size:45;index" json:"ip"`
	IPv6           string         `gorm:"size:128" json:"ipv6"`
	OS             string         `gorm:"size:128" json:"os"`
	OSVersion      string         `gorm:"size:64" json:"os_version"`
	CPU            string         `gorm:"size:128" json:"cpu"`
	CPUCores       int            `gorm:"default:0" json:"cpu_cores"`
	MemoryMB       int            `gorm:"default:0" json:"memory_mb"`
	DiskSizeGB     int            `gorm:"default:0" json:"disk_size_gb"`
	DiskType       string         `gorm:"size:32" json:"disk_type"`
	BandwidthMbps  int            `gorm:"default:0" json:"bandwidth_mbps"`
	TrafficGB      int            `gorm:"default:0" json:"traffic_gb"`
	Location       string         `gorm:"size:256" json:"location"`
	Status         int            `gorm:"default:0;index" json:"status"` // 0=关机 1=运行中 2=故障 3=维护
	PowerStatus    int            `gorm:"default:0" json:"power_status"`
	OwnerID        *uint          `gorm:"index" json:"owner_id"`
	ProductID      *uint          `gorm:"index" json:"product_id"`
	OrderID        *uint          `gorm:"index" json:"order_id"`
	ExpiredAt      *time.Time     `gorm:"index" json:"expired_at"`
	ProvisionedAt  *time.Time     `json:"provisioned_at"`
	Remark         string         `gorm:"type:text" json:"remark"`
	AdminNotes     string         `gorm:"type:text" json:"admin_notes"`
	Tags           datatypes.JSON `gorm:"type:jsonb" json:"tags"`
	Config         datatypes.JSON `gorm:"type:jsonb" json:"config"`
	Metadata       datatypes.JSON `gorm:"type:jsonb" json:"metadata"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type HostOperation struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	HostID     uint       `gorm:"index;not null" json:"host_id"`
	OperatorID uint       `gorm:"index;not null" json:"operator_id"`
	Action     string     `gorm:"size:32;not null" json:"action"` // boot/shutdown/reboot/reinstall
	Params     string     `gorm:"type:text" json:"params"`
	Status     int        `gorm:"default:1" json:"status"` // 1=执行中 2=成功 3=失败
	Result     string     `gorm:"type:text" json:"result"`
	ErrorMsg   string     `gorm:"type:text" json:"error_msg"`
	StartedAt  time.Time  `gorm:"not null" json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type HostService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewHostService(db *gorm.DB, log *logger.Logger) *HostService {
	return &HostService{db: db, log: log}
}

type HostActionRequest struct {
	Action string `json:"action" binding:"required,oneof=boot shutdown reboot reinstall"`
	Params string `json:"params"` // JSON参数，如重装时的OS选择
}

// GetByID returns a single host by ID.
func (s *HostService) GetByID(id uint) (*Host, error) {
	var host Host
	if err := s.db.First(&host, id).Error; err != nil {
		return nil, err
	}
	return &host, nil
}

// GetList returns all hosts with pagination.
func (s *HostService) GetList(page, pageSize int, status *int, keyword string, ownerID *uint) ([]Host, int64, error) {
	var hosts []Host
	var total int64

	query := s.db.Model(&Host{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if keyword != "" {
		q := "%" + keyword + "%"
		query = query.Where("hostname LIKE ? OR ip LIKE ?", q, q)
	}
	if ownerID != nil {
		query = query.Where("owner_id = ?", *ownerID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&hosts).Error; err != nil {
		return nil, 0, err
	}
	return hosts, total, nil
}

// GetUserHosts returns hosts for a specific user.
func (s *HostService) GetUserHosts(userID uint, page, pageSize int) ([]Host, int64, error) {
	var hosts []Host
	var total int64

	query := s.db.Model(&Host{}).Where("owner_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&hosts).Error; err != nil {
		return nil, 0, err
	}
	return hosts, total, nil
}

// PerformAction executes an action on a host.
func (s *HostService) PerformAction(hostID, operatorID uint, req HostActionRequest) (*HostOperation, error) {
	host, err := s.GetByID(hostID)
	if err != nil {
		return nil, err
	}

	// Validate action based on current status
	switch req.Action {
	case "boot":
		if host.Status == 1 {
			return nil, errors.New("host is already running")
		}
	case "shutdown":
		if host.Status == 0 {
			return nil, errors.New("host is already stopped")
		}
	case "reboot":
		if host.Status != 1 {
			return nil, errors.New("host is not running")
		}
	case "reinstall":
		if host.Status == 1 {
			return nil, errors.New("host must be stopped before reinstall")
		}
	}

	operation := &HostOperation{
		HostID:     hostID,
		OperatorID: operatorID,
		Action:     req.Action,
		Params:     req.Params,
		Status:     1,
		StartedAt:  time.Now(),
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(operation).Error; err != nil {
			return err
		}

		// Update host status based on action
		newStatus := host.Status
		switch req.Action {
		case "boot":
			newStatus = 1
		case "shutdown":
			newStatus = 0
		case "reboot":
			newStatus = 1 // Will be 1 after reboot
		case "reinstall":
			newStatus = 3 // 维护中
		}

		return tx.Model(host).Updates(map[string]interface{}{
			"status":        newStatus,
			"power_status":  boolToInt(req.Action == "boot" || req.Action == "reboot"),
		}).Error
	})

	if err != nil {
		return nil, err
	}

	s.log.Infof("host action: host=%d action=%s operator=%d", hostID, req.Action, operatorID)
	return operation, nil
}

// CompleteOperation marks an operation as completed.
func (s *HostService) CompleteOperation(operationID uint, success bool, result, errMsg string) error {
	now := time.Now()
	status := 2
	if !success {
		status = 3
	}

	return s.db.Model(&HostOperation{}).Where("id = ?", operationID).Updates(map[string]interface{}{
		"status":      status,
		"result":      result,
		"error_msg":   errMsg,
		"finished_at": &now,
	}).Error
}

// GetHostOperations returns operations for a host.
func (s *HostService) GetHostOperations(hostID uint, page, pageSize int) ([]HostOperation, int64, error) {
	var ops []HostOperation
	var total int64

	query := s.db.Model(&HostOperation{}).Where("host_id = ?", hostID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&ops).Error; err != nil {
		return nil, 0, err
	}
	return ops, total, nil
}

// GetExpiringHosts returns hosts expiring within N days.
func (s *HostService) GetExpiringHosts(days int) ([]Host, error) {
	var hosts []Host
	future := time.Now().AddDate(0, 0, days)
	if err := s.db.Where("expired_at IS NOT NULL AND expired_at <= ? AND status = 1", future).
		Find(&hosts).Error; err != nil {
		return nil, err
	}
	return hosts, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
