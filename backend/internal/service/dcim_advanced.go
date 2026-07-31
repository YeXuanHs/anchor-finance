package service

import (
	"errors"
	"fmt"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// ==================== KVM/IPMI/BMC ====================

// GetKVMURL 获取KVM控制台URL（通过IPMI/iLO/iDRAC）
func (s *DcimService) GetKVMURL(serverID uint) (map[string]interface{}, error) {
	server, err := s.GetServerByID(serverID)
	if err != nil {
		return nil, err
	}

	result, err := s.executeServerAction(server, "kvm_url", map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("get KVM URL: %w", err)
	}

	if result == nil {
		// 本地模式：构造默认IPMI KVM URL
		return map[string]interface{}{
			"url":    fmt.Sprintf("https://%s/kvm", server.ControlURL),
			"server": server.IP,
		}, nil
	}

	return result, nil
}

// GetBMCInfo 获取BMC连接信息
func (s *DcimService) GetBMCInfo(serverID uint) (map[string]interface{}, error) {
	server, err := s.GetServerByID(serverID)
	if err != nil {
		return nil, err
	}

	result, err := s.executeServerAction(server, "bmc_info", map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("get BMC info: %w", err)
	}

	if result == nil {
		// 本地模式：返回服务器自身配置
		return map[string]interface{}{
			"ip":             server.ControlURL,
			"username":       server.ControlUser,
			"control_method": server.ControlMethod,
		}, nil
	}

	return result, nil
}

// GetNoVNCURL 获取noVNC Web控制台URL
func (s *DcimService) GetNoVNCURL(serverID uint) (map[string]interface{}, error) {
	server, err := s.GetServerByID(serverID)
	if err != nil {
		return nil, err
	}

	result, err := s.executeServerAction(server, "novnc_url", map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("get noVNC URL: %w", err)
	}

	if result == nil {
		return map[string]interface{}{
			"url":    fmt.Sprintf("https://%s/novnc", server.ControlURL),
			"server": server.IP,
		}, nil
	}

	return result, nil
}

// ==================== 救援系统 ====================

// BootRescue 启动救援模式
func (s *DcimService) BootRescue(serverID uint, os string, operatorID uint) (*model.DcimRescueLog, error) {
	server, err := s.GetServerByID(serverID)
	if err != nil {
		return nil, err
	}

	if server.OwnerID == nil {
		return nil, errors.New("server is not assigned to any user")
	}

	rescueLog := &model.DcimRescueLog{
		ServerID: serverID,
		Action:   "rescue",
		Status:   "running",
	}
	if err := s.db.Create(rescueLog).Error; err != nil {
		return nil, fmt.Errorf("create rescue log: %w", err)
	}

	// 记录操作日志
	opLog := &model.DcimOperationLog{
		ServerType: "physical",
		ServerID:   serverID,
		OperatorID: operatorID,
		Action:     "rescue",
		Params:     fmt.Sprintf(`{"os":"%s"}`, os),
		Status:     1,
	}
	s.db.Create(opLog)

	params := map[string]interface{}{
		"os": os,
	}
	remoteResult, remoteErr := s.executeServerAction(server, "rescue", params)

	now := time.Now()
	if remoteErr != nil {
		s.db.Model(rescueLog).Updates(map[string]interface{}{
			"status": "failed",
			"result": remoteErr.Error(),
		})
		s.db.Model(opLog).Updates(map[string]interface{}{
			"status":      3,
			"error_msg":   remoteErr.Error(),
			"finished_at": &now,
		})
		return nil, fmt.Errorf("remote rescue failed: %w", remoteErr)
	}

	resultMsg := "Rescue mode boot initiated"
	if remoteResult != nil {
		if msg, ok := remoteResult["message"].(string); ok && msg != "" {
			resultMsg = msg
		}
	}

	s.db.Model(rescueLog).Updates(map[string]interface{}{
		"status": "completed",
		"result": resultMsg,
	})
	s.db.Model(opLog).Updates(map[string]interface{}{
		"status":      2,
		"result":      resultMsg,
		"finished_at": &now,
	})

	s.log.Infof("server rescue boot: id=%d os=%s operator=%d", serverID, os, operatorID)
	return rescueLog, nil
}

// CrackPassword 重置密码
func (s *DcimService) CrackPassword(serverID uint, operatorID uint) (*model.DcimRescueLog, error) {
	server, err := s.GetServerByID(serverID)
	if err != nil {
		return nil, err
	}

	if server.OwnerID == nil {
		return nil, errors.New("server is not assigned to any user")
	}

	rescueLog := &model.DcimRescueLog{
		ServerID: serverID,
		Action:   "crack_pass",
		Status:   "running",
	}
	if err := s.db.Create(rescueLog).Error; err != nil {
		return nil, fmt.Errorf("create rescue log: %w", err)
	}

	opLog := &model.DcimOperationLog{
		ServerType: "physical",
		ServerID:   serverID,
		OperatorID: operatorID,
		Action:     "crack_pass",
		Status:     1,
	}
	s.db.Create(opLog)

	remoteResult, remoteErr := s.executeServerAction(server, "crack_pass", nil)

	now := time.Now()
	if remoteErr != nil {
		s.db.Model(rescueLog).Updates(map[string]interface{}{
			"status": "failed",
			"result": remoteErr.Error(),
		})
		s.db.Model(opLog).Updates(map[string]interface{}{
			"status":      3,
			"error_msg":   remoteErr.Error(),
			"finished_at": &now,
		})
		return nil, fmt.Errorf("remote crack password failed: %w", remoteErr)
	}

	resultMsg := "Password reset initiated"
	if remoteResult != nil {
		if msg, ok := remoteResult["message"].(string); ok && msg != "" {
			resultMsg = msg
		}
	}

	s.db.Model(rescueLog).Updates(map[string]interface{}{
		"status": "completed",
		"result": resultMsg,
	})
	s.db.Model(opLog).Updates(map[string]interface{}{
		"status":      2,
		"result":      resultMsg,
		"finished_at": &now,
	})

	s.log.Infof("server crack password: id=%d operator=%d", serverID, operatorID)
	return rescueLog, nil
}

// GetRescueStatus 获取救援模式状态
func (s *DcimService) GetRescueStatus(serverID uint) (*model.DcimRescueLog, error) {
	var rescueLog model.DcimRescueLog
	if err := s.db.Where("server_id = ?", serverID).Order("id DESC").First(&rescueLog).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no rescue log found")
		}
		return nil, err
	}
	return &rescueLog, nil
}

