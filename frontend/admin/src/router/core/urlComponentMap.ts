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
  '/customer-add': '/finance/clients/add',
  '/add-customer': '/finance/clients/add',
  '/customer-edit': '/finance/clients/add',
  '/edit-customer': '/finance/clients/add',
  '/customer-detail': '/finance/clients/detail/:id',
  '/customer-view': '/finance/clients/detail',
  '/customer-view/abstract': '/finance/clients/detail',
  '/customer-view/person': '/finance/clients/detail',
  '/customer-view/product-list': '/finance/clients/detail',
  '/customer-view/bill': '/finance/clients/detail',
  '/customer-view/transactions': '/finance/clients/detail',
  '/customer-view/credit': '/finance/clients/detail',
  '/customer-view/tickets': '/finance/clients/detail',
  '/customer-view/log': '/finance/clients/detail',
  '/customer-view/api-overview': '/finance/clients/detail',
  '/customer-view/noticelog': '/finance/clients/detail',
  '/customer-view/annex': '/finance/clients/detail',
  '/customer-view/promotion_plan': '/finance/clients/detail',
  '/customer-view/follow-status': '/finance/clients/detail',
  '/customer-authentication': '/finance/clients/authentication',
  '/customer-resources': '/finance/clients/resource-pool',
  '/customer-track': '/finance/clients/track',
  '/customer-contacts': '/finance/clients/contacts',
  '/customer-remarks': '/finance/clients/remarks',
  '/customer-attachments': '/finance/clients/attachments',
  '/customer-binds': '/finance/clients/binds',
  '/sales-statistics': '/finance/clients/sales-statistics',
  '/customer-promotionplan': '/finance/marketing/promo-plan',
  '/marketing-push': '/finance/advanced/send-message',
  '/affiliate-records': '/finance/affiliate/records',

  // ============================================================
  // 2. 业务 (Business)
  // ============================================================
  '/order-list': '/finance/orders/list',
  '/renewal-order': '/finance/orders/renewal',
  '/dcim-traffic-log': '/finance/orders/traffic-orders',
  '/add-order': '/finance/orders/create',
  '/edit-order': '/finance/orders/create',
  '/order-edit': '/finance/orders/create',
  '/multi-renew': '/finance/orders/multi-renew',
  '/order-upgrades': '/finance/orders/upgrades',
  '/order-create': '/finance/orders/create',
  '/order-detail': '/finance/orders/order-detail/:id',
  '/customer-product': '/finance/clients/services',
  '/customer-cancelreq': '/finance/orders/cancel-requests',
  '/service-details': '/finance/system/service-details',

  // ============================================================
  // 3. 财务 (Finance)
  // ============================================================
  '/business-statement': '/finance/finance-core/accounts',
  '/bill-management': '/finance/orders/invoices',
  '/invoice-items': '/finance/orders/invoice-items',
  '/balance-management': '/finance/finance-core/balance',
  '/credit-management': '/finance/finance-core/credit',
  '/customer-withdrawal': '/finance/finance-core/withdraw',
  '/invoice-audit': '/finance/orders/invoice-audit',
  '/vouchers': '/finance/marketing/vouchers',
  '/contracts_audit': '/finance/advanced/contracts',

  // ============================================================
  // 4. 工单 (Ticket)
  // ============================================================
  '/support-ticket': '/finance/tickets/list',
  '/ticket-detail': '/finance/tickets/detail/:id',
  '/add-support-ticket': '/finance/tickets/create',
  '/ticket-prereply': '/finance/tickets/prereply',
  '/knowledge-base': '/finance/tickets/knowledge-base',
  '/ai-auto-reply': '/finance/tickets/ai-auto-reply',
  '/support-statistics': '/finance/tickets/statistics',
  '/password-config': '/finance/tickets/password-config',

  // ============================================================
  // 5. 功能 (Function)
  // ============================================================
  '/plugins': '/finance/system/plugins',
  '/login-register': '/finance/system/aggregate-login',
  '/statistics-overview': '/finance/statistics/overview',
  '/annual-statistics': '/finance/statistics/annual',
  '/new-customer': '/finance/statistics/new-customers',
  '/product-revenue': '/finance/statistics/product-revenue',
  '/revenue-ranking': '/finance/statistics/revenue-ranking',
  '/timing-results': '/finance/system/crons',
  '/statistics-taskQueue': '/finance/system/task-queue',
  '/database-message': '/finance/system/database-info',
  '/about': '/finance/system/system-info',
  '/dashboard': '/finance/dashboard',

  // ============================================================
  // 6. 资源与商店 (Resources And Stores)
  // ============================================================
  '/app-store': '/finance/marketplace',
  '/api-setup': '/finance/system/api-management',
  '/task-queue': '/finance/system/task-queue',
  '/munual-resource': '/finance/system/servers',
  '/zjmf-api': '/finance/advanced/zjmf-api',
  '/config-edit': '/finance/advanced/zjmf-api',
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
  '/work-order-dept-add': '/finance/tickets/departments',
  '/work-order-status': '/finance/tickets/status',
  '/work-order-rules': '/finance/tickets/deliver',
  '/customer-group': '/finance/clients/groups',
  '/authentication-setting': '/finance/system/authentication-setting',
  '/customer-custom': '/finance/clients/custom-fields',
  '/promotion_plan': '/finance/marketing/promo-plan',
  '/customer-level': '/finance/system/user-levels',
  '/payment-interface': '/finance/system/payment-interface',
  '/promo-code': '/finance/system/promo-code',
  '/currency-settings': '/finance/system/currency-settings',
  '/general-settings/finance': '/finance/system/fund-config',
  '/voucher-setting': '/finance/system/receipt-config',
  '/credit-setting': '/finance/finance-core/credit',
  '/contracts_setting': '/finance/advanced/contracts',

  // 7.3 站务设置
  '/base-info': '/finance/system/site-settings',
  '/theme-template': '/finance/system/theme-template',
  '/custom-template-fields': '/finance/system/custom-fields',
  '/menu_manage': '/finance/system/navigation',
  '/friendly_link': '/finance/content/friendly-links',
  '/service-support': '/finance/content/downloads',
  '/download-group': '/finance/content/download-groups',
  '/news-center': '/finance/content/news',
  '/news-list': '/finance/content/news-list',
  '/news-category': '/finance/content/news-category',
  '/help-center': '/finance/content/help',
  '/help-list': '/finance/content/help-list',
  '/help-category': '/finance/content/help-category',

  // 7.4 系统设置
  '/general-settings': '/finance/system/general',
  '/general-settings/general': '/finance/system/general',
  '/automatic-tasks': '/finance/system/crons',
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

  // 7.5 日志
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
