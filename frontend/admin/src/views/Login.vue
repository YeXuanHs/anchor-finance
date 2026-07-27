<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <el-icon :size="52" color="#0056FF"><Connection /></el-icon>
        <h1>智简魔方</h1>
        <p>管理后台</p>
      </div>

      <el-form ref="formRef" :model="formData" :rules="rules" size="large">
        <el-form-item prop="username">
          <el-input v-model="formData.username" placeholder="请输入管理员用户名" @keydown.enter="handleLogin">
            <template #prefix><el-icon color="#0056FF"><User /></el-icon></template>
          </el-input>
        </el-form-item>

        <el-form-item prop="password">
          <el-input v-model="formData.password" type="password" show-password placeholder="请输入密码" @keydown.enter="handleLogin">
            <template #prefix><el-icon color="#0056FF"><Lock /></el-icon></template>
          </el-input>
        </el-form-item>

        <el-form-item prop="captcha">
          <div class="captcha-row">
            <el-input v-model="formData.captcha" placeholder="请输入验证码" @keydown.enter="handleLogin">
              <template #prefix><el-icon color="#0056FF"><Shield /></el-icon></template>
            </el-input>
            <div class="captcha-img" @click="refreshCaptcha">
              <img :src="captchaUrl" alt="验证码" />
            </div>
          </div>
        </el-form-item>

        <el-form-item>
          <div class="login-options">
            <el-checkbox v-model="formData.remember">记住登录状态</el-checkbox>
          </div>
        </el-form-item>

        <el-button type="primary" class="login-btn" size="large" :loading="loading" @click="handleLogin">
          登 录
        </el-button>
      </el-form>

      <div class="login-footer">
        <el-button text type="primary" @click="goToFrontend">
          <el-icon><Back /></el-icon>返回前台
        </el-button>
      </div>
    </div>

    <div class="login-copyright">
      &copy; {{ new Date().getFullYear() }} 智简魔方 All Rights Reserved
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock, Shield, Back, Connection } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'

const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)

const captchaUrl = ref(`/api/captcha?t=${Date.now()}`)

const formData = reactive({ username: '', password: '', captcha: '', remember: true })

const rules: FormRules = {
  username: { required: true, message: '请输入用户名', trigger: 'blur' },
  password: { required: true, message: '请输入密码', trigger: 'blur' },
  captcha: { required: true, message: '请输入验证码', trigger: 'blur' },
}

function refreshCaptcha() {
  captchaUrl.value = `/api/captcha?t=${Date.now()}`
}

async function handleLogin() {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
  } catch {
    return
  }

  loading.value = true
  try {
    ElMessage.success('登录成功')
    router.push('/admin/dashboard')
  } catch (err: any) {
    ElMessage.error(err?.response?.data?.message || '登录失败，请检查用户名和密码')
    refreshCaptcha()
  } finally {
    loading.value = false
  }
}

function goToFrontend() {
  router.push('/')
}
</script>

<style scoped>
.login-container {
  width: 100%;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #0b1a2e 0%, #0d2744 30%, #103a5c 60%, #0a1f35 100%);
  position: relative;
  overflow: hidden;
}

.login-container::before {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle at 30% 40%, rgba(0, 86, 255, 0.08) 0%, transparent 50%),
              radial-gradient(circle at 70% 60%, rgba(0, 86, 255, 0.05) 0%, transparent 40%);
  animation: floatBg 20s ease-in-out infinite;
}

@keyframes floatBg {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(-2%, -1%); }
}

.login-card {
  width: 420px;
  padding: 48px 40px 36px;
  background: #ffffff;
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.4), 0 0 40px rgba(0, 86, 255, 0.06);
  position: relative;
  z-index: 1;
}

.login-header {
  text-align: center;
  margin-bottom: 36px;
}

.login-header h1 {
  margin: 16px 0 4px;
  font-size: 26px;
  font-weight: 700;
  color: #1a1a2e;
  letter-spacing: 2px;
}

.login-header p {
  color: #8896ab;
  font-size: 14px;
  margin: 0;
}

.captcha-row {
  display: flex;
  gap: 12px;
  width: 100%;
}

.captcha-row .el-input {
  flex: 1;
}

.captcha-img {
  width: 120px;
  height: 40px;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  flex-shrink: 0;
  border: 1px solid #e0e5ec;
  transition: border-color 0.2s;
}

.captcha-img:hover {
  border-color: #0056FF;
}

.captcha-img img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.login-options {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.login-btn {
  width: 100%;
  height: 44px;
  font-size: 16px;
  font-weight: 600;
  border-radius: 10px;
  background: linear-gradient(135deg, #0056FF 0%, #4080FF 100%);
  border: none;
  letter-spacing: 4px;
}

.login-btn:hover {
  opacity: 0.9;
}

.login-footer {
  text-align: center;
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid #f0f2f5;
}

.login-copyright {
  position: relative;
  z-index: 1;
  margin-top: 32px;
  color: rgba(255, 255, 255, 0.3);
  font-size: 12px;
}
</style>
