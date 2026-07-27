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
        <!-- Status Card -->
        <n-card class="status-card" :bordered="false">
          <n-result :status="statusConfig.status" :title="statusConfig.title" :description="statusConfig.description">
            <template #icon>
              <div class="status-icon-wrap" :class="payStatus">
                <n-icon :size="64" :color="statusConfig.iconColor">
                  <component :is="statusConfig.icon" />
                </n-icon>
              </div>
            </template>
            <template #footer>
              <div class="countdown-info" v-if="payStatus === 'success' && countdownSeconds > 0">
                <n-icon size="16" color="#86909c"><TimeOutline /></n-icon>
                <span>{{ countdownSeconds }} 秒后自动跳转到订单详情</span>
              </div>
            </template>
          </n-result>
        </n-card>

        <!-- Order Info Card -->
        <n-card class="order-card" :bordered="false" title="订单信息">
          <template #header-extra>
            <n-tag :type="statusConfig.tagType" :bordered="false">{{ statusConfig.tagText }}</n-tag>
          </template>

          <n-descriptions bordered :column="2" label-placement="left">
            <n-descriptions-item label="订单号">
              <span class="order-no">{{ orderInfo.orderNo }}</span>
            </n-descriptions-item>
            <n-descriptions-item label="商品名称">
              {{ orderInfo.productName }}
            </n-descriptions-item>
            <n-descriptions-item label="支付金额">
              <span class="amount">¥{{ orderInfo.amount.toFixed(2) }}</span>
            </n-descriptions-item>
            <n-descriptions-item label="支付方式">
              <div class="pay-method">
                <n-icon size="18" :color="payMethodConfig.color">
                  <component :is="payMethodConfig.icon" />
                </n-icon>
                <span>{{ orderInfo.payMethod }}</span>
              </div>
            </n-descriptions-item>
            <n-descriptions-item label="支付时间" :span="2">
              {{ orderInfo.payTime }}
            </n-descriptions-item>
            <n-descriptions-item label="订单备注" :span="2">
              {{ orderInfo.remark || '无' }}
            </n-descriptions-item>
          </n-descriptions>
        </n-card>

        <!-- Action Buttons -->
        <div class="action-buttons">
          <n-button size="large" @click="$router.push(`/orders/${orderInfo.orderNo}`)">
            <template #icon>
              <n-icon><DocumentTextOutline /></n-icon>
            </template>
            查看订单详情
          </n-button>
          <n-button type="primary" size="large" @click="$router.push('/')">
            <template #icon>
              <n-icon><HomeOutline /></n-icon>
            </template>
            返回首页
          </n-button>
          <n-button v-if="payStatus === 'fail'" type="warning" size="large" @click="retryPayment">
            <template #icon>
              <n-icon><RefreshOutline /></n-icon>
            </template>
            重新支付
          </n-button>
        </div>

        <!-- Help Tip -->
        <div class="help-tip">
          <n-icon size="16" color="#86909c"><InformationCircleOutline /></n-icon>
          <span>如遇到支付问题，请联系客服：<strong>400-888-8888</strong> 或提交工单获取帮助</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  AnchorOutline,
  DocumentTextOutline,
  HomeOutline,
  RefreshOutline,
  InformationCircleOutline,
  TimeOutline,
  CheckmarkCircleOutline,
  CloseCircleOutline,
  HourglassOutline,
  WalletOutline,
  CardOutline,
  PhonePortraitOutline
} from '@vicons/ionicons5'

const route = useRoute()
const router = useRouter()

// Get status from route query, default to 'success'
const payStatus = ref<'success' | 'fail' | 'processing'>(
  (route.query.status as 'success' | 'fail' | 'processing') || 'success'
)

const countdownSeconds = ref(10)
let countdownTimer: ReturnType<typeof setInterval> | null = null

