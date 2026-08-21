<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <img src="@/assets/logo.png" alt="锚点财务" class="logo" />
        <h1 class="title">锚点财务</h1>
        <p class="subtitle">管理后台</p>
      </div>

      <a-form
        :model="form"
        @submit="handleSubmit"
        layout="vertical"
        class="login-form"
      >
        <a-form-item field="username" label="用户名" :rules="[{ required: true, message: '请输入用户名' }]">
          <a-input
            v-model="form.username"
            placeholder="请输入用户名"
            size="large"
          >
            <template #prefix><icon-user /></template>
          </a-input>
        </a-form-item>

        <a-form-item field="password" label="密码" :rules="[{ required: true, message: '请输入密码' }]">
          <a-input-password
            v-model="form.password"
            placeholder="请输入密码"
            size="large"
          >
            <template #prefix><icon-lock /></template>
          </a-input-password>
        </a-form-item>

        <a-form-item>
          <a-button
            type="primary"
            html-type="submit"
            size="large"
            long
            :loading="loading"
          >
            登录
          </a-button>
        </a-form-item>
      </a-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { Message } from '@arco-design/web-vue'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)

const form = reactive({
  username: '',
  password: ''
})

const handleSubmit = async () => {
  if (!form.username || !form.password) {
    Message.warning('请输入用户名和密码')
    return
  }

  loading.value = true

  try {
    const success = await userStore.loginAction(form.username, form.password)

    if (success) {
      Message.success('登录成功')
      router.push('/')
    } else {
      Message.error('登录失败，请检查用户名和密码')
    }
  } catch (error) {
    Message.error('登录失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped lang="scss">
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
  width: 400px;
  padding: 40px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.logo {
  width: 64px;
  height: 64px;
  margin-bottom: 16px;
}

.title {
  font-size: 28px;
  font-weight: 600;
  color: #1d2129;
  margin: 0 0 8px 0;
}

.subtitle {
  font-size: 14px;
  color: #86909c;
  margin: 0;
}

.login-form {
  :deep(.arco-form-item-label) {
    font-weight: 500;
  }
}
</style>
