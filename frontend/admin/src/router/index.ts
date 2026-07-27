import { createRouter, createWebHistory } from 'vue-router'
import { useAdminStore } from '@/stores/admin'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/admin/login',
      name: 'AdminLogin',
      component: () => import('@/views/Login.vue')
    },
    {
      path: '/admin',
      component: () => import('@/layouts/AdminLayout.vue'),
      meta: { requiresAdmin: true },
      children: [
        {
          path: '',
          redirect: '/admin/dashboard'
        },
        {
          path: 'dashboard',
          name: 'AdminDashboard',
          component: () => import('@/views/Dashboard.vue'),
          meta: { title: '仪表盘' }
        },
        {
          path: 'users',
          name: 'AdminUsers',
          component: () => import('@/views/Users.vue'),
          meta: { title: '用户管理' }
        },
        {
          path: 'agents',
          name: 'AdminAgents',
          component: () => import('@/views/Agents.vue'),
          meta: { title: '代理商管理' }
        },
        {
          path: 'products',
          name: 'AdminProducts',
          component: () => import('@/views/Products.vue'),
          meta: { title: '产品管理' }
        },
        {
          path: 'product-groups',
          name: 'AdminProductGroups',
          component: () => import('@/views/Products.vue'),
          meta: { title: '产品组管理' }
        },
        {
          path: 'orders',
          name: 'AdminOrders',
          component: () => import('@/views/Orders.vue'),
          meta: { title: '订单管理' }
        },
        {
          path: 'invoices',
          name: 'AdminInvoices',
          component: () => import('@/views/Invoices.vue'),
          meta: { title: '账单管理' }
        },
        {
          path: 'tickets',
          name: 'AdminTickets',
          component: () => import('@/views/Tickets.vue'),
          meta: { title: '工单管理' }
        },
        {
          path: 'coupons',
          name: 'AdminCoupons',
          component: () => import('@/views/Coupons.vue'),
          meta: { title: '优惠券管理' }
        },
        {
          path: 'announcements',
          name: 'AdminAnnouncements',
          component: () => import('@/views/Announcements.vue'),
          meta: { title: '公告管理' }
        },
        {
          path: 'email-templates',
          name: 'AdminEmailTemplates',
          component: () => import('@/views/EmailTemplates.vue'),
          meta: { title: '邮件模板' }
        },
        {
          path: 'notifications',
          name: 'AdminNotifications',
          component: () => import('@/views/Notifications.vue'),
          meta: { title: '通知管理' }
        },
        {
          path: 'payments',
          name: 'AdminPayments',
          component: () => import('@/views/Payments.vue'),
          meta: { title: '支付管理' }
        },
        {
          path: 'oauth',
          name: 'AdminOAuth',
          component: () => import('@/views/OAuth.vue'),
          meta: { title: '第三方登录' }
        },
        {
          path: 'reports',
          name: 'AdminReports',
          component: () => import('@/views/Reports.vue'),
          meta: { title: '报表统计' }
        },
        {
          path: 'logs',
          name: 'AdminLogs',
          component: () => import('@/views/Logs.vue'),
          meta: { title: '系统日志' }
        },
        {
          path: 'settings',
          name: 'AdminSettings',
          component: () => import('@/views/Settings.vue'),
          meta: { title: '系统设置' }
        }
      ]
    }
  ]
})

router.beforeEach((to, from, next) => {
  if (to.matched.some(record => record.meta.requiresAdmin)) {
    const adminStore = useAdminStore()
    if (!adminStore.token) {
      next({
        path: '/admin/login',
        query: { redirect: to.fullPath }
      })
    } else {
      next()
    }
  } else {
    next()
  }
})

export default router
