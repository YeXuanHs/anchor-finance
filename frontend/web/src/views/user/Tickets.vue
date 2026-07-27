<template>
  <div class="tickets-page">
    <div class="page-header">
      <h1 class="page-title">工单支持</h1>
      <n-button type="primary" @click="showCreateModal = true">
        <template #icon>
          <n-icon :component="CreateOutline" />
        </template>
        提交工单
      </n-button>
    </div>

    <n-tabs v-model:value="activeFilter" type="line" animated class="filter-tabs">
      <n-tab name="all">
        全部
        <n-badge :value="tickets.length" :max="99" class="tab-badge" />
      </n-tab>
      <n-tab name="open">
        处理中
        <n-badge :value="tickets.filter(t => t.status === 'open').length" :max="99" class="tab-badge" />
      </n-tab>
      <n-tab name="pending">
        等待回复
        <n-badge :value="tickets.filter(t => t.status === 'pending').length" :max="99" class="tab-badge" />
      </n-tab>
      <n-tab name="closed">
        已关闭
        <n-badge :value="tickets.filter(t => t.status === 'closed').length" :max="99" class="tab-badge" />
      </n-tab>
    </n-tabs>

    <div class="ticket-list">
      <div
        v-for="ticket in filteredTickets"
        :key="ticket.id"
        class="ticket-card"
        :class="{ 'ticket-expanded': expandedTicket === ticket.id }"
      >
        <div class="ticket-header" @click="toggleExpand(ticket.id)">
          <div class="ticket-info">
            <div class="ticket-title-row">
              <n-tag :type="statusMap[ticket.status].type" size="small" round>
                {{ statusMap[ticket.status].label }}
              </n-tag>
              <n-tag :type="priorityMap[ticket.priority].type" size="small" round :bordered="false">
                {{ priorityMap[ticket.priority].label }}
              </n-tag>
              <h4 class="ticket-subject">{{ ticket.subject }}</h4>
            </div>
            <div class="ticket-meta">
              <span class="ticket-id">#{{ ticket.id }}</span>
              <span class="ticket-dept">
                <n-icon :component="FolderOutline" size="14" />
                {{ ticket.department }}
              </span>
              <span class="ticket-time">
                <n-icon :component="TimeOutline" size="14" />
                最后回复：{{ ticket.lastReply }}
              </span>
            </div>
          </div>

          <n-icon
            :component="expandedTicket === ticket.id ? ChevronUpOutline : ChevronDownOutline"
            size="20"
            class="expand-icon"
          />
        </div>

        <div v-if="expandedTicket === ticket.id" class="ticket-conversation">
          <n-divider style="margin: 16px 0;" />

          <div class="messages">
            <div
              v-for="msg in ticket.messages"
              :key="msg.id"
              :class="['message', { 'message-admin': msg.isAdmin }]"
            >
              <div class="message-avatar">
                <n-avatar :size="40" round :style="{ background: msg.isAdmin ? '#1890ff' : '#52c41a' }">
                  {{ msg.sender.charAt(0) }}
                </n-avatar>
              </div>
              <div class="message-content">
                <div class="message-header">
                  <span class="message-sender">{{ msg.sender }}</span>
                  <span class="message-time">{{ msg.time }}</span>
                </div>
                <div class="message-body">
                  {{ msg.content }}
                </div>
                <div v-if="msg.attachments && msg.attachments.length > 0" class="message-attachments">
                  <n-tag
                    v-for="file in msg.attachments"
                    :key="file"
                    size="small"
                    @click="handleDownload(file)"
                    style="cursor: pointer;"
                  >
                    <template #icon>
                      <n-icon :component="AttachOutline" />
                    </template>
                    {{ file }}
                  </n-tag>
                </div>
              </div>
            </div>
          </div>

          <div v-if="ticket.status !== 'closed'" class="reply-box">
            <n-input
              v-model:value="replyContent"
              type="textarea"
              placeholder="输入回复内容..."
              :rows="3"
              show-count
              :maxlength="1000"
            />
            <div class="reply-actions">
              <n-upload
                :show-file-list="false"
                @change="handleFileUpload"
              >
                <n-button quaternary>
                  <template #icon>
                    <n-icon :component="AttachOutline" />
                  </template>
                  附件
                </n-button>
              </n-upload>
              <n-button type="primary" @click="handleReply(ticket)">
                <template #icon>
                  <n-icon :component="SendOutline" />
                </template>
                发送
              </n-button>
            </div>
          </div>

          <div v-else class="closed-notice">
            <n-icon :component="InformationCircleOutline" size="16" />
            此工单已关闭，如需继续沟通请重新提交工单
          </div>
        </div>
      </div>

      <div v-if="filteredTickets.length === 0" class="empty-state">
        <n-icon :size="64" :component="ChatbubblesOutline" color="#d9d9d9" />
        <p>暂无工单</p>
        <n-button type="primary" @click="showCreateModal = true">提交工单</n-button>
      </div>
    </div>

    <!-- Create Ticket Modal -->
    <n-modal
      v-model:show="showCreateModal"
      preset="card"
      title="提交工单"
      style="width: 600px; max-width: 90vw;"
      :bordered="false"
      :segmented="{ content: true, footer: true }"
    >
      <n-form ref="createFormRef" :model="createForm" :rules="createRules">
        <n-form-item label="工单部门" path="department">
          <n-select
            v-model:value="createForm.department"
            placeholder="请选择部门"
            :options="departmentOptions"
          />
        </n-form-item>

        <n-form-item label="工单标题" path="subject">
          <n-input
            v-model:value="createForm.subject"
            placeholder="请简要描述您的问题"
          />
        </n-form-item>

        <n-form-item label="优先级" path="priority">
          <n-radio-group v-model:value="createForm.priority">
            <n-radio-button value="low">低</n-radio-button>
            <n-radio-button value="medium">中</n-radio-button>
            <n-radio-button value="high">高</n-radio-button>
            <n-radio-button value="urgent">紧急</n-radio-button>
          </n-radio-group>
        </n-form-item>

        <n-form-item label="问题描述" path="content">
          <n-input
            v-model:value="createForm.content"
            type="textarea"
            placeholder="请详细描述您遇到的问题..."
            :rows="5"
            show-count
            :maxlength="2000"
          />
        </n-form-item>

        <n-form-item label="附件">
          <n-upload
            :file-list="createForm.attachments"
            @update:file-list="createForm.attachments = $event"
            :max="5"
          >
            <n-button>上传附件</n-button>
          </n-upload>
        </n-form-item>
      </n-form>

      <template #footer>
        <n-space justify="end">
          <n-button @click="showCreateModal = false">取消</n-button>
          <n-button type="primary" :loading="creating" @click="handleCreate">
            提交工单
          </n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive } from 'vue'
