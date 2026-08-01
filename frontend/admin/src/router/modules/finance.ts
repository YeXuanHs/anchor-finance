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
      redirect: '/finance/marketing/coupons',
      meta: {
        title: '营销推广',
        icon: 'ep:present'
      },
      children: [
        {
          path: 'coupons',
          name: 'Coupons',
          component: () => import('@/views/finance/marketing/coupons/index.vue'),
          meta: {
            title: '优惠券',
            icon: 'ep:ticket'
          }
        },
        {
          path: 'vouchers',
          name: 'Vouchers',
          component: () => import('@/views/finance/marketing/vouchers/index.vue'),
          meta: {
            title: '代金券',
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
          path: 'cs-widgets',
          name: 'CSWidgets',
          component: () => import('@/views/finance/system/cs-widgets/index.vue'),
          meta: {
            title: '客服悬浮窗',
            icon: 'ep:service'
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
        }
      ]
    }
  ]
}
