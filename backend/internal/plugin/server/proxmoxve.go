package server

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ProxmoxVEConfig ProxmoxVE配置
type ProxmoxVEConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Node     string `json:"node"`
}

// ProxmoxVEModule ProxmoxVE服务器模块
type ProxmoxVEModule struct {
	config *ProxmoxVEConfig
	client *http.Client
}

func init() {
	RegisterModule("proxmoxve", NewProxmoxVEModule)
}

// NewProxmoxVEModule 创建ProxmoxVE模块
func NewProxmoxVEModule(config map[string]interface{}) (ServerModule, error) {
	cfg := &ProxmoxVEConfig{
		Port: 8006,
	}

	if config != nil {
		if host, ok := config["host"].(string); ok {
			cfg.Host = host
		}
		if port, ok := config["port"].(float64); ok {
			cfg.Port = int(port)
		}
		if username, ok := config["username"].(string); ok {
			cfg.Username = username
		}
		if password, ok := config["password"].(string); ok {
			cfg.Password = password
		}
		if node, ok := config["node"].(string); ok {
			cfg.Node = node
		}
	}

	return &ProxmoxVEModule{
		config: cfg,
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}, nil
}

func (m *ProxmoxVEModule) Name() string  { return "proxmoxve" }
func (m *ProxmoxVEModule) Title() string { return "ProxmoxVE" }

// GetConfigOptions 获取配置选项
func (m *ProxmoxVEModule) GetConfigOptions() []ConfigOption {
	return []ConfigOption{
		{Type: "text", Name: "ProxmoxVE面板地址", Key: "panel", Placeholder: "https://pve.example.com:8006/", Description: "ProxmoxVE面板地址，面向用户", Required: true},
		{Type: "text", Name: "节点名称", Key: "node", Placeholder: "pve", Description: "PVE面板内显示的节点名称", Required: true},
		{Type: "text", Name: "CPU数量", Key: "cores", Placeholder: "1", Description: "CPU数量（个）", Required: true},
		{Type: "text", Name: "内存大小", Key: "memory", Placeholder: "512", Description: "内存大小（MiB）", Required: true},
		{Type: "text", Name: "交换分区", Key: "swap", Placeholder: "512", Description: "交换分区大小（MiB）"},
		{Type: "text", Name: "系统镜像", Key: "ostemplate", Placeholder: "local:vztmpl/debian-10.tar.gz", Description: "系统镜像路径"},
		{Type: "text", Name: "磁盘大小", Key: "disksize", Placeholder: "20", Description: "磁盘大小（GiB）", Required: true},
		{Type: "text", Name: "存储位置", Key: "storage", Placeholder: "local-lvm", Description: "磁盘存储位置", Required: true},
		{Type: "text", Name: "桥接网卡", Key: "bridge", Placeholder: "vmbr0", Description: "桥接的网卡"},
		{Type: "text", Name: "IP地址范围", Key: "ip", Placeholder: "172.16.1.", Description: "IP地址范围（最后一位不写）"},
	}
}

// TestConnection 测试连接
func (m *ProxmoxVEModule) TestConnection(ctx context.Context, params *ConnectionParams) error {
	url := fmt.Sprintf("https://%s:%d/api2/json/version", params.Host, params.Port)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	// 添加认证
	req.SetBasicAuth(params.Username, params.Password)

	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("连接失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("认证失败，状态码: %d", resp.StatusCode)
	}

	return nil
}

