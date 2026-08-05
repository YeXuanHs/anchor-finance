<template>
  <div class="wechat-config-page">
    <!-- 微信登录配置 -->
    <el-card shadow="never" v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>微信登录配置</span>
          <el-button type="primary" @click="fetchConfig" :icon="Refresh">刷新</el-button>
        </div>
      </template>

      <el-form :model="loginForm" label-width="140px" :rules="loginRules" ref="loginFormRef">
        <el-divider content-position="left">基础配置</el-divider>

        <el-form-item label="启用微信登录">
          <el-switch v-model="loginForm.allow_wechat" />
          <span class="form-tip">开启后用户可使用微信扫码登录</span>
        </el-form-item>

        <el-form-item label="AppID" prop="wechat_login_appid">
          <el-input v-model="loginForm.wechat_login_appid" placeholder="请输入微信开放平台AppID" style="width: 400px" />
        </el-form-item>

        <el-form-item label="AppSecret" prop="wechat_login_secret">
          <el-input v-model="loginForm.wechat_login_secret" placeholder="请输入微信开放平台AppSecret" show-password style="width: 400px" />
        </el-form-item>

        <el-divider content-position="left">登录注册控制</el-divider>

        <el-form-item label="允许微信注册">
          <el-switch v-model="loginForm.allow_register_wechat" />
          <span class="form-tip">是否允许用户通过微信扫码注册新账号</span>
        </el-form-item>

        <el-form-item label="允许微信登录">
          <el-switch v-model="loginForm.allow_login_wechat" />
          <span class="form-tip">是否允许已绑定用户通过微信扫码登录</span>
        </el-form-item>

        <el-divider content-position="left">回调地址</el-divider>

        <el-form-item label="登录回调地址">
          <div class="callback-url">
            <el-input :model-value="loginCallbackUrl" readonly />
            <el-button type="primary" link @click="handleCopy(loginCallbackUrl)">复制</el-button>
          </div>
          <span class="form-tip">请将此地址配置到微信开放平台</span>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSaveLoginConfig" :loading="loginSaving">保存登录配置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 微信支付配置 -->
    <el-card shadow="never" class="section-card" v-loading="payLoading">
      <template #header>
        <div class="card-header">
          <span>微信支付配置</span>
        </div>
      </template>

      <el-form :model="payForm" label-width="140px" :rules="payRules" ref="payFormRef">
        <el-form-item label="启用微信支付">
          <el-switch v-model="payForm.enabled" />
          <span class="form-tip">开启后用户可使用微信支付进行付款</span>
        </el-form-item>

        <el-form-item label="商户号 (mch_id)" prop="mch_id">
          <el-input v-model="payForm.mch_id" placeholder="请输入微信支付商户号" style="width: 400px" />
        </el-form-item>

        <el-form-item label="商户API密钥" prop="api_key">
          <el-input v-model="payForm.api_key" placeholder="请输入商户API密钥" show-password style="width: 400px" />
        </el-form-item>

        <el-form-item label="商户AppID" prop="app_id">
          <el-input v-model="payForm.app_id" placeholder="请输入关联的商户AppID" style="width: 400px" />
        </el-form-item>

        <el-form-item label="子商户号">
          <el-input v-model="payForm.sub_mch_id" placeholder="服务商模式下的子商户号（可选）" style="width: 400px" />
        </el-form-item>

        <el-form-item label="证书文件">
          <div class="cert-upload">
            <el-upload
              :action="certUploadUrl"
              :headers="uploadHeaders"
              :on-success="handleCertSuccess"
              :before-upload="beforeCertUpload"
              :show-file-list="false"
            >
              <el-button type="primary">上传证书</el-button>
            </el-upload>
            <el-text v-if="payForm.cert_file" type="success" style="margin-left: 12px">
              已上传: {{ payForm.cert_file }}
            </el-text>
          </div>
        </el-form-item>

        <el-divider content-position="left">支付回调</el-divider>

        <el-form-item label="支付回调地址">
          <div class="callback-url">
            <el-input :model-value="payCallbackUrl" readonly />
            <el-button type="primary" link @click="handleCopy(payCallbackUrl)">复制</el-button>
          </div>
          <span class="form-tip">请将此地址配置到微信支付商户平台</span>
        </el-form-item>

        <el-form-item label="退款回调地址">
          <div class="callback-url">
            <el-input :model-value="refundCallbackUrl" readonly />
            <el-button type="primary" link @click="handleCopy(refundCallbackUrl)">复制</el-button>
          </div>
        </el-form-item>

        <el-divider content-position="left">高级设置</el-divider>

        <el-form-item label="支付超时时间">
          <el-input-number v-model="payForm.timeout" :min="5" :max="120" />
          <span class="form-tip">分钟，超过此时间未支付将自动关闭订单</span>
        </el-form-item>

        <el-form-item label="启用沙箱模式">
          <el-switch v-model="payForm.sandbox" />
          <span class="form-tip">调试用，正式环境请关闭</span>
        </el-form-item>

        <el-form-item label="启用Native支付">
          <el-switch v-model="payForm.native_pay" />
          <span class="form-tip">PC端扫码支付</span>
        </el-form-item>

        <el-form-item label="启用JSAPI支付">
          <el-switch v-model="payForm.jsapi_pay" />
          <span class="form-tip">微信内H5支付</span>
        </el-form-item>

        <el-form-item label="启用H5支付">
          <el-switch v-model="payForm.h5_pay" />
          <span class="form-tip">外部浏览器H5支付</span>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSavePayConfig" :loading="paySaving">保存支付配置</el-button>
          <el-button @click="handleTestPay" :loading="testPayLoading">测试连接</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 微信消息推送配置 -->
    <el-card shadow="never" class="section-card">
      <template #header>
        <div class="card-header">
          <span>微信消息推送配置</span>
        </div>
      </template>

      <el-form :model="messageForm" label-width="140px" ref="messageFormRef">
        <el-form-item label="启用消息推送">
          <el-switch v-model="messageForm.enabled" />
          <span class="form-tip">开启后系统将通过微信模板消息通知用户</span>
        </el-form-item>

        <el-form-item label="公众号AppID">
          <el-input v-model="messageForm.official_appid" placeholder="请输入公众号AppID" style="width: 400px" />
        </el-form-item>

        <el-form-item label="公众号AppSecret">
          <el-input v-model="messageForm.official_secret" placeholder="请输入公众号AppSecret" show-password style="width: 400px" />
        </el-form-item>

        <el-divider content-position="left">消息模板</el-divider>

        <el-form-item label="订单通知模板">
          <el-input v-model="messageForm.template_order" placeholder="模板ID" style="width: 400px" />
        </el-form-item>

        <el-form-item label="到期提醒模板">
          <el-input v-model="messageForm.template_expire" placeholder="模板ID" style="width: 400px" />
        </el-form-item>

        <el-form-item label="工单通知模板">
          <el-input v-model="messageForm.template_ticket" placeholder="模板ID" style="width: 400px" />
        </el-form-item>

        <el-form-item label="账单通知模板">
          <el-input v-model="messageForm.template_invoice" placeholder="模板ID" style="width: 400px" />
        </el-form-item>

        <el-form-item label="服务器消息模板">
          <el-input v-model="messageForm.template_server" placeholder="模板ID" style="width: 400px" />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSaveMessageConfig" :loading="messageSaving">保存消息配置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

