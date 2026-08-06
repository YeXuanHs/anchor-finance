import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/api'

interface User {
  id: number
  username: string
  email: string
  phone?: string
  balance: number
  created_at: string
}

export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(null)
  const token = ref<string>(localStorage.getItem('token') || '')

  const isLoggedIn = computed(() => !!token.value)
  const username = computed(() => user.value?.username || '')

  async function login(username: string, password: string, captcha: string) {
    const { data } = await api.post('/api/v1/auth/login', {
      username,
      password,
      captcha
    })
    if (data.code === 0) {
      token.value = data.data.token
      localStorage.setItem('token', data.data.token)
      await fetchProfile()
      return true
    }
    throw new Error(data.message || '登录失败')
  }

  async function register(form: {
    username: string
    email: string
    password: string
    captcha: string
  }) {
    const { data } = await api.post('/api/v1/auth/register', form)
    if (data.code === 0) {
      return true
    }
    throw new Error(data.message || '注册失败')
  }

  async function fetchProfile() {
    if (!token.value) return
    try {
      const { data } = await api.get('/api/v1/user/profile')
      if (data.code === 0) {
        user.value = data.data
      }
    } catch {
      token.value = ''
      localStorage.removeItem('token')
    }
  }

  function setToken(newToken: string) {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  function setUserInfo(userInfo: User) {
    user.value = userInfo
  }

  function logout() {
    user.value = null
    token.value = ''
    localStorage.removeItem('token')
  }

  return {
    user,
    token,
    isLoggedIn,
    username,
    login,
    register,
    fetchProfile,
    setToken,
    setUserInfo,
    logout
  }
})