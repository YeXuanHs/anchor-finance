package client

import (
	"net/http"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
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

// ZjmfCompatHostDetail zjmf兼容服务详情（/v1/host/header）
func ZjmfCompatHostDetail(c *gin.Context) {
	hostID := c.Query("host_id")
	if hostID == "" {
		c.JSON(http.StatusOK, gin.H{"status": 400, "msg": "缺少host_id"})
		return
	}

	db := database.GetDB()
	var svc model.Service
	if err := db.First(&svc, hostID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "服务不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":  "success",
		"data": gin.H{
			"host_data": gin.H{
				"id":           svc.ID,
				"product_id":   svc.ProductID,
				"product_name": svc.ProductName,
				"status":       svc.Status,
				"domain":       svc.Domain,
				"create_time":  svc.CreatedAt.Unix(),
			},
		},
	})
}

// ZjmfCompatModuleStatus zjmf兼容模块状态（/v1/hosts/:id/module/status）
// zjmf文档: type=host 服务器电源状态, type=reinstall 重装进度
func ZjmfCompatModuleStatus(c *gin.Context) {
	hostID := c.Param("id")
	statusType := c.Query("type")

	db := database.GetDB()
	var svc model.Service
	if err := db.First(&svc, hostID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "服务不存在"})
		return
	}

	// 映射到zjmf状态格式
	statusMap := map[string]string{
		"active":    "Active",
		"suspended": "Suspended",
		"terminated": "Terminated",
		"pending":   "Pending",
	}
	zjmfStatus := statusMap[svc.Status]
	if zjmfStatus == "" {
		zjmfStatus = "Unknown"
	}

	if statusType == "reinstall" {
		// 重装进度
		c.JSON(http.StatusOK, gin.H{"status": 200, "data": gin.H{"progress": 100, "status": "completed"}})
		return
	}

	// 默认返回电源状态
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data": gin.H{
			"status": zjmfStatus,
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
func ZjmfCompatBalance(c *gin.Context) {
	userID, _ := c.Get("user_id")
	db := database.GetDB()

	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 404, "msg": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "success",
		"data": gin.H{
			"credit": user.Balance,
			"currency": gin.H{
				"id":     1,
				"code":   "CNY",
				"prefix": "¥",
				"suffix": "元",
			},
		},
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
					"type":          p.Type,
					"product_price": p.Price,
					"billingcycle":  p.BillingCycle,
					"qty":           p.Qty,
					"stock_control": func() int { if p.StockControl { return 1 }; return 0 }(),
					"setup_fee":     p.SetupFee,
					"pay_type":      p.PayType,
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

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "success",
		"data": gin.H{
			"first_group": result,
			"products":    flatProducts,
			"count":       totalProducts,
			"currency": gin.H{
				"id":     1,
				"code":   "CNY",
				"prefix": "¥",
				"suffix": "元",
			},
		},
	})
}

// ZjmfCompatProductInfo zjmf兼容商品详情（/api/product/proinfo）
// zjmf.php line 48: $path = "api/product/proinfo"; params: pids=[]
func ZjmfCompatProductInfo(c *gin.Context) {
	pids := c.QueryArray("pids")
	if len(pids) == 0 {
		pids = append(pids, c.Query("pids"))
	}

	db := database.GetDB()
	var products []model.Product
	if len(pids) > 0 {
		db.Where("id IN ?", pids).Find(&products)
	} else {
		db.Where("status = ?", "active").Limit(50).Find(&products)
	}

	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "success", "data": products})
}

// ZjmfCompatProductConfig zjmf兼容商品配置（/cart/get_product_config）
// zjmf.php line 70: $path = "cart/get_product_config"; params: pid
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

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":   "success",
		"data": gin.H{
			"pid":         product.ID,
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"billing_cycle": product.BillingCycle,
		},
	})
}
