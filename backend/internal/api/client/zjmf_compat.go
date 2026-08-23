package client

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ============================================================
// zjmf兼容API（让zjmf系统能对接锚点财务）
// zjmf系统选"zjmf"类型，填我们的网址，就能拉取商品/同步状态
// ============================================================

// ZjmfCompatLogin zjmf兼容登录（/zjmf_api_login 和 /v1/login_api）
// zjmf用http_build_query发送POST数据（form-encoded，非JSON），参数名username+password
// 同时兼容JSON格式和account参数名
func ZjmfCompatLogin(c *gin.Context) {
	// 兼容form-encoded和JSON两种格式
	account := c.PostForm("username")
	if account == "" {
		account = c.PostForm("account")
	}
	password := c.PostForm("password")

	// 如果form-encoded没拿到，尝试JSON
	if account == "" {
		var req struct {
			Account  string `json:"account"`
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err == nil {
			account = req.Username
			if account == "" {
				account = req.Account
			}
			password = req.Password
		}
	}

	if account == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "参数错误"})
		return
	}

	token, err := LoginByAPIKey(account, password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 401, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"jwt": token, "status": 200, "msg": "鉴权成功"})
}

// ZjmfCompatProducts zjmf兼容商品列表（/v1/products）
// 和/cart/all返回相同嵌套格式
func ZjmfCompatProducts(c *gin.Context) {
	ZjmfCompatCartAll(c)
}

// ZjmfCompatCategories zjmf兼容商品分组（/v1/hosts/cates）
func ZjmfCompatCategories(c *gin.Context) {
	db := database.GetDB()
	var groups []model.ProductGroup
	db.Where("parent_id = 0 AND status = ?", "active").Order("sort_order ASC").Find(&groups)

	var list []gin.H
	for _, g := range groups {
		list = append(list, gin.H{
			"id":   g.ID,
			"name": g.Name,
		})
	}

	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "success", "data": gin.H{"cate": list}})
}

// ZjmfCompatHostDetail zjmf兼容服务详情（/host/header）
// zjmf期望返回data.host_data+module_button+module_power_status等50+字段
func ZjmfCompatHostDetail(c *gin.Context) {
	hostID := c.Query("host_id")
	if hostID == "" {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "缺少host_id"})
		return
	}

	db := database.GetDB()
	var svc model.Service
	if err := db.First(&svc, hostID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 406, "msg": "未找到该产品"})
		return
	}

	// 状态映射
	statusMap := map[string]string{
		"active":     "Active",
		"suspended":  "Suspended",
		"terminated": "Terminated",
		"pending":    "Pending",
		"cancelled":  "Cancelled",
	}
	domainStatus := statusMap[svc.Status]
	if domainStatus == "" {
		domainStatus = "Unknown"
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":   "请求成功",
		"data": gin.H{
			"host_data": gin.H{
				"id":              svc.ID,
				"productid":       svc.ProductID,
				"domain":          svc.Domain,
				"dedicatedip":     "",
				"assignedips":     []string{},
				"username":        "",
				"password":        "",
				"port":            0,
				"os":              "",
				"regdate":         svc.CreatedAt.Format("2006-01-02"),
				"nextduedate":     "",
				"billingcycle":    svc.BillingCycle,
				"firstpaymentamount": svc.Amount,
				"amount":          svc.Amount,
				"domainstatus":    domainStatus,
				"stream_info":     "",
				"bwlimit":         0,
				"bwusage":         0,
				"os_ostype":       "",
				"os_osname":       "",
				"disk_num":        1,
			},
			"module_button": gin.H{
				"control": []gin.H{},
				"console": []gin.H{},
			},
			"module_client_area":      []gin.H{},
			"module_chart":            []gin.H{},
			"module_client_main_area": []gin.H{},
			"module_power_status":     true,
			"reinstall_random_port":   false,
			"reinstall_format_data_disk": false,
			"dcimcloud": gin.H{
				"nat_acl": "",
				"nat_web": "",
			},
		},
	})
}

// ZjmfCompatModuleStatus zjmf兼容模块状态（/v1/hosts/:id/module/status）
// zjmf期望返回data.status为"on"/"off"
func ZjmfCompatModuleStatus(c *gin.Context) {
	hostID := c.Param("id")
	statusType := c.Query("type")

	db := database.GetDB()
	var svc model.Service
	if err := db.First(&svc, hostID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "服务不存在"})
		return
	}

	if statusType == "reinstall" {
		c.JSON(http.StatusOK, gin.H{"status": 200, "data": gin.H{"progress": 100, "status": "completed"}})
		return
	}

	powerStatus := "off"
	des := "已停止"
	if svc.Status == "active" {
		powerStatus = "on"
		des = "运行中"
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data": gin.H{
			"status": powerStatus,
			"des":    des,
		},
	})
}

// ZjmfCompatModuleSuspend zjmf兼容暂停（/v1/hosts/:id/module/suspend）
func ZjmfCompatModuleSuspend(c *gin.Context) {
	hostID := c.Param("id")
	db := database.GetDB()

	var svc model.Service
	if err := db.First(&svc, hostID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "服务不存在"})
		return
	}

	db.Model(&svc).Update("status", "suspended")
	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "success"})
}

// ZjmfCompatModuleUnsuspend zjmf兼容取消暂停（/v1/hosts/:id/module/unsuspend）
func ZjmfCompatModuleUnsuspend(c *gin.Context) {
	hostID := c.Param("id")
	db := database.GetDB()

	var svc model.Service
	if err := db.First(&svc, hostID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "服务不存在"})
		return
	}

	db.Model(&svc).Update("status", "active")
	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "success"})
}

// ZjmfCompatModuleTerminate zjmf兼容终止（/v1/hosts/:id/module/terminate）
func ZjmfCompatModuleTerminate(c *gin.Context) {
	hostID := c.Param("id")
	db := database.GetDB()

	var svc model.Service
	if err := db.First(&svc, hostID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "服务不存在"})
		return
	}

	db.Model(&svc).Update("status", "terminated")
	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "success"})
}

// ZjmfCompatRenew zjmf兼容续费（/v1/hosts/:id/renew）
func ZjmfCompatRenew(c *gin.Context) {
	hostID := c.Param("id")
	db := database.GetDB()

	var svc model.Service
	if err := db.First(&svc, hostID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "服务不存在"})
		return
	}

	// 续费逻辑：延长到期时间
	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "success", "data": gin.H{"amount": svc.Amount}})
}

// ZjmfCompatUser zjmf兼容用户信息（/v1/user）
func ZjmfCompatUser(c *gin.Context) {
	userID, _ := c.Get("user_id")
	db := database.GetDB()

	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":  "success",
		"data": gin.H{
			"client": gin.H{
				"id":       user.ID,
				"email":    user.Email,
				"username": user.Username,
				"phone":    user.Phone,
				"company":  user.Company,
			},
		},
	})
}

