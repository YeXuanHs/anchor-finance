<template>
  <n-card title="工单管理">
    <template #header-extra>
      <n-space>
        <n-select v-model:value="filters.status" :options="statusOptions" placeholder="工单状态" clearable style="width: 140px" />
        <n-select v-model:value="filters.priority" :options="priorityOptions" placeholder="优先级" clearable style="width: 120px" />
        <n-input v-model:value="filters.keyword" placeholder="搜索工单标题/用户" clearable style="width: 200px" @keydown.enter="handleSearch" />
        <n-button type="primary" @click="handleSearch">搜索</n-button>
      </n-space>
    </template>

    <n-data-table :columns="columns" :data="tickets" :loading="loading" :bordered="false" :pagination="pagination" @update:page="handlePageChange" @update:page-size="handlePageSizeChange" />

    <!-- Ticket Detail Drawer -->
    <n-drawer v-model:show="drawerVisible" :width="600">
      <n-drawer-content :title="`工单详情 - ${currentTicket?.title || ''}`">
        <n-descriptions bordered :column="2" label-placement="left">
          <n-descriptions-item label="工单号">{{ currentTicket?.ticketNo }}</n-descriptions-item>
          <n-descriptions-item label="提交用户">{{ currentTicket?.user }}</n-descriptions-item>
          <n-descriptions-item label="优先级">
            <n-tag :type="priorityMap[currentTicket?.priority || '']?.type" size="small">
              {{ priorityMap[currentTicket?.priority || '']?.label }}
            </n-tag>
          </n-descriptions-item>
          <n-descriptions-item label="状态">
            <n-tag :type="statusMap[currentTicket?.status || '']?.type" size="small" round>
              {{ statusMap[currentTicket?.status || '']?.label }}
            </n-tag>
          </n-descriptions-item>
          <n-descriptions-item label="创建时间" :span="2">{{ currentTicket?.createdAt }}</n-descriptions-item>
          <n-descriptions-item label="问题描述" :span="2">{{ currentTicket?.description }}</n-descriptions-item>
        </n-descriptions>

        <n-divider>对话记录</n-divider>
        <div class="conversation-list">
          <div v-for="(msg, idx) in currentTicket?.messages || []" :key="idx" class="message-item" :class="{ admin: msg.isAdmin }">
            <n-avatar :size="32" round :style="{ backgroundColor: msg.isAdmin ? '#18a058' : '#2080f0' }">
              {{ msg.sender.charAt(0) }}
            </n-avatar>
            <div class="message-content">
              <div class="message-header">
                <span class="sender">{{ msg.sender }}</span>
                <span class="time">{{ msg.time }}</span>
              </div>
              <div class="message-text">{{ msg.content }}</div>
            </div>
          </div>
        </div>

        <n-divider>回复</n-divider>
        <n-input v-model:value="replyContent" type="textarea" :rows="3" placeholder="输入回复内容..." />

        <template #footer>
          <n-space>
            <n-select v-model:value="assignAdminId" :options="adminOptions" placeholder="指派给" clearable style="width: 160px" />
            <n-button type="info" @click="handleAssign">指派</n-button>
            <n-button type="success" :disabled="!replyContent.trim()" @click="handleReply">回复</n-button>
            <n-button v-if="currentTicket?.status !== 'closed'" type="warning" @click="handleClose">关闭工单</n-button>
          </n-space>
        </template>
      </n-drawer-content>
    </n-drawer>
  </n-card>
</template>

<script setup lang="ts">
import { h, ref, reactive } from 'vue'
import { useMessage, NTag, NButton, NSpace, type DataTableColumns, type PaginationProps } from 'naive-ui'

const message = useMessage()
const loading = ref(false)
const drawerVisible = ref(false)
const currentTicket = ref<any>(null)
const replyContent = ref('')
const assignAdminId = ref<number | null>(null)

const filters = reactive({
  status: null as string | null,
  priority: null as string | null,
  keyword: '',
})

const statusOptions = [
  { label: '待处理', value: 'open' },
  { label: '处理中', value: 'processing' },
  { label: '已回复', value: 'replied' },
  { label: '已关闭', value: 'closed' },
]

const priorityOptions = [
  { label: '低', value: 'low' },
  { label: '中', value: 'medium' },
  { label: '高', value: 'high' },
  { label: '紧急', value: 'urgent' },
]

const adminOptions = [
  { label: '管理员A', value: 1 },
  { label: '管理员B', value: 2 },
  { label: '管理员C', value: 3 },
]

