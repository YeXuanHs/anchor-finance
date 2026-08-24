package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
)

// DcimCloudClient 魔方云API客户端
// 对接魔方云虚拟化系统，实现自动开小鸡
// 参考zjmf源码: app/common/logic/DcimCloud.php
type DcimCloudClient struct {
	ServerID    uint
	URL         string
	Username    string
	Password    string
	AccessToken string
	HTTPClient  *http.Client
	mu          sync.Mutex
}

// NewDcimCloudClient 从servers表创建魔方云客户端
// server_type='dcimcloud'
func NewDcimCloudClient(serverID uint) (*DcimCloudClient, error) {
	db := database.GetDB()
	var server model.Server
	if err := db.First(&server, serverID).Error; err != nil {
		return nil, fmt.Errorf("服务器不存在")
	}
	if server.ServerType != "dcimcloud" {
		return nil, fmt.Errorf("服务器类型不是dcimcloud")
	}

	scheme := "http"
	if server.Secure {
		scheme = "https"
	}
	baseURL := fmt.Sprintf("%s://%s", scheme, server.Hostname)
	if server.Port > 0 {
		baseURL += fmt.Sprintf(":%d", server.Port)
	}

	return &DcimCloudClient{
		ServerID:   serverID,
		URL:        baseURL,
		Username:   server.Username,
		Password:   server.Password, // 已加密存储
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}, nil
}

// Login 登录魔方云获取access-token
// zjmf: DcimCloud.php:3046 login()
func (c *DcimCloudClient) Login(force bool) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.AccessToken != "" && !force {
		return c.AccessToken, nil
	}

	// 调魔方云登录接口
	loginURL := c.URL + "/v1/login?a=a"
	data := url.Values{}
	data.Set("username", c.Username)
	data.Set("password", c.Password)

	resp, err := c.httpPost(loginURL, data.Encode(), nil)
	if err != nil {
		return "", fmt.Errorf("连接魔方云失败: %v", err)
	}

	if resp["status"] == "success" {
		token, _ := resp["origin"].(string)
		c.AccessToken = strings.Trim(token, "\"")
		return c.AccessToken, nil
	}

	return "", fmt.Errorf("魔方云登录失败: %v", resp["msg"])
}

// Curl 调用魔方云API（带自动重登录）
// zjmf: DcimCloud.php:3028 curl()
func (c *DcimCloudClient) Curl(action string, data map[string]interface{}, method string) (map[string]interface{}, error) {
	token, err := c.Login(false)
	if err != nil {
		return nil, err
	}

	headers := map[string]interface{}{
		"access-token": token,
	}

	apiURL := c.URL + "/v1" + action
	result, err := c.doRequest(apiURL, data, method, headers)
	if err != nil {
		return nil, err
	}

	// 401时自动重登录
	if httpCode, ok := result["http_code"].(int); ok && httpCode == 401 {
		token, err = c.Login(true)
		if err != nil {
			return nil, err
		}
		headers["access-token"] = token
		result, err = c.doRequest(apiURL, data, method, headers)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// ============ 核心操作（参考zjmf DcimCloud.php） ============

// On 开机（DcimCloud.php:1343）
func (c *DcimCloudClient) On(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/on", dcimID), nil, "POST")
}

// Off 关机（DcimCloud.php:1379）
func (c *DcimCloudClient) Off(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/off", dcimID), nil, "POST")
}

// Reboot 重启（DcimCloud.php:1415）
func (c *DcimCloudClient) Reboot(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/reboot", dcimID), nil, "POST")
}

// HardOff 强制关机（DcimCloud.php:1451）
func (c *DcimCloudClient) HardOff(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/hard_off", dcimID), nil, "POST")
}

// HardReboot 强制重启（DcimCloud.php:1487）
func (c *DcimCloudClient) HardReboot(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/hard_reboot", dcimID), nil, "POST")
}

// VNC 获取VNC连接（DcimCloud.php:1523）
func (c *DcimCloudClient) VNC(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/vnc", dcimID), nil, "POST")
}

// Reinstall 重装系统（DcimCloud.php:1585）
func (c *DcimCloudClient) Reinstall(dcimID uint, osName string, port int) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"os":   osName,
		"port": port,
	}
	return c.Curl(fmt.Sprintf("/clouds/%d/reinstall", dcimID), data, "POST")
}

// Rescue 救援模式（DcimCloud.php:2985）
func (c *DcimCloudClient) Rescue(dcimID uint, system string, tempPass string) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"type":      system,
		"temp_pass": tempPass,
	}
	return c.Curl(fmt.Sprintf("/clouds/%d/rescue", dcimID), data, "POST")
}

// ExitRescue 退出救援模式（DcimCloud.php:1268）
func (c *DcimCloudClient) ExitRescue(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/rescue", dcimID), nil, "DELETE")
}

// CrackPassword 重置密码（DcimCloud.php:2505）
func (c *DcimCloudClient) CrackPassword(dcimID uint, newPass string) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"password": newPass,
	}
	return c.Curl(fmt.Sprintf("/clouds/%d/password", dcimID), data, "POST")
}

