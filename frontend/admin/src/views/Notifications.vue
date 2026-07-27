<template>
  <div class="admin-page">
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="card-title">通知管理</span>
          <div class="card-actions">
            <el-select v-model="filterType" placeholder="通知类型" clearable style="width: 130px">
              <el-option label="系统通知" value="system" /><el-option label="订单通知" value="order" />
              <el-option label="工单通知" value="ticket" /><el-option label="营销通知" value="promo" />
            </el-select>
            <el-button type="primary" @click="openModal()">
              <el-icon><Plus /></el-icon>发送通知
            </el-button>
          </div>
        </div>
      </template>

      <el-table :data="filteredNotifications" stripe size="small">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="title" label="标题" show-overflow-tooltip />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="typeMap[row.type]?.type as any" size="small">{{ typeMap[row.type]?.label }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="target" label="目标" width="100">
          <template #default="{ row }">{{ row.target === 'all' ? '全体用户' : row.target }}</template>
        </el-table-column>
        <el-table-column prop="channel" label="渠道" width="120">
          <template #default="{ row }">
            <el-tag v-for="ch in row.channel" :key="ch" size="small" style="margin-right: 4px">{{ channelMap[ch] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 'sent' ? 'success' : row.status === 'pending' ? 'warning' : 'info'" size="small">
              {{ statusMap[row.status] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sentAt" label="发送时间" width="160" />
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-popconfirm title="确认删除该通知？" @confirm="handleDelete(row.id)">
              <template #reference><el-button text type="danger">删除</el-button></template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="modalVisible" title="发送通知" width="600px" destroy-on-close>
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-form-item label="通知标题" prop="title"><el-input v-model="formData.title" placeholder="请输入通知标题" /></el-form-item>
        <el-form-item label="通知类型" prop="type">
          <el-select v-model="formData.type" style="width: 100%">
            <el-option label="系统通知" value="system" /><el-option label="订单通知" value="order" />
            <el-option label="工单通知" value="ticket" /><el-option label="营销通知" value="promo" />
          </el-select>
        </el-form-item>
        <el-form-item label="发送目标">
          <el-select v-model="formData.target" style="width: 100%">
            <el-option label="全体用户" value="all" /><el-option label="活跃用户" value="active" /><el-option label="VIP用户" value="vip" />
          </el-select>
        </el-form-item>
        <el-form-item label="发送渠道">
          <el-checkbox-group v-model="formData.channel">
            <el-checkbox value="site" label="站内信" /><el-checkbox value="email" label="邮件" />
            <el-checkbox value="sms" label="短信" />
          </el-checkbox-group>
        </el-form-item>
        <el-form-item label="通知内容" prop="content">
          <el-input v-model="formData.content" type="textarea" :rows="6" placeholder="通知内容" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="modalVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">发送</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'

const submitting = ref(false)
const modalVisible = ref(false)
const filterType = ref<string | null>(null)
const formRef = ref<FormInstance>()

const typeMap: Record<string, { label: string; type: string }> = {
  system: { label: '系统', type: 'danger' }, order: { label: '订单', type: 'primary' },
  ticket: { label: '工单', type: 'warning' }, promo: { label: '营销', type: 'success' },
}
const channelMap: Record<string, string> = { site: '站内', email: '邮件', sms: '短信' }
const statusMap: Record<string, string> = { sent: '已发送', pending: '待发送', failed: '失败' }

const notifications = ref([
  { id: 1, title: '系统升级通知', type: 'system', target: 'all', channel: ['site', 'email'], status: 'sent', sentAt: '2026-07-27 10:00' },
  { id: 2, title: '订单支付成功提醒', type: 'order', target: 'all', channel: ['site', 'email', 'sms'], status: 'sent', sentAt: '2026-07-27 09:30' },
  { id: 3, title: '新工单提醒', type: 'ticket', target: 'all', channel: ['site'], status: 'sent', sentAt: '2026-07-27 08:00' },
  { id: 4, title: '夏季大促活动', type: 'promo', target: 'active', channel: ['site', 'email'], status: 'pending', sentAt: '-' },
  { id: 5, title: '服务器维护通知', type: 'system', target: 'all', channel: ['site', 'email', 'sms'], status: 'sent', sentAt: '2026-07-26 18:00' },
])

const filteredNotifications = computed(() => {
  if (!filterType.value) return notifications.value
  return notifications.value.filter((n) => n.type === filterType.value)
})

const formData = reactive({ title: '', type: 'system', target: 'all', channel: ['site'] as string[], content: '' })
const rules: FormRules = {
  title: { required: true, message: '请输入通知标题', trigger: 'blur' },
  content: { required: true, message: '请输入通知内容', trigger: 'blur' },
}

function openModal() {
  Object.assign(formData, { title: '', type: 'system', target: 'all', channel: ['site'], content: '' })
  modalVisible.value = true
}

async function handleSubmit() {
  if (!formRef.value) return
  try { await formRef.value.validate() } catch { return }
  submitting.value = true
  try {
    notifications.value.unshift({ id: Date.now(), title: formData.title, type: formData.type, target: formData.target, channel: [...formData.channel], status: 'sent', sentAt: new Date().toLocaleString() })
    ElMessage.success('通知发送成功')
    modalVisible.value = false
  } finally { submitting.value = false }
}

function handleDelete(id: number) { notifications.value = notifications.value.filter((n) => n.id !== id); ElMessage.success('通知已删除') }
</script>

<style scoped>
.card-header { display: flex; align-items: center; justify-content: space-between; }
.card-title { font-size: 16px; font-weight: 600; }
.card-actions { display: flex; align-items: center; gap: 12px; }
</style>
