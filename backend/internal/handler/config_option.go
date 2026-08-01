package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type ConfigOptionHandler struct {
	svc *service.ConfigOptionService
	log *logger.Logger
}

func NewConfigOptionHandler(svc *service.ConfigOptionService, log *logger.Logger) *ConfigOptionHandler {
	return &ConfigOptionHandler{svc: svc, log: log}
}

// GetList returns paginated config options.
func (h *ConfigOptionHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	group := c.Query("group")
	keyword := c.Query("keyword")

	items, total, err := h.svc.GetList(page, pageSize, group, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, items, total, page, pageSize)
}

// GetByGroup returns config options filtered by group.
func (h *ConfigOptionHandler) GetByGroup(c *gin.Context) {
	group := c.Param("group")
	if group == "" {
		response.BadRequest(c, "group is required")
		return
	}

	items, err := h.svc.GetByGroup(group)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, items)
}

// GetGroups returns all distinct option groups.
func (h *ConfigOptionHandler) GetGroups(c *gin.Context) {
	groups, err := h.svc.GetGroups()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, groups)
}

// GetDetail returns a single config option by ID.
func (h *ConfigOptionHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid option id")
		return
	}

	item, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "config option not found")
		return
	}
	response.Success(c, item)
}

// Create creates a config option.
func (h *ConfigOptionHandler) Create(c *gin.Context) {
	var req service.CreateConfigOptionRequest
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

// Update updates a config option.
func (h *ConfigOptionHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid option id")
		return
	}

	var req service.UpdateConfigOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	item, err := h.svc.Update(uint(id), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, item)
}

// Delete deletes a config option.
func (h *ConfigOptionHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid option id")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "config option deleted")
}

// BatchUpdateSort batch-updates sort order for config options.
func (h *ConfigOptionHandler) BatchUpdateSort(c *gin.Context) {
	var req []service.SortItem
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.BatchUpdateSort(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "sort order updated")
}

// UpdateSort updates sort order for a single config option.
func (h *ConfigOptionHandler) UpdateSort(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid option id")
		return
	}

	var req struct {
		SortOrder int `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.UpdateSort(uint(id), req.SortOrder); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "sort order updated")
}

// UpdateValue updates the value of a config option by code.
func (h *ConfigOptionHandler) UpdateValue(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		response.BadRequest(c, "option code is required")
		return
	}

	var req struct {
		Value string `json:"value"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.UpdateValue(code, req.Value); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "option value updated")
}

// BatchUpdateValue batch-updates option values by code.
func (h *ConfigOptionHandler) BatchUpdateValue(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.BatchUpdateValue(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "option values updated")
}

// ─── Product Config Groups ───

// GetProductConfigGroups returns all product config groups (admin).
func (h *ConfigOptionHandler) GetProductConfigGroups(c *gin.Context) {
	groups, err := h.svc.GetProductConfigGroups()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, groups)
}

// CreateProductConfigGroup creates a new product config group (admin).
func (h *ConfigOptionHandler) CreateProductConfigGroup(c *gin.Context) {
	var req service.CreateProductConfigGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	group, err := h.svc.CreateProductConfigGroup(req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, group)
}

// UpdateProductConfigGroup updates a product config group (admin).
func (h *ConfigOptionHandler) UpdateProductConfigGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid group id")
		return
	}

	var req service.UpdateProductConfigGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	group, err := h.svc.UpdateProductConfigGroup(uint(id), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, group)
}

// DeleteProductConfigGroup deletes a product config group (admin).
func (h *ConfigOptionHandler) DeleteProductConfigGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid group id")
		return
	}

	if err := h.svc.DeleteProductConfigGroup(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "product config group deleted")
}

// LinkGroupToProduct links a config group to a product (admin).
func (h *ConfigOptionHandler) LinkGroupToProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid group id")
		return
	}

	var req struct {
		ProductID uint `json:"product_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.LinkGroupToProduct(uint(id), req.ProductID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "group linked to product")
}

// UnlinkGroupFromProduct removes the link between a config group and a product (admin).
func (h *ConfigOptionHandler) UnlinkGroupFromProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid group id")
		return
	}

	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product id")
		return
	}

	if err := h.svc.UnlinkGroupFromProduct(uint(id), uint(productID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "group unlinked from product")
}

// GetProductConfig returns config options for a product (public).
func (h *ConfigOptionHandler) GetProductConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product id")
		return
	}

	result, err := h.svc.GetProductConfigOptionsByProduct(uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, result)
}

// ─── ConfigOptions Admin Methods (from zjmf ConfigOptionsController) ───

// GroupsList returns all global config groups with linked products.
func (h *ConfigOptionHandler) GroupsList(c *gin.Context) {
	order := c.DefaultQuery("order", "id")
	sort := c.DefaultQuery("sort", "desc")
	keywords := c.Query("keywords")

	list, err := h.svc.AdminGroupsList(order, sort, keywords)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"list": list, "type": h.svc.AdminSearchPage()})
}

// SearchPage returns product types for the search page.
func (h *ConfigOptionHandler) SearchPage(c *gin.Context) {
	response.Success(c, gin.H{"type": h.svc.AdminSearchPage()})
}

// CreateGroups returns products grouped by product groups.
func (h *ConfigOptionHandler) CreateGroups(c *gin.Context) {
	typeFilter, _ := strconv.Atoi(c.Query("type"))
	products, err := h.svc.AdminCreateGroupsData(typeFilter)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"products": products})
}

// CreateGroupsPost creates a new config group.
func (h *ConfigOptionHandler) CreateGroupsPost(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Global      int    `json:"global"`
		ProductIDs  []uint `json:"p_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	groupID, err := h.svc.AdminCreateGroup(req.Name, req.Description, req.Global, req.ProductIDs)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"groupid": groupID})
}

