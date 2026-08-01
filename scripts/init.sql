-- ============================================
-- 锚点财务 默认数据
-- 安装时自动导入
-- ============================================

-- 菜单激活配置
INSERT INTO `menu_actives` (`menu_type`, `menuid`) VALUES (1, 1);
INSERT INTO `menu_actives` (`menu_type`, `menuid`) VALUES (2, 1);
INSERT INTO `menu_actives` (`menu_type`, `menuid`) VALUES (3, 1);

-- ============================================
-- 用户中心菜单 (menu_type=1)
-- ============================================
INSERT INTO `navs` (`name`, `url`, `parent_id`, `order`, `fa_icon`, `menu_type`, `nav_type`, `menu_id`, `is_display`) VALUES
('控制台', '/user/dashboard', 0, 0, 'bx bx-home-circle', 1, 0, 1, 1),
('产品与服务', '', 0, 1, 'bx bxs-grid-alt', 1, 0, 1, 1),
('订购产品', '/products', 2, 0, '', 1, 0, 1, 1),
('我的服务', '/user/products', 2, 1, '', 1, 0, 1, 1),
('订单管理', '/user/orders', 2, 2, '', 1, 0, 1, 1),
('产品升降级', '/user/upgrade', 2, 3, '', 1, 0, 1, 1),
('账户管理', '', 0, 2, 'bx bx-user', 1, 0, 1, 1),
('个人信息', '/user/profile', 7, 0, '', 1, 0, 1, 1),
('安全中心', '/user/security', 7, 1, '', 1, 0, 1, 1),
('实名认证', '/user/verification', 7, 2, '', 1, 0, 1, 1),
('消息中心', '/user/system-message', 7, 3, '', 1, 0, 1, 1),
('联系人管理', '/user/contacts', 7, 4, '', 1, 0, 1, 1),
('第三方登录', '/user/oauth-bind', 7, 5, '', 1, 0, 1, 1),
('财务管理', '', 0, 3, 'bx bx-dollar-circle', 1, 0, 1, 1),
('账单列表', '/user/invoices', 14, 0, '', 1, 0, 1, 1),
('账户充值', '/user/wallet', 14, 1, '', 1, 0, 1, 1),
('优惠券', '/user/coupons', 14, 2, '', 1, 0, 1, 1),
('技术支持', '', 0, 4, 'bx bx-detail', 1, 0, 1, 1),
('工单列表', '/user/tickets', 18, 0, '', 1, 0, 1, 1),
('提交工单', '/user/tickets/create', 18, 1, '', 1, 0, 1, 1),
('帮助中心', '/knowledge-base', 18, 2, '', 1, 0, 1, 1),
('资源下载', '/downloads', 18, 3, '', 1, 0, 1, 1),
('新闻中心', '/news', 18, 4, '', 1, 0, 1, 1),
('推介计划', '/user/referral', 0, 5, 'bx bxs-paper-plane', 1, 0, 1, 1),
('交易市场', '/user/marketplace', 0, 6, 'bx bx-store', 1, 0, 1, 1);

-- ============================================
-- www顶部导航 (menu_type=2)
-- ============================================
INSERT INTO `navs` (`name`, `url`, `parent_id`, `order`, `fa_icon`, `menu_type`, `nav_type`, `menu_id`, `is_display`) VALUES
('首页', '/', 0, 0, '', 2, 0, 1, 1),
('产品', '', 0, 1, '', 2, 0, 1, 1),
('云服务器', '/products?group=cloud', 27, 0, '', 2, 0, 1, 1),
('独立服务器', '/products?group=dedicated', 27, 1, '', 2, 0, 1, 1),
('全部产品', '/products', 27, 2, '', 2, 0, 1, 1),
('解决方案', '/solutions', 0, 2, '', 2, 0, 1, 1),
('新闻动态', '/news', 0, 3, '', 2, 0, 1, 1),
('帮助支持', '', 0, 4, '', 2, 0, 1, 1),
('帮助中心', '/help', 33, 0, '', 2, 0, 1, 1),
('知识库', '/knowledge-base', 33, 1, '', 2, 0, 1, 1),
('下载中心', '/downloads', 33, 2, '', 2, 0, 1, 1),
('联系我们', '/contact', 33, 3, '', 2, 0, 1, 1);

-- ============================================
-- www底部导航 (menu_type=3)
-- ============================================
INSERT INTO `navs` (`name`, `url`, `parent_id`, `order`, `fa_icon`, `menu_type`, `nav_type`, `menu_id`, `is_display`) VALUES
('产品服务', '', 0, 0, '', 3, 0, 1, 1),
('帮助支持', '', 0, 1, '', 3, 0, 1, 1),
('帮助中心', '/help', 40, 0, '', 3, 0, 1, 1),
('知识库', '/knowledge-base', 40, 1, '', 3, 0, 1, 1),
('下载中心', '/downloads', 40, 2, '', 3, 0, 1, 1),
('联系我们', '/contact', 40, 3, '', 3, 0, 1, 1),
('关于我们', '', 0, 2, '', 3, 0, 1, 1),
('公司介绍', '/about', 45, 0, '', 3, 0, 1, 1),
('新闻动态', '/news', 45, 1, '', 3, 0, 1, 1),
('解决方案', '/solutions', 45, 2, '', 3, 0, 1, 1);

-- ============================================
-- 默认轮播图 (4张)
-- ============================================
INSERT INTO `banners` (`title`, `description`, `type`, `media_url`, `link_url`, `button_text`, `sort_order`, `status`, `position`) VALUES
('高性能云服务器', '企业级云服务器，99.9% SLA 保障', 'video', '/carousel/1.webm', '/products', '立即选购', 0, 1, 'home'),
('全球节点覆盖', '遍布全球 30+ 数据中心节点', 'video', '/carousel/2.webm', '/products', '立即选购', 1, 1, 'home'),
('专业技术支持', '7×24小时专业技术团队', 'video', '/carousel/3.webm', '/contact', '联系我们', 2, 1, 'home'),
('安全可靠保障', '多层安全防护，数据安全无忧', 'video', '/carousel/4.webm', '/products', '了解更多', 3, 1, 'home');