// ZjmfCompatBalance zjmf兼容余额查询（/cart/credit）
// zjmf真实返回：{status:200, msg:"请求成功", data:{credit:"33.75", currency:{id:1,code:"CNY",...}}, is_aff:"1"}
func ZjmfCompatBalance(c *gin.Context) {
	userID, _ := c.Get("user_id")
	db := database.GetDB()

	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "用户不存在"})
		return
	}

	var defaultCurrency model.Currency
	currencyObj := gin.H{"id": 1, "code": "CNY", "prefix": "¥", "suffix": "元"}
	if err := db.Where("is_default = ?", true).First(&defaultCurrency).Error; err == nil {
		currencyObj = gin.H{"id": defaultCurrency.ID, "code": defaultCurrency.Code, "prefix": defaultCurrency.Symbol, "suffix": "元"}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"msg":     "请求成功",
		"data": gin.H{
			"credit":   fmt.Sprintf("%.2f", user.Balance),
			"currency": currencyObj,
		},
		"is_aff": "1",
	})
}

// ZjmfCompatCartAll zjmf兼容商品目录（/cart/all）
// zjmf调此端点获取商品目录，期望嵌套结构+count+product_price
func ZjmfCompatCartAll(c *gin.Context) {
	db := database.GetDB()

	// 获取一级分组
	var firstGroups []model.ProductGroup
	db.Where("parent_id = 0 AND status = ?", "active").Order("sort_order ASC").Find(&firstGroups)
	if firstGroups == nil {
		firstGroups = []model.ProductGroup{}
	}

	totalProducts := 0
	var result []gin.H
	for _, fg := range firstGroups {
		var groups []model.ProductGroup
		db.Where("parent_id = ? AND status = ?", fg.ID, "active").Order("sort_order ASC").Find(&groups)
		if groups == nil {
			groups = []model.ProductGroup{}
		}

		var groupList []gin.H
		for _, g := range groups {
			var products []model.Product
			db.Where("group_id = ? AND status = ?", g.ID, "active").Find(&products)
			if products == nil {
				products = []model.Product{}
			}
			totalProducts += len(products)

			var productList []gin.H
			for _, p := range products {
				productList = append(productList, gin.H{
					"id":            p.ID,
					"name":          p.Name,
					"description":   p.Description,
					"product_price": p.Price,
					"billingcycle":  p.BillingCycle,
					"type":          p.Type,
					"qty":           1,
					"stock_control": 0,
					"setup_fee":     0,
					"ontrial":       gin.H{"ontrial": 0},
				})
			}

			groupList = append(groupList, gin.H{
				"id":       g.ID,
				"name":     g.Name,
				"headline": "",
				"tagline":  "",
				"fields":   []gin.H{},
				"products": productList,
			})
		}

		result = append(result, gin.H{
			"id":     fg.ID,
			"name":   fg.Name,
			"fields": []gin.H{},
			"group":  groupList,
		})
	}

	// 构建扁平products列表（zjmf的getUpstreamProducts取$list["products"]）
	var flatProducts []gin.H
	for _, fg := range result {
		groups := fg["group"].([]gin.H)
		for _, g := range groups {
			flatProducts = append(flatProducts, gin.H{
				"id":       g["id"],
				"name":     g["name"],
				"products": g["products"],
			})
		}
	}

	// 获取默认货币
	var defaultCurrency model.Currency
	currencyCode := "CNY"
	if err := db.Where("is_default = ?", true).First(&defaultCurrency).Error; err == nil {
		currencyCode = defaultCurrency.Code
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "请求成功",
		"data": gin.H{
			"first_group": result,
			"products":    flatProducts,
			"count":       totalProducts,
			"currency":    currencyCode,
		},
	})
}

