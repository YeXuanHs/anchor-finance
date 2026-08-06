-- ============================================
-- 锚点财务 (AnchorFinance) 初始化脚本
-- 由 install.sh 自动导入
-- 管理员账号、站点名称、后台路径等由 install.sh 动态写入
-- ============================================

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ============================================
-- 用户表
-- ============================================
CREATE TABLE IF NOT EXISTS `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) DEFAULT NULL,
  `username` varchar(50) NOT NULL,
  `email` varchar(100) NOT NULL DEFAULT '',
  `phone` varchar(20) DEFAULT NULL,
  `password` varchar(255) NOT NULL,
  `balance` decimal(10,2) NOT NULL DEFAULT 0.00,
  `credit_limit` decimal(10,2) NOT NULL DEFAULT 0.00,
  `status` tinyint NOT NULL DEFAULT 1,
  `is_admin` tinyint(1) NOT NULL DEFAULT 0,
  `avatar` varchar(255) DEFAULT NULL,
  `company` varchar(100) DEFAULT NULL,
  `address` varchar(255) DEFAULT NULL,
  `city` varchar(50) DEFAULT NULL,
  `state` varchar(50) DEFAULT NULL,
  `postcode` varchar(20) DEFAULT NULL,
  `country` varchar(10) DEFAULT 'CN',
  `language` varchar(10) DEFAULT 'zh-CN',
  `currency` varchar(10) DEFAULT 'CNY',
  `group_id` bigint unsigned DEFAULT NULL,
  `last_login_ip` varchar(45) DEFAULT NULL,
  `last_login_at` datetime DEFAULT NULL,
  `email_verified` tinyint(1) NOT NULL DEFAULT 0,
  `phone_verified` tinyint(1) NOT NULL DEFAULT 0,
  `real_name` varchar(50) DEFAULT NULL,
  `id_card` varchar(30) DEFAULT NULL,
  `real_name_status` tinyint NOT NULL DEFAULT 0 COMMENT '0=未认证 1=已认证 2=待审核',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  UNIQUE KEY `uk_uuid` (`uuid`),
  KEY `idx_email` (`email`),
  KEY `idx_phone` (`phone`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- 产品分组
-- ============================================
CREATE TABLE IF NOT EXISTS `product_groups` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `description` text,
  `slug` varchar(100) DEFAULT NULL,
  `parent_id` bigint unsigned DEFAULT 0,
  `hidden` tinyint(1) NOT NULL DEFAULT 0,
  `show_in_nav` tinyint(1) NOT NULL DEFAULT 1,
  `sort_order` int NOT NULL DEFAULT 0,
  `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- 产品
