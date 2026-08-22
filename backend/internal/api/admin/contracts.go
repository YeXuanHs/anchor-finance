package admin

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetContractList 获取合同列表
// GET /api/admin/contracts
func GetContractList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")
	userID := c.Query("user_id")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetDB()
	query := db.Model(&model.Contract{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	var total int64
	query.Count(&total)

	var contracts []model.Contract
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&contracts)

	if contracts == nil {
		contracts = []model.Contract{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      contracts,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetContractDetail 获取合同详情
// GET /api/admin/contracts/:id
func GetContractDetail(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var contract model.Contract
	if err := db.First(&contract, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "合同不存在", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": contract})
}

// CreateContract 创建合同
// POST /api/admin/contracts
func CreateContract(c *gin.Context) {
	var req struct {
		UserID    uint   `json:"user_id" binding:"required"`
		Title     string `json:"title" binding:"required"`
		Content   string `json:"content"`
		Type      string `json:"type"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 生成合同号
	contractNo := fmt.Sprintf("CON%s%06d", time.Now().Format("20060102"), time.Now().UnixNano()%1000000)

	db := database.GetDB()
	contract := model.Contract{
		UserID:     req.UserID,
		ContractNo: contractNo,
		Title:      req.Title,
		Content:    req.Content,
		Type:       req.Type,
		Status:     "draft",
	}

	if req.StartDate != "" {
		t, _ := time.Parse("2006-01-02", req.StartDate)
		contract.StartDate = &t
	}
	if req.EndDate != "" {
		t, _ := time.Parse("2006-01-02", req.EndDate)
		contract.EndDate = &t
	}

	if err := db.Create(&contract).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "操作失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id":          contract.ID,
			"contract_no": contract.ContractNo,
		},
	})
}

// UpdateContract 更新合同
// PUT /api/admin/contracts/:id
func UpdateContract(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Status  string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	db := database.GetDB()
	var contract model.Contract
	if err := db.First(&contract, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "合同不存在", "data": nil})
		return
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	db.Model(&contract).Updates(updates)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": nil})
}

// DeleteContract 删除合同
// DELETE /api/admin/contracts/:id
func DeleteContract(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var contract model.Contract
	if err := db.First(&contract, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "合同不存在", "data": nil})
		return
	}

	db.Delete(&contract)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功", "data": nil})
}

// SignContract 签署合同
// POST /api/admin/contracts/:id/sign
func SignContract(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var contract model.Contract
	if err := db.First(&contract, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "合同不存在", "data": nil})
		return
	}

	if contract.Status != "pending" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "只有待签署的合同才能签署", "data": nil})
		return
	}

	now := time.Now()
	db.Model(&contract).Updates(map[string]interface{}{
		"status":    "signed",
		"signed_at": &now,
	})

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "签署成功", "data": nil})
}

// CancelContract 取消合同
// POST /api/admin/contracts/:id/cancel
func CancelContract(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var contract model.Contract
	if err := db.First(&contract, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "合同不存在", "data": nil})
		return
	}

	if contract.Status == "signed" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "已签署的合同不能取消", "data": nil})
		return
	}

	db.Model(&contract).Update("status", "cancelled")

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "取消成功", "data": nil})
}

// GetContractTemplateList 获取合同模板列表
// GET /api/admin/contract-templates
func GetContractTemplateList(c *gin.Context) {
	db := database.GetDB()
	var templates []model.ContractTemplate
	db.Where("status = ?", "active").Order("id ASC").Find(&templates)

	if templates == nil {
		templates = []model.ContractTemplate{}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": templates})
}

// CreateContractTemplate 创建合同模板
// POST /api/admin/contract-templates
func CreateContractTemplate(c *gin.Context) {
	var req struct {
		Name      string `json:"name" binding:"required"`
		Content   string `json:"content" binding:"required"`
		Variables string `json:"variables"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	db := database.GetDB()
	template := model.ContractTemplate{
		Name:      req.Name,
		Content:   req.Content,
		Variables: req.Variables,
		Status:    "active",
	}

	if err := db.Create(&template).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "操作失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": template.ID,
		},
	})
}

// UpdateContractTemplate 更新合同模板
// PUT /api/admin/contract-templates/:id
func UpdateContractTemplate(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name      string `json:"name"`
		Content   string `json:"content"`
		Variables string `json:"variables"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	db := database.GetDB()
	var template model.ContractTemplate
	if err := db.First(&template, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "模板不存在", "data": nil})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Variables != "" {
		updates["variables"] = req.Variables
	}

	db.Model(&template).Updates(updates)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": nil})
}

// DeleteContractTemplate 删除合同模板
// DELETE /api/admin/contract-templates/:id
func DeleteContractTemplate(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var template model.ContractTemplate
	if err := db.First(&template, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "模板不存在", "data": nil})
		return
	}

	db.Delete(&template)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功", "data": nil})
}
