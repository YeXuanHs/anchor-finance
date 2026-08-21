import request from '@/utils/request'

// 登录
export function login(data: { username: string; password: string }) {
  return request({
    url: '/admin/login',
    method: 'post',
    data
  })
}

// 登出
export function logout() {
  return request({
    url: '/admin/logout',
    method: 'post'
  })
}

// 获取用户信息
export function getUserInfo() {
  return request({
    url: '/admin/auth/info',
    method: 'get'
  })
}
