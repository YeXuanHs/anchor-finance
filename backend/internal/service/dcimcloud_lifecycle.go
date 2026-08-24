package service

import (
	"fmt"
)

// ============ 服务生命周期操作（参考zjmf DcimCloud.php） ============

// Suspend 暂停服务（DcimCloud.php:1891）
// zjmf: 调 /clouds/{dcimid}/suspend，type=traffic/due/other
func (c *DcimCloudClient) Suspend(dcimID uint, reasonType string) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"type": reasonType, // traffic/due/other
	}
	return c.Curl(fmt.Sprintf("/clouds/%d/suspend", dcimID), data, "POST")
}

// Unsuspend 取消暂停（DcimCloud.php:1943）
// zjmf: 调 /clouds/{dcimid}/unsuspend
func (c *DcimCloudClient) Unsuspend(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/unsuspend", dcimID), nil, "POST")
}

// Terminate 终止/删除（DcimCloud.php:1987）
// zjmf: 调 DELETE /clouds/{dcimid}
func (c *DcimCloudClient) Terminate(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d", dcimID), nil, "DELETE")
}

// Upgrade 升级配置（DcimCloud.php:2078）
// zjmf: 调 /clouds/{dcimid}/upgrade
func (c *DcimCloudClient) Upgrade(dcimID uint, data map[string]interface{}) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/upgrade", dcimID), data, "POST")
}

// ModuleClientArea 客户端控制面板（DcimCloud.php:125）
// zjmf: 返回面板数据（远程信息+NAT信息+流量+安全组等）
func (c *DcimCloudClient) ModuleClientArea(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d", dcimID), nil, "GET")
}

// ModuleClientAreaDetail 客户端面板详情（DcimCloud.php:155）
// zjmf: 根据key返回不同详情（snap_backup/nat_info/security等）
func (c *DcimCloudClient) ModuleClientAreaDetail(dcimID uint, key string) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/%s", dcimID, key), nil, "GET")
}

// ModuleClientButton 客户端控制按钮（DcimCloud.php:2020）
// zjmf: 返回静态按钮列表，不调API
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

// ExecCustomButton 执行自定义按钮（DcimCloud.php:2028）
// zjmf: 支持 resume/exit_rescue/download_rdp
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

// GetCloudStatus 获取云状态（DcimCloud.php:1830）
// zjmf: 调 GET /clouds/{dcimid}/status
func (c *DcimCloudClient) GetCloudStatus(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/status", dcimID), nil, "GET")
}