// CreateAccount 创建LXC容器
func (m *ProxmoxVEModule) CreateAccount(ctx context.Context, params *CreateAccountParams) (*AccountInfo, error) {
	// 生成随机VMID和密码
	vmid := fmt.Sprintf("%d", 50000+time.Now().UnixNano()%50000)
	vmname := "PVELXC" + vmid
	password := generateRandomPassword(8)

	// 解析配置选项
	configOptions := params.ConfigOptions
	cores := configOptions["cores"]
	memory := configOptions["memory"]
	swap := configOptions["swap"]
	storage := configOptions["storage"]
	disksize := configOptions["disksize"]
	bridge := configOptions["bridge"]
	ostemplate := configOptions["ostemplate"]
	ip := configOptions["ip"]
	node := configOptions["node"]

	// 获取认证token
	ticket, err := m.getTicket(ctx, params.Host, params.Port, params.Username, params.Password)
	if err != nil {
		return nil, err
	}

	// 创建用户
	userData := map[string]string{
		"userid":   vmname + "@pve",
		"password": password,
		"enable":   "1",
	}
	if err := m.apiRequest(ctx, params.Host, params.Port, "POST", "/access/users", userData, ticket); err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	// 创建LXC容器
	vmData := map[string]string{
		"vmid":         vmid,
		"hostname":     vmname,
		"password":     password,
		"ostemplate":   ostemplate,
		"cores":        cores,
		"memory":       memory,
		"swap":         swap,
		"storage":      storage,
		"rootfs":       storage + ":" + disksize,
		"net0":         fmt.Sprintf("name=eth0,bridge=%s,ip=%s%d/24,gw=%s1,ip6=auto,firewall=1,type=veth", bridge, ip, 10+time.Now().UnixNano()%224, ip),
		"onboot":       "0",
		"unprivileged": "1",
		"arch":         "amd64",
	}

	if err := m.apiRequest(ctx, params.Host, params.Port, "POST", "/nodes/"+node+"/lxc", vmData, ticket); err != nil {
		return nil, fmt.Errorf("创建容器失败: %w", err)
	}

	// 设置ACL权限
	aclData := map[string]string{
		"path":     "/vms/" + vmid,
		"users":    vmname + "@pve",
		"roles":    "PVEVMUser",
		"propagate": "1",
	}
	if err := m.apiRequest(ctx, params.Host, params.Port, "PUT", "/access/acl", aclData, ticket); err != nil {
		return nil, fmt.Errorf("设置权限失败: %w", err)
	}

	return &AccountInfo{
		Username: vmname,
		Password: password,
		Hostname: vmname,
		OS:       "Linux",
	}, nil
}

// TerminateAccount 删除LXC容器
func (m *ProxmoxVEModule) TerminateAccount(ctx context.Context, params *TerminateAccountParams) error {
	ticket, err := m.getTicket(ctx, params.Host, params.Port, params.Username, params.Password)
	if err != nil {
		return err
	}

	// 获取VMID
	vmid := strings.TrimPrefix(params.Username, "PVELXC")
	node := ""
	if m.config != nil {
		node = m.config.Node
	}

	// 检查状态并停止
	statusURL := fmt.Sprintf("/nodes/%s/lxc/%s/status/current", node, vmid)
	// 尝试停止容器
	stopData := map[string]string{}
	m.apiRequest(ctx, params.Host, params.Port, "POST", fmt.Sprintf("/nodes/%s/lxc/%s/status/stop", node, vmid), stopData, ticket)

	// 等待停止
	time.Sleep(2 * time.Second)

	// 删除容器
	if err := m.apiRequest(ctx, params.Host, params.Port, "DELETE", fmt.Sprintf("/nodes/%s/lxc/%s", node, vmid), nil, ticket); err != nil {
		return fmt.Errorf("删除容器失败: %w", err)
	}

	// 删除用户
	userURL := "/access/users/" + params.Username + "@pve"
	m.apiRequest(ctx, params.Host, params.Port, "DELETE", userURL, nil, ticket)

	return nil
}

// SuspendAccount 暂停账户
func (m *ProxmoxVEModule) SuspendAccount(ctx context.Context, params *SuspendAccountParams) error {
	ticket, err := m.getTicket(ctx, params.Host, params.Port, params.Username, params.Password)
	if err != nil {
		return err
	}

	vmid := strings.TrimPrefix(params.Username, "PVELXC")
	node := ""
	if m.config != nil {
		node = m.config.Node
	}

	// 暂停容器
	suspendData := map[string]string{}
	return m.apiRequest(ctx, params.Host, params.Port, "POST", fmt.Sprintf("/nodes/%s/lxc/%s/status/suspend", node, vmid), suspendData, ticket)
}

