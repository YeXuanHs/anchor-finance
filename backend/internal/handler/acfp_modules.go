package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type ACFPHandler struct {
	svc *service.ACFPService
	log *logger.Logger
}

func NewACFPHandler(svc *service.ACFPService, log *logger.Logger) *ACFPHandler {
	return &ACFPHandler{svc: svc, log: log}
}

// ─── IP历史 ───

func (h *ACFPHandler) GetIPHistory(c *gin.Context) {
	hostID, _ := strconv.ParseUint(c.Param("host_id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 20 }
	items, total, err := h.svc.GetIPHistory(uint(hostID), page, pageSize)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}
	response.Success(c, gin.H{"items": items, "total": total, "page": page})
}

// ─── 限量发售 ───

func (h *ACFPHandler) ListLimitedSales(c *gin.Context) {
	items, err := h.svc.ListLimitedSales()
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}
	response.Success(c, items)
}

func (h *ACFPHandler) SetLimitedSale(c *gin.Context) {
	var item model.ACFPLimitedSale
	if err := c.ShouldBindJSON(&item); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.svc.SetLimitedSale(&item); err != nil {
		response.ServerError(c, "保存失败")
		return
	}
	response.SuccessMsg(c, "保存成功")
}

func (h *ACFPHandler) CheckStock(c *gin.Context) {
	productID, _ := strconv.ParseUint(c.Param("product_id"), 10, 64)
	qty, _ := strconv.Atoi(c.DefaultQuery("qty", "1"))
	ok, remaining, _ := h.svc.CheckStock(uint(productID), qty)
	response.Success(c, gin.H{"available": ok, "remaining": remaining})
}

// ─── 价格锁定 ───

func (h *ACFPHandler) ListPriceLocks(c *gin.Context) {
	items, err := h.svc.ListPriceLocks()
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}
	response.Success(c, items)
}

func (h *ACFPHandler) SetPriceLock(c *gin.Context) {
	var item model.ACFPPriceLock
	if err := c.ShouldBindJSON(&item); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.svc.SetPriceLock(&item); err != nil {
		response.ServerError(c, "保存失败")
		return
	}
	response.SuccessMsg(c, "保存成功")
}

func (h *ACFPHandler) DeletePriceLock(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeletePriceLock(uint(id)); err != nil {
		response.ServerError(c, "删除失败")
		return
	}
	response.SuccessMsg(c, "删除成功")
}

// ─── 操作日志 ───

func (h *ACFPHandler) ListLogs(c *gin.Context) {
	module := c.Query("module")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 20 }
	items, total, err := h.svc.ListLogs(module, page, pageSize)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}
	response.Success(c, gin.H{"items": items, "total": total, "page": page})
}

func (h *ACFPHandler) CleanLogs(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "90"))
	if days < 1 { days = 90 }
	count := h.svc.CleanLogs(days)
	response.Success(c, gin.H{"cleaned": count})
}

// ─── 定时任务状态 ───

func (h *ACFPHandler) GetCronStatuses(c *gin.Context) {
	items, err := h.svc.GetCronStatuses()
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}
	response.Success(c, items)
}

// ─── 实名认证Pro ───

func (h *ACFPHandler) GetCertProConfig(c *gin.Context) {
	cfg := h.svc.GetCertProConfig()
	response.Success(c, cfg)
}

func (h *ACFPHandler) SetCertProConfig(c *gin.Context) {
	var req struct {
		MinAge int `json:"min_age"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if req.MinAge < 1 || req.MinAge > 100 {
		response.BadRequest(c, "年龄范围无效")
		return
	}
	if err := h.svc.SetCertProConfig(req.MinAge); err != nil {
		response.ServerError(c, "保存失败")
		return
	}
	response.SuccessMsg(c, "保存成功")
}

func (h *ACFPHandler) ListMinorCerts(c *gin.Context) {
	items, err := h.svc.ListMinorCerts()
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}
	response.Success(c, items)
}

func (h *ACFPHandler) ScanMinorCerts(c *gin.Context) {
	count, err := h.svc.ScanMinorCerts()
	if err != nil {
		response.ServerError(c, "扫描失败")
		return
	}
	response.Success(c, gin.H{"found": count})
}

// ─── 缓存预热 ───

func (h *ACFPHandler) WarmCache(c *gin.Context) {
	count, err := h.svc.WarmProductCache()
	if err != nil {
		response.ServerError(c, "预热失败")
		return
	}
	response.Success(c, gin.H{"warmed": count})
}

// ─── 批量商品修改 ───

func (h *ACFPHandler) CreateBatchTask(c *gin.Context) {
	var task model.ACFPBatchTask
	if err := c.ShouldBindJSON(&task); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.svc.CreateBatchTask(&task); err != nil {
		response.ServerError(c, "创建失败")
		return
	}
	response.Success(c, task)
}

func (h *ACFPHandler) ListBatchTasks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 20 }
	items, total, err := h.svc.ListBatchTasks(page, pageSize)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}
	response.Success(c, gin.H{"items": items, "total": total, "page": page})
}

func (h *ACFPHandler) ExecuteBatchTask(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.ExecuteBatchTask(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "执行中")
}

// ─── 状态对账 ───

func (h *ACFPHandler) RunStatusSync(c *gin.Context) {
	count, err := h.svc.RunStatusSync()
	if err != nil {
		response.ServerError(c, "对账失败")
		return
	}
	response.Success(c, gin.H{"synced": count})
}

// ─── 通知去重 ───

func (h *ACFPHandler) GetNotifyStats(c *gin.Context) {
	response.Success(c, gin.H{"message": "通知去重运行中"})
}

func (h *ACFPHandler) CleanNotifyEvents(c *gin.Context) {
	count := h.svc.CleanupOldNotifyEvents()
	response.Success(c, gin.H{"cleaned": count})
}
