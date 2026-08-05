-- ============================================================================
-- zjmf (v376) → AnchorFinance 数据库迁移脚本
-- ============================================================================
-- 执行前请确保：
--   1. 已备份 zjmf 数据库
--   2. AnchorFinance 新表结构已由 GORM AutoMigrate 创建
--   3. 本脚本在 zjmf 数据库上下文中执行
--   4. zjmf 原始表数据不会被删除（只读迁移）
-- ============================================================================

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================================
-- 第一部分：创建 AnchorFinance 新增表（zjmf 中不存在的）
-- ============================================================================

-- ----------------------------------------------------------------------------
-- 1.1 user_groups 用户组表
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `user_groups` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(64) NOT NULL,
  `description` text,
  `discount` decimal(5,4) NOT NULL DEFAULT 1.0000,
  `commission_rate` decimal(5,4) NOT NULL DEFAULT 0.0000,
  `is_default` tinyint(1) DEFAULT 0,
  `sort_order` int DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_groups_name` (`name`),
  KEY `idx_user_groups_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 插入默认用户组
INSERT IGNORE INTO `user_groups` (`id`, `name`, `description`, `is_default`, `sort_order`, `created_at`, `updated_at`) VALUES
(1, '默认用户组', '系统默认用户组', 1, 0, NOW(), NOW()),
(2, 'VIP用户', 'VIP用户组', 0, 1, NOW(), NOW()),
(3, '代理用户', '代理用户组', 0, 2, NOW(), NOW());

-- ----------------------------------------------------------------------------
-- 1.2 admins 管理员表
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `admins` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `username` varchar(64) NOT NULL,
  `email` varchar(255) NOT NULL,
  `password_hash` varchar(255) NOT NULL,
  `salt` varchar(32) NOT NULL,
  `nickname` varchar(64) DEFAULT NULL,
  `avatar` varchar(512) DEFAULT NULL,
  `role` varchar(32) NOT NULL DEFAULT 'admin',
  `permissions` json DEFAULT NULL,
  `is_super` tinyint(1) DEFAULT 0,
  `status` smallint NOT NULL DEFAULT 1,
  `last_login_at` datetime(3) DEFAULT NULL,
  `last_login_ip` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_admins_username` (`username`),
  UNIQUE KEY `idx_admins_email` (`email`),
  KEY `idx_admins_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------------------------------------------------------
-- 1.3 product_first_groups 产品一级分组表
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `product_first_groups` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(128) NOT NULL,
  `hidden` tinyint(1) DEFAULT 0,
  `sort_order` int DEFAULT 0,
  `upstream_id` int unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_product_first_groups_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------------------------------------------------------
-- 1.4 user_products 用户产品/服务实例表
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `user_products` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` int unsigned NOT NULL,
  `order_id` int unsigned DEFAULT NULL,
  `product_id` int unsigned NOT NULL,
  `name` varchar(256) NOT NULL,
  `domain` varchar(256) DEFAULT NULL,
  `username` varchar(128) DEFAULT NULL,
  `password` varchar(256) DEFAULT NULL,
  `ip` varchar(64) DEFAULT NULL,
  `dedicated_ip` varchar(64) DEFAULT NULL,
  `hostname` varchar(256) DEFAULT NULL,
  `ns1` varchar(256) DEFAULT NULL,
  `ns2` varchar(256) DEFAULT NULL,
  `billing_cycle` varchar(32) DEFAULT NULL,
  `amount` decimal(20,4) NOT NULL DEFAULT 0.0000,
  `currency` varchar(8) DEFAULT 'CNY',
  `registration_date` datetime(3) DEFAULT NULL,
  `next_due_date` datetime(3) DEFAULT NULL,
  `termination_date` datetime(3) DEFAULT NULL,
  `suspend_reason` varchar(256) DEFAULT NULL,
  `status` smallint NOT NULL DEFAULT 1,
  `provisioning_status` smallint DEFAULT 0,
  `admin_notes` text,
  `notes` text,
  `config_options` json DEFAULT NULL,
  `custom_fields` json DEFAULT NULL,
  `metadata` json DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_products_user_id` (`user_id`),
  KEY `idx_user_products_order_id` (`order_id`),
  KEY `idx_user_products_product_id` (`product_id`),
  KEY `idx_user_products_domain` (`domain`),
  KEY `idx_user_products_next_due_date` (`next_due_date`),
  KEY `idx_user_products_status` (`status`),
  KEY `idx_user_products_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------------------------------------------------------
-- 1.5 transactions 支付交易记录表
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `transactions` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `transaction_no` varchar(64) NOT NULL,
  `user_id` int unsigned NOT NULL,
  `invoice_id` int unsigned DEFAULT NULL,
  `order_id` int unsigned DEFAULT NULL,
  `gateway` varchar(64) NOT NULL,
  `gateway_trans_id` varchar(256) DEFAULT NULL,
  `amount` decimal(20,4) NOT NULL,
  `fee` decimal(20,4) DEFAULT 0.0000,
  `currency` varchar(8) NOT NULL DEFAULT 'CNY',
  `exchange_rate` decimal(16,8) DEFAULT 1.00000000,
  `type` varchar(32) NOT NULL DEFAULT 'payment',
  `status` smallint NOT NULL DEFAULT 0,
  `completed_at` datetime(3) DEFAULT NULL,
  `notes` text,
  `admin_notes` text,
  `ip_address` varchar(64) DEFAULT NULL,
  `callback_data` json DEFAULT NULL,
  `metadata` json DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_transactions_transaction_no` (`transaction_no`),
  KEY `idx_transactions_user_id` (`user_id`),
  KEY `idx_transactions_invoice_id` (`invoice_id`),
  KEY `idx_transactions_order_id` (`order_id`),
  KEY `idx_transactions_gateway_trans_id` (`gateway_trans_id`),
  KEY `idx_transactions_status` (`status`),
  KEY `idx_transactions_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------------------------------------------------------
-- 1.6 attachments 附件表
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `attachments` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `ticket_id` int unsigned DEFAULT NULL,
  `reply_id` int unsigned DEFAULT NULL,
  `uploader_id` int unsigned NOT NULL,
  `file_name` varchar(256) NOT NULL,
  `file_path` varchar(512) NOT NULL,
  `file_size` bigint NOT NULL,
  `mime_type` varchar(128) DEFAULT NULL,
  `storage_driver` varchar(32) DEFAULT 'local',
  `storage_key` varchar(512) DEFAULT NULL,
  `download_count` int DEFAULT 0,
  `is_public` tinyint(1) DEFAULT 0,
  `hash` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_attachments_ticket_id` (`ticket_id`),
  KEY `idx_attachments_reply_id` (`reply_id`),
  KEY `idx_attachments_uploader_id` (`uploader_id`),
  KEY `idx_attachments_hash` (`hash`),
  KEY `idx_attachments_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------------------------------------------------------
-- 1.7 download_categories 下载分类表
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `download_categories` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(128) NOT NULL,
  `description` text,
  `slug` varchar(128) DEFAULT NULL,
  `parent_id` int unsigned DEFAULT NULL,
  `sort_order` int DEFAULT 0,
  `status` smallint DEFAULT 1,
  PRIMARY KEY (`id`),
  KEY `idx_download_categories_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------------------------------------------------------
