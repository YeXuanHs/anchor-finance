package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"anchorfinance/internal/model"

	"gorm.io/datatypes"
)

// ==================== NAT管理 ====================

// GetNATRules 获取云服务器NAT规则列表
func (s *DcimCloudService) GetNATRules(cloudID uint) ([]model.CloudNATRule, error) {
	var rules []model.CloudNATRule
	if err := s.db.Where("cloud_id = ?", cloudID).Order("id DESC").Find(&rules).Error; err != nil {
		return nil, err
	}
	return rules, nil
}

// CreateNATRule 创建NAT端口转发规则
func (s *DcimCloudService) CreateNATRule(cloudID uint, rule *model.CloudNATRule) error {
	server, err := s.GetServerByID(cloudID)
	if err != nil {
		return fmt.Errorf("cloud server not found: %w", err)
	}

	// 调用上游API
	params := map[string]interface{}{
		"name":     rule.Name,
		"protocol": rule.Protocol,
		"ext_port": rule.ExtPort,
		"int_port": rule.IntPort,
		"int_ip":   rule.IntIP,
	}
	if _, err := s.cloudAction(server, "create_nat_rule", params); err != nil {
		return fmt.Errorf("upstream create NAT rule: %w", err)
	}

	rule.CloudID = cloudID
	rule.Status = "active"
	if err := s.db.Create(rule).Error; err != nil {
		return fmt.Errorf("save NAT rule: %w", err)
	}
	return nil
}

// UpdateNATRule 更新NAT规则
func (s *DcimCloudService) UpdateNATRule(ruleID uint, updates map[string]interface{}) error {
	var rule model.CloudNATRule
	if err := s.db.First(&rule, ruleID).Error; err != nil {
		return fmt.Errorf("NAT rule not found: %w", err)
	}

	server, err := s.GetServerByID(rule.CloudID)
	if err != nil {
		return fmt.Errorf("cloud server not found: %w", err)
	}

	// 同步到上游
	upstreamParams := make(map[string]interface{})
	for k, v := range updates {
		upstreamParams[k] = v
	}
	upstreamParams["rule_id"] = ruleID
	if _, err := s.cloudAction(server, "update_nat_rule", upstreamParams); err != nil {
		return fmt.Errorf("upstream update NAT rule: %w", err)
	}

	return s.db.Model(&rule).Updates(updates).Error
}

// DeleteNATRule 删除NAT规则
func (s *DcimCloudService) DeleteNATRule(ruleID uint) error {
	var rule model.CloudNATRule
	if err := s.db.First(&rule, ruleID).Error; err != nil {
		return fmt.Errorf("NAT rule not found: %w", err)
	}

	server, err := s.GetServerByID(rule.CloudID)
	if err != nil {
		return fmt.Errorf("cloud server not found: %w", err)
	}

	// 调用上游删除
	params := map[string]interface{}{
		"rule_id": ruleID,
	}
	if _, err := s.cloudAction(server, "delete_nat_rule", params); err != nil {
		return fmt.Errorf("upstream delete NAT rule: %w", err)
	}

	return s.db.Delete(&rule).Error
}

// ==================== 安全组管理 ====================

