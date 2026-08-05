<template>
  <div class="security-page">
    <div class="page-header">
      <h1 class="page-title">安全设置</h1>
    </div>

    <el-card shadow="never" class="score-card">
      <div class="score-content">
        <div class="score-info">
          <h3>账号安全评分</h3>
          <p>完善安全设置可提高账号安全性</p>
        </div>
        <div class="score-ring">
          <el-progress
            type="circle"
            :percentage="securityScore"
            :width="100"
            :stroke-width="8"
            :color="scoreColor"
          >
            <template #default="{ percentage }">
              <span class="score-text">{{ percentage }}</span>
              <span class="score-unit">分</span>
            </template>
          </el-progress>
        </div>
      </div>
    </el-card>

    <el-card shadow="never" class="options-card">
      <div class="security-list">
        <div class="security-item">
          <div class="security-left">
            <div class="security-icon" style="background: rgba(0,86,255,0.08);">
              <el-icon :size="22" color="#0056FF"><Lock /></el-icon>
            </div>
            <div>
              <h4>登录密码</h4>
              <p>已设置。定期修改密码可提高安全性</p>
            </div>
          </div>
          <el-button @click="showPasswordDialog = true">修改密码</el-button>
        </div>
        <el-divider />

        <div class="security-item">
          <div class="security-left">
            <div class="security-icon" style="background: rgba(82,196,26,0.08);">
              <el-icon :size="22" color="#52c41a"><Key /></el-icon>
            </div>
            <div>
              <h4>两步验证</h4>
              <p>{{ twoFAEnabled ? '已开启' : '未开启' }}。登录时需要验证码二次确认</p>
            </div>
          </div>
          <el-switch v-model="twoFAEnabled" active-text="已开启" inactive-text="已关闭" @change="handleToggle2FA" />
        </div>
        <el-divider />

        <div class="security-item">
          <div class="security-left">
            <div class="security-icon" style="background: rgba(250,140,22,0.08);">
              <el-icon :size="22" color="#fa8c16"><Iphone /></el-icon>
            </div>
            <div>
              <h4>绑定手机</h4>
              <p>{{ phone ? `已绑定：${maskPhone(phone)}` : '未绑定' }}</p>
            </div>
          </div>
          <el-button @click="showPhoneDialog = true">{{ phone ? '更换' : '绑定' }}</el-button>
        </div>
        <el-divider />

        <div class="security-item">
          <div class="security-left">
            <div class="security-icon" style="background: rgba(114,46,209,0.08);">
              <el-icon :size="22" color="#722ed1"><Message /></el-icon>
            </div>
            <div>
              <h4>绑定邮箱</h4>
              <p>{{ email ? `已绑定：${maskEmail(email)}` : '未绑定' }}</p>
            </div>
          </div>
          <el-button @click="showEmailDialog = true">{{ email ? '更换' : '绑定' }}</el-button>
        </div>
        <el-divider />

        <div class="security-item">
          <div class="security-left">
            <div class="security-icon" style="background: rgba(38,166,154,0.08);">
              <el-icon :size="22" color="#26a69a"><Clock /></el-icon>
            </div>
            <div>
              <h4>登录日志</h4>
              <p>查看最近的登录记录</p>
            </div>
          </div>
          <el-button @click="showLoginLog = true">查看</el-button>
        </div>
        <el-divider />

        <div class="security-item">
          <div class="security-left">
            <div class="security-icon" style="background: rgba(245,34,45,0.08);">
              <el-icon :size="22" color="#f5222d"><Key /></el-icon>
            </div>
            <div>
              <h4>API 密钥</h4>
              <p>管理 API 访问密钥</p>
            </div>
          </div>
          <el-button>管理密钥</el-button>
        </div>
      </div>
    </el-card>

    <el-dialog v-model="showLoginLog" title="登录日志" width="640px">
      <el-table :data="loginLogs" stripe>
        <el-table-column prop="time" label="时间" width="170" />
        <el-table-column prop="ip" label="IP 地址" width="140" />
        <el-table-column prop="location" label="登录地点" min-width="120" />
        <el-table-column prop="device" label="设备" min-width="120" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 'success' ? 'success' : 'danger'" size="small" effect="light" round>
              {{ row.statusText }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>

    <el-dialog v-model="showPasswordDialog" title="修改密码" width="460px" destroy-on-close>
      <el-form ref="passwordFormRef" :model="passwordForm" :rules="passwordRules" label-width="80px">
        <el-form-item label="当前密码" prop="oldPassword">
          <el-input v-model="passwordForm.oldPassword" type="password" show-password />
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input v-model="passwordForm.newPassword" type="password" show-password />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input v-model="passwordForm.confirmPassword" type="password" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showPasswordDialog = false">取消</el-button>
        <el-button type="primary" @click="handleChangePassword">确认修改</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showPhoneDialog" title="绑定手机" width="460px" destroy-on-close>
      <el-form label-width="80px">
        <el-form-item label="手机号">
          <el-input v-model="phoneForm.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="验证码">
          <div class="captcha-row">
            <el-input v-model="phoneForm.code" placeholder="请输入验证码" />
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

    <el-dialog v-model="showEmailDialog" title="绑定邮箱" width="460px" destroy-on-close>
      <el-form label-width="80px">
        <el-form-item label="邮箱">
          <el-input v-model="emailForm.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="验证码">
          <div class="captcha-row">
            <el-input v-model="emailForm.code" placeholder="请输入验证码" />
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
import { Lock, Key, Iphone, Message, Clock } from '@element-plus/icons-vue'
import request from '@/utils/request'

const twoFAEnabled = ref(false)
const phone = ref('')
const email = ref('')

async function fetchUserProfile() {
  try {
    const res = await request.get('/api/v1/user/profile')
    const profile = res.data?.data || res.data || {}
    phone.value = profile.phone || ''
    email.value = profile.email || ''
    twoFAEnabled.value = profile.twoFA || profile.two_fa || false
  } catch { /* ignore */ }
}

async function fetchLoginLogs() {
  try {
    const res = await request.get('/api/v1/login-logs')
    loginLogs.value = res.data?.data || res.data || []
  } catch { /* ignore */ }
}

onMounted(() => {
  fetchUserProfile()
  fetchLoginLogs()
})

const showPasswordDialog = ref(false)
const showPhoneDialog = ref(false)
const showEmailDialog = ref(false)
const showLoginLog = ref(false)
const phoneCooldown = ref(0)
const emailCooldown = ref(0)

const passwordFormRef = ref<FormInstance>()
const passwordForm = reactive({ oldPassword: '', newPassword: '', confirmPassword: '' })
const phoneForm = reactive({ phone: '', code: '' })
const emailForm = reactive({ email: '', code: '' })

const passwordRules: FormRules = {
  oldPassword: [{ required: true, message: '请输入当前密码', trigger: 'blur' }],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (_r: any, v: string, cb: any) => {
        v !== passwordForm.newPassword ? cb(new Error('两次密码不一致')) : cb()
      },
      trigger: 'blur'
    }
  ]
}

