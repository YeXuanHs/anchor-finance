<template>
  <div class="payment-result-page">
    <!-- Header -->
    <header class="header">
      <div class="header-inner">
        <router-link to="/" class="logo">
          <div class="logo-icon">
            <n-icon size="24" color="#fff"><AnchorOutline /></n-icon>
          </div>
          <span class="logo-text">锚点财务</span>
        </router-link>
        <nav class="nav-links">
          <router-link to="/" class="nav-link">首页</router-link>
          <router-link to="/products" class="nav-link">产品</router-link>
          <a href="#" class="nav-link">公告</a>
          <a href="#" class="nav-link">帮助</a>
        </nav>
        <div class="header-actions">
          <n-button text @click="$router.push('/login')">登录</n-button>
          <n-button type="primary" round size="small" @click="$router.push('/register')">免费注册</n-button>
        </div>
      </div>
    </header>

    <!-- Breadcrumb -->
    <div class="breadcrumb-bar">
      <div class="breadcrumb-inner">
        <n-breadcrumb>
          <n-breadcrumb-item @click="$router.push('/')">首页</n-breadcrumb-item>
          <n-breadcrumb-item>支付结果</n-breadcrumb-item>
        </n-breadcrumb>
      </div>
    </div>

    <!-- Main Content -->
    <div class="main-content">
      <div class="result-container">
        <!-- Status Result -->
        <n-card class="status-card" :bordered="false">
          <n-result
            :status="resultStatus"
            :title="resultTitle"
            :description="resultDescription"
          >
            <template #icon>
              <div class="status-icon-wrap" :class="statusClass">
                <n-icon :size="64" :color="statusColor">
                  <component :is="statusIcon" />
                </n-icon>
              </div>
            </template>
            <template #footer>
              <div class="countdown-area" v-if="countdownSeconds > 0">
                <span class="countdown-text">{{ countdownSeconds }} 秒后自动跳转至</span>
                <n-button text type="primary" @click="goToTarget">{{ countdownTargetLabel }}</n-button>
              </div>
            </template>
          </n-result>
        </n-card>

        <!-- Order Info -->
        <n-card class="order-card" :bordered="false">
          <template #header>
            <div class="card-header">
              <n-icon size="20" color="#1890ff"><ReceiptOutline /></n-icon>
              <span>订单信息</span>
            </div>
          </template>
          <n-descriptions
            :column="2"
            bordered
            label-placement="left"
            :label-style="{ width: '120px', fontWeight: 600 }"
          >
            <n-descriptions-item label="订单号">
              <span class="order-no">{{ orderInfo.orderNo }}</span>
              <n-button text type="primary" size="tiny" @click="copyOrderNo" style="margin-left: 8px;">
                复制
              </n-button>
            </n-descriptions-item>
            <n-descriptions-item label="商品名称">
              {{ orderInfo.productName }}
            </n-descriptions-item>
            <n-descriptions-item label="支付金额">
              <span class="amount">¥{{ orderInfo.amount.toFixed(2) }}</span>
            </n-descriptions-item>
            <n-descriptions-item label="支付方式">
              <n-tag :type="paymentTagType" size="small" :bordered="false">
                {{ orderInfo.paymentMethod }}
              </n-tag>
            </n-descriptions-item>
            <n-descriptions-item label="支付时间">
              {{ orderInfo.payTime }}
            </n-descriptions-item>
            <n-descriptions-item label="订单状态">
              <n-tag :type="statusTagType" size="small" :bordered="false">
                {{ statusLabel }}
              </n-tag>
            </n-descriptions-item>
          </n-descriptions>
        </n-card>

        <!-- Action Buttons -->
        <div class="action-buttons">
          <n-button size="large" @click="$router.push('/orders')">
            <template #icon>
              <n-icon><ListOutline /></n-icon>
            </template>
            查看订单
          </n-button>
          <n-button type="primary" size="large" @click="$router.push('/')">
            <template #icon>
              <n-icon><HomeOutline /></n-icon>
            </template>
            返回首页
          </n-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  AnchorOutline,
  ReceiptOutline,
  ListOutline,
  HomeOutline,
  CheckmarkCircleOutline,
  CloseCircleOutline,
  TimeOutline,
  AlertCircleOutline
} from '@vicons/ionicons5'

const router = useRouter()
const route = useRoute()

// Payment status: success | fail | processing
const paymentStatus = ref<'success' | 'fail' | 'processing'>(
  (route.query.status as 'success' | 'fail' | 'processing') || 'success'
)

const countdownSeconds = ref(10)
let countdownTimer: ReturnType<typeof setInterval> | null = null

interface OrderInfo {
  orderNo: string
  productName: string
  amount: number
  paymentMethod: string
  payTime: string
}

const orderInfo = ref<OrderInfo>({
  orderNo: (route.query.orderNo as string) || 'AF20251215143025001',
  productName: (route.query.product as string) || '香港云服务器 - 基础型',
  amount: Number(route.query.amount) || 49.00,
  paymentMethod: (route.query.method as string) || '支付宝',
  payTime: (route.query.time as string) || '2025-12-15 14:30:25'
})

const resultStatus = computed(() => {
  const map = {
    success: 'success',
    fail: 'error',
    processing: 'info'
  }
  return map[paymentStatus.value] as 'success' | 'error' | 'info'
})

