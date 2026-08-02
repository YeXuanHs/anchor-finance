package handler

import (
	"fmt"
	"os"
	"path/filepath"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ConfigCertifiHandler struct {
	svc *service.ConfigCertifiService
	log *logger.Logger
	db  *gorm.DB
}

func NewConfigCertifiHandler(svc *service.ConfigCertifiService, log *logger.Logger, db *gorm.DB) *ConfigCertifiHandler {
	return &ConfigCertifiHandler{svc: svc, log: log, db: db}
}

// Get is an alias for GetConfig.
func (h *ConfigCertifiHandler) Get(c *gin.Context) { h.GetConfig(c) }

// Update is an alias for UpdateConfig.
func (h *ConfigCertifiHandler) Update(c *gin.Context) { h.UpdateConfig(c) }

// GetConfig returns the certification configuration.
func (h *ConfigCertifiHandler) GetConfig(c *gin.Context) {
	cfg, err := h.svc.GetConfig()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cfg)
}

// UpdateConfig updates the certification configuration.
func (h *ConfigCertifiHandler) UpdateConfig(c *gin.Context) {
	var req service.CertificationConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.UpdateConfig(req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "certification config updated")
}

// GetCertSetting 获取认证设置页面数据
// GET /admin/config-certifi/setting
func (h *ConfigCertifiHandler) GetCertSetting(c *gin.Context) {
	// 认证配置项列表
	configKeys := []string{
		"certifi_is_upload", "certifi_is_stop", "certifi_stop_day",
		"certifi_open", "certifi_select", "certifi_realname",
		"certifi_isbindphone", "artificial_auto_send_msg",
		"certifi_business_btn", "certifi_business_open",
		"certifi_business_is_upload", "certifi_business_is_author",
		"certifi_business_author_path",
	}

	type ConfigItem struct {
		Setting string `json:"setting"`
		Value   string `json:"value"`
	}
	var items []ConfigItem
	h.db.Table("system_configs").Select("setting, value").Where("setting IN ?", configKeys).Find(&items)

	data := make(map[string]interface{})
	for _, item := range items {
		data[item.Setting] = item.Value
	}

	// 默认值处理
	if _, ok := data["certifi_select"]; !ok || data["certifi_select"] == "" {
		data["certifi_select"] = "artificial"
	}

	// 认证类型选择列表
	data["certifi_select_all"] = gin.H{
		"artificial": "人工审核",
	}

	// 授权书路径URL
	if authorPath, ok := data["certifi_business_author_path"]; ok && authorPath != "" {
		data["certifi_business_author_path_url"] = "/uploads/author/" + authorPath
	} else {
		data["certifi_business_author_path_url"] = ""
	}

	response.Success(c, data)
}

// SaveCertSetting 保存认证设置
// POST /admin/config-certifi/setting
func (h *ConfigCertifiHandler) SaveCertSetting(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// certifi_select 数组转逗号分隔
	if selectVal, ok := req["certifi_select"]; ok {
		if arr, ok := selectVal.([]interface{}); ok {
			var strArr []string
			for _, v := range arr {
				if s, ok := v.(string); ok {
					strArr = append(strArr, s)
				}
			}
			req["certifi_select"] = joinStrings(strArr, ",")
		}
	}

	for key, value := range req {
		// 查找是否已存在
		var count int64
		h.db.Table("system_configs").Where("setting = ?", key).Count(&count)
		if count > 0 {
			h.db.Table("system_configs").Where("setting = ?", key).Update("value", fmt.Sprintf("%v", value))
		} else {
			h.db.Table("system_configs").Create(&map[string]interface{}{
				"setting": key,
				"value":   fmt.Sprintf("%v", value),
			})
		}
	}

	response.SuccessMsg(c, "设置成功")
}

// GetCertTypes 获取认证类型列表
// GET /admin/config-certifi/types
func (h *ConfigCertifiHandler) GetCertTypes(c *gin.Context) {
	certTypes := []gin.H{
		{"name": "人工审核", "value": "artificial"},
		{"name": "支付宝认证", "value": "ali"},
		{"name": "手机号三要素", "value": "phonethree"},
	}

	// 排除人工审核
	var types []gin.H
	for _, t := range certTypes {
		if t["value"] != "artificial" {
			types = append(types, t)
		}
	}

	response.Success(c, types)
}

// DownloadAuthor 下载授权证书
// GET /admin/config-certifi/author-down
func (h *ConfigCertifiHandler) DownloadAuthor(c *gin.Context) {
	var configValue string
	h.db.Table("system_configs").Select("value").Where("setting = ?", "certifi_business_author_path").Scan(&configValue)

	if configValue == "" {
		response.BadRequest(c, "文件资源不存在")
		return
	}

	authorDir := filepath.Join("uploads", "author")
	filePath := filepath.Join(authorDir, configValue)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		response.BadRequest(c, "文件资源不存在")
		return
	}

	c.FileAttachment(filePath, "shouQuan")
}

// DeleteAuthor 删除授权证书
// DELETE /admin/config-certifi/author-del
func (h *ConfigCertifiHandler) DeleteAuthor(c *gin.Context) {
	var configValue string
	h.db.Table("system_configs").Select("value").Where("setting = ?", "certifi_business_author_path").Scan(&configValue)

	if configValue == "" {
		response.BadRequest(c, "文件资源不存在")
		return
	}

	authorDir := filepath.Join("uploads", "author")
	filePath := filepath.Join(authorDir, configValue)

	if _, err := os.Stat(filePath); err == nil {
		os.Remove(filePath)
	}

	h.db.Table("system_configs").Where("setting = ?", "certifi_business_author_path").Update("value", "")
	response.SuccessMsg(c, "删除成功")
}

// GetCertDetail 获取认证详情
// GET /admin/config-certifi/detail
func (h *ConfigCertifiHandler) GetCertDetail(c *gin.Context) {
	type ConfigItem struct {
		Setting string `json:"setting"`
		Value   string `json:"value"`
	}
	var items []ConfigItem
	h.db.Table("system_configs").Where("setting LIKE ?", "certifi%").Find(&items)

	data := make(map[string]interface{})
	for _, item := range items {
		data[item.Setting] = item.Value
	}

	if _, ok := data["certifi_select"]; !ok || data["certifi_select"] == "" {
		data["certifi_select"] = "artificial"
	}

	response.Success(c, data)
}

func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
