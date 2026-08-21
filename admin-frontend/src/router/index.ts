import { createRouter, createWebHistory } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/auth/Login.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/',
      component: MainLayout,
      redirect: '/dashboard',
      meta: { requiresAuth: true },
      children: [
        // 仪表盘
        {
          path: 'dashboard',
          name: 'Dashboard',
          component: () => import('@/views/dashboard/Index.vue'),
          meta: { title: '仪表盘', icon: 'icon-dashboard' }
        },
        // 客户管理
        {
          path: 'customer',
          name: 'Customer',
          redirect: '/customer/list',
          meta: { title: '客户', icon: 'icon-user' },
          children: [
            {
              path: 'list',
              name: 'CustomerList',
              component: () => import('@/views/customer/List.vue'),
              meta: { title: '客户列表' }
            },
            {
              path: 'authentication',
              name: 'CustomerAuthentication',
              component: () => import('@/views/customer/Authentication.vue'),
              meta: { title: '实名认证' }
            },
            {
              path: 'resources',
              name: 'CustomerResources',
              component: () => import('@/views/customer/Resources.vue'),
              meta: { title: '客户资源池' }
            }
          ]
        },
        // 订单管理
        {
          path: 'order',
          name: 'Order',
          redirect: '/order/list',
          meta: { title: '业务', icon: 'icon-file' },
          children: [
            {
              path: 'list',
              name: 'OrderList',
              component: () => import('@/views/order/List.vue'),
              meta: { title: '产品订单' }
            },
            {
              path: 'renewal',
              name: 'RenewalOrder',
              component: () => import('@/views/order/Renewal.vue'),
              meta: { title: '续费订单' }
            },
            {
              path: 'service',
              name: 'ServiceList',
              component: () => import('@/views/order/Service.vue'),
              meta: { title: '业务列表' }
            }
          ]
        },
        // 财务管理
        {
          path: 'finance',
          name: 'Finance',
          redirect: '/finance/transactions',
          meta: { title: '财务', icon: 'icon-money-circle' },
          children: [
            {
              path: 'transactions',
              name: 'Transactions',
              component: () => import('@/views/finance/Transactions.vue'),
              meta: { title: '交易流水' }
            },
            {
              path: 'invoices',
              name: 'Invoices',
              component: () => import('@/views/finance/Invoices.vue'),
              meta: { title: '账单管理' }
            },
            {
              path: 'credit',
              name: 'CreditManagement',
              component: () => import('@/views/finance/Credit.vue'),
              meta: { title: '信用额管理' }
            }
          ]
        },
        // 工单管理
        {
          path: 'ticket',
          name: 'Ticket',
          redirect: '/ticket/list',
          meta: { title: '工单', icon: 'icon-customer-service' },
          children: [
            {
              path: 'list',
              name: 'TicketList',
              component: () => import('@/views/ticket/List.vue'),
              meta: { title: '工单列表' }
            },
            {
              path: 'statistics',
              name: 'TicketStatistics',
              component: () => import('@/views/ticket/Statistics.vue'),
              meta: { title: '工单统计' }
            }
          ]
        },
        // 功能
        {
          path: 'feature',
          name: 'Feature',
          redirect: '/feature/plugins',
          meta: { title: '功能', icon: 'icon-apps' },
          children: [
            {
              path: 'plugins',
              name: 'Plugins',
              component: () => import('@/views/plugin/List.vue'),
              meta: { title: '插件列表' }
            },
            {
              path: 'statistics',
              name: 'Statistics',
              component: () => import('@/views/feature/Statistics.vue'),
              meta: { title: '统计' }
            }
          ]
        },
        // 设置
        {
          path: 'setting',
          name: 'Setting',
          redirect: '/setting/general',
          meta: { title: '设置', icon: 'icon-settings' },
          children: [
            {
              path: 'general',
              name: 'GeneralSettings',
              component: () => import('@/views/setting/General.vue'),
              meta: { title: '常规设置' }
            },
            {
              path: 'payment',
              name: 'PaymentSettings',
              component: () => import('@/views/setting/Payment.vue'),
              meta: { title: '支付接口' }
            },
            {
              path: 'security',
              name: 'SecuritySettings',
              component: () => import('@/views/setting/Security.vue'),
              meta: { title: '安全相关' }
            }
          ]
        }
      ]
    }
  ]
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')

  if (to.meta.requiresAuth !== false && !token) {
    next('/login')
  } else if (to.path === '/login' && token) {
    next('/')
  } else {
    next()
  }
})

export default router
