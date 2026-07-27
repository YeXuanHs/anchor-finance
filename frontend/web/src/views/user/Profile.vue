<template>
  <div class="profile-page">
    <!-- Avatar Section -->
    <n-card :bordered="false" class="page-card">
      <div class="profile-header">
        <div class="avatar-section">
          <n-avatar :size="96" :src="avatarUrl" round class="user-avatar" />
          <n-upload
            :show-file-list="false"
            :max="1"
            accept="image/*"
            @change="handleAvatarUpload"
          >
            <n-button size="small" quaternary class="avatar-edit-btn">
              <template #icon>
                <n-icon :component="CameraOutline" />
              </template>
              更换头像
            </n-button>
          </n-upload>
        </div>
        <div class="profile-summary">
          <h2 class="profile-name">{{ userForm.nickname || '用户' }}</h2>
          <p class="profile-id">ID: {{ userId }}</p>
          <div class="profile-badges">
            <n-tag type="info" size="small" round>普通会员</n-tag>
            <n-tag type="success" size="small" round>已实名</n-tag>
          </div>
        </div>
      </div>
    </n-card>

    <!-- Basic Info -->
    <n-card title="基本信息" :bordered="false" class="page-card" style="margin-top: 16px;">
      <n-form
        ref="profileFormRef"
        :model="userForm"
        :rules="profileRules"
        label-placement="left"
        label-width="80"
      >
        <n-grid :cols="2" :x-gap="24" responsive="screen" item-responsive>
          <n-gi span="2 m:1">
            <n-form-item label="昵称" path="nickname">
              <n-input v-model:value="userForm.nickname" placeholder="请输入昵称" />
            </n-form-item>
          </n-gi>
          <n-gi span="2 m:1">
            <n-form-item label="邮箱" path="email">
              <n-input v-model:value="userForm.email" placeholder="请输入邮箱" />
            </n-form-item>
          </n-gi>
          <n-gi span="2 m:1">
            <n-form-item label="手机号" path="phone">
              <n-input v-model:value="userForm.phone" placeholder="请输入手机号" />
            </n-form-item>
          </n-gi>
        </n-grid>

        <div class="form-actions">
          <n-button type="primary" :loading="saving" @click="handleSaveProfile">
            保存修改
          </n-button>
        </div>
      </n-form>
    </n-card>

    <!-- Security Settings -->
    <n-card title="安全设置" :bordered="false" class="page-card" style="margin-top: 16px;">
      <div class="security-list">
        <div class="security-item">
          <div class="security-info">
            <n-icon :component="LockClosedOutline" size="24" class="security-icon" />
            <div>
              <h4>修改密码</h4>
              <p>定期修改密码可以保护账号安全</p>
            </div>
          </div>
          <n-button @click="showPasswordModal = true">修改密码</n-button>
        </div>

        <n-divider style="margin: 0;" />

        <div class="security-item">
          <div class="security-info">
            <n-icon :component="CallOutline" size="24" class="security-icon" />
            <div>
              <h4>绑定手机</h4>
              <p>{{ userForm.phone ? `已绑定：${maskPhone(userForm.phone)}` : '未绑定' }}</p>
            </div>
          </div>
          <n-button @click="showPhoneModal = true">
            {{ userForm.phone ? '更换' : '绑定' }}
          </n-button>
        </div>

        <n-divider style="margin: 0;" />

        <div class="security-item">
          <div class="security-info">
            <n-icon :component="MailOutline" size="24" class="security-icon" />
            <div>
              <h4>绑定邮箱</h4>
              <p>{{ userForm.email ? `已绑定：${maskEmail(userForm.email)}` : '未绑定' }}</p>
            </div>
          </div>
          <n-button @click="showEmailModal = true">
            {{ userForm.email ? '更换' : '绑定' }}
          </n-button>
        </div>
      </div>
    </n-card>

    <!-- OAuth Bindings -->
    <n-card title="第三方账号绑定" :bordered="false" class="page-card" style="margin-top: 16px;">
      <p class="section-desc">绑定第三方账号可以快速登录</p>

      <div class="oauth-list">
        <div class="oauth-item" v-for="provider in oauthProviders" :key="provider.key">
          <div class="oauth-info">
            <div class="oauth-icon" v-html="provider.svg"></div>
            <div>
              <h4>{{ provider.label }}</h4>
              <p :class="{ bound: provider.bound }">{{ provider.bound ? '已绑定' : '未绑定' }}</p>
            </div>
          </div>
          <n-button
            :type="provider.bound ? 'default' : 'primary'"
            @click="handleOAuth(provider.key)"
          >
            {{ provider.bound ? '解绑' : '绑定' }}
          </n-button>
        </div>
      </div>
    </n-card>

    <!-- Change Password Modal -->
    <n-modal
      v-model:show="showPasswordModal"
      preset="card"
      title="修改密码"
      style="width: 480px; max-width: 90vw;"
      :bordered="false"
    >
      <n-form ref="passwordFormRef" :model="passwordForm" :rules="passwordRules">
        <n-form-item label="当前密码" path="oldPassword">
          <n-input
            v-model:value="passwordForm.oldPassword"
            type="password"
            show-password-on="click"
            placeholder="请输入当前密码"
          />
        </n-form-item>
        <n-form-item label="新密码" path="newPassword">
          <n-input
            v-model:value="passwordForm.newPassword"
            type="password"
            show-password-on="click"
            placeholder="请输入新密码 (至少6位)"
          />
        </n-form-item>
        <n-form-item label="确认密码" path="confirmPassword">
          <n-input
            v-model:value="passwordForm.confirmPassword"
            type="password"
            show-password-on="click"
            placeholder="请确认新密码"
          />
        </n-form-item>
      </n-form>

      <template #footer>
        <n-space justify="end">
          <n-button @click="showPasswordModal = false">取消</n-button>
          <n-button type="primary" :loading="changingPassword" @click="handleChangePassword">
            确认修改
          </n-button>
        </n-space>
      </template>
    </n-modal>

    <!-- Bind Phone Modal -->
    <n-modal
      v-model:show="showPhoneModal"
      preset="card"
      :title="userForm.phone ? '更换手机号' : '绑定手机号'"
      style="width: 480px; max-width: 90vw;"
      :bordered="false"
    >
      <n-form>
        <n-form-item v-if="userForm.phone" label="当前手机">
          <n-input :value="maskPhone(userForm.phone)" disabled />
        </n-form-item>
        <n-form-item label="新手机号">
          <n-input v-model:value="bindPhoneForm.phone" placeholder="请输入手机号" />
        </n-form-item>
        <n-form-item label="验证码">
          <div class="captcha-row">
            <n-input v-model:value="bindPhoneForm.code" placeholder="请输入验证码" />
            <n-button :disabled="phoneCooldown > 0" @click="handleSendPhoneCode">
              {{ phoneCooldown > 0 ? `${phoneCooldown}s` : '发送验证码' }}
            </n-button>
          </div>
        </n-form-item>
      </n-form>

      <template #footer>
        <n-space justify="end">
          <n-button @click="showPhoneModal = false">取消</n-button>
          <n-button type="primary" @click="handleBindPhone">确认</n-button>
        </n-space>
      </template>
    </n-modal>

    <!-- Bind Email Modal -->
    <n-modal
      v-model:show="showEmailModal"
      preset="card"
      :title="userForm.email ? '更换邮箱' : '绑定邮箱'"
      style="width: 480px; max-width: 90vw;"
      :bordered="false"
    >
      <n-form>
        <n-form-item v-if="userForm.email" label="当前邮箱">
          <n-input :value="maskEmail(userForm.email)" disabled />
        </n-form-item>
        <n-form-item label="新邮箱">
          <n-input v-model:value="bindEmailForm.email" placeholder="请输入邮箱" />
        </n-form-item>
        <n-form-item label="验证码">
          <div class="captcha-row">
            <n-input v-model:value="bindEmailForm.code" placeholder="请输入验证码" />
            <n-button :disabled="emailCooldown > 0" @click="handleSendEmailCode">
              {{ emailCooldown > 0 ? `${emailCooldown}s` : '发送验证码' }}
            </n-button>
          </div>
        </n-form-item>
      </n-form>

      <template #footer>
        <n-space justify="end">
          <n-button @click="showEmailModal = false">取消</n-button>
          <n-button type="primary" @click="handleBindEmail">确认</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useMessage } from 'naive-ui'
