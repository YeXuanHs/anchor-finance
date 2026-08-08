-- ============================================================
-- 导入 zjmf 风格菜单结构 (7 个顶级菜单，4 级层级)
-- 基于 auth_rule.sql 分析
-- ============================================================

-- 清空现有菜单
DELETE FROM menus;

-- 重置自增 ID
ALTER TABLE menus AUTO_INCREMENT = 1;

-- ============================================================
-- 1. 客户 (Customer) - 顶级菜单
-- ============================================================
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('客户', 'ep:User', '/customer-list', NULL, 1, 'admin', true, true, '{"CN":"客户","HK":"客戶","US":"Customer"}', NOW(), NOW());

SET @customer_id = LAST_INSERT_ID();

-- 1.1 客户管理
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('客户管理', '', '', @customer_id, 1, 'admin', true, true, '{"CN":"客户管理","HK":"客戶管理","US":"Customer Management"}', NOW(), NOW());

SET @customer_mgmt_id = LAST_INSERT_ID();

INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('客户列表', '', '/customer-list', @customer_mgmt_id, 1, 'admin', true, true, '{"CN":"客户列表","HK":"客戶列表","US":"Customer List"}', NOW(), NOW()),
('实名认证', '', '/customer-authentication', @customer_mgmt_id, 2, 'admin', true, true, '{"CN":"实名认证","HK":"實名認證","US":"Real-name Authentication"}', NOW(), NOW()),
('客户资源池', '', '/customer-resources', @customer_mgmt_id, 3, 'admin', true, true, '{"CN":"客户资源池","HK":"客戶資源池","US":"Customer Resource Pool"}', NOW(), NOW());

-- 1.2 我的业绩
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('我的业绩', '', '/sales-statistics', @customer_id, 2, 'admin', true, true, '{"CN":"我的业绩","HK":"我的業績","US":"My Performance"}', NOW(), NOW());

-- 1.3 运营管理
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('运营管理', '', '', @customer_id, 3, 'admin', true, true, '{"CN":"运营管理","HK":"運營管理","US":"Operation Management"}', NOW(), NOW());

SET @operation_mgmt_id = LAST_INSERT_ID();

INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('推介计划', '', '/customer-promotionplan', @operation_mgmt_id, 1, 'admin', true, true, '{"CN":"推介计划","HK":"推介計劃","US":"Recommendation Plan"}', NOW(), NOW()),
('营销推送', '', '/marketing-push', @operation_mgmt_id, 2, 'admin', true, true, '{"CN":"营销推送","HK":"營銷推送","US":"Marketing Push"}', NOW(), NOW());

-- ============================================================
-- 2. 业务 (Business) - 顶级菜单
-- ============================================================
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('业务', 'ep:ShoppingCart', '/order-list', NULL, 2, 'admin', true, true, '{"CN":"业务","HK":"業務","US":"Business"}', NOW(), NOW());

SET @business_id = LAST_INSERT_ID();

-- 2.1 订单
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('订单', '', '', @business_id, 1, 'admin', true, true, '{"CN":"订单","HK":"訂單","US":"Order"}', NOW(), NOW());

SET @order_id = LAST_INSERT_ID();

INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('产品订单', '', '/order-list', @order_id, 1, 'admin', true, true, '{"CN":"产品订单","HK":"產品訂單","US":"Product Order"}', NOW(), NOW()),
('续费订单', '', '/renewal-order', @order_id, 2, 'admin', true, true, '{"CN":"续费订单","HK":"續費訂單","US":"Renewal Order"}', NOW(), NOW()),
('流量包订单', '', '/dcim-traffic-log', @order_id, 3, 'admin', true, true, '{"CN":"流量包订单","HK":"流量包訂單","US":"Traffic Package Order"}', NOW(), NOW());

-- 2.2 业务
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('业务', '', '', @business_id, 2, 'admin', true, true, '{"CN":"业务","HK":"業務","US":"Business"}', NOW(), NOW());

SET @business_list_id = LAST_INSERT_ID();

INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('业务列表', '', '/customer-product', @business_list_id, 1, 'admin', true, true, '{"CN":"业务列表","HK":"業務列表","US":"Business List"}', NOW(), NOW()),
('产品暂停请求', '', '/customer-cancelreq', @business_list_id, 2, 'admin', true, true, '{"CN":"产品暂停请求","HK":"產品暫停請求","US":"Product Pause Request"}', NOW(), NOW());

-- ============================================================
-- 3. 财务 (Finance) - 顶级菜单
-- ============================================================
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('财务', 'ep:Money', '/business-statement', NULL, 3, 'admin', true, true, '{"CN":"财务","HK":"財務","US":"Finance"}', NOW(), NOW());

SET @finance_id = LAST_INSERT_ID();

-- 3.1 财务记录
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('财务记录', '', '', @finance_id, 1, 'admin', true, true, '{"CN":"财务记录","HK":"財務記錄","US":"Financial Records"}', NOW(), NOW());

SET @finance_records_id = LAST_INSERT_ID();

INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('交易流水', '', '/business-statement', @finance_records_id, 1, 'admin', true, true, '{"CN":"交易流水","HK":"交易流水","US":"Trading Flow"}', NOW(), NOW()),
('账单管理', '', '/bill-management', @finance_records_id, 2, 'admin', true, true, '{"CN":"账单管理","HK":"賬單管理","US":"Bill Management"}', NOW(), NOW()),
('信用额管理', '', '/credit-management', @finance_records_id, 3, 'admin', true, true, '{"CN":"信用额管理","HK":"信用額管理","US":"Credit Management"}', NOW(), NOW());

-- 3.2 审核管理
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('审核管理', '', '', @finance_id, 2, 'admin', true, true, '{"CN":"审核管理","HK":"審核管理","US":"Audit Management"}', NOW(), NOW());

SET @audit_mgmt_id = LAST_INSERT_ID();

INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('提现审核', '', '/customer-withdrawal', @audit_mgmt_id, 1, 'admin', true, true, '{"CN":"提现审核","HK":"提現審核","US":"Withdrawal Review"}', NOW(), NOW());

-- 3.3 发票和合同
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('发票和合同', '', '', @finance_id, 3, 'admin', true, true, '{"CN":"发票和合同","HK":"發票和合同","US":"Invoices and Contracts"}', NOW(), NOW());

SET @invoice_contract_id = LAST_INSERT_ID();

INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('发票列表', '', '/invoice-audit', @invoice_contract_id, 1, 'admin', true, true, '{"CN":"发票列表","HK":"發票列表","US":"Invoice List"}', NOW(), NOW()),
('合同列表', '', '/contracts_audit', @invoice_contract_id, 2, 'admin', true, true, '{"CN":"合同列表","HK":"合同列表","US":"Contracts List"}', NOW(), NOW());

-- ============================================================
-- 4. 工单 (Ticket) - 顶级菜单
-- ============================================================
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('工单', 'ep:Tickets', '/support-ticket', NULL, 4, 'admin', true, true, '{"CN":"工单","HK":"工單","US":"Ticket"}', NOW(), NOW());

SET @ticket_id = LAST_INSERT_ID();

-- 4.1 工单
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('工单', '', '', @ticket_id, 1, 'admin', true, true, '{"CN":"工单","HK":"工單","US":"Ticket"}', NOW(), NOW());

SET @ticket_list_id = LAST_INSERT_ID();

INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('工单列表', '', '/support-ticket', @ticket_list_id, 1, 'admin', true, true, '{"CN":"工单列表","HK":"工單列表","US":"Ticket List"}', NOW(), NOW()),
('工单统计', '', '/support-statistics', @ticket_list_id, 2, 'admin', true, true, '{"CN":"工单统计","HK":"工單統計","US":"Ticket Statistics"}', NOW(), NOW());

-- ============================================================
-- 5. 功能 (Function) - 顶级菜单
-- ============================================================
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('功能', 'ep:DataAnalysis', '/timing-results', NULL, 5, 'admin', true, true, '{"CN":"功能","HK":"功能","US":"Function"}', NOW(), NOW());

SET @function_id = LAST_INSERT_ID();

-- 5.1 插件
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('插件', '', '', @function_id, 1, 'admin', true, true, '{"CN":"插件","HK":"插件","US":"Plugin"}', NOW(), NOW());