const statusMap: Record<string, { label: string; type: 'success' | 'info' | 'warning' | 'error' }> = {
  open: { label: '待处理', type: 'warning' },
  processing: { label: '处理中', type: 'info' },
  replied: { label: '已回复', type: 'success' },
  closed: { label: '已关闭', type: 'error' },
}

const priorityMap: Record<string, { label: string; type: 'success' | 'info' | 'warning' | 'error' }> = {
  low: { label: '低', type: 'success' },
  medium: { label: '中', type: 'info' },
  high: { label: '高', type: 'warning' },
  urgent: { label: '紧急', type: 'error' },
}

const pagination = reactive<PaginationProps>({
  page: 1,
  pageSize: 10,
  itemCount: 20,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
})

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
    { sender: '管理员B', content: 'DNS传播需要24-48小时，请耐心等待。如超过48小时仍未生效请联系。', time: '2024-03-14 17:30', isAdmin: true },
  ] },
  { id: 4, ticketNo: 'TK-2024004', title: '退款申请', user: '赵六', priority: 'urgent', status: 'open', assignee: null, createdAt: '2024-03-15 09:00', description: '购买的服务完全无法使用，要求全额退款。', messages: [
    { sender: '赵六', content: '服务购买后一直无法正常使用，请尽快退款。', time: '2024-03-15 09:00', isAdmin: false },
  ] },
])

const columns: DataTableColumns<any> = [
  { title: '工单号', key: 'ticketNo', width: 130 },
  { title: '标题', key: 'title', ellipsis: { tooltip: true } },
  { title: '用户', key: 'user', width: 80 },
  {
    title: '优先级',
    key: 'priority',
    width: 80,
    render: (row) => {
      const p = priorityMap[row.priority]
      return h(NTag, { type: p.type, size: 'small' }, { default: () => p.label })
    },
  },
  {
    title: '状态',
    key: 'status',
    width: 90,
    render: (row) => {
      const s = statusMap[row.status]
      return h(NTag, { type: s.type, size: 'small', round: true }, { default: () => s.label })
    },
  },
  { title: '指派给', key: 'assignee', width: 100, render: (row) => row.assignee || '未指派' },
  { title: '创建时间', key: 'createdAt', width: 160 },
  {
    title: '操作',
    key: 'actions',
    width: 100,
    render: (row) =>
      h(NButton, { size: 'small', type: 'primary', onClick: () => openDrawer(row) }, { default: () => '处理' }),
  },
]

function openDrawer(ticket: any) {
  currentTicket.value = ticket
  replyContent.value = ''
  assignAdminId.value = null
  drawerVisible.value = true
}

function handleReply() {
  if (!replyContent.value.trim() || !currentTicket.value) return
  currentTicket.value.messages.push({
    sender: '管理员',
    content: replyContent.value.trim(),
    time: new Date().toLocaleString(),
    isAdmin: true,
  })
  currentTicket.value.status = 'replied'
  replyContent.value = ''
  message.success('回复成功')
}

function handleAssign() {
  if (!assignAdminId.value || !currentTicket.value) return
  const admin = adminOptions.find((a) => a.value === assignAdminId.value)
  currentTicket.value.assignee = admin?.label || '管理员'
  currentTicket.value.status = 'processing'
  message.success(`已指派给 ${admin?.label}`)
}

function handleClose() {
  if (currentTicket.value) {
    currentTicket.value.status = 'closed'
    message.success('工单已关闭')
  }
}

function handleSearch() {
  pagination.page = 1
}

function handlePageChange(page: number) {
  pagination.page = page
}

function handlePageSizeChange(pageSize: number) {
  pagination.pageSize = pageSize
  pagination.page = 1
}
</script>

<style scoped>
.conversation-list {
  max-height: 300px;
  overflow-y: auto;
  padding: 8px 0;
}

.message-item {
  display: flex;
  gap: 12px;
  padding: 8px 0;
}

.message-item.admin {
  flex-direction: row-reverse;
}

.message-content {
  max-width: 70%;
}

.message-header {
  display: flex;
  gap: 8px;
  margin-bottom: 4px;
  font-size: 12px;
}

.message-item.admin .message-header {
  flex-direction: row-reverse;
}

.sender {
  font-weight: 600;
  color: #333;
}

.time {
  color: #999;
}

.message-text {
  background: #f5f5f5;
  padding: 8px 12px;
  border-radius: 8px;
  font-size: 14px;
  line-height: 1.5;
}

.message-item.admin .message-text {
  background: #e8f5e9;
}
</style>
