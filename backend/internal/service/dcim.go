package service

import (
	"context"
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

	"gorm.io/gorm"
)

// DcimService DCIM业务逻辑
type DcimService struct {
	db         *gorm.DB
	log        *logger.Logger
	httpClient *http.Client
}

func NewDcimService(db *gorm.DB, log *logger.Logger) *DcimService {
	return &DcimService{
		db:  db,
		log: log,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ==================== 远程控制API ====================

// ipmiResponse IPMI API响应格式
type ipmiResponse struct {
	Result  string          `json:"result"`
	Message string          `json:"msg"`
	Data    json.RawMessage `json:"data"`
}

// bmsResponse BMS API响应格式
type bmsResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// dcimClientResponse DCIM Client API响应格式
type dcimClientResponse struct {
	Result  string          `json:"result"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// ipmiAction 调用IPMI HTTP API（物理服务器管理）
// POST到 {server_url}/index.php?m=api&a={action}
// 认证：username/password在POST数据中
func (s *DcimService) ipmiAction(server *model.DcimServer, action string, params map[string]interface{}) (map[string]interface{}, error) {
	if server.ControlURL == "" {
		return nil, fmt.Errorf("IPMI control URL not configured for server %d", server.ID)
	}

	apiURL := strings.TrimRight(server.ControlURL, "/") + "/index.php?m=api&a=" + url.QueryEscape(action)

	form := url.Values{}
	form.Set("username", server.ControlUser)
	form.Set("password", server.ControlPass)
	for k, v := range params {
		form.Set(k, fmt.Sprintf("%v", v))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("build IPMI request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("IPMI request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("read IPMI response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IPMI HTTP %d: %s", resp.StatusCode, string(body))
	}

	var apiResp ipmiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("parse IPMI response: %w", err)
	}

	if apiResp.Result != "success" && apiResp.Result != "1" {
		return nil, fmt.Errorf("IPMI error: %s", apiResp.Message)
	}

	result := map[string]interface{}{
		"result":  apiResp.Result,
		"message": apiResp.Message,
	}
	if apiResp.Data != nil {
		var data interface{}
		json.Unmarshal(apiResp.Data, &data)
		result["data"] = data
	}
	return result, nil
}

// bmsAction 调用BMS REST API（云裸金属服务器管理）
// POST到 {server_url}/bms/source/{server_id}/{action}
// Headers: access-user, access-token
func (s *DcimService) bmsAction(server *model.DcimServer, action string, params map[string]interface{}) (map[string]interface{}, error) {
	if server.ControlURL == "" {
		return nil, fmt.Errorf("BMS control URL not configured for server %d", server.ID)
	}

	// 从ControlExtra中提取BMS server_id，如果没有则使用服务器ID
	serverID := fmt.Sprintf("%d", server.ID)
	if server.ControlExtra != "" {
		var extra map[string]interface{}
		if err := json.Unmarshal([]byte(server.ControlExtra), &extra); err == nil {
			if sid, ok := extra["bms_server_id"]; ok {
				serverID = fmt.Sprintf("%v", sid)
			}
		}
	}

	apiURL := strings.TrimRight(server.ControlURL, "/") + "/bms/source/" + serverID + "/" + url.PathEscape(action)

	payload, _ := json.Marshal(params)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, strings.NewReader(string(payload)))
	if err != nil {
		return nil, fmt.Errorf("build BMS request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("access-user", server.ControlUser)
	req.Header.Set("access-token", server.ControlPass)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("BMS request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("read BMS response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("BMS HTTP %d: %s", resp.StatusCode, string(body))
	}

	var apiResp bmsResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("parse BMS response: %w", err)
	}

	if apiResp.Code != 0 && apiResp.Code != 200 {
		return nil, fmt.Errorf("BMS error (code=%d): %s", apiResp.Code, apiResp.Message)
	}

	result := map[string]interface{}{
		"code":    apiResp.Code,
		"message": apiResp.Message,
	}
	if apiResp.Data != nil {
		var data interface{}
		json.Unmarshal(apiResp.Data, &data)
		result["data"] = data
	}
	return result, nil
}

// dcimClientAction 调用DCIM Client API（上游DCIM管理）
// POST到 {dcim_client_url}/index.php?a=api&id={dcim_client_id}
// 参数：func, api_user, api_pass, id, 以及action特定参数
func (s *DcimService) dcimClientAction(server *model.DcimServer, funcName string, params map[string]interface{}) (map[string]interface{}, error) {
	if server.ControlURL == "" {
		return nil, fmt.Errorf("DCIM client URL not configured for server %d", server.ID)
	}

	// 从ControlExtra中提取dcim_client_id
	clientID := fmt.Sprintf("%d", server.ID)
	if server.ControlExtra != "" {
		var extra map[string]interface{}
		if err := json.Unmarshal([]byte(server.ControlExtra), &extra); err == nil {
			if cid, ok := extra["dcim_client_id"]; ok {
				clientID = fmt.Sprintf("%v", cid)
			}
		}
	}

	apiURL := strings.TrimRight(server.ControlURL, "/") + "/index.php?a=api&id=" + url.QueryEscape(clientID)

	form := url.Values{}
	form.Set("func", funcName)
	form.Set("api_user", server.ControlUser)
	form.Set("api_pass", server.ControlPass)
	form.Set("id", clientID)
	for k, v := range params {
		form.Set(k, fmt.Sprintf("%v", v))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("build DCIM client request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("DCIM client request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("read DCIM client response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DCIM client HTTP %d: %s", resp.StatusCode, string(body))
	}

	var apiResp dcimClientResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("parse DCIM client response: %w", err)
	}

	if apiResp.Result != "success" {
		return nil, fmt.Errorf("DCIM client error: %s", apiResp.Message)
	}

	result := map[string]interface{}{
		"result":  apiResp.Result,
		"message": apiResp.Message,
	}
	if apiResp.Data != nil {
		var data interface{}
		json.Unmarshal(apiResp.Data, &data)
		result["data"] = data
	}
	return result, nil
}

// zjmfCurl 调用zjmf_api（魔方系统API）
func (s *DcimService) zjmfCurl(server *model.DcimServer, action string, params map[string]interface{}) (map[string]interface{}, error) {
	if server.ControlURL == "" {
		return nil, fmt.Errorf("zjmf API URL not configured for server %d", server.ID)
	}

	apiURL := strings.TrimRight(server.ControlURL, "/") + "/api.php"

	form := url.Values{}
	form.Set("action", action)
	for k, v := range params {
		form.Set(k, fmt.Sprintf("%v", v))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("build zjmf request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("zjmf request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("read zjmf response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("zjmf HTTP %d: %s", resp.StatusCode, string(body))
	}

	var apiResp struct {
		Result string          `json:"result"`
		Msg    string          `json:"msg"`
		Data   json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("parse zjmf response: %w", err)
	}

	if apiResp.Result != "success" {
		return nil, fmt.Errorf("zjmf error: %s", apiResp.Msg)
	}

	result := map[string]interface{}{
		"result": apiResp.Result,
		"msg":    apiResp.Msg,
	}
	if apiResp.Data != nil {
		var data interface{}
		json.Unmarshal(apiResp.Data, &data)
		result["data"] = data
	}
	return result, nil
}

// executeServerAction 根据服务器控制方法执行远程操作
func (s *DcimService) executeServerAction(server *model.DcimServer, action string, params map[string]interface{}) (map[string]interface{}, error) {
	switch server.ControlMethod {
	case "ipmi":
		return s.ipmiAction(server, action, params)
	case "bms":
		return s.bmsAction(server, action, params)
	case "dcim_client":
		return s.dcimClientAction(server, action, params)
	case "zjmf_api":
		return s.zjmfCurl(server, action, params)
	case "local", "":
		return nil, nil // 本地模式不执行远程操作
	default:
		return nil, fmt.Errorf("unsupported control method: %s", server.ControlMethod)
	}
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
		Status:     1,
		Result:     "Server boot initiated",
	}
	s.db.Create(opLog)

	// 执行远程控制
	remoteResult, remoteErr := s.executeServerAction(server, "on", nil)

	now := time.Now()
	if remoteErr != nil {
		s.log.Errorf("remote boot failed for server %d: %v", serverID, remoteErr)
		s.db.Model(opLog).Updates(map[string]interface{}{
			"status":      3,
			"error_msg":   remoteErr.Error(),
			"finished_at": &now,
		})
		return fmt.Errorf("remote boot failed: %w", remoteErr)
	}

	// 更新服务器状态
	s.db.Model(server).Updates(map[string]interface{}{
		"status":       1,
		"power_status": 1,
	})

	resultMsg := "Server booted successfully"
	if remoteResult != nil {
		if msg, ok := remoteResult["message"].(string); ok && msg != "" {
			resultMsg = msg
		}
	}
	s.db.Model(opLog).Updates(map[string]interface{}{
		"status":      2,
		"result":      resultMsg,
		"finished_at": &now,
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

	action := "off"
	if force {
		action = "hard_off"
	}

	opLog := &model.DcimOperationLog{
		ServerType: "physical",
		ServerID:   serverID,
		OperatorID: operatorID,
		Action:     action,
		Status:     1,
		Result:     "Server shutdown initiated",
	}
	s.db.Create(opLog)

	remoteResult, remoteErr := s.executeServerAction(server, action, nil)

	now := time.Now()
	if remoteErr != nil {
		s.log.Errorf("remote shutdown failed for server %d: %v", serverID, remoteErr)
		s.db.Model(opLog).Updates(map[string]interface{}{
			"status":      3,
			"error_msg":   remoteErr.Error(),
			"finished_at": &now,
		})
		return fmt.Errorf("remote shutdown failed: %w", remoteErr)
	}

	s.db.Model(server).Updates(map[string]interface{}{
		"status":       0,
		"power_status": 0,
	})

	resultMsg := "Server shut down successfully"
	if remoteResult != nil {
		if msg, ok := remoteResult["message"].(string); ok && msg != "" {
			resultMsg = msg
		}
	}
	s.db.Model(opLog).Updates(map[string]interface{}{
		"status":      2,
		"result":      resultMsg,
		"finished_at": &now,
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
		Status:     1,
		Result:     "Server reboot initiated",
	}
	s.db.Create(opLog)

	remoteResult, remoteErr := s.executeServerAction(server, "reboot", nil)

	now := time.Now()
	if remoteErr != nil {
		s.log.Errorf("remote reboot failed for server %d: %v", serverID, remoteErr)
		s.db.Model(opLog).Updates(map[string]interface{}{
			"status":      3,
			"error_msg":   remoteErr.Error(),
			"finished_at": &now,
		})
		return fmt.Errorf("remote reboot failed: %w", remoteErr)
	}

	resultMsg := "Server rebooted successfully"
	if remoteResult != nil {
		if msg, ok := remoteResult["message"].(string); ok && msg != "" {
			resultMsg = msg
		}
	}
	s.db.Model(opLog).Updates(map[string]interface{}{
		"status":      2,
		"result":      resultMsg,
		"finished_at": &now,
	})

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

	// 执行远程重装
	params := map[string]interface{}{
		"os": os,
	}
	remoteResult, remoteErr := s.executeServerAction(server, "reinstall", params)

	now := time.Now()
	if remoteErr != nil {
		s.log.Errorf("remote reinstall failed for server %d: %v", serverID, remoteErr)
		s.db.Model(server).Updates(map[string]interface{}{
			"status": 2, // 故障状态
		})
		s.db.Model(opLog).Updates(map[string]interface{}{
			"status":      3,
			"error_msg":   remoteErr.Error(),
			"finished_at": &now,
		})
		return fmt.Errorf("remote reinstall failed: %w", remoteErr)
	}

	resultMsg := "OS reinstall initiated"
	if remoteResult != nil {
		if msg, ok := remoteResult["message"].(string); ok && msg != "" {
			resultMsg = msg
		}
	}

	// 异步等待重装完成（轮询状态）
	go func() {
		s.pollReinstallStatus(server, opLog, os)
	}()

	s.db.Model(opLog).Updates(map[string]interface{}{
		"result": resultMsg,
	})

	s.log.Infof("physical server reinstall started: id=%d os=%s operator=%d", serverID, os, operatorID)
	return nil
}

// pollReinstallStatus 轮询重装状态
func (s *DcimService) pollReinstallStatus(server *model.DcimServer, opLog *model.DcimOperationLog, os string) {
	maxAttempts := 60 // 最多轮询60次，每次间隔10秒
	for i := 0; i < maxAttempts; i++ {
		time.Sleep(10 * time.Second)

		// 查询远程状态
		result, err := s.executeServerAction(server, "sync", nil)
		if err != nil {
			s.log.Warnf("poll reinstall status failed for server %d: %v", server.ID, err)
			continue
		}

		if result != nil {
			if data, ok := result["data"].(map[string]interface{}); ok {
				if status, ok := data["status"].(string); ok {
					if status == "running" || status == "online" {
						// 重装完成
						finishTime := time.Now()
						s.db.Model(server).Updates(map[string]interface{}{
							"status":       1,
							"power_status": 1,
						})
						s.db.Model(opLog).Updates(map[string]interface{}{
							"status":      2,
							"result":      "OS reinstalled successfully",
							"finished_at": &finishTime,
						})
						s.log.Infof("physical server reinstall completed: id=%d os=%s", server.ID, os)
						return
					}
				}
			}
		}
	}

	// 超时，标记为失败
	finishTime := time.Now()
	s.db.Model(server).Updates(map[string]interface{}{
		"status": 2, // 故障状态
	})
	s.db.Model(opLog).Updates(map[string]interface{}{
		"status":      3,
		"error_msg":   "reinstall timeout",
		"finished_at": &finishTime,
	})
	s.log.Errorf("physical server reinstall timeout: id=%d", server.ID)
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

// ListFlowPackets returns paginated flow packets for a server.
func (s *DcimService) ListFlowPackets(serverID uint, page, pageSize int) ([]model.FlowPacket, int64, error) {
	var packets []model.FlowPacket
	var total int64

	q := s.db.Model(&model.FlowPacket{}).Where("server_id = ?", serverID)
	q.Count(&total)

	offset := (page - 1) * pageSize
	if err := q.Offset(offset).Limit(pageSize).Order("id ASC").Find(&packets).Error; err != nil {
		return nil, 0, err
	}
	return packets, total, nil
}

// AddFlowPacket creates a new flow packet.
func (s *DcimService) AddFlowPacket(serverID uint, name string, flow, price float64) (*model.FlowPacket, error) {
	packet := &model.FlowPacket{
		ServerID: serverID,
		Name:     name,
		Flow:     flow,
		Price:    price,
		Status:   1,
	}
	if err := s.db.Create(packet).Error; err != nil {
		return nil, err
	}
	return packet, nil
}

// EditFlowPacket updates a flow packet.
func (s *DcimService) EditFlowPacket(packetID uint, name string, flow, price *float64, status *int) error {
	updates := map[string]interface{}{}
	if name != "" {
		updates["name"] = name
	}
	if flow != nil {
		updates["flow"] = *flow
	}
	if price != nil {
		updates["price"] = *price
	}
	if status != nil {
		updates["status"] = *status
	}
	if len(updates) == 0 {
		return nil
	}
	return s.db.Model(&model.FlowPacket{}).Where("id = ?", packetID).Updates(updates).Error
}

// DeleteFlowPacket deletes a flow packet.
func (s *DcimService) DeleteFlowPacket(packetID uint) error {
	return s.db.Delete(&model.FlowPacket{}, packetID).Error
}

// AssignServer assigns a server to a user.
func (s *DcimService) AssignServer(serverID, userID, productID uint) error {
	updates := map[string]interface{}{
		"owner_id": userID,
	}
	if productID > 0 {
		updates["product_id"] = productID
	}
	return s.db.Model(&model.DcimServer{}).Where("id = ?", serverID).Updates(updates).Error
}

// GetSalesServers returns servers available for sale.
func (s *DcimService) GetSalesServers(page, pageSize int) ([]model.DcimServer, int64, error) {
	var servers []model.DcimServer
	var total int64

	q := s.db.Model(&model.DcimServer{}).Where("owner_id IS NULL")
	q.Count(&total)

	offset := (page - 1) * pageSize
	if err := q.Offset(offset).Limit(pageSize).Order("id ASC").Find(&servers).Error; err != nil {
		return nil, 0, err
	}
	return servers, total, nil
}

// GetNovncInfo returns NoVNC connection info.
func (s *DcimService) GetNovncInfo(serverID uint) (map[string]interface{}, error) {
	var server model.DcimServer
	if err := s.db.First(&server, serverID).Error; err != nil {
		return nil, err
	}

	info := map[string]interface{}{
		"host": server.IP,
		"port": 5900,
		"password": "",
	}
	return info, nil
}

// RefreshServerStatus refreshes server status from remote API.
func (s *DcimService) RefreshServerStatus(serverID uint) error {
	var server model.DcimServer
	if err := s.db.First(&server, serverID).Error; err != nil {
		return err
	}

	// Call remote API to get current status
	result, err := s.executeServerAction(&server, "status", nil)
	if err != nil {
		s.log.Warn("refresh server status failed: %v", err)
		return nil // Non-critical error
	}

	if statusStr, ok := result["status"].(string); ok {
		var status int8
		switch statusStr {
		case "running":
			status = 1
		case "stopped":
			status = 0
		default:
			status = -1
		}
		s.db.Model(&server).Update("status", status)
	}
	return nil
}

// ==================== Buy Records ====================

// DcimBuyRecord DCIM购买记录
type DcimBuyRecord struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	ProductID uint      `json:"product_id"`
	ServerID  uint      `json:"server_id"`
	Amount    float64   `json:"amount"`
	Status    string    `gorm:"size:32" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ListBuyRecords returns paginated buy records.
func (s *DcimService) ListBuyRecords(page, pageSize int, search, order, sort string) ([]DcimBuyRecord, int64, error) {
	var records []DcimBuyRecord
	var total int64

	q := s.db.Model(&DcimBuyRecord{})
	if search != "" {
		q = q.Where("id = ? OR user_id = ?", search, search)
	}

	q.Count(&total)

	offset := (page - 1) * pageSize
	if err := q.Offset(offset).Limit(pageSize).Order(fmt.Sprintf("%s %s", order, sort)).Find(&records).Error; err != nil {
		return nil, 0, err
	}
	return records, total, nil
}

// DeleteBuyRecord deletes a buy record by ID.
func (s *DcimService) DeleteBuyRecord(id uint) error {
	result := s.db.Delete(&DcimBuyRecord{}, id)
	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}
	return result.Error
}

// GetProductList returns all DCIM products.
func (s *DcimService) GetProductList() ([]map[string]interface{}, error) {
	var products []map[string]interface{}
	if err := s.db.Table("products").Select("id, name, price, billing_cycle").
		Where("type = ?", "dcim").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
