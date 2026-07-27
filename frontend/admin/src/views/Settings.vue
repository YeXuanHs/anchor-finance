<template>
  <n-card :bordered="false" rounded>
    <n-tabs v-model:value="activeTab" type="line" animated>
      <!-- 基本设置 -->
      <n-tab-pane name="basic" tab="基本设置">
        <n-form
          :model="basicSettings"
          label-placement="left"
          label-width="120"
          style="max-width: 640px; margin-top: 20px"
        >
          <n-form-item label="站点名称">
            <n-input v-model:value="basicSettings.siteName" placeholder="锚点财务" />
          </n-form-item>
          <n-form-item label="站点Logo">
            <n-upload
              :max="1"
              accept="image/*"
              :default-upload="false"
              @change="handleLogoChange"
            >
              <n-button>
                <template #icon><n-icon><UploadIcon /></n-icon></template>
                上传Logo
              </n-button>
            </n-upload>
          </n-form-item>
          <n-form-item label="站点描述">
            <n-input
              v-model:value="basicSettings.siteDescription"
              type="textarea"
              :rows="3"
              placeholder="站点描述信息"
            />
          </n-form-item>
          <n-form-item label="站点URL">
            <n-input v-model:value="basicSettings.siteUrl" placeholder="https://example.com" />
          </n-form-item>
          <n-form-item label="管理员邮箱">
            <n-input v-model:value="basicSettings.adminEmail" placeholder="admin@example.com" />
          </n-form-item>
          <n-form-item label="时区">
            <n-select v-model:value="basicSettings.timezone" :options="timezoneOptions" />
          </n-form-item>
          <n-form-item label="默认语言">
            <n-select v-model:value="basicSettings.language" :options="languageOptions" />
          </n-form-item>
          <n-form-item label="用户注册">
            <n-switch v-model:value="basicSettings.allowRegistration" />
          </n-form-item>
          <n-form-item label="维护模式">
            <n-switch v-model:value="basicSettings.maintenanceMode" />
          </n-form-item>
          <n-form-item label="维护提示" v-if="basicSettings.maintenanceMode">
            <n-input v-model:value="basicSettings.maintenanceMessage" placeholder="系统维护中，请稍后再试" />
          </n-form-item>
          <n-form-item>
            <n-button type="primary" @click="saveBasic">保存设置</n-button>
          </n-form-item>
        </n-form>
      </n-tab-pane>

      <!-- 邮件设置 -->
      <n-tab-pane name="email" tab="邮件设置">
        <n-form
          :model="emailSettings"
          label-placement="left"
          label-width="120"
          style="max-width: 640px; margin-top: 20px"
        >
          <n-form-item label="邮件驱动">
            <n-select v-model:value="emailSettings.driver" :options="emailDriverOptions" />
          </n-form-item>
          <n-form-item label="SMTP 主机">
            <n-input v-model:value="emailSettings.smtpHost" placeholder="smtp.example.com" />
          </n-form-item>
          <n-form-item label="SMTP 端口">
            <n-input-number v-model:value="emailSettings.smtpPort" :min="1" :max="65535" style="width: 100%" />
          </n-form-item>
          <n-form-item label="加密方式">
            <n-select v-model:value="emailSettings.encryption" :options="encryptionOptions" />
          </n-form-item>
          <n-form-item label="用户名">
            <n-input v-model:value="emailSettings.username" placeholder="noreply@example.com" />
          </n-form-item>
          <n-form-item label="密码">
            <n-input v-model:value="emailSettings.password" type="password" show-password-on="click" placeholder="SMTP密码" />
          </n-form-item>
          <n-form-item label="发件人名称">
            <n-input v-model:value="emailSettings.fromName" placeholder="锚点财务" />
          </n-form-item>
          <n-form-item label="发件人邮箱">
            <n-input v-model:value="emailSettings.fromEmail" placeholder="noreply@example.com" />
          </n-form-item>
          <n-form-item>
            <n-space>
              <n-button type="primary" @click="saveEmail">保存设置</n-button>
              <n-button @click="testEmail">发送测试邮件</n-button>
            </n-space>
          </n-form-item>
        </n-form>
      </n-tab-pane>

      <!-- 短信设置 -->
      <n-tab-pane name="sms" tab="短信设置">
        <n-form
          :model="smsSettings"
          label-placement="left"
          label-width="120"
          style="max-width: 640px; margin-top: 20px"
        >
          <n-form-item label="短信服务商">
            <n-select v-model:value="smsSettings.provider" :options="smsProviderOptions" />
          </n-form-item>
          <n-form-item label="App Key">
            <n-input v-model:value="smsSettings.appKey" placeholder="短信服务商AppKey" />
          </n-form-item>
          <n-form-item label="App Secret">
            <n-input v-model:value="smsSettings.appSecret" type="password" show-password-on="click" placeholder="短信服务商AppSecret" />
          </n-form-item>
          <n-form-item label="短信签名">
            <n-input v-model:value="smsSettings.signName" placeholder="短信签名" />
          </n-form-item>

          <n-divider>短信模板</n-divider>

          <n-form-item label="验证码模板">
            <n-input v-model:value="smsSettings.templateCode" placeholder="SMS_XXXXXX" />
          </n-form-item>
          <n-form-item label="通知模板">
            <n-input v-model:value="smsSettings.notifyTemplateCode" placeholder="SMS_XXXXXX" />
          </n-form-item>
          <n-form-item label="启用场景">
            <n-checkbox-group v-model:value="smsSettings.scenes">
              <n-checkbox value="login" label="登录验证" />
              <n-checkbox value="register" label="注册验证" />
              <n-checkbox value="reset" label="密码重置" />
              <n-checkbox value="notify" label="订单通知" />
            </n-checkbox-group>
          </n-form-item>
          <n-form-item>
            <n-space>
              <n-button type="primary" @click="saveSms">保存设置</n-button>
              <n-button @click="testSms">发送测试短信</n-button>
            </n-space>
          </n-form-item>
        </n-form>
      </n-tab-pane>

      <!-- 支付设置 -->
      <n-tab-pane name="payment" tab="支付设置">
        <n-grid :cols="2" :x-gap="24" style="margin-top: 20px">
          <n-gi>
            <n-card title="支付宝配置" size="small" :bordered="false" style="background:#fafafa;border-radius:12px">
              <n-form :model="paymentSettings.alipay" label-placement="left" label-width="100">
                <n-form-item label="启用支付宝">
                  <n-switch v-model:value="paymentSettings.alipay.enabled" />
                </n-form-item>
                <template v-if="paymentSettings.alipay.enabled">
                  <n-form-item label="应用ID">
                    <n-input v-model:value="paymentSettings.alipay.appId" placeholder="支付宝应用ID" />
                  </n-form-item>
                  <n-form-item label="私钥">
                    <n-input v-model:value="paymentSettings.alipay.privateKey" type="textarea" :rows="3" placeholder="应用私钥" />
                  </n-form-item>
                  <n-form-item label="支付宝公钥">
                    <n-input v-model:value="paymentSettings.alipay.publicKey" type="textarea" :rows="3" placeholder="支付宝公钥" />
                  </n-form-item>
                  <n-form-item label="回调地址">
                    <n-input v-model:value="paymentSettings.alipay.notifyUrl" placeholder="https://example.com/api/pay/alipay/notify" />
                  </n-form-item>
                </template>
              </n-form>
            </n-card>
          </n-gi>
          <n-gi>
            <n-card title="微信支付配置" size="small" :bordered="false" style="background:#fafafa;border-radius:12px">
              <n-form :model="paymentSettings.wechat" label-placement="left" label-width="100">
                <n-form-item label="启用微信">
                  <n-switch v-model:value="paymentSettings.wechat.enabled" />
                </n-form-item>
                <template v-if="paymentSettings.wechat.enabled">
                  <n-form-item label="商户号">
                    <n-input v-model:value="paymentSettings.wechat.mchId" placeholder="微信支付商户号" />
                  </n-form-item>
                  <n-form-item label="API密钥">
                    <n-input v-model:value="paymentSettings.wechat.apiKey" type="password" show-password-on="click" placeholder="API密钥" />
                  </n-form-item>
                  <n-form-item label="证书路径">
                    <n-input v-model:value="paymentSettings.wechat.certPath" placeholder="/path/to/cert.pem" />
                  </n-form-item>
                  <n-form-item label="回调地址">
                    <n-input v-model:value="paymentSettings.wechat.notifyUrl" placeholder="https://example.com/api/pay/wechat/notify" />
                  </n-form-item>
                </template>
              </n-form>
            </n-card>
          </n-gi>
        </n-grid>
        <n-button type="primary" style="margin-top: 20px" @click="savePayment">保存支付设置</n-button>
      </n-tab-pane>

      <!-- 验证码设置 -->
      <n-tab-pane name="captcha" tab="验证码设置">
        <n-form
          :model="captchaSettings"
          label-placement="left"
          label-width="120"
          style="max-width: 640px; margin-top: 20px"
        >
          <n-divider>图形验证码</n-divider>
          <n-form-item label="启用图形验证码">
            <n-switch v-model:value="captchaSettings.imageEnabled" />
          </n-form-item>
          <template v-if="captchaSettings.imageEnabled">
            <n-form-item label="验证码长度">
              <n-input-number v-model:value="captchaSettings.imageLength" :min="4" :max="8" style="width: 160px" />
            </n-form-item>
            <n-form-item label="验证码宽度">
              <n-input-number v-model:value="captchaSettings.imageWidth" :min="80" :max="300" style="width: 160px" />
            </n-form-item>
            <n-form-item label="验证码高度">
              <n-input-number v-model:value="captchaSettings.imageHeight" :min="30" :max="100" style="width: 160px" />
            </n-form-item>
          </template>

          <n-divider>短信验证码</n-divider>
          <n-form-item label="启用短信验证码">
            <n-switch v-model:value="captchaSettings.smsEnabled" />
          </n-form-item>
          <template v-if="captchaSettings.smsEnabled">
            <n-form-item label="验证码长度">
              <n-input-number v-model:value="captchaSettings.smsLength" :min="4" :max="8" style="width: 160px" />
            </n-form-item>
            <n-form-item label="有效期(秒)">
              <n-input-number v-model:value="captchaSettings.smsExpire" :min="60" :max="600" style="width: 160px" />
            </n-form-item>
          </template>

          <n-divider>邮箱验证码</n-divider>
          <n-form-item label="启用邮箱验证码">
            <n-switch v-model:value="captchaSettings.emailEnabled" />
          </n-form-item>
          <template v-if="captchaSettings.emailEnabled">
            <n-form-item label="验证码长度">
              <n-input-number v-model:value="captchaSettings.emailLength" :min="4" :max="8" style="width: 160px" />
            </n-form-item>
            <n-form-item label="有效期(秒)">
              <n-input-number v-model:value="captchaSettings.emailExpire" :min="60" :max="600" style="width: 160px" />
            </n-form-item>
          </template>

          <n-divider>启用场景</n-divider>
          <n-form-item label="验证码场景">
            <n-checkbox-group v-model:value="captchaSettings.scenes">
              <n-checkbox value="login" label="登录" />
              <n-checkbox value="register" label="注册" />
              <n-checkbox value="reset_password" label="找回密码" />
              <n-checkbox value="contact" label="联系我们" />
            </n-checkbox-group>
          </n-form-item>

          <n-form-item>
            <n-button type="primary" @click="saveCaptcha">保存设置</n-button>
          </n-form-item>
        </n-form>
      </n-tab-pane>
    </n-tabs>
  </n-card>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useMessage } from 'naive-ui'
