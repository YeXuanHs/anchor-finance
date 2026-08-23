package admin

import (
	"fmt"
	"strconv"
	"time"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/YeXuanHs/anchor-finance/internal/response"
	"github.com/YeXuanHs/anchor-finance/internal/util"
	"github.com/gin-gonic/gin"
)

// GetDcimServerList DCIM服务器列表（分页+搜索）
// GET /api/admin/dcim-servers
func GetDcimServerList(c *gin.Context) {
	db := database.GetDB()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "15"))
	keyword := c.Query("keyword")
	serverType := c.Query("server_type")
	groupID := c.Query("group_id")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 15
	}

	query := db.Model(&model.Server{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR hostname LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if serverType != "" {
		query = query.Where("server_type = ?", serverType)
	}
	if groupID != "" {
		query = query.Where("group_id = ?", groupID)
	}

	var total int64
	query.Count(&total)

	var servers []model.Server
	query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&servers)

	// 附加dcim扩展信息
	type ServerWithDcim struct {
		model.Server
		DcimServer *model.DcimServer `json:"dcim_server,omitempty"`
	}
	var list []ServerWithDcim
	for _, s := range servers {
		item := ServerWithDcim{Server: s}
		var dcim model.DcimServer
		if err := db.Where("server_id = ?", s.ID).First(&dcim).Error; err == nil {
			item.DcimServer = &dcim
		}
		list = append(list, item)
	}
	if list == nil {
		list = []ServerWithDcim{}
	}

	response.SuccessPage(c, list, total, page, pageSize)
}

