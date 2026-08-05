<template>
  <div class="ticket-detail" v-loading="loading">
    <div class="page-header">
      <el-page-header @back="$router.push('/user/tickets')">
        <template #content>
          <span class="page-title">工单详情</span>
        </template>
      </el-page-header>
    </div>
    
    <div class="content-wrapper" v-if="ticket">
      <!-- 工单信息 -->
      <el-card class="info-card">
        <template #header>
          <div class="card-header">
            <span>{{ ticket.subject }}</span>
            <el-tag :type="getStatusType(ticket.status)">{{ getStatusText(ticket.status) }}</el-tag>
          </div>
        </template>
        <el-descriptions :column="3" border>
          <el-descriptions-item label="工单号">{{ ticket.ticket_no }}</el-descriptions-item>
          <el-descriptions-item label="类型">{{ ticket.type }}</el-descriptions-item>
          <el-descriptions-item label="优先级">
            <el-tag :type="getPriorityType(ticket.priority)" size="small">
              {{ getPriorityText(ticket.priority) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="提交时间">{{ ticket.created_at }}</el-descriptions-item>
          <el-descriptions-item label="最后更新">{{ ticket.updated_at }}</el-descriptions-item>
          <el-descriptions-item label="关联产品">{{ ticket.product_name || '-' }}</el-descriptions-item>
        </el-descriptions>
      </el-card>
      
      <!-- 对话记录 -->
      <el-card class="messages-card">
        <template #header>
          <span>对话记录</span>
        </template>
        <div class="messages-list">
          <div 
            v-for="msg in ticket.messages" 
            :key="msg.id"
            class="message-item"
            :class="{ 'is-admin': msg.is_admin }"
          >
            <div class="message-header">
              <el-avatar :size="36" :class="{ 'admin-avatar': msg.is_admin }">
                {{ msg.is_admin ? '客' : '我' }}
              </el-avatar>
              <div class="message-meta">
                <span class="sender">{{ msg.is_admin ? '客服' : username }}</span>
                <span class="time">{{ msg.created_at }}</span>
              </div>
            </div>
            <div class="message-content">{{ msg.content }}</div>
            <div class="message-attachments" v-if="msg.attachments?.length">
              <div v-for="file in msg.attachments" :key="file.url" class="attachment">
                <el-icon><Document /></el-icon>
                <a :href="file.url" target="_blank">{{ file.name }}</a>
              </div>
            </div>
          </div>
        </div>
      </el-card>
      
      <!-- 回复区域 -->
      <el-card class="reply-card" v-if="ticket.status !== 'closed'">
        <template #header>
          <span>回复</span>
        </template>
        <el-input v-model="replyContent" type="textarea" :rows="4" placeholder="请输入回复内容..." />
        <div class="reply-actions">
          <el-upload
            :action="`/api/v1/tickets/${route.params.id}/attachments`"
            :on-success="handleUploadSuccess"
            :show-file-list="false"
          >
            <el-button>上传附件</el-button>
          </el-upload>
          <el-button type="primary" @click="submitReply" :loading="submitting">提交回复</el-button>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Document } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import request from '@/utils/request'

const route = useRoute()
const userStore = useUserStore()
const loading = ref(false)
const submitting = ref(false)
const ticket = ref<any>(null)
const replyContent = ref('')

const username = computed(() => userStore.username || '用户')

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    open: 'primary',
    replied: 'success',
    closed: 'info'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    open: '待处理',
    replied: '已回复',
    closed: '已关闭'
  }
  return map[status] || status
}

const getPriorityType = (priority: string) => {
  const map: Record<string, string> = {
    low: 'info',
    medium: 'warning',
    high: 'danger'
  }
  return map[priority] || 'info'
}

const getPriorityText = (priority: string) => {
  const map: Record<string, string> = {
    low: '低',
    medium: '中',
    high: '高'
  }
  return map[priority] || priority
}

const handleUploadSuccess = (response: any) => {
  // 处理上传成功
}

const fetchTicket = async () => {
  loading.value = true
  try {
    const { data } = await request.get(`/api/v1/tickets/${route.params.id}`)
    if (data?.data) {
      ticket.value = data.data
    }
  } catch (error) {
    ElMessage.error('获取工单信息失败')
  } finally {
    loading.value = false
  }
}

const submitReply = async () => {
  if (!replyContent.value.trim()) {
    ElMessage.warning('请输入回复内容')
    return
  }
  
  submitting.value = true
  try {
    await request.post(`/api/v1/tickets/${route.params.id}/reply`, {
      content: replyContent.value
    })
    ElMessage.success('回复成功')
    replyContent.value = ''
    fetchTicket()
  } catch (error) {
    ElMessage.error('回复失败')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchTicket()
})
</script>

<style scoped lang="scss">
.ticket-detail {
  .page-header {
    margin-bottom: 24px;
  }
  
  .page-title {
    font-size: 18px;
    font-weight: 600;
  }
  
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  
  .info-card {
    margin-bottom: 20px;
  }
  
  .messages-card {
    margin-bottom: 20px;
    
    .messages-list {
      .message-item {
        padding: 16px;
        margin-bottom: 16px;
        background: #f5f7fa;
        border-radius: 8px;
        
        &.is-admin {
          background: #ecf5ff;
        }
        
        .message-header {
          display: flex;
          align-items: center;
          gap: 12px;
          margin-bottom: 12px;
          
          .admin-avatar {
            background: #409eff;
          }
          
          .message-meta {
            .sender {
              font-weight: 600;
              margin-right: 12px;
            }
            
            .time {
              color: #909399;
              font-size: 12px;
            }
          }
        }
        
        .message-content {
          line-height: 1.6;
          white-space: pre-wrap;
        }
        
        .message-attachments {
          margin-top: 12px;
          padding-top: 12px;
          border-top: 1px solid #eee;
          
          .attachment {
            display: flex;
            align-items: center;
            gap: 8px;
            margin-bottom: 8px;
            
            a {
              color: #409eff;
              text-decoration: none;
            }
          }
        }
      }
    }
  }
  
  .reply-card {
    .reply-actions {
      display: flex;
      justify-content: space-between;
      margin-top: 16px;
    }
  }
}
</style>
