<template>
  <div class="bind-account-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span>账号绑定</span>
        </div>
      </template>

      <div class="bind-section">
        <h3 class="section-title">手机绑定</h3>
        <div class="bind-item">
          <div class="bind-info">
            <el-icon :size="24" color="#409eff"><Iphone /></el-icon>
            <div class="bind-text">
              <div class="bind-label">手机号码</div>
              <div class="bind-value" v-if="phone">{{ phone }}</div>
              <div class="bind-value unbound" v-else>未绑定</div>
            </div>
          </div>
          <el-button v-if="phone" type="default" @click="showPhoneDialog('change')">更换</el-button>
          <el-button v-else type="primary" @click="showPhoneDialog('bind')">绑定</el-button>
        </div>
      </div>

      <el-divider />

      <div class="bind-section">
        <h3 class="section-title">邮箱绑定</h3>
        <div class="bind-item">
          <div class="bind-info">
            <el-icon :size="24" color="#67c23a"><Message /></el-icon>
            <div class="bind-text">
              <div class="bind-label">邮箱地址</div>
              <div class="bind-value" v-if="email">{{ email }}</div>
              <div class="bind-value unbound" v-else>未绑定</div>
            </div>
          </div>
          <el-button v-if="email" type="default" @click="showEmailDialog('change')">更换</el-button>
          <el-button v-else type="primary" @click="showEmailDialog('bind')">绑定</el-button>
        </div>
      </div>

      <el-divider />

      <div class="bind-section">
        <h3 class="section-title">{{ $t('user.bindAccount.qqBotBind') || 'QQ机器人绑定' }}</h3>
        <div class="bind-item">
          <div class="bind-info">
            <el-icon :size="24" color="#12b7f5"><ChatDotRound /></el-icon>
            <div class="bind-text">
              <div class="bind-label">{{ $t('user.bindAccount.qqNumber') || 'QQ号码' }}</div>
              <div class="bind-value" v-if="qqBound">{{ qqNumber }}</div>
              <div class="bind-value unbound" v-else>{{ $t('user.bindAccount.unbound') || '未绑定' }}</div>
            </div>
          </div>
          <el-button v-if="qqBound" type="danger" plain @click="unbindQQ">{{ $t('user.bindAccount.unbind') || '解绑' }}</el-button>
          <el-button v-else type="primary" @click="showQQDialog">{{ $t('user.bindAccount.bind') || '绑定' }}</el-button>
        </div>
        <div class="bind-tips">
          <el-text type="info" size="small">
            {{ $t('user.bindAccount.qqBotTips') || '绑定QQ号码后，可通过QQ机器人接收消息通知。多个QQ号请用英文逗号分隔。' }}
          </el-text>
        </div>
      </div>

      <el-divider />

      <div class="bind-section">
        <h3 class="section-title">第三方账号绑定</h3>
        <div
          v-for="provider in providers"
          :key="provider.name"
          class="bind-item"
        >
          <div class="bind-info">
            <img :src="provider.icon" :alt="provider.label" class="provider-icon" />
            <div class="bind-text">
              <div class="bind-label">{{ provider.label }}</div>
              <div class="bind-value" v-if="provider.bound">
                已绑定：{{ provider.account }}
              </div>
              <div class="bind-value unbound" v-else>未绑定</div>
            </div>
          </div>
          <el-button
            v-if="provider.bound"
            type="danger"
            plain
            @click="unbindProvider(provider)"
          >解绑</el-button>
          <el-button
            v-else
            type="primary"
            plain
            @click="bindProvider(provider)"
          >绑定</el-button>
        </div>
      </div>
    </el-card>

    <!-- 手机绑定对话框 -->
    <el-dialog v-model="phoneDialogVisible" :title="phoneDialogTitle" width="420px">
      <el-form :model="phoneForm" label-width="80px">
        <el-form-item label="手机号码">
          <el-input v-model="phoneForm.phone" placeholder="请输入手机号码" />
        </el-form-item>
        <el-form-item label="验证码">
          <div class="captcha-row">
            <el-input v-model="phoneForm.code" placeholder="请输入验证码" />
            <el-button type="primary" :disabled="phoneCodeCooldown > 0" @click="sendPhoneCode">
              {{ phoneCodeCooldown > 0 ? `${phoneCodeCooldown}s` : '获取验证码' }}
            </el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="phoneDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="confirmPhone">确定</el-button>
      </template>
    </el-dialog>

    <!-- 邮箱绑定对话框 -->
    <el-dialog v-model="emailDialogVisible" :title="emailDialogTitle" width="420px">
      <el-form :model="emailForm" label-width="80px">
        <el-form-item label="邮箱地址">
          <el-input v-model="emailForm.email" placeholder="请输入邮箱地址" />
        </el-form-item>
        <el-form-item label="验证码">
          <div class="captcha-row">
            <el-input v-model="emailForm.code" placeholder="请输入验证码" />
            <el-button type="primary" :disabled="emailCodeCooldown > 0" @click="sendEmailCode">
              {{ emailCodeCooldown > 0 ? `${emailCodeCooldown}s` : '获取验证码' }}
            </el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="emailDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="confirmEmail">确定</el-button>
      </template>
    </el-dialog>

    <!-- QQ机器人绑定对话框 -->
    <el-dialog v-model="qqDialogVisible" :title="$t('user.bindAccount.qqBotDialogTitle') || '绑定QQ机器人'" width="420px">
      <el-form :model="qqForm" label-width="80px">
        <el-form-item :label="$t('user.bindAccount.qqNumber') || 'QQ号码'">
          <el-input
            v-model="qqForm.qq"
            :placeholder="$t('user.bindAccount.qqPlaceholder') || '请输入QQ号码，多个用英文逗号分隔'"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="qqDialogVisible = false">{{ $t('common.cancel') || '取消' }}</el-button>
        <el-button type="primary" :loading="qqSubmitting" @click="confirmBindQQ">{{ $t('common.confirm') || '确定' }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Iphone, Message, ChatDotRound } from '@element-plus/icons-vue'
