package service

import (
	"errors"
	"fmt"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// DcimService DCIM业务逻辑
type DcimService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewDcimService(db *gorm.DB, log *logger.Logger) *DcimService {
	return &DcimService{db: db, log: log}
}

// ==================== 物理服务器 ====================

// GetServerByID 根据ID获取物理服务器
func (s *DcimService) GetServerByID(id uint) (*model.DcimServer, error) {
	var server model.DcimServer
	if err := s.db.Preload("Datacenter").Preload("Owner").First(&server, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("server not found")
		}
		return nil, err
	}
	return &server, nil
}

// GetServerList 获取物理服务器列表
func (s *DcimService) GetServerList(page, pageSize int, keyword string, status *int8, dcID *uint) ([]model.DcimServer, int64, error) {
	var servers []model.DcimServer
	var total int64

	query := s.db.Model(&model.DcimServer{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR ip LIKE ? OR hostname LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if dcID != nil {
		query = query.Where("datacenter_id = ?", *dcID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").
		Preload("Datacenter").Preload("Owner").Find(&servers).Error; err != nil {
		return nil, 0, err
	}

	return servers, total, nil
}

// CreateServer 创建物理服务器
func (s *DcimService) CreateServer(server *model.DcimServer) error {
	var count int64
	s.db.Model(&model.DcimServer{}).Where("ip = ?", server.IP).Count(&count)
	if count > 0 {
		return errors.New("IP address already exists")
	}
	if err := s.db.Create(server).Error; err != nil {
		return fmt.Errorf("create server: %w", err)
	}
	s.log.Infof("physical server created: id=%d ip=%s", server.ID, server.IP)
	return nil
}

// UpdateServer 更新物理服务器
func (s *DcimService) UpdateServer(id uint, updates map[string]interface{}) error {
	result := s.db.Model(&model.DcimServer{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("server not found")
	}
	return nil
}

// DeleteServer 删除物理服务器
func (s *DcimService) DeleteServer(id uint) error {
	var server model.DcimServer
	if err := s.db.First(&server, id).Error; err != nil {
		return errors.New("server not found")
	}
	if server.OwnerID != nil {
		return errors.New("cannot delete server that is assigned to a user")
	}
	if err := s.db.Delete(&server).Error; err != nil {
		return err
	}
	s.log.Infof("physical server deleted: id=%d", id)
	return nil
}

// BootServer 开机
func (s *DcimService) BootServer(serverID uint, operatorID uint) error {
	server, err := s.GetServerByID(serverID)
	if err != nil {
		return err
	}
	if server.Status == 1 {
		return errors.New("server is already running")
	}
	if server.Status == 2 {
		return errors.New("server is in fault state, please check hardware")
	}

	// 记录操作日志
	opLog := &model.DcimOperationLog{
		ServerType: "physical",
		ServerID:   serverID,
		OperatorID: operatorID,
		Action:     "boot",
		Status:     2,
		Result:     "Server boot initiated",
	}
	s.db.Create(opLog)

	// 更新服务器状态
	now := time.Now()
	s.db.Model(server).Updates(map[string]interface{}{
		"status":       1,
		"power_status": 1,
	})

	s.log.Infof("physical server boot: id=%d operator=%d", serverID, operatorID)
	return nil
}

// ShutdownServer 关机
func (s *DcimService) ShutdownServer(serverID uint, force bool, operatorID uint) error {
	server, err := s.GetServerByID(serverID)
	if err != nil {
		return err
	}
	if server.Status == 0 && server.PowerStatus == 0 {
		return errors.New("server is already shut down")
	}

	action := "shutdown"
	if force {
		action = "force_shutdown"
	}

	opLog := &model.DcimOperationLog{
		ServerType: "physical",
		ServerID:   serverID,
		OperatorID: operatorID,
		Action:     action,
		Status:     2,
		Result:     "Server shutdown initiated",
	}
	s.db.Create(opLog)

	s.db.Model(server).Updates(map[string]interface{}{
		"status":       0,
		"power_status": 0,
	})

	s.log.Infof("physical server shutdown: id=%d force=%v operator=%d", serverID, force, operatorID)
	return nil
}

// RebootServer 重启
func (s *DcimService) RebootServer(serverID uint, operatorID uint) error {
	server, err := s.GetServerByID(serverID)
	if err != nil {
		return err
	}
	if server.Status != 1 {
		return errors.New("server is not running")
	}

	opLog := &model.DcimOperationLog{
		ServerType: "physical",
		ServerID:   serverID,
		OperatorID: operatorID,
		Action:     "reboot",
		Status:     2,
		Result:     "Server reboot initiated",
	}
	s.db.Create(opLog)

	s.log.Infof("physical server reboot: id=%d operator=%d", serverID, operatorID)
	return nil
}

// ReinstallServer 重装系统
func (s *DcimService) ReinstallServer(serverID uint, os string, operatorID uint) error {
	server, err := s.GetServerByID(serverID)
	if err != nil {
		return err
	}

	if server.OwnerID == nil {
		return errors.New("server is not assigned to any user")
	}

	// 记录操作
	opLog := &model.DcimOperationLog{
		ServerType: "physical",
		ServerID:   serverID,
		OperatorID: operatorID,
		Action:     "reinstall",
		Params:     fmt.Sprintf(`{"os":"%s"}`, os),
		Status:     1,
	}
	s.db.Create(opLog)

	// 更新状态为维护
	s.db.Model(server).Updates(map[string]interface{}{
		"status": 3,
		"os":     os,
	})

	// 模拟重装完成
	go func() {
		time.Sleep(5 * time.Second)
		s.db.Model(server).Updates(map[string]interface{}{
			"status": 1,
		})
		now := time.Now()
		s.db.Model(opLog).Updates(map[string]interface{}{
			"status":      2,
			"result":      "OS reinstalled successfully",
			"finished_at": &now,
		})
		s.log.Infof("physical server reinstall completed: id=%d os=%s", serverID, os)
	}()

	s.log.Infof("physical server reinstall started: id=%d os=%s operator=%d", serverID, os, operatorID)
	return nil
}

// RenewServer 续费物理服务器
func (s *DcimService) RenewServer(serverID uint, months int, operatorID uint) error {
	server, err := s.GetServerByID(serverID)
	if err != nil {
		return err
	}
	if server.OwnerID == nil {
		return errors.New("server is not assigned to any user")
	}
	if months <= 0 || months > 120 {
		return errors.New("invalid renewal months (1-120)")
	}

	var newExpiredAt time.Time
	if server.ExpiredAt != nil && server.ExpiredAt.After(time.Now()) {
		newExpiredAt = server.ExpiredAt.AddDate(0, months, 0)
	} else {
		newExpiredAt = time.Now().AddDate(0, months, 0)
	}

	s.db.Model(server).Updates(map[string]interface{}{
		"expired_at": &newExpiredAt,
	})

	opLog := &model.DcimOperationLog{
		ServerType: "physical",
		ServerID:   serverID,
		OperatorID: operatorID,
		Action:     "renew",
		Params:     fmt.Sprintf(`{"months":%d}`, months),
		Status:     2,
		Result:     fmt.Sprintf("Renewed until %s", newExpiredAt.Format("2006-01-02")),
	}
	s.db.Create(opLog)

	s.log.Infof("physical server renewed: id=%d months=%d new_expiry=%s", serverID, months, newExpiredAt.Format("2006-01-02"))
	return nil
}

// ==================== 云服务器 ====================

// GetCloudByID 根据ID获取云服务器
func (s *DcimService) GetCloudByID(id uint) (*model.DcimCloud, error) {
	var cloud model.DcimCloud
	if err := s.db.Preload("Datacenter").Preload("Owner").Preload("ParentServer").First(&cloud, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("cloud server not found")
		}
		return nil, err
	}
	return &cloud, nil
}

// GetCloudList 获取云服务器列表
func (s *DcimService) GetCloudList(page, pageSize int, keyword string, status *int8, ownerID *uint) ([]model.DcimCloud, int64, error) {
	var clouds []model.DcimCloud
	var total int64

	query := s.db.Model(&model.DcimCloud{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR ip LIKE ? OR hostname LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if ownerID != nil {
		query = query.Where("owner_id = ?", *ownerID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").
		Preload("Datacenter").Preload("Owner").Find(&clouds).Error; err != nil {
		return nil, 0, err
	}

	return clouds, total, nil
}

// CreateCloud 创建云服务器
func (s *DcimService) CreateCloud(cloud *model.DcimCloud) error {
	var count int64
	s.db.Model(&model.DcimCloud{}).Where("ip = ?", cloud.IP).Count(&count)
	if count > 0 {
		return errors.New("IP address already exists")
	}
	if err := s.db.Create(cloud).Error; err != nil {
		return fmt.Errorf("create cloud server: %w", err)
	}
	s.log.Infof("cloud server created: id=%d ip=%s", cloud.ID, cloud.IP)
	return nil
}

// ProvisionCloud 开通云服务器
func (s *DcimService) ProvisionCloud(cloudID uint, operatorID uint) error {
	cloud, err := s.GetCloudByID(cloudID)
	if err != nil {
		return err
	}
	if cloud.Status == 1 {
		return errors.New("cloud server is already running")
	}

	now := time.Now()
	s.db.Model(cloud).Updates(map[string]interface{}{
		"status":         3,
		"provisioned_at": &now,
	})

	opLog := &model.DcimOperationLog{
		ServerType: "cloud",
		ServerID:   cloudID,
		OperatorID: operatorID,
		Action:     "provision",
		Status:     1,
	}
	s.db.Create(opLog)

	// 模拟开通完成
	go func() {
		time.Sleep(3 * time.Second)
		s.db.Model(cloud).Updates(map[string]interface{}{
			"status":       1,
			"power_status": 1,
		})
		finishTime := time.Now()
		s.db.Model(opLog).Updates(map[string]interface{}{
			"status":      2,
			"result":      "Cloud server provisioned successfully",
			"finished_at": &finishTime,
		})
		s.log.Infof("cloud server provisioned: id=%d", cloudID)
	}()

	s.log.Infof("cloud server provisioning started: id=%d operator=%d", cloudID, operatorID)
	return nil
}

// BootCloud 云服务器开机
func (s *DcimService) BootCloud(cloudID uint, operatorID uint) error {
	cloud, err := s.GetCloudByID(cloudID)
	if err != nil {
		return err
	}
	if cloud.Status == 1 {
		return errors.New("cloud server is already running")
	}

	s.db.Model(cloud).Updates(map[string]interface{}{
		"status":       1,
		"power_status": 1,
	})

	opLog := &model.DcimOperationLog{
		ServerType: "cloud",
		ServerID:   cloudID,
		OperatorID: operatorID,
		Action:     "boot",
		Status:     2,
		Result:     "Cloud server started",
	}
	s.db.Create(opLog)

	s.log.Infof("cloud server boot: id=%d operator=%d", cloudID, operatorID)
	return nil
}

// ShutdownCloud 云服务器关机
func (s *DcimService) ShutdownCloud(cloudID uint, force bool, operatorID uint) error {
	cloud, err := s.GetCloudByID(cloudID)
	if err != nil {
		return err
	}
	if cloud.Status == 0 && cloud.PowerStatus == 0 {
		return errors.New("cloud server is already shut down")
	}

	s.db.Model(cloud).Updates(map[string]interface{}{
		"status":       0,
		"power_status": 0,
	})

	action := "shutdown"
	if force {
		action = "force_shutdown"
	}
	opLog := &model.DcimOperationLog{
		ServerType: "cloud",
		ServerID:   cloudID,
		OperatorID: operatorID,
		Action:     action,
		Status:     2,
		Result:     "Cloud server shut down",
	}
	s.db.Create(opLog)

	s.log.Infof("cloud server shutdown: id=%d force=%v operator=%d", cloudID, force, operatorID)
	return nil
}

// RebootCloud 云服务器重启
func (s *DcimService) RebootCloud(cloudID uint, operatorID uint) error {
	cloud, err := s.GetCloudByID(cloudID)
	if err != nil {
		return err
	}
	if cloud.Status != 1 {
		return errors.New("cloud server is not running")
	}

	opLog := &model.DcimOperationLog{
		ServerType: "cloud",
		ServerID:   cloudID,
		OperatorID: operatorID,
		Action:     "reboot",
		Status:     2,
		Result:     "Cloud server rebooted",
	}
	s.db.Create(opLog)

	s.log.Infof("cloud server reboot: id=%d operator=%d", cloudID, operatorID)
	return nil
}

// ReinstallCloud 云服务器重装系统
func (s *DcimService) ReinstallCloud(cloudID uint, os string, operatorID uint) error {
	cloud, err := s.GetCloudByID(cloudID)
	if err != nil {
		return err
	}

	s.db.Model(cloud).Updates(map[string]interface{}{
		"status": 4,
		"os":     os,
	})

	opLog := &model.DcimOperationLog{
		ServerType: "cloud",
		ServerID:   cloudID,
		OperatorID: operatorID,
		Action:     "reinstall",
		Params:     fmt.Sprintf(`{"os":"%s"}`, os),
		Status:     1,
	}
	s.db.Create(opLog)

	go func() {
		time.Sleep(5 * time.Second)
		s.db.Model(cloud).Updates(map[string]interface{}{
			"status":       1,
			"power_status": 1,
		})
		now := time.Now()
		s.db.Model(opLog).Updates(map[string]interface{}{
			"status":      2,
			"result":      "OS reinstalled successfully",
			"finished_at": &now,
		})
		s.log.Infof("cloud server reinstall completed: id=%d", cloudID)
	}()

	s.log.Infof("cloud server reinstall started: id=%d os=%s operator=%d", cloudID, os, operatorID)
	return nil
}

// RenewCloud 续费云服务器
func (s *DcimService) RenewCloud(cloudID uint, months int, operatorID uint) error {
	cloud, err := s.GetCloudByID(cloudID)
	if err != nil {
		return err
	}
	if months <= 0 || months > 120 {
		return errors.New("invalid renewal months (1-120)")
	}

	var newExpiredAt time.Time
	if cloud.ExpiredAt != nil && cloud.ExpiredAt.After(time.Now()) {
		newExpiredAt = cloud.ExpiredAt.AddDate(0, months, 0)
	} else {
		newExpiredAt = time.Now().AddDate(0, months, 0)
	}

	s.db.Model(cloud).Updates(map[string]interface{}{
		"expired_at": &newExpiredAt,
	})

	opLog := &model.DcimOperationLog{
		ServerType: "cloud",
		ServerID:   cloudID,
		OperatorID: operatorID,
		Action:     "renew",
		Params:     fmt.Sprintf(`{"months":%d}`, months),
		Status:     2,
		Result:     fmt.Sprintf("Renewed until %s", newExpiredAt.Format("2006-01-02")),
	}
	s.db.Create(opLog)

	s.log.Infof("cloud server renewed: id=%d months=%d new_expiry=%s", cloudID, months, newExpiredAt.Format("2006-01-02"))
	return nil
}

// ==================== 机房管理 ====================

// GetDatacenterList 获取机房列表
func (s *DcimService) GetDatacenterList() ([]model.DcimDatacenter, error) {
	var dcs []model.DcimDatacenter
	if err := s.db.Where("status = 1").Order("name ASC").Find(&dcs).Error; err != nil {
		return nil, err
	}
	return dcs, nil
}

// CreateDatacenter 创建机房
func (s *DcimService) CreateDatacenter(dc *model.DcimDatacenter) error {
	var count int64
	s.db.Model(&model.DcimDatacenter{}).Where("code = ?", dc.Code).Count(&count)
	if count > 0 {
		return errors.New("datacenter code already exists")
	}
	return s.db.Create(dc).Error
}

// ==================== 操作日志 ====================

// GetOperationLogs 获取操作日志
func (s *DcimService) GetOperationLogs(serverType string, serverID uint, page, pageSize int) ([]model.DcimOperationLog, int64, error) {
	var logs []model.DcimOperationLog
	var total int64

	query := s.db.Model(&model.DcimOperationLog{})
	if serverType != "" {
		query = query.Where("server_type = ?", serverType)
	}
	if serverID > 0 {
		query = query.Where("server_id = ?", serverID)
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

// GetStats 获取DCIM统计信息
func (s *DcimService) GetStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var physicalTotal, physicalRunning, physicalOff int64
	s.db.Model(&model.DcimServer{}).Count(&physicalTotal)
	s.db.Model(&model.DcimServer{}).Where("status = 1").Count(&physicalRunning)
	s.db.Model(&model.DcimServer{}).Where("status = 0").Count(&physicalOff)
	stats["physical_total"] = physicalTotal
	stats["physical_running"] = physicalRunning
	stats["physical_off"] = physicalOff

	var cloudTotal, cloudRunning, cloudOff int64
	s.db.Model(&model.DcimCloud{}).Count(&cloudTotal)
	s.db.Model(&model.DcimCloud{}).Where("status = 1").Count(&cloudRunning)
	s.db.Model(&model.DcimCloud{}).Where("status = 0").Count(&cloudOff)
	stats["cloud_total"] = cloudTotal
	stats["cloud_running"] = cloudRunning
	stats["cloud_off"] = cloudOff

	var dcCount int64
	s.db.Model(&model.DcimDatacenter{}).Where("status = 1").Count(&dcCount)
	stats["datacenter_count"] = dcCount

	// 即将到期的服务器
	var expiringSoon int64
	s.db.Model(&model.DcimServer{}).Where("expired_at IS NOT NULL AND expired_at <= ? AND owner_id IS NOT NULL", time.Now().AddDate(0, 0, 7)).Count(&expiringSoon)
	s.db.Model(&model.DcimCloud{}).Where("expired_at IS NOT NULL AND expired_at <= ? AND owner_id IS NOT NULL", time.Now().AddDate(0, 0, 7)).Count(&expiringSoon)
	stats["expiring_within_7_days"] = expiringSoon

	return stats, nil
}