// UnsuspendAccount 取消暂停
func (m *ProxmoxVEModule) UnsuspendAccount(ctx context.Context, params *UnsuspendAccountParams) error {
	ticket, err := m.getTicket(ctx, params.Host, params.Port, params.Username, params.Password)
	if err != nil {
		return err
	}

	vmid := strings.TrimPrefix(params.Username, "PVELXC")
	node := ""
	if m.config != nil {
		node = m.config.Node
	}

	// 恢复容器
	resumeData := map[string]string{}
	return m.apiRequest(ctx, params.Host, params.Port, "POST", fmt.Sprintf("/nodes/%s/lxc/%s/status/resume", node, vmid), resumeData, ticket)
}

// GetStatus 获取容器状态
func (m *ProxmoxVEModule) GetStatus(ctx context.Context, params *StatusParams) (*ServerStatus, error) {
	ticket, err := m.getTicket(ctx, params.Host, params.Port, params.Username, params.Password)
	if err != nil {
		return &ServerStatus{Status: "unknown", Desc: "连接失败"}, nil
	}

	vmid := strings.TrimPrefix(params.Username, "PVELXC")
	node := ""
	if m.config != nil {
		node = m.config.Node
	}

	url := fmt.Sprintf("https://%s:%d/api2/json/nodes/%s/lxc/%s/status/current", params.Host, params.Port, node, vmid)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return &ServerStatus{Status: "unknown", Desc: "请求失败"}, nil
	}

	req.Header.Set("Cookie", "PVEAuthCookie="+ticket)

	resp, err := m.client.Do(req)
	if err != nil {
		return &ServerStatus{Status: "unknown", Desc: "获取状态失败"}, nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if data, ok := result["data"].(map[string]interface{}); ok {
		if status, ok := data["status"].(string); ok {
			switch status {
			case "running":
				return &ServerStatus{Status: "on", Desc: "运行中"}, nil
			case "stopped":
				return &ServerStatus{Status: "off", Desc: "已停止"}, nil
			}
		}
	}

	return &ServerStatus{Status: "unknown", Desc: "未知状态"}, nil
}

// GetClientArea 获取客户端信息
func (m *ProxmoxVEModule) GetClientArea(ctx context.Context, params *ClientAreaParams) (map[string]interface{}, error) {
	panel := ""
	if configOptions, ok := params.Config["config_options"].(map[string]string); ok {
		panel = configOptions["panel"]
	}

	return map[string]interface{}{
		"goPanel": map[string]interface{}{
			"name": "控制面板信息",
			"url":  panel,
			"username": params.Username,
			"password": params.Password,
		},
	}, nil
}

// ChangePassword 修改密码
func (m *ProxmoxVEModule) ChangePassword(ctx context.Context, params *ChangePasswordParams) error {
	ticket, err := m.getTicket(ctx, params.Host, params.Port, params.Username, params.Password)
	if err != nil {
		return err
	}

	userURL := "/access/users/" + params.Username + "@pve"
	passwordData := map[string]string{
		"password": params.NewPassword,
	}
	return m.apiRequest(ctx, params.Host, params.Port, "PUT", userURL, passwordData, ticket)
}

// getTicket 获取认证ticket
func (m *ProxmoxVEModule) getTicket(ctx context.Context, host string, port int, username, password string) (string, error) {
	url := fmt.Sprintf("https://%s:%d/api2/json/access/ticket", host, port)
	data := fmt.Sprintf("username=%s&password=%s", username, password)

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := m.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if data, ok := result["data"].(map[string]interface{}); ok {
		if ticket, ok := data["ticket"].(string); ok {
			return ticket, nil
		}
	}

	return "", fmt.Errorf("获取ticket失败")
}

// apiRequest 发送API请求
func (m *ProxmoxVEModule) apiRequest(ctx context.Context, host string, port int, method, path string, data map[string]string, ticket string) error {
	url := fmt.Sprintf("https://%s:%d/api2/json%s", host, port, path)

	var req *http.Request
	var err error

	if data != nil {
		values := ""
		for k, v := range data {
			if values != "" {
				values += "&"
			}
			values += k + "=" + v
		}
		req, err = http.NewRequestWithContext(ctx, method, url, strings.NewReader(values))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
	}

	if err != nil {
		return err
	}

	req.Header.Set("Cookie", "PVEAuthCookie="+ticket)

	resp, err := m.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API请求失败: %s", string(body))
	}

	return nil
}

// generateRandomPassword 生成随机密码
func generateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(time.Nanosecond)
	}
	return string(result)
}
