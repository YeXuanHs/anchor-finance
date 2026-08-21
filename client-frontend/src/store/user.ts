import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login, register, getUserInfo } from '@/api/auth'

export const useUserStore = defineStore('client_user', () => {
  const token = ref(localStorage.getItem('client_token') || '')
  const username = ref('')
  const email = ref('')
  const avatar = ref('')
  const balance = ref(0)
  const userInfo = ref<any>(null)

  const isLoggedIn = computed(() => !!token.value)

  // 登录
  async function loginAction(usernameOrEmail: string, password: string) {
    try {
      const res = await login({ username: usernameOrEmail, password })
      if (res.code === 0) {
        token.value = res.data.token
        localStorage.setItem('client_token', res.data.token)
        await getUserInfoAction()
        return true
      }
      return false
    } catch (error) {
      console.error('Login failed:', error)
      return false
    }
  }

  // 注册
  async function registerAction(data: { username: string; email: string; password: string }) {
    try {
      const res = await register(data)
      if (res.code === 0) {
        return true
      }
      return false
    } catch (error) {
      console.error('Register failed:', error)
      return false
    }
  }

  // 获取用户信息
  async function getUserInfoAction() {
    try {
      const res = await getUserInfo()
      if (res.code === 0) {
        username.value = res.data.username
        email.value = res.data.email
        avatar.value = res.data.avatar || ''
        balance.value = res.data.balance || 0
        userInfo.value = res.data
        return true
      }
      return false
    } catch (error) {
      console.error('Get user info failed:', error)
      return false
    }
  }

  // 登出
  function logout() {
    token.value = ''
    username.value = ''
    email.value = ''
    avatar.value = ''
    balance.value = 0
    userInfo.value = null
    localStorage.removeItem('client_token')
  }

  // 初始化
  async function init() {
    if (token.value) {
      await getUserInfoAction()
    }
  }

  return {
    token,
    username,
    email,
    avatar,
    balance,
    userInfo,
    isLoggedIn,
    loginAction,
    registerAction,
    getUserInfoAction,
    logout,
    init
  }
})
