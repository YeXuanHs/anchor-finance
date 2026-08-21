import request from '@/utils/request'

// 登录
export function login(data: { username: string; password: string }) {
  return request({
    url: '/login',
    method: 'post',
    data
  })
}

// 注册
export function register(data: { username: string; email: string; password: string }) {
  return request({
    url: '/register',
    method: 'post',
    data
  })
}

// 获取用户信息
export function getUserInfo() {
  return request({
    url: '/user/info',
    method: 'get'
  })
}

// 更新用户信息
export function updateUserInfo(data: any) {
  return request({
    url: '/user/info',
    method: 'put',
    data
  })
}

// 修改密码
export function changePassword(data: { old_password: string; new_password: string }) {
  return request({
    url: '/user/password',
    method: 'put',
    data
  })
}

// 重置密码
export function resetPassword(data: { email: string }) {
  return request({
    url: '/password/reset',
    method: 'post',
    data
  })
}