import { useMessage } from 'naive-ui'
import type { FormInst, FormRules } from 'naive-ui'
import {
  CreateOutline,
  FolderOutline,
  TimeOutline,
  ChevronUpOutline,
  ChevronDownOutline,
  AttachOutline,
  InformationCircleOutline,
  SendOutline,
  ChatbubblesOutline
} from '@vicons/ionicons5'

const message = useMessage()

const activeFilter = ref('all')
const expandedTicket = ref<string | null>(null)
const showCreateModal = ref(false)
const creating = ref(false)
const replyContent = ref('')
const createFormRef = ref<FormInst | null>(null)

interface Message {
  id: string
  sender: string
  isAdmin: boolean
  content: string
  time: string
  attachments?: string[]
}

interface Ticket {
  id: string
  subject: string
  department: string
  priority: 'low' | 'medium' | 'high' | 'urgent'
  status: 'open' | 'pending' | 'closed'
  lastReply: string
  messages: Message[]
}

const statusMap: Record<string, { label: string; type: 'success' | 'warning' | 'error' | 'info' | 'default' }> = {
  open: { label: '处理中', type: 'info' },
  pending: { label: '等待回复', type: 'warning' },
  closed: { label: '已关闭', type: 'default' }
}

const priorityMap: Record<string, { label: string; type: 'success' | 'warning' | 'error' | 'info' }> = {
  low: { label: '低', type: 'info' },
  medium: { label: '中', type: 'warning' },
  high: { label: '高', type: 'error' },
  urgent: { label: '紧急', type: 'error' }
}

