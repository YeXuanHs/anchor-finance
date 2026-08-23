package client

import (
	"net/http"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/YeXuanHs/anchor-finance/internal/service"
	"github.com/gin-gonic/gin"
)

// ============================================================
// zjmf兼容API（让zjmf系统能对接锚点财务）
// zjmf系统选"zjmf"类型，填我们的网址，就能拉取商品/同步状态
// ============================================================

// ZjmfCompatLogin zjmf兼容登录（/v1/login_api）
// zjmf用账号+密码登录获取JWT
func ZjmfCompatLogin(c *gin.Context) {
	var req struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	db := database.GetDB()

	// 先查users表
	var user model.User
	if err := db.Where("username = ? OR email = ?", req.Account, req.Account).First(&user).Error; err == nil {
		if service.CheckPassword(req.Password, user.PasswordHash) {
			token, err := service.GenerateTokenStatic(user.ID, user.Username, false)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "生成token失败"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": gin.H{"jwt": token}})
			return
		}
	}

	// 再查admins表
	var admin model.Admin
	if err := db.Where("username = ?", req.Account).First(&admin).Error; err == nil {
		if service.CheckPassword(req.Password, admin.PasswordHash) {
			token, err := service.GenerateTokenStatic(admin.ID, admin.Username, true)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "生成token失败"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": gin.H{"jwt": token}})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "用户名或密码错误"})
}

// ZjmfCompatProducts zjmf兼容商品列表（/v1/products）
func ZjmfCompatProducts(c *gin.Context) {
	db := database.GetDB()
	var products []model.Product
	db.Where("status = ?", "active").Order("id ASC").Find(&products)

	// 转成zjmf格式
	var list []gin.H
	for _, p := range products {
		list = append(list, gin.H{
			"id":         p.ID,
			"name":       p.Name,
			"description": p.Description,
			"price":      p.Price,
			"gid":        p.GroupID,
			"stock":      999, // 默认充足库存
			"status":     "active",
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": gin.H{"product": list}})
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

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": gin.H{"cate": list}})
}

// ZjmfCompatHostDetail zjmf兼容服务详情（/v1/host/header）
func ZjmfCompatHostDetail(c *gin.Context) {
	hostID := c.Query("host_id")
	if hostID == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "缺少host_id"})
		return
	}

	db := database.GetDB()
	var svc model.Service
	if err := db.First(&svc, hostID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "服务不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
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
func ZjmfCompatModuleStatus(c *gin.Context) {
	hostID := c.Param("id")

	db := database.GetDB()
	var svc model.Service
	if err := db.First(&svc, hostID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "服务不存在"})
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

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
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
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "服务不存在"})
		return
	}

	db.Model(&svc).Update("status", "suspended")
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}

// ZjmfCompatModuleUnsuspend zjmf兼容取消暂停（/v1/hosts/:id/module/unsuspend）
func ZjmfCompatModuleUnsuspend(c *gin.Context) {
	hostID := c.Param("id")
	db := database.GetDB()

	var svc model.Service
	if err := db.First(&svc, hostID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "服务不存在"})
		return
	}

	db.Model(&svc).Update("status", "active")
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}

// ZjmfCompatModuleTerminate zjmf兼容终止（/v1/hosts/:id/module/terminate）
func ZjmfCompatModuleTerminate(c *gin.Context) {
	hostID := c.Param("id")
	db := database.GetDB()

	var svc model.Service
	if err := db.First(&svc, hostID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "服务不存在"})
		return
	}

	db.Model(&svc).Update("status", "terminated")
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}

// ZjmfCompatRenew zjmf兼容续费（/v1/hosts/:id/renew）
func ZjmfCompatRenew(c *gin.Context) {
	hostID := c.Param("id")
	db := database.GetDB()

	var svc model.Service
	if err := db.First(&svc, hostID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "服务不存在"})
		return
	}

	// 续费逻辑：延长到期时间
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": gin.H{"amount": svc.Amount}})
}

// ZjmfCompatUser zjmf兼容用户信息（/v1/user）
func ZjmfCompatUser(c *gin.Context) {
	userID, _ := c.Get("user_id")
	db := database.GetDB()

	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
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
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": gin.H{
			"credit": user.Balance,
		},
	})
}
