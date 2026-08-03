<template>
  <div class="user-tastes-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>用户偏好设置</span>
        </div>
      </template>

      <el-form :model="configForm" ref="formRef" label-width="120px" class="config-form">
        <el-divider content-position="left">语言与区域</el-divider>
        <el-form-item label="默认语言">
          <el-select v-model="configForm.language" placeholder="请选择默认语言">
            <el-option label="简体中文" value="zh-cn" />
            <el-option label="English" value="en-us" />
            <el-option label="繁體中文" value="zh-hk" />
            <el-option label="日本語" value="ja" />
            <el-option label="한국어" value="ko" />
          </el-select>
        </el-form-item>
        <el-form-item label="时区">
          <el-select v-model="configForm.timezone" placeholder="请选择时区">
            <el-option label="Asia/Shanghai (UTC+8)" value="Asia/Shanghai" />
            <el-option label="Asia/Tokyo (UTC+9)" value="Asia/Tokyo" />
            <el-option label="America/New_York (UTC-5)" value="America/New_York" />
            <el-option label="America/Los_Angeles (UTC-8)" value="America/Los_Angeles" />
            <el-option label="Europe/London (UTC+0)" value="Europe/London" />
            <el-option label="Europe/Berlin (UTC+1)" value="Europe/Berlin" />
          </el-select>
        </el-form-item>
        <el-form-item label="日期格式">
          <el-select v-model="configForm.date_format" placeholder="请选择日期格式">
            <el-option label="YYYY-MM-DD" value="YYYY-MM-DD" />
            <el-option label="DD/MM/YYYY" value="DD/MM/YYYY" />
            <el-option label="MM/DD/YYYY" value="MM/DD/YYYY" />
            <el-option label="YYYY年MM月DD日" value="YYYY年MM月DD日" />
          </el-select>
        </el-form-item>

        <el-divider content-position="left">通知偏好</el-divider>
        <el-form-item label="通知方式">
          <el-checkbox-group v-model="configForm.notification_methods">
            <el-checkbox value="email">邮件通知</el-checkbox>
            <el-checkbox value="sms">短信通知</el-checkbox>
            <el-checkbox value="wechat">微信通知</el-checkbox>
            <el-checkbox value="site">站内通知</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        <el-form-item label="新工单通知">
          <el-switch v-model="configForm.notify_new_ticket" />
        </el-form-item>
        <el-form-item label="工单回复通知">
          <el-switch v-model="configForm.notify_ticket_reply" />
        </el-form-item>
        <el-form-item label="新订单通知">
          <el-switch v-model="configForm.notify_new_order" />
        </el-form-item>
        <el-form-item label="账单到期提醒">
          <el-switch v-model="configForm.notify_invoice_due" />
        </el-form-item>

        <el-divider content-position="left">界面偏好</el-divider>
        <el-form-item label="每页显示条数">
          <el-select v-model="configForm.page_size" placeholder="请选择">
            <el-option label="10条/页" :value="10" />
            <el-option label="20条/页" :value="20" />
            <el-option label="50条/页" :value="50" />
            <el-option label="100条/页" :value="100" />
          </el-select>
        </el-form-item>
        <el-form-item label="主题">
          <el-radio-group v-model="configForm.theme">
            <el-radio value="light">浅色</el-radio>
            <el-radio value="dark">深色</el-radio>
            <el-radio value="auto">跟随系统</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSave" :loading="saveLoading">保存设置</el-button>
          <el-button @click="handleReset">恢复默认</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/http'

const saveLoading = ref(false)

const configForm = reactive({
  language: 'zh-cn',
  timezone: 'Asia/Shanghai',
  date_format: 'YYYY-MM-DD',
  notification_methods: ['email', 'site'] as string[],
  notify_new_ticket: true,
  notify_ticket_reply: true,
  notify_new_order: true,
  notify_invoice_due: true,
  page_size: 20,
  theme: 'light'
})

const fetchConfig = async () => {
  try {
    const data = await request.get({ url: '/api/admin/user-tastes' })
    if (data) {
      Object.assign(configForm, data)
    }
  } catch (error) {
    console.error('获取用户偏好失败:', error)
  }
}

const handleSave = async () => {
  saveLoading.value = true
  try {
    await request.put({ url: '/api/admin/user-tastes', params: { ...configForm } })
    ElMessage.success('保存成功')
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saveLoading.value = false
  }
}

const handleReset = () => {
  configForm.language = 'zh-cn'
  configForm.timezone = 'Asia/Shanghai'
  configForm.date_format = 'YYYY-MM-DD'
  configForm.notification_methods = ['email', 'site']
  configForm.notify_new_ticket = true
  configForm.notify_ticket_reply = true
  configForm.notify_new_order = true
  configForm.notify_invoice_due = true
  configForm.page_size = 20
  configForm.theme = 'light'
}

onMounted(() => {
  fetchConfig()
})
</script>

<style scoped lang="scss">
.user-tastes-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.config-form {
  max-width: 600px;
}
</style>
