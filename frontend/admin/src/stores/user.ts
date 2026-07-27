import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import request from '@/utils/request'

interface UserInfo {
  id: number
  username: string
  email: string
  avatar?: string
  is_admin: boolean
  permissions: string[]
}

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('admin_token') || '')
  const userInfo = ref<UserInfo | null>(null)

  const isLoggedIn = computed(() => !!token.value)
  const username = computed(() => userInfo.value?.username || '管理员')
  const avatar = computed(() => userInfo.value?.avatar || '')

  async function login(username: string, password: string) {
    const { data } = await request.post('/api/admin/login', { username, password })
    if (data.ok) {
      token.value = data.token
      localStorage.setItem('admin_token', data.token)
      await getUserInfo()
      return true
    }
    throw new Error(data.message || '登录失败')
  }

  async function getUserInfo() {
    try {
      const { data } = await request.get('/api/admin/user/info')
      if (data.ok) {
        userInfo.value = data.data
      }
    } catch (error) {
      console.error('获取用户信息失败:', error)
    }
  }

  function logout() {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('admin_token')
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    username,
    avatar,
    login,
    getUserInfo,
    logout
  }
})
