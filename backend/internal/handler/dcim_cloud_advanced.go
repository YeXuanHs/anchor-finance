package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// ==================== 请求结构体 ====================

type createNATRuleRequest struct {
	Name     string `json:"name" binding:"required,max=256"`
	Protocol string `json:"protocol" binding:"required,oneof=tcp udp"`
	ExtPort  int    `json:"ext_port" binding:"required,min=1,max=65535"`
	IntPort  int    `json:"int_port" binding:"required,min=1,max=65535"`
	IntIP    string `json:"int_ip" binding:"required,ip"`
}

type updateNATRuleRequest struct {
	Name     *string `json:"name"`
	Protocol *string `json:"protocol" binding:"omitempty,oneof=tcp udp"`
	ExtPort  *int    `json:"ext_port" binding:"omitempty,min=1,max=65535"`
	IntPort  *int    `json:"int_port" binding:"omitempty,min=1,max=65535"`
	IntIP    *string `json:"int_ip"`
}

type createSecurityGroupRequest struct {
	Name          string `json:"name" binding:"required,max=256"`
	DefaultAction string `json:"default_action" binding:"required,oneof=accept drop"`
}

type updateSecurityGroupRequest struct {
	Name          *string `json:"name"`
	DefaultAction *string `json:"default_action" binding:"omitempty,oneof=accept drop"`
}

type createSecurityGroupRuleRequest struct {
	Direction string `json:"direction" binding:"required,oneof=in out"`
	Protocol  string `json:"protocol" binding:"required"`
	PortRange string `json:"port_range" binding:"required"`
	Source    string `json:"source" binding:"required"`
	Action    string `json:"action" binding:"required,oneof=accept drop"`
	Priority  int    `json:"priority" binding:"omitempty,gte=0,lte=1000"`
}

type updateSecurityGroupRuleRequest struct {
	Direction *string `json:"direction" binding:"omitempty,oneof=in out"`
	Protocol  *string `json:"protocol"`
	PortRange *string `json:"port_range"`
	Source    *string `json:"source"`
	Action    *string `json:"action" binding:"omitempty,oneof=accept drop"`
	Priority  *int    `json:"priority" binding:"omitempty,gte=0,lte=1000"`
}

type mountISORequest struct {
	ISOID uint `json:"iso_id" binding:"required"`
}

type resetPasswordRequest struct {
	Password string `json:"password" binding:"required,min=6,max=128"`
}

// ==================== NAT管理 ====================

// GetNATRules 获取NAT规则列表
func (h *DcimCloudHandler) GetNATRules(c *gin.Context) {
	cloudID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cloud id")
		return
	}

	rules, err := h.svc.GetNATRules(uint(cloudID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, rules)
}

// CreateNATRule 创建NAT规则
func (h *DcimCloudHandler) CreateNATRule(c *gin.Context) {
	cloudID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cloud id")
		return
	}

	var req createNATRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	rule := &model.CloudNATRule{
		Name:     req.Name,
		Protocol: req.Protocol,
		ExtPort:  req.ExtPort,
		IntPort:  req.IntPort,
		IntIP:    req.IntIP,
	}

	if err := h.svc.CreateNATRule(uint(cloudID), rule); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, rule)
}

// UpdateNATRule 更新NAT规则
func (h *DcimCloudHandler) UpdateNATRule(c *gin.Context) {
	ruleID, err := strconv.ParseUint(c.Param("rule_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid rule id")
		return
	}

	var req updateNATRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Protocol != nil {
		updates["protocol"] = *req.Protocol
	}
	if req.ExtPort != nil {
		updates["ext_port"] = *req.ExtPort
	}
	if req.IntPort != nil {
		updates["int_port"] = *req.IntPort
	}
	if req.IntIP != nil {
		updates["int_ip"] = *req.IntIP
	}

	if len(updates) == 0 {
		response.BadRequest(c, "no fields to update")
		return
	}

	if err := h.svc.UpdateNATRule(uint(ruleID), updates); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "NAT rule updated")
}

// DeleteNATRule 删除NAT规则
func (h *DcimCloudHandler) DeleteNATRule(c *gin.Context) {
	ruleID, err := strconv.ParseUint(c.Param("rule_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid rule id")
		return
	}

	if err := h.svc.DeleteNATRule(uint(ruleID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "NAT rule deleted")
}

// ==================== 安全组管理 ====================

// GetSecurityGroups 获取安全组列表
func (h *DcimCloudHandler) GetSecurityGroups(c *gin.Context) {
	cloudID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cloud id")
		return
	}

	groups, err := h.svc.GetSecurityGroups(uint(cloudID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, groups)
}

// CreateSecurityGroup 创建安全组
func (h *DcimCloudHandler) CreateSecurityGroup(c *gin.Context) {
	cloudID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cloud id")
		return
	}

	var req createSecurityGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	group := &model.CloudSecurityGroup{
		Name:          req.Name,
		DefaultAction: req.DefaultAction,
	}

	if err := h.svc.CreateSecurityGroup(uint(cloudID), group); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, group)
}

// UpdateSecurityGroup 更新安全组
func (h *DcimCloudHandler) UpdateSecurityGroup(c *gin.Context) {
	groupID, err := strconv.ParseUint(c.Param("group_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid group id")
		return
	}

	var req updateSecurityGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.DefaultAction != nil {
		updates["default_action"] = *req.DefaultAction
	}

	if len(updates) == 0 {
		response.BadRequest(c, "no fields to update")
		return
	}

	if err := h.svc.UpdateSecurityGroup(uint(groupID), updates); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "security group updated")
}