// Status 获取状态（DcimCloud.php:1830）
func (c *DcimCloudClient) Status(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/status", dcimID), nil, "GET")
}

// Sync 同步信息（DcimCloud.php:1698）
func (c *DcimCloudClient) Sync(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/info", dcimID), nil, "GET")
}

// RemoteInfo 远程连接信息（DcimCloud.php:1309）
func (c *DcimCloudClient) RemoteInfo(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/remote_info", dcimID), nil, "GET")
}

// CreateSnap 创建快照（DcimCloud.php:318）
func (c *DcimCloudClient) CreateSnap(dcimID uint, name string) (map[string]interface{}, error) {
	data := map[string]interface{}{"name": name}
	return c.Curl(fmt.Sprintf("/clouds/%d/snapshot", dcimID), data, "POST")
}

// DeleteSnap 删除快照（DcimCloud.php:372）
func (c *DcimCloudClient) DeleteSnap(dcimID uint, snapID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/snapshot/%d", dcimID, snapID), nil, "DELETE")
}

// RestoreSnap 恢复快照（DcimCloud.php:419）
func (c *DcimCloudClient) RestoreSnap(dcimID uint, snapID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/snapshot/%d/restore", dcimID, snapID), nil, "POST")
}

// ListSnapBackup 快照/备份列表（DcimCloud.php:288）
func (c *DcimCloudClient) ListSnapBackup(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/snapshot", dcimID), nil, "GET")
}

// AddNatAcl 添加NAT ACL（DcimCloud.php:1088）
func (c *DcimCloudClient) AddNatAcl(dcimID uint, data map[string]interface{}) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/nat_acl", dcimID), data, "POST")
}

// DelNatAcl 删除NAT ACL（DcimCloud.php:1142）
func (c *DcimCloudClient) DelNatAcl(dcimID uint, aclID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/nat_acl/%d", dcimID, aclID), nil, "DELETE")
}

// AddNatWeb 添加NAT Web（DcimCloud.php:1184）
func (c *DcimCloudClient) AddNatWeb(dcimID uint, data map[string]interface{}) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/nat_web", dcimID), data, "POST")
}

// DelNatWeb 删除NAT Web（DcimCloud.php:1226）
func (c *DcimCloudClient) DelNatWeb(dcimID uint, webID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/nat_web/%d", dcimID, webID), nil, "DELETE")
}

// GetTrafficUsage 获取流量使用（DcimCloud.php:2846）
func (c *DcimCloudClient) GetTrafficUsage(dcimID uint, start, end string) (map[string]interface{}, error) {
	data := map[string]interface{}{"start": start, "end": end}
	return c.Curl(fmt.Sprintf("/clouds/%d/traffic", dcimID), data, "GET")
}

// GetChartData 获取图表数据（DcimCloud.php:2357）
func (c *DcimCloudClient) GetChartData(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/chart", dcimID), nil, "GET")
}

// BuyFlowPacket 购买流量包（DcimCloud.php:1771）
func (c *DcimCloudClient) BuyFlowPacket(dcimID uint, packetID uint) (map[string]interface{}, error) {
	data := map[string]interface{}{"packet_id": packetID}
	return c.Curl(fmt.Sprintf("/clouds/%d/flow_packet", dcimID), data, "POST")
}

// CreateSecurityGroup 创建安全组（DcimCloud.php:755）
func (c *DcimCloudClient) CreateSecurityGroup(dcimID uint, name string) (map[string]interface{}, error) {
	data := map[string]interface{}{"name": name}
	return c.Curl(fmt.Sprintf("/clouds/%d/security_group", dcimID), data, "POST")
}

// CreateSecurityRule 创建安全规则（DcimCloud.php:849）
func (c *DcimCloudClient) CreateSecurityRule(dcimID uint, groupID uint, data map[string]interface{}) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/security_group/%d/rule", dcimID, groupID), data, "POST")
}

// ============ HTTP工具方法 ============

func (c *DcimCloudClient) doRequest(apiURL string, data map[string]interface{}, method string, headers map[string]interface{}) (map[string]interface{}, error) {
	var req *http.Request
	var err error

	if method == "GET" && len(data) > 0 {
		params := url.Values{}
		for k, v := range data {
			params.Set(k, fmt.Sprintf("%v", v))
		}
		apiURL += "?" + params.Encode()
		req, err = http.NewRequest("GET", apiURL, nil)
	} else {
		formData := url.Values{}
		for k, v := range data {
			formData.Set(k, fmt.Sprintf("%v", v))
		}
		req, err = http.NewRequest(method, apiURL, strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if result == nil {
		result = map[string]interface{}{}
	}
	result["http_code"] = resp.StatusCode

	return result, nil
}

func (c *DcimCloudClient) httpPost(postURL string, body string, headers map[string]interface{}) (map[string]interface{}, error) {
	req, err := http.NewRequest("POST", postURL, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(respBody, &result)
	if result == nil {
		result = map[string]interface{}{}
	}
	result["http_code"] = resp.StatusCode
	return result, nil
}