import { CloudUploadOutline as UploadIcon } from '@vicons/ionicons5'
import type { UploadFileInfo } from 'naive-ui'

const message = useMessage()
const activeTab = ref('basic')

// ---- Options ----
const timezoneOptions = [
  { label: 'Asia/Shanghai (UTC+8)', value: 'Asia/Shanghai' },
  { label: 'Asia/Tokyo (UTC+9)', value: 'Asia/Tokyo' },
  { label: 'America/New_York (UTC-5)', value: 'America/New_York' },
  { label: 'Europe/London (UTC+0)', value: 'Europe/London' },
  { label: 'UTC', value: 'UTC' },
]

const languageOptions = [
  { label: '简体中文', value: 'zh-CN' },
  { label: 'English', value: 'en' },
]

const emailDriverOptions = [
  { label: 'SMTP', value: 'smtp' },
  { label: 'Sendmail', value: 'sendmail' },
  { label: 'Mailgun', value: 'mailgun' },
  { label: 'SES', value: 'ses' },
]

const encryptionOptions = [
  { label: 'SSL', value: 'ssl' },
  { label: 'TLS', value: 'tls' },
  { label: '无', value: '' },
]

const smsProviderOptions = [
  { label: '阿里云短信', value: 'aliyun' },
  { label: '腾讯云短信', value: 'tencent' },
  { label: '华为云短信', value: 'huawei' },
]