-- ============================================
CREATE TABLE IF NOT EXISTS `products` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `group_id` bigint unsigned DEFAULT NULL,
  `name` varchar(200) NOT NULL,
  `slug` varchar(200) DEFAULT NULL,
  `description` text,
  `type` varchar(50) DEFAULT 'hosting' COMMENT 'hosting|server|cloud|other',
  `price` decimal(10,2) NOT NULL DEFAULT 0.00,
  `currency` varchar(10) DEFAULT 'CNY',
  `billing_cycle` varchar(20) DEFAULT 'monthly' COMMENT 'monthly|quarterly|semi-annually|annually|triennially|free|onetime',
  `pay_type` tinyint NOT NULL DEFAULT 1 COMMENT '1=预付款 2=后付款',
  `auto_setup` tinyint(1) NOT NULL DEFAULT 0,
  `server_group` varchar(50) DEFAULT NULL,
  `stock_control` tinyint(1) NOT NULL DEFAULT 0,
  `stock_qty` int DEFAULT NULL,
  `hidden` tinyint(1) NOT NULL DEFAULT 0,
  `sort_order` int NOT NULL DEFAULT 0,
  `upstream_id` bigint unsigned DEFAULT NULL COMMENT '关联上游供应商',
  `remote_product_id` varchar(64) DEFAULT NULL COMMENT '上游产品ID',
  `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_group_id` (`group_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- 订单
-- ============================================
CREATE TABLE IF NOT EXISTS `orders` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `order_no` varchar(32) NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `product_id` bigint unsigned DEFAULT NULL,
  `amount` decimal(10,2) NOT NULL DEFAULT 0.00,
  `total` decimal(10,2) NOT NULL DEFAULT 0.00,
  `currency` varchar(10) DEFAULT 'CNY',
  `billing_cycle` varchar(20) DEFAULT NULL,
  `payment_method` varchar(50) DEFAULT NULL,
  `status` tinyint NOT NULL DEFAULT 0 COMMENT '0=待付款 1=已付款 2=已取消 3=已完成 4=退款中 5=已退款',
  `promo_code` varchar(50) DEFAULT NULL,
  `notes` text,
  `paid_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_order_no` (`order_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- 订单明细
-- ============================================
CREATE TABLE IF NOT EXISTS `order_items` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `order_id` bigint unsigned NOT NULL,
  `product_id` bigint unsigned DEFAULT NULL,
  `name` varchar(200) NOT NULL,
  `quantity` int NOT NULL DEFAULT 1,
  `unit_price` decimal(10,2) NOT NULL DEFAULT 0.00,
  `amount` decimal(10,2) NOT NULL DEFAULT 0.00,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_order_id` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- 账单
-- ============================================
CREATE TABLE IF NOT EXISTS `invoices` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `invoice_no` varchar(32) NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `order_id` bigint unsigned DEFAULT NULL,
  `sub_total` decimal(10,2) NOT NULL DEFAULT 0.00,
  `tax` decimal(10,2) NOT NULL DEFAULT 0.00,
  `discount` decimal(10,2) NOT NULL DEFAULT 0.00,
  `credit` decimal(10,2) NOT NULL DEFAULT 0.00,
  `total` decimal(10,2) NOT NULL DEFAULT 0.00,
  `paid_amount` decimal(10,2) NOT NULL DEFAULT 0.00,
  `payment_method` varchar(50) DEFAULT NULL,
  `status` tinyint NOT NULL DEFAULT 0 COMMENT '0=未支付 1=已支付 2=已取消 3=已退款',
  `due_date` datetime DEFAULT NULL,
  `paid_at` datetime DEFAULT NULL,
  `notes` text,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_invoice_no` (`invoice_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- 账单明细
-- ============================================
CREATE TABLE IF NOT EXISTS `invoice_items` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `invoice_id` bigint unsigned NOT NULL,
  `type` varchar(20) DEFAULT 'product' COMMENT 'product|addon|domain|other',
  `rel_id` bigint unsigned DEFAULT NULL,
  `description` varchar(500) DEFAULT NULL,
  `quantity` int NOT NULL DEFAULT 1,
  `unit_price` decimal(10,2) NOT NULL DEFAULT 0.00,
  `total` decimal(10,2) NOT NULL DEFAULT 0.00,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_invoice_id` (`invoice_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- 主机/服务实例
-- ============================================
CREATE TABLE IF NOT EXISTS `hosts` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `order_id` bigint unsigned DEFAULT NULL,
  `product_id` bigint unsigned DEFAULT NULL,
  `hostname` varchar(200) DEFAULT NULL,
  `domain` varchar(200) DEFAULT NULL,
  `ip` varchar(45) DEFAULT NULL,
  `dedicated_ip` varchar(45) DEFAULT NULL,
  `os` varchar(100) DEFAULT NULL,
  `cpu` varchar(50) DEFAULT NULL,
  `memory` varchar(50) DEFAULT NULL,
  `disk` varchar(50) DEFAULT NULL,
  `bandwidth` varchar(50) DEFAULT NULL,
  `username` varchar(100) DEFAULT NULL,
  `password` varchar(255) DEFAULT NULL,
  `port` int DEFAULT NULL,
  `billing_cycle` varchar(20) DEFAULT NULL,
  `amount` decimal(10,2) NOT NULL DEFAULT 0.00,
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '0=暂停 1=活跃 2=待开通 3=已删除 4=已过期',
  `next_due_date` datetime DEFAULT NULL,
  `suspended_at` datetime DEFAULT NULL,
  `terminated_at` datetime DEFAULT NULL,
  `notes` text,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- 工单部门
-- ============================================
CREATE TABLE IF NOT EXISTS `departments` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `description` text,
  `slug` varchar(100) DEFAULT NULL,
  `email` varchar(100) DEFAULT NULL,
  `hidden` tinyint(1) NOT NULL DEFAULT 0,
  `sort_order` int NOT NULL DEFAULT 0,
  `status` tinyint NOT NULL DEFAULT 1,
  `auto_close_hours` int DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- 工单
-- ============================================
CREATE TABLE IF NOT EXISTS `tickets` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `ticket_no` varchar(32) NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `department_id` bigint unsigned DEFAULT NULL,
  `subject` varchar(200) NOT NULL,
  `content` text,
  `priority` tinyint NOT NULL DEFAULT 0 COMMENT '0=低 1=中 2=高 3=紧急',
  `status` tinyint NOT NULL DEFAULT 0 COMMENT '0=待回复 1=已回复 2=已关闭 3=待处理',
  `admin_id` bigint unsigned DEFAULT NULL,
  `host_id` bigint unsigned DEFAULT NULL,
  `evaluate_score` tinyint DEFAULT NULL,
  `evaluate_content` text,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_ticket_no` (`ticket_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_department_id` (`department_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- 工单回复
-- ============================================
CREATE TABLE IF NOT EXISTS `ticket_replies` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `ticket_id` bigint unsigned NOT NULL,
  `user_id` bigint unsigned DEFAULT NULL,
  `admin_id` bigint unsigned DEFAULT NULL,
  `content` text NOT NULL,
  `is_internal` tinyint(1) NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_ticket_id` (`ticket_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- 附件
-- ============================================
CREATE TABLE IF NOT EXISTS `attachments` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `ticket_id` bigint unsigned DEFAULT NULL,
  `reply_id` bigint unsigned DEFAULT NULL,
  `uploader_id` bigint unsigned DEFAULT NULL,
  `file_name` varchar(255) NOT NULL,
  `file_path` varchar(500) NOT NULL,
  `file_size` bigint DEFAULT 0,
  `mime_type` varchar(100) DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_ticket_id` (`ticket_id`),
  KEY `idx_reply_id` (`reply_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- 系统配置（统一配置表）
-- ============================================
CREATE TABLE IF NOT EXISTS `system_configs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `group` varchar(50) NOT NULL,
  `key` varchar(100) NOT NULL,
  `value` text,
  `type` varchar(20) DEFAULT 'string' COMMENT 'string|int|bool|json|text',
  `name` varchar(100) DEFAULT NULL,
  `options` text COMMENT '选项(JSON)',
  `sort_order` int NOT NULL DEFAULT 0,
  `status` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_key` (`key`),
  KEY `idx_group` (`group`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- 系统配置默认值（safe defaults，不包含敏感信息）
-- 管理员账号、站点名称、后台路径等由 install.sh 动态写入
-- ============================================

-- 基本设置
INSERT INTO `system_configs` (`group`, `key`, `value`, `type`, `name`, `sort_order`) VALUES
('general', 'company_name', '', 'string', '公司名称', 1),
('general', 'company_email', '', 'string', '公司邮箱', 2),
('general', 'company_phone', '', 'string', '公司电话', 3),
('general', 'company_address', '', 'string', '公司地址', 4),
('general', 'company_profile', '', 'text', '公司简介', 5),
('general', 'domain', '', 'string', '网站域名', 6),
('general', 'system_url', '', 'string', '系统链接', 7),
('general', 'record_no', '', 'string', '备案号', 8),
('general', 'logo_url', '/logo.png', 'string', '默认Logo', 10),
('general', 'logo_url_home', '', 'string', '前台Logo', 11),
('general', 'logo_url_bill', '', 'string', '账单Logo', 12),
('general', 'logo_url_admin', '', 'string', '后台Logo', 13),
('general', 'www_logo', '', 'string', '官网Logo', 14),
('general', 'favicon_url', '', 'string', '网站图标', 15),
('general', 'server_clause_url', '', 'string', '服务条款地址', 20),
('general', 'privacy_clause_url', '', 'string', '隐私条款地址', 21),
('general', 'cancellation_time', '7', 'int', '注销时间(天)', 22)
ON DUPLICATE KEY UPDATE `value` = VALUES(`value`);

-- 安全设置
INSERT INTO `system_configs` (`group`, `key`, `value`, `type`, `name`, `sort_order`) VALUES
('security', 'required_pwstrength', 'alpha_num', 'string', '密码强度要求', 1),
('security', 'invalid_logins_banlength', '30', 'int', '登录失败封禁时长(分钟)', 2),
('security', 'login_error_max_num', '5', 'int', '登录错误次数限制', 3),
('security', 'login_error_switch', '1', 'bool', '登录错误限制开关', 4),
('security', 'home_ip_check', '0', 'bool', '前台登录IP检测', 5),
('security', 'admin_ip_check', '0', 'bool', '后台登录IP检测', 6),
('security', 'admin_path', 'admin', 'string', '后台路径', 10)
ON DUPLICATE KEY UPDATE `value` = VALUES(`value`);

-- 登录注册
INSERT INTO `system_configs` (`group`, `key`, `value`, `type`, `name`, `sort_order`) VALUES
('login', 'allow_phone', '1', 'bool', '允许手机登录', 1),
('login', 'allow_email', '1', 'bool', '允许邮箱登录', 2),
('login', 'allow_id', '0', 'bool', '允许ID登录', 3),
('login', 'allow_register_phone', '1', 'bool', '允许手机注册', 4),
('login', 'allow_register_email', '1', 'bool', '允许邮箱注册', 5),
('login', 'allow_register_wechat', '0', 'bool', '允许微信注册', 6),
('login', 'allow_login_phone', '1', 'bool', '允许手机登录', 7),
('login', 'allow_login_email', '1', 'bool', '允许邮箱登录', 8),
('login', 'allow_login_wechat', '0', 'bool', '允许微信登录', 9),
('login', 'marketing_emails_opt_in', '1', 'bool', '营销邮件默认勾选', 17)
ON DUPLICATE KEY UPDATE `value` = VALUES(`value`);

-- 显示配置
INSERT INTO `system_configs` (`group`, `key`, `value`, `type`, `name`, `sort_order`) VALUES
('display', 'language', 'zh-cn', 'string', '默认语言', 1),
('display', 'allow_user_language', '1', 'bool', '允许用户切换语言', 2),
('display', 'date_format', 'YYYY-MM-DD', 'string', '后台日期格式', 3),
('display', 'default_country', 'CN', 'string', '默认国家', 5),
('display', 'per_page_limit', '20', 'int', '每页显示条数', 6),
('display', 'show_cancel', '1', 'bool', '显示取消按钮', 7),
('display', 'nologin_send_ticket', '0', 'bool', '未登录可发工单', 8),
('display', 'evaluate_ticket', '1', 'bool', '工单评价功能', 9)
ON DUPLICATE KEY UPDATE `value` = VALUES(`value`);

-- 维护模式
INSERT INTO `system_configs` (`group`, `key`, `value`, `type`, `name`, `sort_order`) VALUES
('maintenance', 'main_tenance_mode', '0', 'bool', '维护模式', 1),
('maintenance', 'main_tenance_mode_url', '', 'string', '维护模式重定向URL', 2),
('maintenance', 'main_tenance_mode_message', '系统维护中，请稍后再访问', 'text', '维护模式提示信息', 3)
ON DUPLICATE KEY UPDATE `value` = VALUES(`value`);

-- 代理配置
INSERT INTO `system_configs` (`group`, `key`, `value`, `type`, `name`, `sort_order`) VALUES
('affiliate', 'affiliate_enabled', '0', 'bool', '代理功能开关', 1),
('affiliate', 'affiliate_percent', '10', 'float', '代理佣金比例(%)', 2),
('affiliate', 'affiliate_payout', '100', 'float', '最低提现金额', 4),
('affiliate', 'affiliate_cookie', '30', 'float', 'Cookie有效期(天)', 5)
ON DUPLICATE KEY UPDATE `value` = VALUES(`value`);

-- 充值配置
INSERT INTO `system_configs` (`group`, `key`, `value`, `type`, `name`, `sort_order`) VALUES
('recharge', 'addfunds_enabled', '1', 'bool', '充值功能开关', 1),
('recharge', 'addfunds_minimum', '1', 'float', '单笔最小金额', 2),
('recharge', 'addfunds_maximum', '10000', 'float', '单笔最大金额', 3),
('recharge', 'addfunds_maximum_balance', '50000', 'float', '最高余额', 4)
ON DUPLICATE KEY UPDATE `value` = VALUES(`value`);

-- 信用额度
INSERT INTO `system_configs` (`group`, `key`, `value`, `type`, `name`, `sort_order`) VALUES
('credit', 'credit_limit', '0', 'bool', '信用额度开关', 1),
('credit', 'no_auto_apply_credit', '0', 'bool', '不自动使用信用额', 2)
ON DUPLICATE KEY UPDATE `value` = VALUES(`value`);

-- 发票配置
INSERT INTO `system_configs` (`group`, `key`, `value`, `type`, `name`, `sort_order`) VALUES
('invoice', 'in_circulation_create', '0', 'bool', '循环订单创建发票', 1),
('invoice', 'in_pdf', '1', 'bool', 'PDF发票', 2),
('invoice', 'in_batch_pay', '1', 'bool', '批量支付', 4),
('invoice', 'in_select_payment', '1', 'bool', '选择支付方式', 5)
ON DUPLICATE KEY UPDATE `value` = VALUES(`value`);

-- SEO配置
INSERT INTO `system_configs` (`group`, `key`, `value`, `type`, `name`, `sort_order`) VALUES
('seo', 'seo_keywords', '', 'string', 'SEO关键词', 1),
('seo', 'seo_desc', '', 'text', 'SEO描述', 2)
ON DUPLICATE KEY UPDATE `value` = VALUES(`value`);

-- 高级设置
INSERT INTO `system_configs` (`group`, `key`, `value`, `type`, `name`, `sort_order`) VALUES
('advanced', 'sendmsgtimes', '10', 'int', '每天短信发送次数', 1),
('advanced', 'sendmsgphone', '5', 'int', '每天短信发送手机个数', 2),
('advanced', 'deletelogtime', '90', 'int', '删除日志天数', 3),
('advanced', 'activity_limit', '1000', 'int', '活动日志限制', 4)
ON DUPLICATE KEY UPDATE `value` = VALUES(`value`);

-- Redis 配置（可选，默认关闭）
INSERT INTO `system_configs` (`group`, `key`, `value`, `type`, `name`, `sort_order`) VALUES
('advanced', 'redis_enabled', 'false', 'bool', 'Redis开关', 20),
('advanced', 'redis_host', '127.0.0.1', 'string', 'Redis主机', 21),
('advanced', 'redis_port', '6379', 'string', 'Redis端口', 22),
('advanced', 'redis_password', '', 'string', 'Redis密码', 23),
('advanced', 'redis_db', '0', 'string', 'Redis DB', 24)
ON DUPLICATE KEY UPDATE `value` = VALUES(`value`);

-- 上游同步
INSERT INTO `system_configs` (`group`, `key`, `value`, `type`, `name`, `sort_order`) VALUES
('upstream', 'upstream_sync_interval', '15', 'int', '自动同步间隔(分钟)', 1)
ON DUPLICATE KEY UPDATE `value` = VALUES(`value`);

-- ============================================
-- 默认工单部门
-- ============================================
INSERT INTO `departments` (`name`, `description`, `slug`, `email`, `sort_order`, `status`) VALUES
('技术支持', '技术相关问题', 'tech-support', 'support@example.com', 0, 1),
('财务部门', '账单与支付问题', 'billing', 'billing@example.com', 1, 1),
('销售咨询', '产品购买咨询', 'sales', 'sales@example.com', 2, 1)
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

-- ============================================
-- 默认用户组
-- ============================================
CREATE TABLE IF NOT EXISTS `user_groups` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `description` text,
  `discount` decimal(5,2) NOT NULL DEFAULT 0.00,
  `commission_rate` decimal(5,2) NOT NULL DEFAULT 0.00,
  `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `user_groups` (`name`, `description`, `discount`, `commission_rate`, `status`) VALUES
('默认用户组', '所有用户的默认分组', 0, 0, 1),
('VIP用户', 'VIP会员', 5, 0, 1)
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

-- ============================================
-- 默认支付网关
-- ============================================
CREATE TABLE IF NOT EXISTS `payment_gateways` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `title` varchar(100) NOT NULL,
  `gateway` varchar(50) NOT NULL,
  `code` varchar(50) DEFAULT NULL,
  `config` text,
  `is_enabled` tinyint(1) NOT NULL DEFAULT 1,
  `sort_order` int NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `payment_gateways` (`name`, `title`, `gateway`, `code`, `is_enabled`) VALUES
('alipay', '支付宝', 'alipay', 'alipay', 1),
('wechat', '微信支付', 'wechat', 'wechat', 1),
('balance', '余额支付', 'balance', 'balance', 1)
ON DUPLICATE KEY UPDATE `title` = VALUES(`title`);

-- ============================================
-- 默认邮件模板
-- ============================================
CREATE TABLE IF NOT EXISTS `email_templates` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(50) NOT NULL,
  `name` varchar(100) NOT NULL,
  `subject` varchar(200) DEFAULT NULL,
  `body` text,
  `type` varchar(20) DEFAULT 'email' COMMENT 'email|sms|notice',
  `language` varchar(10) DEFAULT 'zh-CN',
  `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code_lang` (`code`, `language`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `email_templates` (`code`, `name`, `subject`, `body`, `type`, `language`) VALUES
('user_register', '用户注册', '欢迎注册 - {{site_name}}', '<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"></head><body style="margin:0;padding:0;background:#f5f7fa;font-family:Arial,Helvetica,sans-serif;"><table width="100%" cellpadding="0" cellspacing="0" style="background:#f5f7fa;padding:40px 0;"><tr><td align="center"><table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff;border-radius:12px;overflow:hidden;box-shadow:0 2px 12px rgba(0,0,0,0.08);"><tr><td style="background:linear-gradient(135deg,#409eff,#66b1ff);padding:30px 40px;text-align:center;"><h1 style="color:#ffffff;margin:0;font-size:24px;">{{site_name}}</h1></td></tr><tr><td style="padding:40px;"><h2 style="color:#333;margin:0 0 20px;font-size:20px;">欢迎注册！</h2><p style="color:#666;line-height:1.8;margin:0 0 15px;">尊敬的 <strong>{{username}}</strong>，您好！</p><p style="color:#666;line-height:1.8;margin:0 0 25px;">感谢您注册 {{site_name}}，请验证您的邮箱以完成注册。</p><a href="{{verify_url}}" style="display:inline-block;background:#409eff;color:#fff;padding:12px 30px;border-radius:6px;text-decoration:none;font-size:16px;">验证邮箱</a><p style="color:#999;font-size:13px;margin:25px 0 0;">如果按钮无法点击，请复制链接：{{verify_url}}</p></td></tr><tr><td style="background:#f9f9f9;padding:20px 40px;text-align:center;border-top:1px solid #eee;"><p style="color:#999;font-size:12px;margin:0;">© {{year}} {{site_name}} All Rights Reserved</p></td></tr></table></td></tr></table></body></html>', 'email', 'zh-CN'),
('password_reset', '密码重置', '密码重置通知 - {{site_name}}', '<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"></head><body style="margin:0;padding:0;background:#f5f7fa;font-family:Arial,Helvetica,sans-serif;"><table width="100%" cellpadding="0" cellspacing="0" style="background:#f5f7fa;padding:40px 0;"><tr><td align="center"><table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff;border-radius:12px;overflow:hidden;box-shadow:0 2px 12px rgba(0,0,0,0.08);"><tr><td style="background:linear-gradient(135deg,#409eff,#66b1ff);padding:30px 40px;text-align:center;"><h1 style="color:#ffffff;margin:0;font-size:24px;">{{site_name}}</h1></td></tr><tr><td style="padding:40px;"><h2 style="color:#333;margin:0 0 20px;font-size:20px;">密码重置</h2><p style="color:#666;line-height:1.8;margin:0 0 15px;">尊敬的 <strong>{{username}}</strong>，您好！</p><p style="color:#666;line-height:1.8;margin:0 0 25px;">我们收到了您的密码重置请求，请点击下方按钮重置密码（30分钟内有效）。</p><a href="{{reset_link}}" style="display:inline-block;background:#409eff;color:#fff;padding:12px 30px;border-radius:6px;text-decoration:none;font-size:16px;">重置密码</a><p style="color:#999;font-size:13px;margin:25px 0 0;">如果这不是您的操作，请忽略此邮件。</p></td></tr><tr><td style="background:#f9f9f9;padding:20px 40px;text-align:center;border-top:1px solid #eee;"><p style="color:#999;font-size:12px;margin:0;">© {{year}} {{site_name}} All Rights Reserved</p></td></tr></table></td></tr></table></body></html>', 'email', 'zh-CN'),
('ticket_reply', '工单回复', '工单 #{{ticket_no}} 有新回复 - {{site_name}}', '<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"></head><body style="margin:0;padding:0;background:#f5f7fa;font-family:Arial,Helvetica,sans-serif;"><table width="100%" cellpadding="0" cellspacing="0" style="background:#f5f7fa;padding:40px 0;"><tr><td align="center"><table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff;border-radius:12px;overflow:hidden;box-shadow:0 2px 12px rgba(0,0,0,0.08);"><tr><td style="background:linear-gradient(135deg,#409eff,#66b1ff);padding:30px 40px;text-align:center;"><h1 style="color:#ffffff;margin:0;font-size:24px;">{{site_name}}</h1></td></tr><tr><td style="padding:40px;"><h2 style="color:#333;margin:0 0 20px;font-size:20px;">工单回复通知</h2><p style="color:#666;line-height:1.8;margin:0 0 15px;">尊敬的 <strong>{{username}}</strong>，您好！</p><p style="color:#666;line-height:1.8;margin:0 0 10px;">您的工单 <strong>#{{ticket_no}}</strong> 有新回复：</p><div style="background:#f5f7fa;border-left:4px solid #409eff;padding:15px 20px;margin:15px 0 25px;border-radius:4px;"><p style="color:#333;margin:0;line-height:1.6;">{{reply_content}}</p></div><a href="{{ticket_url}}" style="display:inline-block;background:#409eff;color:#fff;padding:12px 30px;border-radius:6px;text-decoration:none;font-size:16px;">查看工单</a></td></tr><tr><td style="background:#f9f9f9;padding:20px 40px;text-align:center;border-top:1px solid #eee;"><p style="color:#999;font-size:12px;margin:0;">© {{year}} {{site_name}} All Rights Reserved</p></td></tr></table></td></tr></table></body></html>', 'email', 'zh-CN'),
('order_created', '订单创建', '订单 #{{order_no}} 已创建 - {{site_name}}', '<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"></head><body style="margin:0;padding:0;background:#f5f7fa;font-family:Arial,Helvetica,sans-serif;"><table width="100%" cellpadding="0" cellspacing="0" style="background:#f5f7fa;padding:40px 0;"><tr><td align="center"><table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff;border-radius:12px;overflow:hidden;box-shadow:0 2px 12px rgba(0,0,0,0.08);"><tr><td style="background:linear-gradient(135deg,#409eff,#66b1ff);padding:30px 40px;text-align:center;"><h1 style="color:#ffffff;margin:0;font-size:24px;">{{site_name}}</h1></td></tr><tr><td style="padding:40px;"><h2 style="color:#333;margin:0 0 20px;font-size:20px;">订单已创建</h2><p style="color:#666;line-height:1.8;margin:0 0 15px;">尊敬的 <strong>{{username}}</strong>，您好！</p><p style="color:#666;line-height:1.8;margin:0 0 20px;">您的订单已创建，请尽快完成支付。</p><table width="100%" style="background:#f5f7fa;border-radius:8px;padding:20px;margin:0 0 25px;"><tr><td style="color:#999;font-size:14px;padding:5px 0;">订单编号</td><td style="color:#333;font-size:14px;font-weight:bold;text-align:right;">{{order_no}}</td></tr><tr><td style="color:#999;font-size:14px;padding:5px 0;">订单金额</td><td style="color:#ff4d4f;font-size:18px;font-weight:bold;text-align:right;">{{order_amount}}</td></tr></table><a href="{{order_url}}" style="display:inline-block;background:#409eff;color:#fff;padding:12px 30px;border-radius:6px;text-decoration:none;font-size:16px;">立即支付</a></td></tr><tr><td style="background:#f9f9f9;padding:20px 40px;text-align:center;border-top:1px solid #eee;"><p style="color:#999;font-size:12px;margin:0;">© {{year}} {{site_name}} All Rights Reserved</p></td></tr></table></td></tr></table></body></html>', 'email', 'zh-CN'),
('invoice_created', '账单创建', '账单 #{{invoice_no}} 已生成 - {{site_name}}', '<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"></head><body style="margin:0;padding:0;background:#f5f7fa;font-family:Arial,Helvetica,sans-serif;"><table width="100%" cellpadding="0" cellspacing="0" style="background:#f5f7fa;padding:40px 0;"><tr><td align="center"><table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff;border-radius:12px;overflow:hidden;box-shadow:0 2px 12px rgba(0,0,0,0.08);"><tr><td style="background:linear-gradient(135deg,#409eff,#66b1ff);padding:30px 40px;text-align:center;"><h1 style="color:#ffffff;margin:0;font-size:24px;">{{site_name}}</h1></td></tr><tr><td style="padding:40px;"><h2 style="color:#333;margin:0 0 20px;font-size:20px;">账单通知</h2><p style="color:#666;line-height:1.8;margin:0 0 15px;">尊敬的 <strong>{{username}}</strong>，您好！</p><p style="color:#666;line-height:1.8;margin:0 0 20px;">您有一笔新账单，请及时支付。</p><table width="100%" style="background:#f5f7fa;border-radius:8px;padding:20px;margin:0 0 25px;"><tr><td style="color:#999;font-size:14px;padding:5px 0;">账单编号</td><td style="color:#333;font-size:14px;font-weight:bold;text-align:right;">{{invoice_no}}</td></tr><tr><td style="color:#999;font-size:14px;padding:5px 0;">应付金额</td><td style="color:#ff4d4f;font-size:18px;font-weight:bold;text-align:right;">{{invoice_amount}}</td></tr><tr><td style="color:#999;font-size:14px;padding:5px 0;">到期时间</td><td style="color:#333;font-size:14px;text-align:right;">{{due_date}}</td></tr></table><a href="{{invoice_url}}" style="display:inline-block;background:#409eff;color:#fff;padding:12px 30px;border-radius:6px;text-decoration:none;font-size:16px;">立即支付</a></td></tr><tr><td style="background:#f9f9f9;padding:20px 40px;text-align:center;border-top:1px solid #eee;"><p style="color:#999;font-size:12px;margin:0;">© {{year}} {{site_name}} All Rights Reserved</p></td></tr></table></td></tr></table></body></html>', 'email', 'zh-CN')
ON DUPLICATE KEY UPDATE `body` = VALUES(`body`);

-- ============================================
-- 导航菜单
-- ============================================
CREATE TABLE IF NOT EXISTS `navs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `url` varchar(255) DEFAULT NULL,
  `parent_id` bigint unsigned NOT NULL DEFAULT 0,
  `order` int NOT NULL DEFAULT 0,
  `fa_icon` varchar(100) DEFAULT NULL,
  `menu_type` tinyint NOT NULL DEFAULT 1 COMMENT '1=用户中心侧栏 2=www顶部 3=www底部',
  `nav_type` tinyint NOT NULL DEFAULT 0,
  `menu_id` bigint unsigned NOT NULL DEFAULT 1,
  `is_display` tinyint(1) NOT NULL DEFAULT 1,
  `target` varchar(20) DEFAULT '_self',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 菜单激活配置
CREATE TABLE IF NOT EXISTS `menu_actives` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `menu_type` tinyint NOT NULL DEFAULT 1,
  `menuid` bigint unsigned NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- 轮播图
-- ============================================
CREATE TABLE IF NOT EXISTS `banners` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(200) DEFAULT NULL,
  `description` text,
  `type` varchar(20) DEFAULT 'image' COMMENT 'image|video',
  `media_url` varchar(500) DEFAULT NULL,
  `link_url` varchar(500) DEFAULT NULL,
  `button_text` varchar(50) DEFAULT NULL,
  `open_new` tinyint(1) NOT NULL DEFAULT 0,
  `sort_order` int NOT NULL DEFAULT 0,
  `status` tinyint NOT NULL DEFAULT 1,
  `start_time` datetime DEFAULT NULL,
  `end_time` datetime DEFAULT NULL,
  `position` varchar(20) DEFAULT 'home',
  `click_count` int NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_position` (`position`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- 语言
-- ============================================
CREATE TABLE IF NOT EXISTS `languages` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(10) NOT NULL,
  `name` varchar(50) NOT NULL,
  `flag` varchar(10) DEFAULT NULL,
  `is_default` tinyint(1) NOT NULL DEFAULT 0,
  `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- 其他核心表
-- ============================================

-- 支付记录
CREATE TABLE IF NOT EXISTS `payments` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `invoice_id` bigint unsigned DEFAULT NULL,
  `trade_no` varchar(64) DEFAULT NULL,
  `gateway` varchar(50) DEFAULT NULL,
  `amount` decimal(10,2) NOT NULL DEFAULT 0.00,
  `status` tinyint NOT NULL DEFAULT 0 COMMENT '0=待确认 1=成功 2=失败',
  `paid_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_trade_no` (`trade_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 优惠码
CREATE TABLE IF NOT EXISTS `promo_codes` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(50) NOT NULL,
  `type` varchar(20) NOT NULL DEFAULT 'percentage' COMMENT 'percentage|fixed',
  `value` decimal(10,2) NOT NULL DEFAULT 0.00,
  `cycles` varchar(20) DEFAULT 'once' COMMENT 'once|recurring|forever',
  `applies_to` text COMMENT '适用产品ID(JSON)',
  `max_times` int DEFAULT NULL,
  `used_times` int NOT NULL DEFAULT 0,
  `start_date` datetime DEFAULT NULL,
  `expire_date` datetime DEFAULT NULL,
  `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 交易记录
CREATE TABLE IF NOT EXISTS `transactions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `transaction_no` varchar(32) NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `invoice_id` bigint unsigned DEFAULT NULL,
  `gateway` varchar(50) DEFAULT NULL,
  `amount` decimal(10,2) NOT NULL DEFAULT 0.00,
  `type` varchar(20) DEFAULT 'payment' COMMENT 'payment|refund|credit|debit',
  `status` tinyint NOT NULL DEFAULT 0,
  `notes` text,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_transaction_no` (`transaction_no`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- OAuth绑定
CREATE TABLE IF NOT EXISTS `oauth_accounts` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `provider` varchar(50) NOT NULL,
  `openid` varchar(200) DEFAULT NULL,
  `unionid` varchar(200) DEFAULT NULL,
  `nickname` varchar(100) DEFAULT NULL,
  `avatar` varchar(500) DEFAULT NULL,
  `raw_data` text,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_provider_openid` (`provider`, `openid`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 系统日志
CREATE TABLE IF NOT EXISTS `system_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `level` varchar(10) DEFAULT 'info',
  `module` varchar(50) DEFAULT NULL,
  `message` text,
  `details` text,
  `user_id` bigint unsigned DEFAULT NULL,
  `ip` varchar(45) DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_level` (`level`),
  KEY `idx_module` (`module`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 下载文件
CREATE TABLE IF NOT EXISTS `download_files` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `category_id` bigint unsigned DEFAULT NULL,
  `title` varchar(200) NOT NULL,
  `description` text,
  `file_path` varchar(500) DEFAULT NULL,
  `file_size` bigint DEFAULT 0,
  `download_count` int NOT NULL DEFAULT 0,
  `sort_order` int NOT NULL DEFAULT 0,
  `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 联系人
CREATE TABLE IF NOT EXISTS `contacts` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `name` varchar(100) NOT NULL,
  `email` varchar(100) DEFAULT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `company` varchar(100) DEFAULT NULL,
  `address` varchar(255) DEFAULT NULL,
  `is_default` tinyint(1) NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 通知
CREATE TABLE IF NOT EXISTS `notifications` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `type` varchar(50) DEFAULT NULL,
  `title` varchar(200) DEFAULT NULL,
  `content` text,
  `is_read` tinyint(1) NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_is_read` (`is_read`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 上游供应商
CREATE TABLE IF NOT EXISTS `upstream_providers` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `type` varchar(50) NOT NULL COMMENT 'zjmf|v10|whmcs|custom|anchorfinance',
  `api_url` varchar(500) DEFAULT NULL,
  `api_key` varchar(255) DEFAULT NULL,
  `config` text,
  `status` tinyint NOT NULL DEFAULT 1,
  `last_sync_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 上游产品映射
CREATE TABLE IF NOT EXISTS `upstream_products` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `provider_id` bigint unsigned NOT NULL,
  `local_product_id` bigint unsigned DEFAULT NULL,
  `upstream_id` varchar(64) NOT NULL COMMENT '上游产品ID',
  `remote_product_id` varchar(64) DEFAULT NULL,
  `name` varchar(200) DEFAULT NULL,
  `synced_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_provider_id` (`provider_id`),
  KEY `idx_local_product_id` (`local_product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 知识库分类
CREATE TABLE IF NOT EXISTS `knowledge_base_categories` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `description` text,
  `parent_id` bigint unsigned DEFAULT 0,
  `sort_order` int NOT NULL DEFAULT 0,
  `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 知识库文章
CREATE TABLE IF NOT EXISTS `knowledge_base_articles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `category_id` bigint unsigned DEFAULT NULL,
  `title` varchar(200) NOT NULL,
  `content` longtext,
  `tags` varchar(500) DEFAULT NULL,
  `keywords` varchar(500) DEFAULT NULL,
  `view_count` int NOT NULL DEFAULT 0,
  `helpful_count` int NOT NULL DEFAULT 0,
  `sort_order` int NOT NULL DEFAULT 0,
  `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 新闻分类
CREATE TABLE IF NOT EXISTS `news_categories` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `parent_id` bigint unsigned NOT NULL DEFAULT 0,
  `title` varchar(100) NOT NULL,
  `name` varchar(50) NOT NULL,
  `slug` varchar(50) DEFAULT NULL,
  `sort_order` int NOT NULL DEFAULT 0,
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_slug` (`slug`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT IGNORE INTO `news_categories` (`id`, `name`, `title`, `slug`, `sort_order`) VALUES
(1, '未分类', '未分类', 'uncategorized', 0)
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

-- 新闻/公告
CREATE TABLE IF NOT EXISTS `news` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `category_id` bigint unsigned NOT NULL DEFAULT 0,
  `title` varchar(255) NOT NULL,
  `slug` varchar(255) DEFAULT NULL,
  `summary` varchar(500) DEFAULT NULL,
  `content` longtext,
  `cover_image` varchar(255) DEFAULT NULL,
  `keywords` varchar(255) DEFAULT NULL,
  `view_count` int NOT NULL DEFAULT 0,
  `is_published` tinyint(1) NOT NULL DEFAULT 1,
  `is_sticky` tinyint(1) NOT NULL DEFAULT 0,
  `published_at` datetime DEFAULT NULL,
  `admin_id` bigint unsigned DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_is_published` (`is_published`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;

-- ============================================
-- 补充：验证码配置（对齐zjmf）
-- ============================================
INSERT INTO `system_configs` (`group`, `key`, `value`, `type`, `name`, `sort_order`) VALUES
('captcha', 'is_captcha', '0', 'bool', '验证码总开关', 1),
('captcha', 'captcha_length', '5', 'int', '验证码长度', 2),
('captcha', 'captcha_combination', '1', 'int', '验证码组合（1=数字 2=字母 3=混合）', 3),
('captcha', 'allow_register_email_captcha', '1', 'bool', '邮箱注册验证码', 10),
('captcha', 'allow_register_phone_captcha', '1', 'bool', '手机注册验证码', 11),
('captcha', 'allow_login_phone_captcha', '1', 'bool', '手机登录验证码', 12),
('captcha', 'allow_login_email_captcha', '0', 'bool', '邮箱登录验证码', 13),
('captcha', 'allow_login_code_captcha', '1', 'bool', '账号登录验证码', 14),
('captcha', 'allow_login_id_captcha', '0', 'bool', 'ID登录验证码', 15),
('captcha', 'allow_phone_forgetpwd_captcha', '1', 'bool', '手机找回密码验证码', 16),
('captcha', 'allow_email_forgetpwd_captcha', '1', 'bool', '邮箱找回密码验证码', 17),
('captcha', 'allow_resetpwd_captcha', '1', 'bool', '重置密码验证码', 18),
('captcha', 'allow_phone_bind_captcha', '1', 'bool', '手机绑定验证码', 19),
('captcha', 'allow_email_bind_captcha', '1', 'bool', '邮箱绑定验证码', 20),
('captcha', 'allow_cancel_captcha', '1', 'bool', '注销账号验证码', 21),
('captcha', 'allow_login_admin_captcha', '0', 'bool', '后台登录验证码', 22)
ON DUPLICATE KEY UPDATE `value` = VALUES(`value`);

-- ============================================
-- 补充：二次验证配置（对齐zjmf）
-- ============================================
INSERT INTO `system_configs` (`group`, `key`, `value`, `type`, `name`, `sort_order`) VALUES
('security', 'second_verify', '0', 'bool', '二次验证开关', 20),
('security', 'second_verify_action', '', 'string', '二次验证操作', 21),
('security', 'second_verify_action_type', 'email', 'string', '二次验证类型', 22),
('security', 'second_verify_home', '0', 'bool', '前台二次验证', 23),
('security', 'second_verify_admin', '0', 'bool', '后台二次验证', 24),
('security', 'second_verify_action_home', '', 'string', '前台二次验证操作', 25),
('security', 'second_verify_action_home_type', '', 'string', '前台二次验证类型', 26)
ON DUPLICATE KEY UPDATE `value` = VALUES(`value`);

-- ============================================
-- 补充：短信/邮件发送开关（对齐zjmf）
-- ============================================
INSERT INTO `system_configs` (`group`, `key`, `value`, `type`, `name`, `sort_order`) VALUES
('advanced', 'allow_sms_send', '1', 'bool', '允许发送短信', 30),
('advanced', 'allow_email_send', '1', 'bool', '允许发送邮件', 31)
ON DUPLICATE KEY UPDATE `value` = VALUES(`value`);

-- ============================================
-- 补充：更多邮件模板（对齐zjmf）
-- ============================================
INSERT INTO `email_templates` (`code`, `name`, `subject`, `body`, `type`, `language`) VALUES
('product_welcome', '产品开通', '您购买的{{product_name}}已开通 - {{site_name}}', '<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"></head><body style="margin:0;padding:0;background:#f5f7fa;font-family:Arial,Helvetica,sans-serif;"><table width="100%" cellpadding="0" cellspacing="0" style="background:#f5f7fa;padding:40px 0;"><tr><td align="center"><table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff;border-radius:12px;overflow:hidden;box-shadow:0 2px 12px rgba(0,0,0,0.08);"><tr><td style="background:linear-gradient(135deg,#67c23a,#85ce61);padding:30px 40px;text-align:center;"><h1 style="color:#ffffff;margin:0;font-size:24px;">🎉 产品已开通</h1></td></tr><tr><td style="padding:40px;"><p style="color:#666;line-height:1.8;margin:0 0 15px;">尊敬的 <strong>{{username}}</strong>，您好！</p><p style="color:#666;line-height:1.8;margin:0 0 20px;">您购买的 <strong>{{product_name}}</strong> 产品已成功开通，感谢使用！</p><div style="background:#f0f9eb;border:1px solid #e1f3d8;border-radius:8px;padding:20px;margin:0 0 25px;"><p style="color:#67c23a;margin:0;font-weight:bold;">✅ 产品信息</p><p style="color:#666;margin:10px 0 0;">产品名称：{{product_name}}<br>开通时间：{{open_time}}<br>到期时间：{{expire_time}}</p></div><a href="{{product_url}}" style="display:inline-block;background:#409eff;color:#fff;padding:12px 30px;border-radius:6px;text-decoration:none;font-size:16px;">管理产品</a></td></tr><tr><td style="background:#f9f9f9;padding:20px 40px;text-align:center;border-top:1px solid #eee;"><p style="color:#999;font-size:12px;margin:0;">© {{year}} {{site_name}} All Rights Reserved</p></td></tr></table></td></tr></table></body></html>', 'email', 'zh-CN'),
('product_suspend', '产品暂停', '产品因未实名被暂停 - {{site_name}}', '<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"></head><body style="margin:0;padding:0;background:#f5f7fa;font-family:Arial,Helvetica,sans-serif;"><table width="100%" cellpadding="0" cellspacing="0" style="background:#f5f7fa;padding:40px 0;"><tr><td align="center"><table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff;border-radius:12px;overflow:hidden;box-shadow:0 2px 12px rgba(0,0,0,0.08);"><tr><td style="background:linear-gradient(135deg,#e6a23c,#f0c78a);padding:30px 40px;text-align:center;"><h1 style="color:#ffffff;margin:0;font-size:24px;">⚠️ 产品已暂停</h1></td></tr><tr><td style="padding:40px;"><p style="color:#666;line-height:1.8;margin:0 0 15px;">尊敬的 <strong>{{username}}</strong>，您好！</p><p style="color:#666;line-height:1.8;margin:0 0 20px;">您的产品 <strong>{{product_name}}</strong> 因未完成实名认证已被暂停，完成实名后将自动恢复。</p><div style="background:#fdf6ec;border:1px solid #faecd8;border-radius:8px;padding:20px;margin:0 0 25px;"><p style="color:#e6a23c;margin:0;font-weight:bold;">📋 处理方式</p><p style="color:#666;margin:10px 0 0;">请登录后台完成实名认证，认证通过后产品将自动恢复运行。</p></div><a href="{{verify_url}}" style="display:inline-block;background:#e6a23c;color:#fff;padding:12px 30px;border-radius:6px;text-decoration:none;font-size:16px;">前往实名认证</a></td></tr><tr><td style="background:#f9f9f9;padding:20px 40px;text-align:center;border-top:1px solid #eee;"><p style="color:#999;font-size:12px;margin:0;">© {{year}} {{site_name}} All Rights Reserved</p></td></tr></table></td></tr></table></body></html>', 'email', 'zh-CN'),
('product_terminate', '产品删除', '产品已到期删除 - {{site_name}}', '<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"></head><body style="margin:0;padding:0;background:#f5f7fa;font-family:Arial,Helvetica,sans-serif;"><table width="100%" cellpadding="0" cellspacing="0" style="background:#f5f7fa;padding:40px 0;"><tr><td align="center"><table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff;border-radius:12px;overflow:hidden;box-shadow:0 2px 12px rgba(0,0,0,0.08);"><tr><td style="background:linear-gradient(135deg,#f56c6c,#f89898);padding:30px 40px;text-align:center;"><h1 style="color:#ffffff;margin:0;font-size:24px;">产品已删除</h1></td></tr><tr><td style="padding:40px;"><p style="color:#666;line-height:1.8;margin:0 0 15px;">尊敬的 <strong>{{username}}</strong>，您好！</p><p style="color:#666;line-height:1.8;margin:0 0 20px;">您的产品 <strong>{{product_name}}</strong> 因到期已被删除。</p><p style="color:#666;line-height:1.8;margin:0 0 25px;">如有疑问，请联系客服。</p><a href="{{contact_url}}" style="display:inline-block;background:#409eff;color:#fff;padding:12px 30px;border-radius:6px;text-decoration:none;font-size:16px;">联系客服</a></td></tr><tr><td style="background:#f9f9f9;padding:20px 40px;text-align:center;border-top:1px solid #eee;"><p style="color:#999;font-size:12px;margin:0;">© {{year}} {{site_name}} All Rights Reserved</p></td></tr></table></td></tr></table></body></html>', 'email', 'zh-CN'),
('invoice_reminder', '账单提醒', '账单 #{{invoice_no}} 即将到期 - {{site_name}}', '<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"></head><body style="margin:0;padding:0;background:#f5f7fa;font-family:Arial,Helvetica,sans-serif;"><table width="100%" cellpadding="0" cellspacing="0" style="background:#f5f7fa;padding:40px 0;"><tr><td align="center"><table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff;border-radius:12px;overflow:hidden;box-shadow:0 2px 12px rgba(0,0,0,0.08);"><tr><td style="background:linear-gradient(135deg,#e6a23c,#f0c78a);padding:30px 40px;text-align:center;"><h1 style="color:#ffffff;margin:0;font-size:24px;">⏰ 账单提醒</h1></td></tr><tr><td style="padding:40px;"><p style="color:#666;line-height:1.8;margin:0 0 15px;">尊敬的 <strong>{{username}}</strong>，您好！</p><p style="color:#666;line-height:1.8;margin:0 0 20px;">您的账单即将到期，请及时支付以避免服务中断。</p><table width="100%" style="background:#fdf6ec;border:1px solid #faecd8;border-radius:8px;padding:20px;margin:0 0 25px;"><tr><td style="color:#999;font-size:14px;padding:5px 0;">账单编号</td><td style="color:#333;font-weight:bold;text-align:right;">{{invoice_no}}</td></tr><tr><td style="color:#999;font-size:14px;padding:5px 0;">应付金额</td><td style="color:#ff4d4f;font-size:18px;font-weight:bold;text-align:right;">{{invoice_amount}}</td></tr><tr><td style="color:#999;font-size:14px;padding:5px 0;">到期时间</td><td style="color:#e6a23c;font-weight:bold;text-align:right;">{{due_date}}</td></tr></table><a href="{{invoice_url}}" style="display:inline-block;background:#409eff;color:#fff;padding:12px 30px;border-radius:6px;text-decoration:none;font-size:16px;">立即支付</a></td></tr><tr><td style="background:#f9f9f9;padding:20px 40px;text-align:center;border-top:1px solid #eee;"><p style="color:#999;font-size:12px;margin:0;">© {{year}} {{site_name}} All Rights Reserved</p></td></tr></table></td></tr></table></body></html>', 'email', 'zh-CN'),
('invoice_overdue', '账单逾期', '账单 #{{invoice_no}} 已逾期 - {{site_name}}', '<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"></head><body style="margin:0;padding:0;background:#f5f7fa;font-family:Arial,Helvetica,sans-serif;"><table width="100%" cellpadding="0" cellspacing="0" style="background:#f5f7fa;padding:40px 0;"><tr><td align="center"><table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff;border-radius:12px;overflow:hidden;box-shadow:0 2px 12px rgba(0,0,0,0.08);"><tr><td style="background:linear-gradient(135deg,#f56c6c,#f89898);padding:30px 40px;text-align:center;"><h1 style="color:#ffffff;margin:0;font-size:24px;">账单已逾期</h1></td></tr><tr><td style="padding:40px;"><p style="color:#666;line-height:1.8;margin:0 0 15px;">尊敬的 <strong>{{username}}</strong>，您好！</p><p style="color:#666;line-height:1.8;margin:0 0 20px;">您的账单已逾期，请尽快支付以避免服务暂停。</p><table width="100%" style="background:#fef0f0;border:1px solid #fde2e2;border-radius:8px;padding:20px;margin:0 0 25px;"><tr><td style="color:#999;font-size:14px;padding:5px 0;">账单编号</td><td style="color:#333;font-weight:bold;text-align:right;">{{invoice_no}}</td></tr><tr><td style="color:#999;font-size:14px;padding:5px 0;">逾期金额</td><td style="color:#ff4d4f;font-size:18px;font-weight:bold;text-align:right;">{{invoice_amount}}</td></tr></table><a href="{{invoice_url}}" style="display:inline-block;background:#f56c6c;color:#fff;padding:12px 30px;border-radius:6px;text-decoration:none;font-size:16px;">立即支付</a></td></tr><tr><td style="background:#f9f9f9;padding:20px 40px;text-align:center;border-top:1px solid #eee;"><p style="color:#999;font-size:12px;margin:0;">© {{year}} {{site_name}} All Rights Reserved</p></td></tr></table></td></tr></table></body></html>', 'email', 'zh-CN'),
('ticket_created', '工单创建', '工单 #{{ticket_no}} 已创建 - {{site_name}}', '<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"></head><body style="margin:0;padding:0;background:#f5f7fa;font-family:Arial,Helvetica,sans-serif;"><table width="100%" cellpadding="0" cellspacing="0" style="background:#f5f7fa;padding:40px 0;"><tr><td align="center"><table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff;border-radius:12px;overflow:hidden;box-shadow:0 2px 12px rgba(0,0,0,0.08);"><tr><td style="background:linear-gradient(135deg,#409eff,#66b1ff);padding:30px 40px;text-align:center;"><h1 style="color:#ffffff;margin:0;font-size:24px;">{{site_name}}</h1></td></tr><tr><td style="padding:40px;"><h2 style="color:#333;margin:0 0 20px;font-size:20px;">工单已创建</h2><p style="color:#666;line-height:1.8;margin:0 0 15px;">尊敬的 <strong>{{username}}</strong>，您好！</p><p style="color:#666;line-height:1.8;margin:0 0 20px;">您的工单已创建，我们将尽快处理。</p><div style="background:#f5f7fa;border-left:4px solid #409eff;padding:15px 20px;margin:0 0 25px;border-radius:4px;"><p style="color:#999;font-size:13px;margin:0 0 5px;">工单 #{{ticket_no}}</p><p style="color:#333;font-weight:bold;margin:0;">{{ticket_subject}}</p></div><a href="{{ticket_url}}" style="display:inline-block;background:#409eff;color:#fff;padding:12px 30px;border-radius:6px;text-decoration:none;font-size:16px;">查看工单</a></td></tr><tr><td style="background:#f9f9f9;padding:20px 40px;text-align:center;border-top:1px solid #eee;"><p style="color:#999;font-size:12px;margin:0;">© {{year}} {{site_name}} All Rights Reserved</p></td></tr></table></td></tr></table></body></html>', 'email', 'zh-CN'),
('ticket_closed', '工单关闭', '工单 #{{ticket_no}} 已关闭 - {{site_name}}', '<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"></head><body style="margin:0;padding:0;background:#f5f7fa;font-family:Arial,Helvetica,sans-serif;"><table width="100%" cellpadding="0" cellspacing="0" style="background:#f5f7fa;padding:40px 0;"><tr><td align="center"><table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff;border-radius:12px;overflow:hidden;box-shadow:0 2px 12px rgba(0,0,0,0.08);"><tr><td style="background:linear-gradient(135deg,#909399,#b1b3b8);padding:30px 40px;text-align:center;"><h1 style="color:#ffffff;margin:0;font-size:24px;">{{site_name}}</h1></td></tr><tr><td style="padding:40px;"><h2 style="color:#333;margin:0 0 20px;font-size:20px;">工单已关闭</h2><p style="color:#666;line-height:1.8;margin:0 0 15px;">尊敬的 <strong>{{username}}</strong>，您好！</p><p style="color:#666;line-height:1.8;margin:0 0 20px;">您的工单 <strong>#{{ticket_no}}</strong> 已关闭。</p><p style="color:#666;line-height:1.8;margin:0 0 25px;">如有问题请重新提交工单。</p><a href="{{ticket_url}}" style="display:inline-block;background:#409eff;color:#fff;padding:12px 30px;border-radius:6px;text-decoration:none;font-size:16px;">查看详情</a></td></tr><tr><td style="background:#f9f9f9;padding:20px 40px;text-align:center;border-top:1px solid #eee;"><p style="color:#999;font-size:12px;margin:0;">© {{year}} {{site_name}} All Rights Reserved</p></td></tr></table></td></tr></table></body></html>', 'email', 'zh-CN'),
('order_paid', '订单支付', '订单 #{{order_no}} 支付成功 - {{site_name}}', '<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"></head><body style="margin:0;padding:0;background:#f5f7fa;font-family:Arial,Helvetica,sans-serif;"><table width="100%" cellpadding="0" cellspacing="0" style="background:#f5f7fa;padding:40px 0;"><tr><td align="center"><table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff;border-radius:12px;overflow:hidden;box-shadow:0 2px 12px rgba(0,0,0,0.08);"><tr><td style="background:linear-gradient(135deg,#67c23a,#85ce61);padding:30px 40px;text-align:center;"><h1 style="color:#ffffff;margin:0;font-size:24px;">✅ 支付成功</h1></td></tr><tr><td style="padding:40px;"><p style="color:#666;line-height:1.8;margin:0 0 15px;">尊敬的 <strong>{{username}}</strong>，您好！</p><p style="color:#666;line-height:1.8;margin:0 0 20px;">您的订单已支付成功，我们将尽快为您处理。</p><table width="100%" style="background:#f0f9eb;border:1px solid #e1f3d8;border-radius:8px;padding:20px;margin:0 0 25px;"><tr><td style="color:#999;font-size:14px;padding:5px 0;">订单编号</td><td style="color:#333;font-weight:bold;text-align:right;">{{order_no}}</td></tr><tr><td style="color:#999;font-size:14px;padding:5px 0;">支付金额</td><td style="color:#67c23a;font-size:18px;font-weight:bold;text-align:right;">{{order_amount}}</td></tr></table><a href="{{order_url}}" style="display:inline-block;background:#409eff;color:#fff;padding:12px 30px;border-radius:6px;text-decoration:none;font-size:16px;">查看订单</a></td></tr><tr><td style="background:#f9f9f9;padding:20px 40px;text-align:center;border-top:1px solid #eee;"><p style="color:#999;font-size:12px;margin:0;">© {{year}} {{site_name}} All Rights Reserved</p></td></tr></table></td></tr></table></body></html>', 'email', 'zh-CN')
ON DUPLICATE KEY UPDATE `body` = VALUES(`body`);

-- ============================================
-- 补充：导航组（对齐zjmf的shd_nav_group）
-- ============================================
CREATE TABLE IF NOT EXISTS `nav_groups` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `description` text,
  `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `nav_groups` (`id`, `name`) VALUES
(1, '云服务器'),
(2, '独立服务器'),
(3, '其他')
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

-- ============================================
-- 补充：更完整的导航菜单（对齐zjmf的shd_nav）
-- ============================================

-- 菜单激活
INSERT INTO `menu_actives` (`menu_type`, `menuid`) VALUES
(1, 1), (2, 1), (3, 1)
ON DUPLICATE KEY UPDATE `menuid` = VALUES(`menuid`);

-- 用户中心侧栏 (menu_type=1)
INSERT INTO `navs` (`id`, `name`, `url`, `parent_id`, `order`, `fa_icon`, `menu_type`, `nav_type`, `menu_id`, `is_display`) VALUES
(1, '用户中心', '/user/dashboard', 0, 0, 'bx bx-home-circle', 1, 0, 1, 1),
(2, '产品与服务', '', 0, 1, 'bx bxs-grid-alt', 1, 0, 1, 1),
(3, '订购产品', '/products', 2, 0, '', 1, 0, 1, 1),
(4, '我的服务', '/user/products', 2, 1, '', 1, 0, 1, 1),
(5, '订单管理', '/user/orders', 2, 2, '', 1, 0, 1, 1),
(6, '产品升降级', '/user/upgrade', 2, 3, '', 1, 0, 1, 1),
(7, '账户管理', '', 0, 2, 'bx bx-user', 1, 0, 1, 1),
(8, '个人信息', '/user/profile', 7, 0, '', 1, 0, 1, 1),
(9, '安全中心', '/user/security', 7, 1, '', 1, 0, 1, 1),
(10, '实名认证', '/user/verification', 7, 2, '', 1, 0, 1, 1),
(11, '消息中心', '/user/system-message', 7, 3, '', 1, 0, 1, 1),
(12, '联系人管理', '/user/contacts', 7, 4, '', 1, 0, 1, 1),
(13, '第三方登录', '/user/oauth-bind', 7, 5, '', 1, 0, 1, 1),
(14, '系统日志', '/user/system-log', 7, 6, '', 1, 0, 1, 1),
(15, '财务管理', '', 0, 3, 'bx bx-dollar-circle', 1, 0, 1, 1),
(16, '账单列表', '/user/invoices', 15, 0, '', 1, 0, 1, 1),
(17, '账户充值', '/user/wallet', 15, 1, '', 1, 0, 1, 1),
(18, '优惠码', '/user/coupons', 15, 2, '', 1, 0, 1, 1),
(19, '发票管理', '', 15, 3, '', 1, 0, 1, 1),
(20, '发票列表', '/user/invoices/list', 19, 0, '', 1, 0, 1, 1),
(21, '发票抬头', '/user/invoices/title', 19, 1, '', 1, 0, 1, 1),
(22, '技术支持', '', 0, 4, 'bx bx-detail', 1, 0, 1, 1),
(23, '工单列表', '/user/tickets', 22, 0, '', 1, 0, 1, 1),
(24, '提交工单', '/user/tickets/create', 22, 1, '', 1, 0, 1, 1),
(25, '帮助中心', '/knowledge-base', 22, 2, '', 1, 0, 1, 1),
(26, '资源下载', '/downloads', 22, 3, '', 1, 0, 1, 1),
(27, '新闻中心', '/news', 22, 4, '', 1, 0, 1, 1),
(28, '推介计划', '/user/referral', 0, 5, 'bx bxs-paper-plane', 1, 0, 1, 1)
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

-- www顶部导航 (menu_type=2)
INSERT INTO `navs` (`id`, `name`, `url`, `parent_id`, `order`, `fa_icon`, `menu_type`, `nav_type`, `menu_id`, `is_display`) VALUES
(100, '首页', '/', 0, 0, '', 2, 0, 1, 1),
(101, '产品', '', 0, 1, '', 2, 0, 1, 1),
(102, '云服务器', '/products?group=cloud', 101, 0, '', 2, 0, 1, 1),
(103, '独立服务器', '/products?group=dedicated', 101, 1, '', 2, 0, 1, 1),
(104, '全部产品', '/products', 101, 2, '', 2, 0, 1, 1),
(105, '解决方案', '/solutions', 0, 2, '', 2, 0, 1, 1),
(106, '新闻动态', '/news', 0, 3, '', 2, 0, 1, 1),
(107, '帮助支持', '', 0, 4, '', 2, 0, 1, 1),
(108, '帮助中心', '/help', 107, 0, '', 2, 0, 1, 1),
(109, '知识库', '/knowledge-base', 107, 1, '', 2, 0, 1, 1),
(110, '下载中心', '/downloads', 107, 2, '', 2, 0, 1, 1),
(111, '联系我们', '/contact', 107, 3, '', 2, 0, 1, 1)
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

-- www底部导航 (menu_type=3)
INSERT INTO `navs` (`id`, `name`, `url`, `parent_id`, `order`, `fa_icon`, `menu_type`, `nav_type`, `menu_id`, `is_display`) VALUES
(200, '产品服务', '', 0, 0, '', 3, 0, 1, 1),
(201, '帮助支持', '', 0, 1, '', 3, 0, 1, 1),
(202, '帮助中心', '/help', 201, 0, '', 3, 0, 1, 1),
(203, '知识库', '/knowledge-base', 201, 1, '', 3, 0, 1, 1),
(204, '下载中心', '/downloads', 201, 2, '', 3, 0, 1, 1),
(205, '联系我们', '/contact', 201, 3, '', 3, 0, 1, 1),
(206, '关于我们', '', 0, 2, '', 3, 0, 1, 1),
(207, '公司介绍', '/about', 206, 0, '', 3, 0, 1, 1),
(208, '新闻动态', '/news', 206, 1, '', 3, 0, 1, 1),
(209, '解决方案', '/solutions', 206, 2, '', 3, 0, 1, 1)
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

-- ============================================
-- 补充：默认货币
-- ============================================
CREATE TABLE IF NOT EXISTS `currencies` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(10) NOT NULL,
  `prefix` varchar(10) DEFAULT NULL,
  `suffix` varchar(10) DEFAULT NULL,
  `format` varchar(20) DEFAULT NULL,
  `rate` decimal(10,4) NOT NULL DEFAULT 1.0000,
  `is_default` tinyint(1) NOT NULL DEFAULT 0,
  `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `currencies` (`code`, `prefix`, `suffix`, `format`, `rate`, `is_default`, `status`) VALUES
('CNY', '¥', '元', 'prefix', 1.0000, 1, 1),
('USD', '$', '', 'prefix', 7.2000, 0, 1)
ON DUPLICATE KEY UPDATE `prefix` = VALUES(`prefix`);
