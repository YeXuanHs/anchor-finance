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

-- ============================================
-- 补充：审计日志表（对应 model/audit_log.go）
-- ============================================
CREATE TABLE IF NOT EXISTS `audit_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned DEFAULT NULL,
  `username` varchar(100) DEFAULT '',
  `user_type` varchar(20) DEFAULT '' COMMENT 'admin|client|system',
  `action` varchar(100) DEFAULT '',
  `description` varchar(500) DEFAULT '',
  `module` varchar(50) DEFAULT '',
  `controller` varchar(50) DEFAULT '',
  `method` varchar(50) DEFAULT '',
  `ip` varchar(50) DEFAULT '',
  `user_agent` varchar(500) DEFAULT '',
  `request_data` text,
  `response_code` int DEFAULT 0,
  `target_type` varchar(50) DEFAULT '',
  `target_id` bigint unsigned DEFAULT 0,
  `duration` bigint DEFAULT 0 COMMENT '耗时(毫秒)',
  `status` varchar(20) DEFAULT '' COMMENT 'success|failed',
  `remark` text,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_action` (`action`),
  KEY `idx_target_id` (`target_id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- 补充：上游同步日志表（对应 model/upstream.go UpstreamSyncLog）
-- ============================================
CREATE TABLE IF NOT EXISTS `upstream_sync_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `upstream_id` bigint unsigned DEFAULT NULL,
  `action` varchar(50) DEFAULT '',
  `status` varchar(20) DEFAULT '' COMMENT 'success|failed',
  `message` text,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_upstream_id` (`upstream_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- 修复：upstream_providers 表添加 is_active 列
-- 注意：MySQL 不支持 ADD COLUMN IF NOT EXISTS，
-- 如果列已存在会报错，可忽略
-- ============================================
ALTER TABLE `upstream_providers` ADD COLUMN `is_active` tinyint(1) NOT NULL DEFAULT 1 AFTER `config`;

-- ============================================
-- 修复：upstream_products 表对齐模型定义
-- ============================================
ALTER TABLE `upstream_products` ADD COLUMN `config` json DEFAULT NULL AFTER `remote_product_id`;

-- ============================================
-- 后台菜单表 (对齐 zjmf 7 个顶级菜单，4 级层级)
-- ============================================
CREATE TABLE IF NOT EXISTS `menus` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL COMMENT '菜单名称',
  `icon` varchar(128) DEFAULT '' COMMENT '图标',
  `url` varchar(256) DEFAULT '' COMMENT '跳转链接',
  `parent_id` bigint unsigned DEFAULT NULL COMMENT '父级ID',
  `sort_order` int NOT NULL DEFAULT 0 COMMENT '排序',
  `type` varchar(16) NOT NULL DEFAULT 'admin' COMMENT 'admin/client/top/side/bottom',
  `is_visible` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否显示',
  `is_active` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否启用',
  `permission` varchar(128) DEFAULT '' COMMENT '权限标识',
  `target` varchar(16) DEFAULT '_self' COMMENT '_self/_blank',
  `badge` varchar(32) DEFAULT '' COMMENT '角标文字',
  `badge_type` varchar(16) DEFAULT '' COMMENT 'dot/number/text',
  `language_map` json DEFAULT NULL COMMENT '多语言 {"CN":"客户","HK":"客戶","US":"Customer"}',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_type` (`type`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='后台菜单表';

-- ============================================
-- 菜单默认数据 (对齐 zjmf 7 个顶级菜单)
-- ============================================

-- 1. 客户 (Customer)
INSERT INTO `menus` (`id`, `name`, `icon`, `url`, `parent_id`, `sort_order`, `type`, `is_visible`, `is_active`, `language_map`) VALUES
(1, '客户', 'ep:User', '/customer-list', NULL, 1, 'admin', 1, 1, '{"CN":"客户","HK":"客戶","US":"Customer"}'),
(2, '客户管理', '', '', 1, 1, 'admin', 1, 1, '{"CN":"客户管理","HK":"客戶管理","US":"Customer Management"}'),
(3, '客户列表', '', '/customer-list', 2, 1, 'admin', 1, 1, '{"CN":"客户列表","HK":"客戶列表","US":"Customer List"}'),
(4, '实名认证', '', '/customer-authentication', 2, 2, 'admin', 1, 1, '{"CN":"实名认证","HK":"實名認證","US":"Real-name Authentication"}'),
(5, '客户资源池', '', '/customer-resources', 2, 3, 'admin', 1, 1, '{"CN":"客户资源池","HK":"客戶資源池","US":"Customer Resource Pool"}'),
(6, '我的业绩', '', '/sales-statistics', 1, 2, 'admin', 1, 1, '{"CN":"我的业绩","HK":"我的業績","US":"My Performance"}'),
(7, '运营管理', '', '', 1, 3, 'admin', 1, 1, '{"CN":"运营管理","HK":"運營管理","US":"Operation Management"}'),
(8, '推介计划', '', '/customer-promotionplan', 7, 1, 'admin', 1, 1, '{"CN":"推介计划","HK":"推介計劃","US":"Recommendation Plan"}'),
(9, '营销推送', '', '/marketing-push', 7, 2, 'admin', 1, 1, '{"CN":"营销推送","HK":"營銷推送","US":"Marketing Push"}');

-- 2. 业务 (Business)
INSERT INTO `menus` (`id`, `name`, `icon`, `url`, `parent_id`, `sort_order`, `type`, `is_visible`, `is_active`, `language_map`) VALUES
(10, '业务', 'ep:ShoppingCart', '/order-list', NULL, 2, 'admin', 1, 1, '{"CN":"业务","HK":"業務","US":"Business"}'),
(11, '订单', '', '', 10, 1, 'admin', 1, 1, '{"CN":"订单","HK":"訂單","US":"Order"}'),
(12, '产品订单', '', '/order-list', 11, 1, 'admin', 1, 1, '{"CN":"产品订单","HK":"產品訂單","US":"Product Order"}'),
(13, '续费订单', '', '/renewal-order', 11, 2, 'admin', 1, 1, '{"CN":"续费订单","HK":"續費訂單","US":"Renewal Order"}'),
(14, '流量包订单', '', '/dcim-traffic-log', 11, 3, 'admin', 1, 1, '{"CN":"流量包订单","HK":"流量包訂單","US":"Traffic Package Order"}'),
(15, '业务', '', '', 10, 2, 'admin', 1, 1, '{"CN":"业务","HK":"業務","US":"Business"}'),
(16, '业务列表', '', '/customer-product', 15, 1, 'admin', 1, 1, '{"CN":"业务列表","HK":"業務列表","US":"Business List"}'),
(17, '产品暂停请求', '', '/customer-cancelreq', 15, 2, 'admin', 1, 1, '{"CN":"产品暂停请求","HK":"產品暫停請求","US":"Product Pause Request"}');

-- 3. 财务 (Finance)
INSERT INTO `menus` (`id`, `name`, `icon`, `url`, `parent_id`, `sort_order`, `type`, `is_visible`, `is_active`, `language_map`) VALUES
(18, '财务', 'ep:Money', '/business-statement', NULL, 3, 'admin', 1, 1, '{"CN":"财务","HK":"財務","US":"Finance"}'),
(19, '财务记录', '', '', 18, 1, 'admin', 1, 1, '{"CN":"财务记录","HK":"財務記錄","US":"Financial Records"}'),
(20, '交易流水', '', '/business-statement', 19, 1, 'admin', 1, 1, '{"CN":"交易流水","HK":"交易流水","US":"Trading Flow"}'),
(21, '账单管理', '', '/bill-management', 19, 2, 'admin', 1, 1, '{"CN":"账单管理","HK":"賬單管理","US":"Bill Management"}'),
(22, '信用额管理', '', '/credit-management', 19, 3, 'admin', 1, 1, '{"CN":"信用额管理","HK":"信用額管理","US":"Credit Management"}'),
(23, '审核管理', '', '', 18, 2, 'admin', 1, 1, '{"CN":"审核管理","HK":"審核管理","US":"Audit Management"}'),
(24, '提现审核', '', '/customer-withdrawal', 23, 1, 'admin', 1, 1, '{"CN":"提现审核","HK":"提現審核","US":"Withdrawal Review"}'),
(25, '发票和合同', '', '', 18, 3, 'admin', 1, 1, '{"CN":"发票和合同","HK":"發票和合同","US":"Invoices and Contracts"}'),
(26, '发票列表', '', '/invoice-audit', 25, 1, 'admin', 1, 1, '{"CN":"发票列表","HK":"發票列表","US":"Invoice List"}'),
(27, '合同列表', '', '/contracts_audit', 25, 2, 'admin', 1, 1, '{"CN":"合同列表","HK":"合同列表","US":"Contracts List"}');

-- 4. 工单 (Ticket)
INSERT INTO `menus` (`id`, `name`, `icon`, `url`, `parent_id`, `sort_order`, `type`, `is_visible`, `is_active`, `language_map`) VALUES
(28, '工单', 'ep:Tickets', '/support-ticket', NULL, 4, 'admin', 1, 1, '{"CN":"工单","HK":"工單","US":"Ticket"}'),
(29, '工单', '', '', 28, 1, 'admin', 1, 1, '{"CN":"工单","HK":"工單","US":"Ticket"}'),
(30, '工单列表', '', '/support-ticket', 29, 1, 'admin', 1, 1, '{"CN":"工单列表","HK":"工單列表","US":"Ticket List"}'),
(31, '工单统计', '', '/support-statistics', 29, 2, 'admin', 1, 1, '{"CN":"工单统计","HK":"工單統計","US":"Ticket Statistics"}');

-- 5. 功能 (Function)
INSERT INTO `menus` (`id`, `name`, `icon`, `url`, `parent_id`, `sort_order`, `type`, `is_visible`, `is_active`, `language_map`) VALUES
(32, '功能', 'ep:DataAnalysis', '/timing-results', NULL, 5, 'admin', 1, 1, '{"CN":"功能","HK":"功能","US":"Function"}'),
(33, '插件', '', '', 32, 1, 'admin', 1, 1, '{"CN":"插件","HK":"插件","US":"Plugin"}'),
(34, '插件列表', '', '/plugins', 33, 1, 'admin', 1, 1, '{"CN":"插件列表","HK":"插件列表","US":"Plugin List"}'),
(35, '系统状态', '', '', 32, 2, 'admin', 1, 1, '{"CN":"系统状态","HK":"系統狀態","US":"System Status"}'),
(36, '数据库状态', '', '/database-message', 35, 1, 'admin', 1, 1, '{"CN":"数据库状态","HK":"數據庫狀態","US":"Database Status"}'),
(37, '任务队列', '', '/statistics-taskQueue', 35, 2, 'admin', 1, 1, '{"CN":"任务队列","HK":"任務隊列","US":"Task Queue"}'),
(38, '定时任务状态', '', '/timing-results', 35, 3, 'admin', 1, 1, '{"CN":"定时任务状态","HK":"定時任務狀態","US":"Timed Task Status"}'),
(39, '统计', '', '', 32, 3, 'admin', 1, 1, '{"CN":"统计","HK":"統計","US":"Statistics"}'),
(40, '年度收入统计', '', '/annual-statistics', 39, 1, 'admin', 1, 1, '{"CN":"年度收入统计","HK":"年度收入統計","US":"Annual Income Statistics"}'),
(41, '新客户', '', '/new-customer', 39, 2, 'admin', 1, 1, '{"CN":"新客户","HK":"新客戶","US":"New Customer"}'),
(42, '产品收入', '', '/product-revenue', 39, 3, 'admin', 1, 1, '{"CN":"产品收入","HK":"產品收入","US":"Product Revenue"}'),
(43, '收入排名', '', '/revenue-ranking', 39, 4, 'admin', 1, 1, '{"CN":"收入排名","HK":"收入排名","US":"Revenue Ranking"}');

-- 6. 资源与商店 (Resources And Stores)
INSERT INTO `menus` (`id`, `name`, `icon`, `url`, `parent_id`, `sort_order`, `type`, `is_visible`, `is_active`, `language_map`) VALUES
(44, '资源与商店', 'ep:Shop', '/app-store', NULL, 6, 'admin', 1, 1, '{"CN":"资源与商店","HK":"資源與商店","US":"Resources And Stores"}'),
(45, '应用商店', '', '', 44, 1, 'admin', 1, 1, '{"CN":"应用商店","HK":"應用商店","US":"App Store"}'),
(46, '我的应用', '', '/app-store', 45, 1, 'admin', 1, 1, '{"CN":"我的应用","HK":"我的應用","US":"My Application"}'),
(47, '上下游', '', '', 44, 2, 'admin', 1, 1, '{"CN":"上下游","HK":"上下游","US":"Upstream And Downstream"}'),
(48, '下游管理', '', '', 47, 1, 'admin', 1, 1, '{"CN":"下游管理","HK":"下游管理","US":"Downstream Management"}'),
(49, 'API设置', '', '/api-setup', 48, 1, 'admin', 1, 1, '{"CN":"API设置","HK":"API設置","US":"API Settings"}'),
(50, '任务队列', '', '/task-queue', 48, 2, 'admin', 1, 1, '{"CN":"任务队列","HK":"任務隊列","US":"Task Queue"}'),
(51, '上游资源', '', '', 47, 2, 'admin', 1, 1, '{"CN":"上游资源","HK":"上游資源","US":"Upstream Resources"}'),
(52, '服务器列表', '', '/munual-resource', 51, 1, 'admin', 1, 1, '{"CN":"服务器列表","HK":"服務器列表","US":"Server List"}'),
(53, '供应商管理', '', '/zjmf-api', 51, 2, 'admin', 1, 1, '{"CN":"供应商管理","HK":"供應商管理","US":"Supplier Management"}'),
(54, '商品管理', '', '/commodity-list', 51, 3, 'admin', 1, 1, '{"CN":"商品管理","HK":"商品管理","US":"Commodity Management"}'),
(55, '产品管理', '', '/commodity-product', 51, 4, 'admin', 1, 1, '{"CN":"产品管理","HK":"產品管理","US":"Product Management"}'),
(56, '任务队列', '', '/commodity-taskQueue', 51, 5, 'admin', 1, 1, '{"CN":"任务队列","HK":"任務隊列","US":"Task Queue"}'),
(57, '订单列表', '', '/supplier-order-list', 51, 6, 'admin', 1, 1, '{"CN":"订单列表","HK":"訂單列表","US":"Order List"}'),
(58, '续费订单', '', '/supplier-renewal-order', 51, 7, 'admin', 1, 1, '{"CN":"续费订单","HK":"續費訂單","US":"Renewal Order"}');

-- 7. 设置 (Settings) - 4 级菜单结构
INSERT INTO `menus` (`id`, `name`, `icon`, `url`, `parent_id`, `sort_order`, `type`, `is_visible`, `is_active`, `language_map`) VALUES
(59, '设置', 'ep:Setting', '/set', NULL, 7, 'admin', 1, 1, '{"CN":"设置","HK":"設置","US":"Settings"}'),
-- 7.1 商品设置
(60, '商品设置', '', '', 59, 1, 'admin', 1, 1, '{"CN":"商品设置","HK":"商品設置","US":"Commodity Settings"}'),
(61, '商品配置', '', '', 60, 1, 'admin', 1, 1, '{"CN":"商品配置","HK":"商品配置","US":"Commodity Configuration"}'),
(62, '商品管理', '', '/product-server', 61, 1, 'admin', 1, 1, '{"CN":"商品管理","HK":"商品管理","US":"Commodity Management"}'),
(63, '流量包管理', '', '/dcim-traffic', 61, 2, 'admin', 1, 1, '{"CN":"流量包管理","HK":"流量包管理","US":"Traffic Package Management"}'),
(64, '自动化接口', '', '', 60, 2, 'admin', 1, 1, '{"CN":"自动化接口","HK":"自動化接口","US":"Automation Interface"}'),
(65, '通用接口', '', '/server-settings', 64, 1, 'admin', 1, 1, '{"CN":"通用接口","HK":"通用接口","US":"Universal Interface"}'),
(66, '魔方DCIM', '', '/dcim', 64, 2, 'admin', 1, 1, '{"CN":"魔方DCIM","HK":"魔方DCIM","US":"Cube DCIM"}'),
(67, '魔方云', '', '/zjmfcloud', 64, 3, 'admin', 1, 1, '{"CN":"魔方云","HK":"魔方雲","US":"Magic Cube Cloud"}'),
(68, '设置', '', '', 60, 3, 'admin', 1, 1, '{"CN":"设置","HK":"設置","US":"Settings"}'),
(69, '全局可配置项', '', '/configurable-option', 68, 1, 'admin', 1, 1, '{"CN":"全局可配置项","HK":"全局可配置項","US":"Globally Configurable Items"}'),
(70, '商品订购设置', '', '/order-product', 68, 2, 'admin', 1, 1, '{"CN":"商品订购设置","HK":"商品訂購設置","US":"Product Order Setting"}'),
-- 7.2 基础设置
(71, '基础设置', '', '', 59, 2, 'admin', 1, 1, '{"CN":"基础设置","HK":"基礎設置","US":"Basic Settings"}'),
(72, '工单设置', '', '', 71, 1, 'admin', 1, 1, '{"CN":"工单设置","HK":"工單設置","US":"Ticket Settings"}'),
(73, '工单部门', '', '/work-order-dept', 72, 1, 'admin', 1, 1, '{"CN":"工单部门","HK":"工單部門","US":"Ticket Department"}'),
(74, '工单状态', '', '/work-order-status', 72, 2, 'admin', 1, 1, '{"CN":"工单状态","HK":"工單狀態","US":"Ticket Status"}'),
(75, '工单传递', '', '/work-order-rules', 72, 3, 'admin', 1, 1, '{"CN":"工单传递","HK":"工單傳遞","US":"Ticket Transfer"}'),
(76, '客户设置', '', '', 71, 2, 'admin', 1, 1, '{"CN":"客户设置","HK":"客戶設置","US":"Customer Settings"}'),
(77, '客户分组与折扣', '', '/customer-group', 76, 1, 'admin', 1, 1, '{"CN":"客户分组与折扣","HK":"客戶分組與折扣","US":"Customer Grouping and Discount"}'),
(78, '实名设置', '', '/authentication-setting', 76, 2, 'admin', 1, 1, '{"CN":"实名设置","HK":"實名設置","US":"Identity Verification"}'),
(79, '自定义客户字段', '', '/customer-custom', 76, 3, 'admin', 1, 1, '{"CN":"自定义客户字段","HK":"自定義客戶字段","US":"Custom Customer Field"}'),
(80, '推介设置', '', '/promotion_plan', 76, 4, 'admin', 1, 1, '{"CN":"推介设置","HK":"推介設置","US":"Recommendation Settings"}'),
(81, '客户等级', '', '/customer-level', 76, 5, 'admin', 1, 1, '{"CN":"客户等级","HK":"客戶等級","US":"Customer Level"}'),
(82, '财务设置', '', '', 71, 3, 'admin', 1, 1, '{"CN":"财务设置","HK":"財務設置","US":"Financial Settings"}'),
(83, '支付接口', '', '/payment-interface', 82, 1, 'admin', 1, 1, '{"CN":"支付接口","HK":"支付接口","US":"Payment Interface"}'),
(84, '优惠码', '', '/promo-code', 82, 2, 'admin', 1, 1, '{"CN":"优惠码","HK":"優惠碼","US":"Promotion Code"}'),
(85, '货币配置', '', '/currency-settings', 82, 3, 'admin', 1, 1, '{"CN":"货币配置","HK":"貨幣配置","US":"Currency Configuration"}'),
(86, '充值与财务', '', '/general-settings/finance', 82, 4, 'admin', 1, 1, '{"CN":"充值与财务","HK":"充值與財務","US":"Recharge and Finance"}'),
(87, '发票设置', '', '/voucher-setting', 82, 5, 'admin', 1, 1, '{"CN":"发票设置","HK":"發票設置","US":"Invoice Settings"}'),
(88, '信用额设置', '', '/credit-setting', 82, 6, 'admin', 1, 1, '{"CN":"信用额设置","HK":"信用額設定","US":"Credit Limit Setting"}'),
(89, '合同设置', '', '/contracts_setting', 82, 7, 'admin', 1, 1, '{"CN":"合同设置","HK":"合同設置","US":"Contracts Settings"}'),
-- 7.3 站务设置
(90, '站务设置', '', '', 59, 3, 'admin', 1, 1, '{"CN":"站务设置","HK":"站務設置","US":"Station Service Settings"}'),
(91, '显示设置', '', '/base-info', 90, 1, 'admin', 1, 1, '{"CN":"显示设置","HK":"顯示設置","US":"Display Settings"}'),
(92, '文件下载', '', '/service-support', 90, 2, 'admin', 1, 1, '{"CN":"文件下载","HK":"文件下載","US":"File Download"}'),
(93, '新闻中心', '', '/news-center', 90, 3, 'admin', 1, 1, '{"CN":"新闻中心","HK":"新聞中心","US":"News Center"}'),
(94, '帮助中心', '', '/help-center', 90, 4, 'admin', 1, 1, '{"CN":"帮助中心","HK":"幫助中心","US":"Help Center"}'),
-- 7.4 系统设置
(95, '系统设置', '', '', 59, 4, 'admin', 1, 1, '{"CN":"系统设置","HK":"系統設置","US":"System Settings"}'),
(96, '基础设置', '', '', 95, 1, 'admin', 1, 1, '{"CN":"基础设置","HK":"基礎設置","US":"Basic Settings"}'),
(97, '常规设置', '', '/general-settings/general', 96, 1, 'admin', 1, 1, '{"CN":"常规设置","HK":"常規設置","US":"General Settings"}'),
(98, '定时任务', '', '/automatic-tasks', 96, 2, 'admin', 1, 1, '{"CN":"定时任务","HK":"定時任務","US":"Timed Task"}'),
(99, '注册登录', '', '/login-register', 96, 3, 'admin', 1, 1, '{"CN":"注册登录","HK":"註冊登錄","US":"Register and Login"}'),
(100, '第三方登录', '', '/third-login', 96, 4, 'admin', 1, 1, '{"CN":"第三方登录","HK":"第三方登錄","US":"Third Party Login"}'),
(101, '人员管理', '', '', 95, 2, 'admin', 1, 1, '{"CN":"人员管理","HK":"人員管理","US":"Personnel Management"}'),
(102, '员工管理', '', '/admin-management', 101, 1, 'admin', 1, 1, '{"CN":"员工管理","HK":"員工管理","US":"Staff Management"}'),
(103, '分组权限', '', '/permissions-managment', 101, 2, 'admin', 1, 1, '{"CN":"分组权限","HK":"分組權限","US":"Group Permission"}'),
(104, '销售设置', '', '/sales-management', 101, 3, 'admin', 1, 1, '{"CN":"销售设置","HK":"銷售設置","US":"Sales Settings"}'),
(105, '短信邮件设置', '', '', 95, 3, 'admin', 1, 1, '{"CN":"短信邮件设置","HK":"短信郵件設置","US":"SMS Mail Settings"}'),
(106, '接口设置', '', '/sms-template/sms', 105, 1, 'admin', 1, 1, '{"CN":"接口设置","HK":"接口設置","US":"Interface Settings"}'),
(107, '邮件模板', '', '/email-list', 105, 2, 'admin', 1, 1, '{"CN":"邮件模板","HK":"郵件模板","US":"Mail Template"}'),
(108, '短信模板', '', '/sms-template-index', 105, 3, 'admin', 1, 1, '{"CN":"短信模板","HK":"短信模板","US":"SMS Template"}'),
(109, '发送设置', '', '/sms-send-settings', 105, 4, 'admin', 1, 1, '{"CN":"发送设置","HK":"發送設置","US":"Send Settings"}'),
(110, '安全相关', '', '', 95, 4, 'admin', 1, 1, '{"CN":"安全相关","HK":"安全相關","US":"Security Related"}'),
(111, '黑名单列表', '', '/black-list', 110, 1, 'admin', 1, 1, '{"CN":"黑名单列表","HK":"黑名單列表","US":"Blacklist List"}'),
(112, '验证码设置', '', '/general-settings/captcha', 110, 2, 'admin', 1, 1, '{"CN":"验证码设置","HK":"驗證碼設置","US":"Verification Code Settings"}'),
(113, '二次验证', '', '/twice-confirm', 110, 3, 'admin', 1, 1, '{"CN":"二次验证","HK":"二次驗證","US":"Secondary Verification"}'),
(114, '系统相关', '', '', 95, 5, 'admin', 1, 1, '{"CN":"系统相关","HK":"系統相關","US":"System Related"}'),
(115, '系统升级', '', '/system-message', 114, 1, 'admin', 1, 1, '{"CN":"系统升级","HK":"系統升級","US":"System Upgrade"}'),
(116, '关于', '', '/about', 114, 2, 'admin', 1, 1, '{"CN":"关于","HK":"關於","US":"About"}'),
(117, '日志', '', '', 95, 6, 'admin', 1, 1, '{"CN":"日志","HK":"日誌","US":"Log"}'),
(118, '系统日志', '', '/system-log', 117, 1, 'admin', 1, 1, '{"CN":"系统日志","HK":"系統日誌","US":"System Log"}'),
(119, '管理员登录日志', '', '/system-admin-log', 117, 2, 'admin', 1, 1, '{"CN":"管理员登录日志","HK":"管理員登錄日誌","US":"Administrator Login Log"}'),
(120, '邮件日志', '', '/email-log', 117, 3, 'admin', 1, 1, '{"CN":"邮件日志","HK":"郵件日誌","US":"Mail Log"}'),
(121, '短信日志', '', '/sms-log', 117, 4, 'admin', 1, 1, '{"CN":"短信日志","HK":"短信日誌","US":"SMS Log"}'),
(122, '站内信日志', '', '/station-letter-log', 117, 5, 'admin', 1, 1, '{"CN":"站内信日志","HK":"站內信日誌","US":"Site Message Log"}'),
(123, '定时任务日志', '', '/automatic-task-log', 117, 6, 'admin', 1, 1, '{"CN":"定时任务日志","HK":"定時任務日誌","US":"Timed Task Log"}'),
(124, 'API日志', '', '/api-log', 117, 7, 'admin', 1, 1, '{"CN":"API日志","HK":"API日誌","US":"API Log"}'),
(125, '日志清理', '', '/log-cleanup', 117, 8, 'admin', 1, 1, '{"CN":"日志清理","HK":"日誌清理","US":"Log Cleanup"}')
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);
