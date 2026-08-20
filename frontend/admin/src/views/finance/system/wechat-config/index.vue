<template>
  <div class="wechat-config-page">
    <!-- 微信登录配置 -->
    <el-card shadow="never" v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ $t('wechatConfig.title') }}</span>
          <el-button type="primary" @click="fetchConfig" :icon="Refresh">{{ $t('wechatConfig.refresh') }}</el-button>
        </div>
      </template>

      <el-form :model="loginForm" label-width="140px" :rules="loginRules" ref="loginFormRef">
        <el-divider content-position="left">{{ $t('wechatConfig.basicConfig') }}</el-divider>

        <el-form-item :label="$t('wechatConfig.enableWechatLogin')">
          <el-switch v-model="loginForm.allow_wechat" />
          <span class="form-tip">{{ $t('wechatConfig.enableWechatLoginTip') }}</span>
        </el-form-item>

        <el-form-item label="AppID" prop="wechat_login_appid">
          <el-input v-model="loginForm.wechat_login_appid" :placeholder="$t('wechatConfig.appIdPlaceholder')" style="width: 400px" />
        </el-form-item>

        <el-form-item label="AppSecret" prop="wechat_login_secret">
          <el-input v-model="loginForm.wechat_login_secret" :placeholder="$t('wechatConfig.appSecretPlaceholder')" show-password style="width: 400px" />
        </el-form-item>

        <el-divider content-position="left">{{ $t('wechatConfig.loginRegisterControl') }}</el-divider>

        <el-form-item :label="$t('wechatConfig.allowWechatRegister')">
          <el-switch v-model="loginForm.allow_register_wechat" />
          <span class="form-tip">{{ $t('wechatConfig.allowWechatRegisterTip') }}</span>
        </el-form-item>

        <el-form-item :label="$t('wechatConfig.allowWechatLogin')">
          <el-switch v-model="loginForm.allow_login_wechat" />
          <span class="form-tip">{{ $t('wechatConfig.allowWechatLoginTip') }}</span>
        </el-form-item>

        <el-divider content-position="left">{{ $t('wechatConfig.callbackUrl') }}</el-divider>

        <el-form-item :label="$t('wechatConfig.loginCallback')">
          <div class="callback-url">
            <el-input :model-value="loginCallbackUrl" readonly />
            <el-button type="primary" link @click="handleCopy(loginCallbackUrl)">{{ $t('wechatConfig.copy') }}</el-button>
          </div>
          <span class="form-tip">{{ $t('wechatConfig.callbackTip') }}</span>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSaveLoginConfig" :loading="loginSaving">{{ $t('wechatConfig.saveLoginConfig') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 微信支付配置 -->
    <el-card shadow="never" class="section-card" v-loading="payLoading">
      <template #header>
        <div class="card-header">
          <span>{{ $t('wechatConfig.wechatPayConfig') }}</span>
        </div>
      </template>

      <el-form :model="payForm" label-width="140px" :rules="payRules" ref="payFormRef">
        <el-form-item :label="$t('wechatConfig.enableWechatPay')">
          <el-switch v-model="payForm.enabled" />
          <span class="form-tip">{{ $t('wechatConfig.enableWechatPayTip') }}</span>
        </el-form-item>

        <el-form-item :label="$t('wechatConfig.merchantId')" prop="mch_id">
          <el-input v-model="payForm.mch_id" :placeholder="$t('wechatConfig.merchantIdPlaceholder')" style="width: 400px" />
        </el-form-item>

        <el-form-item :label="$t('wechatConfig.apiKeyPlaceholder')" prop="api_key">
          <el-input v-model="payForm.api_key" :placeholder="$t('wechatConfig.apiKeyPlaceholder')" show-password style="width: 400px" />
        </el-form-item>

        <el-form-item :label="$t('wechatConfig.merchantAppId')" prop="app_id">
          <el-input v-model="payForm.app_id" :placeholder="$t('wechatConfig.merchantAppIdPlaceholder')" style="width: 400px" />
        </el-form-item>

        <el-form-item :label="$t('wechatConfig.subMerchantPlaceholder')">
          <el-input v-model="payForm.sub_mch_id" :placeholder="$t('wechatConfig.subMerchantPlaceholder')" style="width: 400px" />
        </el-form-item>

        <el-form-item :label="$t('wechatConfig.certFile')">
          <div class="cert-upload">
            <el-upload
              :action="certUploadUrl"
              :headers="uploadHeaders"
              :on-success="handleCertSuccess"
              :before-upload="beforeCertUpload"
              :show-file-list="false"
            >
              <el-button type="primary">{{ $t('wechatConfig.uploadCert') }}</el-button>
            </el-upload>
            <el-text v-if="payForm.cert_file" type="success" style="margin-left: 12px">
              {{ $t('wechatConfig.uploaded') }}: {{ payForm.cert_file }}
            </el-text>
          </div>
        </el-form-item>

        <el-divider content-position="left">{{ $t('wechatConfig.payCallback') }}</el-divider>

        <el-form-item :label="$t('wechatConfig.payCallback')">
          <div class="callback-url">
            <el-input :model-value="payCallbackUrl" readonly />
            <el-button type="primary" link @click="handleCopy(payCallbackUrl)">{{ $t('wechatConfig.copy') }}</el-button>
          </div>
          <span class="form-tip">{{ $t('wechatConfig.payCallbackTip') }}</span>
        </el-form-item>

        <el-form-item :label="$t('wechatConfig.refundCallback')">
          <div class="callback-url">
            <el-input :model-value="refundCallbackUrl" readonly />
            <el-button type="primary" link @click="handleCopy(refundCallbackUrl)">{{ $t('wechatConfig.copy') }}</el-button>
          </div>
        </el-form-item>

        <el-divider content-position="left">{{ $t('wechatConfig.advancedSettings') }}</el-divider>

        <el-form-item :label="$t('wechatConfig.payTimeout')">
          <el-input-number v-model="payForm.timeout" :min="5" :max="120" />
          <span class="form-tip">{{ $t('wechatConfig.payTimeoutTip') }}</span>
        </el-form-item>

        <el-form-item :label="$t('wechatConfig.enableSandbox')">
          <el-switch v-model="payForm.sandbox" />
          <span class="form-tip">{{ $t('wechatConfig.sandboxTip') }}</span>
        </el-form-item>

        <el-form-item :label="$t('wechatConfig.enableNativePay')">
          <el-switch v-model="payForm.native_pay" />
          <span class="form-tip">{{ $t('wechatConfig.nativePayTip') }}</span>
        </el-form-item>

        <el-form-item :label="$t('wechatConfig.enableJsapiPay')">
          <el-switch v-model="payForm.jsapi_pay" />
          <span class="form-tip">{{ $t('wechatConfig.jsapiPayTip') }}</span>
        </el-form-item>

        <el-form-item :label="$t('wechatConfig.enableH5Pay')">
          <el-switch v-model="payForm.h5_pay" />
          <span class="form-tip">{{ $t('wechatConfig.h5PayTip') }}</span>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSavePayConfig" :loading="paySaving">{{ $t('wechatConfig.savePayConfig') }}</el-button>
          <el-button @click="handleTestPay" :loading="testPayLoading">{{ $t('wechatConfig.testConnection') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 微信消息推送配置 -->
    <el-card shadow="never" class="section-card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('wechatConfig.wechatMessageConfig') }}</span>
        </div>
      </template>

      <el-form :model="messageForm" label-width="140px" ref="messageFormRef">
        <el-form-item :label="$t('wechatConfig.enableMessagePush')">
          <el-switch v-model="messageForm.enabled" />
          <span class="form-tip">{{ $t('wechatConfig.enableMessagePushTip') }}</span>
        </el-form-item>

        <el-form-item :label="$t('wechatConfig.officialAppIdPlaceholder')">
          <el-input v-model="messageForm.official_appid" :placeholder="$t('wechatConfig.officialAppIdPlaceholder')" style="width: 400px" />
        </el-form-item>

        <el-form-item :label="$t('wechatConfig.officialSecretPlaceholder')">
          <el-input v-model="messageForm.official_secret" :placeholder="$t('wechatConfig.officialSecretPlaceholder')" show-password style="width: 400px" />
        </el-form-item>

        <el-divider content-position="left">{{ $t('wechatConfig.messageTemplate') }}</el-divider>

        <el-form-item :label="$t('wechatConfig.orderTemplate')">
          <el-input v-model="messageForm.template_order" :placeholder="$t('wechatConfig.templateIdPlaceholder')" style="width: 400px" />
        </el-form-item>

        <el-form-item :label="$t('wechatConfig.expireTemplate')">
          <el-input v-model="messageForm.template_expire" :placeholder="$t('wechatConfig.templateIdPlaceholder')" style="width: 400px" />
        </el-form-item>

        <el-form-item :label="$t('wechatConfig.ticketTemplate')">
          <el-input v-model="messageForm.template_ticket" :placeholder="$t('wechatConfig.templateIdPlaceholder')" style="width: 400px" />
        </el-form-item>

        <el-form-item :label="$t('wechatConfig.invoiceTemplate')">
          <el-input v-model="messageForm.template_invoice" :placeholder="$t('wechatConfig.templateIdPlaceholder')" style="width: 400px" />
        </el-form-item>

        <el-form-item :label="$t('wechatConfig.serverTemplate')">
          <el-input v-model="messageForm.template_server" :placeholder="$t('wechatConfig.templateIdPlaceholder')" style="width: 400px" />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSaveMessageConfig" :loading="messageSaving">{{ $t('wechatConfig.saveMessageConfig') }}</el-button>
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
import { $t } from '@/locales'
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

const loginRules = computed<FormRules>(() => ({
  wechat_login_appid: [{ required: true, message: $t('wechatConfig.appIdPlaceholder'), trigger: 'blur' }],
  wechat_login_secret: [{ required: true, message: $t('wechatConfig.appSecretPlaceholder'), trigger: 'blur' }]
}))

const payRules = computed<FormRules>(() => ({
  mch_id: [{ required: true, message: $t('wechatConfig.merchantIdPlaceholder'), trigger: 'blur' }],
  api_key: [{ required: true, message: $t('wechatConfig.apiKeyPlaceholder'), trigger: 'blur' }],
  app_id: [{ required: true, message: $t('wechatConfig.merchantAppIdPlaceholder'), trigger: 'blur' }]
}))

const baseUrl = computed(() => window.location.origin)
const loginCallbackUrl = computed(() => `${baseUrl.value}/wechat_login_handle`)
const payCallbackUrl = computed(() => `${baseUrl.value}/api/payment/callback/wechat`)
const refundCallbackUrl = computed(() => `${baseUrl.value}/api/payment/refund/callback/wechat`)
const certUploadUrl = computed(() => `${import.meta.env.VITE_API_URL || ''}/api/admin/wechat/pay/cert`)
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
    ElMessage.error($t('wechatConfig.fetchConfigFailed'))
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
    ElMessage.success($t('wechatConfig.certUploadSuccess'))
  }
}

const beforeCertUpload = (file: File) => {
  const isValidType = file.name.endsWith('.pem') || file.name.endsWith('.p12')
  if (!isValidType) {
    ElMessage.error($t('wechatConfig.certTypeError'))
    return false
  }
  const isLt5M = file.size / 1024 / 1024 < 5
  if (!isLt5M) {
    ElMessage.error($t('wechatConfig.certSizeError'))
    return false
  }
  return true
}

const handleCopy = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success($t('wechatConfig.copied'))
  } catch {
    const input = document.createElement('input')
    input.value = text
    document.body.appendChild(input)
    input.select()
    document.execCommand('copy')
    document.body.removeChild(input)
    ElMessage.success($t('wechatConfig.copied'))
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