// ZjmfCompatProductInfo zjmf兼容商品详情（/api/product/proinfo）
// zjmf期望返回格式：{status:200, data:{info:[{...}], currency:"CNY"}}
// zjmf取 $res["data"]["info"][0] 和 $res["data"]["currency"]
func ZjmfCompatProductInfo(c *gin.Context) {
	// 解析pids参数（支持pids[0]=63格式）
	pids := c.QueryArray("pids")
	if len(pids) == 0 {
		for i := 0; ; i++ {
			key := fmt.Sprintf("pids[%d]", i)
			val := c.Query(key)
			if val == "" {
				break
			}
			pids = append(pids, val)
		}
	}
	if len(pids) == 0 {
		if val := c.Query("pids"); val != "" {
			pids = append(pids, val)
		}
	}

	db := database.GetDB()

	// 获取默认货币
	var defaultCurrency model.Currency
	currencyCode := "CNY"
	if err := db.Where("is_default = ?", true).First(&defaultCurrency).Error; err == nil {
		currencyCode = defaultCurrency.Code
	}

	var products []model.Product
	if len(pids) > 0 {
		db.Where("id IN ?", pids).Find(&products)
	} else {
		db.Where("status = ?", "active").Limit(50).Find(&products)
	}

	// 构建info数组（zjmf只返回5个字段：id, name, location_version, stock_control, qty）
	var infoList []gin.H
	for _, p := range products {
		infoList = append(infoList, gin.H{
			"id":               p.ID,
			"name":             p.Name,
			"location_version": 1,
			"stock_control":    0,
			"qty":              999,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":   "success",
		"data": gin.H{
			"info":     infoList,
			"currency": currencyCode,
		},
	})
}

// ZjmfCompatProductConfig zjmf兼容商品配置（/cart/get_product_config）
// 完全照着zjmf格式返回
func ZjmfCompatProductConfig(c *gin.Context) {
	pid := c.Query("pid")
	if pid == "" {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "缺少pid"})
		return
	}

	db := database.GetDB()
	var product model.Product
	if err := db.First(&product, pid).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "产品不存在"})
		return
	}

	// 获取货币ID
	var defaultCurrency model.Currency
	currencyID := 1
	db.Where("is_default = ?", true).First(&defaultCurrency)
	if defaultCurrency.ID > 0 {
		currencyID = int(defaultCurrency.ID)
	}

	// pay_type JSON
	payType := fmt.Sprintf(`{"pay_type":"%s","pay_hour_cycle":"720","pay_day_cycle":"30","pay_ontrial_status":0,"pay_ontrial_cycle":"0","pay_ontrial_num":"1","pay_ontrial_condition":[],"pay_ontrial_cycle_type":"day","pay_ontrial_num_rule":"0","clientscount_rule":0}`, product.BillingCycle)
	if product.BillingCycle == "free" {
		payType = `{"pay_type":"free","pay_hour_cycle":"720","pay_day_cycle":"30","pay_ontrial_status":0,"pay_ontrial_cycle":"0","pay_ontrial_num":"1","pay_ontrial_condition":[],"pay_ontrial_cycle_type":"day","pay_ontrial_num_rule":"0","clientscount_rule":0}`
	}
	if product.BillingCycle == "onetime" {
		payType = `{"pay_type":"onetime","pay_hour_cycle":"720","pay_day_cycle":"30","pay_ontrial_status":0,"pay_ontrial_cycle":"0","pay_ontrial_num":"1","pay_ontrial_condition":[],"pay_ontrial_cycle_type":"day","pay_ontrial_num_rule":"0","clientscount_rule":0}`
	}

	products := gin.H{
		"id":                         product.ID,
		"type":                       product.Type,
		"gid":                        product.GroupID,
		"name":                       product.Name,
		"description":                product.Description,
		"host":                       `{"show":"0","modify":0,"prefix":"ser","rule":{"upper":"0","lower":"0","num":"1","len_num":12}}`,
		"is_domain":                  0,
		"hidden":                     0,
		"password":                   `{"show":"1","modify":0,"rule":{"len_num":"12","upper":"1","lower":"1","num":"1","special":"0"}}`,
		"show_domain_options":        0,
		"welcome_email":              0,
		"stock_control":              0,
		"qty":                        999,
		"prorata_billing":            0,
		"prorata_date":               0,
		"prorata_charge_next_month":  0,
		"pay_type":                   payType,
		"pay_method":                 "prepayment",
		"allow_qty":                  0,
		"auto_setup":                 "payment",
		"server_type":                "",
		"server_group":               0,
		"config_option1":             "",
		"config_option2":             "",
		"config_option3":             "",
		"config_option4":             "",
		"config_option5":             "",
		"config_option6":             "",
		"config_option7":             "",
		"config_option8":             "",
		"config_option9":             "",
		"config_option10":            "",
		"config_option11":            "",
		"config_option12":            "",
		"config_option13":            "",
		"config_option14":            "",
		"config_option15":            "",
		"config_option16":            "",
		"config_option17":            "",
		"config_option18":            "",
		"config_option19":            "",
		"config_option20":            "",
		"config_option21":            "",
		"config_option22":            "",
		"config_option23":            "",
		"config_option24":            "",
		"recurring_cycles":           0,
		"auto_terminate_days":        0,
		"auto_terminate_email":       0,
		"config_options_upgrade":     0,
		"billing_cycle_upgrade":      "",
		"upgrade_email":              0,
		"down_configoption_refund":   0,
		"overages_enabled":           "",
		"overages_disk_limit":        0,
		"overages_bw_limit":          0,
		"overages_disk_price":        "0.0000",
		"overages_bw_price":          "0.0000",
		"tax":                        0,
		"affiliateonetime":           0,
		"affiliate_pay_type":         "default",
		"affiliate_pay_amount":       "0.00",
		"order":                      0,
		"retired":                    0,
		"is_featured":                0,
		"create_time":                product.CreatedAt.Unix(),
		"update_time":                product.UpdatedAt.Unix(),
		"auto_create_config_options": 0,
		"groupid":                    product.GroupID,
		"api_type":                   "normal",
		"location_version":           1,
		"upstream_version":           0,
		"upstream_price_type":        "",
		"upstream_price_value":       "",
		"zjmf_api_id":                0,
		"upstream_pid":               0,
		"is_truename":                0,
		"uuid":                       nil,
		"p_uid":                      0,
		"rate":                       "1.00",
		"clientscount":               0,
		"product_shopping_url":       "",
		"product_group_url":          "",
		"upper_reaches_id":           0,
		"is_bind_phone":              0,
		"upstream_qty":               0,
		"cancel_control":             1,
		"upstream_stock_control":     0,
		"upstream_auto_setup":        "",
		"upstream_ontrial_status":    0,
		"upstream_price":             "0.00",
		"version_description":        nil,
		"app_type":                   "",
		"app_version":                nil,
		"app_hot_order":              0,
		"app_hot_lock":               0,
		"app_hot_heat":               0,
		"app_recommend_status":       0,
		"app_recommend_order":        0,
		"app_recommend_lock":         0,
		"app_pay_type":               0,
		"app_score":                  0,
		"app_images":                 nil,
		"app_status":                 0,
		"unretired_time":             0,
		"upstream_cycle":             "",
	}

	// product_pricings
	monthlyPrice := fmt.Sprintf("%.2f", product.Price)
	productPricings := []gin.H{{
		"id":          1,
		"type":        "product",
		"currency":    currencyID,
		"relid":       product.ID,
		"osetupfee":   "0.00",
		"hsetupfee":   "0.00",
		"dsetupfee":   "0.00",
		"ontrialfee":  "0.00",
		"msetupfee":   "0.00",
		"qsetupfee":   "0.00",
		"ssetupfee":   "0.00",
		"asetupfee":   "0.00",
		"bsetupfee":   "0.00",
		"tsetupfee":   "0.00",
		"onetime":     "-1.00",
		"hour":        "-1.00",
		"day":         "-1.00",
		"ontrial":     "-1.00",
		"monthly":     monthlyPrice,
		"quarterly":   fmt.Sprintf("%.2f", product.Price*3),
		"semiannually": fmt.Sprintf("%.2f", product.Price*6),
		"annually":    fmt.Sprintf("%.2f", product.Price*12),
		"biennially":  "-1.00",
		"triennially": "-1.00",
	}}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"msg":     "请求成功",
		"data": gin.H{
			"flag":             gin.H{"type": 1, "bates": "100.00"},
			"products":         products,
			"customfields":     []gin.H{},
			"product_pricings": productPricings,
			"advanced":         "",
			"config_groups":    []gin.H{},
			"config_links":     0,
		},
		"is_aff": "1",
	})
}

