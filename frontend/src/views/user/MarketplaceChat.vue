<template>
  <div class="marketplace-chat">
    <div class="chat-sidebar">
      <div class="sidebar-header">
        <el-button @click="$router.back()" text size="small">
          <el-icon><ArrowLeft /></el-icon>
          返回
        </el-button>
        <h3>消息</h3>
      </div>

      <div class="session-list" v-loading="loadingSessions">
        <div v-if="sessions.length === 0" class="empty-tip">
          <el-empty description="暂无消息" :image-size="60" />
        </div>
        <div
          v-for="session in sessions"
          :key="session.id"
          class="session-item"
          :class="{ active: isActiveSession(session) }"
          @click="selectSession(session)"
        >
          <el-avatar :size="36" class="session-avatar">
            {{ getOtherUser(session)?.username?.charAt(0) || '?' }}
          </el-avatar>
          <div class="session-info">
            <div class="session-top">
              <span class="session-name">{{ getOtherUser(session)?.username || '未知用户' }}</span>
              <span class="session-time">{{ formatTime(session.last_message_at) }}</span>
            </div>
            <div class="session-bottom">
              <span class="session-product">{{ session.listing?.product_name || '主机' }}</span>
              <el-badge
                v-if="session.unread_count > 0"
                :value="session.unread_count"
                class="unread-badge"
              />
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="chat-main">
      <div v-if="!currentSession" class="no-chat">
        <el-empty description="选择一个会话开始聊天" />
      </div>

      <template v-else>
        <div class="chat-header">
          <div class="chat-user">
            <el-avatar :size="32">{{ otherUser?.username?.charAt(0) || '?' }}</el-avatar>
            <div class="user-info">
              <span class="name">{{ otherUser?.username }}</span>
              <span class="product" v-if="currentSession.listing">
                {{ currentSession.listing.product_name }}
              </span>
            </div>
          </div>
          <div class="chat-actions">
            <el-button size="small" @click="viewListing">查看商品</el-button>
          </div>
        </div>

        <div class="chat-messages" ref="messagesContainer" v-loading="loadingMessages">
          <div class="messages-wrapper">
            <div
              v-for="msg in messages"
              :key="msg.id"
              class="message-item"
              :class="{ 'is-mine': msg.sender_id === currentUserId }"
            >
              <el-avatar :size="28" class="msg-avatar">
                {{ msg.sender_id === currentUserId ? '我' : (otherUser?.username?.charAt(0) || '?') }}
              </el-avatar>
              <div class="msg-content">
                <div class="msg-bubble">{{ msg.content }}</div>
                <div class="msg-time">{{ formatTime(msg.created_at) }}</div>
              </div>
            </div>
          </div>
        </div>

        <div class="chat-input">
          <el-input
            v-model="inputMessage"
            type="textarea"
            :rows="2"
            placeholder="输入消息..."
            resize="none"
            @keydown.enter.exact.prevent="sendMessage"
          />
          <div class="input-actions">
            <el-button type="primary" :loading="sending" @click="sendMessage">
              发送
            </el-button>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import request from '@/utils/request'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const currentUserId = computed(() => userStore.userInfo?.id)

const loadingSessions = ref(false)
const loadingMessages = ref(false)
const sessions = ref<any[]>([])
const messages = ref<any[]>([])
const currentSession = ref<any>(null)
const inputMessage = ref('')
const sending = ref(false)
const messagesContainer = ref<HTMLElement>()

const otherUser = computed(() => {
  if (!currentSession.value) return null
  const session = currentSession.value
  const otherId = session.user1_id === currentUserId.value ? session.user2_id : session.user1_id
  // 从消息中获取用户信息
  const msg = messages.value.find(m => m.sender_id === otherId)
  return msg?.sender || { username: '用户' }
})

onMounted(() => {
  fetchSessions()

  // 从路由参数进入特定会话
  const listingId = route.params.listing_id
  const userId = route.params.user_id
  if (listingId && userId) {
    loadDirectChat(Number(listingId), Number(userId))
  }
})

async function fetchSessions() {
  loadingSessions.value = true
  try {
    const res = await request.get('/api/v1/marketplace/chat-sessions')
    sessions.value = res.data || []
  } catch (e) {
    console.error(e)
  } finally {
    loadingSessions.value = false
  }
}

async function loadDirectChat(listingId: number, userId: number) {
  // 创建或加载直接聊天
  try {
    const res = await request.get(`/api/v1/marketplace/messages/${listingId}/${userId}`, {
      params: { page: 1, page_size: 50 }
    })
    messages.value = res.data?.list || []

    // 创建临时会话对象
    currentSession.value = {
      listing_id: listingId,
      user1_id: currentUserId.value,
      user2_id: userId,
      listing: { product_name: '主机' }
    }

    await fetchSessions()
    scrollToBottom()
  } catch (e) {
    console.error(e)
  }
}

