package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type FriendlyLinkHandler struct {
	svc *service.FriendlyLinkService
	log *logger.Logger
}

func NewFriendlyLinkHandler(svc *service.FriendlyLinkService, log *logger.Logger) *FriendlyLinkHandler {
	return &FriendlyLinkHandler{svc: svc, log: log}
}

func (h *FriendlyLinkHandler) AdminGetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status, _ := strconv.Atoi(c.DefaultQuery("status", "-1"))
	group := c.Query("group")
	items, total, err := h.svc.GetList(page, pageSize, status, group)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, items, total, page, pageSize)
}

func (h *FriendlyLinkHandler) AdminCreate(c *gin.Context) {
	var link model.FriendlyLink
	if err := c.ShouldBindJSON(&link); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.Create(&link); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, link)
}

func (h *FriendlyLinkHandler) AdminUpdate(c *gin.Context) {
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
	if err := h.svc.Update(uint(id), updates); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "updated")
}

func (h *FriendlyLinkHandler) AdminDelete(c *gin.Context) {
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

func (h *FriendlyLinkHandler) GetActive(c *gin.Context) {
	group := c.DefaultQuery("group", "default")
	items, err := h.svc.GetActive(group)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, items)
}
