<template>
  <div class="payment-result-page">
    <!-- Header -->
    <header class="page-header">
      <div class="header-inner">
        <router-link to="/" class="logo">
          <div class="logo-icon">
            <el-icon :size="22" color="#fff"><Monitor /></el-icon>
          </div>
          <span class="logo-text">{{ siteName }}</span>
        </router-link>
        <nav class="nav-links">
          <router-link to="/" class="nav-link">首页</router-link>
          <router-link to="/products" class="nav-link">产品</router-link>
          <router-link to="/" class="nav-link">公告</router-link>
          <router-link to="/" class="nav-link">帮助</router-link>
        </nav>
        <div class="header-actions">
          <el-button text @click="$router.push('/login')">登录</el-button>
          <el-button class="btn-gradient" round size="small" @click="$router.push('/register')">免费注册</el-button>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <div class="main-content">
      <div class="result-container">
        <!-- Status Result -->
        <div class="status-card">
          <el-result
            :icon="resultIcon"
            :title="resultTitle"
            :sub-title="resultDescription"
          >
            <template #icon>
              <div class="status-icon-wrap" :class="paymentStatus">
                <el-icon :size="56" :color="statusColor">
                  <component :is="statusIcon" />
                </el-icon>
              </div>
            </template>
            <template #extra>
              <div class="countdown-area" v-if="countdownSeconds > 0">
                <span class="countdown-text">{{ countdownSeconds }} 秒后自动跳转至</span>
                <el-button text type="primary" @click="goToTarget" class="countdown-link">
                  {{ countdownTargetLabel }}
                </el-button>
              </div>
            </template>
          </el-result>
        </div>

        <!-- Order Info -->
        <div class="order-card">
          <h3 class="card-title">
            <el-icon :size="18" color="#0056FF"><Document /></el-icon>
            订单信息
          </h3>
          <el-descriptions
            :column="2"
            border
            label-placement="left"
            :label-style="{ width: '120px', fontWeight: 600 }"
          >
            <el-descriptions-item label="订单号">
              <span class="order-no">{{ orderInfo.orderNo }}</span>
              <el-button text type="primary" size="small" @click="copyOrderNo" style="margin-left: 8px">
                复制
              </el-button>
            </el-descriptions-item>
            <el-descriptions-item label="商品名称">
              {{ orderInfo.productName }}
            </el-descriptions-item>
            <el-descriptions-item label="支付金额">
              <span class="amount">¥{{ orderInfo.amount.toFixed(2) }}</span>
            </el-descriptions-item>
            <el-descriptions-item label="支付方式">
              <el-tag :type="paymentTagType" size="small" effect="plain" round>
                {{ orderInfo.paymentMethod }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="支付时间">
              {{ orderInfo.payTime }}
            </el-descriptions-item>
            <el-descriptions-item label="订单状态">
              <el-tag :type="statusTagType" size="small" effect="plain" round>
                {{ statusLabel }}
              </el-tag>
            </el-descriptions-item>
          </el-descriptions>
        </div>

        <!-- Action Buttons -->
        <div class="action-buttons">
          <el-button size="large" round @click="$router.push('/user/orders')" class="action-btn">
            <el-icon style="margin-right: 6px"><List /></el-icon>
            查看订单
          </el-button>
          <el-button class="btn-gradient action-btn" size="large" round @click="$router.push('/')">
            <el-icon style="margin-right: 6px"><HomeFilled /></el-icon>
            返回首页
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Monitor,
  CircleCheckFilled,
  CircleCloseFilled,
  Loading,
  Document,
  List,
  HomeFilled
} from '@element-plus/icons-vue'
import request from '@/utils/request'

const router = useRouter()
const route = useRoute()

// Payment status: success | fail | processing
const paymentStatus = ref<'success' | 'fail' | 'processing'>(
  (route.query.status as 'success' | 'fail' | 'processing') || 'success'
)

const countdownSeconds = ref(10)
let countdownTimer: ReturnType<typeof setInterval> | null = null

const siteName = ref('')

const fetchSiteSettings = async () => {
  try {
    const res = await request.get('/api/v1/settings/public')
    if (res?.data?.site_name) {
      siteName.value = res.data.site_name
    }
  } catch {
    // Use empty string as fallback
  }
}

interface OrderInfo {
  orderNo: string
  productName: string
  amount: number
  paymentMethod: string
  payTime: string
}

