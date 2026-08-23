package admin

import (
	"net/http"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// domainMapping Go端域概念名 → PHP端实际目录名
// MD定义：payment/upstream/verification，PHP实际：gateways/servers/certification
var domainMapping = map[string][]string{
	"payment":       {"payment", "gateways"},
	"upstream":      {"upstream", "servers"},
	"verification":  {"verification", "certification"},
	"sms":           {"sms"},
	"mail":          {"mail"},
	"oauth":         {"oauth"},
	"captcha":       {"captcha"},
	"addons":        {"addons"},
}

// getPluginsByDomain 通用函数：按域获取已启用插件列表（支持域名映射）
func getPluginsByDomain(domain string) []model.Plugin {
	db := database.GetDB()
	var plugins []model.Plugin

	domains := []string{domain}
	if mapped, ok := domainMapping[domain]; ok {
		domains = mapped
	}

	db.Where("domain IN ? AND status = ?", domains, "active").Order("name ASC").Find(&plugins)
	if plugins == nil {
		plugins = []model.Plugin{}
	}
	return plugins
}

// GetPaymentGateways 获取支付网关列表
// GET /api/admin/payment-gateways
func GetPaymentGateways(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": getPluginsByDomain("payment")})
}

// GetSMSProviders 获取短信提供商列表
// GET /api/admin/sms-providers
func GetSMSProviders(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": getPluginsByDomain("sms")})
}

// GetMailProviders 获取邮件提供商列表
// GET /api/admin/mail-providers
func GetMailProviders(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": getPluginsByDomain("mail")})
}

// GetCertificationProviders 获取实名认证提供商列表
// GET /api/admin/certification-providers
func GetCertificationProviders(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": getPluginsByDomain("verification")})
}

// GetServerModules 获取服务器开通模块列表
// GET /api/admin/server-modules
func GetServerModules(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": getPluginsByDomain("upstream")})
}
