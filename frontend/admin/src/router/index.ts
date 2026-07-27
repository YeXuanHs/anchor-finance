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
          component: () => import('@/views/Dashboard.vue')
        },
        {
          path: 'users',
          name: 'AdminUsers',
          component: () => import('@/views/Users.vue')
        },
        {
          path: 'products',
          name: 'AdminProducts',
          component: () => import('@/views/Products.vue')
        },
        {
          path: 'product-groups',
          name: 'AdminProductGroups',
          component: () => import('@/views/Products.vue')
        },
        {
          path: 'orders',
          name: 'AdminOrders',
          component: () => import('@/views/Orders.vue')
        },
        {
          path: 'invoices',
          name: 'AdminInvoices',
          component: () => import('@/views/Invoices.vue')
        },
        {
          path: 'tickets',
          name: 'AdminTickets',
          component: () => import('@/views/Tickets.vue')
        },
        {
          path: 'coupons',
          name: 'AdminCoupons',
          component: () => import('@/views/Coupons.vue')
        },
        {
          path: 'announcements',
          name: 'AdminAnnouncements',
          component: () => import('@/views/Announcements.vue')
        },
        {
          path: 'payments',
          name: 'AdminPayments',
          component: () => import('@/views/Payments.vue')
        },
        {
          path: 'oauth',
          name: 'AdminOAuth',
          component: () => import('@/views/OAuth.vue')
        },
        {
          path: 'settings',
          name: 'AdminSettings',
          component: () => import('@/views/Settings.vue')
        },
        {
          path: 'logs',
          name: 'AdminLogs',
          component: () => import('@/views/Logs.vue')
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