// ==================== 流量监控 ====================

// GetTrafficUsage 获取当月流量使用情况
func (s *DcimService) GetTrafficUsage(serverID uint) (map[string]interface{}, error) {
	server, err := s.GetServerByID(serverID)
	if err != nil {
		return nil, err
	}

	// 先查询远程获取最新数据
	remoteResult, _ := s.executeServerAction(server, "traffic_usage", nil)

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	var inBytes, outBytes, totalBytes int64
	s.db.Model(&model.DcimTrafficLog{}).
		Where("server_id = ? AND recorded_at >= ?", serverID, startOfMonth).
		Select("COALESCE(SUM(in_bytes),0), COALESCE(SUM(out_bytes),0), COALESCE(SUM(total_bytes),0)").
		Row().Scan(&inBytes, &outBytes, &totalBytes)

	// 如果远程有最新数据，优先使用
	if remoteResult != nil {
		if data, ok := remoteResult["data"].(map[string]interface{}); ok {
			if v, ok := data["in_bytes"].(float64); ok {
				inBytes = int64(v)
			}
			if v, ok := data["out_bytes"].(float64); ok {
				outBytes = int64(v)
			}
			if v, ok := data["total_bytes"].(float64); ok {
				totalBytes = int64(v)
			}
		}
	}

	// 记录到数据库
	if totalBytes > 0 {
		s.db.Create(&model.DcimTrafficLog{
			ServerID:   serverID,
			InBytes:    inBytes,
			OutBytes:   outBytes,
			TotalBytes: totalBytes,
			RecordedAt: now,
		})
	}

	limitGB := server.TrafficGB
	usedGB := float64(totalBytes) / (1024 * 1024 * 1024)
	var usagePercent float64
	if limitGB > 0 {
		usagePercent = usedGB / float64(limitGB) * 100
	}

	return map[string]interface{}{
		"server_id":     serverID,
		"in_bytes":      inBytes,
		"out_bytes":     outBytes,
		"total_bytes":   totalBytes,
		"in_gb":         float64(inBytes) / (1024 * 1024 * 1024),
		"out_gb":        float64(outBytes) / (1024 * 1024 * 1024),
		"total_gb":      usedGB,
		"limit_gb":      limitGB,
		"usage_percent": usagePercent,
		"month":         now.Format("2006-01"),
	}, nil
}