import request from '@/utils/request'

const phone = ref('')
const email = ref('')
const submitting = ref(false)

// QQ机器人绑定
const qqBound = ref(false)
const qqNumber = ref('')
const qqDialogVisible = ref(false)
const qqSubmitting = ref(false)
const qqForm = ref({ qq: '' })

// 手机绑定
const phoneDialogVisible = ref(false)
const phoneAction = ref<'bind' | 'change'>('bind')
const phoneForm = ref({ phone: '', code: '' })
const phoneCodeCooldown = ref(0)
const phoneDialogTitle = computed(() => phoneAction.value === 'bind' ? '绑定手机' : '更换手机')

// 邮箱绑定
const emailDialogVisible = ref(false)
const emailAction = ref<'bind' | 'change'>('bind')
const emailForm = ref({ email: '', code: '' })
const emailCodeCooldown = ref(0)
const emailDialogTitle = computed(() => emailAction.value === 'bind' ? '绑定邮箱' : '更换邮箱')

// 第三方绑定
const providers = ref([
  { name: 'wechat', label: '微信', icon: '/assets/oauth/wechat.svg', bound: false, account: '' },
  { name: 'qq', label: 'QQ', icon: '/assets/oauth/qq.svg', bound: false, account: '' },
  { name: 'github', label: 'GitHub', icon: '/assets/oauth/github.svg', bound: false, account: '' },
  { name: 'google', label: 'Google', icon: '/assets/oauth/google.svg', bound: false, account: '' }
])

const showPhoneDialog = (action: 'bind' | 'change') => {
  phoneAction.value = action
  phoneForm.value = { phone: '', code: '' }
  phoneDialogVisible.value = true
}

const showEmailDialog = (action: 'bind' | 'change') => {
  emailAction.value = action
  emailForm.value = { email: '', code: '' }
  emailDialogVisible.value = true
}

const startCooldown = (type: 'phone' | 'email') => {
  const cooldown = type === 'phone' ? phoneCodeCooldown : emailCodeCooldown
  cooldown.value = 60
  const timer = setInterval(() => {
    cooldown.value--
    if (cooldown.value <= 0) clearInterval(timer)
  }, 1000)
}

const sendPhoneCode = async () => {
  if (!phoneForm.value.phone) {
    ElMessage.warning('请输入手机号码')
    return
  }
  try {
    await request.post('/api/v1/sms/send', { phone: phoneForm.value.phone })
    startCooldown('phone')
    ElMessage.success('验证码已发送')
  } catch (e: any) {
    ElMessage.error(e.response?.data?.message || '发送失败，请稍后重试')
  }
}

const sendEmailCode = async () => {
  if (!emailForm.value.email) {
    ElMessage.warning('请输入邮箱地址')
    return
  }
  try {
    await request.post('/api/v1/email/send', { email: emailForm.value.email })
    startCooldown('email')
    ElMessage.success('验证码已发送')
  } catch (e: any) {
    ElMessage.error(e.response?.data?.message || '发送失败，请稍后重试')
  }
}

const showQQDialog = () => {
  qqForm.value.qq = qqNumber.value || ''
  qqDialogVisible.value = true
}