const securityScore = computed(() => {
  let score = 40
  if (phone.value) score += 20
  if (email.value) score += 20
  if (twoFAEnabled.value) score += 20
  return score
})

const scoreColor = computed(() => {
  if (securityScore.value >= 80) return '#52c41a'
  if (securityScore.value >= 60) return '#fa8c16'
  return '#f5222d'
})

const loginLogs = ref<any[]>([])

function maskPhone(p: string) { return p.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2') }
function maskEmail(e: string) { const [n, d] = e.split('@'); return n.slice(0, 3) + '***@' + d }

async function handleChangePassword() {
  if (!passwordFormRef.value) return
  try {
    await passwordFormRef.value.validate()
    await request.post('/api/v1/user/change-password', {
      oldPassword: passwordForm.oldPassword,
      newPassword: passwordForm.newPassword
    })
    showPasswordDialog.value = false
    ElMessage.success('密码已修改')
    passwordForm.oldPassword = ''
    passwordForm.newPassword = ''
    passwordForm.confirmPassword = ''
  } catch {}
}

async function handleBindPhone() {
  if (!phoneForm.phone || !phoneForm.code) { ElMessage.warning('请填写手机号和验证码'); return }
  try {
    await request.post('/api/v1/user/bind-phone', { phone: phoneForm.phone, code: phoneForm.code })
    phone.value = phoneForm.phone
    showPhoneDialog.value = false
    phoneForm.phone = ''
    phoneForm.code = ''
    ElMessage.success('绑定成功')
  } catch { ElMessage.error('绑定失败') }
}

async function handleBindEmail() {
  if (!emailForm.email || !emailForm.code) { ElMessage.warning('请填写邮箱和验证码'); return }
  try {
    await request.post('/api/v1/user/bind-email', { email: emailForm.email, code: emailForm.code })
    email.value = emailForm.email
    showEmailDialog.value = false
    emailForm.email = ''
    emailForm.code = ''
    ElMessage.success('绑定成功')
  } catch { ElMessage.error('绑定失败') }
}

async function handleSendPhoneCode() {
  if (!phoneForm.phone) { ElMessage.warning('请输入手机号'); return }
  try {
    await request.post('/api/v1/sms/send', { phone: phoneForm.phone, type: 'bind_phone' })
    phoneCooldown.value = 60
    const timer = setInterval(() => { phoneCooldown.value--; if (phoneCooldown.value <= 0) clearInterval(timer) }, 1000)
    ElMessage.success('验证码已发送')
  } catch { ElMessage.error('发送失败') }
}

async function handleSendEmailCode() {
  if (!emailForm.email) { ElMessage.warning('请输入邮箱'); return }
  try {
    await request.post('/api/v1/email/send', { email: emailForm.email, type: 'bind_email' })
    emailCooldown.value = 60
    const timer = setInterval(() => { emailCooldown.value--; if (emailCooldown.value <= 0) clearInterval(timer) }, 1000)
    ElMessage.success('验证码已发送')
  } catch { ElMessage.error('发送失败') }
}

async function handleToggle2FA(val: boolean) {
  try {
    const endpoint = val ? '/api/v1/user/2fa/enable' : '/api/v1/user/2fa/disable'
    await request.post(endpoint)
    ElMessage.success(val ? '已开启两步验证' : '已关闭两步验证')
  } catch {
    twoFAEnabled.value = !val
    ElMessage.error('操作失败')
  }
}
</script>

<style scoped>
.security-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
  max-width: 800px;
}

.page-title {
  font-size: 20px;
  font-weight: 700;
  color: #303133;
  margin: 0;
}

.score-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
  background: linear-gradient(135deg, #f5f7fa, #EBF3FD);
}

.score-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.score-info h3 { font-size: 18px; font-weight: 600; color: #303133; margin: 0 0 6px 0; }
.score-info p { font-size: 14px; color: #909399; margin: 0; }

.score-text { font-size: 28px; font-weight: 700; color: #303133; }
.score-unit { font-size: 12px; color: #909399; margin-left: 2px; }

.options-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
  background: #fff;
}

.security-list { display: flex; flex-direction: column; }

.security-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 0;
}

.security-left { display: flex; align-items: center; gap: 16px; }

.security-icon {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.security-left h4 { font-size: 15px; font-weight: 600; color: #303133; margin: 0 0 4px 0; }
.security-left p { font-size: 13px; color: #909399; margin: 0; }
.captcha-row { display: flex; gap: 12px; width: 100%; }

@media (max-width: 768px) {
  .score-content { flex-direction: column; gap: 16px; text-align: center; }
  .security-item { flex-direction: column; gap: 12px; align-items: flex-start; }
}
</style>
