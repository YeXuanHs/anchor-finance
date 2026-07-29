package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type DcimCloudHandler struct {
	svc *service.DcimCloudService
	log *logger.Logger
}

func NewDcimCloudHandler(svc *service.DcimCloudService, log *logger.Logger) *DcimCloudHandler {
	return &DcimCloudHandler{svc: svc, log: log}
}

// GetServers returns a list of DCIM cloud servers.
func (h *DcimCloudHandler) GetServers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	var status *int8
	if s := c.Query("status"); s != "" {
		v, _ := strconv.ParseInt(s, 10, 8)
		st := int8(v)
		status = &st
	}

	servers, total, err := h.svc.GetServerList(page, pageSize, keyword, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, servers, total, page, pageSize)
}

// GetServer returns a single DCIM cloud server.
func (h *DcimCloudHandler) GetServer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server id")
		return
	}

	server, err := h.svc.GetServerByID(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, server)
}

// CreateServer creates a new DCIM cloud server.
func (h *DcimCloudHandler) CreateServer(c *gin.Context) {
	var req struct {
		Name          string  `json:"name" binding:"required"`
		Hostname      string  `json:"hostname"`
		IP            string  `json:"ip" binding:"required"`
		Username      string  `json:"username"`
		Password      string  `json:"password"`
		Secure        int8    `json:"secure"`
		Disabled      int8    `json:"disabled"`
		UserPrefix    string  `json:"user_prefix"`
		AccountType   string  `json:"account_type"`
		DatacenterID  uint    `json:"datacenter_id"`
		CPU           int     `json:"cpu"`
		MemoryMB      int     `json:"memory_mb"`
		DiskSizeGB    int     `json:"disk_size_gb"`
		BandwidthMbps int     `json:"bandwidth_mbps"`
		TrafficGB     int     `json:"traffic_gb"`
		OS            string  `json:"os"`
		VirtualType   string  `json:"virtual_type"`
		PriceMonthly  float64 `json:"price_monthly"`
		Remark        string  `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	server := &model.DcimCloudServer{
		Name:          req.Name,
		Hostname:      req.Hostname,
		IP:            req.IP,
		Username:      req.Username,
		Password:      req.Password,
		Secure:        req.Secure,
		Disabled:      req.Disabled,
		UserPrefix:    req.UserPrefix,
		AccountType:   req.AccountType,
		DatacenterID:  req.DatacenterID,
		CPU:           req.CPU,
		MemoryMB:      req.MemoryMB,
		DiskSizeGB:    req.DiskSizeGB,
		BandwidthMbps: req.BandwidthMbps,
		TrafficGB:     req.TrafficGB,
		OS:            req.OS,
		VirtualType:   req.VirtualType,
		PriceMonthly:  req.PriceMonthly,
		Remark:        req.Remark,
		Status:        1,
	}

	if err := h.svc.CreateServer(server); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, server)
}

// UpdateServer updates an existing DCIM cloud server.
func (h *DcimCloudHandler) UpdateServer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server id")
		return
	}

	var req struct {
		Name          *string  `json:"name"`
		Hostname      *string  `json:"hostname"`
		IP            *string  `json:"ip"`
		Username      *string  `json:"username"`
		Password      *string  `json:"password"`
		Secure        *int8    `json:"secure"`
		Disabled      *int8    `json:"disabled"`
		UserPrefix    *string  `json:"user_prefix"`
		AccountType   *string  `json:"account_type"`
		DatacenterID  *uint    `json:"datacenter_id"`
		CPU           *int     `json:"cpu"`
		MemoryMB      *int     `json:"memory_mb"`
		DiskSizeGB    *int     `json:"disk_size_gb"`
		BandwidthMbps *int     `json:"bandwidth_mbps"`
		TrafficGB     *int     `json:"traffic_gb"`
		OS            *string  `json:"os"`
		VirtualType   *string  `json:"virtual_type"`
		PriceMonthly  *float64 `json:"price_monthly"`
		Remark        *string  `json:"remark"`
	}
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
	if req.IP != nil {
		updates["ip"] = *req.IP
	}
	if req.Username != nil {
		updates["username"] = *req.Username
	}
	if req.Password != nil {
		updates["password"] = *req.Password
	}
	if req.Secure != nil {
		updates["secure"] = *req.Secure
	}
	if req.Disabled != nil {
		updates["disabled"] = *req.Disabled
	}
	if req.UserPrefix != nil {
		updates["user_prefix"] = *req.UserPrefix
	}
	if req.AccountType != nil {
		updates["account_type"] = *req.AccountType
	}
	if req.DatacenterID != nil {
		updates["datacenter_id"] = *req.DatacenterID
	}
	if req.CPU != nil {
		updates["cpu"] = *req.CPU
	}
	if req.MemoryMB != nil {
		updates["memory_mb"] = *req.MemoryMB
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
	if req.VirtualType != nil {
		updates["virtual_type"] = *req.VirtualType
	}
	if req.PriceMonthly != nil {
		updates["price_monthly"] = *req.PriceMonthly
	}
	if req.Remark != nil {
		updates["remark"] = *req.Remark
	}

	if err := h.svc.UpdateServer(uint(id), updates); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "server updated")
}

// DeleteServer deletes a DCIM cloud server.
func (h *DcimCloudHandler) DeleteServer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server id")
		return
	}

	if err := h.svc.DeleteServer(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "server deleted")
}

// TestConnection tests connection to a DCIM cloud server.
func (h *DcimCloudHandler) TestConnection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server id")
		return
	}

	if err := h.svc.TestConnection(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "connection successful"})
}

// SyncServer syncs server info from remote.
func (h *DcimCloudHandler) SyncServer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server id")
		return
	}

	if err := h.svc.SyncServer(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "server synced")
}

// GetOperationLogs returns operation logs.
func (h *DcimCloudHandler) GetOperationLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var serverID uint
	if sid := c.Query("server_id"); sid != "" {
		v, _ := strconv.ParseUint(sid, 10, 64)
		serverID = uint(v)
	}

	logs, total, err := h.svc.GetOperationLogs(serverID, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, logs, total, page, pageSize)
}
