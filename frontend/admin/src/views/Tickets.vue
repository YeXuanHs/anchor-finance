<template>
  <div class="admin-page">
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="card-title">工单管理</span>
          <div class="card-actions">
            <el-select v-model="filters.status" placeholder="工单状态" clearable style="width: 130px">
              <el-option v-for="o in statusOptions" :key="o.value" :label="o.label" :value="o.value" />
            </el-select>
            <el-select v-model="filters.priority" placeholder="优先级" clearable style="width: 110px">
              <el-option v-for="o in priorityOptions" :key="o.value" :label="o.label" :value="o.value" />
            </el-select>
            <el-input v-model="filters.keyword" placeholder="搜索工单标题/用户" clearable style="width: 200px" @keydown.enter="handleSearch">
              <template #prefix><el-icon><Search /></el-icon></template>
            </el-input>
            <el-button type="primary" @click="handleSearch">搜索</el-button>
          </div>
        </div>
      </template>

      <el-table :data="tickets" v-loading="loading" stripe size="small">
        <el-table-column prop="ticketNo" label="工单号" width="130" />
        <el-table-column prop="title" label="标题" show-overflow-tooltip />
        <el-table-column prop="user" label="用户" width="80" />
        <el-table-column prop="priority" label="优先级" width="80">
          <template #default="{ row }">
            <el-tag :type="priorityMap[row.priority]?.type as any" size="small">{{ priorityMap[row.priority]?.label }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="statusMap[row.status]?.type as any" size="small" round>{{ statusMap[row.status]?.label }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="assignee" label="指派给" width="100">
          <template #default="{ row }">{{ row.assignee || '未指派' }}</template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="160" />
        <el-table-column label="操作" width="80" fixed="right">
          <template #default="{ row }">
            <el-button text type="primary" @click="openDrawer(row)">处理</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="20"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @current-change="(p: number) => pagination.page = p"
          @size-change="(s: number) => { pagination.pageSize = s; pagination.page = 1 }"
        />
      </div>
    </el-card>

    <!-- Ticket Detail Drawer -->
    <el-drawer v-model="drawerVisible" :title="`工单详情 - ${currentTicket?.title || ''}`" size="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="工单号">{{ currentTicket?.ticketNo }}</el-descriptions-item>
        <el-descriptions-item label="提交用户">{{ currentTicket?.user }}</el-descriptions-item>
        <el-descriptions-item label="优先级">
          <el-tag :type="priorityMap[currentTicket?.priority || '']?.type as any" size="small">{{ priorityMap[currentTicket?.priority || '']?.label }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="statusMap[currentTicket?.status || '']?.type as any" size="small" round>{{ statusMap[currentTicket?.status || '']?.label }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间" :span="2">{{ currentTicket?.createdAt }}</el-descriptions-item>
        <el-descriptions-item label="问题描述" :span="2">{{ currentTicket?.description }}</el-descriptions-item>
      </el-descriptions>

      <el-divider>对话记录</el-divider>
      <div class="conversation-list">
        <div v-for="(msg, idx) in currentTicket?.messages || []" :key="idx" class="message-item" :class="{ admin: msg.isAdmin }">
          <el-avatar :size="32" :style="{ backgroundColor: msg.isAdmin ? '#52c41a' : '#0056FF' }">{{ msg.sender.charAt(0) }}</el-avatar>
          <div class="message-content">
            <div class="message-header">
              <span class="sender">{{ msg.sender }}</span>
              <span class="time">{{ msg.time }}</span>
            </div>
            <div class="message-text">{{ msg.content }}</div>
          </div>
        </div>
      </div>

      <el-divider>回复</el-divider>
      <el-input v-model="replyContent" type="textarea" :rows="3" placeholder="输入回复内容..." />

      <template #footer>
        <el-space>
          <el-select v-model="assignAdminId" placeholder="指派给" clearable style="width: 140px">
            <el-option v-for="o in adminOptions" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
          <el-button type="primary" @click="handleAssign">指派</el-button>
          <el-button type="success" :disabled="!replyContent.trim()" @click="handleReply">回复</el-button>
          <el-button v-if="currentTicket?.status !== 'closed'" type="warning" @click="handleClose">关闭工单</el-button>
        </el-space>
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'

const loading = ref(false)
const drawerVisible = ref(false)
const currentTicket = ref<any>(null)
const replyContent = ref('')
const assignAdminId = ref<number | null>(null)

const filters = reactive({ status: null as string | null, priority: null as string | null, keyword: '' })
const pagination = reactive({ page: 1, pageSize: 10 })

const statusOptions = [
  { label: '待处理', value: 'open' }, { label: '处理中', value: 'processing' },
  { label: '已回复', value: 'replied' }, { label: '已关闭', value: 'closed' },
]
const priorityOptions = [
  { label: '低', value: 'low' }, { label: '中', value: 'medium' },
  { label: '高', value: 'high' }, { label: '紧急', value: 'urgent' },
]
const adminOptions = [
  { label: '管理员A', value: 1 }, { label: '管理员B', value: 2 }, { label: '管理员C', value: 3 },
]

const statusMap: Record<string, { label: string; type: string }> = {
  open: { label: '待处理', type: 'warning' }, processing: { label: '处理中', type: 'primary' },
  replied: { label: '已回复', type: 'success' }, closed: { label: '已关闭', type: 'danger' },
}
const priorityMap: Record<string, { label: string; type: string }> = {
  low: { label: '低', type: 'success' }, medium: { label: '中', type: 'info' },
  high: { label: '高', type: 'warning' }, urgent: { label: '紧急', type: 'danger' },
}

const tickets = ref([
  { id: 1, ticketNo: 'TK-2024001', title: '无法登录控制面板', user: '张三', priority: 'high', status: 'open', assignee: null, createdAt: '2024-03-15 14:30', description: '登录时提示密码错误，但我确认密码是正确的。', messages: [
    { sender: '张三', content: '我尝试了多次登录都不行，显示密码错误。', time: '2024-03-15 14:30', isAdmin: false },
  ] },
  { id: 2, ticketNo: 'TK-2024002', title: '服务器访问速度慢', user: '李四', priority: 'medium', status: 'processing', assignee: '管理员A', createdAt: '2024-03-15 10:20', description: '最近一周服务器响应时间明显变长。', messages: [
    { sender: '李四', content: '从上周开始网站打开要10秒以上。', time: '2024-03-15 10:20', isAdmin: false },
    { sender: '管理员A', content: '已收到，正在排查服务器负载情况。', time: '2024-03-15 11:00', isAdmin: true },
  ] },
  { id: 3, ticketNo: 'TK-2024003', title: '域名解析问题', user: '王五', priority: 'low', status: 'replied', assignee: '管理员B', createdAt: '2024-03-14 16:45', description: '新注册的域名无法解析到我的服务器。', messages: [
    { sender: '王五', content: '域名已经注册24小时了，但nslookup还是查不到。', time: '2024-03-14 16:45', isAdmin: false },
    { sender: '管理员B', content: 'DNS传播需要24-48小时，请耐心等待。', time: '2024-03-14 17:30', isAdmin: true },
  ] },
  { id: 4, ticketNo: 'TK-2024004', title: '退款申请', user: '赵六', priority: 'urgent', status: 'open', assignee: null, createdAt: '2024-03-15 09:00', description: '购买的服务完全无法使用，要求全额退款。', messages: [
    { sender: '赵六', content: '服务购买后一直无法正常使用，请尽快退款。', time: '2024-03-15 09:00', isAdmin: false },
  ] },
])

function openDrawer(ticket: any) {
  currentTicket.value = ticket
  replyContent.value = ''
  assignAdminId.value = null
  drawerVisible.value = true
}

function handleReply() {
  if (!replyContent.value.trim() || !currentTicket.value) return
  currentTicket.value.messages.push({ sender: '管理员', content: replyContent.value.trim(), time: new Date().toLocaleString(), isAdmin: true })
  currentTicket.value.status = 'replied'
  replyContent.value = ''
  ElMessage.success('回复成功')
}

function handleAssign() {
  if (!assignAdminId.value || !currentTicket.value) return
  const admin = adminOptions.find((a) => a.value === assignAdminId.value)
  currentTicket.value.assignee = admin?.label || '管理员'
  currentTicket.value.status = 'processing'
  ElMessage.success(`已指派给 ${admin?.label}`)
}

function handleClose() {
  if (currentTicket.value) {
    currentTicket.value.status = 'closed'
    ElMessage.success('工单已关闭')
  }
}

function handleSearch() { pagination.page = 1 }
</script>

<style scoped>
.card-header { display: flex; align-items: center; justify-content: space-between; flex-wrap: wrap; gap: 12px; }
.card-title { font-size: 16px; font-weight: 600; }
.card-actions { display: flex; align-items: center; gap: 12px; }
.conversation-list { max-height: 300px; overflow-y: auto; padding: 8px 0; }
.message-item { display: flex; gap: 12px; padding: 8px 0; }
.message-item.admin { flex-direction: row-reverse; }
.message-content { max-width: 70%; }
.message-header { display: flex; gap: 8px; margin-bottom: 4px; font-size: 12px; }
.message-item.admin .message-header { flex-direction: row-reverse; }
.sender { font-weight: 600; color: #333; }
.time { color: #999; }
.message-text { background: #f5f5f5; padding: 8px 12px; border-radius: 8px; font-size: 14px; line-height: 1.5; }
.message-item.admin .message-text { background: #e8f5e9; }
</style>
