package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type RbacPageHandler struct {
	svc *service.RbacPageService
	log *logger.Logger
}

func NewRbacPageHandler(svc *service.RbacPageService, log *logger.Logger) *RbacPageHandler {
	return &RbacPageHandler{svc: svc, log: log}
}

// GetPages returns a list of RBAC pages.
func (h *RbacPageHandler) GetPages(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	module := c.Query("module")
	keyword := c.Query("keyword")

	pages, total, err := h.svc.GetList(page, pageSize, module, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, pages, total, page, pageSize)
}

// GetPage returns a single RBAC page.
func (h *RbacPageHandler) GetPage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid page id")
		return
	}

	page, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "page not found")
		return
	}
	response.Success(c, page)
}

// GetTree returns RBAC pages as a tree structure.
func (h *RbacPageHandler) GetTree(c *gin.Context) {
	module := c.Query("module")

	tree, err := h.svc.GetTree(module)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, tree)
}

// CreatePage creates a new RBAC page.
func (h *RbacPageHandler) CreatePage(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Path     string `json:"path" binding:"required"`
		Method   string `json:"method"`
		ParentID *uint  `json:"parent_id"`
		Module   string `json:"module"`
		Remark   string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if req.Method == "" {
		req.Method = "GET"
	}

	page := &model.RbacPage{
		Name:     req.Name,
		Path:     req.Path,
		Method:   req.Method,
		ParentID: req.ParentID,
		Module:   req.Module,
		Status:   1,
		Remark:   req.Remark,
	}

	if err := h.svc.Create(page); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, page)
}

// UpdatePage updates an existing RBAC page.
func (h *RbacPageHandler) UpdatePage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid page id")
		return
	}

	var req struct {
		Name     string `json:"name"`
		Path     string `json:"path"`
		Method   string `json:"method"`
		ParentID *uint  `json:"parent_id"`
		Module   string `json:"module"`
		Status   *int16 `json:"status"`
		Remark   string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Path != "" {
		updates["path"] = req.Path
	}
	if req.Method != "" {
		updates["method"] = req.Method
	}
	if req.ParentID != nil {
		updates["parent_id"] = *req.ParentID
	}
	if req.Module != "" {
		updates["module"] = req.Module
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}

	if err := h.svc.Update(uint(id), updates); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "page updated")
}

// DeletePage deletes an RBAC page.
func (h *RbacPageHandler) DeletePage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid page id")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "page deleted")
}

// GetAuthTree returns the auth rule tree for RBAC pages.
func (h *RbacPageHandler) GetAuthTree(c *gin.Context) {
	tree, err := h.svc.GetAuthTree()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, tree)
}
