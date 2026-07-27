<template>
  <div class="login-page">
    <div class="login-bg">
      <div class="bg-shape bg-shape-1"></div>
      <div class="bg-shape bg-shape-2"></div>
      <div class="bg-shape bg-shape-3"></div>
    </div>

    <div class="login-card">
      <div class="login-header">
        <div class="logo-wrapper">
          <div class="logo-icon">
            <svg viewBox="0 0 24 24" width="32" height="32" fill="none">
              <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5" stroke="#1890ff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </div>
          <h1 class="logo-text">锚点财务</h1>
        </div>
        <p class="login-subtitle">欢迎回来，请登录您的账户</p>
      </div>

      <n-tabs v-model:value="activeTab" type="segment" animated class="login-tabs">
        <n-tab-pane name="password" tab="密码登录">
          <n-form ref="passwordFormRef" :model="passwordForm" :rules="passwordRules" class="login-form">
            <n-form-item path="username">
              <n-input
                v-model:value="passwordForm.username"
                placeholder="请输入用户名/邮箱"
                size="large"
                :input-props="{ autocomplete: 'username' }"
              >
                <template #prefix>
                  <n-icon :component="PersonOutline" color="#1890ff" />
                </template>
              </n-input>
            </n-form-item>

            <n-form-item path="password">
              <n-input
                v-model:value="passwordForm.password"
                type="password"
                show-password-on="click"
                placeholder="请输入密码"
                size="large"
                :input-props="{ autocomplete: 'current-password' }"
              >
                <template #prefix>
                  <n-icon :component="LockClosedOutline" color="#1890ff" />
                </template>
              </n-input>
            </n-form-item>

            <n-form-item path="captcha">
              <div class="captcha-row">
                <n-input
                  v-model:value="passwordForm.captcha"
                  placeholder="请输入验证码"
                  size="large"
                  @keyup.enter="handlePasswordLogin"
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

            <div class="login-options">
              <n-checkbox v-model:checked="passwordForm.remember">记住我</n-checkbox>
              <router-link to="/forgot-password" class="link-text">忘记密码?</router-link>
            </div>

            <n-button
              type="primary"
              block
              size="large"
              :loading="loading"
              class="login-btn"
              @click="handlePasswordLogin"
            >
              登录
            </n-button>
          </n-form>
        </n-tab-pane>

        <n-tab-pane name="sms" tab="验证码登录">
          <n-form ref="smsFormRef" :model="smsForm" :rules="smsRules" class="login-form">
            <n-form-item path="phone">
              <n-input
                v-model:value="smsForm.phone"
                placeholder="请输入手机号"
                size="large"
              >
                <template #prefix>
                  <n-icon :component="CallOutline" color="#1890ff" />
                </template>
              </n-input>
            </n-form-item>

            <n-form-item path="imageCaptcha">
              <div class="captcha-row">
                <n-input
                  v-model:value="smsForm.imageCaptcha"
                  placeholder="请输入图形验证码"
                  size="large"
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

            <n-form-item path="smsCode">
              <div class="captcha-row">
                <n-input
                  v-model:value="smsForm.smsCode"
                  placeholder="请输入短信验证码"
                  size="large"
                  @keyup.enter="handleSmsLogin"
                >
                  <template #prefix>
                    <n-icon :component="ChatboxOutline" color="#1890ff" />
                  </template>
                </n-input>
                <n-button
                  size="large"
                  :disabled="smsCooldown > 0"
                  :loading="sendingSms"
                  class="sms-btn"
                  @click="handleSendSms"
                >
                  {{ smsCooldown > 0 ? `${smsCooldown}s` : '获取验证码' }}
                </n-button>
              </div>
            </n-form-item>

            <n-button
              type="primary"
              block
              size="large"
              :loading="loading"
              class="login-btn"
              @click="handleSmsLogin"
            >
              登录
            </n-button>
          </n-form>
        </n-tab-pane>
      </n-tabs>

      <div class="login-divider">
        <span>其他登录方式</span>
      </div>

      <div class="third-party-login">
        <n-tooltip trigger="hover">
          <template #trigger>
            <n-button circle size="large" class="third-party-btn wechat" @click="handleThirdParty('wechat')">
              <template #icon>
                <svg viewBox="0 0 24 24" width="20" height="20" fill="#07c160">
                  <path d="M8.691 2.188C3.891 2.188 0 5.476 0 9.53c0 2.212 1.17 4.203 3.002 5.55a.59.59 0 0 1 .213.665l-.39 1.48c-.019.07-.048.141-.048.213 0 .163.13.295.29.295a.326.326 0 0 0 .167-.054l1.903-1.114a.864.864 0 0 1 .717-.098 10.16 10.16 0 0 0 2.837.403c.276 0 .543-.027.811-.05-.857-2.578.157-4.972 1.932-6.446 1.703-1.415 3.882-1.98 5.853-1.838-.576-3.583-4.196-6.348-8.596-6.348zM5.785 5.991c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 0 1-1.162 1.178A1.17 1.17 0 0 1 4.623 7.17c0-.651.52-1.18 1.162-1.18zm5.813 0c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 0 1-1.162 1.178 1.17 1.17 0 0 1-1.162-1.178c0-.651.52-1.18 1.162-1.18zm5.34 2.867c-1.797-.052-3.746.512-5.28 1.786-1.72 1.428-2.687 3.72-1.78 6.22.942 2.453 3.666 4.229 6.884 4.229.826 0 1.622-.12 2.361-.336a.722.722 0 0 1 .598.082l1.584.926a.272.272 0 0 0 .14.047c.134 0 .24-.111.24-.247 0-.06-.023-.12-.038-.177l-.327-1.233a.582.582 0 0 1-.023-.156.49.49 0 0 1 .201-.398C23.024 18.48 24 16.82 24 14.98c0-3.21-2.931-5.837-6.656-6.088V8.89c-.135-.01-.27-.027-.407-.032zm-2.53 3.274c.535 0 .969.44.969.982a.976.976 0 0 1-.969.983.976.976 0 0 1-.969-.983c0-.542.434-.982.97-.982zm4.844 0c.535 0 .969.44.969.982a.976.976 0 0 1-.969.983.976.976 0 0 1-.969-.983c0-.542.434-.982.97-.982z"/>
                </svg>
              </template>
            </n-button>
          </template>
          微信登录
        </n-tooltip>
        <n-tooltip trigger="hover">
          <template #trigger>
            <n-button circle size="large" class="third-party-btn qq" @click="handleThirdParty('qq')">
              <template #icon>
                <svg viewBox="0 0 24 24" width="20" height="20" fill="#12b7f5">
                  <path d="M12.003 2c-2.265 0-6.29 1.364-6.29 7.325v1.195S3.55 14.96 3.55 17.474c0 .665.17 1.025.281 1.025.114 0 .902-.484 1.748-2.072 0 0-.18 2.197 1.904 3.967 0 0-1.77.495-1.77 1.182 0 .686 4.078.43 6.29.43 2.239 0 6.29.256 6.29-.43 0-.687-1.77-1.182-1.77-1.182 2.085-1.77 1.905-3.967 1.905-3.967.845 1.588 1.634 2.072 1.746 2.072.111 0 .283-.36.283-1.025 0-2.514-2.166-6.954-2.166-6.954V9.325C18.29 3.364 14.268 2 12.003 2z"/>
                </svg>
              </template>
            </n-button>
          </template>
          QQ登录
        </n-tooltip>
        <n-tooltip trigger="hover">
          <template #trigger>
            <n-button circle size="large" class="third-party-btn github" @click="handleThirdParty('github')">
              <template #icon>
                <svg viewBox="0 0 24 24" width="20" height="20" fill="#333">
                  <path d="M12 .297c-6.63 0-12 5.373-12 12 0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61C4.422 18.07 3.633 17.7 3.633 17.7c-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12"/>
                </svg>
              </template>
            </n-button>
          </template>
          GitHub登录
        </n-tooltip>
        <n-tooltip trigger="hover">
          <template #trigger>
            <n-button circle size="large" class="third-party-btn google" @click="handleThirdParty('google')">
              <template #icon>
                <svg viewBox="0 0 24 24" width="20" height="20">
                  <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 0 1-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z" fill="#4285F4"/>
                  <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/>
                  <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05"/>
                  <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/>
                </svg>
              </template>
            </n-button>
          </template>
          Google登录
        </n-tooltip>
      </div>

      <div class="login-footer">
        还没有账号？<router-link to="/register" class="link-text">立即注册</router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useMessage } from 'naive-ui'
