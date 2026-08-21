import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login, logout, getUserInfo } from '@/api/auth'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const username = ref('')
  const avatar = ref('')
  const role = ref('')
  const permissions = ref<string[]>([])

  // 登录
  async function loginAction(username: string, password: string) {
    try {
      const res = await login({ username, password })
      if (res.code === 0) {
        token.value = res.data.token
        localStorage.setItem('token', res.data.token)
        return true
      }
      return false
    } catch (error) {
      console.error('Login failed:', error)
      return false
    }
  }

  // 获取用户信息
  async function getUserInfoAction() {
    try {
      const res = await getUserInfo()
      if (res.code === 0) {
        username.value = res.data.username
        avatar.value = res.data.avatar || ''
        role.value = res.data.role || ''
        permissions.value = res.data.permissions || []
        return true
      }
      return false
    } catch (error) {
      console.error('Get user info failed:', error)
      return false
    }
  }

  // 登出
  async function logoutAction() {
    try {
      await logout()
    } catch (error) {
      console.error('Logout failed:', error)
    } finally {
      logout()
    }
  }

  // 清除用户信息
  function logout() {
    token.value = ''
    username.value = ''
    avatar.value = ''
    role.value = ''
    permissions.value = []
    localStorage.removeItem('token')
  }

  return {
    token,
    username,
    avatar,
    role,
    permissions,
    loginAction,
    getUserInfoAction,
    logoutAction,
    logout
  }
})
