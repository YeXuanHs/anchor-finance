<template>
  <div class="settings-page">
    <el-card class="admin-card" shadow="never">
      <el-tabs v-model="activeTab">
        <!-- Basic Settings -->
        <el-tab-pane label="基本设置" name="basic">
          <el-form :model="basicSettings" label-width="120px" style="max-width: 640px; margin-top: 20px">
            <el-form-item label="站点名称">
              <el-input v-model="basicSettings.siteName" placeholder="锚点财务" />
            </el-form-item>
            <el-form-item label="站点Logo">
              <el-upload action="#" :auto-upload="false" :limit="1" accept="image/*">
                <el-button :icon="Upload">上传Logo</el-button>
              </el-upload>
            </el-form-item>
            <el-form-item label="站点描述">
              <el-input v-model="basicSettings.siteDescription" type="textarea" :rows="3" placeholder="站点描述信息" />
            </el-form-item>
            <el-form-item label="站点URL">
              <el-input v-model="basicSettings.siteUrl" placeholder="https://example.com" />
            </el-form-item>
            <el-form-item label="管理员邮箱">
              <el-input v-model="basicSettings.adminEmail" placeholder="admin@example.com" />
            </el-form-item>
            <el-form-item label="时区">
              <el-select v-model="basicSettings.timezone" style="width: 100%">
                <el-option v-for="opt in timezoneOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
              </el-select>
            </el-form-item>
            <el-form-item label="默认语言">
              <el-select v-model="basicSettings.language" style="width: 100%">
                <el-option v-for="opt in languageOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
              </el-select>
            </el-form-item>
            <el-form-item label="用户注册">
              <el-switch v-model="basicSettings.allowRegistration" />
            </el-form-item>
            <el-form-item label="维护模式">
              <el-switch v-model="basicSettings.maintenanceMode" />
            </el-form-item>
            <el-form-item label="维护提示" v-if="basicSettings.maintenanceMode">
              <el-input v-model="basicSettings.maintenanceMessage" placeholder="系统维护中，请稍后再试" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveBasic">保存设置</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- Email Settings -->
        <el-tab-pane label="邮件设置" name="email">
          <el-form :model="emailSettings" label-width="120px" style="max-width: 640px; margin-top: 20px">
            <el-form-item label="邮件驱动">
              <el-select v-model="emailSettings.driver" style="width: 100%">
                <el-option v-for="opt in emailDriverOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
              </el-select>
            </el-form-item>
            <el-form-item label="SMTP 主机">
              <el-input v-model="emailSettings.smtpHost" placeholder="smtp.example.com" />
            </el-form-item>
            <el-form-item label="SMTP 端口">
              <el-input-number v-model="emailSettings.smtpPort" :min="1" :max="65535" style="width: 100%" />
            </el-form-item>
            <el-form-item label="加密方式">
              <el-select v-model="emailSettings.encryption" style="width: 100%">
                <el-option v-for="opt in encryptionOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
              </el-select>
            </el-form-item>
            <el-form-item label="用户名">
              <el-input v-model="emailSettings.username" placeholder="noreply@example.com" />
            </el-form-item>
            <el-form-item label="密码">
              <el-input v-model="emailSettings.password" type="password" show-password placeholder="SMTP密码" />
            </el-form-item>
            <el-form-item label="发件人名称">
              <el-input v-model="emailSettings.fromName" placeholder="锚点财务" />
            </el-form-item>
            <el-form-item label="发件人邮箱">
              <el-input v-model="emailSettings.fromEmail" placeholder="noreply@example.com" />
            </el-form-item>
            <el-form-item>
              <el-space>
                <el-button type="primary" @click="saveEmail">保存设置</el-button>
                <el-button @click="testEmail">发送测试邮件</el-button>
              </el-space>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- SMS Settings -->
        <el-tab-pane label="短信设置" name="sms">
          <el-form :model="smsSettings" label-width="120px" style="max-width: 640px; margin-top: 20px">
            <el-form-item label="短信服务商">
              <el-select v-model="smsSettings.provider" style="width: 100%">
                <el-option v-for="opt in smsProviderOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
              </el-select>
            </el-form-item>
            <el-form-item label="App Key">
              <el-input v-model="smsSettings.appKey" placeholder="短信服务商AppKey" />
            </el-form-item>
            <el-form-item label="App Secret">
              <el-input v-model="smsSettings.appSecret" type="password" show-password placeholder="短信服务商AppSecret" />
            </el-form-item>
            <el-form-item label="短信签名">
              <el-input v-model="smsSettings.signName" placeholder="短信签名" />
            </el-form-item>
            <el-divider>短信模板</el-divider>
            <el-form-item label="验证码模板">
              <el-input v-model="smsSettings.templateCode" placeholder="SMS_XXXXXX" />
            </el-form-item>
            <el-form-item label="通知模板">
              <el-input v-model="smsSettings.notifyTemplateCode" placeholder="SMS_XXXXXX" />
            </el-form-item>
            <el-form-item label="启用场景">
              <el-checkbox-group v-model="smsSettings.scenes">
                <el-checkbox value="login" label="登录验证" />
                <el-checkbox value="register" label="注册验证" />
                <el-checkbox value="reset" label="密码重置" />
                <el-checkbox value="notify" label="订单通知" />
              </el-checkbox-group>
            </el-form-item>
            <el-form-item>
              <el-space>
                <el-button type="primary" @click="saveSms">保存设置</el-button>
                <el-button @click="testSms">发送测试短信</el-button>
              </el-space>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- Payment Settings -->
        <el-tab-pane label="支付设置" name="payment">
          <el-row :gutter="24" style="margin-top: 20px">
            <el-col :span="12">
              <el-card shadow="never" style="background: #fafafa; border-radius: 12px">
                <template #header>支付宝配置</template>
                <el-form :model="paymentSettings.alipay" label-width="100px">
                  <el-form-item label="启用支付宝">
                    <el-switch v-model="paymentSettings.alipay.enabled" />
                  </el-form-item>
                  <template v-if="paymentSettings.alipay.enabled">
                    <el-form-item label="应用ID">
                      <el-input v-model="paymentSettings.alipay.appId" placeholder="支付宝应用ID" />
                    </el-form-item>
                    <el-form-item label="私钥">
                      <el-input v-model="paymentSettings.alipay.privateKey" type="textarea" :rows="3" placeholder="应用私钥" />
                    </el-form-item>
                    <el-form-item label="支付宝公钥">
                      <el-input v-model="paymentSettings.alipay.publicKey" type="textarea" :rows="3" placeholder="支付宝公钥" />
                    </el-form-item>
                    <el-form-item label="回调地址">
                      <el-input v-model="paymentSettings.alipay.notifyUrl" placeholder="https://example.com/api/pay/alipay/notify" />
                    </el-form-item>
                  </template>
                </el-form>
              </el-card>
            </el-col>
            <el-col :span="12">
              <el-card shadow="never" style="background: #fafafa; border-radius: 12px">
                <template #header>微信支付配置</template>
                <el-form :model="paymentSettings.wechat" label-width="100px">
                  <el-form-item label="启用微信">
                    <el-switch v-model="paymentSettings.wechat.enabled" />
                  </el-form-item>
                  <template v-if="paymentSettings.wechat.enabled">
                    <el-form-item label="商户号">
                      <el-input v-model="paymentSettings.wechat.mchId" placeholder="微信支付商户号" />
                    </el-form-item>
                    <el-form-item label="API密钥">
                      <el-input v-model="paymentSettings.wechat.apiKey" type="password" show-password placeholder="API密钥" />
                    </el-form-item>
                    <el-form-item label="证书路径">
                      <el-input v-model="paymentSettings.wechat.certPath" placeholder="/path/to/cert.pem" />
                    </el-form-item>
                    <el-form-item label="回调地址">
                      <el-input v-model="paymentSettings.wechat.notifyUrl" placeholder="https://example.com/api/pay/wechat/notify" />
                    </el-form-item>
                  </template>
                </el-form>
              </el-card>
            </el-col>
          </el-row>
          <el-button type="primary" style="margin-top: 20px" @click="savePayment">保存支付设置</el-button>
        </el-tab-pane>

        <!-- Captcha Settings -->
        <el-tab-pane label="验证码设置" name="captcha">
          <el-form :model="captchaSettings" label-width="120px" style="max-width: 640px; margin-top: 20px">
            <el-divider>图形验证码</el-divider>
            <el-form-item label="启用图形验证码">
              <el-switch v-model="captchaSettings.imageEnabled" />
            </el-form-item>
            <template v-if="captchaSettings.imageEnabled">
              <el-form-item label="验证码长度">
                <el-input-number v-model="captchaSettings.imageLength" :min="4" :max="8" />
              </el-form-item>
              <el-form-item label="验证码宽度">
                <el-input-number v-model="captchaSettings.imageWidth" :min="80" :max="300" />
              </el-form-item>
              <el-form-item label="验证码高度">
                <el-input-number v-model="captchaSettings.imageHeight" :min="30" :max="100" />
              </el-form-item>
            </template>

            <el-divider>短信验证码</el-divider>
            <el-form-item label="启用短信验证码">
              <el-switch v-model="captchaSettings.smsEnabled" />
            </el-form-item>
            <template v-if="captchaSettings.smsEnabled">
              <el-form-item label="验证码长度">
                <el-input-number v-model="captchaSettings.smsLength" :min="4" :max="8" />
              </el-form-item>
              <el-form-item label="有效期(秒)">
                <el-input-number v-model="captchaSettings.smsExpire" :min="60" :max="600" />
              </el-form-item>
            </template>

            <el-divider>邮箱验证码</el-divider>
            <el-form-item label="启用邮箱验证码">
              <el-switch v-model="captchaSettings.emailEnabled" />
            </el-form-item>
            <template v-if="captchaSettings.emailEnabled">
              <el-form-item label="验证码长度">
                <el-input-number v-model="captchaSettings.emailLength" :min="4" :max="8" />
              </el-form-item>
              <el-form-item label="有效期(秒)">
                <el-input-number v-model="captchaSettings.emailExpire" :min="60" :max="600" />
              </el-form-item>
            </template>

            <el-divider>启用场景</el-divider>
            <el-form-item label="验证码场景">
              <el-checkbox-group v-model="captchaSettings.scenes">
                <el-checkbox value="login" label="登录" />
                <el-checkbox value="register" label="注册" />
                <el-checkbox value="reset_password" label="找回密码" />
                <el-checkbox value="contact" label="联系我们" />
              </el-checkbox-group>
            </el-form-item>

            <el-form-item>
              <el-button type="primary" @click="saveCaptcha">保存设置</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { Upload } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

definePageMeta({
  layout: 'admin',
})

const activeTab = ref('basic')

// Options
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

// Basic Settings
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

// Email Settings
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

// SMS Settings
const smsSettings = reactive({
  provider: 'aliyun',
  appKey: '',
  appSecret: '',
  signName: '锚点财务',
  templateCode: '',
  notifyTemplateCode: '',
  scenes: ['login', 'register'],
})

// Payment Settings
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

// Captcha Settings
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

// Handlers
function saveBasic() { ElMessage.success('基本设置已保存') }
function saveEmail() { ElMessage.success('邮件设置已保存') }
function testEmail() { ElMessage.info('测试邮件发送中...') }
function saveSms() { ElMessage.success('短信设置已保存') }
function testSms() { ElMessage.info('测试短信发送中...') }
function savePayment() { ElMessage.success('支付设置已保存') }
function saveCaptcha() { ElMessage.success('验证码设置已保存') }
</script>

<style scoped>
.settings-page {
  padding: 0;
}
</style>
