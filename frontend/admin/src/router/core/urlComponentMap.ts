/**
 * zjmf URL → Vue component path mapping
 *
 * Generated from init.sql menus table.
 * Used by ComponentLoader to resolve menu URLs to actual Vue components.
 *
 * Component base: @/views/
 */

export const urlComponentMap: Record<string, string> = {
  // ============================================================
  // 1. 客户 (Customer)
  // ============================================================
  '/customer-list': '/finance/clients/list',
  '/customer-authentication': '/finance/clients/authentication',
  '/customer-resources': '/finance/clients/resource-pool',
  '/sales-statistics': '/finance/clients/sales-statistics',
  '/customer-promotionplan': '/finance/marketing/promo-plan',
  '/marketing-push': '/finance/advanced/send-message',

  // ============================================================
  // 2. 业务 (Business)
  // ============================================================
  '/order-list': '/finance/orders/list',
  '/renewal-order': '/finance/orders/renewal',
  '/dcim-traffic-log': '/finance/orders/traffic-orders',
  '/customer-product': '/finance/clients/services',
  '/customer-cancelreq': '/finance/orders/cancel-requests',

  // ============================================================
  // 3. 财务 (Finance)
  // ============================================================
  '/business-statement': '/finance/finance-core/accounts',
  '/bill-management': '/finance/finance-core/accounts',
  '/credit-management': '/finance/finance-core/credit',
  '/customer-withdrawal': '/finance/finance-core/withdraw',
  '/invoice-audit': '/finance/orders/invoice-audit',
  '/contracts_audit': '/finance/advanced/contracts',

  // ============================================================
  // 4. 工单 (Ticket)
  // ============================================================
  '/support-ticket': '/finance/tickets/list',
  '/support-statistics': '/finance/tickets/statistics',

  // ============================================================
  // 5. 功能 (Function)
  // ============================================================
  '/plugins': '/finance/system/plugins',
  '/database-message': '/finance/system/database-info',
  '/statistics-taskQueue': '/finance/system/task-queue',
  '/timing-results': '/finance/system/crons',
  '/annual-statistics': '/finance/statistics/annual',
  '/new-customer': '/finance/statistics/new-customers',
  '/product-revenue': '/finance/statistics/product-revenue',
  '/revenue-ranking': '/finance/statistics/revenue-ranking',

  // ============================================================
  // 6. 资源与商店 (Resources And Stores)
  // ============================================================
  '/app-store': '/finance/marketplace',
  '/api-setup': '/finance/system/api-management',
  '/task-queue': '/finance/system/task-queue',
  '/munual-resource': '/finance/system/servers',
  '/zjmf-api': '/finance/advanced/zjmf-api',
  '/commodity-list': '/finance/products/list',
  '/commodity-product': '/finance/products/groups',
  '/commodity-taskQueue': '/finance/system/task-queue',
  '/supplier-order-list': '/finance/orders/list',
  '/supplier-renewal-order': '/finance/orders/renewal',

  // ============================================================
  // 7. 设置 (Settings)
  // ============================================================

  // 7.1 商品设置
  '/product-server': '/finance/products/servers',
  '/dcim-traffic': '/finance/products/traffic-config',
  '/server-settings': '/finance/system/config-servers',
  '/dcim': '/finance/system/dcim',
  '/zjmfcloud': '/finance/system/dcim-cloud',
  '/configurable-option': '/finance/system/config-options',
  '/order-product': '/finance/products/advanced-options',

  // 7.2 基础设置
  '/work-order-dept': '/finance/tickets/departments',
  '/work-order-status': '/finance/tickets/status',
  '/work-order-rules': '/finance/tickets/deliver',
  '/customer-group': '/finance/clients/groups',
  '/authentication-setting': '/finance/system/authentication-setting',
  '/customer-custom': '/finance/clients/custom-fields',
  '/promotion_plan': '/finance/marketing/promo-plan',
  '/customer-level': '/finance/system/user-levels',
  '/payment-interface': '/finance/system/payment-interface',
  '/promo-code': '/finance/marketing/promo-codes',
  '/currency-settings': '/finance/system/currencies',
  '/general-settings/finance': '/finance/system/fund-config',
  '/voucher-setting': '/finance/system/receipt-config',
  '/credit-setting': '/finance/finance-core/credit',
  '/contracts_setting': '/finance/advanced/contracts',

  // 7.3 站务设置
  '/base-info': '/finance/system/site-settings',
  '/service-support': '/finance/content/downloads',
  '/news-center': '/finance/content/news',
  '/help-center': '/finance/content/help',

  // 7.4 系统设置
  '/general-settings/general': '/finance/system/general',
  '/automatic-tasks': '/finance/system/crons',
  '/login-register': '/finance/system/aggregate-login',
  '/third-login': '/finance/system/third-login',
  '/admin-management': '/finance/system/admins',
  '/permissions-managment': '/finance/system/roles',
  '/sales-management': '/finance/advanced/sales',
  '/sms-template/sms': '/finance/system/sms-templates',
  '/email-list': '/finance/system/email-templates',
  '/sms-template-index': '/finance/system/sms-templates',
  '/sms-send-settings': '/finance/system/sms-templates',
  '/black-list': '/finance/clients/blacklist',
  '/general-settings/captcha': '/finance/system/captcha',
  '/twice-confirm': '/finance/system/twice-confirm',
  '/system-message': '/finance/system/system-messages',
  '/about': '/finance/system/system-info',

  // 7.4 日志
  '/system-log': '/finance/system/logs',
  '/system-admin-log': '/finance/system/login-logs',
  '/email-log': '/finance/system/email-templates', // TODO: no dedicated email log component
  '/sms-log': '/finance/clients/sms-log',
  '/station-letter-log': '/finance/system/messages',
  '/automatic-task-log': '/finance/system/cron-url',
  '/api-log': '/finance/system/operation-logs',
  '/log-cleanup': '/finance/system/log-cleanup',

  // Top-level settings redirect
  '/set': '/finance/dashboard',
}