// EditGroups returns data for editing a config group.
func (h *ConfigOptionHandler) EditGroups(c *gin.Context) {
	gid, _ := strconv.ParseUint(c.Query("gid"), 10, 64)
	if gid == 0 {
		response.BadRequest(c, "gid is required")
		return
	}
	typeFilter, _ := strconv.Atoi(c.Query("type"))

	data, err := h.svc.AdminEditGroupData(uint(gid), typeFilter)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, data)
}

// EditGroupsPost updates a config group.
func (h *ConfigOptionHandler) EditGroupsPost(c *gin.Context) {
	var req struct {
		GID         uint           `json:"gid" binding:"required"`
		Name        string         `json:"name"`
		Description string         `json:"description"`
		ProductIDs  []uint         `json:"productlinks"`
		Order       map[string]int `json:"order"`
		Hidden      map[string]int `json:"hidden"`
		Upgrade     map[string]int `json:"upgrade"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.AdminEditGroupPost(req.GID, req.Name, req.Description, req.ProductIDs, req.Order, req.Hidden, req.Upgrade); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "config group updated")
}

// AddOptionsPage returns data for the add options page.
func (h *ConfigOptionHandler) AddOptionsPage(c *gin.Context) {
	gid, _ := strconv.ParseUint(c.Query("gid"), 10, 64)
	if gid == 0 {
		response.BadRequest(c, "gid is required")
		return
	}
	pid, _ := strconv.ParseUint(c.Query("pid"), 10, 64)

	data, err := h.svc.AdminAddOptionsPage(uint(gid), uint(pid))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, data)
}

// AddOptions creates a new config option with a sub-option.
func (h *ConfigOptionHandler) AddOptions(c *gin.Context) {
	var req struct {
		GID           uint   `json:"gid" binding:"required"`
		OptionName    string `json:"option_name" binding:"required"`
		OptionType    int    `json:"option_type"`
		Notes         string `json:"notes"`
		QtyStage      int    `json:"qty_stage"`
		LinkagePID    uint   `json:"linkage_pid"`
		AddOptionName string `json:"addoptionname" binding:"required"`
		AddSortOrder  int    `json:"addsortorder"`
		AddHidden     int    `json:"addhidden"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	cid, err := h.svc.AdminAddOption(req.GID, req.OptionName, req.OptionType, req.Notes, req.QtyStage, req.LinkagePID, req.AddOptionName, req.AddSortOrder, req.AddHidden)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"cid": cid})
}

// DeleteSubOptions deletes a config sub-option.
func (h *ConfigOptionHandler) DeleteSubOptions(c *gin.Context) {
	subID, _ := strconv.ParseUint(c.Query("subid"), 10, 64)
	if subID == 0 {
		response.BadRequest(c, "subid is required")
		return
	}

	if err := h.svc.AdminDeleteSubOption(uint(subID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "sub-option deleted")
}

// DeleteOptions deletes a config option and all its sub-options.
func (h *ConfigOptionHandler) DeleteOptions(c *gin.Context) {
	cid, _ := strconv.ParseUint(c.Query("cid"), 10, 64)
	if cid == 0 {
		response.BadRequest(c, "cid is required")
		return
	}

	if err := h.svc.AdminDeleteOption(uint(cid)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "option deleted")
}

// DeleteGroups deletes a config group and all its options.
func (h *ConfigOptionHandler) DeleteGroups(c *gin.Context) {
	gid, _ := strconv.ParseUint(c.Query("gid"), 10, 64)
	if gid == 0 {
		response.BadRequest(c, "gid is required")
		return
	}

	if err := h.svc.AdminDeleteGroup(uint(gid)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "group deleted")
}

// DuplicateGroups returns all global groups for duplication.
func (h *ConfigOptionHandler) DuplicateGroups(c *gin.Context) {
	groups, err := h.svc.AdminDuplicateGroups()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"groups": groups})
}