const orderInfo = ref<OrderInfo>({
  orderNo: (route.query.orderNo as string) || '',
  productName: (route.query.product as string) || '',
  amount: Number(route.query.amount) || 0,
  paymentMethod: (route.query.method as string) || '',
  payTime: (route.query.time as string) || ''
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

const resultIcon = computed(() => {
  const map = {
    success: 'success',
    fail: 'error',
    processing: 'warning'
  }
  return map[paymentStatus.value]
})

const statusColor = computed(() => {
  const map = {
    success: '#00b42a',
    fail: '#f53f3f',
    processing: '#0056FF'
  }
  return map[paymentStatus.value]
})

const statusIcon = computed(() => {
  const map = {
    success: CircleCheckFilled,
    fail: CircleCloseFilled,
    processing: Loading
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
    fail: 'danger',
    processing: 'warning'
  }
  return map[paymentStatus.value] as 'success' | 'danger' | 'warning'
})

const paymentTagType = computed(() => {
  // Style mapping uses English keys; display names come from API/i18n
  const methods: Record<string, 'info' | 'success' | 'warning' | 'primary'> = {
    alipay: 'primary',
    wechat: 'success',
    qqpay: 'primary',
    balance: 'warning'
  }
  // Reverse lookup: Chinese name -> English key
  const methodNameToKey: Record<string, string> = {
    '支付宝': 'alipay',
    '微信支付': 'wechat',
    'QQ钱包': 'qqpay',
    '余额支付': 'balance'
  }
  const raw = orderInfo.value.paymentMethod
  const key = methodNameToKey[raw] || raw
  return methods[key] || 'info'
})

const countdownTargetLabel = computed(() => {
  return paymentStatus.value === 'success' ? '查看订单' : '返回首页'
})

function goToTarget() {
  if (paymentStatus.value === 'success') {
    router.push('/user/orders')
  } else {
    router.push('/')
  }
}

function copyOrderNo() {
  navigator.clipboard.writeText(orderInfo.value.orderNo).then(() => {
    ElMessage.success('已复制')
  }).catch(() => {
    const textarea = document.createElement('textarea')
    textarea.value = orderInfo.value.orderNo
    document.body.appendChild(textarea)
    textarea.select()
    document.execCommand('copy')
    document.body.removeChild(textarea)
    ElMessage.success('已复制')
  })
}

onMounted(() => {
  fetchSiteSettings()
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
  background: #f5f7fa;
}

/* Header */
.page-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 60px;
  background: #fff;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
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
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #0056FF 0%, #4080FF 100%);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-text {
  font-size: 17px;
  font-weight: 700;
  color: #1a3a5c;
}

.nav-links {
  display: flex;
  gap: 28px;
}

.nav-link {
  color: #666;
  text-decoration: none;
  font-size: 14px;
  font-weight: 500;
  transition: color 0.2s;
  padding: 4px 0;
}

.nav-link:hover,
.nav-link.active {
  color: #0056FF;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

/* Gradient Button */
.btn-gradient {
  background: linear-gradient(135deg, #0056FF 0%, #4080FF 100%) !important;
  border: none !important;
  color: #fff !important;
  font-weight: 500;
}

.btn-gradient:hover {
  opacity: 0.9;
}

/* Main Content */
.main-content {
  padding-top: 100px;
  padding-bottom: 40px;
}

.result-container {
  max-width: 720px;
  margin: 0 auto;
  padding: 0 24px;
}

/* Status Card */
.status-card {
  background: #fff;
  border-radius: 16px;
  padding: 24px 32px 32px;
  text-align: center;
  margin-bottom: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.status-icon-wrap {
  width: 88px;
  height: 88px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
}

.status-icon-wrap.success {
  background: linear-gradient(135deg, #e8ffea 0%, #b3f0b8 100%);
}

.status-icon-wrap.fail {
  background: linear-gradient(135deg, #ffe8e8 0%, #fcb8b8 100%);
}

.status-icon-wrap.processing {
  background: linear-gradient(135deg, #EBF3FD 0%, #d6e6ff 100%);
}

:deep(.el-result__title) {
  font-size: 22px;
  font-weight: 700;
  color: #1a3a5c;
}

:deep(.el-result__subtitle) {
  font-size: 14px;
  color: #999;
  line-height: 1.6;
}

.countdown-area {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  margin-top: 12px;
}

.countdown-text {
  font-size: 13px;
  color: #999;
}

.countdown-link {
  font-size: 13px;
}

/* Order Card */
.order-card {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.card-title {
  font-size: 15px;
  font-weight: 600;
  color: #1a3a5c;
  margin-bottom: 18px;
  display: flex;
  align-items: center;
  gap: 8px;
}

:deep(.el-descriptions__label) {
  font-weight: 600;
  color: #666;
}

:deep(.el-descriptions__content) {
  color: #1a3a5c;
}

.order-no {
  font-family: 'Courier New', monospace;
  font-weight: 600;
  color: #1a3a5c;
}

.amount {
  font-size: 16px;
  font-weight: 700;
  color: #ff6b35;
}

/* Action Buttons */
.action-buttons {
  display: flex;
  justify-content: center;
  gap: 16px;
}

.action-btn {
  min-width: 140px;
  height: 44px;
  font-size: 14px;
}

/* Responsive */
@media (max-width: 768px) {
  .nav-links {
    display: none;
  }

  .action-buttons {
    flex-direction: column;
  }

  .action-buttons .el-button {
    width: 100%;
  }
}
</style>