// ---- Basic Settings ----
const basicSettings = reactive({
  siteName: '锚点财务',
  siteUrl: 'https://anchorfinance.com',
  siteDescription: '锚点财务 - 专业财务管理系统',
  adminEmail: 'admin@anchorfinance.com',
  timezone: 'Asia/Shanghai',
  language: 'zh-CN',
  allowRegistration: true,
  maintenanceMode: false,
  maintenanceMessage: '系统维护中，请稍后再试',
})

// ---- Email Settings ----
const emailSettings = reactive({
  driver: 'smtp',
  smtpHost: 'smtp.example.com',
  smtpPort: 465,
  encryption: 'ssl',
  username: '',
  password: '',
  fromName: '锚点财务',
  fromEmail: 'noreply@anchorfinance.com',
})

// ---- SMS Settings ----
const smsSettings = reactive({
  provider: 'aliyun',
  appKey: '',
  appSecret: '',
  signName: '锚点财务',
  templateCode: '',
  notifyTemplateCode: '',
  scenes: ['login', 'register'],
})

// ---- Payment Settings ----
const paymentSettings = reactive({
  alipay: {
    enabled: false,
    appId: '',
    privateKey: '',
    publicKey: '',
    notifyUrl: 'https://anchorfinance.com/api/pay/alipay/notify',
  },
  wechat: {
    enabled: false,
    mchId: '',
    apiKey: '',
    certPath: '',
    notifyUrl: 'https://anchorfinance.com/api/pay/wechat/notify',
  },
})

// ---- Captcha Settings ----
const captchaSettings = reactive({
  imageEnabled: true,
  imageLength: 4,
  imageWidth: 120,
  imageHeight: 40,
  smsEnabled: false,
  smsLength: 6,
  smsExpire: 300,
  emailEnabled: false,
  emailLength: 6,
  emailExpire: 300,
  scenes: ['login', 'register'],
})

// ---- Handlers ----
function handleLogoChange(options: { fileList: UploadFileInfo[] }) {
  // TODO: Upload logo file
  message.info('Logo上传功能待接入')
}

function saveBasic() {
  // TODO: API call
  message.success('基本设置已保存')
}

function saveEmail() {
  // TODO: API call
  message.success('邮件设置已保存')
}

function testEmail() {
  // TODO: API call
  message.info('测试邮件发送中...')
}

function saveSms() {
  // TODO: API call
  message.success('短信设置已保存')
}

function testSms() {
  // TODO: API call
  message.info('测试短信发送中...')
}

function savePayment() {
  // TODO: API call
  message.success('支付设置已保存')
}

function saveCaptcha() {
  // TODO: API call
  message.success('验证码设置已保存')
}
</script>

<style scoped>
.n-card {
  border-radius: 12px;
}
</style>
