<template>
  <div class="ticket-settings-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('ticketSetting.title') }}</span>
          <el-button type="primary" :loading="saving" @click="handleSave">{{ $t('ticketSetting.saveSettings') }}</el-button>
        </div>
      </template>

      <el-form :model="formData" label-width="150px">
        <el-divider content-position="left">{{ $t('ticketSetting.basicSettings') }}</el-divider>
        <el-form-item :label="$t('ticketSetting.ticketPrefix')">
          <el-input v-model="formData.ticket_prefix" :placeholder="$t('ticketSetting.ticketPrefixPlaceholder')" style="width: 200px" />
        </el-form-item>
        <el-form-item :label="$t('ticketSetting.autoCloseHours')">
          <el-input-number v-model="formData.auto_close_hours" :min="0" :max="720" />
          <span class="form-tip">{{ $t('ticketSetting.autoCloseHoursTip') }}</span>
        </el-form-item>
        <el-form-item :label="$t('ticketSetting.clientCanClose')">
          <el-switch v-model="formData.client_can_close" />
        </el-form-item>
        <el-form-item :label="$t('ticketSetting.clientCanRate')">
          <el-switch v-model="formData.client_can_rate" />
        </el-form-item>

        <el-divider content-position="left">{{ $t('ticketSetting.notificationSettings') }}</el-divider>
        <el-form-item :label="$t('ticketSetting.notifyNewTicket')">
          <el-switch v-model="formData.notify_new_ticket" />
        </el-form-item>
        <el-form-item :label="$t('ticketSetting.notifyReply')">
          <el-switch v-model="formData.notify_reply" />
        </el-form-item>
        <el-form-item :label="$t('ticketSetting.notifyAdminEmail')">
          <el-switch v-model="formData.notify_admin_email" />
        </el-form-item>

        <el-divider content-position="left">{{ $t('ticketSetting.ticketStatus') }}</el-divider>
        <el-form-item :label="$t('ticketSetting.statusWaiting')">
          <el-input v-model="formData.status_waiting" :placeholder="$t('ticketSetting.statusWaitingPlaceholder')" style="width: 200px" />
        </el-form-item>
        <el-form-item :label="$t('ticketSetting.statusReplied')">
          <el-input v-model="formData.status_replied" :placeholder="$t('ticketSetting.statusRepliedPlaceholder')" style="width: 200px" />
        </el-form-item>
        <el-form-item :label="$t('ticketSetting.statusClosed')">
          <el-input v-model="formData.status_closed" :placeholder="$t('ticketSetting.statusClosedPlaceholder')" style="width: 200px" />
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

const saving = ref(false)
const formData = reactive({
  ticket_prefix: 'TK-',
  auto_close_hours: 72,
  client_can_close: true,
  client_can_rate: true,
  notify_new_ticket: true,
  notify_reply: true,
  notify_admin_email: true,
  status_waiting: $t('ticketSetting.statusWaitingDefault'),
  status_replied: $t('ticketSetting.statusRepliedDefault'),
  status_closed: $t('ticketSetting.statusClosedDefault')
})

const fetchSettings = async () => {
  try {
    const data = await request.get({ url: '/api/admin/settings' })
    Object.assign(formData, data)
  } catch (error) {
    console.error($t('ticketSetting.fetchFailed'), error)
  }
}

const handleSave = async () => {
  saving.value = true
  try {
    await request.put({ url: '/api/admin/settings', data: formData })
    ElMessage.success($t('common.saveSuccess'))
  } catch (error) {
    console.error($t('ticketSetting.saveFailed'), error)
  } finally {
    saving.value = false
  }
}

onMounted(() => { fetchSettings() })
</script>

<style scoped lang="scss">
.ticket-settings-page { padding: 16px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.form-tip { margin-left: 10px; font-size: 12px; color: #86909C; }
:deep(.el-divider__text) { font-weight: 600; color: #1D2129; }
</style>