const confirmBindQQ = async () => {
  if (!qqForm.value.qq.trim()) {
    ElMessage.warning('请输入QQ号码')
    return
  }
  qqSubmitting.value = true
  try {
    await request.post('/api/v1/interflow/bind', { qq: qqForm.value.qq })
    qqNumber.value = qqForm.value.qq
    qqBound.value = true
    qqDialogVisible.value = false
    ElMessage.success('QQ机器人绑定成功')
  } catch (e: any) {
    ElMessage.error(e.response?.data?.message || '绑定失败')
  } finally {
    qqSubmitting.value = false
  }
}

const unbindQQ = async () => {
  try {
    await ElMessageBox.confirm('确定要解绑QQ机器人吗？', '确认解绑', {
      type: 'warning'
    })
    await request.post('/api/v1/interflow/unbind')
    qqBound.value = false
    qqNumber.value = ''
    ElMessage.success('解绑成功')
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.response?.data?.message || '解绑失败')
    }
  }
}

const confirmPhone = async () => {
  if (!phoneForm.value.phone || !phoneForm.value.code) {
    ElMessage.warning('请填写完整信息')
    return
  }
  submitting.value = true
  try {
    await request.post('/api/v1/user/bind-phone', {
    phone.value = phoneForm.value.phone
    phoneDialogVisible.value = false
    ElMessage.success('手机绑定成功')
  } catch (e: any) {
    ElMessage.error(e.response?.data?.message || '绑定失败')
  } finally {
    submitting.value = false
  }
}

const confirmEmail = async () => {
  if (!emailForm.value.email || !emailForm.value.code) {
    ElMessage.warning('请填写完整信息')
    return
  }
  submitting.value = true
  try {
    await request.post('/api/v1/user/bind-email', {
    email.value = emailForm.value.email
    emailDialogVisible.value = false
    ElMessage.success('邮箱绑定成功')
  } catch (e: any) {
    ElMessage.error(e.response?.data?.message || '绑定失败')
  } finally {
    submitting.value = false
  }
}

const bindProvider = async (provider: any) => {
  try {
    const { data } = await request.get(`/api/v1/oauth/${provider.name}`)
    if (data?.data?.url) {
      window.location.href = data.data.url
    } else {
      ElMessage.warning('暂不支持该第三方登录')
    }
  } catch (e: any) {
    ElMessage.error(e.response?.data?.message || '获取授权链接失败')
  }
}

const unbindProvider = async (provider: any) => {
  try {
    await ElMessageBox.confirm(`确定要解绑${provider.label}吗？`, '确认解绑', {
      type: 'warning'
    })
    await request.post('/api/v1/oauth/unbind', { provider: provider.name })
    provider.bound = false
    provider.account = ''
    ElMessage.success('解绑成功')
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.response?.data?.message || '解绑失败')
    }
  }
}

onMounted(async () => {
  try {
    const { data } = await request.get('/api/v1/oauth/accounts')
    if (data?.data) {
      const accounts = data.data
      accounts.forEach((account: any) => {
        const provider = providers.value.find(p => p.name === account.provider)
        if (provider) {
          provider.bound = true
          provider.account = account.account || account.username || ''
        }
      })
    }
  } catch {}
  try {
    const { data } = await request.get('/api/v1/user/profile')
    if (data?.data) {
      phone.value = data.data.phone || ''
      email.value = data.data.email || ''
    }
  } catch {}
  try {
    const { data } = await request.get('/api/v1/interflow/bind-info')
    if (data?.data?.qq) {
      qqNumber.value = data.data.qq.replace(/,\s*$/, '')
      qqBound.value = !!qqNumber.value
    }
  } catch {}
})
</script>

<style scoped lang="scss">
.bind-account-page {
  .bind-section {
    .section-title {
      font-size: 16px;
      font-weight: 600;
      color: #303133;
      margin: 0 0 16px 0;
    }
  }

  .bind-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 20px;
    border: 1px solid #ebeef5;
    border-radius: 8px;
    margin-bottom: 12px;

    &:hover {
      border-color: #c0c4cc;
    }

    .bind-info {
      display: flex;
      align-items: center;
      gap: 16px;

      .provider-icon {
        width: 24px;
        height: 24px;
      }

      .bind-text {
        .bind-label {
          font-weight: 500;
          margin-bottom: 4px;
        }

        .bind-value {
          font-size: 14px;
          color: #67c23a;

          &.unbound {
            color: #909399;
          }
        }
      }
    }
  }

  .captcha-row {
    display: flex;
    gap: 10px;
    width: 100%;
  }

  .bind-tips {
    margin-top: 8px;
    padding: 0 20px;
  }
}
</style>