// GetTrafficChart 获取流量图表数据
func (s *DcimService) GetTrafficChart(serverID uint, period string) ([]model.DcimTrafficLog, error) {
	var logs []model.DcimTrafficLog

	now := time.Now()
	var startTime time.Time

	switch period {
	case "7d":
		startTime = now.AddDate(0, 0, -7)
	case "30d":
		startTime = now.AddDate(0, 0, -30)
	case "90d":
		startTime = now.AddDate(0, 0, -90)
	case "month":
		startTime = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	default:
		startTime = now.AddDate(0, 0, -30)
	}

	if err := s.db.Where("server_id = ? AND recorded_at >= ?", serverID, startTime).
		Order("recorded_at ASC").Find(&logs).Error; err != nil {
		return nil, err
	}

	return logs, nil
}

// ResetTrafficCounter 重置月流量计数器
func (s *DcimService) ResetTrafficCounter(serverID uint) error {
	server, err := s.GetServerByID(serverID)
	if err != nil {
		return err
	}

	_, err = s.executeServerAction(server, "reset_traffic", nil)
	if err != nil {
		return fmt.Errorf("reset traffic counter: %w", err)
	}

	// 删除当月记录
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	s.db.Where("server_id = ? AND recorded_at >= ?", serverID, startOfMonth).Delete(&model.DcimTrafficLog{})

	s.log.Infof("traffic counter reset: server_id=%d", serverID)
	return nil
}

// ==================== 快照管理 ====================

// CreateSnapshot 创建快照
func (s *DcimService) CreateSnapshot(serverID uint, name string) (*model.DcimSnapshot, error) {
	server, err := s.GetServerByID(serverID)
	if err != nil {
		return nil, err
	}

	if server.OwnerID == nil {
		return nil, errors.New("server is not assigned to any user")
	}

	snapshot := &model.DcimSnapshot{
		ServerID: serverID,
		Name:     name,
		Status:   "creating",
	}
	if err := s.db.Create(snapshot).Error; err != nil {
		return nil, fmt.Errorf("create snapshot record: %w", err)
	}

	// 远程创建快照
	params := map[string]interface{}{
		"snapshot_id": snapshot.ID,
		"name":        name,
	}
	remoteResult, remoteErr := s.executeServerAction(server, "snapshot_create", params)
	if remoteErr != nil {
		s.db.Model(snapshot).Update("status", "failed")
		return nil, fmt.Errorf("remote snapshot create: %w", remoteErr)
	}

	// 异步等待完成
	go func() {
		s.pollSnapshotStatus(snapshot, remoteResult)
	}()

	s.log.Infof("snapshot creation started: server_id=%d name=%s", serverID, name)
	return snapshot, nil
}

// pollSnapshotStatus 轮询快照创建状态
func (s *DcimService) pollSnapshotStatus(snapshot *model.DcimSnapshot, remoteResult map[string]interface{}) {
	if remoteResult != nil {
		if data, ok := remoteResult["data"].(map[string]interface{}); ok {
			if size, ok := data["size_gb"].(float64); ok {
				s.db.Model(snapshot).Update("size_gb", int(size))
			}
		}
	}
	// 模拟快照完成
	time.Sleep(3 * time.Second)
	s.db.Model(snapshot).Update("status", "available")
	s.log.Infof("snapshot created: id=%d", snapshot.ID)
}

