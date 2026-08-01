<template>
  <div class="page-container">
    <art-card title="消息配置" shadow="never">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="邮件配置" name="email">
          <el-form :model="emailConfig" label-width="120px" style="max-width: 600px">
            <el-form-item label="SMTP服务器">
              <el-input v-model="emailConfig.host" />
            </el-form-item>
            <el-form-item label="SMTP端口">
              <el-input-number v-model="emailConfig.port" :min="1" :max="65535" />
            </el-form-item>
            <el-form-item label="发件人邮箱">
              <el-input v-model="emailConfig.from" />
            </el-form-item>
            <el-form-item label="邮箱密码">
              <el-input v-model="emailConfig.password" type="password" show-password />
            </el-form-item>
            <el-form-item label="加密方式">
              <el-select v-model="emailConfig.encryption">
                <el-option label="SSL" value="ssl" />
                <el-option label="TLS" value="tls" />
                <el-option label="无" value="" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveEmail">保存</el-button>
              <el-button @click="handleTestEmail">测试发送</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="短信配置" name="sms">
          <el-form :model="smsConfig" label-width="120px" style="max-width: 600px">
            <el-form-item label="服务商">
              <el-select v-model="smsConfig.provider">
                <el-option label="阿里云" value="aliyun" />
                <el-option label="腾讯云" value="tencent" />
                <el-option label="互亿无线" value="ihuyi" />
              </el-select>
            </el-form-item>
            <el-form-item label="AppKey">
              <el-input v-model="smsConfig.app_key" />
            </el-form-item>
            <el-form-item label="AppSecret">
              <el-input v-model="smsConfig.app_secret" type="password" show-password />
            </el-form-item>
            <el-form-item label="签名">
              <el-input v-model="smsConfig.sign" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveSms">保存</el-button>
              <el-button @click="handleTestSms">测试发送</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="微信配置" name="wechat">
          <el-form :model="wechatConfig" label-width="120px" style="max-width: 600px">
            <el-form-item label="AppID">
              <el-input v-model="wechatConfig.app_id" />
            </el-form-item>
            <el-form-item label="AppSecret">
              <el-input v-model="wechatConfig.app_secret" type="password" show-password />
            </el-form-item>
            <el-form-item label="模板ID">
              <el-input v-model="wechatConfig.template_id" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSaveWechat">保存</el-button>
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

const activeTab = ref('email')
const emailConfig = ref({ host: '', port: 465, from: '', password: '', encryption: 'ssl' })
const smsConfig = ref({ provider: 'aliyun', app_key: '', app_secret: '', sign: '' })
const wechatConfig = ref({ app_id: '', app_secret: '', template_id: '' })

const fetchConfig = async () => {
  try {
    const { data } = await request.get('/admin/config/messages')
    if (data?.data) {
      // 填充配置
    }
  } catch (error) {
    console.error(error)
  }
}

const handleSaveEmail = async () => {
  try {
    await request.put('/admin/config/messages/email', emailConfig.value)
    ElMessage.success('保存成功')
  } catch (error) {
    console.error(error)
  }
}

const handleTestEmail = async () => {
  try {
    await request.post('/admin/config/messages/email/test')
    ElMessage.success('测试邮件已发送')
  } catch (error) {
    console.error(error)
  }
}

const handleSaveSms = async () => {
  try {
    await request.put('/admin/config/messages/sms', smsConfig.value)
    ElMessage.success('保存成功')
  } catch (error) {
    console.error(error)
  }
}

const handleTestSms = async () => {
  try {
    await request.post('/admin/config/messages/sms/test')
    ElMessage.success('测试短信已发送')
  } catch (error) {
    console.error(error)
  }
}

const handleSaveWechat = async () => {
  try {
    await request.put('/admin/config/messages/wechat', wechatConfig.value)
    ElMessage.success('保存成功')
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
