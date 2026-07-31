package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// DcimAdvancedHandler DCIM高级操作handler
type DcimAdvancedHandler struct {
	dcimSvc *service.DcimService
	log     *logger.Logger
}

func NewDcimAdvancedHandler(dcimSvc *service.DcimService, log *logger.Logger) *DcimAdvancedHandler {
	return &DcimAdvancedHandler{dcimSvc: dcimSvc, log: log}
}

// ==================== 请求结构体 ====================

type rescueRequest struct {
	OS string `json:"os" binding:"required"`
}

type snapshotRequest struct {
	Name string `json:"name" binding:"required,max=256"`
}

type backupRequest struct {
	Name string `json:"name" binding:"required,max=256"`
	Type string `json:"type" binding:"omitempty,oneof=manual auto"`
}

// ==================== KVM/IPMI/BMC ====================

// GetKVMURL 获取KVM控制台URL
func (h *DcimAdvancedHandler) GetKVMURL(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server id")
		return
	}

	result, err := h.dcimSvc.GetKVMURL(uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// GetBMCInfo 获取BMC连接信息
func (h *DcimAdvancedHandler) GetBMCInfo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server id")
		return
	}

	result, err := h.dcimSvc.GetBMCInfo(uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// GetNoVNCURL 获取noVNC Web控制台URL
func (h *DcimAdvancedHandler) GetNoVNCURL(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server id")
		return
	}

	result, err := h.dcimSvc.GetNoVNCURL(uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// ==================== 救援系统 ====================

// BootRescue 启动救援模式
func (h *DcimAdvancedHandler) BootRescue(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server id")
		return
	}

	var req rescueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	operatorID := c.GetUint("user_id")
	rescueLog, err := h.dcimSvc.BootRescue(uint(id), req.OS, operatorID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, rescueLog)
}

// CrackPassword 重置密码
func (h *DcimAdvancedHandler) CrackPassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server id")
		return
	}

	operatorID := c.GetUint("user_id")
	rescueLog, err := h.dcimSvc.CrackPassword(uint(id), operatorID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, rescueLog)
}

// GetRescueStatus 获取救援模式状态
func (h *DcimAdvancedHandler) GetRescueStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server id")
		return
	}

	result, err := h.dcimSvc.GetRescueStatus(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, result)
}

// ==================== 流量监控 ====================

// GetTrafficUsage 获取流量使用情况
func (h *DcimAdvancedHandler) GetTrafficUsage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server id")
		return
	}

	result, err := h.dcimSvc.GetTrafficUsage(uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// GetTrafficChart 获取流量图表数据
func (h *DcimAdvancedHandler) GetTrafficChart(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server id")
		return
	}

	period := c.DefaultQuery("period", "30d")

	logs, err := h.dcimSvc.GetTrafficChart(uint(id), period)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, logs)
}

// ResetTrafficCounter 重置流量计数器
func (h *DcimAdvancedHandler) ResetTrafficCounter(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server id")
		return
	}

	if err := h.dcimSvc.ResetTrafficCounter(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "traffic counter reset")
}

// ==================== 快照管理 ====================

// CreateSnapshot 创建快照
func (h *DcimAdvancedHandler) CreateSnapshot(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server id")
		return
	}

	var req snapshotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	snapshot, err := h.dcimSvc.CreateSnapshot(uint(id), req.Name)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, snapshot)
}

// RestoreSnapshot 从快照恢复
func (h *DcimAdvancedHandler) RestoreSnapshot(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid snapshot id")
		return
	}

	operatorID := c.GetUint("user_id")
	if err := h.dcimSvc.RestoreSnapshot(uint(id), operatorID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "snapshot restore initiated")
}

// DeleteSnapshot 删除快照
func (h *DcimAdvancedHandler) DeleteSnapshot(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid snapshot id")
		return
	}

	if err := h.dcimSvc.DeleteSnapshot(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "snapshot deleted")
}

// GetSnapshots 获取快照列表
func (h *DcimAdvancedHandler) GetSnapshots(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server id")
		return
	}

	snapshots, err := h.dcimSvc.GetSnapshots(uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, snapshots)
}

// ==================== 备份管理 ====================

// CreateBackup 创建备份
func (h *DcimAdvancedHandler) CreateBackup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server id")
		return
	}

	var req backupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	backup, err := h.dcimSvc.CreateBackup(uint(id), req.Name, req.Type)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, backup)
}

// RestoreBackup 从备份恢复
func (h *DcimAdvancedHandler) RestoreBackup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid backup id")
		return
	}

	operatorID := c.GetUint("user_id")
	if err := h.dcimSvc.RestoreBackup(uint(id), operatorID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "backup restore initiated")
}

// DeleteBackup 删除备份
func (h *DcimAdvancedHandler) DeleteBackup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid backup id")
		return
	}

	if err := h.dcimSvc.DeleteBackup(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "backup deleted")
}

// GetBackups 获取备份列表
func (h *DcimAdvancedHandler) GetBackups(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server id")
		return
	}

	backups, err := h.dcimSvc.GetBackups(uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, backups)
}

// ==================== 电源管理增强 ====================

// GetPowerStatus 获取实时电源状态
func (h *DcimAdvancedHandler) GetPowerStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server id")
		return
	}

	result, err := h.dcimSvc.GetPowerStatus(uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// RefreshPowerStatus 强制刷新电源状态
func (h *DcimAdvancedHandler) RefreshPowerStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server id")
		return
	}

	result, err := h.dcimSvc.RefreshPowerStatus(uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// ==================== 重装增强 ====================

// GetReinstallStatus 获取重装进度
func (h *DcimAdvancedHandler) GetReinstallStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server id")
		return
	}

	result, err := h.dcimSvc.GetReinstallStatus(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, result)
}

// CancelReinstall 取消重装
func (h *DcimAdvancedHandler) CancelReinstall(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server id")
		return
	}

	operatorID := c.GetUint("user_id")
	if err := h.dcimSvc.CancelReinstall(uint(id), operatorID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "reinstall cancelled")
}

// GetOSList 获取可用操作系统列表
func (h *DcimAdvancedHandler) GetOSList(c *gin.Context) {
	osList, err := h.dcimSvc.GetOSList()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, osList)
}

// ==================== 救护日志 ====================

// GetRescueLogs 获取救援操作日志
func (h *DcimAdvancedHandler) GetRescueLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var serverID uint
	if sid := c.Query("server_id"); sid != "" {
		v, err := strconv.ParseUint(sid, 10, 64)
		if err == nil {
			serverID = uint(v)
		}
	}

	logs, total, err := h.dcimSvc.GetRescueLogs(serverID, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, logs, total, page, pageSize)
}
