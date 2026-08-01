<template>
  <div class="ai-shopping-widget" v-if="enabled">
    <!-- 浮窗按钮 -->
    <div class="widget-trigger" @click="toggleChat" :class="{ active: isOpen }">
      <n-icon size="28">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
        </svg>
      </n-icon>
      <span class="widget-badge" v-if="unreadCount > 0">{{ unreadCount }}</span>
    </div>

    <!-- 聊天窗口 -->
    <transition name="slide-up">
      <div class="widget-chat" v-if="isOpen">
        <!-- 头部 -->
        <div class="chat-header">
          <div class="header-info">
            <n-icon size="20" color="#1890ff">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
              </svg>
            </n-icon>
            <span>{{ title }}</span>
          </div>
          <n-button text @click="toggleChat">
            <template #icon>
              <n-icon><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18" /><line x1="6" y1="6" x2="18" y2="18" /></svg></n-icon>
            </template>
          </n-button>
        </div>

        <!-- 消息列表 -->
        <div class="chat-messages" ref="messagesRef">
          <div v-for="(msg, index) in messages" :key="index" :class="['message', msg.role]">
            <div class="message-avatar" v-if="msg.role === 'assistant'">
              <n-icon size="20" color="#1890ff">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M12 2a10 10 0 0 1 10 10c0 5.52-4.48 10-10 10S2 17.52 2 12 6.48 2 12 2z" />
                  <path d="M8 14s1.5 2 4 2 4-2 4-2" />
                  <line x1="9" y1="9" x2="9.01" y2="9" />
                  <line x1="15" y1="9" x2="15.01" y2="9" />
                </svg>
              </n-icon>
            </div>
            <div class="message-content">
              <div v-html="formatMessage(msg.content)"></div>
            </div>
            <div class="message-avatar" v-if="msg.role === 'user'">
              <n-icon size="20" color="#52c41a">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
                  <circle cx="12" cy="7" r="4" />
                </svg>
              </n-icon>
            </div>
          </div>
          <div v-if="loading" class="message assistant">
            <div class="message-avatar">
              <n-icon size="20" color="#1890ff">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M12 2a10 10 0 0 1 10 10c0 5.52-4.48 10-10 10S2 17.52 2 12 6.48 2 12 2z" />
                </svg>
              </n-icon>
            </div>
            <div class="message-content typing">
              <span class="dot"></span>
              <span class="dot"></span>
              <span class="dot"></span>
            </div>
          </div>
        </div>

        <!-- 输入框 -->
        <div class="chat-input">
          <n-input
            v-model:value="inputMessage"
            placeholder="输入消息..."
            :disabled="loading"
            @keyup.enter="sendMessage"
            size="small"
          />
          <n-button type="primary" size="small" :disabled="!inputMessage.trim() || loading" @click="sendMessage">
            发送
          </n-button>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick, watch } from 'vue'
import { useRoute } from 'vue-router'
import request from '@/utils/http'

const route = useRoute()

const enabled = ref(false)
const title = ref('AI导购')
const isOpen = ref(false)
const loading = ref(false)
const inputMessage = ref('')
const messages = ref<Array<{ role: string; content: string }>>([])
const messagesRef = ref<HTMLElement>()
const unreadCount = ref(0)
const sessionId = ref('')

// 页面上下文
const pageContext = ref('')

// 获取配置
const fetchConfig = async () => {
  try {
    const res = await request.get('/api/v1/ai-shopping/config')
    if (res.data) {
      enabled.value = res.data.ai_enabled === '1'
      title.value = res.data.widget_title || 'AI导购'
    }
  } catch {
    // 静默失败
  }
}

// 生成会话ID
const generateSessionId = () => {
  return 'shopping_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9)
}

// 检测页面上下文
const detectPageContext = () => {
  const path = route.path
  const query = route.query

  if (path.startsWith('/products')) {
    if (query.group) {
      pageContext.value = `用户正在浏览商品分组：${query.group}`
    } else if (query.id) {
      pageContext.value = `用户正在查看商品详情，商品ID：${query.id}`
    } else {
      pageContext.value = '用户正在浏览商品列表页'
    }
  } else if (path.startsWith('/order') || path.startsWith('/cart')) {
    pageContext.value = '用户正在订单/购物车页面'
  } else {
    pageContext.value = ''
  }
}

// 切换聊天窗口
const toggleChat = () => {
  isOpen.value = !isOpen.value
  if (isOpen.value) {
    unreadCount.value = 0
    if (messages.value.length === 0) {
      // 发送欢迎消息
      messages.value.push({
        role: 'assistant',
        content: '你好，我是 AI 导购，可以帮你对比商品、选配置、估算价格。直接告诉我你的需求即可。'
      })
    }
    nextTick(() => scrollToBottom())
  }
}