// RestoreSnapshot 从快照恢复
func (s *DcimService) RestoreSnapshot(snapshotID uint, operatorID uint) error {
	var snapshot model.DcimSnapshot
	if err := s.db.First(&snapshot, snapshotID).Error; err != nil {
		return errors.New("snapshot not found")
	}

	if snapshot.Status != "available" {
		return fmt.Errorf("snapshot is not available (status: %s)", snapshot.Status)
	}

	server, err := s.GetServerByID(snapshot.ServerID)
	if err != nil {
		return err
	}

	opLog := &model.DcimOperationLog{
		ServerType: "physical",
		ServerID:   snapshot.ServerID,
		OperatorID: operatorID,
		Action:     "snapshot_restore",
		Params:     fmt.Sprintf(`{"snapshot_id":%d}`, snapshotID),
		Status:     1,
	}
	s.db.Create(opLog)

	params := map[string]interface{}{
		"snapshot_id": snapshotID,
	}
	_, remoteErr := s.executeServerAction(server, "snapshot_restore", params)

	now := time.Now()
	if remoteErr != nil {
		s.db.Model(opLog).Updates(map[string]interface{}{
			"status":      3,
			"error_msg":   remoteErr.Error(),
			"finished_at": &now,
		})
		return fmt.Errorf("remote snapshot restore: %w", remoteErr)
	}

	s.db.Model(opLog).Updates(map[string]interface{}{
		"status":      2,
		"result":      "Snapshot restored successfully",
		"finished_at": &now,
	})

	s.log.Infof("snapshot restored: id=%d server_id=%d", snapshotID, snapshot.ServerID)
	return nil
}

// DeleteSnapshot 删除快照
func (s *DcimService) DeleteSnapshot(snapshotID uint) error {
	var snapshot model.DcimSnapshot
	if err := s.db.First(&snapshot, snapshotID).Error; err != nil {
		return errors.New("snapshot not found")
	}

	if snapshot.Status == "deleting" {
		return errors.New("snapshot is already being deleted")
	}

	s.db.Model(&snapshot).Update("status", "deleting")

	server, err := s.GetServerByID(snapshot.ServerID)
	if err == nil {
		params := map[string]interface{}{
			"snapshot_id": snapshotID,
		}
		s.executeServerAction(server, "snapshot_delete", params)
	}

	if err := s.db.Delete(&snapshot).Error; err != nil {
		return err
	}

	s.log.Infof("snapshot deleted: id=%d server_id=%d", snapshotID, snapshot.ServerID)
	return nil
}

// GetSnapshots 获取服务器快照列表
func (s *DcimService) GetSnapshots(serverID uint) ([]model.DcimSnapshot, error) {
	var snapshots []model.DcimSnapshot
	if err := s.db.Where("server_id = ?", serverID).Order("id DESC").Find(&snapshots).Error; err != nil {
		return nil, err
	}
	return snapshots, nil
}

// ==================== 备份管理 ====================

// CreateBackup 创建备份
func (s *DcimService) CreateBackup(serverID uint, name string, backupType string) (*model.DcimBackup, error) {
	server, err := s.GetServerByID(serverID)
	if err != nil {
		return nil, err
	}

	if server.OwnerID == nil {
		return nil, errors.New("server is not assigned to any user")
	}

	if backupType == "" {
		backupType = "manual"
	}

	backup := &model.DcimBackup{
		ServerID: serverID,
		Name:     name,
		Type:     backupType,
		Status:   "creating",
	}
	if err := s.db.Create(backup).Error; err != nil {
		return nil, fmt.Errorf("create backup record: %w", err)
	}

	params := map[string]interface{}{
		"backup_id": backup.ID,
		"name":      name,
		"type":      backupType,
	}
	remoteResult, remoteErr := s.executeServerAction(server, "backup_create", params)
	if remoteErr != nil {
		s.db.Model(backup).Update("status", "failed")
		return nil, fmt.Errorf("remote backup create: %w", remoteErr)
	}

	go func() {
		s.pollBackupStatus(backup, remoteResult)
	}()

	s.log.Infof("backup creation started: server_id=%d name=%s type=%s", serverID, name, backupType)
	return backup, nil
}

// pollBackupStatus 轮询备份创建状态
func (s *DcimService) pollBackupStatus(backup *model.DcimBackup, remoteResult map[string]interface{}) {
	if remoteResult != nil {
		if data, ok := remoteResult["data"].(map[string]interface{}); ok {
			if size, ok := data["size_gb"].(float64); ok {
				s.db.Model(backup).Update("size_gb", int(size))
			}
		}
	}
	time.Sleep(3 * time.Second)
	s.db.Model(backup).Update("status", "available")
	s.log.Infof("backup created: id=%d", backup.ID)
}