// GetSecurityGroups 获取云服务器安全组列表
func (s *DcimCloudService) GetSecurityGroups(cloudID uint) ([]model.CloudSecurityGroup, error) {
	var groups []model.CloudSecurityGroup
	if err := s.db.Where("cloud_id = ?", cloudID).Order("id DESC").Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

// CreateSecurityGroup 创建安全组
func (s *DcimCloudService) CreateSecurityGroup(cloudID uint, group *model.CloudSecurityGroup) error {
	server, err := s.GetServerByID(cloudID)
	if err != nil {
		return fmt.Errorf("cloud server not found: %w", err)
	}

	params := map[string]interface{}{
		"name":           group.Name,
		"default_action": group.DefaultAction,
	}
	if _, err := s.cloudAction(server, "create_security_group", params); err != nil {
		return fmt.Errorf("upstream create security group: %w", err)
	}

	group.CloudID = cloudID
	if err := s.db.Create(group).Error; err != nil {
		return fmt.Errorf("save security group: %w", err)
	}
	return nil
}

// UpdateSecurityGroup 更新安全组
func (s *DcimCloudService) UpdateSecurityGroup(groupID uint, updates map[string]interface{}) error {
	var group model.CloudSecurityGroup
	if err := s.db.First(&group, groupID).Error; err != nil {
		return fmt.Errorf("security group not found: %w", err)
	}

	server, err := s.GetServerByID(group.CloudID)
	if err != nil {
		return fmt.Errorf("cloud server not found: %w", err)
	}

	upstreamParams := make(map[string]interface{})
	for k, v := range updates {
		upstreamParams[k] = v
	}
	upstreamParams["group_id"] = groupID
	if _, err := s.cloudAction(server, "update_security_group", upstreamParams); err != nil {
		return fmt.Errorf("upstream update security group: %w", err)
	}

	return s.db.Model(&group).Updates(updates).Error
}

// DeleteSecurityGroup 删除安全组
func (s *DcimCloudService) DeleteSecurityGroup(groupID uint) error {
	var group model.CloudSecurityGroup
	if err := s.db.First(&group, groupID).Error; err != nil {
		return fmt.Errorf("security group not found: %w", err)
	}

	server, err := s.GetServerByID(group.CloudID)
	if err != nil {
		return fmt.Errorf("cloud server not found: %w", err)
	}

	params := map[string]interface{}{
		"group_id": groupID,
	}
	if _, err := s.cloudAction(server, "delete_security_group", params); err != nil {
		return fmt.Errorf("upstream delete security group: %w", err)
	}

	// 同时删除关联规则
	if err := s.db.Where("group_id = ?", groupID).Delete(&model.CloudSecurityGroupRule{}).Error; err != nil {
		s.log.Warnf("delete security group rules: %v", err)
	}

	return s.db.Delete(&group).Error
}

// GetSecurityGroupRules 获取安全组规则列表
func (s *DcimCloudService) GetSecurityGroupRules(groupID uint) ([]model.CloudSecurityGroupRule, error) {
	var rules []model.CloudSecurityGroupRule
	if err := s.db.Where("group_id = ?", groupID).Order("priority ASC, id ASC").Find(&rules).Error; err != nil {
		return nil, err
	}
	return rules, nil
}

// AddSecurityGroupRule 添加安全组规则
func (s *DcimCloudService) AddSecurityGroupRule(groupID uint, rule *model.CloudSecurityGroupRule) error {
	var group model.CloudSecurityGroup
	if err := s.db.First(&group, groupID).Error; err != nil {
		return fmt.Errorf("security group not found: %w", err)
	}

	server, err := s.GetServerByID(group.CloudID)
	if err != nil {
		return fmt.Errorf("cloud server not found: %w", err)
	}

	params := map[string]interface{}{
		"group_id":  groupID,
		"direction": rule.Direction,
		"protocol":  rule.Protocol,
		"port_range": rule.PortRange,
		"source":    rule.Source,
		"action":    rule.Action,
		"priority":  rule.Priority,
	}
	if _, err := s.cloudAction(server, "add_security_group_rule", params); err != nil {
		return fmt.Errorf("upstream add security group rule: %w", err)
	}

	rule.GroupID = groupID
	if err := s.db.Create(rule).Error; err != nil {
		return fmt.Errorf("save security group rule: %w", err)
	}
	return nil
}

// DeleteSecurityGroupRule 删除安全组规则
func (s *DcimCloudService) DeleteSecurityGroupRule(ruleID uint) error {
	var rule model.CloudSecurityGroupRule
	if err := s.db.First(&rule, ruleID).Error; err != nil {
		return fmt.Errorf("security group rule not found: %w", err)
	}

	var group model.CloudSecurityGroup
	if err := s.db.First(&group, rule.GroupID).Error; err != nil {
		return fmt.Errorf("security group not found: %w", err)
	}

	server, err := s.GetServerByID(group.CloudID)
	if err != nil {
		return fmt.Errorf("cloud server not found: %w", err)
	}

	params := map[string]interface{}{
		"rule_id": ruleID,
	}
	if _, err := s.cloudAction(server, "delete_security_group_rule", params); err != nil {
		return fmt.Errorf("upstream delete security group rule: %w", err)
	}

	return s.db.Delete(&rule).Error
}

// ==================== ISO管理 ====================

// GetISOList 获取可用ISO镜像列表
func (s *DcimCloudService) GetISOList() ([]model.CloudISO, error) {
	var isos []model.CloudISO
	if err := s.db.Order("id DESC").Find(&isos).Error; err != nil {
		return nil, err
	}
	return isos, nil
}

// MountISO 挂载ISO到云服务器
func (s *DcimCloudService) MountISO(cloudID, isoID uint) error {
	server, err := s.GetServerByID(cloudID)
	if err != nil {
		return fmt.Errorf("cloud server not found: %w", err)
	}

	var iso model.CloudISO
	if err := s.db.First(&iso, isoID).Error; err != nil {
		return fmt.Errorf("ISO not found: %w", err)
	}

	params := map[string]interface{}{
		"iso_id": isoID,
		"iso_url": iso.URL,
	}
	if _, err := s.cloudAction(server, "mount_iso", params); err != nil {
		return fmt.Errorf("upstream mount ISO: %w", err)
	}

	// 更新ISO状态
	s.db.Model(&iso).Update("status", "mounted")
	return nil
}

// UnmountISO 卸载ISO
func (s *DcimCloudService) UnmountISO(cloudID uint) error {
	server, err := s.GetServerByID(cloudID)
	if err != nil {
		return fmt.Errorf("cloud server not found: %w", err)
	}

	params := map[string]interface{}{}
	if _, err := s.cloudAction(server, "unmount_iso", params); err != nil {
		return fmt.Errorf("upstream unmount ISO: %w", err)
	}

	// 将该服务器上挂载的ISO状态重置
	s.db.Model(&model.CloudISO{}).Where("status = ?", "mounted").Update("status", "available")
	return nil
}

// ==================== VNC控制台 ====================

// GetVNCURL 获取VNC控制台连接URL
func (s *DcimCloudService) GetVNCURL(cloudID uint) (string, error) {
	server, err := s.GetServerByID(cloudID)
	if err != nil {
		return "", fmt.Errorf("cloud server not found: %w", err)
	}

	result, err := s.cloudAction(server, "get_vnc_url", nil)
	if err != nil {
		return "", fmt.Errorf("upstream get VNC URL: %w", err)
	}

	if data, ok := result["data"].(map[string]interface{}); ok {
		if vncURL, ok := data["url"].(string); ok {
			return vncURL, nil
		}
	}
	return "", fmt.Errorf("VNC URL not found in response")
}

// GetVNCPage 获取VNC Web页面数据
func (s *DcimCloudService) GetVNCPage(cloudID uint) (map[string]interface{}, error) {
	server, err := s.GetServerByID(cloudID)
	if err != nil {
		return nil, fmt.Errorf("cloud server not found: %w", err)
	}

	result, err := s.cloudAction(server, "get_vnc_page", nil)
	if err != nil {
		return nil, fmt.Errorf("upstream get VNC page: %w", err)
	}

	if data, ok := result["data"].(map[string]interface{}); ok {
		return data, nil
	}
	return nil, fmt.Errorf("VNC page data not found in response")
}

// ==================== 监控图表 ====================

// GetResourceChart 获取资源监控图表数据
func (s *DcimCloudService) GetResourceChart(cloudID uint, period string) ([]model.CloudChart, error) {
	// 先尝试从本地数据库获取
	var charts []model.CloudChart
	var since time.Time

	switch period {
	case "1h":
		since = time.Now().Add(-1 * time.Hour)
	case "6h":
		since = time.Now().Add(-6 * time.Hour)
	case "24h":
		since = time.Now().Add(-24 * time.Hour)
	case "7d":
		since = time.Now().Add(-7 * 24 * time.Hour)
	case "30d":
		since = time.Now().Add(-30 * 24 * time.Hour)
	default:
		since = time.Now().Add(-24 * time.Hour)
	}

	if err := s.db.Where("cloud_id = ? AND timestamp >= ?", cloudID, since).
		Order("timestamp ASC").Find(&charts).Error; err != nil {
		return nil, err
	}

	// 本地无数据则尝试从上游获取
	if len(charts) == 0 {
		server, err := s.GetServerByID(cloudID)
		if err != nil {
			return nil, err
		}

		params := map[string]interface{}{
			"period": period,
		}
		result, err := s.cloudAction(server, "get_resource_chart", params)
		if err != nil {
			return nil, fmt.Errorf("upstream get resource chart: %w", err)
		}

		if data, ok := result["data"].([]interface{}); ok {
			for _, item := range data {
				if m, ok := item.(map[string]interface{}); ok {
					chart := model.CloudChart{
						CloudID: cloudID,
					}
					if ts, ok := m["timestamp"].(string); ok {
						t, _ := time.Parse(time.RFC3339, ts)
						chart.Timestamp = t
					}
					if v, ok := m["cpu_rate"].(float64); ok {
						chart.CPURate = v
					}
					if v, ok := m["memory_rate"].(float64); ok {
						chart.MemoryRate = v
					}
					if v, ok := m["disk_rate"].(float64); ok {
						chart.DiskRate = v
					}
					if v, ok := m["net_in"].(float64); ok {
						chart.NetIn = int64(v)
					}
					if v, ok := m["net_out"].(float64); ok {
						chart.NetOut = int64(v)
					}
					charts = append(charts, chart)
				}
			}
			// 批量保存到本地
			if len(charts) > 0 {
				s.db.CreateInBatches(charts, 100)
			}
		}
	}

	return charts, nil
}

// GetResourceInfo 获取当前资源使用情况
func (s *DcimCloudService) GetResourceInfo(cloudID uint) (map[string]interface{}, error) {
	server, err := s.GetServerByID(cloudID)
	if err != nil {
		return nil, fmt.Errorf("cloud server not found: %w", err)
	}

	result, err := s.cloudAction(server, "get_resource_info", nil)
	if err != nil {
		return nil, fmt.Errorf("upstream get resource info: %w", err)
	}

	if data, ok := result["data"].(map[string]interface{}); ok {
		return data, nil
	}
	return map[string]interface{}{}, nil
}

// ==================== 流量包管理 ====================

// GetFlowPackets 获取流量包列表
func (s *DcimCloudService) GetFlowPackets(cloudID uint) ([]model.CloudFlowPacket, error) {
	var packets []model.CloudFlowPacket
	if err := s.db.Where("cloud_id = ?", cloudID).Order("id DESC").Find(&packets).Error; err != nil {
		return nil, err
	}
	return packets, nil
}

// BuyFlowPacket 购买流量包
func (s *DcimCloudService) BuyFlowPacket(cloudID, packetID uint) error {
	server, err := s.GetServerByID(cloudID)
	if err != nil {
		return fmt.Errorf("cloud server not found: %w", err)
	}

	params := map[string]interface{}{
		"packet_id": packetID,
	}
	if _, err := s.cloudAction(server, "buy_flow_packet", params); err != nil {
		return fmt.Errorf("upstream buy flow packet: %w", err)
	}

	// 记录购买的流量包
	packet := &model.CloudFlowPacket{
		CloudID:   cloudID,
		Status:    "active",
		ExpiredAt: time.Now().Add(30 * 24 * time.Hour), // 默认30天
	}
	if err := s.db.Create(packet).Error; err != nil {
		return fmt.Errorf("save flow packet: %w", err)
	}
	return nil
}

// GetFlowPacketUsage 获取当前流量使用情况
func (s *DcimCloudService) GetFlowPacketUsage(cloudID uint) (map[string]interface{}, error) {
	server, err := s.GetServerByID(cloudID)
	if err != nil {
		return nil, fmt.Errorf("cloud server not found: %w", err)
	}

	result, err := s.cloudAction(server, "get_flow_usage", nil)
	if err != nil {
		return nil, fmt.Errorf("upstream get flow usage: %w", err)
	}

	if data, ok := result["data"].(map[string]interface{}); ok {
		return data, nil
	}
	return map[string]interface{}{}, nil
}

// ==================== 附加操作 ====================

// GetCloudStatus 获取云服务器实时状态
func (s *DcimCloudService) GetCloudStatus(cloudID uint) (map[string]interface{}, error) {
	server, err := s.GetServerByID(cloudID)
	if err != nil {
		return nil, fmt.Errorf("cloud server not found: %w", err)
	}

	result, err := s.cloudAction(server, "get_status", nil)
	if err != nil {
		return nil, fmt.Errorf("upstream get status: %w", err)
	}

	if data, ok := result["data"].(map[string]interface{}); ok {
		return data, nil
	}
	return map[string]interface{}{}, nil
}

// ResetCloudPassword 重置云服务器密码
func (s *DcimCloudService) ResetCloudPassword(cloudID uint, newPass string) error {
	server, err := s.GetServerByID(cloudID)
	if err != nil {
		return fmt.Errorf("cloud server not found: %w", err)
	}

	params := map[string]interface{}{
		"password": newPass,
	}
	if _, err := s.cloudAction(server, "reset_password", params); err != nil {
		return fmt.Errorf("upstream reset password: %w", err)
	}

	// 更新本地密码记录
	updates := map[string]interface{}{
		"password": newPass,
	}
	return s.db.Model(&model.DcimCloudServer{}).Where("id = ?", cloudID).Updates(updates).Error
}

// SaveChartRecord 保存监控数据点（供定时任务或回调使用）
func (s *DcimCloudService) SaveChartRecord(chart *model.CloudChart) error {
	return s.db.Create(chart).Error
}

// ==================== 上游API调用 ====================

// cloudAction 调用魔方云上游API
func (s *DcimCloudService) cloudAction(server *model.DcimCloudServer, action string, params map[string]interface{}) (map[string]interface{}, error) {
	if server.IP == "" {
		return nil, fmt.Errorf("cloud server IP not configured for server %d", server.ID)
	}

	apiURL := fmt.Sprintf("https://%s/api/v1/cloud/%s", server.IP, action)

	form := url.Values{}
	for k, v := range params {
		form.Set(k, fmt.Sprintf("%v", v))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var req *http.Request
	var err error

	if len(params) > 0 {
		req, err = http.NewRequestWithContext(ctx, "POST", apiURL, strings.NewReader(form.Encode()))
		if err != nil {
			return nil, fmt.Errorf("build cloud request: %w", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req, err = http.NewRequestWithContext(ctx, "GET", apiURL, nil)
		if err != nil {
			return nil, fmt.Errorf("build cloud request: %w", err)
		}
	}

	if server.Username != "" {
		req.SetBasicAuth(server.Username, server.Password)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cloud API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("read cloud API response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cloud API HTTP %d: %s", resp.StatusCode, string(body))
	}

	var apiResp struct {
		Result  string          `json:"result"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &apiResp); err != nil {
		// 尝试直接解析为map
		var raw map[string]interface{}
		if err2 := json.Unmarshal(body, &raw); err2 == nil {
			return raw, nil
		}
		return nil, fmt.Errorf("parse cloud API response: %w", err)
	}

	if apiResp.Result != "success" && apiResp.Result != "" {
		return nil, fmt.Errorf("cloud API error: %s", apiResp.Message)
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

// UpdateSecurityGroupRule 更新安全组规则
func (s *DcimCloudService) UpdateSecurityGroupRule(ruleID uint, updates map[string]interface{}) error {
	var rule model.CloudSecurityGroupRule
	if err := s.db.First(&rule, ruleID).Error; err != nil {
		return fmt.Errorf("security group rule not found: %w", err)
	}

	var group model.CloudSecurityGroup
	if err := s.db.First(&group, rule.GroupID).Error; err != nil {
		return fmt.Errorf("security group not found: %w", err)
	}

	server, err := s.GetServerByID(group.CloudID)
	if err != nil {
		return fmt.Errorf("cloud server not found: %w", err)
	}

	upstreamParams := make(map[string]interface{})
	for k, v := range updates {
		upstreamParams[k] = v
	}
	upstreamParams["rule_id"] = ruleID
	if _, err := s.cloudAction(server, "update_security_group_rule", upstreamParams); err != nil {
		return fmt.Errorf("upstream update security group rule: %w", err)
	}

	return s.db.Model(&rule).Updates(updates).Error
}

// UpdateISO 更新ISO信息
func (s *DcimCloudService) UpdateISO(isoID uint, updates map[string]interface{}) error {
	return s.db.Model(&model.CloudISO{}).Where("id = ?", isoID).Updates(updates).Error
}

// CreateISO 创建ISO记录
func (s *DcimCloudService) CreateISO(iso *model.CloudISO) error {
	iso.Status = "available"
	return s.db.Create(iso).Error
}

// DeleteISO 删除ISO记录
func (s *DcimCloudService) DeleteISO(isoID uint) error {
	return s.db.Delete(&model.CloudISO{}, isoID).Error
}

// SyncFlowPacketsFromUpstream 从上游同步流量包数据
func (s *DcimCloudService) SyncFlowPacketsFromUpstream(cloudID uint) error {
	server, err := s.GetServerByID(cloudID)
	if err != nil {
		return fmt.Errorf("cloud server not found: %w", err)
	}

	result, err := s.cloudAction(server, "get_flow_packets", nil)
	if err != nil {
		return fmt.Errorf("upstream get flow packets: %w", err)
	}

	data, ok := result["data"].([]interface{})
	if !ok {
		return nil
	}

	for _, item := range data {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		packet := model.CloudFlowPacket{
			CloudID: cloudID,
		}
		if v, ok := m["name"].(string); ok {
			packet.Name = v
		}
		if v, ok := m["size_gb"].(float64); ok {
			packet.SizeGB = int(v)
		}
		if v, ok := m["used_gb"].(float64); ok {
			packet.UsedGB = int(v)
		}
		if v, ok := m["status"].(string); ok {
			packet.Status = v
		}

		// 更新或创建
		s.db.Where("cloud_id = ? AND name = ?", cloudID, packet.Name).
			Assign(packet).
			FirstOrCreate(&packet)
	}

	return nil
}

// SyncSecurityGroupsFromUpstream 从上游同步安全组数据
func (s *DcimCloudService) SyncSecurityGroupsFromUpstream(cloudID uint) error {
	server, err := s.GetServerByID(cloudID)
	if err != nil {
		return fmt.Errorf("cloud server not found: %w", err)
	}

	result, err := s.cloudAction(server, "get_security_groups", nil)
	if err != nil {
		return fmt.Errorf("upstream get security groups: %w", err)
	}

	data, ok := result["data"].([]interface{})
	if !ok {
		return nil
	}

	for _, item := range data {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		group := model.CloudSecurityGroup{
			CloudID: cloudID,
		}
		if v, ok := m["name"].(string); ok {
			group.Name = v
		}
		if v, ok := m["default_action"].(string); ok {
			group.DefaultAction = v
		}
		if rulesData, ok := m["rules"]; ok {
			if rulesJSON, err := json.Marshal(rulesData); err == nil {
				group.Rules = datatypes.JSON(rulesJSON)
			}
		}

		s.db.Where("cloud_id = ? AND name = ?", cloudID, group.Name).
			Assign(group).
			FirstOrCreate(&group)
	}

	return nil
}
