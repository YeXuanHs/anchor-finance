package service

import "fmt"

// ============ 非关键补充方法（参考zjmf DcimCloud.php） ============

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

// CreateBackup 创建备份（DcimCloud.php:466）
func (c *DcimCloudClient) CreateBackup(dcimID uint, name string) (map[string]interface{}, error) {
	data := map[string]interface{}{"name": name}
	return c.Curl(fmt.Sprintf("/clouds/%d/backup", dcimID), data, "POST")
}

// DeleteBackup 删除备份（DcimCloud.php:520）
func (c *DcimCloudClient) DeleteBackup(dcimID uint, backupID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/backup/%d", dcimID, backupID), nil, "DELETE")
}

// RestoreBackup 恢复备份（DcimCloud.php:567）
func (c *DcimCloudClient) RestoreBackup(dcimID uint, backupID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/backup/%d/restore", dcimID, backupID), nil, "POST")
}

// ListBackup 备份列表（DcimCloud.php:288 快照+备份共用）
func (c *DcimCloudClient) ListBackup(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/backup", dcimID), nil, "GET")
}

// MountIso 挂载ISO（DcimCloud.php:614）
func (c *DcimCloudClient) MountIso(dcimID uint, isoPath string) (map[string]interface{}, error) {
	data := map[string]interface{}{"iso": isoPath}
	return c.Curl(fmt.Sprintf("/clouds/%d/iso/mount", dcimID), data, "POST")
}

// UnmountIso 卸载ISO（DcimCloud.php:661）
func (c *DcimCloudClient) UnmountIso(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/iso/unmount", dcimID), nil, "POST")
}

// SetBootOrder 设置启动顺序（DcimCloud.php:708）
func (c *DcimCloudClient) SetBootOrder(dcimID uint, order string) (map[string]interface{}, error) {
	data := map[string]interface{}{"boot_order": order}
	return c.Curl(fmt.Sprintf("/clouds/%d/boot_order", dcimID), data, "POST")
}

// DelSecurityGroup 删除安全组（DcimCloud.php:911）
func (c *DcimCloudClient) DelSecurityGroup(dcimID uint, groupID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/security_group/%d", dcimID, groupID), nil, "DELETE")
}

// ShowSecurityRules 显示安全规则（DcimCloud.php:807）
func (c *DcimCloudClient) ShowSecurityRules(dcimID uint, groupID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/security_group/%d/rule", dcimID, groupID), nil, "GET")
}

// DelSecurityRule 删除安全规则（DcimCloud.php:971）
func (c *DcimCloudClient) DelSecurityRule(dcimID uint, ruleID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/security_rule/%d", dcimID, ruleID), nil, "DELETE")
}

// LinkSecurityGroup 关联安全组（DcimCloud.php:1039）
func (c *DcimCloudClient) LinkSecurityGroup(dcimID uint, groupID uint) (map[string]interface{}, error) {
	data := map[string]interface{}{"group_id": groupID}
	return c.Curl(fmt.Sprintf("/clouds/%d/security_group/link", dcimID), data, "POST")
}

// GetNatInfo 获取NAT信息（DcimCloud.php:238）
func (c *DcimCloudClient) GetNatInfo(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/nat", dcimID), nil, "GET")
}

// ModuleAdminButton 管理员控制按钮（DcimCloud.php:2537）
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

// ManagePanel 管理面板（DcimCloud.php:2547）
func (c *DcimCloudClient) ManagePanel(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/panel", dcimID), nil, "GET")
}

// SavePanelPass 保存面板密码（DcimCloud.php:2880）
func (c *DcimCloudClient) SavePanelPass(dcimID uint, password string) (map[string]interface{}, error) {
	data := map[string]interface{}{"password": password}
	return c.Curl(fmt.Sprintf("/clouds/%d/panel/password", dcimID), data, "POST")
}

// ResetFlow 重置流量（DcimCloud.php:2902）
func (c *DcimCloudClient) ResetFlow(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/flow/reset", dcimID), nil, "POST")
}

// SupportReinstallRandomPort 是否支持随机端口重装（DcimCloud.php:265）
func (c *DcimCloudClient) SupportReinstallRandomPort(dcimID uint) (map[string]interface{}, error) {
	return c.Curl(fmt.Sprintf("/clouds/%d/reinstall/random_port", dcimID), nil, "GET")
}

// IsoList ISO列表（DcimCloud.php:614 挂载前获取可用ISO）
func (c *DcimCloudClient) IsoList() (map[string]interface{}, error) {
	return c.Curl("/iso", nil, "GET")
}