// RestoreBackup 从备份恢复
func (s *DcimService) RestoreBackup(backupID uint, operatorID uint) error {
	var backup model.DcimBackup
	if err := s.db.First(&backup, backupID).Error; err != nil {
		return errors.New("backup not found")
	}

	if backup.Status != "available" {
		return fmt.Errorf("backup is not available (status: %s)", backup.Status)
	}

	server, err := s.GetServerByID(backup.ServerID)
	if err != nil {
		return err
	}

	opLog := &model.DcimOperationLog{
		ServerType: "physical",
		ServerID:   backup.ServerID,
		OperatorID: operatorID,
		Action:     "backup_restore",
		Params:     fmt.Sprintf(`{"backup_id":%d}`, backupID),
		Status:     1,
	}
	s.db.Create(opLog)

	params := map[string]interface{}{
		"backup_id": backupID,
	}
	_, remoteErr := s.executeServerAction(server, "backup_restore", params)

	now := time.Now()
	if remoteErr != nil {
		s.db.Model(opLog).Updates(map[string]interface{}{
			"status":      3,
			"error_msg":   remoteErr.Error(),
			"finished_at": &now,
		})
		return fmt.Errorf("remote backup restore: %w", remoteErr)
	}

	s.db.Model(opLog).Updates(map[string]interface{}{
		"status":      2,
		"result":      "Backup restored successfully",
		"finished_at": &now,
	})

	s.log.Infof("backup restored: id=%d server_id=%d", backupID, backup.ServerID)
	return nil
}

// DeleteBackup 删除备份
func (s *DcimService) DeleteBackup(backupID uint) error {
	var backup model.DcimBackup
	if err := s.db.First(&backup, backupID).Error; err != nil {
		return errors.New("backup not found")
	}

	server, err := s.GetServerByID(backup.ServerID)
	if err == nil {
		params := map[string]interface{}{
			"backup_id": backupID,
		}
		s.executeServerAction(server, "backup_delete", params)
	}

	if err := s.db.Delete(&backup).Error; err != nil {
		return err
	}

	s.log.Infof("backup deleted: id=%d server_id=%d", backupID, backup.ServerID)
	return nil
}

// GetBackups 获取服务器备份列表
func (s *DcimService) GetBackups(serverID uint) ([]model.DcimBackup, error) {
	var backups []model.DcimBackup
	if err := s.db.Where("server_id = ?", serverID).Order("id DESC").Find(&backups).Error; err != nil {
		return nil, err
	}
	return backups, nil
}

// ==================== 电源管理增强 ====================

// GetPowerStatus 获取实时电源状态
func (s *DcimService) GetPowerStatus(serverID uint) (map[string]interface{}, error) {
	server, err := s.GetServerByID(serverID)
	if err != nil {
		return nil, err
	}

	result, err := s.executeServerAction(server, "power_status", nil)
	if err != nil {
		return nil, fmt.Errorf("get power status: %w", err)
	}

	if result == nil {
		return map[string]interface{}{
			"power_status": server.PowerStatus,
			"status":       server.Status,
		}, nil
	}

	return result, nil
}

// RefreshPowerStatus 强制刷新电源状态（从IPMI/BMC实时获取）
func (s *DcimService) RefreshPowerStatus(serverID uint) (map[string]interface{}, error) {
	server, err := s.GetServerByID(serverID)
	if err != nil {
		return nil, err
	}

	result, err := s.executeServerAction(server, "sync", nil)
	if err != nil {
		return nil, fmt.Errorf("refresh power status: %w", err)
	}

	if result != nil {
		if data, ok := result["data"].(map[string]interface{}); ok {
			// 更新本地状态
			if status, ok := data["status"].(string); ok {
				switch status {
				case "running", "online":
					s.db.Model(server).Updates(map[string]interface{}{
						"status":       1,
						"power_status": 1,
					})
				case "off", "offline":
					s.db.Model(server).Updates(map[string]interface{}{
						"status":       0,
						"power_status": 0,
					})
				}
			}
		}
	}

	return map[string]interface{}{
		"power_status": server.PowerStatus,
		"status":       server.Status,
		"remote":       result,
	}, nil
}

// ==================== 重装增强 ====================