const departmentOptions = [
  { label: '技术支持', value: '技术支持' },
  { label: '账单问题', value: '账单问题' },
  { label: '产品咨询', value: '产品咨询' },
  { label: '投诉建议', value: '投诉建议' },
  { label: '其他', value: '其他' }
]

const createForm = reactive({
  department: null as string | null,
  priority: 'medium',
  subject: '',
  content: '',
  attachments: [] as any[]
})

const createRules: FormRules = {
  department: { required: true, message: '请选择部门', trigger: 'change' },
  subject: { required: true, message: '请输入标题', trigger: 'blur' },
  content: { required: true, message: '请输入问题描述', trigger: 'blur' }
}

const tickets = ref<Ticket[]>([
  {
    id: 'TK20260725001',
    subject: '云服务器无法远程连接',
    department: '技术支持',
    priority: 'high',
    status: 'open',
    lastReply: '2026-07-25 14:30',
    messages: [
      {
        id: '1',
        sender: '我',
        isAdmin: false,
        content: '我的云服务器今天突然无法远程连接了，ping也不通，请帮忙检查一下。',
        time: '2026-07-25 10:00'
      },
      {
        id: '2',
        sender: '技术支持',
        isAdmin: true,
        content: '您好，我们已经收到您的工单。经检查，您的服务器因安全策略被临时限制。我们已解除限制，请稍后重试。如仍有问题请回复。',
        time: '2026-07-25 14:30'
      }
    ]
  },
  {
    id: 'TK20260720002',
    subject: '账单金额异常',
    department: '账单问题',
    priority: 'medium',
    status: 'pending',
    lastReply: '2026-07-20 16:00',
    messages: [
      {
        id: '1',
        sender: '我',
        isAdmin: false,
        content: '这个月的账单金额比上个月多了很多，但是我没有升级过配置，请帮忙核实。',
        time: '2026-07-20 09:00'
      },
      {
        id: '2',
        sender: '账单专员',
        isAdmin: true,
        content: '您好，经核实您上月有额外的流量费用产生。详细账单已发送至您的邮箱，请查收。如有疑问请回复。',
        time: '2026-07-20 16:00'
      }
    ]
  },
  {
    id: 'TK20260715003',
    subject: 'SSL证书安装问题',
    department: '技术支持',
    priority: 'low',
    status: 'closed',
    lastReply: '2026-07-16 10:00',
    messages: [
      {
        id: '1',
        sender: '我',
        isAdmin: false,
        content: '我购买了SSL证书，但是不知道怎么安装到Nginx上，能提供教程吗？',
        time: '2026-07-15 14:00',
        attachments: ['error-screenshot.png']
      },
      {
        id: '2',
        sender: '技术支持',
        isAdmin: true,
        content: '您好，以下是Nginx安装SSL证书的步骤：\n1. 将证书文件上传到服务器\n2. 编辑Nginx配置文件\n3. 添加SSL相关配置\n4. 重启Nginx服务\n\n详细教程请参考帮助文档。',
        time: '2026-07-16 10:00'
      }
    ]
  }
])

const filteredTickets = computed(() => {
  if (activeFilter.value === 'all') return tickets.value
  return tickets.value.filter((t) => t.status === activeFilter.value)
})

function toggleExpand(id: string) {
  expandedTicket.value = expandedTicket.value === id ? null : id
}

function handleReply(ticket: Ticket) {
  if (!replyContent.value.trim()) {
    message.warning('请输入回复内容')
    return
  }

  ticket.messages.push({
    id: String(Date.now()),
    sender: '我',
    isAdmin: false,
    content: replyContent.value,
    time: new Date().toLocaleString('zh-CN')
  })

  ticket.lastReply = new Date().toLocaleString('zh-CN')
  replyContent.value = ''
  message.success('回复已发送')
}

