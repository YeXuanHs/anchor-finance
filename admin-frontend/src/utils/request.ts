import axios from 'axios'
import { Message } from '@arco-design/web-vue'
import router from '@/router'

const request = axios.create({
  baseURL: '/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response) => {
    const res = response.data

    // 如果返回的code不是0，认为是错误
    if (res.code !== undefined && res.code !== 0) {
      Message.error(res.message || '请求失败')

      // 401: 未授权
      if (res.code === 401) {
        localStorage.removeItem('token')
        router.push('/login')
      }

      return Promise.reject(new Error(res.message || '请求失败'))
    }

    return res
  },
  (error) => {
    console.error('请求错误:', error)

    if (error.response) {
      const status = error.response.status

      switch (status) {
        case 401:
          Message.error('未授权，请重新登录')
          localStorage.removeItem('token')
          router.push('/login')
          break
        case 403:
          Message.error('没有权限访问')
          break
        case 404:
          Message.error('请求的资源不存在')
          break
        case 500:
          Message.error('服务器内部错误')
          break
        default:
          Message.error(`请求失败: ${status}`)
      }
    } else if (error.message.includes('timeout')) {
      Message.error('请求超时')
    } else {
      Message.error('网络错误')
    }

    return Promise.reject(error)
  }
)

export default request