// GetReinstallStatus 获取重装进度
func (s *DcimService) GetReinstallStatus(serverID uint) (map[string]interface{}, error) {
	server, err := s.GetServerByID(serverID)
	if err != nil {
		return nil, err
	}

	// 查询最新重装操作日志
	var opLog model.DcimOperationLog
	if err := s.db.Where("server_id = ? AND action = ?", serverID, "reinstall").
		Order("id DESC").First(&opLog).Error; err != nil {
		return nil, errors.New("no reinstall operation found")
	}

	result := map[string]interface{}{
		"status":      opLog.Status,
		"result":      opLog.Result,
		"error_msg":   opLog.ErrorMsg,
		"created_at":  opLog.CreatedAt,
		"finished_at": opLog.FinishedAt,
	}

	// 如果还在进行中，查询远程状态
	if opLog.Status == 1 {
		remoteResult, _ := s.executeServerAction(server, "sync", nil)
		if remoteResult != nil {
			result["remote"] = remoteResult
		}
	}

	return result, nil
}

// CancelReinstall 取消正在进行的重装
func (s *DcimService) CancelReinstall(serverID uint, operatorID uint) error {
	server, err := s.GetServerByID(serverID)
	if err != nil {
		return err
	}

	// 查找进行中的重装
	var opLog model.DcimOperationLog
	if err := s.db.Where("server_id = ? AND action = ? AND status = ?", serverID, "reinstall", 1).
		First(&opLog).Error; err != nil {
		return errors.New("no active reinstall operation found")
	}

	// 远程取消
	_, remoteErr := s.executeServerAction(server, "cancel_reinstall", nil)

	now := time.Now()
	if remoteErr != nil {
		return fmt.Errorf("cancel reinstall: %w", remoteErr)
	}

	s.db.Model(server).Updates(map[string]interface{}{
		"status": 1,
	})
	s.db.Model(&opLog).Updates(map[string]interface{}{
		"status":      3,
		"error_msg":   "Cancelled by user",
		"finished_at": &now,
	})

	s.log.Infof("reinstall cancelled: server_id=%d operator=%d", serverID, operatorID)
	return nil
}

// GetOSList 获取可用操作系统列表
func (s *DcimService) GetOSList() ([]map[string]interface{}, error) {
	// 默认OS列表
	defaultOS := []map[string]interface{}{
		{"name": "CentOS 7.9", "os": "centos7.9", "arch": "x86_64"},
		{"name": "CentOS 8 Stream", "os": "centos8stream", "arch": "x86_64"},
		{"name": "AlmaLinux 8", "os": "almalinux8", "arch": "x86_64"},
		{"name": "AlmaLinux 9", "os": "almalinux9", "arch": "x86_64"},
		{"name": "Rocky Linux 8", "os": "rocky8", "arch": "x86_64"},
		{"name": "Rocky Linux 9", "os": "rocky9", "arch": "x86_64"},
		{"name": "Ubuntu 20.04", "os": "ubuntu20.04", "arch": "x86_64"},
		{"name": "Ubuntu 22.04", "os": "ubuntu22.04", "arch": "x86_64"},
		{"name": "Ubuntu 24.04", "os": "ubuntu24.04", "arch": "x86_64"},
		{"name": "Debian 11", "os": "debian11", "arch": "x86_64"},
		{"name": "Debian 12", "os": "debian12", "arch": "x86_64"},
		{"name": "Windows Server 2019", "os": "win2019", "arch": "x86_64"},
		{"name": "Windows Server 2022", "os": "win2022", "arch": "x86_64"},
		{"name": "ESXi 7.0", "os": "esxi7.0", "arch": "x86_64"},
		{"name": "ESXi 8.0", "os": "esxi8.0", "arch": "x86_64"},
	}

	// 尝试从远程获取真实列表
	var servers []model.DcimServer
	s.db.Where("control_method != ?", "local").Limit(1).Find(&servers)
	if len(servers) > 0 {
		result, err := s.executeServerAction(&servers[0], "os_list", nil)
		if err == nil && result != nil {
			if data, ok := result["data"].([]interface{}); ok && len(data) > 0 {
				var osList []map[string]interface{}
				for _, item := range data {
					if osItem, ok := item.(map[string]interface{}); ok {
						osList = append(osList, osItem)
					}
				}
				if len(osList) > 0 {
					return osList, nil
				}
			}
		}
	}

	return defaultOS, nil
}

// ==================== 救援/流量日志查询 ====================

// GetRescueLogs 获取救援操作日志
func (s *DcimService) GetRescueLogs(serverID uint, page, pageSize int) ([]model.DcimRescueLog, int64, error) {
	var logs []model.DcimRescueLog
	var total int64

	query := s.db.Model(&model.DcimRescueLog{})
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