import type { FormInst, FormRules } from 'naive-ui'
import {
  CameraOutline,
  LockClosedOutline,
  CallOutline,
  MailOutline
} from '@vicons/ionicons5'

const message = useMessage()

const userId = ref('100001')
const avatarUrl = ref('')
const saving = ref(false)
const changingPassword = ref(false)

const showPasswordModal = ref(false)
const showPhoneModal = ref(false)
const showEmailModal = ref(false)

const phoneCooldown = ref(0)
const emailCooldown = ref(0)

const profileFormRef = ref<FormInst | null>(null)
const passwordFormRef = ref<FormInst | null>(null)

const userForm = reactive({
  nickname: '张三',
  email: 'zhangsan@example.com',
  phone: '13800008888'
})

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const bindPhoneForm = reactive({
  phone: '',
  code: ''
})

const bindEmailForm = reactive({
  email: '',
  code: ''
})

const oauthBindings = reactive({
  wechat: true,
  qq: false,
  github: true,
  google: false
})

const oauthProviders = computed(() => [
  {
    key: 'wechat',
    label: '微信',
    bound: oauthBindings.wechat,
    svg: '<svg viewBox="0 0 24 24" width="28" height="28" fill="#07c160"><path d="M8.691 2.188C3.891 2.188 0 5.476 0 9.53c0 2.212 1.17 4.203 3.002 5.55a.59.59 0 0 1 .213.665l-.39 1.48c-.019.07-.048.141-.048.213 0 .163.13.295.29.295a.326.326 0 0 0 .167-.054l1.903-1.114a.864.864 0 0 1 .717-.098 10.16 10.16 0 0 0 2.837.403c.276 0 .543-.027.811-.05-.857-2.578.157-4.972 1.932-6.446 1.703-1.415 3.882-1.98 5.853-1.838-.576-3.583-4.196-6.348-8.596-6.348zM5.785 5.991c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 0 1-1.162 1.178A1.17 1.17 0 0 1 4.623 7.17c0-.651.52-1.18 1.162-1.18zm5.813 0c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 0 1-1.162 1.178 1.17 1.17 0 0 1-1.162-1.178c0-.651.52-1.18 1.162-1.18zm5.34 2.867c-1.797-.052-3.746.512-5.28 1.786-1.72 1.428-2.687 3.72-1.78 6.22.942 2.453 3.666 4.229 6.884 4.229.826 0 1.622-.12 2.361-.336a.722.722 0 0 1 .598.082l1.584.926a.272.272 0 0 0 .14.047c.134 0 .24-.111.24-.247 0-.06-.023-.12-.038-.177l-.327-1.233a.582.582 0 0 1-.023-.156.49.49 0 0 1 .201-.398C23.024 18.48 24 16.82 24 14.98c0-3.21-2.931-5.952-7.062-6.122zM14.57 13.39c.535 0 .969.44.969.982a.976.976 0 0 1-.969.983.976.976 0 0 1-.969-.983c0-.542.434-.982.97-.982zm4.844 0c.535 0 .969.44.969.982a.976.976 0 0 1-.969.983.976.976 0 0 1-.969-.983c0-.542.434-.982.97-.982z"/></svg>'
  },
  {
    key: 'qq',
    label: 'QQ',
    bound: oauthBindings.qq,
    svg: '<svg viewBox="0 0 24 24" width="28" height="28" fill="#12b7f5"><path d="M12.003 2c-2.265 0-6.29 1.364-6.29 7.325v1.195S3.55 14.96 3.55 17.474c0 .665.17 1.025.281 1.025.114 0 .902-.484 1.748-2.072 0 0-.18 2.197 1.904 3.967 0 0-1.77.495-1.77 1.182 0 .686 4.078.43 6.29.43 2.239 0 6.29.256 6.29-.43 0-.687-1.77-1.182-1.77-1.182 2.085-1.77 1.905-3.967 1.905-3.967.845 1.588 1.634 2.072 1.746 2.072.111 0 .283-.36.283-1.025 0-2.514-2.166-6.954-2.166-6.954V9.325C18.29 3.364 14.268 2 12.003 2z"/></svg>'
  },
  {
    key: 'github',
    label: 'GitHub',
    bound: oauthBindings.github,
    svg: '<svg viewBox="0 0 24 24" width="28" height="28" fill="#333"><path d="M12 .297c-6.63 0-12 5.373-12 12 0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61C4.422 18.07 3.633 17.7 3.633 17.7c-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12"/></svg>'
  },
  {
    key: 'google',
    label: 'Google',
    bound: oauthBindings.google,
    svg: '<svg viewBox="0 0 24 24" width="28" height="28"><path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 0 1-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z" fill="#4285F4"/><path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/><path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05"/><path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/></svg>'
  }
])

