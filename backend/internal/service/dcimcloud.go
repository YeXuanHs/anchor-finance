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
	"github.com/YeXuanHs/anchor-finance/internal/util"
)

// DcimCloudClient 魔方云API客户端
// 对接魔方云虚拟化系统，实现自动开小鸡
// 参考zjmf源码: app/common/logic/DcimCloud.php
// 所有路径从zjmf源码逐个确认，basecurl自动加/v1前缀
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

	// 解密密码（zjmf用aesPasswordDecode）
	password := server.Password
	if decrypted, err := util.DecryptAES(server.Password); err == nil {
		password = decrypted
	}

	return &DcimCloudClient{
		ServerID:   serverID,
		URL:        baseURL,
		Username:   server.Username,
		Password:   password,
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
// basecurl自动拼接 /v1 + action
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

// ============ 核心操作（路径全部从zjmf DcimCloud.php确认） ============

// On 开机（DcimCloud.php:1362）
func (c *DcimCloudClient) On(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/on", dcimID), nil, "POST")
}

// Off 关机（DcimCloud.php:1398）
func (c *DcimCloudClient) Off(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/off", dcimID), nil, "POST")
}

// Reboot 重启（DcimCloud.php:1434）
func (c *DcimCloudClient) Reboot(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/reboot", dcimID), nil, "POST")
}

// HardOff 强制关机（DcimCloud.php:1470）注意：zjmf用hardoff不是hard_off
func (c *DcimCloudClient) HardOff(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/hardoff", dcimID), nil, "POST")
}

// HardReboot 强制重启（DcimCloud.php:1506）
func (c *DcimCloudClient) HardReboot(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/hard_reboot", dcimID), nil, "POST")
}

// VNC 获取VNC连接（DcimCloud.php:1542）
func (c *DcimCloudClient) VNC(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/vnc", dcimID), nil, "POST")
}

// Reinstall 重装系统（DcimCloud.php:1657）注意：zjmf用PUT不是POST
func (c *DcimCloudClient) Reinstall(dcimID uint, data map[string]interface{}) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/reinstall", dcimID), data, "PUT")
}

// Rescue 救援模式（DcimCloud.php:3014）
func (c *DcimCloudClient) Rescue(dcimID uint, system string, tempPass string) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"type":      system,
		"temp_pass": tempPass,
	}
	return c.Curl(fmt.Sprintf("/clouds/%d/rescue", dcimID), data, "POST")
}

// ExitRescue 退出救援模式（DcimCloud.php:1290）
func (c *DcimCloudClient) ExitRescue(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/rescue", dcimID), nil, "DELETE")
}

// CrackPassword 重置密码（DcimCloud.php:2524）注意：zjmf用PUT不是POST
func (c *DcimCloudClient) CrackPassword(dcimID uint, newPass string) (map[string]interface{}, error) {
	data := map[string]interface{}{"password": newPass}
	return c.Curl(fmt.Sprintf("/clouds/%d/password", dcimID), data, "PUT")
}

// Status 获取状态（DcimCloud.php:1849）
func (c *DcimCloudClient) Status(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/status", dcimID), nil, "GET")
}

// Sync 同步信息（DcimCloud.php:1717）注意：zjmf用GET /clouds/{id}不是/info
func (c *DcimCloudClient) Sync(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d", dcimID), nil, "GET")
}

// RemoteInfo 远程连接信息（DcimCloud.php:1331）
func (c *DcimCloudClient) RemoteInfo(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d", dcimID), nil, "GET")
}

// GetCloudStatus 获取云状态（DcimCloud.php:1849）
func (c *DcimCloudClient) GetCloudStatus(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/status", dcimID), nil, "GET")
}

// ============ 快照操作（路径从zjmf确认） ============

// CreateSnap 创建快照（DcimCloud.php:318）
// zjmf: /disks/{disk_id}/snapshots?type=snap
func (c *DcimCloudClient) CreateSnap(diskID uint, name string) (map[string]interface{}, error) {
	data := map[string]interface{}{"type": "snap", "name": name}
	return c.Curl(fmt.Sprintf("/disks/%d/snapshots", diskID), data, "POST")
}

// DeleteSnap 删除快照（DcimCloud.php:400）
// zjmf: DELETE /snapshots/{snap_id}
func (c *DcimCloudClient) DeleteSnap(snapID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/snapshots/%d", snapID), nil, "DELETE")
}