function handleDownload(filename: string) {
  message.info(`下载附件：${filename}`)
}

function handleFileUpload() {
  message.info('附件上传功能')
}

async function handleCreate() {
  try {
    await createFormRef.value?.validate()
    creating.value = true

    const newTicket: Ticket = {
      id: `TK${Date.now()}`,
      subject: createForm.subject,
      department: createForm.department!,
      priority: createForm.priority as any,
      status: 'open',
      lastReply: new Date().toLocaleString('zh-CN'),
      messages: [
        {
          id: '1',
          sender: '我',
          isAdmin: false,
          content: createForm.content,
          time: new Date().toLocaleString('zh-CN')
        }
      ]
    }

    tickets.value.unshift(newTicket)
    showCreateModal.value = false
    message.success('工单已提交')

    createForm.department = null
    createForm.priority = 'medium'
    createForm.subject = ''
    createForm.content = ''
    createForm.attachments = []
  } catch {
    // Validation failed
  } finally {
    creating.value = false
  }
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
  font-size: 22px;
  font-weight: 700;
  color: #262626;
  margin: 0;
}

.filter-tabs {
  background: #fff;
  border-radius: 12px;
  padding: 4px 16px;
  border: 1px solid #f0f0f0;
}

.tab-badge {
  margin-left: 6px;
}

.ticket-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.ticket-card {
  background: #fff;
  border: 1px solid #f0f0f0;
  border-radius: 12px;
  padding: 20px;
  transition: all 0.3s ease;
}

.ticket-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.ticket-expanded {
  border-color: #1890ff;
  box-shadow: 0 4px 16px rgba(24, 144, 255, 0.12);
}

.ticket-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  cursor: pointer;
}

.ticket-info {
  flex: 1;
}

.ticket-title-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.ticket-subject {
  font-size: 16px;
  font-weight: 600;
  color: #262626;
  margin: 0;
}

.ticket-meta {
  display: flex;
  gap: 16px;
  font-size: 13px;
  color: #8c8c8c;
}

.ticket-meta span {
  display: flex;
  align-items: center;
  gap: 4px;
}

.ticket-id {
  font-family: monospace;
}

.expand-icon {
  color: #999;
  transition: transform 0.2s;
  flex-shrink: 0;
  margin-top: 4px;
}

.ticket-conversation {
  margin-top: 8px;
}

.messages {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.message {
  display: flex;
  gap: 12px;
}

.message-admin {
  flex-direction: row-reverse;
}

.message-avatar {
  flex-shrink: 0;
}

.message-content {
  flex: 1;
  max-width: 80%;
}

.message-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.message-admin .message-header {
  flex-direction: row-reverse;
}

.message-sender {
  font-size: 14px;
  font-weight: 600;
  color: #262626;
}

.message-time {
  font-size: 12px;
  color: #8c8c8c;
}

.message-body {
  background: #f5f7fa;
  padding: 12px 16px;
  border-radius: 12px;
  font-size: 14px;
  line-height: 1.6;
  color: #333;
  white-space: pre-wrap;
}

.message-admin .message-body {
  background: #e6f7ff;
}

.message-attachments {
  display: flex;
  gap: 8px;
  margin-top: 8px;
}

.reply-box {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
}

.reply-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 12px;
}

.closed-notice {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 16px;
  padding: 12px 16px;
  background: #f5f5f5;
  border-radius: 8px;
  font-size: 14px;
  color: #8c8c8c;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  padding: 80px 0;
  background: #fff;
  border-radius: 12px;
  border: 1px solid #f0f0f0;
}

.empty-state p {
  margin: 0;
  color: #8c8c8c;
  font-size: 14px;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 12px;
    align-items: flex-start;
  }

  .ticket-header {
    flex-direction: column;
    gap: 12px;
  }

  .ticket-meta {
    flex-direction: column;
    gap: 4px;
  }

  .ticket-title-row {
    flex-wrap: wrap;
  }

  .message-content {
    max-width: 90%;
  }
}
</style>
