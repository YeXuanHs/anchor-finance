package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type AnnounceHandler struct {
	svc *service.AnnounceService
	log *logger.Logger
}

func NewAnnounceHandler(svc *service.AnnounceService, log *logger.Logger) *AnnounceHandler {
	return &AnnounceHandler{svc: svc, log: log}
}

func (h *AnnounceHandler) AdminGetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status, _ := strconv.Atoi(c.DefaultQuery("status", "-1"))
	announceType := c.Query("type")
	items, total, err := h.svc.GetList(page, pageSize, status, announceType)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, items, total, page, pageSize)
}

func (h *AnnounceHandler) AdminGetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	item, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "not found")
		return
	}
	response.Success(c, item)
}

func (h *AnnounceHandler) AdminCreate(c *gin.Context) {
	var req service.CreateAnnounceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	item, err := h.svc.Create(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, item)
}

func (h *AnnounceHandler) AdminUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	item, err := h.svc.Update(uint(id), updates)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, item)
}

func (h *AnnounceHandler) AdminDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	if err := h.svc.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "deleted")
}

func (h *AnnounceHandler) GetActive(c *gin.Context) {
	items, err := h.svc.GetActive()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, items)
}