// RestoreSnap 恢复快照（DcimCloud.php:447）
// zjmf: POST /snapshots/{snap_id}/restore?hostid={dcimid}
func (c *DcimCloudClient) RestoreSnap(snapID uint, dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/snapshots/%d/restore?hostid=%d", snapID, dcimID), nil, "POST")
}

// ListSnapBackup 快照列表（DcimCloud.php:306）
// zjmf: GET /clouds/{dcimid}/snapshots
func (c *DcimCloudClient) ListSnapBackup(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/snapshots", dcimID), nil, "GET")
}

// ============ 备份操作（路径从zjmf确认） ============

// CreateBackup 创建备份（DcimCloud.php:501）
// zjmf: /disks/{disk_id}/snapshots?type=backup
func (c *DcimCloudClient) CreateBackup(diskID uint, name string) (map[string]interface{}, error) {
	data := map[string]interface{}{"type": "backup", "name": name}
	return c.Curl(fmt.Sprintf("/disks/%d/snapshots", diskID), data, "POST")
}

// DeleteBackup 删除备份（DcimCloud.php:548）
// zjmf: DELETE /snapshots/{snap_id}（和快照共用同一API）
func (c *DcimCloudClient) DeleteBackup(snapID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/snapshots/%d", snapID), nil, "DELETE")
}

// RestoreBackup 恢复备份（DcimCloud.php:595）
// zjmf: POST /snapshots/{snap_id}/restore?hostid={dcimid}
func (c *DcimCloudClient) RestoreBackup(snapID uint, dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/snapshots/%d/restore?hostid=%d", snapID, dcimID), nil, "POST")
}

// ============ ISO操作（路径从zjmf确认） ============

// MountIso 挂载ISO（DcimCloud.php:642）
// zjmf: POST /clouds/{dcimid}/iso
func (c *DcimCloudClient) MountIso(dcimID uint, isoID string) (map[string]interface{}, error) {
	data := map[string]interface{}{"iso": isoID}
	return c.Curl(fmt.Sprintf("/clouds/%d/iso", dcimID), data, "POST")
}

// UnmountIso 卸载ISO（DcimCloud.php:689）
// zjmf: DELETE /clouds/{dcimid}/iso
func (c *DcimCloudClient) UnmountIso(dcimID uint, isoID string) (map[string]interface{}, error) {
	data := map[string]interface{}{"iso": isoID}
	return c.Curl(fmt.Sprintf("/clouds/%d/iso", dcimID), data, "DELETE")
}

// IsoList 获取可用ISO列表（DcimCloud.php:185）
// zjmf: GET /node_isos?id={node_id}&type=node
func (c *DcimCloudClient) IsoList(nodeID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/node_isos?id=%d&type=node", nodeID), nil, "GET")
}

// SetBootOrder 设置启动顺序（DcimCloud.php:736）
// zjmf: PUT /clouds/{dcimid}（注意是PUT不是POST）
func (c *DcimCloudClient) SetBootOrder(dcimID uint, bootOrder string) (map[string]interface{}, error) {
	data := map[string]interface{}{"bootorder": bootOrder}
	return c.Curl(fmt.Sprintf("/clouds/%d", dcimID), data, "PUT")
}

// ============ 安全组操作（路径从zjmf确认） ============

// CreateSecurityGroup 创建安全组（DcimCloud.php:788）
// zjmf: POST /security_groups
func (c *DcimCloudClient) CreateSecurityGroup(data map[string]interface{}) (map[string]interface{}, error) {
	return c.Curl("/security_groups", data, "POST")
}

// DelSecurityGroup 删除安全组（DcimCloud.php:952）
// zjmf: DELETE /security_groups/{group_id}
func (c *DcimCloudClient) DelSecurityGroup(groupID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/security_groups/%d", groupID), nil, "DELETE")
}

// ShowSecurityRules 显示安全规则（DcimCloud.php:830）
// zjmf: GET /security_groups/{group_id}/rules
func (c *DcimCloudClient) ShowSecurityRules(groupID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/security_groups/%d/rules", groupID), nil, "GET")
}