const profileRules: FormRules = {
  nickname: { required: true, message: '请输入昵称', trigger: 'blur' },
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '邮箱格式不正确', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '手机号格式不正确', trigger: 'blur' }
  ]
}

const passwordRules: FormRules = {
  oldPassword: { required: true, message: '请输入当前密码', trigger: 'blur' },
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (_rule: any, value: string) => {
        return value === passwordForm.newPassword ? true : new Error('两次密码不一致')
      },
      trigger: 'blur'
    }
  ]
}

function maskPhone(phone: string) {
  return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')
}

function maskEmail(email: string) {
  const [name, domain] = email.split('@')
  return name.slice(0, 3) + '***@' + domain
}

async function handleSaveProfile() {
  try {
    await profileFormRef.value?.validate()
    saving.value = true
    await new Promise((resolve) => setTimeout(resolve, 1000))
    message.success('个人信息已更新')
  } catch {
    // Validation failed
  } finally {
    saving.value = false
  }
}

async function handleChangePassword() {
  try {
    await passwordFormRef.value?.validate()
    changingPassword.value = true
    await new Promise((resolve) => setTimeout(resolve, 1000))
    showPasswordModal.value = false
    message.success('密码已修改')
    passwordForm.oldPassword = ''
    passwordForm.newPassword = ''
    passwordForm.confirmPassword = ''
  } catch {
    // Validation failed
  } finally {
    changingPassword.value = false
  }
}