// CreateDcimServer 添加DCIM服务器
// POST /api/admin/dcim-servers
func CreateDcimServer(c *gin.Context) {
	db := database.GetDB()

	var req struct {
		Name           string  `json:"name" binding:"required"`
		Hostname       string  `json:"hostname" binding:"required"`
		Username       string  `json:"username"`
		Password       string  `json:"password"`
		Port           int     `json:"port"`
		Secure         bool    `json:"secure"`
		AccessHash     string  `json:"access_hash"`
		ServerType     string  `json:"server_type"`
		GroupID        uint    `json:"group_id"`
		Disabled       bool    `json:"disabled"`
		// DCIM扩展字段
		Area           string  `json:"area"`
		BillType       string  `json:"bill_type"`
		FlowRemind     string  `json:"flow_remind"`
		ReinstallTimes int     `json:"reinstall_times"`
		BuyTimes       int     `json:"buy_times"`
		ReinstallPrice float64 `json:"reinstall_price"`
		Auth           string  `json:"auth"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if req.Port == 0 {
		req.Port = 443
	}
	if req.ServerType == "" {
		req.ServerType = "dcim"
	}
	if req.ReinstallTimes == 0 {
		req.ReinstallTimes = 3
	}
	if req.BuyTimes == 0 {
		req.BuyTimes = 1
	}
	if req.BillType == "" {
		req.BillType = "month"
	}

	// 加密密码
	encryptedPassword := ""
	if req.Password != "" {
		enc, err := util.EncryptAES(req.Password)
		if err != nil {
			response.ServerError(c, "密码加密失败: "+err.Error())
			return
		}
		encryptedPassword = enc
	}

	server := model.Server{
		Name:       req.Name,
		Hostname:   req.Hostname,
		Username:   req.Username,
		Password:   encryptedPassword,
		Port:       req.Port,
		Secure:     req.Secure,
		AccessHash: req.AccessHash,
		ServerType: req.ServerType,
		GroupID:    req.GroupID,
		Disabled:   req.Disabled,
	}

	if err := db.Create(&server).Error; err != nil {
		response.ServerError(c, "创建服务器失败: "+err.Error())
		return
	}

	// 创建DCIM扩展记录
	dcim := model.DcimServer{
		ServerID:       server.ID,
		Auth:           req.Auth,
		Area:           req.Area,
		BillType:       req.BillType,
		FlowRemind:     req.FlowRemind,
		ReinstallTimes: req.ReinstallTimes,
		BuyTimes:       req.BuyTimes,
		ReinstallPrice: req.ReinstallPrice,
		APIStatus:      0,
	}
	db.Create(&dcim)

	response.Success(c, gin.H{
		"id":         server.ID,
		"name":       server.Name,
		"hostname":   server.Hostname,
		"port":       server.Port,
		"server_type": server.ServerType,
		"created_at": server.CreatedAt,
	})
}

// GetDcimServer 服务器详情
// GET /api/admin/dcim-servers/:id
func GetDcimServer(c *gin.Context) {
	db := database.GetDB()
	id := c.Param("id")

	var server model.Server
	if err := db.First(&server, id).Error; err != nil {
		response.NotFound(c, "服务器不存在")
		return
	}

	var dcim model.DcimServer
	db.Where("server_id = ?", server.ID).First(&dcim)

	// 脱敏密码（解密后只显示部分）
	passwordMasked := ""
	if server.Password != "" {
		if dec, err := util.DecryptAES(server.Password); err == nil {
			passwordMasked = util.MaskSecret(dec)
		}
	}

	response.Success(c, gin.H{
		"id":             server.ID,
		"name":           server.Name,
		"hostname":       server.Hostname,
		"username":       server.Username,
		"password_mask":  passwordMasked,
		"port":           server.Port,
		"secure":         server.Secure,
		"access_hash":    server.AccessHash,
		"server_type":    server.ServerType,
		"group_id":       server.GroupID,
		"disabled":       server.Disabled,
		"link_status":    server.LinkStatus,
		"created_at":     server.CreatedAt,
		"updated_at":     server.UpdatedAt,
		"dcim_server":    dcim,
	})
}

// UpdateDcimServer 更新服务器
// PUT /api/admin/dcim-servers/:id
func UpdateDcimServer(c *gin.Context) {
	db := database.GetDB()
	id := c.Param("id")

	var server model.Server
	if err := db.First(&server, id).Error; err != nil {
		response.NotFound(c, "服务器不存在")
		return
	}

	var req struct {
		Name           string  `json:"name"`
		Hostname       string  `json:"hostname"`
		Username       string  `json:"username"`
		Password       string  `json:"password"`
		Port           int     `json:"port"`
		Secure         *bool   `json:"secure"`
		AccessHash     string  `json:"access_hash"`
		ServerType     string  `json:"server_type"`
		GroupID        *uint   `json:"group_id"`
		Disabled       *bool   `json:"disabled"`
		// DCIM扩展字段
		Area           string  `json:"area"`
		BillType       string  `json:"bill_type"`
		FlowRemind     string  `json:"flow_remind"`
		ReinstallTimes *int    `json:"reinstall_times"`
		BuyTimes       *int    `json:"buy_times"`
		ReinstallPrice *float64 `json:"reinstall_price"`
		Auth           string  `json:"auth"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 更新基础字段
	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Hostname != "" {
		updates["hostname"] = req.Hostname
	}
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Password != "" {
		enc, err := util.EncryptAES(req.Password)
		if err != nil {
			response.ServerError(c, "密码加密失败: "+err.Error())
			return
		}
		updates["password"] = enc
	}
	if req.Port > 0 {
		updates["port"] = req.Port
	}
	if req.Secure != nil {
		updates["secure"] = *req.Secure
	}
	if req.AccessHash != "" {
		updates["access_hash"] = req.AccessHash
	}
	if req.ServerType != "" {
		updates["server_type"] = req.ServerType
	}
	if req.GroupID != nil {
		updates["group_id"] = *req.GroupID
	}
	if req.Disabled != nil {
		updates["disabled"] = *req.Disabled
	}

	if len(updates) > 0 {
		db.Model(&server).Updates(updates)
	}

	// 更新DCIM扩展
	var dcim model.DcimServer
	result := db.Where("server_id = ?", server.ID).First(&dcim)
	dcimUpdates := map[string]interface{}{}
	if req.Auth != "" {
		dcimUpdates["auth"] = req.Auth
	}
	if req.Area != "" {
		dcimUpdates["area"] = req.Area
	}
	if req.BillType != "" {
		dcimUpdates["bill_type"] = req.BillType
	}
	if req.FlowRemind != "" {
		dcimUpdates["flow_remind"] = req.FlowRemind
	}
	if req.ReinstallTimes != nil {
		dcimUpdates["reinstall_times"] = *req.ReinstallTimes
	}
	if req.BuyTimes != nil {
		dcimUpdates["buy_times"] = *req.BuyTimes
	}
	if req.ReinstallPrice != nil {
		dcimUpdates["reinstall_price"] = *req.ReinstallPrice
	}

	if len(dcimUpdates) > 0 {
		if result.Error != nil {
			// 不存在则创建
			dcim = model.DcimServer{
				ServerID:       server.ID,
				Auth:           req.Auth,
				Area:           req.Area,
				BillType:       req.BillType,
				FlowRemind:     req.FlowRemind,
				ReinstallTimes: 3,
				BuyTimes:       1,
				ReinstallPrice: 0,
			}
			if req.ReinstallTimes != nil {
				dcim.ReinstallTimes = *req.ReinstallTimes
			}
			if req.BuyTimes != nil {
				dcim.BuyTimes = *req.BuyTimes
			}
			if req.ReinstallPrice != nil {
				dcim.ReinstallPrice = *req.ReinstallPrice
			}
			db.Create(&dcim)
		} else {
			db.Model(&dcim).Updates(dcimUpdates)
		}
	}

	response.SuccessMsg(c, "更新成功")
}

// DeleteDcimServer 删除服务器
// DELETE /api/admin/dcim-servers/:id
func DeleteDcimServer(c *gin.Context) {
	db := database.GetDB()
	id := c.Param("id")

	var server model.Server
	if err := db.First(&server, id).Error; err != nil {
		response.NotFound(c, "服务器不存在")
		return
	}

	// 检查是否有关联的服务
	var serviceCount int64
	db.Model(&model.Service{}).Where("server_id = ?", server.ID).Count(&serviceCount)
	if serviceCount > 0 {
		response.BadRequest(c, fmt.Sprintf("该服务器关联了 %d 个服务，请先解除关联", serviceCount))
		return
	}

	// 删除DCIM扩展记录
	db.Where("server_id = ?", server.ID).Delete(&model.DcimServer{})
	// 软删除服务器
	db.Delete(&server)

	response.SuccessMsg(c, "删除成功")
}

// TestDcimServer 测试连接
// POST /api/admin/dcim-servers/:id/test
func TestDcimServer(c *gin.Context) {
	db := database.GetDB()
	id := c.Param("id")

	var server model.Server
	if err := db.First(&server, id).Error; err != nil {
		response.NotFound(c, "服务器不存在")
		return
	}

	// 尝试连接IPMI（解密密码后测试）
	scheme := "https"
	if !server.Secure {
		scheme = "http"
	}
	testURL := fmt.Sprintf("%s://%s:%d", scheme, server.Hostname, server.Port)

	// 这里可以实现实际的IPMI连接测试
	// 目前模拟测试：检查配置是否完整
	if server.Hostname == "" {
		db.Model(&server).Update("link_status", false)
		response.Error(c, 400, "服务器地址为空，无法连接")
		return
	}

	// 模拟连接测试（真实环境应调用IPMI API）
	// TODO: 实现真实的IPMI连接测试
	connected := true

	// 更新连接状态
	db.Model(&server).Update("link_status", connected)
	// 更新DCIM API状态
	db.Model(&model.DcimServer{}).Where("server_id = ?", server.ID).Update("api_status", 1)

	if connected {
		response.Success(c, gin.H{
			"connected": true,
			"url":       testURL,
			"msg":       "连接成功",
		})
	} else {
		response.Error(c, 500, "连接失败，请检查IPMI配置")
	}
}

// RefreshDcimServerStatus 刷新状态
// POST /api/admin/dcim-servers/:id/refresh-status
func RefreshDcimServerStatus(c *gin.Context) {
	db := database.GetDB()
	id := c.Param("id")

	var server model.Server
	if err := db.First(&server, id).Error; err != nil {
		response.NotFound(c, "服务器不存在")
		return
	}

	// TODO: 实现真实的IPMI状态查询
	// 模拟状态刷新
	apiStatus := 0
	if server.Hostname != "" && !server.Disabled {
		apiStatus = 1
	}

	// 更新连接状态
	linked := apiStatus == 1
	db.Model(&server).Update("link_status", linked)

	// 更新DCIM API状态
	now := time.Now()
	db.Model(&model.DcimServer{}).Where("server_id = ?", server.ID).Updates(map[string]interface{}{
		"api_status": apiStatus,
		"updated_at": now,
	})

	response.Success(c, gin.H{
		"link_status": linked,
		"api_status":  apiStatus,
		"updated_at":  now,
	})
}

// GetDcimServerOptions 获取服务器选项（用于下拉选择）
// GET /api/admin/dcim-servers/options
func GetDcimServerOptions(c *gin.Context) {
	db := database.GetDB()

	var servers []model.Server
	db.Where("disabled = ?", false).Order("name ASC").Find(&servers)

	var options []gin.H
	for _, s := range servers {
		options = append(options, gin.H{
			"id":   s.ID,
			"name": s.Name,
			"host": s.Hostname,
		})
	}
	if options == nil {
		options = []gin.H{}
	}

	response.Success(c, options)
}
