<template>
  <div class="profile-page">
    <el-card shadow="never" class="page-card">
      <div class="profile-header">
        <div class="avatar-section">
          <el-avatar :size="80" class="user-avatar">{{ userInitial }}</el-avatar>
          <el-upload :show-file-list="false" :auto-upload="false" accept="image/*" @change="handleAvatarUpload">
            <el-button size="small" text type="primary">
              <el-icon><Camera /></el-icon>更换头像
            </el-button>
          </el-upload>
        </div>
        <div class="profile-summary">
          <h2 class="profile-name">{{ userForm.nickname || '用户' }}</h2>
          <p class="profile-id">ID: {{ userId }}</p>
          <div class="profile-badges">
            <el-tag type="primary" size="small" effect="plain" round>普通会员</el-tag>
            <el-tag type="success" size="small" effect="plain" round>已实名</el-tag>
          </div>
        </div>
      </div>
    </el-card>

    <el-card shadow="never" class="page-card">
      <template #header>
        <span class="card-title">基本信息</span>
      </template>
      <el-form
        ref="profileFormRef"
        :model="userForm"
        :rules="profileRules"
        label-width="80px"
        label-position="left"
        class="profile-form"
      >
        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="userForm.nickname" placeholder="请输入昵称" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="userForm.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="userForm.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="saving" @click="handleSaveProfile">保存修改</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card shadow="never" class="page-card">
      <template #header>
        <span class="card-title">安全设置</span>
      </template>
      <div class="security-list">
        <div class="security-item">
          <div class="security-info">
            <el-icon :size="22" color="#0056FF"><Lock /></el-icon>
            <div>
              <h4>修改密码</h4>
              <p>定期修改密码可以保护账号安全</p>
            </div>
          </div>
          <el-button @click="showPasswordDialog = true">修改密码</el-button>
        </div>
        <el-divider style="margin: 12px 0;" />
        <div class="security-item">
          <div class="security-info">
            <el-icon :size="22" color="#0056FF"><Iphone /></el-icon>
            <div>
              <h4>绑定手机</h4>
              <p>{{ userForm.phone ? `已绑定：${maskPhone(userForm.phone)}` : '未绑定' }}</p>
            </div>
          </div>
          <el-button @click="showPhoneDialog = true">{{ userForm.phone ? '更换' : '绑定' }}</el-button>
        </div>
        <el-divider style="margin: 12px 0;" />
        <div class="security-item">
          <div class="security-info">
            <el-icon :size="22" color="#0056FF"><Message /></el-icon>
            <div>
              <h4>绑定邮箱</h4>
              <p>{{ userForm.email ? `已绑定：${maskEmail(userForm.email)}` : '未绑定' }}</p>
            </div>
          </div>
          <el-button @click="showEmailDialog = true">{{ userForm.email ? '更换' : '绑定' }}</el-button>
        </div>
      </div>
    </el-card>

    <!-- Dialogs -->
    <el-dialog v-model="showPasswordDialog" title="修改密码" width="460px" destroy-on-close>
      <el-form ref="passwordFormRef" :model="passwordForm" :rules="passwordRules" label-width="80px">
        <el-form-item label="当前密码" prop="oldPassword">
          <el-input v-model="passwordForm.oldPassword" type="password" show-password placeholder="请输入当前密码" />
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input v-model="passwordForm.newPassword" type="password" show-password placeholder="请输入新密码 (至少6位)" />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input v-model="passwordForm.confirmPassword" type="password" show-password placeholder="请确认新密码" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showPasswordDialog = false">取消</el-button>
        <el-button type="primary" :loading="changingPassword" @click="handleChangePassword">确认修改</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showPhoneDialog" :title="userForm.phone ? '更换手机号' : '绑定手机号'" width="460px" destroy-on-close>
      <el-form label-width="80px">
        <el-form-item v-if="userForm.phone" label="当前手机">
          <el-input :model-value="maskPhone(userForm.phone)" disabled />
        </el-form-item>
        <el-form-item label="新手机号">
          <el-input v-model="bindPhoneForm.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="验证码">
          <div class="captcha-row">
            <el-input v-model="bindPhoneForm.code" placeholder="请输入验证码" />
            <el-button :disabled="phoneCooldown > 0" @click="handleSendPhoneCode">
              {{ phoneCooldown > 0 ? `${phoneCooldown}s` : '发送验证码' }}
            </el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showPhoneDialog = false">取消</el-button>
        <el-button type="primary" @click="handleBindPhone">确认</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showEmailDialog" :title="userForm.email ? '更换邮箱' : '绑定邮箱'" width="460px" destroy-on-close>
      <el-form label-width="80px">
        <el-form-item v-if="userForm.email" label="当前邮箱">
          <el-input :model-value="maskEmail(userForm.email)" disabled />
        </el-form-item>
        <el-form-item label="新邮箱">
          <el-input v-model="bindEmailForm.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="验证码">
          <div class="captcha-row">
            <el-input v-model="bindEmailForm.code" placeholder="请输入验证码" />
            <el-button :disabled="emailCooldown > 0" @click="handleSendEmailCode">
              {{ emailCooldown > 0 ? `${emailCooldown}s` : '发送验证码' }}
            </el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEmailDialog = false">取消</el-button>
        <el-button type="primary" @click="handleBindEmail">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { Camera, Lock, Iphone, Message } from '@element-plus/icons-vue'