function handleSendPhoneCode() {
  if (!bindPhoneForm.phone) {
    message.warning('请输入手机号')
    return
  }
  phoneCooldown.value = 60
  const timer = setInterval(() => {
    phoneCooldown.value--
    if (phoneCooldown.value <= 0) clearInterval(timer)
  }, 1000)
  message.success('验证码已发送')
}

function handleSendEmailCode() {
  if (!bindEmailForm.email) {
    message.warning('请输入邮箱')
    return
  }
  emailCooldown.value = 60
  const timer = setInterval(() => {
    emailCooldown.value--
    if (emailCooldown.value <= 0) clearInterval(timer)
  }, 1000)
  message.success('验证码已发送')
}

function handleBindPhone() {
  showPhoneModal.value = false
  message.success('手机号已绑定')
}

function handleBindEmail() {
  showEmailModal.value = false
  message.success('邮箱已绑定')
}

function handleOAuth(platform: string) {
  const isBound = oauthBindings[platform as keyof typeof oauthBindings]
  if (isBound) {
    oauthBindings[platform as keyof typeof oauthBindings] = false
    message.success(`已解绑${oauthProviders.value.find(p => p.key === platform)?.label}`)
  } else {
    oauthBindings[platform as keyof typeof oauthBindings] = true
    message.success(`已绑定${oauthProviders.value.find(p => p.key === platform)?.label}`)
  }
}

function handleAvatarUpload() {
  message.info('头像上传功能')
}

onMounted(() => {
  // Fetch user profile
})
</script>

<style scoped>
.profile-page {
  display: flex;
  flex-direction: column;
  gap: 0;
  max-width: 900px;
}

.page-card {
  border-radius: 12px;
  border: 1px solid #f0f0f0;
}

.profile-header {
  display: flex;
  align-items: center;
  gap: 32px;
}

.avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.user-avatar {
  border: 4px solid #f0f5ff;
  box-shadow: 0 4px 12px rgba(24, 144, 255, 0.2);
}

.profile-summary {
  flex: 1;
}

.profile-name {
  font-size: 24px;
  font-weight: 700;
  color: #262626;
  margin: 0 0 4px 0;
}

.profile-id {
  font-size: 14px;
  color: #8c8c8c;
  font-family: monospace;
  margin: 0 0 12px 0;
}

.profile-badges {
  display: flex;
  gap: 8px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 8px;
}

.security-list {
  display: flex;
  flex-direction: column;
}

.security-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 0;
}

.security-info {
  display: flex;
  align-items: center;
  gap: 16px;
}

.security-icon {
  color: #1890ff;
}

.security-info h4 {
  font-size: 15px;
  font-weight: 600;
  color: #262626;
  margin: 0 0 4px 0;
}

.security-info p {
  font-size: 13px;
  color: #8c8c8c;
  margin: 0;
}

.section-desc {
  font-size: 14px;
  color: #8c8c8c;
  margin: 0 0 16px 0;
}

.oauth-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.oauth-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: #fafafa;
  border-radius: 12px;
  transition: background 0.2s;
}

.oauth-item:hover {
  background: #f0f5ff;
}

.oauth-info {
  display: flex;
  align-items: center;
  gap: 16px;
}

.oauth-icon {
  display: flex;
  align-items: center;
  justify-content: center;
}

.oauth-info h4 {
  font-size: 15px;
  font-weight: 600;
  color: #262626;
  margin: 0 0 4px 0;
}

.oauth-info p {
  font-size: 13px;
  color: #8c8c8c;
  margin: 0;
}

.oauth-info p.bound {
  color: #52c41a;
}

.captcha-row {
  display: flex;
  gap: 12px;
  width: 100%;
}

@media (max-width: 768px) {
  .profile-header {
    flex-direction: column;
    text-align: center;
  }

  .profile-badges {
    justify-content: center;
  }

  .security-item {
    flex-direction: column;
    gap: 12px;
    align-items: flex-start;
  }

  .oauth-item {
    flex-direction: column;
    gap: 12px;
    align-items: flex-start;
  }
}
</style>
