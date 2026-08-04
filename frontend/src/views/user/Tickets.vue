<template>
  <div class="tickets-page">
    <div class="page-header">
      <h1 class="page-title">工单列表</h1>
      <el-button type="primary" @click="showCreateDialog = true">
        <el-icon><Plus /></el-icon>提交工单
      </el-button>
    </div>

    <el-radio-group v-model="statusFilter" class="status-filter">
      <el-radio-button value="all">全部</el-radio-button>
      <el-radio-button value="open">处理中</el-radio-button>
      <el-radio-button value="replied">已回复</el-radio-button>
      <el-radio-button value="closed">已关闭</el-radio-button>
    </el-radio-group>

    <div class="ticket-list">
      <el-card
        v-for="ticket in filteredTickets"
        :key="ticket.id"
        shadow="never"
        class="ticket-card"
        @click="handleView(ticket)"
      >
        <div class="ticket-header">
          <div class="ticket-left">
            <el-tag :type="getStatusType(ticket.status)" size="small" effect="light" round>
              {{ ticket.statusText }}
            </el-tag>
            <span class="ticket-id">#{{ ticket.id }}</span>
          </div>
          <span class="ticket-time">{{ ticket.updatedAt }}</span>
        </div>
        <h3 class="ticket-title">{{ ticket.title }}</h3>
        <p class="ticket-desc">{{ ticket.description }}</p>
        <div class="ticket-footer">
          <el-tag type="info" size="small" effect="plain">{{ ticket.department }}</el-tag>
          <span class="ticket-priority">
            优先级：<el-tag :type="getPriorityType(ticket.priority)" size="small" effect="plain" round>{{ ticket.priorityText }}</el-tag>
          </span>
        </div>
      </el-card>
    </div>

    <el-empty v-if="filteredTickets.length === 0" description="暂无工单">
      <el-button type="primary" @click="showCreateDialog = true">提交工单</el-button>
    </el-empty>

    <el-dialog v-model="showCreateDialog" title="提交工单" width="560px" destroy-on-close>
      <el-form :model="newTicket" label-width="80px">
        <el-form-item label="工单标题" required>
          <el-input v-model="newTicket.title" placeholder="请简要描述您的问题" />
        </el-form-item>
        <el-form-item label="部门" required>
          <el-select v-model="newTicket.department" placeholder="请选择部门" style="width: 100%;">
            <el-option label="技术支持" value="技术支持" />
            <el-option label="财务部门" value="财务部门" />
            <el-option label="售前咨询" value="售前咨询" />
            <el-option label="投诉建议" value="投诉建议" />
          </el-select>
        </el-form-item>
        <el-form-item label="优先级">
          <el-select v-model="newTicket.priority" style="width: 100%;">
            <el-option label="低" value="low" />
            <el-option label="中" value="medium" />
            <el-option label="高" value="high" />
            <el-option label="紧急" value="urgent" />
          </el-select>
        </el-form-item>
        <el-form-item label="问题描述" required>
          <el-input v-model="newTicket.description" type="textarea" :rows="5" placeholder="请详细描述您遇到的问题" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="handleCreateTicket">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import request from '@/utils/request'

const statusFilter = ref('all')
const showCreateDialog = ref(false)
const loading = ref(false)

interface Ticket {
  id: number
  title: string
  description: string
  status: string
  statusText: string
  department: string
  priority: string
  priorityText: string
  updatedAt: string
}

const tickets = ref<Ticket[]>([])

const newTicket = reactive({ title: '', department: '', priority: 'medium', description: '' })

onMounted(async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/v2/tickets')
    tickets.value = data.data?.list || data.list || data.data || []
  } catch (e) { console.error(e) } finally { loading.value = false }
})

const filteredTickets = computed(() => {
  if (statusFilter.value === 'all') return tickets.value
  return tickets.value.filter(t => t.status === statusFilter.value)
})

function getStatusType(status: string) {
  const map: Record<string, 'success' | 'warning' | 'info' | 'danger'> = {
    open: 'warning', replied: 'success', closed: 'info'
  }
  return map[status] || 'info'
}

function getPriorityType(priority: string) {
  const map: Record<string, 'success' | 'warning' | 'info' | 'danger'> = {
    low: 'info', medium: 'warning', high: 'danger', urgent: 'danger'
  }
  return map[priority] || 'info'
}

function handleView(ticket: Ticket) { ElMessage.info(`查看工单：#${ticket.id}`) }

async function handleCreateTicket() {
  if (!newTicket.title || !newTicket.description) { ElMessage.warning('请填写完整信息'); return }
  try {
    await request.post('/api/v2/tickets', {
    showCreateDialog.value = false
    ElMessage.success('工单已提交')
    newTicket.title = ''
    newTicket.department = ''
    newTicket.priority = 'medium'
    newTicket.description = ''
    const { data } = await request.get('/api/v2/tickets')
    tickets.value = data.data?.list || data.list || data.data || []
  } catch (e: any) { ElMessage.error(e?.message || '提交失败，请重试') }
}
</script>

<style scoped>
.tickets-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.page-title {
  font-size: 20px;
  font-weight: 700;
  color: #303133;
  margin: 0;
}

.status-filter { align-self: flex-start; }

.ticket-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.ticket-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
  cursor: pointer;
  transition: all 0.2s;
  background: #fff;
}

.ticket-card:hover {
  border-color: #0056FF;
  box-shadow: 0 2px 12px rgba(0,86,255,0.08);
}

.ticket-card :deep(.el-card__body) {
  padding: 18px 20px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.ticket-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.ticket-left { display: flex; align-items: center; gap: 10px; }
.ticket-id { font-size: 13px; color: #909399; font-family: 'Monaco', 'Menlo', monospace; }
.ticket-time { font-size: 12px; color: #c0c4cc; }
.ticket-title { font-size: 16px; font-weight: 600; color: #303133; margin: 0; }
.ticket-desc {
  font-size: 14px;
  color: #606266;
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.ticket-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.ticket-priority {
  font-size: 13px;
  color: #909399;
  display: flex;
  align-items: center;
  gap: 4px;
}

@media (max-width: 768px) {
  .page-header { flex-direction: column; gap: 12px; align-items: flex-start; }
}
</style>