const resultTitle = computed(() => {
  const map = {
    success: '支付成功',
    fail: '支付失败',
    processing: '支付处理中'
  }
  return map[paymentStatus.value]
})

const resultDescription = computed(() => {
  const map = {
    success: '您的订单已支付成功，服务将在1-5分钟内自动开通，请注意查收邮件通知。',
    fail: '很抱歉，您的支付未成功完成，请重新尝试支付或联系客服获取帮助。',
    processing: '您的支付正在处理中，请稍候，我们将在确认后为您开通服务。'
  }
  return map[paymentStatus.value]
})

const statusClass = computed(() => paymentStatus.value)

const statusColor = computed(() => {
  const map = {
    success: '#18a058',
    fail: '#d03050',
    processing: '#1890ff'
  }
  return map[paymentStatus.value]
})

const statusIcon = computed(() => {
  const map = {
    success: CheckmarkCircleOutline,
    fail: CloseCircleOutline,
    processing: TimeOutline
  }
  return map[paymentStatus.value]
})

const statusLabel = computed(() => {
  const map = {
    success: '已支付',
    fail: '支付失败',
    processing: '处理中'
  }
  return map[paymentStatus.value]
})

const statusTagType = computed(() => {
  const map = {
    success: 'success',
    fail: 'error',
    processing: 'warning'
  }
  return map[paymentStatus.value] as 'success' | 'error' | 'warning'
})

const paymentTagType = computed(() => {
  const methods: Record<string, 'info' | 'success' | 'warning'> = {
    '支付宝': 'info',
    '微信支付': 'success',
    '银行卡': 'warning'
  }
  return methods[orderInfo.value.paymentMethod] || 'info'
})

const countdownTargetLabel = computed(() => {
  return paymentStatus.value === 'success' ? '查看订单' : '返回首页'
})

function goToTarget() {
  if (paymentStatus.value === 'success') {
    router.push('/orders')
  } else {
    router.push('/')
  }
}

function copyOrderNo() {
  navigator.clipboard.writeText(orderInfo.value.orderNo).catch(() => {
    // fallback
    const textarea = document.createElement('textarea')
    textarea.value = orderInfo.value.orderNo
    document.body.appendChild(textarea)
    textarea.select()
    document.execCommand('copy')
    document.body.removeChild(textarea)
  })
}

onMounted(() => {
  countdownTimer = setInterval(() => {
    if (countdownSeconds.value > 0) {
      countdownSeconds.value--
    } else {
      goToTarget()
    }
  }, 1000)
})

onUnmounted(() => {
  if (countdownTimer) {
    clearInterval(countdownTimer)
  }
})
</script>

<style scoped>
.payment-result-page {
  min-height: 100vh;
  background: #f7f8fa;
}

/* Header */
.header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 64px;
  background: #fff;
  box-shadow: 0 1px 8px rgba(0, 0, 0, 0.06);
  z-index: 100;
}

.header-inner {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 24px;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  text-decoration: none;
}

.logo-icon {
  width: 34px;
  height: 34px;
  background: linear-gradient(135deg, #1890ff, #096dd9);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-text {
  font-size: 18px;
  font-weight: 700;
  color: #1d2129;
}

.nav-links {
  display: flex;
  gap: 32px;
}

.nav-link {
  color: #4e5969;
  text-decoration: none;
  font-size: 15px;
  font-weight: 500;
  transition: color 0.2s;
}

.nav-link:hover,
.nav-link.active {
  color: #1890ff;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

/* Breadcrumb */
.breadcrumb-bar {
  background: #fff;
  border-bottom: 1px solid #f0f1f5;
  margin-top: 64px;
}

.breadcrumb-inner {
  max-width: 1400px;
  margin: 0 auto;
  padding: 14px 24px;
}

/* Main Content */
.main-content {
  padding: 40px 24px;
}

.result-container {
  max-width: 720px;
  margin: 0 auto;
}

/* Status Card */
.status-card {
  margin-bottom: 24px;
  border-radius: 12px;
  text-align: center;
  padding: 20px 0;
}

.status-icon-wrap {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
}

.status-icon-wrap.success {
  background: linear-gradient(135deg, #e8f8ef, #b7ebd1);
}

.status-icon-wrap.fail {
  background: linear-gradient(135deg, #fde8eb, #fcc6cb);
}

.status-icon-wrap.processing {
  background: linear-gradient(135deg, #e6f7ff, #bae7ff);
}

.countdown-area {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  margin-top: 16px;
}

.countdown-text {
  font-size: 14px;
  color: #86909c;
}

/* Order Card */
.order-card {
  margin-bottom: 24px;
  border-radius: 12px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: #1d2129;
}

.order-no {
  font-family: 'Courier New', monospace;
  font-weight: 600;
  color: #1d2129;
}

.amount {
  font-size: 18px;
  font-weight: 700;
  color: #ff7a45;
}

/* Action Buttons */
.action-buttons {
  display: flex;
  justify-content: center;
  gap: 16px;
}

/* Responsive */
@media (max-width: 768px) {
  .nav-links {
    display: none;
  }

  .action-buttons {
    flex-direction: column;
  }

  .action-buttons .n-button {
    width: 100%;
  }
}
</style>
