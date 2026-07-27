<template>
  <div class="forgot-page">
    <div class="forgot-bg">
      <div class="bg-shape bg-shape-1"></div>
      <div class="bg-shape bg-shape-2"></div>
      <div class="bg-shape bg-shape-3"></div>
    </div>

    <div class="forgot-card">
      <div class="forgot-header">
        <div class="logo-wrapper">
          <div class="logo-icon">
            <svg viewBox="0 0 24 24" width="32" height="32" fill="none">
              <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5" stroke="#1890ff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </div>
          <h1 class="logo-text">找回密码</h1>
        </div>
        <p class="forgot-subtitle">请按照步骤重置您的密码</p>
      </div>

      <n-steps :current="currentStep" :status="stepStatus" class="steps-bar">
        <n-step title="输入账号" />
        <n-step title="验证身份" />
        <n-step title="设置新密码" />
        <n-step title="完成" />
      </n-steps>

      <!-- Step 1: 输入账号 -->
      <div v-if="currentStep === 1" class="step-content">
        <n-form ref="accountFormRef" :model="accountForm" :rules="accountRules" class="forgot-form">
          <n-form-item path="account" label="邮箱或手机号">
            <n-input
              v-model:value="accountForm.account"
              placeholder="请输入邮箱或手机号"
              size="large"
              @keyup.enter="handleNextStep1"
            >
              <template #prefix>
                <n-icon :component="PersonOutline" color="#1890ff" />
              </template>
            </n-input>
          </n-form-item>

          <n-form-item path="captcha">
            <div class="captcha-row">
              <n-input
                v-model:value="accountForm.captcha"
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

          <n-button
            type="primary"
            block
            size="large"
            :loading="loading"
            class="submit-btn"
            @click="handleNextStep1"
          >
            下一步
          </n-button>
        </n-form>
      </div>

      <!-- Step 2: 验证身份 -->
      <div v-if="currentStep === 2" class="step-content">
        <n-form ref="verifyFormRef" :model="verifyForm" :rules="verifyRules" class="forgot-form">
          <div class="verify-hint">
            验证码已发送至 <span class="highlight">{{ maskedAccount }}</span>
          </div>

          <n-form-item path="verifyCode" label="验证码">
            <div class="captcha-row">
              <n-input
                v-model:value="verifyForm.verifyCode"
                placeholder="请输入验证码"
                size="large"
                @keyup.enter="handleNextStep2"
              >
                <template #prefix>
                  <n-icon :component="ShieldCheckmarkOutline" color="#1890ff" />
                </template>
              </n-input>
              <n-button
                size="large"
                :disabled="cooldown > 0"
                :loading="sendingCode"
                class="sms-btn"
                @click="handleResendCode"
              >
                {{ cooldown > 0 ? `${cooldown}s` : '重新发送' }}
              </n-button>
            </div>
          </n-form-item>

          <div class="step-actions">
            <n-button size="large" @click="currentStep = 1">上一步</n-button>
            <n-button
              type="primary"
              size="large"
              :loading="loading"
              @click="handleNextStep2"
            >
              验证
            </n-button>
          </div>
        </n-form>
      </div>

      <!-- Step 3: 设置新密码 -->
      <div v-if="currentStep === 3" class="step-content">
        <n-form ref="passwordFormRef" :model="passwordForm" :rules="passwordRules" class="forgot-form">
          <n-form-item path="newPassword" label="新密码">
            <n-input
              v-model:value="passwordForm.newPassword"
              type="password"
              show-password-on="click"
              placeholder="请输入新密码（至少6位）"
              size="large"
            >
              <template #prefix>
                <n-icon :component="LockClosedOutline" color="#1890ff" />
              </template>
            </n-input>
          </n-form-item>

          <n-form-item path="confirmPassword" label="确认密码">
            <n-input
              v-model:value="passwordForm.confirmPassword"
              type="password"
              show-password-on="click"
              placeholder="请再次输入新密码"
              size="large"
              @keyup.enter="handleNextStep3"
            >
              <template #prefix>
                <n-icon :component="LockClosedOutline" color="#1890ff" />
              </template>
            </n-input>
          </n-form-item>

          <div class="step-actions">
            <n-button size="large" @click="currentStep = 2">上一步</n-button>
            <n-button
              type="primary"
              size="large"
              :loading="loading"
              @click="handleNextStep3"
            >
              重置密码
            </n-button>
          </div>
        </n-form>
      </div>

      <!-- Step 4: 完成 -->
      <div v-if="currentStep === 4" class="step-content step-result">
        <n-result status="success" title="密码重置成功" description="您的密码已成功重置，请使用新密码登录">
          <template #footer>
            <n-button type="primary" size="large" @click="goToLogin">
              前往登录
            </n-button>
          </template>
        </n-result>
      </div>

      <div class="forgot-footer">
        想起密码了？<router-link to="/login" class="link-text">返回登录</router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import type { FormInst, FormRules, StepsProps } from 'naive-ui'
import {
  PersonOutline,
  ImageOutline,
  RefreshOutline,
  LockClosedOutline,
  ShieldCheckmarkOutline
} from '@vicons/ionicons5'

const router = useRouter()
const message = useMessage()

const currentStep = ref(1)
const stepStatus = ref<StepsProps['status']>('process')
const loading = ref(false)
const sendingCode = ref(false)
const cooldown = ref(0)
const captchaUrl = ref('')
let cooldownTimer: ReturnType<typeof setInterval> | null = null