// DeleteSecurityGroup 删除安全组
func (h *DcimCloudHandler) DeleteSecurityGroup(c *gin.Context) {
	groupID, err := strconv.ParseUint(c.Param("group_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid group id")
		return
	}

	if err := h.svc.DeleteSecurityGroup(uint(groupID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "security group deleted")
}

// GetSecurityGroupRules 获取安全组规则列表
func (h *DcimCloudHandler) GetSecurityGroupRules(c *gin.Context) {
	groupID, err := strconv.ParseUint(c.Param("group_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid group id")
		return
	}

	rules, err := h.svc.GetSecurityGroupRules(uint(groupID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, rules)
}

// AddSecurityGroupRule 添加安全组规则
func (h *DcimCloudHandler) AddSecurityGroupRule(c *gin.Context) {
	groupID, err := strconv.ParseUint(c.Param("group_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid group id")
		return
	}

	var req createSecurityGroupRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	rule := &model.CloudSecurityGroupRule{
		Direction: req.Direction,
		Protocol:  req.Protocol,
		PortRange: req.PortRange,
		Source:    req.Source,
		Action:    req.Action,
		Priority:  req.Priority,
	}

	if err := h.svc.AddSecurityGroupRule(uint(groupID), rule); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, rule)
}

// UpdateSecurityGroupRule 更新安全组规则
func (h *DcimCloudHandler) UpdateSecurityGroupRule(c *gin.Context) {
	ruleID, err := strconv.ParseUint(c.Param("rule_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid rule id")
		return
	}

	var req updateSecurityGroupRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Direction != nil {
		updates["direction"] = *req.Direction
	}
	if req.Protocol != nil {
		updates["protocol"] = *req.Protocol
	}
	if req.PortRange != nil {
		updates["port_range"] = *req.PortRange
	}
	if req.Source != nil {
		updates["source"] = *req.Source
	}
	if req.Action != nil {
		updates["action"] = *req.Action
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}

	if len(updates) == 0 {
		response.BadRequest(c, "no fields to update")
		return
	}

	if err := h.svc.UpdateSecurityGroupRule(uint(ruleID), updates); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "security group rule updated")
}

// DeleteSecurityGroupRule 删除安全组规则
func (h *DcimCloudHandler) DeleteSecurityGroupRule(c *gin.Context) {
	ruleID, err := strconv.ParseUint(c.Param("rule_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid rule id")
		return
	}

	if err := h.svc.DeleteSecurityGroupRule(uint(ruleID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "security group rule deleted")
}

// ==================== ISO管理 ====================

// GetISOList 获取ISO镜像列表
func (h *DcimCloudHandler) GetISOList(c *gin.Context) {
	isos, err := h.svc.GetISOList()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, isos)
}

// MountISO 挂载ISO
func (h *DcimCloudHandler) MountISO(c *gin.Context) {
	cloudID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cloud id")
		return
	}

	var req mountISORequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.MountISO(uint(cloudID), req.ISOID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "ISO mounted")
}

// UnmountISO 卸载ISO
func (h *DcimCloudHandler) UnmountISO(c *gin.Context) {
	cloudID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cloud id")
		return
	}

	if err := h.svc.UnmountISO(uint(cloudID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "ISO unmounted")
}

// ==================== VNC控制台 ====================

// GetVNCURL 获取VNC控制台URL
func (h *DcimCloudHandler) GetVNCURL(c *gin.Context) {
	cloudID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cloud id")
		return
	}

	vncURL, err := h.svc.GetVNCURL(uint(cloudID))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"url": vncURL})
}

// GetVNCPage 获取VNC Web页面数据
func (h *DcimCloudHandler) GetVNCPage(c *gin.Context) {
	cloudID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cloud id")
		return
	}

	data, err := h.svc.GetVNCPage(uint(cloudID))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

// ==================== 监控图表 ====================

// GetResourceChart 获取资源监控图表
func (h *DcimCloudHandler) GetResourceChart(c *gin.Context) {
	cloudID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cloud id")
		return
	}

	period := c.DefaultQuery("period", "24h")

	charts, err := h.svc.GetResourceChart(uint(cloudID), period)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, charts)
}

// GetResourceInfo 获取当前资源使用情况
func (h *DcimCloudHandler) GetResourceInfo(c *gin.Context) {
	cloudID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cloud id")
		return
	}

	info, err := h.svc.GetResourceInfo(uint(cloudID))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, info)
}

// ==================== 流量包管理 ====================

// GetFlowPackets 获取流量包列表
func (h *DcimCloudHandler) GetFlowPackets(c *gin.Context) {
	cloudID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cloud id")
		return
	}

	packets, err := h.svc.GetFlowPackets(uint(cloudID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, packets)
}

// BuyFlowPacket 购买流量包
func (h *DcimCloudHandler) BuyFlowPacket(c *gin.Context) {
	cloudID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cloud id")
		return
	}

	var req struct {
		PacketID uint `json:"packet_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.BuyFlowPacket(uint(cloudID), req.PacketID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "flow packet purchased")
}

// GetFlowPacketUsage 获取流量使用情况
func (h *DcimCloudHandler) GetFlowPacketUsage(c *gin.Context) {
	cloudID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cloud id")
		return
	}

	usage, err := h.svc.GetFlowPacketUsage(uint(cloudID))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, usage)
}

// ==================== 附加操作 ====================

// GetCloudStatus 获取云服务器实时状态
func (h *DcimCloudHandler) GetCloudStatus(c *gin.Context) {
	cloudID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cloud id")
		return
	}

	status, err := h.svc.GetCloudStatus(uint(cloudID))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, status)
}

// ResetCloudPassword 重置云服务器密码
func (h *DcimCloudHandler) ResetCloudPassword(c *gin.Context) {
	cloudID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cloud id")
		return
	}

	var req resetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.ResetCloudPassword(uint(cloudID), req.Password); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "password reset")
}
