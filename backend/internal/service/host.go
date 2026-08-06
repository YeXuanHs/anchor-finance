package service

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/upstream"

	"gorm.io/datatypes"
	"gorm.io/gorm"
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
	Rack           string         `gorm:"size:32" json:"rack"`
	RackPosition   string         `gorm:"size:16" json:"rack_position"`
	Status         int            `gorm:"default:0;index" json:"status"` // 0=关机 1=运行中 2=故障 3=维护
	PowerStatus    int            `gorm:"default:0" json:"power_status"`
	OwnerID        *uint          `gorm:"index" json:"owner_id"`
	ProductID      *uint          `gorm:"index" json:"product_id"`
	OrderID        *uint          `gorm:"index" json:"order_id"`
	ExpiredAt      *time.Time     `gorm:"index" json:"expired_at"`
	ProvisionedAt  *time.Time     `json:"provisioned_at"`
	Remark         string         `gorm:"type:text" json:"remark"`
	AdminNotes     string         `gorm:"type:text" json:"admin_notes"`
	Tags           datatypes.JSON `gorm:"type:json" json:"tags"`
	Config         datatypes.JSON `gorm:"type:json" json:"config"`
	Metadata       datatypes.JSON `gorm:"type:json" json:"metadata"`
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
	db           *gorm.DB
	log          *logger.Logger
	upstreamSvc  *UpstreamService
}