const loading = ref(false)
const payLoading = ref(false)
const loginSaving = ref(false)
const paySaving = ref(false)
const messageSaving = ref(false)
const testPayLoading = ref(false)

const loginFormRef = ref<FormInstance>()
const payFormRef = ref<FormInstance>()
const messageFormRef = ref<FormInstance>()

const loginForm = reactive({
  allow_wechat: false,
  wechat_login_appid: '',
  wechat_login_secret: '',
  allow_register_wechat: false,
  allow_login_wechat: false
})

const payForm = reactive({
  enabled: false,
  mch_id: '',
  api_key: '',
  app_id: '',
  sub_mch_id: '',
  cert_file: '',
  timeout: 30,
  sandbox: false,
  native_pay: true,
  jsapi_pay: true,
  h5_pay: false
})

const messageForm = reactive({
  enabled: false,
  official_appid: '',
  official_secret: '',
  template_order: '',
  template_expire: '',
  template_ticket: '',
  template_invoice: '',
  template_server: ''
})

const loginRules: FormRules = {
  wechat_login_appid: [{ required: true, message: '请输入AppID', trigger: 'blur' }],
  wechat_login_secret: [{ required: true, message: '请输入AppSecret', trigger: 'blur' }]
}

const payRules: FormRules = {
  mch_id: [{ required: true, message: '请输入商户号', trigger: 'blur' }],
  api_key: [{ required: true, message: '请输入API密钥', trigger: 'blur' }],
  app_id: [{ required: true, message: '请输入商户AppID', trigger: 'blur' }]
}

