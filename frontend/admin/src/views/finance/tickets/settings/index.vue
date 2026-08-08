<template>
  <div class="ticket-settings-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>工单设置</span>
          <el-button type="primary" :loading="saving" @click="handleSave">保存设置</el-button>
        </div>
      </template>

      <el-form :model="formData" label-width="150px">
        <el-divider content-position="left">基础设置</el-divider>
        <el-form-item label="工单前缀">
          <el-input v-model="formData.ticket_prefix" placeholder="如 TK-" style="width: 200px" />
        </el-form-item>
        <el-form-item label="自动关闭时间">
          <el-input-number v-model="formData.auto_close_hours" :min="0" :max="720" />
          <span class="form-tip">小时（0表示不自动关闭）</span>
        </el-form-item>
        <el-form-item label="客户可关闭工单">
          <el-switch v-model="formData.client_can_close" />
        </el-form-item>
        <el-form-item label="客户可评价工单">
          <el-switch v-model="formData.client_can_rate" />
        </el-form-item>

        <el-divider content-position="left">通知设置</el-divider>
        <el-form-item label="新工单通知">
          <el-switch v-model="formData.notify_new_ticket" />
        </el-form-item>
        <el-form-item label="回复通知客户">
          <el-switch v-model="formData.notify_reply" />
        </el-form-item>
        <el-form-item label="邮件通知管理员">
          <el-switch v-model="formData.notify_admin_email" />
        </el-form-item>

        <el-divider content-position="left">工单状态</el-divider>
        <el-form-item label="等待回复状态名">
          <el-input v-model="formData.status_waiting" placeholder="待回复" style="width: 200px" />
        </el-form-item>
        <el-form-item label="已回复状态名">
          <el-input v-model="formData.status_replied" placeholder="已回复" style="width: 200px" />
        </el-form-item>
        <el-form-item label="已关闭状态名">
          <el-input v-model="formData.status_closed" placeholder="已关闭" style="width: 200px" />
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
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
  status_waiting: '待回复',
  status_replied: '已回复',
  status_closed: '已关闭'
})

const fetchSettings = async () => {
  try {
    const data = await request.get({ url: '/api/admin/settings/tickets' })
    Object.assign(formData, data)
  } catch (error) {
    console.error('获取工单设置失败:', error)
  }
}

const handleSave = async () => {
  saving.value = true
  try {
    await request.post({ url: '/api/admin/settings/tickets', data: formData })
    ElMessage.success('保存成功')
  } catch (error) {
    console.error('保存失败:', error)
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