SET @plugin_id = LAST_INSERT_ID();

INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('插件列表', '', '/plugins', @plugin_id, 1, 'admin', true, true, '{"CN":"插件列表","HK":"插件列表","US":"Plugin List"}', NOW(), NOW());

-- 5.2 系统状态
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('系统状态', '', '', @function_id, 2, 'admin', true, true, '{"CN":"系统状态","HK":"系統狀態","US":"System Status"}', NOW(), NOW());

SET @system_status_id = LAST_INSERT_ID();

INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('数据库状态', '', '/database-message', @system_status_id, 1, 'admin', true, true, '{"CN":"数据库状态","HK":"數據庫狀態","US":"Database Status"}', NOW(), NOW()),
('任务队列', '', '/statistics-taskQueue', @system_status_id, 2, 'admin', true, true, '{"CN":"任务队列","HK":"任務隊列","US":"Task Queue"}', NOW(), NOW()),
('定时任务状态', '', '/timing-results', @system_status_id, 3, 'admin', true, true, '{"CN":"定时任务状态","HK":"定時任務狀態","US":"Timed Task Status"}', NOW(), NOW());

-- 5.3 统计
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('统计', '', '', @function_id, 3, 'admin', true, true, '{"CN":"统计","HK":"統計","US":"Statistics"}', NOW(), NOW());

SET @statistics_id = LAST_INSERT_ID();

INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('年度收入统计', '', '/annual-statistics', @statistics_id, 1, 'admin', true, true, '{"CN":"年度收入统计","HK":"年度收入統計","US":"Annual Income Statistics"}', NOW(), NOW()),
('新客户', '', '/new-customer', @statistics_id, 2, 'admin', true, true, '{"CN":"新客户","HK":"新客戶","US":"New Customer"}', NOW(), NOW()),
('产品收入', '', '/product-revenue', @statistics_id, 3, 'admin', true, true, '{"CN":"产品收入","HK":"產品收入","US":"Product Revenue"}', NOW(), NOW()),
('收入排名', '', '/revenue-ranking', @statistics_id, 4, 'admin', true, true, '{"CN":"收入排名","HK":"收入排名","US":"Revenue Ranking"}', NOW(), NOW());

-- ============================================================
-- 6. 资源与商店 (Resources And Stores) - 顶级菜单
-- ============================================================
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('资源与商店', 'ep:Shop', '/app-store', NULL, 6, 'admin', true, true, '{"CN":"资源与商店","HK":"資源與商店","US":"Resources And Stores"}', NOW(), NOW());

SET @store_id = LAST_INSERT_ID();

-- 6.1 应用商店
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('应用商店', '', '', @store_id, 1, 'admin', true, true, '{"CN":"应用商店","HK":"應用商店","US":"App Store"}', NOW(), NOW());

SET @app_store_id = LAST_INSERT_ID();

INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('我的应用', '', '/app-store', @app_store_id, 1, 'admin', true, true, '{"CN":"我的应用","HK":"我的應用","US":"My Application"}', NOW(), NOW());

-- 6.2 上下游
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('上下游', '', '', @store_id, 2, 'admin', true, true, '{"CN":"上下游","HK":"上下游","US":"Upstream And Downstream"}', NOW(), NOW());

SET @upstream_downstream_id = LAST_INSERT_ID();

-- 6.2.1 下游管理
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('下游管理', '', '', @upstream_downstream_id, 1, 'admin', true, true, '{"CN":"下游管理","HK":"下游管理","US":"Downstream Management"}', NOW(), NOW());

SET @downstream_id = LAST_INSERT_ID();

INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('API设置', '', '/api-setup', @downstream_id, 1, 'admin', true, true, '{"CN":"API设置","HK":"API設置","US":"API Settings"}', NOW(), NOW()),
('任务队列', '', '/task-queue', @downstream_id, 2, 'admin', true, true, '{"CN":"任务队列","HK":"任務隊列","US":"Task Queue"}', NOW(), NOW());

-- 6.2.2 上游资源
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('上游资源', '', '', @upstream_downstream_id, 2, 'admin', true, true, '{"CN":"上游资源","HK":"上游資源","US":"Upstream Resources"}', NOW(), NOW());

