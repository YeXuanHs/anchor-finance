import { AppRouteRecord } from '@/types/router'

/**
 * 财务系统路由模块
 */
export const financeRoutes: AppRouteRecord = {
  path: '/finance',
  name: 'Finance',
  redirect: '/finance/dashboard',
  meta: {
    title: '财务管理',
    icon: 'ep:money',
    order: 1
  },
  children: [
    {
      path: 'dashboard',
      name: 'FinanceDashboard',
      component: () => import('@/views/finance/dashboard/index.vue'),
      meta: {
        title: '财务概览',
        icon: 'ep:data-analysis',
        affix: true
      }
    },
    {
      path: 'clients',
      name: 'Clients',
      redirect: '/finance/clients/list',
      meta: {
        title: '客户管理',
        icon: 'ep:user'
      },
      children: [
        {
          path: 'list',
          name: 'ClientList',
          component: () => import('@/views/finance/clients/list/index.vue'),
          meta: {
            title: '客户列表',
            icon: 'ep:list'
          }
        },
        {
          path: 'groups',
          name: 'ClientGroups',
          component: () => import('@/views/finance/clients/groups/index.vue'),
          meta: {
            title: '客户分组',
            icon: 'ep:folder'
          }
        },
        {
          path: 'track',
          name: 'ClientTrack',
          component: () => import('@/views/finance/clients/track/index.vue'),
          meta: {
            title: '客户跟踪',
            icon: 'ep:position'
          }
        },
        {
          path: 'contacts',
          name: 'ClientContacts',
          component: () => import('@/views/finance/clients/contacts/index.vue'),
          meta: {
            title: '联系人',
            icon: 'ep:phone'
          }
        },
        {
          path: 'services',
          name: 'ClientServices',
          component: () => import('@/views/finance/clients/services/index.vue'),
          meta: {
            title: '客户服务',
            icon: 'ep:service'
          }
        },
        {
          path: 'remarks',
          name: 'ClientRemarks',
          component: () => import('@/views/finance/clients/remarks/index.vue'),
          meta: {
            title: '用户备注',
            icon: 'ep:edit'
          }
        },
        {
          path: 'authentication',
          name: 'ClientAuthentication',
          component: () => import('@/views/finance/clients/authentication/index.vue'),
          meta: {
            title: '实名认证',
            icon: 'ep:stamp'
          }
        },
        {
          path: 'resources',
          name: 'ClientResources',
          component: () => import('@/views/finance/clients/resources/index.vue'),
          meta: {
            title: '资源管理',
            icon: 'ep:box'
          }
        },
        {
          path: 'resource-pool',
          name: 'ResourcePool',
          component: () => import('@/views/finance/clients/resource-pool/index.vue'),
          meta: {
            title: '资源池管理',
            icon: 'ep:coin'
          }
        },
        {
          path: 'sales-statistics',
          name: 'SalesStatistics',
          component: () => import('@/views/finance/clients/sales-statistics/index.vue'),
          meta: {
            title: '销售统计',
            icon: 'ep:data-analysis'
          }
        },
        {
          path: 'custom-fields',
          name: 'CustomFields',
          component: () => import('@/views/finance/clients/custom-fields/index.vue'),
          meta: {
            title: '自定义字段',
            icon: 'ep:edit'
          }
        },
        {
          path: 'affiliate-config',
          name: 'AffiliateConfig',
          component: () => import('@/views/finance/clients/affiliate-config/index.vue'),
          meta: {
            title: '推介计划配置',
            icon: 'ep:connection'
          }
        },
        {
          path: 'notification-log',
          name: 'NotificationLog',
          component: () => import('@/views/finance/clients/notification-log/index.vue'),
          meta: {
            title: '通知日志',
            icon: 'ep:bell'
          }
        },
        {
          path: 'attachments',
          name: 'ClientAttachments',
          component: () => import('@/views/finance/clients/attachments/index.vue'),
          meta: {
            title: '附件管理',
            icon: 'ep:paperclip'
          }
        },
        {
          path: 'crm',
          name: 'ClientCRM',
          component: () => import('@/views/finance/clients/crm/index.vue'),
          meta: {
            title: '客户CRM',
            icon: 'ep:user'
          }
        },
        {
          path: 'email-view',
          name: 'EmailView',
          component: () => import('@/views/finance/clients/email-view/index.vue'),
          meta: {
            title: '邮件查看',
            icon: 'ep:message'
          }
        },
        {
          path: 'add',
          name: 'ClientAdd',
          component: () => import('@/views/finance/clients/add/index.vue'),
          meta: {
            title: '添加客户',
            icon: 'ep:plus',
            isHide: true
          }
        },
        {
          path: 'blacklist',
          name: 'ClientBlacklist',
          component: () => import('@/views/finance/clients/blacklist/index.vue'),
          meta: {
            title: '黑名单管理',
            icon: 'ep:circle-close-filled'
          }
        },
        {
          path: 'sms-log',
          name: 'SmsLog',
          component: () => import('@/views/finance/clients/sms-log/index.vue'),
          meta: {
            title: '短信日志',
            icon: 'ep:message'
          }
        },
        {
          path: 'binds',
          name: 'ClientBinds',
          component: () => import('@/views/finance/clients/binds/index.vue'),
          meta: {
            title: '账号绑定',
            icon: 'ep:connection'
          }
        },
        {
          path: 'detail/:id',
          name: 'ClientDetail',
          component: () => import('@/views/finance/clients/detail/index.vue'),
          meta: {
            title: '客户详情',
            icon: 'ep:user',
            isHide: true
          }
        },
        {
          path: 'custom-field-groups',
          name: 'CustomFieldGroups',
          component: () => import('@/views/finance/clients/custom-field-groups/index.vue'),
          meta: {
            title: '自定义字段组',
            icon: 'ep:files'
          }
        }
      ]
    },
    {
      path: 'products',
      name: 'Products',
      redirect: '/finance/products/list',
      meta: {
        title: '产品管理',
        icon: 'ep:goods'
      },
      children: [
        {
          path: 'list',
          name: 'ProductList',
          component: () => import('@/views/finance/products/list/index.vue'),
          meta: {
            title: '产品列表',
            icon: 'ep:list'
          }
        },
        {
          path: 'groups',
          name: 'ProductGroups',
          component: () => import('@/views/finance/products/groups/index.vue'),
          meta: {
            title: '产品分组',
            icon: 'ep:folder'
          }
        },
        {
          path: 'pricing',
          name: 'ProductPricing',
          component: () => import('@/views/finance/products/pricing/index.vue'),
          meta: {
            title: '定价管理',
            icon: 'ep:price-tag'
          }
        },
        {
          path: 'ai-shopping',
          name: 'AIShopping',
          component: () => import('@/views/finance/products/ai-shopping/index.vue'),
          meta: {
            title: 'AI购物助手',
            icon: 'ep:shopping-cart-full'
          }
        },
        {
          path: 'advanced-options',
          name: 'AdvancedOptions',
          component: () => import('@/views/finance/products/advanced-options/index.vue'),
          meta: {
            title: '高级配置项',
            icon: 'ep:setting'
          }
        },
        {
          path: 'traffic-config',
          name: 'TrafficConfig',
          component: () => import('@/views/finance/products/traffic-config/index.vue'),
          meta: {
            title: '流量包配置',
            icon: 'ep:connection'
          }
        },
        {
          path: 'product-diverts',
          name: 'ProductDiverts',
          component: () => import('@/views/finance/products/product-diverts/index.vue'),
          meta: {
            title: '产品转移/分流',
            icon: 'ep:sort'
          }
        },
        {
          path: 'product-transfers',
          name: 'ProductTransfers',
          component: () => import('@/views/finance/products/product-transfers/index.vue'),
          meta: {
            title: '产品转移配置',
            icon: 'ep:switch'
          }
        },
        {
          path: 'first-groups',
          name: 'ProductFirstGroups',
          component: () => import('@/views/finance/products/first-groups/index.vue'),
          meta: {
            title: '一级产品分组',
            icon: 'ep:folder-opened'
          }
        },
        {
          path: 'config-groups',
          name: 'ProductConfigGroups',
          component: () => import('@/views/finance/products/config-groups/index.vue'),
          meta: {
            title: '产品配置组',
            icon: 'ep:setting'
          }
        },
        {
          path: 'interflows',
          name: 'Interflows',
          component: () => import('@/views/finance/products/interflows/index.vue'),
          meta: {
            title: '互通管理',
            icon: 'ep:connection'
          }
        },
        {
          path: 'downloads',
          name: 'ProductDownloads',
          component: () => import('@/views/finance/products/downloads/index.vue'),
          meta: {
            title: '产品下载管理',
            icon: 'ep:download'
          }
        }
      ]
    },
    {
      path: 'orders',
      name: 'Orders',
      redirect: '/finance/orders/list',
      meta: {
        title: '订单管理',
        icon: 'ep:document'
      },
      children: [
        {
          path: 'list',
          name: 'OrderList',
          component: () => import('@/views/finance/orders/list/index.vue'),
          meta: {
            title: '订单列表',
            icon: 'ep:list'
          }
        },
        {
          path: 'invoices',
          name: 'Invoices',
          component: () => import('@/views/finance/orders/invoices/index.vue'),
          meta: {
            title: '账单管理',
            icon: 'ep:tickets'
          }
        },
        {
          path: 'invoice-items',
          name: 'InvoiceItems',
          component: () => import('@/views/finance/orders/invoice-items/index.vue'),
          meta: {
            title: '账单明细',
            icon: 'ep:list'
          }
        },
        {
          path: 'cancel-requests',
          name: 'CancelRequests',
          component: () => import('@/views/finance/orders/cancel-requests/index.vue'),
          meta: {
            title: '取消请求',
            icon: 'ep:warning-filled'
          }
        },
        {
          path: 'traffic-orders',
          name: 'TrafficOrders',
          component: () => import('@/views/finance/orders/traffic-orders/index.vue'),
          meta: {
            title: '流量包订单',
            icon: 'ep:connection'
          }
        },
        {
          path: 'renewal',
          name: 'RenewalOrders',
          component: () => import('@/views/finance/orders/renewal/index.vue'),
          meta: {
            title: '续费订单',
            icon: 'ep:refresh'
          }
        },
        {
          path: 'invoice-audit',
          name: 'InvoiceAudit',
          component: () => import('@/views/finance/orders/invoice-audit/index.vue'),
          meta: {
            title: '发票审核',
            icon: 'ep:document-checked'
          }
        },
        {
          path: 'invoice-detail/:id',
          name: 'InvoiceDetail',
          component: () => import('@/views/finance/orders/invoice-detail/index.vue'),
          meta: {
            title: '账单详情',
            icon: 'ep:document',
            isHide: true
          }
        },
        {
          path: 'order-detail/:id',
          name: 'OrderDetail',
          component: () => import('@/views/finance/orders/order-detail/index.vue'),
          meta: {
            title: '订单详情',
            icon: 'ep:document',
            isHide: true
          }
        },
        {
          path: 'create',
          name: 'CreateOrder',
          component: () => import('@/views/finance/orders/create/index.vue'),
          meta: {
            title: '创建订单',
            icon: 'ep:plus',
            isHide: true
          }
        },
        {
          path: 'multi-renew',
          name: 'MultiRenew',
          component: () => import('@/views/finance/orders/multi-renew/index.vue'),
          meta: {
            title: '批量续费',
            icon: 'ep:refresh'
          }
        },
        {
          path: 'upgrades',
          name: 'Upgrades',
          component: () => import('@/views/finance/orders/upgrades/index.vue'),
          meta: {
            title: '升级管理',
            icon: 'ep:top'
          }
        }
      ]
    },
    {
      path: 'tickets',
      name: 'Tickets',
      redirect: '/finance/tickets/list',
      meta: {
        title: '工单系统',
        icon: 'ep:chat-dot-round'
      },
      children: [
        {
          path: 'list',
          name: 'TicketList',
          component: () => import('@/views/finance/tickets/list/index.vue'),
          meta: {
            title: '工单列表',
            icon: 'ep:list'
          }
        },
        {
          path: 'departments',
          name: 'TicketDepartments',
          component: () => import('@/views/finance/tickets/departments/index.vue'),
          meta: {
            title: '部门管理',
            icon: 'ep:office-building'
          }
        },
        {
          path: 'prereply',
          name: 'TicketPrereply',
          component: () => import('@/views/finance/tickets/prereply/index.vue'),
          meta: {
            title: '预回复',
            icon: 'ep:chat-line-round'
          }
        },
        {
          path: 'deliver',
          name: 'TicketDeliver',
          component: () => import('@/views/finance/tickets/deliver/index.vue'),
          meta: {
            title: '工单传递',
            icon: 'ep:promotion'
          }
        },
        {
          path: 'status',
          name: 'TicketStatus',
          component: () => import('@/views/finance/tickets/status/index.vue'),
          meta: {
            title: '工单状态',
            icon: 'ep:flag'
          }
        },
        {
          path: 'knowledge-base',
          name: 'KnowledgeBase',
          component: () => import('@/views/finance/tickets/knowledge-base/index.vue'),
          meta: {
            title: '知识库',
            icon: 'ep:collection'
          }
        },
        {
          path: 'ai-auto-reply',
          name: 'AIAutoReply',
          component: () => import('@/views/finance/tickets/ai-auto-reply/index.vue'),
          meta: {
            title: 'AI自动回复',
            icon: 'ep:chat-line-round'
          }
        },
        {
          path: 'statistics',
          name: 'TicketStatistics',
          component: () => import('@/views/finance/tickets/statistics/index.vue'),
          meta: {
            title: '工单统计',
            icon: 'ep:data-analysis'
          }
        },
        {
          path: 'detail/:id',
          name: 'TicketDetail',
          component: () => import('@/views/finance/tickets/detail/index.vue'),
          meta: {
            title: '工单详情',
            icon: 'ep:document',
            isHide: true
          }
        },
        {
          path: 'password-config',
          name: 'TicketPasswordConfig',
          component: () => import('@/views/finance/tickets/password-config/index.vue'),
          meta: {
            title: '密码配置',
            icon: 'ep:key'
          }
        }
      ]
    },
    {
      path: 'finance-core',
      name: 'FinanceCore',
      redirect: '/finance/finance-core/balance',
      meta: {
        title: '资金管理',
        icon: 'ep:wallet'
      },
      children: [
        {
          path: 'balance',
          name: 'Balance',
          component: () => import('@/views/finance/finance-core/balance/index.vue'),
          meta: {
            title: '余额管理',
            icon: 'ep:coin'
          }
        },
        {
          path: 'credit',
          name: 'Credit',
          component: () => import('@/views/finance/finance-core/credit/index.vue'),
          meta: {
            title: '信用额度',
            icon: 'ep:credit-card'
          }
        },
        {
          path: 'withdraw',
          name: 'Withdraw',
          component: () => import('@/views/finance/finance-core/withdraw/index.vue'),
          meta: {
            title: '提现管理',
            icon: 'ep:download'
          }
        },
        {
          path: 'accounts',
          name: 'Accounts',
          component: () => import('@/views/finance/finance-core/accounts/index.vue'),
          meta: {
            title: '交易流水',
            icon: 'ep:money'
          }
        }
      ]
    },
    {
      path: 'affiliate',
      name: 'Affiliate',
      redirect: '/finance/affiliate/list',
      meta: {
        title: '推介计划',
        icon: 'ep:connection'
      },
      children: [
        {
          path: 'list',
          name: 'AffiliateList',
          component: () => import('@/views/finance/affiliate/list/index.vue'),
          meta: {
            title: '推介列表',
            icon: 'ep:list'
          }
        },
        {
          path: 'records',
          name: 'AffiliateRecords',
          component: () => import('@/views/finance/affiliate/records/index.vue'),
          meta: {
            title: '佣金记录',
            icon: 'ep:list'
          }
        }
      ]
    },
    {
      path: 'marketing',
      name: 'Marketing',
      redirect: '/finance/marketing/promo-codes',
      meta: {
        title: '营销推广',
        icon: 'ep:present'
      },
      children: [
        {
          path: 'coupons',
          name: 'Coupons',
          redirect: '/finance/marketing/promo-codes',
          meta: {
            title: '优惠券(已废弃)',
            icon: 'ep:ticket'
          }
        },
        {
          path: 'vouchers',
          name: 'Vouchers',
          component: () => import('@/views/finance/marketing/vouchers/index.vue'),
          meta: {
            title: '发票管理',
            icon: 'ep:coin'
          }
        },
        {
          path: 'promo-codes',
          name: 'PromoCodes',
          component: () => import('@/views/finance/marketing/promo-codes/index.vue'),
          meta: {
            title: '优惠码',
            icon: 'ep:ticket'
          }
        }
      ]
    },
    {
      path: 'content',
      name: 'Content',
      redirect: '/finance/content/announcements',
      meta: {
        title: '内容管理',
        icon: 'ep:document'
      },
      children: [
        {
          path: 'announcements',
          name: 'Announcements',
          component: () => import('@/views/finance/content/announcements/index.vue'),
          meta: {
            title: '公告管理',
            icon: 'ep:bell'
          }
        },
        {
          path: 'news',
          name: 'News',
          component: () => import('@/views/finance/content/news/index.vue'),
          meta: {
            title: '新闻管理',
            icon: 'ep:reading'
          }
        },
        {
          path: 'knowledge',
          name: 'Knowledge',
          component: () => import('@/views/finance/content/knowledge/index.vue'),
          meta: {
            title: '知识库',
            icon: 'ep:collection'
          }
        },
        {
          path: 'downloads',
          name: 'Downloads',
          component: () => import('@/views/finance/content/downloads/index.vue'),
          meta: {
            title: '下载中心',
            icon: 'ep:download'
          }
        },
        {
          path: 'community',
          name: 'Community',
          component: () => import('@/views/finance/content/community/index.vue'),
          meta: {
            title: '社区管理',
            icon: 'ep:chat-dot-round'
          }
        },
        {
          path: 'help',
          name: 'HelpCenter',
          component: () => import('@/views/finance/content/help/index.vue'),
          meta: {
            title: '帮助中心',
            icon: 'ep:help'
          }
        },
        {
          path: 'link-cause',
          name: 'LinkCause',
          component: () => import('@/views/finance/content/link-cause/index.vue'),
          meta: {
            title: '关联原因',
            icon: 'ep:connection'
          }
        },
        {
          path: 'link-knowledge',
          name: 'LinkKnowledge',
          component: () => import('@/views/finance/content/link-knowledge/index.vue'),
          meta: {
            title: '关联知识库',
            icon: 'ep:collection'
          }
        },
        {
          path: 'friendly-links',
          name: 'FriendlyLinks',
          component: () => import('@/views/finance/content/friendly-links/index.vue'),
          meta: {
            title: '友情链接',
            icon: 'ep:link'
          }
        }
      ]
    },
    {
      path: 'system',
      name: 'System',
      redirect: '/finance/system/general',
      meta: {
        title: '系统设置',
        icon: 'ep:setting'
      },
      children: [
        {
          path: 'general',
          name: 'GeneralConfig',
          component: () => import('@/views/finance/system/general/index.vue'),
          meta: {
            title: '常规设置',
            icon: 'ep:tools'
          }
        },
        {
          path: 'email-templates',
          name: 'EmailTemplates',
          component: () => import('@/views/finance/system/email-templates/index.vue'),
          meta: {
            title: '邮件模板',
            icon: 'ep:message'
          }
        },
        {
          path: 'admins',
          name: 'Admins',
          component: () => import('@/views/finance/system/admins/index.vue'),
          meta: {
            title: '管理员管理',
            icon: 'ep:user-filled'
          }
        },
        {
          path: 'roles',
          name: 'Roles',
          component: () => import('@/views/finance/system/roles/index.vue'),
          meta: {
            title: '角色权限',
            icon: 'ep:lock'
          }
        },
        {
          path: 'currencies',
          name: 'Currencies',
          component: () => import('@/views/finance/system/currencies/index.vue'),
          meta: {
            title: '货币管理',
            icon: 'ep:coin'
          }
        },
        {
          path: 'crons',
          name: 'Crons',
          component: () => import('@/views/finance/system/crons/index.vue'),
          meta: {
            title: '定时任务',
            icon: 'ep:clock'
          }
        },
        {
          path: 'payment-gateways',
          name: 'PaymentGateways',
          component: () => import('@/views/finance/system/payment-gateways/index.vue'),
          meta: {
            title: '支付管理',
            icon: 'ep:credit-card'
          }
        },
        {
          path: 'plugins',
          name: 'Plugins',
          component: () => import('@/views/finance/system/plugins/index.vue'),
          meta: {
            title: '插件管理',
            icon: 'ep:connection'
          }
        },
        {
          path: 'oauth-providers',
          name: 'OAuthProviders',
          component: () => import('@/views/finance/system/oauth-providers/index.vue'),
          meta: {
            title: 'OAuth登录',
            icon: 'ep:platform'
          }
        },
        {
          path: 'logs',
          name: 'Logs',
          component: () => import('@/views/finance/system/logs/index.vue'),
          meta: {
            title: '系统日志',
            icon: 'ep:list'
          }
        },
        {
          path: 'operation-logs',
          name: 'OperationLogs',
          component: () => import('@/views/finance/system/operation-logs/index.vue'),
          meta: {
            title: '操作日志',
            icon: 'ep:document'
          }
        },
        {
          path: 'user-levels',
          name: 'UserLevels',
          component: () => import('@/views/finance/system/user-levels/index.vue'),
          meta: {
            title: '用户等级',
            icon: 'ep:medal'
          }
        },
        {
          path: 'servers',
          name: 'Servers',
          component: () => import('@/views/finance/system/servers/index.vue'),
          meta: {
            title: '供给配置',
            icon: 'ep:monitor'
          }
        },
        {
          path: 'config-servers',
          name: 'ConfigServers',
          component: () => import('@/views/finance/system/config-servers/index.vue'),
          meta: {
            title: '服务器管理',
            icon: 'ep:monitor'
          }
        },
        {
          path: 'config-options',
          name: 'ConfigOptions',
          component: () => import('@/views/finance/system/config-options/index.vue'),
          meta: {
            title: '配置选项',
            icon: 'ep:setting'
          }
        },
        {
          path: 'hosts',
          name: 'Hosts',
          component: () => import('@/views/finance/system/hosts/index.vue'),
          meta: {
            title: '主机管理',
            icon: 'ep:cpu'
          }
        },
        {
          path: 'dcim',
          name: 'DCIM',
          component: () => import('@/views/finance/system/dcim/index.vue'),
          meta: {
            title: 'DCIM管理',
            icon: 'ep:platform'
          }
        },
        {
          path: 'dcim-cloud',
          name: 'DcimCloud',
          component: () => import('@/views/finance/system/dcim-cloud/index.vue'),
          meta: {
            title: '魔方云对接',
            icon: 'ep:cloudy'
          }
        },
        {
          path: 'rules',
          name: 'Rules',
          component: () => import('@/views/finance/system/rules/index.vue'),
          meta: {
            title: '规则管理',
            icon: 'ep:document-checked'
          }
        },
        {
          path: 'uploads',
          name: 'Uploads',
          component: () => import('@/views/finance/system/uploads/index.vue'),
          meta: {
            title: '上传管理',
            icon: 'ep:upload-filled'
          }
        },
        {
          path: 'user-tastes',
          name: 'UserTastes',
          component: () => import('@/views/finance/system/user-tastes/index.vue'),
          meta: {
            title: '用户偏好',
            icon: 'ep:setting'
          }
        },
        {
          path: 'maintenance',
          name: 'Maintenance',
          component: () => import('@/views/finance/system/maintenance/index.vue'),
          meta: {
            title: '维护模式',
            icon: 'ep:warning'
          }
        },
        {
          path: 'email-suffixes',
          name: 'EmailSuffixes',
          component: () => import('@/views/finance/system/email-suffixes/index.vue'),
          meta: {
            title: '邮箱后缀白名单',
            icon: 'ep:message'
          }
        },
        {
          path: 'cs-chat',
          name: 'CSChat',
          component: () => import('@/views/finance/system/cs-chat/index.vue'),
          meta: {
            title: '客服聊天系统',
            icon: 'ep:service'
          }
        },
        {
          path: 'marketplace',
          name: 'MarketplaceConfig',
          component: () => import('@/views/finance/marketplace/index.vue'),
          meta: {
            title: '交易市场',
            icon: 'ep:shop'
          }
        },
        {
          path: 'menus',
          name: 'Menus',
          component: () => import('@/views/finance/system/menus/index.vue'),
          meta: {
            title: '菜单管理',
            icon: 'ep:menu'
          }
        },
        {
          path: 'messages',
          name: 'Messages',
          component: () => import('@/views/finance/system/messages/index.vue'),
          meta: {
            title: '消息配置',
            icon: 'ep:message'
          }
        },
        {
          path: 'site-settings',
          name: 'SiteSettings',
          component: () => import('@/views/finance/system/site-settings/index.vue'),
          meta: {
            title: '站点设置',
            icon: 'ep:house'
          }
        },
        {
          path: 'languages',
          name: 'Languages',
          component: () => import('@/views/finance/system/languages/index.vue'),
          meta: {
            title: '语言管理',
            icon: 'ep:connection'
          }
        },
        {
          path: 'cron-url',
          name: 'CronUrl',
          component: () => import('@/views/finance/system/cron-url/index.vue'),
          meta: {
            title: 'URL定时任务',
            icon: 'ep:link'
          }
        },
        {
          path: 'provisions',
          name: 'Provisions',
          component: () => import('@/views/finance/system/provisions/index.vue'),
          meta: {
            title: '供应管理',
            icon: 'ep:box'
          }
        },
        {
          path: 'log-records',
          name: 'LogRecords',
          component: () => import('@/views/finance/system/log-records/index.vue'),
          meta: {
            title: '日志记录',
            icon: 'ep:document'
          }
        },
        {
          path: 'run-map',
          name: 'RunMap',
          component: () => import('@/views/finance/system/run-map/index.vue'),
          meta: {
            title: '运行映射',
            icon: 'ep:map-location'
          }
        },
        {
          path: 'system-info',
          name: 'SystemInfo',
          component: () => import('@/views/finance/system/system-info/index.vue'),
          meta: {
            title: '系统信息',
            icon: 'ep:info-filled'
          }
        },
        {
          path: 'database-info',
          name: 'DatabaseInfo',
          component: () => import('@/views/finance/system/database-info/index.vue'),
          meta: {
            title: '数据库信息',
            icon: 'ep:coin'
          }
        },
        {
          path: 'api-management',
          name: 'ApiManagement',
          component: () => import('@/views/finance/system/api-management/index.vue'),
          meta: {
            title: 'API管理',
            icon: 'ep:key'
          }
        },
        {
          path: 'sms-templates',
          name: 'SmsTemplates',
          component: () => import('@/views/finance/system/sms-templates/index.vue'),
          meta: {
            title: '短信模板',
            icon: 'ep:message'
          }
        },
        {
          path: 'security-config',
          name: 'SecurityConfig',
          component: () => import('@/views/finance/system/security-config/index.vue'),
          meta: {
            title: '安全配置',
            icon: 'ep:shield'
          }
        },
        {
          path: 'fund-config',
          name: 'FundConfig',
          component: () => import('@/views/finance/system/fund-config/index.vue'),
          meta: {
            title: '资金设置',
            icon: 'ep:wallet'
          }
        },
        {
          path: 'receipt-config',
          name: 'ReceiptConfig',
          component: () => import('@/views/finance/system/receipt-config/index.vue'),
          meta: {
            title: '发票设置',
            icon: 'ep:tickets'
          }
        },
        {
          path: 'captcha',
          name: 'CaptchaConfig',
          component: () => import('@/views/finance/system/captcha/index.vue'),
          meta: {
            title: '验证码配置',
            icon: 'ep:lock'
          }
        },
        {
          path: 'log-cleanup',
          name: 'LogCleanup',
          component: () => import('@/views/finance/system/log-cleanup/index.vue'),
          meta: {
            title: '日志清理',
            icon: 'ep:delete'
          }
        },
        {
          path: 'data-migration',
          name: 'DataMigration',
          component: () => import('@/views/finance/system/data-migration/index.vue'),
          meta: {
            title: '数据迁移',
            icon: 'ep:upload'
          }
        },
        {
          path: 'twice-confirm',
          name: 'TwiceConfirm',
          component: () => import('@/views/finance/system/twice-confirm/index.vue'),
          meta: {
            title: '二次确认',
            icon: 'ep:circle-check'
          }
        },
        {
          path: 'task-queue',
          name: 'TaskQueue',
          component: () => import('@/views/finance/system/task-queue/index.vue'),
          meta: {
            title: '任务队列',
            icon: 'ep:list'
          }
        },
        {
          path: 'service-details',
          name: 'ServiceDetails',
          component: () => import('@/views/finance/system/service-details/index.vue'),
          meta: {
            title: '服务详情管理',
            icon: 'ep:service'
          }
        },
        {
          path: 'login-logs',
          name: 'LoginLogs',
          component: () => import('@/views/finance/system/login-logs/index.vue'),
          meta: {
            title: '登录日志',
            icon: 'ep:document'
          }
        },
        {
          path: 'notification-templates',
          name: 'NotificationTemplates',
          component: () => import('@/views/finance/system/notification-templates/index.vue'),
          meta: {
            title: '通知模板',
            icon: 'ep:bell'
          }
        },
        {
          path: 'system-messages',
          name: 'SystemMessages',
          component: () => import('@/views/finance/system/system-messages/index.vue'),
          meta: {
            title: '系统消息',
            icon: 'ep:message'
          }
        },
        {
          path: 'rbac-pages',
          name: 'RbacPages',
          component: () => import('@/views/finance/system/rbac-pages/index.vue'),
          meta: {
            title: '页面权限',
            icon: 'ep:lock'
          }
        },
        {
          path: 'nav-groups',
          name: 'NavGroups',
          component: () => import('@/views/finance/system/nav-groups/index.vue'),
          meta: {
            title: '导航分组',
            icon: 'ep:menu'
          }
        },
        {
          path: 'rule-middle',
          name: 'RuleMiddle',
          component: () => import('@/views/finance/system/rule-middle/index.vue'),
          meta: {
            title: '规则中间件',
            icon: 'ep:operation'
          }
        },
        {
          path: 'backup',
          name: 'Backup',
          component: () => import('@/views/finance/system/backup/index.vue'),
          meta: {
            title: '备份管理',
            icon: 'ep:folder-checked'
          }
        },
        {
          path: 'aggregate-login',
          name: 'AggregateLogin',
          component: () => import('@/views/finance/system/aggregate-login/index.vue'),
          meta: {
            title: '聚合登录管理',
            icon: 'ep:platform'
          }
        },
        {
          path: 'server-groups',
          name: 'ServerGroups',
          component: () => import('@/views/finance/system/server-groups/index.vue'),
          meta: {
            title: '服务器组管理',
            icon: 'ep:monitor'
          }
        },
        {
          path: 'menu-groups',
          name: 'MenuGroups',
          component: () => import('@/views/finance/system/menu-groups/index.vue'),
          meta: {
            title: '菜单分组管理',
            icon: 'ep:menu'
          }
        },
        {
          path: 'sms-batch',
          name: 'SmsBatch',
          component: () => import('@/views/finance/system/sms-batch/index.vue'),
          meta: {
            title: '短信批量发送',
            icon: 'ep:message'
          }
        },
        {
          path: 'oauth-bind',
          name: 'OAuthBind',
          component: () => import('@/views/finance/system/oauth-bind/index.vue'),
          meta: {
            title: 'OAuth账号绑定',
            icon: 'ep:connection'
          }
        },
        {
          path: 'cache-clear',
          name: 'CacheClear',
          component: () => import('@/views/finance/system/cache-clear/index.vue'),
          meta: {
            title: '缓存清理管理',
            icon: 'ep:delete'
          }
        }
      ]
    },
    {
      path: 'advanced',
      name: 'Advanced',
      redirect: '/finance/advanced/contracts',
      meta: {
        title: '高级功能',
        icon: 'ep:star'
      },
      children: [
        {
          path: 'contracts',
          name: 'Contracts',
          component: () => import('@/views/finance/advanced/contracts/index.vue'),
          meta: {
            title: '合同管理',
            icon: 'ep:document'
          }
        },
        {
          path: 'client-care',
          name: 'ClientCare',
          component: () => import('@/views/finance/advanced/client-care/index.vue'),
          meta: {
            title: '客户关怀',
            icon: 'ep:present'
          }
        },
        {
          path: 'sales',
          name: 'Sales',
          component: () => import('@/views/finance/advanced/sales/index.vue'),
          meta: {
            title: '销售管理',
            icon: 'ep:trend-charts'
          }
        },
        {
          path: 'agents',
          name: 'Agents',
          component: () => import('@/views/finance/advanced/agents/index.vue'),
          meta: {
            title: '代理管理',
            icon: 'ep:connection'
          }
        },
        {
          path: 'upstream',
          name: 'Upstream',
          component: () => import('@/views/finance/advanced/upstream/index.vue'),
          meta: {
            title: '上游管理',
            icon: 'ep:upload'
          }
        },
        {
          path: 'reports',
          name: 'Reports',
          component: () => import('@/views/finance/advanced/reports/index.vue'),
          meta: {
            title: '报表统计',
            icon: 'ep:data-analysis'
          }
        },
        {
          path: 'upstream-providers',
          name: 'UpstreamProviders',
          component: () => import('@/views/finance/advanced/upstream-providers/index.vue'),
          meta: {
            title: '上游供应商',
            icon: 'ep:connection'
          }
        },
        {
          path: 'send-message',
          name: 'SendMessage',
          component: () => import('@/views/finance/advanced/send-message/index.vue'),
          meta: {
            title: '群发消息',
            icon: 'ep:promotion'
          }
        },
        {
          path: 'batch-sync',
          name: 'BatchSync',
          component: () => import('@/views/finance/advanced/batch-sync/index.vue'),
          meta: {
            title: '批量同步',
            icon: 'ep:refresh'
          }
        },
        {
          path: 'zjmf-api',
          name: 'ZjmfApi',
          component: () => import('@/views/finance/advanced/zjmf-api/index.vue'),
          meta: {
            title: '智简魔方API',
            icon: 'ep:connection'
          }
        }
      ]
    },
    {
      path: 'statistics',
      name: 'Statistics',
      redirect: '/finance/statistics/overview',
      meta: {
        title: '统计分析',
        icon: 'ep:data-analysis'
      },
      children: [
        {
          path: 'overview',
          name: 'StatisticsOverview',
          component: () => import('@/views/finance/statistics/overview/index.vue'),
          meta: {
            title: '统计概览',
            icon: 'ep:pie-chart'
          }
        },
        {
          path: 'revenue-ranking',
          name: 'RevenueRanking',
          component: () => import('@/views/finance/statistics/revenue-ranking/index.vue'),
          meta: {
            title: '收入排行',
            icon: 'ep:trophy'
          }
        },
        {
          path: 'annual',
          name: 'AnnualStatistics',
          component: () => import('@/views/finance/statistics/annual/index.vue'),
          meta: {
            title: '年度统计',
            icon: 'ep:calendar'
          }
        },
        {
          path: 'product-revenue',
          name: 'ProductRevenue',
          component: () => import('@/views/finance/statistics/product-revenue/index.vue'),
          meta: {
            title: '产品收入',
            icon: 'ep:goods'
          }
        },
        {
          path: 'new-customers',
          name: 'NewCustomers',
          component: () => import('@/views/finance/statistics/new-customers/index.vue'),
          meta: {
            title: '新客户统计',
            icon: 'ep:user'
          }
        }
      ]
    }
  ]
}
