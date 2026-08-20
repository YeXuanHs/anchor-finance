<template>
  <div class="user-tastes-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('page.userTastes.title') }}</span>
        </div>
      </template>

      <el-form :model="configForm" ref="formRef" label-width="120px" class="config-form">
        <el-divider content-position="left">{{ $t('page.userTastes.languageAndRegion') }}</el-divider>
        <el-form-item :label="$t('page.userTastes.defaultLanguage')">
          <el-select v-model="configForm.language" :placeholder="$t('page.userTastes.selectLanguage')">
            <el-option :label="$t('page.userTastes.simplifiedChinese')" value="zh-cn" />
            <el-option label="English" value="en-us" />
            <el-option :label="$t('page.userTastes.traditionalChinese')" value="zh-hk" />
            <el-option label="日本語" value="ja" />
            <el-option label="한국어" value="ko" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('page.userTastes.timezone')">
          <el-select v-model="configForm.timezone" :placeholder="$t('page.userTastes.selectTimezone')">
            <el-option label="Asia/Shanghai (UTC+8)" value="Asia/Shanghai" />
            <el-option label="Asia/Tokyo (UTC+9)" value="Asia/Tokyo" />
            <el-option label="America/New_York (UTC-5)" value="America/New_York" />
            <el-option label="America/Los_Angeles (UTC-8)" value="America/Los_Angeles" />
            <el-option label="Europe/London (UTC+0)" value="Europe/London" />
            <el-option label="Europe/Berlin (UTC+1)" value="Europe/Berlin" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('page.userTastes.dateFormat')">
          <el-select v-model="configForm.date_format" :placeholder="$t('page.userTastes.selectDateFormat')">
            <el-option label="YYYY-MM-DD" value="YYYY-MM-DD" />
            <el-option label="DD/MM/YYYY" value="DD/MM/YYYY" />
            <el-option label="MM/DD/YYYY" value="MM/DD/YYYY" />
            <el-option :label="$t('page.userTastes.dateFormatCn')" value="YYYY年MM月DD日" />
          </el-select>
        </el-form-item>

        <el-divider content-position="left">{{ $t('page.userTastes.notificationPreferences') }}</el-divider>
        <el-form-item :label="$t('page.userTastes.notificationMethod')">
          <el-checkbox-group v-model="configForm.notification_methods">
            <el-checkbox value="email">{{ $t('page.userTastes.emailNotification') }}</el-checkbox>
            <el-checkbox value="sms">{{ $t('page.userTastes.smsNotification') }}</el-checkbox>
            <el-checkbox value="wechat">{{ $t('page.userTastes.wechatNotification') }}</el-checkbox>
            <el-checkbox value="site">{{ $t('page.userTastes.siteNotification') }}</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        <el-form-item :label="$t('page.userTastes.newTicketNotification')">
          <el-switch v-model="configForm.notify_new_ticket" />
        </el-form-item>
        <el-form-item :label="$t('page.userTastes.ticketReplyNotification')">
          <el-switch v-model="configForm.notify_ticket_reply" />
        </el-form-item>
        <el-form-item :label="$t('page.userTastes.newOrderNotification')">
          <el-switch v-model="configForm.notify_new_order" />
        </el-form-item>
        <el-form-item :label="$t('page.userTastes.invoiceDueReminder')">
          <el-switch v-model="configForm.notify_invoice_due" />
        </el-form-item>

        <el-divider content-position="left">{{ $t('page.userTastes.uiPreferences') }}</el-divider>
        <el-form-item :label="$t('page.userTastes.pageSize')">
          <el-select v-model="configForm.page_size" :placeholder="$t('page.userTastes.pleaseSelect')">
            <el-option :label="'10' + $t('page.userTastes.itemsPerPage')" :value="10" />
            <el-option :label="'20' + $t('page.userTastes.itemsPerPage')" :value="20" />
            <el-option :label="'50' + $t('page.userTastes.itemsPerPage')" :value="50" />
            <el-option :label="'100' + $t('page.userTastes.itemsPerPage')" :value="100" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('page.userTastes.theme')">
          <el-radio-group v-model="configForm.theme">
            <el-radio value="light">{{ $t('page.userTastes.light') }}</el-radio>
            <el-radio value="dark">{{ $t('page.userTastes.dark') }}</el-radio>
            <el-radio value="auto">{{ $t('page.userTastes.followSystem') }}</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSave" :loading="saveLoading">{{ $t('page.userTastes.save') }}</el-button>
          <el-button @click="handleReset">{{ $t('page.userTastes.resetDefault') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { $t } from '@/locales'
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
    ElMessage.success($t('page.userTastes.saveSuccess'))
  } catch (error) {
    ElMessage.error($t('page.userTastes.saveFailed'))
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