import request from '@/utils/request'

const userId = ref('')
const saving = ref(false)
const changingPassword = ref(false)
const showPasswordDialog = ref(false)
const showPhoneDialog = ref(false)
const showEmailDialog = ref(false)
const phoneCooldown = ref(0)
const emailCooldown = ref(0)
const profileFormRef = ref<FormInstance>()
const passwordFormRef = ref<FormInstance>()
const userInitial = computed(() => (userForm.nickname || '用户').charAt(0))

const userForm = reactive({ nickname: '', email: '', phone: '' })

async function fetchProfile() {
  try {
    const res = await request.get('/api/v2/user/profile')
    const profile = res.data?.data || res.data || {}
    userForm.nickname = profile.nickname || ''
    userForm.email = profile.email || ''
    userForm.phone = profile.phone || ''
    userId.value = profile.id || profile.user_id || ''
  } catch { /* ignore */ }
}

onMounted(() => { fetchProfile() })
const passwordForm = reactive({ oldPassword: '', newPassword: '', confirmPassword: '' })
const bindPhoneForm = reactive({ phone: '', code: '' })
const bindEmailForm = reactive({ email: '', code: '' })

const profileRules: FormRules = {
  nickname: [{ required: true, message: '请输入昵称', trigger: 'blur' }],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '邮箱格式不正确', trigger: 'blur' }
  ],
  phone: [{ pattern: /^1[3-9]\d{9}$/, message: '手机号格式不正确', trigger: 'blur' }]
}

const passwordRules: FormRules = {
  oldPassword: [{ required: true, message: '请输入当前密码', trigger: 'blur' }],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (_rule: any, value: string, callback: any) => {
        if (value !== passwordForm.newPassword) callback(new Error('两次密码不一致'))
        else callback()
      },
      trigger: 'blur'
    }
  ]
}

