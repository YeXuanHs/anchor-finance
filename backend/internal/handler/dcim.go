package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type DcimHandler struct {
	dcimSvc *service.DcimService
	log     *logger.Logger
}

func NewDcimHandler(dcimSvc *service.DcimService, log *logger.Logger) *DcimHandler {
	return &DcimHandler{dcimSvc: dcimSvc, log: log}
}

// ==================== 请求结构体 ====================

type createPhysicalServerRequest struct {
	Name          string  `json:"name" binding:"required,max=128"`
	Hostname      string  `json:"hostname" binding:"omitempty,max=128"`
	IP            string  `json:"ip" binding:"required,ip"`
	IPv6          string  `json:"ipv6" binding:"omitempty"`
	MAC           string  `json:"mac" binding:"omitempty"`
	DatacenterID  uint    `json:"datacenter_id" binding:"required"`
	Rack          string  `json:"rack" binding:"omitempty"`
	RackPosition  string  `json:"rack_position" binding:"omitempty"`
	CPU           string  `json:"cpu" binding:"omitempty"`
	CPUCores      int     `json:"cpu_cores" binding:"omitempty,gte=0"`
	MemoryMB      int     `json:"memory_mb" binding:"omitempty,gte=0"`
	DiskType      string  `json:"disk_type" binding:"omitempty,oneof=SSD HDD NVMe"`
	DiskSizeGB    int     `json:"disk_size_gb" binding:"omitempty,gte=0"`
	BandwidthMbps int     `json:"bandwidth_mbps" binding:"omitempty,gte=0"`
	TrafficGB     int     `json:"traffic_gb" binding:"omitempty,gte=0"`
	OS            string  `json:"os" binding:"omitempty"`
	Remark        string  `json:"remark" binding:"omitempty"`
	Tags          string  `json:"tags" binding:"omitempty"`
}

type updatePhysicalServerRequest struct {
	Name          *string `json:"name" binding:"omitempty,max=128"`
	Hostname      *string `json:"hostname"`
	IPv6          *string `json:"ipv6"`
	MAC           *string `json:"mac"`
	DatacenterID  *uint   `json:"datacenter_id"`
	Rack          *string `json:"rack"`
	RackPosition  *string `json:"rack_position"`
	CPU           *string `json:"cpu"`
	CPUCores      *int    `json:"cpu_cores" binding:"omitempty,gte=0"`
	MemoryMB      *int    `json:"memory_mb" binding:"omitempty,gte=0"`
	DiskType      *string `json:"disk_type" binding:"omitempty,oneof=SSD HDD NVMe"`
	DiskSizeGB    *int    `json:"disk_size_gb" binding:"omitempty,gte=0"`
	BandwidthMbps *int    `json:"bandwidth_mbps" binding:"omitempty,gte=0"`
	TrafficGB     *int    `json:"traffic_gb" binding:"omitempty,gte=0"`
	OS            *string `json:"os"`
	Remark        *string `json:"remark"`
	Tags          *string `json:"tags"`
	OwnerID       *uint   `json:"owner_id"`
	ExpiredAt     *string `json:"expired_at"`
}

type createCloudServerRequest struct {
	Name           string  `json:"name" binding:"required,max=128"`
	Hostname       string  `json:"hostname" binding:"omitempty,max=128"`
	IP             string  `json:"ip" binding:"required,ip"`
	IPv6           string  `json:"ipv6" binding:"omitempty"`
	DatacenterID   uint    `json:"datacenter_id" binding:"required"`
	CPU            int     `json:"cpu" binding:"required,gte=1"`
	MemoryMB       int     `json:"memory_mb" binding:"required,gte=1"`
	DiskSizeGB     int     `json:"disk_size_gb" binding:"required,gte=1"`
	BandwidthMbps  int     `json:"bandwidth_mbps" binding:"omitempty,gte=0"`
	TrafficGB      int     `json:"traffic_gb" binding:"omitempty,gte=0"`
	OS             string  `json:"os" binding:"omitempty"`
	VirtualType    string  `json:"virtual_type" binding:"omitempty,oneof=KVM VMware LXC"`
	ParentServerID *uint   `json:"parent_server_id"`
	PriceMonthly   float64 `json:"price_monthly" binding:"omitempty,gte=0"`
	Remark         string  `json:"remark" binding:"omitempty"`
}

type serverActionRequest struct {
	Force bool `json:"force"`
}

type reinstallRequest struct {
	OS string `json:"os" binding:"required"`
}

type renewRequest struct {
	Months int `json:"months" binding:"required,gte=1,lte=120"`
}

// ==================== 物理服务器 ====================