// CreateSecurityRule 创建安全规则（DcimCloud.php:892）
// zjmf: POST /security_groups/{group_id}/rules
func (c *DcimCloudClient) CreateSecurityRule(groupID uint, data map[string]interface{}) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/security_groups/%d/rules", groupID), data, "POST")
}

// DelSecurityRule 删除安全规则（DcimCloud.php:1020）
// zjmf: DELETE /security_group_rules/{rule_id}
func (c *DcimCloudClient) DelSecurityRule(ruleID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/security_group_rules/%d", ruleID), nil, "DELETE")
}

// LinkSecurityGroup 关联安全组（DcimCloud.php:1069）
// zjmf: POST /security_groups/{group_id}/links
func (c *DcimCloudClient) LinkSecurityGroup(groupID uint, data map[string]interface{}) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/security_groups/%d/links", groupID), data, "POST")
}

// ============ NAT操作（路径从zjmf确认） ============

// GetNatInfo 获取NAT信息（DcimCloud.php:256+260）
// zjmf: GET /clouds/{dcimid}/nat_acl + /clouds/{dcimid}/nat_web
func (c *DcimCloudClient) GetNatInfo(dcimID uint) (map[string]interface{}, error) {
	acl, _ := c.Curl(fmt.Sprintf("/clouds/%d/nat_acl?per_page=1&sort=asc", dcimID), nil, "GET")
	web, _ := c.Curl(fmt.Sprintf("/clouds/%d/nat_web?per_page=1&sort=asc", dcimID), nil, "GET")
	return map[string]interface{}{"nat_acl": acl, "nat_web": web}, nil
}

// AddNatAcl 添加NAT ACL（DcimCloud.php:1117）
// zjmf: POST /clouds/{dcimid}/nat_acl
func (c *DcimCloudClient) AddNatAcl(dcimID uint, data map[string]interface{}) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/nat_acl", dcimID), data, "POST")
}

// DelNatAcl 删除NAT ACL（DcimCloud.php:1165）
// zjmf: DELETE /nat_acl/{acl_id}?hostid={dcimid}
func (c *DcimCloudClient) DelNatAcl(aclID uint, dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/nat_acl/%d?hostid=%d", aclID, dcimID), nil, "DELETE")
}

// AddNatWeb 添加NAT Web（DcimCloud.php:1207）
// zjmf: POST /clouds/{dcimid}/nat_web
func (c *DcimCloudClient) AddNatWeb(dcimID uint, data map[string]interface{}) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/nat_web", dcimID), data, "POST")
}

// DelNatWeb 删除NAT Web（DcimCloud.php:1249）
// zjmf: DELETE /nat_web/{web_id}?hostid={dcimid}
func (c *DcimCloudClient) DelNatWeb(webID uint, dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/nat_web/%d?hostid=%d", webID, dcimID), nil, "DELETE")
}

// ============ 流量操作（路径从zjmf确认） ============

// GetTrafficUsage 获取流量使用（DcimCloud.php:2846-2878）
// zjmf: GET /clouds/{dcimid}/flow_data，时间戳用毫秒级，单位GB
func (c *DcimCloudClient) GetTrafficUsage(dcimID uint, start string, end string) (map[string]interface{}, error) {
	// zjmf: strtotime($start . " 00:00:00") . "000"（毫秒级时间戳）
	startTs := fmt.Sprintf("%d000", time.Now().AddDate(0, 0, -30).Unix())
	endTs := fmt.Sprintf("%d000", time.Now().Unix())
	if start != "" {
		if t, err := time.Parse("2006-01-02", start); err == nil {
			startTs = fmt.Sprintf("%d000", t.Unix())
		}
	}
	if end != "" {
		if t, err := time.Parse("2006-01-02", end); err == nil {
			endTs = fmt.Sprintf("%d000", t.Add(24*time.Hour-time.Second).Unix())
		}
	}
	data := map[string]interface{}{
		"type":       2,
		"start_time": startTs,
		"end_time":   endTs,
		"unit":       "GB",
	}
	return c.Curl(fmt.Sprintf("/clouds/%d/flow_data", dcimID), data, "GET")
}

// GetChartData 获取图表数据（DcimCloud.php:2379）
// zjmf: GET /statistics?host_id={dcimid}&type={type}&start={start}&end={end}
func (c *DcimCloudClient) GetChartData(dcimID uint, chartType string, start string, end string) (map[string]interface{}, error) {
	data := map[string]interface{}{"host_id": dcimID, "type": chartType, "start": start, "end": end}
	return c.Curl("/statistics", data, "GET")
}