func NewHostService(db *gorm.DB, log *logger.Logger, upstreamSvc *UpstreamService) *HostService {
	return &HostService{db: db, log: log, upstreamSvc: upstreamSvc}
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

		// Call upstream API if host has upstream provider configured
		upstreamErr := s.callUpstreamAction(host, req.Action, req.Params)
		if upstreamErr != nil {
			s.log.Errorf("upstream action failed for host %d: %v", hostID, upstreamErr)
			now := time.Now()
			tx.Model(operation).Updates(map[string]interface{}{
				"status":      3,
				"error_msg":   upstreamErr.Error(),
				"finished_at": &now,
			})
			return fmt.Errorf("upstream action failed: %w", upstreamErr)
		}

		// Update host status based on action
		newStatus := host.Status
		switch req.Action {
		case "boot":
			newStatus = 1
		case "shutdown":
			newStatus = 0
		case "reboot":
			newStatus = 1
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

// callUpstreamAction calls the upstream provider API to perform a host action.
// Returns nil if no upstream is configured (manually managed host).
func (s *HostService) callUpstreamAction(host *Host, action, params string) error {
	if len(host.Metadata) == 0 {
		return nil // manually managed host, no upstream
	}

	var meta map[string]interface{}
	if err := json.Unmarshal(host.Metadata, &meta); err != nil {
		return nil // invalid metadata
	}

	providerIDRaw, hasProvider := meta["upstream_provider_id"]
	productIDRaw, hasProduct := meta["upstream_product_id"]
	if !hasProvider || !hasProduct {
		return nil // no upstream configured
	}

	providerID := uint(providerIDRaw.(float64))
	remoteProductID := fmt.Sprintf("%v", productIDRaw)

	provider, err := s.upstreamSvc.GetProviderByID(providerID)
	if err != nil {
		return fmt.Errorf("upstream provider %d not found: %w", providerID, err)
	}

	client, err := upstream.NewClient(provider)
	if err != nil {
		return fmt.Errorf("create upstream client: %w", err)
	}

	// Use TestConnection as a health check before sending the action
	connResult, err := client.TestConnection()
	if err != nil {
		return fmt.Errorf("upstream connection test failed: %w", err)
	}
	if !connResult.OK {
		return fmt.Errorf("upstream not reachable: %s", connResult.Message)
	}

	// Dispatch the action to the upstream API based on provider type
	providerType := strings.ToLower(provider.Type)
	var actionErr error
	switch providerType {
	case "zjmf", "zjmfv3":
		actionErr = s.callZJMFHostAction(provider, host, remoteProductID, action, params)
	case "whmcs":
		actionErr = s.callWHMCSHostAction(provider, host, remoteProductID, action, params)
	case "v10":
		actionErr = s.callV10HostAction(provider, host, remoteProductID, action, params)
	default:
		actionErr = s.callCustomHostAction(provider, host, remoteProductID, action, params)
	}

	if actionErr != nil {
		return fmt.Errorf("upstream action dispatch failed: %w", actionErr)
	}

	s.log.Infof("upstream action dispatched: host=%d provider=%d remote=%s action=%s",
		host.ID, providerID, remoteProductID, action)
	return nil
}

// callZJMFHostAction dispatches a host action to a ZJMF-compatible panel.
func (s *HostService) callZJMFHostAction(provider *model.UpstreamProvider, host *Host, remoteID, action, params string) error {
	apiAction := mapHostActionToZJMF(action)
	apiURL := strings.TrimRight(provider.APIURL, "/") + "/api.php"

	form := url.Values{}
	form.Set("action", apiAction)
	form.Set("vps_id", remoteID)

	// ZJMF sign: md5(action + vps_id + api_key)
	sign := md5Sum(apiAction + remoteID + provider.APIKey)
	form.Set("sign", sign)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Post(apiURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("zjmf host action request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("zjmf host action http %d: %s", resp.StatusCode, string(body))
	}

	var apiResp struct {
		Result string `json:"result"`
		Msg    string `json:"msg"`
	}
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return fmt.Errorf("zjmf parse response: %w", err)
	}
	if apiResp.Result != "success" {
		return fmt.Errorf("zjmf host action error: %s", apiResp.Msg)
	}
	return nil
}

// callWHMCSHostAction dispatches a host action to a WHMCS-compatible API.
func (s *HostService) callWHMCSHostAction(provider *model.UpstreamProvider, host *Host, remoteID, action, params string) error {
	whmcsAction := mapHostActionToWHMCS(action)
	if whmcsAction == "" {
		return fmt.Errorf("whmcs does not support host action: %s", action)
	}

	cfg := map[string]interface{}{}
	if provider.Config != nil {
		cfg = provider.Config
	}
	identifier := provider.APIKey
	if id, ok := cfg["identifier"].(string); ok && id != "" {
		identifier = id
	}
	secret, _ := cfg["secret"].(string)

	form := url.Values{}
	form.Set("action", whmcsAction)
	form.Set("identifier", identifier)
	form.Set("secret", secret)
	form.Set("responsetype", "json")
	form.Set("serviceid", remoteID)

	apiURL := strings.TrimRight(provider.APIURL, "/") + "/includes/api.php"
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Post(apiURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("whmcs host action request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("whmcs host action http %d: %s", resp.StatusCode, string(body))
	}

	var apiResp struct {
		Result  string `json:"result"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return fmt.Errorf("whmcs parse response: %w", err)
	}
	if apiResp.Result != "success" {
		return fmt.Errorf("whmcs host action error: %s", apiResp.Message)
	}
	return nil
}

// callV10HostAction dispatches a host action to a V10 panel.
func (s *HostService) callV10HostAction(provider *model.UpstreamProvider, host *Host, remoteID, action, params string) error {
	apiPath := mapHostActionToV10(action)
	if apiPath == "" {
		return fmt.Errorf("v10 does not support host action: %s", action)
	}

	payload := map[string]interface{}{
		"server_id": remoteID,
	}
	if params != "" {
		var extra map[string]interface{}
		if json.Unmarshal([]byte(params), &extra) == nil {
			for k, v := range extra {
				payload[k] = v
			}
		}
	}
	bodyBytes, _ := json.Marshal(payload)

	token := provider.APIKey
	if t, ok := provider.Config["token"].(string); ok && t != "" {
		token = t
	}

	apiURL := strings.TrimRight(provider.APIURL, "/") + apiPath
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(string(bodyBytes)))
	if err != nil {
		return fmt.Errorf("v10 build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("v10 host action request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return fmt.Errorf("v10 host action http %d: %s", resp.StatusCode, string(respBody))
	}

	var apiResp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return fmt.Errorf("v10 parse response: %w", err)
	}
	if apiResp.Code != 0 && apiResp.Code != 200 {
		return fmt.Errorf("v10 host action error: %s", apiResp.Message)
	}
	return nil
}

// callCustomHostAction dispatches a host action to a custom upstream API.
func (s *HostService) callCustomHostAction(provider *model.UpstreamProvider, host *Host, remoteID, action, params string) error {
	payload := map[string]interface{}{
		"action":     action,
		"server_id":  remoteID,
		"host_id":    host.ID,
		"hostname":   host.Hostname,
		"ip":         host.IP,
	}
	if params != "" {
		var extra map[string]interface{}
		if json.Unmarshal([]byte(params), &extra) == nil {
			payload["params"] = extra
		} else {
			payload["params"] = params
		}
	}
	bodyBytes, _ := json.Marshal(payload)

	apiURL := strings.TrimRight(provider.APIURL, "/") + fmt.Sprintf("/api/host/%s", action)
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(string(bodyBytes)))
	if err != nil {
		return fmt.Errorf("custom build request: %w", err)
	}
	if provider.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+provider.APIKey)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("custom host action request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return fmt.Errorf("custom host action http %d: %s", resp.StatusCode, string(respBody))
	}

	var apiResp map[string]interface{}
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return fmt.Errorf("custom parse response: %w", err)
	}
	if success, ok := apiResp["success"].(bool); ok && !success {
		msg, _ := apiResp["message"].(string)
		return fmt.Errorf("custom host action error: %s", msg)
	}
	return nil
}

// mapHostActionToZJMF maps host actions to ZJMF API action names.
func mapHostActionToZJMF(action string) string {
	switch action {
	case "boot":
		return "start"
	case "shutdown":
		return "stop"
	case "reboot":
		return "restart"
	case "reinstall":
		return "reinstall"
	default:
		return action
	}
}

// mapHostActionToWHMCS maps host actions to WHMCS API action names.
func mapHostActionToWHMCS(action string) string {
	switch action {
	case "boot":
		return "ModuleStart"
	case "shutdown":
		return "ModuleStop"
	case "reboot":
		return "ModuleRestart"
	case "reinstall":
		return "ModuleReinstall"
	default:
		return ""
	}
}

// mapHostActionToV10 maps host actions to V10 API paths.
func mapHostActionToV10(action string) string {
	switch action {
	case "boot":
		return "/api/server/start"
	case "shutdown":
		return "/api/server/stop"
	case "reboot":
		return "/api/server/restart"
	case "reinstall":
		return "/api/server/reinstall"
	default:
		return ""
	}
}

// md5Sum computes the MD5 hex digest of a string.
func md5Sum(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// GetBillingInfo returns billing information for a host.
func (s *HostService) GetBillingInfo(hostID uint) (map[string]interface{}, error) {
	var host Host
	if err := s.db.First(&host, hostID).Error; err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"host_id":     host.ID,
		"hostname":    host.Hostname,
		"ip":          host.IP,
		"expired_at":  host.ExpiredAt,
		"status":      host.Status,
	}

	// Get related invoice if exists
	if host.OrderID != nil {
		var invoice struct {
			ID     uint    `json:"id"`
			Total  float64 `json:"total"`
			Status int     `json:"status"`
		}
		s.db.Table("invoices").
			Select("id, total, status").
			Where("order_id = ?", *host.OrderID).
			Order("id DESC").
			Limit(1).
			Find(&invoice)
		result["latest_invoice"] = invoice
	}

	return result, nil
}

// GetDownloadFiles returns download files for a host.
func (s *HostService) GetDownloadFiles(hostID uint) ([]map[string]interface{}, error) {
	var downloads []map[string]interface{}

	// Get downloads from the downloads table
	s.db.Table("downloads").
		Where("host_id = ? OR product_id = (SELECT product_id FROM hosts WHERE id = ?)", hostID, hostID).
		Find(&downloads)

	return downloads, nil
}

// UpdateRemark updates the remark for a host.
func (s *HostService) UpdateRemark(hostID uint, remark string) error {
	return s.db.Model(&Host{}).Where("id = ?", hostID).Update("remark", remark).Error
}