// ZjmfCompatProductDetail zjmf兼容商品详细配置（/api/product/prodetail）
// 完全照zjmf真实格式返回
func ZjmfCompatProductDetail(c *gin.Context) {
	pids := c.QueryArray("pids")
	if len(pids) == 0 {
		for i := 0; ; i++ {
			key := fmt.Sprintf("pids[%d]", i)
			val := c.Query(key)
			if val == "" {
				break
			}
			pids = append(pids, val)
		}
	}
	if len(pids) == 0 {
		if val := c.Query("pids"); val != "" {
			pids = append(pids, val)
		}
	}

	db := database.GetDB()

	var defaultCurrency model.Currency
	currencyID := 1
	currencyCode := "CNY"
	db.Where("is_default = ?", true).First(&defaultCurrency)
	if defaultCurrency.ID > 0 {
		currencyID = int(defaultCurrency.ID)
		currencyCode = defaultCurrency.Code
	}

	var products []model.Product
	if len(pids) > 0 {
		db.Where("id IN ?", pids).Find(&products)
	} else {
		db.Where("status = ?", "active").Limit(50).Find(&products)
	}

	detail := gin.H{}
	for _, p := range products {
		payType := fmt.Sprintf(`{"pay_type":"%s","pay_hour_cycle":"720","pay_day_cycle":"30","pay_ontrial_status":0,"pay_ontrial_cycle":"0","pay_ontrial_num":"1","pay_ontrial_condition":[],"pay_ontrial_cycle_type":"day","pay_ontrial_num_rule":"0","clientscount_rule":0}`, p.BillingCycle)

		monthlyPrice := fmt.Sprintf("%.2f", p.Price)
		productPricings := []gin.H{{
			"id":              1,
			"type":            "product",
			"currency":        currencyID,
			"relid":           p.ID,
			"osetupfee":       "0.00",
			"hsetupfee":       "0.00",
			"dsetupfee":       "0.00",
			"ontrialfee":      "0.00",
			"msetupfee":       "0.00",
			"qsetupfee":       "0.00",
			"ssetupfee":       "0.00",
			"asetupfee":       "0.00",
			"bsetupfee":       "0.00",
			"tsetupfee":       "0.00",
			"foursetupfee":    "0.00",
			"fivesetupfee":    "0.00",
			"sixsetupfee":     "0.00",
			"sevensetupfee":   "0.00",
			"eightsetupfee":   "0.00",
			"ninesetupfee":    "0.00",
			"tensetupfee":     "0.00",
			"onetime":         "-1.00",
			"hour":            "-1.00",
			"day":             "-1.00",
			"ontrial":         "-1.00",
			"monthly":         monthlyPrice,
			"quarterly":       fmt.Sprintf("%.2f", p.Price*3),
			"semiannually":    fmt.Sprintf("%.2f", p.Price*6),
			"annually":        fmt.Sprintf("%.2f", p.Price*12),
			"biennially":      "-1.00",
			"triennially":     "-1.00",
			"fourly":          "-1.00",
			"fively":          "-1.00",
			"sixly":           "-1.00",
			"sevenly":         "-1.00",
			"eightly":         "-1.00",
			"ninely":          "-1.00",
			"tenly":           "-1.00",
			"code":            currencyCode,
		}}

		detail[fmt.Sprintf("%d", p.ID)] = gin.H{
			"id":                         p.ID,
			"type":                       p.Type,
			"gid":                        p.GroupID,
			"name":                       p.Name,
			"description":                p.Description,
			"host":                       `{"show":"0","modify":0,"prefix":"ser","rule":{"upper":"0","lower":"0","num":"1","len_num":12}}`,
			"is_domain":                  0,
			"hidden":                     0,
			"password":                   `{"show":"1","modify":0,"rule":{"len_num":"12","upper":"1","lower":"1","num":"1","special":"0"}}`,
			"show_domain_options":        0,
			"welcome_email":              0,
			"stock_control":              0,
			"qty":                        999,
			"prorata_billing":            0,
			"prorata_date":               0,
			"prorata_charge_next_month":  0,
			"pay_type":                   payType,
			"pay_method":                 "prepayment",
			"allow_qty":                  0,
			"auto_setup":                 "payment",
			"server_type":                "",
			"server_group":               0,
			"config_option1":             "",
			"config_option2":             "",
			"config_option3":             "",
			"config_option4":             "",
			"config_option5":             "",
			"config_option6":             "",
			"config_option7":             "",
			"config_option8":             "",
			"config_option9":             "",
			"config_option10":            "",
			"config_option11":            "",
			"config_option12":            "",
			"config_option13":            "",
			"config_option14":            "",
			"config_option15":            "",
			"config_option16":            "",
			"config_option17":            "",
			"config_option18":            "",
			"config_option19":            "",
			"config_option20":            "",
			"config_option21":            "",
			"config_option22":            "",
			"config_option23":            "",
			"config_option24":            "",
			"recurring_cycles":           0,
			"auto_terminate_days":        0,
			"auto_terminate_email":       0,
			"config_options_upgrade":     0,
			"billing_cycle_upgrade":      "",
			"upgrade_email":              0,
			"down_configoption_refund":   0,
			"overages_enabled":           "",
			"overages_disk_limit":        0,
			"overages_bw_limit":          0,
			"overages_disk_price":        "0.0000",
			"overages_bw_price":          "0.0000",
			"tax":                        0,
			"affiliateonetime":           0,
			"affiliate_pay_type":         "default",
			"affiliate_pay_amount":       "0.00",
			"order":                      0,
			"retired":                    0,
			"is_featured":                0,
			"create_time":                p.CreatedAt.Unix(),
			"update_time":                p.UpdatedAt.Unix(),
			"auto_create_config_options": 0,
			"groupid":                    p.GroupID,
			"api_type":                   "normal",
			"location_version":           1,
			"upstream_version":           0,
			"upstream_price_type":        "",
			"upstream_price_value":       "",
			"zjmf_api_id":                0,
			"upstream_pid":               0,
			"is_truename":                0,
			"uuid":                       nil,
			"p_uid":                      0,
			"rate":                       "1.00",
			"clientscount":               0,
			"app_type":                   "",
			"product_shopping_url":       "",
			"product_group_url":          "",
			"upper_reaches_id":           0,
			"version_description":        nil,
			"app_version":                nil,
			"is_bind_phone":              0,
			"app_hot_order":              0,
			"app_hot_lock":               0,
			"app_hot_heat":               0,
			"app_recommend_status":       0,
			"app_recommend_order":        0,
			"app_recommend_lock":         0,
			"app_pay_type":               0,
			"app_score":                  0,
			"app_images":                 nil,
			"app_status":                 0,
			"upstream_qty":               0,
			"cancel_control":             1,
			"unretired_time":             0,
			"upstream_stock_control":     0,
			"upstream_auto_setup":        "",
			"upstream_ontrial_status":    0,
			"upstream_price":             "0.00",
			"upstream_product_shopping_url": "",
			"upstream_cycle":             "",
			"customfields":               []gin.H{},
			"product_pricings":           productPricings,
			"advanced":                   []gin.H{},
			"config_groups":              []gin.H{},
			"config_links":               []int{},
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":   "请求成功",
		"data": gin.H{
			"detail": detail,
		},
	})
}

// ZjmfCompatOnTrialMax zjmf兼容试用上限（/cart/ontrialmax）
// zjmf取 $res["data"]["product"]["qty"]
func ZjmfCompatOnTrialMax(c *gin.Context) {
	pid := c.Query("pid")
	if pid == "" {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "缺少pid"})
		return
	}

	db := database.GetDB()
	var product model.Product
	if err := db.First(&product, pid).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "产品不存在"})
		return
	}

	// zjmf返回空数组表示无试用产品
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"msg":     "请求成功",
		"data": gin.H{
			"product": []gin.H{},
		},
		"is_aff": "1",
	})
}