// BuyFlowPacket 购买流量包（DcimCloud.php:1774）
// zjmf: PUT /clouds/{dcimid}/temp_traffic
func (c *DcimCloudClient) BuyFlowPacket(dcimID uint, data map[string]interface{}) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/temp_traffic", dcimID), data, "PUT")
}

// ResetFlow 重置流量（DcimCloud.php:2904）
// zjmf: POST /clouds/{dcimid}/reset_traffic
func (c *DcimCloudClient) ResetFlow(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/reset_traffic", dcimID), nil, "POST")
}

// GetNetInfo 获取网络信息（DcimCloud.php:1777）
// zjmf: GET /net_info?host_id={dcimid}
func (c *DcimCloudClient) GetNetInfo(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/net_info?host_id=%d", dcimID), nil, "GET")
}

// ============ 创建/开通操作（路径从zjmf确认） ============

// CreateAccount 开通虚拟机（DcimCloud.php:2787）
// zjmf: POST /clouds
func (c *DcimCloudClient) CreateAccount(params map[string]interface{}) (map[string]interface{}, error) {
	return c.Curl("/clouds", params, "POST")
}

// GetOsList 获取操作系统列表（DcimCloud.php:98）
// zjmf: GET /image?per_page=9999&sort=asc
func (c *DcimCloudClient) GetOsList() (map[string]interface{}, error) {
	return c.Curl("/image?per_page=9999&sort=asc", nil, "GET")
}

// GetAreaList 获取区域列表（DcimCloud.php:116）
// zjmf: GET /areas?sort=asc&list_type=all
func (c *DcimCloudClient) GetAreaList() (map[string]interface{}, error) {
	return c.Curl("/areas?sort=asc&list_type=all", nil, "GET")
}

// GetCommonConfig 获取通用配置（DcimCloud.php:278）
// zjmf: GET /common_config
func (c *DcimCloudClient) GetCommonConfig() (map[string]interface{}, error) {
	return c.Curl("/common_config", nil, "GET")
}

// GetVpcNetworks 获取VPC网络列表（DcimCloud.php:2652）
// zjmf: GET /vpc_networks
func (c *DcimCloudClient) GetVpcNetworks() (map[string]interface{}, error) {
	return c.Curl("/vpc_networks", nil, "GET")
}

// CreateUser 创建用户（DcimCloud.php:2595）
// zjmf: POST /user
func (c *DcimCloudClient) CreateUser(data map[string]interface{}) (map[string]interface{}, error) {
	return c.Curl("/user", data, "POST")
}

// CheckUser 检查用户（DcimCloud.php:2596）
// zjmf: GET /user/check?username={username}
func (c *DcimCloudClient) CheckUser(username string) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/user/check?username=%s", username), nil, "GET")
}

// GetSecurityGroupProtocols 获取安全组协议列表（DcimCloud.php:208）
// zjmf: GET /security_group_rule_protocols
func (c *DcimCloudClient) GetSecurityGroupProtocols() (map[string]interface{}, error) {
	return c.Curl("/security_group_rule_protocols", nil, "GET")
}

// GetSecurityGroups 获取安全组列表（DcimCloud.php:201）
// zjmf: GET /security_groups?list_type=all&user={user_id}&type={type}
func (c *DcimCloudClient) GetSecurityGroups(userID uint, cloudType string) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/security_groups?list_type=all&user=%d&type=%s", userID, cloudType), nil, "GET")
}

// ============ 升级操作（路径从zjmf确认） ============

// Upgrade 升级配置（DcimCloud.php:2093-2334）
// zjmf: PUT /clouds/{dcimid} 更新配置 + POST /clouds/{dcimid}/on 开机
func (c *DcimCloudClient) Upgrade(dcimID uint, data map[string]interface{}) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d", dcimID), data, "PUT")
}

// ============ 面板操作（路径从zjmf确认） ============

// ManagePanel 管理面板（DcimCloud.php:2547）
func (c *DcimCloudClient) ManagePanel(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d", dcimID), nil, "GET")
}

