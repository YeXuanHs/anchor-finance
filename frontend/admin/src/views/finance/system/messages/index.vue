<template>
  <div class="page-container">
    <art-card :title="$t('messageConfig.title')" shadow="never">
      <el-tabs v-model="activeTab">
        <el-tab-pane :label="$t('messageConfig.emailConfig')" name="email">
          <el-form :model="emailConfig" label-width="120px" style="max-width: 600px">
            <el-form-item :label="$t('messageConfig.smtpServer')">
              <el-input v-model="emailConfig.host" />
            </el-form-item>
            <el-form-item :label="$t('messageConfig.smtpPort')">
              <el-input-number v-model="emailConfig.port" :min="1" :max="65535" />
            </el-form-item>
            <el-form-item :label="$t('messageConfig.senderEmail')">
              <el-input v-model="emailConfig.from" />
            </el-form-item>
            <el-form-item :label="$t('messageConfig.emailPassword')">
              <el-input v-model="emailConfig.password" type="password" show-password />
            </el-form-item>
            <el-form-item :label="$t('messageConfig.encryption')">
              <el-select v-model="emailConfig.encryption">
                <el-option label="SSL" value="ssl" />
                <el-option label="TLS" value="tls" />
                <el-option :label="$t('messageConfig.none')" value="" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveEmail">{{ $t('messageConfig.save') }}</el-button>
              <el-button @click="handleTestEmail">{{ $t('messageConfig.testSend') }}</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane :label="$t('messageConfig.smsConfig')" name="sms">
          <el-form :model="smsConfig" label-width="120px" style="max-width: 600px">
            <el-form-item :label="$t('messageConfig.provider')">
              <el-select v-model="smsConfig.provider">
                <el-option :label="$t('messageConfig.aliyun')" value="aliyun" />
                <el-option :label="$t('messageConfig.tencent')" value="tencent" />
                <el-option :label="$t('messageConfig.ihuyi')" value="ihuyi" />
              </el-select>
            </el-form-item>
            <el-form-item label="AppKey">
              <el-input v-model="smsConfig.app_key" />
            </el-form-item>
            <el-form-item label="AppSecret">
              <el-input v-model="smsConfig.app_secret" type="password" show-password />
            </el-form-item>
            <el-form-item :label="$t('messageConfig.sign')">
              <el-input v-model="smsConfig.sign" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveSms">{{ $t('messageConfig.save') }}</el-button>
              <el-button @click="handleTestSms">{{ $t('messageConfig.testSend') }}</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane :label="$t('messageConfig.wechatConfig')" name="wechat">
          <el-form :model="wechatConfig" label-width="120px" style="max-width: 600px">
            <el-form-item label="AppID">
              <el-input v-model="wechatConfig.app_id" />
            </el-form-item>
            <el-form-item label="AppSecret">
              <el-input v-model="wechatConfig.app_secret" type="password" show-password />
            </el-form-item>
            <el-form-item :label="$t('messageConfig.templateId')">
              <el-input v-model="wechatConfig.template_id" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveWechat">{{ $t('messageConfig.save') }}</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </art-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const activeTab = ref('email')
const emailConfig = ref({ host: '', port: 465, from: '', password: '', encryption: 'ssl' })
const smsConfig = ref({ provider: 'aliyun', app_key: '', app_secret: '', sign: '' })
const wechatConfig = ref({ app_id: '', app_secret: '', template_id: '' })

const fetchConfig = async () => {
  try {
    const res = await request.get({ url: '/api/admin/config/messages' })
    if (res) {
      // 填充配置
    }
  } catch (error) {
    console.error(error)
  }
}

const handleSaveEmail = async () => {
  try {
    await request.put({ url: '/api/admin/config/messages/email', params: emailConfig.value })
    ElMessage.success($t('messageConfig.saveSuccess'))
  } catch (error) {
    console.error(error)
  }
}

const handleTestEmail = async () => {
  try {
    await request.post({ url: '/api/admin/config/messages/email/test' })
    ElMessage.success($t('messageConfig.testEmailSent'))
  } catch (error) {
    console.error(error)
  }
}

const handleSaveSms = async () => {
  try {
    await request.put({ url: '/api/admin/config/messages/sms', params: smsConfig.value })
    ElMessage.success($t('messageConfig.saveSuccess'))
  } catch (error) {
    console.error(error)
  }
}

const handleTestSms = async () => {
  try {
    await request.post({ url: '/api/admin/config/messages/sms/test' })
    ElMessage.success($t('messageConfig.testSmsSent'))
  } catch (error) {
    console.error(error)
  }
}

const handleSaveWechat = async () => {
  try {
    await request.put({ url: '/api/admin/config/messages/wechat', params: wechatConfig.value })
    ElMessage.success($t('messageConfig.saveSuccess'))
  } catch (error) {
    console.error(error)
  }
}

onMounted(() => fetchConfig())
</script>

<style scoped lang="scss">
.page-container {
  padding: 20px;
}
</style>