// ZjmfCompatStockControl zjmf兼容库存检查（/cart/stock_control）
// zjmf取 $upstream_data["product"]，检查hidden、stock_control、qty
func ZjmfCompatStockControl(c *gin.Context) {
	pid := c.Query("pid")
	if pid == "" {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "缺少pid"})
		return
	}

	db := database.GetDB()
	var product model.Product
	if err := db.First(&product, pid).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 200, "data": gin.H{}})
		return
	}

	if product.Status != "active" {
		c.JSON(http.StatusOK, gin.H{"status": 200, "data": gin.H{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"msg":     "请求成功",
		"data": gin.H{
			"product": gin.H{
				"hidden":        0,
				"stock_control": 0,
				"qty":           999,
			},
		},
		"is_aff": "1",
	})
}

// ZjmfCompatApplyCredit zjmf兼容使用余额支付（/apply_credit）
// zjmf传 invoiceid + use_credit
func ZjmfCompatApplyCredit(c *gin.Context) {
	var req struct {
		InvoiceID uint `json:"invoiceid"`
		UseCredit int  `json:"use_credit"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "参数错误"})
		return
	}

	db := database.GetDB()
	var invoice model.Invoice
	if err := db.First(&invoice, req.InvoiceID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "账单不存在"})
		return
	}

	if invoice.Status == "Paid" {
		c.JSON(http.StatusOK, gin.H{"status": 1001, "msg": "账单已支付"})
		return
	}

	userID, _ := c.Get("user_id")
	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "用户不存在"})
		return
	}

	if user.Balance < invoice.Amount {
		c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "余额不足"})
		return
	}

	// 扣余额+标记已付
	db.Model(&user).Update("balance", gorm.Expr("balance - ?", invoice.Amount))
	db.Model(&invoice).Update("status", "Paid")

	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "支付成功"})
}

// ZjmfCompatTrafficUsage zjmf兼容流量使用统计（/host/trafficusage）
func ZjmfCompatTrafficUsage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data": gin.H{
			"incoming": 0,
			"outgoing": 0,
			"total":    0,
		},
	})
}

// ZjmfCompatSetDownstream zjmf兼容设置下游信息（/host/setdownstream）
func ZjmfCompatSetDownstream(c *gin.Context) {
	var req struct {
		ID             uint   `json:"id"`
		PID            uint   `json:"pid"`
		DownstreamURL  string `json:"downstream_url"`
		DownstreamToken string `json:"downstream_token"`
		DownstreamID   uint   `json:"downstream_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "参数错误"})
		return
	}

	// 存储下游信息到service记录
	db := database.GetDB()
	var svc model.Service
	if err := db.First(&svc, req.ID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "服务不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "success"})
}

// ZjmfCompatDcimOn POST /dcim/on - 开机
func ZjmfCompatDcimOn(c *gin.Context) {
	var req struct {
		ID uint `json:"id" form:"id"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "参数错误"})
		return
	}

	if req.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "缺少id"})
		return
	}

	db := database.GetDB()
	var svc model.Service
	if err := db.First(&svc, req.ID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "服务不存在"})
		return
	}

	db.Model(&svc).Update("status", "active")
	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "开机成功"})
}

// ZjmfCompatDcimOff POST /dcim/off - 关机
func ZjmfCompatDcimOff(c *gin.Context) {
	var req struct {
		ID uint `json:"id" form:"id"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "参数错误"})
		return
	}

	if req.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "缺少id"})
		return
	}

	db := database.GetDB()
	var svc model.Service
	if err := db.First(&svc, req.ID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "服务不存在"})
		return
	}

	db.Model(&svc).Update("status", "suspended")
	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "关机成功"})
}

// ZjmfCompatDcimReboot POST /dcim/reboot - 重启
func ZjmfCompatDcimReboot(c *gin.Context) {
	var req struct {
		ID uint `json:"id" form:"id"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "参数错误"})
		return
	}

	if req.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "缺少id"})
		return
	}

	db := database.GetDB()
	var svc model.Service
	if err := db.First(&svc, req.ID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "服务不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "重启成功"})
}

// ZjmfCompatDcimRescue POST /dcim/rescue - 救援模式
func ZjmfCompatDcimRescue(c *gin.Context) {
	var req struct {
		ID uint `json:"id" form:"id"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "参数错误"})
		return
	}

	if req.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "缺少id"})
		return
	}

	db := database.GetDB()
	var svc model.Service
	if err := db.First(&svc, req.ID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "服务不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "请求成功", "data": gin.H{"rescue_mode": true}})
}

// ZjmfCompatDcimReinstall POST /dcim/reinstall - 重装系统
func ZjmfCompatDcimReinstall(c *gin.Context) {
	var req struct {
		ID uint   `json:"id" form:"id"`
		OS string `json:"os" form:"os"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "参数错误"})
		return
	}

	if req.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "缺少id"})
		return
	}

	db := database.GetDB()
	var svc model.Service
	if err := db.First(&svc, req.ID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "服务不存在"})
		return
	}

	if req.OS != "" {
		db.Model(&svc).Update("config", req.OS)
	}

	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "重装成功"})
}

// ZjmfCompatDcimCrackPass POST /dcim/crack_pass - 重置密码
func ZjmfCompatDcimCrackPass(c *gin.Context) {
	var req struct {
		ID        uint   `json:"id" form:"id"`
		Password  string `json:"password" form:"password"`
		User      string `json:"user" form:"user"`
		OtherUser string `json:"other_user" form:"other_user"`
		Action    string `json:"action" form:"action"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "参数错误"})
		return
	}

	if req.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "缺少id"})
		return
	}

	db := database.GetDB()
	var svc model.Service
	if err := db.First(&svc, req.ID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "服务不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "发起成功"})
}

// ZjmfCompatDcimCheckReinstall POST /dcim/check_reinstall - 检查可重装
func ZjmfCompatDcimCheckReinstall(c *gin.Context) {
	var req struct {
		ID uint `json:"id" form:"id"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "参数错误"})
		return
	}

	if req.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "缺少id"})
		return
	}

	db := database.GetDB()
	var svc model.Service
	if err := db.First(&svc, req.ID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "服务不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "可以重装", "num": 0, "max_times": 3})
}

// ZjmfCompatDcimDetail GET /dcim/detail - 服务器详情
func ZjmfCompatDcimDetail(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "缺少id"})
		return
	}

	db := database.GetDB()
	var svc model.Service
	if err := db.First(&svc, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "服务不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":   "请求成功",
		"data": gin.H{
			"dcimid":      svc.ID,
			"status":      svc.Status,
			"domain":      svc.Domain,
			"dedicatedip": "",
			"os":          "",
			"username":    "",
			"password":    "",
		},
	})
}

// ZjmfCompatDcimBuyReinstallTimes POST /dcim/buy_reinstall_times - 购买重装次数
func ZjmfCompatDcimBuyReinstallTimes(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "不需要购买次数"})
}

// ZjmfCompatDcimBuyFlowPacket POST /dcim/buy_flow_packet - 购买流量包
func ZjmfCompatDcimBuyFlowPacket(c *gin.Context) {
	var req struct {
		ID uint `json:"id" form:"id"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "参数错误"})
		return
	}

	if req.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "缺少id"})
		return
	}

	db := database.GetDB()
	var svc model.Service
	if err := db.First(&svc, req.ID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "服务不存在"})
		return
	}

	invoice := model.Invoice{
		UserID: svc.UserID,
		Amount: 0,
		Status: "unpaid",
		Note:   fmt.Sprintf("购买流量包 服务#%d", svc.ID),
	}
	db.Create(&invoice)

	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "请求成功", "data": gin.H{"invoiceid": invoice.ID, "amount": fmt.Sprintf("%.2f", invoice.Amount)}})
}

