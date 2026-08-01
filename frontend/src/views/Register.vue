<template>
  <div class="register-page">
    <div class="register-bg">
      <div class="bg-shape bg-shape-1"></div>
      <div class="bg-shape bg-shape-2"></div>
      <div class="bg-shape bg-shape-3"></div>
    </div>

    <div class="register-card">
      <div class="register-header">
        <div class="logo-wrapper">
          <div class="logo-icon">
            <svg viewBox="0 0 24 24" width="32" height="32" fill="none">
              <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5" stroke="#1890ff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </div>
          <h1 class="logo-text">锚点财务</h1>
        </div>
        <p class="register-subtitle">{{ $t('auth.registerTitle') }}</p>
      </div>

      <n-form ref="formRef" :model="form" :rules="rules" class="register-form">
        <n-form-item path="username">
          <n-input
            v-model:value="form.username"
            placeholder="请输入用户名"
            size="large"
          >
            <template #prefix>
              <n-icon :component="PersonOutline" color="#1890ff" />
            </template>
          </n-input>
        </n-form-item>

        <n-form-item path="email">
          <n-input
            v-model:value="form.email"
            placeholder="请输入邮箱"
            size="large"
          >
            <template #prefix>
              <n-icon :component="MailOutline" color="#1890ff" />
            </template>
          </n-input>
        </n-form-item>

        <n-form-item path="password">
          <div class="password-field">
            <n-input
              v-model:value="form.password"
              type="password"
              show-password-on="click"
              placeholder="请输入密码（至少6位）"
              size="large"
              @input="checkPasswordStrength"
            >
              <template #prefix>
                <n-icon :component="LockClosedOutline" color="#1890ff" />
              </template>
            </n-input>
            <div class="password-strength" v-if="form.password">
              <div class="strength-bars">
                <div
                  class="strength-bar"
                  :class="passwordStrength >= 1 ? strengthLevel : ''"
                ></div>
                <div
                  class="strength-bar"
                  :class="passwordStrength >= 2 ? strengthLevel : ''"
                ></div>
                <div
                  class="strength-bar"
                  :class="passwordStrength >= 3 ? strengthLevel : ''"
                ></div>
              </div>
              <span class="strength-text" :class="strengthLevel">
                {{ strengthLabel }}
              </span>
            </div>
          </div>
        </n-form-item>

        <n-form-item path="confirmPassword">
          <n-input
            v-model:value="form.confirmPassword"
            type="password"
            show-password-on="click"
            placeholder="请再次输入密码"
            size="large"
          >
            <template #prefix>
              <n-icon :component="LockClosedOutline" color="#1890ff" />
            </template>
          </n-input>
        </n-form-item>

        <n-form-item path="captcha">
          <div class="captcha-row">
            <n-input
              v-model:value="form.captcha"
              placeholder="请输入验证码"
              size="large"
              @keyup.enter="handleRegister"
            >
              <template #prefix>
                <n-icon :component="ImageOutline" color="#1890ff" />
              </template>
            </n-input>
            <div class="captcha-image" @click="refreshCaptcha">
              <img v-if="captchaUrl" :src="captchaUrl" alt="验证码" />
              <div v-else class="captcha-placeholder">
                <n-icon :component="RefreshOutline" />
              </div>
            </div>
          </div>
        </n-form-item>

        <div class="agreement">
          <n-checkbox v-model:checked="agreed">
            我已阅读并同意
            <a href="#" class="link-text" @click.prevent>《服务条款》</a>
          </n-checkbox>
        </div>

        <n-button
          type="primary"
          block
          size="large"
          :loading="loading"
          :disabled="!agreed"
          class="register-btn"
          @click="handleRegister"
        >
          {{ $t('auth.register') }}
        </n-button>
      </n-form>

      <div class="register-lang">
        <LanguageSwitch />
      </div>

      <div class="register-footer">
        {{ $t('auth.hasAccount') }}<router-link to="/login" class="link-text">{{ $t('auth.loginNow') }}</router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import type { FormInst, FormRules } from 'naive-ui'
import {
  PersonOutline,
  LockClosedOutline,
  MailOutline,
  ImageOutline,
  RefreshOutline
} from '@vicons/ionicons5'
import { useUserStore } from '@/stores/user'
import LanguageSwitch from '@/components/LanguageSwitch.vue'

const router = useRouter()
const message = useMessage()
const userStore = useUserStore()

const formRef = ref<FormInst | null>(null)
const loading = ref(false)
const captchaUrl = ref('')
const agreed = ref(false)
const passwordStrength = ref(0)

const form = ref({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  captcha: ''
})

const strengthLevel = computed(() => {
  if (passwordStrength.value <= 1) return 'weak'
  if (passwordStrength.value === 2) return 'medium'
  return 'strong'
})

const strengthLabel = computed(() => {
  if (passwordStrength.value <= 1) return '弱'
  if (passwordStrength.value === 2) return '中'
  return '强'
})

function checkPasswordStrength() {
  const pwd = form.value.password
  if (!pwd) {
    passwordStrength.value = 0
    return
  }

  let score = 0
  if (pwd.length >= 6) score++
  if (pwd.length >= 10) score++
  if (/[a-z]/.test(pwd) && /[A-Z]/.test(pwd)) score++
  if (/\d/.test(pwd)) score++
  if (/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(pwd)) score++

  if (score <= 2) passwordStrength.value = 1
  else if (score <= 3) passwordStrength.value = 2
  else passwordStrength.value = 3
}