SET @upstream_id = LAST_INSERT_ID();

INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('服务器列表', '', '/munual-resource', @upstream_id, 1, 'admin', true, true, '{"CN":"服务器列表","HK":"服務器列表","US":"Server List"}', NOW(), NOW()),
('供应商管理', '', '/zjmf-api', @upstream_id, 2, 'admin', true, true, '{"CN":"供应商管理","HK":"供應商管理","US":"Supplier Management"}', NOW(), NOW()),
('商品管理', '', '/commodity-list', @upstream_id, 3, 'admin', true, true, '{"CN":"商品管理","HK":"商品管理","US":"Commodity Management"}', NOW(), NOW()),
('产品管理', '', '/commodity-product', @upstream_id, 4, 'admin', true, true, '{"CN":"产品管理","HK":"產品管理","US":"Product Management"}', NOW(), NOW()),
('任务队列', '', '/commodity-taskQueue', @upstream_id, 5, 'admin', true, true, '{"CN":"任务队列","HK":"任務隊列","US":"Task Queue"}', NOW(), NOW()),
('订单列表', '', '/supplier-order-list', @upstream_id, 6, 'admin', true, true, '{"CN":"订单列表","HK":"訂單列表","US":"Order List"}', NOW(), NOW()),
('续费订单', '', '/supplier-renewal-order', @upstream_id, 7, 'admin', true, true, '{"CN":"续费订单","HK":"續費訂單","US":"Renewal Order"}', NOW(), NOW());

-- ============================================================
-- 7. 设置 (Settings) - 顶级菜单 (4 级层级)
-- ============================================================
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('设置', 'ep:Setting', '/set', NULL, 7, 'admin', true, true, '{"CN":"设置","HK":"設置","US":"Settings"}', NOW(), NOW());

SET @settings_id = LAST_INSERT_ID();

-- 7.1 商品设置
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('商品设置', '', '', @settings_id, 1, 'admin', true, true, '{"CN":"商品设置","HK":"商品設置","US":"Commodity Settings"}', NOW(), NOW());

SET @product_settings_id = LAST_INSERT_ID();

-- 7.1.1 商品配置
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('商品配置', '', '', @product_settings_id, 1, 'admin', true, true, '{"CN":"商品配置","HK":"商品配置","US":"Commodity Configuration"}', NOW(), NOW());

SET @product_config_id = LAST_INSERT_ID();

INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('商品管理', '', '/product-server', @product_config_id, 1, 'admin', true, true, '{"CN":"商品管理","HK":"商品管理","US":"Commodity Management"}', NOW(), NOW()),
('流量包管理', '', '/dcim-traffic', @product_config_id, 2, 'admin', true, true, '{"CN":"流量包管理","HK":"流量包管理","US":"Traffic Package Management"}', NOW(), NOW());

-- 7.1.2 自动化接口
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('自动化接口', '', '', @product_settings_id, 2, 'admin', true, true, '{"CN":"自动化接口","HK":"自動化接口","US":"Automation Interface"}', NOW(), NOW());

SET @auto_interface_id = LAST_INSERT_ID();

INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('通用接口', '', '/server-settings', @auto_interface_id, 1, 'admin', true, true, '{"CN":"通用接口","HK":"通用接口","US":"Universal Interface"}', NOW(), NOW()),
('魔方DCIM', '', '/dcim', @auto_interface_id, 2, 'admin', true, true, '{"CN":"魔方DCIM","HK":"魔方DCIM","US":"Cube DCIM"}', NOW(), NOW()),
('魔方云', '', '/zjmfcloud', @auto_interface_id, 3, 'admin', true, true, '{"CN":"魔方云","HK":"魔方雲","US":"Magic Cube Cloud"}', NOW(), NOW());

-- 7.1.3 设置
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('设置', '', '', @product_settings_id, 3, 'admin', true, true, '{"CN":"设置","HK":"設置","US":"Settings"}', NOW(), NOW());

SET @product_setting_id = LAST_INSERT_ID();

INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('全局可配置项', '', '/configurable-option', @product_setting_id, 1, 'admin', true, true, '{"CN":"全局可配置项","HK":"全局可配置項","US":"Globally Configurable Items"}', NOW(), NOW()),
('商品订购设置', '', '/order-product', @product_setting_id, 2, 'admin', true, true, '{"CN":"商品订购设置","HK":"商品訂購設置","US":"Product Order Setting"}', NOW(), NOW());

-- 7.2 基础设置
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('基础设置', '', '', @settings_id, 2, 'admin', true, true, '{"CN":"基础设置","HK":"基礎設置","US":"Basic Settings"}', NOW(), NOW());

SET @basic_settings_id = LAST_INSERT_ID();

-- 7.2.1 工单设置
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('工单设置', '', '', @basic_settings_id, 1, 'admin', true, true, '{"CN":"工单设置","HK":"工單設置","US":"Ticket Settings"}', NOW(), NOW());

SET @ticket_settings_id = LAST_INSERT_ID();

INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('工单部门', '', '/work-order-dept', @ticket_settings_id, 1, 'admin', true, true, '{"CN":"工单部门","HK":"工單部門","US":"Ticket Department"}', NOW(), NOW()),
('工单状态', '', '/work-order-status', @ticket_settings_id, 2, 'admin', true, true, '{"CN":"工单状态","HK":"工單狀態","US":"Ticket Status"}', NOW(), NOW()),
('工单传递', '', '/work-order-rules', @ticket_settings_id, 3, 'admin', true, true, '{"CN":"工单传递","HK":"工單傳遞","US":"Ticket Transfer"}', NOW(), NOW());

-- 7.2.2 客户设置
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('客户设置', '', '', @basic_settings_id, 2, 'admin', true, true, '{"CN":"客户设置","HK":"客戶設置","US":"Customer Settings"}', NOW(), NOW());

SET @customer_settings_id = LAST_INSERT_ID();

INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('客户分组与折扣', '', '/customer-group', @customer_settings_id, 1, 'admin', true, true, '{"CN":"客户分组与折扣","HK":"客戶分組與折扣","US":"Customer Grouping and Discount"}', NOW(), NOW()),
('实名设置', '', '/authentication-setting', @customer_settings_id, 2, 'admin', true, true, '{"CN":"实名设置","HK":"實名設置","US":"Identity Verification"}', NOW(), NOW()),
('自定义客户字段', '', '/customer-custom', @customer_settings_id, 3, 'admin', true, true, '{"CN":"自定义客户字段","HK":"自定義客戶字段","US":"Custom Customer Field"}', NOW(), NOW()),
('推介设置', '', '/promotion_plan', @customer_settings_id, 4, 'admin', true, true, '{"CN":"推介设置","HK":"推介設置","US":"Recommendation Settings"}', NOW(), NOW()),
('客户等级', '', '/customer-level', @customer_settings_id, 5, 'admin', true, true, '{"CN":"客户等级","HK":"客戶等級","US":"Customer Level"}', NOW(), NOW());

-- 7.2.3 财务设置
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('财务设置', '', '', @basic_settings_id, 3, 'admin', true, true, '{"CN":"财务设置","HK":"財務設置","US":"Financial Settings"}', NOW(), NOW());

SET @finance_settings_id = LAST_INSERT_ID();

INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('支付接口', '', '/payment-interface', @finance_settings_id, 1, 'admin', true, true, '{"CN":"支付接口","HK":"支付接口","US":"Payment Interface"}', NOW(), NOW()),
('优惠码', '', '/promo-code', @finance_settings_id, 2, 'admin', true, true, '{"CN":"优惠码","HK":"優惠碼","US":"Promotion Code"}', NOW(), NOW()),
('货币配置', '', '/currency-settings', @finance_settings_id, 3, 'admin', true, true, '{"CN":"货币配置","HK":"貨幣配置","US":"Currency Configuration"}', NOW(), NOW()),
('充值与财务', '', '/general-settings/finance', @finance_settings_id, 4, 'admin', true, true, '{"CN":"充值与财务","HK":"充值與財務","US":"Recharge and Finance"}', NOW(), NOW()),
('发票设置', '', '/voucher-setting', @finance_settings_id, 5, 'admin', true, true, '{"CN":"发票设置","HK":"發票設置","US":"Invoice Settings"}', NOW(), NOW()),
('信用额设置', '', '/credit-setting', @finance_settings_id, 6, 'admin', true, true, '{"CN":"信用额设置","HK":"信用額設定","US":"Credit Limit Setting"}', NOW(), NOW()),
('合同设置', '', '/contracts_setting', @finance_settings_id, 7, 'admin', true, true, '{"CN":"合同设置","HK":"合同設置","US":"Contracts Settings"}', NOW(), NOW());