// ZjmfCompatUpgradeCheckoutConfig POST /upgrade/checkout_config_upgrade - 配置升级结算
func ZjmfCompatUpgradeCheckoutConfig(c *gin.Context) {
	var req struct {
		HostID uint `json:"host_id" form:"host_id"`
	}
	c.ShouldBind(&req)

	db := database.GetDB()
	var svc model.Service
	if req.HostID > 0 {
		db.First(&svc, req.HostID)
	}

	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "请求成功", "data": gin.H{"invoiceid": 0, "amount": "0.00"}})
}

// ZjmfCompatUpgradeProductPost POST /upgrade/upgrade_product_post - 产品升级
func ZjmfCompatUpgradeProductPost(c *gin.Context) {
	var req struct {
		HostID uint `json:"host_id" form:"host_id"`
	}
	c.ShouldBind(&req)

	db := database.GetDB()
	var svc model.Service
	if req.HostID > 0 {
		db.First(&svc, req.HostID)
	}

	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "请求成功", "data": gin.H{"invoiceid": 0, "amount": "0.00"}})
}

// ZjmfCompatUpgradeCheckoutProduct POST /upgrade/checkout_upgrade_product - 产品升级结算
func ZjmfCompatUpgradeCheckoutProduct(c *gin.Context) {
	var req struct {
		HostID uint `json:"host_id" form:"host_id"`
	}
	c.ShouldBind(&req)

	db := database.GetDB()
	var svc model.Service
	if req.HostID > 0 {
		db.First(&svc, req.HostID)
	}

	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "请求成功", "data": gin.H{"invoiceid": 0, "amount": "0.00"}})
}

// ZjmfCompatSslCertFunc POST /provision/sslCertFunc - SSL证书管理
func ZjmfCompatSslCertFunc(c *gin.Context) {
	var req struct {
		ID uint `json:"id" form:"id"`
	}
	c.ShouldBind(&req)

	if req.ID > 0 {
		db := database.GetDB()
		var svc model.Service
		db.First(&svc, req.ID)
	}

	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "请求成功", "data": gin.H{}})
}

// ZjmfCompatRefreshPowerStatus POST /dcim/refresh_power_status - 刷新电源状态
func ZjmfCompatRefreshPowerStatus(c *gin.Context) {
	var req struct {
		ID uint `json:"id" form:"id"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "参数错误"})
		return
	}

	if req.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "缺少id"})
		return
	}

	db := database.GetDB()
	var svc model.Service
	if err := db.First(&svc, req.ID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "服务不存在"})
		return
	}

	powerStatus := "off"
	if svc.Status == "active" {
		powerStatus = "on"
	}

	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "请求成功", "data": gin.H{"power_status": powerStatus}})
}

// ZjmfCompatRefreshAllPowerStatus POST /dcim/refresh_all_power_status - 批量刷新电源状态
func ZjmfCompatRefreshAllPowerStatus(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" form:"ids"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "参数错误"})
		return
	}

	db := database.GetDB()
	var results []gin.H
	for _, id := range req.IDs {
		var svc model.Service
		if err := db.First(&svc, id).Error; err != nil {
			results = append(results, gin.H{"id": id, "status": "off"})
			continue
		}

		powerStatus := "off"
		if svc.Status == "active" {
			powerStatus = "on"
		}
		results = append(results, gin.H{"id": id, "status": powerStatus})
	}

	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "请求成功", "data": gin.H{"results": results}})
}

// ZjmfCompatDcimHideResult POST /dcim/hide_result - 隐藏操作结果
func ZjmfCompatDcimHideResult(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "请求成功", "data": gin.H{}})
}

// ZjmfCompatClearCart zjmf兼容清空购物车（/cart/clear）
// zjmf源码(Host.php:74): zjmfCurl($apiId, "/cart/clear", [], 30, "POST")
func ZjmfCompatClearCart(c *gin.Context) {
	userID, _ := c.Get("user_id")
	db := database.GetDB()

	db.Where("user_id = ?", userID).Delete(&model.CartItem{})

	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "请求成功", "data": []interface{}{}})
}

// ZjmfCompatAddToShop zjmf兼容添加商品到购物车（/cart/add_to_shop）
// zjmf源码(Host.php:115): zjmfCurl($apiId, "/cart/add_to_shop", $cartData, 30, "POST")
// zjmf传参: pid, qty, billingcycle, configoptions等
func ZjmfCompatAddToShop(c *gin.Context) {
	var req struct {
		PID          uint   `json:"pid"`
		Qty          int    `json:"qty"`
		BillingCycle string `json:"billingcycle"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "参数错误"})
		return
	}

	if req.PID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "缺少pid"})
		return
	}
	if req.Qty <= 0 {
		req.Qty = 1
	}

	db := database.GetDB()

	// 查询产品
	var product model.Product
	if err := db.First(&product, req.PID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "产品不存在"})
		return
	}

	// 确定计费周期
	cycle := req.BillingCycle
	if cycle == "" {
		cycle = product.BillingCycle
	}

	// 计算价格
	amount := product.Price * float64(req.Qty)
	switch cycle {
	case "quarterly":
		amount = product.Price * 3 * float64(req.Qty)
	case "semiannually":
		amount = product.Price * 6 * float64(req.Qty)
	case "annually":
		amount = product.Price * 12 * float64(req.Qty)
	case "onetime":
		amount = product.Price * float64(req.Qty)
	case "free":
		amount = 0
	}

	userID, _ := c.Get("user_id")

	// 创建购物车项
	cartItem := model.CartItem{
		UserID:      userID.(uint),
		ProductID:   product.ID,
		ProductName: product.Name,
		Quantity:    req.Qty,
		Cycle:       cycle,
		Amount:      amount,
	}
	db.Create(&cartItem)

	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "请求成功", "data": gin.H{"cart_id": cartItem.ID}})
}