const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度为3-20个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    {
      validator: (_rule: any, value: string) => {
        return value === form.value.password
      },
      message: '两次输入的密码不一致',
      trigger: 'blur'
    }
  ],
  captcha: { required: true, message: '请输入验证码', trigger: 'blur' }
}

function refreshCaptcha() {
  captchaUrl.value = `/api/captcha?t=${Date.now()}`
}

onMounted(() => {
  refreshCaptcha()
})

async function handleRegister() {
  if (!agreed.value) {
    message.warning('请先阅读并同意服务条款')
    return
  }
  try {
    await formRef.value?.validate()
    loading.value = true
    await userStore.register({
      username: form.value.username,
      email: form.value.email,
      password: form.value.password,
      captcha: form.value.captcha
    })
    message.success('注册成功，请登录')
    router.push('/login')
  } catch (error: any) {
    message.error(error.message || '注册失败')
    refreshCaptcha()
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #e8f4fd 0%, #f0f7ff 40%, #ffffff 100%);
  position: relative;
  overflow: hidden;
  padding: 40px 0;
}

.register-bg {
  position: absolute;
  inset: 0;
  pointer-events: none;
}

.bg-shape {
  position: absolute;
  border-radius: 50%;
  opacity: 0.12;
}

.bg-shape-1 {
  width: 600px;
  height: 600px;
  background: linear-gradient(135deg, #1890ff, #40a9ff);
  top: -200px;
  right: -100px;
  animation: float1 20s ease-in-out infinite;
}

.bg-shape-2 {
  width: 400px;
  height: 400px;
  background: linear-gradient(135deg, #096dd9, #1890ff);
  bottom: -150px;
  left: -100px;
  animation: float2 25s ease-in-out infinite;
}

.bg-shape-3 {
  width: 200px;
  height: 200px;
  background: linear-gradient(135deg, #40a9ff, #69c0ff);
  top: 50%;
  left: 10%;
  animation: float3 18s ease-in-out infinite;
}

@keyframes float1 {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(-30px, 30px); }
}

@keyframes float2 {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(20px, -20px); }
}

@keyframes float3 {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(-15px, -25px); }
}

.register-card {
  width: 420px;
  padding: 40px;
  background: rgba(255, 255, 255, 0.96);
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(24, 144, 255, 0.12), 0 0 0 1px rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(20px);
  position: relative;
  z-index: 1;
  transition: box-shadow 0.3s ease;
}

.register-card:hover {
  box-shadow: 0 24px 64px rgba(24, 144, 255, 0.18), 0 0 0 1px rgba(255, 255, 255, 0.9);
}

.register-header {
  text-align: center;
  margin-bottom: 28px;
}

.logo-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.logo-icon {
  width: 56px;
  height: 56px;
  background: linear-gradient(135deg, #e8f4fd, #d6e8fa);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-text {
  font-size: 26px;
  font-weight: 700;
  background: linear-gradient(135deg, #1890ff, #096dd9);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.register-subtitle {
  color: #8c8c8c;
  font-size: 14px;
  margin: 0;
}

.register-form {
  margin-bottom: 8px;
}

.password-field {
  width: 100%;
}

.password-strength {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 8px;
}

.strength-bars {
  display: flex;
  gap: 4px;
  flex: 1;
}

.strength-bar {
  height: 4px;
  flex: 1;
  border-radius: 2px;
  background: #e8e8e8;
  transition: background 0.3s ease;
}

.strength-bar.weak {
  background: #ff4d4f;
}

.strength-bar.medium {
  background: #faad14;
}

.strength-bar.strong {
  background: #52c41a;
}

.strength-text {
  font-size: 12px;
  font-weight: 500;
  min-width: 20px;
  transition: color 0.3s ease;
}

.strength-text.weak {
  color: #ff4d4f;
}

.strength-text.medium {
  color: #faad14;
}

.strength-text.strong {
  color: #52c41a;
}

.captcha-row {
  display: flex;
  gap: 12px;
  width: 100%;
}

.captcha-image {
  width: 120px;
  height: 40px;
  cursor: pointer;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #e8e8e8;
  flex-shrink: 0;
  transition: border-color 0.3s;
}

.captcha-image:hover {
  border-color: #1890ff;
}

.captcha-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.captcha-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f0f5ff;
  color: #1890ff;
}

.agreement {
  margin-bottom: 20px;
}

.link-text {
  color: #1890ff;
  text-decoration: none;
  transition: color 0.3s;
}

.link-text:hover {
  color: #40a9ff;
}

.register-btn {
  height: 44px;
  font-size: 16px;
  font-weight: 500;
  border-radius: 12px;
  background: linear-gradient(135deg, #1890ff, #096dd9);
  border: none;
  transition: all 0.3s;
}

.register-btn:hover:not(:disabled) {
  background: linear-gradient(135deg, #40a9ff, #1890ff);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(24, 144, 255, 0.4);
}

.register-lang {
  display: flex;
  justify-content: center;
  margin-bottom: 16px;
}

.register-footer {
  text-align: center;
  margin-top: 24px;
  color: #8c8c8c;
  font-size: 14px;
}
</style>