// 发送消息
const sendMessage = async () => {
  const msg = inputMessage.value.trim()
  if (!msg || loading.value) return

  messages.value.push({ role: 'user', content: msg })
  inputMessage.value = ''
  loading.value = true

  nextTick(() => scrollToBottom())

  try {
    const res = await request.post(`/api/v1/ai-shopping/chat/${sessionId.value}`, {
      message: msg,
      page_context: pageContext.value
    })

    if (res.data?.reply) {
      // 处理多消息分段
      const replyParts = res.data.reply.split('[[[MSG]]]').filter((p: string) => p.trim())
      for (const part of replyParts) {
        messages.value.push({ role: 'assistant', content: part.trim() })
      }
    }

    if (!isOpen.value) {
      unreadCount.value++
    }
  } catch {
    messages.value.push({
      role: 'assistant',
      content: '抱歉，AI 暂时无法响应，请稍后重试。'
    })
  } finally {
    loading.value = false
    nextTick(() => scrollToBottom())
  }
}

// 格式化消息（支持Markdown链接）
const formatMessage = (content: string) => {
  // 简单的Markdown链接转换
  return content
    .replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" target="_blank" rel="noopener">$1</a>')
    .replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>')
    .replace(/\n/g, '<br>')
}

// 滚动到底部
const scrollToBottom = () => {
  if (messagesRef.value) {
    messagesRef.value.scrollTop = messagesRef.value.scrollHeight
  }
}

// 监听路由变化，更新页面上下文
watch(() => route.path, () => {
  detectPageContext()
})

watch(() => route.query, () => {
  detectPageContext()
})

onMounted(() => {
  sessionId.value = generateSessionId()
  fetchConfig()
  detectPageContext()
})
</script>

<style scoped>
.ai-shopping-widget {
  position: fixed;
  bottom: 24px;
  right: 24px;
  z-index: 9999;
}

.widget-trigger {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: linear-gradient(135deg, #1890ff 0%, #096dd9 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(24, 144, 255, 0.4);
  transition: all 0.3s ease;
  position: relative;
}

.widget-trigger:hover {
  transform: scale(1.1);
  box-shadow: 0 6px 16px rgba(24, 144, 255, 0.5);
}

.widget-trigger.active {
  background: linear-gradient(135deg, #ff4d4f 0%, #cf1322 100%);
  box-shadow: 0 4px 12px rgba(255, 77, 79, 0.4);
}

.widget-badge {
  position: absolute;
  top: -4px;
  right: -4px;
  background: #ff4d4f;
  color: white;
  font-size: 12px;
  min-width: 20px;
  height: 20px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 6px;
}

.widget-chat {
  position: absolute;
  bottom: 70px;
  right: 0;
  width: 380px;
  height: 500px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.chat-header {
  padding: 16px;
  background: linear-gradient(135deg, #1890ff 0%, #096dd9 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-info {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.message {
  display: flex;
  gap: 8px;
  max-width: 85%;
}

.message.user {
  align-self: flex-end;
  flex-direction: row-reverse;
}

.message.assistant {
  align-self: flex-start;
}

.message-avatar {
  flex-shrink: 0;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #f0f2f5;
  display: flex;
  align-items: center;
  justify-content: center;
}

.message.user .message-avatar {
  background: #f6ffed;
}

.message-content {
  padding: 10px 14px;
  border-radius: 12px;
  background: #f0f2f5;
  line-height: 1.5;
  font-size: 14px;
}

.message.user .message-content {
  background: linear-gradient(135deg, #1890ff 0%, #096dd9 100%);
  color: white;
}

.message-content :deep(a) {
  color: #1890ff;
  text-decoration: underline;
}

.message.user .message-content :deep(a) {
  color: #bae7ff;
}

.typing {
  display: flex;
  gap: 4px;
  padding: 12px 16px;
}

.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #1890ff;
  animation: bounce 1.4s infinite ease-in-out;
}

.dot:nth-child(1) { animation-delay: -0.32s; }
.dot:nth-child(2) { animation-delay: -0.16s; }

@keyframes bounce {
  0%, 80%, 100% { transform: scale(0); }
  40% { transform: scale(1); }
}

.chat-input {
  padding: 12px;
  border-top: 1px solid #f0f0f0;
  display: flex;
  gap: 8px;
}

.chat-input .n-input {
  flex: 1;
}

.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.3s ease;
}

.slide-up-enter-from,
.slide-up-leave-to {
  opacity: 0;
  transform: translateY(20px);
}
</style>