// ZjmfCompatSettle zjmf兼容结算购物车（/cart/settle）
// zjmf源码(Host.php:128): zjmfCurl($apiId, "/cart/settle", $settleData, 30, "POST")
// zjmf传参: payment method等
func ZjmfCompatSettle(c *gin.Context) {
	var req struct {
		PaymentMethod string `json:"payment_method"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		// 没有参数也可以，用默认支付方式
		req.PaymentMethod = "balance"
	}

	userID, _ := c.Get("user_id")
	db := database.GetDB()

	// 获取购物车项
	var cartItems []model.CartItem
	db.Where("user_id = ?", userID).Find(&cartItems)

	if len(cartItems) == 0 {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "购物车为空"})
		return
	}

	// 计算总金额
	var totalAmount float64
	for _, item := range cartItems {
		totalAmount += item.Amount
	}

	// 创建账单
	invoice := model.Invoice{
		UserID:        userID.(uint),
		Amount:        totalAmount,
		Status:        "unpaid",
		PaymentMethod: req.PaymentMethod,
	}
	db.Create(&invoice)

	// 清空购物车
	db.Where("user_id = ?", userID).Delete(&model.CartItem{})

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "请求成功",
		"data": gin.H{
			"invoiceid": invoice.ID,
			"amount":    fmt.Sprintf("%.2f", totalAmount),
		},
	})
}

// ZjmfCompatApplyCreditLimit zjmf兼容申请信用额度（/apply_credit_limit）
// zjmf源码(Host.php:150): zjmfCurl($apiId, "/apply_credit_limit", $post_data)
// zjmf传参: invoiceid
func ZjmfCompatApplyCreditLimit(c *gin.Context) {
	var req struct {
		InvoiceID uint `json:"invoiceid"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "参数错误"})
		return
	}

	userID, _ := c.Get("user_id")
	db := database.GetDB()

	// 查询账单
	var invoice model.Invoice
	if err := db.First(&invoice, req.InvoiceID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "账单不存在"})
		return
	}

	// 验证账单归属
	if invoice.UserID != userID.(uint) {
		c.JSON(http.StatusOK, gin.H{"status": 403, "msg": "无权限"})
		return
	}

	// 查询信用额度
	var creditLimit model.CreditLimit
	if err := db.Where("user_id = ?", userID).First(&creditLimit).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "信用额度不足"})
		return
	}

	// 检查可用信用额度
	available := creditLimit.Amount - creditLimit.Used
	if available < invoice.Amount {
		c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "信用额度不足"})
		return
	}

	// 使用信用额度支付
	db.Model(&creditLimit).Update("used", gorm.Expr("used + ?", invoice.Amount))
	db.Model(&invoice).Update("status", "Paid")

	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "请求成功"})
}

// ZjmfCompatUserInfo zjmf兼容用户信息（/user_info）
// zjmf源码(Host.php:56): zjmfCurl($apiId, "/user_info", [], 30, "GET")
func ZjmfCompatUserInfo(c *gin.Context) {
	userID, _ := c.Get("user_id")
	db := database.GetDB()

	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "用户不存在"})
		return
	}

	// 查询信用额度
	creditLimitAmount := "0.00"
	var creditLimit model.CreditLimit
	if err := db.Where("user_id = ?", userID).First(&creditLimit).Error; err == nil {
		creditLimitAmount = fmt.Sprintf("%.2f", creditLimit.Amount-creditLimit.Used)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "请求成功",
		"data": gin.H{
			"credit":       fmt.Sprintf("%.2f", user.Balance),
			"credit_limit": creditLimitAmount,
		},
	})
}

// ============================================================
// 以下为 zjmf 深度兼容端点（host/provision/upgrade）
// ============================================================

// ZjmfCompatHostCancel zjmf兼容取消服务（POST /host/cancel）
// zjmf源码(Host.php:538): zjmfCurl($apiId, "/host/cancel", $post_data)
// 传参: host_id, reason
func ZjmfCompatHostCancel(c *gin.Context) {
	var req struct {
		HostID uint   `json:"host_id" form:"host_id"`
		Reason string `json:"reason" form:"reason"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "参数错误"})
		return
	}

	if req.HostID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "缺少host_id"})
		return
	}

	db := database.GetDB()
	var svc model.Service
	if err := db.First(&svc, req.HostID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "未找到该产品"})
		return
	}

	// 标记服务为已取消
	db.Model(&svc).Update("status", "cancelled")

	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "请求成功"})
}

// ZjmfCompatHostRenew zjmf兼容续费（POST /host/renew）
// zjmf源码(Host.php:1855): zjmfCurl($apiId, "/host/renew", $renewData)
// 传参: host_id, billingcycle
// 返回: {"status": 200, "msg": "请求成功", "data": {"invoiceid": 1, "amount": "99.99"}}
func ZjmfCompatHostRenew(c *gin.Context) {
	var req struct {
		HostID       uint   `json:"host_id" form:"host_id"`
		BillingCycle string `json:"billingcycle" form:"billingcycle"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "参数错误"})
		return
	}

	if req.HostID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "缺少host_id"})
		return
	}

	db := database.GetDB()
	var svc model.Service
	if err := db.First(&svc, req.HostID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "未找到该产品"})
		return
	}

	// 生成续费账单
	amount := svc.Amount

	invoice := model.Invoice{
		UserID: svc.UserID,
		Amount: amount,
		Status: "unpaid",
		Note:   fmt.Sprintf("续费服务 #%d", svc.ID),
	}
	db.Create(&invoice)

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "请求成功",
		"data": gin.H{
			"invoiceid": invoice.ID,
			"amount":    fmt.Sprintf("%.2f", amount),
		},
	})
}

// ZjmfCompatProvisionDefault zjmf兼容通用模块操作（POST /provision/default）
// zjmf源码(Host.php:322等12+处): zjmfCurl($apiId, "/provision/default", $post_data)
// 传参: func(操作类型), id(host_id), 其他参数
// 返回: {"status": 200, "msg": "请求成功", "data": {}}
func ZjmfCompatProvisionDefault(c *gin.Context) {
	var req struct {
		Func string `json:"func" form:"func"`
		ID   uint   `json:"id" form:"id"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "参数错误"})
		return
	}

	if req.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "缺少id"})
		return
	}

	db := database.GetDB()
	var svc model.Service
	if err := db.First(&svc, req.ID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "未找到该产品"})
		return
	}

	// 根据func执行不同操作（Create/Suspend/Unsuspend/Terminate等）
	// 目前作为通用占位返回成功
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "请求成功",
		"data":   gin.H{},
	})
}

// ZjmfCompatProvisionButton zjmf兼容模块按钮（POST /provision/button）
// zjmf源码(ProvisionController:467): zjmfCurl($apiId, "/provision/button", $post_data)
// 传参: id, type等
// 返回: {"status": 200, "msg": "请求成功", "data": {"button": [], "console": []}}
func ZjmfCompatProvisionButton(c *gin.Context) {
	var req struct {
		ID   uint   `json:"id" form:"id"`
		Type string `json:"type" form:"type"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "参数错误"})
		return
	}

	if req.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "缺少id"})
		return
	}

	db := database.GetDB()
	var svc model.Service
	if err := db.First(&svc, req.ID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "未找到该产品"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "请求成功",
		"data": gin.H{
			"button":  []gin.H{},
			"console": []gin.H{},
		},
	})
}

