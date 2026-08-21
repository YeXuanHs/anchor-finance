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
      path: '/register',
      name: 'Register',
      component: () => import('@/views/auth/Register.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/',
      component: MainLayout,
      redirect: '/home',
      meta: { requiresAuth: true },
      children: [
        // 首页
        {
          path: 'home',
          name: 'Home',
          component: () => import('@/views/home/Index.vue'),
          meta: { title: '首页' }
        },
        // 产品
        {
          path: 'products',
          name: 'Products',
          component: () => import('@/views/product/List.vue'),
          meta: { title: '产品中心' }
        },
        {
          path: 'products/:id',
          name: 'ProductDetail',
          component: () => import('@/views/product/Detail.vue'),
          meta: { title: '产品详情' }
        },
        // 服务
        {
          path: 'services',
          name: 'Services',
          component: () => import('@/views/service/List.vue'),
          meta: { title: '我的服务' }
        },
        {
          path: 'services/:id',
          name: 'ServiceDetail',
          component: () => import('@/views/service/Detail.vue'),
          meta: { title: '服务详情' }
        },
        // 订单
        {
          path: 'orders',
          name: 'Orders',
          component: () => import('@/views/order/List.vue'),
          meta: { title: '我的订单' }
        },
        {
          path: 'orders/:id',
          name: 'OrderDetail',
          component: () => import('@/views/order/Detail.vue'),
          meta: { title: '订单详情' }
        },
        // 工单
        {
          path: 'tickets',
          name: 'Tickets',
          component: () => import('@/views/ticket/List.vue'),
          meta: { title: '我的工单' }
        },
        {
          path: 'tickets/create',
          name: 'CreateTicket',
          component: () => import('@/views/ticket/Create.vue'),
          meta: { title: '提交工单' }
        },
        {
          path: 'tickets/:id',
          name: 'TicketDetail',
          component: () => import('@/views/ticket/Detail.vue'),
          meta: { title: '工单详情' }
        },
        // 财务
        {
          path: 'finance',
          name: 'Finance',
          component: () => import('@/views/finance/Index.vue'),
          meta: { title: '财务中心' }
        },
        {
          path: 'finance/invoices',
          name: 'Invoices',
          component: () => import('@/views/finance/Invoices.vue'),
          meta: { title: '我的账单' }
        },
        {
          path: 'finance/transactions',
          name: 'Transactions',
          component: () => import('@/views/finance/Transactions.vue'),
          meta: { title: '交易记录' }
        },
        // 用户中心
        {
          path: 'user',
          name: 'User',
          component: () => import('@/views/user/Index.vue'),
          meta: { title: '个人中心' }
        },
        {
          path: 'user/profile',
          name: 'Profile',
          component: () => import('@/views/user/Profile.vue'),
          meta: { title: '个人资料' }
        },
        {
          path: 'user/security',
          name: 'Security',
          component: () => import('@/views/user/Security.vue'),
          meta: { title: '安全设置' }
        }
      ]
    }
  ]
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('client_token')

  if (to.meta.requiresAuth !== false && !token) {
    next('/login')
  } else if ((to.path === '/login' || to.path === '/register') && token) {
    next('/')
  } else {
    next()
  }
})

export default router