// SavePanelPass 保存面板密码（DcimCloud.php:2880-2900）
// zjmf存到customfields表，我们存到service.Config JSON的panel_password字段
func (c *DcimCloudClient) SavePanelPass(serviceID uint, password string) (map[string]interface{}, error) {
	db := database.GetDB()
	var svc model.Service
	if err := db.First(&svc, serviceID).Error; err != nil {
		return nil, fmt.Errorf("服务不存在")
	}
	// 解析现有config，添加panel_password
	config := make(map[string]interface{})
	if svc.Config != "" {
		json.Unmarshal([]byte(svc.Config), &config)
	}
	config["panel_password"] = password
	configBytes, _ := json.Marshal(config)
	db.Model(&svc).Update("config", string(configBytes))
	return map[string]interface{}{"status": "success"}, nil
}

// SupportReinstallRandomPort 是否支持随机端口重装（DcimCloud.php:265）
func (c *DcimCloudClient) SupportReinstallRandomPort(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/reinstall/random_port", dcimID), nil, "GET")
}

// ============ 生命周期操作（路径从zjmf确认） ============

// Suspend 暂停服务（DcimCloud.php:1917）
// zjmf: POST /clouds/{dcimid}/suspend
func (c *DcimCloudClient) Suspend(dcimID uint, reasonType string) (map[string]interface{}, error) {
	data := map[string]interface{}{"type": reasonType}
	return c.Curl(fmt.Sprintf("/clouds/%d/suspend", dcimID), data, "POST")
}

// Unsuspend 取消暂停（DcimCloud.php:1964）
// zjmf: POST /clouds/{dcimid}/unsuspend
func (c *DcimCloudClient) Unsuspend(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/unsuspend", dcimID), nil, "POST")
}

// Terminate 终止/删除（DcimCloud.php:2006）
// zjmf: DELETE /clouds/{dcimid}
func (c *DcimCloudClient) Terminate(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d", dcimID), nil, "DELETE")
}

// ModuleClientArea 客户端控制面板（DcimCloud.php:125-230）
// zjmf: 根据key调用不同API获取面板数据
func (c *DcimCloudClient) ModuleClientArea(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d", dcimID), nil, "GET")
}

