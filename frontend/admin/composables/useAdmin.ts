import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface AdminUser {
  id: number
  username: string
  email: string
  avatar?: string
  role: string
  lastLogin?: string
}

export const useAdminStore = defineStore('admin', () => {
  const token = ref<string>('')
  const adminInfo = ref<AdminUser | null>(null)
  const sidebarCollapsed = ref(false)

  const isLoggedIn = computed(() => !!token.value)

  // Initialize from localStorage on client side
  if (import.meta.client) {
    token.value = localStorage.getItem('admin_token') || ''
  }

  async function login(username: string, password: string) {
    // Mock login - replace with actual API call
    const mockAdmin: AdminUser = {
      id: 1,
      username: username,
      email: 'admin@anchorfinance.com',
      role: 'super_admin',
      lastLogin: new Date().toISOString(),
    }

    token.value = 'mock_token_' + Date.now()
    adminInfo.value = mockAdmin

    if (import.meta.client) {
      localStorage.setItem('admin_token', token.value)
    }

    await navigateTo('/admin/dashboard')
  }

  function logout() {
    token.value = ''
    adminInfo.value = null

    if (import.meta.client) {
      localStorage.removeItem('admin_token')
    }

    navigateTo('/admin/login')
  }

  async function fetchAdminInfo() {
    if (!token.value) return

    // Mock fetch - replace with actual API call
    adminInfo.value = {
      id: 1,
      username: 'Admin',
      email: 'admin@anchorfinance.com',
      role: 'super_admin',
      lastLogin: new Date().toISOString(),
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
