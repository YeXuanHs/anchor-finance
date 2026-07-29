<template>
  <div class="ticket-detail page-container">
    <div class="page-header">
      <el-button @click="$router.back()">
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>
      <h2>工单详情</h2>
    </div>
    
    <div class="detail-grid">
      <div class="main-content">
        <div class="art-card">
          <h3>{{ ticket.title }}</h3>
          <div class="ticket-meta">
            <span>工单号：{{ ticket.ticket_no }}</span>
            <span>创建时间：{{ ticket.created_at }}</span>
          </div>
          
          <div class="messages-list">
            <div v-for="msg in ticket.messages" :key="msg.id" class="message-item" :class="msg.type">
              <div class="message-header">
                <span class="sender">{{ msg.sender }}</span>
                <span class="time">{{ msg.time }}</span>
              </div>
              <div class="message-content">{{ msg.content }}</div>
            </div>
          </div>
          
          <div class="reply-area">
            <el-input
              v-model="replyContent"
              type="textarea"
              :rows="4"
              placeholder="输入回复内容..."
            />
            <div class="reply-actions">
              <el-upload
                action="#"
                :auto-upload="false"
                :show-file-list="false"
              >
                <el-button>
                  <el-icon><Upload /></el-icon>
                  附件
                </el-button>
              </el-upload>
              <el-button type="primary" @click="submitReply" :loading="submitting">
                提交回复
              </el-button>
            </div>
          </div>
        </div>
      </div>
      
      <div class="side-info">
        <div class="art-card">
          <h3>工单信息</h3>
          <div class="info-item">
            <span class="label">状态</span>
            <span class="value">
              <el-tag :type="getStatusType(ticket.status)">
                {{ getStatusText(ticket.status) }}
              </el-tag>
            </span>
          </div>
          <div class="info-item">
            <span class="label">优先级</span>
            <span class="value">
              <el-tag :type="getPriorityType(ticket.priority)">
                {{ ticket.priority }}
              </el-tag>
            </span>
          </div>
          <div class="info-item">
            <span class="label">部门</span>
            <span class="value">{{ ticket.department }}</span>
          </div>
          <div class="info-item">
            <span class="label">用户</span>
            <span class="value">{{ ticket.user }}</span>
          </div>
        </div>
        
        <div class="art-card">
          <h3>操作</h3>
          <div class="actions">
            <el-button @click="changeStatus('processing')" v-if="ticket.status === 'open'">
              开始处理
            </el-button>
            <el-button type="success" @click="changeStatus('closed')" v-if="ticket.status !== 'closed'">
              关闭工单
            </el-button>
            <el-button @click="changePriority">修改优先级</el-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const route = useRoute()
const replyContent = ref('')
const submitting = ref(false)

const ticket = ref({
  ticket_no: route.params.id,
  title: '服务器无法连接',
  status: 'open',
  priority: '高',
  department: '技术支持',
  user: '张三',
  created_at: '2026-07-27 18:00:00',
  messages: [
    {
      id: 1,
      type: 'user',
      sender: '张三',
      time: '2026-07-27 18:00:00',
      content: '我的服务器从今天下午开始无法连接，IP是 192.168.1.100，请帮忙检查一下。'
    }
  ]
})

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    open: 'warning',
    processing: 'primary',
    replied: 'success',
    closed: 'info'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    open: '待处理',
    processing: '处理中',
    replied: '已回复',
    closed: '已关闭'
  }
  return map[status] || status
}

const getPriorityType = (priority: string) => {
  const map: Record<string, string> = {
    '高': 'danger',
    '中': 'warning',
    '低': 'info'
  }
  return map[priority] || 'info'
}

const fetchTicket = async () => {
  try {
    const { data } = await request.get(`/admin/api/v1/tickets/${route.params.id}`)
    if (data.data) {
      Object.assign(ticket.value, data.data)
    }
  } catch {}
}

const submitReply = async () => {
  if (!replyContent.value.trim()) return
  
  submitting.value = true
  try {
    await request.post(`/admin/api/v1/tickets/${route.params.id}/reply`, { content: replyContent.value })
    ElMessage.success('回复成功')
    replyContent.value = ''
    fetchTicket()
  } catch {
    ElMessage.error('回复失败')
  } finally {
    submitting.value = false
  }
}

const changeStatus = async (status: string) => {
  try {
    await request.put(`/admin/api/v1/tickets/${route.params.id}/status`, { status })
    ticket.value.status = status
    ElMessage.success('状态已更新')
  } catch {
    ElMessage.error('操作失败')
  }
}

const changePriority = async () => {
  try {
    await request.put(`/admin/api/v1/tickets/${route.params.id}/priority`, { priority: ticket.value.priority })
    ElMessage.success('优先级已更新')
  } catch {
    ElMessage.error('操作失败')
  }
}

onMounted(() => {
  fetchTicket()
})
</script>

<style scoped lang="scss">
.ticket-detail {
  .page-header {
    display: flex;
    align-items: center;
    gap: 16px;
    margin-bottom: 20px;
    
    h2 {
      margin: 0;
      font-size: 20px;
    }
  }
  
  .detail-grid {
    display: grid;
    grid-template-columns: 2fr 1fr;
    gap: 20px;
    
    @media (max-width: 1200px) {
      grid-template-columns: 1fr;
    }
  }
  
  .main-content {
    h3 {
      font-size: 18px;
      font-weight: 600;
      margin: 0 0 12px;
    }
    
    .ticket-meta {
      display: flex;
      gap: 20px;
      color: var(--text-secondary);
      font-size: 14px;
      margin-bottom: 20px;
      padding-bottom: 16px;
      border-bottom: 1px solid var(--border-color);
    }
    
    .messages-list {
      margin-bottom: 20px;
    }
    
    .message-item {
      padding: 16px;
      border-radius: 12px;
      margin-bottom: 16px;
      
      &.user {
        background: #f5f7fa;
      }
      
      &.admin {
        background: var(--primary-bg);
        margin-left: 40px;
      }
      
      .message-header {
        display: flex;
        justify-content: space-between;
        margin-bottom: 8px;
        
        .sender {
          font-weight: 600;
        }
        
        .time {
          color: var(--text-secondary);
          font-size: 13px;
        }
      }
      
      .message-content {
        line-height: 1.8;
      }
    }
    
    .reply-area {
      .reply-actions {
        display: flex;
        justify-content: space-between;
        margin-top: 12px;
      }
    }
  }
  
  .side-info {
    .art-card {
      margin-bottom: 20px;
      
      h3 {
        font-size: 16px;
        font-weight: 600;
        margin: 0 0 16px;
        padding-bottom: 12px;
        border-bottom: 1px solid var(--border-color);
      }
      
      .info-item {
        display: flex;
        padding: 12px 0;
        border-bottom: 1px solid #f5f5f5;
        
        &:last-child {
          border-bottom: none;
        }
        
        .label {
          width: 80px;
          color: var(--text-secondary);
        }
        
        .value {
          flex: 1;
        }
      }
      
      .actions {
        display: flex;
        flex-direction: column;
        gap: 12px;
        
        .el-button {
          width: 100%;
        }
      }
    }
  }
}
</style>