// ModuleClientAreaDetail 客户端面板详情（DcimCloud.php:155-230）
// zjmf: 根据key调用不同API
func (c *DcimCloudClient) ModuleClientAreaDetail(dcimID uint, key string) (map[string]interface{}, error) {
	// zjmf DcimCloud.php:155-236，根据key返回不同面板数据
	switch key {
	case "snapshot":
		// zjmf: 获取云主机详情+快照列表+磁盘信息
		cloud, err := c.Curl(fmt.Sprintf("/clouds/%d", dcimID), nil, "GET")
		if err != nil {
			return nil, err
		}
		snaps, _ := c.Curl(fmt.Sprintf("/clouds/%d/snapshots?per_page=100", dcimID), nil, "GET")
		return map[string]interface{}{
			"status": 200,
			"data": map[string]interface{}{
				"list":           snaps["data"],
				"disk":           cloud["data"],
				"support_snap":   true,
				"support_backup": true,
				"host_type":      "host",
			},
		}, nil
	case "setting":
		// zjmf: 获取ISO列表+启动顺序+云主机类型
		cloud, err := c.Curl(fmt.Sprintf("/clouds/%d", dcimID), nil, "GET")
		if err != nil {
			return nil, err
		}
		nodeID := uint(0)
		if data, ok := cloud["data"].(map[string]interface{}); ok {
			if nid, ok := data["node_id"].(float64); ok {
				nodeID = uint(nid)
			}
		}
		isos, _ := c.Curl(fmt.Sprintf("/node_isos?id=%d&type=node", nodeID), nil, "GET")
		isoList := []map[string]interface{}{}
		if isosData, ok := isos["data"].([]interface{}); ok {
			for _, v := range isosData {
				if group, ok := v.(map[string]interface{}); ok {
					if info, ok := group["info"].([]interface{}); ok {
						for _, item := range info {
							if iso, ok := item.(map[string]interface{}); ok {
								isoList = append(isoList, map[string]interface{}{
									"id":   iso["id"],
									"name": iso["name"],
								})
							}
						}
					}
				}
			}
		}
		return map[string]interface{}{
			"status": 200,
			"data": map[string]interface{}{
				"iso":       cloud["data"],
				"bootorder": cloud["data"],
				"iso2":      isoList,
				"host_type": "host",
			},
		}, nil
	case "security_groups":
		// zjmf: 获取安全组列表+安全协议
		cloud, err := c.Curl(fmt.Sprintf("/clouds/%d", dcimID), nil, "GET")
		if err != nil {
			return nil, err
		}
		userID := uint(0)
		cloudType := "host"
		if data, ok := cloud["data"].(map[string]interface{}); ok {
			if uid, ok := data["user_id"].(float64); ok {
				userID = uint(uid)
			}
			if t, ok := data["type"].(string); ok {
				cloudType = t
			}
		}
		groups, _ := c.GetSecurityGroups(userID, cloudType)
		protocols, _ := c.GetSecurityGroupProtocols()
		return map[string]interface{}{
			"status": 200,
			"data": map[string]interface{}{
				"list":      groups["data"],
				"used":      cloud["data"],
				"protocols": protocols["data"],
				"host_type": cloudType,
			},
		}, nil
	case "nat_acl":
		// zjmf: 获取NAT ACL列表（DcimCloud.php:215）
		data, err := c.Curl(fmt.Sprintf("/clouds/%d/nat_acl?list_type=all", dcimID), nil, "GET")
		if err != nil {
			return nil, err
		}
		natHostIP := ""
		if d, ok := data["data"].(map[string]interface{}); ok {
			if ip, ok := d["nat_host_ip"].(string); ok {
				natHostIP = ip
			}
		}
		return map[string]interface{}{
			"status": 200,
			"data": map[string]interface{}{
				"list":        data["data"],
				"nat_host_ip": natHostIP,
			},
		}, nil
	case "nat_web":
		// zjmf: 获取NAT Web列表（DcimCloud.php:218）
		data, err := c.Curl(fmt.Sprintf("/clouds/%d/nat_web?list_type=all", dcimID), nil, "GET")
		if err != nil {
			return nil, err
		}
		natHostIP := ""
		if d, ok := data["data"].(map[string]interface{}); ok {
			if ip, ok := d["nat_host_ip"].(string); ok {
				natHostIP = ip
			}
		}
		return map[string]interface{}{
			"status": 200,
			"data": map[string]interface{}{
				"list":        data["data"],
				"nat_host_ip": natHostIP,
			},
		}, nil
	default:
		return c.Curl(fmt.Sprintf("/clouds/%d", dcimID), nil, "GET")
	}
}

// ModuleClientButton 客户端控制按钮（DcimCloud.php:2020-2027）
func (c *DcimCloudClient) ModuleClientButton() map[string]interface{} {
	return map[string]interface{}{
		"control": []map[string]interface{}{
			{"type": "default", "func": "on", "name": "开机"},
			{"type": "default", "func": "off", "name": "关机"},
			{"type": "default", "func": "reboot", "name": "重启"},
			{"type": "default", "func": "hard_off", "name": "硬关机"},
			{"type": "default", "func": "hard_reboot", "name": "硬重启"},
			{"type": "default", "func": "reinstall", "name": "重装系统"},
			{"type": "default", "func": "crack_pass", "name": "重置密码"},
			{"type": "default", "func": "rescue_system", "name": "救援系统"},
			{"type": "custom", "func": "exit_rescue", "name": "退出救援系统"},
		},
		"console": []map[string]interface{}{
			{"type": "default", "func": "vnc", "name": "VNC"},
		},
	}
}

// ModuleAdminButton 管理员控制按钮（DcimCloud.php:2537-2545）
func (c *DcimCloudClient) ModuleAdminButton() []map[string]interface{} {
	return []map[string]interface{}{
		{"type": "default", "func": "create", "name": "开通"},
		{"type": "default", "func": "suspend", "name": "暂停"},
		{"type": "default", "func": "unsuspend", "name": "解除暂停"},
		{"type": "default", "func": "terminate", "name": "删除"},
		{"type": "default", "func": "on", "name": "开机"},
		{"type": "default", "func": "off", "name": "关机"},
		{"type": "default", "func": "reboot", "name": "重启"},
		{"type": "default", "func": "hard_off", "name": "硬关机"},
		{"type": "default", "func": "hard_reboot", "name": "硬重启"},
		{"type": "default", "func": "reinstall", "name": "重装系统"},
		{"type": "default", "func": "crack_pass", "name": "重置密码"},
		{"type": "default", "func": "rescue_system", "name": "救援系统"},
		{"type": "custom", "func": "exit_rescue", "name": "退出救援系统"},
		{"type": "default", "func": "vnc", "name": "VNC"},
		{"type": "default", "func": "sync", "name": "拉取信息"},
	}
}