const accountFormRef = ref<FormInst | null>(null)
const verifyFormRef = ref<FormInst | null>(null)
const passwordFormRef = ref<FormInst | null>(null)

const accountForm = ref({
  account: '',
  captcha: ''
})

const verifyForm = ref({
  verifyCode: ''
})

const passwordForm = ref({
  newPassword: '',
  confirmPassword: ''
})

const accountRules: FormRules = {
  account: [
    { required: true, message: '请输入邮箱或手机号', trigger: 'blur' }
  ],
  captcha: { required: true, message: '请输入验证码', trigger: 'blur' }
}

const verifyRules: FormRules = {
  verifyCode: { required: true, message: '请输入验证码', trigger: 'blur' }
}

const passwordRules: FormRules = {
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    {
      validator: (_rule: any, value: string) => {
        if (value !== passwordForm.value.newPassword) {
          return new Error('两次密码输入不一致')
        }
        return true
      },
      trigger: 'blur'
    }
  ]
}

const maskedAccount = computed(() => {
  const account = accountForm.value.account
  if (!account) return ''
  if (account.includes('@')) {
    const [name, domain] = account.split('@')
    return name.slice(0, 2) + '***@' + domain
  }
  return account.slice(0, 3) + '****' + account.slice(-4)
})

function refreshCaptcha() {
  captchaUrl.value = `/api/captcha?t=${Date.now()}`
}

function startCooldown() {
  cooldown.value = 60
  cooldownTimer = setInterval(() => {
    cooldown.value--
    if (cooldown.value <= 0) {
      if (cooldownTimer) clearInterval(cooldownTimer)
    }
  }, 1000)
}

onMounted(() => {
  refreshCaptcha()
})

onUnmounted(() => {
  if (cooldownTimer) clearInterval(cooldownTimer)
})

async function handleNextStep1() {
  try {
    await accountFormRef.value?.validate()
    loading.value = true
    // TODO: 调用发送验证码API
    await new Promise(resolve => setTimeout(resolve, 1000))
    message.success('验证码已发送')
    currentStep.value = 2
    startCooldown()
  } catch (error: any) {
    if (error?.message) message.error(error.message)
  } finally {
    loading.value = false
  }
}

async function handleResendCode() {
  try {
    sendingCode.value = true
    // TODO: 调用重新发送验证码API
    await new Promise(resolve => setTimeout(resolve, 1000))
    message.success('验证码已重新发送')
    startCooldown()
  } catch {
    message.error('发送失败，请稍后重试')
  } finally {
    sendingCode.value = false
  }
}

async function handleNextStep2() {
  try {
    await verifyFormRef.value?.validate()
    loading.value = true
    // TODO: 调用验证验证码API
    await new Promise(resolve => setTimeout(resolve, 1000))
    message.success('验证成功')
    currentStep.value = 3
  } catch {
    message.error('验证码错误')
  } finally {
    loading.value = false
  }
}

async function handleNextStep3() {
  try {
    await passwordFormRef.value?.validate()
    loading.value = true
    // TODO: 调用重置密码API
    await new Promise(resolve => setTimeout(resolve, 1000))
    message.success('密码重置成功')
    stepStatus.value = 'finish'
    currentStep.value = 4
  } catch (error: any) {
    if (error?.message) message.error(error.message)
  } finally {
    loading.value = false
  }
}

function goToLogin() {
  router.push('/login')
}
</script>

<style scoped>
.forgot-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #e8f4fd 0%, #f0f7ff 40%, #ffffff 100%);
  position: relative;
  overflow: hidden;
}

.forgot-bg {
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

.forgot-card {
  width: 480px;
  padding: 40px;
  background: rgba(255, 255, 255, 0.96);
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(24, 144, 255, 0.12), 0 0 0 1px rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(20px);
  position: relative;
  z-index: 1;
  transition: box-shadow 0.3s ease;
}

.forgot-card:hover {
  box-shadow: 0 24px 64px rgba(24, 144, 255, 0.18), 0 0 0 1px rgba(255, 255, 255, 0.9);
}

.forgot-header {
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

.forgot-subtitle {
  color: #8c8c8c;
  font-size: 14px;
  margin: 0;
}

.steps-bar {
  margin-bottom: 32px;
}

.step-content {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.forgot-form {
  margin-top: 8px;
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

.verify-hint {
  text-align: center;
  color: #595959;
  font-size: 14px;
  margin-bottom: 20px;
  padding: 12px;
  background: #f0f7ff;
  border-radius: 8px;
}

.verify-hint .highlight {
  color: #1890ff;
  font-weight: 500;
}

.step-actions {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  margin-top: 8px;
}

.step-actions .n-button {
  flex: 1;
  height: 44px;
  font-size: 15px;
  border-radius: 12px;
}

.step-actions .n-button[type="button"]:last-child {
  background: linear-gradient(135deg, #1890ff, #096dd9);
  border: none;
}

.step-result {
  padding: 24px 0;
}

.submit-btn {
  height: 44px;
  font-size: 16px;
  font-weight: 500;
  border-radius: 12px;
  background: linear-gradient(135deg, #1890ff, #096dd9);
  border: none;
  transition: all 0.3s;
}

.submit-btn:hover {
  background: linear-gradient(135deg, #40a9ff, #1890ff);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(24, 144, 255, 0.4);
}

.forgot-footer {
  text-align: center;
  color: #8c8c8c;
  font-size: 14px;
  margin-top: 24px;
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
</style>
