package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type DcimCloudService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewDcimCloudService(db *gorm.DB, log *logger.Logger) *DcimCloudService {
	return &DcimCloudService{db: db, log: log}
}

// GetServerList returns paginated DCIM cloud servers.
func (s *DcimCloudService) GetServerList(page, pageSize int, keyword string, status *int8) ([]model.DcimCloudServer, int64, error) {
	var servers []model.DcimCloudServer
	var total int64

	query := s.db.Model(&model.DcimCloudServer{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR ip LIKE ? OR hostname LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&servers).Error; err != nil {
		return nil, 0, err
	}
	return servers, total, nil
}

// GetServerByID returns a single server by ID.
func (s *DcimCloudService) GetServerByID(id uint) (*model.DcimCloudServer, error) {
	var server model.DcimCloudServer
	if err := s.db.First(&server, id).Error; err != nil {
		return nil, err
	}
	return &server, nil
}

// CreateServer creates a new DCIM cloud server.
func (s *DcimCloudService) CreateServer(server *model.DcimCloudServer) error {
	return s.db.Create(server).Error
}

// UpdateServer updates a DCIM cloud server.
func (s *DcimCloudService) UpdateServer(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.DcimCloudServer{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteServer deletes a DCIM cloud server.
func (s *DcimCloudService) DeleteServer(id uint) error {
	return s.db.Delete(&model.DcimCloudServer{}, id).Error
}

// TestConnection tests connection to a DCIM cloud server.
func (s *DcimCloudService) TestConnection(id uint) error {
	server, err := s.GetServerByID(id)
	if err != nil {
		return fmt.Errorf("server not found: %w", err)
	}

	apiURL := fmt.Sprintf("https://%s/api/v1/system/status", server.IP)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}

	if server.Username != "" {
		req.SetBasicAuth(server.Username, server.Password)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("authentication failed: HTTP %d", resp.StatusCode)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("server returned HTTP %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
	s.log.Infof("DCIM cloud connection test succeeded: server=%s status=%d body=%s", server.IP, resp.StatusCode, string(body))
	return nil
}

// dcimServerInfo represents the remote server info from DCIM API.
type dcimServerInfo struct {
	Hostname      string `json:"hostname"`
	IP            string `json:"ip"`
	IPv6          string `json:"ipv6"`
	Status        string `json:"status"`
	OS            string `json:"os"`
	CPU           int    `json:"cpu"`
	MemoryMB      int    `json:"memory_mb"`
	DiskSizeGB    int    `json:"disk_size_gb"`
	BandwidthMbps int    `json:"bandwidth_mbps"`
	TrafficGB     int    `json:"traffic_gb"`
	VirtualType   string `json:"virtual_type"`
}

// SyncServer syncs server info from remote API.
func (s *DcimCloudService) SyncServer(id uint) error {
	server, err := s.GetServerByID(id)
	if err != nil {
		return fmt.Errorf("server not found: %w", err)
	}

	apiURL := fmt.Sprintf("https://%s/api/v1/servers/%s/info", server.IP, server.IP)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}

	if server.Username != "" {
		req.SetBasicAuth(server.Username, server.Password)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("sync request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("sync API returned HTTP %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 8192))
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	var remoteInfo dcimServerInfo
	if err := json.Unmarshal(body, &remoteInfo); err != nil {
		return fmt.Errorf("parse server info: %w", err)
	}

	updates := map[string]interface{}{
		"last_sync_at": time.Now(),
	}

	if remoteInfo.Hostname != "" {
		updates["hostname"] = remoteInfo.Hostname
	}
	if remoteInfo.IPv6 != "" {
		updates["ipv6"] = remoteInfo.IPv6
	}
	if remoteInfo.OS != "" {
		updates["os"] = remoteInfo.OS
	}
	if remoteInfo.CPU > 0 {
		updates["cpu"] = remoteInfo.CPU
	}
	if remoteInfo.MemoryMB > 0 {
		updates["memory_mb"] = remoteInfo.MemoryMB
	}
	if remoteInfo.DiskSizeGB > 0 {
		updates["disk_size_gb"] = remoteInfo.DiskSizeGB
	}
	if remoteInfo.BandwidthMbps > 0 {
		updates["bandwidth_mbps"] = remoteInfo.BandwidthMbps
	}
	if remoteInfo.TrafficGB > 0 {
		updates["traffic_gb"] = remoteInfo.TrafficGB
	}
	if remoteInfo.VirtualType != "" {
		updates["virtual_type"] = remoteInfo.VirtualType
	}

	switch remoteInfo.Status {
	case "running", "online":
		updates["status"] = int8(1)
	case "stopped", "offline":
		updates["status"] = int8(2)
	case "error":
		updates["status"] = int8(4)
	case "creating":
		updates["status"] = int8(3)
	}

	if err := s.db.Model(&model.DcimCloudServer{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return fmt.Errorf("update server: %w", err)
	}

	s.log.Infof("DCIM cloud server synced: id=%d ip=%s", id, server.IP)
	return nil
}

// GetOperationLogs returns operation logs for a server.
func (s *DcimCloudService) GetOperationLogs(serverID uint, page, pageSize int) ([]model.DcimCloudOperationLog, int64, error) {
	var logs []model.DcimCloudOperationLog
	var total int64

	query := s.db.Model(&model.DcimCloudOperationLog{})
	if serverID > 0 {
		query = query.Where("server_id = ?", serverID)
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
