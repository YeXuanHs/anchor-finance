import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { adminApi } from '@/api'
import router from '@/router'

export interface AdminUser {
  id: number
  username: string
  email: string
  avatar?: string
  role: string
  lastLogin?: string
}

export interface DashboardStats {
  totalUsers: number
  totalOrders: number
  totalRevenue: number
  openTickets: number
}

export const useAdminStore = defineStore('admin', () => {
  const token = ref<string>(localStorage.getItem('admin_token') || '')
  const adminInfo = ref<AdminUser | null>(null)
  const sidebarCollapsed = ref(false)

  const isLoggedIn = computed(() => !!token.value)

  async function login(username: string, password: string) {
    const res = await adminApi.login(username, password)
    token.value = res.data.token
    adminInfo.value = res.data.admin
    localStorage.setItem('admin_token', res.data.token)
    router.push('/admin/dashboard')
  }

  function logout() {
    token.value = ''
    adminInfo.value = null
    localStorage.removeItem('admin_token')
    router.push('/admin/login')
  }

  async function fetchAdminInfo() {
    if (!token.value) return
    try {
      const res = await adminApi.getProfile()
      adminInfo.value = res.data
    } catch {
      logout()
    }
  }

  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  return {
    token,
    adminInfo,
    sidebarCollapsed,
    isLoggedIn,
    login,
    logout,
    fetchAdminInfo,
    toggleSidebar,
  }
})