-- 1.8 product_pricings 产品定价表
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `product_pricings` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `product_id` int unsigned NOT NULL,
  `cycle` varchar(32) NOT NULL,
  `price` decimal(20,4) NOT NULL,
  `setup_fee` decimal(20,4) DEFAULT 0.0000,
  `currency` varchar(8) DEFAULT 'CNY',
  `sort_order` int DEFAULT 0,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_product_pricings_product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------------------------------------------------------
-- 1.9 order_notes 订单备注表
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `order_notes` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `order_id` int unsigned NOT NULL,
  `admin_id` int unsigned DEFAULT NULL,
  `content` text NOT NULL,
  `is_private` tinyint(1) DEFAULT 1,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_order_notes_order_id` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------------------------------------------------------
-- 1.10 product_downloads 产品下载关联表
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `product_downloads` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `product_id` int unsigned NOT NULL,
  `download_id` int unsigned NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_product_downloads_product_id` (`product_id`),
  KEY `idx_product_downloads_download_id` (`download_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------------------------------------------------------
-- 1.11 login_logs 登录日志表
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `login_logs` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` int unsigned NOT NULL,
  `ip` varchar(64) NOT NULL,
  `user_agent` varchar(512) DEFAULT NULL,
  `location` varchar(256) DEFAULT NULL,
  `device` varchar(128) DEFAULT NULL,
  `status` smallint NOT NULL,
  `remark` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_login_logs_user_id` (`user_id`),
  KEY `idx_login_logs_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------------------------------------------------------
-- 1.12 ticket_transfer_logs 工单转移日志表
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `ticket_transfer_logs` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `ticket_id` int unsigned NOT NULL,
  `from_dept_id` int unsigned DEFAULT NULL,
  `to_dept_id` int unsigned DEFAULT NULL,
  `from_agent_id` int unsigned DEFAULT NULL,
  `to_agent_id` int unsigned DEFAULT NULL,
  `operator_id` int unsigned NOT NULL,
  `reason` text,
  PRIMARY KEY (`id`),
  KEY `idx_ticket_transfer_logs_ticket_id` (`ticket_id`),
  KEY `idx_ticket_transfer_logs_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------------------------------------------------------
-- 1.13 ticket_notes 工单内部备注表
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `ticket_notes` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `ticket_id` int unsigned NOT NULL,
  `admin_id` int unsigned NOT NULL,
  `content` text NOT NULL,
  `attachment` text,
  PRIMARY KEY (`id`),
  KEY `idx_ticket_notes_ticket_id` (`ticket_id`),
  KEY `idx_ticket_notes_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------------------------------------------------------
-- 1.14 ticket_custom_fields 工单自定义字段表
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `ticket_custom_fields` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `type` varchar(32) NOT NULL,
  `rel_id` int unsigned NOT NULL,
  `field_name` varchar(128) NOT NULL,
  `field_type` varchar(32) DEFAULT 'text',
  `description` varchar(255) DEFAULT NULL,
  `field_option` text,
  `reg_expr` varchar(255) DEFAULT NULL,
  `admin_only` tinyint DEFAULT 0,
  `required` tinyint DEFAULT 0,
  `sort_order` int DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `idx_ticket_custom_fields_type_rel_id` (`type`, `rel_id`),
  KEY `idx_ticket_custom_fields_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------------------------------------------------------
-- 1.15 ticket_custom_field_values 工单自定义字段值表
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `ticket_custom_field_values` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `field_id` int unsigned NOT NULL,
  `rel_id` int unsigned NOT NULL,
  `value` text,
  PRIMARY KEY (`id`),
  KEY `idx_ticket_custom_field_values_field_id` (`field_id`),
  KEY `idx_ticket_custom_field_values_rel_id` (`rel_id`),
  KEY `idx_ticket_custom_field_values_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------------------------------------------------------
-- 1.16 agents 代理表
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `agents` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` int unsigned NOT NULL,
  `agent_no` varchar(20) NOT NULL,
  `parent_id` int unsigned DEFAULT NULL,
  `level` int DEFAULT 1,
  `commission_rate` decimal(5,2) DEFAULT 10.00,
  `balance` decimal(10,2) DEFAULT 0.00,
  `total_earned` decimal(10,2) DEFAULT 0.00,
  `status` smallint DEFAULT 1,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_agents_user_id` (`user_id`),
  UNIQUE KEY `idx_agents_agent_no` (`agent_no`),
  KEY `idx_agents_parent_id` (`parent_id`),
  KEY `idx_agents_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------------------------------------------------------
-- 1.17 agent_commissions 代理佣金记录表
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `agent_commissions` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `agent_id` int unsigned NOT NULL,
  `order_id` int unsigned DEFAULT NULL,
  `user_id` int unsigned DEFAULT NULL,
  `amount` decimal(10,2) NOT NULL,
  `rate` decimal(5,2) DEFAULT NULL,
  `status` smallint DEFAULT 1,
  PRIMARY KEY (`id`),
  KEY `idx_agent_commissions_agent_id` (`agent_id`),
  KEY `idx_agent_commissions_order_id` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------------------------------------------------------
-- 1.18 custom_fields 自定义字段表
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `custom_fields` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(128) NOT NULL,
  `label` varchar(256) NOT NULL,
  `type` varchar(32) NOT NULL,
  `group_name` varchar(32) DEFAULT NULL,
  `required` tinyint(1) DEFAULT 0,
  `default_val` text,
  `options` json DEFAULT NULL,
  `validation` varchar(256) DEFAULT NULL,
  `placeholder` varchar(256) DEFAULT NULL,
  `help_text` varchar(512) DEFAULT NULL,
  `sort_order` int DEFAULT 0,
  `enabled` tinyint(1) DEFAULT 1,
  PRIMARY KEY (`id`),
  KEY `idx_custom_fields_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------------------------------------------------------
-- 1.19 custom_field_values 自定义字段值表
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `custom_field_values` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `field_id` int unsigned NOT NULL,
  `owner_id` int unsigned NOT NULL,
  `owner_type` varchar(32) DEFAULT NULL,
  `value` text,
  PRIMARY KEY (`id`),
  KEY `idx_custom_field_values_field_id` (`field_id`),
  KEY `idx_custom_field_values_owner_id` (`owner_id`),
  KEY `idx_custom_field_values_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------------------------------------------------------
-- 1.20 host_operations 主机操作记录表
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `host_operations` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `host_id` int unsigned NOT NULL,
  `operator_id` int unsigned NOT NULL,
  `action` varchar(32) NOT NULL,
  `params` text,
  `status` smallint DEFAULT 1,
  `result` text,
  `error_msg` text,
  `started_at` datetime(3) NOT NULL,
  `finished_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_host_operations_host_id` (`host_id`),
  KEY `idx_host_operations_operator_id` (`operator_id`),
  KEY `idx_host_operations_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------------------------------------------------------
-- 1.21 languages 语言表
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `languages` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `code` varchar(16) NOT NULL,
  `name` varchar(64) NOT NULL,
  `flag` varchar(8) DEFAULT NULL,
  `is_default` tinyint(1) DEFAULT 0,
  `status` smallint DEFAULT 1,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_languages_code` (`code`),
  KEY `idx_languages_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------------------------------------------------------
-- 1.22 banners 轮播图表
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `banners` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `title` varchar(256) DEFAULT NULL,
  `description` text,
  `type` varchar(16) DEFAULT 'image',
  `media_url` varchar(512) DEFAULT NULL,
  `link_url` varchar(512) DEFAULT NULL,
  `button_text` varchar(64) DEFAULT NULL,
  `sort_order` int DEFAULT 0,
  `status` smallint DEFAULT 1,
  `position` varchar(32) DEFAULT 'home',
  PRIMARY KEY (`id`),
  KEY `idx_banners_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------------------------------------------------------
-- 1.23 navs 导航菜单表
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `navs` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(128) NOT NULL,
  `url` varchar(512) DEFAULT NULL,
  `parent_id` int unsigned DEFAULT 0,
  `sort_order` int DEFAULT 0,
  `fa_icon` varchar(64) DEFAULT NULL,
  `menu_type` smallint DEFAULT 1,
  `nav_type` smallint DEFAULT 0,
  `menu_id` int unsigned DEFAULT 1,
  `is_display` tinyint(1) DEFAULT 1,
  PRIMARY KEY (`id`),
  KEY `idx_navs_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------------------------------------------------------
-- 1.24 menu_actives 菜单激活表
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `menu_actives` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `menu_type` smallint DEFAULT NULL,
  `menuid` int unsigned DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------------------------------------------------------
-- 1.25 email_template_logs 邮件模板发送日志表
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `email_template_logs` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `template_id` int unsigned NOT NULL,
  `recipient` varchar(256) NOT NULL,
  `subject` varchar(256) DEFAULT NULL,
  `content` text,
  `type` varchar(16) NOT NULL,
  `status` smallint DEFAULT 1,
  `error` text,
  `sent_at` datetime(3) DEFAULT NULL,
  `operator_id` int unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_email_template_logs_template_id` (`template_id`),
  KEY `idx_email_template_logs_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------------------------------------------------------
-- 1.26 promo_code_logs 优惠码使用记录表
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `promo_code_logs` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `promo_id` int unsigned NOT NULL,
  `user_id` int unsigned NOT NULL,
  `order_id` int unsigned DEFAULT NULL,
  `amount` decimal(10,2) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_promo_code_logs_promo_id` (`promo_id`),
  KEY `idx_promo_code_logs_user_id` (`user_id`),
  KEY `idx_promo_code_logs_order_id` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------------------------------------------------------
-- 1.27 upstream_sync_logs 上游同步日志表
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `upstream_sync_logs` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `upstream_id` int unsigned DEFAULT NULL,
  `action` varchar(50) DEFAULT NULL,
  `status` varchar(20) DEFAULT NULL,
  `message` text,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_upstream_sync_logs_upstream_id` (`upstream_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------------------------------------------------------
-- 1.28 server_templates 服务器配置模板表
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `server_templates` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(128) NOT NULL,
  `code` varchar(64) NOT NULL,
  `type` varchar(32) NOT NULL,
  `description` text,
  `config` json NOT NULL,
  `sort_order` int DEFAULT 0,
  `status` smallint DEFAULT 1,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_server_templates_code` (`code`),
  KEY `idx_server_templates_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------------------------------------------------------
-- 1.29 server_products 服务器产品关联表
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `server_products` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `server_config_id` int unsigned NOT NULL,
  `product_id` int unsigned NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_server_products_server_config_id` (`server_config_id`),
  KEY `idx_server_products_product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------------------------------------------------------
-- 1.30 accounts 交易流水表
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `accounts` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `uid` int unsigned NOT NULL,
  `currency` varchar(16) DEFAULT 'CNY',
  `gateway` varchar(64) DEFAULT NULL,
  `pay_time` bigint DEFAULT NULL,
  `update_time` bigint DEFAULT NULL,
  `create_time` bigint DEFAULT NULL,
  `delete_time` bigint DEFAULT 0,
  `description` text,
  `trans_id` varchar(128) DEFAULT NULL,
  `invoice_id` int unsigned DEFAULT NULL,
  `amount_in` decimal(20,4) DEFAULT 0.0000,
  `fees` decimal(20,4) DEFAULT 0.0000,
  `amount_out` decimal(20,4) DEFAULT 0.0000,
  `rate` decimal(20,8) DEFAULT 1.00000000,
  `refund` int DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `idx_accounts_uid` (`uid`),
  KEY `idx_accounts_invoice_id` (`invoice_id`),
  KEY `idx_accounts_trans_id` (`trans_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- ============================================================================
-- 第二部分：数据迁移（从 shd_* 表 → AF 表）
-- ============================================================================

-- ----------------------------------------------------------------------------
-- 2.1 用户迁移：shd_clients → users
-- 字段映射：
--   id → id (保留原ID)
--   uuid → uuid
--   username → username
--   email → email
--   phonenumber → phone
--   password → password_hash (zjmf密码格式需应用层二次处理)
--   '' → salt (zjmf无独立salt字段，需应用层处理)
--   credit → balance
--   status → status (1→1, 0→0, 其他→0)
--   currency → currency
--   groupid → group_id
--   create_time(int) → created_at(FROM_UNIXTIME)
--   update_time(int) → updated_at(FROM_UNIXTIME)
--   delete_time(int) → deleted_at(FROM_UNIXTIME, 0→NULL)
-- ----------------------------------------------------------------------------
INSERT IGNORE INTO `users` (
  `id`, `uuid`, `username`, `email`, `password_hash`, `salt`,
  `phone`, `currency`, `balance`, `group_id`, `status`, `is_admin`,
  `created_at`, `updated_at`, `deleted_at`
)
SELECT
  `id`,
  IFNULL(`uuid`, ''),
  `username`,
  IFNULL(`email`, CONCAT('user_', `id`, '@migrated.local')),
  IFNULL(`password`, ''),
  '',  -- salt: zjmf无独立salt，留空待应用层处理
  IFNULL(`phonenumber`, ''),
  IFNULL(`currency`, 'CNY'),
  IFNULL(`credit`, 0),
  IFNULL(`groupid`, 1),
  CASE
    WHEN `status` = 1 THEN 1
    WHEN `status` = 0 THEN 0
    ELSE 0
  END,
  0,  -- is_admin: zjmf客户端用户都不是管理员
  FROM_UNIXTIME(IFNULL(`create_time`, UNIX_TIMESTAMP())),
  FROM_UNIXTIME(IFNULL(`update_time`, UNIX_TIMESTAMP())),
  CASE
    WHEN `delete_time` IS NOT NULL AND `delete_time` > 0
    THEN FROM_UNIXTIME(`delete_time`)
    ELSE NULL
  END
FROM `shd_clients`
WHERE `id` IS NOT NULL;

-- ----------------------------------------------------------------------------
-- 2.2 产品组迁移：shd_product_groups → product_groups
-- 字段映射：
--   id → id
--   name → name
--   description → description
--   hidden → status (hidden=1→status=0, hidden=0→status=1)
--   '' → slug (zjmf无slug，用name生成)
--   NULL → parent_id
-- ----------------------------------------------------------------------------
INSERT IGNORE INTO `product_groups` (
  `id`, `name`, `description`, `slug`, `status`,
  `created_at`, `updated_at`
)
SELECT
  `id`,
  `name`,
  IFNULL(`description`, ''),
  CONCAT('group-', `id`),
  CASE WHEN `hidden` = 1 THEN 0 ELSE 1 END,
  NOW(), NOW()
FROM `shd_product_groups`
WHERE `id` IS NOT NULL;

-- ----------------------------------------------------------------------------
-- 2.3 产品迁移：shd_products → products
-- 字段映射：
--   id → id
--   gid → group_id
--   name → name
--   description → description
--   type → type (保留原值)
--   auto_setup → auto_setup (1→true, 0→false)
--   stock_control → stock_control (1→true, 0→false)
--   pay_type → billing_cycle (需应用层映射)
--   '' → slug (zjmf无slug)
--   0 → price (zjmf产品价格在配置中，需应用层处理)
--   'CNY' → currency
-- ----------------------------------------------------------------------------
INSERT IGNORE INTO `products` (
  `id`, `group_id`, `name`, `slug`, `description`, `type`,
  `auto_setup`, `stock_control`, `billing_cycle`, `currency`,
  `price`, `status`, `created_at`, `updated_at`
)
SELECT
  `id`,
  IFNULL(`gid`, 0),
  `name`,
  CONCAT('product-', `id`),
  IFNULL(`description`, ''),
  IFNULL(`type`, 'other'),
  IF(`auto_setup` = 1, 1, 0),
  IF(`stock_control` = 1, 1, 0),
  IFNULL(`pay_type`, 'monthly'),
  'CNY',
  0,  -- price: zjmf产品价格通常在配置中
  1,  -- status: 默认上架
  NOW(), NOW()
FROM `shd_products`
WHERE `id` IS NOT NULL;

-- ----------------------------------------------------------------------------
-- 2.4 订单迁移：shd_orders → orders
-- 字段映射：
--   id → id
--   ordernum → order_no
--   uid → user_id
--   productid → product_id (注：zjmf中productid可能为0)
--   amount → amount
--   amount → total (zjmf无total字段，用amount代替)
--   payment → payment_method
--   promo_code → promo_code_id (需查找对应ID)
--   status → status (需映射：Pending→0, Active→1, Cancelled→4, Fraud→7)
--   create_time(int) → created_at(FROM_UNIXTIME)
-- ----------------------------------------------------------------------------
INSERT IGNORE INTO `orders` (
  `id`, `order_no`, `user_id`, `product_id`, `amount`, `total`,
  `currency`, `payment_method`, `status`, `payment_status`,
  `created_at`, `updated_at`
)
SELECT
  `o`.`id`,
  IFNULL(`o`.`ordernum`, CONCAT('ORD-', `o`.`id`)),
  `o`.`uid`,
  IFNULL(`o`.`productid`, 0),
  IFNULL(`o`.`amount`, 0),
  IFNULL(`o`.`amount`, 0),
  'CNY',
  IFNULL(`o`.`payment`, ''),
  CASE
    WHEN `o`.`status` = 'Pending' THEN 0
    WHEN `o`.`status` = 'Active' THEN 1
    WHEN `o`.`status` = 'Cancelled' THEN 4
    WHEN `o`.`status` = 'Fraud' THEN 7
    WHEN `o`.`status` = 'Completed' THEN 3
    ELSE 0
  END,
  CASE
    WHEN `o`.`status` = 'Active' THEN 1
    WHEN `o`.`status` = 'Completed' THEN 1
    ELSE 0
  END,
  FROM_UNIXTIME(IFNULL(`o`.`create_time`, UNIX_TIMESTAMP())),
  NOW()
FROM `shd_orders` `o`
WHERE `o`.`id` IS NOT NULL;

-- ----------------------------------------------------------------------------
-- 2.5 账单迁移：shd_invoices → invoices
-- 字段映射：
--   id → id
--   invoice_num → invoice_no
--   uid → user_id
--   subtotal → sub_total
--   credit → credit
--   tax → tax
--   total → total
--   status → status (需映射：Unpaid→0, Paid→1, Cancelled→3, Refunded→4, Overdue→5)
--   payment → payment_method
--   due_time(int) → due_date(FROM_UNIXTIME)
--   paid_time(int) → paid_at(FROM_UNIXTIME)
--   0 → paid_amount (zjmf无此字段，Paid状态设为total)
-- ----------------------------------------------------------------------------
INSERT IGNORE INTO `invoices` (
  `id`, `invoice_no`, `user_id`, `sub_total`, `credit`, `tax`, `total`,
  `paid_amount`, `status`, `payment_method`, `currency`,
  `due_date`, `paid_at`, `created_at`, `updated_at`
)
SELECT
  `i`.`id`,
  IFNULL(`i`.`invoice_num`, CONCAT('INV-', `i`.`id`)),
  `i`.`uid`,
  IFNULL(`i`.`subtotal`, 0),
  IFNULL(`i`.`credit`, 0),
  IFNULL(`i`.`tax`, 0),
  IFNULL(`i`.`total`, 0),
  CASE
    WHEN `i`.`status` = 'Paid' THEN IFNULL(`i`.`total`, 0)
    ELSE 0
  END,
  CASE
    WHEN `i`.`status` = 'Unpaid' THEN 0
    WHEN `i`.`status` = 'Paid' THEN 1
    WHEN `i`.`status` = 'Cancelled' THEN 3
    WHEN `i`.`status` = 'Refunded' THEN 4
    WHEN `i`.`status` = 'Overdue' THEN 5
    WHEN `i`.`status` = 'Collections' THEN 5
    ELSE 0
  END,
  IFNULL(`i`.`payment`, ''),
  'CNY',
  CASE
    WHEN `i`.`due_time` IS NOT NULL AND `i`.`due_time` > 0
    THEN FROM_UNIXTIME(`i`.`due_time`)
    ELSE NULL
  END,
  CASE
    WHEN `i`.`paid_time` IS NOT NULL AND `i`.`paid_time` > 0
    THEN FROM_UNIXTIME(`i`.`paid_time`)
    ELSE NULL
  END,
  FROM_UNIXTIME(IFNULL(`i`.`due_time`, UNIX_TIMESTAMP())),
  NOW()
FROM `shd_invoices` `i`
WHERE `i`.`id` IS NOT NULL;

-- ----------------------------------------------------------------------------
-- 2.6 账单明细迁移：shd_invoice_items → invoice_items
-- 字段映射：
--   id → id
--   invoice_id → invoice_id
--   type → type
--   relid → rel_id
--   description → description
--   amount → total (zjmf用amount，AF用total)
--   amount → unit_price
--   1 → quantity
-- ----------------------------------------------------------------------------
INSERT IGNORE INTO `invoice_items` (
  `id`, `invoice_id`, `type`, `rel_id`, `description`,
  `quantity`, `unit_price`, `total`, `created_at`, `updated_at`
)
SELECT
  `ii`.`id`,
  `ii`.`invoice_id`,
  IFNULL(`ii`.`type`, 'product'),
  IFNULL(`ii`.`relid`, 0),
  IFNULL(`ii`.`description`, ''),
  1,
  IFNULL(`ii`.`amount`, 0),
  IFNULL(`ii`.`amount`, 0),
  NOW(), NOW()
FROM `shd_invoice_items` `ii`
WHERE `ii`.`id` IS NOT NULL;

-- ----------------------------------------------------------------------------
-- 2.7 主机/服务迁移：shd_host → user_products
-- zjmf的shd_host对应AF的user_products（用户已购服务实例）
-- 字段映射：
--   id → id
--   uid → user_id
--   orderid → order_id
--   productid → product_id
--   domain → domain
--   domainstatus → status (需映射)
--   billingcycle(JSON text) → billing_cycle
--   amount → amount
--   nextduedate → next_due_date (需转换)
--   serverid → (保留到config_options)
--   '' → name (zjmf无name，用domain代替)
-- ----------------------------------------------------------------------------
INSERT IGNORE INTO `user_products` (
  `id`, `user_id`, `order_id`, `product_id`, `name`, `domain`,
  `billing_cycle`, `amount`, `currency`, `status`,
  `next_due_date`, `created_at`, `updated_at`
)
SELECT
  `h`.`id`,
  `h`.`uid`,
  IFNULL(`h`.`orderid`, 0),
  IFNULL(`h`.`productid`, 0),
  IFNULL(`h`.`domain`, CONCAT('service-', `h`.`id`)),
  IFNULL(`h`.`domain`, ''),
  -- billingcycle: zjmf存储为JSON text，提取值
  CASE
    WHEN `h`.`billingcycle` IS NOT NULL AND `h`.`billingcycle` != ''
    THEN `h`.`billingcycle`
    ELSE 'monthly'
  END,
  IFNULL(`h`.`amount`, 0),
  'CNY',
  -- domainstatus映射：Active→1, Suspended→2, Pending→3, Terminated→4, Cancelled→6, Expired→5
  CASE
    WHEN `h`.`domainstatus` = 'Active' THEN 1
    WHEN `h`.`domainstatus` = 'Suspended' THEN 2
    WHEN `h`.`domainstatus` = 'Pending' THEN 3
    WHEN `h`.`domainstatus` = 'Terminated' THEN 4
    WHEN `h`.`domainstatus` = 'Cancelled' THEN 6
    WHEN `h`.`domainstatus` = 'Expired' THEN 5
    WHEN `h`.`domainstatus` = 'Fraud' THEN 4
    ELSE 3
  END,
  CASE
    WHEN `h`.`nextduedate` IS NOT NULL AND `h`.`nextduedate` != '' AND `h`.`nextduedate` != '0000-00-00'
    THEN `h`.`nextduedate`
    ELSE NULL
  END,
  NOW(), NOW()
FROM `shd_host` `h`
WHERE `h`.`id` IS NOT NULL;

-- ----------------------------------------------------------------------------
-- 2.8 工单部门迁移：shd_ticket_department → departments
-- 字段映射：
--   id → id
--   name → name
--   description → description
--   email → email
--   hidden → status (hidden=1→0, hidden=0→1)
--   '' → slug
-- ----------------------------------------------------------------------------
INSERT IGNORE INTO `departments` (
  `id`, `name`, `description`, `email`, `slug`, `status`,
  `created_at`, `updated_at`
)
SELECT
  `id`,
  `name`,
  IFNULL(`description`, ''),
  IFNULL(`email`, ''),
  CONCAT('dept-', `id`),
  CASE WHEN `hidden` = 1 THEN 0 ELSE 1 END,
  NOW(), NOW()
FROM `shd_ticket_department`
WHERE `id` IS NOT NULL;

-- ----------------------------------------------------------------------------
-- 2.9 工单迁移：shd_ticket → tickets
-- 字段映射：
--   id → id
--   tid → ticket_no
--   uid → user_id
--   dptid → department_id
--   title → subject
--   content → (zjmf的content在ticket表中，AF的content在reply中)
--   status → status (需映射：Open→0, Answered→1, Closed→2, Customer-Reply→0)
--   priority → priority (Low→0, Medium→1, High→2, Emergency→3)
--   admin_id → assigned_to
-- ----------------------------------------------------------------------------
INSERT IGNORE INTO `tickets` (
  `id`, `ticket_no`, `user_id`, `department_id`, `subject`,
  `priority`, `status`, `assigned_to`, `source`,
  `created_at`, `updated_at`
)
SELECT
  `t`.`id`,
  IFNULL(`t`.`tid`, CONCAT('TK-', `t`.`id`)),
  `t`.`uid`,
  IFNULL(`t`.`dptid`, 1),
  IFNULL(`t`.`title`, '无标题'),
  CASE
    WHEN `t`.`priority` = 'Low' THEN 0
    WHEN `t`.`priority` = 'Medium' THEN 1
    WHEN `t`.`priority` = 'High' THEN 2
    WHEN `t`.`priority` = 'Emergency' THEN 3
    ELSE 1
  END,
  CASE
    WHEN `t`.`status` = 'Open' THEN 0
    WHEN `t`.`status` = 'Answered' THEN 1
    WHEN `t`.`status` = 'Closed' THEN 2
    WHEN `t`.`status` = 'Customer-Reply' THEN 0
    WHEN `t`.`status` = 'In Progress' THEN 3
    WHEN `t`.`status` = 'Resolved' THEN 4
    ELSE 0
  END,
  CASE
    WHEN `t`.`admin_id` IS NOT NULL AND `t`.`admin_id` > 0
    THEN `t`.`admin_id`
    ELSE NULL
  END,
  'web',
  NOW(), NOW()
FROM `shd_ticket` `t`
WHERE `t`.`id` IS NOT NULL;

-- ----------------------------------------------------------------------------
-- 2.10 工单回复迁移：shd_ticket_reply → ticket_replies
-- 字段映射：
--   id → id
--   tid → ticket_id (注意：zjmf中tid是工单编号，需关联到ticket.id)
--   uid → user_id (非0时)
--   admin_id → admin_id (非0时)
--   content → content
--   attachment → (迁移到attachments表)
-- ----------------------------------------------------------------------------
INSERT IGNORE INTO `ticket_replies` (
  `id`, `ticket_id`, `user_id`, `admin_id`, `content`,
  `is_internal`, `created_at`, `updated_at`
)
SELECT
  `r`.`id`,
  `r`.`tid`,  -- 注意：zjmf的tid就是ticket的id
  CASE
    WHEN `r`.`uid` IS NOT NULL AND `r`.`uid` > 0 THEN `r`.`uid`
    ELSE NULL
  END,
  CASE
    WHEN `r`.`admin_id` IS NOT NULL AND `r`.`admin_id` > 0 THEN `r`.`admin_id`
    ELSE NULL
  END,
  IFNULL(`r`.`content`, ''),
  0,  -- is_internal
  NOW(), NOW()
FROM `shd_ticket_reply` `r`
WHERE `r`.`id` IS NOT NULL;

-- ----------------------------------------------------------------------------
-- 2.11 服务器配置迁移：shd_servers → server_configs
-- 字段映射：
--   id → id
--   gid → (关联到server_config)
--   name → name
--   ip_address → (存入metadata)
--   hostname → (存入metadata)
--   username → (存入metadata)
--   password → (存入metadata，加密)
--   type → type
--   port → (存入metadata)
--   '' → code (zjmf无code)
-- ----------------------------------------------------------------------------
INSERT IGNORE INTO `server_configs` (
  `id`, `name`, `code`, `type`, `provider`,
  `location`, `status`, `metadata`,
  `created_at`, `updated_at`
)
SELECT
  `s`.`id`,
  IFNULL(`s`.`name`, CONCAT('server-', `s`.`id`)),
  CONCAT('srv-', `s`.`id`),
  IFNULL(`s`.`type`, 'vps'),
  IFNULL(`s`.`hostname`, ''),
  IFNULL(`s`.`ip_address`, ''),
  1,  -- status: 默认启用
  JSON_OBJECT(
    'ip_address', IFNULL(`s`.`ip_address`, ''),
    'hostname', IFNULL(`s`.`hostname`, ''),
    'username', IFNULL(`s`.`username`, ''),
    'port', IFNULL(`s`.`port`, 22),
    'gid', IFNULL(`s`.`gid`, 0)
  ),
  NOW(), NOW()
FROM `shd_servers` `s`
WHERE `s`.`id` IS NOT NULL;

-- ----------------------------------------------------------------------------
-- 2.12 系统配置迁移：shd_configuration → system_configs
-- zjmf是简单key-value，AF是group+key+value
-- 映射规则：根据配置key的前缀/含义分配group
-- ----------------------------------------------------------------------------
INSERT IGNORE INTO `system_configs` (`key`, `value`, `group`, `type`, `name`, `created_at`, `updated_at`)
SELECT
  `setting`,
  IFNULL(`value`, ''),
  -- 根据setting名称自动分组
  CASE
    -- 基本设置
    WHEN `setting` IN ('company_name', 'company_url', 'system_url', 'domain',
                        'logo', 'logo_url', 'favicon_url', 'www_logo',
                        'record_no', 'address', 'phonenumber', 'email',
                        'site_name', 'site_url', 'site_description', 'site_keywords',
                        'site_icp', 'site_copyright', 'contact_phone', 'contact_email',
                        'contact_address', 'contact_qq', 'sales_phone', 'support_phone',
                        'sales_email', 'work_time')
    THEN 'general'
    -- 安全设置
    WHEN `setting` LIKE '%captcha%' OR `setting` LIKE '%recaptcha%'
         OR `setting` LIKE '%password%' OR `setting` LIKE '%login_error%'
         OR `setting` LIKE '%ip_check%' OR `setting` LIKE '%banlength%'
         OR `setting` IN ('admin_path')
    THEN 'security'
    -- 登录注册
    WHEN `setting` LIKE '%register%' OR `setting` LIKE '%login%'
         OR `setting` LIKE '%allow_phone%' OR `setting` LIKE '%allow_email%'
    THEN 'login'
    -- 显示配置
    WHEN `setting` IN ('language', 'default_language', 'default_timezone',
                        'date_format', 'per_page_limit', 'nologin_send_ticket')
    THEN 'display'
    -- 模板配置
    WHEN `setting` LIKE '%template%' OR `setting` LIKE '%theme%'
         OR `setting` LIKE '%header%' OR `setting` LIKE '%footer%'
    THEN 'template'
    -- 邮件配置
    WHEN `setting` LIKE '%smtp%' OR `setting` LIKE '%mail%'
         OR `setting` LIKE '%email_%' OR `setting` LIKE '%phpmail%'
    THEN 'email'
    -- 短信配置
    WHEN `setting` LIKE '%sms%' OR `setting` LIKE '%aliyun%'
         OR `setting` LIKE '%tencent%' OR `setting` LIKE '%submail%'
    THEN 'sms'
    -- 支付配置
    WHEN `setting` LIKE '%payment%' OR `setting` LIKE '%alipay%'
         OR `setting` LIKE '%wechat_pay%' OR `setting` LIKE '%epay%'
    THEN 'payment'
    -- 代理配置
    WHEN `setting` LIKE '%affiliate%'
    THEN 'affiliate'
    -- 发票配置
    WHEN `setting` LIKE '%invoice%'
    THEN 'invoice'
    -- 通知配置
    WHEN `setting` LIKE '%notify%' OR `setting` LIKE '%notification%'
    THEN 'notification'
    -- SEO配置
    WHEN `setting` LIKE '%seo%'
    THEN 'seo'
    -- 维护模式
    WHEN `setting` LIKE '%maintenance%' OR `setting` LIKE '%main_tenance%'
    THEN 'maintenance'
    -- 其他
    ELSE 'advanced'
  END,
  'string',  -- type: 默认string
  `setting`,  -- name: 用setting名作为显示名
  NOW(), NOW()
FROM `shd_configuration`
WHERE `setting` IS NOT NULL AND `setting` != '';

-- ----------------------------------------------------------------------------
-- 2.13 支付网关迁移：shd_payment_gateways → payment_gateways
-- zjmf的gateway列对应AF的name，setting是JSON配置
-- ----------------------------------------------------------------------------
INSERT IGNORE INTO `payment_gateways` (
  `name`, `title`, `gateway`, `code`, `config`, `is_enabled`,
  `created_at`, `updated_at`
)
SELECT
  `gateway`,
  `gateway`,  -- title默认同name
  `gateway`,
  `gateway`,
  IFNULL(`setting`, '{}'),
  1,
  NOW(), NOW()
FROM `shd_payment_gateways`
WHERE `gateway` IS NOT NULL AND `gateway` != '';

-- ----------------------------------------------------------------------------
-- 2.14 优惠码迁移：shd_promo_code → promo_codes
-- 字段映射：
--   id → id
--   code → code
--   type → type (Percentage→percent, Fixed→fixed, Free→free)
--   value → value
--   recurring → recurring (1→true, 0→false)
--   maxuses → max_times
-- ----------------------------------------------------------------------------
INSERT IGNORE INTO `promo_codes` (
  `id`, `code`, `type`, `value`, `recurring`, `max_times`,
  `status`, `start_time`, `expiration_time`,
  `created_at`, `updated_at`
)
SELECT
  `id`,
  `code`,
  CASE
    WHEN `type` = 'Percentage' THEN 'percent'
    WHEN `type` = 'Fixed' THEN 'fixed'
    WHEN `type` = 'Free' THEN 'free'
    ELSE IFNULL(`type`, 'percent')
  END,
  IFNULL(`value`, 0),
  IF(`recurring` = 1, 1, 0),
  IFNULL(`maxuses`, 0),
  1,  -- status: 默认启用
  UNIX_TIMESTAMP(),  -- start_time
  UNIX_TIMESTAMP() + 365 * 86400,  -- expiration_time: 默认一年后
  NOW(), NOW()
FROM `shd_promo_code`
WHERE `id` IS NOT NULL;

-- ----------------------------------------------------------------------------
-- 2.15 自定义字段迁移：shd_customfields → custom_fields + custom_field_values
-- zjmf的customfields同时定义字段和值
-- 先迁移字段定义（type+relid=0的为字段定义）
-- ----------------------------------------------------------------------------
INSERT IGNORE INTO `custom_fields` (
  `id`, `name`, `label`, `type`, `group_name`,
  `created_at`, `updated_at`
)
SELECT
  `id`,
  `fieldname`,
  `fieldname`,
  IFNULL(`fieldtype`, 'text'),
  CASE
    WHEN `type` = 'product' THEN 'product'
    WHEN `type` = 'client' THEN 'client'
    WHEN `type` = 'host' THEN 'host'
    WHEN `type` = 'ticket' THEN 'ticket'
    WHEN `type` = 'order' THEN 'order'
    ELSE 'product'
  END,
  NOW(), NOW()
FROM `shd_customfields`
WHERE `fieldname` IS NOT NULL AND `fieldname` != '';

-- 迁移自定义字段值
INSERT IGNORE INTO `custom_field_values` (
  `field_id`, `owner_id`, `owner_type`, `value`,
  `created_at`, `updated_at`
)
SELECT
  `id`,
  IFNULL(`relid`, 0),
  CASE
    WHEN `type` = 'product' THEN 'product'
    WHEN `type` = 'client' THEN 'client'
    WHEN `type` = 'host' THEN 'host'
    WHEN `type` = 'ticket' THEN 'ticket'
    WHEN `type` = 'order' THEN 'order'
    ELSE 'product'
  END,
  '',  -- zjmf自定义字段的值需要从其他表获取
  NOW(), NOW()
FROM `shd_customfields`
WHERE `relid` IS NOT NULL AND `relid` > 0;

-- ----------------------------------------------------------------------------
-- 2.16 上游供应商迁移：shd_zjmf_finance_api → upstream_providers
-- zjmf的上游API配置 → AF的上游供应商表
-- ----------------------------------------------------------------------------
INSERT IGNORE INTO `upstream_providers` (
  `id`, `name`, `type`, `api_url`, `api_key`, `config`,
  `is_active`, `created_at`, `updated_at`
)
SELECT
  `id`,
  IFNULL(`hostname`, CONCAT('upstream-', `id`)),
  'zjmfv3',
  IFNULL(`hostname`, ''),
  IFNULL(`password`, ''),
  JSON_OBJECT(
    'username', IFNULL(`username`, ''),
    'upstream_uid', IFNULL(`upstream_uid`, 0)
  ),
  1,
  NOW(), NOW()
FROM `shd_zjmf_finance_api`
WHERE `id` IS NOT NULL;

-- ----------------------------------------------------------------------------
-- 2.17 系统日志迁移：shd_system_log → system_logs
-- 字段映射：
--   id → id
--   create_time → created_at (FROM_UNIXTIME)
--   uid → user_id
--   type → module
--   description → message
-- ----------------------------------------------------------------------------
INSERT IGNORE INTO `system_logs` (
  `id`, `level`, `module`, `message`, `user_id`, `created_at`
)
SELECT
  `id`,
  'info',
  IFNULL(`type`, 'system'),
  IFNULL(`description`, ''),
  IFNULL(`uid`, 0),
  FROM_UNIXTIME(IFNULL(`create_time`, UNIX_TIMESTAMP()))
FROM `shd_system_log`
WHERE `id` IS NOT NULL;

-- ----------------------------------------------------------------------------
-- 2.18 邮件模板迁移：shd_email_templates → email_templates
-- 字段映射：
--   id → id
--   type → type (保留)
--   name → name
--   subject → subject
--   message → body
--   '' → code (zjmf无code)
-- ----------------------------------------------------------------------------
INSERT IGNORE INTO `email_templates` (
  `id`, `code`, `name`, `subject`, `body`, `type`, `language`,
  `status`, `created_at`, `updated_at`
)
SELECT
  `id`,
  CONCAT('tpl-', `id`),
  IFNULL(`name`, CONCAT('template-', `id`)),
  IFNULL(`subject`, ''),
  IFNULL(`message`, ''),
  IFNULL(`type`, 'email'),
  'zh-CN',
  1,
  NOW(), NOW()
FROM `shd_email_templates`
WHERE `id` IS NOT NULL;

-- ----------------------------------------------------------------------------
-- 2.19 代理/推介迁移：shd_affiliates → agents
-- 字段映射：
--   id → id
--   uid → user_id
--   visitors → (存入metadata，AF无此字段)
--   registcount → (存入metadata)
--   payamount → total_earned
-- ----------------------------------------------------------------------------
INSERT IGNORE INTO `agents` (
  `id`, `user_id`, `agent_no`, `commission_rate`,
  `total_earned`, `status`, `created_at`, `updated_at`
)
SELECT
  `a`.`id`,
  `a`.`uid`,
  CONCAT('AG-', LPAD(`a`.`id`, 6, '0')),
  10.00,  -- 默认佣金比例
  IFNULL(`a`.`payamount`, 0),
  1,
  NOW(), NOW()
FROM `shd_affiliates` `a`
WHERE `a`.`uid` IS NOT NULL AND `a`.`uid` > 0;

-- ----------------------------------------------------------------------------
-- 2.20 下载文件迁移：shd_downloads → download_files
-- 字段映射：
--   id → id
--   category → category_id (需查找或创建分类)
--   title → title
--   filename → file_path
-- ----------------------------------------------------------------------------
-- 先创建默认下载分类
INSERT IGNORE INTO `download_categories` (`id`, `name`, `description`, `slug`, `status`, `created_at`, `updated_at`)
VALUES (1, '默认分类', '系统默认下载分类', 'default', 1, NOW(), NOW());

INSERT IGNORE INTO `download_files` (
  `id`, `category_id`, `title`, `file_path`, `file_size`, `download_count`,
  `is_published`, `created_at`, `updated_at`
)
SELECT
  `d`.`id`,
  IFNULL(`d`.`category`, 1),
  IFNULL(`d`.`title`, CONCAT('file-', `d`.`id`)),
  IFNULL(`d`.`filename`, ''),
  0,  -- file_size: zjmf未记录
  0,
  1,
  NOW(), NOW()
FROM `shd_downloads` `d`
WHERE `d`.`id` IS NOT NULL;

-- ----------------------------------------------------------------------------
-- 2.21 交易流水迁移（如zjmf有accounts表）
-- 注意：如果zjmf使用shd_accounts或其他表名存储交易流水
-- 需根据实际情况调整表名
-- ----------------------------------------------------------------------------


-- ============================================================================
-- 第三部分：创建默认管理员账号
-- ============================================================================

-- 创建默认管理员（密码需首次登录后修改）
-- 默认密码: admin123 (bcrypt hash)
INSERT IGNORE INTO `admins` (
  `id`, `username`, `email`, `password_hash`, `salt`,
  `nickname`, `role`, `is_super`, `status`,
  `created_at`, `updated_at`
) VALUES (
  1, 'admin', 'admin@localhost',
  '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
  '',
  '系统管理员', 'super_admin', 1, 1,
  NOW(), NOW()
);

-- 创建默认语言
INSERT IGNORE INTO `languages` (`code`, `name`, `flag`, `is_default`, `status`, `created_at`, `updated_at`) VALUES
('zh-CN', '中文简体', 'CN', 1, 1, NOW(), NOW()),
('en-US', 'English', 'US', 0, 1, NOW(), NOW()),
('zh-TW', '中文繁體', 'TW', 0, 1, NOW(), NOW());


-- ============================================================================
-- 第四部分：更新自增ID起始值（避免ID冲突）
-- ============================================================================

-- 确保后续新记录的ID不会与迁移数据冲突
-- 获取各表最大ID并更新AUTO_INCREMENT

SET @max_users = (SELECT IFNULL(MAX(`id`), 0) + 100 FROM `users`);
SET @sql = CONCAT('ALTER TABLE `users` AUTO_INCREMENT = ', @max_users);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @max_products = (SELECT IFNULL(MAX(`id`), 0) + 100 FROM `products`);
SET @sql = CONCAT('ALTER TABLE `products` AUTO_INCREMENT = ', @max_products);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @max_orders = (SELECT IFNULL(MAX(`id`), 0) + 100 FROM `orders`);
SET @sql = CONCAT('ALTER TABLE `orders` AUTO_INCREMENT = ', @max_orders);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @max_invoices = (SELECT IFNULL(MAX(`id`), 0) + 100 FROM `invoices`);
SET @sql = CONCAT('ALTER TABLE `invoices` AUTO_INCREMENT = ', @max_invoices);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @max_tickets = (SELECT IFNULL(MAX(`id`), 0) + 100 FROM `tickets`);
SET @sql = CONCAT('ALTER TABLE `tickets` AUTO_INCREMENT = ', @max_tickets);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @max_user_products = (SELECT IFNULL(MAX(`id`), 0) + 100 FROM `user_products`);
SET @sql = CONCAT('ALTER TABLE `user_products` AUTO_INCREMENT = ', @max_user_products);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;


-- ============================================================================
-- 第五部分：数据验证查询
-- ============================================================================

-- 以下查询用于验证迁移结果，不影响数据

-- 验证用户数
SELECT 'users' AS table_name, COUNT(*) AS migrated_count FROM `users`
UNION ALL
SELECT 'products', COUNT(*) FROM `products`
UNION ALL
SELECT 'product_groups', COUNT(*) FROM `product_groups`
UNION ALL
SELECT 'orders', COUNT(*) FROM `orders`
UNION ALL
SELECT 'invoices', COUNT(*) FROM `invoices`
UNION ALL
SELECT 'invoice_items', COUNT(*) FROM `invoice_items`
UNION ALL
SELECT 'user_products', COUNT(*) FROM `user_products`
UNION ALL
SELECT 'tickets', COUNT(*) FROM `tickets`
UNION ALL
SELECT 'ticket_replies', COUNT(*) FROM `ticket_replies`
UNION ALL
SELECT 'departments', COUNT(*) FROM `departments`
UNION ALL
SELECT 'payment_gateways', COUNT(*) FROM `payment_gateways`
UNION ALL
SELECT 'promo_codes', COUNT(*) FROM `promo_codes`
UNION ALL
SELECT 'email_templates', COUNT(*) FROM `email_templates`
UNION ALL
SELECT 'system_configs', COUNT(*) FROM `system_configs`
UNION ALL
SELECT 'system_logs', COUNT(*) FROM `system_logs`
UNION ALL
SELECT 'server_configs', COUNT(*) FROM `server_configs`
UNION ALL
SELECT 'upstream_providers', COUNT(*) FROM `upstream_providers`
UNION ALL
SELECT 'agents', COUNT(*) FROM `agents`
UNION ALL
SELECT 'download_files', COUNT(*) FROM `download_files`
UNION ALL
SELECT 'admins', COUNT(*) FROM `admins`;

SET FOREIGN_KEY_CHECKS = 1;

-- ============================================================================
-- 迁移完成！
-- ============================================================================
-- 后续步骤：
--   1. 检查上面的验证查询结果
--   2. 对于zjmf密码，需要在应用层做兼容处理（读取时识别hash格式）
--   3. 配置表(system_configs)中的值可能需要根据AF逻辑重新映射
--   4. 产品价格(shd_products无price字段)需要从zjmf的定价配置表补充
--   5. 自定义字段值(custom_field_values)的value需要从zjmf相关表补充
--   6. 建议在测试环境先执行一次，验证无误后再在生产环境执行
-- ============================================================================