import type { FormInst, FormRules } from 'naive-ui'
import {
  PersonOutline,
  LockClosedOutline,
  ImageOutline,
  RefreshOutline,
  CallOutline,
  ChatboxOutline
} from '@vicons/ionicons5'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const route = useRoute()
const message = useMessage()
const userStore = useUserStore()

const activeTab = ref('password')
const loading = ref(false)
const sendingSms = ref(false)
const captchaUrl = ref('')
const smsCooldown = ref(0)
let cooldownTimer: ReturnType<typeof setInterval> | null = null

const passwordFormRef = ref<FormInst | null>(null)
const smsFormRef = ref<FormInst | null>(null)

const passwordForm = ref({
  username: '',
  password: '',
  captcha: '',
  remember: false
})

const smsForm = ref({
  phone: '',
  imageCaptcha: '',
  smsCode: ''
})

const passwordRules: FormRules = {
  username: { required: true, message: '请输入用户名', trigger: 'blur' },
  password: { required: true, message: '请输入密码', trigger: 'blur' },
  captcha: { required: true, message: '请输入验证码', trigger: 'blur' }
}

const smsRules: FormRules = {
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ],
  imageCaptcha: { required: true, message: '请输入图形验证码', trigger: 'blur' },
  smsCode: { required: true, message: '请输入短信验证码', trigger: 'blur' }
}

