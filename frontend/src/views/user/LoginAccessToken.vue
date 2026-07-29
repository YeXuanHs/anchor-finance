<template>
  <div class="access-token-page">
    <div class="page-header">
      <h2>令牌登录</h2>
      <p>使用 Access Token 快速登录</p>
    </div>
    <el-card>
      <el-form :model="form" label-width="100px">
        <el-form-item label="Access Token">
          <el-input v-model="form.token" placeholder="请输入Access Token" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleLogin" :loading="loading">登录</el-button>
        </el-form-item>
      </el-form>
      <el-divider />
      <h4>如何获取 Token？</h4>
      <ol>
        <li>前往 <router-link to="/user/api-manage">API管理</router-link> 页面</li>
        <li>创建或查看您的 API Key</li>
        <li>使用 API Key 作为 Access Token 登录</li>
      </ol>
    </el-card>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const router = useRouter()
const loading = ref(false)
const form = ref({ token: '' })

const handleLogin = async () => {
  if (!form.value.token) { ElMessage.warning('请输入Token'); return }
  loading.value = true
  try {
    const { data } = await request.post('/api/auth/access-token', { token: form.value.token })
    localStorage.setItem('token', data.data.token)
    ElMessage.success('登录成功')
    router.push('/user')
  } catch { ElMessage.error('Token无效或已过期') } finally { loading.value = false }
}
</script>
<style scoped lang="scss">
.access-token-page { .page-header { margin-bottom: 24px; h2 { font-size: 20px; color: #1a365d; } p { color: #6b7280; margin-top: 4px; } } }
ol { padding-left: 20px; li { color: #4b5563; line-height: 2; a { color: #2563eb; } } }
</style>