const baseUrl = computed(() => window.location.origin)
const loginCallbackUrl = computed(() => `${baseUrl.value}/wechat_login_handle`)
const payCallbackUrl = computed(() => `${baseUrl.value}/api/payment/callback/wechat`)
const refundCallbackUrl = computed(() => `${baseUrl.value}/api/payment/refund/callback/wechat`)
const certUploadUrl = computed(() => `${import.meta.env.VITE_API_URL}/api/admin/wechat/pay/cert`)
const uploadHeaders = computed(() => ({
  Authorization: localStorage.getItem('token') || ''
}))

const fetchConfig = async () => {
  loading.value = true
  try {
    const res = await request.get({ url: '/api/admin/wechat/config' })
    if (res) {
      Object.assign(loginForm, {
        allow_wechat: res.allow_wechat ?? false,
        wechat_login_appid: res.wechat_login_appid ?? '',
        wechat_login_secret: res.wechat_login_secret ?? '',
        allow_register_wechat: res.allow_register_wechat ?? false,
        allow_login_wechat: res.allow_login_wechat ?? false
      })
    }
  } catch {
    ElMessage.error('获取微信登录配置失败')
  } finally {
    loading.value = false
  }
}

const fetchPayConfig = async () => {
  payLoading.value = true
  try {
    const res = await request.get({ url: '/api/admin/wechat/pay/config' })
    if (res) Object.assign(payForm, res)
  } catch { /* ignore */ } finally {
    payLoading.value = false
  }
}

const fetchMessageConfig = async () => {
  try {
    const res = await request.get({ url: '/api/admin/wechat/message/config' })
    if (res) Object.assign(messageForm, res)
  } catch { /* ignore */ }
}

const handleSaveLoginConfig = async () => {
  if (!loginFormRef.value) return
  await loginFormRef.value.validate(async (valid) => {
    if (!valid) return
    loginSaving.value = true
    try {
      await request.put({
        url: '/api/admin/wechat/config',
        data: loginForm,
        showSuccessMessage: true
      })
    } catch { /* error handled by request */ } finally {
      loginSaving.value = false
    }
  })
}

const handleSavePayConfig = async () => {
  if (!payFormRef.value) return
  await payFormRef.value.validate(async (valid) => {
    if (!valid) return
    paySaving.value = true
    try {
      await request.put({
        url: '/api/admin/wechat/pay/config',
        data: payForm,
        showSuccessMessage: true
      })
    } catch { /* error handled by request */ } finally {
      paySaving.value = false
    }
  })
}

const handleSaveMessageConfig = async () => {
  messageSaving.value = true
  try {
    await request.put({
      url: '/api/admin/wechat/message/config',
      data: messageForm,
      showSuccessMessage: true
    })
  } catch { /* error handled by request */ } finally {
    messageSaving.value = false
  }
}

const handleTestPay = async () => {
  testPayLoading.value = true
  try {
    await request.post({ url: '/api/admin/wechat/pay/test', showSuccessMessage: true })
  } catch { /* error handled by request */ } finally {
    testPayLoading.value = false
  }
}

const handleCertSuccess = (response: any) => {
  if (response?.data?.filename) {
    payForm.cert_file = response.data.filename
    ElMessage.success('证书上传成功')
  }
}

const beforeCertUpload = (file: File) => {
  const isValidType = file.name.endsWith('.pem') || file.name.endsWith('.p12')
  if (!isValidType) {
    ElMessage.error('请上传 .pem 或 .p12 格式的证书文件')
    return false
  }
  const isLt5M = file.size / 1024 / 1024 < 5
  if (!isLt5M) {
    ElMessage.error('证书文件大小不能超过 5MB')
    return false
  }
  return true
}

const handleCopy = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('已复制到剪贴板')
  } catch {
    const input = document.createElement('input')
    input.value = text
    document.body.appendChild(input)
    input.select()
    document.execCommand('copy')
    document.body.removeChild(input)
    ElMessage.success('已复制到剪贴板')
  }
}

onMounted(() => {
  fetchConfig()
  fetchPayConfig()
  fetchMessageConfig()
})
</script>

<style scoped lang="scss">
.wechat-config-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.section-card {
  margin-top: 20px;
}

.form-tip {
  margin-left: 12px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.callback-url {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  max-width: 500px;
}

.cert-upload {
  display: flex;
  align-items: center;
}
</style>
