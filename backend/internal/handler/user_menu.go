package handler

import (
	"anchorfinance/internal/model"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserMenuHandler 用户菜单处理器
type UserMenuHandler struct {
	db *gorm.DB
}

func NewUserMenuHandler(db *gorm.DB) *UserMenuHandler {
	return &UserMenuHandler{db: db}
}

// GetUserMenus 获取用户中心菜单
func (h *UserMenuHandler) GetUserMenus(c *gin.Context) {
	// 从数据库读取菜单
	menus, err := model.GetNavTree(h.db, 1) // 1=用户中心
	if err != nil || len(menus) == 0 {
		// 如果数据库没有菜单，返回默认菜单
		response.Success(c, h.getDefaultMenus())
		return
	}

	// 根据配置过滤菜单
	menus = h.filterByConfig(menus)

	response.Success(c, menus)
}

// GetTopNav 获取www顶部导航
func (h *UserMenuHandler) GetTopNav(c *gin.Context) {
	menus, err := model.GetNavTree(h.db, 2) // 2=www头部
	if err != nil || len(menus) == 0 {
		response.Success(c, h.getDefaultTopNav())
		return
	}
	response.Success(c, menus)
}

// GetBottomNav 获取www底部导航
func (h *UserMenuHandler) GetBottomNav(c *gin.Context) {
	menus, err := model.GetNavTree(h.db, 3) // 3=www尾部
	if err != nil || len(menus) == 0 {
		response.Success(c, h.getDefaultBottomNav())
		return
	}
	response.Success(c, menus)
}

// filterByConfig 根据配置过滤菜单
func (h *UserMenuHandler) filterByConfig(menus []*model.MenuItem) []*model.MenuItem {
	var result []*model.MenuItem

	for _, m := range menus {
		// 检查交易市场是否启用
		if m.URL == "/user/marketplace" || m.URL == "marketplace" {
			if !h.isMarketplaceEnabled() {
				continue
			}
		}

		// 递归过滤子菜单
		if len(m.Children) > 0 {
			m.Children = h.filterByConfig(m.Children)
		}

		result = append(result, m)
	}
	return result
}

// isMarketplaceEnabled 检查交易市场是否启用
func (h *UserMenuHandler) isMarketplaceEnabled() bool {
	var config model.MarketplaceConfig
	if err := h.db.First(&config).Error; err != nil {
		return false
	}
	return config.Enabled
}

// getDefaultMenus 获取默认菜单（当数据库没有数据时）
func (h *UserMenuHandler) getDefaultMenus() []*model.MenuItem {
	return []*model.MenuItem{
		{
			ID:     1,
			Name:   "控制台",
			FaIcon: "bx bx-home-circle",
			URL:    "/user/dashboard",
		},
		{
			ID:     2,
			Name:   "产品与服务",
			FaIcon: "bx bxs-grid-alt",
			Children: []*model.MenuItem{
				{ID: 8, Name: "订购产品", URL: "/products"},
				{ID: 9, Name: "我的服务", URL: "/user/products"},
				{ID: 10, Name: "订单管理", URL: "/user/orders"},
			},
		},
		{
			ID:     3,
			Name:   "账户管理",
			FaIcon: "bx bx-user",
			Children: []*model.MenuItem{
				{ID: 12, Name: "个人信息", URL: "/user/profile"},
				{ID: 13, Name: "安全中心", URL: "/user/security"},
				{ID: 14, Name: "实名认证", URL: "/user/verification"},
				{ID: 15, Name: "消息中心", URL: "/user/system-message"},
				{ID: 16, Name: "联系人管理", URL: "/user/contacts"},
				{ID: 17, Name: "第三方登录", URL: "/user/oauth-bind"},
			},
		},
		{
			ID:     4,
			Name:   "财务管理",
			FaIcon: "bx bx-dollar-circle",
			Children: []*model.MenuItem{
				{ID: 19, Name: "账单列表", URL: "/user/invoices"},
				{ID: 24, Name: "账户充值", URL: "/user/wallet"},
				{ID: 25, Name: "优惠券", URL: "/user/coupons"},
			},
		},
		{
			ID:     5,
			Name:   "技术支持",
			FaIcon: "bx bx-detail",
			Children: []*model.MenuItem{
				{ID: 26, Name: "工单列表", URL: "/user/tickets"},
				{ID: 27, Name: "提交工单", URL: "/user/tickets/create"},
				{ID: 28, Name: "帮助中心", URL: "/knowledge-base"},
				{ID: 29, Name: "资源下载", URL: "/downloads"},
				{ID: 30, Name: "新闻中心", URL: "/news"},
			},
		},
		{
			ID:     7,
			Name:   "推介计划",
			FaIcon: "bx bxs-paper-plane",
			URL:    "/user/referral",
		},
	}
}

// getDefaultTopNav 获取默认顶部导航
func (h *UserMenuHandler) getDefaultTopNav() []*model.MenuItem {
	return []*model.MenuItem{
		{ID: 26, Name: "首页", URL: "/"},
		{
			ID:   27,
			Name: "产品",
			Children: []*model.MenuItem{
				{ID: 28, Name: "云服务器", URL: "/products?group=cloud"},
				{ID: 29, Name: "独立服务器", URL: "/products?group=dedicated"},
				{ID: 30, Name: "全部产品", URL: "/products"},
			},
		},
		{ID: 31, Name: "解决方案", URL: "/solutions"},
		{ID: 32, Name: "新闻动态", URL: "/news"},
		{
			ID:   33,
			Name: "帮助支持",
			Children: []*model.MenuItem{
				{ID: 34, Name: "帮助中心", URL: "/help"},
				{ID: 35, Name: "知识库", URL: "/knowledge-base"},
				{ID: 36, Name: "下载中心", URL: "/downloads"},
				{ID: 37, Name: "联系我们", URL: "/contact"},
			},
		},
	}
}

// getDefaultBottomNav 获取默认底部导航
func (h *UserMenuHandler) getDefaultBottomNav() []*model.MenuItem {
	return []*model.MenuItem{
		{
			ID:   38,
			Name: "产品服务",
			Children: []*model.MenuItem{},
		},
		{
			ID:   40,
			Name: "帮助支持",
			Children: []*model.MenuItem{
				{ID: 41, Name: "帮助中心", URL: "/help"},
				{ID: 42, Name: "知识库", URL: "/knowledge-base"},
				{ID: 43, Name: "下载中心", URL: "/downloads"},
				{ID: 44, Name: "联系我们", URL: "/contact"},
			},
		},
		{
			ID:   45,
			Name: "关于我们",
			Children: []*model.MenuItem{
				{ID: 46, Name: "公司介绍", URL: "/about"},
				{ID: 47, Name: "新闻动态", URL: "/news"},
				{ID: 48, Name: "解决方案", URL: "/solutions"},
			},
		},
	}
}