function maskPhone(phone: string) { return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2') }
function maskEmail(email: string) { const [name, domain] = email.split('@'); return name.slice(0, 3) + '***@' + domain }

async function handleSaveProfile() {
  if (!profileFormRef.value) return
  try {
    await profileFormRef.value.validate()
    saving.value = true
    await request.put('/api/v2/user/profile', { nickname: userForm.nickname, email: userForm.email, phone: userForm.phone })
    ElMessage.success('个人信息已更新')
  } catch (e: any) { ElMessage.error(e?.message || '保存失败，请重试') } finally { saving.value = false }
}

async function handleChangePassword() {
  if (!passwordFormRef.value) return
  try {
    await passwordFormRef.value.validate()
    changingPassword.value = true
    await request.post('/api/v2/user/change-password', { oldPassword: passwordForm.oldPassword, newPassword: passwordForm.newPassword })
    showPasswordDialog.value = false
    ElMessage.success('密码已修改')
    passwordForm.oldPassword = ''
    passwordForm.newPassword = ''
    passwordForm.confirmPassword = ''
  } catch (e: any) { ElMessage.error(e?.message || '密码修改失败，请重试') } finally { changingPassword.value = false }
}

function handleSendPhoneCode() {
  if (!bindPhoneForm.phone) { ElMessage.warning('请输入手机号'); return }
  phoneCooldown.value = 60
  const timer = setInterval(() => { phoneCooldown.value--; if (phoneCooldown.value <= 0) clearInterval(timer) }, 1000)
  ElMessage.success('验证码已发送')
}

function handleSendEmailCode() {
  if (!bindEmailForm.email) { ElMessage.warning('请输入邮箱'); return }
  emailCooldown.value = 60
  const timer = setInterval(() => { emailCooldown.value--; if (emailCooldown.value <= 0) clearInterval(timer) }, 1000)
  ElMessage.success('验证码已发送')
}

async function handleBindPhone() {
  if (!bindPhoneForm.phone || !bindPhoneForm.code) { ElMessage.warning('请填写手机号和验证码'); return }
  try {
    await request.post('/api/v2/user/bind-phone', { phone: bindPhoneForm.phone, code: bindPhoneForm.code })
    userForm.phone = bindPhoneForm.phone
    showPhoneDialog.value = false
    bindPhoneForm.phone = ''
    bindPhoneForm.code = ''
    ElMessage.success('手机号已绑定')
  } catch { ElMessage.error('绑定失败') }
}

async function handleBindEmail() {
  if (!bindEmailForm.email || !bindEmailForm.code) { ElMessage.warning('请填写邮箱和验证码'); return }
  try {
    await request.post('/api/v2/user/bind-email', { email: bindEmailForm.email, code: bindEmailForm.code })
    userForm.email = bindEmailForm.email
    showEmailDialog.value = false
    bindEmailForm.email = ''
    bindEmailForm.code = ''
    ElMessage.success('邮箱已绑定')
  } catch { ElMessage.error('绑定失败') }
}
function handleAvatarUpload() { ElMessage.info('头像上传功能') }
</script>

<style scoped>
.profile-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
  max-width: 860px;
}

.page-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
  background: #fff;
}

.card-title {
  font-size: 15px;
  font-weight: 600;
  color: #303133;
}

.profile-header {
  display: flex;
  align-items: center;
  gap: 28px;
}

.avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
}

.user-avatar {
  background: linear-gradient(135deg, #0056FF, #4080FF);
  color: #fff;
  font-weight: 700;
  font-size: 22px;
  border: 4px solid #EBF3FD;
}

.profile-summary { flex: 1; }
.profile-name { font-size: 20px; font-weight: 700; color: #303133; margin: 0 0 4px 0; }
.profile-id { font-size: 14px; color: #909399; font-family: monospace; margin: 0 0 12px 0; }
.profile-badges { display: flex; gap: 8px; }
.profile-form { max-width: 480px; }

.security-list { display: flex; flex-direction: column; }

.security-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
}

.security-info { display: flex; align-items: center; gap: 14px; }
.security-info h4 { font-size: 15px; font-weight: 600; color: #303133; margin: 0 0 4px 0; }
.security-info p { font-size: 13px; color: #909399; margin: 0; }
.captcha-row { display: flex; gap: 12px; width: 100%; }

@media (max-width: 768px) {
  .profile-header { flex-direction: column; text-align: center; }
  .profile-badges { justify-content: center; }
  .security-item { flex-direction: column; gap: 12px; align-items: flex-start; }
}
</style>
