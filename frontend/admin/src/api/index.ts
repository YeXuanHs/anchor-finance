import axios from 'axios'
import { useAdminStore } from '@/stores/admin'

const http = axios.create({
  baseURL: '/api',
  timeout: 15000,
})

http.interceptors.request.use((config) => {
  const store = useAdminStore()
  if (store.token) {
    config.headers.Authorization = `Bearer ${store.token}`
  }
  return config
})

http.interceptors.response.use(
  (res) => res.data,
  (err) => {
    if (err.response?.status === 401) {
      const store = useAdminStore()
      store.logout()
    }
    return Promise.reject(err)
  }
)

// ── Admin Auth ──
export const adminApi = {
  login: (username: string, password: string) =>
    http.post('/admin/login', { username, password }),
  getProfile: () => http.get('/admin/profile'),
}

// ── Dashboard ──
export const dashboardApi = {
  getStats: () => http.get('/admin/dashboard/stats'),
  getRevenueChart: (params?: { period?: string }) =>
    http.get('/admin/dashboard/revenue', { params }),
  getRecentOrders: (limit?: number) =>
    http.get('/admin/dashboard/recent-orders', { params: { limit } }),
  getOrderStatusStats: () => http.get('/admin/dashboard/order-status'),
}

// ── Users ──
export const usersApi = {
  list: (params: any) => http.get('/admin/users', { params }),
  get: (id: number) => http.get(`/admin/users/${id}`),
  create: (data: any) => http.post('/admin/users', data),
  update: (id: number, data: any) => http.put(`/admin/users/${id}`, data),
  delete: (id: number) => http.delete(`/admin/users/${id}`),
  toggleStatus: (id: number) => http.patch(`/admin/users/${id}/toggle-status`),
}

// ── Products ──
export const productsApi = {
  list: (params?: any) => http.get('/admin/products', { params }),
  get: (id: number) => http.get(`/admin/products/${id}`),
  create: (data: any) => http.post('/admin/products', data),
  update: (id: number, data: any) => http.put(`/admin/products/${id}`, data),
  delete: (id: number) => http.delete(`/admin/products/${id}`),
  getGroups: () => http.get('/admin/product-groups'),
  createGroup: (data: any) => http.post('/admin/product-groups', data),
  updateGroup: (id: number, data: any) => http.put(`/admin/product-groups/${id}`, data),
  deleteGroup: (id: number) => http.delete(`/admin/product-groups/${id}`),
}

// ── Orders ──
export const ordersApi = {
  list: (params?: any) => http.get('/admin/orders', { params }),
  get: (id: number) => http.get(`/admin/orders/${id}`),
  updateStatus: (id: number, status: string) =>
    http.patch(`/admin/orders/${id}/status`, { status }),
  delete: (id: number) => http.delete(`/admin/orders/${id}`),
}

// ── Invoices ──
export const invoicesApi = {
  list: (params?: any) => http.get('/admin/invoices', { params }),
  get: (id: number) => http.get(`/admin/invoices/${id}`),
  markPaid: (id: number) => http.patch(`/admin/invoices/${id}/mark-paid`),
  delete: (id: number) => http.delete(`/admin/invoices/${id}`),
}

// ── Tickets ──
export const ticketsApi = {
  list: (params?: any) => http.get('/admin/tickets', { params }),
  get: (id: number) => http.get(`/admin/tickets/${id}`),
  assign: (id: number, adminId: number) =>
    http.patch(`/admin/tickets/${id}/assign`, { adminId }),
  reply: (id: number, content: string) =>
    http.post(`/admin/tickets/${id}/reply`, { content }),
  close: (id: number) => http.patch(`/admin/tickets/${id}/close`),
}

// ── Announcements ──
export const announcementsApi = {
  list: (params?: any) => http.get('/admin/announcements', { params }),
  get: (id: number) => http.get(`/admin/announcements/${id}`),
  create: (data: any) => http.post('/admin/announcements', data),
  update: (id: number, data: any) => http.put(`/admin/announcements/${id}`, data),
  delete: (id: number) => http.delete(`/admin/announcements/${id}`),
}

// ── Settings ──
export const settingsApi = {
  get: () => http.get('/admin/settings'),
  update: (data: any) => http.put('/admin/settings', data),
  testEmail: () => http.post('/admin/settings/test-email'),
  testSms: () => http.post('/admin/settings/test-sms'),
}

export default http
