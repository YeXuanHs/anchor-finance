import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

NProgress.configure({ showSpinner: false })

const router = createRouter({
  history: createWebHistory('/admin/'),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/login/index.vue'),
      meta: { title: '登录' }
    },
    {
      path: '/',
      component: () => import('@/layouts/AdminLayout.vue'),
      redirect: '/dashboard',
      children: [
        {
          path: 'dashboard',
          name: 'Dashboard',
          component: () => import('@/views/dashboard/index.vue'),
          meta: { title: '仪表盘', icon: 'Odometer' }
        },
        // 用户管理
        {
          path: 'users',
          name: 'Users',
          component: () => import('@/views/users/index.vue'),
          meta: { title: '用户管理', icon: 'User' }
        },
        {
          path: 'users/:id',
          name: 'UserDetail',
          component: () => import('@/views/users/detail.vue'),
          meta: { title: '用户详情', hidden: true }
        },
        {
          path: 'users/contacts',
          name: 'UserContacts',
          component: () => import('@/views/users/contacts.vue'),
          meta: { title: '联系人管理', hidden: true }
        },
        {
          path: 'users/groups',
          name: 'UserGroups',
          component: () => import('@/views/users/groups.vue'),
          meta: { title: '用户分组', hidden: true }
        },
        {
          path: 'users/levels',
          name: 'UserLevels',
          component: () => import('@/views/users/levels.vue'),
          meta: { title: '用户等级', hidden: true }
        },
        {
          path: 'users/certification',
          name: 'UserCertification',
          component: () => import('@/views/users/certification.vue'),
          meta: { title: '实名认证', hidden: true }
        },
        // 产品管理
        {
          path: 'products',
          name: 'Products',
          component: () => import('@/views/products/index.vue'),
          meta: { title: '产品管理', icon: 'Box' }
        },
        {
          path: 'products/groups',
          name: 'ProductGroups',
          component: () => import('@/views/products/groups.vue'),
          meta: { title: '产品分组', hidden: true }
        },
        {
          path: 'products/categories',
          name: 'ProductCategories',
          component: () => import('@/views/products/categories.vue'),
          meta: { title: '产品分类', hidden: true }
        },
        {
          path: 'products/advanced-options',
          name: 'ProductAdvancedOptions',
          component: () => import('@/views/products/advanced-options.vue'),
          meta: { title: '高级配置', hidden: true }
        },
        // 订单管理
        {
          path: 'orders',
          name: 'Orders',
          component: () => import('@/views/orders/index.vue'),
          meta: { title: '订单管理', icon: 'Document' }
        },
        {
          path: 'orders/:id',
          name: 'OrderDetail',
          component: () => import('@/views/orders/detail.vue'),
          meta: { title: '订单详情', hidden: true }
        },
        // 财务管理
        {
          path: 'finance',
          name: 'Finance',
          redirect: '/finance/invoices',
          meta: { title: '财务管理', icon: 'Wallet' },
          children: [
            {
              path: 'invoices',
              name: 'Invoices',
              component: () => import('@/views/finance/invoices.vue'),
              meta: { title: '账单管理' }
            },
            {
              path: 'transactions',
              name: 'Transactions',
              component: () => import('@/views/finance/transactions.vue'),
              meta: { title: '交易记录' }
            },
            {
              path: 'payments',
              name: 'Payments',
              component: () => import('@/views/finance/payments.vue'),
              meta: { title: '支付方式' }
            },
            {
              path: 'vouchers',
              name: 'FinanceVouchers',
              component: () => import('@/views/finance/vouchers.vue'),
              meta: { title: '代金券' }
            },
            {
              path: 'credit',
              name: 'Credit',
              component: () => import('@/views/finance/credit.vue'),
              meta: { title: '信用额度' }
            },
            {
              path: 'currencies',
              name: 'Currencies',
              component: () => import('@/views/finance/currencies.vue'),
              meta: { title: '货币管理' }
            }
          ]
        },
        // 工单管理
        {
          path: 'tickets',
          name: 'Tickets',
          component: () => import('@/views/tickets/index.vue'),
          meta: { title: '工单管理', icon: 'Tickets' }
        },
        {
          path: 'tickets/:id',
          name: 'TicketDetail',
          component: () => import('@/views/tickets/detail.vue'),
          meta: { title: '工单详情', hidden: true }
        },
        {
          path: 'tickets/departments',
          name: 'TicketDepartments',
          component: () => import('@/views/tickets/departments.vue'),
          meta: { title: '工单部门', hidden: true }
        },
        {
          path: 'tickets/prereply',
          name: 'TicketPrereply',
          component: () => import('@/views/tickets/prereply.vue'),
          meta: { title: '预设回复', hidden: true }
        },
        {
          path: 'tickets/status',
          name: 'TicketStatus',
          component: () => import('@/views/tickets/status.vue'),
          meta: { title: '工单状态', hidden: true }
        },
        // 服务器管理
        {
          path: 'servers',
          name: 'Servers',
          redirect: '/servers/physical',
          meta: { title: '服务器管理', icon: 'Monitor' },
          children: [
            {
              path: 'physical',
              name: 'PhysicalServers',
              component: () => import('@/views/servers/physical.vue'),
              meta: { title: '物理服务器' }
            },
            {
              path: 'cloud',
              name: 'CloudServers',
              component: () => import('@/views/servers/cloud.vue'),
              meta: { title: '云服务器' }
            },
            {
              path: 'datacenters',
              name: 'Datacenters',
              component: () => import('@/views/servers/datacenters.vue'),
              meta: { title: '机房管理' }
            },
            {
              path: 'upstream',
              name: 'Upstream',
              component: () => import('@/views/servers/upstream.vue'),
              meta: { title: '上游管理' }
            },
            {
              path: 'batch-sync',
              name: 'BatchSync',
              component: () => import('@/views/servers/batch-sync.vue'),
              meta: { title: '批量同步' }
            }
          ]
        },
        // 内容管理
        {
          path: 'content',
          name: 'Content',
          redirect: '/content/news',
          meta: { title: '内容管理', icon: 'Document' },
          children: [
            {
              path: 'news',
              name: 'News',
              component: () => import('@/views/content/news.vue'),
              meta: { title: '新闻公告' }
            },
            {
              path: 'knowledge',
              name: 'Knowledge',
              component: () => import('@/views/content/knowledge.vue'),
              meta: { title: '知识库' }
            },
            {
              path: 'downloads',
              name: 'Downloads',
              component: () => import('@/views/content/downloads.vue'),
              meta: { title: '下载中心' }
            },
            {
              path: 'banners',
              name: 'Banners',
              component: () => import('@/views/content/banners.vue'),
              meta: { title: '轮播图' }
            },
            {
              path: 'solutions',
              name: 'Solutions',
              component: () => import('@/views/content/solutions.vue'),
              meta: { title: '解决方案' }
            },
            {
              path: 'pages',
              name: 'CustomPages',
              component: () => import('@/views/content/pages.vue'),
              meta: { title: '自定义页面' }
            }
          ]
        },
        // 营销管理
        {
          path: 'marketing',
          name: 'Marketing',
          redirect: '/marketing/coupons',
          meta: { title: '营销管理', icon: 'Present' },
          children: [
            {
              path: 'coupons',
              name: 'Coupons',
              component: () => import('@/views/marketing/coupons.vue'),
              meta: { title: '优惠券' }
            },
            {
              path: 'promotions',
              name: 'Promotions',
              component: () => import('@/views/marketing/promotions.vue'),
              meta: { title: '促销活动' }
            },
            {
              path: 'affiliate',
              name: 'Affiliate',
              component: () => import('@/views/marketing/affiliate.vue'),
              meta: { title: '推介计划' }
            },
            {
              path: 'referral',
              name: 'Referral',
              component: () => import('@/views/marketing/referral.vue'),
              meta: { title: '推介管理' }
            },
            {
              path: 'vouchers',
              name: 'MarketingVouchers',
              component: () => import('@/views/marketing/vouchers.vue'),
              meta: { title: '代金券活动' }
            }
          ]
        },
        // 系统设置
        {
          path: 'system',
          name: 'System',
          redirect: '/system/general',
          meta: { title: '系统设置', icon: 'Setting' },
          children: [
            {
              path: 'general',
              name: 'General',
              component: () => import('@/views/system/general.vue'),
              meta: { title: '基本设置' }
            },
            {
              path: 'notification',
              name: 'Notification',
              component: () => import('@/views/system/notification.vue'),
              meta: { title: '通知设置' }
            },
            {
              path: 'security',
              name: 'Security',
              component: () => import('@/views/system/security.vue'),
              meta: { title: '安全设置' }
            },
            {
              path: 'oauth',
              name: 'OAuth',
              component: () => import('@/views/system/oauth.vue'),
              meta: { title: 'OAuth配置' }
            },
            {
              path: 'payment',
              name: 'Payment',
              component: () => import('@/views/system/payment.vue'),
              meta: { title: '支付配置' }
            },
            {
              path: 'email',
              name: 'Email',
              component: () => import('@/views/system/email.vue'),
              meta: { title: '邮件配置' }
            },
            {
              path: 'sms',
              name: 'SMS',
              component: () => import('@/views/system/sms.vue'),
              meta: { title: '短信配置' }
            },
            {
              path: 'cron',
              name: 'Cron',
              component: () => import('@/views/system/cron.vue'),
              meta: { title: '定时任务' }
            },
            {
              path: 'certification',
              name: 'SystemCertification',
              component: () => import('@/views/system/certification.vue'),
              meta: { title: '认证配置' }
            },
            {
              path: 'servers',
              name: 'SystemServers',
              component: () => import('@/views/system/servers.vue'),
              meta: { title: '服务器配置' }
            },
            {
              path: 'advanced',
              name: 'Advanced',
              component: () => import('@/views/system/advanced.vue'),
              meta: { title: '高级设置' }
            },
            {
              path: 'menu',
              name: 'Menu',
              component: () => import('@/views/system/menu.vue'),
              meta: { title: '菜单管理' }
            },
            {
              path: 'rbac',
              name: 'RBAC',
              component: () => import('@/views/system/rbac.vue'),
              meta: { title: '权限管理' }
            },
            {
              path: 'roles',
              name: 'Roles',
              component: () => import('@/views/system/roles.vue'),
              meta: { title: '角色管理' }
            }
          ]
        },
        // 日志管理
        {
          path: 'logs',
          name: 'Logs',
          component: () => import('@/views/logs.vue'),
          meta: { title: '操作日志', icon: 'Notebook' }
        },
        // 插件管理
        {
          path: 'plugins',
          name: 'Plugins',
          component: () => import('@/views/plugins/index.vue'),
          meta: { title: '插件管理', icon: 'Connection' }
        },
        {
          path: 'plugins/store',
          name: 'PluginStore',
          component: () => import('@/views/plugins/store.vue'),
          meta: { title: '应用商店', hidden: true }
        },
        // 报表统计
        {
          path: 'reports',
          name: 'Reports',
          component: () => import('@/views/reports/index.vue'),
          meta: { title: '报表统计', icon: 'DataAnalysis' }
        },
        // 客户关怀
        {
          path: 'client-care',
          name: 'ClientCare',
          component: () => import('@/views/client-care/index.vue'),
          meta: { title: '客户关怀', icon: 'Star' }
        },
        // 社区管理
        {
          path: 'community',
          name: 'Community',
          component: () => import('@/views/community/index.vue'),
          meta: { title: '社区管理', icon: 'ChatDotRound' }
        },
        // URL定时任务
        {
          path: 'cron-url',
          name: 'CronUrl',
          component: () => import('@/views/cron-url/index.vue'),
          meta: { title: 'URL定时任务', icon: 'Link' }
        },
        // 魔方云对接
        {
          path: 'dcim-cloud',
          name: 'DcimCloud',
          component: () => import('@/views/dcim-cloud/index.vue'),
          meta: { title: '魔方云对接', icon: 'Cloudy' }
        },
        // 获取用户
        {
          path: 'get-user',
          name: 'GetUser',
          component: () => import('@/views/get-user/index.vue'),
          meta: { title: '获取用户', icon: 'UserFilled' }
        },
        // 关联原因
        {
          path: 'link-cause',
          name: 'LinkCause',
          component: () => import('@/views/link-cause/index.vue'),
          meta: { title: '关联原因', icon: 'Connection' }
        },
        // 关联知识
        {
          path: 'link-knowledge',
          name: 'LinkKnowledge',
          component: () => import('@/views/link-knowledge/index.vue'),
          meta: { title: '关联知识', icon: 'Document' }
        },
        // DCIM管理
        {
          path: 'dcim',
          name: 'DCIM',
          component: () => import('@/views/dcim/index.vue'),
          meta: { title: 'DCIM管理', icon: 'Cpu' }
        },
        // 高级配置
        {
          path: 'advanced-options',
          name: 'AdvancedOptions',
          component: () => import('@/views/advanced-options/index.vue'),
          meta: { title: '高级配置', icon: 'Tools' }
        },
        {
          path: ':pathMatch(.*)*',
          name: 'NotFound',
          component: () => import('@/views/not-found.vue')
        }
      ]
    }
  ]
})

router.beforeEach((to, from, next) => {
  NProgress.start()

  if (to.path !== '/login') {
    const userStore = useUserStore()
    if (!userStore.token) {
      next(`/login?redirect=${to.path}`)
      return
    }
  }

  document.title = `${to.meta.title || ''} - 锚点财务管理后台`
  next()
})

router.afterEach(() => {
  NProgress.done()
})

export default router
