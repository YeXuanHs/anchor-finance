import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'Home',
      component: () => import('@/views/Home.vue')
    },
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/Login.vue')
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('@/views/Register.vue')
    },
    {
      path: '/products',
      name: 'Products',
      component: () => import('@/views/Products.vue')
    },
    {
      path: '/products/:id',
      name: 'ProductDetail',
      component: () => import('@/views/ProductDetail.vue')
    },
    {
      path: '/cart',
      name: 'Cart',
      component: () => import('@/views/Cart.vue')
    },
    {
      path: '/user',
      component: () => import('@/layouts/UserLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          redirect: '/user/dashboard'
        },
        {
          path: 'dashboard',
          name: 'UserDashboard',
          component: () => import('@/views/user/Dashboard.vue')
        },
        {
          path: 'products',
          name: 'UserProducts',
          component: () => import('@/views/user/Products.vue')
        },
        {
          path: 'orders',
          name: 'UserOrders',
          component: () => import('@/views/user/Orders.vue')
        },
        {
          path: 'invoices',
          name: 'UserInvoices',
          component: () => import('@/views/user/Invoices.vue')
        },
        {
          path: 'tickets',
          name: 'UserTickets',
          component: () => import('@/views/user/Tickets.vue')
        },
        {
          path: 'profile',
          name: 'UserProfile',
          component: () => import('@/views/user/Profile.vue')
        },
        {
          path: 'coupons',
          name: 'UserCoupons',
          component: () => import('@/views/user/Coupons.vue')
        },
        {
          path: 'referral',
          name: 'UserReferral',
          component: () => import('@/views/user/Referral.vue')
        },
        {
          path: 'verification',
          name: 'UserVerification',
          component: () => import('@/views/user/Verification.vue')
        },
        {
          path: 'security',
          name: 'UserSecurity',
          component: () => import('@/views/user/Profile.vue')
        }
      ]
    }
  ]
})

// 路由守卫
router.beforeEach((to, from, next) => {
  if (to.matched.some(record => record.meta.requiresAuth)) {
    const userStore = useUserStore()
    if (!userStore.token) {
      next({
        path: '/login',
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