-- 7.3 站务设置
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('站务设置', '', '', @settings_id, 3, 'admin', true, true, '{"CN":"站务设置","HK":"站務設置","US":"Station Service Settings"}', NOW(), NOW());

SET @station_settings_id = LAST_INSERT_ID();

INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('显示设置', '', '/base-info', @station_settings_id, 1, 'admin', true, true, '{"CN":"显示设置","HK":"顯示設置","US":"Display Settings"}', NOW(), NOW()),
('文件下载', '', '/service-support', @station_settings_id, 2, 'admin', true, true, '{"CN":"文件下载","HK":"文件下載","US":"File Download"}', NOW(), NOW()),
('新闻中心', '', '/news-center', @station_settings_id, 3, 'admin', true, true, '{"CN":"新闻中心","HK":"新聞中心","US":"News Center"}', NOW(), NOW()),
('帮助中心', '', '/help-center', @station_settings_id, 4, 'admin', true, true, '{"CN":"帮助中心","HK":"幫助中心","US":"Help Center"}', NOW(), NOW());

-- 7.4 系统设置
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('系统设置', '', '', @settings_id, 4, 'admin', true, true, '{"CN":"系统设置","HK":"系統設置","US":"System Settings"}', NOW(), NOW());

SET @system_settings_id = LAST_INSERT_ID();

-- 7.4.1 基础设置
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('基础设置', '', '', @system_settings_id, 1, 'admin', true, true, '{"CN":"基础设置","HK":"基礎設置","US":"Basic Settings"}', NOW(), NOW());

SET @sys_basic_id = LAST_INSERT_ID();

INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('常规设置', '', '/general-settings/general', @sys_basic_id, 1, 'admin', true, true, '{"CN":"常规设置","HK":"常規設置","US":"General Settings"}', NOW(), NOW()),
('定时任务', '', '/automatic-tasks', @sys_basic_id, 2, 'admin', true, true, '{"CN":"定时任务","HK":"定時任務","US":"Timed Task"}', NOW(), NOW()),
('注册登录', '', '/login-register', @sys_basic_id, 3, 'admin', true, true, '{"CN":"注册登录","HK":"註冊登錄","US":"Register and Login"}', NOW(), NOW()),
('第三方登录', '', '/third-login', @sys_basic_id, 4, 'admin', true, true, '{"CN":"第三方登录","HK":"第三方登錄","US":"Third Party Login"}', NOW(), NOW());

-- 7.4.2 人员管理
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('人员管理', '', '', @system_settings_id, 2, 'admin', true, true, '{"CN":"人员管理","HK":"人員管理","US":"Personnel Management"}', NOW(), NOW());

SET @personnel_id = LAST_INSERT_ID();

INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('员工管理', '', '/admin-management', @personnel_id, 1, 'admin', true, true, '{"CN":"员工管理","HK":"員工管理","US":"Staff Management"}', NOW(), NOW()),
('分组权限', '', '/permissions-managment', @personnel_id, 2, 'admin', true, true, '{"CN":"分组权限","HK":"分組權限","US":"Group Permission"}', NOW(), NOW()),
('销售设置', '', '/sales-management', @personnel_id, 3, 'admin', true, true, '{"CN":"销售设置","HK":"銷售設置","US":"Sales Settings"}', NOW(), NOW());

-- 7.4.3 短信邮件设置
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('短信邮件设置', '', '', @system_settings_id, 3, 'admin', true, true, '{"CN":"短信邮件设置","HK":"短信郵件設置","US":"SMS Mail Settings"}', NOW(), NOW());

SET @sms_email_id = LAST_INSERT_ID();

INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('接口设置', '', '/sms-template/sms', @sms_email_id, 1, 'admin', true, true, '{"CN":"接口设置","HK":"接口設置","US":"Interface Settings"}', NOW(), NOW()),
('邮件模板', '', '/email-list', @sms_email_id, 2, 'admin', true, true, '{"CN":"邮件模板","HK":"郵件模板","US":"Mail Template"}', NOW(), NOW()),
('短信模板', '', '/sms-template-index', @sms_email_id, 3, 'admin', true, true, '{"CN":"短信模板","HK":"短信模板","US":"SMS Template"}', NOW(), NOW()),
('发送设置', '', '/sms-send-settings', @sms_email_id, 4, 'admin', true, true, '{"CN":"发送设置","HK":"發送設置","US":"Send Settings"}', NOW(), NOW());

-- 7.4.4 安全相关
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('安全相关', '', '', @system_settings_id, 4, 'admin', true, true, '{"CN":"安全相关","HK":"安全相關","US":"Security Related"}', NOW(), NOW());

SET @security_id = LAST_INSERT_ID();

INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('黑名单列表', '', '/black-list', @security_id, 1, 'admin', true, true, '{"CN":"黑名单列表","HK":"黑名單列表","US":"Blacklist List"}', NOW(), NOW()),
('验证码设置', '', '/general-settings/captcha', @security_id, 2, 'admin', true, true, '{"CN":"验证码设置","HK":"驗證碼設置","US":"Verification Code Settings"}', NOW(), NOW()),
('二次验证', '', '/twice-confirm', @security_id, 3, 'admin', true, true, '{"CN":"二次验证","HK":"二次驗證","US":"Secondary Verification"}', NOW(), NOW());

-- 7.4.5 系统相关
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('系统相关', '', '', @system_settings_id, 5, 'admin', true, true, '{"CN":"系统相关","HK":"系統相關","US":"System Related"}', NOW(), NOW());

SET @sys_related_id = LAST_INSERT_ID();

INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('系统升级', '', '/system-message', @sys_related_id, 1, 'admin', true, true, '{"CN":"系统升级","HK":"系統升級","US":"System Upgrade"}', NOW(), NOW()),
('关于', '', '/about', @sys_related_id, 2, 'admin', true, true, '{"CN":"关于","HK":"關於","US":"About"}', NOW(), NOW());

-- 7.4.6 日志
INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('日志', '', '', @system_settings_id, 6, 'admin', true, true, '{"CN":"日志","HK":"日誌","US":"Log"}', NOW(), NOW());

SET @log_id = LAST_INSERT_ID();

INSERT INTO menus (name, icon, url, parent_id, sort_order, type, is_visible, is_active, language_map, created_at, updated_at) VALUES
('系统日志', '', '/system-log', @log_id, 1, 'admin', true, true, '{"CN":"系统日志","HK":"系統日誌","US":"System Log"}', NOW(), NOW()),
('管理员登录日志', '', '/system-admin-log', @log_id, 2, 'admin', true, true, '{"CN":"管理员登录日志","HK":"管理員登錄日誌","US":"Administrator Login Log"}', NOW(), NOW()),
('邮件日志', '', '/email-log', @log_id, 3, 'admin', true, true, '{"CN":"邮件日志","HK":"郵件日誌","US":"Mail Log"}', NOW(), NOW()),
('短信日志', '', '/sms-log', @log_id, 4, 'admin', true, true, '{"CN":"短信日志","HK":"短信日誌","US":"SMS Log"}', NOW(), NOW()),
('站内信日志', '', '/station-letter-log', @log_id, 5, 'admin', true, true, '{"CN":"站内信日志","HK":"站內信日誌","US":"Site Message Log"}', NOW(), NOW()),
('定时任务日志', '', '/automatic-task-log', @log_id, 6, 'admin', true, true, '{"CN":"定时任务日志","HK":"定時任務日誌","US":"Timed Task Log"}', NOW(), NOW()),
('API日志', '', '/api-log', @log_id, 7, 'admin', true, true, '{"CN":"API日志","HK":"API日誌","US":"API Log"}', NOW(), NOW()),
('日志清理', '', '/log-cleanup', @log_id, 8, 'admin', true, true, '{"CN":"日志清理","HK":"日誌清理","US":"Log Cleanup"}', NOW(), NOW());
