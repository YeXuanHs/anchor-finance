<template>
  <div class="system-message-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span>系统消息</span>
          <el-button @click="markAllRead">全部已读</el-button>
        </div>
      </template>

      <el-table :data="messages" style="width: 100%" v-loading="loading">
        <el-table-column prop="created_at" label="时间" width="180" />
        <el-table-column prop="title" label="标题">
          <template #default="{ row }">
            <div class="message-title" :class="{ unread: !row.read }">
              <el-badge is-dot :hidden="row.read" />
              <span @click="viewMessage(row)">{{ row.title }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getTypeTag(row.type)">{{ getTypeText(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <el-button size="small" @click="viewMessage(row)">查看</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 消息详情对话框 -->
    <el-dialog v-model="showDialog" title="消息详情" width="600px">
      <div class="message-detail">
        <h3>{{ currentMessage?.title }}</h3>
        <div class="message-time">{{ currentMessage?.created_at }}</div>
        <div class="message-content" v-html="currentMessage?.content" />
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import request from '@/utils/request'

const loading = ref(false)
const messages = ref<any[]>([])
const showDialog = ref(false)
const currentMessage = ref<any>(null)

onMounted(async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/v1/messages')
    messages.value = data.data?.list || data.list || data.data || []
  } catch (e) { console.error(e) } finally { loading.value = false }
})

const getTypeTag = (type: string) => {
  const map: Record<string, string> = {
    system: 'primary',
    order: 'success',
    finance: 'warning',
    security: 'danger'
  }
  return map[type] || 'info'
}

const getTypeText = (type: string) => {
  const map: Record<string, string> = {
    system: '系统',
    order: '订单',
    finance: '财务',
    security: '安全'
  }
  return map[type] || type
}

const viewMessage = (message: any) => {
  currentMessage.value = message
  message.read = true
  showDialog.value = true
}

const markAllRead = () => {
  messages.value.forEach(m => m.read = true)
}
</script>

<style scoped lang="scss">
.system-message-page {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .message-title {
    display: flex;
    align-items: center;
    gap: 8px;
    cursor: pointer;

    &.unread span {
      font-weight: bold;
    }
  }

  .message-detail {
    h3 {
      margin: 0 0 10px;
    }

    .message-time {
      color: #909399;
      font-size: 14px;
      margin-bottom: 20px;
    }

    .message-content {
      line-height: 1.8;
    }
  }
}
</style>