const statusConfig = computed(() => {
  switch (payStatus.value) {
    case 'success':
      return {
        status: 'success' as const,
        title: '支付成功',
        description: '您的订单已支付成功，我们将尽快为您处理。',
        icon: CheckmarkCircleOutline,
        iconColor: '#18a058',
        tagType: 'success' as const,
        tagText: '已支付'
      }
    case 'fail':
      return {
        status: 'error' as const,
        title: '支付失败',
        description: '很抱歉，您的订单支付未成功，请重新尝试支付。',
        icon: CloseCircleOutline,
        iconColor: '#d03050',
        tagType: 'error' as const,
        tagText: '支付失败'
      }
    case 'processing':
      return {
        status: 'info' as const,
        title: '支付处理中',
        description: '您的订单正在处理中，请稍后查看支付结果。',
        icon: HourglassOutline,
        iconColor: '#1890ff',
        tagType: 'warning' as const,
        tagText: '处理中'
      }
    default:
      return {
        status: 'info' as const,
        title: '支付处理中',
        description: '您的订单正在处理中，请稍后查看支付结果。',
        icon: HourglassOutline,
        iconColor: '#1890ff',
        tagType: 'warning' as const,
        tagText: '处理中'
      }
  }
})

interface OrderInfo {
  orderNo: string
  productName: string
  amount: number
  payMethod: string
  payTime: string
  remark: string
}

const orderInfo = ref<OrderInfo>({
  orderNo: (route.query.orderNo as string) || 'AF20260727001234',
  productName: (route.query.product as string) || '香港云服务器 - 基础型',
  amount: Number(route.query.amount) || 49.00,
  payMethod: (route.query.payMethod as string) || '微信支付',
  payTime: (route.query.payTime as string) || '2026-07-27 10:05:30',
  remark: (route.query.remark as string) || ''
})

const payMethodConfig = computed(() => {
  const method = orderInfo.value.payMethod
  if (method.includes('微信')) {
    return { icon: PhonePortraitOutline, color: '#07c160' }
  } else if (method.includes('支付宝')) {
    return { icon: CardOutline, color: '#1677ff' }
  } else {
    return { icon: WalletOutline, color: '#1890ff' }
  }
})

function startCountdown() {
  if (payStatus.value !== 'success') return
  countdownTimer = setInterval(() => {
    countdownSeconds.value--
    if (countdownSeconds.value <= 0) {
      clearInterval(countdownTimer!)
      router.push(`/orders/${orderInfo.value.orderNo}`)
    }
  }, 1000)
}

function retryPayment() {
  router.push(`/checkout?orderNo=${orderInfo.value.orderNo}`)
}

onMounted(() => {
  startCountdown()
})

onBeforeUnmount(() => {
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
  display: flex;
  justify-content: center;
}

.result-container {
  width: 100%;
  max-width: 720px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* Status Card */
.status-card {
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
  padding: 16px 0;
}

.status-icon-wrap {
  width: 96px;
  height: 96px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
}

.status-icon-wrap.success {
  background: linear-gradient(135deg, #e8f8e8, #d4f0d4);
}

.status-icon-wrap.fail {
  background: linear-gradient(135deg, #fde8e8, #f8d4d4);
}

.status-icon-wrap.processing {
  background: linear-gradient(135deg, #e6f7ff, #bae7ff);
}

.countdown-info {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  color: #86909c;
  font-size: 14px;
  margin-top: 8px;
}

/* Order Card */
.order-card {
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
}

.order-no {
  font-family: 'SF Mono', 'Monaco', 'Menlo', monospace;
  font-weight: 600;
  color: #1890ff;
}

.amount {
  font-size: 20px;
  font-weight: 700;
  color: #ff7a45;
}

.pay-method {
  display: flex;
  align-items: center;
  gap: 6px;
}

/* Action Buttons */
.action-buttons {
  display: flex;
  justify-content: center;
  gap: 16px;
  flex-wrap: wrap;
}

/* Help Tip */
.help-tip {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 13px;
  color: #86909c;
  padding: 16px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
}

.help-tip strong {
  color: #1890ff;
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