// DuplicateGroupPost duplicates a config group.
func (h *ConfigOptionHandler) DuplicateGroupPost(c *gin.Context) {
	var req struct {
		GID     uint   `json:"gid" binding:"required"`
		NewName string `json:"newname" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.AdminDuplicateGroupPost(req.GID, req.NewName); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "group duplicated")
}

// ConfigOptionsCheckOs returns products with OS-type config options.
func (h *ConfigOptionHandler) ConfigOptionsCheckOs(c *gin.Context) {
	products, err := h.svc.AdminCheckOS()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"products": products})
}

// ConfigOptionsOs returns OS options for a product.
func (h *ConfigOptionHandler) ConfigOptionsOs(c *gin.Context) {
	pid, _ := strconv.ParseUint(c.Query("pid"), 10, 64)
	if pid == 0 {
		response.BadRequest(c, "pid is required")
		return
	}

	data, err := h.svc.AdminGetOS(uint(pid))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

// EditConfig returns config option detail with sub-options.
func (h *ConfigOptionHandler) EditConfig(c *gin.Context) {
	cid, _ := strconv.ParseUint(c.Query("cid"), 10, 64)
	if cid == 0 {
		response.BadRequest(c, "cid is required")
		return
	}
	pid, _ := strconv.ParseUint(c.Query("pid"), 10, 64)

	data, err := h.svc.AdminEditConfig(uint(cid), uint(pid))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, data)
}

// GetNextLinkAgeList returns the next level of linkage options.
func (h *ConfigOptionHandler) GetNextLinkAgeList(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	if id == 0 {
		response.BadRequest(c, "id is required")
		return
	}

	data, err := h.svc.AdminGetNextLinkAgeList(uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, data)
}

// EditConfigPost updates a config option and its sub-options.
func (h *ConfigOptionHandler) EditConfigPost(c *gin.Context) {
	var req struct {
		CID             uint                                    `json:"cid" binding:"required"`
		OptionName      string                                  `json:"configoptionname"`
		OptionType      int                                     `json:"configoptiontype"`
		Notes           string                                  `json:"notes"`
		QtyMin          int                                     `json:"qtyminimum"`
		QtyMax          int                                     `json:"qtymaximum"`
		QtyStage        int                                     `json:"qty_stage"`
		IsDiscount      int                                     `json:"is_discount"`
		Senior          int                                     `json:"senior"`
		Unit            string                                  `json:"unit"`
		SubUpdates      map[uint]map[string]interface{}          `json:"sub_updates"`
		Prices          map[string]map[string]map[string]interface{} `json:"price"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.AdminEditConfigPost(req.CID, req.OptionName, req.OptionType, req.Notes, req.QtyMin, req.QtyMax, req.QtyStage, req.IsDiscount, req.Senior, req.Unit, req.SubUpdates, req.Prices); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "config option updated")
}

// SaveLinkAgeLevel saves a linkage level config option.
func (h *ConfigOptionHandler) SaveLinkAgeLevel(c *gin.Context) {
	var req struct {
		GID             uint   `json:"gid" binding:"required"`
		OptionName      string `json:"option_name" binding:"required"`
		OptionType      int    `json:"option_type"`
		Notes           string `json:"notes"`
		LinkagePID      uint   `json:"linkage_pid"`
		SubOptionName   string `json:"option_sub_name" binding:"required"`
		SubLinkagePID   uint   `json:"sub_linkage_pid"`
		Hidden          int    `json:"hidden"`
		OptionID        uint   `json:"option_id"`
		SubOptionID     uint   `json:"sub_option_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	optID, subID, err := h.svc.AdminSaveLinkAgeLevel(req.GID, req.OptionName, req.OptionType, req.Notes, req.LinkagePID, req.SubOptionName, req.SubLinkagePID, req.Hidden, req.OptionID, req.SubOptionID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"option_id": optID, "sub_option_id": subID})
}

// SaveConfigOptionInfo saves basic config option info.
func (h *ConfigOptionHandler) SaveConfigOptionInfo(c *gin.Context) {
	var req struct {
		GID          uint   `json:"gid" binding:"required"`
		OptionName   string `json:"option_name" binding:"required"`
		OptionType   int    `json:"option_type"`
		Notes        string `json:"notes"`
		QtyStage     int    `json:"qty_stage"`
		LinkagePID   uint   `json:"linkage_pid"`
		IsDiscount   int    `json:"is_discount"`
		Unit         string `json:"unit"`
		Senior       int    `json:"senior"`
		OptionID     uint   `json:"option_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	optID, err := h.svc.AdminSaveConfigOptionInfo(req.GID, req.OptionName, req.OptionType, req.Notes, req.QtyStage, req.LinkagePID, req.IsDiscount, req.Unit, req.Senior, req.OptionID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"option_id": optID})
}

// SaveLinkAgeOrder saves sort order for linkage sub-options.
func (h *ConfigOptionHandler) SaveLinkAgeOrder(c *gin.Context) {
	var req struct {
		SubIDs []uint `json:"sub_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.AdminSaveLinkAgeOrder(req.SubIDs); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "sort order saved")
}

// DelLinkAgeSub deletes a linkage sub-option.
func (h *ConfigOptionHandler) DelLinkAgeSub(c *gin.Context) {
	subID, _ := strconv.ParseUint(c.Query("sub_id"), 10, 64)
	if subID == 0 {
		response.BadRequest(c, "sub_id is required")
		return
	}

	if err := h.svc.AdminDelLinkAgeSub(uint(subID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "linkage sub deleted")
}