// GetServerList 获取物理服务器列表
func (h *DcimHandler) GetServerList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	var status *int8
	if s := c.Query("status"); s != "" {
		v, err := strconv.ParseInt(s, 10, 8)
		if err == nil {
			st := int8(v)
			status = &st
		}
	}

	var dcID *uint
	if d := c.Query("datacenter_id"); d != "" {
		v, err := strconv.ParseUint(d, 10, 64)
		if err == nil {
			id := uint(v)
			dcID = &id
		}
	}

	servers, total, err := h.dcimSvc.GetServerList(page, pageSize, keyword, status, dcID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, servers, total, page, pageSize)
}

// GetServerDetail 获取物理服务器详情
func (h *DcimHandler) GetServerDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server id")
		return
	}

	server, err := h.dcimSvc.GetServerByID(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, server)
}

// CreateServer 创建物理服务器
func (h *DcimHandler) CreateServer(c *gin.Context) {
	var req createPhysicalServerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	server := &model.DcimServer{
		Name:          req.Name,
		Hostname:      req.Hostname,
		IP:            req.IP,
		IPv6:          req.IPv6,
		MAC:           req.MAC,
		DatacenterID:  req.DatacenterID,
		Rack:          req.Rack,
		RackPosition:  req.RackPosition,
		CPU:           req.CPU,
		CPUCores:      req.CPUCores,
		MemoryMB:      req.MemoryMB,
		DiskType:      req.DiskType,
		DiskSizeGB:    req.DiskSizeGB,
		BandwidthMbps: req.BandwidthMbps,
		TrafficGB:     req.TrafficGB,
		OS:            req.OS,
		Remark:        req.Remark,
		Tags:          req.Tags,
	}

	if err := h.dcimSvc.CreateServer(server); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, server)
}

// UpdateServer 更新物理服务器
func (h *DcimHandler) UpdateServer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server id")
		return
	}

	var req updatePhysicalServerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Hostname != nil {
		updates["hostname"] = *req.Hostname
	}
	if req.IPv6 != nil {
		updates["ipv6"] = *req.IPv6
	}
	if req.MAC != nil {
		updates["mac"] = *req.MAC
	}
	if req.DatacenterID != nil {
		updates["datacenter_id"] = *req.DatacenterID
	}
	if req.Rack != nil {
		updates["rack"] = *req.Rack
	}
	if req.RackPosition != nil {
		updates["rack_position"] = *req.RackPosition
	}
	if req.CPU != nil {
		updates["cpu"] = *req.CPU
	}
	if req.CPUCores != nil {
		updates["cpu_cores"] = *req.CPUCores
	}
	if req.MemoryMB != nil {
		updates["memory_mb"] = *req.MemoryMB
	}
	if req.DiskType != nil {
		updates["disk_type"] = *req.DiskType
	}
	if req.DiskSizeGB != nil {
		updates["disk_size_gb"] = *req.DiskSizeGB
	}
	if req.BandwidthMbps != nil {
		updates["bandwidth_mbps"] = *req.BandwidthMbps
	}
	if req.TrafficGB != nil {
		updates["traffic_gb"] = *req.TrafficGB
	}
	if req.OS != nil {
		updates["os"] = *req.OS
	}
	if req.Remark != nil {
		updates["remark"] = *req.Remark
	}
	if req.Tags != nil {
		updates["tags"] = *req.Tags
	}
	if req.OwnerID != nil {
		updates["owner_id"] = *req.OwnerID
	}
	if req.ExpiredAt != nil {
		updates["expired_at"] = *req.ExpiredAt
	}

	if len(updates) == 0 {
		response.BadRequest(c, "no fields to update")
		return
	}

	if err := h.dcimSvc.UpdateServer(uint(id), updates); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "server updated")
}

// DeleteServer 删除物理服务器
func (h *DcimHandler) DeleteServer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server id")
		return
	}

	if err := h.dcimSvc.DeleteServer(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "server deleted")
}

// BootServer 物理服务器开机
func (h *DcimHandler) BootServer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server id")
		return
	}

	operatorID := c.GetUint("user_id")
	if err := h.dcimSvc.BootServer(uint(id), operatorID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "server boot initiated")
}

// ShutdownServer 物理服务器关机
func (h *DcimHandler) ShutdownServer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server id")
		return
	}

	var req serverActionRequest
	c.ShouldBindJSON(&req)

	operatorID := c.GetUint("user_id")
	if err := h.dcimSvc.ShutdownServer(uint(id), req.Force, operatorID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "server shutdown initiated")
}

// RebootServer 物理服务器重启
func (h *DcimHandler) RebootServer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server id")
		return
	}

	operatorID := c.GetUint("user_id")
	if err := h.dcimSvc.RebootServer(uint(id), operatorID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "server reboot initiated")
}

// ReinstallServer 物理服务器重装系统
func (h *DcimHandler) ReinstallServer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server id")
		return
	}

	var req reinstallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	operatorID := c.GetUint("user_id")
	if err := h.dcimSvc.ReinstallServer(uint(id), req.OS, operatorID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "server reinstall initiated")
}