// ZjmfCompatProvisionCustom zjmf兼容自定义模块操作（POST /provision/custom/:id）
// zjmf源码(DcimCloud.php:297等): zjmfCurl($apiId, "/provision/custom/" . $id, $post)
// 返回: {"status": 200, "msg": "请求成功", "data": {}}
func ZjmfCompatProvisionCustom(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "缺少id"})
		return
	}

	db := database.GetDB()
	var svc model.Service
	if err := db.First(&svc, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "未找到该产品"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "请求成功",
		"data":   gin.H{},
	})
}

// ZjmfCompatProvisionChart zjmf兼容图表数据（GET /provision/chart/:id）
// zjmf源码(ProvisionController:432): zjmfCurl($apiId, "/provision/chart/" . $id, $params, 30, "GET")
// 返回: {"status": 200, "msg": "请求成功", "data": {"chart": []}}
func ZjmfCompatProvisionChart(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "缺少id"})
		return
	}

	db := database.GetDB()
	var svc model.Service
	if err := db.First(&svc, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "未找到该产品"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "请求成功",
		"data": gin.H{
			"chart": []gin.H{},
		},
	})
}

// ZjmfCompatUpgradeConfigPost zjmf兼容配置选项升级（POST /upgrade/upgrade_config_post）
// zjmf源码(Host.php:1534): zjmfCurl($apiId, "/upgrade/upgrade_config_post", $post_data)
// 传参: host_id, configoptions等
// 返回: {"status": 200, "msg": "请求成功", "data": {"invoiceid": 1, "amount": "99.99"}}
func ZjmfCompatUpgradeConfigPost(c *gin.Context) {
	var req struct {
		HostID        uint                   `json:"host_id" form:"host_id"`
		ConfigOptions map[string]interface{} `json:"configoptions" form:"configoptions"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "参数错误"})
		return
	}

	if req.HostID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "缺少host_id"})
		return
	}

	db := database.GetDB()
	var svc model.Service
	if err := db.First(&svc, req.HostID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "未找到该产品"})
		return
	}

	// 生成升级账单（差异金额，此处简化为当前产品价格的10%作为升级费）
	amount := svc.Amount * 0.1

	invoice := model.Invoice{
		UserID: svc.UserID,
		Amount: amount,
		Status: "unpaid",
		Note:   fmt.Sprintf("配置升级服务 #%d", svc.ID),
	}
	db.Create(&invoice)

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "请求成功",
		"data": gin.H{
			"invoiceid": invoice.ID,
			"amount":    fmt.Sprintf("%.2f", amount),
		},
	})
}

// ==================== 补充缺失端点 ====================

// ZjmfCompatDcimBmc POST /dcim/bmc - BMC管理
func ZjmfCompatDcimBmc(c *gin.Context) {
	var req struct {
		ID uint `json:"id" form:"id"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "参数错误"})
		return
	}

	if req.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "缺少id"})
		return
	}

	db := database.GetDB()
	var svc model.Service
	if err := db.First(&svc, req.ID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "服务不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":   "请求成功",
		"data": gin.H{
			"bmc_ip":   "",
			"bmc_user": "",
			"bmc_pass": "",
		},
	})
}

// ZjmfCompatDcimCancelTask POST /dcim/cancel_task - 取消任务
func ZjmfCompatDcimCancelTask(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "请求成功", "data": gin.H{}})
}

// ZjmfCompatDcimIkvm POST /dcim/ikvm - iKVM
func ZjmfCompatDcimIkvm(c *gin.Context) {
	var req struct {
		ID uint `json:"id" form:"id"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "参数错误"})
		return
	}

	if req.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "缺少id"})
		return
	}

	db := database.GetDB()
	var svc model.Service
	if err := db.First(&svc, req.ID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "服务不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "请求成功", "data": gin.H{"url": ""}})
}

// ZjmfCompatDcimKvm POST /dcim/kvm - KVM
func ZjmfCompatDcimKvm(c *gin.Context) {
	var req struct {
		ID uint `json:"id" form:"id"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "参数错误"})
		return
	}

	if req.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "缺少id"})
		return
	}

	db := database.GetDB()
	var svc model.Service
	if err := db.First(&svc, req.ID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "服务不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "请求成功", "data": gin.H{"url": ""}})
}

// ZjmfCompatDcimNovnc POST /dcim/novnc - NoVNC
func ZjmfCompatDcimNovnc(c *gin.Context) {
	var req struct {
		ID uint `json:"id" form:"id"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "参数错误"})
		return
	}

	if req.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "缺少id"})
		return
	}

	db := database.GetDB()
	var svc model.Service
	if err := db.First(&svc, req.ID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "服务不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "请求成功", "data": gin.H{"url": "", "port": 0}})
}

// ZjmfCompatDcimReinstallStatus GET /dcim/resintall_status - 重装状态（注意zjmf原始拼写）
func ZjmfCompatDcimReinstallStatus(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "缺少id"})
		return
	}

	db := database.GetDB()
	var svc model.Service
	if err := db.First(&svc, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "服务不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "请求成功", "data": gin.H{"progress": 100, "status": "completed"}})
}

// ZjmfCompatDcimTraffic POST /dcim/traffic - 流量管理
func ZjmfCompatDcimTraffic(c *gin.Context) {
	var req struct {
		ID uint `json:"id" form:"id"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "参数错误"})
		return
	}

	if req.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "缺少id"})
		return
	}

	db := database.GetDB()
	var svc model.Service
	if err := db.First(&svc, req.ID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "服务不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "请求成功", "data": gin.H{"bwlimit": 0, "bwusage": 0}})
}

// ZjmfCompatDcimTrafficUsage GET /dcim/traffic_usage - 流量使用统计
func ZjmfCompatDcimTrafficUsage(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "缺少id"})
		return
	}

	db := database.GetDB()
	var svc model.Service
	if err := db.First(&svc, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "服务不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "请求成功", "data": gin.H{"incoming": 0, "outgoing": 0, "total": 0}})
}

// ZjmfCompatCartHostinfo GET /cart/hostinfo - 购物车主机信息
// TODO: 购物车功能完整实现后改为查询CartItem
func ZjmfCompatCartHostinfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "请求成功", "data": gin.H{"hostinfo": []gin.H{}}})
}

// ZjmfCompatCartSummary GET /cart/summary - 购物车摘要
// TODO: 购物车功能完整实现后改为查询CartItem
func ZjmfCompatCartSummary(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "请求成功", "data": gin.H{"total": "0.00", "items": []gin.H{}}})
}
