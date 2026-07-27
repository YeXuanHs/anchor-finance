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
        // 订单管理
        {
          path: 'orders',
          name: 'Orders',
          component: () => import('@/views/orders/index.vue'),
          meta: { title: '订单管理', icon: 'Document' }
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
            }
          ]
        },
        // 日志管理
        {
          path: 'logs',
          name: 'Logs',
          component: () => import('@/views/logs/index.vue'),
          meta: { title: '操作日志', icon: 'Notebook' }
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