// RenewServer 物理服务器续费
func (h *DcimHandler) RenewServer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server id")
		return
	}

	var req renewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	operatorID := c.GetUint("user_id")
	if err := h.dcimSvc.RenewServer(uint(id), req.Months, operatorID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "server renewed")
}

// ==================== 云服务器 ====================

// GetCloudList 获取云服务器列表
func (h *DcimHandler) GetCloudList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	var status *int8
	if s := c.Query("status"); s != "" {
		v, err := strconv.ParseInt(s, 10, 8)
		if err == nil {
			st := int8(v)
			status = &st
		}
	}

	var ownerID *uint
	if o := c.Query("owner_id"); o != "" {
		v, err := strconv.ParseUint(o, 10, 64)
		if err == nil {
			id := uint(v)
			ownerID = &id
		}
	}

	clouds, total, err := h.dcimSvc.GetCloudList(page, pageSize, keyword, status, ownerID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, clouds, total, page, pageSize)
}

// GetCloudDetail 获取云服务器详情
func (h *DcimHandler) GetCloudDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cloud server id")
		return
	}

	cloud, err := h.dcimSvc.GetCloudByID(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, cloud)
}

// CreateCloud 创建云服务器
func (h *DcimHandler) CreateCloud(c *gin.Context) {
	var req createCloudServerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	cloud := &model.DcimCloud{
		Name:           req.Name,
		Hostname:       req.Hostname,
		IP:             req.IP,
		IPv6:           req.IPv6,
		DatacenterID:   req.DatacenterID,
		CPU:            req.CPU,
		MemoryMB:       req.MemoryMB,
		DiskSizeGB:     req.DiskSizeGB,
		BandwidthMbps:  req.BandwidthMbps,
		TrafficGB:      req.TrafficGB,
		OS:             req.OS,
		VirtualType:    req.VirtualType,
		ParentServerID: req.ParentServerID,
		PriceMonthly:   req.PriceMonthly,
		Remark:         req.Remark,
	}

	if err := h.dcimSvc.CreateCloud(cloud); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, cloud)
}

// ProvisionCloud 开通云服务器
func (h *DcimHandler) ProvisionCloud(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cloud server id")
		return
	}

	operatorID := c.GetUint("user_id")
	if err := h.dcimSvc.ProvisionCloud(uint(id), operatorID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "cloud server provisioning started")
}

// BootCloud 云服务器开机
func (h *DcimHandler) BootCloud(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cloud server id")
		return
	}

	operatorID := c.GetUint("user_id")
	if err := h.dcimSvc.BootCloud(uint(id), operatorID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "cloud server boot initiated")
}

// ShutdownCloud 云服务器关机
func (h *DcimHandler) ShutdownCloud(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cloud server id")
		return
	}

	var req serverActionRequest
	c.ShouldBindJSON(&req)

	operatorID := c.GetUint("user_id")
	if err := h.dcimSvc.ShutdownCloud(uint(id), req.Force, operatorID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "cloud server shutdown initiated")
}

// RebootCloud 云服务器重启
func (h *DcimHandler) RebootCloud(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cloud server id")
		return
	}

	operatorID := c.GetUint("user_id")
	if err := h.dcimSvc.RebootCloud(uint(id), operatorID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "cloud server reboot initiated")
}

// ReinstallCloud 云服务器重装系统
func (h *DcimHandler) ReinstallCloud(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cloud server id")
		return
	}

	var req reinstallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	operatorID := c.GetUint("user_id")
	if err := h.dcimSvc.ReinstallCloud(uint(id), req.OS, operatorID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "cloud server reinstall initiated")
}

// RenewCloud 云服务器续费
func (h *DcimHandler) RenewCloud(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cloud server id")
		return
	}

	var req renewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	operatorID := c.GetUint("user_id")
	if err := h.dcimSvc.RenewCloud(uint(id), req.Months, operatorID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "cloud server renewed")
}

// ==================== 通用 ====================

// GetDatacenters 获取机房列表
func (h *DcimHandler) GetDatacenters(c *gin.Context) {
	dcs, err := h.dcimSvc.GetDatacenterList()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, dcs)
}

// GetOperationLogs 获取操作日志
func (h *DcimHandler) GetOperationLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	serverType := c.Query("server_type")

	var serverID uint
	if sid := c.Query("server_id"); sid != "" {
		v, err := strconv.ParseUint(sid, 10, 64)
		if err == nil {
			serverID = uint(v)
		}
	}

	logs, total, err := h.dcimSvc.GetOperationLogs(serverType, serverID, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, logs, total, page, pageSize)
}

// GetDatacenterList is an alias for GetDatacenters.
func (h *DcimHandler) GetDatacenterList(c *gin.Context) { h.GetDatacenters(c) }

// GetStats 获取DCIM统计信息
func (h *DcimHandler) GetStats(c *gin.Context) {
	stats, err := h.dcimSvc.GetStats()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, stats)
}