// ExecCustomButton 执行自定义按钮（DcimCloud.php:2028-2077）
// zjmf: 支持 exit_rescue/download_rdp
func (c *DcimCloudClient) ExecCustomButton(dcimID uint, funcName string) (map[string]interface{}, error) {
	switch funcName {
	case "exit_rescue":
		return c.ExitRescue(dcimID)
	case "download_rdp":
		return c.Curl(fmt.Sprintf("/clouds/%d/download_rdp", dcimID), nil, "POST")
	default:
		return map[string]interface{}{"status": 400, "msg": "不支持该功能"}, nil
	}
}

// ModuleAllowFunction 允许的功能列表（DcimCloud.php:284）
func (c *DcimCloudClient) ModuleAllowFunction() []string {
	return []string{
		"listSnapBackup", "createSnap", "deleteSnap", "restoreSnap",
		"createBackup", "deleteBackup", "restoreBackup",
		"mountIso", "unmountIso", "setBootOrder",
		"createSecurityGroup", "delSecurityGroup", "showSecurityRules",
		"createSecurityRule", "delSecurityRule", "linkSecurityGroup",
		"addNatAcl", "delNatAcl", "addNatWeb", "delNatWeb",
		"exitRescue", "remoteInfo",
	}
}

// ============ 磁盘/网络高级操作（路径从zjmf确认） ============

// AdjustBandwidth 带宽调整（DcimCloud.php:2209）
// zjmf: PUT /clouds/{dcimid}/bw
func (c *DcimCloudClient) AdjustBandwidth(dcimID uint, bw int) (map[string]interface{}, error) {
	data := map[string]interface{}{"bw": bw}
	return c.Curl(fmt.Sprintf("/clouds/%d/bw", dcimID), data, "PUT")
}

// AdjustIP IP数量调整（DcimCloud.php:2215）
// zjmf: PUT /clouds/{dcimid}/ip
func (c *DcimCloudClient) AdjustIP(dcimID uint, num int, ipGroup string) (map[string]interface{}, error) {
	data := map[string]interface{}{"num": num, "ip_group": ipGroup}
	return c.Curl(fmt.Sprintf("/clouds/%d/ip", dcimID), data, "PUT")
}

// CreateDisk 创建数据盘（DcimCloud.php:2256）
// zjmf: POST /clouds/{dcimid}/disks
func (c *DcimCloudClient) CreateDisk(dcimID uint, size int, storeID int, driver string) (map[string]interface{}, error) {
	data := map[string]interface{}{"size": size, "store": storeID, "driver": driver}
	return c.Curl(fmt.Sprintf("/clouds/%d/disks", dcimID), data, "POST")
}

// GetStores 获取存储列表（DcimCloud.php:2252）
// zjmf: GET /clouds/{dcimid}/stores
func (c *DcimCloudClient) GetStores(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/stores", dcimID), nil, "GET")
}

// ResizeDisk 磁盘扩容（DcimCloud.php:2249）
// zjmf: PUT /disks/{disk_id}
func (c *DcimCloudClient) ResizeDisk(diskID uint, newSize int) (map[string]interface{}, error) {
	data := map[string]interface{}{"size": newSize}
	return c.Curl(fmt.Sprintf("/disks/%d", diskID), data, "PUT")
}

// AdjustIPv6 IPv6调整（DcimCloud.php:2267）
// zjmf: PUT /clouds/{dcimid}/ipv6
func (c *DcimCloudClient) AdjustIPv6(dcimID uint, num int) (map[string]interface{}, error) {
	data := map[string]interface{}{"num": num}
	return c.Curl(fmt.Sprintf("/clouds/%d/ipv6", dcimID), data, "PUT")
}

// AdjustConfig 调整配置（DcimCloud.php:2206）
// zjmf: PUT /clouds/{dcimid}（CPU/内存/快照数/备份数等）
func (c *DcimCloudClient) AdjustConfig(dcimID uint, data map[string]interface{}) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d", dcimID), data, "PUT")
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