async function selectSession(session: any) {
  currentSession.value = session
  await fetchMessages()
}

async function fetchMessages() {
  if (!currentSession.value) return

  loadingMessages.value = true
  try {
    const otherId = currentSession.value.user1_id === currentUserId.value
      ? currentSession.value.user2_id
      : currentSession.value.user1_id

    const res = await request.get(
      `/api/v1/marketplace/messages/${currentSession.value.listing_id}/${otherId}`,
      { params: { page: 1, page_size: 100 } }
    )
    messages.value = res.data?.list || []
    scrollToBottom()
  } catch (e) {
    console.error(e)
  } finally {
    loadingMessages.value = false
  }
}

async function sendMessage() {
  if (!inputMessage.value.trim() || !currentSession.value) return

  sending.value = true
  try {
    const otherId = currentSession.value.user1_id === currentUserId.value
      ? currentSession.value.user2_id
      : currentSession.value.user1_id

    await request.post('/api/v1/marketplace/messages', {
      receiver_id: otherId,
      listing_id: currentSession.value.listing_id,
      content: inputMessage.value.trim()
    })

    inputMessage.value = ''
    await fetchMessages()
    await fetchSessions()
  } catch (e: any) {
    ElMessage.error(e.message || '发送失败')
  } finally {
    sending.value = false
  }
}

function viewListing() {
  if (currentSession.value?.listing_id) {
    router.push(`/user/marketplace?listing=${currentSession.value.listing_id}`)
  }
}

function isActiveSession(session: any) {
  return currentSession.value &&
    currentSession.value.listing_id === session.listing_id &&
    ((currentSession.value.user1_id === session.user1_id && currentSession.value.user2_id === session.user2_id) ||
     (currentSession.value.user1_id === session.user2_id && currentSession.value.user2_id === session.user1_id))
}

function getOtherUser(session: any) {
  const otherId = session.user1_id === currentUserId.value ? session.user2_id : session.user1_id
  return { id: otherId, username: `用户${otherId}` }
}

function scrollToBottom() {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

function formatTime(date: string): string {
  if (!date) return ''
  const d = new Date(date)
  const now = new Date()
  const diff = now.getTime() - d.getTime()

  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return `${Math.floor(diff / 60000)}分钟前`
  if (d.toDateString() === now.toDateString()) {
    return d.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  }
  return d.toLocaleDateString('zh-CN', { month: '2-digit', day: '2-digit' })
}
</script>

<style scoped lang="scss">
.marketplace-chat {
  display: flex;
  height: calc(100vh - 120px);
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
}

.chat-sidebar {
  width: 300px;
  border-right: 1px solid #eee;
  display: flex;
  flex-direction: column;
}

.sidebar-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 16px;
  border-bottom: 1px solid #eee;

  h3 {
    margin: 0;
    font-size: 16px;
  }
}

.session-list {
  flex: 1;
  overflow-y: auto;
}

.session-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  cursor: pointer;
  transition: background 0.2s;

  &:hover {
    background: #f5f7fa;
  }

  &.active {
    background: #ecf5ff;
  }
}

.session-avatar {
  flex-shrink: 0;
  background: #409eff;
  color: #fff;
}

.session-info {
  flex: 1;
  min-width: 0;
}

.session-top {
  display: flex;
  justify-content: space-between;
  margin-bottom: 4px;
}

.session-name {
  font-size: 14px;
  font-weight: 500;
}

.session-time {
  font-size: 11px;
  color: #999;
}

.session-bottom {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.session-product {
  font-size: 12px;
  color: #666;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.chat-main {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.no-chat {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid #eee;
}

.chat-user {
  display: flex;
  align-items: center;
  gap: 10px;

  .user-info {
    display: flex;
    flex-direction: column;
  }

  .name {
    font-weight: 500;
  }

  .product {
    font-size: 12px;
    color: #666;
  }
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.messages-wrapper {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.message-item {
  display: flex;
  gap: 8px;
  max-width: 70%;

  &.is-mine {
    flex-direction: row-reverse;
    margin-left: auto;

    .msg-bubble {
      background: #409eff;
      color: #fff;
    }

    .msg-time {
      text-align: right;
    }
  }
}

.msg-avatar {
  flex-shrink: 0;
  background: #e4e7ed;
}

.msg-content {
  display: flex;
  flex-direction: column;
}

.msg-bubble {
  background: #f5f7fa;
  padding: 10px 14px;
  border-radius: 8px;
  font-size: 14px;
  line-height: 1.5;
  word-break: break-word;
}

.msg-time {
  font-size: 11px;
  color: #999;
  margin-top: 4px;
}

.chat-input {
  padding: 12px 16px;
  border-top: 1px solid #eee;

  .input-actions {
    display: flex;
    justify-content: flex-end;
    margin-top: 8px;
  }
}

.empty-tip {
  padding: 40px 0;
}
</style>