function refreshCaptcha() {
  captchaUrl.value = `/api/captcha?t=${Date.now()}`
}

onMounted(() => {
  refreshCaptcha()
})

onUnmounted(() => {
  if (cooldownTimer) clearInterval(cooldownTimer)
})

async function handlePasswordLogin() {
  try {
    await passwordFormRef.value?.validate()
    loading.value = true
    await userStore.login(passwordForm.value.username, passwordForm.value.password, passwordForm.value.captcha)
    message.success('登录成功')
    const redirect = (route.query.redirect as string) || '/user/dashboard'
    router.push(redirect)
  } catch (error: any) {
    message.error(error.message || '登录失败')
    refreshCaptcha()
  } finally {
    loading.value = false
  }
}

async function handleSendSms() {
  try {
    await smsFormRef.value?.validate(['phone', 'imageCaptcha'])
    sendingSms.value = true
    // TODO: 调用发送短信API
    await new Promise(resolve => setTimeout(resolve, 1000))
    message.success('验证码已发送')
    smsCooldown.value = 60
    cooldownTimer = setInterval(() => {
      smsCooldown.value--
      if (smsCooldown.value <= 0) {
        if (cooldownTimer) clearInterval(cooldownTimer)
      }
    }, 1000)
  } catch {
    message.error('请先填写手机号和图形验证码')
  } finally {
    sendingSms.value = false
  }
}

async function handleSmsLogin() {
  try {
    await smsFormRef.value?.validate()
    loading.value = true
    // TODO: 调用短信登录API
    await new Promise(resolve => setTimeout(resolve, 1000))
    message.success('登录成功')
    const redirect = (route.query.redirect as string) || '/user/dashboard'
    router.push(redirect)
  } catch (error: any) {
    message.error(error.message || '登录失败')
  } finally {
    loading.value = false
  }
}

function handleThirdParty(platform: string) {
  message.info(`${platform} 登录功能开发中`)
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #e8f4fd 0%, #f0f7ff 40%, #ffffff 100%);
  position: relative;
  overflow: hidden;
}

.login-bg {
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

.login-card {
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

.login-card:hover {
  box-shadow: 0 24px 64px rgba(24, 144, 255, 0.18), 0 0 0 1px rgba(255, 255, 255, 0.9);
}

.login-header {
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

.login-subtitle {
  color: #8c8c8c;
  font-size: 14px;
  margin: 0;
}

.login-tabs {
  margin-bottom: 20px;
}

.login-form {
  margin-top: 20px;
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

.sms-btn {
  width: 120px;
  flex-shrink: 0;
}

.login-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.link-text {
  color: #1890ff;
  text-decoration: none;
  font-size: 14px;
  transition: color 0.3s;
}

.link-text:hover {
  color: #40a9ff;
}

.login-btn {
  height: 44px;
  font-size: 16px;
  font-weight: 500;
  border-radius: 12px;
  background: linear-gradient(135deg, #1890ff, #096dd9);
  border: none;
  transition: all 0.3s;
}

.login-btn:hover {
  background: linear-gradient(135deg, #40a9ff, #1890ff);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(24, 144, 255, 0.4);
}

.login-divider {
  display: flex;
  align-items: center;
  margin: 24px 0;
  color: #bfbfbf;
  font-size: 13px;
}

.login-divider::before,
.login-divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: linear-gradient(90deg, transparent, #e8e8e8, transparent);
}

.login-divider span {
  padding: 0 16px;
}

.third-party-login {
  display: flex;
  justify-content: center;
  gap: 16px;
  margin-bottom: 24px;
}

.third-party-btn {
  width: 44px;
  height: 44px;
  border: 1px solid #e8e8e8;
  background: #fff;
  transition: all 0.3s;
}

.third-party-btn.wechat:hover {
  border-color: #07c160;
  background: #f0fff8;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(7, 193, 96, 0.2);
}

.third-party-btn.qq:hover {
  border-color: #12b7f5;
  background: #f0f9ff;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(18, 183, 245, 0.2);
}

.third-party-btn.github:hover {
  border-color: #333;
  background: #f5f5f5;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.third-party-btn.google:hover {
  border-color: #4285F4;
  background: #f0f4ff;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(66, 133, 244, 0.2);
}

.login-footer {
  text-align: center;
  color: #8c8c8c;
  font-size: 14px;
}
</style>
